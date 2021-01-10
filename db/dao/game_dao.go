package dao

import (
	"battleship/db/mongodb"
	"battleship/dto"
	"battleship/model"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameDao interface {
	Insert(game model.Game) (id string, err error)
	GetOne(gameId string) (game model.Game, err error)
	Update(game model.Game) error
}

type GameDaoImpl struct {
}

func NewGameDaoImpl() GameDaoImpl {
	return GameDaoImpl{}
}

func (r GameDaoImpl) Insert(game model.Game) (id string, err error) {
	one, err := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionGame).InsertOne(context.TODO(), game)
	if err != nil {
		log.Warn().Str("gameId", game.Id.Hex()).Err(err).Msg("cannot insert Game")
		return "", dto.ParseError(err)
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r GameDaoImpl) Update(game model.Game) error {
	res, err := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionGame).
		UpdateOne(context.TODO(), bson.M{"_id": game.Id}, bson.D{{"$set", game}})
	if err != nil {
		log.Warn().Str("gameId", game.Id.Hex()).Err(err).Msg("cannot update Game")
		return dto.ParseError(err)
	}
	log.Debug().Int64("matched_count", res.MatchedCount).Int64("modified_count", res.ModifiedCount).
		Str("game_id", game.Id.Hex()).Msg("")
	return nil
}

func (r GameDaoImpl) GetOne(gameId string) (game model.Game, err error) {
	hex, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot convert to objectId")
		return game, dto.ParseError(err)
	}
	filter := bson.D{{"_id", hex}}
	one := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionGame).FindOne(context.TODO(), filter)
	err = one.Decode(&game)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot decode Game")
	}
	return game, dto.ParseError(err)
}
