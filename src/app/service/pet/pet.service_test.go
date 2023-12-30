package pet

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	mock "github.com/isd-sgcu/johnjud-gateway/src/mocks/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	image_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PetServiceTest struct {
	suite.Suite
	Pets           []*proto.Pet
	Pet            *proto.Pet
	PetReq         *proto.Pet
	UpdatePetReq   *proto.UpdatePetRequest
	PetDto         *dto.PetDto
	UpdatePetDto   *dto.UpdatePetDto
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InternalErr    *dto.ResponseErr

	Images        []*image_proto.Image
	ImageUrls     []string
	ImagesList    [][]*image_proto.Image
	ImageUrlsList [][]string
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {
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
			ImageUrls:    []string{},
			Background:   faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
		}
		// ? wait for image mock
		// var images []*image_proto.Image
		// var imageUrls []string
		// for i := 0; i < 3; i++ {
		// 	url := faker.URL()
		// 	images = append(images, &image_proto.Image{
		// 		Id:       faker.UUIDDigit(),
		// 		PetId:    pet.Id,
		// 		ImageUrl: url,
		// 	})
		// 	imageUrls = append(imageUrls, url)
		// }
		// t.ImagesList = append(t.ImagesList, images)
		// t.ImageUrlsList = append(t.ImageUrlsList, imageUrls)
		pets = append(pets, pet)
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetReq = &proto.Pet{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		ImageUrls:    t.Pet.ImageUrls,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}

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

	t.UpdatePetDto = &dto.UpdatePetDto{
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

	t.UpdatePetReq = &proto.UpdatePetRequest{
		Pet: &proto.Pet{
			Id:           t.Pet.Id,
			Type:         t.Pet.Type,
			Species:      t.Pet.Species,
			Name:         t.Pet.Name,
			Birthdate:    t.Pet.Birthdate,
			Gender:       t.Pet.Gender,
			Habit:        t.Pet.Habit,
			Caption:      t.Pet.Caption,
			Status:       t.Pet.Status,
			ImageUrls:    t.Pet.ImageUrls,
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
		Message:    constant.ServiceDownMessage,
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    constant.PetNotFoundMessage,
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}
}

func (t *PetServiceTest) TestFindAllSuccess() {}

func (t *PetServiceTest) TestFindOneSuccess() {
	want := t.Pet

	c := &mock.ClientMock{}
	c.On("FindOne", &proto.FindOnePetRequest{Id: t.Pet.Id}).Return(&proto.FindOnePetResponse{Pet: want}, nil)

	service := NewService(c)

	actual, err := service.FindOne(t.Pet.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestFindOneNotFound() {
	want := t.NotFoundErr

	c := &mock.ClientMock{}
	c.On("FindOne", &proto.FindOnePetRequest{Id: t.Pet.Id}).Return(nil, status.Error(codes.NotFound, "Pet not found"))

	service := NewService(c)

	actual, err := service.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *PetServiceTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ClientMock{}
	c.On("FindOne", &proto.FindOnePetRequest{Id: t.Pet.Id}).Return(nil, errors.New(constant.ServiceDownMessage))

	service := NewService(c)

	actual, err := service.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *PetServiceTest) TestCreateSuccess() {
	want := t.Pet

	c := &mock.ClientMock{}
	c.On("Create", t.PetReq).Return(&proto.CreatePetResponse{Pet: want}, nil)

	service := NewService(c)
	actual, err := service.Create(&dto.CreatePetDto{Pet: t.PetDto})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestCreateGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ClientMock{}
	c.On("Create", t.PetReq).Return(nil, errors.New("Service is down"))

	service := NewService(c)

	actual, err := service.Create(&dto.CreatePetDto{Pet: t.PetDto})

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *PetServiceTest) TestUpdateSuccess() {
	want := t.Pet

	c := &mock.ClientMock{}
	c.On("Update", t.UpdatePetReq).Return(&proto.UpdatePetResponse{Pet: want}, nil)

	service := NewService(c)

	actual, err := service.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestUpdateNotFound() {
	want := t.NotFoundErr

	c := &mock.ClientMock{}
	c.On("Update", t.UpdatePetReq).Return(nil, status.Error(codes.NotFound, "Pet not found"))

	service := NewService(c)

	actual, err := service.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *PetServiceTest) TestUpdateGrpcErr() {
	want := t.ServiceDownErr

	c := &mock.ClientMock{}
	c.On("Update", t.UpdatePetReq).Return(nil, errors.New("Service iis down"))

	service := NewService(c)

	actual, err := service.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *PetServiceTest) TestDeleteSuccess() {
	c := &mock.ClientMock{}
	c.On("Delete", &proto.DeletePetRequest{Id: t.Pet.Id}).Return(&proto.DeletePetResponse{Success: true}, nil)

	service := NewService(c)

	actual, err := service.Delete(t.Pet.Id)

	assert.Nil(t.T(), err)
	assert.True(t.T(), actual)
}

func (t *PetServiceTest) TestDeleteNotFound() {
	want := t.NotFoundErr

	c := &mock.ClientMock{}
	c.On("Delete", &proto.DeletePetRequest{Id: t.Pet.Id}).Return(nil, status.Error(codes.NotFound, "Pet not found"))

	service := NewService(c)

	actual, err := service.Delete(t.Pet.Id)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *PetServiceTest) TestDeleteGrpcErr() {
	want := t.NotFoundErr

	c := &mock.ClientMock{}
	c.On("Delete", &proto.DeletePetRequest{Id: t.Pet.Id}).Return(nil, status.Error(codes.NotFound, "Pet not found"))

	service := NewService(c)

	actual, err := service.Delete(t.Pet.Id)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), want, err)
}
