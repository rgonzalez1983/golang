package entity

type JsonResponse struct {
	Message    string `json:"Message,omitempty"`
	StatusCode int
	Data       interface{} `json:"Data,omitempty"`
}
