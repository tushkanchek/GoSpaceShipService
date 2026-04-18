package converter

import (
	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"
)

func OrderToApi(order *model.Order) *orderV1.Order {
	var transactionUUID orderV1.OptUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptUUID(*order.TransactionUUID)
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		pm := PaymentMethodModelToApi(*order.PaymentMethod)
		paymentMethod = orderV1.NewOptPaymentMethod(pm)
	}

	return &orderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      float32(order.TotalPrice),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderV1.OrderStatus(order.OrderStatus),
	}
}

func PaymentMethodApiToModel(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func PaymentMethodModelToApi(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
