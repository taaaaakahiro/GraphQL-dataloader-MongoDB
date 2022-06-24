package loader

import "github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/infrastracture/persistence"

type Loaders struct {
	UserLoader *UserLoader
}

func NewLoader(r *persistence.Repositories) *Loaders {
	return &Loaders{
		UserLoader: NewUserLoader(r),
	}
}
