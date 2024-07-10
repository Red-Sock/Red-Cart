package domain

const (
	DbErrorMsg = "Проблема с доступом к базе данных"
)

type StartMessagePayload struct {
	UserCart

	Msg string
}

type UserCart struct {
	User User
	Cart Cart
}
