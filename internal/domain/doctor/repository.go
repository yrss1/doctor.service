package doctor

import "context"

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
	Add(ctx context.Context, doctor Entity) (id string, err error)
}
