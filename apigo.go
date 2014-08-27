package main

import (
	"fmt"
	//"os"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	"html/template"
	"encoding/json"
)

var upperhomeTemplate = template.Must(template.New("image").ParseFiles("layout/home.html"))


func homeHandler(w http.ResponseWriter, r *http.Request) {
tempErr := upperhomeTemplate.ExecuteTemplate(w, "home", nil)
if tempErr != nil {
	http.Error(w, tempErr.Error(), http.StatusInternalServerError)
}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	/*addr := "sun"//r.FormValue("str")
	safeAddr := url.QueryEscape(addr)*/
	fullUrl := "https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=e7ef66cea848474a3e1fe3de117f4670&tags=newyork,travel&extras=url_m&format=json&nojsoncallback=1&min_taken_date=1388534400&sort=relevance"

	fmt.Println("fullUrl", fullUrl)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
		return
	}
	defer resp.Body.Close()
	fmt.Println("resp:", resp)
	fmt.Println("respbody", resp.Body)

	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}
	fmt.Println("body:", body)

	type FlickrResponse struct {
		Photos struct {
			Photo []struct{ Url_m string }
		}
	}

	var f FlickrResponse

	errr := json.Unmarshal(body, &f)
	if errr != nil {
		log.Fatal(errr)
	}

	type ImageDisplay struct {
		Images []string
	}

	var a []string


	for i:=0;i<4;i++{
	photourl := f.Photos.Photo[i].Url_m
	a = append(a, photourl)
	}

	b := ImageDisplay{Images: a}




tempErr := upperhomeTemplate.ExecuteTemplate(w, "home", b)
if tempErr != nil {
	http.Error(w, tempErr.Error(), http.StatusInternalServerError)
}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/image/", imageHandler)

	http.ListenAndServe(":8000", nil)
}
