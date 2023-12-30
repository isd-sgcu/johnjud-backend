package pet

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	imageMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/image"
	mock "github.com/isd-sgcu/johnjud-gateway/src/mocks/pet"
	"github.com/rs/zerolog/log"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PetHandlerTest struct {
	suite.Suite
	Pet            *proto.Pet
	Pets           []*proto.Pet
	PetDto         *dto.PetDto
	UpdatedPetReq  *dto.UpDatePetDto
	BindErr        *dto.ResponseErr
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestPetHandler(t *testing.T) {
	suite.Run(t, new(PetHandlerTest))
}

func (t *PetHandlerTest) SetupTest() {
	var pets []*proto.Pet
	for i := 0; i <= 3; i++ {
		pet := &proto.Pet{
			Id:           faker.UUIDDigit(),
			Type:         faker.Word(),
			Species:      faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       proto.Gender(rand.Intn(1) + 1),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       proto.PetStatus(rand.Intn(1) + 1),
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

	t.UpdatedPetReq = &dto.UpDatePetDto{
		Pet: &dto.PetDto{
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
		},
	}

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

func (t *PetHandlerTest) TestFindOneSuccess() {
	want := t.Pet

	petService := &mock.ServiceMock{}
	imageService := &imageMock.ServiceMock{}

	petService.On("FindOne", t.Pet.Id).Return(want, nil)
	imageService.On("FindByPetId", t.Pet.Id).Return(t.Pet.ImageUrls, nil)

	context := &mock.ContextMock{}
	context.On("ID").Return(t.Pet.Id, nil)

	validator, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
	}

	h := NewHandler(petService, imageService, validator)
	h.FindOne(context)

	assert.Equal(t.T(), want, context.V)
	assert.Equal(t.T(), http.StatusOK, context.StatusCode)
}

func (t *PetHandlerTest) TestFindOneNotFoundErr() {
	want := t.NotFoundErr

	petService := &mock.ServiceMock{}
	imageService := &imageMock.ServiceMock{}

	petService.On("FindOne", t.Pet.Id).Return(nil, t.NotFoundErr)
	imageService.On("FindByPetId", t.Pet.Id).Return(nil, t.NotFoundErr)

	validator, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
		return
	}

	context := &mock.ContextMock{}
	context.On("ID").Return(t.Pet.Id, nil)

	h := NewHandler(petService, imageService, validator)
	h.FindOne(context)

	assert.Equal(t.T(), want, context.V)
	assert.Equal(t.T(), http.StatusNotFound, context.StatusCode)
}

func (t *PetHandlerTest) TestFindOneInternalErr() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Invalid ID",
		Data:       nil,
	}

	petService := &mock.ServiceMock{}
	imageService := &imageMock.ServiceMock{}

	petService.On("FindOne", t.Pet.Id).Return(nil, t.ServiceDownErr)
	imageService.On("FindByPetId", t.Pet.Id).Return(nil, t.NotFoundErr)

	c := &mock.ContextMock{}
	c.On("ID").Return("", errors.New("Cannot parse id"))

	validator, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
		return
	}

	h := NewHandler(petService, imageService, validator)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusInternalServerError, c.StatusCode)
}

func (t *PetHandlerTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	petService := &mock.ServiceMock{}
	imageService := &imageMock.ServiceMock{}

	imageService.On("FindByPetId", t.Pet.Id).Return(nil, t.NotFoundErr)
	petService.On("FindOne", t.Pet.Id).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{}
	c.On("ID").Return(t.Pet.Id, nil)

	validator, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
		return
	}

	h := NewHandler(petService, imageService, validator)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.StatusCode)
}

func (t *PetHandlerTest) TestCreateSuccess() {
	want := t.Pet

	petService := &mock.ServiceMock{}
	imageService := &imageMock.ServiceMock{}

	petService.On("Create", t.PetDto).Return(want, nil)

	c := &mock.ContextMock{}
	c.On("Bind", &dto.PetDto{}).Return(t.PetDto, nil)

	validator, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
		return
	}

	h := NewHandler(petService, imageService, validator)
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusCreated, c.StatusCode)
}

func (t *PetHandlerTest) TestCreateGrpcErr() {
	want := t.ServiceDownErr

	petService := &mock.ServiceMock{}
	imageService := &imageMock.ServiceMock{}

	imageService.On("FindByPetId", t.Pet.Id).Return(nil, t.ServiceDownErr)
	petService.On("FindOne", t.Pet.Id).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{}
	c.On("Bind", &dto.PetDto{}).Return(t.PetDto, nil)

	validator, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
		return
	}

	h := NewHandler(petService, imageService, validator)
	// TODO: fix this
	h.Create(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.StatusCode)
}
