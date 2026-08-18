package dto

import "notesapp/entity"

type UserResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ToUserResponse(userEntity entity.UserEntity) UserResponse {
	return UserResponse{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}
}
