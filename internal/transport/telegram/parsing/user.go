package parsing

import (
	"github.com/Red-Sock/go_tg/model"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

func ToDomainUser(msgIn *model.MessageIn) domain.User {
	return domain.User{
		Id:        msgIn.From.ID,
		UserName:  msgIn.From.UserName,
		FirstName: msgIn.From.FirstName,
		LastName:  msgIn.From.LastName,
	}
}
