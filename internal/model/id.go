package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ID string

func (id ID) String() string {
	return string(id)
}

func ParseId(str string) (ID, error) {
	return ID(str), nil
}

func ParseIdFromObjectId(oid primitive.ObjectID) ID {
	return ID(oid.Hex())
}

func (id ID) ObjectId() primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id.String())
	return oid
}
