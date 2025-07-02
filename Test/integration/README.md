# Tests de Integración

Este directorio contiene los tests de integración para la API de ia-boilerplate-go, implementados usando Cucumber/Gherkin con el framework Godog.

## Filosofías del Framework de Testing

### 1. **Gestión Autónoma de Recursos**
- Todos los recursos creados durante los tests son automáticamente rastreados
- Limpieza automática al final de cada escenario
- Prevención de contaminación entre tests

### 2. **Autenticación Automática**
- Login automático al inicio de cada escenario
- Tokens de acceso manejados globalmente
- Headers de autorización agregados automáticamente

### 3. **Variables Dinámicas**
- Generación de valores únicos para evitar conflictos
- Sustitución de variables en URLs y payloads
- Persistencia de valores entre pasos del escenario

### 4. **Validación Robusta**
- Verificación de códigos de estado HTTP
- Validación de estructura JSON de respuestas
- Manejo de errores y casos edge

## Estructura de Archivos

```
Test/integration/
├── main_test.go              # Configuración principal de tests
├── steps.go                  # Implementación de pasos Gherkin
├── README.md                 # Este archivo
└── features/                 # Archivos de features Gherkin
    ├── auth.feature          # Tests de autenticación
    ├── users.feature         # Tests de usuarios, roles y dispositivos
    ├── medicine.feature      # Tests de medicamentos
    ├── icd-cie.feature       # Tests de códigos ICD-CIE
    ├── device-info.feature   # Tests de información de dispositivos
    └── error-handling.feature # Tests de manejo de errores
```

## Archivos de Features

### 1. **auth.feature**
Tests de autenticación y autorización:
- Login con credenciales válidas/inválidas
- Refresh de tokens
- Acceso a endpoints protegidos sin autenticación

### 2. **users.feature**
Tests completos de gestión de usuarios:
- CRUD de roles de usuario
- CRUD de usuarios
- CRUD de dispositivos asociados a usuarios
- Búsquedas paginadas y por propiedades

### 3. **medicine.feature**
Tests de gestión de medicamentos:
- CRUD de medicamentos
- Validación de códigos EAN únicos
- Búsquedas avanzadas
- Manejo de campos requeridos

### 4. **icd-cie.feature**
Tests de códigos ICD-CIE:
- CRUD de registros ICD-CIE
- Búsquedas con filtros múltiples
- Validación de propiedades de búsqueda
- Paginación y casos edge

### 5. **device-info.feature**
Tests de información de dispositivos:
- Endpoint de información de dispositivo
- Health check autenticado
- Verificación de middleware de dispositivos

### 6. **error-handling.feature**
Tests de manejo de errores:
- Casos de autenticación fallida
- IDs inválidos
- Campos requeridos faltantes
- Payloads JSON malformados
- Casos edge de paginación

## Ejecución de Tests

### Opción 1: Script Automatizado (Recomendado)

```bash
# Ejecutar todos los tests
./scripts/run-all-integration-tests.bash

# Ejecutar tests específicos
./scripts/run-all-integration-tests.bash -f auth.feature
./scripts/run-all-integration-tests.bash -f users.feature

# Ejecutar con Docker
./scripts/run-all-integration-tests.bash -d -v

# Ejecutar con tags específicos
./scripts/run-all-integration-tests.bash -t @smoke

# Modo verbose
./scripts/run-all-integration-tests.bash -v
```

### Opción 2: Comando Directo

```bash
# Ejecutar todos los tests
go test -tags=integration ./Test/integration/...

# Ejecutar con verbose
go test -v -tags=integration ./Test/integration/...

# Ejecutar feature específico
INTEGRATION_FEATURE_FILE=auth.feature go test -tags=integration ./Test/integration/...

# Ejecutar con tags específicos
INTEGRATION_SCENARIO_TAGS=@smoke go test -tags=integration ./Test/integration/...
```

### Opción 3: Docker Compose

```bash
# Ejecutar tests con Docker
docker-compose run --rm app go test -tags=integration ./Test/integration/...

# Ejecutar con verbose
docker-compose run --rm app go test -v -tags=integration ./Test/integration/...
```

## Variables de Entorno

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `INTEGRATION_FEATURE_FILE` | Ejecutar solo un archivo de feature | `auth.feature` |
| `INTEGRATION_SCENARIO_TAGS` | Ejecutar solo escenarios con tags específicos | `@smoke` |
| `INTEGRATION_TEST_MODE` | Modo de testing activado | `true` |

## Estructura de un Escenario

```gherkin
Scenario: TC01 - Create a new user successfully
  Given I generate a unique alias as "newUserUsername"
  And I generate a unique alias as "newUserEmail"
  When I send a POST request to "/api/users" with body:
    """
    {
      "username": "${newUserUsername}",
      "email": "${newUserEmail}@test.com",
      "password": "securePassword123",
      "roleId": 1,
      "enabled": true
    }
    """
  Then the response code should be 201
  And the JSON response should contain key "id"
  And I save the JSON response key "id" as "userID"
```

## Pasos Disponibles

### Pasos Given (Configuración)
- `I generate a unique alias as "varName"`
- `I generate a unique EAN code as "varName"`
- `I clear the authentication token`
- `I am authenticated as a user`

### Pasos When (Acciones)
- `I send a GET request to "path"`
- `I send a POST request to "path" with body:`
- `I send a PUT request to "path" with body:`
- `I send a DELETE request to "path"`

### Pasos Then (Validaciones)
- `the response code should be 200`
- `the JSON response should contain key "keyName"`
- `the JSON response should contain "field": "value"`
- `the JSON response should contain error "error": "message"`
- `I save the JSON response key "key" as "varName"`

## Gestión de Recursos

### Creación Automática
Los recursos creados durante los tests son automáticamente rastreados:

```go
// En steps.go
func trackResource(path string) {
    // Rastrea recursos para limpieza posterior
}
```

### Limpieza Automática
Al final de cada escenario, todos los recursos creados son eliminados:

```go
// En steps.go
func InitializeScenario(ctx *godog.ScenarioContext) {
    // Setup y teardown automático
}
```

## Debugging

### Modo Verbose
```bash
go test -v -tags=integration ./Test/integration/...
```

### Logs Detallados
Los tests incluyen logs detallados que muestran:
- URLs de requests
- Headers enviados
- Códigos de respuesta
- Cuerpo de respuestas
- Variables generadas

### Variables de Debug
```bash
# Habilitar logs de debug
export DEBUG=true
go test -tags=integration ./Test/integration/...
```

## Mejores Prácticas

### 1. **Nombres Únicos**
Siempre usa generadores de valores únicos:
```gherkin
Given I generate a unique alias as "testUser"
```

### 2. **Validación Completa**
Valida tanto el código de respuesta como el contenido:
```gherkin
Then the response code should be 201
And the JSON response should contain key "id"
And the JSON response should contain "username": "${testUser}"
```

### 3. **Manejo de Errores**
Incluye tests para casos de error:
```gherkin
Scenario: Attempt to create user with missing fields
  When I send a POST request to "/api/users" with body:
    """
    {
      "firstName": "John"
    }
    """
  Then the response code should be 400
  And the JSON response should contain key "error"
```

### 4. **Limpieza de Recursos**
Los recursos se limpian automáticamente, pero puedes limpiar manualmente:
```gherkin
When I send a DELETE request to "/api/users/${userID}"
Then the response code should be 200
```

## Troubleshooting

### Problemas Comunes

1. **Error de conexión a base de datos**
   - Verifica que Docker Compose esté corriendo
   - Revisa las variables de entorno de conexión

2. **Tests fallando por recursos existentes**
   - Ejecuta con la opción `-c` para limpiar antes
   - Verifica que no haya tests corriendo en paralelo

3. **Errores de autenticación**
   - Verifica que las credenciales de test sean correctas
   - Revisa que el servidor esté corriendo

4. **Timeouts en tests**
   - Aumenta el timeout en la configuración
   - Verifica la conectividad de red

### Logs de Debug
```bash
# Habilitar logs detallados
export GODOG_DEBUG=true
go test -v -tags=integration ./Test/integration/...
```

## Contribución

### Agregar Nuevos Tests

1. **Crear archivo de feature**:
   ```bash
   touch Test/integration/features/nueva-funcionalidad.feature
   ```

2. **Implementar pasos** (si es necesario):
   - Agregar funciones en `steps.go`
   - Registrar en `InitializeScenario`

3. **Ejecutar tests**:
   ```bash
   ./scripts/run-all-integration-tests.bash -f nueva-funcionalidad.feature
   ```

### Convenciones de Nomenclatura

- **Archivos de feature**: `kebab-case.feature`
- **Escenarios**: `TC01 - Descripción del test`
- **Variables**: `camelCase` o `snake_case`
- **Tags**: `@smoke`, `@regression`, `@critical`

## Integración Continua

### GitHub Actions
```yaml
- name: Run Integration Tests
  run: |
    docker-compose up -d
    ./scripts/run-all-integration-tests.bash -d -v
```

### Jenkins Pipeline
```groovy
stage('Integration Tests') {
    steps {
        sh './scripts/run-all-integration-tests.bash -d -v'
    }
}
```

## Recursos Adicionales

- [Documentación de Godog](https://github.com/cucumber/godog)
- [Sintaxis Gherkin](https://cucumber.io/docs/gherkin/)
- [Testing en Go](https://golang.org/pkg/testing/) 