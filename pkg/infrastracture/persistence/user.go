package persistence

import (
	"context"

	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "user"
)

type UserRepo struct {
	col *mongo.Collection
}

var _ repository.IUserRepository = (*UserRepo)(nil)

func NewUserRepository(db *mongo.Database) *UserRepo {
	return &UserRepo{
		col: db.Collection(userCollection),
	}
}

func (r *UserRepo) ListUsers(ctx context.Context) ([]entity.User, error) {

	return []entity.User{}, nil
}

func (r *UserRepo) User(ctx context.Context, useId int) (entity.User, error) {

	return entity.User{}, nil
}
