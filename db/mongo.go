package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"chatapp/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client          *mongo.Client
	RoomsCollection *mongo.Collection
	MsgsCollection  *mongo.Collection
)

func InitMongo() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Ping database
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Warning: Could not ping MongoDB at %s: %v", uri, err)
	} else {
		log.Printf("Successfully connected to MongoDB at %s", uri)
	}

	Client = client
	database := client.Database("chatapp")
	RoomsCollection = database.Collection("rooms")
	MsgsCollection = database.Collection("messages")

	_, _ = RoomsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "rname", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	_, _ = MsgsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "room", Value: 1},
			{Key: "created_at", Value: -1},
		},
	})

	return nil
}

func CreateRoom(room model.Room) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if room.CreatedAt.IsZero() {
		room.CreatedAt = time.Now()
	}

	_, err := RoomsCollection.InsertOne(ctx, room)
	return err
}

func GetRoomByName(name string) (*model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var room model.Room
	err := RoomsCollection.FindOne(ctx, bson.M{"rname": name}).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func SaveMessage(msg model.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	_, err := MsgsCollection.InsertOne(ctx, msg)
	return err
}

func GetRoomHistory(roomName string, limit int64) ([]model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
	cursor, err := MsgsCollection.Find(ctx, bson.M{"room": roomName}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var msgs []model.Message
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}

	// Reverse to return chronologically (oldest to newest)
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nil
}
