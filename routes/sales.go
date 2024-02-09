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

func GetSales(c *gin.Context) {
	reqSales := []model.Sales{}

	if data := config.DB.Preload(clause.Associations).Find(&reqSales); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := []model.SalesResponse{}

	for _, s := range reqSales {
		//:: GET DATA DISTRIBUTOR AND DATA PRODUCT
		distributor, product, err := GetDataDistributorProduct(s.Inventory.DistributorID, s.Inventory.ProductID)
		if err != nil {
			stringSlice := []string{}
			message := err.Error()

			helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
			return
		}

		salesResponse := model.SalesResponse{
			ID:          s.ID,
			Member:      s.Member.Name,
			Chasier:     s.Chasier.Name,
			Inventory:   s.Inventory.Name,
			Distributor: distributor,
			Product:     product,
			Quantity:    s.Quantity,
			Price:       uint(s.Inventory.Price),
		}

		response = append(response, salesResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully get data sales",
		"status":  http.StatusOK,
	})
}

func InsertSales(c *gin.Context) {
	validate := validator.New()
	reqSalesParam := model.Sales{}

	if err := c.BindJSON(&reqSalesParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqSalesParam); errs != nil {
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

	//:: IDENTIFICATION DATA SALES
	reqData := model.Sales{
		InventoryID: reqSalesParam.InventoryID,
		MemberID:    reqSalesParam.MemberID,
		Quantity:    reqSalesParam.Quantity,
		ChasierID:   reqSalesParam.ChasierID,
	}

	//:: INSERT DATA SALES
	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: UPDATE DATA STOCK INVENTORY
	stockCount, err := UpdDataStockInventory(reqSalesParam.InventoryID, reqSalesParam.Quantity, true)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"stock_now", strconv.FormatInt(int64(stockCount), 10)}
	message := "Successfully Insert Sales and Update Stock in Inventory"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetSalesById(c *gin.Context) {
	reqSales := model.Sales{}

	id := c.Param("id")

	if data := config.DB.Preload(clause.Associations).First(&reqSales, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: GET DATA DISTRIBUTOR AND DATA PRODUCT
	distributor, product, err := GetDataDistributorProduct(reqSales.Inventory.DistributorID, reqSales.Inventory.ProductID)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.SalesResponse{
		ID:          reqSales.ID,
		Member:      reqSales.Member.Name,
		Chasier:     reqSales.Chasier.Name,
		Inventory:   reqSales.Inventory.Name,
		Distributor: distributor,
		Product:     product,
		Quantity:    reqSales.Quantity,
		Price:       uint(reqSales.Inventory.Price),
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully get data sales",
		"status":  http.StatusOK,
	})
}

func DeleteSales(c *gin.Context) {
	id := c.Param("id")

	reqSales := model.Sales{}

	//:: DELETE MASTER INVENTORY
	if data := config.DB.Delete(&reqSales, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "Successfully delete sales",
		"status":  http.StatusOK,
	})
}

func UpdateSales(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	//:: IDENTIFYING MODEL INVENTORY
	reqSales := model.Sales{}

	//:: GET BODY REQUEST JSON
	if err := c.BindJSON(&reqSales); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED USER
	if errs := validate.Struct(&reqSales); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: UPDATE SALES
	dataUpdate := config.DB.Model(reqSales).Where("id = ?", id).Updates(reqSales)

	if dataUpdate.Error != nil {
		stringSlice := []string{}
		message := dataUpdate.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: GET DATA DISTRIBUTOR AND DATA PRODUCT
	if data := config.DB.Preload(clause.Associations).First(&reqSales, "id = ?", id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	distributor, product, err := GetDataDistributorProduct(reqSales.Inventory.DistributorID, reqSales.Inventory.ProductID)
	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.SalesResponse{
		ID:          reqSales.ID,
		Member:      reqSales.Member.Name,
		Chasier:     reqSales.Chasier.Name,
		Inventory:   reqSales.Inventory.Name,
		Distributor: distributor,
		Product:     product,
		Quantity:    reqSales.Quantity,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Successfully update inventory",
		"status":  http.StatusOK,
	})
}

func GetDataDistributorProduct(distributorId, productId uint) (string, string, error) {
	reqDistributor := model.Distributor{}
	if data := config.DB.Preload(clause.Associations).First(&reqDistributor, "id = ?", distributorId); data.Error != nil {
		return "", "", data.Error
	}

	//:: GET DATA PRODUCT
	reqProduct := model.Product{}
	if data := config.DB.Preload(clause.Associations).First(&reqProduct, "id = ?", productId); data.Error != nil {
		return "", "", data.Error
	}

	return reqDistributor.Name, reqProduct.Name, nil
}

func UpdDataStockInventory(inventoryId, quantityNow uint, isSales bool) (uint, error) {
	reqDataInv := model.Inventory{}
	if dataInv := config.DB.First(&reqDataInv, "id = ?", inventoryId); dataInv.Error != nil {
		return 0, dataInv.Error
	}

	//:: UPDATE DATA STOCK INVENTORY
	var stockCount uint
	if isSales { //:: IF SALES STOCK DECREASES
		stockCount = reqDataInv.Stock - quantityNow
	} else { //:: IF PURCHASES STOCK INCREASES
		stockCount = reqDataInv.Stock + quantityNow
	}

	reqUpdInv := model.Inventory{
		Stock: stockCount,
	}

	if updInv := config.DB.Where("id = ?", inventoryId).Updates(reqUpdInv); updInv.Error != nil {
		return 0, updInv.Error
	}

	return stockCount, nil
}
