package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rwcarlsen/goexif/exif"
	//"gopkg.in/mgo.v2"
	"bytes"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type User struct {
	UserId         bson.ObjectId `bson:"_id"`
	FirstName      string        `bson:"firstname"`
	LastName       string        `bson:"lastname"`
	Email          string        `bson:"email"`
	Password       string        `bson:"password"`
	ProfilePicture string        `bson:"profilepic"`
	Albums         []Album       `bson:"albums"`
	GoogleId       string        `bson:"gId"`
	FacebookId     string        `bson:"fId"`
	TwitterId      string        `bson:"tId"`
	Id             string        `bson:"userId"`
}

type Album struct {
	Id          bson.ObjectId `bson:"_id"`
	AlbumId     string        `bson:"albumId"`
	Owner       string        `bson:"owner"`
	OwnerName   string        `bson:"ownerName"`
	Name        string        `bson:"albumname"`
	Description string        `bson:"description"`
	Photo       []Photo       `bson:"photos"`
	Video       []Video       `bson:"albums"`
}

type Photo struct {
	Id          bson.ObjectId  `bson:"_id"`
	PhotoId     string         `bson:"photoId"`
	Owner       string         `bson:"owner"`
	OwnerName   string         `bson:"ownerName"`
	AlbumId     string         `bson:"albumId"`
	URL         string         `bson:"url"`
	Description string         `bson:"description"`
	Location    Location       `bson:"location"`
	Timestamp   string         `bson:"timestamp"`
	Views       int            `bson:"views"`
	Tags        []string       `bson:"tags"`
	Comments    []PhotoComment `bson:"comments"`
}

type Video struct {
	Id          bson.ObjectId  `bson:"_id"`
	VideoId     string         `bson:"videoId"`
	Owner       string         `bson:"owner"`
	OwnerName   string         `bson:"ownerName"`
	AlbumId     string         `bson:"albumId"`
	URL         string         `bson:"url"`
	Description string         `bson:"description"`
	Location    Location       `bson:"location"`
	Timestamp   string         `bson:"timestamp"`
	Views       int            `bson:"views"`
	Tags        []string       `bson:"tags"`
	Comments    []PhotoComment `bson:"comments"`
}

type Location struct {
	Name      string `bson:"locationName"`
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type PhotoComment struct {
	User      string `bson:"userName"`
	UserId    string `bson:"userId"`
	Body      string `bson:"comment"`
	Timestamp string `bson:"time"`
}

type PhotoContainer struct {
	Categories []Photo
}

type Tag struct {
	Name   string  `bson:"tag"`
	Photos []Photo `bson:"photos"`
	Videos []Video `bson:"videos"`
}

type DisplayPhotos struct {
	Name   string  `bson:"name"`
	Photos []Photo `bson:"photos"`
	Videos []Video `bson:"videos"`
}

type FlickrTag struct {
	Tags struct {
		Source string `json:"source"`
		Tag    []struct {
			Content string `json:"_content"`
		} `json:"tag"`
	} `json:"tags"`
	Stat string `json:"stat"`
}

var router = mux.NewRouter()

var authKey = []byte("NCDIUyd78DBCSJBlcsd783")

// Encryption Key
var encKey = []byte("nckdajKBDSY6778FDV891bdf")

var store = sessions.NewCookieStore(authKey, encKey)

var dbConnection *MongoDBConn

//add(dbConnection, name, password) ->add to db
//find(dbConnection, name) ->find in db

func main() {
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/logout", handleLogout)
	router.HandleFunc("/register", handleRegister)
	router.HandleFunc("/authenticated", handleAuthenticated)
	router.HandleFunc("/pictures", handlePictures)
	router.HandleFunc("/videos", handleVideos)
	router.HandleFunc("/flickrNews", handleFlickrNews)
	router.HandleFunc("/albums", handleAlbums)
	router.HandleFunc("/upload", handleUpload)
	router.HandleFunc("/uploadPic", uploadHandler)
	router.HandleFunc("/saveComment", handleComments)
	router.HandleFunc("/flickr", handleFlickr)
	router.HandleFunc("/tag", handleTag)
	router.HandleFunc("/tagCloud", createTagCloud)
	router.HandleFunc("/checkLogIn", checkLoggedIn)
	router.HandleFunc("/saveFile", handleSaveImage)
	router.HandleFunc("/createAlbum", handleCreateAlbum)
	router.HandleFunc("/user", handleUserProfile)
	router.HandleFunc("/upvote", handleUpvote)
	router.HandleFunc("/cmsHome", handleCms)
	router.HandleFunc("/flickrTest", handleFlickrTest)
	router.HandleFunc("/delete", handleDelete)
	authenticateGoogle()
	authenticateFacebook()
	authenticateTwitter()

	dbConnection = NewMongoDBConn()
	_ = dbConnection.connect()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println("in delete")

	picture := r.FormValue("pic")
	album := r.FormValue("album")
	owner := r.FormValue("owner")
	cType := r.FormValue("cType")

	fmt.Println("in delete", picture, " ", album, " ", cType)

	user := findUser(dbConnection, owner)

	var al int
	photo := Photo{}
	video := Video{}
	for i := range user.Albums {
		if user.Albums[i].AlbumId == album {
			al = i
			break
		}
	}

	if cType == "image" {
		var pic int

		for i := range user.Albums[al].Photo {
			if user.Albums[al].Photo[i].PhotoId == picture {
				pic = i
				break
			}
		}
		photo = user.Albums[al].Photo[pic]
		user.Albums[al].Photo = append(user.Albums[al].Photo[:pic], user.Albums[al].Photo[pic+1:]...)

	} else {
		var vid int

		for i := range user.Albums[al].Video {
			if user.Albums[al].Video[i].VideoId == picture {
				vid = i
				break
			}
		}
		video = user.Albums[al].Video[vid]
		user.Albums[al].Video = append(user.Albums[al].Video[:vid], user.Albums[al].Video[vid+1:]...)

	}
	deleteFromOthers(dbConnection, photo, video)

	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"_id": user.UserId}, bson.M{"$set": bson.M{"albums": user.Albums}})
	if err != nil {
		fmt.Fprintf(w, "No")
	}
	fmt.Println(picture)
	resp := "Yes_" + picture

	fmt.Fprintf(w, resp)

}

func handleFlickrTest(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		"George",
	}

	var doc bytes.Buffer
	t, _ := template.ParseFiles("test.html")
	t.Execute(&doc, data)
	s := doc.String()

	//fmt.Println(s)
	fmt.Fprintf(w, s)

	/*session2, err := mgo.Dial("localhost:27018")
	//DialWithTimeout(fwd.localAddr, 10*time.Minute)
	//mongodb://imdcserv1.dcs.gla.ac.uk/gmsTry
	if err != nil {
		fmt.Println(err)
	}
	//defer session2.Close()
	fmt.Println(session2)
	admindb := session2.DB("gmsTry")
	fmt.Println(admindb)
	/*
		err = admindb.Login("gms", "rdm$248")
		if err != nil {
			fmt.Println(err)
		}

		coll := session2.DB("gmsTry").C("gmsNewsScottish")
		var result string
		err = coll.Find(bson.M{"source": "http://www.theguardian.com", "url": "http://www.theguardian.com/sport/2014/aug/04/australian-athletes-funding-commonwealth-games"}).One(&result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result) */
}

func handleVideos(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(dbConnection, currentUser)

	authenticated, _ := template.ParseFiles("videos.html")
	authenticated.Execute(w, u)

}

func handleFlickrNews(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(dbConnection, currentUser)

	authenticated, _ := template.ParseFiles("flickrNews.html")
	authenticated.Execute(w, u)

}

func handleCms(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(dbConnection, currentUser)

	if u == nil {
		u = &User{}
	}

	var p DisplayPhotos
	c := dbConnection.session.DB("gmsTry").C("displayPhotos")
	err := c.Find(bson.M{"name": "views"}).One(&p)
	if err != nil {
		fmt.Println("could not get most viewed photos")
	}

	var recent DisplayPhotos
	c = dbConnection.session.DB("gmsTry").C("displayPhotos")
	err = c.Find(bson.M{"name": "recent"}).One(&recent)
	if err != nil {
		fmt.Println("could not get most viewed photos")
	}

	data := struct {
		P DisplayPhotos
		R DisplayPhotos
		U User
	}{
		p,
		recent,
		*u,
	}

	authenticated, _ := template.ParseFiles("cmsHome.html")
	authenticated.Execute(w, data)
}

func handleUpvote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	picId := r.FormValue("picId")
	albumId := r.FormValue("albumId")
	owner := r.FormValue("picOwner")
	cType := r.FormValue("cType")

	user := findUser(dbConnection, owner)

	var al int
	photo := Photo{}
	video := Video{}

	for i := range user.Albums {
		if user.Albums[i].AlbumId == albumId {
			al = i
			break
		}
	}

	if cType == "image" {
		var pic int

		for i := range user.Albums[al].Photo {
			if user.Albums[al].Photo[i].PhotoId == picId {
				pic = i
				break
			}
		}

		user.Albums[al].Photo[pic].Views = user.Albums[al].Photo[pic].Views + 1
		photo = user.Albums[al].Photo[pic]
	} else {
		var vid int

		for i := range user.Albums[al].Video {
			if user.Albums[al].Video[i].VideoId == picId {
				vid = i
				break
			}
		}
		user.Albums[al].Video[vid].Views = user.Albums[al].Video[vid].Views + 1
		video = user.Albums[al].Video[vid]
	}

	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"_id": user.UserId}, bson.M{"$set": bson.M{"albums": user.Albums}})
	updateTagDB(photo, video, dbConnection)
	updateMostViewed(photo, video, dbConnection)
	if err != nil {
		fmt.Println("could not update comments in tag db")
		fmt.Println(err)
		fmt.Fprintf(w, "No")
	}
	if cType == "image" {
		fmt.Fprintf(w, "Yes_"+strconv.Itoa(photo.Views))
	} else {
		fmt.Fprintf(w, "Yes_"+strconv.Itoa(video.Views))
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in login")
	r.ParseForm()
	email := r.FormValue("email")
	pass := r.FormValue("pass")
	c := find(dbConnection, email)

	if c == nil {
		fmt.Fprintf(w, "No")
	} else {
		if c.Password == pass {
			session, _ := store.Get(r, "cookie")
			session.Values["user"] = c.Id
			session.Save(r, w)
			fmt.Fprintf(w, "Yes")
		} else {
			fmt.Fprintf(w, "No")
		}
	}
}

func handleAuthenticated(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(dbConnection, currentUser)

	authenticated, _ := template.ParseFiles("pictures2.html")
	authenticated.Execute(w, u)
}

func tagAlgo(u string) string {
	grepCmd, err := exec.Command("/bin/sh", "run.sh", u).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error")
	}

	return string(grepCmd)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	if session.Values["user"] == nil {
		session.Values["user"] = ""
		session.Save(r, w)
	}

	authenticated, _ := template.ParseFiles("gmsHome.html")
	authenticated.Execute(w, session.Values["user"].(string))

}

func handleUserProfile(w http.ResponseWriter, r *http.Request) {
	u := r.URL.RawQuery
	user := findUser(dbConnection, u)

	session, _ := store.Get(r, "cookie")
	u = session.Values["user"].(string)
	currentUser := findUser(dbConnection, u)

	data := struct {
		U   User
		NEW User
	}{
		*currentUser,
		*user,
	}

	authenticated, _ := template.ParseFiles("otherUsers.html")
	authenticated.Execute(w, data)

}

func handleCreateAlbum(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	description := r.FormValue("description")

	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	c := findUser(dbConnection, currentUser)

	albumId := createAlbum(name, description, c.Email, dbConnection)

	fmt.Fprintf(w, albumId)
}

func handleSaveImage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	file, _, err := r.FormFile("uploadData")

	if err != nil {
		fmt.Println(w, err)
		fmt.Fprintf(w, "No")
		return
	}

	id := bson.NewObjectId()
	fileName := "./resources/images/userUploaded/" + id.Hex()

	dst, err := os.Create(fileName)
	defer dst.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "No")
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "No")
		return
	}

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Yes_"+fileName+"_nil_nil")
		return
	}

	x, err := exif.Decode(f)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Yes_"+fileName+"_nil_nil")
		return
	}

	if x == nil {
		fmt.Println("x is nil")
		fmt.Fprintf(w, "Yes_"+fileName+"_nil_nil")

	} else {

		lat, long, err := x.LatLong()
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Yes_"+fileName+"_nil_nil")
		} else {

			fmt.Fprintf(w, "Yes_"+fileName+"_"+strconv.FormatFloat(lat, 'f', -1, 64)+"_"+strconv.FormatFloat(long, 'f', -1, 64))
		}
	}

}

func handleLogout(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie")
	session.Values["user"] = ""
	session.Save(r, w)
	u := findUser(dbConnection, session.Values["user"].(string))

	if u == nil {
		u = &User{}
	}
	http.Redirect(w, r, "/cmsHome", http.StatusFound)
}

func checkLoggedIn(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")

	if session.Values["user"].(string) == "" {
		fmt.Fprintf(w, "No")
	} else {
		message := "Yes," + findUser(dbConnection, session.Values["user"].(string)).FirstName
		fmt.Fprintf(w, message)
	}
}

func createTagCloud(w http.ResponseWriter, r *http.Request) {
	result := getAllTags(dbConnection)
	var tags string
	var max = 0
	for tag := range result {
		if len(result[tag].Photos)+len(result[tag].Videos) > max {
			max = len(result[tag].Photos) + len(result[tag].Videos)
		}

		tags += result[tag].Name + " " + strconv.Itoa(len(result[tag].Photos)+len(result[tag].Videos)) + ","
	}
	tags += "maximum " + strconv.Itoa(max)
	fmt.Fprintf(w, tags)

}

func handleTag(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RawQuery
	tag := findByTag(dbConnection, url)
	session, _ := store.Get(r, "cookie")
	u := session.Values["user"].(string)
	currentUser := findUser(dbConnection, u)

	data := struct {
		T Tag
		U User
	}{
		*tag,
		*currentUser,
	}

	displaySameTagPhoto, _ := template.ParseFiles("taggedPictures2.html")
	displaySameTagPhoto.Execute(w, data)

}

func handleFlickr(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url1 := r.FormValue("url1")
	url2 := r.FormValue("url2")
	tag := r.FormValue("tags")
	var tags = ""

	tagList := strings.Split(tag, ",")
	for tag := range tagList {

		res, err := http.Get(url1 + tagList[tag] + url2)
		if err != nil {
			fmt.Println(err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			fmt.Println(err.Error())
		}

		resp := string(body)
		resp = resp[14 : len(resp)-1]

		var data FlickrTag
		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			fmt.Println("error when unmarshalling JSON response from Flickr" + err.Error())
		}

		for tag := 0; tag < 4; tag++ {
			tags = tags + data.Tags.Tag[tag].Content + ","
		}

	}

	if tags == "" {
		tags = tagAlgo(tag)
	}

	fmt.Fprintf(w, tags)

}

func handleRegister(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	fname := r.FormValue("first")
	lname := r.FormValue("last")
	email := r.FormValue("email")
	pass := r.FormValue("pass")

	id := bson.NewObjectId()

	albums := createDefaultAlbum(id.Hex(), fname+" "+lname, "")

	newUser := User{id, fname, lname, email, pass, "./resources/images/userUploaded/default.gif", albums, "", "", "", id.Hex()}
	add(dbConnection, newUser)

	c := find(dbConnection, email)

	if c == nil {
		fmt.Fprintf(w, "No")
	} else {

		session, _ := store.Get(r, "cookie")
		session.Values["user"] = c.Id
		session.Save(r, w)
		fmt.Fprintf(w, "Yes")
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	image := r.FormValue("imageURL")
	caption := r.FormValue("caption")
	cType := r.FormValue("contentType")
	album := r.FormValue("albumSelect")
	location := r.FormValue("location")
	lng := r.FormValue("lng")
	lat := r.FormValue("lat")
	locationN := r.FormValue("locality")
	if location == "" {
		lng = ""
		lat = ""
		locationN = ""
	}

	streetN := r.FormValue("formatted_address")
	streetN = strings.Split(streetN, ",")[0]
	tags := r.FormValue("tagList")

	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)

	c := uploadToAlbum(cType, image, caption, album, lng, lat, streetN+", "+locationN, tags, currentUser)

	authenticated, _ := template.ParseFiles("pictures2.html")
	authenticated.Execute(w, c)
}

func uploadToAlbum(cType string, filename string, caption string, album string, lng string, lat string, loc string, tags string, user string) *User {

	var location = *new(Location)
	if loc != "" && lat != "" && lng != "" {
		location = Location{loc, lat, lng}
	}

	t := make([]string, 0)
	if tags != "" {
		t = parseTags(tags, filename)
	}

	currentUser := findUser(dbConnection, user)
	id := bson.NewObjectId()
	var al int
	p := Photo{}
	v := Video{}
	for i := range currentUser.Albums {
		if currentUser.Albums[i].AlbumId == album {
			al = i
			break
		}
	}

	if cType == "image" {
		p = Photo{id, id.Hex(), currentUser.Id, currentUser.FirstName + " " + currentUser.LastName, currentUser.Albums[al].AlbumId, filename, caption, location, time.Now().Local().Format("2006-01-02"), 0, t, make([]PhotoComment, 1)}
		currentUser.Albums[al].Photo = append(currentUser.Albums[al].Photo, p)
		addTags(dbConnection, t, p, Video{})
	} else {
		v = Video{id, id.Hex(), currentUser.Id, currentUser.FirstName + " " + currentUser.LastName, currentUser.Albums[al].AlbumId, filename, caption, location, time.Now().Local().Format("2006-01-02"), 0, t, make([]PhotoComment, 1)}
		currentUser.Albums[al].Video = append(currentUser.Albums[al].Video, v)
		addTags(dbConnection, t, Photo{}, v)
	}

	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"email": currentUser.Email}, bson.M{"$set": bson.M{"albums": currentUser.Albums}})
	if err != nil {

		fmt.Println("error while trying to update in upload to album")
	}

	updateMostRecent(p, v, dbConnection)
	return currentUser

}

func parseTags(tags string, filename string) []string {
	tags = strings.ToLower(tags)
	s := strings.Split(tags, ",")

	return s
}

func handlePictures(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(dbConnection, currentUser)

	authenticated, _ := template.ParseFiles("pictures2.html")
	authenticated.Execute(w, u)

}

func handleAlbums(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery

	session, _ := store.Get(r, "cookie")
	user := session.Values["user"].(string)
	currentUser := findUser(dbConnection, user)

	if query == "" {
		authenticated, _ := template.ParseFiles("albums.html")
		authenticated.Execute(w, currentUser)
	} else {

		var al Album
		for i := range currentUser.Albums {

			if currentUser.Albums[i].AlbumId == query {
				al = currentUser.Albums[i]
				break
			}
		}

		data := struct {
			A Album
			U User
		}{
			al,
			*currentUser,
		}

		authenticated, _ := template.ParseFiles("albumDetail.html")
		authenticated.Execute(w, data)

	}

}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	u := session.Values["user"].(string)
	currentUser := findUser(dbConnection, u)

	authenticated, _ := template.ParseFiles("upload2.html")
	authenticated.Execute(w, currentUser)
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	comment := r.FormValue("comment")
	picture := r.FormValue("pic")
	album := r.FormValue("album")
	owner := r.FormValue("owner")
	cType := r.FormValue("cType")

	var user *User
	user = findUser(dbConnection, owner)

	session, _ := store.Get(r, "cookie")
	user2 := session.Values["user"].(string)

	currentUser := findUser(dbConnection, user2)
	com := PhotoComment{currentUser.FirstName + " " + currentUser.LastName, currentUser.Id, comment, time.Now().Local().Format("2006-01-02")}

	var al int
	photo := Photo{}
	video := Video{}

	for i := range user.Albums {
		if user.Albums[i].AlbumId == album {
			al = i
			break
		}
	}

	if cType == "image" {
		var pic int

		for i := range user.Albums[al].Photo {
			if user.Albums[al].Photo[i].PhotoId == picture {
				pic = i
				break
			}
		}
		user.Albums[al].Photo[pic].Comments = append(user.Albums[al].Photo[pic].Comments, com)
		photo = user.Albums[al].Photo[pic]
	} else {
		var vid int

		for i := range user.Albums[al].Video {
			if user.Albums[al].Video[i].VideoId == picture {
				vid = i
				break
			}
		}
		user.Albums[al].Video[vid].Comments = append(user.Albums[al].Video[vid].Comments, com)
		video = user.Albums[al].Video[vid]
	}

	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"_id": user.UserId}, bson.M{"$set": bson.M{"albums": user.Albums}})
	if err != nil {
		fmt.Fprintf(w, "No")
	}

	updateTagDB(photo, video, dbConnection)

	response := com.Body + "_" + com.User + "_" + com.Timestamp
	fmt.Fprintf(w, "Yes_"+response)
}
