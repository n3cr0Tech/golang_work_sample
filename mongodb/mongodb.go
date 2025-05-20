package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	types "example.com/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var db *mongo.Database

func Init(url string, dbName string, userCollections string){	
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(url))
	if err != nil {
		panic(err)
	}
	db = mongoClient.Database(dbName)

	CreateCollections(userCollections)
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	
}

func GetRecord(collectionName string, v interface{})(*types.User, error){
	data, _ := bson.Marshal(v)
	coll := db.Collection(collectionName)		
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

func CreateCollections(collectionsName string){		
	command := bson.D{{"create", collectionsName}}
	var result bson.M
	if err := db.RunCommand(context.TODO(), command).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Created collections: ", collectionsName)
}

func EnsureRegisterUser(collectionName string, uniqueKey interface{}, data interface{}) bool{
	res := true
	if UserExists(collectionName, uniqueKey, data){
		return false
	}
	res = UpsertRecord(collectionName, uniqueKey, data)
	return res
}

func UserExists(collectionName string, uniqueKey interface{}, data interface{}) bool{
	recordAlreadyExists, recErr := GetRecord(collectionName, uniqueKey)
	res := true
	if recErr != nil{
		panic(recErr)
		res = false
	}
	if recordAlreadyExists == nil{
		res = false
	}else{
		fmt.Println("Record already exists for ", uniqueKey)		
	}
	return res
}

func UpsertRecord(collectionName string, uniqueKey interface{}, data interface{}) bool{
	bsonData := convertInterfaceToBsonMap(data)
	upsertIsSuccess := true

	coll := db.Collection(collectionName)	
	
	update := bson.M{
		"$set": bsonData,
	}
	filter, _ := bson.Marshal(uniqueKey)	
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	res, err := coll.UpdateOne(context.TODO(), filter, update, &opts)
	if err != nil{
		panic(err)
		upsertIsSuccess = false
	}
	fmt.Printf("Upserted Succesfully: %s\n", res)
	return upsertIsSuccess
}

func DeleteRecord(collectionName string, indexKey string, indexVal interface{}){	
	filter := bson.D{{indexKey, bson.D{{"$eq", indexVal}}}}
	coll := db.Collection(collectionName)	
	_, err := coll.DeleteOne(context.TODO(), filter, nil)
	if err != nil{
		panic(err)
	}
	fmt.Println("Successfully deleted record: ", indexKey, ": ", indexVal);
}

func convertInterfaceToBsonMap(goMap interface{} ) bson.M{
	bsonMap, ok := goMap.(map[string]interface{})
	if !ok{
		panic("FAILED to convert goMap: ");		
	}
	bsonData := bson.M(bsonMap)
	return bsonData
}