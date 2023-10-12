package service

type Storage interface {
	User() UserService
	Cart() CartService
}
