# Preguntas Frecuentes

## General

### ¿Qué es VM Manager?
VM Manager es una plataforma web desarrollada en Golang que automatiza la gestión de máquinas virtuales en VirtualBox. Permite crear máquinas virtuales, generar llaves SSH RSA, crear discos multiconexión y gestionar acceso remoto sin entrada de usuario.

### ¿Qué versión de Go necesito?
Mínimo Go 1.21 o superior. Puedes descargarla desde https://golang.org/dl/

### ¿Funciona en Linux o Mac?
El código está escrito de forma portátil en Go, pero algunos comandos de VirtualBox pueden variar. El PATH de configuración está optimizado para Windows.

### ¿Es gratis?
Sí, es un proyecto educativo de la Universidad del Quindío.

---

## Instalación

### No me aparece el comando `vboxmanage`
- Instala VirtualBox desde https://www.virtualbox.org/wiki/Downloads
- Verifica que está en `C:\Program Files\Oracle\VirtualBox`
- O agrega la ruta a tu PATH de Windows

### ¿Dónde se guardan las llaves SSH?
En: `C:\Users\[TuUsuario]\.vm_manager\ssh_keys\`

Estructura:
```
.vm_manager/
└── ssh_keys/
    ├── [machine_id]/
    │   └── root/
    │       ├── id_rsa (privada)
    │       └── id_rsa.pub (pública)
    └── [another_machine]/
        └── username/
            ├── id_rsa
            └── id_rsa.pub
```

### ¿Cómo cambio el puerto del servidor?
En `main.go`, cambia:
```go
e.Logger.Fatal(e.Start(":8080"))
```

A:
```go
e.Logger.Fatal(e.Start(":9000"))  // O el puerto que quieras
```

### ¿Qué pasa si el puerto 8080 ya está en uso?
Usa PowerShell:
```powershell
netstat -ano | findstr :8080
taskkill /PID [numero] /F
```

---

## Uso de la Aplicación

### ¿Qué es una "Máquina Virtual Base"?
Es la máquina virtual original en VirtualBox que usarás como plantilla. De ella se clonarán las máquinas de usuario.

### ¿Qué es un "Disco Multiconexión"?
Es un disco VirtualBox (formato VDI) que solo se puede crear una vez pero puede conectarse a múltiples máquinas virtuales. Una máquina de usuario usa este disco.

### ¿Obligatorio crear llaves de Root primero?
Sí. Las llaves de root deben existir antes de crear un disco multiconexión, porque el disco es parte de la infraestructura de la máquina base.

### ¿Puedo eliminar una máquina de usuario?
Sí, al hacer clic en "Eliminar VM", se eliminará completamente de VirtualBox. El disco multiconexión permanecerá y podrá usarse para otra máquina.

### ¿Qué se descarga cuando digo "Descargar llaves"?
Se descarga el archivo `id_rsa` que contiene la llave privada RSA de 1024 bits. Guárdalo en lugar seguro.

---

## Llaves SSH

### ¿Por qué RSA 1024 bits?
Es un requisito del proyecto. Aunque más lentas que las modernas (2048+), son suficientes para propósitos educativos.

### ¿Cómo conecto vía SSH?
Con la llave privada descargada:
```bash
ssh -i id_rsa -p 2200 root@localhost
```

Para usuarios:
```bash
ssh -i id_rsa -p 2200 username@localhost
```

### ¿Qué pasa si pierdo la llave privada?
Debes eliminar la máquina y crear una nueva. Las llaves no se pueden recuperar.

### ¿Las llaves están encriptadas?
Actualmente no, se guardan en texto plano. En producción, considera usar encriptación con `openssl`.

---

## VirtualBox

### ¿Las máquinas base deben estar apagadas?
Sí, la mayoría de operaciones requieren que estén apagadas para clonación.

### ¿Qué distribuciones soporta?
Cualquier distribución Debian-based instalada y configurada en VirtualBox. Se sugieren Ubuntu, Debian, Linux Mint.

### ¿Cómo configuro SSH antes de usar VM Manager?
1. Instala OpenSSH en cada máquina base:
   ```bash
   sudo apt-get update
   sudo apt-get install openssh-server
   ```
2. Asegúrate de que el servicio esté activo
3. Conoce la contraseña de root

### ¿Necesito acceso de red entre máquinas?
No es requisito para crear las máquinas. Pero para acceso SSH remoto, yes necesitarás conectividad.

### ¿Qué tipo de discos necesito?
Los discos deben ser de tipo normal (no dinámicos). Especifica esto al crear máquinas base.

---

## Problemas Comunes

### Error: "No such file or directory: templates"
Verifica que existe la carpeta `templates/` con `dashboard.html` en la raíz del proyecto.

```bash
dir templates
```

### Error: "Port address already in use"
Cambia el puerto en `main.go` o mata el proceso:
```powershell
netstat -ano | findstr :8080
taskkill /PID [numero] /F
```

### Error: "Machine not found"
La máquina no existe en VirtualBox. Verifica:
```bash
vboxmanage list vms
```

### Las llaves no se guardan
Verifica permisos en `C:\Users\[usuario]\.vm_manager\`:
```powershell
ls -Force C:\Users\$env:USERNAME\.vm_manager\
```

### No puedo conectar vía SSH
1. Verifica que SSH está activo en la VM
2. Verifica el puerto correcto
3. Revisa los logs del servidor

---

## Desarrollo

### ¿Cómo agrego una nueva funcionalidad?
Ver `DEVELOPMENT.md` para patrones de desarrollo.

### ¿Cómo hago tests?
```bash
go test ./...
go test -v ./...  # Con detalles
```

### ¿Cómo actualizo dependencias?
```bash
go get -u
go mod tidy
```

### ¿Dónde está el código de frontend?
En `templates/dashboard.html` está todo el HTML/CSS/JavaScript del dashboard.

---

## Performance y Escalabilidad

### ¿Cuántas máquinas puede gestionar?
Teóricamente ilimitadas, pero el dashboard se hace lento con muchas. Para producción, considerauses una base de datos.

### ¿Puede ejecutarse en paralelo?
Sí, usa `sync.RWMutex` para sincronización thread-safe.

### ¿Almacena datos en BD?
Actualmente, no. Usa almacenamiento en memoria. Al reiniciar, se pierde todo.

### ¿Cómo persisto los datos?
Por ahora, manual con VirtualBox. Para automatización, ver `DEVELOPMENT.md` sobre agregar SQLite.

---

## Seguridad

### ¿Las contraseñas se almacenan?
No. Solo se usan para la autenticación SSH inicial (no implementada en esta versión).

### ¿Las llaves están seguras?
Se guardan localmente con permisos 0600 (solo lectura para usuario). Protege la carpeta `.vm_manager`.

### ¿Hay autenticación en la web?
Actualmente no. Se asume acceso local. Para producción, agrega autenticación.

### ¿Se transmiten datos encriptados?
No (HTTP). Para producción, usa HTTPS con certificados SSL.

---

## Contacto y Soporte

### ¿Encontré un bug?
Verifica la sección Troubleshooting. Si persiste, revisa los logs.

### ¿Tengo una sugerencia?
Ver `DEVELOPMENT.md` para mejoras futuras.

### ¿Documentación oficial?
- [Echo Framework](https://echo.labstack.com/) 
- [Golang Crypto](https://pkg.go.dev/golang.org/x/crypto)
- [VirtualBox docs](https://www.virtualbox.org/wiki/Documentation)

---

## Notas Finales

- ✅ Proyecto educativo
- ✅ Basado en Golang + Echo
- ✅ Interfaz web simple
- ✅ Llaves SSH RSA 1024-bit
- ⏳ No hay persistencia en BD (en desarrollo)
- ⏳ No hay autenticación web (planificado)
- ⚠️ Usa en entorno local o seguro

---

¿Más preguntas? Revisa el README.md y DEVELOPMENT.md
