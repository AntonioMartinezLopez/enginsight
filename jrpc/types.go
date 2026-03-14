package jrpc

const (
	MethodCount = "Count"
)

type CountRequest struct {
	Message string `json:"msg"`
}

type CountResponse struct {
	Count int `json:"count"`
}
