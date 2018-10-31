package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/config"
	"github.com/ninjadotorg/handshake-dispatcher/models"
	"github.com/ninjadotorg/handshake-dispatcher/services"
)

type IDVerification struct{}

type IDVerificationEmail struct {
	Subject string
	Content string
}

var emailContent = [3]IDVerificationEmail{{
	Subject: "Ninja Coin - Successful ID Verification",
	Content: `<p>Dear %s</p>
	<p>Thank you for signing up with <a href="%s/coin">Ninja Coin</a>!</p>
	<p>Your ID has been verified. Your current daily limit is 500 USD with various payment method such as cash, credit card and bank transfer.</p>
	<p>If you wish to upgrade to your next tier, which is up to 5000 USD daily, please <a href="%s/me/profile">submit your selfie here</a>.</p>
	<p>This is not you who verify the account? Inform us immediately!</p>
	<p>Have fun at the dojo!<br />Ninja Team</p>`,
}, {
	Subject: "Ninja Coin - Successful Account Verification",
	Content: `<p>Dear %s</p>
	<p>Your account verification is completed. Now you can buy up to 5,000 USD everyday with various payment method such as cash, credit card and bank transfer.</p>
	<p>If you are receiving this email and have never signed up with us, please inform us immediately.</p>
	<p>Have fun at the dojo!<br />Ninja Team</p>`,
}, {
	Subject: "Your account verification has been rejected",
	Content: `<p>Dear %s,</p>
	<p>We are sorry inform you that your account verification has been rejected. The reason for this is %s. Please <a href="%s/me/profile">click here</a> to try again.</p>
	<p>If you think this is a mistake, please contact our hotline +84975504082 for further assistance</p>
	<p>Have fun at the dojo<br />Ninja Team</p>`,
}}

func (i IDVerification) List(c *gin.Context) {
	db := models.Database()
	var listIDVerification []models.IDVerification
	status := 0

	if filterStatus, err := strconv.Atoi(c.Query("status")); err == nil {
		status = filterStatus
	}

	query := db.Where("status = ?", status)

	if filterUserID, err := strconv.Atoi(c.Query("uid")); err == nil {
		query = query.Where("user_id = ?", filterUserID)
	}

	errDb := query.Find(&listIDVerification).Error

	if errDb != nil {
		resp := JsonResponse{0, "Unable to load list", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "Success", &listIDVerification}
	c.JSON(http.StatusOK, resp)
}

func (i IDVerification) Get(c *gin.Context) {
	var userModel models.User
	user, _ := c.Get("User")
	userModel = user.(models.User)

	db := models.Database()
	var existsIDVerification models.IDVerification
	existsIDVerificationErr := db.Where("user_id = ?", userModel.ID).First(&existsIDVerification).Error

	if existsIDVerificationErr != nil {
		resp := JsonResponse{0, "Not found", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "Success", &existsIDVerification}
	c.JSON(http.StatusOK, resp)
}

func (i IDVerification) UpdateStatus(c *gin.Context) {
	id, convErr := strconv.Atoi(c.DefaultPostForm("id", "-1"))
	conf := config.GetConfig()

	if convErr != nil || id < 0 {
		resp := JsonResponse{0, "Invalid id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	status, convErr := strconv.Atoi(c.DefaultPostForm("status", "0"))

	if convErr != nil || status < -1 || status > 1 {
		resp := JsonResponse{0, "Invalid status", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	db := models.Database()
	var idVerificationItem models.IDVerification
	errDb := db.Where("id = ?", id).First(&idVerificationItem).Error

	if errDb != nil {
		resp := JsonResponse{0, "Invalid id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	var user models.User
	errDb = db.Where("id = ?", idVerificationItem.UserID).First(&user).Error

	if errDb != nil {
		resp := JsonResponse{0, "Could not found user related to this id verification", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	idVerificationItem.Status = status
	user.IDVerified = status
	if status == 1 {
		user.IDVerificationLevel++
		idVerificationItem.Level++
	}

	if db.Save(&idVerificationItem).Error != nil || db.Save(&user).Error != nil {
		resp := JsonResponse{0, "Could not update status for this id verification. Please try again", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	userEmail := idVerificationItem.Email
	userFullName := idVerificationItem.Name
	if userEmail != "" && user.IDVerificationLevel > 0 {
		mailClient := services.MailService{}
		var subject string
		var content string
		workingDomain := conf.GetString("working_domain")
		if status == 1 {
			emailContentToSend := emailContent[user.IDVerificationLevel-1]
			subject = emailContentToSend.Subject
			if user.IDVerificationLevel == 1 {
				content = fmt.Sprintf(emailContentToSend.Content, userFullName, workingDomain, workingDomain)
			} else {
				content = fmt.Sprintf(emailContentToSend.Content, userFullName)
			}
		} else if status == -1 {
			emailContentToSend := emailContent[2]
			subject = emailContentToSend.Subject
			content = fmt.Sprintf(emailContentToSend.Content, userFullName, "Invalid Documents", workingDomain)
		}
		if subject != "" {
			go mailClient.Send("dojo@ninja.org", userEmail, subject, content)
		}
	}

	resp := JsonResponse{1, "Success", nil}
	c.JSON(http.StatusOK, resp)
}
