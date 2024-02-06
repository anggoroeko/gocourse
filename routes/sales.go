package routes

import (
	"go_pos_v1_2/config"
	"go_pos_v1_2/helper"
	model "go_pos_v1_2/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetSales(c *gin.Context) {
	reqSales := []model.Sales{}

	if data := config.DB.Preload(clause.Associations).Find(&reqSales); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	// GetSalesResponse := []model.SalesResponse{}

	// for _, p := range reqSales {
	// 	//:: GET DATA INVENTORY
	// 	inventory, distributor, product, err := GetDataInventory(p.ID)
	// 	if err != nil {
	// 		stringSlice := []string{}
	// 		message := err.Error()

	// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
	// 		return
	// 	}

	// 	productResp := model.SalesResponse{
	// 		// ID:                  p.ID,
	// 		// Name:                p.Name,
	// 		// DistributorID:       p.DistributorID,
	// 		// ProductID:           p.ProductID,
	// 		// DistributorResponse: distributor,
	// 		// ProductResponse:     product,
	// 		// Stock:               p.Stock,
	// 		// Price:               p.Price,
	// 		// CreatedAt:           inventory.CreatedAt,
	// 		// UpdatedAt:           inventory.UpdatedAt,
	// 	}

	// 	GetSalesResponse = append(GetSalesResponse, productResp)
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Successfully",
	// 	"data":    GetSalesResponse,
	// })
}

// func InsertInventory(c *gin.Context) {
// 	validate := validator.New()
// 	reqSalesParam := model.Sales{}

// 	if err := c.BindJSON(&reqSalesParam); err != nil {
// 		stringSlice := []string{}
// 		message := err.Error()
// 		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
// 		return
// 	}

// 	//:: CHECK REQUIRED/VALIDATED USER
// 	if errs := validate.Struct(&reqSalesParam); errs != nil {
// 		// validateResp := []model.ErrValidationResp{}
// 		// for _, err := range errs.(validator.ValidationErrors) {
// 		// Access the field name causing the validation error
// 		// fieldName := err.Field()
// 		// valDateResp := model.ErrValidationResp{
// 		// 	Name: fieldName,
// 		// }
// 		// validateResp = append(validateResp, valDateResp)
// 		// }
// 		// stringSlice := errs.(validator.ValidationErrors)
// 		// helper.JsonResponseMap(stringSlice, message, http.StatusInternalServerError, c)
// 		stringSlice := []string{}
// 		message := errs.Error()
// 		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
// 		return
// 	}

// 	reqData := model.Sales{
// 		Name:          reqSalesParam.Name,
// 		DistributorID: reqSalesParam.DistributorID,
// 		ProductID:     reqSalesParam.ProductID,
// 		Stock:         reqSalesParam.Stock,
// 		Price:         reqSalesParam.Price,
// 	}

// 	data := config.DB.Create(&reqData)

// 	if data.Error != nil {
// 		stringSlice := []string{}
// 		message := data.Error.Error()
// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
// 	stringSlice := []string{"name", reqSalesParam.Name}
// 	message := "Successfully Insert Inventory"

// 	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
// }

// func GetSalesById(c *gin.Context) {
// 	reqSales := model.Sales{}

// 	id := c.Param("id")

// 	if data := config.DB.Preload(clause.Associations).First(&reqSales, "id = ?", id); data.Error != nil {
// 		stringSlice := []string{}
// 		message := data.Error.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: CHANGE STRING TO UINT64
// 	uint64, err := strconv.ParseUint(id, 10, 0)
// 	if err != nil {
// 		stringSlice := []string{}
// 		message := err.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: IDENTIFICATION UINT64 AS UINT
// 	idN := uint(uint64)

// 	//:: GET DATA INVENTORY
// 	inventory, distributor, product, err := GetDataInventory(idN)
// 	if err != nil {
// 		stringSlice := []string{}
// 		message := err.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: FINAL RESPONSE
// 	response := model.SalesResponse{
// 		ID:                  inventory.ID,
// 		Name:                reqSales.Name,
// 		DistributorID:       reqSales.DistributorID,
// 		ProductID:           reqSales.ProductID,
// 		DistributorResponse: distributor,
// 		ProductResponse:     product,
// 		Stock:               reqSales.Stock,
// 		Price:               reqSales.Price,
// 		CreatedAt:           inventory.CreatedAt,
// 		UpdatedAt:           inventory.UpdatedAt,
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":    response,
// 		"message": "Successfully get data inventory",
// 		"status":  http.StatusOK,
// 	})
// }

// func DeleteInventory(c *gin.Context) {
// 	id := c.Param("id")

// 	reqSales := model.Sales{}

// 	//:: DELETE MASTER INVENTORY
// 	if data := config.DB.Delete(&reqSales, id); data.Error != nil {
// 		stringSlice := []string{}
// 		message := data.Error.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":    nil,
// 		"message": "Successfully update inventory",
// 		"status":  http.StatusOK,
// 	})
// }

// func UpdateInventory(c *gin.Context) {
// 	validate := validator.New()
// 	id := c.Param("id")

// 	//:: IDENTIFYING MODEL INVENTORY
// 	reqSales := model.Sales{}

// 	//:: GET BODY REQUEST JSON
// 	if err := c.BindJSON(&reqSales); err != nil {
// 		stringSlice := []string{}
// 		message := err.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: CHECK REQUIRED/VALIDATED USER
// 	if errs := validate.Struct(&reqSales); errs != nil {
// 		stringSlice := []string{}
// 		message := errs.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
// 		c.Abort()
// 		return
// 	}

// 	//:: UPDATE INVENTORY
// 	dataUpdate := config.DB.Model(reqSales).Where("id = ?", id).Updates(reqSales)

// 	if dataUpdate.Error != nil {
// 		stringSlice := []string{}
// 		message := dataUpdate.Error.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: CHANGE STRING TO UINT64
// 	uint64, err := strconv.ParseUint(id, 10, 0)
// 	if err != nil {
// 		stringSlice := []string{}
// 		message := err.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: IDENTIFICATION UINT64 AS UINT
// 	idN := uint(uint64)

// 	//:: GET DATA INVENTORY
// 	inventory, distributor, product, err := GetDataInventory(idN)

// 	if err != nil {
// 		stringSlice := []string{}
// 		message := err.Error()

// 		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
// 		return
// 	}

// 	//:: FINAL RESPONSE
// 	response := model.SalesResponse{
// 		ID:                  inventory.ID,
// 		Name:                reqSales.Name,
// 		DistributorID:       reqSales.DistributorID,
// 		ProductID:           reqSales.ProductID,
// 		DistributorResponse: distributor,
// 		ProductResponse:     product,
// 		Stock:               reqSales.Stock,
// 		Price:               reqSales.Price,
// 		CreatedAt:           inventory.CreatedAt,
// 		UpdatedAt:           inventory.UpdatedAt,
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":    response,
// 		"message": "Successfully update inventory",
// 		"status":  http.StatusOK,
// 	})
// }

// func GetDataInventory(id uint) (model.SalesResponse, model.DistributorResponse, model.ProductResponse, error) {
// 	reqSales := model.Sales{}
// 	dataInventory := config.DB.Preload(clause.Associations).Where("id=?", id).First(&reqSales)

// 	//:: GET DATA DISTRIBUTOR FROM INVENTORY
// 	if dataInventory.Error != nil {
// 		return model.SalesResponse{}, model.DistributorResponse{}, model.ProductResponse{}, dataInventory.Error
// 	}

// 	inventory := model.SalesResponse{
// 		ID:        reqSales.ID,
// 		CreatedAt: reqSales.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UpdatedAt: reqSales.UpdatedAt.Format("2006-01-02 15:04:05"),
// 	}

// 	distributor := model.DistributorResponse{
// 		ID:        reqSales.Distributor.ID,
// 		Name:      reqSales.Distributor.Name,
// 		CreatedAt: reqSales.Distributor.CreatedAt,
// 		UpdatedAt: reqSales.Distributor.UpdatedAt,
// 	}

// 	//:: GET DATA PRODUCT FROM INVENTORY
// 	product := model.ProductResponse{
// 		ID:        reqSales.Product.ID,
// 		Name:      reqSales.Product.Name,
// 		CreatedAt: reqSales.Product.CreatedAt,
// 		UpdatedAt: reqSales.Product.UpdatedAt,
// 	}

// 	return inventory, distributor, product, nil
// }
