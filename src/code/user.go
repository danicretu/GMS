package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rwcarlsen/goexif/exif"
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
	Downvote    int            `bson:"downvote"`
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
}

type DisplayPhotos struct {
	Name   string  `bson:"name"`
	Photos []Photo `bson:"photos"`
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

//var currentUser *User

//add(dbConnection, name, password) ->add to db
//find(dbConnection, name) ->find in db

func main() {
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/logout", handleLogout)
	router.HandleFunc("/register", handleRegister)
	router.HandleFunc("/authenticated", handleAuthenticated)
	router.HandleFunc("/pictures", handlePictures)
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
	authenticateGoogle()
	authenticateFacebook()
	authenticateTwitter()

	dbConnection = NewMongoDBConn()
	_ = dbConnection.connect()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func handleUpvote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	picId := r.FormValue("picId")
	albumId := r.FormValue("albumId")
	owner := r.FormValue("picOwner")

	user := findUser(dbConnection, owner)

	fmt.Println(picId, "*****************in handle upvote")

	var al int

	for i := range user.Albums {
		if user.Albums[i].AlbumId == albumId {
			al = i
			break
		}
	}

	var pic int

	for i := range user.Albums[al].Photo {
		if user.Albums[al].Photo[i].PhotoId == picId {
			pic = i
			break
		}
	}

	fmt.Println(user.Albums[al].Photo[pic], "*****************in handle upvote")

	user.Albums[al].Photo[pic].Views = user.Albums[al].Photo[pic].Views + 1

	fmt.Println(user.Albums[al].Photo[pic].Views, "*****************in handle upvote views")

	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"_id": user.UserId}, bson.M{"$set": bson.M{"albums": user.Albums}})
	updateTagDB(user.Albums[al].Photo[pic], dbConnection)
	updateMostViewed(user.Albums[al].Photo[pic], dbConnection)
	if err != nil {
		fmt.Println("could not update comments in tag db")
		fmt.Println(err)
		fmt.Fprintf(w, "No")
	}
	fmt.Fprintf(w, "Yes_"+strconv.Itoa(user.Albums[al].Photo[pic].Views))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
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

	authenticated, _ := template.ParseFiles("authenticated2.html")
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
	authenticated.Execute(w, session.Values["user"])

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
	logout, _ := template.ParseFiles("gmsHome.html")
	logout.Execute(w, session.Values["user"])
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
		if len(result[tag].Photos) > max {
			max = len(result[tag].Photos)
		}
		tags += result[tag].Name + " " + strconv.Itoa(len(result[tag].Photos)) + ","
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

	c := uploadToAlbum(image, caption, album, lng, lat, streetN+", "+locationN, tags, currentUser)

	authenticated, _ := template.ParseFiles("authenticated2.html")
	authenticated.Execute(w, c)
}

func uploadToAlbum(filename string, caption string, album string, lng string, lat string, loc string, tags string, user string) *User {

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
	for i := range currentUser.Albums {
		if currentUser.Albums[i].AlbumId == album {
			al = i
			break
		}
	}

	photo := Photo{id, id.Hex(), currentUser.Id, currentUser.FirstName + " " + currentUser.LastName, currentUser.Albums[al].AlbumId, filename, caption, location, time.Now().Local().Format("2006-01-02"), 0, 0, t, make([]PhotoComment, 1)}
	addTags(dbConnection, t, photo)

	currentUser.Albums[al].Photo = append(currentUser.Albums[al].Photo, photo)
	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"email": currentUser.Email}, bson.M{"$set": bson.M{"albums": currentUser.Albums}})
	if err != nil {

		fmt.Println("error while trying to update in upload to album")
	}
	updateMostRecent(photo, dbConnection)
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

	var user *User
	user = findUser(dbConnection, owner)

	session, _ := store.Get(r, "cookie")
	user2 := session.Values["user"].(string)

	var al int

	for i := range user.Albums {
		if user.Albums[i].AlbumId == album {
			al = i
			break
		}
	}

	var pic int

	for i := range user.Albums[al].Photo {
		if user.Albums[al].Photo[i].PhotoId == picture {
			pic = i
			break
		}
	}
	currentUser := findUser(dbConnection, user2)
	com := PhotoComment{currentUser.FirstName + " " + currentUser.LastName, currentUser.Id, comment, time.Now().Local().Format("2006-01-02")}

	user.Albums[al].Photo[pic].Comments = append(user.Albums[al].Photo[pic].Comments, com)
	photo := user.Albums[al].Photo[pic]
	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"_id": user.UserId}, bson.M{"$set": bson.M{"albums": user.Albums}})
	if err != nil {
		fmt.Fprintf(w, "No")
	}

	photo = user.Albums[al].Photo[pic]
	updateTagDB(photo, dbConnection)

	response := com.Body + "_" + com.User + "_" + com.Timestamp
	fmt.Fprintf(w, "Yes_"+response)
}
