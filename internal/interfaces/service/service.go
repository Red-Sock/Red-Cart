package service

type Service interface {
	User() UserService
	Cart() CartService
	CartItem() CartItem
}
