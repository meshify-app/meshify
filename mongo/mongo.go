package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

///
/// Mongo DB primitives
///

// Serialize write interface to disk
func Serialize(id string, col string, c interface{}) error {
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

	collection := client.Database("meshify").Collection(col)

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
func Deserialize(id string, col string, t reflect.Type) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection(col)

	findstr := fmt.Sprintf("{\"id\":\"%s\"}", id)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	switch t.String() {
	case "model.Host":
		var c *model.Host
		err = collection.FindOne(ctx, filter).Decode(&c)
		return c, nil

	case "model.User":
		var c *model.User
		err = collection.FindOne(ctx, filter).Decode(&c)
		return c, nil

	case "model.Mesh":
		var c *model.Mesh
		err = collection.FindOne(ctx, filter).Decode(&c)
		return c, nil
	}

	log.Infof("reflect.TypeOf(t) = %v", t.String())

	return nil, nil
}

// DeleteHost removes the given id from the given collection
func DeleteHost(id string, col string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection(col)

	findstr := fmt.Sprintf("{\"id\":\"%s\"}", id)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	collection.FindOneAndDelete(ctx, filter)

	return nil
}

// Delete removes the given id from the given collection
func Delete(id string, ident string, col string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection(col)

	findstr := fmt.Sprintf("{\"%s\":\"%s\"}", ident, id)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	collection.FindOneAndDelete(ctx, filter)

	return nil
}

// ReadAllHosts from MongoDB
func ReadAllHosts() []*model.Host {
	hosts := make([]*model.Host, 0)

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
			var host *model.Host
			err = cursor.Decode(&host)
			if err == nil {
				hosts = append(hosts, host)
			}
		}

	}

	return hosts

}

// ReadAllMeshes from MongoDB
func ReadAllMeshes() []*model.Mesh {
	meshes := make([]*model.Mesh, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection("mesh")

	cursor, err := collection.Find(ctx, bson.D{})

	if err == nil {

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var mesh *model.Mesh
			err = cursor.Decode(&mesh)
			if err == nil {
				meshes = append(meshes, mesh)
			}
		}

	}

	return meshes

}

// ReadAllUsers from MongoDB
func ReadAllUsers() []*model.User {
	users := make([]*model.User, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection("users")

	cursor, err := collection.Find(ctx, bson.D{})

	if err == nil {

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user *model.User
			err = cursor.Decode(&user)
			if err == nil {
				users = append(users, user)
			}
		}

	}

	return users

}
