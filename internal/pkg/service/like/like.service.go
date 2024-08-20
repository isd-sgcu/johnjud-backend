package like

import "github.com/isd-sgcu/johnjud-gateway/internal/dto"

type Service interface {
	FindByUserId(string) ([]*dto.LikeResponse, *dto.ResponseErr)
	Create(*dto.CreateLikeRequest) (*dto.LikeResponse, *dto.ResponseErr)
	Delete(string) (*dto.DeleteLikeResponse, *dto.ResponseErr)
}
