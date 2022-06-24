package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/infrastracture/persistence"
)

type UserLoader struct {
	Loader *dataloader.Loader
}

func NewUserLoader(r *persistence.Repositories) *UserLoader {
	return &UserLoader{
		Loader: dataloader.NewBatchedLoader(r.User.GetUsers),
	}
}

func (l *Loaders) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return nil, nil
}
