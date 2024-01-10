package image

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	client proto.ImageServiceClient
}

func NewService(client proto.ImageServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindByPetId(petId string) ([]*dto.ImageResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindByPetId(ctx, &proto.FindImageByPetIdRequest{PetId: petId})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Str("service", "image").
			Str("module", "find by pet id").
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
	return utils.ProtoToDtoList(res.Images), nil
}

func (s *Service) Upload(in *dto.UploadImageRequest) (*dto.ImageResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := utils.CreateDtoToProto(in)
	res, errRes := s.client.Upload(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "image").
			Str("module", "upload").
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
	return utils.ProtoToDto(res.Image), nil
}

func (s *Service) Delete(id string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &proto.DeleteImageRequest{
		Id: id,
	}

	res, errRes := s.client.Delete(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "image").
			Str("module", "delete").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.ImageNotFoundMessage,
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
	return &dto.DeleteImageResponse{
		Success: res.Success,
	}, nil
}

func (s *Service) AssignPet(in *dto.AssignPetRequest) (*dto.AssignPetResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &proto.AssignPetRequest{
		Ids:   in.Ids,
		PetId: in.PetId,
	}

	res, errRes := s.client.AssignPet(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "image").
			Str("module", "assign pet").
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
	return &dto.AssignPetResponse{
		Success: res.Success,
	}, nil
}
