package telegram

import (
	"context"

	"github.com/Red-Sock/go_tg"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/delete_cart"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/delete_item"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/edit/rename"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/purge_cart"

	"github.com/Red-Sock/Red-Cart/internal/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/edit"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/edit/increment"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/start"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/version"
)

type Server struct {
	bot go_tg.TgApi
}

func NewServer(cfg *config.Config, bot go_tg.TgApi, srv service.Storage) (s *Server) {
	s = &Server{
		bot: bot,
	}

	{
		// Add handlers here
		s.bot.AddCommandHandler(version.New(cfg))
		s.bot.AddCommandHandler(start.New(srv.User(), srv.Cart()))

		s.bot.AddCommandHandler(cart.New(srv.User(), srv.Cart()))

		s.bot.AddCommandHandler(edit.New(srv.User(), srv.Cart()))

		s.bot.AddCommandHandler(rename.New(srv.User(), srv.Cart()))
		s.bot.AddCommandHandler(increment.New())
		s.bot.AddCommandHandler(delete_cart.New(srv.Cart()))
		s.bot.AddCommandHandler(delete_item.New(srv.Item(), srv.Cart()))

		s.bot.AddCommandHandler(purge_cart.New(srv.Cart()))

		s.bot.SetDefaultCommandHandler(handlers.NewDefaultCommandHandler(srv.User(), srv.Cart(), srv.Item()))
	}

	return s
}

func (s *Server) Start(_ context.Context) error {
	s.bot.Start()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.bot.Stop()
	return nil
}
