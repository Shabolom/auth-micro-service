package dto

type RegisterRequest struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
	Name      string
	Age       int
}

type LoginRequest struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}
