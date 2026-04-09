# Resumen del Proyecto - VM Manager

## Visión General

Este proyecto implementa una **plataforma web en Golang** para automatizar la gestión de máquinas virtuales en VirtualBox con autenticación SSH mediante llaves RSA de 1024 bits, evitando interacción manual del usuario.

**Universidad:** Universidad del Quindío  
**Proyecto:** Gestión Automatizada de Máquinas Virtuales en la Nube  
**Fecha:** 2026  

---

## Contenido Entregado

### 1. **Estructura de Proyecto Completa**

```
C:\Users\murde\Parcial#2parte1Nube\
├── main.go                      # Punto de entrada (Echo + Templates)
├── go.mod                       # Dependencias Go
├── dashboard.go                 # Conectador original (vacío, reemplazado)
│
├── models/
│   └── models.go               # Estructuras: BaseMachine, UserMachine, MediaDisk, SSHKey
│
├── handlers/
│   ├── dashboard.go            # Controlador de dashboard y máquinas
│   └── disks.go                # Controlador de discos multiconexión
│
├── services/
│   ├── virtualbox.go           # Servicio VirtualBox (comando vboxmanage)
│   └── ssh.go                  # Servicio SSH (generación y almacenamiento de llaves)
│
├── config/
│   └── config.go               # Configuración de aplicación
│
├── utils/
│   └── helpers.go              # Utilidades: RSA, validaciones, archivos
│
├── templates/
│   └── dashboard.html          # Interfaz web HTML/CSS/JavaScript
│
├── README.md                   # Documentación principal
├── DEVELOPMENT.md              # Guía de desarrollo y arquitectura
├── BUILD.md                    # Instrucciones de compilación
├── FAQ.md                      # Preguntas frecuentes
└── .gitignore                  # Control de versiones
```

---

## Características Implementadas

### ✅ 1. Gestión de Máquinas Virtuales Base

- **Crear máquinas base**: Nombre, descripción, registro automático en VirtualBox
- **Estado**: Seguimiento de si está corriendo o apagada
- **Asociación de recursos**: Vinculación de llaves y discos
- **Almacenamiento**: En memoria (VMManager - basado en mapas)

**Endpoints:**
- `POST /machines/base/add` - Crear máquina
- `GET /` - Ver dashboard

### ✅ 2. Generación de Llaves SSH RSA

**Para Root (Máquinas Base):**
- Generación de par RSA 1024-bit
- Almacenamiento en: `.vm_manager/ssh_keys/[machine_id]/root/`
- Descarga de llave privada
- Marcado de estado "RootKeysReady"

**Para Usuarios (Máquinas de Usuario):**
- Generación independiente por usuario
- Almacenamiento organizado por máquina y usuario
- Descarga segura de llaves privadas
- Permisos restrictivos: 0600 para privadas

**Endpoints:**
- `POST /machines/base/:id/keys/root` - Crear llaves root
- `GET /machines/base/:id/keys/root/download` - Descargar
- `POST /machines/user/:vmid/keys` - Crear llaves usuario
- `GET /machines/user/:vmid/keys/:username/download` - Descargar

### ✅ 3. Discos Multiconexión

- **Crear discos**: Solo después de configurar llaves root
- **Requisito previo**: Máquina base con RootKeysReady = true
- **Conexión/Desconexión**: Control de qué máquina usa el disco
- **Eliminación**: Solo si no está conectado
- **Disponibilidad**: Reutilizable para múltiples máquinas de usuario

**Endpoints:**
- `POST /machines/base/:id/disks/create` - Crear disco
- `POST /disks/:diskid/connect/:vmid` - Conectar
- `POST /disks/:diskid/disconnect` - Desconectar
- `DELETE /disks/:diskid` - Eliminar

### ✅ 4. Máquinas Virtuales de Usuario

- **Creación**: Desde disco multiconexión, con nombre y dueño
- **Clonación**: Basada en máquina virtual base
- **Asociación**: Vinculación a disco y base
- **Gestión de usuario**: Creación de usuario con llaves SSH
- **Eliminación**: Limpieza completa en VirtualBox

**Endpoints:**
- `POST /machines/user/add/:diskid` - Crear VM usuario
- `DELETE /machines/user/:vmid` - Eliminar VM

### ✅ 5. Interfaz de Dashboard Web

**Características UI:**
- Diseño responsive con CSS Grid
- 4 secciones principales:
  1. Agregar máquina base (formulario)
  2. Máquinas virtuales base (listado + controles)
  3. Discos multiconexión (listado + controles)
  4. Máquinas virtuales de usuario (listado + controles)

**Elementos interactivos:**
- Formularios validados
- Botones con estados (colores según tipo)
- Indicadores de estado (Online/Offline/Listo/Pendiente)
- Descargas de llaves
- Eliminación con confirmación

**Estilo:**
- Gradiente purpura moderno
- Cards para cada recurso
- Iconos emoji para visualización
- Responsive para móvil

### ✅ 6. Validaciones y Controles

**Validación de entrada:**
- Nombres de máquinas: No vacíos, máx 50 caracteres, alfanuméricos + guiones/guiones_bajos
- Nombres de disco: No vacíos
- Datos de usuario: Nombre de usuario requerido

**Controles lógicos:**
- No crear disco sin llaves root
- No crear máquina de usuario sin disco
- No conectar disco ya conectado
- No eliminar disco conectado
- No eliminar máquina sin existir

**Estados habilitados/deshabilitados:**
- Botones de descarga solo si existen llaves
- Formularios de creación solo si requisitos cumplidos
- Botones de eliminación deshabilitados según estado

---

## Arquitectura Técnica

### Stack Tecnológico

```
Presentación:    HTML5 + CSS3 + JavaScript Vanilla
Backend:         Golang 1.21+
Framework Web:   Echo v4 (routing, middleware)
Criptografía:    golang.org/x/crypto (RSA)
Hipervisor:      VirtualBox (vboxmanage)
Sistema:         Windows (optimizado)
Almacenamiento:  Memoria (MapasGo) + Filesystem (Llaves SSH)
```

### Patrón de Arquitectura: MVC + Servicios

```
REQUEST
  ↓
Handlers (MVC - Controller)
  ├─ Validación
  ├─ Parseo de datos
  └─ Delegación a servicios
    ↓
Services (Business Logic)
  ├─ VirtualBoxService (vboxmanage CLI)
  ├─ SSHService (generación de llaves)
  └─ Utils (helpers)
    ↓
Models (Data Layer)
  ├─ En memoria (sync.Map equivalente)
  ├─ Filesystem (llaves SSH)
  └─ VirtualBox (máquinas reales)
    ↓
RESPONSE
```

### Sincronización

- **RWMutex**: Protege acceso a VMManager
- **Patrón**: Lock para escritura, RLock para lectura
- **Concurrencia**: Thread-safe para operaciones paralelas

---

## Flujos de Negocio Implementados

### Flujo 1: Crear Máquina Base

```
1. Usuario ingresa nombre + descripción
2. POST /machines/base/add
3. Validar nombre
4. Ejecutar: vboxmanage createvm --name [name] --ostype Debian --register
5. Guardar en VMManager.BaseMachines
6. Mostrar en dashboard
```

### Flujo 2: Generar Llaves Root

```
1. Usuario clicks "Crear Llaves Root"
2. POST /machines/base/:id/keys/root
3. Verificar máquina existe
4. GenarateRSAKeyPair(1024):
   - Generar private key con crypto/rsa
   - Encodear a PEM
   - Generar public key
   - Encod with crypto/ssh
5. SaveKeyToFile() a ~/.vm_manager/ssh_keys/:id/root/
6. Marcar RootKeysReady = true
7. Volver a dashboard
```

### Flujo 3: Crear Disco Multiconexión

```
1. User clicks "Crear Disco Multiconexión" (en máquina con llaves root)
2. POST /machines/base/:id/disks/create
3. Validar:
   - Máquina existe
   - RootKeysReady = true
   - Nombre válido
4. Ejecutar: vboxmanage createmedium disk --filename [path] --size [MB]
5. Crear objeto MediaDisk
6. Guardar en VMManager.MediaDisks
7. Asociar con máquina base
8. Mostrar en dashboard
```

### Flujo 4: Crear Máquina de Usuario

```
1. Usuario ingresa nombre + owner + descripción + elige disco
2. POST /machines/user/add/:diskid
3. Validar:
   - Disco existe
   - No está conectado
4. Ejecutar: vboxmanage clonevm [base] --name [newname] --register
5. Crear objeto UserMachine
6. Asociar con disco y máquina base
7. Guardar en VMManager
8. Mostrar en dashboard
```

### Flujo 5: Crear Llaves de Usuario

```
1. User clicks "Crear Llaves Usuario"
2. Ingresa nombre de usuario
3. POST /machines/user/:vmid/keys
4. GenerateRSAKeyPair(1024)
5. SaveKeyToFile() a ~/.vm_manager/ssh_keys/:vmid/:username/
6. Marcar UserMachine.KeysReady = true
7. Return "Éxito"
```

---

## Casos de Uso Cubiertos

### Caso 1: Dashboard ✅
- Visualizar máquinas base, discos y máquinas de usuario
- Mostrar estado de cada recurso
- Proporcionar controles contextuales

### Caso 2: Agregar Máquinas Base ✅
- Formulario con nombre y descripción
- Crear llaves de root (botón)
- Descargar llaves de root (descarga)
- Crear disco multiconexión (botón)

### Caso 3: Crear Discos Multiconexión ✅
- Requisito: Llaves root configuradas
- Crear disco (botón)
- Gestionar acciones: conectar, desconectar, eliminar

### Caso 4: Crear Máquinas de Usuario ✅
- Basadas en disco multiconexión
- Crear usuario con llaves SSH (botón)
- Descargar llaves de usuario (descarga)
- Eliminar máquina (botón)

---

## Validaciones Implementadas

| Elemento | Validación |
|----------|-----------|
| Nombre máquina | No vacío, máx 50 car, alfanuméricos + - _ |
| Nombre disco | No vacío |
| Nombre usuario | No vacío |
| Llaves RSA | 1024 bits |
| Permisos archivo | 0600 privada, 0644 pública |
| Estados transición | Máquina → necesita llaves → disco → VM usuario |
| Eliminación | No disco conectado, máquina puede estar running |

---

## Seguridad Implementada

✅ Validación de entrada  
✅ Permisos de archivo restrictivos  
✅ Llaves en carpeta aislada  
✅ Sin exposición de rutas internas  
✅ Sin inyección SQL (sin BD)  
✅Sincronización thread-safe  

⚠️ No implementado (futuro):
- Autenticación web
- HTTPS
- Encriptación de llaves
- Rate limiting
- Logging detallado

---

## Archivos Entregados

### Código Fuente (5 carpetas)
- **models/**: 1 archivo (structs)
- **handlers/**: 2 archivos (controladores)
- **services/**: 2 archivos (VirtualBox + SSH)
- **config/**: 1 archivo (configuración)
- **utils/**: 1 archivo (utilidades)

### Archivos de Configuración
- **main.go**: Punto de entrada
- **go.mod**: Dependencias
- **.gitignore**: Control de versiones

### Interfaz
- **templates/dashboard.html**: UI completa

### Documentación (4 archivos)
- **README.md**: Guía de uso
- **DEVELOPMENT.md**: Arquitectura y desarrollo
- **BUILD.md**: Compilación y deployment
- **FAQ.md**: Preguntas frecuentes

---

## Cómo Usar

### 1. Preparación
```bash
# Descargar dependencias
go mod download

# Compilar
go build -o vm-manager.exe
```

### 2. Ejecución
```bash
# Ejecutar
.\vm-manager.exe

# Abrir navegador
http://localhost:8080
```

### 3. Flujo de Trabajo
1. Crear máquina base
2. Crear llaves de root
3. Crear disco multiconexión
4. Crear máquina de usuario
5. Crear llaves de usuario
6. Descargar llaves para SSH remoto

---

## Próximas Mejoras (Roadmap)

### Fase 2: Persistencia
- [ ] SQLite para guardar estado
- [ ] Migración de datos
- [ ] Backup/Restore

### Fase 3: Autenticación
- [ ] Login web
- [ ] Roles (admin, usuario)
- [ ] Auditoría de acciones

### Fase 4: Monitoreo
- [ ] Estado en tiempo real
- [ ] Logs de operaciones
- [ ] Alertas

### Fase 5: Infraestructura
- [ ] HTTPS /SSL
- [ ] Docker
- [ ] API REST separada

---

## Testing

```bash
# Tests unitarios (estructura preparada)
go test ./...

# Coverage
go test -cover ./...
```

---

## Despliegue

### Desarrollo
```bash
go run main.go
```

### Producción
```powershell
go build -ldflags "-s -w" -o vm-manager.exe
# Copiar a C:\Program Files\VMManager\
# Registrar en tareas programadas
```

---

## Conclusiones

✅ **Completado**: Plataforma web funcional para gestionar VirtualBox en Windows  
✅ **Código limpio**: Estructura modular y reutilizable  
✅ **Documentado**: Guías técnicas y de usuario  
✅ **Escalable**: Base sólida para expansión  
✅ **Seguro**: Validaciones y prácticas de seguridad básicas  

---

## Contacto / Autor

**Proyecto**: VM Manager - Universidad del Quindío  
**Lenguaje**: Golang  
**Versión**: 1.0.0  
**Fecha**: 2026  

---

## Licencia

Proyecto educativo - Libre para uso y modificación.
