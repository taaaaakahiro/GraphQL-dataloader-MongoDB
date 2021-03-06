package v1

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/infrastracture/persistence"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	repo   *persistence.Repositories
	query  *handler.Server
}

func NewHandler(logger *zap.Logger, repo *persistence.Repositories, query *handler.Server) *Handler {
	return &Handler{
		logger: logger,
		repo:   repo,
		query:  query,
	}
}
