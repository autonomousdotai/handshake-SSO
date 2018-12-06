package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ninjadotorg/handshake-dispatcher/config"
)

type HookService struct{}

// UserModelHooks :
func (s HookService) UserModelHooks(typeChange string, userID uint, metaData string, email string, name string) {
	jsonData := make(map[string]interface{})
	jsonData["user_id"] = userID
	jsonData["name"] = name
	jsonData["email"] = email
	jsonData["meta_data"] = metaData
	jsonData["type_change"] = typeChange
	jsonValue, _ := json.Marshal(jsonData)

	c := config.GetConfig()
	services := c.GetStringMapString("user_hook_services")
	// Send all user model's event to services
	for key, value := range services {
		fmt.Printf("Start call to hook services '%s': %s \n", key, value)
		endpoint := value
		request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		b, _ := ioutil.ReadAll(response.Body)

		var data map[string]interface{}
		json.Unmarshal(b, &data)
		fmt.Println("====== Result ======")
		fmt.Println(data)
		result, ok := data["status"]
		message, _ := data["message"]

		if !ok || !(float64(1) == result) {
			fmt.Println(errors.New(message.(string)))
			return
		}
	}
	return
}
