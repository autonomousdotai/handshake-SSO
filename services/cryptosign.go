package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ninjadotorg/handshake-dispatcher/utils"
)

type CryptosignService struct{}

// UserModelHooks :
func (s CryptosignService) UserModelHooks(typeChange string, userID uint, metaData string, email string) {
	endpoint, _ := utils.GetServicesEndpoint("cryptosign")
	jsonData := make(map[string]interface{})
	jsonData["user_id"] = userID
	jsonData["email"] = email
	jsonData["meta_data"] = metaData
	jsonData["type_change"] = typeChange
	jsonValue, _ := json.Marshal(jsonData)

	endpoint = fmt.Sprintf("%s/user/hook/dispatcher", endpoint)
	fmt.Println(endpoint)
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

	if ok && (float64(1) == result) {
		return
	}
	fmt.Println(errors.New(message.(string)))
	return
}
