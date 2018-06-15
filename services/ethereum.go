package services

import (
	"github.com/ninjadotorg/handshake-dispatcher/utils"
	"io/ioutil"
	"net/http"
	"fmt"
	"encoding/json"
)

type EthereumService struct{}

func (s EthereumService) FreeEther(userId string, to string, value string, networkId string) (bool, string) {
	endpoint, _ := utils.GetServicesEndpoint("ethereum")

    endpoint = fmt.Sprintf("%s/free-ether?to_address=%s&value=%s&network_id=%s", endpoint, to, value, networkId)

	request, _ := http.NewRequest("POST", endpoint, nil)
	request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Uid", userId)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return false, err.Error()
	}

	b, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(b, &data)

	status, ok := data["status"]
    message, _ := data["message"]
    
    if ok && (float64(1) == status) {
        rData := data["data"].(map[string]interface{})
        return true, rData["hash"].(string)
    } else {
        return false, message.(string)
    }
}

func (s EthereumService) FreeToken(userId string, to string, amount string, networkId string) (bool, string) {
	endpoint, _ := utils.GetServicesEndpoint("ethereum")

    endpoint = fmt.Sprintf("%s/free-token?to_address=%s&amount=%s&network_id=%s", endpoint, to, amount, networkId)

	request, _ := http.NewRequest("POST", endpoint, nil)
	request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Uid", userId)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return false, err.Error() 
	}

	b, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(b, &data)

	status, ok := data["status"]
    message, _ := data["message"]
    
    if ok && (float64(1) == status) {
        rData := data["data"].(map[string]interface{})
        return true, rData["hash"].(string)
    } else {
        return false, message.(string)
    }
}
