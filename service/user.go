package service

import (
	"context"
	"ki-d-assignment/dto"
	"ki-d-assignment/entity"
	"ki-d-assignment/helpers"
	"ki-d-assignment/repository"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	//FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	//Verify(ctx context.Context, email string, password string) (bool, error)
	//CheckUser(ctx context.Context, email string) ( bool, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (error)
	UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) (error)
	//MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func(us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error) {
	var user entity.User
	
	user.username_AES = userDTO.Username
	user.username_DES = userDTO.Username
	user.username_RC4 = userDTO.Username
	user.password_AES = userDTO.Password
	user.password_DES = userDTO.Password
	user.password_RC4 = userDTO.Password

	result, err := u.UserRepository.RegisterUser(ctx, user)
		if err != nil {
			return entity.User{}, err
		}
	
		return result, nil
}

func(us *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	result, err := u.UserRepository.GetUserByID(ctx, userID)
		if err != nil {
			return entity.User{}, err
		}
	
		return result, nil
}

// func(us *userService) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
// 	return us.userRepository.FindUserByEmail(ctx, email)
// }

// func(us *userService) Verify(ctx context.Context, email string, password string) (bool, error) {
// 	res, err := us.userRepository.FindUserByEmail(ctx, email)
// 	if err != nil {
// 		return false, err
// 	}
// 	CheckPassword, err := helpers.CheckPassword(res.Password, []byte(password))
// 	if err != nil {
// 		return false, err
// 	}
// 	if res.Email == email && CheckPassword {
// 		return true, nil
// 	}
// 	return false, nil
// }

// func(us *userService) CheckUser(ctx context.Context, email string) (bool, error) {
// 	result, err := us.userRepository.FindUserByEmail(ctx, email)
// 	if err != nil {
// 		return false, err
// 	}

// 	if result.Email == "" {
// 		return false, nil
// 	}
// 	return true, nil
// }

func(us *userService) DeleteUser(ctx context.Context, userID uuid.UUID) (error) {
	err := u.UserRepository.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func(us *userService) UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) (error) {
	var user entity.User

	user.ID = userDTO.ID
	user.username_AES = userDTO.Username
	user.username_DES = userDTO.Username
	user.username_RC4 = userDTO.Username
	user.password_AES = userDTO.Password
	user.password_DES = userDTO.Password
	user.password_RC4 = userDTO.Password
	
	result, err := u.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return entity.User{}, err
	}
	
	return result, nil
}

// func(us *userService) MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error) {
// 	return us.userRepository.FindUserByID(ctx, userID)
// }