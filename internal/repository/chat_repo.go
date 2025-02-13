package repository

import (
	"context"

	"github.com/umardev500/common/database"
	"github.com/umardev500/gochat/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository interface {
	FindChats(ctx context.Context, jid, csid string) ([]domain.Chat, error)
}

type chatRepository struct {
	mongDB *database.Mongo
}

func NewChatRepository(db *database.Mongo) ChatRepository {
	return &chatRepository{
		mongDB: db,
	}
}

func (r *chatRepository) FindChats(ctx context.Context, jid, csid string) ([]domain.Chat, error) {
	coll := r.mongDB.Db.Collection("messages")

	filter := bson.D{
		{Key: "csid", Value: csid},
	}

	if jid != "" {
		filter = append(filter, bson.E{Key: "jid", Value: jid})
	}

	aggPipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: filter},
		},
		bson.D{
			{Key: "$unwind", Value: "$messages"},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "jid", Value: "$jid"},
					{Key: "csid", Value: "$csid"},
				}},
				{Key: "status", Value: bson.D{
					{Key: "$last", Value: "$status"},
				}},
				{Key: "unread", Value: bson.D{
					{Key: "$last", Value: "$unread"},
				}},
				{Key: "message", Value: bson.D{
					{Key: "$last", Value: "$messages"},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "jid", Value: "$_id.jid"},
				{Key: "csid", Value: "$_id.csid"},
				{Key: "status", Value: 1},
				{Key: "unread", Value: 1},
				{Key: "message", Value: 1},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "message.timestamp", Value: -1},
			}},
		},
	}

	cur, err := coll.Aggregate(ctx, aggPipeline)
	if err != nil {
		return nil, err
	}

	var chats = make([]domain.Chat, 0)
	err = cur.All(ctx, &chats)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
