package controllers

import (
	"net/http"
	"strconv"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/models"
)

type StoreController struct{}


func (u StoreController) Create(c *gin.Context) {

	status, convErr := strconv.Atoi(c.DefaultPostForm("status", "1"))

	if convErr != nil {
		resp := JsonResponse{0, "Invalid id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}


	wallets_receive := c.DefaultPostForm("wallets_receive", "_")

	if wallets_receive == "_" {
		resp := JsonResponse{0, "Invalid wallets_receive", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	confirm_url := c.DefaultPostForm("confirm_url", "_")

	if confirm_url == "_" {
		resp := JsonResponse{0, "Invalid confirm_url", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	store_id := c.DefaultPostForm("store_id", "_")
	if store_id == "_" {
		resp := JsonResponse{0, "Invalid store_id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	store_name := c.DefaultPostForm("store_name", "_")
	if store_name == "_" {
		resp := JsonResponse{0, "Invalid store_name", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// check name:
	var _store models.Store
	errDb := models.Database().Where("store_id = ?", store_id).First(&_store).Error

	if errDb == nil {
		resp := JsonResponse{0, "store_id already exists ", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	db := models.Database()

	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)

	store := models.Store{
				Status: status,
				StoreName: store_name,
				StoreID: store_id,
				WalletsReceive: wallets_receive,
				ConfirmURL: confirm_url,
				UserID: userModel.ID,
	}

	errDb = db.Save(&store).Error

	if errDb != nil {
		resp := JsonResponse{0, "Unable to create your store", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "", nil}
	c.JSON(http.StatusOK, resp)
}


func (i StoreController) List(c *gin.Context) {

	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)

	db := models.Database()
	var listStore []models.Store

	errDb := db.Where("user_id = ?", userModel.ID).Find(&listStore).Error

	if errDb != nil {
		resp := JsonResponse{0, "Unable to load list", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "Success", &listStore}
	c.JSON(http.StatusOK, resp)
}

func (i StoreController) Detail(c *gin.Context) {

	store_id := c.DefaultQuery("store_id", "_")

	if store_id == "_" {
		resp := JsonResponse{0, "Invalid store_id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}
	db := models.Database()
	var store models.Store
	errDb := db.Where("store_id = ?", store_id).First(&store).Error

	if errDb != nil {
		resp := JsonResponse{0, "Invalid store_id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}
	resp := JsonResponse{1, "Success", &store}
	c.JSON(http.StatusOK, resp)

}


func (i StoreController) UpdateStore(c *gin.Context) {

	id, convErr := strconv.Atoi(c.DefaultPostForm("id", "-1"))

	if convErr != nil || id < 0 {
		resp := JsonResponse{0, "Invalid id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	wallets_receive := c.DefaultPostForm("wallets_receive", "_")
	confirm_url := c.DefaultPostForm("confirm_url", "_")
	store_id := c.DefaultPostForm("store_id", "_")
	store_name := c.DefaultPostForm("store_name", "_")
	status, _ := strconv.Atoi(c.DefaultPostForm("status", "-1"))


	db := models.Database()
	var store models.Store
	errDb := db.Where("id = ?", id).First(&store).Error

	if errDb != nil {
		resp := JsonResponse{0, "Invalid id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	if status > -1 {
		store.Status = status
	}
	if wallets_receive != "_" {
		store.WalletsReceive = wallets_receive
	}
	if confirm_url != "_" {
		store.ConfirmURL = confirm_url
	}
	if store_id != "_" {

		var _store models.Store
		errDb := models.Database().Where("store_id = ? AND id != ?", store_id, store.ID).First(&_store).Error

		if errDb == nil {
			resp := JsonResponse{0, "store_id already exists ", nil}
			c.JSON(http.StatusOK, resp)
			return
		} else {
			store.StoreID = store_id
		}
	}
	if store_name != "_" {
		store.StoreName = store_name
	}


	var user models.User
	errDb = db.Where("id = ?", store.UserID).First(&user).Error

	if errDb != nil {
		log.Println("Error", errDb.Error())
		resp := JsonResponse{0, "Could not found user related to this store", nil}
		c.JSON(http.StatusOK, resp)
		return
	}


	if db.Save(&store).Error != nil  {
		resp := JsonResponse{0, "Could not update status for this store. Please try again", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "Success", nil}
	c.JSON(http.StatusOK, resp)
}
