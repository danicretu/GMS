package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"time"
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

func addTags(m *MongoDBConn, tags []string, photo Photo) {

	c := m.session.DB("gmsTry").C("tags")
	for tag := range tags {
		result := Tag{}
		err := c.Find(bson.M{"tag": tags[tag]}).One(&result)
		if err != nil {
			fmt.Println("error while finding tag", tags[tag], "-inserting new tag in database")
			result.Name = tags[tag]
			result.Photos = make([]Photo, 1)
			result.Photos[0] = photo
			err2 := c.Insert(result)
			if err2 != nil {
				fmt.Println("error while adding tag ", result.Name)
			}
		} else {
			result.Photos = append(result.Photos, photo)
			err = c.Update(bson.M{"tag": result.Name}, bson.M{"$set": bson.M{"photos": result.Photos}})
			if err != nil {
				fmt.Println("error while trying to update tag ", result.Name)
			}
		}

	}

}

func findByTag(m *MongoDBConn, tag string) *Tag {
	c := m.session.DB("gmsTry").C("tags")
	result := Tag{}
	err := c.Find(bson.M{"tag": tag}).One(&result)
	if err != nil {
		fmt.Println("Error finding tag")
		fmt.Println(err)
		return nil
	}
	return &result
}

func getAllTags(m *MongoDBConn) []Tag {
	c := m.session.DB("gmsTry").C("tags")
	var result []Tag
	err := c.Find(nil).All(&result)
	if err != nil {
		fmt.Println("Error finding tag")
		fmt.Println(err)
		return nil
	}

	return result
}

func find(m *MongoDBConn, email string) *User {
	result := User{}
	c := m.session.DB("gmsTry").C("user")
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func findUser(m *MongoDBConn, id string) *User {
	result := User{}
	c := m.session.DB("gmsTry").C("user")
	err := c.Find(bson.M{"userId": id}).One(&result)
	if err != nil {
		return nil

	}

	return &result
}

func createDefaultAlbum(ownerId string, ownerName string, picture string) []Album {
	albums := make([]Album, 1)
	id := bson.NewObjectId()

	photos := make([]Photo, 0)

	album := Album{id, id.Hex(), ownerId, ownerName, "Default Album", "", photos}
	albums[0] = album

	return albums
}

func createAlbum(name string, description string, email string, m *MongoDBConn) string {
	user := find(m, email)

	id := bson.NewObjectId()
	album := Album{id, id.Hex(), user.Id, user.FirstName + " " + user.LastName, name, description, make([]Photo, 0)}

	user.Albums = append(user.Albums, album)
	err := m.session.DB("gmsTry").C("user").Update(bson.M{"email": user.Email}, bson.M{"$set": bson.M{"albums": user.Albums}})
	if err != nil {
		fmt.Println(err)
	}
	return id.Hex()
}

func updateTagDB(photo Photo, m *MongoDBConn) {
	tags := photo.Tags
	for tag := range tags {
		query := bson.M{
			"tag":            tags[tag],
			"photos.photoId": photo.PhotoId,
		}

		update := bson.M{
			"$set": bson.M{
				"photos.$.comments": photo.Comments,
				"photos.$.views":    photo.Views,
			},
		}

		err := m.session.DB("gmsTry").C("tags").Update(query, update)
		if err != nil {
			fmt.Println("could not update comments in tag db")
		}
	}
}

func updateMostViewed(photo Photo, m *MongoDBConn) {
	var photos []Photo
	err := m.session.DB("gmsTry").C("mostViewed").Find(nil).All(&photos)
	if err != nil {
		fmt.Println("could not get all users")
	} else {
		fmt.Println("////////////////////////////////")
		fmt.Println(photos)
		if len(photos) == 0 {
			c := m.session.DB("gmsTry").C("mostViewed")
			err := c.Insert(photo)
			if err != nil {
				fmt.Println("could not insert photo into most viewed")
			}
		} else if len(photos) < 5 && len(photos) > 0 {

		}
	}
}

func updateMostRecent(photo Photo, m *MongoDBConn) {
	var p DisplayPhotos
	c := m.session.DB("gmsTry").C("displayPhotos")
	err := c.Find(bson.M{"name": "recent"}).One(&p)

	if err != nil {
		p.Name = "recent"
		p.Photos = make([]Photo, 1)
		p.Photos[0] = photo
		err = c.Insert(p)
		if err != nil {
			fmt.Println("could not insert photo into most recent")
			fmt.Println(err)
		}
		return

	} else if len(p.Photos) < 5 {
		fmt.Println(p.Photos)
		p.Photos = append(p.Photos, photo)
	} else {
		p.Photos = p.Photos[1:]
		p.Photos = append(p.Photos, photo)
	}
	err = c.Update(bson.M{"name": "recent"}, bson.M{"$set": bson.M{"photos": p.Photos}})
	if err != nil {
		fmt.Println(err)
	}
}
