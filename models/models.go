package models

import "time"

// MachineType define los tipos de máquinas virtuales
type MachineType string

const (
	BaseMachine MachineType = "base"
	UserMachine MachineType = "user"
)

// BaseMachine representa una máquina virtual base
type BaseMachine struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	CreatedAt     time.Time               `json:"created_at"`
	IsRunning     bool                    `json:"is_running"`
	RootKeysReady bool                    `json:"root_keys_ready"`
	Disks         map[string]*MediaDisk   `json:"disks"`
	Users         map[string]*UserMachine `json:"users"`
	Port          int                     `json:"port"`
}

// MediaDisk representa un disco multiconexión
type MediaDisk struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"` // en MB
	CreatedAt   time.Time `json:"created_at"`
	IsConnected bool      `json:"is_connected"`
	ConnectedTo string    `json:"connected_to"` // ID de la máquina a la que está conectado
	BaseMachine string    `json:"base_machine"` // ID de la máquina base
	IsReady     bool      `json:"is_ready"`
	Location    string    `json:"location"`
}

// UserMachine representa una máquina virtual de usuario
type UserMachine struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	IsRunning    bool      `json:"is_running"`
	Owner        string    `json:"owner"`
	BaseMachine  string    `json:"base_machine"`  // ID de la máquina base de la que fue clonada
	AttachedDisk string    `json:"attached_disk"` // ID del disco multiconexión
	UserPassword string    `json:"user_password"` // Contraseña del usuario
	Port         int       `json:"port"`
	KeysReady    bool      `json:"keys_ready"`
}

// SSHKey representa una par de llaves RSA
type SSHKey struct {
	Name       string    `json:"name"`
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key"`
	CreatedAt  time.Time `json:"created_at"`
}

// VMManager es el gestor principal de máquinas virtuales
type VMManager struct {
	BaseMachines map[string]*BaseMachine
	UserMachines map[string]*UserMachine
	MediaDisks   map[string]*MediaDisk
	SSHKeys      map[string]*SSHKey
	RootKeys     map[string]*SSHKey
}

// Dashboard contiene la información del tablero
type Dashboard struct {
	BaseMachines []BaseMachine
	MediaDisks   []MediaDisk
	UserMachines []UserMachine
}
