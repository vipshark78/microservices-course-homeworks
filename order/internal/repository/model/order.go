package model

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   string
	Status          string
}

type CreateOrder struct {
	UserUUID  string
	PartUuids []string
}

const (
	OrderStatusPENDINGPAYMENT = "PENDING_PAYMENT"
	OrderStatusPAID           = "PAID"
	OrderStatusCANCELLED      = "CANCELLED"
)

const (
	PaymentMethodUNKNOWN       = "UNKNOWN"
	PaymentMethodCARD          = "CARD"
	PaymentMethodSBP           = "SBP"
	PaymentMethodCREDITCARD    = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY = "INVESTOR_MONEY"
)
