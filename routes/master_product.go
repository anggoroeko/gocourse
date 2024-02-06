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

func GetProduct(c *gin.Context) {
	reqProduct := []model.Product{}

	if data := config.DB.Preload(clause.Associations).Find(&reqProduct); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	getProductResponse := []model.ProductResponse{}

	for _, p := range reqProduct {
		product := model.ProductResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}

		getProductResponse = append(getProductResponse, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    getProductResponse,
	})
}

func InsertProduct(c *gin.Context) {
	reqProductParam := model.Product{}
	if err := c.BindJSON(&reqProductParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	reqData := model.Product{
		Name: reqProductParam.Name,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"name", reqProductParam.Name}
	message := "Successfully Insert User"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetProductById(c *gin.Context) {
	reqProduct := model.Product{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqProduct, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	response := model.ProductResponse{
		ID:        reqProduct.ID,
		Name:      reqProduct.Name,
		CreatedAt: reqProduct.CreatedAt,
		UpdatedAt: reqProduct.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    response,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	reqProduct := model.Product{}
	if data := config.DB.Delete(&reqProduct, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{}
	message := "Successfully deleted product"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func UpdateProduct(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	reqProduct := model.Product{}
	currentTime := time.Now()

	//: CONVERT STRING TO UINT
	uint64Value, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	if err := c.BindJSON(&reqProduct); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqProduct); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: IDENTIFICATION UINT AS UINT64
	idVal := uint(uint64Value)

	//:: REQUEST USER DATA
	reqData := model.Product{
		Name: reqProduct.Name,
	}

	//:: UPDATE DATA
	data := config.DB.Model(&reqProduct).Where("id = ?", id).Updates(reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: ROLE RESPONSE
	dataRole := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqProduct)

	if dataRole.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.ProductResponse{
		ID:        idVal,
		Name:      reqProduct.Name,
		CreatedAt: reqProduct.CreatedAt,
		UpdatedAt: currentTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update product",
		"status":  http.StatusOK,
	})
}
