package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ConnectDisk conecta un disco a una máquina virtual
func (ac *AppContext) ConnectDisk(c echo.Context) error {
	diskID := c.Param("diskid")
	vmID := c.Param("vmid")

	ac.mu.RLock()
	disk, diskExists := ac.Manager.MediaDisks[diskID]
	ac.mu.RUnlock()

	if !diskExists {
		return c.String(http.StatusNotFound, "Disco no encontrado")
	}

	ac.mu.Lock()
	disk.IsConnected = true
	disk.ConnectedTo = vmID
	ac.mu.Unlock()

	return c.String(http.StatusOK, "Disco conectado")
}

// DisconnectDisk desconecta un disco de una máquina virtual
func (ac *AppContext) DisconnectDisk(c echo.Context) error {
	diskID := c.Param("diskid")

	ac.mu.RLock()
	disk, diskExists := ac.Manager.MediaDisks[diskID]
	ac.mu.RUnlock()

	if !diskExists {
		return c.String(http.StatusNotFound, "Disco no encontrado")
	}

	ac.mu.Lock()
	disk.IsConnected = false
	disk.ConnectedTo = ""
	ac.mu.Unlock()

	return c.String(http.StatusOK, "Disco desconectado")
}

// DeleteDisk elimina un disco multiconexión
func (ac *AppContext) DeleteDisk(c echo.Context) error {
	diskID := c.Param("diskid")

	ac.mu.RLock()
	disk, diskExists := ac.Manager.MediaDisks[diskID]
	ac.mu.RUnlock()

	if !diskExists {
		return c.String(http.StatusNotFound, "Disco no encontrado")
	}

	if disk.IsConnected {
		return c.String(http.StatusBadRequest, "El disco está conectado. Desconéctelo primero")
	}

	// Aquí se eliminaría el disco físicamente en VirtualBox
	// ac.VBoxService.DeleteDisk(disk.Location)

	ac.mu.Lock()
	delete(ac.Manager.MediaDisks, diskID)
	ac.mu.Unlock()

	return c.Redirect(http.StatusSeeOther, "/")
}
