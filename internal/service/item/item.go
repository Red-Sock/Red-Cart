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

func (s *Service) Delete(ctx context.Context, cartId int64, itemName string) error {
	return s.itemRepo.Delete(ctx, cartId, itemName)
}

func (s *Service) Check(ctx context.Context, cartId int64, itemName string) error {
	return s.itemRepo.Check(ctx, cartId, itemName)
}

func (s *Service) Uncheck(ctx context.Context, cartId int64, itemName string) error {
	return s.itemRepo.Uncheck(ctx, cartId, itemName)
}
