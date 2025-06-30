package domain

import (
	"encoding/json"
	"time"
)

func UnmarshalOrder(data []byte) (Order, error) {
	var r Order
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Order) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Order struct {
	OrderUid          string    `json:"order_uid"          validate:"required"`
	TrackNumber       string    `json:"track_number"       validate:"required"`
	Entry             string    `json:"entry"              validate:"required"`
	Delivery          Delivery  `json:"delivery"           validate:"required, dive"`
	Payment           Payment   `json:"payment"            validate:"required, dive"`
	Items             []Item    `json:"items"              validate:"required, dive"`
	Locale            string    `json:"locale"             validate:"required"`
	InternalSignature string    `json:"internal_signature" validate:"required"`
	CustomerID        string    `json:"customer_id"        validate:"required"`
	DeliveryService   string    `json:"delivery_service"   validate:"required"`
	Shardkey          string    `json:"shardkey"`
	SmID              int64     `json:"sm_id"              validate:"required"`
	DateCreated       time.Time `json:"date_created"       validate:"required"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"    validate:"required"`
	Phone   string `json:"phone"   validate:"required"`
	Zip     string `json:"zip"     validate:"required"`
	City    string `json:"city"    validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"  validate:"required"`
	Email   string `json:"email"   validate:"required"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id"      validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int64  `json:"price"        validate:"required, gt=0"`
	Rid         string `json:"rid"          validate:"required"`
	Name        string `json:"name"         validate:"required"`
	Sale        int64  `json:"sale"         validate:"required, gte=0"`
	Size        string `json:"size"         validate:"required"`
	TotalPrice  int64  `json:"total_price"  validate:"required, gt=0"`
	NmID        int64  `json:"nm_id"        validate:"required"`
	Brand       string `json:"brand"        validate:"required"`
	Status      int64  `json:"status"       validate:"required"`
}

type Payment struct {
	Transaction  string `json:"transaction"   validate:"required"`
	RequestID    string `json:"request_id"    validate:"required"`
	Currency     string `json:"currency"      validate:"required"`
	Provider     string `json:"provider"      validate:"required"`
	Amount       int64  `json:"amount"        validate:"required, gt=0"`
	PaymentDt    int64  `json:"payment_dt"    validate:"required"`
	Bank         string `json:"bank"          validate:"required"`
	DeliveryCost int64  `json:"delivery_cost" validate:"required, gte=0"`
	GoodsTotal   int64  `json:"goods_total"   validate:"required, gt=0"`
	CustomFee    int64  `json:"custom_fee"    validate:"required, gte=0"`
}
