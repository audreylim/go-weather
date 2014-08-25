package main 

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
)

const homeTemplate = `
<html>
<form action="/image" method="POST">
<input type="text"></input>
<button type="submit" value="submit" name="flickr"></button>
</form>
</html>
`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, homeTemplate)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("str")
	safeAddr := url.QueryEscape(addr)
	fullUrl := fmt.Sprintf("https://api.flickr.com/services/rest?api_key=%s&format=json&tags=%s&content_type=1&nojsoncallback=1", os.Getenv("FLICKR_APIKEY&extras=url_m"), safeAddr)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
		return
	}
	defer resp.Body.Close()

	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}

	res := make(map[string]map[string][]map[string]string)

	json.Unmarshal(body, &res)

	photourl, _ := res["photos"]["photo"][0]["url_m"]

	http.Redirect(w, r, photourl, 302)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/image/", imageHandler)

	http.ListenAndServe(":8000", nil)
}

