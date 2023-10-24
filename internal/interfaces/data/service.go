package data

type Storage interface {
	User() Users
	Cart() Carts
	CartsItem() CartsItem
}
