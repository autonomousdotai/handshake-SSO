package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/models"
	"github.com/ninjadotorg/handshake-dispatcher/services"
	"github.com/ninjadotorg/handshake-dispatcher/utils"
	"github.com/ninjadotorg/handshake-dispatcher/config"
)

type VerifierController struct{}

func (s VerifierController) SendPhoneVerification(c *gin.Context) {
	phone := c.DefaultQuery("phone", "")
	countryCode := c.DefaultQuery("country", "")
	locale := c.DefaultQuery("locale", "en")

	twilioClient := services.TwilioService{}
	success, err := twilioClient.SendVerification(countryCode, phone, locale)
	if err != nil {
		resp := JsonResponse{0, "Send verification failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{0, "Send verification failed", nil}
	if success {
		resp = JsonResponse{1, "", nil}
	}

	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) CheckPhoneVerification(c *gin.Context) {
	phone := c.DefaultQuery("phone", "")
	countryCode := c.DefaultQuery("country", "")
	code := c.DefaultQuery("code", "")

	twilioClient := services.TwilioService{}
	success, err := twilioClient.CheckVerification(countryCode, phone, code)
	if err != nil {
		resp := JsonResponse{0, "Check verification failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{0, "Phone verified failed", nil}
	if success {
		resp = JsonResponse{1, "", nil}
	}

	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) SendEmailVerification(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	isNeedEmail := c.DefaultQuery("isNeedEmail", "1")

	err := utils.ValidateFormat(email)
	if err != nil {
		resp := JsonResponse{0, "Invalid email", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	code := utils.RandomVerificationCode()

	var userModel models.User
	user, _ := c.Get("User")
	userModel = user.(models.User)

	var md map[string]interface{}
	if userModel.Metadata != "" {
		json.Unmarshal([]byte(userModel.Metadata), &md)
	} else {
		md = map[string]interface{}{}
	}

	md["verification-code"] = map[string]interface{}{
		"time": time.Now(),
		"code": code,
	}

	metadata, _ := json.Marshal(md)
	userModel.Metadata = string(metadata)
	userModel.Email = email
	dbErr := models.Database().Save(&userModel).Error
	if dbErr != nil {
		log.Println("Send verification failed", dbErr.Error())
		resp := JsonResponse{0, "Send verification failed", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	mailClient := services.MailService{}

	subject := "Email verification"
	content := fmt.Sprintf(EMAIL_VERIFICATION_TEMPLATE, fmt.Sprint(code))

	if isNeedEmail == "1" {
		go mailClient.Send("dojo@ninja.org", email, subject, content)
	}

	resp := JsonResponse{1, "", nil}
	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) CheckEmailVerification(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	code := c.DefaultQuery("code", "")

	var userModel models.User
	user, _ := c.Get("User")
	userModel = user.(models.User)

	var md map[string]interface{}
	if userModel.Metadata != "" {
		json.Unmarshal([]byte(userModel.Metadata), &md)
	} else {
		md = map[string]interface{}{}
	}

	verificationCode, hasCode := md["verification-code"]

	if !hasCode {
		resp := JsonResponse{0, "Email verified failed", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	realCode := (verificationCode.(map[string]interface{}))["code"]
	if fmt.Sprint(realCode) != fmt.Sprint(code) {
		resp := JsonResponse{0, "Email verified failed", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	delete(md, "verification-code")
	metadata, _ := json.Marshal(md)
	userModel.Metadata = string(metadata)
	userModel.Email = email
	dbErr := models.Database().Save(&userModel).Error
	if dbErr != nil {
		resp := JsonResponse{0, "Email verified failed", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	resp := JsonResponse{1, "", nil}
	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) CheckRedeemCodeVerification(c *gin.Context) {
	code := c.DefaultQuery("code", "")

	conf := config.GetConfig()
	apiVerifyRedeemCode := conf.GetString("autonomous_api")

	endpoint := apiVerifyRedeemCode + "promotion-program-api/verify-promotion-code?promotion_code=%s"
	uri := fmt.Sprintf(endpoint, code)

	request, _ := http.NewRequest("POST", uri, nil)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(b, &data)

	c.JSON(http.StatusOK, data)
}
func (s VerifierController) ActiveRedeemCode(c *gin.Context) {

	// 1 check code first:

	code := c.DefaultQuery("code", "")

	conf := config.GetConfig()
	apiVerifyRedeemCode := conf.GetString("autonomous_api")

	endpoint := apiVerifyRedeemCode + "promotion-program-api/verify-promotion-code?promotion_code=%s"
	uri := fmt.Sprintf(endpoint, code)


	request, _ := http.NewRequest("POST", uri, nil)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(b, &data)


	if data["status"].(float64) != 1 {
		c.JSON(http.StatusOK, data)
		return
	}


	fmt.Println ("active + transfer ========================")

	toAddress := c.DefaultQuery("to-address", "")


	fiatAmountValue := ((data["data"].(map[string]interface{}))["amount"]).(float64)
	currency := c.DefaultQuery("currency", "")

	log.Println("fiatAmountValue", fiatAmountValue)

	if toAddress == "" {
		resp := JsonResponse{0, "to-address invalid", nil}
		c.JSON(http.StatusOK, resp)
		return
	}
	if currency == "" {
		resp := JsonResponse{0, "currency invalid", nil}
		c.JSON(http.StatusOK, resp)
		return
	}



	fmt.Println ("transfer----------------------------------")
	exchangeAPI, found := utils.GetForwardingEndpoint("exchange")
	log.Println(exchangeAPI, found)

	endpoint = exchangeAPI + "/internal/redeem"

	jsonData := make(map[string]interface{})
	jsonData["address"] = toAddress
	jsonData["fiat_amount"] = fiatAmountValue
	jsonData["currency"] = currency
	jsonData["ref_data"] = "wallet-giftcard-redeem"
	jsonValue, _ := json.Marshal(jsonData)

	request, _ = http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))

	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, _ = ioutil.ReadAll(response.Body)

	var dataExchange map[string]interface{}
	json.Unmarshal(b, &dataExchange)

	log.Println("dataExchange=>", dataExchange)

	if dataExchange["status"].(float64) != 1 {
		c.JSON(http.StatusOK, dataExchange)
		return
	}
	c.JSON(http.StatusOK, dataExchange)


	//fmt.Println ("active------------------------------------")
	//endpoint = apiVerifyRedeemCode + "promotion-program-api/redeem-code?promotion_code=%s"
	//uri = fmt.Sprintf(endpoint, code)
	//
	//
	//request, _ = http.NewRequest("POST", uri, nil)
	//
	//client = &http.Client{}
	//response, err = client.Do(request)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//b, _ = ioutil.ReadAll(response.Body)
	//
	//var dataActive map[string]interface{}
	//json.Unmarshal(b, &dataActive)
	//c.JSON(http.StatusOK, dataActive)

}

const EMAIL_VERIFICATION_TEMPLATE = `<html>
<body>
<p>
    Hey Ninja,
</p>
<p>
    Here's your email verification code: <b>%s</b>
</p>
<p>
    Just tap and you're in.
</p>
<p>
    You look like a winner.
</p>
<p>
    Ninja.org<br/>
    Join the conversation at <a href="https://t.me/ninja_org">t.me/ninja_org</a>
</p>
<p>
</p>
</body>
</html>
`
