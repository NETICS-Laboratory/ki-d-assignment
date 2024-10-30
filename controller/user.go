package controller

import (
	"ki-d-assignment/common"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"

	// "ki-d-assignment/entity"
	"ki-d-assignment/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	// GetAllUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	// DeleteUser(ctx *gin.Context)
	// UpdateUser(ctx *gin.Context)
	MeUser(ctx *gin.Context)
	MeUserDecrypted(ctx *gin.Context)
	DecryptUserIDCard(ctx *gin.Context)
	RequestAccess(ctx *gin.Context)
	GetAccessRequests(ctx *gin.Context)
	UpdateAccessRequestStatus(ctx *gin.Context)
}

type userController struct {
	jwtService  service.JWTService
	userService service.UserService
}

func NewUserController(us service.UserService, jwts service.JWTService) UserController {
	return &userController{
		userService: us,
		jwtService:  jwts,
	}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	var user dto.UserCreateDto
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	checkUser, _ := uc.userService.CheckUser(ctx.Request.Context(), user.Username)
	if checkUser {
		res := common.BuildErrorResponse("User Sudah Terdaftar", "false", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menambahkan User", result)
	ctx.JSON(http.StatusOK, res)
}

// func (uc *userController) GetAllUser(ctx *gin.Context) {
// 	result, err := uc.userService.GetAllUser(ctx.Request.Context())
// 	if err != nil {
// 		res := common.BuildErrorResponse("Gagal Mendapatkan List User", err.Error(), common.EmptyObj{})
// 		ctx.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	res := common.BuildResponse(true, "Berhasil Mendapatkan List User", result)
// 	ctx.JSON(http.StatusOK, res)
// }

func (uc *userController) LoginUser(ctx *gin.Context) {
	var userLoginDTO dto.UserLoginDTO
	err := ctx.ShouldBind(&userLoginDTO)
	if err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, _ := uc.userService.Verify(ctx.Request.Context(), userLoginDTO.Username, userLoginDTO.Password)
	if !res {
		response := common.BuildErrorResponse("Gagal Login", "Username atau Password Salah", common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := uc.userService.FindUserByUsername(ctx.Request.Context(), userLoginDTO.Username)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	token := uc.jwtService.GenerateToken(user.ID)
	userResponse := entity.Authorization{
		Token: token,
	}

	response := common.BuildResponse(true, "Berhasil Login", userResponse)
	ctx.JSON(http.StatusOK, response)
}

// func (uc *userController) DeleteUser(ctx *gin.Context) {
// 	token := ctx.MustGet("token").(string)
// 	userID, err := uc.jwtService.GetUserIDByToken(token)
// 	if err != nil {
// 		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 		return
// 	}
// 	err = uc.userService.DeleteUser(ctx.Request.Context(), userID)
// 	if err != nil {
// 		res := common.BuildErrorResponse("Gagal Menghapus User", err.Error(), common.EmptyObj{})
// 		ctx.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	res := common.BuildResponse(true, "Berhasil Menghapus User", common.EmptyObj{})
// 	ctx.JSON(http.StatusOK, res)
// }

// func (uc *userController) UpdateUser(ctx *gin.Context) {
// 	var user dto.UserUpdateDto
// 	err := ctx.ShouldBind(&user)
// 	if err != nil {
// 		res := common.BuildErrorResponse("Gagal Mengupdate User", err.Error(), common.EmptyObj{})
// 		ctx.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	token := ctx.MustGet("token").(string)
// 	userID, err := uc.jwtService.GetUserIDByToken(token)
// 	if err != nil {
// 		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 		return
// 	}

// 	user.ID = userID
// 	err = uc.userService.UpdateUser(ctx.Request.Context(), user)
// 	if err != nil {
// 		res := common.BuildErrorResponse("Gagal Mengupdate User", err.Error(), common.EmptyObj{})
// 		ctx.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	res := common.BuildResponse(true, "Berhasil Mengupdate User", common.EmptyObj{})
// 	ctx.JSON(http.StatusOK, res)
// }

func (uc *userController) MeUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	result, err := uc.userService.MeUser(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) MeUserDecrypted(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	result, err := uc.userService.MeUserDecrypted(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) DecryptUserIDCard(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = uc.userService.DecryptUserIDCard(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Melakukan Dekripsi File", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Melakukan Dekripsi File", true)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) RequestAccess(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var requestData dto.AccessRequestCreateDto
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	requestedUser, err := uc.userService.FindUserByUsername(ctx.Request.Context(), requestData.RequestedUsername)
	if err != nil {
		response := common.BuildErrorResponse("User Tidak Ditemukan", "User Yang Diminta Aksesnya Tidak Ditemukan.", nil)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	if userID == requestedUser.ID {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Anda Tidak Dapat Merequest File Anda Sendiri", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	request, err := uc.userService.RequestAccess(ctx.Request.Context(), userID, requestedUser.ID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Membuat Akses Request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	responseData := dto.AccessRequestResponseDto{
		ID:              request.ID,
		UserID:          request.UserID,
		RequestedUserID: request.RequestedUserID,
		Status:          request.Status,
	}
	res := common.BuildResponse(true, "Akses Request Berhasil Dibuat", responseData)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetAccessRequests(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	requestType := ctx.Query("type")
	var result []entity.AccessRequest

	if requestType == "received" {
		result, err = uc.userService.GetReceivedAccessRequests(ctx.Request.Context(), userID)
	} else {
		result, err = uc.userService.GetSentAccessRequests(ctx.Request.Context(), userID)
	}

	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Akses Request", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if len(result) == 0 {
		res := common.BuildResponse(true, "Tidak Ada Access Request", []entity.AccessRequest{})
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan Akses Request", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateAccessRequestStatus(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	requestID, err := uuid.Parse(ctx.Param("request_id"))
	if err != nil {
		response := common.BuildErrorResponse("Invalid Request ID", "Request ID Tidak Valid", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var statusDto dto.AccessRequestChangeStatusDto
	if err := ctx.ShouldBindJSON(&statusDto); err != nil {
		response := common.BuildErrorResponse("Validation Error", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := validator.New().Struct(&statusDto); err != nil {
		response := common.BuildErrorResponse("Validation Error", "Status must be one of: pending, approved, denied", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.UpdateAccessRequestStatus(ctx.Request.Context(), userID, requestID, statusDto.Status)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Mengupdate Access Request Status", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.BuildResponse(true, "Access Request Status Berhasil Diupdate", nil)
	ctx.JSON(http.StatusOK, response)
}
