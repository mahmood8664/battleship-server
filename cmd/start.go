package cmd

import (
	"battleship/db/mongodb"
	"battleship/http"
	"battleship/socket"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var startCMD = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		client := connectToMongo()
		defer client.Close()
		go http.StartHttpServer()
		socket.StartSocketServer()
	},
}

func connectToMongo() *mongodb.Client {
	//mongodb
	client, err := mongodb.CreateMongoClient()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	err = client.Client.Ping(context.TODO(), &readpref.ReadPref{})
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	return client
}
