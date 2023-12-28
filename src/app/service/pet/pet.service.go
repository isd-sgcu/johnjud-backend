package pet

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
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
			Message:    "Service is down",
			Data:       nil,
		}
	}
	return res.Pets, nil
}

func (s *Service) FindOne(request *dto.FindOnePetDto) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindOne(ctx, &proto.FindOnePetRequest{Id: request.Id})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "pet").
					Str("module", "find one").
					Str("pet_id", request.Id).
					Msg("Not found")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				log.Error().
					Err(errRes).
					Str("service", "pet").
					Str("pet_id", request.Id).
					Msg("Invaild pet id")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "pet").
					Str("pet_id", request.Id).
					Msg("Error while connecting to service")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
		log.Error().
			Err(errRes).
			Str("service", "group").
			Str("module", "find one").
			Str("per_id", request.Id).
			Msg("Error while connecting to service")
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}
	log.Info().
		Str("service", "pet").
		Str("module", "find one").
		Str("pet_id", request.Id).
		Msg("Find pet success")

	return res.Pet, nil
}

func (s *Service) Create(in *dto.CreatePetDto) (ressult *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	petDto := &proto.Pet{
		Type:         in.Pet.Type,
		Species:      in.Pet.Species,
		Name:         in.Pet.Name,
		Birthdate:    in.Pet.Birthdate,
		Gender:       proto.Gender(in.Pet.Gender),
		Habit:        in.Pet.Habit,
		Caption:      in.Pet.Caption,
		Status:       proto.PetStatus(in.Pet.Status),
		ImageUrls:    []string{},
		IsSterile:    in.Pet.IsSterile,
		IsVaccinated: in.Pet.IsVaccinated,
		IsVisible:    in.Pet.IsVisible,
		IsClubPet:    in.Pet.IsClubPet,
		Background:   in.Pet.Background,
		Address:      in.Pet.Address,
		Contact:      in.Pet.Contact,
	}

	res, errRes := s.petClient.Create(ctx, &proto.CreatePetRequest{Pet: petDto})
	if errRes != nil {
		log.Error().
			Err(errRes).
			Str("service", "user").
			Str("module", "create").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}
	return res.Pet, nil
}

func (s *Service) Update(id string, in *dto.UpDatePetDto) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	petReq := &proto.UpdatePetRequest{
		Pet: &proto.Pet{
			Id:           id,
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
		},
	}

	res, errRes := s.petClient.Update(ctx, petReq)
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "pet").
					Str("module", "update").
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "update").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}
	return res.Pet, nil
}

func (s *Service) Delete(id string) (result bool, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.Delete(ctx, &proto.DeletePetRequest{})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return false, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "pet").
					Str("module", "delete").
					Msg("Error while connecting to service")

				return false, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "delete").
			Msg("Error while connecting to service")
		return false, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}
	return res.Success, nil
}
