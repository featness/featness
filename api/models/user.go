package models

import (
	"fmt"
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

func Users() (*mgo.Session, *mgo.Collection, error) {
	conn, _, err := Conn()

	if err != nil {
		return nil, nil, err
	}

	return UsersWithConn(conn)
}

func UsersWithConn(conn *mgo.Session) (*mgo.Session, *mgo.Collection, error) {
	return conn, conn.DB("featness").C("users"), nil
}

func GetOrCreateUser(provider string, accessToken string, userID string, name string, imageURL string) (*User, error) {
	conn, usersColl, err := Users()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if provider == "" {
		return nil, fmt.Errorf("Can't create user with empty provider.")
	}

	if accessToken == "" {
		return nil, fmt.Errorf("Can't create user with empty accessToken.")
	}

	if userID == "" {
		return nil, fmt.Errorf("Can't create user with empty userID.")
	}

	if name == "" {
		return nil, fmt.Errorf("Can't create user with empty name.")
	}

	user, err := GetUserByUserId(userID)
	if err != nil {
		user = &User{bson.NewObjectId(), provider, accessToken, name, userID, time.Now(), imageURL}
		err = usersColl.Insert(user)

		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func FindUsersWithIdLike(name string) (*[]User, error) {
	conn, usersColl, err := Users()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	users := []User{}

	findByUserId := bson.M{"userid": bson.M{"$regex": bson.RegEx{name, ""}}}
	findByUserName := bson.M{"name": bson.M{"$regex": bson.RegEx{name, ""}}}
	query := bson.M{"$or": []bson.M{findByUserId, findByUserName}}

	err = usersColl.Find(query).All(&users)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func GetUserByUserId(userId string) (*User, error) {
	conn, usersColl, err := Users()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	user := &User{}
	err = usersColl.Find(bson.M{"userid": userId}).One(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
