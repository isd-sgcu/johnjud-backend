package image

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

type Service interface {
	FindByPetId(string) ([]*proto.Image, *dto.ResponseErr)
	Upload(*dto.ImageDto) (*proto.Image, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}
