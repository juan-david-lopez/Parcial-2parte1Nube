package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

// GenerateRSAKeyPair genera un par de llaves RSA
func GenerateRSAKeyPair(bits int) (string, string, error) {
	// Generar llave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("error generando llave privada: %w", err)
	}

	// Codificar llave privada
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Generar llave pública
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("error generando llave pública: %w", err)
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

	return string(publicKeyBytes), string(privateKeyPEM), nil
}

// SaveKeyToFile guarda una llave en un archivo
func SaveKeyToFile(filePath string, content string) error {
	// Crear directorio si no existe
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	// Escribir archivo
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creando archivo: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	// Establecer permisos de archivo
	if strings.Contains(filePath, "private") {
		os.Chmod(filePath, 0600)
	} else {
		os.Chmod(filePath, 0644)
	}

	return nil
}

// GenerateID crea un ID único
func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// ValidateMachineName valida el nombre de una máquina
func ValidateMachineName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("el nombre de la máquina no puede estar vacío")
	}
	if len(name) > 50 {
		return fmt.Errorf("el nombre de la máquina no puede exceder 50 caracteres")
	}
	// Solo caracteres alfanuméricos, guiones y guiones bajos
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return fmt.Errorf("el nombre contiene caracteres no permitidos")
		}
	}
	return nil
}

// FileExists verifica si un archivo existe
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureDir crea un directorio si no existe
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
