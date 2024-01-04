package pet

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	image_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	petClient proto.PetServiceClient
}

func NewService(petClient proto.PetServiceClient) *Service {
	return &Service{
		petClient: petClient,
	}
}

func (s *Service) FindAll() (result []*proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindAll(ctx, &proto.FindAllPetRequest{})
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
	_ = &image_proto.Image{}
	return res.Pets, nil
}

func (s *Service) FindOne(id string) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindOne(ctx, &proto.FindOnePetRequest{Id: id})
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
	return res.Pet, nil
}

func (s *Service) Create(in *dto.CreatePetRequest) (ressult *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := DtoToRaw(in.Pet)

	res, errRes := s.petClient.Create(ctx, &proto.CreatePetRequest{Pet: request})
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
	return res.Pet, nil
}

func (s *Service) Update(id string, in *dto.UpdatePetRequest) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &proto.UpdatePetRequest{
		Pet: &proto.Pet{
			Id:           id,
			Type:         in.Pet.Type,
			Species:      in.Pet.Species,
			Name:         in.Pet.Name,
			Birthdate:    in.Pet.Birthdate,
			Gender:       proto.Gender(in.Pet.Gender),
			Habit:        in.Pet.Habit,
			Caption:      in.Pet.Caption,
			Images:       []*image_proto.Image{},
			Status:       proto.PetStatus(in.Pet.Status),
			IsSterile:    *in.Pet.IsSterile,
			IsVaccinated: *in.Pet.IsSterile,
			IsVisible:    *in.Pet.IsVaccinated,
			IsClubPet:    *in.Pet.IsClubPet,
			Background:   in.Pet.Background,
			Address:      in.Pet.Address,
			Contact:      in.Pet.Contact,
			AdoptBy:      in.Pet.AdoptBy,
		},
	}

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
	return res.Pet, nil
}

func (s *Service) Delete(id string) (result bool, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.Delete(ctx, &proto.DeletePetRequest{
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
			return false, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}
		case codes.Unavailable:
			return false, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		}
		return false, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InternalErrorMessage,
			Data:       nil,
		}
	}
	return res.Success, nil
}

func (s *Service) ChangeView(id string, in *dto.ChangeViewPetRequest) (result bool, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.ChangeView(ctx, &proto.ChangeViewPetRequest{
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
			return false, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}
		case codes.Unavailable:
			return false, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return false, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}
	return res.Success, nil
}

func DtoToRaw(in *dto.PetDto) *proto.Pet {
	return &proto.Pet{
		Id:           in.Id,
		Type:         in.Type,
		Species:      in.Species,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       proto.Gender(in.Gender),
		Habit:        in.Habit,
		Caption:      in.Caption,
		Images:       []*image_proto.Image{},
		Status:       proto.PetStatus(in.Status),
		IsSterile:    *in.IsSterile,
		IsVaccinated: *in.IsVaccinated,
		IsVisible:    *in.IsVisible,
		IsClubPet:    *in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
		AdoptBy:      in.AdoptBy,
	}
}

func RawToDto(in *proto.Pet) *dto.PetDto {
	return &dto.PetDto{
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
}
