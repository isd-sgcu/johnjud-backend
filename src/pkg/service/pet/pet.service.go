package pet

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
)

type Service interface {
	FindAll() ([]*proto.Pet, *dto.ResponseErr)
	FindOne(string) (*proto.Pet, *dto.ResponseErr)
	Create(*dto.CreatePetRequest) (*proto.Pet, *dto.ResponseErr)
	Update(string, *dto.UpdatePetRequest) (*proto.Pet, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
	ChangeView(string, *dto.ChangeViewPetRequest) (bool, *dto.ResponseErr)
	Adopt(*dto.AdoptDto) (bool, *dto.ResponseErr)
}
