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
	VerifyFileSignature(ctx *gin.Context)
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
		response := common.BuildErrorResponse("Failed to process request", "Invalid token", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var fileDTO dto.FileCreateDto
	if err := ctx.ShouldBind(&fileDTO); err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	uploadedFile, err := fc.fileService.UploadFile(ctx.Request.Context(), fileDTO, userID)
	if err != nil {
		res := common.BuildErrorResponse("Failed to upload file", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "File uploaded successfully", uploadedFile)
	ctx.JSON(http.StatusOK, res)
}

func (fc *fileController) GetUserFiles(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := fc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Failed to process request", "Invalid token", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	files, err := fc.fileService.GetUserFiles(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Failed to retrieve files", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "Files retrieved successfully", files)
	ctx.JSON(http.StatusOK, res)
}

func (fc *fileController) GetUserFileDecrypted(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := fc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Failed to process request", "Invalid token", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var fileDTO dto.FileDecryptByIDDto
	if err := ctx.ShouldBind(&fileDTO); err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	file, err := fc.fileService.GetUserFileDecryptedByID(ctx.Request.Context(), fileDTO.ID, userID)
	if err != nil {
		res := common.BuildErrorResponse("Failed to retrieve decrypted file", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "Decrypted file retrieved successfully", file)
	ctx.JSON(http.StatusOK, res)
}

func (fc *fileController) VerifyFileSignature(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := fc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Failed to process request", "Invalid token", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	var fileDTO dto.FileVerifySignatureDto
	if err := ctx.ShouldBind(&fileDTO); err != nil {
		res := common.BuildErrorResponse("Validation Error", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	check, err := fc.fileService.CheckFileDigitalSignature(ctx.Request.Context(), fileDTO.FileID, userID, fileDTO.Signature)

	if err != nil {
		res := common.BuildErrorResponse("Failed to verify file signature", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := common.BuildResponse(true, "File signature verified successfully", check)
	ctx.JSON(http.StatusOK, res)
}
