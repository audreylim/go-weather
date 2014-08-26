package main 

import (
	"fmt"
	//"os"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	"encoding/json"
)

const homeTemplate = `
<html>
<form action="/image" method="POST">
<input type="text" name="str"></input>
<button type="submit" value="submit" name="flickr">tag</button>
</form>
</html>
`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, homeTemplate)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	/*addr := "sun"//r.FormValue("str")
	safeAddr := url.QueryEscape(addr)*/
	fullUrl := "https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=e7ef66cea848474a3e1fe3de117f4670&tags=summer&extras=url_m&per_page=1&format=json&nojsoncallback=1"

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

	var f interface{}
	errr := json.Unmarshal(body, &f)
	if errr != nil {
		log.Fatal(errr)
	}
	fmt.Println(f)

	m := f.(map[string]interface{})





	fmt.Println(m)

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is a string")
		case map[string]interface{}:
			fmt.Println(k, "is a map")
			for i, u := range vv {
            	switch uu := u.(type) {
            	case []map[string]interface{}:
            		fmt.Println(i, "is another map", uu)
        		default:
        			fmt.Println(i, u)
            	}

    
        	}
        	
        default:
        fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

	

/*
fmt.Println("photourl", photourl)

	http.Redirect(w, r, photourl, 302)*/
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/image/", imageHandler)

	http.ListenAndServe(":8000", nil)
}

