package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	mock "github.com/isd-sgcu/johnjud-gateway/src/mocks/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/rs/zerolog/log"
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
			ImageUrls:    []string{},
			Background:   faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
		}

		pets = append(pets, pet)
	}

	t.Pet = t.Pets[0]

	t.PetDto = &dto.PetDto{
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
}

func (t *PetHandlerTest) TestFindOnePet() {
	want := t.Pet

	srv := new(mock.ServiceMock)
	imageSrv := new(mock.ServiceMock)
	srv.On("FindOne", t.Pet.Id).Return(want, nil)

	c := &mock.ContextMock{}
	c.On("ID").Return(t.Pet.Id, nil)

	v, err := validator.NewValidator()
	if err != nil {
		log.Error().Err(err).
			Str("handler", "pet").
			Msg("Err creating validator")
	}

	h := NewHandler(srv, v)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}
