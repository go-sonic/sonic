package wp

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (u *UserHandler) List(ctx *gin.Context) (interface{}, error) {
	allUser, err := u.UserService.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}

	userDTOList := make([]*dto.User, len(allUser))
	for _, user := range allUser {
		userDTO := u.UserService.ConvertToDTO(ctx, user)
		userDTOList = append(userDTOList, userDTO)
	}
	return userDTOList, nil
}
