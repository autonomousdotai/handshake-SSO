package services

import (
	"github.com/ninjadotorg/handshake-dispatcher/utils"
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"
	"encoding/json"
	"mime/multipart"
)

type MailService struct{}

func (s MailService) Send(from string, to string, subject string, content string) (success bool, err error) {
	endpoint, _ := utils.GetServicesEndpoint("mail")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormField("from")
	part.Write([]byte(from))
	part, _ = writer.CreateFormField("to[]")
	part.Write([]byte(to))
	part, _ = writer.CreateFormField("subject")
	part.Write([]byte(subject))
	part, _ = writer.CreateFormField("body")
	part.Write([]byte(content))
	writer.Close()

	request, _ := http.NewRequest("POST", endpoint, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

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
	status, _ := data["status"].(float64)
	success = status >= 1

	return
}
