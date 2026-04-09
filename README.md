# Gestor de Máquinas Virtuales en Golang

Plataforma web para automatizar la gestión de máquinas virtuales en VirtualBox con autenticación SSH basada en llaves RSA.

## Características

- **Gestión de Máquinas Base**: Crear y gestionar máquinas virtuales base
- **Creación de Llaves SSH**: Generar llaves RSA (1024 bits) para root y usuarios
- **Discos Multiconexión**: Crear y gestionar discos que pueden conectarse a múltiples VMs
- **Máquinas de Usuario**: Crear máquinas virtuales a partir de discos multiconexión
- **Control Remoto SSH**: Acceso remoto sin interacción del usuario mediante llaves
- **Interfaz Web**: Dashboard intuitivo con controles habilitados/deshabilitados según estado

## Requisitos Previos

1. **Go 1.21 o superior**
2. **VirtualBox instalado y configurado**
3. **SSH configurado en las máquinas virtuales base**
4. **Contraseña de root de las máquinas virtuales (conocida)**

## Estructura del Proyecto

```
vm-manager/
├── main.go                     # Punto de entrada
├── go.mod                      # Dependencias del módulo
├── models/
│   └── models.go              # Estructuras de datos
├── config/
│   └── config.go              # Configuración de la aplicación
├── handlers/
│   ├── dashboard.go           # Controlador del dashboard
│   └── disks.go               # Controlador de discos
├── services/
│   ├── virtualbox.go          # Servicio de VirtualBox
│   └── ssh.go                 # Servicio SSH
├── utils/
│   └── helpers.go             # Funciones auxiliares
└── templates/
    └── dashboard.html         # Interfaz web
```

## Instalación

### 1. Clonar o descargar el proyecto

```bash
cd C:\Users\[usuario]\Parcial#2parte1Nube
```

### 2. Instalar dependencias

```bash
go mod download
go mod tidy
```

### 3. Compilar la aplicación

```bash
go build -o vm-manager.exe
```

O ejecutar directamente sin compilar:

```bash
go run main.go
```

## Uso

### Iniciar el servidor

```bash
# Compilado
.\vm-manager.exe

# O directamente
go run main.go
```

El servidor estará disponible en: **http://localhost:8080**

### Configuración Previa

1. **Prepare máquinas virtuales base en VirtualBox**:
   - Al menos 2 máquinas virtuales
   - Distribuciones Debian derivadas diferentes
   - SSH configurado y funcionando
   - Discos tipo normal
   - Apagadas antes de usar

2. **Conozca las contraseñas de root** de cada máquina

3. **Rutas de VirtualBox**:
   - Máquinas: `C:\Users\[usuario]\VirtualBox VMs`
   - Discos: Se almacenarán en la carpeta configurable

## Flujo de Trabajo

### 1. Crear Máquina Virtual Base

1. Ve a la sección "Agregar Máquina Virtual Base"
2. Ingresa nombre y descripción
3. Haz clic en "Crear Máquina Base"
4. En la sección de máquinas base, haz clic en "Crear Llaves Root"
5. Descarga las llaves si lo deseas

### 2. Crear Disco Multiconexión

1. En la máquina base con llaves creadas, haz clic en "Crear Disco Multiconexión"
2. Ingresa el nombre del disco
3. El disco se creará y estará disponible en la sección "Discos Multiconexión"

### 3. Crear Máquina Virtual de Usuario

1. En la sección "Discos Multiconexión", selecciona un disco disponible
2. Ingresa: nombre de máquina, descripción y dueño
3. Haz clic en "Crear VM Usuario"
4. La máquina aparecerá en "Máquinas Virtuales de Usuario"
5. Crea llaves de usuario (genera llaves SSH para el usuario)
6. Descarga las llaves si lo deseas

### 4. Gestionar Discos

- **Conectar/Desconectar**: Usa los botones en la tarjeta del disco
- **Eliminar**: Solo es posible si el disco no está conectado

## Controles de Seguridad

- ✅ Validación de nombres de máquinas
- ✅ Llaves RSA de 1024 bits
- ✅ Permisos de archivo restrictivos (0600 para privadas)
- ✅ Evita eliminación de elementos en uso
- ✅ Control de estados y transiciones

## Consideraciones Importantes

1. **Llaves SSH**:
   - Se almacenan localmente en: `%USERPROFILE%\.vm_manager\ssh_keys\`
   - Protege estas carpetas
   - Descarga las llaves privadas para backup

2. **Máquinas Virtuales**:
   - Las máquinas deben estar apagadas antes de operaciones
   - El sistema no inicia máquinas automáticamente
   - Los puertos SSH se asignan incrementalmente (2200, 2201, etc.)

3. **Discos Multiconexión**:
   - Son independientes de se conecten a máquinas
   - Se pueden compartir entre múltiples usuarios
   - El disco debe estar desconectado antes de eliminarlo

## Troubleshooting

### Error: "VirtualBox no encontrado"
- Instala VirtualBox desde https://www.virtualbox.org/
- Asegúrate de que `vboxmanage` esté en el PATH

### Error: "Puerto SSH en uso"
- El puerto puede estar siendo usado por otra aplicación
- Aumenta el número de puerto manualmente

### No se cargan las templates
- Verifica que existe la carpeta `templates/`
- Asegúrate de que `dashboard.html` está en la carpeta

## Próximas Mejoras

- [ ] Base de datos para persistencia
- [ ] Autenticación de usuarios para la web
- [ ] Monitoreo de máquinas virtuales en tiempo real
- [ ] Copias de seguridad automáticas
- [ ] API REST completa
- [ ] Interfaz frontend mejorada
- [ ] Soporte para más formatos de disco

## Contribuciones

Este proyecto es educativo. Siéntete libre de modificarlo y mejorarlo.

## Licencia

Proyecto educativo - Universidad del Quindío
