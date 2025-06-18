package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Deba00407/basic-web-server/database"
	schemamodels "github.com/Deba00407/basic-web-server/schema-models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type usersData struct {
	Name     string
	Email    string
	Username string
}

func RegisterUser(user schemamodels.User) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the user already exists by decoding result
	var existing schemamodels.User
	err := database.Collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"username": user.Username},
		},
	}).Decode(&existing)

	if err != nil && err != mongo.ErrNoDocuments {
		// Some DB error happened
		return primitive.NilObjectID, fmt.Errorf("DB error: %w", err)
	}
	if err == nil {
		// User found
		return primitive.NilObjectID, errors.New("user already exists")
	}

	// If we reached here, user does not exist — proceed to insert
	newUser, err := database.Collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("DB insert failed: %w", err)
	}

	return newUser.InsertedID, nil
}

func GetAllRegisteredUsers() ([]usersData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var registeredUsers []usersData

	// Get a cursor pointing to all matching docs
	cursor, err := database.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("users not found")
	}
	defer cursor.Close(ctx)

	// Decode each document from cursor
	for cursor.Next(ctx) {
		var user usersData
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}

		registeredUsers = append(registeredUsers, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return registeredUsers, nil

}
