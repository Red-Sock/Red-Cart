package add

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
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
	commandFromTg := in.Args
	if len(commandFromTg) < 3 {
		out.SendMessage(response.NewMessage("Чтобы добавить товар в корзину воспользуйтесь командой /add_item {id} {товар_1} {товар_2}\n" +
			"Пример: /add_item 2 беляши кола сникерс"))
		return
	}

	id, _ := strconv.Atoi(commandFromTg[0])
	if id == 0 {
		out.SendMessage(response.NewMessage("Идентификатор корзины должен быть целочисленным и положительным"))
		return
	}

	items := make([]domain.Item, len(commandFromTg)-1)
	for _, itemName := range commandFromTg[1:] {
		items = append(items, domain.Item{
			Name:   itemName,
			Amount: 1,
		})
	}
	//err := h.cartService.AddToDefaultCart(in.Ctx, items, int64(id), in.From.ID)
	//if err != nil {
	//	out.SendMessage(response.NewMessage(err.Error()))
	//	return
	//}

	out.SendMessage(response.NewMessage("Предметы были успешно добавлены в корзину!"))
}
