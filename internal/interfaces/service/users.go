package service

type UserService interface {
	Start(id int64) (message string)
}
