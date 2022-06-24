package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/graph/model"
)

func (r *messageResolver) User(ctx context.Context, obj *model.Message) (*model.User, error) {
	// case by no dataloader
	userId, err := strconv.Atoi(obj.UserID)
	if err != nil {
		return nil, err
	}
	entityUser, err := r.Repo.User.User(ctx, userId)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:   strconv.Itoa(entityUser.Id),
		Name: entityUser.Name,
	}
	return user, nil
	// case by use dataloader
}

func (r *mutationResolver) CreateMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {
	userId, err := strconv.Atoi(input.UserID)
	if err != nil {
		return nil, err
	}
	_, err = r.Repo.User.User(ctx, userId)
	if err != nil {
		// not exist etc...
		return nil, errors.New("user error. " + err.Error())
	}

	entityMessage := &entity.Message{
		UserId:  userId,
		Message: input.Message,
	}
	err = r.Repo.Message.CreateMessage(ctx, entityMessage)
	if err != nil {
		return nil, err
	}
	result := &model.Message{
		Message: input.Message,
		ID:      strconv.Itoa(entityMessage.Id),
		UserID:  input.UserID,
	}
	return result, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	entityUser := &entity.User{
		Name: input.UserName,
	}
	err := r.Repo.User.CreateUser(ctx, entityUser)
	if err != nil {
		return nil, err
	}
	result := &model.User{
		ID:   strconv.Itoa(entityUser.Id),
		Name: input.UserName,
	}
	return result, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, input model.DeleteMessage) (*model.Message, error) {
	messageId, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}
	// entityMessage := &entity.Message{
	// 	Id: messageId,
	// }
	err = r.Repo.Message.DeleteMessage(ctx, messageId)
	if err != nil {
		return nil, err
	}
	result := &model.Message{
		ID: input.ID,
	}

	return result, nil
}

func (r *mutationResolver) UpdateMessage(ctx context.Context, input *model.UpdateMessage) (*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	entities, err := r.Repo.User.ListUsers(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	users := make([]*model.User, 0)
	for _, u := range entities {
		users = append(users, &model.User{
			ID:   strconv.Itoa(u.Id),
			Name: u.Name,
		})
	}
	return users, nil
}

func (r *queryResolver) Messages(ctx context.Context, userID string) ([]*model.Message, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}
	entities, err := r.Repo.Message.ListMessages(ctx, id)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	for _, entity := range entities {
		messages = append(messages, &model.Message{
			ID:      strconv.Itoa(entity.Id),
			Message: entity.Message,
			UserID:  userID,
		})
	}
	return messages, nil
}

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type messageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
