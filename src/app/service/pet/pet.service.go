package pet

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/pet"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	petproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imgproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	petClient petproto.PetServiceClient
}

func NewService(petClient petproto.PetServiceClient) *Service {
	return &Service{
		petClient: petClient,
	}
}

func (s *Service) FindAll() (result []*dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindAll(ctx, &petproto.FindAllPetRequest{})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "find all").
			Msg(st.Message())
		switch st.Code() {
		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		}
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InternalErrorMessage,
			Data:       nil,
		}
	}
	imagesList := utils.MockImageList(len(res.Pets))
	findAllResponse := RawToDtoList(res.Pets, imagesList)
	return findAllResponse, nil
}

func (s *Service) FindOne(id string) (result *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindOne(ctx, &petproto.FindOnePetRequest{Id: id})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "find one").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}

		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}
	images := utils.MockImageList(1)[0]
	findOneResponse := RawToDto(res.Pet, images)
	return findOneResponse, nil
}

func (s *Service) Create(in *dto.CreatePetRequest) (ressult *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := CreateDtoToRaw(in)

	res, errRes := s.petClient.Create(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "create").
			Msg(st.Message())
		switch st.Code() {
		case codes.InvalidArgument:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    constant.InvalidArgument,
				Data:       nil,
			}
		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}
	images := utils.MockImageList(1)[0]
	createPetResponse := RawToDto(res.Pet, images)
	return createPetResponse, nil
}

func (s *Service) Update(id string, in *dto.UpdatePetRequest) (result *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := UpdateDtoToRaw(id, in)

	res, errRes := s.petClient.Update(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "update").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}
		case codes.InvalidArgument:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    constant.InvalidArgument,
				Data:       nil,
			}
		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}
	images := utils.MockImageList(1)[0]
	updatePetResponse := RawToDto(res.Pet, images)
	return updatePetResponse, nil
}

func (s *Service) Delete(id string) (result *dto.DeleteResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.Delete(ctx, &petproto.DeletePetRequest{
		Id: id,
	})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "delete").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return &dto.DeleteResponse{
					Success: false,
				}, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    constant.PetNotFoundMessage,
					Data:       nil,
				}
		case codes.Unavailable:
			return &dto.DeleteResponse{
					Success: false,
				}, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    constant.UnavailableServiceMessage,
					Data:       nil,
				}
		}
		return &dto.DeleteResponse{
				Success: false,
			}, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
	}
	return &dto.DeleteResponse{
		Success: res.Success,
	}, nil
}

func (s *Service) ChangeView(id string, in *dto.ChangeViewPetRequest) (result *dto.ChangeViewPetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.ChangeView(ctx, &petproto.ChangeViewPetRequest{
		Id:      id,
		Visible: in.Visible,
	})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "change view").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return &dto.ChangeViewPetResponse{
					Success: false,
				}, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    constant.PetNotFoundMessage,
					Data:       nil,
				}
		case codes.Unavailable:
			return &dto.ChangeViewPetResponse{
					Success: false,
				}, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    constant.UnavailableServiceMessage,
					Data:       nil,
				}
		default:
			return &dto.ChangeViewPetResponse{
					Success: false,
				}, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    constant.InternalErrorMessage,
					Data:       nil,
				}
		}
	}
	return &dto.ChangeViewPetResponse{
		Success: res.Success,
	}, nil
}

func CreateDtoToRaw(in *dto.CreatePetRequest) *petproto.CreatePetRequest {
	return &petproto.CreatePetRequest{
		Pet: &petproto.Pet{
			Type:         in.Type,
			Species:      in.Species,
			Name:         in.Name,
			Birthdate:    in.Birthdate,
			Gender:       petproto.Gender(in.Gender),
			Habit:        in.Habit,
			Caption:      in.Caption,
			Images:       []*imgproto.Image{},
			Status:       petproto.PetStatus(in.Status),
			IsSterile:    *in.IsSterile,
			IsVaccinated: *in.IsVaccinated,
			IsVisible:    *in.IsVisible,
			IsClubPet:    *in.IsClubPet,
			Background:   in.Background,
			Address:      in.Address,
			Contact:      in.Contact,
			AdoptBy:      in.AdoptBy,
		},
	}
}

func UpdateDtoToRaw(id string, in *dto.UpdatePetRequest) *petproto.UpdatePetRequest {
	req := &petproto.UpdatePetRequest{
		Pet: &petproto.Pet{
			Id:         id,
			Type:       in.Type,
			Species:    in.Species,
			Name:       in.Name,
			Birthdate:  in.Birthdate,
			Gender:     petproto.Gender(in.Gender),
			Habit:      in.Habit,
			Caption:    in.Caption,
			Images:     []*imgproto.Image{},
			Status:     petproto.PetStatus(in.Status),
			Background: in.Background,
			Address:    in.Address,
			Contact:    in.Contact,
			AdoptBy:    in.AdoptBy,
		},
	}

	if in.IsClubPet == nil {
		req.Pet.IsClubPet = false
	} else {
		req.Pet.IsClubPet = *in.IsClubPet
	}

	if in.IsSterile == nil {
		req.Pet.IsSterile = false
	} else {
		req.Pet.IsSterile = *in.IsSterile
	}

	if in.IsVaccinated == nil {
		req.Pet.IsVaccinated = false
	} else {
		req.Pet.IsVaccinated = *in.IsVaccinated
	}

	if in.IsVisible == nil {
		req.Pet.IsVisible = false
	} else {
		req.Pet.IsVisible = *in.IsVisible
	}

	return req
}

func RawToDto(in *petproto.Pet, images []*imgproto.Image) *dto.PetResponse {
	pet := &dto.PetResponse{
		Id:           in.Id,
		Type:         in.Type,
		Species:      in.Species,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       pet.Gender(in.Gender),
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       pet.Status(in.Status),
		IsSterile:    &in.IsSterile,
		IsVaccinated: &in.IsVaccinated,
		IsVisible:    &in.IsVisible,
		IsClubPet:    &in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
		AdoptBy:      in.AdoptBy,
		Images:       extractImages(images),
	}
	return pet
}

func RawToDtoList(in []*petproto.Pet, imagesList [][]*imgproto.Image) []*dto.PetResponse {
	var resp []*dto.PetResponse
	for i, p := range in {
		pet := &dto.PetResponse{
			Id:           p.Id,
			Type:         p.Type,
			Species:      p.Species,
			Name:         p.Name,
			Birthdate:    p.Birthdate,
			Gender:       pet.Gender(p.Gender),
			Habit:        p.Habit,
			Caption:      p.Caption,
			Status:       pet.Status(p.Status),
			IsSterile:    &p.IsSterile,
			IsVaccinated: &p.IsVaccinated,
			IsVisible:    &p.IsVisible,
			IsClubPet:    &p.IsClubPet,
			Background:   p.Background,
			Address:      p.Address,
			Contact:      p.Contact,
			AdoptBy:      p.AdoptBy,
			Images:       extractImages(imagesList[i]),
		}
		resp = append(resp, pet)
	}
	return resp
}

func extractImages(images []*imgproto.Image) []dto.ImageResponse {
	var result []dto.ImageResponse
	for _, img := range images {
		result = append(result, dto.ImageResponse{
			Id:  img.Id,
			Url: img.ImageUrl,
		})
	}
	return result
}
