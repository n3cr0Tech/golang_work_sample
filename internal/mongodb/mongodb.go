package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	types "golang_work_sample/internal/types"
	"golang_work_sample/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient(url string, user string, pwd string, dbName string) (*MongoClient, error) {
	var err error
	fullURL := url
	if user != "" && pwd != "" {
		// Insert credentials into URI if they are not already there
		// This is a simple insertion assuming the URI starts with mongodb://
		if len(url) > 10 && url[:10] == "mongodb://" {
			fullURL = fmt.Sprintf("mongodb://%s:%s@%s", user, pwd, url[10:])
		}
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(fullURL))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MongoDB: %w", err)
	}

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	return &MongoClient{
		client: mongoClient,
		db:     mongoClient.Database(dbName),
	}, nil
}

func (m *MongoClient) GetRecord(collectionName string, v interface{}) (*types.User, error) {
	fmt.Printf("-- Getting record: %v from collection: %s\n", v, collectionName)
	data, _ := bson.Marshal(v)
	coll := m.db.Collection(collectionName)
	var result bson.M
	err := coll.FindOne(context.TODO(), data).
		Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the name %s\n", v)
		return nil, err
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return nil, err
	}

	var userRecord types.User
	if err := json.Unmarshal(jsonData, &userRecord); err != nil {
		return nil, err
	}
	return &userRecord, nil
}

func (m *MongoClient) CreateCollections(collectionsName string) {
	command := bson.D{{Key: "create", Value: collectionsName}}
	var result bson.M
	if err := m.db.RunCommand(context.TODO(), command).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Created collections: ", collectionsName)
}

func (m *MongoClient) EnsureRegisterUser(collectionName string, uniqueKey interface{}, data interface{}) bool {
	res := true
	// Extract username from uniqueKey if possible, or expect uniqueKey to be the filter
	// For now, let's assume we can get the username from the uniqueKey map or similar
	// But the request is to refactor UserExists to take a username string.
	// Let's look at how EnsureRegisterUser is called.

	// In auth.go: filter := map[string]interface{}{"username": username}
	// return mongoDB.EnsureRegisterUser("users", filter, newUser)

	username := ""
	if m, ok := uniqueKey.(map[string]interface{}); ok {
		if u, ok := m["username"].(string); ok {
			username = u
		}
	} else if m, ok := uniqueKey.(map[string]string); ok {
		if u, ok := m["username"]; ok {
			username = u
		}
	}

	if m.UserExists(username) {
		return false
	}
	res = m.UpsertRecord(collectionName, uniqueKey, data)
	return res
}

func (m *MongoClient) UserExists(username string) bool {
	recordIndex := map[string]string{"username": username}
	userRecord, _ := m.GetRecord(utils.EnvEntries["MONGO_USERS_DB"], recordIndex)
	res := true
	if userRecord == nil {
		res = false
	} else {
		fmt.Println("Record already exists for ", username)
	}
	return res
}

func (m *MongoClient) UpsertRecord(collectionName string, uniqueKey interface{}, data interface{}) bool {
	bsonData := convertInterfaceToBsonMap(data)
	upsertIsSuccess := true

	coll := m.db.Collection(collectionName)

	update := bson.M{
		"$set": bsonData,
	}
	filter, _ := bson.Marshal(uniqueKey)
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	res, err := coll.UpdateOne(context.TODO(), filter, update, &opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upserted Succesfully: %+v\n", res)
	return upsertIsSuccess
}

func (m *MongoClient) DeleteRecord(collectionName string, indexKey string, indexVal interface{}) {
	filter := bson.D{{Key: indexKey, Value: bson.D{{Key: "$eq", Value: indexVal}}}}
	coll := m.db.Collection(collectionName)
	_, err := coll.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully deleted record: ", indexKey, ": ", indexVal)
}

func convertInterfaceToBsonMap(goMap interface{}) bson.M {
	bsonMap, ok := goMap.(map[string]interface{})
	if !ok {
		panic("FAILED to convert goMap: ")
	}
	bsonData := bson.M(bsonMap)
	return bsonData
}
