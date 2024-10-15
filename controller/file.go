package controller

import (
	"ki-d-assignment/common"
	"ki-d-assignment/dto"
	"ki-d-assignment/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileController interface {
	UploadFile(ctx *gin.Context)
	GetUserFiles(ctx *gin.Context)
	GetUserFileDecrypted(ctx *gin.Context)
}

type fileController struct {
	jwtService  service.JWTService
	userService service.UserService
	fileService service.FileService
}

func NewFileController(fs service.FileService, us service.UserService, jwts service.JWTService) FileController {
	return &fileController{
		fileService: fs,
		userService: us,
		jwtService:  jwts,
	}
}

func (fc *fileController) UploadFile(ctx *gin.Context) {

	token := ctx.MustGet("token").(string)
	userID, err := fc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var fileDTO dto.FileCreateDto
	err = ctx.ShouldBind(&fileDTO)
	if err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	uploadedFile, err := fc.fileService.UploadFile(ctx.Request.Context(), fileDTO, userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupload File", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "File Berhasil Diunggah", uploadedFile)
	ctx.JSON(http.StatusOK, res)
}

func (fc *fileController) GetUserFiles(ctx *gin.Context) {

	token := ctx.MustGet("token").(string)
	userID, err := fc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	files, err := fc.fileService.GetUserFiles(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan File", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "File Ditemukan", files)
	ctx.JSON(http.StatusOK, res)
}

func (fc *fileController) GetUserFileDecrypted(ctx *gin.Context) {

	token := ctx.MustGet("token").(string)
	userID, err := fc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var fileDTO dto.FileDecryptByIDDto
	err = ctx.ShouldBind(&fileDTO)
	if err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	files, err := fc.fileService.GetUserFileDecryptedByID(ctx.Request.Context(), fileDTO.ID, userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan File", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "File Ditemukan", files)
	ctx.JSON(http.StatusOK, res)
}