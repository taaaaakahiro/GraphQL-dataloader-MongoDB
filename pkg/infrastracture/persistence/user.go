package persistence

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	users := make([]entity.User, 0)
	srt := bson.D{
		primitive.E{Key: "id", Value: -1},
	}
	opt := options.Find().SetSort(srt)
	cur, err := r.col.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var user entity.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) User(ctx context.Context, userId int) (entity.User, error) {
	user := entity.User{}
	flt := bson.D{
		primitive.E{Key: "id", Value: userId},
	}
	opt := options.FindOne()
	err := r.col.FindOne(ctx, flt, opt).Decode(&user)
	if err == mongo.ErrNoDocuments {
		log.Println("Documents not found")
	}
	return user, err
}

func (r *UserRepo) CreateUser(ctx context.Context, user *entity.User) error {
	maxEntity := entity.User{}
	srt := bson.D{
		primitive.E{Key: "id", Value: -1},
	}
	opt := options.FindOne().SetSort(srt)
	err := r.col.FindOne(ctx, bson.D{}, opt).Decode(&maxEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.Id = 1
		} else {
			return err
		}
	} else {
		user.Id = maxEntity.Id + 1
	}
	_, err = r.col.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// for dataloader 〜途中〜
func (r *UserRepo) GetUsers(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	output := make([]*dataloader.Result, len(keys))

	userIds := make([]interface{}, len(keys))
	for i, key := range keys {
		userId, err := strconv.Atoi(key.String())
		if err != nil {
			log.Printf("%+v", err)
			err := fmt.Errorf("user error %s", err.Error())
			output[0] = &dataloader.Result{Data: nil, Error: err}
			return output
		}
		userIds[i] = userId
	}
	return nil
}
