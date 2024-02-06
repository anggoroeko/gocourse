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

func GetInventory(c *gin.Context) {
	reqInventory := []model.Inventory{}

	if data := config.DB.Preload(clause.Associations).Find(&reqInventory); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	getInventoryResponse := []model.InventoryResponse{}

	for _, p := range reqInventory {
		//:: GET DATA INVENTORY
		inventory, distributor, product, err := GetDataInventory(p.ID)
		if err != nil {
			stringSlice := []string{}
			message := err.Error()

			helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
			return
		}

		productResp := model.InventoryResponse{
			ID:                  p.ID,
			Name:                p.Name,
			DistributorID:       p.DistributorID,
			ProductID:           p.ProductID,
			DistributorResponse: distributor,
			ProductResponse:     product,
			Stock:               p.Stock,
			Price:               p.Price,
			CreatedAt:           inventory.CreatedAt,
			UpdatedAt:           inventory.UpdatedAt,
		}

		getInventoryResponse = append(getInventoryResponse, productResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    getInventoryResponse,
	})
}

func InsertInventory(c *gin.Context) {
	validate := validator.New()
	reqInventoryParam := model.Inventory{}

	if err := c.BindJSON(&reqInventoryParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqInventoryParam); errs != nil {
		// validateResp := []model.ErrValidationResp{}
		// for _, err := range errs.(validator.ValidationErrors) {
		// Access the field name causing the validation error
		// fieldName := err.Field()
		// valDateResp := model.ErrValidationResp{
		// 	Name: fieldName,
		// }
		// validateResp = append(validateResp, valDateResp)
		// }
		// stringSlice := errs.(validator.ValidationErrors)
		// helper.JsonResponseMap(stringSlice, message, http.StatusInternalServerError, c)
		stringSlice := []string{}
		message := errs.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	reqData := model.Inventory{
		Name:          reqInventoryParam.Name,
		DistributorID: reqInventoryParam.DistributorID,
		ProductID:     reqInventoryParam.ProductID,
		Stock:         reqInventoryParam.Stock,
		Price:         reqInventoryParam.Price,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"name", reqInventoryParam.Name}
	message := "Successfully Insert Inventory"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetInventoryById(c *gin.Context) {
	reqInventory := model.Inventory{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqInventory, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHANGE STRING TO UINT64
	uint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: IDENTIFICATION UINT64 AS UINT
	idN := uint(uint64)

	//:: GET DATA INVENTORY
	inventory, distributor, product, err := GetDataInventory(idN)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.InventoryResponse{
		ID:                  inventory.ID,
		Name:                reqInventory.Name,
		DistributorID:       reqInventory.DistributorID,
		ProductID:           reqInventory.ProductID,
		DistributorResponse: distributor,
		ProductResponse:     product,
		Stock:               reqInventory.Stock,
		Price:               reqInventory.Price,
		CreatedAt:           inventory.CreatedAt,
		UpdatedAt:           inventory.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully get data inventory",
		"status":  http.StatusOK,
	})
}

func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	reqInventory := model.Inventory{}

	//:: DELETE MASTER INVENTORY
	if data := config.DB.Delete(&reqInventory, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "Successfully update inventory",
		"status":  http.StatusOK,
	})
}

func UpdateInventory(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	//:: IDENTIFYING MODEL INVENTORY
	reqInventory := model.Inventory{}

	//:: GET BODY REQUEST JSON
	if err := c.BindJSON(&reqInventory); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqInventory); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: UPDATE INVENTORY
	dataUpdate := config.DB.Model(reqInventory).Where("id = ?", id).Updates(reqInventory)

	if dataUpdate.Error != nil {
		stringSlice := []string{}
		message := dataUpdate.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHANGE STRING TO UINT64
	uint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: IDENTIFICATION UINT64 AS UINT
	idN := uint(uint64)

	//:: GET DATA INVENTORY
	inventory, distributor, product, err := GetDataInventory(idN)

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.InventoryResponse{
		ID:                  inventory.ID,
		Name:                reqInventory.Name,
		DistributorID:       reqInventory.DistributorID,
		ProductID:           reqInventory.ProductID,
		DistributorResponse: distributor,
		ProductResponse:     product,
		Stock:               reqInventory.Stock,
		Price:               reqInventory.Price,
		CreatedAt:           inventory.CreatedAt,
		UpdatedAt:           inventory.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update inventory",
		"status":  http.StatusOK,
	})
}

func GetDataInventory(id uint) (model.InventoryResponse, model.DistributorResponse, model.ProductResponse, error) {
	reqInventory := model.Inventory{}
	dataInventory := config.DB.Preload(clause.Associations).Where("id=?", id).First(&reqInventory)

	//:: GET DATA DISTRIBUTOR FROM INVENTORY
	if dataInventory.Error != nil {
		return model.InventoryResponse{}, model.DistributorResponse{}, model.ProductResponse{}, dataInventory.Error
	}

	inventory := model.InventoryResponse{
		ID:        reqInventory.ID,
		CreatedAt: reqInventory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqInventory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	distributor := model.DistributorResponse{
		ID:        reqInventory.Distributor.ID,
		Name:      reqInventory.Distributor.Name,
		CreatedAt: reqInventory.Distributor.CreatedAt,
		UpdatedAt: reqInventory.Distributor.UpdatedAt,
	}

	//:: GET DATA PRODUCT FROM INVENTORY
	product := model.ProductResponse{
		ID:        reqInventory.Product.ID,
		Name:      reqInventory.Product.Name,
		CreatedAt: reqInventory.Product.CreatedAt,
		UpdatedAt: reqInventory.Product.UpdatedAt,
	}

	return inventory, distributor, product, nil
}
