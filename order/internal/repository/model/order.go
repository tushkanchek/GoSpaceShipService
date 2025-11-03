package model

type PaymentMethod int32

const (
	PaymentMethodUNKNOWN PaymentMethod = 0
	PaymentMethodCARD PaymentMethod = 1
	PaymentMethodSBP PaymentMethod = 2
	PaymentMethodCreditCard PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)
type OrderStatus string

const (
	OrderStatusUNKNOWN        OrderStatus = "UNKNOWN"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type Order struct{
	OrderUUID string
	UserUUID string
	PartUuids []string 
	TotalPrice float64
	TransactionUUID *string   
	PaymentMethod *PaymentMethod
	OrderStatus OrderStatus
}