package controllers

type JSONSwagger struct {
}

type generalResponse struct {
	Messages []string `json:"messages" example:"error description,other error description"`
}
