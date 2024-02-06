package routes

import (
	"go_pos_v1_2/config"
	"go_pos_v1_2/helper"
	model "go_pos_v1_2/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm/clause"
)

func GetUser(c *gin.Context) {
	reqUser := []model.User{}

	if data := config.DB.Preload(clause.Associations).Find(&reqUser); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	response := []model.UserResponse{}

	for _, u := range reqUser {
		role := model.RoleResponse{
			ID:        u.Role.ID,
			Name:      u.Role.Name,
			CreatedAt: u.Role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.Role.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		user := model.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			RoleID:    u.RoleID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			Role:      role,
		}

		response = append(response, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    response,
	})
}

func RegisterUser(c *gin.Context) {
	validate := validator.New()
	reqUserParam := model.User{}
	if err := c.BindJSON(&reqUserParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqUserParam); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: CHANGE PASSWORD TO HASH
	password := reqUserParam.Password
	hash, err := model.HashPassword(password) // ignore error for the sake of simplicity

	if err != nil {
		stringSlice := []string{}
		message := "Failed insert user, hash password not generated"

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: INSERT DATA USER
	reqData := model.User{
		Name:     reqUserParam.Name,
		RoleID:   reqUserParam.RoleID,
		Username: reqUserParam.Username,
		Password: hash,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"username", reqUserParam.Username}
	message := "Success Insert User"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetUserById(c *gin.Context) {
	reqUser := model.User{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqUser, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	role := model.RoleResponse{
		ID:        reqUser.Role.ID,
		Name:      reqUser.Role.Name,
		CreatedAt: reqUser.Role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqUser.Role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	response := model.UserResponse{
		ID:        reqUser.ID,
		Name:      reqUser.Name,
		RoleID:    reqUser.RoleID,
		Username:  reqUser.Username,
		CreatedAt: reqUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		Role:      role,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    response,
	})
}

func GenerateToken(c *gin.Context) {
	request := model.TokenRequest{}
	user := model.User{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"error":   err.Error(),
		})

		c.Abort()
		return
	}

	//:: CHECK USERNAME
	checkEmail := config.DB.Where("username = ?", request.Username).First(&user)

	if checkEmail.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Username not found",
			"error":   checkEmail.Error.Error(),
		})

		c.Abort()
		return
	}

	//:: CHECK PASSWORD
	credentialError := user.CheckPasswordHash(request.Password)

	if credentialError != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Password Not match",
			"error":   credentialError.Error(),
		})

		c.Abort()
		return
	}

	//:: GENERATE TOKEN
	tokenString, err := helper.GenerateJWT(user.Username, user.Username, user.RoleID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})

		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"token":   tokenString,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	reqUser := model.User{}
	if data := config.DB.Delete(&reqUser, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{}
	message := "Success deleted user"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func UpdateUser(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	reqUser := model.User{}
	currentTime := time.Now()

	//: CONVERT STRING TO UINT
	uint64Value, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: GET BODY JSON
	if err := c.BindJSON(&reqUser); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqUser); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: IDENTIFICATION UINT AS UINT64
	idVal := uint(uint64Value)

	//:: SET HASH PASSWORD
	password := reqUser.Password
	hash, err := model.HashPassword(password) // ignore error for the sake of simplicity

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: REQUEST USER DATA
	reqData := model.User{
		Name:     reqUser.Name,
		RoleID:   reqUser.RoleID,
		Username: reqUser.Username,
		Password: hash,
	}

	//:: UPDATE DATA
	data := config.DB.Model(&reqUser).Where("id = ?", id).Updates(reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: ROLE RESPONSE
	dataRole := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqUser)

	if dataRole.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	role := model.RoleResponse{
		ID:        reqUser.Role.ID,
		Name:      reqUser.Role.Name,
		CreatedAt: reqUser.Role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqUser.Role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	//:: FINAL RESPONSE
	response := model.UserResponse{
		ID:        idVal,
		Name:      reqUser.Name,
		RoleID:    reqUser.RoleID,
		Username:  reqUser.Username,
		CreatedAt: reqUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: currentTime.Format("2006-01-02 15:04:05"),
		Role:      role,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Success update user",
		"status":  http.StatusOK,
	})
}
