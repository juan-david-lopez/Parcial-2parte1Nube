package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"vm-manager/config"
	"vm-manager/handlers"
	"vm-manager/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer implementa echo.Renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Render renderiza un template
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Crear servicios
	vboxService := &services.VirtualBoxService{}
	sshService := services.NewSSHService(filepath.Join(os.Getenv("USERPROFILE"), ".vm_manager", "ssh_keys"))

	// Crear contexto de la aplicación
	appCtx := handlers.NewAppContext(vboxService, sshService)

	// Crear instancia de Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Cargar templates
	templates, err := template.ParseGlob(filepath.Join(cfg.Templates, "*.html"))
	if err != nil {
		fmt.Printf("Error cargando templates: %v\n", err)
	}

	e.Renderer = &TemplateRenderer{
		templates: templates,
	}

	// Rutas
	// Dashboard
	e.GET("/", appCtx.GetDashboard)

	// Base Machines
	e.POST("/machines/base/add", appCtx.AddBaseMachine)
	e.POST("/machines/base/:id/keys/root", appCtx.CreateRootKeys)
	e.GET("/machines/base/:id/keys/root/download", appCtx.DownloadRootKeys)
	e.POST("/machines/base/:id/disks/create", appCtx.CreateMediaDisk)

	// User Machines
	e.POST("/machines/user/add/:diskid", appCtx.AddUserMachine)
	e.POST("/machines/user/:vmid/keys", appCtx.CreateUserSSHKeys)
	e.GET("/machines/user/:vmid/keys/:username/download", appCtx.DownloadUserKeys)
	e.DELETE("/machines/user/:vmid", appCtx.DeleteUserMachine)

	// Media Disks
	e.POST("/disks/:diskid/connect/:vmid", appCtx.ConnectDisk)
	e.POST("/disks/:diskid/disconnect", appCtx.DisconnectDisk)
	e.DELETE("/disks/:diskid", appCtx.DeleteDisk)

	// Iniciar servidor
	fmt.Printf("Servidor iniciado en http://localhost:8080\n")
	e.Logger.Fatal(e.Start(":8080"))
}
