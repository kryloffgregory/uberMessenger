package messages

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	DBName = "messenger"
	CollectionName = "messages"
)


type DAO struct {
	client *mongo.Client
	db *mongo.Database
	collection *mongo.Collection
}

func NewDAO(ctx context.Context, client *mongo.Client) (*DAO, error) {
	db := client.Database(DBName)
	collection:=db.Collection(CollectionName)

	indexOptions := options.Index().SetUnique(false)
	indexKeys := bsonx.MDoc{
		"chatId": bsonx.Int32(1),
		"time":bsonx.Int32(-1),
	}

	noteIndexModel := mongo.IndexModel{
		Options: indexOptions,
		Keys:    indexKeys,
	}

	_, err := collection.Indexes().CreateOne(ctx, noteIndexModel)
	if err != nil {
		return nil, err
	}

	return &DAO{
		client:client,
		db:db,
		collection:collection,
	}, nil
}

func (dao *DAO) GetMessagesByChat(ctx context.Context, chatID primitive.ObjectID, limit int, offset int) ([]*Message, error) {
	filter := bson.D{{"chatId", chatID}}
	options:=options.Find()
	options.SetSort(bson.D{{"time", -1}})
	options.SetLimit(int64(limit))
	options.SetSkip(int64(offset))

	cursor, err := dao.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var result []*Message

	for cursor.Next(ctx) {
		var message *Message
		if err:=cursor.Decode(&message); err!=nil {
			return nil, err
		}
		result = append(result, message)
	}

	return result, nil
}

func (dao *DAO) InitJunk(ctx context.Context) error{
	userID1,err:=primitive.ObjectIDFromHex("5f78829a44202661a33d787a")
	if err!=nil {
		return nil
	}

	userID2,err:=primitive.ObjectIDFromHex("5f78829a44202661a33d787b")
	if err!=nil {
		return nil
	}

	chat,err:=primitive.ObjectIDFromHex("5f788850ddf089e25fa8ada6")
	if err!=nil {
		return nil
	}

	msgs:=[]*Message{
		{
			ID:     primitive.NewObjectID(),
			From:   userID2,
			ChatID: chat,
			Text:   "300",
			Time:   time.Now().Add(-time.Hour),
		},
		{
			ID:     primitive.NewObjectID(),
			From:   userID1,
			ChatID: chat,
			Text:   "Отсоси у тракториста",
			Time:   time.Now(),
		},
	}

	if _, err:= dao.collection.InsertOne(ctx, msgs[0]); err!=nil {
		return err
	}

	if _, err:= dao.collection.InsertOne(ctx, msgs[1]); err!=nil {
		return err
	}

	return nil
}