package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
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
	Link        string
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

	fmt.Print(url)

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

	fmt.Print(string(body))

	json.Unmarshal(body, &user)
	fmt.Println(user.Id)
	fmt.Println(user.Given_Name)
	fmt.Println(user.Family_Name)
	fmt.Println(user.Picture)
	fmt.Println(user.Locale)
	fmt.Println(user.Name)
	fmt.Println("nope")

	userInfoTemplate.Execute(w, &user)
}
