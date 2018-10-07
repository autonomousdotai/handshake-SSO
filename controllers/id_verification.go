package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/models"
)

type IDVerification struct{}

func (i IDVerification) List(c *gin.Context) {
	db := models.Database()
	var listIDVerification []models.IDVerification

	errDb := db.Where("status = 0").Find(&listIDVerification).Error

	if errDb != nil {
		resp := JsonResponse{0, "Unable to load list", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "Success", &listIDVerification}
	c.JSON(http.StatusOK, resp)
}

func (i IDVerification) UpdateStatus(c *gin.Context) {
	id, convErr := strconv.Atoi(c.DefaultPostForm("id", "-1"))

	if convErr != nil || id < 0 {
		resp := JsonResponse{0, "Invalid id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	status, convErr := strconv.Atoi(c.DefaultPostForm("status", "0"))

	if convErr != nil || status < 0 || status > 1 {
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

	idVerificationItem.Status = status

	var user models.User
	errDb = db.Where("id = ?", idVerificationItem.UserID).First(&user).Error

	if errDb != nil {
		resp := JsonResponse{0, "Could not found user related to this id verification", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	user.IDVerified = status

	if db.Save(&idVerificationItem).Error != nil || db.Save(&user).Error != nil {
		resp := JsonResponse{0, "Could not update status for this id verification. Please try again", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "Success", nil}
	c.JSON(http.StatusOK, resp)
}
