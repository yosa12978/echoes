package cache

import (
	"context"

	"github.com/yosa12978/echoes/types"
)

type Link interface {
	GetLinks(ctx context.Context) ([]types.Link, error)
	GetLinkById(ctx context.Context, id string) (*types.Link, error)

	Create(ctx context.Context, link types.Link) error
	Update(ctx context.Context, id string, link types.Link) error
	Delete(ctx context.Context, id string) error
}
