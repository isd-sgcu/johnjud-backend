package pet

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
)

type Service interface {
	FindAll(*dto.FindAllPetRequest, bool) (*dto.FindAllPetResponse, *dto.ResponseErr)
	FindOne(string) (*dto.PetResponse, *dto.ResponseErr)
	Create(*dto.CreatePetRequest) (*dto.PetResponse, *dto.ResponseErr)
	Update(string, *dto.UpdatePetRequest) (*dto.PetResponse, *dto.ResponseErr)
	Delete(string) (*dto.DeleteResponse, *dto.ResponseErr)
	ChangeView(string, *dto.ChangeViewPetRequest) (*dto.ChangeViewPetResponse, *dto.ResponseErr)
	Adopt(string, *dto.AdoptByRequest) (*dto.AdoptByResponse, *dto.ResponseErr)
}
