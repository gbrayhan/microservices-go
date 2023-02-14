// Package errors defines the domain errors used in the application.
package errors

// GormErr is a struct that contains the error number and message for Gorm errors
type GormErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}
