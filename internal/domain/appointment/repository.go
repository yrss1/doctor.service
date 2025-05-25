package appointment

import (
	"context"
)

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
	Add(ctx context.Context, doctor Entity) (id string, err error)
	Get(ctx context.Context, id string) (dest Entity, err error)
	Cancel(ctx context.Context, id string) (err error)
	ListByUserID(ctx context.Context, id string) ([]EntityView, error)
	UpdateMeetingURL(ctx context.Context, id string, meetingURL string) error
}
