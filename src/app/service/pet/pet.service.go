package pet

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/pet"
	petproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
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
	findAllResponse := utils.RawToDtoList(res.Pets, imagesList)
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
	findOneResponse := utils.RawToDto(res.Pet, images)
	return findOneResponse, nil
}

func (s *Service) Create(in *dto.CreatePetRequest) (ressult *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := utils.CreateDtoToRaw(in)

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
	createPetResponse := utils.RawToDto(res.Pet, images)
	return createPetResponse, nil
}

func (s *Service) Update(id string, in *dto.UpdatePetRequest) (result *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := utils.UpdateDtoToRaw(id, in)

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
	updatePetResponse := utils.RawToDto(res.Pet, images)
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
