package image

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	FindAll() ([]*dto.ImageResponse, *dto.ResponseErr)
	FindByPetId(string) ([]*dto.ImageResponse, *dto.ResponseErr)
	Upload(*dto.UploadImageRequest) (*dto.ImageResponse, *dto.ResponseErr)
	Delete(string) (*dto.DeleteImageResponse, *dto.ResponseErr)
	DeleteByPetId(string) (*dto.DeleteImageResponse, *dto.ResponseErr)
	AssignPet(*dto.AssignPetRequest) (*dto.AssignPetResponse, *dto.ResponseErr)
}

type serviceImpl struct {
	client proto.ImageServiceClient
}

func NewService(client proto.ImageServiceClient) Service {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) FindAll() ([]*dto.ImageResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindAll(ctx, &proto.FindAllImageRequest{})
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Str("service", "image").
			Str("module", "find all").
			Msg(st.Message())
		switch st.Code() {
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
	return ProtoToDtoList(res.Images), nil
}

func (s *serviceImpl) FindByPetId(petId string) ([]*dto.ImageResponse, *dto.ResponseErr) {
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
	return ProtoToDtoList(res.Images), nil
}

func (s *serviceImpl) Upload(in *dto.UploadImageRequest) (*dto.ImageResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := CreateDtoToProto(in)
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
	return ProtoToDto(res.Image), nil
}

func (s *serviceImpl) Delete(id string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
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

func (s *serviceImpl) DeleteByPetId(petId string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &proto.DeleteByPetIdRequest{
		PetId: petId,
	}

	res, errRes := s.client.DeleteByPetId(ctx, request)
	if errRes != nil {
		st, _ := status.FromError(errRes)
		log.Error().
			Err(errRes).
			Str("service", "image").
			Str("module", "delete by pet id").
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

func (s *serviceImpl) AssignPet(in *dto.AssignPetRequest) (*dto.AssignPetResponse, *dto.ResponseErr) {
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
		case codes.InvalidArgument:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    constant.InvalidArgumentMessage,
				Data:       nil,
			}
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
