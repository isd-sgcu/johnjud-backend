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
	client proto.PetServiceClient
}

func NewService(client proto.PetServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindAll(req *dto.FindOnePetDto) (result *proto.Pet, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindOne(ctx, &proto.FindOnePetRequest{Id: request.Id})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "pet").
					Str("module", "find one").
					Str("pet_id", req.Id).
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
					Str("pet_id", req.Id).
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
					Str("pet_id", req.Id).
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
			Str("per_id", req.Id).
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
		Str("pet_id", req.Id).
		Msg("Find pet success")
	return res.Pet, nil
}
