# Endpoints de Búsqueda - Medicine y User

Este documento describe los nuevos endpoints de búsqueda implementados para las entidades `Medicine` y `User`.

## Endpoints Disponibles

### Medicine

#### 1. Búsqueda Paginada
```
GET /medicine/search
```

**Parámetros de consulta:**
- `page` (opcional): Número de página (default: 1)
- `pageSize` (opcional): Tamaño de página (default: 10)
- `sortBy` (opcional): Campo(s) para ordenar (múltiples valores permitidos)
- `sortDirection` (opcional): Dirección de ordenamiento (`asc` o `desc`, default: `asc`)

**Filtros LIKE (búsqueda parcial):**
- `name_like`: Búsqueda parcial en el nombre
- `description_like`: Búsqueda parcial en la descripción
- `eanCode_like`: Búsqueda parcial en el código EAN
- `laboratory_like`: Búsqueda parcial en el laboratorio

**Filtros de coincidencia exacta:**
- `name_match`: Coincidencia exacta en el nombre (múltiples valores)
- `description_match`: Coincidencia exacta en la descripción (múltiples valores)
- `eanCode_match`: Coincidencia exacta en el código EAN (múltiples valores)
- `laboratory_match`: Coincidencia exacta en el laboratorio (múltiples valores)

**Filtros de rango de fechas:**
- `createdAt_start`: Fecha de inicio para createdAt (formato RFC3339)
- `createdAt_end`: Fecha de fin para createdAt (formato RFC3339)
- `updatedAt_start`: Fecha de inicio para updatedAt (formato RFC3339)
- `updatedAt_end`: Fecha de fin para updatedAt (formato RFC3339)

**Ejemplo de uso:**
```
GET /medicine/search?page=1&pageSize=10&name_like=aspirin&sortBy=name&sortDirection=asc
```

**Respuesta:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Aspirin",
      "description": "Pain reliever",
      "eanCode": "1234567890123",
      "laboratory": "Bayer",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "pageSize": 10,
  "totalPages": 1,
  "filters": {
    "likeFilters": {
      "name": ["aspirin"]
    },
    "matches": {},
    "dateRanges": [],
    "sortBy": ["name"],
    "sortDirection": "asc",
    "page": 1,
    "pageSize": 10
  }
}
```

#### 2. Búsqueda por Propiedad
```
GET /medicine/search/property
```

**Parámetros de consulta:**
- `property` (requerido): Propiedad a buscar (`name`, `description`, `eanCode`, `laboratory`)
- `searchText` (requerido): Texto a buscar

**Ejemplo de uso:**
```
GET /medicine/search/property?property=name&searchText=asp
```

**Respuesta:**
```json
["Aspirin", "Aspartame"]
```

### User

#### 1. Búsqueda Paginada
```
GET /user/search
```

**Parámetros de consulta:**
- `page` (opcional): Número de página (default: 1)
- `pageSize` (opcional): Tamaño de página (default: 10)
- `sortBy` (opcional): Campo(s) para ordenar (múltiples valores permitidos)
- `sortDirection` (opcional): Dirección de ordenamiento (`asc` o `desc`, default: `asc`)

**Filtros LIKE (búsqueda parcial):**
- `userName_like`: Búsqueda parcial en el nombre de usuario
- `email_like`: Búsqueda parcial en el email
- `firstName_like`: Búsqueda parcial en el nombre
- `lastName_like`: Búsqueda parcial en el apellido
- `status_like`: Búsqueda parcial en el estado

**Filtros de coincidencia exacta:**
- `userName_match`: Coincidencia exacta en el nombre de usuario (múltiples valores)
- `email_match`: Coincidencia exacta en el email (múltiples valores)
- `firstName_match`: Coincidencia exacta en el nombre (múltiples valores)
- `lastName_match`: Coincidencia exacta en el apellido (múltiples valores)
- `status_match`: Coincidencia exacta en el estado (múltiples valores)

**Filtros de rango de fechas:**
- `createdAt_start`: Fecha de inicio para createdAt (formato RFC3339)
- `createdAt_end`: Fecha de fin para createdAt (formato RFC3339)
- `updatedAt_start`: Fecha de inicio para updatedAt (formato RFC3339)
- `updatedAt_end`: Fecha de fin para updatedAt (formato RFC3339)

**Ejemplo de uso:**
```
GET /user/search?page=1&pageSize=10&email_like=john&sortBy=createdAt&sortDirection=desc
```

**Respuesta:**
```json
{
  "data": [
    {
      "id": 1,
      "user": "john_doe",
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "status": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "pageSize": 10,
  "totalPages": 1,
  "filters": {
    "likeFilters": {
      "email": ["john"]
    },
    "matches": {},
    "dateRanges": [],
    "sortBy": ["createdAt"],
    "sortDirection": "desc",
    "page": 1,
    "pageSize": 10
  }
}
```

#### 2. Búsqueda por Propiedad
```
GET /user/search/property
```

**Parámetros de consulta:**
- `property` (requerido): Propiedad a buscar (`userName`, `email`, `firstName`, `lastName`, `status`)
- `searchText` (requerido): Texto a buscar

**Ejemplo de uso:**
```
GET /user/search/property?property=email&searchText=john
```

**Respuesta:**
```json
["john@example.com", "johnny@example.com"]
```

## Características de los Endpoints

### Búsqueda Paginada
- **Paginación**: Control de página y tamaño de página
- **Filtros LIKE**: Búsqueda parcial con `ILIKE` (case-insensitive)
- **Filtros de coincidencia exacta**: Búsqueda exacta con `IN`
- **Filtros de rango de fechas**: Búsqueda por rangos de fechas
- **Ordenamiento**: Múltiples campos con dirección configurable
- **Respuesta estructurada**: Incluye metadatos de paginación y filtros aplicados

### Búsqueda por Propiedad
- **Búsqueda de valores únicos**: Retorna valores distintos de una propiedad
- **Límite de resultados**: Máximo 20 resultados
- **Validación de propiedades**: Solo propiedades válidas permitidas
- **Búsqueda parcial**: Usa `ILIKE` para búsqueda case-insensitive

## Autenticación

Todos los endpoints requieren autenticación JWT. Incluye el token en el header:
```
Authorization: Bearer <your-jwt-token>
```

## Manejo de Errores

Los endpoints manejan los siguientes tipos de errores:
- **400 Bad Request**: Parámetros inválidos o faltantes
- **401 Unauthorized**: Token JWT inválido o faltante
- **500 Internal Server Error**: Errores internos del servidor

## Notas de Implementación

- Los filtros de fecha deben estar en formato RFC3339
- Los filtros LIKE son case-insensitive
- Los filtros de coincidencia exacta pueden recibir múltiples valores
- El ordenamiento soporta múltiples campos
- La paginación comienza desde la página 1 