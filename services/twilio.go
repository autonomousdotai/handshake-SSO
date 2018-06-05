package services

import (
	"github.com/ninjadotorg/handshake-dispatcher/utils"
	"fmt"
	"bytes"
	"net/http"
	"github.com/ninjadotorg/handshake-dispatcher/config"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

type TwilioService struct {}

func (s TwilioService) SendVerification(countryCode string, phoneNumber string, locale string) (success bool, err error) {
	endpoint, _ := utils.GetServicesEndpoint("twilio")
	uri := fmt.Sprintf("%s/protected/json/phones/verification/start", endpoint)

	c := config.GetConfig()
	apiKey := c.GetString("twilio_key")

	form := url.Values{
		"via": {"sms"},
		"phone_number": {phoneNumber},
		"country_code": {countryCode},
		"locale": {locale},
		"code_length": {"6"},
	}

	body := bytes.NewBufferString(form.Encode())
	request, _ := http.NewRequest("POST", uri, body)
	request.Header.Set("X-Authy-API-Key", apiKey)

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

	success = data["success"].(bool)

	return
}

func (s TwilioService) CheckVerification(countryCode string, phoneNumber string, code string) (success bool, err error) {
	endpoint, _ := utils.GetServicesEndpoint("twilio")
	uri := fmt.Sprintf("%s/protected/json/phones/verification/check", endpoint)
	uri = fmt.Sprintf("%s?phone_number=%s&country_code=%s&verification_code=%s", uri, phoneNumber, countryCode, code)

	c := config.GetConfig()
	apiKey := c.GetString("twilio_key")
	request, _ := http.NewRequest("GET", uri, nil)
	request.Header.Set("X-Authy-API-Key", apiKey)

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

	success = data["success"].(bool)

	return
}