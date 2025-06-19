#!/bin/bash

# Script para generar cobertura de código excluyendo main.go

echo "Generando cobertura de código..."

# Ejecutar tests con cobertura
go test -coverprofile=coverage.out ./...

# Filtrar main.go del archivo de cobertura
grep -v "main.go" coverage.out > coverage_filtered.out

# Generar reporte HTML
go tool cover -html=coverage_filtered.out -o coverage.html

# Mostrar porcentaje de cobertura
echo "Cobertura de código (excluyendo main.go):"
go tool cover -func=coverage_filtered.out | tail -1

echo "Reporte HTML generado en coverage.html" 