package controllers

import (
	"net/http"
	"strconv"
	"log"
	"math"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/models"
)

type ActivityLogController struct{}

const (
	defaultLimit = "100"
	defaultPage  = "1"
	defaultOrder = "desc"
)

type Pagination struct {
	Limit  int
	Page   int
	LastID int
	Order  string
}

func (u ActivityLogController) Create(c *gin.Context) {


	name := c.DefaultPostForm("name", "_")
	if name == "_" {
		resp := JsonResponse{0, "Invalid name", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	action := c.DefaultPostForm("action", "_")
	if action == "_" {
		resp := JsonResponse{0, "Invalid store_name", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	description := c.DefaultPostForm("description", "")
	path := c.DefaultPostForm("path", "")
	host := c.DefaultPostForm("host", "")
	method := c.DefaultPostForm("method", "")
	user_agent := c.DefaultPostForm("user_agent", "")



	db := models.Database()

	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)

	activity_log := models.ActivityLog{
				Name: name,
				Action: action,
				Description: description,
				Path: path,
				Host: host,
				Method: method,
				UserAgent: user_agent,
				UserID: userModel.ID,
	}

	errDb := db.Save(&activity_log).Error

	if errDb != nil {
		resp := JsonResponse{0, "Unable to create your activity_log", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "", nil}
	c.JSON(http.StatusOK, resp)
}


func (i ActivityLogController) List(c *gin.Context) {

	limitQuery := c.DefaultQuery("limit", defaultLimit)
	pageQuery := c.DefaultQuery("page", defaultPage)

	lastIDQuery := c.Query("last_id")

	pOrder := c.DefaultQuery("order", defaultOrder)

	log.Println("order", pOrder)


	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		resp := JsonResponse{0, "Unable to load list", err}
		c.JSON(http.StatusOK, resp)
		return
	}

	pLimit := int(math.Max(1, math.Min(10000, float64(limit))))

	db := models.Database()
	var listActivityLog []models.ActivityLog

	count := 0
	db.Model(&models.ActivityLog{}).Count(&count)

	log.Println("count", count)

	if lastIDQuery != "" {
		lastID, err := strconv.Atoi(lastIDQuery)
		if err != nil {
			log.Println("err", err.Error())
			resp := JsonResponse{0, "Unable to load list", nil}
			c.JSON(http.StatusOK, resp)
			return
		}

		pLastID := int(math.Max(0, float64(lastID)))

		if pOrder == "asc" {
			errDb := db.Where("id > ?", pLastID).Limit(pLimit).Order("id asc").Find(&listActivityLog).Error
			if errDb != nil {
				resp := JsonResponse{0, "Unable to load list", nil}
				c.JSON(http.StatusOK, resp)
				return
			}
			resp := JsonResponse{1, "Success", &listActivityLog}
			c.JSON(http.StatusOK, resp)
		}

		errDb := db.Where("id < ?", pLastID).Limit(pLimit).Order("id desc").Find(&listActivityLog).Error
		if errDb != nil {
			resp := JsonResponse{0, "Unable to load list", nil}
			c.JSON(http.StatusOK, resp)
			return
		}
		resp := JsonResponse{1, "Success", &listActivityLog}
		c.JSON(http.StatusOK, resp)
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		resp := JsonResponse{0, "Unable to load list", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	pPage := int(math.Max(1, float64(page)))

	errDb := db.Offset(limit * (pPage - 1)).Limit(pLimit).Order("id " + pOrder).Find(&listActivityLog).Error
	if errDb != nil {
		resp := JsonResponse{0, "Unable to load list", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	var data = make(map[string]interface{})
	data["page"] = page
    	data["page_size"] = count
	data["data"] = &listActivityLog

	resp := JsonResponse{1, "Success", data}
	c.JSON(http.StatusOK, resp)

}

func (i ActivityLogController) Detail(c *gin.Context) {

	id, _ := strconv.Atoi(c.DefaultPostForm("type", "-1"))

	if id == -1 {
		resp := JsonResponse{0, "Invalid activity log", nil}
		c.JSON(http.StatusOK, resp)
		return
	}
	db := models.Database()
	var activityLog models.ActivityLog
	errDb := db.Where("id = ?", id).First(&activityLog).Error

	if errDb != nil {
		resp := JsonResponse{0, "Invalid log id", nil}
		c.JSON(http.StatusOK, resp)
		return
	}
	resp := JsonResponse{1, "Success", &activityLog}
	c.JSON(http.StatusOK, resp)
}
