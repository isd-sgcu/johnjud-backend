package pet

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"

	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/pet"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/src/pkg/service/image"
	petproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	petClient    petproto.PetServiceClient
	imageService imageSvc.Service
}

func NewService(petClient petproto.PetServiceClient, imageService imageSvc.Service) *Service {
	return &Service{
		petClient:    petClient,
		imageService: imageService,
	}
}

func (s *Service) FindAll(in *dto.FindAllPetRequest, isAdmin bool) (result *dto.FindAllPetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.FindAll(ctx, utils.FindAllDtoToProto(in, isAdmin))
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
	log.Info().
		Str("service", "pet").
		Str("module", "petClient.FindAll").Interface("res", res).Msg("")

	images, errSvc := s.imageService.FindAll()
	if errSvc != nil {
		return nil, errSvc
	}

	imagesList := utils.ImageList(images)
	findAllDto := utils.ProtoToDtoList(res.Pets, imagesList, isAdmin)
	metaData := utils.MetadataProtoToDto(res.Metadata)
	log.Info().
		Str("service", "pet").
		Str("module", "findAllDto").Interface("dto", findAllDto).Msg("")

	return &dto.FindAllPetResponse{
		Pets:     findAllDto,
		Metadata: metaData,
	}, nil
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

	imgRes, imgErrRes := s.imageService.FindByPetId(res.Pet.Id)
	if imgErrRes != nil {
		return nil, imgErrRes
	}

	findOneResponse := utils.ProtoToDto(res.Pet, imgRes)
	return findOneResponse, nil
}

func (s *Service) Create(in *dto.CreatePetRequest) (result *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := utils.CreateDtoToProto(in)

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
				Message:    constant.InvalidArgumentMessage,
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

	_, assignErr := s.imageService.AssignPet(&dto.AssignPetRequest{
		Ids:   in.Images,
		PetId: res.Pet.Id,
	})
	if assignErr != nil {
		return nil, assignErr
	}

	imgRes, imgErrRes := s.imageService.FindByPetId(res.Pet.Id)
	if imgErrRes != nil {
		return nil, imgErrRes
	}

	createPetResponse := utils.ProtoToDto(res.Pet, imgRes)
	return createPetResponse, nil
}

func (s *Service) Update(id string, in *dto.UpdatePetRequest) (result *dto.PetResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := utils.UpdateDtoToProto(id, in)

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
				Message:    constant.InvalidArgumentMessage,
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

	images, errSvc := s.imageService.FindByPetId(res.Pet.Id)
	if errSvc != nil {
		return nil, errSvc
	}

	updatePetResponse := utils.ProtoToDto(res.Pet, images)
	return updatePetResponse, nil
}

func (s *Service) Delete(id string) (result *dto.DeleteResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, errSvc := s.imageService.DeleteByPetId(id)
	if errSvc != nil {
		return nil, errSvc
	}

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

func (s *Service) Adopt(petId string, in *dto.AdoptByRequest) (result *dto.AdoptByResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.petClient.AdoptPet(ctx, &petproto.AdoptPetRequest{
		UserId: in.UserID,
		PetId:  petId,
	})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "pet").
			Str("module", "adopt").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return nil,
				&dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    constant.PetNotFoundMessage,
					Data:       nil,
				}
		case codes.Unavailable:
			return nil,
				&dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    constant.UnavailableServiceMessage,
					Data:       nil,
				}
		default:
			return nil,
				&dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    constant.InternalErrorMessage,
					Data:       nil,
				}
		}
	}
	return &dto.AdoptByResponse{
		Success: res.Success,
	}, nil
}
