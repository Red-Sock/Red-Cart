package telegram

import (
	"context"

	"github.com/Red-Sock/go_tg/client"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/add"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/create"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/version"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/welcome"
)

type Server struct {
	bot *client.Bot
}

func NewServer(cfg *config.Config, bot *client.Bot, srv service.Storage) (s *Server) {
	s = &Server{
		bot: bot,
	}

	{
		// Add handlers here
		s.bot.AddCommandHandler(version.New(cfg))
		s.bot.AddCommandHandler(welcome.New(srv.User()))

		s.bot.AddCommandHandler(add.New(srv.Cart()))
		s.bot.AddCommandHandler(create.New(srv.Cart()))

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
