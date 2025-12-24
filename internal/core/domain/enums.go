package domain

// Location Enums
type LocationType string

const (
	LocationTypeCountry      LocationType = "COUNTRY"
	LocationTypeState        LocationType = "STATE"
	LocationTypeCity         LocationType = "CITY"
	LocationTypeMunicipality LocationType = "MUNICIPALITY"
	LocationTypeParish       LocationType = "PARISH"
	LocationTypeNeighborhood LocationType = "NEIGHBORHOOD"
)

// User and Security Enums
type GenderType string

const (
	GenderMale   GenderType = "M"
	GenderFemale GenderType = "F"
	GenderOther  GenderType = "O"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusLocked    UserStatus = "LOCKED"
)

type AuditActionType string

const (
	AuditActionLogin  AuditActionType = "LOGIN"
	AuditActionLogout AuditActionType = "LOGOUT"
	AuditActionCreate AuditActionType = "CREATE"
	AuditActionRead   AuditActionType = "READ"
	AuditActionUpdate AuditActionType = "UPDATE"
	AuditActionDelete AuditActionType = "DELETE"
	AuditActionApprove AuditActionType = "APPROVE"
	AuditActionReject  AuditActionType = "REJECT"
	AuditActionExport  AuditActionType = "EXPORT"
)

// Employee and Payroll Enums
type EmployeeStatus string

const (
	EmployeeStatusActive     EmployeeStatus = "ACTIVE"
	EmployeeStatusInactive   EmployeeStatus = "INACTIVE"
	EmployeeStatusOnLeave    EmployeeStatus = "ON_LEAVE"
	EmployeeStatusTerminated EmployeeStatus = "TERMINATED"
)

type PayrollConceptType string

const (
	PayrollConceptAllowance            PayrollConceptType = "ALLOWANCE"
	PayrollConceptDeduction            PayrollConceptType = "DEDUCTION"
	PayrollConceptEmployerContribution PayrollConceptType = "EMPLOYER_CONTRIBUTION"
)

type PayrollPeriodStatus string

const (
	PayrollPeriodStatusOpen       PayrollPeriodStatus = "OPEN"
	PayrollPeriodStatusProcessing PayrollPeriodStatus = "PROCESSING"
	PayrollPeriodStatusProcessed  PayrollPeriodStatus = "PROCESSED"
	PayrollPeriodStatusPaid       PayrollPeriodStatus = "PAID"
	PayrollPeriodStatusClosed     PayrollPeriodStatus = "CLOSED"
)

// Customer Enums
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "INDIVIDUAL"
	CustomerTypeBusiness   CustomerType = "BUSINESS"
)

type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "ACTIVE"
	CustomerStatusInactive  CustomerStatus = "INACTIVE"
	CustomerStatusSuspended CustomerStatus = "SUSPENDED"
)

// Product and Inventory Enums
type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "ACTIVE"
	ProductStatusInactive     ProductStatus = "INACTIVE"
	ProductStatusDiscontinued ProductStatus = "DISCONTINUED"
)

type MovementType string

const (
	MovementTypeIn                  MovementType = "IN"
	MovementTypeOut                 MovementType = "OUT"
	MovementTypeAdjustment          MovementType = "ADJUSTMENT"
	MovementTypeTransfer            MovementType = "TRANSFER"
	MovementTypeReservation         MovementType = "RESERVATION"
	MovementTypeReservationRelease  MovementType = "RESERVATION_RELEASE"
)

type StockStatus string

const (
	StockStatusCritical    StockStatus = "CRITICAL"
	StockStatusLow         StockStatus = "LOW"
	StockStatusNormal      StockStatus = "NORMAL"
	StockStatusOverstocked StockStatus = "OVERSTOCKED"
)

// Sales Enums
type SaleType string

const (
	SaleTypeCash      SaleType = "CASH"
	SaleTypeCredit    SaleType = "CREDIT"
	SaleTypeReservation SaleType = "RESERVATION"
	SaleTypePreOrder  SaleType = "PRE_ORDER"
)

type SaleStatus string

const (
	SaleStatusDraft          SaleStatus = "DRAFT"
	SaleStatusCompleted      SaleStatus = "COMPLETED"
	SaleStatusCancelled      SaleStatus = "CANCELLED"
	SaleStatusPendingPayment SaleStatus = "PENDING_PAYMENT"
	SaleStatusReserved       SaleStatus = "RESERVED"
)

type PaymentMethod string

const (
	PaymentMethodCash            PaymentMethod = "CASH"
	PaymentMethodBankTransfer    PaymentMethod = "BANK_TRANSFER"
	PaymentMethodCreditCard      PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard       PaymentMethod = "DEBIT_CARD"
	PaymentMethodMobilePayment   PaymentMethod = "MOBILE_PAYMENT"
	PaymentMethodForeignCurrency PaymentMethod = "FOREIGN_CURRENCY"
	PaymentMethodMixed           PaymentMethod = "MIXED"
)

// Purchase Enums
type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft             PurchaseOrderStatus = "DRAFT"
	PurchaseOrderStatusPending           PurchaseOrderStatus = "PENDING"
	PurchaseOrderStatusApproved          PurchaseOrderStatus = "APPROVED"
	PurchaseOrderStatusPartiallyReceived PurchaseOrderStatus = "PARTIALLY_RECEIVED"
	PurchaseOrderStatusReceived          PurchaseOrderStatus = "RECEIVED"
	PurchaseOrderStatusCancelled         PurchaseOrderStatus = "CANCELLED"
)

// Accounts Enums
type AccountStatus string

const (
	AccountStatusPending        AccountStatus = "PENDING"
	AccountStatusPartiallyPaid  AccountStatus = "PARTIALLY_PAID"
	AccountStatusPaid           AccountStatus = "PAID"
	AccountStatusOverdue        AccountStatus = "OVERDUE"
	AccountStatusCancelled      AccountStatus = "CANCELLED"
)

// Expense Enums
type ExpenseApprovalStatus string

const (
	ExpenseApprovalStatusPending  ExpenseApprovalStatus = "PENDING"
	ExpenseApprovalStatusApproved ExpenseApprovalStatus = "APPROVED"
	ExpenseApprovalStatusRejected ExpenseApprovalStatus = "REJECTED"
)

// Currency Enums
type CurrencyCode string

const (
	CurrencyVES CurrencyCode = "VES"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyEUR CurrencyCode = "EUR"
)

type ExchangeRateSource string

const (
	ExchangeRateSourceBCV      ExchangeRateSource = "BCV"
	ExchangeRateSourceParallel ExchangeRateSource = "PARALLEL"
	ExchangeRateSourceManual   ExchangeRateSource = "MANUAL"
	ExchangeRateSourceOfficial ExchangeRateSource = "OFFICIAL"
)

// Back to School Enums
type SchoolLevel string

const (
	SchoolLevelPreschool   SchoolLevel = "PRESCHOOL"
	SchoolLevelPrimary     SchoolLevel = "PRIMARY"
	SchoolLevelMiddleSchool SchoolLevel = "MIDDLE_SCHOOL"
	SchoolLevelHighSchool  SchoolLevel = "HIGH_SCHOOL"
	SchoolLevelUniversity  SchoolLevel = "UNIVERSITY"
)

type SchoolListStatus string

const (
	SchoolListStatusDraft     SchoolListStatus = "DRAFT"
	SchoolListStatusPublished SchoolListStatus = "PUBLISHED"
	SchoolListStatusActive    SchoolListStatus = "ACTIVE"
	SchoolListStatusArchived  SchoolListStatus = "ARCHIVED"
)

type ReservationStatus string

const (
	ReservationStatusPending           ReservationStatus = "PENDING"
	ReservationStatusConfirmed         ReservationStatus = "CONFIRMED"
	ReservationStatusPartiallyFulfilled ReservationStatus = "PARTIALLY_FULFILLED"
	ReservationStatusFulfilled         ReservationStatus = "FULFILLED"
	ReservationStatusCancelled         ReservationStatus = "CANCELLED"
	ReservationStatusExpired           ReservationStatus = "EXPIRED"
)

type PreOrderStatus string

const (
	PreOrderStatusPending       PreOrderStatus = "PENDING"
	PreOrderStatusConfirmed     PreOrderStatus = "CONFIRMED"
	PreOrderStatusInPreparation PreOrderStatus = "IN_PREPARATION"
	PreOrderStatusReady         PreOrderStatus = "READY"
	PreOrderStatusDelivered     PreOrderStatus = "DELIVERED"
	PreOrderStatusCancelled     PreOrderStatus = "CANCELLED"
)

type NotificationType string

const (
	NotificationTypeEmail    NotificationType = "EMAIL"
	NotificationTypeSMS      NotificationType = "SMS"
	NotificationTypeWhatsApp NotificationType = "WHATSAPP"
	NotificationTypePush     NotificationType = "PUSH"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent    NotificationStatus = "SENT"
	NotificationStatusFailed  NotificationStatus = "FAILED"
	NotificationStatusRead    NotificationStatus = "READ"
)
