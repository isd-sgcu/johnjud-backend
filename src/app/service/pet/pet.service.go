package pet

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
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
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "find all").
			Msg("Error while find all pets")
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    constant.UnavailableServiceMessage,
			Data:       nil,
		}
	}
	log.Info().
		Str("service", "pet").
		Str("module", "find all").
		Msg("Find pet success")
	return res.Pets, nil
}

func (s *Service) FindOne(id string) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindOne(ctx, &proto.FindOnePetRequest{Id: id})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		switch st.Code() {
		case codes.NotFound:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "find one").
				Str("pet_id", id).
				Msg("Not found")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}
		default:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("pet_id", id).
				Msg("Error while connecting to service")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		}
	}
	log.Info().
		Str("service", "pet").
		Str("module", "find one").
		Str("pet_id", id).
		Msg("Find pet success")
	return res.Pet, nil
}

func (s *Service) Create(in *dto.CreatePetDto) (ressult *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := DtoToRaw(in.Pet)

	res, errRes := s.petClient.Create(ctx, &proto.CreatePetRequest{Pet: request})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		switch st.Code() {
		case codes.InvalidArgument:
			log.Error().
				Err(errRes).
				Str("service", "user").
				Str("module", "create").
				Msg(constant.InvalidArgument)
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    constant.InvalidArgument,
				Data:       nil,
			}
		case codes.Unavailable:
			log.Error().
				Err(errRes).
				Str("service", "user").
				Str("module", "create").
				Msg(constant.UnavailableServiceMessage)
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			log.Error().
				Err(errRes).
				Str("service", "user").
				Str("module", "create").
				Msg(constant.InternalErrorMessage)

			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}
	return res.Pet, nil
}

func (s *Service) Update(id string, in *dto.UpdatePetDto) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &proto.UpdatePetRequest{
		Pet: DtoToRaw(in.Pet),
	}

	res, errRes := s.petClient.Update(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		switch st.Code() {
		case codes.NotFound:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "update").
				Msg(constant.PetNotFoundMessage)
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}
		case codes.InvalidArgument:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "update").
				Msg(constant.InvalidArgument)
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    constant.InvalidArgument,
				Data:       nil,
			}
		case codes.Unavailable:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "update").
				Msg("Error while connecting to service")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "update").
				Msg(constant.InternalErrorMessage)
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

		switch st.Code() {
		case codes.NotFound:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "delete").
				Msg(constant.PetNotFoundMessage)
			return false, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}

		case codes.Unavailable:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "delete").
				Msg(constant.UnavailableServiceMessage)
			return false, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		}
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "delete").
			Msg(constant.InternalErrorMessage)
		return false, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InternalErrorMessage,
			Data:       nil,
		}
	}
	return res.Success, nil
}

func (s *Service) ChangeView(id string, in *dto.ChangeViewPetDto) (result bool, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.ChangeView(ctx, &proto.ChangeViewPetRequest{
		Id:      id,
		Visible: in.Visible,
	})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		switch st.Code() {
		case codes.NotFound:
			return false, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.PetNotFoundMessage,
				Data:       nil,
			}
		case codes.Unavailable:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "change view").
				Msg(constant.UnavailableServiceMessage)
			return false, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			log.Error().
				Err(errRes).
				Str("service", "pet").
				Str("module", "change view").
				Msg(constant.InternalErrorMessage)
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
		Status:       proto.PetStatus(in.Status),
		ImageUrls:    []string{},
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		IsClubPet:    in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
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
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		IsClubPet:    in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
	}
}
