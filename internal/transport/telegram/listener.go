package telegram

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/tgapi"
	"github.com/Red-Sock/Red-Cart/internal/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/start"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/version"
)

type Server struct {
	bot tgapi.TgApi
}

func NewServer(cfg *config.Config, bot tgapi.TgApi, srv service.Storage) (s *Server) {
	s = &Server{
		bot: bot,
	}

	{
		// Add handlers here
		s.bot.AddCommandHandler(version.New(cfg))
		s.bot.AddCommandHandler(start.New(srv.User(), srv.Cart()))

		//s.bot.AddCommandHandler(add.New(srv.Cart()))
		//s.bot.AddCommandHandler(create.New(srv.Cart()))
		//s.bot.AddCommandHandler(checkout.New(srv.Cart()))
		//s.bot.AddCommandHandler(accept.New(srv.Cart()))
		//s.bot.AddCommandHandler(decline.New(srv.Cart()))

		s.bot.SetDefaultCommandHandler(handlers.NewDefaultCommandHandler(srv.User(), srv.Cart()))
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
