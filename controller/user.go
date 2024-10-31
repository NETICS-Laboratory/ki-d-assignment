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
	LoginUser(ctx *gin.Context)
	MeUser(ctx *gin.Context)
	MeUserDecrypted(ctx *gin.Context)
	DecryptUserIDCard(ctx *gin.Context)

	RequestAccess(ctx *gin.Context)
	GetAccessRequests(ctx *gin.Context)
	UpdateAccessRequestStatus(ctx *gin.Context)

	DecryptKeys(ctx *gin.Context)
	AccessPrivateData(ctx *gin.Context)
}

type userController struct {
	jwtService  service.JWTService
	userService service.UserService
	fileService service.FileService
}

func NewUserController(us service.UserService, jwts service.JWTService, fs service.FileService) UserController {
	return &userController{
		userService: us,
		jwtService:  jwts,
		fileService: fs,
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

func (uc *userController) DecryptKeys(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var decryptRequest dto.RSAKeyDecryptDto
	if err := ctx.ShouldBindJSON(&decryptRequest); err != nil {
		response := common.BuildErrorResponse("Validation Error", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	decryptedKey, decryptedKey8Byte, err := uc.userService.DecryptKeys(ctx.Request.Context(), userID, decryptRequest.EncryptedKey, decryptRequest.EncryptedKey8Byte)
	if err != nil {
		res := common.BuildErrorResponse("Decryption Error", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Respond with decrypted keys
	responseData := dto.RSAKeyDecryptResponseDto{
		DecryptedKey:      decryptedKey,
		DecryptedKey8Byte: decryptedKey8Byte,
	}
	response := common.BuildResponse(true, "Keys decrypted successfully", responseData)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) AccessPrivateData(ctx *gin.Context) {
	var accessKeys dto.AccessPrivateDataRequestDto

	if err := ctx.ShouldBindJSON(&accessKeys); err != nil {
		response := common.BuildErrorResponse("Validation Error", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token := ctx.MustGet("token").(string)
	// TODO: Change for later
	_, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	requestedUser, err := uc.userService.FindUserByUsername(ctx.Request.Context(), accessKeys.RequestedUserUsername)
	if err != nil {
		response := common.BuildErrorResponse("Gagal mendapatkan user", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	resfile, err := uc.fileService.GetRequestedUserData(ctx.Request.Context(), requestedUser, accessKeys.SecretKey, accessKeys.SecretKey8Byte)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Mendapatkan File", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	resdata, err := uc.userService.MeUserDecrypted(ctx.Request.Context(), requestedUser.ID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	responseData := dto.AccessPrivateDataResponseDto{
		RequestedUser: dto.UserRequestDecryptedDto{
			ID:       resdata.ID,
			Username: resdata.Username,
			Name:     resdata.Name,
			Email:    resdata.Email,
			NoTelp:   resdata.NoTelp,
			Address:  resdata.Address,
			ID_Card:  resdata.ID_Card,
		},
		Files: resfile,
	}

	response := common.BuildResponse(true, "Data Berhasil Didapatkan", responseData)
	ctx.JSON(http.StatusOK, response)
}
