package telegram

import (
	"context"

	"github.com/Red-Sock/go_tg"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/default_handler"
)

type Server struct {
	bot go_tg.TgApi
}

func NewServer(_ config.Config, bot go_tg.TgApi, srv service.Service) (s *Server) {
	s = &Server{
		bot: bot,
	}

	{
		hs := newHandlerStore(srv)

		s.bot.AddCommandHandler(hs.handlers[commands.Start])

		s.bot.AddCommandHandler(hs.handlers[commands.Cart])

		s.bot.AddCommandHandler(hs.handlers[commands.EditItem])
		s.bot.AddCommandHandler(hs.handlers[commands.RenameItem])
		s.bot.AddCommandHandler(hs.handlers[commands.IncrementItemCount])
		s.bot.AddCommandHandler(hs.handlers[commands.CheckItem])
		s.bot.AddCommandHandler(hs.handlers[commands.UncheckItem])

		s.bot.AddCommandHandler(hs.handlers[commands.ClearMenu])
		s.bot.AddCommandHandler(hs.handlers[commands.DeleteItem])

		s.bot.AddCommandHandler(hs.handlers[commands.CartSetting])

		s.bot.AddCommandHandler(hs.handlers[commands.Purge])

		s.bot.SetDefaultCommandHandler(default_handler.NewDefaultCommandHandler(srv, hs.handlers))
	}

	return s
}

func (s *Server) Start(_ context.Context) error {
	err := s.bot.Start()
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.bot.Stop()

	return nil
}
