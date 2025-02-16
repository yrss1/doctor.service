package doctor

import "context"

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
}
