package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	
)

type PaymentMethod int32

const (
	PaymentMethodUNKNOWN       PaymentMethod = 0
	PaymentMethodCARD          PaymentMethod = 1
	PaymentMethodSBP           PaymentMethod = 2
	PaymentMethodCreditCard    PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)

func (p *PaymentMethod) Scan(src any) error{
	switch v := src.(type){
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

func (p PaymentMethod) Value() (driver.Value, error){
	switch p{
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

func (o *OrderStatus) Scan(src any) error{
	switch v := src.(type){
	case string:
		switch v{
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

func (o OrderStatus) Value() (driver.Value, error){
	return string(o), nil
}





// Scan реализует sql.Scanner (для чтения)
// func (u *UUIDArray) Scan(src any) error {
//     if src == nil {
//         *u = []uuid.UUID{}
//         return nil
//     }

//     var s string
//     switch v := src.(type) {
//     case string:
//         s = v
//     case []byte:
//         s = string(v)
//     default:
//         return fmt.Errorf("UUIDArray: expected string or []byte, got %T", src)
//     }

//     s = strings.Trim(s, "{}")
//     if s == "" {
//         *u = []uuid.UUID{}
//         return nil
//     }

//     parts := strings.Split(s, ",")
//     *u = make([]uuid.UUID, 0, len(parts))

//     for _, p := range parts {
//         p = strings.TrimSpace(p)
//         p = strings.Trim(p, `"`)  // <-- важно, чтобы убрать кавычки
//         if p == "" {
//             continue
//         }
//         id, err := uuid.Parse(p)
//         if err != nil {
//             return fmt.Errorf("UUIDArray: parse error for '%s': %w", p, err)
//         }
//         *u = append(*u, id)
//     }

//     return nil
// }


// Value реализует driver.Valuer (для записи)
// func (u UUIDArray) Value() (driver.Value, error) {
//     if len(u) == 0 {
//         return "{}", nil
//     }
//     strs := make([]string, len(u))
//     for i, id := range u {
//         strs[i] = id.String()
//     }
//     return "{" + strings.Join(strs, ",") + "}", nil
// }


type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	OrderStatus     OrderStatus
}


