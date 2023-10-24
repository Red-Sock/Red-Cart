package add

import (
	"strconv"
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/add_item"

type Handler struct {
	cartService service.CartService
}

func New(service service.CartService) *Handler {
	return &Handler{
		cartService: service,
	}
}

func (h *Handler) GetDescription() string {
	return "add Item"
}

func (h *Handler) GetCommand() string {
	return Command
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	commandFromTg := strings.Split(in.Text, " ")
	var outMsg string = "testMsg"
	//TODO не знаю на сколько это правильно, просто пытался избавиться от дублирования кода на перевода из строки в число
	var id int64
	var ok bool

	switch len(commandFromTg) {
	case 1, 2:
		outMsg = "Чтобы добавить товар в корзину воспользуйтесь коммандой /add_item {id} {товар_1} {товар_2}\n" +
			"Пример: /add_item 2 беляши кола сникерс"
		out.SendMessage(response.NewMessage(outMsg))
		return
	default:
		if len(commandFromTg) > 2 {
			id, ok = checkIdForInteger(commandFromTg[1])
			if !ok {
				outMsg = "Идентификатор корзины должен быть целочисленным и положительным"
				out.SendMessage(response.NewMessage(outMsg))
				return
			}
		}

	}
	_, err := h.cartService.GetByCartId(in.Ctx, id)
	if err != nil {
		out.SendMessage(response.NewMessage(outMsg))
		return
	}

	out.SendMessage(response.NewMessage(outMsg))
}

func checkIdForInteger(id string) (int64, bool) {
	newId, err := strconv.Atoi(id)

	if err != nil {
		return 0, false
	}

	return int64(newId), true
}
