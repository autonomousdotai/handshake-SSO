package services

import (
    "fmt"
    "bytes"
    "errors"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "github.com/ninjadotorg/handshake-dispatcher/utils"
)

type FCMService struct {}

func (s FCMService) Notify(jsonData map[string]interface{}) (bool, error) {
    endpoint, _ := utils.GetServicesEndpoint("fcm")
    jsonValue, _ := json.Marshal(jsonData)
    
    endpoint = fmt.Sprintf("%s/send", endpoint)

    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)

    var data map[string]interface{}
    json.Unmarshal(b, &data)

    result, ok := data["status"]
    message, _ := data["message"]

    if ok && (float64(1) == result) {
        return true, nil
    } else {
        return false, errors.New(message.(string)) 
    }
}

