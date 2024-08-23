package image

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/client/bucket"
	"github.com/isd-sgcu/johnjud-backend/constant"
	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"github.com/isd-sgcu/johnjud-backend/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	FindAll() ([]*dto.ImageResponse, *dto.ResponseErr)
	FindByPetId(petID string) ([]*dto.ImageResponse, *dto.ResponseErr)
	Upload(request *dto.UploadImageRequest) (*dto.ImageResponse, *dto.ResponseErr)
	Delete(id string) (*dto.DeleteImageResponse, *dto.ResponseErr)
	DeleteByPetId(petID string) (*dto.DeleteImageResponse, *dto.ResponseErr)
	AssignPet(request *dto.AssignPetRequest) (*dto.AssignPetResponse, *dto.ResponseErr)
}

type serviceImpl struct {
	client     bucket.Client
	repository Repository
	random     utils.RandomUtil
}

func NewService(client bucket.Client, repository Repository, random utils.RandomUtil) Service {
	return &serviceImpl{
		client:     client,
		repository: repository,
		random:     random,
	}
}

func (s *serviceImpl) FindAll() ([]*dto.ImageResponse, *dto.ResponseErr) {
	var images []*model.Image

	err := s.repository.FindAll(&images)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "find all").
			Msg("Error finding all images")

		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return RawToDtoList(&images), nil
}

func (s *serviceImpl) FindByPetId(petID string) ([]*dto.ImageResponse, *dto.ResponseErr) {
	var images []*model.Image

	err := s.repository.FindByPetId(petID, &images)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "find by petId").
			Str("petId", petID).
			Msg("Error finding image by pet id from repo")
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NotFoundError(constant.ImageNotFoundErrorMessage)
		}

		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return RawToDtoList(&images), nil
}

func (s *serviceImpl) Upload(req *dto.UploadImageRequest) (*dto.ImageResponse, *dto.ResponseErr) {
	if req.PetId != "" {
		_, err := uuid.Parse(req.PetId)
		if err != nil {
			log.Error().Err(err).
				Str("service", "image").
				Str("module", "upload").
				Str("petId", req.PetId).
				Msg(constant.PetIdNotUUIDErrorMessage)

			return nil, dto.BadRequestError(constant.PetIdNotUUIDErrorMessage)
		}
	}

	randomString, err := s.random.GenerateRandomString(10)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "upload").
			Str("petId", req.PetId).
			Msg("Error while generating random string")
		return nil, dto.InternalServerError("Error while generating random string")
	}

	imageUrl, objectKey, err := s.client.Upload(req.File, randomString+"_"+req.Filename)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "upload").
			Str("petId", req.PetId).
			Msg(constant.UploadToBucketErrorMessage)

		// return nil, status.Error(codes.Internal, constant.UploadToBucketErrorMessage)
		return nil, dto.InternalServerError(constant.UploadToBucketErrorMessage)
	}

	raw, _ := DtoToRaw(&dto.ImageResponse{
		PetId:     req.PetId,
		Url:       imageUrl,
		ObjectKey: objectKey,
	})

	err = s.repository.Create(raw)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "upload").
			Str("petId", req.PetId).
			Msg(constant.CreateImageErrorMessage)

		return nil, dto.InternalServerError(constant.CreateImageErrorMessage)
	}

	return RawToDto(raw), nil
}

func (s *serviceImpl) AssignPet(req *dto.AssignPetRequest) (*dto.AssignPetResponse, *dto.ResponseErr) {
	petId, err := uuid.Parse(req.PetId)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "assign pet").
			Str("petId", req.PetId).
			Msg(constant.PrimaryKeyRequiredErrorMessage)

		return nil, dto.BadRequestError(constant.PrimaryKeyRequiredErrorMessage)
	}

	for _, id := range req.Ids {
		err = s.repository.Update(id, &model.Image{
			PetID: &petId,
		})
		if err == nil {
			continue
		}

		log.Error().Err(err).
			Str("service", "image").
			Str("module", "assign pet").
			Str("petId", req.PetId).
			Msg("Error updating image in repo")

		if strings.Contains(err.Error(), gorm.ErrForeignKeyViolated.Error()) {
			return nil, dto.NotFoundError(constant.PetIdNotFoundErrorMessage)
		}
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, dto.NotFoundError(constant.ImageNotFoundErrorMessage)
		default:
			return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
		}
	}

	return &dto.AssignPetResponse{Success: true}, nil
}

func (s *serviceImpl) Delete(id string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
	var image model.Image

	err := s.repository.FindOne(id, &image)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete").
			Str("id", id).
			Msg("Error finding image from repo")
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NotFoundError(constant.ImageNotFoundErrorMessage)
		}

		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	err = s.client.Delete(image.ObjectKey)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete").
			Str("id", id).
			Msg(constant.DeleteFromBucketErrorMessage)

		return nil, dto.InternalServerError(constant.DeleteFromBucketErrorMessage)
	}

	err = s.repository.Delete(id)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete").
			Str("id", id).
			Msg(constant.DeleteImageErrorMessage)

		return nil, dto.InternalServerError(constant.DeleteImageErrorMessage)
	}

	return &dto.DeleteImageResponse{Success: true}, nil
}

func (s *serviceImpl) DeleteByPetId(petID string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
	var images []*model.Image

	err := s.repository.FindByPetId(petID, &images)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete by pet id").
			Str("pet id", petID).
			Msg("Error finding image from repo")
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NotFoundError(constant.ImageNotFoundErrorMessage)
		}

		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	imageObjectKeys := ExtractImageObjectKeys(images)
	err = s.client.DeleteMany(imageObjectKeys)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete by pet id").
			Interface("image object keys", imageObjectKeys).
			Msg(constant.DeleteFromBucketErrorMessage)

		return nil, dto.InternalServerError(constant.DeleteFromBucketErrorMessage)
	}

	imageIds := ExtractImageIds(images)
	err = s.repository.DeleteMany(imageIds)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete by pet id").
			Interface("image ids", imageIds).
			Msg(constant.DeleteImageErrorMessage)

		return nil, dto.InternalServerError(constant.DeleteImageErrorMessage)
	}

	return &dto.DeleteImageResponse{Success: true}, nil
}

func DtoToRaw(in *dto.ImageResponse) (result *model.Image, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	petId, err := uuid.Parse(in.PetId)
	if err != nil {
		return &model.Image{
			Base: model.Base{
				ID:        id,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			PetID:     nil,
			ImageUrl:  in.Url,
			ObjectKey: in.ObjectKey,
		}, nil
	}

	return &model.Image{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		PetID:     &petId,
		ImageUrl:  in.Url,
		ObjectKey: in.ObjectKey,
	}, nil
}

func RawToDtoList(in *[]*model.Image) []*dto.ImageResponse {
	var result []*dto.ImageResponse
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *model.Image) *dto.ImageResponse {
	var id string
	var petId string
	if in.ID != uuid.Nil {
		id = in.ID.String()
	}
	if in.PetID != nil {
		petId = in.PetID.String()
	}

	return &dto.ImageResponse{
		Id:        id,
		PetId:     petId,
		Url:       in.ImageUrl,
		ObjectKey: in.ObjectKey,
	}
}

func ExtractImageIds(in []*model.Image) []string {
	var imageIds []string
	for _, image := range in {
		imageIds = append(imageIds, image.ID.String())
	}

	return imageIds
}

func ExtractImageObjectKeys(in []*model.Image) []string {
	var imageObjectKeys []string
	for _, image := range in {
		imageObjectKeys = append(imageObjectKeys, image.ObjectKey)
	}

	return imageObjectKeys
}
