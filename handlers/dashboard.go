package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"vm-manager/models"
	"vm-manager/services"

	"github.com/labstack/echo/v4"
)

// AppContext contiene el contexto de la aplicación
type AppContext struct {
	VBoxService *services.VirtualBoxService
	SSHService  *services.SSHService
	Manager     *models.VMManager
	mu          sync.RWMutex
}

// NewAppContext crea un nuevo contexto de aplicación
func NewAppContext(vboxService *services.VirtualBoxService, sshService *services.SSHService) *AppContext {
	return &AppContext{
		VBoxService: vboxService,
		SSHService:  sshService,
		Manager: &models.VMManager{
			BaseMachines: make(map[string]*models.BaseMachine),
			UserMachines: make(map[string]*models.UserMachine),
			MediaDisks:   make(map[string]*models.MediaDisk),
			SSHKeys:      make(map[string]*models.SSHKey),
			RootKeys:     make(map[string]*models.SSHKey),
		},
	}
}

// GetDashboard retorna el dashboard
func (ac *AppContext) GetDashboard(c echo.Context) error {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	baseMachines := make([]models.BaseMachine, 0)
	for _, m := range ac.Manager.BaseMachines {
		baseMachines = append(baseMachines, *m)
	}

	userMachines := make([]models.UserMachine, 0)
	for _, m := range ac.Manager.UserMachines {
		userMachines = append(userMachines, *m)
	}

	mediaDisks := make([]models.MediaDisk, 0)
	for _, d := range ac.Manager.MediaDisks {
		mediaDisks = append(mediaDisks, *d)
	}

	dashboard := models.Dashboard{
		BaseMachines: baseMachines,
		MediaDisks:   mediaDisks,
		UserMachines: userMachines,
	}

	return c.Render(http.StatusOK, "dashboard.html", dashboard)
}

// AddBaseMachine agrega una máquina virtual base
func (ac *AppContext) AddBaseMachine(c echo.Context) error {
	name := c.FormValue("name")
	description := c.FormValue("description")

	if name == "" {
		return c.String(http.StatusBadRequest, "El nombre de la máquina es requerido")
	}

	machine, err := ac.VBoxService.CreateBaseMachine(name, description)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creando máquina: "+err.Error())
	}

	ac.mu.Lock()
	ac.Manager.BaseMachines[machine.ID] = machine
	ac.mu.Unlock()

	return c.Redirect(http.StatusSeeOther, "/")
}

// CreateRootKeys crea las llaves SSH para el root de una máquina base
func (ac *AppContext) CreateRootKeys(c echo.Context) error {
	machineID := c.Param("id")

	ac.mu.RLock()
	machine, exists := ac.Manager.BaseMachines[machineID]
	ac.mu.RUnlock()

	if !exists {
		return c.String(http.StatusNotFound, "Máquina no encontrada")
	}

	// Crear llaves
	sshKey, err := ac.SSHService.CreateRootSSHKeys(machineID, 1024)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creando llaves: "+err.Error())
	}

	ac.mu.Lock()
	machine.RootKeysReady = true
	ac.Manager.RootKeys[machineID] = sshKey
	ac.mu.Unlock()

	return c.String(http.StatusOK, "Llaves de root creadas exitosamente")
}

// DownloadRootKeys descarga las llaves SSH para el root
func (ac *AppContext) DownloadRootKeys(c echo.Context) error {
	machineID := c.Param("id")

	ac.mu.RLock()
	sshKey, exists := ac.Manager.RootKeys[machineID]
	ac.mu.RUnlock()

	if !exists {
		return c.String(http.StatusNotFound, "Llaves no encontradas")
	}

	c.Response().Header().Set("Content-Disposition", `attachment; filename="id_rsa"`)
	c.Response().Header().Set("Content-Type", "text/plain")
	return c.String(http.StatusOK, sshKey.PrivateKey)
}

// CreateMediaDisk crea un disco multiconexión
func (ac *AppContext) CreateMediaDisk(c echo.Context) error {
	machineID := c.Param("id")
	diskName := c.FormValue("disk_name")
	diskSizeStr := c.FormValue("disk_size")

	if diskName == "" {
		return c.String(http.StatusBadRequest, "Nombre de disco requerido")
	}

	ac.mu.RLock()
	machine, exists := ac.Manager.BaseMachines[machineID]
	ac.mu.RUnlock()

	if !exists {
		return c.String(http.StatusNotFound, "Máquina no encontrada")
	}

	if !machine.RootKeysReady {
		return c.String(http.StatusBadRequest, "Primero debe crear las llaves de root")
	}

	// Convertir diskSize a int64 con valor por defecto de 10GB
	var size int64 = 10240
	if diskSizeStr != "" {
		parsedSize, err := strconv.ParseInt(diskSizeStr, 10, 64)
		if err == nil {
			size = parsedSize
		}
	}

	disk, err := ac.VBoxService.CreateMediaDisk(diskName, size, machineID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creando disco: "+err.Error())
	}

	ac.mu.Lock()
	ac.Manager.MediaDisks[disk.ID] = disk
	machine.Disks[disk.ID] = disk
	ac.mu.Unlock()

	return c.Redirect(http.StatusSeeOther, "/")
}

// AddUserMachine agrega una máquina virtual de usuario
func (ac *AppContext) AddUserMachine(c echo.Context) error {
	diskID := c.Param("diskid")
	name := c.FormValue("name")
	description := c.FormValue("description")
	owner := c.FormValue("owner")

	if name == "" || owner == "" {
		return c.String(http.StatusBadRequest, "Datos incompletos")
	}

	ac.mu.RLock()
	disk, diskExists := ac.Manager.MediaDisks[diskID]
	ac.mu.RUnlock()

	if !diskExists {
		return c.String(http.StatusNotFound, "Disco no encontrado")
	}

	baseMachineID := disk.BaseMachine

	machine, err := ac.VBoxService.CreateUserMachine(name, description, baseMachineID, diskID, owner)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creando máquina: "+err.Error())
	}

	ac.mu.Lock()
	ac.Manager.UserMachines[machine.ID] = machine
	if baseMachine, exists := ac.Manager.BaseMachines[baseMachineID]; exists {
		baseMachine.Users[machine.ID] = machine
	}
	ac.mu.Unlock()

	return c.Redirect(http.StatusSeeOther, "/")
}

// CreateUserSSHKeys crea llaves SSH para un usuario
func (ac *AppContext) CreateUserSSHKeys(c echo.Context) error {
	machineID := c.Param("vmid")
	username := c.FormValue("username")

	if username == "" {
		return c.String(http.StatusBadRequest, "Nombre de usuario requerido")
	}

	// Crear llaves
	sshKey, err := ac.SSHService.CreateUserSSHKeys(machineID, username, 1024)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creando llaves: "+err.Error())
	}

	ac.mu.Lock()
	ac.Manager.SSHKeys[machineID+"_"+username] = sshKey
	if machine, exists := ac.Manager.UserMachines[machineID]; exists {
		machine.KeysReady = true
	}
	ac.mu.Unlock()

	return c.String(http.StatusOK, "Llaves de usuario creadas")
}

// DownloadUserKeys descarga las llaves SSH para un usuario
func (ac *AppContext) DownloadUserKeys(c echo.Context) error {
	machineID := c.Param("vmid")
	username := c.Param("username")

	key, err := ac.SSHService.DownloadPrivateKey(machineID, username)
	if err != nil {
		return c.String(http.StatusNotFound, "Llaves no encontradas")
	}

	c.Response().Header().Set("Content-Disposition", `attachment; filename="id_rsa"`)
	c.Response().Header().Set("Content-Type", "text/plain")
	return c.Blob(http.StatusOK, "text/plain", key)
}

// DeleteUserMachine elimina una máquina virtual de usuario
func (ac *AppContext) DeleteUserMachine(c echo.Context) error {
	machineID := c.Param("vmid")

	ac.mu.RLock()
	machine, exists := ac.Manager.UserMachines[machineID]
	ac.mu.RUnlock()

	if !exists {
		return c.String(http.StatusNotFound, "Máquina no encontrada")
	}

	// Detener si está corriendo
	if machine.IsRunning {
		ac.VBoxService.StopMachine(machine.Name)
	}

	// Eliminar máquina
	if err := ac.VBoxService.DeleteMachine(machine.Name); err != nil {
		return c.String(http.StatusInternalServerError, "Error eliminando máquina: "+err.Error())
	}

	ac.mu.Lock()
	delete(ac.Manager.UserMachines, machineID)
	ac.mu.Unlock()

	return c.Redirect(http.StatusSeeOther, "/")
}
