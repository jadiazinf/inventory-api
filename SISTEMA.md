# Sistema de Gestión de Inventario - Bazaar Araira

## 📖 Índice

- [Razón de Ser del Proyecto](#razón-de-ser-del-proyecto)
- [Descripción del Sistema](#descripción-del-sistema)
- [Módulos del Core](#módulos-del-core)
- [Reglas de Negocio](#reglas-de-negocio)
- [Flujos de Trabajo](#flujos-de-trabajo)
- [Modelos del Dominio](#modelos-del-dominio)

---

## 🎯 Razón de Ser del Proyecto

### Contexto del Negocio

**Bazaar Araira** es un comercio venezolano especializado en la venta de útiles escolares, papelería y productos generales. El negocio opera con las siguientes características:

- **Temporadas altas**: Inicio del año escolar (agosto-octubre)
- **Cliente objetivo**: Familias con hijos en edad escolar, instituciones educativas
- **Modalidades de venta**: Contado, crédito, reservas con depósito
- **Inventario diverso**: Miles de SKUs con alta rotación temporal
- **Gestión multi-almacén**: Tienda física y bodega central

### Problemática

Antes de este sistema, el negocio enfrentaba:

1. **Control manual de inventario** → Errores frecuentes, stock desactualizado
2. **Gestión informal de créditos** → Difícil seguimiento de cuentas por cobrar
3. **Reservas en cuadernos** → Pérdida de depósitos, confusión con apartados
4. **Sin trazabilidad** → Imposible rastrear ventas, devoluciones o movimientos
5. **Listas escolares manuales** → Proceso tedioso y propenso a errores
6. **Falta de reportes** → Decisiones basadas en intuición, no en datos

### Objetivos del Sistema

1. **Digitalizar operaciones** del negocio completo
2. **Automatizar procesos** críticos (ventas, inventario, notificaciones)
3. **Centralizar información** de clientes, productos y transacciones
4. **Generar reportes** para toma de decisiones
5. **Mejorar experiencia del cliente** con reservas y seguimiento de créditos
6. **Optimizar inventario** para temporadas escolares
7. **Facilitar gestión de listas escolares** por grado y escuela

---

## 🏢 Descripción del Sistema

### Visión General

Sistema integral de gestión empresarial (ERP simplificado) diseñado específicamente para comercios de útiles escolares y papelería, con capacidades de:

- Gestión completa de inventario multi-almacén
- Ventas en efectivo y crédito
- Cuentas por cobrar con seguimiento de pagos
- Sistema de reservas con depósitos
- Gestión de clientes y sus hijos (para útiles escolares)
- Listas escolares personalizables por grado
- Notificaciones automatizadas
- Reportes y analíticas

### Características Especializadas

#### 1. **Gestión de Útiles Escolares**

- Productos categorizados por nivel escolar (preescolar, primaria, bachillerato, universidad)
- Listas escolares predefinidas por grado
- Sugerencias de productos según edad del niño
- Inventario preparado para temporadas escolares

#### 2. **Sistema de Reservas con Depósito**

- Cliente reserva productos con un depósito
- Productos quedan apartados en inventario (estado "reservado")
- Fecha de vencimiento configurable
- Notificaciones automáticas antes del vencimiento
- Liberación automática de inventario si vence

#### 3. **Ventas a Crédito y Cuentas por Cobrar**

- Ventas a crédito con plazos configurables
- Generación automática de cuenta por cobrar
- Registro de pagos parciales y totales
- Estados automáticos: Pendiente → Parcialmente Pagado → Pagado
- Identificación de cuentas vencidas
- Historial completo de pagos

#### 4. **Gestión de Clientes Multi-dimensional**

- Clientes individuales y empresariales (escuelas, empresas)
- Registro de hijos con fecha de nacimiento y grado escolar
- Puntos de lealtad
- Historial de compras
- Preferencias de productos

### Usuarios del Sistema

1. **Administrador**: Acceso total, configuración del sistema
2. **Vendedor**: Ventas, reservas, consulta de inventario
3. **Cajero**: Registro de pagos, cuentas por cobrar
4. **Almacenista**: Movimientos de inventario, recepciones
5. **Gerente**: Reportes, analíticas, supervisión

---

## 🧩 Módulos del Core

### 1. **Módulo de Productos** (`product_service.go`)

**Responsabilidades:**
- CRUD de productos
- Gestión de categorías
- Actualización de precios con historial
- Búsqueda y filtrado
- Productos escolares con niveles educativos

**Entidades Principales:**
- `Product`: Producto con SKU, precio, stock mínimo/máximo
- `Category`: Categorías jerárquicas
- `UnitOfMeasure`: Unidades de medida
- `ProductPriceHistory`: Historial de cambios de precio

**Casos de Uso:**
- Crear producto nuevo
- Actualizar precio (registra en historial)
- Buscar productos por nombre, SKU o categoría
- Marcar producto como útil escolar
- Asignar niveles educativos a productos

---

### 2. **Módulo de Inventario** (`inventory_service.go`)

**Responsabilidades:**
- Control de stock por almacén
- Movimientos de inventario (entrada, salida, ajuste, transferencia)
- Cantidad disponible vs reservada vs en tránsito
- Alertas de stock mínimo
- Trazabilidad completa

**Entidades Principales:**
- `Inventory`: Stock por producto y almacén
- `InventoryMovement`: Registro de movimientos
- `Warehouse`: Almacenes/bodegas

**Tipos de Movimiento:**
- **INBOUND**: Entrada de mercancía (compra, producción, devolución)
- **OUTBOUND**: Salida de mercancía (venta, pérdida, donación)
- **ADJUSTMENT**: Ajuste de inventario (conteo físico)
- **TRANSFER**: Transferencia entre almacenes
- **RESERVATION**: Apartado para reserva
- **RELEASE**: Liberación de reserva

**Casos de Uso:**
- Registrar entrada de mercancía
- Registrar venta (salida automática)
- Ajustar inventario tras conteo físico
- Transferir entre almacenes
- Verificar disponibilidad antes de venta/reserva
- Alertar cuando stock < stock mínimo

---

### 3. **Módulo de Clientes** (`customer_handler.go`)

**Responsabilidades:**
- CRUD de clientes
- Gestión de hijos (para útiles escolares)
- Puntos de lealtad
- Consulta de historial de compras

**Entidades Principales:**
- `Customer`: Cliente (persona o empresa)
- `CustomerChild`: Hijos del cliente con edad y grado escolar
- `Location`: Ubicación geográfica del cliente

**Tipos de Cliente:**
- Individual (persona natural)
- Empresarial (escuelas, empresas)

**Casos de Uso:**
- Registrar nuevo cliente
- Agregar hijos con fecha de nacimiento y grado
- Acumular puntos de lealtad por compras
- Consultar historial de compras
- Buscar por CI/RIF

---

### 4. **Módulo de Ventas** (`sale_service.go`)

**Responsabilidades:**
- Crear ventas en efectivo y crédito
- Generar facturas
- Cálculo de impuestos
- Descuentos
- Cancelación de ventas

**Entidades Principales:**
- `Sale`: Venta con total, impuestos, descuentos
- `SaleDetail`: Detalle de productos vendidos
- `AccountsReceivable`: Cuenta por cobrar (si es venta a crédito)

**Tipos de Venta:**
- **CASH**: Venta de contado
- **CREDIT**: Venta a crédito
- **LAYAWAY**: Venta por apartado (reserva + pago final)

**Estados de Venta:**
- **PENDING**: En proceso
- **COMPLETED**: Completada y pagada
- **CANCELLED**: Cancelada
- **REFUNDED**: Reembolsada

**Casos de Uso:**
- Crear venta de contado
- Crear venta a crédito (genera cuenta por cobrar automáticamente)
- Cancelar venta (revierte inventario si aplica)
- Aplicar descuentos
- Consultar ventas diarias

---

### 5. **Módulo de Cuentas por Cobrar** (`accounts_receivable_service.go`)

**Responsabilidades:**
- Seguimiento de deudas de clientes
- Registro de pagos parciales y totales
- Cálculo automático de saldos
- Identificación de cuentas vencidas

**Entidades Principales:**
- `AccountsReceivable`: Cuenta por cobrar
- `CustomerPayment`: Pagos realizados

**Estados:**
- **PENDING**: Sin pagar
- **PARTIALLY_PAID**: Pago parcial
- **PAID**: Totalmente pagado
- **OVERDUE**: Vencido

**Casos de Uso:**
- Registrar pago parcial (actualiza estado a PARTIALLY_PAID)
- Registrar pago completo (actualiza estado a PAID)
- Consultar cuentas vencidas
- Generar reporte de cuentas por cobrar
- Consultar historial de pagos de un cliente

**Lógica de Pagos (Transaccional):**
```
1. Recibir pago
2. Crear registro de pago
3. Actualizar monto pagado en cuenta
4. Recalcular saldo
5. Actualizar estado automáticamente:
   - Si pago total >= deuda → PAID
   - Si 0 < pago < deuda → PARTIALLY_PAID
```

---

### 6. **Módulo de Reservas** (`reservation_service.go`)

**Responsabilidades:**
- Crear reservas con depósito
- Apartar inventario
- Confirmar y cumplir reservas
- Notificaciones automáticas
- Cancelación y devolución

**Entidades Principales:**
- `Reservation`: Reserva con depósito y fecha de vencimiento
- `ReservationItem`: Items reservados
- `PreOrder`: Pre-orden de productos no disponibles

**Estados de Reserva:**
- **PENDING**: Creada, esperando confirmación
- **CONFIRMED**: Confirmada por cliente
- **FULFILLED**: Cumplida (venta completada)
- **CANCELLED**: Cancelada
- **EXPIRED**: Vencida

**Casos de Uso:**
- Crear reserva (aparta inventario, registra depósito)
- Confirmar reserva (cliente acepta)
- Cumplir reserva (convertir en venta, aplicar depósito como pago)
- Cancelar reserva (libera inventario, maneja devolución de depósito)
- Enviar recordatorio automático antes de vencer

---

### 7. **Módulo de Notificaciones** (`notification_service.go`)

**Responsabilidades:**
- Envío de notificaciones por email/SMS
- Confirmaciones de reserva
- Recordatorios de vencimiento
- Notificaciones de pre-órdenes listas

**Entidades Principales:**
- `CustomerNotification`: Notificación a enviar

**Tipos de Notificación:**
- **EMAIL**: Correo electrónico
- **SMS**: Mensaje de texto
- **PUSH**: Notificación push (app móvil)
- **WHATSAPP**: Mensaje de WhatsApp

**Estados:**
- **PENDING**: Pendiente de envío
- **SENT**: Enviada
- **FAILED**: Falló el envío
- **READ**: Leída por el cliente

**Casos de Uso:**
- Enviar confirmación al crear reserva
- Enviar recordatorio 48h antes de vencimiento
- Notificar cuando pre-orden está lista
- Envío masivo de promociones

---

### 8. **Módulo de Listas Escolares** (Dominio definido)

**Responsabilidades:**
- Crear listas por grado escolar
- Asignar productos a listas
- Sugerir alternativas más económicas
- Calcular costo total de lista

**Entidades Principales:**
- `SchoolSupplyList`: Lista escolar por grado
- `SchoolSupplyListItem`: Item de la lista
- `ListItemAlternative`: Productos alternativos

**Niveles Escolares:**
- PRESCHOOL (Preescolar)
- PRIMARY (Primaria)
- MIDDLE_SCHOOL (Secundaria)
- HIGH_SCHOOL (Bachillerato)
- UNIVERSITY (Universidad)

**Casos de Uso (Pendiente de implementar):**
- Crear lista para "3er grado primaria"
- Agregar 20 lápices, 5 cuadernos, etc.
- Sugerir alternativa económica para marcadores
- Calcular precio total de la lista

---

## ⚖️ Reglas de Negocio

### Productos

1. **SKU único**: No pueden existir dos productos con el mismo SKU
2. **Precio mínimo**: El precio de venta debe ser ≥ precio de costo
3. **Stock mínimo**: Si stock < stock mínimo → Alerta
4. **Cambio de precio**: Todo cambio de precio se registra en historial
5. **Productos escolares**: Deben tener al menos un nivel escolar asignado

### Inventario

1. **No vender sin stock**: No se permite venta si cantidad disponible < cantidad solicitada
2. **Reservas apartan stock**: Cantidad reservada no está disponible para venta
3. **Movimientos trazables**: Todo movimiento debe tener referencia (venta, compra, ajuste)
4. **Stock nunca negativo**: Validación para evitar stock negativo
5. **Transferencias balanceadas**: Salida de almacén A = Entrada a almacén B

### Clientes

1. **Identificación única**: CI/RIF único por cliente
2. **Hijos con edad**: Fecha de nacimiento obligatoria para hijos
3. **Puntos de lealtad**: Se acumulan por monto de compra (configurable)
4. **Cliente empresarial**: Debe tener razón social
5. **Cliente individual**: Debe tener nombre y apellido

### Ventas

1. **Cliente requerido**: Toda venta debe tener cliente asociado
2. **Venta a crédito**: Solo clientes autorizados pueden comprar a crédito
3. **Descuento máximo**: Descuento no puede superar el total de la venta
4. **Cancelación**: Solo se pueden cancelar ventas en estado COMPLETED
5. **Impuestos**: IVA aplicable según configuración de producto

### Cuentas por Cobrar

1. **Pago no excede deuda**: Monto de pago ≤ saldo pendiente
2. **Actualización automática**: Estado se actualiza automáticamente al registrar pago
3. **Vencimiento**: Cuenta vencida si fecha_vencimiento < hoy y estado != PAID
4. **Intereses**: Cuentas vencidas pueden generar intereses (configurable)
5. **Una cuenta por venta**: Cada venta a crédito genera exactamente una cuenta por cobrar

### Reservas

1. **Depósito mínimo**: Depósito ≥ 30% del total (configurable)
2. **Plazo máximo**: Reserva no puede durar más de 30 días (configurable)
3. **Stock disponible**: Solo se puede reservar si hay stock disponible
4. **Vencimiento automático**: Si pasa fecha límite sin confirmar → Estado EXPIRED
5. **Devolución de depósito**: Si cancelación por parte del negocio → devolución total

---

## 🔄 Flujos de Trabajo

### Flujo 1: Venta de Contado

```
1. Cliente llega con productos
2. Vendedor escanea productos
3. Sistema verifica stock disponible
4. Sistema calcula total + impuestos - descuentos
5. Cliente paga
6. Sistema registra venta (tipo: CASH, estado: COMPLETED)
7. Sistema descuenta inventario automáticamente
8. Sistema imprime factura
```

### Flujo 2: Venta a Crédito

```
1. Cliente solicita compra a crédito
2. Sistema verifica que cliente esté autorizado para crédito
3. Vendedor crea venta (tipo: CREDIT)
4. Sistema crea cuenta por cobrar automáticamente:
   - Total: monto de la venta
   - Saldo: monto total
   - Fecha vencimiento: hoy + plazo de crédito
   - Estado: PENDING
5. Sistema descuenta inventario
6. Sistema genera factura
7. Sistema envía notificación al cliente
```

### Flujo 3: Pago de Cuenta por Cobrar

```
1. Cliente llega a pagar deuda
2. Cajero busca cuenta por cobrar
3. Cajero registra monto del pago
4. Sistema ejecuta transacción:
   a. Crear registro de pago
   b. Actualizar monto_pagado
   c. Recalcular saldo
   d. Actualizar estado:
      - Si pago completo → PAID
      - Si pago parcial → PARTIALLY_PAID
5. Sistema imprime recibo de pago
6. Si cuenta queda en PAID → enviar notificación de agradecimiento
```

### Flujo 4: Crear Reserva

```
1. Cliente solicita reservar productos
2. Sistema verifica disponibilidad de todos los items
3. Cliente paga depósito (mínimo 30%)
4. Sistema crea reserva:
   - Estado: PENDING
   - Depósito registrado
   - Fecha vencimiento: hoy + 15 días
5. Sistema aparta inventario (marca como "reservado")
6. Sistema envía confirmación por email/WhatsApp
7. Sistema programa recordatorio para 2 días antes de vencer
```

### Flujo 5: Cumplir Reserva

```
1. Cliente regresa para completar compra
2. Vendedor busca reserva
3. Cliente paga saldo restante (Total - Depósito)
4. Sistema:
   a. Crea venta con total completo
   b. Aplica depósito como pago inicial
   c. Registra pago del saldo
   d. Actualiza estado de reserva a FULFILLED
   e. Descuenta inventario (libera "reservado" y descuenta "disponible")
5. Sistema imprime factura
6. Sistema envía agradecimiento
```

### Flujo 6: Reserva Vencida

```
1. Sistema ejecuta tarea programada diaria
2. Sistema busca reservas con estado PENDING donde fecha_vencimiento < hoy
3. Para cada reserva vencida:
   a. Actualizar estado a EXPIRED
   b. Liberar inventario reservado
   c. Procesar devolución de depósito (según política)
   d. Enviar notificación al cliente
```

### Flujo 7: Crear Lista Escolar para Cliente

```
1. Cliente llega con lista de útiles de la escuela
2. Vendedor registra:
   - Grado escolar del niño
   - Nombre de la escuela
3. Sistema busca si existe lista predefinida para ese grado
4. Si existe:
   a. Carga lista predefinida
   b. Vendedor ajusta cantidades según lista física
5. Si no existe:
   a. Vendedor crea lista desde cero
6. Sistema:
   a. Calcula total
   b. Muestra alternativas económicas
   c. Indica productos sin stock
7. Cliente decide:
   - Comprar todo (venta normal)
   - Reservar con depósito
   - Pre-ordenar faltantes
```

---

## 📊 Modelos del Dominio

### Modelo de Producto

```go
type Product struct {
    ProductID      uuid.UUID
    SKU            string          // Único
    Name           string
    Description    *string
    CategoryID     *uuid.UUID
    UnitID         *uuid.UUID
    CostPrice      *float64
    SellingPrice   float64         // Precio de venta
    PriceCurrency  CurrencyCode    // VES, USD
    MinStock       int             // Stock mínimo
    MaxStock       int             // Stock máximo
    HasTax         bool            // ¿Aplica IVA?
    TaxPercentage  float64         // % de IVA
    Status         ProductStatus   // ACTIVE, INACTIVE, DISCONTINUED

    // Útiles escolares
    IsSchoolSupply bool
    GradeLevels    []SchoolLevel   // PRESCHOOL, PRIMARY, etc.
    SeasonalDemand bool            // ¿Demanda estacional?

    // Auditoría
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      *time.Time
    CreatedBy      *uuid.UUID
    UpdatedBy      *uuid.UUID
}
```

### Modelo de Venta

```go
type Sale struct {
    SaleID         uuid.UUID
    InvoiceNumber  string          // Único
    CustomerID     *uuid.UUID
    SalespersonID  *uuid.UUID
    StoreID        *uuid.UUID
    SaleDate       time.Time

    // Montos
    SubTotal       float64         // Subtotal sin impuestos
    TaxAmount      float64         // Monto de impuestos
    DiscountAmount float64         // Descuento aplicado
    TotalAmount    float64         // Total final
    Currency       CurrencyCode

    // Tipo y estado
    SaleType       SaleType        // CASH, CREDIT, LAYAWAY
    Status         SaleStatus      // PENDING, COMPLETED, CANCELLED

    // Método de pago
    PaymentMethod  *PaymentMethod  // CASH, CARD, TRANSFER, etc.

    // Relaciones
    Details        []SaleDetail
    AccountReceivable *AccountsReceivable
}

type SaleDetail struct {
    DetailID       uuid.UUID
    SaleID         uuid.UUID
    ProductID      uuid.UUID
    Quantity       float64
    UnitPrice      float64         // Precio al momento de venta
    Discount       float64
    TaxAmount      float64
    Subtotal       float64
    Total          float64
}
```

### Modelo de Cuenta por Cobrar

```go
type AccountsReceivable struct {
    ReceivableID   uuid.UUID
    SaleID         *uuid.UUID
    CustomerID     uuid.UUID

    // Montos
    TotalAmount    float64         // Deuda total
    PaidAmount     float64         // Pagado hasta ahora
    Balance        float64         // Saldo pendiente
    Currency       CurrencyCode

    // Fechas
    DueDate        time.Time       // Fecha límite de pago
    Status         AccountStatus   // PENDING, PARTIALLY_PAID, PAID, OVERDUE

    // Relaciones
    Payments       []CustomerPayment
}

type CustomerPayment struct {
    PaymentID      uuid.UUID
    ReceivableID   uuid.UUID
    Amount         float64
    Currency       CurrencyCode
    PaymentDate    time.Time
    PaymentMethod  PaymentMethod
    Reference      *string         // # de referencia bancaria
}
```

### Modelo de Reserva

```go
type Reservation struct {
    ReservationID     uuid.UUID
    ReservationNumber string        // Único, ej: RES-2024-00001
    CustomerID        uuid.UUID

    // Montos
    TotalAmount       float64
    DepositAmount     float64       // Depósito pagado
    Balance           float64       // Saldo pendiente
    Currency          CurrencyCode

    // Fechas
    ReservationDate   time.Time     // Fecha de creación
    ExpirationDate    time.Time     // Fecha límite

    // Estado
    Status            ReservationStatus // PENDING, CONFIRMED, FULFILLED, CANCELLED, EXPIRED

    // Relaciones
    Items             []ReservationItem
    Notifications     []CustomerNotification
}

type ReservationItem struct {
    ItemID         uuid.UUID
    ReservationID  uuid.UUID
    ProductID      uuid.UUID
    Quantity       float64
    UnitPrice      float64
    Subtotal       float64
}
```

### Modelo de Cliente

```go
type Customer struct {
    CustomerID     uuid.UUID
    TaxID          string          // CI/RIF, único
    CustomerType   CustomerType    // INDIVIDUAL, BUSINESS

    // Persona individual
    FirstName      *string
    LastName       *string
    DateOfBirth    *time.Time

    // Empresa
    BusinessName   *string
    TradeName      *string

    // Contacto
    Email          *string
    Phone          *string
    LocationID     *uuid.UUID
    Address        *string

    // Lealtad
    LoyaltyPoints  int
    LoyaltyLevel   *LoyaltyLevel   // BRONZE, SILVER, GOLD, PLATINUM

    // Estado
    Status         CustomerStatus  // ACTIVE, INACTIVE, BLOCKED

    // Relaciones
    Children       []CustomerChild
    Sales          []Sale
    Reservations   []Reservation
}

type CustomerChild struct {
    ChildID        uuid.UUID
    CustomerID     uuid.UUID
    FirstName      string
    LastName       string
    DateOfBirth    *time.Time
    SchoolLevel    *SchoolLevel    // PRESCHOOL, PRIMARY, etc.
    Grade          *string         // "3er grado"
    SchoolName     *string
}
```

---

## 🎨 Casos de Uso Especiales

### Caso 1: Temporada Escolar

**Escenario**: Agosto-Septiembre, inicio del año escolar

**Acciones del Sistema**:
1. Incrementar stock de productos escolares
2. Crear listas escolares predefinidas por grado
3. Activar promociones para paquetes escolares
4. Permitir pre-órdenes para productos agotados
5. Enviar notificaciones a clientes con hijos en edad escolar
6. Generar reportes de demanda por nivel educativo

### Caso 2: Cliente VIP con Crédito

**Escenario**: Cliente frecuente con historial de pago impecable

**Acciones del Sistema**:
1. Clasificar cliente como VIP (Gold/Platinum)
2. Otorgar plazo de crédito extendido (45-60 días vs 30 días)
3. Aplicar descuentos automáticos
4. Priorizar en notificaciones de nuevos productos
5. Permitir reservas sin depósito (solo VIP Platinum)

### Caso 3: Producto Agotado con Alta Demanda

**Escenario**: Producto requerido pero sin stock

**Acciones del Sistema**:
1. Permitir crear pre-orden
2. Registrar cantidad requerida
3. Notificar a proveedor automáticamente
4. Cliente paga depósito
5. Al llegar mercancía: notificar a todos los clientes con pre-orden
6. Prioridad de entrega según orden de pre-orden

---

**Documento vivo** - Se actualiza conforme evoluciona el sistema

**Autor**: Jesus Dicen
**Última Actualización**: 2025-12-23
