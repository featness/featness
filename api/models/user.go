package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name     string        `json:"name"`
	UserID   string        `json:"userid"`
	JoinedAt time.Time     `json:"joinedat"`
	ImageURL string        `json:"imageurl"`
}

func Users() *mgo.Collection {
	conn, _ := Conn()
	return conn.C("users")
}

func GetOrCreateUser(name string, userID string, imageURL string) (*User, error) {
	usersColl := Users()

	_, err := usersColl.Upsert(
		bson.M{"userid": userID},
		&User{bson.NewObjectId(), name, userID, time.Now(), imageURL},
	)

	user := &User{}
	err = usersColl.Find(bson.M{"userid": userID}).One(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
