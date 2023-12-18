package dto

type RequestWithID struct {
	ID string `json:"id" validate:"required"`
}

type ResponseWithID struct {
	ID string `json:"id"`
}

type ResponseOK struct {
	Success string `json:"success"`
}
