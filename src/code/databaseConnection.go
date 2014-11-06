package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBConn struct {
	session *mgo.Session
}

func NewMongoDBConn() *MongoDBConn {
	return &MongoDBConn{}
}

func (m *MongoDBConn) connect() *mgo.Session {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	fmt.Println("connect")
	m.session = session
	return m.session
}

func add(m *MongoDBConn, user User) {

	c := m.session.DB("gmsTry").C("user")
	err := c.Insert(user)
	if err != nil {
		panic(err)
	}

}

func find(m *MongoDBConn, email string) *User {
	result := User{}
	c := m.session.DB("gmsTry").C("user")
	fmt.Println(email)
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return nil

	}

	return &result
}

func createDefaultAlbum() []Album {
	albums := make([]Album, 1)
	//id := bson.NewObjectId()
	album := Album{"Default Album", "", make([]Photo, 1)}
	albums[0] = album

	return albums
}
