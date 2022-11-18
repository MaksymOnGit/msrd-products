package models

type PostgreSqlEvent struct {
	After     *Value `json:"after"`
	Operation string `json:"op"`
}

type Value struct {
	ProductId      string  `json:"product_id"`
	ActualQuantity float32 `json:"quantity_actual"`
}
