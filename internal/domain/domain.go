package domain

type Order struct {
	OrderUid          string   `json:"order_uid"          db:"order_uid"          validate:"required"`
	TrackNumber       string   `json:"track_number"       db:"track_number"       validate:"required"`
	Entry             string   `json:"entry"              db:"entry"              validate:"required"`
	Delivery          Delivery `json:"delivery"                                   validate:"required"`
	Payment           Payment  `json:"payment"                                    validate:"required"`
	Items             []Item   `json:"items"                                      validate:"required,dive"`
	Locale            string   `json:"locale"             db:"locale"             validate:"required"`
	InternalSignature string   `json:"internal_signature" db:"internal_signature" validate:"required"`
	CustomerID        string   `json:"customer_id"        db:"customer_id"        validate:"required"`
	DeliveryService   string   `json:"delivery_service"   db:"delivery_service"   validate:"required"`
	Shardkey          string   `json:"shardkey"           db:"shardkey"`
	SmID              int64    `json:"sm_id"              db:"sm_id"              validate:"required"`
	DateCreated       string   `json:"date_created" db:"date_created"`
	OofShard          string   `json:"oof_shard"          db:"oof_shard"`
}

type Delivery struct {
	OrderUid string `db:"order_uid"`

	Name    string `json:"name"    db:"name"    validate:"required"`
	Phone   string `json:"phone"   db:"phone"   validate:"required"`
	Zip     string `json:"zip"     db:"zip"     validate:"required"`
	City    string `json:"city"    db:"city"    validate:"required"`
	Address string `json:"address" db:"address" validate:"required"`
	Region  string `json:"region"  db:"region"  validate:"required"`
	Email   string `json:"email"   db:"email"   validate:"required"`
}

type Item struct {
	ID       int    `db:"id"`
	OrderUid string `db:"order_uid"`

	ChrtID      int64  `json:"chrt_id"      db:"chrt_id"      validate:"required"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required"`
	Price       int64  `json:"price"        db:"price"        validate:"gt=0"`
	Rid         string `json:"rid"          db:"rid"          validate:"required"`
	Name        string `json:"name"         db:"name"         validate:"required"`
	Sale        int64  `json:"sale"         db:"sale"         validate:"gte=0"`
	Size        string `json:"size"         db:"size"         validate:"required"`
	TotalPrice  int64  `json:"total_price"  db:"total_price"  validate:"gt=0"`
	NmID        int64  `json:"nm_id"        db:"nm_id"        validate:"required"`
	Brand       string `json:"brand"        db:"brand"        validate:"required"`
	Status      int64  `json:"status"       db:"status"       validate:"required"`
}

type Payment struct {
	OrderUid string `db:"order_uid"`

	Transaction  string `json:"transaction"   db:"transaction"   validate:"required"`
	RequestID    string `json:"request_id"    db:"request_id"    validate:"required"`
	Currency     string `json:"currency"      db:"currency"      validate:"required"`
	Provider     string `json:"provider"      db:"provider"      validate:"required"`
	Amount       int64  `json:"amount"        db:"amount"        validate:"gt=0"`
	PaymentDt    int64  `json:"payment_dt"    db:"payment_dt"    validate:"required"`
	Bank         string `json:"bank"          db:"bank"          validate:"required"`
	DeliveryCost int64  `json:"delivery_cost" db:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int64  `json:"goods_total"   db:"goods_total"   validate:"gt=0"`
	CustomFee    int64  `json:"custom_fee"    db:"custom_fee"    validate:"gte=0"`
}
