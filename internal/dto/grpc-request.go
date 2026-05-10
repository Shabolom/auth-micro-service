package dto

type RegisterRequest struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type LoginRequest struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type CreateUserDescriptionRequest struct {
	Name string
	Age  int
}

type UpdateUsersRequest struct {
	Id   string
	Mail string
	Name string
	Age  int32
}
