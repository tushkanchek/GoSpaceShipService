package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	OrderStatus     OrderStatus
}

type PaymentMethod int32

const (
	PaymentMethodUNKNOWN       PaymentMethod = 0
	PaymentMethodCARD          PaymentMethod = 1
	PaymentMethodSBP           PaymentMethod = 2
	PaymentMethodCreditCard    PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)

// Реализуем интерфейс Scaner, для извлечения PaymentMethod из postgres
func (p *PaymentMethod) Scan(src any) error {
	switch v := src.(type) {
	case string:
		switch v {
		case "CARD":
			*p = PaymentMethodCARD
		case "SBP":
			*p = PaymentMethodSBP
		case "Credit Card":
			*p = PaymentMethodCreditCard
		case "Investor Money":
			*p = PaymentMethodInvestorMoney
		default:
			*p = PaymentMethodUNKNOWN
		}
	default:
		return fmt.Errorf("unsupported type for PaymentMethod: %T", v)
	}

	return nil
}

// Реализуем интерфейс Value, для вставки PaymentMethod в postgres
func (p PaymentMethod) Value() (driver.Value, error) {
	switch p {
	case PaymentMethodSBP:
		return "SBP", nil
	case PaymentMethodInvestorMoney:
		return "InvestorMoney", nil
	case PaymentMethodCreditCard:
		return "Credit Card", nil
	case PaymentMethodCARD:
		return "CARD", nil
	default:
		return "UNKNOWN", nil
	}
}

type OrderStatus string

const (
	OrderStatusUNKNOWN        OrderStatus = "UNKNOWN"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

// Реализуем интерфейс Scaner, для извлечения OrderStatus из postgres
func (o *OrderStatus) Scan(src any) error {
	switch v := src.(type) {
	case string:
		switch v {
		case "PENDING_PAYMENT":
			*o = OrderStatusPENDINGPAYMENT
		case "PAID":
			*o = OrderStatusPAID
		case "CANCELLED":
			*o = OrderStatusCANCELLED
		default:
			*o = OrderStatusUNKNOWN
		}
	default:
		return fmt.Errorf("unsupported type for PaymentMethod: %T", v)
	}
	return nil
}

// Реализуем интерфейс Value, для вставки OrderStatus в postgres
func (o OrderStatus) Value() (driver.Value, error) {
	return string(o), nil
}
