package services

import (
    "fmt"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
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
    
    endpoint = fmt.Sprintf("%s/select", endpoint)

    request, _ := http.NewRequest("POST", endpoint , bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

    b, _ := ioutil.ReadAll(response.Body)
    
    fmt.Println(string(b))

    var data map[string]interface{}
    json.Unmarshal(b, &data)

    wrapData := make(map[string]interface{})
    handshakes := []map[string]interface{}{}
    
    fmt.Println("Start parse Results")
    for k, v := range data["Results"].(map[string]interface{}) {
        fmt.Println("Parse key", k)
        if k == "NumFound" {
            wrapData["total"] = v
        }
        if k == "Start" {
            wrapData["index"] = v
        }
        if k == "Collection" {
            collections, colOk := v.([]map[string]interface{})
            fmt.Println("Start parse Collections", colOk)
            if colOk {
                for _, collection := range collections {
                    fmt.Println("start parse collection item")
                    fields := collection["Fields"].(map[string]interface{})
                    handshake := make(map[string]interface{})
                    fmt.Println("start extract data")
                    for k3, v3 := range fields {
                        fmt.Println("extract field", k3);
                        if k3 != "version" {
                            handshake[CleanSolrName(k3)] = v3;
                        }
                    }
                    fmt.Println("add to handshakes")
                    handshakes = append(handshakes, handshake)
                }
            }
        }
    }
    fmt.Println("end parse result")
    wrapData["handshakes"] = handshakes

    return wrapData, nil
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
 
    endpoint = fmt.Sprintf("%s/update", endpoint)

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
   
    endpoint = fmt.Sprintf("%s/update", endpoint)

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
 
    endpoint = fmt.Sprintf("%s/update", endpoint)

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
