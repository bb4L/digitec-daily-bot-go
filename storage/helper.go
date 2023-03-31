// Package storage provides an abstraction to interact with the storage of the data used by the bot.
//
// Currently the only storage provider is [MongoDB]
//
// [MongoDB]: https://www.mongodb.com/
package storage

import (
	"context"
	"log"
	"os"

	settings "github.com/bb4L/digitec-daily-bot-go/settings"
	"github.com/bb4L/digitec-daily-bot-go/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StorageHelper is handling the interactions with the storage
type StorageHelper struct {
	connectionString string
	dbName           string
	client           *mongo.Client
	db               *mongo.Database
}

var logger = log.New(os.Stdout, "[storage] ", log.Ldate|log.Ltime|log.Lmsgprefix)

// Connect to the database with the given dbName and connectionString
//
// The created connection will be stored in the [StorageHelper]
func (storage *StorageHelper) Connect(dbName string, connectionString string) {
	storage.dbName = dbName
	storage.connectionString = connectionString
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(storage.connectionString))
	if err != nil {
		logger.Panicln(err)
	}
	storage.client = client
	storage.db = client.Database(storage.dbName)
}

// Disconnect disconnects from the database
func (storage *StorageHelper) Disconnect() {
	logger.Println("disconnect mongodb")
	if err := storage.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

// GetCurrentOfferString returns the stringtemplate for the current offer
func (storage *StorageHelper) GetCurrentOfferString() string {
	return getOneFromDB[types.TaskSettings](*storage, "task").CurrentTextTemplate
}

// GetTaskSettings returns the [types.TaskSettings] for the task
func (storage *StorageHelper) GetTaskSettings() types.TaskSettings {
	return getOneFromDB[types.TaskSettings](*storage, "task")
}

// GetSettings returns the settings for the bot
func (storage *StorageHelper) GetSettings() settings.Settings {
	return getOneFromDB[settings.Settings](*storage, "settings")
}

// AddUser adds a user to the subscribers of the bot
func (storage *StorageHelper) AddUser(newUser int64) {
	settings := storage.GetTaskSettings()
	for _, user := range settings.Users {
		if user == newUser {
			return
		}
	}
	settings.Users = append(settings.Users, newUser)
	updateManyToDB(*storage, "task", settings)
}

// RemoveUser removes a user from the subscribers
func (storage *StorageHelper) RemoveUser(userToRemove int64) {
	settings := storage.GetTaskSettings()
	newUsers := []int64{}
	for _, user := range settings.Users {
		if user == userToRemove {
			continue
		}
		newUsers = append(newUsers, user)
	}
	settings.Users = newUsers
	updateManyToDB(*storage, "task", settings)
}
