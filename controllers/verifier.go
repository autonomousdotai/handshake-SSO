package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/models"
	"github.com/ninjadotorg/handshake-dispatcher/services"
	"github.com/ninjadotorg/handshake-dispatcher/utils"
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
		"email": email,
		"time":  time.Now(),
		"code":  code,
	}

	metadata, _ := json.Marshal(md)
	userModel.Metadata = string(metadata)
	userModel.Email = email
	userModel.Verified = 0
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

	realEmail := (verificationCode.(map[string]interface{}))["email"]
	if fmt.Sprint(realEmail) != fmt.Sprint(email) {
		resp := JsonResponse{0, "Email verified failed", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	delete(md, "verification-code")
	metadata, _ := json.Marshal(md)
	userModel.Metadata = string(metadata)
	userModel.Email = email
	userModel.Verified = 1

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
	Join us on Telegram: <a href="https://t.me/ninja_org">t.me/ninja_org</a><br/>
	Find us on <a href="https://www.facebook.com/ninjadotorg">Facebook</a> | <a href="https://twitter.com/ninjadotorg">Twitter</a>
</p>
<p>
</p>
</body>
</html>
`
