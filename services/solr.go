package services

import (
    "fmt"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "github.com/autonomousdotai/handshake-dispatcher/config"
)

type SolrService struct {}

func (s SolrService) List(t string, q []string, offset int, limit int) (map[string]interface{}, error) {
    jsonData := make(map[string]interface{})
    if q != nil {
        params := make(map[string]interface{})
        params["q"] = q
        
        jsonData["params"] = params
    }
    jsonData["Start"] = offset
    jsonData["Rows"] = limit

    endpoint := GetSolrEndpoint(t)
    jsonValue, _ := json.Marshal(jsonData)
    
    request, _ := http.NewRequest("GET", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, data)

    return data, nil
}

func (s SolrService) Create(t string, d map[string]string) (bool, error) {
    jsonData := make(map[string]interface{})
    add := make(map[string]string)
    for k, v := range d {
        add[k] = v
    }
    jsonData["add"] = add

    endpoint := GetSolrEndpoint(t)
    jsonValue, _ := json.Marshal(jsonData)
    
    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, data)

    result, ok := data["Success"]
    
    if ok {
        return result.(bool), nil
    } else {
        return false, nil
    }
}

func (s SolrService) Update(t string, d map[string]string) (bool, error) {
    jsonData := make(map[string]interface{})
    update := make(map[string]string)
    for k, v := range d {
        update[k] = v
    }
    jsonData["update"] = update

    endpoint := GetSolrEndpoint(t)
    jsonValue, _ := json.Marshal(jsonData)
    
    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, data)

    result, ok := data["Success"]

    if ok {
        return result.(bool), nil
    } else {
        return false, nil
    }
}

func (s SolrService) Delete(t string, id string) (bool, error) {
    jsonData := make(map[string]interface{})
    delete := make(map[string]string)
    delete["id"] = id
    jsonData["delete"] = delete

    endpoint := GetSolrEndpoint(t)
    jsonValue, _ := json.Marshal(jsonData)
    
    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, data)

    result, ok := data["Success"]

    if ok {
        return result.(bool), nil
    } else {
        return false, nil
    }
}

func GetSolrEndpoint(t string) string {
    conf := config.GetConfig()
    var endpoint string
    
    for ex, ep := range conf.GetStringMap("services") {
        if ex == "solr" {
            endpoint = ep.(string)
            break
        }
    }

    if len(endpoint) > 0 {
        endpoint = fmt.Sprintf("%s/%s", endpoint, t)    
    }

    return endpoint
}
