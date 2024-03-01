package item

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type Service struct {
	itemRepo domain.ItemRepo
}

func New(itemRepo domain.ItemRepo) *Service {
	return &Service{
		itemRepo: itemRepo,
	}
}

func (s *Service) UpdateName(ctx context.Context, cartID int64, oldName, newName string) error {
	return s.itemRepo.ChangeName(ctx, cartID, oldName, newName)
}
