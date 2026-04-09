package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	// Server
	Port string
	Host string

	// VirtualBox
	VBoxPath         string
	VBoxDisksPath    string
	VBoxMachinesPath string

	// SSH
	SSHKeySize int
	SSHKeyType string

	// Application
	Templates  string
	PublicPath string
}

// Load carga la configuración del sistema
func Load() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		Port:             ":8080",
		Host:             "localhost",
		VBoxPath:         filepath.Join(homeDir, ".VirtualBox"),
		VBoxMachinesPath: filepath.Join(homeDir, "VirtualBox VMs"),
		VBoxDisksPath:    filepath.Join(homeDir, "VirtualBox VMs", "Disks"),
		SSHKeySize:       1024,
		SSHKeyType:       "rsa",
		Templates:        "./templates",
		PublicPath:       "./public",
	}
}

// GetSSHKeyPath retorna la ruta donde se almacenarán las llaves SSH
func (c *Config) GetSSHKeyPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".vm_manager", "ssh_keys")
}
