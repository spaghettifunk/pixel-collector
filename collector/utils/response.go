package utils

// Response is a wrapper object to send back to the client
type Response struct {
	Data interface{} `json:"data"`
}

// ResponseError is a wrapper object to send back to the client when
// there is an error in the API
type ResponseError struct {
	Message string `json:"message"`
}
