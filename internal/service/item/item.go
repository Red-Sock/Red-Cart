package item

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"

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
	err := s.itemRepo.ChangeName(ctx, cartID, oldName, newName)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, cartId int64, itemName string) error {
	err := s.itemRepo.Delete(ctx, cartId, itemName)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (s *Service) Check(ctx context.Context, cartId int64, itemName string) error {
	err := s.itemRepo.Check(ctx, cartId, itemName)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (s *Service) Uncheck(ctx context.Context, cartId int64, itemName string) error {
	err := s.itemRepo.Uncheck(ctx, cartId, itemName)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
