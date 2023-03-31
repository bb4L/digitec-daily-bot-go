// A bot to notify you about the daily offer on [Digitec].
//
// It is in no way associated with Digitec Galaxus AG.
//
// Configuration:
//   - set the database connection string by setting the environment variable "MongoDB" or passing it with the flag "-db"
//
// [Digitec]: https://digitec.ch
package main

import (
	"flag"
	"log"
	"os"

	"github.com/bb4L/digitec-daily-bot-go/storage"
	"github.com/bb4L/digitec-daily-bot-go/task"
	"github.com/bb4L/digitec-daily-bot-go/telegram"
)

var logger = log.New(os.Stdout, "[main] ", log.Ldate|log.Ltime|log.Lmsgprefix)

func main() {
	logger.Println("starting digitec-daily-bot")
	offerChannel := make(chan string)
	mongoDbConnectionString := os.Getenv("MongoDB")
	var mongoDBFlag string
	flag.StringVar(&mongoDBFlag, "db", "none", "Specify the mongoDB connection string, it overrides the envirnoment variabel \"MongoDB\".")

	if mongoDbConnectionString == "" && mongoDBFlag == "none" {
		logger.Panicln("mongo db connection string not set please set env variable \"MongoDB\"")
	}

	if mongoDbConnectionString == "" {
		mongoDbConnectionString = mongoDBFlag
	}

	storage := new(storage.StorageHelper)
	storage.Connect("digitecbot", mongoDbConnectionString)

	task.SetupValues(storage)

	go task.StartTask(offerChannel)

	telegram.RunBot(storage.GetSettings().BotToken, offerChannel, storage)
}
