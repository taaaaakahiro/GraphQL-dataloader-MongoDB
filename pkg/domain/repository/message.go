package repository

import (
	"context"

	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/domain/entity"
)

type IMessageRepository interface {
	ListMessages(ctx context.Context, userId int) ([]entity.Message, error)
	CreateMessage(ctx context.Context, message *entity.Message) error
	UpdateMessage(ctx context.Context, message *entity.Message) error
	DeleteMessage(ctx context.Context, messageId int) error
}
