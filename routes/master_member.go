package routes

import (
	"go_pos_v1_2/config"
	"go_pos_v1_2/helper"
	model "go_pos_v1_2/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm/clause"
)

func GetMember(c *gin.Context) {
	reqMember := []model.Member{}

	if data := config.DB.Preload(clause.Associations).Find(&reqMember); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	getMemberResponse := []model.MemberResponse{}

	for _, p := range reqMember {
		product := model.MemberResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}

		getMemberResponse = append(getMemberResponse, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    getMemberResponse,
	})
}

func InsertMember(c *gin.Context) {
	reqMemberParam := model.Member{}
	if err := c.BindJSON(&reqMemberParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	reqData := model.Member{
		Name: reqMemberParam.Name,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"name", reqMemberParam.Name}
	message := "Successfully Insert User"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetMemberById(c *gin.Context) {
	reqMember := model.Member{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqMember, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	response := model.MemberResponse{
		ID:        reqMember.ID,
		Name:      reqMember.Name,
		CreatedAt: reqMember.CreatedAt,
		UpdatedAt: reqMember.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    response,
	})
}

func DeleteMember(c *gin.Context) {
	id := c.Param("id")

	reqMember := model.Member{}
	if data := config.DB.Delete(&reqMember, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{}
	message := "Successfully deleted member"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func UpdateMember(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	reqMember := model.Member{}
	// currentTime := time.Now()

	//: CONVERT STRING TO UINT
	uint64Value, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	if err := c.BindJSON(&reqMember); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqMember); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: IDENTIFICATION UINT AS UINT64
	idVal := uint(uint64Value)

	//:: REQUEST USER DATA
	reqData := model.Member{
		Name: reqMember.Name,
	}

	//:: UPDATE DATA
	data := config.DB.Model(&reqMember).Where("id = ?", id).Updates(reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: ROLE RESPONSE
	dataRole := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqMember)

	if dataRole.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.MemberResponse{
		ID:        idVal,
		Name:      reqMember.Name,
		CreatedAt: reqMember.CreatedAt,
		UpdatedAt: reqMember.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update member",
		"status":  http.StatusOK,
	})
}
