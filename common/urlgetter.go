package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

func GetJSON(url string, json_map interface{}) error {
	resp, err := http.Get(url)

	if resp.StatusCode != 200 {
		return errors.New("데이터를 불러오지 못했습니다")
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//fmt.Printf("%#v\n", resp)

	dec := json.NewDecoder(resp.Body)
	if dec == nil {
		return errors.New("Failed to start decoding JSON data")
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
