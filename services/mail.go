package services

import (
	"github.com/ninjadotorg/handshake-dispatcher/utils"
	"bytes"
	"io/ioutil"
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
)

type MailService struct {}

func (s MailService) Send(from string, to string, subject string, content string) (success bool, err error) {
	endpoint, _ := utils.GetServicesEndpoint("mail")

	form := url.Values{
		"from": {from},
		"to[]": {to},
		"subject": {subject},
		"body": {content},
	}

	body := bytes.NewBufferString(form.Encode())
	request, _ := http.NewRequest("POST", endpoint, body)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(b, &data)

	fmt.Println(data)

	success = data["status"].(int) == 1

	return
}
