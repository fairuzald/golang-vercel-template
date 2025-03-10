package interfaces

import (
	"context"

	"golang-template/app/core/entity"
	"golang-template/infrastructure/logger"
)

type Service[T entity.Entity, CreateDTO any, UpdateDTO any] interface {
	Create(ctx context.Context, dto CreateDTO) (T, error)

	Update(ctx context.Context, id string, dto UpdateDTO) (T, error)

	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (T, error)

	List(ctx context.Context, paginator Paginator) ([]T, int64, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BaseService provides a common implementation for services
type BaseService[T entity.Entity, CreateDTO any, UpdateDTO any] struct {
	Repository Repository[T]
	Logger     logger.Logger
}
