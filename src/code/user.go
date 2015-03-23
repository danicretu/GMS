package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rwcarlsen/goexif/exif"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type User struct {
	UserId     bson.ObjectId `bson:"_id"`
	FirstName  string        `bson:"firstname"`
	LastName   string        `bson:"lastname"`
	Email      string        `bson:"email"`
	Password   string        `bson:"password"`
	GoogleId   string        `bson:"gId"`
	FacebookId string        `bson:"fId"`
	TwitterId  string        `bson:"tId"`
	Id         string        `bson:"userId"`
}

type Album struct {
	Id        bson.ObjectId `bson:"_id"`
	AlbumId   string        `bson:"albumId"`
	Owner     string        `bson:"owner"`
	OwnerName string        `bson:"ownerName"`
	Name      string        `bson:"albumname"`
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

type FlickrImage struct {
	//PhotoID     string
	URL         string
	ImageName   string   `bson:"imageName"`
	Description string   `bson:"description"`
	TimeStamp   string   `bson:"timeStamp"`
	Keywords    []string `bson:"keywords"`
}

type FlickrImage1 struct {
	//PhotoID     string
	URL         string
	ImageName   string  `bson:"imageName"`
	Description string  `bson:"description"`
	TimeStamp   string  `bson:"datePosted"`
	Location    string  `bson:"exifLocation"`
	Latitude    float32 `bson:"latitude"`
	Longitude   float32 `bson:"longitude"`

	Keywords []string `bson:"keywords"`
}

type News struct {
	Title        string `bson:"title"`
	URL          string `bson:"url"`
	ImageName    string
	ImageCaption string
	ImageUrl     string
	Images       []NewsImage `bson:"images"`
}

type NewsImage struct {
	Name    string `bson:"name"`
	Caption string `bson:"caption"`
}

type Response struct {
	Name    string
	Content string
}

type AlbumStruct struct {
	Name    string
	AlbumId string
	Photo   string
}

type MapImage struct {
	Id     bson.ObjectId `bson:"_id"`
	User   bson.ObjectId `bson:"user"`
	Lat    string        `bson:"lat"`
	Lon    string        `bson:"lon"`
	Street string        `bson:"street"`
	URL    string        `bson:"url"`
	Places []string      `bson:"places"`
}

type CwgImage struct {
	Lat      string   `bson:"latitude"`
	Lon      string   `bson:"longitude"`
	Location string   `bson:"location"`
	Events   []string `bson:"events"`
	Photos   int      `bson:"photoCount"`
}

type RecommendedUser struct {
	User      string             `bson:"user"`
	Recommend []RecommendedPlace `bson:"Recommend"`
}

type RecommendedPlace struct {
	Vicinity   string `bson:"vicinity"`
	Category   string `bson:"placeCategory"`
	Popularity string `bson:"popularity"`
	Lon        string `bson:"lon"`
	Lat        string `bson:"lat"`
	Name       string `bson:"placeName"`
}

type TrendingUser struct {
	User     string          `bson:"user"`
	Trending []TrendingPlace `bson:"Trending"`
}

type TrendingPlace struct {
	Popularity string `bson:"pop"`
	Name       string `bson:"locationName"`
	URL        string `bson:"url"`
	Lon        string `bson:"lon"`
	Lat        string `bson:"lat"`
}

type TrendingAll struct {
	Lon        string `bson:"lon"`
	Lat        string `bson:"lat"`
	Loc        string `bson:"loc"`
	URL        string `bson:"url"`
	Popularity string `bson:"popularity"`
}

type FlickrStatPage struct {
	Date         string `bson:"date"`
	FlickrNumber string `bson:"totalImages"`
	FlickrSize   string `bson:"totalSize"`
}

type EachStatPage struct {
	Date      string           `bson:"date"`
	TotalSize string           `bson:"totalSize"`
	TotalNews string           `bson:"totalNews"`
	List      []SourceStatPage `bson:"list"`
}

type TwitterStatPage struct {
	Date          string `bson:"date"`
	TwitterNumber string `bson:"twitterNumber"`
	TwitterSize   string `bson:"twitterSize"`
}

type SourceStatPage struct {
	Source    string `bson:"source"`
	TotalNews string `bson:"totalNews"`
	TotalSize string `bson:"totalSize"`
}

type ScotlandStatPage struct {
	Quantity                  string
	AllTotalWholeNews         string
	AllTotalWholeSize         string
	AllTotalAverage           string
	StartDate                 string
	SelectedMonth             string
	TotalWholeNews            string
	TotalWholeSize            string
	Average                   string
	Labels                    []string
	Dates                     []string
	TotalNews                 []string
	TotalSize                 []string
	TwitterEachDayTotalTweets []string
	TwitterEachDayTotalSize   []string
	TwitterMonthTotalTweets   string
	TwitterMonthTotalSize     string
	TwitterAllTotalTweets     string
	TwitterAllTotalSize       string
	ScotTotalNews             []string
	ScotTotalSize             []string
	ScotCountTotalNews        string
	ScotCountTotalSize        string
	ScotAverage               string
	ETTotalNews               []string
	ETTotalSize               []string
	ETCountTotalNews          string
	ETCountTotalSize          string
	ETAverage                 string
	BBCTotalNews              []string
	BBCTotalSize              []string
	BBCCountTotalNews         string
	BBCCountTotalSize         string
	BBCAverage                string
	DRTotalNews               []string
	DRTotalSize               []string
	DRCountTotalNews          string
	DRCountTotalSize          string
	DRAverage                 string
	IndiTotalNews             []string
	IndiTotalSize             []string
	IndiCountTotalNews        string
	IndiCountTotalSize        string
	IndiAverage               string
	GuardTotalNews            []string
	GuardTotalSize            []string
	GuardCountTotalNews       string
	GuardCountTotalSize       string
	GuardAverage              string
	CourierTotalNews          []string
	CourierTotalSize          []string
	CourierCountTotalNews     string
	CourierCountTotalSize     string
	CourierAverage            string

	ExpressTotalNews      []string
	ExpressTotalSize      []string
	ExpressCountTotalNews string
	ExpressCountTotalSize string
	ExpressAverage        string

	EvExpressTotalNews      []string
	EvExpressTotalSize      []string
	EvExpressCountTotalNews string
	EvExpressCountTotalSize string
	EvExpressAverage        string

	FlickrEachDayTotalImages []string
	FlickrEachDayTotalSize   []string
	FlickrMonthTotalImages   string
	FlickrMonthTotalSize     string
	FlickrAllTotalImages     string
	FlickrAllTotalSize       string

	Sources []string
}

var router = mux.NewRouter()

var authKey = []byte("NCDIUyd78DBCSJBlcsd783")

// Encryption Key
var encKey = []byte("nckdajKBDSY6778FDV891bdf")

var store = sessions.NewCookieStore(authKey, encKey)

var dbConnection *MongoDBConn

var templates = template.Must(template.New("test").Funcs(funcMap).ParseFiles("index.html", "gmsFlickrStat.html", "gmsHome.html", "gmsToday.html", "stat.html", "statScotland.html", "news.html", "detailNews.html", "today.html"))

var funcMap = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"inc": func(i int) int {
		return i + 1
	},
	"mod": func(i int) int {
		return i % 2
	},
}

var sess *mgo.Session

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
	router.HandleFunc("/flickrCwg", handleFlickrNews)
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
	router.HandleFunc("/delete", handleDelete)
	router.HandleFunc("/retrieveTag", handleMainTag)
	router.HandleFunc("/retrieveUser", handleMainUser)
	router.HandleFunc("/retrieveFlickrNews", handleMainFlickr)
	router.HandleFunc("/flickrImages", handleFlickrGeneral)
	router.HandleFunc("/mapImages", handleMapImages)
	router.HandleFunc("/CWGmapImages", handleCWGMapImages)
	router.HandleFunc("/statScotland", statHandlerScotland)
	router.HandleFunc("/statRangeScotland", statScotlandHandlerByRange)
	router.HandleFunc("/resetPass", handlePassReset)

	authenticateGoogle()
	authenticateFacebook()
	authenticateTwitter()

	http.Handle("/resources/flickr/", http.StripPrefix("/resources/flickr/", http.FileServer(http.Dir("/local/imcd1/gms/flickrData"))))
	http.Handle("/resources/news/", http.StripPrefix("/resources/news/", http.FileServer(http.Dir("/local/imcd1/gms/gmsNewsImages"))))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.Handle("/", router)
	http.ListenAndServe(":8892", nil)
}

func statScotlandHandlerByRange(w http.ResponseWriter, r *http.Request) {

	startDate := r.URL.Query()["startDate"][0]
	endDate := r.URL.Query()["endDate"][0]
	fmt.Println(startDate, "-", endDate)

	tempStartDate := strings.Split(startDate, "-")
	tempEndDate := strings.Split(endDate, "-")
	fmt.Println("TempStart", tempStartDate)
	fmt.Println("temopend", tempEndDate)
	dayRange := calculateDays(tempStartDate, tempEndDate)
	fmt.Println(dayRange)

	//session, err := mgo.Dial("localhost")
	//session, err := mgo.Dial("imcdserv1.dcs.gla.ac.uk")

	/*mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"imcdserv1.dcs.gla.ac.uk"},
		Timeout:  60 * time.Second,
		Database: "gmsTry",
		Username: "gms",
		Password: "rdm$248",
	}
	//session, err := mgo.Dial("mongodb://gms:rdm$248@imcdserv1.dcs.gla.ac.uk")
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true) */

	dbConnection = NewMongoDBConn()
	session := dbConnection.connect()

	c := session.DB("gmsTry").C("gmsNewsStatScotland")

	result := []EachStatPage{}
	resultTotal := []EachStatPage{}

	AllTotalWholeNews := 0
	AllTotalWholeSize := 0
	AllTotalAverage := 0.0

	err := c.Find(bson.M{}).All(&resultTotal)
	for _, element := range resultTotal {
		//fmt.Println("index:", index)
		newsSize, err := strconv.Atoi(element.TotalSize)
		if err != nil {
			// Invalid string
			fmt.Println("Error error!")
		}

		newsNumber, err := strconv.Atoi(element.TotalNews)
		if err != nil {
			// Invalid string
			fmt.Println("Error error!")
		}

		AllTotalWholeSize += newsSize
		AllTotalWholeNews += newsNumber
	}
	AllTotalWholeSize = (AllTotalWholeSize / 1024) / 1024
	AllTotalAverage = float64(AllTotalWholeSize) / float64(AllTotalWholeNews)

	labels := []string{}
	TotalNews := []string{}
	TotalSize := []string{}

	monthCount := 0
	prevMonth := ""
	monthTotalNews := 0
	monthTotalSize := 0

	TotalWholeNews := 0
	TotalWholeSize := 0
	Average := 0.0
	Dates := []string{}

	ScotTotalNews := []string{}
	ScotTotalSize := []string{}
	ScotAverage := ""
	scotMonthTotalNews := 0
	scotMonthTotalSize := 0

	ETTotalNews := []string{}
	ETTotalSize := []string{}
	ETAverage := ""
	etMonthTotalNews := 0
	etMonthTotalSize := 0

	BBCTotalNews := []string{}
	BBCTotalSize := []string{}
	BBCAverage := ""
	bbcMonthTotalNews := 0
	bbcMonthTotalSize := 0

	DRTotalNews := []string{}
	DRTotalSize := []string{}
	DRAverage := ""
	drMonthTotalNews := 0
	drMonthTotalSize := 0

	IndiTotalNews := []string{}
	IndiTotalSize := []string{}
	IndiAverage := ""
	indiMonthTotalNews := 0
	indiMonthTotalSize := 0

	GuardTotalNews := []string{}
	GuardTotalSize := []string{}
	GuardAverage := ""
	guardMonthTotalNews := 0
	guardMonthTotalSize := 0

	CourierTotalNews := []string{}
	CourierTotalSize := []string{}
	CourierAverage := ""
	courierMonthTotalNews := 0
	courierMonthTotalSize := 0

	ExpressTotalNews := []string{}
	ExpressTotalSize := []string{}
	ExpressAverage := ""
	expressMonthTotalNews := 0
	expressMonthTotalSize := 0

	EvExpressTotalNews := []string{}
	EvExpressTotalSize := []string{}
	EvExpressAverage := ""
	evExpressMonthTotalNews := 0
	evExpressMonthTotalSize := 0
	quantity := ""
	if strings.EqualFold(tempStartDate[1], tempEndDate[1]) {
		for i := 0; i < len(dayRange); i++ {
			quantity = "KB"
			currentDate := dayRange[i]
			Dates = append(Dates, currentDate)
			labels = append(labels, currentDate)

			err = c.Find(bson.M{"date": bson.M{"$regex": currentDate, "$options": "i"}}).All(&result)

			if err != nil {
				log.Fatal(err)
			}

			if len(result) > 0 {
				tempMonthTotalNews, err := strconv.Atoi(result[0].TotalNews)
				if err != nil {
					fmt.Println("Error error!")
				}

				tempMonthTotalSize, err := strconv.Atoi(result[0].TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}

				TotalNews = append(TotalNews, strconv.Itoa(tempMonthTotalNews))
				TotalSize = append(TotalSize, strconv.Itoa((tempMonthTotalSize / 1024)))

				for k := 0; k < len(result[0].List); k++ {
					innerElement := result[0].List[k]
					if strings.Contains(innerElement.Source, "http://www.bbc.co.uk") {

						tempBBCMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempBBCMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						BBCTotalNews = append(BBCTotalNews, strconv.Itoa(tempBBCMonthTotalNews))
						BBCTotalSize = append(BBCTotalSize, strconv.Itoa((tempBBCMonthTotalSize / 1024)))

					}
					if strings.Contains(innerElement.Source, "Evening Times") {
						tempETMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempETMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}
						ETTotalNews = append(ETTotalNews, strconv.Itoa(tempETMonthTotalNews))
						ETTotalSize = append(ETTotalSize, strconv.Itoa((tempETMonthTotalSize / 1024)))
					}
					if strings.Contains(innerElement.Source, "http://www.theguardian.com") {
						tempGuardMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempGuardMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						GuardTotalNews = append(GuardTotalNews, strconv.Itoa(tempGuardMonthTotalNews))
						GuardTotalSize = append(GuardTotalSize, strconv.Itoa((tempGuardMonthTotalSize / 1024)))
					}
					if strings.Contains(innerElement.Source, "http://www.scotsman.com") {
						tempScotMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempScotMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						ScotTotalNews = append(ScotTotalNews, strconv.Itoa(tempScotMonthTotalNews))
						ScotTotalSize = append(ScotTotalSize, strconv.Itoa((tempScotMonthTotalSize / 1024)))
					}
					if strings.Contains(innerElement.Source, "http://www.dailyrecord.co.uk") {
						tempDRMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempDRMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						DRTotalNews = append(DRTotalNews, strconv.Itoa(tempDRMonthTotalNews))
						DRTotalSize = append(DRTotalSize, strconv.Itoa((tempDRMonthTotalSize / 1024)))
					}
					if strings.Contains(innerElement.Source, "http://www.independent.co.uk/") {
						tempIndiMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempIndiMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						IndiTotalNews = append(IndiTotalNews, strconv.Itoa(tempIndiMonthTotalNews))
						IndiTotalSize = append(IndiTotalSize, strconv.Itoa((tempIndiMonthTotalSize / 1024)))
					}

					if strings.Contains(innerElement.Source, "http://www.express.co.uk") {
						tempExpressMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempExpressMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						ExpressTotalNews = append(ExpressTotalNews, strconv.Itoa(tempExpressMonthTotalNews))
						ExpressTotalSize = append(ExpressTotalSize, strconv.Itoa((tempExpressMonthTotalSize / 1024)))
					}

					if strings.Contains(innerElement.Source, "http://www.eveningexpress.co.uk/") {
						tempEvExpressMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempEvExpressMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						EvExpressTotalNews = append(EvExpressTotalNews, strconv.Itoa(tempEvExpressMonthTotalNews))
						EvExpressTotalSize = append(EvExpressTotalSize, strconv.Itoa((tempEvExpressMonthTotalSize / 1024)))
					}

					if strings.Contains(innerElement.Source, "http://www.thecourier.co.uk/") {
						tempCourierMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempCourierMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						CourierTotalNews = append(CourierTotalNews, strconv.Itoa(tempCourierMonthTotalNews))
						CourierTotalSize = append(CourierTotalSize, strconv.Itoa((tempCourierMonthTotalSize / 1024)))
					}
				}
			}
		}
	} else {
		for i := 0; i < len(dayRange); i++ {
			quantity = "MB"
			currentDate := dayRange[i]
			Dates = append(Dates, currentDate)
			splitCurrentDate := strings.Split(currentDate, "/")

			currentMonth := splitCurrentDate[1]
			currentYear := splitCurrentDate[2]
			fmt.Println("previous month", prevMonth)

			if currentMonth != prevMonth {
				if monthCount != 0 {
					TotalNews = append(TotalNews, strconv.Itoa(monthTotalNews))
					TotalSize = append(TotalSize, strconv.Itoa((monthTotalSize/1024)/1024))

					BBCTotalNews = append(BBCTotalNews, strconv.Itoa(bbcMonthTotalNews))
					BBCTotalSize = append(BBCTotalSize, strconv.Itoa((bbcMonthTotalSize/1024)/1024))

					ETTotalNews = append(ETTotalNews, strconv.Itoa(etMonthTotalNews))
					ETTotalSize = append(ETTotalSize, strconv.Itoa((etMonthTotalSize/1024)/1024))

					DRTotalNews = append(DRTotalNews, strconv.Itoa(drMonthTotalNews))
					DRTotalSize = append(DRTotalSize, strconv.Itoa((drMonthTotalSize/1024)/1024))

					ScotTotalNews = append(ScotTotalNews, strconv.Itoa(scotMonthTotalNews))
					ScotTotalSize = append(ScotTotalSize, strconv.Itoa((scotMonthTotalSize/1024)/1024))

					GuardTotalNews = append(GuardTotalNews, strconv.Itoa(guardMonthTotalNews))
					GuardTotalSize = append(GuardTotalSize, strconv.Itoa((guardMonthTotalSize/1024)/1024))

					IndiTotalNews = append(IndiTotalNews, strconv.Itoa(indiMonthTotalNews))
					IndiTotalSize = append(IndiTotalSize, strconv.Itoa((indiMonthTotalSize/1024)/1024))

					CourierTotalNews = append(CourierTotalNews, strconv.Itoa(courierMonthTotalNews))
					CourierTotalSize = append(CourierTotalSize, strconv.Itoa((courierMonthTotalSize/1024)/1024))

					ExpressTotalNews = append(ExpressTotalNews, strconv.Itoa(expressMonthTotalNews))
					ExpressTotalSize = append(ExpressTotalSize, strconv.Itoa((expressMonthTotalSize/1024)/1024))

					EvExpressTotalNews = append(EvExpressTotalNews, strconv.Itoa(evExpressMonthTotalNews))
					EvExpressTotalSize = append(EvExpressTotalSize, strconv.Itoa((evExpressMonthTotalSize/1024)/1024))
				}
				monthCount++
				prevMonth = currentMonth
				monthTotalNews = 0
				monthTotalSize = 0

				expressMonthTotalNews = 0
				expressMonthTotalSize = 0

				evExpressMonthTotalNews = 0
				evExpressMonthTotalSize = 0

				courierMonthTotalNews = 0
				courierMonthTotalSize = 0

				bbcMonthTotalNews = 0
				bbcMonthTotalSize = 0

				etMonthTotalNews = 0
				etMonthTotalSize = 0

				drMonthTotalNews = 0
				drMonthTotalSize = 0

				scotMonthTotalNews = 0
				scotMonthTotalSize = 0

				guardMonthTotalNews = 0
				guardMonthTotalSize = 0

				indiMonthTotalNews = 0
				indiMonthTotalSize = 0

				labels = append(labels, currentMonth+"-"+currentYear)
			}

			err = c.Find(bson.M{"date": bson.M{"$regex": currentDate, "$options": "i"}}).All(&result)

			if err != nil {
				log.Fatal(err)
			}

			if len(result) > 0 {
				tempMonthTotalNews, err := strconv.Atoi(result[0].TotalNews)
				if err != nil {
					fmt.Println("Error error!")
				}

				tempMonthTotalSize, err := strconv.Atoi(result[0].TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				monthTotalNews += tempMonthTotalNews
				monthTotalSize += tempMonthTotalSize

				for k := 0; k < len(result[0].List); k++ {
					innerElement := result[0].List[k]
					if strings.Contains(innerElement.Source, "http://www.bbc.co.uk") {

						tempBBCMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempBBCMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						bbcMonthTotalNews += tempBBCMonthTotalNews
						bbcMonthTotalSize += tempBBCMonthTotalSize

					}
					if strings.Contains(innerElement.Source, "Evening Times") {
						tempETMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempETMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						etMonthTotalNews += tempETMonthTotalNews
						etMonthTotalSize += tempETMonthTotalSize
					}
					if strings.Contains(innerElement.Source, "http://www.theguardian.com") {
						tempGuardMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempGuardMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						guardMonthTotalNews += tempGuardMonthTotalNews
						guardMonthTotalSize += tempGuardMonthTotalSize
					}
					if strings.Contains(innerElement.Source, "http://www.scotsman.com") {
						tempScotMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempScotMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						scotMonthTotalNews += tempScotMonthTotalNews
						scotMonthTotalSize += tempScotMonthTotalSize
					}
					if strings.Contains(innerElement.Source, "http://www.dailyrecord.co.uk") {
						tempDRMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempDRMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						drMonthTotalNews += tempDRMonthTotalNews
						drMonthTotalSize += tempDRMonthTotalSize
					}
					if strings.Contains(innerElement.Source, "http://www.independent.co.uk/") {
						tempIndiMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempIndiMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						indiMonthTotalNews += tempIndiMonthTotalNews
						indiMonthTotalSize += tempIndiMonthTotalSize
					}

					if strings.Contains(innerElement.Source, "http://www.express.co.uk") {
						tempExpressMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempExpressMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						expressMonthTotalNews += tempExpressMonthTotalNews
						expressMonthTotalSize += tempExpressMonthTotalSize
					}

					if strings.Contains(innerElement.Source, "http://www.eveningexpress.co.uk/") {
						tempEvExpressMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempEvExpressMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						evExpressMonthTotalNews += tempEvExpressMonthTotalNews
						evExpressMonthTotalSize += tempEvExpressMonthTotalSize
					}

					if strings.Contains(innerElement.Source, "http://www.thecourier.co.uk/") {
						tempCourierMonthTotalNews, err := strconv.Atoi(innerElement.TotalNews)
						if err != nil {
							fmt.Println("Error error!")
						}

						tempCourierMonthTotalSize, err := strconv.Atoi(innerElement.TotalSize)
						if err != nil {
							fmt.Println("Error error!")
						}

						courierMonthTotalNews += tempCourierMonthTotalNews
						courierMonthTotalSize += tempCourierMonthTotalSize
					}
				}
			}

			if strings.EqualFold(currentDate, (tempEndDate[2] + "/" + tempEndDate[1] + "/" + tempEndDate[0])) {
				TotalNews = append(TotalNews, strconv.Itoa(monthTotalNews))

				TotalSize = append(TotalSize, strconv.Itoa((monthTotalSize/1024)/1024))

				BBCTotalNews = append(BBCTotalNews, strconv.Itoa(bbcMonthTotalNews))
				BBCTotalSize = append(BBCTotalSize, strconv.Itoa((bbcMonthTotalSize/1024)/1024))

				ETTotalNews = append(ETTotalNews, strconv.Itoa(etMonthTotalNews))
				ETTotalSize = append(ETTotalSize, strconv.Itoa((etMonthTotalSize/1024)/1024))

				DRTotalNews = append(DRTotalNews, strconv.Itoa(drMonthTotalNews))
				DRTotalSize = append(DRTotalSize, strconv.Itoa((drMonthTotalSize/1024)/1024))

				ScotTotalNews = append(ScotTotalNews, strconv.Itoa(scotMonthTotalNews))
				ScotTotalSize = append(ScotTotalSize, strconv.Itoa((scotMonthTotalSize/1024)/1024))

				GuardTotalNews = append(GuardTotalNews, strconv.Itoa(guardMonthTotalNews))
				GuardTotalSize = append(GuardTotalSize, strconv.Itoa((guardMonthTotalSize/1024)/1024))

				IndiTotalNews = append(IndiTotalNews, strconv.Itoa(indiMonthTotalNews))
				IndiTotalSize = append(IndiTotalSize, strconv.Itoa((indiMonthTotalSize/1024)/1024))

				CourierTotalNews = append(CourierTotalNews, strconv.Itoa(courierMonthTotalNews))
				CourierTotalSize = append(CourierTotalSize, strconv.Itoa((courierMonthTotalSize/1024)/1024))

				ExpressTotalNews = append(ExpressTotalNews, strconv.Itoa(expressMonthTotalNews))
				ExpressTotalSize = append(ExpressTotalSize, strconv.Itoa((expressMonthTotalSize/1024)/1024))

				EvExpressTotalNews = append(EvExpressTotalNews, strconv.Itoa(evExpressMonthTotalNews))
				EvExpressTotalSize = append(EvExpressTotalSize, strconv.Itoa((evExpressMonthTotalSize/1024)/1024))
			}

		}
	}

	for _, value := range TotalNews {
		newsNumber, err := strconv.Atoi(value)
		if err != nil {
			// Invalid string
			fmt.Println("Error error!")
		}
		TotalWholeNews += newsNumber
	}
	for _, value := range TotalSize {
		newsSize, err := strconv.Atoi(value)
		if err != nil {
			// Invalid string
			fmt.Println("Error error!")
		}
		TotalWholeSize += newsSize
	}

	SelectedMonth := tempStartDate[2] + "/" + tempStartDate[1] + "/" + tempStartDate[0] + " - " + tempEndDate[2] + "/" + tempEndDate[1] + "/" + tempEndDate[0]
	//fmt.Println("Labels", labels)
	//fmt.Println("monthTotalNews:", TotalNews)
	//fmt.Println("monthTotalSize:", TotalSize)

	fmt.Println("Labels", labels)
	fmt.Println("Express monthTotalNews:", ETTotalNews)
	fmt.Println("Express monthTotalSize:", ETTotalSize)

	BBCCountTotalSize := 0
	BBCCountTotalNews := 0

	for _, val := range BBCTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		BBCCountTotalSize += size
	}

	for _, val := range BBCTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		BBCCountTotalNews += news
	}

	BBCAverage = strconv.FormatFloat(float64(BBCCountTotalSize)/float64(BBCCountTotalNews), 'f', 2, 64)

	ScotCountTotalSize := 0
	ScotCountTotalNews := 0

	for _, val := range ScotTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ScotCountTotalSize += size
	}

	for _, val := range ScotTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ScotCountTotalNews += news
	}

	ScotAverage = strconv.FormatFloat(float64(ScotCountTotalSize)/float64(ScotCountTotalNews), 'f', 2, 64)

	DRCountTotalSize := 0
	DRCountTotalNews := 0

	for _, val := range DRTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		DRCountTotalSize += size
	}

	for _, val := range DRTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		DRCountTotalNews += news
	}

	DRAverage = strconv.FormatFloat(float64(DRCountTotalSize)/float64(DRCountTotalNews), 'f', 2, 64)

	ETCountTotalSize := 0
	ETCountTotalNews := 0

	for _, val := range ETTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ETCountTotalSize += size
	}

	for _, val := range ETTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ETCountTotalNews += news
	}

	ETAverage = strconv.FormatFloat(float64(ETCountTotalSize)/float64(ETCountTotalNews), 'f', 2, 64)

	GuardCountTotalSize := 0
	GuardCountTotalNews := 0

	for _, val := range GuardTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		GuardCountTotalSize += size
	}

	for _, val := range GuardTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		GuardCountTotalNews += news
	}

	GuardAverage = strconv.FormatFloat(float64(GuardCountTotalSize)/float64(GuardCountTotalNews), 'f', 2, 64)

	IndiCountTotalSize := 0
	IndiCountTotalNews := 0

	for _, val := range IndiTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		IndiCountTotalSize += size
	}

	for _, val := range IndiTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		IndiCountTotalNews += news
	}

	IndiAverage = strconv.FormatFloat(float64(IndiCountTotalSize)/float64(IndiCountTotalNews), 'f', 2, 64)

	CourierCountTotalSize := 0
	CourierCountTotalNews := 0

	for _, val := range CourierTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		CourierCountTotalSize += size
	}

	for _, val := range CourierTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		CourierCountTotalNews += news
	}

	CourierAverage = strconv.FormatFloat(float64(CourierCountTotalSize)/float64(CourierCountTotalNews), 'f', 2, 64)

	ExpressCountTotalSize := 0
	ExpressCountTotalNews := 0

	for _, val := range ExpressTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ExpressCountTotalSize += size
	}

	for _, val := range ExpressTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ExpressCountTotalNews += news
	}

	ExpressAverage = strconv.FormatFloat(float64(ExpressCountTotalSize)/float64(ExpressCountTotalNews), 'f', 2, 64)

	EvExpressCountTotalSize := 0
	EvExpressCountTotalNews := 0

	for _, val := range EvExpressTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		EvExpressCountTotalSize += size
	}

	for _, val := range EvExpressTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		EvExpressCountTotalNews += news
	}

	TwitterEachDayTotalTweets := make([]string, 31)
	TwitterEachDayTotalSize := make([]string, 31)
	TwitterMonthTotalTweets := 0
	TwitterMonthTotalSize := 0

	TwitterAllTotalTweets := 0
	TwitterAllTotalSize := 0
	prevMonth = ""
	monthCount = 0
	FlickrEachDayTotalImages := []string{}
	FlickrEachDayTotalSize := []string{}
	FlickrMonthTotalImages := 0
	FlickrMonthTotalSize := 0
	FlickrAllTotalImages := 0
	FlickrAllTotalSize := 0

	labels = nil

	c = session.DB("gmsTry").C("statTest")

	flickrResult := []FlickrStatPage{}
	err = c.Find(bson.M{}).All(&flickrResult)

	for _, element := range flickrResult {
		//fmt.Println("index:", index)
		imageSize, err := strconv.Atoi(element.FlickrSize)
		if err != nil {
			// Invalid string
			fmt.Println("Error error!")
		}

		imageNumber, err := strconv.Atoi(element.FlickrNumber)
		if err != nil {
			// Invalid string
			fmt.Println("Error error!")
		}

		FlickrAllTotalImages += imageNumber
		FlickrAllTotalSize += imageSize
	}
	//AllTotalWholeSize = (AllTotalWholeSize / 1024) / 1024
	//AllTotalAverage = float64(AllTotalWholeSize) / float64(AllTotalWholeNews)

	fmt.Println("Total images", FlickrAllTotalImages)
	fmt.Println("Total size", FlickrAllTotalSize)

	if strings.EqualFold(tempStartDate[1], tempEndDate[1]) {
		for i := 0; i < len(dayRange); i++ {
			quantity = "KB"
			currentDate := dayRange[i]
			Dates = append(Dates, currentDate)
			labels = append(labels, currentDate)
			err = c.Find(bson.M{"date": bson.M{"$regex": currentDate, "$options": "i"}}).All(&flickrResult)
			fmt.Println("label sert 2", labels)
			if err != nil {
				log.Fatal(err)
			}

			if len(flickrResult) > 0 {
				FlickrMonthTotalImages, err := strconv.Atoi(flickrResult[0].FlickrNumber)
				if err != nil {
					fmt.Println("Error error!")
				}

				FlickrMonthTotalSize, err := strconv.Atoi(flickrResult[0].FlickrSize)
				if err != nil {
					fmt.Println("Error error!")
				}

				FlickrEachDayTotalImages = append(FlickrEachDayTotalImages, strconv.Itoa(FlickrMonthTotalImages))
				fmt.Println("number", FlickrEachDayTotalImages)
				FlickrEachDayTotalSize = append(FlickrEachDayTotalSize, strconv.Itoa(FlickrMonthTotalSize))
				fmt.Println("size", FlickrEachDayTotalSize)
			}
		}
	} else {
		for i := 0; i < len(dayRange); i++ {
			quantity = "MB"
			currentDate := dayRange[i]
			Dates = append(Dates, currentDate)
			splitCurrentDate := strings.Split(currentDate, "/")
			currentMonth := splitCurrentDate[1]
			currentYear := splitCurrentDate[2]
			fmt.Println("current month", currentMonth)
			fmt.Println("previous month", prevMonth)
			if currentMonth != prevMonth {
				if monthCount != 0 {
					FlickrEachDayTotalImages = append(FlickrEachDayTotalImages, strconv.Itoa(FlickrMonthTotalImages))
					FlickrEachDayTotalSize = append(FlickrEachDayTotalSize, strconv.Itoa(FlickrMonthTotalSize))
					fmt.Println("xx", FlickrEachDayTotalImages)
					fmt.Println("yy", FlickrEachDayTotalSize)
				}

				monthCount++
				prevMonth = currentMonth
				FlickrMonthTotalImages = 0
				FlickrMonthTotalSize = 0

				labels = append(labels, currentMonth+"-"+currentYear)
			}

			err = c.Find(bson.M{"date": bson.M{"$regex": currentDate, "$options": "i"}}).All(&flickrResult)

			if err != nil {
				log.Fatal(err)
			}

			if len(flickrResult) > 0 {
				tempMonthTotalImages, err := strconv.Atoi(flickrResult[0].FlickrNumber)
				if err != nil {
					fmt.Println("Error error!")
				}

				tempMonthTotalSizeF, err := strconv.Atoi(flickrResult[0].FlickrSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				FlickrMonthTotalImages += tempMonthTotalImages
				FlickrMonthTotalSize += tempMonthTotalSizeF
				fmt.Println("xx1", FlickrEachDayTotalImages)
				fmt.Println("yy1", FlickrEachDayTotalSize)

			}

			if strings.EqualFold(currentDate, (tempEndDate[2] + "/" + tempEndDate[1] + "/" + tempEndDate[0])) {

				FlickrEachDayTotalImages = append(FlickrEachDayTotalImages, strconv.Itoa(FlickrMonthTotalImages))
				FlickrEachDayTotalSize = append(FlickrEachDayTotalSize, strconv.Itoa(FlickrMonthTotalSize))

			}
		}

	}

	fmt.Println("Labels", labels)
	fmt.Println("monthTotalNews:", FlickrEachDayTotalImages)
	fmt.Println("monthTotalSize:", FlickrEachDayTotalSize)

	EvExpressAverage = strconv.FormatFloat(float64(EvExpressCountTotalSize)/float64(EvExpressCountTotalNews), 'f', 2, 64)
	Average = float64(TotalWholeSize) / float64(TotalWholeNews)

	Sources := []string{"BBC", "Courier", "Express", "Guardian", "Scotsman", "Independent", "Daily Record", "Evening Times", "Evening Express"}

	statPage := ScotlandStatPage{quantity, strconv.Itoa(AllTotalWholeNews), strconv.Itoa(AllTotalWholeSize), strconv.FormatFloat(AllTotalAverage, 'f', 2, 64), startDate, SelectedMonth,
		strconv.Itoa(TotalWholeNews), strconv.Itoa(TotalWholeSize), strconv.FormatFloat(Average, 'f', 2, 64), labels, Dates, TotalNews, TotalSize,
		TwitterEachDayTotalTweets, TwitterEachDayTotalSize, strconv.Itoa(TwitterMonthTotalTweets), strconv.Itoa(TwitterMonthTotalSize), strconv.Itoa(TwitterAllTotalTweets), strconv.Itoa(TwitterAllTotalSize),
		ScotTotalNews, ScotTotalSize, strconv.Itoa(ScotCountTotalNews), strconv.Itoa(ScotCountTotalSize), ScotAverage,
		ETTotalNews, ETTotalSize, strconv.Itoa(ETCountTotalNews), strconv.Itoa(ETCountTotalSize), ETAverage,
		BBCTotalNews, BBCTotalSize, strconv.Itoa(BBCCountTotalNews), strconv.Itoa(BBCCountTotalSize), BBCAverage,
		DRTotalNews, DRTotalSize, strconv.Itoa(DRCountTotalNews), strconv.Itoa(DRCountTotalSize), DRAverage,
		IndiTotalNews, IndiTotalSize, strconv.Itoa(IndiCountTotalNews), strconv.Itoa(IndiCountTotalSize), IndiAverage,
		GuardTotalNews, GuardTotalSize, strconv.Itoa(GuardCountTotalNews), strconv.Itoa(GuardCountTotalSize), GuardAverage,
		CourierTotalNews, CourierTotalSize, strconv.Itoa(CourierCountTotalNews), strconv.Itoa(CourierCountTotalSize), CourierAverage,
		ExpressTotalNews, ExpressTotalSize, strconv.Itoa(ExpressCountTotalNews), strconv.Itoa(ExpressCountTotalSize), ExpressAverage,
		EvExpressTotalNews, EvExpressTotalSize, strconv.Itoa(EvExpressCountTotalNews), strconv.Itoa(EvExpressCountTotalSize), EvExpressAverage,
		FlickrEachDayTotalImages, FlickrEachDayTotalSize, strconv.Itoa(FlickrMonthTotalImages), strconv.Itoa(FlickrMonthTotalSize), strconv.Itoa(FlickrAllTotalImages), strconv.Itoa(FlickrAllTotalSize),
		Sources}

	defer session.Close()

	renderStatScotlandTemplate(w, "statScotland", &statPage)
}

func calculateDays(startDate []string, endDate []string) []string {
	startMonth, err := strconv.Atoi(startDate[1])
	if err != nil {
		// Invalid string
		fmt.Println("Error error!")
	}

	endMonth, err := strconv.Atoi(endDate[1])
	if err != nil {
		// Invalid string
		fmt.Println("Error error!")
	}

	startDay, err := strconv.Atoi(startDate[2])
	if err != nil {
		// Invalid string
		fmt.Println("Error error!")
	}

	endDay, err := strconv.Atoi(endDate[2])
	if err != nil {
		// Invalid string
		fmt.Println("Error error!")
	}

	startYear, err := strconv.Atoi(startDate[0])
	if err != nil {
		// Invalid string
		fmt.Println("Error error!")
	}

	endYear, err := strconv.Atoi(endDate[0])
	if err != nil {
		// Invalid string
		fmt.Println("Error error!")
	}

	currentDay := startDay
	currentMonth := startMonth
	currentYear := startYear

	dayRange := []string{}
	for {
		if currentDay < 10 {
			if currentMonth < 10 {
				dayRange = append(dayRange, "0"+strconv.Itoa(currentDay)+"/"+"0"+strconv.Itoa(currentMonth)+"/"+strconv.Itoa(currentYear))
			} else {
				dayRange = append(dayRange, "0"+strconv.Itoa(currentDay)+"/"+strconv.Itoa(currentMonth)+"/"+strconv.Itoa(currentYear))
			}
		} else {
			if currentMonth < 10 {
				dayRange = append(dayRange, strconv.Itoa(currentDay)+"/"+"0"+strconv.Itoa(currentMonth)+"/"+strconv.Itoa(currentYear))
			} else {
				dayRange = append(dayRange, strconv.Itoa(currentDay)+"/"+strconv.Itoa(currentMonth)+"/"+strconv.Itoa(currentYear))
			}
		}

		if currentDay == endDay && currentMonth == endMonth && currentYear == endYear {
			break
		}
		currentDay++
		if currentMonth == 1 && currentDay == 32 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 2 && isLeap(currentYear) && currentDay == 30 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 2 && !isLeap(currentYear) && currentDay == 29 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 3 && currentDay == 32 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 4 && currentDay == 31 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 5 && currentDay == 32 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 6 && currentDay == 31 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 7 && currentDay == 32 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 8 && currentDay == 32 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 9 && currentDay == 31 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 10 && currentDay == 32 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 11 && currentDay == 31 {
			currentDay = 1
			currentMonth++
		}
		if currentMonth == 12 && currentDay == 32 {
			currentDay = 1
			currentMonth = 1
			currentYear++
		}
	}
	return dayRange
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func statHandlerScotland(w http.ResponseWriter, r *http.Request) {
	flag := false
	for key := range r.URL.Query() {
		if strings.HasPrefix(key, "month") {
			flag = true
		}
	}

	//session, err := mgo.Dial("localhost")
	//session, err := mgo.Dial("imcdserv1.dcs.gla.ac.uk")
	/*mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"imcdserv1.dcs.gla.ac.uk"},
		Timeout:  60 * time.Second,
		Database: "gmsTry",
		Username: "gms",
		Password: "rdm$248",
	}
	//session, err := mgo.Dial("mongodb://gms:rdm$248@imcdserv1.dcs.gla.ac.uk")
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	} */

	//defer session.Close()
	//currenttime.Format("02/01/2006")
	// Optional. Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)

	dbConnection = NewMongoDBConn()
	session := dbConnection.connect()

	c := session.DB("gmsTry").C("gmsNewsStatScotland")

	result := []EachStatPage{}
	resultTotal := []EachStatPage{}
	currentMonth := []string{}
	StartDate := ""
	if flag {
		cMonth := "" + r.URL.Query()["month"][0]
		currentMonth = strings.Split(cMonth, ",")
		StartDate = currentMonth[1] + "-" + currentMonth[0] + "-" + "01"
	} else {
		currentTime := time.Now().Local().Format("02/01/2006")

		currentMonth = strings.Split(currentTime, "/")
		StartDate = currentMonth[2] + "-" + currentMonth[1] + "-" + currentMonth[0]
	}
	labels := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	TotalNews := make([]string, 31)
	TotalSize := make([]string, 31)

	ScotTotalNews := make([]string, 31)
	ScotTotalSize := make([]string, 31)
	ScotAverage := ""

	ETTotalNews := make([]string, 31)
	ETTotalSize := make([]string, 31)
	ETAverage := ""

	BBCTotalNews := make([]string, 31)
	BBCTotalSize := make([]string, 31)
	BBCAverage := ""

	DRTotalNews := make([]string, 31)
	DRTotalSize := make([]string, 31)
	DRAverage := ""

	IndiTotalNews := make([]string, 31)
	IndiTotalSize := make([]string, 31)
	IndiAverage := ""

	GuardTotalNews := make([]string, 31)
	GuardTotalSize := make([]string, 31)
	GuardAverage := ""

	CourierTotalNews := make([]string, 31)
	CourierTotalSize := make([]string, 31)
	CourierAverage := ""

	ExpressTotalNews := make([]string, 31)
	ExpressTotalSize := make([]string, 31)
	ExpressAverage := ""

	EvExpressTotalNews := make([]string, 31)
	EvExpressTotalSize := make([]string, 31)
	EvExpressAverage := ""

	AllTotalWholeNews := 0
	AllTotalWholeSize := 0
	AllTotalAverage := 0.0

	err := c.Find(bson.M{}).All(&resultTotal)
	for _, element := range resultTotal {
		newsSize, err := strconv.Atoi(element.TotalSize)
		if err != nil {
			fmt.Println("Error error!")
		}

		newsNumber, err := strconv.Atoi(element.TotalNews)
		if err != nil {
			fmt.Println("Error error!")
		}

		AllTotalWholeSize += newsSize
		AllTotalWholeNews += newsNumber
	}

	AllTotalWholeSize = (AllTotalWholeSize / 1024) / 1024
	AllTotalAverage = float64(AllTotalWholeSize) / float64(AllTotalWholeNews)

	if flag {
		err = c.Find(bson.M{"date": bson.M{"$regex": "/" + currentMonth[0] + "/" + currentMonth[1], "$options": "i"}}).All(&result)
	} else {
		err = c.Find(bson.M{"date": bson.M{"$regex": "/" + currentMonth[1] + "/" + currentMonth[2], "$options": "i"}}).All(&result)
	}

	if err != nil {
		log.Fatal(err)
	}
	Dates := []string{}

	TotalWholeNews := 0
	TotalWholeSize := 0
	Average := 0.0
	SelectedMonth := ""
	SelectedMonthIndex := 0
	for index, element := range result {
		timeStamp := element.Date
		s := strings.Split(timeStamp, " ")
		timeStamp = s[0]
		s = strings.Split(timeStamp, "/")
		timeStamp = s[2] + "-" + s[1] + "-" + s[0]
		result[index].Date = timeStamp
		Dates = append(Dates, timeStamp)

		dayIndex, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("Error error!")
		}

		newsSize, err := strconv.Atoi(element.TotalSize)
		if err != nil {
			fmt.Println("Error error!")
		}

		newsNumber, err := strconv.Atoi(element.TotalNews)
		if err != nil {
			fmt.Println("Error error!")
		}

		TotalWholeSize += (newsSize / 1024) / 1024
		TotalWholeNews += newsNumber
		MonthName := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October",
			"November", "December"}
		prevMonthIndex := ""
		if flag {
			prevMonthIndex = currentMonth[0]
		} else {
			prevMonthIndex = currentMonth[1]
		}
		monthIndex, err := strconv.Atoi(prevMonthIndex)
		if err != nil {
			fmt.Println("Error error!")
		}

		SelectedMonth = MonthName[monthIndex-1]
		SelectedMonthIndex = monthIndex

		TotalNews[dayIndex-1] = element.TotalNews
		TotalSize[dayIndex-1] = strconv.Itoa((newsSize / 1024) / 1024)

		for k := 0; k < len(element.List); k++ {
			innerElement := element.List[k]
			if strings.Contains(innerElement.Source, "http://www.bbc.co.uk") {
				BBCTotalNews[dayIndex-1] = innerElement.TotalNews

				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					// Invalid string
					fmt.Println("Error error!")
				}
				BBCTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}
			if strings.Contains(innerElement.Source, "Evening Times") {
				ETTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					// Invalid string
					fmt.Println("Error error!")
				}
				ETTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}
			if strings.Contains(innerElement.Source, "http://www.theguardian.com") {
				GuardTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				GuardTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}
			if strings.Contains(innerElement.Source, "http://www.scotsman.com") {
				ScotTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				ScotTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}
			if strings.Contains(innerElement.Source, "http://www.dailyrecord.co.uk") {
				DRTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				DRTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}
			if strings.Contains(innerElement.Source, "http://www.independent.co.uk/") {
				IndiTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				IndiTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}

			if strings.Contains(innerElement.Source, "http://www.thecourier.co.uk/") {
				CourierTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				CourierTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}

			if strings.Contains(innerElement.Source, "http://www.express.co.uk") {
				ExpressTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				ExpressTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}

			if strings.Contains(innerElement.Source, "http://www.eveningexpress.co.uk/") {
				EvExpressTotalNews[dayIndex-1] = innerElement.TotalNews
				size, err := strconv.Atoi(innerElement.TotalSize)
				if err != nil {
					fmt.Println("Error error!")
				}
				EvExpressTotalSize[dayIndex-1] = strconv.Itoa(size / 1024)
			}
		}
	}

	for l := len(result); l < 31; l++ {
		TotalNews[l] = "0"
		TotalSize[l] = "0"

		ScotTotalNews[l] = "0"
		ScotTotalSize[l] = "0"

		ETTotalNews[l] = "0"
		ETTotalSize[l] = "0"

		BBCTotalNews[l] = "0"
		BBCTotalSize[l] = "0"

		DRTotalNews[l] = "0"
		DRTotalSize[l] = "0"

		IndiTotalNews[l] = "0"
		IndiTotalSize[l] = "0"

		GuardTotalNews[l] = "0"
		GuardTotalSize[l] = "0"

		CourierTotalNews[l] = "0"
		CourierTotalSize[l] = "0"

		ExpressTotalNews[l] = "0"
		ExpressTotalSize[l] = "0"

		EvExpressTotalNews[l] = "0"
		EvExpressTotalSize[l] = "0"
	}

	for index, val := range TotalNews {
		if strings.EqualFold(val, "") {
			TotalNews[index] = "0"
		}
	}

	for index2, val2 := range TotalSize {
		if strings.EqualFold(val2, "") {
			TotalSize[index2] = "0"
		}
	}

	for index, val := range BBCTotalNews {
		if strings.EqualFold(val, "") {
			BBCTotalNews[index] = "0"
		}
	}

	for index2, val2 := range BBCTotalSize {
		if strings.EqualFold(val2, "") {
			BBCTotalSize[index2] = "0"
		}
	}

	for index, val := range ScotTotalNews {
		if strings.EqualFold(val, "") {
			ScotTotalNews[index] = "0"
		}
	}

	for index2, val2 := range ScotTotalSize {
		if strings.EqualFold(val2, "") {
			ScotTotalSize[index2] = "0"
		}
	}

	for index, val := range ETTotalNews {
		if strings.EqualFold(val, "") {
			ETTotalNews[index] = "0"
		}
	}

	for index2, val2 := range ETTotalSize {
		if strings.EqualFold(val2, "") {
			ETTotalSize[index2] = "0"
		}
	}

	for index, val := range DRTotalNews {
		if strings.EqualFold(val, "") {
			DRTotalNews[index] = "0"
		}
	}

	for index2, val2 := range DRTotalSize {
		if strings.EqualFold(val2, "") {
			DRTotalSize[index2] = "0"
		}
	}

	for index, val := range GuardTotalNews {
		if strings.EqualFold(val, "") {
			GuardTotalNews[index] = "0"
		}
	}

	for index2, val2 := range GuardTotalSize {
		if strings.EqualFold(val2, "") {
			GuardTotalSize[index2] = "0"
		}
	}

	for index, val := range IndiTotalNews {
		if strings.EqualFold(val, "") {
			IndiTotalNews[index] = "0"
		}
	}

	for index2, val2 := range IndiTotalSize {
		if strings.EqualFold(val2, "") {
			IndiTotalSize[index2] = "0"
		}
	}

	for index, val := range CourierTotalNews {
		if strings.EqualFold(val, "") {
			CourierTotalNews[index] = "0"
		}
	}

	for index2, val2 := range CourierTotalSize {
		if strings.EqualFold(val2, "") {
			CourierTotalSize[index2] = "0"
		}
	}

	for index, val := range ExpressTotalNews {
		if strings.EqualFold(val, "") {
			ExpressTotalNews[index] = "0"
		}
	}

	for index2, val2 := range ExpressTotalSize {
		if strings.EqualFold(val2, "") {
			ExpressTotalSize[index2] = "0"
		}
	}

	for index, val := range EvExpressTotalNews {
		if strings.EqualFold(val, "") {
			EvExpressTotalNews[index] = "0"
		}
	}

	for index2, val2 := range EvExpressTotalSize {
		if strings.EqualFold(val2, "") {
			EvExpressTotalSize[index2] = "0"
		}
	}

	BBCCountTotalSize := 0
	BBCCountTotalNews := 0

	for _, val := range BBCTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		BBCCountTotalSize += size
	}

	for _, val := range BBCTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		BBCCountTotalNews += news
	}

	BBCAverage = strconv.FormatFloat(float64(BBCCountTotalSize)/float64(BBCCountTotalNews), 'f', 2, 64)

	ScotCountTotalSize := 0
	ScotCountTotalNews := 0

	for _, val := range ScotTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ScotCountTotalSize += size
	}

	for _, val := range ScotTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ScotCountTotalNews += news
	}

	ScotAverage = strconv.FormatFloat(float64(ScotCountTotalSize)/float64(ScotCountTotalNews), 'f', 2, 64)

	DRCountTotalSize := 0
	DRCountTotalNews := 0

	for _, val := range DRTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		DRCountTotalSize += size
	}

	for _, val := range DRTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		DRCountTotalNews += news
	}

	DRAverage = strconv.FormatFloat(float64(DRCountTotalSize)/float64(DRCountTotalNews), 'f', 2, 64)

	ETCountTotalSize := 0
	ETCountTotalNews := 0

	for _, val := range ETTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ETCountTotalSize += size
	}

	for _, val := range ETTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ETCountTotalNews += news
	}

	ETAverage = strconv.FormatFloat(float64(ETCountTotalSize)/float64(ETCountTotalNews), 'f', 2, 64)

	GuardCountTotalSize := 0
	GuardCountTotalNews := 0

	for _, val := range GuardTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		GuardCountTotalSize += size
	}

	for _, val := range GuardTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		GuardCountTotalNews += news
	}

	GuardAverage = strconv.FormatFloat(float64(GuardCountTotalSize)/float64(GuardCountTotalNews), 'f', 2, 64)

	IndiCountTotalSize := 0
	IndiCountTotalNews := 0

	for _, val := range IndiTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		IndiCountTotalSize += size
	}

	for _, val := range IndiTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		IndiCountTotalNews += news
	}

	IndiAverage = strconv.FormatFloat(float64(IndiCountTotalSize)/float64(IndiCountTotalNews), 'f', 2, 64)

	CourierCountTotalSize := 0
	CourierCountTotalNews := 0

	for _, val := range CourierTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		CourierCountTotalSize += size
	}

	for _, val := range CourierTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		CourierCountTotalNews += news
	}

	CourierAverage = strconv.FormatFloat(float64(CourierCountTotalSize)/float64(CourierCountTotalNews), 'f', 2, 64)

	ExpressCountTotalSize := 0
	ExpressCountTotalNews := 0

	for _, val := range ExpressTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ExpressCountTotalSize += size
	}

	for _, val := range ExpressTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		ExpressCountTotalNews += news
	}

	ExpressAverage = strconv.FormatFloat(float64(ExpressCountTotalSize)/float64(ExpressCountTotalNews), 'f', 2, 64)

	EvExpressCountTotalSize := 0
	EvExpressCountTotalNews := 0

	for _, val := range EvExpressTotalSize {
		size, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		EvExpressCountTotalSize += size
	}

	for _, val := range EvExpressTotalNews {
		news, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error error!")
		}
		EvExpressCountTotalNews += news
	}

	TwitterEachDayTotalTweets := make([]string, 31)
	TwitterEachDayTotalSize := make([]string, 31)
	TwitterMonthTotalTweets := 0
	TwitterMonthTotalSize := 0

	TwitterAllTotalTweets := 0
	TwitterAllTotalSize := 0

	for i := 0; i < 31; i++ {
		TwitterEachDayTotalTweets[i] = "0"
		TwitterEachDayTotalSize[i] = "0"
	}

	c = session.DB("gmsTry").C("twitterStat")

	twitterResult := []TwitterStatPage{}

	err = c.Find(bson.M{}).All(&twitterResult)

	for _, element := range twitterResult {
		date := element.Date
		//fmt.Println("Date  ", date)
		s := strings.Split(date, "/")
		intDate, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("Error1!")
		}

		intMonth, err := strconv.Atoi(s[1])
		if err != nil {
			fmt.Println("Error2!")
		}

		if SelectedMonthIndex == intMonth {
			for i := 0; i < 31; i++ {
				if i == (intDate - 1) {
					TwitterEachDayTotalTweets[i] = element.TwitterNumber

					//fmt.Println("Vlue of i: ", i,  ", Twitter Number ", element.TwitterNumber, ", Twitter Size", element.TwitterSize)
					intTwitterTotalTweets, err := strconv.Atoi(element.TwitterNumber)
					if err != nil {
						fmt.Println("Error3!")
					}
					TwitterMonthTotalTweets += intTwitterTotalTweets

					intTwitterTotalSize, err := strconv.Atoi(element.TwitterSize)
					if err != nil {
						fmt.Println("Error4!")
					}
					TwitterEachDayTotalSize[i] = strconv.Itoa((intTwitterTotalSize / 1024) / 1024)
					TwitterMonthTotalSize += intTwitterTotalSize
				} /*else {
					TwitterEachDayTotalTweets[i] = "0"
					TwitterEachDayTotalSize [i] = "0"
				}*/
			}
		}

		intTwitterAllTotalTweets, err := strconv.Atoi(element.TwitterNumber)
		if err != nil {
			fmt.Println("Error5!")
		}
		TwitterAllTotalTweets += intTwitterAllTotalTweets

		intTwitterAllTotalSize, err := strconv.Atoi(element.TwitterSize)
		if err != nil {
			fmt.Println("Error6!")
		}
		TwitterAllTotalSize += intTwitterAllTotalSize
	}
	fmt.Println("Twitter Each Day Total Size:", TwitterEachDayTotalSize)
	TwitterAllTotalSize = (TwitterAllTotalSize / 1024) / 1024

	EvExpressAverage = strconv.FormatFloat(float64(EvExpressCountTotalSize)/float64(EvExpressCountTotalNews), 'f', 2, 64)

	Average = float64(TotalWholeSize) / float64(TotalWholeNews)

	FlickrEachDayTotalImages := make([]string, 31)
	FlickrEachDayTotalSize := make([]string, 31)
	FlickrMonthTotalImages := 0
	FlickrMonthTotalSize := 0
	FlickrAllTotalImages := 0
	FlickrAllTotalSize := 0

	for i := 0; i < 31; i++ {
		FlickrEachDayTotalImages[i] = "0"
		FlickrEachDayTotalSize[i] = "0"

	}

	c = session.DB("gmsTry").C("statTest")
	flickrResult := []FlickrStatPage{}
	err = c.Find(bson.M{}).All(&flickrResult)

	for _, element := range flickrResult {
		date := element.Date
		fmt.Println("Date  ", date)
		s := strings.Split(date, "/")
		intDate, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("Error1!")
		}
		//fmt.Println(element.FlickrNumber)

		intMonth, err := strconv.Atoi(s[1])
		if err != nil {
			fmt.Println("Error2!")
		}

		if SelectedMonthIndex == intMonth {
			//fmt.Println(SelectedMonthIndex)
			for i := 0; i < 31; i++ {
				if i == (intDate - 1) {
					FlickrEachDayTotalImages[i] = element.FlickrNumber

					//fmt.Println("Value of i: ", i,  ", FlickrNumber ", element.FlickrNumber, ", Flickr Size", element.FlickrSize)
					intFlickrTotalImages, err := strconv.Atoi(element.FlickrNumber)
					if err != nil {
						fmt.Println("Error3!")
					}
					FlickrMonthTotalImages += intFlickrTotalImages
					//fmt.Println(FlickrMonthTotalImages)
					FlickrEachDayTotalSize[i] = element.FlickrSize

					intFlickrTotalSize, err := strconv.Atoi(element.FlickrSize)
					if err != nil {
						fmt.Println("Error4!")
					}
					//FlickrEachDayTotalSize [i] = intFlickrTotalSize
					//strconv.Itoa((intFlickrTotalSize / 1024) / 1024)
					//fmt.Println(FlickrEachDayTotalSize [i] )
					FlickrMonthTotalSize += intFlickrTotalSize

				} /*else {
					TwitterEachDayTotalTweets[i] = "0"
					TwitterEachDayTotalSize [i] = "0"
				}*/
			}
		}

		intFlickrAllTotalImages, err := strconv.Atoi(element.FlickrNumber)
		if err != nil {
			fmt.Println("Error5!")
		}
		FlickrAllTotalImages += intFlickrAllTotalImages

		intFlickrAllTotalSize, err := strconv.Atoi(element.FlickrSize)
		if err != nil {
			fmt.Println("Error6!")
		}
		FlickrAllTotalSize += intFlickrAllTotalSize
	}

	/*TwitterTotalNews := make([]string, 31)
	TwitterTotalSize := make([]string, 31)
	TwitterAverage := ""
	TwitterPath := "/home/ripul/twitter/"

	for i := 0; i < 31; i++{
		twitterFileName
	}*/
	FlickrAllTotalSize = (FlickrAllTotalSize / 1024) / 1024

	//fmt.Println("Flickr Each Day Total Images:", FlickrEachDayTotalImages)
	//fmt.Println("Flickr Each Day Total size:", FlickrEachDayTotalSize)
	//fmt.Println("Flickr  Total Images:", FlickrAllTotalImages)
	//fmt.Println("Flickr  Total size:", FlickrAllTotalSize)
	//fmt.Println("Flickr Each Month Total Images:", FlickrEachDayTotalImages)
	//fmt.Println("Flickr Each Month Total size:", FlickrEachDayTotalSize)

	/*TwitterTotalNews := make([]string, 31)
	TwitterTotalSize := make([]string, 31)
	TwitterAverage := ""
	TwitterPath := "/home/ripul/twitter/"

	for i := 0; i < 31; i++{
		twitterFileName
	}*/

	Sources := []string{"BBC", "Courier", "Express", "Guardian", "Scotsman", "Independent", "Daily Record", "Evening Times", "Evening Express"}
	statPage := ScotlandStatPage{"MB", strconv.Itoa(AllTotalWholeNews), strconv.Itoa(AllTotalWholeSize), strconv.FormatFloat(AllTotalAverage, 'f', 2, 64), StartDate, SelectedMonth, strconv.Itoa(TotalWholeNews), strconv.Itoa(TotalWholeSize), strconv.FormatFloat(Average, 'f', 2, 64), labels, Dates, TotalNews, TotalSize,
		TwitterEachDayTotalTweets, TwitterEachDayTotalSize, strconv.Itoa(TwitterMonthTotalTweets), strconv.Itoa((TwitterMonthTotalSize / 1024) / 1024), strconv.Itoa(TwitterAllTotalTweets), strconv.Itoa(TwitterAllTotalSize),
		ScotTotalNews, ScotTotalSize, strconv.Itoa(ScotCountTotalNews), strconv.Itoa(ScotCountTotalSize), ScotAverage,
		ETTotalNews, ETTotalSize, strconv.Itoa(ETCountTotalNews), strconv.Itoa(ETCountTotalSize), ETAverage,
		BBCTotalNews, BBCTotalSize, strconv.Itoa(BBCCountTotalNews), strconv.Itoa(BBCCountTotalSize), BBCAverage,
		DRTotalNews, DRTotalSize, strconv.Itoa(DRCountTotalNews), strconv.Itoa(DRCountTotalSize), DRAverage,
		IndiTotalNews, IndiTotalSize, strconv.Itoa(IndiCountTotalNews), strconv.Itoa(IndiCountTotalSize), IndiAverage,
		GuardTotalNews, GuardTotalSize, strconv.Itoa(GuardCountTotalNews), strconv.Itoa(GuardCountTotalSize), GuardAverage,
		CourierTotalNews, CourierTotalSize, strconv.Itoa(CourierCountTotalNews), strconv.Itoa(CourierCountTotalSize), CourierAverage,
		ExpressTotalNews, ExpressTotalSize, strconv.Itoa(ExpressCountTotalNews), strconv.Itoa(ExpressCountTotalSize), ExpressAverage,
		EvExpressTotalNews, EvExpressTotalSize, strconv.Itoa(EvExpressCountTotalNews), strconv.Itoa(EvExpressCountTotalSize), EvExpressAverage,
		FlickrEachDayTotalImages, FlickrEachDayTotalSize, strconv.Itoa(FlickrMonthTotalImages), strconv.Itoa((FlickrMonthTotalSize / 1024) / 1024), strconv.Itoa(FlickrAllTotalImages), strconv.Itoa(FlickrAllTotalSize),
		Sources}

	defer session.Close()

	renderStatScotlandTemplate(w, "statScotland", &statPage)
}

func renderStatScotlandTemplate(w http.ResponseWriter, tmpl string, p *ScotlandStatPage) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func handleCWGMapImages(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := ""
	var doc bytes.Buffer

	location := r.FormValue("location")
	var pics []FlickrImage1
	st := r.FormValue("start")
	start, _ := strconv.Atoi(st)
	t, _ := template.ParseFiles("flickrHelper.html")

	fmt.Println(location)

	if location == "" {

		var pics []CwgImage
		pics = getCwgMapImages()

		b, err := json.Marshal(pics)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "%s", b)
	} else {
		if strings.HasPrefix(location, "getTags_") {
			tag := location[8:]
			pics = getFlickrMain(tag, "", start, "location", "")
			data := struct {
				Tag string
				//P   []FlickrImage
				P        []FlickrImage1
				PageIP   int
				PageIN   int
				Type     string
				Function string
			}{
				location,
				pics,
				start - 1,
				start + 1,
				"",
				"getMoreMapImages",
			}
			t.Execute(&doc, data)

		} else {

			loc := strings.Replace(location, "_", " ", -1)
			fmt.Println("in else " + loc)
			pics = getFlickrMain("", "", start, "location", loc)
			fmt.Println(pics)

			data := struct {
				Tag string
				//P   []FlickrImage
				P        []FlickrImage1
				PageIP   int
				PageIN   int
				Type     string
				Function string
			}{
				location,
				pics,
				start - 1,
				start + 1,
				"",
				"getMoreMapImages",
			}
			t.Execute(&doc, data)
		}

		s = doc.String()

		fmt.Fprintf(w, s)

	}
}

func handleMapImages(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	var trending []TrendingPlace
	var recomm []RecommendedPlace
	var heat []FlickrImage1
	var heatPhoto []Photo
	var trendingAll []TrendingAll

	if session.Values["user"] == nil {
		//no
		trendingAll = getTrendingAll()
		heat = getFlickrMap()

		flickrData := struct {
			//Heat []MapImage //replace with below
			Heat              []FlickrImage1
			TrendingMarker    []TrendingPlace
			TrendingMarkerAll []TrendingAll
			RecommendedMarker []RecommendedPlace
			User              string
		}{
			heat,
			trending,
			trendingAll,
			recomm,
			"no",
		}

		b, err := json.Marshal(flickrData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, "%s", b)
		return
	} else if session.Values["user"].(string) == "" {
		//no
		trendingAll = getTrendingAll()
		heat = getFlickrMap()

		flickrData := struct {
			//Heat []MapImage //replace with below
			Heat              []FlickrImage1
			TrendingMarker    []TrendingPlace
			TrendingMarkerAll []TrendingAll
			RecommendedMarker []RecommendedPlace
			User              string
		}{
			heat,
			trending,
			trendingAll,
			recomm,
			"no",
		}

		b, err := json.Marshal(flickrData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, "%s", b)
		return

	} else {
		//recomm = getRecommImages("550b107abc2e8b25cc000132")
		//trending = getTrendingImages("550b107abc2e8b25cc000132")

		recomm = getRecommImages(session.Values["user"].(string))
		trending = getTrendingImages(session.Values["user"].(string))

		dbConnection = NewMongoDBConn()
		sess := dbConnection.connect()
		c := sess.DB(db_name).C("photos")

		err := c.Find(bson.M{"owner": session.Values["user"].(string)}).All(&heatPhoto)
		if err != nil {
			fmt.Println("error while finding recommended places")
			fmt.Println(err)
		}
		defer sess.Close()

		flickrData := struct {
			//Heat []MapImage //replace with below
			Heat              []Photo
			TrendingMarker    []TrendingPlace
			TrendingMarkerAll []TrendingAll
			RecommendedMarker []RecommendedPlace
			User              string
		}{
			heatPhoto,
			trending,
			trendingAll,
			recomm,
			"yes",
		}

		b, err := json.Marshal(flickrData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, "%s", b)
		return
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	picture := r.FormValue("pic")
	//album := r.FormValue("album")
	//owner := r.FormValue("owner")
	cType := r.FormValue("cType")

	deleteFromOthers(picture, cType)

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	if cType == "image" {
		err := sess.DB(db_name).C("photos").Remove(bson.M{"photoId": picture})
		if err != nil {
			fmt.Println(err)
			defer sess.Close()
			fmt.Fprintf(w, "No")
			return
		}
	} else {
		err := sess.DB(db_name).C("videos").Remove(bson.M{"videoId": picture})
		if err != nil {
			fmt.Println(err)
			defer sess.Close()
			fmt.Fprintf(w, "No")
			return
		}
	}

	deleteFromOthers(picture, cType)
	resp := "Yes_" + picture

	defer sess.Close()

	fmt.Fprintf(w, resp)

}
func flickrHelper(flickr string, start int) string {
	s := ""
	var doc bytes.Buffer
	images := getFlickrImages(flickr, start)
	flickrData := struct {
		Tag    string
		PageIP int
		PageIN int
		P      []FlickrImage
	}{
		flickr,
		start - 1,
		start + 1,
		images,
	}
	t, _ := template.ParseFiles("pictureHelper.html")
	t.Execute(&doc, flickrData)
	s = doc.String()
	return s
}

func newsHelper(guardian string, start int) string {
	s := ""
	var doc bytes.Buffer
	news := getNews(guardian, start)
	newsData := struct {
		Tag    string
		PageIP int
		PageIN int
		N      []News
	}{
		guardian,
		start - 1,
		start + 1,
		news,
	}
	t, _ := template.ParseFiles("newsHelper.html")
	t.Execute(&doc, newsData)
	s = doc.String()
	return s
}

func getImages(request string, init string, temp string, start int, cType string) string {
	s := ""
	var doc bytes.Buffer
	var doc1 bytes.Buffer

	if request == "start" {

		photos := getFlickrMain("", "", start, cType, "")

		data := struct {
			//P []FlickrImage
			P []FlickrImage1
		}{
			photos,
		}

		t, _ := template.ParseFiles(temp)
		t.Execute(&doc, data)

	} else if strings.HasPrefix(request, "getTags") {
		response := make([]Response, 2)
		input := ""
		if strings.HasPrefix(request, "getTags_") {
			input = request[8:]
		} else {
			input = request[7:]
		}
		photos := getFlickrMain(input, "", start, cType, "")

		if len(photos) > 0 {

			data := struct {
				Tag string
				//P   []FlickrImage
				P        []FlickrImage1
				PageIP   int
				PageIN   int
				Type     string
				Function string
			}{
				input,
				photos,
				start - 1,
				start + 1,
				"and",
				"flickrMenu",
			}

			t, _ := template.ParseFiles(temp)
			t.Execute(&doc, data)
			s = doc.String()

			tags := photos[0].Keywords
			tagString := ""
			if len(tags) < 15 {
				for img := range photos {
					if len(tags) > 15 {
						break
					}
					if img > 0 {
						for tag := range photos[img].Keywords {
							flag := false
							for existing := range tags {
								if tags[existing] == photos[img].Keywords[tag] || strings.ToLower(tags[existing]) == photos[img].Keywords[tag] {
									flag = true
								}
							}
							if flag == false {
								tags = append(tags, photos[img].Keywords[tag])
								tagString += "," + photos[img].Keywords[tag]
							}
						}
					} else {
						for tag := range tags {
							if tag != len(tags)-1 {
								tagString += tags[tag] + ","
							} else {
								tagString += tags[tag]
							}
						}
					}
				}
			} else {
				for tag := range tags {
					if tag != len(tags)-1 {
						tagString += tags[tag] + ","
					} else {
						tagString += tags[tag]
					}
				}
			}

			response[0].Name = "tags"
			response[0].Content = tagString
			response[1].Name = "pics"
			response[1].Content = s

			b, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err)
			}

			//ret := (string)b

			return string(b)
		} else {
			return "No content found with requested tag"
		}

	} else {
		response := make([]Response, 2)

		if cType == "and" {
			photos := getFlickrMain(request, init, start, cType, "")
			if len(photos) > 0 {

				data := struct {
					Tag string
					//P   []FlickrImage
					P        []FlickrImage1
					PageIP   int
					PageIN   int
					Type     string
					Function string
				}{
					request,
					photos,
					start - 1,
					start + 1,
					"and",
					"flickrMenu",
				}

				t, _ := template.ParseFiles(temp)
				t.Execute(&doc, data)
				s = doc.String()

				response[0].Name = "and"
				response[0].Content = s
				response[1].Name = "or"
				response[1].Content = ""

				b, err := json.Marshal(response)
				if err != nil {
					fmt.Println(err)
				}

				//ret := (string)b

				return string(b)

			} else {
				return "No content found with requested tag"
			}
		} else if cType == "or" {

			photos := getFlickrMain(request, init, start, cType, "")
			if len(photos) > 0 {

				data := struct {
					Tag string
					//P   []FlickrImage
					P        []FlickrImage1
					PageIP   int
					PageIN   int
					Type     string
					Function string
				}{
					request,
					photos,
					start - 1,
					start + 1,
					"or",
					"flickrMenu",
				}

				t, _ := template.ParseFiles(temp)
				t.Execute(&doc, data)
				s = doc.String()

				response[0].Name = "and"
				response[0].Content = ""
				response[1].Name = "or"
				response[1].Content = s

				b, err := json.Marshal(response)
				if err != nil {
					fmt.Println(err)
				}

				//ret := (string)b

				return string(b)

			} else {
				return "No content found with requested tag"
			}

		} else {
			photos := getFlickrMain(request, init, start, "and", "")
			sAnd := ""
			sOr := ""
			if len(photos) > 0 {

				data := struct {
					Tag string
					//P   []FlickrImage
					P        []FlickrImage1
					PageIP   int
					PageIN   int
					Type     string
					Function string
				}{
					request,
					photos,
					start - 1,
					start + 1,
					"and",
					"flickrMenu",
				}

				t, _ := template.ParseFiles(temp)
				t.Execute(&doc, data)
				sAnd = doc.String()

			} else {
				sAnd = "No content found with requested tags"
			}

			photos = getFlickrMain(request, init, start, "or", "")
			if len(photos) > 0 {

				data := struct {
					Tag string
					//P   []FlickrImage
					P        []FlickrImage1
					PageIP   int
					PageIN   int
					Type     string
					Function string
				}{
					request,
					photos,
					start - 1,
					start + 1,
					"or",
					"flickrMenu",
				}

				t, _ := template.ParseFiles(temp)
				t.Execute(&doc1, data)
				sOr = doc1.String()
			} else {
				sOr = "No content found with requested tags"
			}

			response[0].Name = "and"
			response[0].Content = sAnd
			response[1].Name = "or"
			response[1].Content = sOr

			b, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err)
			}

			//ret := (string)b

			return string(b)

		}

	}

	s = doc.String()

	return s
}

func handleFlickrGeneral(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	request := r.FormValue("req")
	init := r.FormValue("init")
	st := r.FormValue("start")
	cType := r.FormValue("cType")
	s := ""
	//var doc bytes.Buffer
	var start int

	//response := make([]Response, 2)
	if st == "" {
		start = 0
	} else {
		start, _ = strconv.Atoi(st)
	}

	if request == "start" {
		s = getImages(request, "", "flickrImages.html", start, "")
		fmt.Fprintf(w, s)
	} else if strings.HasPrefix(request, "getTags") {
		s = getImages(request, "", "flickrHelper.html", start, "")

		fmt.Fprintf(w, s)

	} else {
		s = getImages(request, init, "flickrHelper.html", start, cType)

		fmt.Fprintf(w, s)

	}
}

func handleFlickrNews(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	request := r.FormValue("req")
	st := r.FormValue("start")
	cType := r.FormValue("cType")
	s := ""
	var doc bytes.Buffer
	var start int

	response := make([]Response, 2)
	if st == "" {
		start = 0
	} else {
		start, _ = strconv.Atoi(st)
	}

	if request == "start" {

		t, _ := template.ParseFiles("flickrNews.html")
		t.Execute(&doc, nil)
		s = doc.String()
		tags := "Scotland 10761,Glasgow 2740,Swimming 132,Cycling 155,Weightlifting 101,Wales 74,Gymnastics 81,Netball 53,London 43,Boxing 55,Tennis 30,Triathlon 36,Wrestling 6,Diving 30,Squash 12,Photography 8,India 15,Badminton 6,maximum 10761"
		response[0].Name = "html"
		response[0].Content = s
		response[1].Name = "tags"
		response[1].Content = tags

		b, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Printf("%s", b)
		fmt.Fprintf(w, "%s", b)

	} else {
		flickr := ""
		guardian := ""

		if strings.HasPrefix(request, "tag_") {
			flickr = request[4:]
			a := []rune(flickr)
			a[0] = unicode.ToUpper(a[0])
			guardian = string(a)
		} else {
			guardian = request
			flickr = strings.ToLower(request)
		}
		if cType == "image" {

			response[0].Name = "flickr"
			response[0].Content = flickrHelper(flickr, start)
			response[1].Name = "news"
			response[1].Content = ""

		} else if cType == "news" {

			response[1].Name = "news"
			response[1].Content = newsHelper(guardian, start)
			response[0].Name = "flickr"
			response[0].Content = ""
		} else {

			response[1].Name = "news"
			response[1].Content = newsHelper(guardian, start)
			response[0].Name = "flickr"
			response[0].Content = flickrHelper(flickr, start)
		}

		b, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Printf("%s", b)
		fmt.Fprintf(w, "%s", b)

	}
}

func handleVideos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	start, _ := strconv.Atoi(r.FormValue("req"))
	limit := 8
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(currentUser)

	response := make([]Response, 1)
	s := ""
	var doc bytes.Buffer

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var videos []Video
	err := sess.DB(db_name).C("videos").Find(bson.M{"owner": u.Id}).Skip(start * limit).Limit(limit).All(&videos)

	if len(videos) > 0 || start == 0 {
		data := struct {
			PageP int
			PageN int
			Video []Video
		}{
			start - 1,
			start + 1,
			videos,
		}

		//fmt.Println(photos)

		t, _ := template.ParseFiles("videosTemplate.html")
		if t == nil {
			fmt.Println("no template - videosTemplate")
		}
		t.Execute(&doc, data)
		s = doc.String()
	} else {
		s = ""
	}

	response[0].Name = "ownVideos"
	response[0].Content = s

	//fmt.Println(s)

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	defer sess.Close()

	//fmt.Printf("%s", b)
	fmt.Fprintf(w, "%s", b)
	return

	//authenticated, _ := template.ParseFiles("videos.html")
	//authenticated.Execute(w, u)

}

func handleCms(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	u := &User{}
	if session.Values["user"] == nil {
		u = &User{}
	} else if session.Values["user"].(string) == "" {
		u = &User{}
	} else {
		u = findUser(session.Values["user"].(string))
	}

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var p DisplayPhotos
	c := sess.DB(db_name).C("displayPhotos")
	err := c.Find(bson.M{"name": "views"}).One(&p)
	if err != nil {
		fmt.Println("could not get most viewed photos")
	}

	var recent DisplayPhotos
	c = sess.DB(db_name).C("displayPhotos")
	err = c.Find(bson.M{"name": "recent"}).One(&recent)
	if err != nil {
		fmt.Println("could not get most viewed photos")
	}

	flickrImages := getFlickrImages("", 0)

	news := getNews("", 0)

	data := struct {
		P      DisplayPhotos
		R      DisplayPhotos
		Flickr []FlickrImage
		N      []News
		U      User
	}{
		p,
		recent,
		flickrImages,
		news,
		*u,
	}

	defer sess.Close()
	authenticated, _ := template.ParseFiles("cmsHome.html")
	authenticated.Execute(w, data)

}

func handleUpvote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	picId := r.FormValue("picId")
	//albumId := r.FormValue("albumId")
	//owner := r.FormValue("picOwner")
	cType := r.FormValue("cType")

	//user := findUser(dbConnection, owner)

	//var al int
	photo := Photo{}
	video := Video{}

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	if cType == "image" {
		err := sess.DB(db_name).C("photos").Find(bson.M{"photoId": picId}).One(&photo)
		photo.Views = photo.Views + 1
		err = sess.DB(db_name).C("photos").Update(bson.M{"photoId": picId}, bson.M{"$set": bson.M{"views": photo.Views}})
		if err != nil {
			fmt.Println("could not update photos in tag db")
			fmt.Println(err)
			defer sess.Close()
			fmt.Fprintf(w, "No")
		}
	} else {
		err := sess.DB(db_name).C("videos").Find(bson.M{"videoId": picId}).One(&video)
		video.Views = video.Views + 1
		err = sess.DB(db_name).C("videos").Update(bson.M{"videoId": picId}, bson.M{"$set": bson.M{"views": video.Views}})
		if err != nil {
			fmt.Println("could not update views in videos db")
			fmt.Println(err)
			defer sess.Close()
			fmt.Fprintf(w, "No")
		}
	}

	updateTagDB(photo, video)
	updateMostViewed(photo, video)
	updateMostRecent(photo, video)

	defer sess.Close()

	if cType == "image" {
		fmt.Fprintf(w, "Yes_"+strconv.Itoa(photo.Views))
	} else {
		fmt.Fprintf(w, "Yes_"+strconv.Itoa(video.Views))
	}

}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	r.ParseForm()
	email := r.FormValue("email")
	pass := r.FormValue("pass")
	c := find(email)
	if c == nil {
		defer sess.Close()
		fmt.Fprintf(w, "No")
	} else {
		if c.Password == pass {
			session, _ := store.Get(r, "cookie")
			session.Values["user"] = c.Id
			session.Save(r, w)

			defer sess.Close()
			fmt.Fprintf(w, "Yes_"+c.FirstName)
		} else {
			defer sess.Close()
			fmt.Fprintf(w, "No")
		}
	}
}

func handleAuthenticated(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	u := findUser(currentUser)
	var photos []Photo

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	err := sess.DB(db_name).C("photos").Find(bson.M{"owner": u.Id}).Skip(0).Limit(8).All(&photos)
	if err != nil {
		fmt.Println("could not get images from DB")
	}

	photoData := struct {
		FirstName string
		PageN     int
		PageP     int
		Photo     []Photo
	}{
		u.FirstName,
		1,
		1,
		photos,
	}
	defer sess.Close()
	authenticated, _ := template.ParseFiles("pictures2.html")
	authenticated.Execute(w, photoData)
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

	authenticated, _ := template.ParseFiles("index.html")
	authenticated.Execute(w, session.Values["user"].(string))

}

func handleMainUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	user := findUser(currentUser)

	u := r.URL.RawQuery
	user2 := findUser(u)

	var photos []Photo
	var videos []Video

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	sess.DB(db_name).C("photos").Find(bson.M{"owner": u}).Skip(0).Limit(8).All(&photos)
	sess.DB(db_name).C("videos").Find(bson.M{"owner": u}).Skip(0).Limit(8).All(&videos)

	photoData := struct {
		FirstName string
		UserName  string
		PageIN    int
		PageIP    int
		PageVN    int
		PageVP    int
		User      string
		Photo     []Photo
		Video     []Video
	}{
		user.FirstName,
		user2.FirstName,
		1,
		-1,
		1,
		-1,
		u,
		photos,
		videos,
	}

	defer sess.Close()

	authenticated, _ := template.ParseFiles("otherUsers.html")
	authenticated.Execute(w, photoData)

}

func handleMainFlickr(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	user := findUser(currentUser)

	u := r.URL.RawQuery

	flickr := u
	a := []rune(flickr)
	a[0] = unicode.ToUpper(a[0])
	guardian := string(a)

	images := getFlickrImages(flickr, 0)
	news := getNews(guardian, 0)
	data := struct {
		Tag       string
		FirstName string
		Tags      string
		P         []FlickrImage
		N         []News
		PageIP    int
		PageIN    int
	}{
		u,
		user.FirstName,
		"Scotland 10761,Glasgow 2740,Swimming 132,Cycling 155,Weightlifting 101,Wales 74,Gymnastics 81,Netball 53,London 43,Boxing 55,Tennis 30,Triathlon 36,Wrestling 6,Diving 30,Squash 12,Photography 8,India 15,Badminton 6,maximum 10761",
		images,
		news,
		-1,
		1,
	}

	authenticated, _ := template.ParseFiles("flickr2.html")
	authenticated.Execute(w, data)
}

func handleUserProfile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := r.FormValue("user")
	start := r.FormValue("start")
	cType := r.FormValue("cType")
	nModP := r.FormValue("nModP")
	//nModR := r.FormValue("nModR")
	st, _ := strconv.Atoi(start)
	nMod, _ := strconv.Atoi(nModP)
	nMod += 1
	limit := 8
	flag := true

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var photos []Photo
	var videos []Video
	var sti int
	var stv int
	var doc bytes.Buffer
	s := ""

	if t == "" {
		sess.DB(db_name).C("photos").Find(bson.M{"owner": t}).Skip(0).Limit(3).All(&photos)
		sess.DB(db_name).C("videos").Find(bson.M{"owner": t}).Skip(0).Limit(3).All(&videos)
		sti = 0
		stv = 0

	} else {

		if cType == "" {
			err := sess.DB(db_name).C("photos").Find(bson.M{"owner": t}).Skip(st * limit).Limit(limit).All(&photos)
			if err != nil {
				fmt.Println(err)
			}

			err = sess.DB(db_name).C("videos").Find(bson.M{"owner": t}).Skip(st * limit).Limit(limit).All(&videos)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(photos)
			sti = 0
			stv = 0
		} else if cType == "image" {
			err := sess.DB(db_name).C("photos").Find(bson.M{"owner": t}).Skip(st * limit).Limit(limit).All(&photos)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(photos)

			if len(photos) == 0 {
				flag = false
			}

			err = sess.DB(db_name).C("videos").Find(bson.M{"owner": t}).Skip(nMod * limit).Limit(limit).All(&videos)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(photos)
			sti = st
			stv = nMod
		} else {
			err := sess.DB(db_name).C("photos").Find(bson.M{"owner": t}).Skip(nMod * limit).Limit(limit).All(&photos)
			if err != nil {
				fmt.Println(err)
			}

			err = sess.DB(db_name).C("videos").Find(bson.M{"owner": t}).Skip(st * limit).Limit(limit).All(&videos)
			if err != nil {
				fmt.Println(err)
			}

			if len(videos) == 0 {
				flag = false
			}

			sti = nMod
			stv = st
		}
	}

	if flag == true {

		photoData := struct {
			Owner  string
			PageIN int
			PageIP int
			PageVN int
			PageVP int
			User   string
			Photo  []Photo
			Video  []Video
		}{
			findUser(t).FirstName,
			sti + 1,
			sti - 1,
			stv + 1,
			stv - 1,
			t,
			photos,
			videos,
		}

		temp, _ := template.ParseFiles("photoVideoTemplate.html")
		if temp == nil {
			fmt.Println("no template photo video template")
		}

		temp.Execute(&doc, photoData)
		s = doc.String()
	} else {
		s = ""
	}

	defer sess.Close()

	fmt.Fprintf(w, s)

}

func handleCreateAlbum(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	//description := r.FormValue("description")

	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	c := findUser(currentUser)

	albumId := createAlbum(name, c.Id, c.FirstName+" "+c.LastName)

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
	u := findUser(session.Values["user"].(string))

	if u == nil {
		u = &User{}
	}
	http.Redirect(w, r, "/cmsHome", http.StatusFound)
}

func checkLoggedIn(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")

	if session.Values["user"] == nil {
		fmt.Fprintf(w, "No")
	} else if session.Values["user"].(string) == "" {
		fmt.Fprintf(w, "No")
	} else {
		message := "Yes," + findUser(session.Values["user"].(string)).FirstName
		fmt.Fprintf(w, message)
	}
}

func createTagCloud(w http.ResponseWriter, r *http.Request) {
	result := getAllTags()
	var tags string
	var max = 0

	fmt.Println(len(result), "   length of result")

	if len(result) > 30 {
		for i := len(result); i > 0; i-- {
			for j := 0; j < i-1; j++ {
				if len(result[j].Photos)+len(result[j].Videos) > len(result[j+1].Photos)+len(result[j+1].Videos) {
					help := result[j]
					result[j] = result[j+1]
					result[j+1] = help
				}
			}
		}

		result = result[len(result)-30:]
	}

	fmt.Println(len(result), "   length of result")

	dest := make([]Tag, len(result))
	perm := rand.Perm(len(result))
	for i, v := range perm {
		dest[v] = result[i]
	}

	for tag := range dest {
		if len(dest[tag].Photos)+len(dest[tag].Videos) > max {
			max = len(dest[tag].Photos) + len(dest[tag].Videos)
		}

		tags += dest[tag].Name + " " + strconv.Itoa(len(dest[tag].Photos)+len(dest[tag].Videos)) + ","
	}
	tags += "maximum " + strconv.Itoa(max)

	fmt.Println(tags)
	fmt.Fprintf(w, tags)

}

func handleTag(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := r.FormValue("tag")
	start := r.FormValue("start")
	cType := r.FormValue("cType")
	nModP := r.FormValue("nModP")
	//nModR := r.FormValue("nModR")
	st, _ := strconv.Atoi(start)
	nMod, _ := strconv.Atoi(nModP)
	nMod += 1
	limit := 8

	var photos []Photo
	var videos []Video
	var sti int
	var stv int
	var doc bytes.Buffer
	s := ""
	flag := true

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	tagStruct := findByTag(t)

	if cType == "" {
		if len(tagStruct.Photos) < limit {
			photos = tagStruct.Photos[:len(tagStruct.Photos)]
		} else {
			photos = tagStruct.Photos[:limit]
		}

		if len(tagStruct.Videos) < limit {
			videos = tagStruct.Videos[:len(tagStruct.Videos)]
		} else {
			videos = tagStruct.Videos[:limit]
		}

		sti = 0
		stv = 0
	} else if cType == "image" {
		if st*limit > len(tagStruct.Photos) {
			st -= 1
			photos = tagStruct.Photos[st*limit : len(tagStruct.Photos)]
		} else {
			if st*limit+limit > len(tagStruct.Photos) {
				photos = tagStruct.Photos[st*limit : len(tagStruct.Photos)]
			} else {
				photos = tagStruct.Photos[st*limit : st*limit+limit]
			}
		}
		if len(tagStruct.Videos) < nMod*limit {
			nMod -= 1
			videos = tagStruct.Videos[nMod*limit : len(tagStruct.Videos)]
		} else {
			if nMod*limit+limit > len(tagStruct.Videos) {

				videos = tagStruct.Videos[nMod*limit : len(tagStruct.Videos)]
			} else {
				videos = tagStruct.Videos[nMod*limit : nMod*limit+limit]
			}

		}

		sti = st
		stv = nMod
	} else {

		if st*limit > len(tagStruct.Videos) {
			st -= 1
			videos = tagStruct.Videos[st*limit : len(tagStruct.Videos)]
		} else {
			if st*limit+limit > len(tagStruct.Videos) {
				videos = tagStruct.Videos[st*limit : len(tagStruct.Videos)]
			} else {
				videos = tagStruct.Videos[st*limit : st*limit+limit]
			}
		}
		if len(tagStruct.Photos) < nMod*limit {
			nMod -= 1
			photos = tagStruct.Photos[nMod*limit : len(tagStruct.Photos)]
		} else {
			if nMod*limit+limit > len(tagStruct.Photos) {

				photos = tagStruct.Photos[nMod*limit : len(tagStruct.Photos)]
			} else {
				photos = tagStruct.Photos[nMod*limit : nMod*limit+limit]
			}

		}

		sti = nMod
		stv = st

	}

	if flag == true {

		photoData := struct {
			PageIN int
			PageIP int
			PageVN int
			PageVP int
			Tag    string
			Photo  []Photo
			Video  []Video
		}{
			sti + 1,
			sti - 1,
			stv + 1,
			stv - 1,
			t,
			photos,
			videos,
		}

		temp, _ := template.ParseFiles("tagContentTemplate.html")
		if temp == nil {
			fmt.Println("no template tag content template")
		}

		temp.Execute(&doc, photoData)
		s = doc.String()
	} else {
		s = ""
	}

	defer sess.Close()

	fmt.Fprintf(w, s)
}

func handleMainTag(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	user := findUser(currentUser)

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	u := r.URL.RawQuery
	tagStruct := findByTag(u)

	var photos []Photo
	var videos []Video

	limit := 8

	if len(tagStruct.Photos) < limit {
		photos = tagStruct.Photos[:len(tagStruct.Photos)]
	} else {
		photos = tagStruct.Photos[:limit]
	}

	if len(tagStruct.Videos) < limit {
		videos = tagStruct.Videos[:len(tagStruct.Videos)]
	} else {
		videos = tagStruct.Videos[:limit]
	}

	photoData := struct {
		FirstName string
		PageIN    int
		PageIP    int
		PageVN    int
		PageVP    int
		Tag       string
		Photo     []Photo
		Video     []Video
	}{
		user.FirstName,
		1,
		-1,
		1,
		-1,
		u,
		photos,
		videos,
	}
	defer sess.Close()

	authenticated, _ := template.ParseFiles("taggedPictures2.html")
	authenticated.Execute(w, photoData)
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

	createDefaultAlbum(id.Hex(), fname+" "+lname)

	newUser := User{id, fname, lname, email, pass, "", "", "", id.Hex()}
	add(newUser)

	c := find(email)

	if c == nil {
		fmt.Fprintf(w, "No")
	} else {

		session, _ := store.Get(r, "cookie")
		session.Values["user"] = c.Id
		session.Save(r, w)
		fmt.Fprintf(w, "Yes")
	}

}

func handlePassReset(w http.ResponseWriter, r *http.Request) {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()
	r.ParseForm()

	email := r.FormValue("email")
	pass := r.FormValue("pass")

	fmt.Println(pass)

	err := sess.DB(db_name).C("user").Update(bson.M{"email": email}, bson.M{"$set": bson.M{"password": pass}})
	if err != nil {
		fmt.Println(err)
	}
	c := find(email)

	fmt.Println(c)
	if c == nil {
		defer sess.Close()
		fmt.Fprintf(w, "No")
	} else {
		defer sess.Close()
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
	lng := r.FormValue("lng")
	lat := r.FormValue("lat")
	locationN := r.FormValue("locality")
	/*if loc == "" {
		lng = ""
		lat = ""
		locationN = ""
	} */

	streetN := r.FormValue("formatted_address")
	streetN = strings.Split(streetN, ",")[0]
	tags := r.FormValue("tagList")

	var location = *new(Location)
	if lat != "" && lng != "" {
		location = Location{streetN + ", " + locationN, lat, lng}
	}

	t := make([]string, 0)
	if tags != "" {
		t = parseTags(tags, image)
	}

	id := bson.NewObjectId()
	p := Photo{}
	v := Video{}

	session, _ := store.Get(r, "cookie")
	user := session.Values["user"].(string)
	currentUser := findUser(user)

	//c := uploadToAlbum(cType, image, caption, album, lng, lat, streetN+", "+locationN, tags, currentUser)

	if cType == "image" {
		p = Photo{id, id.Hex(), currentUser.Id, currentUser.FirstName + " " + currentUser.LastName, album, image, caption, location, time.Now().Local().Format("02/01/2006"), 0, t, make([]PhotoComment, 1)}
		addTags(t, p, Video{})

		dbConnection = NewMongoDBConn()
		sess := dbConnection.connect()
		c := sess.DB(db_name).C("photos")
		err := c.Insert(p)
		if err != nil {
			panic(err)
		}
		defer sess.Close()
	} else {
		v = Video{id, id.Hex(), currentUser.Id, currentUser.FirstName + " " + currentUser.LastName, album, image, caption, location, time.Now().Local().Format("02/01/2006"), 0, t, make([]PhotoComment, 1)}
		addTags(t, Photo{}, v)
		dbConnection = NewMongoDBConn()
		sess := dbConnection.connect()
		c := sess.DB(db_name).C("videos")
		err := c.Insert(v)
		if err != nil {
			panic(err)
		}

		defer sess.Close()

	}

	insertInMostRecent(p, v)

}

func parseTags(tags string, filename string) []string {
	tags = strings.ToLower(tags)
	s := strings.Split(tags, ",")

	return s
}

func getPictures(collName string, field string, userId string, templateName string, start int) string {
	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	s := ""
	var doc bytes.Buffer
	var photos []Photo
	limit := 8
	err := sess.DB(db_name).C(collName).Find(bson.M{field: userId}).Skip(start * limit).Limit(limit).All(&photos)

	if err != nil {
		fmt.Println(err)
	}
	if len(photos) > 0 || start == 0 {
		photoData := struct {
			PageN int
			PageP int
			Photo []Photo
		}{
			start + 1,
			start - 1,
			photos,
		}

		t, _ := template.ParseFiles(templateName)
		if t == nil {
			fmt.Println("no template", templateName)
		}

		t.Execute(&doc, photoData)
		s = doc.String()
	} else {
		s = ""
	}

	defer sess.Close()
	return s

}

func handlePictures(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	start := r.FormValue("req")
	s, _ := strconv.Atoi(start)

	session, _ := store.Get(r, "cookie")
	currentUser := session.Values["user"].(string)
	response := make([]Response, 1)

	response[0].Name = "ownPictures"
	response[0].Content = getPictures("photos", "owner", currentUser, "pictureTemplate.html", s)

	//fmt.Println(s)

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("%s", b)
	fmt.Fprintf(w, "%s", b)
	return

}

func handleAlbums(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.FormValue("albumId")
	start := r.FormValue("start")
	cType := r.FormValue("cType")
	nModP := r.FormValue("nModP")
	//nModR := r.FormValue("nModR")
	st, _ := strconv.Atoi(start)
	nMod, _ := strconv.Atoi(nModP)
	nMod += 1
	limit := 8

	session, _ := store.Get(r, "cookie")
	user := session.Values["user"].(string)
	currentUser := findUser(user)
	response := make([]Response, 1)
	s := ""
	var doc bytes.Buffer
	flag := true

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	if query == "" {

		var al []Album
		err := sess.DB(db_name).C("albums").Find(bson.M{"owner": currentUser.Id}).All(&al)
		if err != nil {
			fmt.Println(err)
		}
		albums := make([]AlbumStruct, len(al))

		for i := range al {
			albums[i].Name = al[i].Name
			albums[i].AlbumId = al[i].AlbumId
			var photo Photo
			err = sess.DB(db_name).C("photos").Find(bson.M{"albumId": albums[i].AlbumId}).One(&photo)
			if err != nil {
				fmt.Println(err)
			}
			albums[i].Photo = photo.URL
		}

		data := struct {
			Page   string
			Albums []AlbumStruct
		}{
			"0",
			albums,
		}
		t, _ := template.ParseFiles("albumTemplate.html")
		if t == nil {
			fmt.Println("no template album template")
		}
		t.Execute(&doc, data)
		s = doc.String()
		response[0].Name = "ownAlbums"
	} else {

		var photos []Photo
		var videos []Video
		var sti int
		var stv int

		if cType == "" {
			err := sess.DB(db_name).C("photos").Find(bson.M{"albumId": query}).Skip(st * limit).Limit(limit).All(&photos)
			if err != nil {
				fmt.Println(err)
			}

			err = sess.DB(db_name).C("videos").Find(bson.M{"albumId": query}).Skip(st * limit).Limit(limit).All(&videos)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(photos)
			sti = 0
			stv = 0
		} else if cType == "image" {
			err := sess.DB(db_name).C("photos").Find(bson.M{"albumId": query}).Skip(st * limit).Limit(limit).All(&photos)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(photos)

			if len(photos) == 0 {
				flag = false
			}

			err = sess.DB(db_name).C("videos").Find(bson.M{"albumId": query}).Skip(nMod * limit).Limit(limit).All(&videos)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(photos)
			sti = st
			stv = nMod
		} else {
			err := sess.DB(db_name).C("photos").Find(bson.M{"albumId": query}).Skip(nMod * limit).Limit(limit).All(&photos)
			if err != nil {
				fmt.Println(err)
			}

			err = sess.DB(db_name).C("videos").Find(bson.M{"albumId": query}).Skip(st * limit).Limit(limit).All(&videos)
			if err != nil {
				fmt.Println(err)
			}

			if len(videos) == 0 {
				flag = false
			}

			sti = nMod
			stv = st
		}

		if flag == true {

			photoData := struct {
				PageIN  int
				PageIP  int
				PageVN  int
				PageVP  int
				AlbumId string
				Photo   []Photo
				Video   []Video
			}{
				sti + 1,
				sti - 1,
				stv + 1,
				stv - 1,
				query,
				photos,
				videos,
			}

			temp, _ := template.ParseFiles("albumDetailTemplate.html")
			if temp == nil {
				fmt.Println("no template album detail template")
			}

			temp.Execute(&doc, photoData)
			s = doc.String()
		} else {
			s = ""
		}

		response[0].Name = "albumDetail"
	}

	response[0].Content = s

	//fmt.Println(s)

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	defer sess.Close()

	//fmt.Printf("%s", b)
	fmt.Fprintf(w, "%s", b)

	return

}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie")
	u := session.Values["user"].(string)
	currentUser := findUser(u)
	response := make([]Response, 1)
	s := ""
	var doc bytes.Buffer

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	var albums []Album
	err := sess.DB(db_name).C("albums").Find(bson.M{"owner": currentUser.Id}).All(&albums)
	if err != nil {
		fmt.Println(err)
	}

	data := struct {
		Albums []Album
	}{
		albums,
	}
	t, _ := template.ParseFiles("upload2.html")
	t.Execute(&doc, data)
	s = doc.String()

	fmt.Println(s)
	response[0].Name = "upload"
	response[0].Content = s

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	defer sess.Close()

	fmt.Fprintf(w, "%s", b)
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	comment := r.FormValue("comment")
	picture := r.FormValue("pic")
	cType := r.FormValue("cType")

	session, _ := store.Get(r, "cookie")
	user2 := session.Values["user"].(string)

	currentUser := findUser(user2)
	com := PhotoComment{currentUser.FirstName + " " + currentUser.LastName, currentUser.Id, comment, time.Now().Local().Format("02/01/2006")}

	photo := Photo{}
	video := Video{}

	dbConnection = NewMongoDBConn()
	sess := dbConnection.connect()

	if cType == "image" {
		err := sess.DB(db_name).C("photos").Find(bson.M{"photoId": picture}).One(&photo)
		photo.Comments = append(photo.Comments, com)
		err = sess.DB(db_name).C("photos").Update(bson.M{"photoId": picture}, bson.M{"$set": bson.M{"comments": photo.Comments}})
		if err != nil {
			fmt.Println("could not update photos in tag db")
			fmt.Println(err)
			defer sess.Close()
			fmt.Fprintf(w, "No")
			return
		}
	} else {
		err := sess.DB(db_name).C("videos").Find(bson.M{"videoId": picture}).One(&video)
		video.Comments = append(video.Comments, com)
		err = sess.DB(db_name).C("videos").Update(bson.M{"videoId": picture}, bson.M{"$set": bson.M{"comments": video.Comments}})
		if err != nil {
			fmt.Println("could not update views in videos db")
			fmt.Println(err)
			defer sess.Close()
			fmt.Fprintf(w, "No")
			return
		}
	}

	updateTagDB(photo, video)
	updateMostRecent(photo, video)
	updateMostViewed(photo, video)
	defer sess.Close()
	response := com.Body + "_" + com.User + "_" + com.Timestamp
	fmt.Fprintf(w, "Yes_"+response)
}
