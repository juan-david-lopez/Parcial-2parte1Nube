package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"vm-manager/models"
	"vm-manager/utils"
)

// SSHService gestiona las llaves SSH y conexiones
type SSHService struct {
	mu       sync.RWMutex
	keysPath string
}

// NewSSHService crea un nuevo servicio SSH
func NewSSHService(keysPath string) *SSHService {
	utils.EnsureDir(keysPath)
	return &SSHService{
		keysPath: keysPath,
	}
}

// CreateRootSSHKeys crea las llaves SSH para el usuario root
func (s *SSHService) CreateRootSSHKeys(machineID string, keySize int) (*models.SSHKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generar par de llaves RSA
	publicKeyStr, privateKeyStr, err := utils.GenerateRSAKeyPair(keySize)
	if err != nil {
		return nil, fmt.Errorf("error generando llaves RSA: %w", err)
	}

	// Crear estructura SSHKey
	sshKey := &models.SSHKey{
		Name:       fmt.Sprintf("root_%s", machineID),
		PublicKey:  publicKeyStr,
		PrivateKey: privateKeyStr,
		CreatedAt:  time.Now(),
	}

	// Guardar llaves en archivos
	keyDir := filepath.Join(s.keysPath, machineID, "root")
	utils.EnsureDir(keyDir)

	publicKeyPath := filepath.Join(keyDir, "id_rsa.pub")
	privateKeyPath := filepath.Join(keyDir, "id_rsa")

	if err := utils.SaveKeyToFile(publicKeyPath, publicKeyStr); err != nil {
		return nil, err
	}

	if err := utils.SaveKeyToFile(privateKeyPath, privateKeyStr); err != nil {
		return nil, err
	}

	return sshKey, nil
}

// CreateUserSSHKeys crea las llaves SSH para un usuario
func (s *SSHService) CreateUserSSHKeys(machineID, username string, keySize int) (*models.SSHKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generar par de llaves RSA
	publicKeyStr, privateKeyStr, err := utils.GenerateRSAKeyPair(keySize)
	if err != nil {
		return nil, fmt.Errorf("error generando llaves RSA: %w", err)
	}

	// Crear estructura SSHKey
	sshKey := &models.SSHKey{
		Name:       fmt.Sprintf("%s_%s", username, machineID),
		PublicKey:  publicKeyStr,
		PrivateKey: privateKeyStr,
		CreatedAt:  time.Now(),
	}

	// Guardar llaves en archivos
	keyDir := filepath.Join(s.keysPath, machineID, username)
	utils.EnsureDir(keyDir)

	publicKeyPath := filepath.Join(keyDir, "id_rsa.pub")
	privateKeyPath := filepath.Join(keyDir, "id_rsa")

	if err := utils.SaveKeyToFile(publicKeyPath, publicKeyStr); err != nil {
		return nil, err
	}

	if err := utils.SaveKeyToFile(privateKeyPath, privateKeyStr); err != nil {
		return nil, err
	}

	return sshKey, nil
}

// GetPublicKey obtiene la llave pública de un usuario
func (s *SSHService) GetPublicKey(machineID, username string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pubKeyPath := filepath.Join(s.keysPath, machineID, username, "id_rsa.pub")
	data, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return "", fmt.Errorf("error leyendo llave pública: %w", err)
	}

	return string(data), nil
}

// GetPrivateKey obtiene la llave privada de un usuario
func (s *SSHService) GetPrivateKey(machineID, username string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	privKeyPath := filepath.Join(s.keysPath, machineID, username, "id_rsa")
	data, err := os.ReadFile(privKeyPath)
	if err != nil {
		return "", fmt.Errorf("error leyendo llave privada: %w", err)
	}

	return string(data), nil
}

// DownloadPrivateKey retorna la llave privada para descargar
func (s *SSHService) DownloadPrivateKey(machineID, username string) ([]byte, error) {
	keyStr, err := s.GetPrivateKey(machineID, username)
	if err != nil {
		return nil, err
	}

	return []byte(keyStr), nil
}
