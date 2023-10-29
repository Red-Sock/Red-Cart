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
	return "Adds item to cart"
}

func (h *Handler) GetCommand() string {
	return Command
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	commandFromTg := strings.Split(in.Text, " ")
	var outMsg string = "testMsg"

	if len(commandFromTg) < 3 {
		outMsg = "Чтобы добавить товар в корзину воспользуйтесь коммандой /add_item {id} {товар_1} {товар_2}\n" +
			"Пример: /add_item 2 беляши кола сникерс"
		out.SendMessage(response.NewMessage(outMsg))
		return
	}

	id, _ := strconv.Atoi(commandFromTg[1])
	if id == 0 {
		outMsg = "Идентификатор корзины должен быть целочисленным и положительным"
		out.SendMessage(response.NewMessage(outMsg))
		return
	}

	err := h.cartService.AddCartItems(in.Ctx, commandFromTg[2:], int64(id), in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
	}
	outMsg = "Предметы были успешно добавлены в корзину!"

	out.SendMessage(response.NewMessage(outMsg))
}
