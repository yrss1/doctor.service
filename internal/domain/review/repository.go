package review

import "context"

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
	Add(ctx context.Context, doctor Entity) (id string, err error)
	Get(ctx context.Context, id string) (dest Entity, err error)
	Delete(ctx context.Context, id string) (err error)
}
