package pet

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
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

// TODO: change reutnr type to []*dto.PetRespone
func (s *Service) FindAll() (result []*petproto.Pet, err *dto.ResponseErr) {
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
	// findAllResponse := RawToDtoList(res, nil)
	return res.Pets, nil
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
	findOneResponse := RawToDto(res.Pet)
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
	fmt.Println(res)
	createPetResponse := RawToDto(res.Pet)
	return createPetResponse, nil
}

func (s *Service) Update(id string, in *dto.UpdatePetRequest) (result *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := UpdateDtoToRaw(in)

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
	fmt.Println(res)
	updatePetResponse := RawToDto(res.Pet)
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

func UpdateDtoToRaw(in *dto.UpdatePetRequest) *petproto.UpdatePetRequest {
	return &petproto.UpdatePetRequest{
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

func RawToDto(in *petproto.Pet) *dto.PetResponse {
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
	}
	return pet
}

// TODO: add `Images` in &dto.PetResponse
func RawToDtoList(in []*petproto.Pet, images [][]*imgproto.Image) []*dto.PetResponse {
	var resp []*dto.PetResponse
	for _, p := range in {
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
		}
		resp = append(resp, pet)
	}
	return resp
}
