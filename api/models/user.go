package models

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
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

	user := &User{}
	err = usersColl.Find(bson.M{"userid": userID}).One(user)
	if err != nil {
		return nil, err
	}

	if user == nil {
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

	userId, err := regexp.Compile(fmt.Sprintf("/%s/", name))
	if err != nil {
		return nil, err
	}

	users := &[]User{}
	err = usersColl.Find(bson.M{"userid": userId}).One(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}
