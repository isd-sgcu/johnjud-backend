package image

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

type Service interface {
	FindByPetId(id string) ([]*proto.Image, *dto.ResponseErr)
	Upload(in *dto.ImageDto) (*proto.Image, *dto.ResponseErr)
	Delete(id string) (bool, *dto.ResponseErr)
}
