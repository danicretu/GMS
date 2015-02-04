package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"net/http"
)

// variables used during oauth protocol flow of authentication
var (
	codeF  = ""
	tokenF = ""
)

var oauthCfgF = &oauth.Config{

	ClientId: "853319121375165",

	ClientSecret: "36e6f25d4b121b11ace2d812c12c67af",

	AuthURL: "https://www.facebook.com/dialog/oauth",

	TokenURL: "https://graph.facebook.com/oauth/access_token",

	RedirectURL: "http://localhost:8080/oauth2callbackF",
}

type UserF struct {
	Id          string
	Name        string
	Given_Name  string
	Family_Name string
	Picture     string
	Locale      string
}

const profileInfoURLF = "https://graph.facebook.com/me?"

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
This app is now authenticated to access your Facebook user info. <img src = {{.Picture}}> <br />Your details are:<br />
{{.Name}} 
</body></html>
`))

func authenticateFacebook() {
	http.HandleFunc("/authorizeFacebook", handleAuthorize)

	//Facebook will redirect to this page to return your code, so handle it appropriately
	http.HandleFunc("/oauth2callbackF", handleOAuth2Callback)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	//Get the Facebook URL which shows the Authentication page to the user
	url := oauthCfgF.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Facebook server
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	//Get the code from the response
	code := r.FormValue("code")

	t := &oauth.Transport{Config: oauthCfgF}

	// Exchange the received code for a token
	t.Exchange(code)

	//now get user data based on the Transport which has the token
	resp, err := t.Client().Get(profileInfoURLF)

	if err != nil {
		panic(err.Error())
	}

	var user UserF

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &user)

	var existing *User
	dbConnection.session.DB("gmsTry").C("user").Find(bson.M{"fId": user.Id}).One(&existing)

	if existing != nil {
		currentUser = existing
	} else {

		id := bson.NewObjectId()
		albums := createDefaultAlbum(id.Hex(), user.Given_Name+" "+user.Family_Name, "")

		newUser := User{id, user.Given_Name, user.Family_Name, "", "", "", albums, "", user.Id, "", id.Hex()}
		add(dbConnection, newUser)
		currentUser = &newUser

	}

	authenticated, _ := template.ParseFiles("authenticated2.html")
	authenticated.Execute(w, currentUser)

}
