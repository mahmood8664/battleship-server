package dao

import (
	"battleship/db/mongodb"
	"battleship/dto"
	"battleship/model"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDao interface {
	Insert(user model.User) (id string, err error)
	GetOne(id string) (user model.User, err error)
}

type UserDaoImpl struct {
}

func NewUserDaoImpl() *UserDaoImpl {
	return &UserDaoImpl{}
}

func (r *UserDaoImpl) GetOne(id string) (user model.User, err error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warn().Str("userId", id).Err(err).Msg("cannot convert to ObjectId")
		return user, dto.ParseError(err)
	}

	result := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionUser).FindOne(context.TODO(), model.User{Id: hex})

	err = result.Decode(&user)
	if err != nil {
		log.Warn().Str("userId", id).Err(err).Msg("cannot decode user")
	}
	return user, dto.ParseError(err)
}

func (r *UserDaoImpl) Insert(user model.User) (id string, err error) {
	one, err := mongodb.DB.Client.Database(mongodb.BattleshipDb).Collection(mongodb.CollectionUser).InsertOne(context.TODO(), user)
	if err != nil {
		log.Warn().Str("userId", id).Err(err).Msg("cannot insert user")
		return "", dto.ParseError(err)
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), nil
}
