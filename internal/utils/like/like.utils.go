package like

import (
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	likeProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
)

func ProtoToDto(in *likeProto.Like) *dto.LikeResponse {
	return &dto.LikeResponse{
		UserID: in.UserId,
		PetID:  in.PetId,
	}
}

func ProtoToDtoList(in []*likeProto.Like) []*dto.LikeResponse {
	var res []*dto.LikeResponse
	for _, i := range in {
		res = append(res, &dto.LikeResponse{
			UserID: i.UserId,
			PetID:  i.PetId,
		})
	}
	return res
}

func CreateDtoToProto(in *dto.CreateLikeRequest) *likeProto.CreateLikeRequest {
	return &likeProto.CreateLikeRequest{
		Like: &likeProto.Like{
			UserId: in.UserID,
			PetId:  in.PetID,
		},
	}
}
