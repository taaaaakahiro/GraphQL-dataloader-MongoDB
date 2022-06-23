package repository

import (
	"context"

	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/entity"
)

type IUserRepository interface {
	ListUsers(ctx context.Context) ([]entity.User, error)
	User(ctx context.Context, userId int) (entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}
