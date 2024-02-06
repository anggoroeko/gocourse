package routes

import (
	"go_pos_v1_2/config"
	"go_pos_v1_2/helper"
	model "go_pos_v1_2/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm/clause"
)

func GetChasier(c *gin.Context) {
	reqChasier := []model.Chasier{}

	data := config.DB.Preload(clause.Associations).Find(&reqChasier)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	response := []model.ChasierResponse{}

	for _, u := range reqChasier {
		//:: GET RESPONSE OF USER AND ROLE
		user, err := GetUserRoleMap(int(u.User.RoleID))

		if err != nil {
			stringSlice := []string{}
			message := err.Error()

			helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
			return
		}

		Chasier := model.ChasierResponse{
			ID:        u.ID,
			Name:      u.Name,
			UserID:    u.UserID,
			User:      user,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		response = append(response, Chasier)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    response,
	})
}

func InsertChasier(c *gin.Context) {
	validate := validator.New()
	reqChasierParam := model.Chasier{}

	//:: GET BODY JSON
	if err := c.BindJSON(&reqChasierParam); err != nil {
		stringSlice := []string{}
		message := err.Error()
		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED Chasier
	if errs := validate.Struct(&reqChasierParam); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	reqData := model.Chasier{
		Name:   reqChasierParam.Name,
		UserID: reqChasierParam.UserID,
	}

	data := config.DB.Create(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()
		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: PUT ALL RESPONSE TO HELPER : ON DEVELOPMENT
	stringSlice := []string{"Chasiername", reqChasierParam.Name}
	message := "Success Insert Chasier"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func GetChasierById(c *gin.Context) {
	reqChasier := model.Chasier{}

	id := c.Param("id")

	//:: USER AND ROLE RESPONSE
	dataUser := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqChasier)

	if dataUser.Error != nil {
		stringSlice := []string{}
		message := dataUser.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: GET RESPONSE OF USER AND ROLE
	// user, err := GetUserRole(reqChasier)
	user, err := GetUserRoleMap(int(reqChasier.User.RoleID))

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.ChasierResponse{
		ID:        reqChasier.ID,
		Name:      reqChasier.Name,
		UserID:    reqChasier.UserID,
		CreatedAt: reqChasier.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqChasier.UpdatedAt.Format("2006-01-02 15:04:05"),
		User:      user,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    response,
	})
}

func DeleteChasier(c *gin.Context) {
	id := c.Param("id")

	reqChasier := model.Chasier{}
	if data := config.DB.Delete(&reqChasier, id); data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	stringSlice := []string{}
	message := "Success deleted Chasier"

	helper.JsonResponse(stringSlice, message, http.StatusOK, c)
}

func UpdateChasier(c *gin.Context) {
	validate := validator.New()
	id := c.Param("id")

	reqChasier := model.Chasier{}
	// currentTime := time.Now()

	if err := c.BindJSON(&reqChasier); err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: CHECK REQUIRED/VALIDATED Chasier
	if errs := validate.Struct(&reqChasier); errs != nil {
		stringSlice := []string{}
		message := errs.Error()

		helper.JsonResponse(stringSlice, message, http.StatusBadRequest, c)
		c.Abort()
		return
	}

	//:: REQUEST Chasier DATA
	reqData := model.Chasier{
		Name:   reqChasier.Name,
		UserID: reqChasier.UserID,
	}

	//:: UPDATE DATA
	data := config.DB.Model(&reqChasier).Where("id = ?", id).Updates(&reqData)

	if data.Error != nil {
		stringSlice := []string{}
		message := data.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: USER AND ROLE RESPONSE
	dataUser := config.DB.Preload(clause.Associations).Where("id = ?", id).First(&reqChasier)

	if dataUser.Error != nil {
		stringSlice := []string{}
		message := dataUser.Error.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: GET RESPONSE OF USER AND ROLE
	// user, err := GetUserRole(reqChasier)
	user, err := GetUserRoleMap(int(reqChasier.User.RoleID))

	if err != nil {
		stringSlice := []string{}
		message := err.Error()

		helper.JsonResponse(stringSlice, message, http.StatusInternalServerError, c)
		return
	}

	//:: FINAL RESPONSE
	response := model.ChasierResponse{
		ID:        reqChasier.ID,
		Name:      reqChasier.Name,
		UserID:    reqChasier.UserID,
		CreatedAt: reqChasier.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqChasier.UpdatedAt.Format("2006-01-02 15:04:05"),
		User:      user,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Success update Chasier",
		"status":  http.StatusOK,
	})
}

func GetUserRole(reqChasier model.Chasier) (model.UserResponse, error) {
	reqUserRole := model.User{}
	dataUserRole := config.DB.Preload(clause.Associations).Where("role_id = ?", reqChasier.User.RoleID).First(&reqUserRole)

	if dataUserRole.Error != nil {
		return model.UserResponse{}, dataUserRole.Error
	}

	//:: ROLE RESPONSE
	Role := model.RoleResponse{
		ID:        reqUserRole.Role.ID,
		Name:      reqUserRole.Role.Name,
		CreatedAt: reqUserRole.Role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqUserRole.Role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	//:: USER RESPONSE
	User := model.UserResponse{
		ID:        reqChasier.User.ID,
		Name:      reqChasier.User.Name,
		RoleID:    reqChasier.User.RoleID,
		Username:  reqChasier.User.Username,
		CreatedAt: reqChasier.User.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqChasier.User.UpdatedAt.Format("2006-01-02 15:04:05"),
		Role:      Role,
	}

	return User, nil
}

func GetUserRoleMap(roleId int) (model.UserResponse, error) {
	reqUserRole := model.User{}
	dataUserRole := config.DB.Preload(clause.Associations).Where("role_id = ?", roleId).First(&reqUserRole)

	if dataUserRole.Error != nil {
		return model.UserResponse{}, dataUserRole.Error
	}

	//:: ROLE RESPONSE
	Role := model.RoleResponse{
		ID:        reqUserRole.Role.ID,
		Name:      reqUserRole.Role.Name,
		CreatedAt: reqUserRole.Role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqUserRole.Role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	//:: USER RESPONSE
	User := model.UserResponse{
		ID:        reqUserRole.ID,
		Name:      reqUserRole.Name,
		RoleID:    reqUserRole.RoleID,
		Username:  reqUserRole.Username,
		CreatedAt: reqUserRole.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reqUserRole.UpdatedAt.Format("2006-01-02 15:04:05"),
		Role:      Role,
	}

	return User, nil
}
