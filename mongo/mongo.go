package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Serialize write interface to disk
func Serialize(id string, c interface{}) error {
	//b, err := json.MarshalIndent(c, "", "  ")
	//if err != nil {
	//	return err
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	data, err := json.Marshal(c)
	//	json := fmt.Sprintf("%v", user)
	var b interface{}
	err = bson.UnmarshalExtJSON([]byte(data), true, &b)

	collection := client.Database("meshify").Collection("hosts")

	findstr := fmt.Sprintf("{\"id\":\"%s\"}", id)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	update := bson.M{
		"$set": b,
	}

	opts := options.Update().SetUpsert(true)
	res, err := collection.UpdateOne(ctx, filter, update, opts)

	//	if res != nil && res.Err != nil {
	//		collection.InsertOne(ctx, b)
	//	}

	log.Infof("Res: %v", res)
	return err
}

// Deserialize read interface from disk
func Deserialize(id string) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection("hosts")

	findstr := fmt.Sprintf("{\"id\":\"%s\"}", id)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	var c *model.Client
	err = collection.FindOne(ctx, filter).Decode(&c)

	return c, nil
}

// DeleteClient removes the given client id
func DeleteClient(id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection("hosts")

	findstr := fmt.Sprintf("{\"id\":\"%s\"}", id)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	collection.FindOneAndDelete(ctx, filter)

	return nil
}

// ReadAllClients from MongoDB
func ReadAllClients() []*model.Client {
	clients := make([]*model.Client, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection("hosts")

	cursor, err := collection.Find(ctx, bson.D{})

	if err == nil {

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var client *model.Client
			err = cursor.Decode(&client)
			if err == nil {
				clients = append(clients, client)
			}
		}

	}

	return clients

}
