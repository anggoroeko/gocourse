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

func GetPurchase(c *gin.Context) {
	reqPurchase := []model.Purchase{}

	if data := config.DB.Preload(clause.Associations).Find(&reqPurchase); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	getPurchaseResponse := []model.PurchaseResponse{}

	for _, p := range reqPurchase {
		//:: GET DATA Purchase
		inventory, err := GetDataPurchaseInventory(p.InventoryID)
		if err != nil {
			stringSlice := []string{}
			message := err.Error()

			helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
			return
		}

		purchaseResp := model.PurchaseResponse{
			ID:          p.ID,
			InventoryID: p.InventoryID,
			Date:        p.Date,
			Quantity:    p.Quantity,
			Description: p.Description,
			Inventory:   inventory,
			CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		getPurchaseResponse = append(getPurchaseResponse, purchaseResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"data":    getPurchaseResponse,
	})
}

func InsertPurchase(c *gin.Context) {
	validate := validator.New()
	reqPurchaseParam := model.Purchase{}

	if err := c.BindJSON(&reqPurchaseParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqPurchaseParam); errs != nil {
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

	reqData := model.Purchase{
		InventoryID: reqPurchaseParam.InventoryID,
		Date:        reqPurchaseParam.Date,
		Quantity:    reqPurchaseParam.Quantity,
		Description: reqPurchaseParam.Description,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: UPDATE DATA STOCK INVENTORY
	stockCount, err := UpdDataStockInventory(reqPurchaseParam.InventoryID, reqPurchaseParam.Quantity, false)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"stock_now", strconv.FormatInt(int64(stockCount), 10)}
	message := "Successfully Insert Purchase"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetPurchaseById(c *gin.Context) {
	reqPurchase := model.Purchase{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqPurchase, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHANGE STRING TO UINT64
	idInventory := reqPurchase.InventoryID

	//:: GET DATA INVENTORY
	inventory, err := GetDataPurchaseInventory(idInventory)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.PurchaseResponse{
		ID:          reqPurchase.ID,
		Date:        reqPurchase.Date,
		InventoryID: reqPurchase.InventoryID,
		Inventory:   inventory,
		Quantity:    reqPurchase.Quantity,
		Description: reqPurchase.Description,
		CreatedAt:   reqPurchase.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   reqPurchase.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully get data Purchase",
		"status":  http.StatusOK,
	})
}

func DeletePurchase(c *gin.Context) {
	id := c.Param("id")

	reqPurchase := model.Purchase{}

	//:: DELETE MASTER Purchase
	if data := config.DB.Delete(&reqPurchase, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "Successfully delete Purchase",
		"status":  http.StatusOK,
	})
}

func UpdatePurchase(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	//:: IDENTIFYING MODEL Purchase
	reqPurchase := model.Purchase{}

	//:: GET BODY REQUEST JSON
	if err := c.BindJSON(&reqPurchase); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqPurchase); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: UPDATE Purchase
	dataUpdate := config.DB.Model(reqPurchase).Where("id = ?", id).Updates(reqPurchase)

	if dataUpdate.Error != nil {
		stringSlice := []string{}
		message := dataUpdate.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHANGE STRING TO UINT64
	idInventory := reqPurchase.InventoryID

	//:: GET DATA INVENTORY
	inventory, err := GetDataPurchaseInventory(idInventory)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.PurchaseResponse{
		ID:        reqPurchase.ID,
		Date:      reqPurchase.Date,
		Inventory: inventory,
		Quantity:  reqPurchase.Quantity,
		CreatedAt: reqPurchase.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqPurchase.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update Purchase",
		"status":  http.StatusOK,
	})
}

func GetDataPurchaseInventory(idInventory uint) (model.InventoryResponse, error) {
	reqPurchaseInventory := model.Inventory{}
	dataPurchaseInventory := config.DB.Preload(clause.Associations).Where("id=?", idInventory).First(&reqPurchaseInventory)

	if dataPurchaseInventory.Error != nil {
		return model.InventoryResponse{}, dataPurchaseInventory.Error
	}

	//:: GET DATA DISTRIBUTOR FROM INVENTORY
	distributor := model.DistributorResponse{
		ID:        reqPurchaseInventory.Distributor.ID,
		Name:      reqPurchaseInventory.Distributor.Name,
		CreatedAt: reqPurchaseInventory.Distributor.CreatedAt,
		UpdatedAt: reqPurchaseInventory.Distributor.UpdatedAt,
	}

	//:: GET DATA PRODUCT FROM INVENTORY
	product := model.ProductResponse{
		ID:        reqPurchaseInventory.Product.ID,
		Name:      reqPurchaseInventory.Product.Name,
		CreatedAt: reqPurchaseInventory.Product.CreatedAt,
		UpdatedAt: reqPurchaseInventory.Product.UpdatedAt,
	}

	inventory := model.InventoryResponse{
		ID:                  reqPurchaseInventory.ID,
		ProductID:           reqPurchaseInventory.ProductID,
		Name:                reqPurchaseInventory.Name,
		DistributorID:       reqPurchaseInventory.DistributorID,
		Stock:               reqPurchaseInventory.Stock,
		Price:               reqPurchaseInventory.Price,
		ProductResponse:     product,
		DistributorResponse: distributor,
		CreatedAt:           reqPurchaseInventory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           reqPurchaseInventory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return inventory, nil
}
