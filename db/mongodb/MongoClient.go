package mongodb

import (
	"battleship/config"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	BattleshipDb   = "battleship"
	CollectionGame = "game"
	CollectionUser = "user"
)

var (
	DB Client
)

type Client struct {
	Client *mongo.Client
	ctx    context.Context
}

func CreateMongoClient() (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.C.MongoDB.URL))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	DB = Client{
		Client: client,
		ctx:    ctx,
	}
	return &DB, err
}

func (r *Client) Close() {
	_ = r.Client.Disconnect(r.ctx)
}
