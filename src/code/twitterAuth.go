package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mrjones/oauth"
	"gopkg.in/mgo.v2/bson"
	//"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var tokens map[string]*oauth.RequestToken
var c *oauth.Consumer

type UserT struct {
	Id_Str string `json:"id_str"`
	Name   string `json:"name"`
}

func authenticateTwitter() {
	tokens = make(map[string]*oauth.RequestToken)

	var consumerKey *string = flag.String(
		"consumerkey",
		"HWbeqMdUyIbwIPRPUu83U42oP",
		"Consumer Key from Twitter. See: https://dev.twitter.com/apps/new")

	var consumerSecret *string = flag.String(
		"consumersecret",
		"lnKO1ncCiWaRKykbtyqbw9p2DPAB5S1dHrbNQkMZyrERxwDWA4",
		"Consumer Secret from Twitter. See: https://dev.twitter.com/apps/new")

	flag.Parse()

	c = oauth.NewConsumer(
		*consumerKey,
		*consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)
	//c.Debug(true)

	http.HandleFunc("/authorizeTwitter", RedirectUserToTwitter)
	http.HandleFunc("/oauth2callbackT", GetTwitterToken)
}

func RedirectUserToTwitter(w http.ResponseWriter, r *http.Request) {
	tokenUrl := fmt.Sprintf("http://%s/oauth2callbackT", r.Host)
	token, requestUrl, err := c.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	tokens[token.Token] = token
	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
}

func GetTwitterToken(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	accessToken, err := c.AuthorizeToken(tokens[tokenKey], verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	response, err := c.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{"count": "1"},
		accessToken)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(bits))

	var user UserT

	json.Unmarshal(bits, &user)

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var existing *User
	sess.DB(db_name).C("user").Find(bson.M{"tId": user.Id_Str}).One(&existing)
	session, _ := store.Get(r, "cookie")
	if existing != nil && existing.Id != "" {
		session.Values["user"] = existing.Id
		session.Save(r, w)
	} else {

		id := bson.NewObjectId()

		newUser := User{id, user.Name, "", "", "", "", "", user.Id_Str, user.Id_Str}
		add(newUser)
		createDefaultAlbum(newUser.Id, user.Name)

		session.Values["user"] = newUser.Id
		session.Save(r, w)

	}

	defer sess.Close()

	http.Redirect(w, r, "/authenticated", http.StatusFound)
	return
}
