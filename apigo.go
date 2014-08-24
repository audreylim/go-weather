package main 

import (
	"fmt"
)

func ShowImage() {
	addr := r.FormValue("str")
	safeAddr := url.QueryEscape(addr)
	fullUrl := fmt.Sprintf("https://api.flickr.com/services/rest?api_key=??&format=json&tags=%s", safeAddr)

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
	if dataReadErr := nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}
	
	res := make the maps////
	json.Unmarshal(body, &res)


}