package helpers

import (
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	"github.com/sirupsen/logrus"
)

func DeleteIncomingMessage(msgIn *model.MessageIn, out interfaces.Chat) {
	if msgIn.IsCallback {
		return
	}

	err := out.SendMessage(&response.DeleteMessage{
		ChatId:    msgIn.Chat.ID,
		MessageId: int64(msgIn.MessageID),
	})
	if err != nil {
		logrus.Errorf("error deleting incoming message: %s", err)
	}
}
