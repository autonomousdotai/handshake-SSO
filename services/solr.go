package services

import (
    "fmt"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
    "github.com/autonomousdotai/handshake-dispatcher/utils"
)

type SolrService struct {}

func (s SolrService) Init() {
    
}

func (s SolrService) List(t string, q []string, offset int, limit int, sort []string) (map[string]interface{}, error) {
    jsonData := make(map[string]interface{})
    
    params := make(map[string]interface{})
    if q != nil {
        params["q"] = q    
    }

    if sort != nil {
        params["sort"] = sort
    }

    jsonData["Params"] = params
    jsonData["Start"] = offset
    jsonData["Rows"] = limit

    endpoint, _ := utils.GetServicesEndpoint("solr")
    jsonValue, _ := json.Marshal(jsonData)
    
    endpoint = fmt.Sprintf("%s/%s/select", endpoint, t)

    request, _ := http.NewRequest("POST", endpoint , bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

    b, _ := ioutil.ReadAll(response.Body)

    var data map[string]interface{}
    json.Unmarshal(b, &data)

    wrapData := make(map[string]interface{})
    handshakes := []map[string]interface{}{}
    
    for k, v := range data["Results"].(map[string]interface{}) {
        if k == "NumFound" {
            wrapData["total"] = v
        }
        if k == "Start" {
            wrapData["index"] = v
        }
        if k == "Collection" {
            collections, colOk := v.([]interface{})
            if colOk {
                for _, item := range collections {
                    collection := item.(map[string]interface{})
                    fields := collection["Fields"].(map[string]interface{})
                    handshake := make(map[string]interface{})
                    for k3, v3 := range fields {
                        if k3 != "_version_" {
                            handshake[CleanSolrName(k3)] = v3;
                        }
                    }
                    handshakes = append(handshakes, handshake)
                }
            }
        }
    }
    wrapData["handshakes"] = handshakes

    return wrapData, nil
}

func (s SolrService) Create(t string, d map[string]interface{}) (bool, error) {
    jsonData := make(map[string]interface{})
    jsonData["add"] = []map[string]interface{}{d}

    endpoint, _ := utils.GetServicesEndpoint("solr")
    jsonValue, _ := json.Marshal(jsonData)
    
    endpoint = fmt.Sprintf("%s/%s/update", endpoint, t)
 
    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    
    response, err := client.Do(request)
    if err != nil {
        fmt.Println("Error response", err)
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, &data)

    result, ok := data["Success"]
    
    if ok {
        fmt.Println("has key Success")
        return result.(bool), nil
    } else {
        return false, nil
    }
}

func (s SolrService) Update(t string, d map[string]interface{}) (bool, error) {
    jsonData := make(map[string]interface{})
    jsonData["add"] = []map[string]interface{}{d}

    endpoint, _ := utils.GetServicesEndpoint("solr")
    jsonValue, _ := json.Marshal(jsonData)
   
    endpoint = fmt.Sprintf("%s/%s/update", endpoint, t)

    fmt.Println(endpoint, jsonData)

    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, &data)

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
    
    endpoint, _ := utils.GetServicesEndpoint("solr")
    jsonValue, _ := json.Marshal(jsonData)
 
    endpoint = fmt.Sprintf("%s/%s/update", endpoint, t)

    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return false, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    var data map[string]interface{}
    json.Unmarshal(b, &data)

    result, ok := data["Success"]

    if ok {
        return result.(bool), nil
    } else {
        return false, nil
    }
}

func CleanSolrName(name string) string {
    ignoreSuffixes := []string{"_i", "_is", "_s", "_ss", "_l", "_ls", "_t", "_txt", "_b", "_bs", "_f", "_fs", "_d", "_ds", "_str", "_dt", "_dts", "_p", "_srpt", "_dpf", "_dpi", "_dps"}

    ignorePrefixes := []string{"attr_"}
 
    result := name
    for _, suffix := range(ignoreSuffixes) {
        if strings.HasSuffix(result, suffix) {
            result = result[:len(result)-len(suffix)]
        }
    }

    for _, prefix := range(ignorePrefixes) {
        if strings.HasPrefix(result, prefix) {
            result = result[len(prefix):len(result)-len(prefix)]
        }
    }
    return result
}
