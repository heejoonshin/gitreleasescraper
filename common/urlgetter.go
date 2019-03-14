package common

import (
	"encoding/json"
	"net/http"
)

func GetJSON(url string,json_map interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Printf("%#v\n", resp)

	dec := json.NewDecoder(resp.Body)
	if dec == nil {
		panic("Failed to start decoding JSON data")
	}

	//json_map := make(map[string]interface{})
	err = dec.Decode(&json_map)
	if err != nil {
		panic(err)
	}
	//fmt.Println(json_map["tag_name"])


	//fmt.Printf("%v\n", json_map)
	return nil
}