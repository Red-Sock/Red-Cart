package telegram

import (
	"context"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/welcome"

	"github.com/Red-Sock/go_tg/client"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/version"
)

type Server struct {
	bot *client.Bot
}

func NewServer(cfg *config.Config, bot *client.Bot, usrSrv service.UserService) (s *Server) {
	s = &Server{
		bot: bot,
	}

	{
		// Add handlers here
		s.bot.AddCommandHandler(version.New(cfg))
		s.bot.AddCommandHandler(welcome.New(usrSrv))

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
