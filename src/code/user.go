package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"net/http"
	"os"
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
	Id             string        `bson:"userId"`
}

type Album struct {
	Id         bson.ObjectId `bson:"_id"`
	AlbumId    string        `bson:"albumId"`
	Owner      string        `bson:"owner"`
	OwnerName  string        `bson:"ownerName"`
	Name       string        `bson:"albumname"`
	Desciption string        `bson:"description"`
	Photo      []Photo       `bson:"photos"`
}

type Photo struct {
	Id          bson.ObjectId  `bson:"_id"`
	PhotoId     string         `bson:"photoId"`
	Owner       string         `bson:"owner"`
	OwnerName   string         `bson:"ownerName"`
	URL         string         `bson:"url"`
	Description string         `bson:"description"`
	Location    Location       `bson:"location"`
	Timestamp   string         `bson:"timestamp"`
	Upvote      int            `bson:"upvote"`
	Downvote    int            `bson:"downvote"`
	Tags        []Tag          `bson:"tags"`
	Comments    []PhotoComment `bson:"comments"`
}

type Location struct {
	Name      string `bson:"locationName"`
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type Tag struct {
	Tags []string `bson:"tags"`
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

var dbConnection *MongoDBConn

var currentUser *User

//add(dbConnection, name, password) ->add to db
//find(dbConnection, name) ->find in db

func login() {

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/authenticated", handleAuthenticated)
	http.HandleFunc("/pictures", handlePictures)
	http.HandleFunc("/albums", handleAlbums)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/uploadPic", uploadHandler)
	http.HandleFunc("/saveComment", handleComments)
	http.HandleFunc("/auth", authenticate)
	authenticateGoogle()
	authenticateFacebook()

	dbConnection = NewMongoDBConn()
	_ = dbConnection.connect()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

}

func authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------------------------------------")
	authenticated, _ := template.ParseFiles("authenticated.html")
	authenticated.Execute(w, currentUser)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	fmt.Println("meeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	r.ParseForm()
	email := r.FormValue("email")
	pass := r.FormValue("pass")
	fmt.Println("meeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	fmt.Println(r)
	fmt.Println(r.Body)
	fmt.Println(email, pass, "*****************************************************")

	c := find(dbConnection, email)

	if c == nil {
		fmt.Fprintf(w, "No")
	} else {
		if c.Password == pass {
			currentUser = c
			fmt.Fprintf(w, "Yes")
		} else {
			fmt.Fprintf(w, "No")
		}
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {

	fname := r.FormValue("first")
	lname := r.FormValue("last")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	pass2 := r.FormValue("confirmPassword")

	id := bson.NewObjectId()

	albums := createDefaultAlbum(id.Hex(), fname+" "+lname, "")

	newUser := User{id, fname, lname, email, pass, albums[0].Photo[0].URL, albums, "", "", id.Hex()}

	if pass == pass2 {
		fmt.Println(email)
		add(dbConnection, newUser)

		c := find(dbConnection, email)
		currentUser = c
		authenticated, _ := template.ParseFiles("authenticated.html")
		authenticated.Execute(w, currentUser)
	}

}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("in upload go")

	switch r.Method {
	//GET displays the upload form.
	case "GET":
		//display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//get the multipart reader for the request.
		reader, err := r.MultipartReader()
		id := bson.NewObjectId()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("no uploaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaad")
			fmt.Fprintf(w, "No")
		}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}

			fileName := "./resources/images/userUploaded/" + id.Hex()

			dst, err := os.Create(fileName)
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			uploadToAlbum(fileName, id)

		}
		fmt.Println("uploaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaad")
		fmt.Fprintf(w, "YES")
		/*
			authenticated, _ := template.ParseFiles("authenticated.html")
			authenticated.Execute(w, currentUser)
			//display success message.
			//display(w, "upload", "Upload successful.") */
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func uploadToAlbum(filename string, id bson.ObjectId) {

	user := find(dbConnection, currentUser.Email)
	location := Location{"Glasgow", "1", "2"}

	photo := Photo{id, id.Hex(), currentUser.Id, currentUser.FirstName + " " + currentUser.LastName, filename, "feet", location, time.Now().Local().Format("2006-01-02"), 0, 0, make([]Tag, 1), make([]PhotoComment, 1)}

	fmt.Println(user)

	fmt.Println("***********")

	user.Albums[0].Photo = append(user.Albums[0].Photo, photo)
	currentUser.Albums[0].Photo = append(currentUser.Albums[0].Photo, photo)

	fmt.Println(user)
	fmt.Println("***********************")
	fmt.Println(currentUser)
	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"email": user.Email}, bson.M{"$set": bson.M{"albums": user.Albums}})
	if err != nil {

		fmt.Println("***************")
		fmt.Println("error while trying to update2")
	}

}

func handleAuthenticated(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := template.ParseFiles("authenticated.html")
	authenticated.Execute(w, currentUser)
}

func handlePictures(w http.ResponseWriter, r *http.Request) {

	authenticated, _ := template.ParseFiles("pictures.html")
	authenticated.Execute(w, currentUser)

}

func handleAlbums(w http.ResponseWriter, r *http.Request) {

}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := template.ParseFiles("upload.html")
	authenticated.Execute(w, currentUser)
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	picture := r.FormValue("pictureNumber")
	album := r.FormValue("albumNumber")
	owner := r.FormValue("owner")

	var user *User
	user = findUser(dbConnection, owner)

	fmt.Println(comment, picture, owner, album)

	fmt.Println(owner)
	fmt.Println(comment, picture, owner, album)
	fmt.Println(user)
	fmt.Println(comment, picture, owner, album)

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

	fmt.Println(al, pic)

	com := PhotoComment{currentUser.FirstName + " " + currentUser.LastName, currentUser.Id, comment, time.Now().Local().Format("2006-01-02")}

	fmt.Println(com)

	user.Albums[al].Photo[pic].Comments = append(user.Albums[al].Photo[pic].Comments, com)

	fmt.Println(user)
	err := dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"_id": user.UserId}, bson.M{"$set": bson.M{"albums": user.Albums}})
	if err != nil {
		panic(err)
	}

	authenticated, _ := template.ParseFiles("pictures.html")
	authenticated.Execute(w, currentUser)
}
