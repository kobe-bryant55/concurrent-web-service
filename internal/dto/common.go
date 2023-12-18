package dto

type RequestWithID struct {
	ID string `json:"id" validate:"required"`
}

type ResponseWithID struct {
	ID string `json:"id"`
}
