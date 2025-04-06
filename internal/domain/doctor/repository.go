package doctor

import "context"

type Repository interface {
	ListWithSchedules(ctx context.Context) (dest []Entity, err error)
	GetWithSchedules(ctx context.Context, id string) (Entity, error)
	Delete(ctx context.Context, id string) (err error)
	SearchWithSchedules(ctx context.Context, filter Entity) ([]Entity, error)
}
