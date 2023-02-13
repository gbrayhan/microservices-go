// Package controllers contains the common functions and structures for the controllers
package controllers

type JSONSwagger struct {
}

// MessageResponse is a struct that contains the response body for the message
type MessageResponse struct {
	Message string `json:"message"`
}
