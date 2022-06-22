package graph

import "github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/infrastracture/persistence"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo *persistence.Repositories
}
