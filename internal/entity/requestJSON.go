package entity

type UpdatePerson struct {
	ID     string `json:"id"`
	Values Person `json:"values"`
}

type DeletePerson struct {
	ID string `json:"id"`
}
