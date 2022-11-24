package documentStatusModels

type PostgreSqlEvent struct {
	After     *Value `json:"after"`
	Operation string `json:"op"`
}

type Value struct {
	DocumentId string `json:"document_id"`
	Status     string `json:"status"`
}
