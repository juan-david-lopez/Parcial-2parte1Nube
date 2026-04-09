# Guía de Desarrollo - Gestor de Máquinas Virtuales

## Arquitectura General

La aplicación sigue una arquitectura en capas:

```
┌─────────────────────────────────────┐
│      Interfaz Web (HTML/CSS)        │
├─────────────────────────────────────┤
│   Controladores (handlers)          │  ← Rutas HTTP, validación básica
├─────────────────────────────────────┤
│   Servicios (services)              │  ← Lógica de negocio
├─────────────────────────────────────┤
│   Modelos (models)                  │  ← Estructuras de datos
└─────────────────────────────────────┘
       ↓              ↓              ↓
    VirtualBox     SSH Keys       Filesystem
```

## Componentes

### 1. **Models** (`models/models.go`)

Define las estructuras de datos principales:

- `BaseMachine`: Máquina virtual base para clonar
- `UserMachine`: Máquina virtual creada para un usuario
- `MediaDisk`: Disco multiconexión compartido
- `SSHKey`: Par de llaves RSA
- `VMManager`: Contenedor de todas las máquinas y discos

### 2. **Handlers** (`handlers/`)

Gestiona las peticiones HTTP:

- `handlers/dashboard.go`: Controlador principal
  - `GetDashboard()`: Renderiza el dashboard
  - `AddBaseMachine()`: Crear máquina base
  - `CreateRootKeys()`: Generar llaves de root
  - `AddUserMachine()`: Crear máquina de usuario
  - `CreateUserSSHKeys()`: Generar llaves de usuario

- `handlers/disks.go`: Gestión de discos
  - `ConnectDisk()`: Conectar disco a máquina
  - `DisconnectDisk()`: Desconectar disco
  - `DeleteDisk()`: Eliminar disco

### 3. **Services** (`services/`)

Contiene la lógica de negocio:

- `services/virtualbox.go`: `VirtualBoxService`
  - Envuelve comandos `vboxmanage`
  - Crear/eliminar máquinas
  - Gestionar discos
  - Control de máquinas (iniciar, detener)

- `services/ssh.go`: `SSHService`
  - Generar llaves RSA
  - Guardar llaves en archivos
  - Descargar llaves privadas

### 4. **Utils** (`utils/helpers.go`)

Funciones auxiliares:

- `GenerateRSAKeyPair()`: Generar par de llaves
- `SaveKeyToFile()`: Guardar llaves con permisos
- `GenerateID()`: ID único para entidades
- `ValidateMachineName()`: Validar nombres
- `FileExists()`: Verificar archivos
- `EnsureDir()`: Crear directorios

### 5. **Config** (`config/config.go`)

Configuración de la aplicación:

- Rutas de VirtualBox
- Rutas de plantillas
- Tamaño de llaves SSH
- Puerto del servidor

## Flujo de Datos

### Crear Máquina Base

```
GET /  (mostrar dashboard)
  ↓
User llena formulario "Agregar Máquina Virtual Base"
  ↓
POST /machines/base/add
  ↓
Handler → validar input
  ↓
VBoxService → CreateBaseMachine()
  ↓
Ejecutar: vboxmanage createvm
  ↓
VMManager → guardar en BaseMachines
  ↓
Redirect → GET / (actualizar dashboard)
```

### Crear Llaves de Root

```
POST /machines/base/:id/keys/root
  ↓
Handler → validar máquina existe
  ↓
SSHService → CreateRootSSHKeys()
  ↓
GenerateRSAKeyPair() → generar llaves
  ↓
SaveKeyToFile() → guardar en ~/.vm_manager/ssh_keys/
  ↓
VMManager → marcar RootKeysReady = true
  ↓
Retornar nombre de archivo o error
```

### Crear Disco Multiconexión

```
POST /machines/base/:id/disks/create
  ↓
Handler → validar:
  - Máquina existe
  - RootKeysReady = true
  - Nombre del disco válido
  ↓
VBoxService → CreateMediaDisk()
  ↓
Ejecutar: vboxmanage createmedium
  ↓
VMManager → guardar disco
            → asociar con BaseMachine
  ↓
Redirect → actualizar dashboard
```

### Crear Máquina de Usuario

```
POST /machines/user/add/:diskid
  ↓
Handler → validar:
  - Disco existe
  - Disco no está conectado
  - Datos de usuario válidos
  ↓
VBoxService → CreateUserMachine()
  ↓
Ejecutar: vboxmanage clonevm (clonar base)
  ↓
VMManager → guardar máquina
           → asociar disco
           → asociar con base
  ↓
Redirect → actualizar dashboard
```

## Extensiones Comunes

### Agregar Nueva Funcionalidad

1. **Agregar un Handler**:

```go
// En handlers/dashboard.go o un archivo nuevo
func (ac *AppContext) NuevaFuncion(c echo.Context) error {
    // Obtener parámetros
    param := c.Param("id")
    
    // Validar
    ac.mu.RLock()
    item, exists := ac.Manager.Items[param]
    ac.mu.RUnlock()
    
    if !exists {
        return c.String(http.StatusNotFound, "No encontrado")
    }
    
    // Usar servicio
    err := ac.VBoxService.AlgunaOperacion(item)
    if err != nil {
        return c.String(http.StatusInternalServerError, err.Error())
    }
    
    // Actualizar estado
    ac.mu.Lock()
    item.Status = "updated"
    ac.mu.Unlock()
    
    // Retornar
    return c.String(http.StatusOK, "Éxito")
}
```

2. **Registrar la Ruta** (en `main.go`):

```go
e.GET("/ruta/:id", appCtx.NuevaFuncion)
```

3. **Agregar Controles HTML** (en `templates/dashboard.html`):

```html
<form method="POST" action="/ruta/id">
    <button type="submit" class="btn btn-primary">Ejecutar</button>
</form>
```

### Agregar Nueva Propiedad a Máquina

1. Modificar `models/models.go`:

```go
type BaseMachine struct {
    // ... existentes ...
    NuevaPropiedad string `json:"nueva_propiedad"`
}
```

2. Usar en handler:

```go
machine.NuevaPropiedad = value
```

3. Mostrar en template:

```html
<p>{{.NuevaPropiedad}}</p>
```

### Agregar Nuevo Servicio VirtualBox

1. Crear método en `VirtualBoxService`:

```go
func (vb *VirtualBoxService) NuevaOperacion(param string) error {
    vb.mu.Lock()
    defer vb.mu.Unlock()
    
    cmd := exec.Command("vboxmanage", "comando", param)
    return cmd.Run()
}
```

2. Usar en handler:

```go
if err := ac.VBoxService.NuevaOperacion(param); err != nil {
    return c.String(http.StatusInternalServerError, err.Error())
}
```

## Sincronización y Concurrencia

La aplicación usa `sync.RWMutex` para proteger:

- Acceso a `ac.Manager` (lectura frecuente, escritura ocasional)
- Operaciones de VirtualBox (escritura frecuente)
- Operaciones SSH (lectura y escritura de archivos)

### Patrón de uso:

```go
// Lectura
ac.mu.RLock()
item := ac.Manager.Items[id]
ac.mu.RUnlock()

// Escritura
ac.mu.Lock()
ac.Manager.Items[id] = newItem
delete(ac.Manager.Items, oldId)
ac.mu.Unlock()
```

## Manejo de Errores

La aplicación retorna códigos HTTP estándar:

| Código | Significado |
|--------|------------|
| 200 | Éxito |
| 400 | Bad Request (validación) |
| 404 | Not Found |
| 500 | Server Error |

## Testing

Para agregar tests:

```go
// machine_test.go
package handlers

import (
    "testing"
    // ...
)

func TestAddBaseMachine(t *testing.T) {
    // Preparar
    appCtx := NewAppContext(/* ... */)
    
    // Ejecutar
    result := appCtx.AddBaseMachine(testContext)
    
    // Verificar
    if result != nil {
        t.Errorf("Error: %v", result)
    }
}
```

## Debugging

### Habilitar logs

En `main.go`:

```go
e.Debug = true
```

### Ver estado de máquinas

```bash
vboxmanage list vms
vboxmanage list runningvms
```

### Ver llaves creadas

```bash
ls %USERPROFILE%\.vm_manager\ssh_keys\
```

## Performance

- `RWMutex` para proteger lectura concurrente
- Comando vboxmanage ejecutados en segundo plano
- Templates HTML pre-compilados
- Minimizar bloqueos con `defer Unlock()`

## Dependencias Clave

- **Echo v4**: Framework web
- **golang.org/x/crypto**: Cifrado SSH
- **Standard library**: Manejo de procesos, filesystems

## Notas Importantes

1. Las máquinas virtuales deben estar apagadas antes de ciertas operaciones
2. VirtualBox debe estar instalado y `vboxmanage` accessible
3. Las llaves privadas se almacenan localmente - protege esta carpeta
4. Los puertos SSH se asignan secuencialmente
5. Los nombres deben cumplir validaciones

## Mejoras Futuras

- [ ] Persistencia en BD (SQLite/PostgreSQL)
- [ ] Autenticación para la web
-[ ] Historial de operaciones
- [ ] Notificaciones en tiempo real
- [ ] Métricas de uso
- [ ] Copias de seguridad automáticas
- [ ] CLI adicional además de web
