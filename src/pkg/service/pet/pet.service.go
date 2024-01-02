package pet

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
)

type Service interface {
	FindAll() (result []*proto.Pet, err *dto.ResponseErr)
	FindOne(id string) (result *proto.Pet, err *dto.ResponseErr)
	Create(in *dto.CreatePetRequest) (ressult *proto.Pet, err *dto.ResponseErr)
	Update(id string, in *dto.UpdatePetRequest) (result *proto.Pet, err *dto.ResponseErr)
	Delete(id string) (result bool, err *dto.ResponseErr)
	ChangeView(id string, in *dto.ChangeViewPetRequest) (result bool, err *dto.ResponseErr)
}
