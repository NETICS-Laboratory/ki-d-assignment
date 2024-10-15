package controller

import (
	"ki-d-assignment/common"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"

	// "ki-d-assignment/entity"
	"ki-d-assignment/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
