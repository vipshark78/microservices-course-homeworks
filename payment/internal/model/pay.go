package model

type PayOrder struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod string
}

const (
	PaymentMethodUNKNOWN       = "UNKNOWN"
	PaymentMethodCARD          = "CARD"
	PaymentMethodSBP           = "SBP"
	PaymentMethodCREDITCARD    = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY = "INVESTOR_MONEY"
)
