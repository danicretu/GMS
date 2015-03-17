package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	//"time"
	"math/rand"
	"strings"
)

type MongoDBConn struct {
	session *mgo.Session
}

func NewMongoDBConn() *MongoDBConn {
	return &MongoDBConn{}
}

//var db_name = "gmsTry"

var db_name = "ugc"
var flickrDB = "gmsTry"

func (m *MongoDBConn) connect() *mgo.Session {
	session, err := mgo.Dial("mongodb://ugc:ugc_pass@imcdserv1.dcs.gla.ac.uk/ugc")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	m.session = session
	return m.session
}

func (m *MongoDBConn) connectFlickr() *mgo.Session {
	session, err := mgo.Dial("mongodb://gms:rdm$248@imcdserv1.dcs.gla.ac.uk/gmsTry")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	m.session = session
	return m.session
}

func add(user User) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	c := sess.DB(db_name).C("user")
	err := c.Insert(user)
	if err != nil {
		panic(err)
	}

	defer sess.Close()

}

func addTags(tags []string, photo Photo, video Video) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	c := sess.DB(db_name).C("tags")
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
	defer sess.Close()

}

func findByTag(tag string) *Tag {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	c := sess.DB(db_name).C("tags")
	result := Tag{}
	err := c.Find(bson.M{"tag": tag}).One(&result)
	if err != nil {
		fmt.Println("Error finding tag")
		fmt.Println(err)
		defer sess.Close()
		return nil
	}

	defer sess.Close()
	return &result

}

func getAllTags() []Tag {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	c := sess.DB(db_name).C("tags")
	var result []Tag
	err := c.Find(nil).All(&result)
	if err != nil {
		fmt.Println("Error finding tag")
		fmt.Println(err)
		defer sess.Close()
		return nil
	}
	defer sess.Close()

	return result
}

func find(email string) *User {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	result := User{}
	c := sess.DB(db_name).C("user")
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		defer sess.Close()
		return nil
	}
	defer sess.Close()
	return &result
}

func findUser(id string) *User {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	result := User{}
	c := sess.DB(db_name).C("user")
	err := c.Find(bson.M{"userId": id}).One(&result)
	if err != nil {
		defer sess.Close()
		return nil

	}
	defer sess.Close()
	return &result
}
func getFlickrMap() []FlickrImage1 {
	dbConn := NewMongoDBConn()
	sess1 := dbConn.connectFlickr()
	c := sess1.DB(flickrDB).C("gmsFlickr1")

	var flickrImage []FlickrImage1

	err := c.Find(bson.M{"$and": []bson.M{bson.M{"latitude": bson.M{"$ne": 0}}, bson.M{"longitude": bson.M{"$ne": 0}}}}).All(&flickrImage)
	if err != nil {
		fmt.Println(err)
	}
	defer sess1.Close()

	return flickrImage
}

func getFlickrMain(tag string, tag2 string, start int, cType string, location string) []FlickrImage1 {

	source := "/resources/flickr/"
	dbConn := NewMongoDBConn()
	sess1 := dbConn.connectFlickr()
	c := sess1.DB(flickrDB).C("gmsFlickr1")
	c1 := sess1.DB(flickrDB).C("gmsFlickrCWGUpdated")
	//c1 := sess1.DB(flickrDB).C("flickrCWG")
	var flickrImage []FlickrImage1

	limit := 8
	if cType == "location" {
		if location != "" {
			err := c1.Find(bson.M{"exifLocation": location}).Skip(start * limit).Limit(limit).All(&flickrImage)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			var myarr = []string{tag}
			err := c1.Find(bson.M{"keywords": bson.M{"$all": myarr}}).Skip(start * limit).Limit(limit).All(&flickrImage)
			if err != nil {
				fmt.Println(err)
			}
		}

		for i := range flickrImage {

			date := strings.Split(flickrImage[i].TimeStamp, " ")
			t := strings.Split(date[0], "/")
			folderName := t[0] + "_" + t[1] + "_" + t[2]
			flickrImage[i].URL = source + folderName + "/" + flickrImage[i].ImageName
			k := 0
			for k < len(flickrImage[i].Keywords) {
				if strings.Contains(flickrImage[i].Keywords[k], "|") {
					flickrImage[i].Keywords = append(flickrImage[i].Keywords[:k], flickrImage[i].Keywords[k+1:]...)
					k -= 1
				} else {
					k += 1
				}
			}
			tim, _ := time.Parse("02/01/2006 15:04:05", flickrImage[i].TimeStamp)
			flickrImage[i].TimeStamp = tim.Format("02/01/2006")
		}
		defer sess1.Close()
		return flickrImage
	}
	if tag == "" {
		records, _ := c.Find(bson.M{}).Count()
		a := rand.Intn(records - limit)
		//err := c.Find(bson.M{"source": "https://www.flickr.com"}).Skip(10).Limit(8).All(&flickrImage)
		err := c.Find(bson.M{}).Skip(a).Limit(limit).All(&flickrImage)
		if err != nil {
			fmt.Println(err)
		}
	} else if tag != "" && tag2 == "" {
		var myarr = []string{tag}

		records, _ := c.Find(bson.M{"keywords": bson.M{"$all": myarr}}).Count()
		random := rand.Intn(records)

		//err := c.Find(bson.M{"source": "https://www.flickr.com", "keywords": bson.M{"$all": myarr}}).Skip(10).Limit(8).All(&flickrImage)
		err := c.Find(bson.M{"keywords": bson.M{"$all": myarr}}).Skip(random).Limit(limit).All(&flickrImage)
		if err != nil {
			fmt.Println(err)
		}
	} else if tag == tag2 {
		var myarr = []string{tag}
		err := c.Find(bson.M{"keywords": bson.M{"$all": myarr}}).Skip(start * limit).Limit(limit).All(&flickrImage)
		//err := c.Find(bson.M{"source": "https://www.flickr.com", "keywords": bson.M{"$all": myarr}}).Skip(start).Limit(8).All(&flickrImage)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		var myarr = []string{tag, tag2}
		if cType == "and" {
			err := c.Find(bson.M{"keywords": bson.M{"$all": myarr}}).Skip(start * limit).Limit(limit).All(&flickrImage)
			if err != nil {
				fmt.Println(err)
			}
		} else if cType == "or" {
			err := c.Find(bson.M{"keywords": bson.M{"$in": myarr}}).Skip(start * limit).Limit(limit).All(&flickrImage)
			if err != nil {
				fmt.Println(err)
			}
		}
		//err := c.Find(bson.M{"source": "https://www.flickr.com", "keywords": bson.M{"$all": myarr}}).Skip(start).Limit(8).All(&flickrImage)

	}

	for i := range flickrImage {

		date := strings.Split(flickrImage[i].TimeStamp, " ")
		t := strings.Split(date[0], "/")
		folderName := t[0] + "_" + t[1] + "_" + t[2]
		flickrImage[i].URL = source + folderName + "/" + flickrImage[i].ImageName
		k := 0
		for k < len(flickrImage[i].Keywords) {

			if strings.Contains(flickrImage[i].Keywords[k], "|") {
				flickrImage[i].Keywords = append(flickrImage[i].Keywords[:k], flickrImage[i].Keywords[k+1:]...)
				k -= 1
			} else {
				k += 1
			}
		}
		tim, _ := time.Parse("02/01/2006 15:04:05", flickrImage[i].TimeStamp)
		flickrImage[i].TimeStamp = tim.Format("02/01/2006")
	}
	defer sess1.Close()
	return flickrImage

}

func getFlickrImages(tag string, start int) []FlickrImage {
	source := "/resources/flickr/"
	dbConn := NewMongoDBConn()
	sess1 := dbConn.connectFlickr()
	c := sess1.DB(flickrDB).C("gmsNewsScottish")
	limit := 8
	var flickrImage []FlickrImage
	var myarr = []string{tag}
	if start == 0 && tag == "" {
		records, _ := c.Find(bson.M{"source": "https://www.flickr.com"}).Count()
		start = rand.Intn(records - limit)
		err := c.Find(bson.M{"source": "https://www.flickr.com"}).Skip(start).Limit(limit).All(&flickrImage)
		if err != nil {
			fmt.Println(err)
		}
	} else {

		err := c.Find(bson.M{"source": "https://www.flickr.com", "keywords": bson.M{"$all": myarr}}).Skip(start * limit).Limit(limit).All(&flickrImage)
		if err != nil {
			fmt.Println(err)
		}
	}

	for i := range flickrImage {

		date := strings.Split(flickrImage[i].TimeStamp, " ")
		t := strings.Split(date[0], "/")
		folderName := t[0] + "_" + t[1] + "_" + t[2]
		flickrImage[i].URL = source + folderName + "/" + flickrImage[i].ImageName
		k := 0
		for k < len(flickrImage[i].Keywords) {
			if strings.Contains(flickrImage[i].Keywords[k], "|") {
				flickrImage[i].Keywords = append(flickrImage[i].Keywords[:k], flickrImage[i].Keywords[k+1:]...)
				k -= 1
			} else {
				k += 1
			}
		}
		tim, _ := time.Parse("02/01/2006 15:04:05", flickrImage[i].TimeStamp)
		flickrImage[i].TimeStamp = tim.Format("02/01/2006")
	}
	defer sess1.Close()
	return flickrImage
}

func getNews(tag string, start int) []News {
	dbConn := NewMongoDBConn()
	sess1 := dbConn.connectFlickr()
	c := sess1.DB(flickrDB).C("gmsNewsScottish")
	var news []News
	var newsKey = []string{tag}
	limit := 8

	if start == 0 && tag == "" {
		records, _ := c.Find(bson.M{"source": "http://www.theguardian.com"}).Count()
		start = rand.Intn(records - limit)
		err := c.Find(bson.M{"source": "http://www.theguardian.com"}).Skip(start).Limit(limit).All(&news)
		if err != nil {
			fmt.Println(err)
		}
	} else {

		err := c.Find(bson.M{"source": "http://www.theguardian.com", "keywords": bson.M{"$all": newsKey}}).Skip(start * limit).Limit(limit).All(&news)
		if err != nil {
			fmt.Println(err)
		}
	}

	for i := range news {
		if news[i].Images != nil {
			news[i].ImageName = news[i].Images[0].Name
			news[i].ImageUrl = news[i].Images[0].Name
			news[i].ImageCaption = news[i].Images[0].Caption
		}

	}

	/*for i := range flickrImage {
		date := strings.Split(flickrImage[i].TimeStamp, " ")
		t := strings.Split(date[0], "/")
		folderName := t[0] + "_" + t[1] + "_" + t[2]
		flickrImage[i].URL = source + folderName + "/" + flickrImage[i].ImageName
	} */
	defer sess1.Close()
	return news
}

func createDefaultAlbum(ownerId string, ownerName string) {
	id := bson.NewObjectId()

	album := Album{id, id.Hex(), ownerId, ownerName, "Default Album"}
	dbConnection = NewMongoDBConn()
	sess = dbConnection.connect()
	c := sess.DB(db_name).C("albums")
	err := c.Insert(album)
	if err != nil {
		panic(err)
	}

	defer sess.Close()
}

func createAlbum(name string, uId string, uName string) string {
	id := bson.NewObjectId()
	album := Album{id, id.Hex(), uId, uName, name}
	dbConnection = NewMongoDBConn()
	sess = dbConnection.connect()
	c := sess.DB(db_name).C("albums")
	err := c.Insert(album)
	if err != nil {
		panic(err)
	}

	defer sess.Close()

	return album.AlbumId
}

func deleteFromDisplay(content string, cType string) {
	var p DisplayPhotos
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	c := sess.DB(db_name).C("displayPhotos")
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
	defer sess.Close()

}

func deleteFromTag(content string, cType string) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var t Tag
	if cType == "image" {
		var photo Photo
		err := sess.DB(db_name).C("photos").Find(bson.M{"photoId": content}).One(&photo)
		if err != nil {
			fmt.Println(err)
		}

		for r := range photo.Tags {

			err = sess.DB(db_name).C("tags").Find(bson.M{"tag": photo.Tags[r]}).One(&t)
			for x := range t.Photos {
				if t.Photos[x].PhotoId == content {
					t.Photos = append(t.Photos[:x], t.Photos[x+1:]...)
					break
				}

			}
			err = sess.DB(db_name).C("tags").Update(bson.M{"tag": photo.Tags[r]}, bson.M{"$set": bson.M{"photos": t.Photos}})
		}

	} else {
		var video Video
		err := sess.DB(db_name).C("videos").Find(bson.M{"videoId": content}).One(&video)
		if err != nil {
			fmt.Println(err)
		}

		for r := range video.Tags {

			err = sess.DB(db_name).C("tags").Find(bson.M{"tag": video.Tags[r]}).One(&t)
			for x := range t.Videos {
				if t.Videos[x].VideoId == content {
					t.Videos = append(t.Videos[:x], t.Videos[x+1:]...)
					break
				}

			}
			err = sess.DB(db_name).C("tags").Update(bson.M{"tag": video.Tags[r]}, bson.M{"$set": bson.M{"videos": t.Videos}})
		}
	}

	defer sess.Close()
}

func deleteFromOthers(content string, cType string) {
	deleteFromDisplay(content, cType)
	deleteFromTag(content, cType)

}

func updateTagDB(photo Photo, video Video) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

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
			err := sess.DB(db_name).C("tags").Update(query, update)
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

			err := sess.DB(db_name).C("tags").Update(query, update)
			if err != nil {
				fmt.Println("could not update comments in tag db")
			}

		}
	}

	defer sess.Close()
}

func updateMostViewed(photo Photo, video Video) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var p DisplayPhotos
	c := sess.DB(db_name).C("displayPhotos")
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
		if len(p.Photos) < 8 {
			flag := false
			for m := range p.Photos {
				if p.Photos[m].PhotoId == photo.PhotoId {
					p.Photos[m].Views = photo.Views
					p.Photos[m].Comments = photo.Comments
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
					p.Photos[m].Comments = photo.Comments
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
		if len(p.Videos) < 8 {
			flag := false
			for m := range p.Videos {
				if p.Videos[m].VideoId == video.VideoId {
					p.Videos[m].Views = video.Views
					p.Videos[m].Comments = video.Comments
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

	defer sess.Close()
}

func insertInMostRecent(photo Photo, video Video) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var p DisplayPhotos
	c := sess.DB(db_name).C("displayPhotos")
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
		defer sess.Close()
		return

	} else if photo.PhotoId != "" {
		if len(p.Photos) < 8 {
			fmt.Println(p.Photos)
			p.Photos = append(p.Photos, photo)
		} else {
			p.Photos = p.Photos[1:]
			p.Photos = append(p.Photos, photo)
		}
		err = c.Update(bson.M{"name": "recent"}, bson.M{"$set": bson.M{"photos": p.Photos}})
	} else {
		if len(p.Videos) < 8 {
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

	defer sess.Close()
}

func updateMostRecent(photo Photo, video Video) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	if photo.PhotoId != "" {
		query := bson.M{
			"name":           "recent",
			"photos.photoId": photo.PhotoId,
		}

		update := bson.M{
			"$set": bson.M{
				"photos.$.comments": photo.Comments,
				"photos.$.views":    photo.Views,
			},
		}
		err := sess.DB(db_name).C("displayPhotos").Update(query, update)
		if err != nil {
			fmt.Println("could not update comments in recent db")
		} else {
			fmt.Println("updated in photo recent")
		}
	} else {

		query := bson.M{
			"name":           "recent",
			"videos.videoId": video.VideoId,
		}

		update := bson.M{
			"$set": bson.M{
				"videos.$.comments": video.Comments,
				"videos.$.views":    video.Views,
			},
		}

		err := sess.DB(db_name).C("displayPhotos").Update(query, update)
		if err != nil {
			fmt.Println("could not update comments in ecent db")
		} else {
			fmt.Println("updated in recent db")
		}

	}
	defer sess.Close()

}

func getMapImages(user string) []MapImage {
	dbConn := NewMongoDBConn()
	sess1 := dbConn.connectFlickr()
	var pics []MapImage
	if user == "" {
		err := sess1.DB(flickrDB).C("locationDB").Find(bson.M{}).All(&pics)
		if err != nil {
			fmt.Println("could not get map images from db")
		}
	} else {
		err := sess1.DB(flickrDB).C("locationDB").Find(bson.M{"user": bson.ObjectIdHex(user)}).All(&pics)
		if err != nil {
			fmt.Println("could not get map images for user from DB")
		}
	}
	defer sess1.Close()
	return pics
}

func getCwgMapImages() []CwgImage {
	dbConn := NewMongoDBConn()
	sess1 := dbConn.connectFlickr()
	var pics []CwgImage
	err := sess1.DB(flickrDB).C("cwgLocations").Find(bson.M{}).All(&pics)
	if err != nil {
		fmt.Println("could not get CWG images from db")
	}
	defer sess1.Close()
	return pics
}
