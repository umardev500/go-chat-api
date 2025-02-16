package repository

import (
	"context"

	"github.com/umardev500/common/database"
	"github.com/umardev500/gochat/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatRepository interface {
	CheckChatIsExist(ctx context.Context, jid, csid string) error
	FindChats(ctx context.Context, jid, csid string, status *domain.ChatStatus) ([]domain.Chat, error)
	CreateChat(ctx context.Context, jid, csid string, chat interface{}) (bool, error)
	PushMessage(ctx context.Context, jid, csid string, message interface{}) error
	UpdateUnread(ctx context.Context, jid, csid string, value int64) error
}

type chatRepository struct {
	mongDB *database.Mongo
	coll   *mongo.Collection
}

func NewChatRepository(db *database.Mongo) ChatRepository {
	return &chatRepository{
		mongDB: db,
		coll:   db.Db.Collection("messages"),
	}
}

func (r *chatRepository) CheckChatIsExist(ctx context.Context, jid, csid string) error {
	filter := bson.M{"jid": jid, "csid": csid}
	err := r.coll.FindOne(ctx, filter).Err()
	return err
}

func (r *chatRepository) FindChats(ctx context.Context, jid, csid string, status *domain.ChatStatus) ([]domain.Chat, error) {
	coll := r.mongDB.Db.Collection("messages")

	filter := bson.D{
		{Key: "csid", Value: csid},
	}

	if jid != "" {
		filter = append(filter, bson.E{Key: "jid", Value: jid})
	}

	if status != nil {
		filter = append(filter, bson.E{Key: "status", Value: status})
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

func (r *chatRepository) CreateChat(ctx context.Context, jid, csid string, chat interface{}) (bool, error) {
	filter := bson.M{"jid": jid, "csid": csid}

	err := r.coll.FindOne(ctx, filter).Err()
	if err == nil {
		// Chat exists
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		// Chat does not exist, so create it
		_, err = r.coll.InsertOne(ctx, chat)
		return false, err
	}

	// Return error if something went wrong
	return false, err
}

func (r *chatRepository) PushMessage(ctx context.Context, jid, csid string, message interface{}) error {
	filter := bson.D{
		{Key: "jid", Value: jid},
		{Key: "csid", Value: csid},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "messages", Value: message},
		}},
	}

	_, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) UpdateUnread(ctx context.Context, jid, csid string, value int64) error {
	filter := bson.D{
		{Key: "jid", Value: jid},
		{Key: "csid", Value: csid},
	}

	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "unread", Value: value},
		}},
	}

	_, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
