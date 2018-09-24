package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ninjadotorg/handshake-dispatcher/config"
	"github.com/ninjadotorg/handshake-dispatcher/daos"
	"github.com/ninjadotorg/handshake-dispatcher/models"
	"github.com/ninjadotorg/handshake-dispatcher/utils"
)

type UserController struct{}

func (u UserController) SignUp(c *gin.Context) {
	config := config.GetConfig()
	UUID, passpharse, err := utils.HashNewUID(config.GetString("secret_key"))

	if err != nil {
		resp := JsonResponse{0, "Sign up failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	ref := c.Query("ref")

	db := models.Database()

	user := models.User{UUID: UUID, Username: UUID}
	if ref != "" {
		refUser := models.User{}
		refErr := db.Where("username = ?", ref).First(&refUser).Error

		if refErr == nil {
			user.RefID = refUser.ID
		}
	}

	errDb := db.Create(&user).Error

	if errDb != nil {
		resp := JsonResponse{0, "Sign up failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// implement another logic
	go ExchangeSignUp(user.ID, user.RefID)

	resp := JsonResponse{1, "", map[string]interface{}{"passpharse": passpharse}}
	c.JSON(http.StatusOK, resp)
	return
}

func (u UserController) Profile(c *gin.Context) {
	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)
	userModel.UUID = ""

	resp := JsonResponse{1, "", userModel}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) Username(c *gin.Context) {
	userId := c.Param("id")
	user := models.User{}
	errDb := models.Database().Where("id = ?", userId).First(&user).Error

	if errDb != nil {
		resp := JsonResponse{0, "Can't found user", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	resp := JsonResponse{1, "", user.Username}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) UsernameExist(c *gin.Context) {
	username := c.DefaultQuery("username", "_")

	if username == "_" {
		resp := JsonResponse{0, "Invalid Username", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	var userModel models.User
	user, _ := c.Get("User")
	userModel = user.(models.User)

	var _u models.User
	errDb := models.Database().Where("username = ? AND id != ?", username, userModel.ID).First(&_u).Error

	var result bool

	if errDb != nil {
		log.Println("Error", errDb.Error())
		result = false
	} else {
		result = true
	}

	resp := JsonResponse{1, "", result}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) UpdateProfile(c *gin.Context) {
	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)

	email := c.DefaultPostForm("email", "_")
	name := c.DefaultPostForm("name", "_")
	username := c.DefaultPostForm("username", "_")
	rwas := c.DefaultPostForm("reward_wallet_addresses", "_")
	was := c.DefaultPostForm("wallet_addresses", "_")
	phone := c.DefaultPostForm("phone", "_")
	ft := c.DefaultPostForm("fcm_token", "_")
	avatar, avatarErr := c.FormFile("avatar")

	log.Println(email, name, username, rwas, phone, ft)

	oldUsername := userModel.Username

	if email != "_" {
		userModel.Email = email
	}
	if username != "_" {
		userModel.Username = username
	}
	if name != "_" {
		userModel.Name = name
	}
	if rwas != "_" {
		userModel.RewardWalletAddresses = rwas
	}
	if was != "_" {
		userModel.WalletAddresses = was
	}
	if phone != "_" {
		userModel.Phone = phone
	}
	if ft != "_" {
		userModel.FCMToken = ft
	}

	if avatarErr == nil {
		uploadImageFolder := "user"
		fileName := avatar.Filename
		imageExt := strings.Split(fileName, ".")[1]
		fileNameImage := fmt.Sprintf("avatar-%d-image-%s.%s", userModel.ID, time.Now().Format("20060102150405"), imageExt)
		path := uploadImageFolder + "/" + fileNameImage

		success, _ := uploadService.Upload(path, avatar)
		if !success {
			resp := JsonResponse{0, "Update profile failed: upload file error", nil}
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		userModel.Avatar = path
	}

	db := models.Database()
	dbErr := db.Save(&userModel).Error

	if dbErr != nil {
		log.Println("Error", dbErr.Error())
		resp := JsonResponse{0, "Update profile failed.", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	// implement another logic
	if oldUsername != userModel.Username {
		go ExchangeUpdateProfile(userModel.ID, userModel.Username)
	}

	userModel.UUID = ""
	log.Println(userModel)
	resp := JsonResponse{1, "", userModel}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) FreeRinkebyEther(c *gin.Context) {
	var userModel models.User
	user, _ := c.Get("User")
	userModel = user.(models.User)

	address := c.DefaultQuery("address", "_")

	if address == "_" {
		resp := JsonResponse{0, "Invalid address", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	var md map[string]interface{}
	if userModel.Metadata != "" {
		json.Unmarshal([]byte(userModel.Metadata), &md)
	} else {
		md = map[string]interface{}{}
	}

	var status bool
	var message string
	shouldRequest := false

	rinkeby, ok := md["free-rinkeby"]
	if ok {
		status = false
		message = fmt.Sprintf("Your free eth transaction is %s", rinkeby.(map[string]interface{})["hash"])
	} else {
		shouldRequest = true
	}

	if shouldRequest {
		value := "1"
		status, message = ethereumService.FreeEther(fmt.Sprint(userModel.ID), address, value, "rinkeby")
		if status {
			md["free-rinkeby"] = map[string]interface{}{
				"address": address,
				"value":   value,
				"hash":    message,
				"time":    time.Now().UTC().Unix(),
			}

			metadata, _ := json.Marshal(md)
			userModel.Metadata = string(metadata)
			dbErr := models.Database().Save(&userModel).Error
			if dbErr != nil {
				status = false
				message = dbErr.Error()
			} else {
				status = true
			}
		}
	}

	resp := JsonResponse{1, message, status}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) CompleteProfile(c *gin.Context) {
	var status bool
	var message string
	var user models.User

	userModel, _ := c.Get("User")
	user = userModel.(models.User)

	conf := config.GetConfig()

	env := conf.GetString("env")
	network := "rinkeby"
	if env == "prod" {
		network = "mainnet"
	}

	log.Println("Start after update profile", user.ID)

	status = false
	// valid user
	if user.Email != "" {
		var md map[string]interface{}
		if user.Metadata != "" {
			json.Unmarshal([]byte(user.Metadata), &md)
		} else {
			md = map[string]interface{}{}
		}

		completeProfile, ok := md["complete-profile"]
		// not received token.
		if !ok {
			log.Println("Yay, User don't receive token yet")
			var wallets map[string]interface{}
			if user.RewardWalletAddresses != "" {
				log.Println("Yay, User have reward wallet address", user.RewardWalletAddresses)
				json.Unmarshal([]byte(user.RewardWalletAddresses), &wallets)

				ethWallet, hasEthWallet := wallets["ETH"]

				if hasEthWallet {
					log.Println("Yay, User has eth wallet.")
					amount := "80"
					fmt.Println("WTF 11")
					address := ((ethWallet.(map[string]interface{}))["address"]).(string)
					fmt.Println("WTF 1")
					tokenStatus, hash := ethereumService.FreeToken(fmt.Sprint(user.ID), address, amount, network)
					log.Println("Receive token result", tokenStatus, hash)
					if tokenStatus {
						md["complete-profile"] = map[string]interface{}{
							"address": address,
							"amount":  amount,
							"hash":    hash,
							"time":    time.Now().UTC().Unix(),
						}

						metadata, _ := json.Marshal(md)
						user.Metadata = string(metadata)
						dbErr := models.Database().Save(&user).Error
						if dbErr != nil {
							log.Println(dbErr.Error())
							message = fmt.Sprintf("Complete Profile Token fail: %s", hash)
						} else {
							status = true
							message = fmt.Sprintf("Your complete profile token transaction is %s", hash)

							go mailService.SendCompleteProfile(user.Email, user.Username, hash)

							if user.RefID != 0 {
								log.Println("This user has referrer", user.RefID)
								go FreeTokenReferrer(fmt.Sprint(user.ID), fmt.Sprint(user.RefID), network)
							}
						}
					} else {
						message = fmt.Sprintf("Complete Profile Token fail: %s", hash)
					}
				} else {
					message = "User does not have ETH reward wallet"
				}
			} else {
				message = "User is not updated reward wallet addresses"
			}
		} else {
			message = fmt.Sprintf("Your complete profile token transaction is %s", completeProfile.(map[string]interface{})["hash"])
		}
	} else {
		message = "User is not complete profile yet"
	}

	resp := JsonResponse{1, message, status}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) Referred(c *gin.Context) {
	var user models.User

	userModel, _ := c.Get("User")
	user = userModel.(models.User)

	var md map[string]interface{}
	if user.Metadata != "" {
		json.Unmarshal([]byte(user.Metadata), &md)
	} else {
		md = map[string]interface{}{}
	}

	referral_total := 0
	referral_amount := 0
	firstbet_total := 0
	firstbet_amount := 0

	referrals, hasReferrals := md["referrals"]

	if hasReferrals {
		for key, _ := range referrals.(map[string]interface{}) {
			if strings.HasPrefix(key, "bonus") {
				referral_total += 1
				referral_amount += 20
			}
			if strings.HasPrefix(key, "firstbet") {
				firstbet_total += 1
				firstbet_total += 20
			}
		}
	}

	data := map[string]interface{}{
		"referral": map[string]interface{}{
			"total":  referral_total,
			"amount": referral_amount,
		},
		"firstbet": map[string]interface{}{
			"total":  firstbet_total,
			"amount": firstbet_amount,
		},
	}

	resp := JsonResponse{1, "", data}
	c.JSON(http.StatusOK, resp)
}

// Subscribe : collect user email
func (u UserController) Subscribe(c *gin.Context) {
	var userModel models.SubscribedUser
	email := c.DefaultPostForm("email", "_")
	product := c.DefaultPostForm("product", "_")
	productType := c.DefaultPostForm("type", "_")

	err := utils.ValidateFormat(email)
	if err != nil {
		resp := JsonResponse{0, "Invalid email.", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	if email != "_" {
		userModel.Email = email
	}
	if product != "_" {
		userModel.Product = product
	}

	if productType != "_" {
		userModel.ProductType = productType
	}

	db := models.Database()
	dbErr := db.Save(&userModel).Error

	if dbErr != nil {
		log.Println("Error", dbErr.Error())
		resp := JsonResponse{0, "Subscribe failed.", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	log.Println(userModel)
	resp := JsonResponse{1, "", userModel}

	switch product {
	case "cash":
		go mailService.SendCashEmail(email)
	case "prediction":
		go mailService.SendPredictionEmail(email)
	case "wallet":
		go mailService.SendWalletEmail(email)
	case "whisper":
		go mailService.SendWhisperEmail(email)
	case "chrome_extension":
		go mailService.SendChromeExtensionEmail(email)
	}

	c.JSON(http.StatusOK, resp)
}

// CountSubscribedUsers : count how many user subscribed a product
func (u UserController) CountSubscribedUsers(c *gin.Context) {
	product := c.DefaultPostForm("product", "")

	if !utils.ValidateProduct(product) {
		resp := JsonResponse{0, "Invalid product.", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	var userDAO = &daos.SubscribedUserDAO{}
	count, err := userDAO.CountUsersByProduct(product)
	if err != nil {
		resp := JsonResponse{0, err.Error(), nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	resp := JsonResponse{1, "", count}
	c.JSON(http.StatusOK, resp)
}

func (u UserController) Notification(c *gin.Context) {
	var data map[string]interface{}

	err := c.BindJSON(&data)

	if err != nil {
		fmt.Println(err)
		resp := JsonResponse{0, "Invalid params", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	to, hasTo := data["to"]

	fmt.Println(data)

	if !hasTo {
		resp := JsonResponse{0, "Invalid params", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	user := models.User{}
	errDb := models.Database().Where("wallet_addresses LIKE ?", fmt.Sprintf("%%%s%%", to)).First(&user).Error

	if errDb != nil {
		resp := JsonResponse{0, "User is not found.", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	if user.FCMToken == "" {
		resp := JsonResponse{0, "Invalid fcm token", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	data["to"] = user.FCMToken

	jsonData := map[string]interface{}{
		"data": data,
	}

	status, err := fcmService.Notify(jsonData)

	if !status {
		log.Println(err.Error())
		resp := JsonResponse{0, "Send notification failed.", nil}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	resp := JsonResponse{1, "", nil}
	c.JSON(http.StatusOK, resp)
}

func ExchangeSignUp(userId uint, refId uint) {
	jsonData := make(map[string]interface{})
	jsonData["id"] = userId
	jsonData["refId"] = refId

	endpoint, found := utils.GetForwardingEndpoint("exchange")
	log.Println(endpoint, found)
	jsonValue, _ := json.Marshal(jsonData)

	endpoint = fmt.Sprintf("%s/%s", endpoint, "user/profile")

	request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(request)
	if err != nil {
		log.Println("call exchange failed ", err)
	} else {
		log.Println("call exchange on SignUp success")
	}
}

func ExchangeUpdateProfile(userId uint, username string) {
	jsonData := make(map[string]interface{})
	jsonData["id"] = userId

	endpoint, found := utils.GetForwardingEndpoint("exchange")
	log.Println(endpoint, found)
	jsonValue, _ := json.Marshal(jsonData)

	endpoint = fmt.Sprintf("%s/%s?alias=%s", endpoint, "user/profile", username)

	request, _ := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(request)
	if err != nil {
		log.Println("call exchange failed ", err)
	} else {
		log.Println("call exchange on SignUp success")
	}
}

// FreeTokenReferrer :
func FreeTokenReferrer(userId string, refId string, network string) {
	log.Println("start free token referrer", userId, refId, network)
	ref := models.User{}
	errDb := models.Database().Where("id = ?", refId).First(&ref).Error

	if errDb != nil {
		log.Println("Get referrer failed.")
	} else {
		var refMd map[string]interface{}
		if ref.Metadata != "" {
			json.Unmarshal([]byte(ref.Metadata), &refMd)
		} else {
			refMd = map[string]interface{}{}
		}

		referrals, hasReferrals := refMd["referrals"]
		if !hasReferrals {
			referrals = map[string]interface{}{}
		}

		aReferrals := referrals.(map[string]interface{})

		bonusKey := fmt.Sprintf("bonus%s", userId)

		_, hasBonus := aReferrals[bonusKey]
		if !hasBonus {
			var refWallets map[string]interface{}
			if ref.RewardWalletAddresses != "" {
				json.Unmarshal([]byte(ref.RewardWalletAddresses), &refWallets)

				ethWallet, hasEthWallet := refWallets["ETH"]

				if hasEthWallet {
					amount := "20"
					address := ((ethWallet.(map[string]interface{}))["address"]).(string)

					time.Sleep(2 * time.Second)

					status, hash := ethereumService.FreeToken(fmt.Sprint(ref.ID), address, amount, network)
					log.Println("status", status, hash)
					if status {
						aReferrals[bonusKey] = map[string]interface{}{
							"address": address,
							"amount":  amount,
							"hash":    hash,
							"time":    time.Now().UTC().Unix(),
						}

						refMd["referrals"] = aReferrals
						metadata, _ := json.Marshal(refMd)
						ref.Metadata = string(metadata)
						dbErr := models.Database().Save(&ref).Error
						if dbErr != nil {
							log.Println(dbErr.Error())
						}
						log.Println(ref)

						go mailService.SendCompleteReferrer(ref.Email, ref.Username, hash)
					}
				}
			}
		}
	}
}

// CheckEmailExist : check email exist in system or not
func (u UserController) CheckEmailExist(c *gin.Context) {

	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)

	if userModel.Email == "" {
		resp := JsonResponse{0, "", map[string]interface{}{"email_existed": 0}}
		c.JSON(http.StatusOK, resp)
		c.Abort()
		return
	}

	resp := JsonResponse{1, "", map[string]interface{}{"email_existed": 1}}
	c.JSON(http.StatusOK, resp)
	c.Abort()
}
