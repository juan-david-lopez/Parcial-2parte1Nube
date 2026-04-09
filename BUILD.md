# Instrucciones de Compilación y Deployment

## Compilación Rápida

### Windows PowerShell

```powershell
# Descargar dependencias
go mod download

# Compilar
go build -o vm-manager.exe

# Ejecutar
.\vm-manager.exe
```

### Compilación con Variables de Construcción

```powershell
# Incluir versión y fecha de compilación
$version = "1.0.0"
$buildDate = (Get-Date).ToString("2006-01-02")
go build -ldflags "-X main.Version=$version -X main.BuildDate=$buildDate" -o vm-manager.exe
```

### Compilación Optimizada (Producción)

```powershell
# Compilación optimizada, sin información de debug
go build -ldflags "-s -w" -o vm-manager.exe

# El ejecutable será más pequeño (~40% menos)
```

## Instalación de Dependencias

### Primera vez

```powershell
go mod init vm-manager
go get github.com/labstack/echo/v4
go get golang.org/x/crypto
go mod tidy
```

### Actualizar dependencias

```powershell
go get -u
go mod tidy
```

## Ejecución

### Modo desarrollo

```powershell
# Con recarga automática (requiere 'air')
go install github.com/cosmtrek/air@latest
air

# O sin air
go run main.go
```

### Modo producción

```powershell
# Compilado
.\vm-manager.exe

# En segundo plano (PS v5+)
Start-Process .\vm-manager.exe -WindowStyle Minimized
```

## Verificación de Instalación

```powershell
# Verificar Go
go version

# Verificar VirtualBox
vboxmanage --version

# Verificar SSH
ssh -V

# Probar compilación
go build -v
```

## Estructura de Carpetas en Producción

```
C:\Program Files\VMManager\
├── vm-manager.exe
├── config.yaml (opcional)
├── templates\
│   └── dashboard.html
└── logs\
    └── access.log
```

## Variables de Entorno (Opcional)

```powershell
# Configurar variables de entorno
$env:VBOX_PATH = "C:\Program Files\Oracle\VirtualBox"
$env:SSH_KEY_PATH = "C:\Users\[usuario]\.vm_manager\ssh_keys"

# Persistir variables
[Environment]::SetEnvironmentVariable("VBOX_PATH", "C:\Program Files\Oracle\VirtualBox", "User")
```

## Troubleshooting de Compilación

### Error: "GO111MODULE must be enabled"

```powershell
$env:GO111MODULE = "on"
```

### Error: "command not found: vboxmanage"

VirtualBox no está en el PATH. Instala desde:
https://www.virtualbox.org/wiki/Downloads

### Port ya en uso

```powershell
# Encontrar proceso en puerto 8080
netstat -ano | findstr :8080

# Matar proceso
taskkill /PID [PID] /F
```

## Testing

### Ejecutar tests

```powershell
go test ./...
```

### Tests verbosos

```powershell
go test -v ./...
```

### Cobertura de tests

```powershell
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Build Script (build.ps1)

```powershell
# Build.ps1
param(
    [string]$Version = "1.0.0",
    [string]$Output = "vm-manager.exe"
)

Write-Host "Building VM Manager v$Version..."

# Limpiar build anterior
Remove-Item $Output -ErrorAction SilentlyContinue

# Verificar dependencias
Write-Host "Downloading dependencies..."
go mod download

# Compilar
Write-Host "Compiling..."
$buildDate = (Get-Date).ToString("2006-01-02")
go build -ldflags "-X main.Version=$Version -X main.BuildDate=$buildDate" -o $Output

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful: $Output" -ForegroundColor Green
    $size = (Get-Item $Output).Length / 1MB
    Write-Host "Size: $([Math]::Round($size, 2)) MB"
} else {
    Write-Host "✗ Build failed!" -ForegroundColor Red
    exit 1
}
```

Uso:
```powershell
.\build.ps1 -Version "1.0.0"
```

## Deployment Checklist

- [ ] Compilar con `go build`
- [ ] Copiar `vm-manager.exe` a carpeta de instalación
- [ ] Copiar carpeta `templates/` 
- [ ] Verificar ruta de VirtualBox
- [ ] Crear carpeta `.vm_manager` en USERPROFILE
- [ ] Asignar permisos de lectura/escritura
- [ ] Probar con `vm-manager.exe`
- [ ] Verificar acceso a http://localhost:8080
- [ ] Validar que VirtualBox es accesible
- [ ] Hacer backup de llaves SSH

## Automatizar Inicio

### Crear atajo en Inicio

```powershell
$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\VMManager.lnk")
$Shortcut.TargetPath = "C:\Program Files\VMManager\vm-manager.exe"
$Shortcut.Save()
```

### Tarea programada (opcional)

```powershell
# Registrar tarea para iniciar al boot
$Principal = New-ScheduledTaskPrincipal -UserId $env:USERNAME -RunLevel Highest
$Trigger = New-ScheduledTaskTrigger -AtStartup
$Action = New-ScheduledTaskAction -Execute "C:\Program Files\VMManager\vm-manager.exe"
Register-ScheduledTask -TaskName "VMManager" -Principal $Principal -Trigger $Trigger -Action $Action
```

## Monitoreo

### Logs de Echo

```powershell
# Ver logs en tiempo real
Get-Content logs\access.log -Wait
```

### Monitoreo de Puerto

```powershell
# Verificar puerto abierto
Test-NetConnection -ComputerName localhost -Port 8080
```

## Backup

```powershell
# Backup de llaves SSH
$backupPath = "C:\Backups\SSH_Keys_$(Get-Date -f 'yyyyMMdd_HHmmss')"
Copy-Item "$env:USERPROFILE\.vm_manager\ssh_keys" $backupPath -Recurse

Write-Host "Backup created: $backupPath"
```

## Notas de Versión

- **v1.0.0**: Release inicial
  - Gestión de máquinas base
  - Generación de llaves SSH
  - Creación de discos multiconexión
  - Dashboard web básico
