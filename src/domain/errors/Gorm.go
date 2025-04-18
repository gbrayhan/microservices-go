package errors

type GormErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}
