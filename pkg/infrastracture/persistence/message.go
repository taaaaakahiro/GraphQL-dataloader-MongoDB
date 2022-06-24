package persistence

import (
	"context"
	"fmt"

	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	messageCollection = "message"
)

type MessageRepo struct {
	col *mongo.Collection
}

var _ repository.IMessageRepository = (*MessageRepo)(nil)

func NewMessageRepository(db *mongo.Database) *MessageRepo {
	return &MessageRepo{
		col: db.Collection(messageCollection),
	}
}

func (r MessageRepo) ListMessages(ctx context.Context, userId int) ([]entity.Message, error) {
	messages := make([]entity.Message, 0)
	srt := bson.D{
		primitive.E{Key: "id", Value: -1},
	}
	opt := options.Find().SetSort(srt)
	flt := bson.D{
		primitive.E{Key: "userId", Value: userId},
	}
	cur, err := r.col.Find(ctx, flt, opt)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var message entity.Message
		err := cur.Decode(&message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r MessageRepo) CreateMessage(ctx context.Context, message *entity.Message) error {
	maxEntity := entity.Message{}
	srt := bson.D{
		primitive.E{Key: "id", Value: -1},
	}
	opt := options.FindOne().SetSort(srt)
	err := r.col.FindOne(ctx, bson.D{}, opt).Decode(&maxEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			message.Id = 1
		} else {
			return err
		}
	} else {
		message.Id = maxEntity.Id + 1
	}
	_, err = r.col.InsertOne(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (r *MessageRepo) UpdateMessage(ctx context.Context, message *entity.Message) error {
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"id": message.Id},
		bson.D{
			{"$set", bson.D{{"message", message.Message}}},
		},
	)
	if err != nil {
		return err
	}
	fmt.Println("update success")
	return nil
}

func (r *MessageRepo) DeleteMessage(ctx context.Context, messageId int) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"id": messageId})
	if err != nil {
		return err
	}
	return nil
}
