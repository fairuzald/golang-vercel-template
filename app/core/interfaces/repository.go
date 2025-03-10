package interfaces

import (
	"context"

	"golang-template/app/core/entity"
)

// Repository defines the base interface for all repositories
type Repository[T entity.Entity] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity T) (T, error)

	// Update updates an existing entity
	Update(ctx context.Context, entity T) (T, error)

	// Delete deletes an entity by ID
	Delete(ctx context.Context, id string) error

	// GetByID gets an entity by ID
	GetByID(ctx context.Context, id string) (T, error)

	// List gets entities with optional pagination and filters
	List(ctx context.Context, page, limit int, filters map[string]interface{}) ([]T, int64, error)

	// Count counts entities with optional filters
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)

	// Exists checks if an entity exists with the given filters
	Exists(ctx context.Context, filters map[string]interface{}) (bool, error)

	// Transaction executes a function within a transaction
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// Paginator interface for pagination functionality
type Paginator interface {
	// GetPage returns the current page number
	GetPage() int

	// GetLimit returns the page size
	GetLimit() int

	// GetOffset returns the offset for the current page
	GetOffset() int

	// GetSort returns the sort field and direction
	GetSort() string

	// GetFilter returns the filter parameters
	GetFilter() map[string]interface{}
}

// NewPaginator creates a new paginator with default values
func NewPaginator(page, limit int, sort string, filter map[string]interface{}) Paginator {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	return &paginator{
		page:   page,
		limit:  limit,
		sort:   sort,
		filter: filter,
	}
}

// paginator implements the Paginator interface
type paginator struct {
	page   int
	limit  int
	sort   string
	filter map[string]interface{}
}

func (p *paginator) GetPage() int {
	return p.page
}

func (p *paginator) GetLimit() int {
	return p.limit
}

func (p *paginator) GetOffset() int {
	return (p.page - 1) * p.limit
}

func (p *paginator) GetSort() string {
	return p.sort
}

func (p *paginator) GetFilter() map[string]interface{} {
	return p.filter
}
