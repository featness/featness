package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	Id          bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Provider    string        `json:"provider"`
	AccessToken string        `json:"accessToken"`
	Name        string        `json:"name"`
	UserId      string        `json:"userid"`
	JoinedAt    time.Time     `json:"joinedat"`
	ImageUrl    string        `json:"imageurl"`
}

func Users() *mgo.Collection {
	conn, _ := Conn()
	return conn.C("users")
}

func GetOrCreateUser(provider string, accessToken string, userID string, name string, imageURL string) (*User, error) {
	usersColl := Users()

	_, err := usersColl.Upsert(
		bson.M{"userid": userID},
		&User{bson.NewObjectId(), provider, accessToken, name, userID, time.Now(), imageURL},
	)

	user := &User{}
	err = usersColl.Find(bson.M{"userid": userID}).One(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
