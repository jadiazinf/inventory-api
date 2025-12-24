# Sistema de Gestión de Inventario - Bazaar Araira

API REST desarrollada en Go para la gestión integral de inventario, ventas, clientes y cuentas por cobrar del Bazaar Araira, con funcionalidades especializadas para útiles escolares.

## 📋 Tabla de Contenidos

- [Requisitos Previos](#requisitos-previos)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Arquitectura](#arquitectura)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Ejecución](#ejecución)
- [Testing](#testing)
- [Endpoints API](#endpoints-api)
- [Autenticación](#autenticación)
- [Base de Datos](#base-de-datos)

## 🔧 Requisitos Previos

- **Go**: versión 1.23 o superior
- **PostgreSQL**: versión 14 o superior
- **Firebase**: proyecto configurado para autenticación
- **Git**: para clonar el repositorio

## 📁 Estructura del Proyecto

```
backend/
├── api/                              # Aplicación principal
│   ├── cmd/
│   │   └── main.go                   # Punto de entrada de la aplicación
│   │
│   ├── internal/
│   │   ├── core/                     # Núcleo de la lógica de negocio
│   │   │   ├── domain/               # Entidades del dominio
│   │   │   │   ├── base.go           # Modelos base (timestamps, soft delete)
│   │   │   │   ├── customer.go       # Cliente, CustomerChild
│   │   │   │   ├── employee.go       # Empleados, roles, permisos
│   │   │   │   ├── enums.go          # Enumeraciones del sistema
│   │   │   │   ├── inventory.go      # Inventario, movimientos, almacenes
│   │   │   │   ├── location.go       # Ubicaciones geográficas
│   │   │   │   ├── product.go        # Productos, categorías, precios
│   │   │   │   ├── reservation.go    # Reservas y pre-órdenes
│   │   │   │   ├── sale.go           # Ventas, detalles, cuentas por cobrar
│   │   │   │   ├── school.go         # Listas escolares
│   │   │   │   ├── store.go          # Tiendas
│   │   │   │   └── user.go           # Usuarios del sistema
│   │   │   │
│   │   │   ├── ports/                # Interfaces (Hexagonal Architecture)
│   │   │   │   ├── repositories/     # Interfaces de repositorios
│   │   │   │   │   ├── customer_repository.go
│   │   │   │   │   ├── inventory_repository.go
│   │   │   │   │   ├── product_repository.go
│   │   │   │   │   ├── reservation_repository.go
│   │   │   │   │   ├── sale_repository.go
│   │   │   │   │   └── user_repository.go
│   │   │   │   │
│   │   │   │   └── services/         # Interfaces de servicios
│   │   │   │       ├── accounts_receivable_service.go
│   │   │   │       ├── inventory_service.go
│   │   │   │       ├── notification_service.go
│   │   │   │       ├── product_service.go
│   │   │   │       ├── reservation_service.go
│   │   │   │       └── sale_service.go
│   │   │   │
│   │   │   └── services/             # Implementación de servicios
│   │   │       ├── accounts_receivable_service.go
│   │   │       ├── inventory_service.go
│   │   │       ├── notification_service.go
│   │   │       ├── product_service.go
│   │   │       ├── reservation_service.go
│   │   │       ├── sale_service.go
│   │   │       └── utils.go
│   │   │
│   │   ├── adapters/                 # Adaptadores externos
│   │   │   ├── http/
│   │   │   │   ├── dto/              # Data Transfer Objects
│   │   │   │   │   ├── customer_dto.go
│   │   │   │   │   ├── inventory_dto.go
│   │   │   │   │   ├── product_dto.go
│   │   │   │   │   ├── reservation_dto.go
│   │   │   │   │   ├── response.go
│   │   │   │   │   └── sale_dto.go
│   │   │   │   │
│   │   │   │   ├── handlers/         # Controladores HTTP
│   │   │   │   │   ├── customer_handler.go
│   │   │   │   │   ├── inventory_handler.go
│   │   │   │   │   ├── product_handler.go
│   │   │   │   │   ├── reservation_handler.go
│   │   │   │   │   └── sale_handler.go
│   │   │   │   │
│   │   │   │   └── middleware/       # Middleware HTTP
│   │   │   │       └── auth_middleware.go
│   │   │   │
│   │   │   └── repository/
│   │   │       └── postgres/         # Implementación PostgreSQL
│   │   │           ├── accounts_receivable_repository.go
│   │   │           ├── customer_child_repository.go
│   │   │           ├── customer_repository.go
│   │   │           ├── inventory_repository.go
│   │   │           ├── product_repository.go
│   │   │           ├── reservation_repository.go
│   │   │           ├── sale_repository.go
│   │   │           ├── user_repository.go
│   │   │           ├── *_test.go     # Tests unitarios
│   │   │           └── utils.go
│   │   │
│   │   ├── api/                      # Configuración API
│   │   │   ├── server.go             # Servidor Fiber
│   │   │   └── routes.go             # Definición de rutas
│   │   │
│   │   ├── config/                   # Configuración
│   │   │   └── config.go
│   │   │
│   │   ├── platform/                 # Infraestructura
│   │   │   ├── database/             # Conexión a base de datos
│   │   │   │   └── postgres.go
│   │   │   ├── firebase/             # Cliente Firebase
│   │   │   │   └── firebase.go
│   │   │   └── i18n/                 # Internacionalización
│   │   │       ├── i18n.go
│   │   │       └── locales/
│   │   │           ├── en.json
│   │   │           └── es.json
│   │   │
│   │   ├── common/                   # Utilidades comunes
│   │   │   └── errors/               # Manejo centralizado de errores
│   │   │       └── errors.go
│   │   │
│   │   └── health/                   # Health checks
│   │       └── handler.go
│   │
│   ├── go.mod                        # Dependencias Go
│   └── go.sum                        # Checksums de dependencias
│
├── database/                         # Scripts de base de datos
│   └── schema.sql                    # Schema completo con datos semilla
│
├── SISTEMA.md                        # Documentación del sistema
└── ESTADO_PROYECTO.md                # Estado actual y análisis DOFA
```

## 🏗️ Arquitectura

### Clean Architecture (Arquitectura Hexagonal)

El proyecto sigue los principios de Clean Architecture para mantener el código desacoplado, testeable y mantenible:

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer (Fiber)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Handlers   │  │  Middleware  │  │     DTOs     │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                   Service Layer                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Business Logic (Use Cases)                      │  │
│  │  - ProductService                                │  │
│  │  - SaleService                                   │  │
│  │  - ReservationService                            │  │
│  │  - InventoryService                              │  │
│  │  - AccountsReceivableService                     │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                  Domain Layer                           │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Entities & Business Rules                       │  │
│  │  - Product, Customer, Sale, Reservation, etc.    │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Ports (Interfaces)                              │  │
│  │  - Repository Interfaces                         │  │
│  │  - Service Interfaces                            │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│              Infrastructure Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  PostgreSQL  │  │   Firebase   │  │     i18n     │  │
│  │ Repositories │  │     Auth     │  │              │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Flujo de Datos

1. **Request** → Handler (HTTP)
2. **Handler** → DTO validation
3. **Handler** → Service (Business Logic)
4. **Service** → Repository (Data Access)
5. **Repository** → Database
6. **Database** → Repository → Service → Handler → Response

### Principios Aplicados

- **Separación de responsabilidades**: Cada capa tiene una responsabilidad específica
- **Inversión de dependencias**: Las capas internas no conocen las externas
- **Inyección de dependencias**: Manual via constructores
- **Repository Pattern**: Abstracción del acceso a datos
- **DTO Pattern**: Separación entre modelos de dominio y API

## 🚀 Instalación

### 1. Clonar el Repositorio

```bash
git clone <repository-url>
cd backend/api
```

### 2. Instalar Dependencias

```bash
go mod download
```

### 3. Configurar PostgreSQL

```bash
# Crear base de datos
createdb inventory

# Ejecutar schema
psql -U postgres -d inventory -f ../database/schema.sql
```

### 4. Configurar Firebase

1. Crear proyecto en [Firebase Console](https://console.firebase.google.com/)
2. Habilitar Authentication
3. Descargar credenciales (Service Account)
4. Guardar como `firebase-credentials.json` en la raíz del proyecto

## ⚙️ Configuración

### Variables de Entorno

Crear archivo `.env` en `/backend/api/`:

```env
# Servidor
APP_PORT=3000

# Base de Datos PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=tu_usuario
DB_PASSWORD=tu_contraseña
DB_NAME=inventory

# Firebase Authentication
FIREBASE_CREDENTIALS=firebase-credentials.json
```

### Estructura de Configuración

La configuración se carga desde el archivo `internal/config/config.go` que lee las variables de entorno usando `godotenv`.

## 🏃 Ejecución

### Modo Desarrollo

```bash
# Ejecutar directamente
go run cmd/main.go
```

### Modo Producción

```bash
# Compilar
go build -o server cmd/main.go

# Ejecutar binario
./server
```

### Verificar que el Servidor está Corriendo

```bash
# Health check
curl http://localhost:3000/health

# Response esperado:
# {"message":"Service is healthy","status":"ok"}
```

## 🧪 Testing

### Ejecutar Todos los Tests

```bash
go test ./...
```

### Tests por Paquete

```bash
# Repositorios
go test ./internal/adapters/repository/postgres/... -v

# Servicios (pendiente)
go test ./internal/core/services/... -v

# Handlers (pendiente)
go test ./internal/adapters/http/handlers/... -v
```

### Tests con Cobertura

```bash
go test -cover ./...
```

### Nota sobre Tests de Repositorio

Los tests de repositorio requieren PostgreSQL debido al uso de:

- Tipos UUID nativos
- Enums personalizados
- Funciones PostgreSQL (`uuid_generate_v4()`)

SQLite no es compatible para testing in-memory.

## 📡 Endpoints API

### Health Check

```http
GET /health
GET /api/v1/health
GET /api/v1/greet/:name
```

### Productos

```http
GET    /api/v1/products              # Listar (público)
GET    /api/v1/products/search       # Buscar (público)
GET    /api/v1/products/:id          # Ver detalle (público)
GET    /api/v1/products/sku/:sku     # Buscar por SKU (público)
POST   /api/v1/products              # Crear (requiere auth)
PUT    /api/v1/products/:id          # Actualizar (requiere auth)
PUT    /api/v1/products/:id/price    # Actualizar precio (requiere auth)
DELETE /api/v1/products/:id          # Eliminar (requiere auth)
```

### Clientes

```http
GET    /api/v1/customers                      # Listar (requiere auth)
GET    /api/v1/customers/:id                  # Ver detalle (requiere auth)
GET    /api/v1/customers/:id/with-children    # Ver con hijos (requiere auth)
GET    /api/v1/customers/tax-id/:taxId        # Buscar por CI/RIF (requiere auth)
POST   /api/v1/customers                      # Crear (requiere auth)
PUT    /api/v1/customers/:id                  # Actualizar (requiere auth)
DELETE /api/v1/customers/:id                  # Eliminar (requiere auth)

# Gestión de Hijos
POST   /api/v1/customers/:id/children         # Agregar hijo (requiere auth)
GET    /api/v1/customers/:id/children         # Listar hijos (requiere auth)
PUT    /api/v1/customers/:id/children/:childId # Actualizar hijo (requiere auth)
DELETE /api/v1/customers/:id/children/:childId # Eliminar hijo (requiere auth)

# Puntos de Lealtad
PUT    /api/v1/customers/:id/loyalty-points   # Actualizar puntos (requiere auth)
```

### Ventas

```http
GET    /api/v1/sales                    # Listar (requiere auth)
GET    /api/v1/sales/daily              # Ventas diarias (requiere auth)
GET    /api/v1/sales/:id                # Ver detalle (requiere auth)
GET    /api/v1/sales/invoice/:invoice   # Buscar por factura (requiere auth)
POST   /api/v1/sales                    # Crear venta (requiere auth)
POST   /api/v1/sales/credit             # Crear venta a crédito (requiere auth)
POST   /api/v1/sales/:id/cancel         # Cancelar venta (requiere auth)
```

### Cuentas por Cobrar

```http
GET    /api/v1/accounts-receivable/:id           # Ver cuenta (requiere auth)
POST   /api/v1/accounts-receivable/:id/payments  # Registrar pago (requiere auth)
```

### Reservas

```http
GET    /api/v1/reservations                # Listar (requiere auth)
GET    /api/v1/reservations/:id            # Ver detalle (requiere auth)
GET    /api/v1/reservations/number/:number # Buscar por número (requiere auth)
POST   /api/v1/reservations                # Crear (requiere auth)
POST   /api/v1/reservations/:id/confirm    # Confirmar (requiere auth)
POST   /api/v1/reservations/:id/fulfill    # Cumplir (requiere auth)
POST   /api/v1/reservations/:id/cancel     # Cancelar (requiere auth)
```

### Inventario

```http
GET    /api/v1/inventory/product/:productId/warehouse/:warehouseId  # Stock específico (requiere auth)
GET    /api/v1/inventory/warehouse/:warehouseId                     # Todo el inventario (requiere auth)
GET    /api/v1/inventory/product/:productId                         # Por producto (requiere auth)
GET    /api/v1/inventory/check-availability                         # Verificar disponibilidad (requiere auth)

# Movimientos
POST   /api/v1/inventory/movements/inbound     # Entrada (requiere auth)
POST   /api/v1/inventory/movements/outbound    # Salida (requiere auth)
POST   /api/v1/inventory/movements/adjustment  # Ajuste (requiere auth)
GET    /api/v1/inventory/movements/product/:productId    # Por producto (requiere auth)
GET    /api/v1/inventory/movements/warehouse/:warehouseId # Por almacén (requiere auth)
```

## 🔐 Autenticación

### Firebase Authentication

Los endpoints marcados con **(requiere auth)** necesitan un token JWT de Firebase en el header:

```http
Authorization: Bearer <firebase-id-token>
```

### Obtener Token (Frontend)

```javascript
// Firebase Web SDK
const user = firebase.auth().currentUser;
const token = await user.getIdToken();

// Hacer request
fetch('http://localhost:3000/api/v1/products', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

### Middleware de Autenticación

El middleware en `internal/adapters/http/middleware/auth_middleware.go`:

1. Valida el token con Firebase
2. Extrae el UID del usuario
3. Busca el usuario en la base de datos
4. Agrega el usuario al contexto de la petición

## 💾 Base de Datos

### Schema

El schema completo se encuentra en `/backend/database/schema.sql` e incluye:

- **49 tablas** con relaciones completas
- **15+ enums personalizados** (tipos de venta, estados, roles, etc.)
- **Triggers** para actualización automática de timestamps
- **Índices optimizados** para búsquedas frecuentes
- **Full-text search** configurado para español
- **Funciones PostgreSQL** para lógica de negocio
- **Datos semilla** para desarrollo

### Migraciones

Actualmente el proyecto usa un único archivo SQL. Para migraciones incrementales se recomienda usar:

- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [pressly/goose](https://github.com/pressly/goose)

### Soft Deletes

Todas las tablas principales implementan soft delete mediante la columna `deleted_at`. Los registros eliminados:

- No aparecen en queries normales (GORM los filtra automáticamente)
- Pueden recuperarse si es necesario
- Mantienen integridad referencial

## 📦 Dependencias Principales

```go
require (
    github.com/gofiber/fiber/v2 v2.52.10          // Framework web
    gorm.io/gorm v1.25.12                         // ORM
    gorm.io/driver/postgres v1.5.9                // Driver PostgreSQL
    github.com/google/uuid v1.6.0                 // UUIDs
    github.com/goccy/go-json v0.10.2              // JSON alta performance
    firebase.google.com/go/v4 v4.13.0             // Firebase Admin SDK
    github.com/joho/godotenv v1.5.1               // Variables de entorno
    github.com/go-playground/validator/v10 v10.22.1 // Validación
    github.com/nicksnyder/go-i18n/v2 v2.4.0       // Internacionalización
)
```

## 🛠️ Stack Tecnológico

- **Lenguaje**: Go 1.23
- **Framework Web**: Fiber v2
- **ORM**: GORM
- **Base de Datos**: PostgreSQL 14+
- **Autenticación**: Firebase Authentication
- **Testing**: go test + testify
- **Documentación**: Swagger (instalado, pendiente de configurar)

## 📚 Recursos Adicionales

- **Documentación del Sistema**: Ver [SISTEMA.md](SISTEMA.md)
- **Estado del Proyecto**: Ver [ESTADO_PROYECTO.md](ESTADO_PROYECTO.md)
- **Fiber Documentation**: https://docs.gofiber.io/
- **GORM Documentation**: https://gorm.io/docs/
- **Firebase Admin SDK**: https://firebase.google.com/docs/admin/setup

## 🤝 Contribución

Para contribuir al proyecto:

1. Seguir la arquitectura hexagonal establecida
2. Mantener cobertura de tests alta
3. Documentar cambios significativos
4. Usar conventional commits
5. Respetar las convenciones de código Go

## 📄 Licencia

Propietario - Bazaar Araira

---

**Desarrollado por**: Jesus Diaz
**Asistencia de IA**: Claude Sonnet 4.5
**Versión**: 1.0.0
**Última Actualización**: 2025-12-23
