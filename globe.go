package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AllApiData struct {
	Images  []string
	Weather *WeatherData
}

type WeatherData struct {
	Temp string
	City string
	Icon string
}

var dispdata AllApiData
var celsiusNum string
var imagesArray []string
var rainOrShine string
var RANDi int

var cityLibrary = []string{
	"Tokyo",
	"Paris",
	"Singapore",
	"London",
	"Shanghai",
	"Beijing",
	"Seoul",
	"Mumbai",
	"Washington",
	"Bangkok",
	"Hanoi",
	"Toronto",
	"Atlanta",
	"Rome",
	"Milan",
	"Edinburgh",
	"Vienna",
	"Prague",
	"Stockholm",
	"Vancouver",
	"Barcelona",
	"Sydney",
	"Istanbul",
	"Hokkaido",
	"Santiago",
	"Valencia",
	"Peru",
	"Moscow",
	"Florence",
	"Berlin",
	"Auckland",
	"Kyoto",
}

//API funcs

//doc for Flickr API: https://www.flickr.com/services/api/flickr.photos.search.html
func ImageDisplay() {
	reqUrl := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%s&tags=%s&extras=url_m&format=json&nojsoncallback=1&sort=relevance",
		os.Getenv("FLICKR_APIKEY"),
		cityLibrary[RANDi])

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ReadAll: ", err)
		return
	}

	type FlickrResponse struct {
		Photos struct {
			Photo []struct {
				Id, Secret, Server string
				Farm               int
			}
		}
	}

	var f FlickrResponse
	if err := json.Unmarshal(body, &f); err != nil {
		log.Fatalf("unmarshal err: %v", err)
	}
	imagesArray = []string{} //resets previous array on click
	v := rand.Perm(100)[:27] //get different numbers
	for i := 0; i < 27; i++ {
		b := v[i]
		respUrl := "https://farm" + strconv.Itoa(f.Photos.Photo[b].Farm) + ".staticflickr.com/" + f.Photos.Photo[b].Server + "/" + f.Photos.Photo[b].Id + "_" + f.Photos.Photo[b].Secret + "_q.jpg"

		imagesArray = append(imagesArray, respUrl)
	}
}

//doc for weather API: http://openweathermap.org/weather-data#current
func WeatherDisplay() {
	reqUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&APPID=%s",
		cityLibrary[RANDi],
		os.Getenv("WEATHER_APIKEY"),
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ReadAll: ", err)
		return
	}

	type WeatherResponse struct {
		Main struct {
			Temp float64
		}
		Weather []struct {
			Icon string
		}
	}
	var f WeatherResponse
	if err := json.Unmarshal(body, &f); err != nil {
		log.Fatalf("unmarshal weatherdisplay %v", err)
	}

	celsiusNum = fmt.Sprintf("%.1f", f.Main.Temp-273.15) //formula to get celsius
	rainOrShine = fmt.Sprintf("http://openweathermap.org/img/w/%s.png", f.Weather[0].Icon)
}

//handler and template
func homeHandler(w http.ResponseWriter, r *http.Request) {
	RANDi = rand.Intn(len(cityLibrary))
	ImageDisplay()
	WeatherDisplay()
	dispdata = AllApiData{Images: imagesArray, Weather: &WeatherData{Temp: celsiusNum, City: cityLibrary[RANDi], Icon: rainOrShine}}
	renderTemplate(w, "home", dispdata)
}

func renderTemplate(w http.ResponseWriter, tmpl string, structdata AllApiData) {
	t := template.Must(template.New("image").ParseFiles("layout/home.html"))
	if err := t.ExecuteTemplate(w, tmpl, structdata); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("layout"))))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
