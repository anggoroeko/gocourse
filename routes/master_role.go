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

func InsertRole(c *gin.Context) {
	reqRole := model.Role{}
	if err := c.BindJSON(&reqRole); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	if data := config.DB.Create(&reqRole); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{"name", reqRole.Name}
	message := "Successfully insert role"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetRole(c *gin.Context) {
	reqRole := []model.Role{}

	if data := config.DB.Find(&reqRole); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	roleResponse := []model.RoleResponse{}

	for _, role := range reqRole {
		rR := model.RoleResponse{
			ID:        role.ID,
			Name:      role.Name,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		roleResponse = append(roleResponse, rR)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully",
		"data":    roleResponse,
	})
}

func GetRoleById(c *gin.Context) {
	reqRole := model.Role{}

	id := c.Param("id")

	if data := config.DB.First(&reqRole, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusNotFound, c)
		return
	}

	response := model.RoleResponse{
		ID:        reqRole.ID,
		Name:      reqRole.Name,
		CreatedAt: reqRole.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqRole.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully",
		"data":    response,
	})
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	reqRole := model.Role{}
	if data := config.DB.Delete(&reqRole, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{}
	message := "Successfullyfully deleted role"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func UpdateRole(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	reqRole := model.Role{}
	currentTime := time.Now()

	//: CONVERT STRING TO UINT
	uint64Value, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	if err := c.BindJSON(&reqRole); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqRole); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: REQUEST USER DATA
	reqData := model.Role{
		Name: reqRole.Name,
	}

	//:: UPDATE DATA
	data := config.DB.Model(&reqRole).Where("id = ?", id).Updates(reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: ROLE RESPONSE
	dataRole := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqRole)

	if dataRole.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: IDENTIFICATION UINT AS UINT64
	idVal := uint(uint64Value)

	//:: FINAL RESPONSE
	response := model.RoleResponse{
		ID:        idVal,
		Name:      reqRole.Name,
		CreatedAt: reqRole.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: currentTime.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update role",
		"status":  http.StatusOK,
	})
}
