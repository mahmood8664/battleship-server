package dao

import (
	"battleship/db/mongodb"
	"battleship/dto"
	"battleship/model"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GameEventDao interface {
	Insert(event model.GameEvent) (id string, err error)
	FindMany(gameId string) (events []model.GameEvent, err error)
	FindManyByType(gameId string, eventType model.GameEventType) (events []model.GameEvent, err error)
	GetLast(GameId string) (event model.GameEvent, err error)
}

type GameEventDaoImpl struct {
}

func NewEventGameDaoImpl() GameEventDaoImpl {
	return GameEventDaoImpl{}
}

func (r GameEventDaoImpl) Insert(event model.GameEvent) (id string, err error) {
	one, err := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionGameEvent).InsertOne(context.TODO(), event)
	if err != nil {
		log.Warn().Str("event_type", string(event.Type)).Err(err).Msg("cannot insert Game")
		return "", dto.ParseError(err)
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r GameEventDaoImpl) FindMany(gameId string) (events []model.GameEvent, err error) {
	events = []model.GameEvent{}
	hex, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot convert to objectId")
		return events, dto.ParseError(err)
	}
	filter := bson.D{{"game_id", hex}}
	opts := options.Find()
	opts.SetSort(bson.D{{"time", -1}})
	many, err := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionGameEvent).Find(context.TODO(), filter, opts)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot find game events")
	}
	err = many.All(context.TODO(), &events)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot decode Game events")
	}
	return events, dto.ParseError(err)
}

func (r GameEventDaoImpl) FindManyByType(gameId string, eventType model.GameEventType) (events []model.GameEvent, err error) {
	events = []model.GameEvent{}
	hex, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot convert to objectId")
		return events, dto.ParseError(err)
	}
	filter := bson.D{{"game_id", hex}, {"type", eventType}}
	opts := options.Find()
	opts.SetSort(bson.D{{"time", -1}})
	many, err := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionGameEvent).Find(context.TODO(), filter, opts)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot find game events")
	}
	err = many.All(context.TODO(), &events)
	if err != nil {
		log.Warn().Str("gameId", gameId).Err(err).Msg("cannot decode Game events")
	}
	return events, dto.ParseError(err)
}

func (r GameEventDaoImpl) GetLast(GameId string) (event model.GameEvent, err error) {
	return model.GameEvent{}, nil
}
