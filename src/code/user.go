package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"net/http"
	"os"
)

type User struct {
	FirstName      string  `bson:"firstname"`
	LastName       string  `bson:"lastname"`
	Email          string  `bson:"email"`
	Password       string  `bson:"password"`
	ProfilePicture string  `bson:"profilepic"`
	Albums         []Album `bson:"albums"`
}

type Album struct {
	Name       string  `bson:"albumname"`
	Desciption string  `bson:"description"`
	Photo      []Photo `bson:"photos"`
}

type Photo struct {
	Owner       string    `bson:"owner"`
	URL         string    `bson:"url"`
	Description string    `bson:"description"`
	Location    Location  `bson:"location"`
	Timestamp   string    `bson:"timestamp"`
	Upvote      int       `bson:"upvote"`
	Downvote    int       `bson:"downvote"`
	Tags        []Tag     `bson:"tags"`
	Comments    []Comment `bson:"comments"`
}

type Location struct {
	Name      string `bson:"locationName"`
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type Tag struct {
	Tags []string `bson:"tags"`
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
	//http.HandleFunc("/saveComment", handleComments)
	authenticateGoogle()
	authenticateFacebook()

	dbConnection = NewMongoDBConn()
	_ = dbConnection.connect()

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	pass := r.FormValue("pass")

	fmt.Println(email, pass)

	c := find(dbConnection, email)

	if c == nil {
		authenticated, _ := template.ParseFiles("wrongCredentials.html")
		authenticated.Execute(w, c)
	} else {

		if c.Password == pass {
			currentUser = c
			authenticated, _ := template.ParseFiles("authenticated.html")
			authenticated.Execute(w, currentUser)
		}
	}

}

func handleRegister(w http.ResponseWriter, r *http.Request) {

	fname := r.FormValue("first")
	lname := r.FormValue("last")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	pass2 := r.FormValue("confirmPassword")

	albums := createDefaultAlbum()

	newUser := User{fname, lname, email, pass, "", albums}

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

	fmt.Println("heeeeeeereeeeeeeeeeeee")

	switch r.Method {
	//GET displays the upload form.
	case "GET":
		//display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//get the multipart reader for the request.
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			fileName := "./resources/images/userUploaded/" + part.FileName()

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
			if currentUser.ProfilePicture == "" {
				currentUser.ProfilePicture = fileName
				err = dbConnection.session.DB("gmsTry").C("user").Update(bson.M{"email": currentUser.Email}, bson.M{"$set": bson.M{"profilepic": fileName}})
				if err != nil {
					fmt.Println("***************")
					fmt.Println("error while trying to update")
				}
			}
			uploadToAlbum(fileName)

		}
		authenticated, _ := template.ParseFiles("authenticated.html")
		authenticated.Execute(w, currentUser)
		//display success message.
		//display(w, "upload", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func uploadToAlbum(filename string) {

	user := find(dbConnection, currentUser.Email)
	location := Location{"Glasgow", "1", "2"}

	photo := Photo{currentUser.FirstName + " " + currentUser.LastName, filename, "feet", location, "", 0, 0, make([]Tag, 1), make([]Comment, 1)}

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

/*
func handleComments(w http.ResponseWriter, r *http.Request) {
	//comment := r.FormValue("comment")
	//picture := r.FormValue("pictureNumber")

	var user *User

	user = find(dbConnection, currentUser.Email)

	fmt.Println("--------")

	//var photos []Photo
	//photos = user.Albums[0].Photo

	//var album Album
	//album = user.Albums[0]

} */

// Start the authorization process
