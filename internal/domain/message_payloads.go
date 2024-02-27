package domain

const DbErrorMsg = "Проблема с доступом к базе данных"

type StartMessagePayload struct {
	User User
	Cart Cart

	Msg string
}
