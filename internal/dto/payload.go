package dto

type PayloadCollection struct {
	Payloads []Payload `json:"payloads" validate:"required"`
}

type OperationType string

const (
	OpPut    OperationType = "OpPut"
	OpCreate OperationType = "OpCreate"
	OpDelete OperationType = "OpDelete"
)

type Payload struct {
	OperationType `json:"operationType" validate:"required"`
	Data          interface{} `json:"data"          validate:"required"`
}
