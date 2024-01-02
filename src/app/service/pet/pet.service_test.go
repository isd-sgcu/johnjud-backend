package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	petMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/client/pet"
	petProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PetServiceTest struct {
	suite.Suite
	Pets                  []*petProto.Pet
	Pet                   *petProto.Pet
	PetReq                *petProto.Pet
	PetNotVisible         *petProto.Pet
	UpdatePetReq          *petProto.UpdatePetRequest
	ChangeViewPetReq      *petProto.ChangeViewPetRequest
	PetDto                *dto.PetDto
	CreatePetDto          *dto.CreatePetDto
	UpdatePetDto          *dto.UpdatePetDto
	NotFoundErr           *dto.ResponseErr
	UnavailableServiceErr *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
	InternalErr           *dto.ResponseErr
	ChangeViewedPetDto    *dto.ChangeViewPetDto

	Images        []*imageProto.Image
	ImageUrls     []string
	ImagesList    [][]*imageProto.Image
	ImageUrlsList [][]string
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {
	var pets []*petProto.Pet
	for i := 0; i <= 3; i++ {
		pet := &petProto.Pet{
			Id:           faker.UUIDDigit(),
			Type:         faker.Word(),
			Species:      faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       petProto.Gender(rand.Intn(1) + 1),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       petProto.PetStatus(rand.Intn(1) + 1),
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

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetReq = &petProto.Pet{
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

	t.PetNotVisible = &petProto.Pet{
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
		IsVisible:    false,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}

	t.PetDto = RawToDto(t.Pet)

	t.CreatePetDto = &dto.CreatePetDto{
		Pet: RawToDto(t.PetReq),
	}

	t.UpdatePetDto = &dto.UpdatePetDto{
		Pet: RawToDto(t.Pet),
	}

	t.UpdatePetReq = &petProto.UpdatePetRequest{
		Pet: t.Pet,
	}

	t.ChangeViewedPetDto = &dto.ChangeViewPetDto{
		Visible: false,
	}

	t.ChangeViewPetReq = &petProto.ChangeViewPetRequest{
		Id:      t.Pet.Id,
		Visible: false,
	}

	t.UnavailableServiceErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    constant.PetNotFoundMessage,
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	t.InvalidArgumentErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidArgument,
		Data:       nil,
	}
}

func (t *PetServiceTest) TestFindAllSuccess() {
	protoReq := &petProto.FindAllPetRequest{}
	protoResp := &petProto.FindAllPetResponse{
		Pets: t.Pets,
	}

	expected := t.Pets

	client := petMock.PetClientMock{}
	client.On("FindAll", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindAll()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestFindAllUnavailableServiceError() {
	protoReq := &petProto.FindAllPetRequest{}

	expected := t.UnavailableServiceErr

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	client := petMock.PetClientMock{}
	client.On("FindAll", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindAll()

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestFindOneSuccess() {
	protoReq := &petProto.FindOnePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petProto.FindOnePetResponse{
		Pet: t.Pet,
	}

	expected := t.Pet

	client := petMock.PetClientMock{}
	client.On("FindOne", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestFindOneNotFoundError() {
	protoReq := &petProto.FindOnePetRequest{
		Id: t.Pet.Id,
	}

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := petMock.PetClientMock{}
	client.On("FindOne", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestFindOneUnavailableServiceError() {
	protoReq := &petProto.FindOnePetRequest{
		Id: t.Pet.Id,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := petMock.PetClientMock{}
	client.On("FindOne", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateSuccess() {
	protoReq := &petProto.CreatePetRequest{
		Pet: t.PetReq,
	}
	protoResp := &petProto.CreatePetResponse{
		Pet: t.Pet,
	}

	expected := t.Pet

	client := &petMock.PetClientMock{}
	client.On("Create", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestCreateInvalidArgumentError() {
	protoReq := &petProto.CreatePetRequest{
		Pet: t.PetReq,
	}

	expected := t.InvalidArgumentErr

	clientErr := status.Error(codes.InvalidArgument, constant.InvalidArgument)

	client := &petMock.PetClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateInternalError() {
	protoReq := &petProto.CreatePetRequest{
		Pet: t.PetReq,
	}

	expected := t.InternalErr

	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)

	client := &petMock.PetClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateUnavailableServiceError() {
	protoReq := &petProto.CreatePetRequest{
		Pet: t.PetReq,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petMock.PetClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestUpdateSuccess() {
	protoReq := t.UpdatePetReq
	protoResp := &petProto.UpdatePetResponse{
		Pet: t.Pet,
	}

	expected := t.Pet

	client := &petMock.PetClientMock{}
	client.On("Update", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestUpdateNotFound() {
	protoReq := t.UpdatePetReq
	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petMock.PetClientMock{}
	client.On("Update", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestUpdateUnavailableServiceError() {
	protoReq := t.UpdatePetReq
	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petMock.PetClientMock{}
	client.On("Update", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestDeleteSuccess() {
	protoReq := &petProto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petProto.DeletePetResponse{
		Success: true,
	}

	expected := true

	client := &petMock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Delete(t.Pet.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestDeleteNotFound() {
	protoReq := &petProto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petProto.DeletePetResponse{
		Success: false,
	}
	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petMock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Pet.Id)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestDeleteServiceUnavailableError() {
	protoReq := &petProto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petProto.DeletePetResponse{
		Success: false,
	}
	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petMock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Pet.Id)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestChangeViewSuccess() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petProto.ChangeViewPetResponse{
		Success: true,
	}

	client := &petMock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.Nil(t.T(), err)
	assert.True(t.T(), actual)
}

func (t *PetServiceTest) TestChangeViewNotFoundError() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petProto.ChangeViewPetResponse{
		Success: false,
	}

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petMock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestChangeViewUnavailableServiceError() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petProto.ChangeViewPetResponse{
		Success: false,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petMock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.False(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}
