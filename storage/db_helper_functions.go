package storage

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var emptyFilter = bson.D{{}}

func updateManyToDB[T interface{}](storageHandler StorageHelper, collection string, data T) {
	createCollectionIfMissing(storageHandler, collection)

	coll := storageHandler.db.Collection(collection)

	log.Printf("update data: %s", collection)
	log.Println(data)

	result, err := coll.ReplaceOne(context.TODO(), emptyFilter, data, options.Replace().SetUpsert(true))
	if err != nil {
		panic(err)
	}
	log.Printf("save %s updated: %d", collection, result)
}

func getOneFromDB[T interface{}](storageHandler StorageHelper, collection string) T {
	createCollectionIfMissing(storageHandler, collection)
	coll := storageHandler.db.Collection(collection)
	result := coll.FindOne(context.TODO(), emptyFilter)

	var value T

	if e := result.Decode(&value); e != nil {
		logger.Printf("decode error collection %s\n", collection)
		if !errors.Is(e, mongo.ErrNoDocuments) {
			log.Panicln(e)
		}
	}
	logger.Println("return value:")
	logger.Println(value)
	return value
}

func createCollectionIfMissing(storageHandler StorageHelper, collection string) {
	collections, err := storageHandler.db.ListCollectionNames(context.TODO(), emptyFilter)
	if err != nil {
		panic(err)
	}
	exists := false

	for _, coll := range collections {
		if coll != collection {
			continue
		}
		exists = true
		break
	}

	if !exists {
		err = storageHandler.db.CreateCollection(context.TODO(), collection)
		if err != nil {
			panic(err)
		}
	}

}
