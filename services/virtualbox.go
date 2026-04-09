package services

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"vm-manager/models"
)

// VirtualBoxService gestiona las operaciones de VirtualBox
type VirtualBoxService struct {
	mu sync.RWMutex
	// Almacenaremos referencias a máquinas
}

// CreateBaseMachine crea una nueva máquina virtual base
func (vb *VirtualBoxService) CreateBaseMachine(name, description string) (*models.BaseMachine, error) {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	// Validar nombre
	if err := validateMachineName(name); err != nil {
		return nil, err
	}

	machine := &models.BaseMachine{
		ID:          generateID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		IsRunning:   false,
		Disks:       make(map[string]*models.MediaDisk),
		Users:       make(map[string]*models.UserMachine),
		Port:        2200,
	}

	// Crear máquina en VirtualBox
	if err := vb.createVMInVirtualBox(name, "Debian"); err != nil {
		return nil, fmt.Errorf("error creando máquina en VirtualBox: %w", err)
	}

	return machine, nil
}

// CreateMediaDisk crea un disco multiconexión
func (vb *VirtualBoxService) CreateMediaDisk(name string, sizeInMB int64, baseMachineID string) (*models.MediaDisk, error) {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	if name == "" {
		return nil, fmt.Errorf("el nombre del disco no puede estar vacío")
	}

	disk := &models.MediaDisk{
		ID:          generateID(),
		Name:        name,
		Size:        sizeInMB,
		CreatedAt:   time.Now(),
		IsConnected: false,
		BaseMachine: baseMachineID,
		IsReady:     true,
		Location:    fmt.Sprintf("/path/to/disk/%s.vdi", name),
	}

	// Crear disco en VirtualBox
	if err := vb.createVDIDisk(disk.Location, sizeInMB); err != nil {
		return nil, fmt.Errorf("error creando disco VDI: %w", err)
	}

	return disk, nil
}

// CreateUserMachine crea una máquina virtual de usuario
func (vb *VirtualBoxService) CreateUserMachine(name, description, baseMachineID, diskID, owner string) (*models.UserMachine, error) {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	if err := validateMachineName(name); err != nil {
		return nil, err
	}

	machine := &models.UserMachine{
		ID:           generateID(),
		Name:         name,
		Description:  description,
		CreatedAt:    time.Now(),
		IsRunning:    false,
		Owner:        owner,
		BaseMachine:  baseMachineID,
		AttachedDisk: diskID,
		Port:         2201,
		KeysReady:    false,
	}

	// Crear nueva VM clonando from base
	if err := vb.cloneVM(baseMachineID, name); err != nil {
		return nil, fmt.Errorf("error clonando máquina: %w", err)
	}

	return machine, nil
}

// StartMachine inicia una máquina virtual
func (vb *VirtualBoxService) StartMachine(machineName string) error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	cmd := exec.Command("vboxmanage", "startvm", machineName, "--type", "headless")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error iniciando máquina: %w", err)
	}

	return nil
}

// StopMachine detiene una máquina virtual
func (vb *VirtualBoxService) StopMachine(machineName string) error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	cmd := exec.Command("vboxmanage", "controlvm", machineName, "poweroff")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error deteniendo máquina: %w", err)
	}

	return nil
}

// DeleteMachine elimina una máquina virtual
func (vb *VirtualBoxService) DeleteMachine(machineName string) error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	// Detener máquina si está corriendo
	exec.Command("vboxmanage", "controlvm", machineName, "poweroff").Run()

	cmd := exec.Command("vboxmanage", "unregistervm", machineName, "--delete")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error eliminando máquina: %w", err)
	}

	return nil
}

// AttachDisk conecta un disco a una máquina virtual
func (vb *VirtualBoxService) AttachDisk(machineName string, diskPath string, controller string, port int, device int) error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	cmd := exec.Command("vboxmanage", "storageattach", machineName,
		"--storagectl", controller,
		"--port", fmt.Sprintf("%d", port),
		"--device", fmt.Sprintf("%d", device),
		"--type", "hdd",
		"--medium", diskPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error adjuntando disco: %w", err)
	}

	return nil
}

// DetachDisk desconecta un disco de una máquina virtual
func (vb *VirtualBoxService) DetachDisk(machineName string, controller string, port int, device int) error {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	cmd := exec.Command("vboxmanage", "storagedetach", machineName,
		"--storagectl", controller,
		"--port", fmt.Sprintf("%d", port),
		"--device", fmt.Sprintf("%d", device))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error desconectando disco: %w", err)
	}

	return nil
}

// Private methods

func (vb *VirtualBoxService) createVMInVirtualBox(name, osType string) error {
	cmd := exec.Command("vboxmanage", "createvm", "--name", name, "--ostype", osType, "--register")
	return cmd.Run()
}

func (vb *VirtualBoxService) createVDIDisk(location string, sizeInMB int64) error {
	cmd := exec.Command("vboxmanage", "createmedium", "disk", "--filename", location, "--size", fmt.Sprintf("%d", sizeInMB))
	return cmd.Run()
}

func (vb *VirtualBoxService) cloneVM(baseMachineID, newMachineName string) error {
	cmd := exec.Command("vboxmanage", "clonevm", baseMachineID, "--name", newMachineName, "--register")
	return cmd.Run()
}

func validateMachineName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("el nombre no puede estar vacío")
	}
	if len(name) > 50 {
		return fmt.Errorf("el nombre es demasiado largo")
	}
	return nil
}

func generateID() string {
	// Implementar usando utils
	return fmt.Sprintf("vm_%d", time.Now().UnixNano())
}
