package domain

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"math"
)

var validate *validator.Validate

const PaymentAmountCheckTag = "payment_amount_check"
const ItemPriceCheckTag = "item_price_check"
const GoodsTotalCheckTag = "goods_total_check"

func init() {
	validate = validator.New()
	validate.RegisterStructValidation(orderValidator, Order{})
	if err := validate.RegisterValidation(ItemPriceCheckTag, itemPriceValidator); err != nil {
		return
	}
	if err := validate.RegisterValidation(PaymentAmountCheckTag, paymentAmountValidator); err != nil {
		return
	}
}

func ValidateOrder(order *Order) error {
	return validate.Struct(order)
}

func orderValidator(sl validator.StructLevel) {
	order := sl.Current().Interface().(Order)
	var sumTotalPrice int64
	for i, item := range order.Items {
		err := validate.VarWithValue(item.TotalPrice, item, ItemPriceCheckTag)
		if err != nil {
			sl.ReportError(
				item.TotalPrice,
				fmt.Sprintf("items[%d].total_price", i),
				"TotalPrice",
				ItemPriceCheckTag,
				"",
			)
		}
		sumTotalPrice += item.TotalPrice
	}

	if order.Payment.GoodsTotal != sumTotalPrice {
		sl.ReportError(
			order.Payment.GoodsTotal,
			"payment.goods_total",
			"GoodsTotal",
			GoodsTotalCheckTag,
			"",
		)
	}
	if err := validate.Var(order.Payment, PaymentAmountCheckTag); err != nil {
		sl.ReportError(
			order.Payment.Amount,
			"payment.amount",
			"Amount",
			PaymentAmountCheckTag,
			"",
		)
	}
}

func itemPriceValidator(fl validator.FieldLevel) bool {
	totalPrice := fl.Field().Int()
	item := fl.Parent().Interface().(Item)
	expected := int64(math.Round(float64(item.Price*(100-item.Sale)) / 100))
	return totalPrice == expected
}

func paymentAmountValidator(fl validator.FieldLevel) bool {
	payment := fl.Field().Interface().(Payment)
	expectedSum := payment.GoodsTotal + payment.DeliveryCost + payment.CustomFee
	return payment.Amount == expectedSum
}
