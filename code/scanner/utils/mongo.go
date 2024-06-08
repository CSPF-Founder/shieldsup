package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func IsValidObjectId(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
