#!/usr/bin/env powershell
# QUICK START - Inicia rápidamente

Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║          VM Manager - Gestor de Máquinas Virtuales            ║" -ForegroundColor Cyan
Write-Host "║                   GUÍA DE INICIO RÁPIDO                        ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Verificar requisitos
Write-Host "Verificando requisitos..." -ForegroundColor Yellow
$checks = @{
    "Go" = "go version"
    "VirtualBox CLI" = "vboxmanage --version"
    "SSH" = "ssh -V"
}

foreach ($check in $checks.GetEnumerator()) {
    Write-Host "Verificando $($check.Name)..." -NoNewline
    try {
        $result = Invoke-Expression $check.Value 2>&1 | Select-Object -First 1
        Write-Host "OK" -ForegroundColor Green
    } catch {
        Write-Host "NO ENCONTRADO" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "Descargar dependencias..." -ForegroundColor Yellow
go mod download

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error descargando dependencias" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Compilando aplicación..." -ForegroundColor Yellow
go build -o vm-manager.exe

if ($LASTEXITCODE -ne 0) {
    Write-Host " Error compilando" -ForegroundColor Red
    exit 1
}

Write-Host "Compilación exitosa" -ForegroundColor Green
Write-Host ""

Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║                    🚀 LISTO PARA EJECUTAR                      ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""

Write-Host "Comando para iniciar:" -ForegroundColor Cyan
Write-Host "  .\vm-manager.exe" -ForegroundColor White
Write-Host ""

Write-Host "URL del dashboard:" -ForegroundColor Cyan
Write-Host "  http://localhost:8080" -ForegroundColor White
Write-Host ""

Write-Host "Documentación:" -ForegroundColor Cyan
Write-Host "README.md - Guía de usuario" -ForegroundColor Gray
Write-Host "DEVELOPMENT.md - Arquitectura" -ForegroundColor Gray
Write-Host "BUILD.md - Deployment" -ForegroundColor Gray
Write-Host "FAQ.md - Preguntas frecuentes" -ForegroundColor Gray
Write-Host "FILES_INDEX.md - Índice de archivos" -ForegroundColor Gray
Write-Host ""

Write-Host "¿Ejecutar ahora?" -ForegroundColor Yellow
$response = Read-Host "Escribe 'si' para iniciar (o cualquier otra tecla para salir)"

if ($response -eq "si") {
    Write-Host ""
    Write-Host "Iniciando VM Manager..." -ForegroundColor Green
    Write-Host ""
    .\vm-manager.exe
} else {
    Write-Host ""
    Write-Host "Ejecución guardada para después" -ForegroundColor Gray
    Write-Host "   Inicia con:  .\vm-manager.exe" -ForegroundColor Gray
}
