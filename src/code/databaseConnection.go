package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"time"
	"strings"
)

type MongoDBConn struct {
	session *mgo.Session
}

func NewMongoDBConn() *MongoDBConn {
	return &MongoDBConn{}
}

var db_name = "db_name"
var flickrDB = "gmsTry"

func (m *MongoDBConn) connect() *mgo.Session {
	session, err := mgo.Dial("mongodb://gms:rdm$248@imcdserv1.dcs.gla.ac.uk/gms")
	if err != nil {
		panic(err)
	}

	fmt.Println("connect")
	m.session = session
	return m.session
}

func (m *MongoDBConn) connectFlickr() *mgo.Session {
	session, err := mgo.Dial("mongodb://gms:rdm$248@imcdserv1.dcs.gla.ac.uk/gms")
	if err != nil {
		panic(err)
	}

	fmt.Println("connect")
	m.session = session
	return m.session
}

func add(m *MongoDBConn, user User) {

	c := m.session.DB(db_name).C("user")
	err := c.Insert(user)
	if err != nil {
		panic(err)
	}

}

func addTags(m *MongoDBConn, tags []string, photo Photo, video Video) {

	c := m.session.DB(db_name).C("tags")
	for tag := range tags {
		result := Tag{}
		err := c.Find(bson.M{"tag": tags[tag]}).One(&result)
		if err != nil {
			fmt.Println("error while finding tag", tags[tag], "-inserting new tag in database")
			result.Name = tags[tag]
			result.Videos = make([]Video, 0)
			result.Photos = make([]Photo, 0)

			if photo.PhotoId != "" {
				result.Photos = append(result.Photos, photo)
			} else {
				result.Videos = append(result.Videos, video)
			}

			err2 := c.Insert(result)
			if err2 != nil {
				fmt.Println("error while adding tag ", result.Name)
			}
		} else {
			if photo.PhotoId != "" {
				result.Photos = append(result.Photos, photo)
				err = c.Update(bson.M{"tag": result.Name}, bson.M{"$set": bson.M{"photos": result.Photos}})
			} else {
				result.Videos = append(result.Videos, video)
				err = c.Update(bson.M{"tag": result.Name}, bson.M{"$set": bson.M{"videos": result.Videos}})
			}
			if err != nil {
				fmt.Println("error while trying to update tag ", result.Name)
			}
		}

	}

}

func findByTag(m *MongoDBConn, tag string) *Tag {
	c := m.session.DB(db_name).C("tags")
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
	c := m.session.DB(db_name).C("tags")
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
	c := m.session.DB(db_name).C("user")
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func findUser(m *MongoDBConn, id string) *User {
	result := User{}
	c := m.session.DB(db_name).C("user")
	err := c.Find(bson.M{"userId": id}).One(&result)
	if err != nil {
		return nil

	}

	return &result
}

func getFlickrImages(m *MongoDBConn, tag string) []FlickrImage {
	source := "/home/dani/go-programs/flickrData/"
	dbConnection = NewMongoDBConn()
	_ = dbConnection.connectFlickr()
	c := dbConnection.session.DB("gmsTry").C("gmsNewsScottish")
	var flickrImage []FlickrImage
	fmt.Println("the keyword is...", tag)
	var myarr = []string{tag}

	err := c.Find(bson.M{"source": "https://www.flickr.com", "keywords": bson.M{"$all": myarr}}).Skip(0).Limit(8).All(&flickrImage)
	if err != nil {
		fmt.Println(err)
	}

	for i := range flickrImage {
		date := strings.Split(flickrImage[i].TimeStamp, " ")
		t := strings.Split(date[0], "/")
		folderName := t[0] + "_" + t[1] + "_" + t[2]
		flickrImage[i].URL = source + folderName + "/" + flickrImage[i].ImageName
	}

	return flickrImage
}

func getNews(m *MongoDBConn, tag string) []News {
	dbConnection = NewMongoDBConn()
	_ = dbConnection.connectFlickr()
	c := dbConnection.session.DB("gmsTry").C("gmsNewsScottish")
	var news []News
	var newsKey = []string{tag}

	err := c.Find(bson.M{"source": "http://www.theguardian.com", "keywords": bson.M{"$all": newsKey}}).Skip(0).Limit(8).All(&news)
	if err != nil {
		fmt.Println(err)
	}

	for i := range news {
		if news[i].Images != nil {
			news[i].ImageName = news[i].Images[0].Name
			news[i].ImageUrl = news[i].Images[0].Name
			news[i].ImageCaption = news[i].Images[0].Caption
		}

	}

	fmt.Println(news)

	return news
}

func createDefaultAlbum(m *MongoDBConn, ownerId string, ownerName string) {
	id := bson.NewObjectId()

	album := Album{id, id.Hex(), ownerId, ownerName, "Default Album"}

	c := m.session.DB(db_name).C("albums")
	err := c.Insert(album)
	if err != nil {
		panic(err)
	}
}

func createAlbum(name string, uId string, uName string, m *MongoDBConn) string {
	id := bson.NewObjectId()
	album := Album{id, id.Hex(), uId, uName, name}

	c := m.session.DB(db_name).C("albums")
	err := c.Insert(album)
	if err != nil {
		panic(err)
	}

	return album.AlbumId
}

func deleteFromDisplay(m *MongoDBConn, content string, cType string) {
	var p DisplayPhotos
	c := m.session.DB(db_name).C("displayPhotos")
	err := c.Find(bson.M{"name": "views"}).One(&p)

	var r DisplayPhotos
	err = c.Find(bson.M{"name": "recent"}).One(&r)

	if cType == "image" {
		for i := range p.Photos {
			if p.Photos[i].PhotoId == content {
				p.Photos = append(p.Photos[:i], p.Photos[i+1:]...)
				err = c.Update(bson.M{"name": "views"}, bson.M{"$set": bson.M{"photos": p.Photos}})
				if err != nil {
					fmt.Println("could not delete from views ", err)
				}
				break
			}
		}

		for i := range r.Photos {
			if r.Photos[i].PhotoId == content {
				r.Photos = append(r.Photos[:i], r.Photos[i+1:]...)
				err = c.Update(bson.M{"name": "recent"}, bson.M{"$set": bson.M{"photos": r.Photos}})
				if err != nil {
					fmt.Println("could not delete from recent ", err)
				}
				break
			}
		}

	} else {
		for i := range p.Videos {
			if p.Videos[i].VideoId == content {
				p.Videos = append(p.Videos[:i], p.Videos[i+1:]...)
				err = c.Update(bson.M{"name": "views"}, bson.M{"$set": bson.M{"videos": p.Videos}})
				if err != nil {
					fmt.Println("could not delete from views ", err)
				}
				break
			}
		}

		for i := range r.Videos {
			if r.Videos[i].VideoId == content {
				r.Videos = append(r.Videos[:i], r.Videos[i+1:]...)
				err = c.Update(bson.M{"name": "recent"}, bson.M{"$set": bson.M{"videos": r.Videos}})
				if err != nil {
					fmt.Println("could not delete from views ", err)
				}
				break
			}
		}

	}

}

func deleteFromTag(m *MongoDBConn, content string, cType string) {
	var t Tag
	if cType == "image" {
		var photo Photo
		err := m.session.DB(db_name).C("photos").Find(bson.M{"photoId": content}).One(&photo)
		if err != nil {
			fmt.Println(err)
		}

		for r := range photo.Tags {

			err = dbConnection.session.DB(db_name).C("tags").Find(bson.M{"tag": photo.Tags[r]}).One(&t)
			for x := range t.Photos {
				if t.Photos[x].PhotoId == content {
					t.Photos = append(t.Photos[:x], t.Photos[x+1:]...)
					break
				}

			}
			err = dbConnection.session.DB(db_name).C("tags").Update(bson.M{"tag": photo.Tags[r]}, bson.M{"$set": bson.M{"photos": t.Photos}})
		}

	} else {
		var video Video
		err := m.session.DB(db_name).C("videos").Find(bson.M{"videoId": content}).One(&video)
		if err != nil {
			fmt.Println(err)
		}

		for r := range video.Tags {

			err = dbConnection.session.DB(db_name).C("tags").Find(bson.M{"tag": video.Tags[r]}).One(&t)
			for x := range t.Videos {
				if t.Videos[x].VideoId == content {
					t.Videos = append(t.Videos[:x], t.Videos[x+1:]...)
					break
				}

			}
			err = dbConnection.session.DB(db_name).C("tags").Update(bson.M{"tag": video.Tags[r]}, bson.M{"$set": bson.M{"videos": t.Videos}})
		}
	}
}

func deleteFromOthers(m *MongoDBConn, content string, cType string) {
	deleteFromDisplay(m, content, cType)
	deleteFromTag(m, content, cType)

}

func updateTagDB(photo Photo, video Video, m *MongoDBConn) {
	tags := photo.Tags
	for tag := range tags {
		if photo.PhotoId != "" {
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
			err := m.session.DB(db_name).C("tags").Update(query, update)
			if err != nil {
				fmt.Println("could not update comments in tag db")
			}
		} else {

			query := bson.M{
				"tag":            tags[tag],
				"videos.videoId": video.VideoId,
			}

			update := bson.M{
				"$set": bson.M{
					"videos.$.comments": video.Comments,
					"videos.$.views":    video.Views,
				},
			}

			err := m.session.DB(db_name).C("tags").Update(query, update)
			if err != nil {
				fmt.Println("could not update comments in tag db")
			}

		}
	}
}

func updateMostViewed(photo Photo, video Video, m *MongoDBConn) {
	var p DisplayPhotos
	c := m.session.DB(db_name).C("displayPhotos")
	err := c.Find(bson.M{"name": "views"}).One(&p)

	if err != nil {
		p.Name = "views"
		p.Photos = make([]Photo, 0)
		p.Videos = make([]Video, 0)
		if photo.PhotoId != "" {
			p.Photos = append(p.Photos, photo)
		} else {
			p.Videos = append(p.Videos, video)
		}
		err = c.Insert(p)
		if err != nil {
			fmt.Println("could not insert photo into most recent")
			fmt.Println(err)
		}
		return
	} else if photo.PhotoId != "" {
		if len(p.Photos) < 5 {
			flag := false
			for m := range p.Photos {
				if p.Photos[m].PhotoId == photo.PhotoId {
					p.Photos[m].Views = photo.Views
					flag = true
				}
			}
			if flag == false {
				p.Photos = append(p.Photos, photo)
			}
		} else {
			flag := false
			low := p.Photos[0].Views
			index := 0
			for m := range p.Photos {
				if p.Photos[m].PhotoId == photo.PhotoId {
					p.Photos[m].Views = photo.Views
					flag = true
				}

				if p.Photos[m].Views < low {
					low = p.Photos[m].Views
					index = m
				}
			}
			if flag == false {
				if photo.Views > p.Photos[index].Views {
					p.Photos[index] = photo
				}
			}
		}
		err = c.Update(bson.M{"name": "views"}, bson.M{"$set": bson.M{"photos": p.Photos}})
	} else {
		if len(p.Videos) < 5 {
			flag := false
			for m := range p.Videos {
				if p.Videos[m].VideoId == video.VideoId {
					p.Videos[m].Views = video.Views
					flag = true
				}
			}
			if flag == false {
				p.Videos = append(p.Videos, video)
			}
		} else {
			flag := false
			low := p.Videos[0].Views
			index := 0
			for m := range p.Videos {
				if p.Videos[m].VideoId == video.VideoId {
					p.Videos[m].Views = video.Views
					flag = true
				}

				if p.Videos[m].Views < low {
					low = p.Videos[m].Views
					index = m
				}
			}
			if flag == false {
				if video.Views > p.Videos[index].Views {
					p.Videos[index] = video
				}
			}
		}
		fmt.Println("in update most viewed")
		err = c.Update(bson.M{"name": "views"}, bson.M{"$set": bson.M{"videos": p.Videos}})

	}

	if err != nil {
		fmt.Println(err)
	}
}

func updateMostRecent(photo Photo, video Video, m *MongoDBConn) {
	var p DisplayPhotos
	c := m.session.DB(db_name).C("displayPhotos")
	err := c.Find(bson.M{"name": "recent"}).One(&p)

	if err != nil {
		p.Name = "recent"
		p.Photos = make([]Photo, 0)
		p.Videos = make([]Video, 0)
		if photo.PhotoId != "" {
			p.Photos = append(p.Photos, photo)
		} else {
			p.Videos = append(p.Videos, video)
		}
		err = c.Insert(p)
		if err != nil {
			fmt.Println("could not insert photo into most recent")
			fmt.Println(err)
		}
		return

	} else if photo.PhotoId != "" {
		if len(p.Photos) < 5 {
			fmt.Println(p.Photos)
			p.Photos = append(p.Photos, photo)
		} else {
			p.Photos = p.Photos[1:]
			p.Photos = append(p.Photos, photo)
		}
		err = c.Update(bson.M{"name": "recent"}, bson.M{"$set": bson.M{"photos": p.Photos}})
	} else {
		if len(p.Videos) < 5 {
			fmt.Println(p.Videos)
			p.Videos = append(p.Videos, video)
		} else {
			p.Videos = p.Videos[1:]
			p.Videos = append(p.Videos, video)
		}
		err = c.Update(bson.M{"name": "recent"}, bson.M{"$set": bson.M{"videos": p.Videos}})
	}
	if err != nil {
		fmt.Println(err)
	}
}
