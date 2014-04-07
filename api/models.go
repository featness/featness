package api

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

type Team struct {
	Id      bson.ObjectId   `json:"id"        bson:"_id,omitempty"`
	Name    string          `json:"name"`
	Members []bson.ObjectId `json:"members"`
}

func Users() *mgo.Collection {
	conn, _ := Conn()
	return conn.C("users")
}

func Teams() *mgo.Collection {
	conn, _ := Conn()
	return conn.C("teams")
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

func GetTeamsFor(member bson.ObjectId) ([]Team, error) {
	teams := []Team{}
	err := Teams().Find(bson.M{"members": member}).All(&teams)

	if err != nil {
		return nil, err
	}

	return teams, nil
}

func GetOrCreateTeam(name string, members ...User) (*Team, error) {
	usersBson := make([]bson.ObjectId, len(members))
	for i, v := range members {
		usersBson[i] = v.Id
	}

	teamsColl := Teams()

	_, err := teamsColl.Upsert(
		bson.M{"name": name},
		&Team{bson.NewObjectId(), name, usersBson},
	)

	if err != nil {
		return nil, err
	}

	team := &Team{}
	err = teamsColl.Find(bson.M{"name": name}).One(team)
	return team, nil
}
