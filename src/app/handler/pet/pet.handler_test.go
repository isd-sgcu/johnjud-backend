package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	mock_router "github.com/isd-sgcu/johnjud-gateway/src/mocks/router"
	mock_image "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/image"
	mock_pet "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/pet"
	mock_validator "github.com/isd-sgcu/johnjud-gateway/src/mocks/validator"

	pet_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/suite"
)

type PetHandlerTest struct {
	suite.Suite
	Pet                  *pet_proto.Pet
	Pets                 []*pet_proto.Pet
	PetDto               *dto.PetDto
	CreatePetRequest     *dto.CreatePetRequest
	ChangeViewPetRequest *dto.ChangeViewPetRequest
	UpdatePetRequest     *dto.UpdatePetRequest
	BindErr              *dto.ResponseErr
	NotFoundErr          *dto.ResponseErr
	ServiceDownErr       *dto.ResponseErr
	InternalErr          *dto.ResponseErr
}

func TestPetHandler(t *testing.T) {
	suite.Run(t, new(PetHandlerTest))
}

func (t *PetHandlerTest) SetupTest() {
	var pets []*pet_proto.Pet
	for i := 0; i <= 3; i++ {
		pet := &pet_proto.Pet{
			Id:           faker.UUIDDigit(),
			Type:         faker.Word(),
			Species:      faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       pet_proto.Gender(rand.Intn(1) + 1),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       pet_proto.PetStatus(rand.Intn(1) + 1),
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			IsClubPet:    true,
			ImageUrls:    []string{""},
			Background:   faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
		}

		pets = append(pets, pet)
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetDto = &dto.PetDto{
		Id:           t.Pet.Id,
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       pet.Gender(t.Pet.Gender),
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       pet.Status(t.Pet.Status),
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}

	t.CreatePetRequest = &dto.CreatePetRequest{
		Pet: &dto.PetDto{},
	}

	t.UpdatePetRequest = &dto.UpdatePetRequest{
		Pet: &dto.PetDto{},
	}

	t.ChangeViewPetRequest = &dto.ChangeViewPetRequest{}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Pet not found",
		Data:       nil,
	}

	t.BindErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid ID",
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Data:       nil,
	}
}

func (t *PetHandlerTest) TestFindAllSuccess() {
	findAllResponse := t.Pets

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	petSvc.EXPECT().FindAll().Return(findAllResponse, nil)
	context.EXPECT().JSON(http.StatusOK, findAllResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.FindAll(context)
}

func (t *PetHandlerTest) TestFindOneSuccess() {
	findOneResponse := t.Pet

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(findOneResponse, nil)
	context.EXPECT().JSON(http.StatusOK, findOneResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestFindOneNotFoundErr() {
	findOneResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(nil, findOneResponse)
	context.EXPECT().JSON(http.StatusNotFound, findOneResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestFindOneGrpcErr() {
	findOneResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(nil, findOneResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, findOneResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestCreateSuccess() {
	createErrorResponse := t.Pet

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.CreatePetRequest).Return(nil)
	petSvc.EXPECT().Create(t.CreatePetRequest).Return(createErrorResponse, nil)
	context.EXPECT().JSON(http.StatusCreated, createErrorResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Create(context)
}

func (t *PetHandlerTest) TestCreateGrpcErr() {
	createErrorResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.CreatePetRequest).Return(nil)
	petSvc.EXPECT().Create(t.CreatePetRequest).Return(nil, createErrorResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, createErrorResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Create(context)
}

func (t *PetHandlerTest) TestUpdateSuccess() {
	updateResponse := t.Pet

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(updateResponse, nil)
	context.EXPECT().JSON(http.StatusOK, updateResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestUpdateNotFound() {
	updateResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(nil, updateResponse)
	context.EXPECT().JSON(http.StatusNotFound, updateResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestUpdateGrpcErr() {
	updateResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(nil, updateResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, updateResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestDeleteSuccess() {
	deleteResponse := true

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, nil)
	context.EXPECT().JSON(http.StatusOK, deleteResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}
func (t *PetHandlerTest) TestDeleteNotFound() {
	deleteResponse := false

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestDeleteGrpcErr() {
	deleteResponse := false

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestChangeViewSuccess() {
	changeViewResponse := true

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(changeViewResponse, nil)
	context.EXPECT().JSON(http.StatusOK, changeViewResponse)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestChangeViewNotFound() {
	changeViewResponse := false

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(changeViewResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestChangeViewGrpcErr() {
	changeViewResponse := false

	controller := gomock.NewController(t.T())

	petSvc := mock_pet.NewMockService(controller)
	imageSvc := mock_image.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(changeViewResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}
