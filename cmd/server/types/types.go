package types

import "net/http"

type ServerConfig struct {
	Host string
	Port int
}

type Server struct {
	Config ServerConfig
	Router *http.ServeMux
}

type Middleware func(http.Handler) http.Handler

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Result string `json:"result"`
}


