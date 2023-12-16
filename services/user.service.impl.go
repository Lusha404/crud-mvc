package services

import (
	"context"
	"crud-gin/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userServicesImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(ctx context.Context, userCollection *mongo.Collection) UserService {
	return &userServicesImpl{
		ctx:            ctx,
		userCollection: userCollection,
	}
}

func (u *userServicesImpl) CreateUser(user *models.User) error {
	_, err := u.userCollection.InsertOne(u.ctx, user)
	return err
}

func (u *userServicesImpl) GetUser(name string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "user_name", Value: name}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err

}

func (u *userServicesImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.userCollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(u.ctx)
	if len(users) == 0 {
		return nil, fmt.Errorf("documents not found")
	}
	return users, nil
}

func (u *userServicesImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "user_name", Value: user.Name}}
	update := bson.D{bson.E{Key: "$set",
		Value: bson.D{bson.E{Key: "user_name",
			Value: user.Name}, bson.E{Key: "user_age",
			Value: user.Age}, bson.E{Key: "user_address",
			Value: user.Address}}}}
	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return fmt.Errorf("no matched document found for update\n")
	}
	return nil
}

func (u *userServicesImpl) DeleteUser(name string) error {
	filter := bson.D{bson.E{Key: "user_name", Value: name}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return fmt.Errorf("no matched document for delete")
	}
	return nil
}
