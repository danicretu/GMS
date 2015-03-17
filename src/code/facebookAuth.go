package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	//"github.com/gorilla/sessions"
	"fmt"
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

	ClientId: "583627198440666",

	ClientSecret: "0554c954561004b6b0f646f31b6dc387",

	AuthURL: "https://www.facebook.com/dialog/oauth",

	TokenURL: "https://graph.facebook.com/oauth/access_token",

	RedirectURL: "http://mirugc.dcs.gla.ac.uk/oauth2callbackF",
}

type UserF struct {
	Id          string `json:"id"`
	Given_Name  string `json:"first_name"`
	Family_Name string `json:"last_name"`
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

	fmt.Println(string(body), "    *********************")

	json.Unmarshal(body, &user)

	fmt.Println(user, "    *********************")

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var existing *User
	sess.DB(db_name).C("user").Find(bson.M{"fId": user.Id}).One(&existing)
	session, _ := store.Get(r, "cookie")

	if existing != nil && existing.Id != "" {

		session.Values["user"] = existing.Id
		session.Save(r, w)
	} else {

		fmt.Println("in else " + db_name)
		id := bson.NewObjectId()

		newUser := User{id, user.Given_Name, user.Family_Name, "", "", "", user.Id, "", user.Id}
		add(newUser)
		createDefaultAlbum(newUser.Id, user.Given_Name+" "+user.Family_Name)

		session.Values["user"] = newUser.Id
		session.Save(r, w)

	}

	defer sess.Close()

	http.Redirect(w, r, "/authenticated", http.StatusFound)
	return

}
