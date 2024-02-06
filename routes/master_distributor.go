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

func GetDistributor(c *gin.Context) {
	reqDistributor := []model.Distributor{}

	if data := config.DB.Preload(clause.Associations).Find(&reqDistributor); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	getDistributorResponse := []model.DistributorResponse{}

	for _, p := range reqDistributor {
		product := model.DistributorResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}

		getDistributorResponse = append(getDistributorResponse, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    getDistributorResponse,
	})
}

func InsertDistributor(c *gin.Context) {
	reqDistributorParam := model.Distributor{}
	if err := c.BindJSON(&reqDistributorParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	reqData := model.Distributor{
		Name: reqDistributorParam.Name,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"name", reqDistributorParam.Name}
	message := "Successfully Insert User"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetDistributorById(c *gin.Context) {
	reqDistributor := model.Distributor{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqDistributor, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	response := model.DistributorResponse{
		ID:        reqDistributor.ID,
		Name:      reqDistributor.Name,
		CreatedAt: reqDistributor.CreatedAt,
		UpdatedAt: reqDistributor.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    response,
	})
}

func DeleteDistributor(c *gin.Context) {
	id := c.Param("id")

	reqDistributor := model.Distributor{}
	if data := config.DB.Delete(&reqDistributor, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{}
	message := "Successfully deleted product"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func UpdateDistributor(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	reqDistributor := model.Distributor{}
	// currentTime := time.Now()

	//: CONVERT STRING TO UINT
	uint64Value, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	if err := c.BindJSON(&reqDistributor); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqDistributor); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: IDENTIFICATION UINT64 AS UINT
	idVal := uint(uint64Value)

	//:: REQUEST USER DATA
	reqData := model.Distributor{
		Name: reqDistributor.Name,
	}

	//:: UPDATE DATA
	data := config.DB.Model(&reqDistributor).Where("id = ?", id).Updates(reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: ROLE RESPONSE
	dataRole := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqDistributor)

	if dataRole.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.DistributorResponse{
		ID:        idVal,
		Name:      reqDistributor.Name,
		CreatedAt: reqDistributor.CreatedAt,
		UpdatedAt: reqDistributor.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update product",
		"status":  http.StatusOK,
	})
}
