package pet

import (
	"errors"
	"fmt"

	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/isd-sgcu/johnjud-backend/internal/image"
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

type Service interface {
	FindAll(req *dto.FindAllPetRequest, isAdmin bool) (*dto.FindAllPetResponse, *dto.ResponseErr)
	FindOne(id string) (*dto.PetResponse, *dto.ResponseErr)
	Create(req *dto.CreatePetRequest) (*dto.PetResponse, *dto.ResponseErr)
	Update(id string, req *dto.UpdatePetRequest) (*dto.PetResponse, *dto.ResponseErr)
	Delete(id string) (*dto.DeleteResponse, *dto.ResponseErr)
	ChangeView(id string, req *dto.ChangeViewPetRequest) (*dto.ChangeViewPetResponse, *dto.ResponseErr)
	Adopt(id string, req *dto.AdoptByRequest) (*dto.AdoptByResponse, *dto.ResponseErr)
}

type serviceImpl struct {
	repository   Repository
	imageService image.Service
}

func NewService(repository Repository, imageService image.Service) Service {
	return &serviceImpl{repository: repository, imageService: imageService}
}

func (s *serviceImpl) Delete(id string) (*dto.DeleteResponse, *dto.ResponseErr) {
	err := s.repository.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.NotFoundError("pet not found")
		}
		return nil, dto.InternalServerError("internal error")
	}
	return &dto.DeleteResponse{Success: true}, nil
}

func (s *serviceImpl) Update(id string, req *dto.UpdatePetRequest) (*dto.PetResponse, *dto.ResponseErr) {
	raw := UpdateDtoToModel(req)

	err := s.repository.Update(id, raw)
	if err != nil {
		return nil, dto.NotFoundError("pet not found")
	}

	images, apperr := s.imageService.FindByPetId(id)
	if apperr != nil {
		return nil, dto.InternalServerError("error querying image service")
	}

	return RawToDto(raw, images), nil
}

func (s *serviceImpl) ChangeView(id string, req *dto.ChangeViewPetRequest) (*dto.ChangeViewPetResponse, *dto.ResponseErr) {
	petData, apperr := s.FindOne(id)
	if apperr != nil {
		return nil, apperr
	}
	pet, err := DtoToRaw(petData)
	if err != nil {
		return nil, dto.InternalServerError("error converting dto to raw")
	}
	pet.IsVisible = req.Visible

	err = s.repository.Update(id, pet)
	if err != nil {
		return nil, dto.NotFoundError("pet not found")
	}

	return &dto.ChangeViewPetResponse{Success: true}, nil
}

func (s *serviceImpl) FindAll(req *dto.FindAllPetRequest, isAdmin bool) (*dto.FindAllPetResponse, *dto.ResponseErr) {
	var pets []*model.Pet
	imagesList := make(map[string][]*dto.ImageResponse)
	metaData := dto.FindAllMetadata{}

	err := s.repository.FindAll(&pets, isAdmin)
	if err != nil {
		log.Error().Err(err).Str("service", "event").Str("module", "find all").Msg("Error while querying all events")
		return nil, dto.InternalServerError("error querying all pets")
	}

	FilterPet(&pets, req)
	PaginatePets(&pets, req.Page, req.PageSize, &metaData)

	for _, pet := range pets {
		images, err := s.imageService.FindByPetId(pet.ID.String())
		if err != nil {
			return nil, dto.InternalServerError("error querying image service")
		}
		imagesList[pet.ID.String()] = images
	}
	petWithImages, err := RawToDtoList(&pets, imagesList, req)
	if err != nil {
		return nil, dto.InternalServerError(fmt.Sprintf("error converting raw to dto list: %v", err))
	}

	// images, errSvc := s.imageService.FindAll()
	// if errSvc != nil {
	// 	return nil, errSvc
	// }

	// allImagesList := ImageList(images)
	// findAllDto := ProtoToDtoList(res.Pets, allImagesList, isAdmin)
	// metaData := MetadataProtoToDto(res.Metadata)

	return &dto.FindAllPetResponse{Pets: petWithImages, Metadata: &metaData}, nil
}

func (s *serviceImpl) FindOne(id string) (*dto.PetResponse, *dto.ResponseErr) {
	var pet model.Pet

	err := s.repository.FindOne(id, &pet)
	if err != nil {
		log.Error().Err(err).
			Str("service", "pet").Str("module", "find one").Str("id", id).Msg("Not found")
		return nil, dto.NotFoundError("pet not found")
	}

	images, apperr := s.imageService.FindByPetId(id)
	if apperr != nil {
		return nil, apperr
	}

	return RawToDto(&pet, images), nil
}

func (s *serviceImpl) Create(req *dto.CreatePetRequest) (*dto.PetResponse, *dto.ResponseErr) {
	raw := CreateDtoToModel(req)

	err := s.repository.Create(raw)
	if err != nil {
		return nil, dto.InternalServerError("failed to create pet")
	}

	assignReq := &dto.AssignPetRequest{
		PetId: raw.ID.String(),
		Ids:   req.Images,
	}
	_, apperr := s.imageService.AssignPet(assignReq)
	if apperr != nil {
		return nil, apperr
	}

	images, apperr := s.imageService.FindByPetId(raw.ID.String())
	if apperr != nil {
		return nil, apperr
	}

	return RawToDto(raw, images), nil
}

func (s *serviceImpl) Adopt(id string, req *dto.AdoptByRequest) (*dto.AdoptByResponse, *dto.ResponseErr) {
	dtoPet, apperr := s.FindOne(id)
	if apperr != nil {
		return nil, apperr
	}

	pet, err := DtoToRaw(dtoPet)
	if err != nil {
		return nil, dto.InternalServerError("error converting dto to raw")
	}
	pet.Owner = req.UserID

	err = s.repository.Update(id, pet)
	if err != nil {
		return nil, dto.NotFoundError("pet not found")
	}

	return &dto.AdoptByResponse{Success: true}, nil
}
