package graph

//go:generate go run github.com/99designs/gqlgen generate

import "github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/infrastracture/persistence"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo *persistence.Repositories
}
