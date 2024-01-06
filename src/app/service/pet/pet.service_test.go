package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/pet"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	petmock "github.com/isd-sgcu/johnjud-gateway/src/mocks/client/pet"
	petproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imgproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PetServiceTest struct {
	suite.Suite
	Pets                  []*petproto.Pet
	Pet                   *petproto.Pet
	PetNotVisible         *petproto.Pet
	UpdatePetReq          *petproto.UpdatePetRequest
	CreatePetReq          *petproto.CreatePetRequest
	ChangeViewPetReq      *petproto.ChangeViewPetRequest
	DeletePetReq          *petproto.DeletePetRequest
	AdoptReq              *petproto.AdoptPetRequest
	PetDto                *dto.PetResponse
	CreatePetDto          *dto.CreatePetRequest
	UpdatePetDto          *dto.UpdatePetRequest
	NotFoundErr           *dto.ResponseErr
	UnavailableServiceErr *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
	InternalErr           *dto.ResponseErr
	ChangeViewedPetDto    *dto.ChangeViewPetRequest

	Images     []*imgproto.Image
	ImagesList [][]*imgproto.Image
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {
	imagesList := utils.MockImageList(3)
	t.ImagesList = imagesList
	t.Images = imagesList[0]

	var pets []*petproto.Pet
	for i := 0; i <= 3; i++ {
		pet := &petproto.Pet{
			Id:           faker.UUIDDigit(),
			Type:         faker.Word(),
			Species:      faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       petproto.Gender(rand.Intn(1) + 1),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Images:       imagesList[i],
			Status:       petproto.PetStatus(rand.Intn(1) + 1),
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			IsClubPet:    true,
			Background:   faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
			AdoptBy:      faker.UUIDDigit(),
		}

		pets = append(pets, pet)
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetNotVisible = &petproto.Pet{
		Id:           t.Pet.Id,
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Images:       t.Pet.Images,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    false,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
		AdoptBy:      t.Pet.AdoptBy,
	}

	t.PetDto = utils.ProtoToDto(t.Pet, t.Pet.Images)

	t.CreatePetDto = &dto.CreatePetRequest{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       pet.Gender(t.Pet.Gender),
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Images:       []string{},
		Status:       pet.Status(t.Pet.Status),
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		IsClubPet:    &t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
		AdoptBy:      t.Pet.AdoptBy,
	}

	t.UpdatePetDto = &dto.UpdatePetRequest{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       pet.Gender(t.Pet.Gender),
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Images:       []string{},
		Status:       pet.Status(t.Pet.Status),
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		IsClubPet:    &t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
		AdoptBy:      t.Pet.AdoptBy,
	}

	t.CreatePetReq = utils.CreateDtoToProto(t.CreatePetDto)
	t.UpdatePetReq = utils.UpdateDtoToProto(t.Pet.Id, t.UpdatePetDto)

	t.ChangeViewPetReq = &petproto.ChangeViewPetRequest{
		Id:      t.Pet.Id,
		Visible: false,
	}

	t.ChangeViewedPetDto = &dto.ChangeViewPetRequest{
		Visible: false,
	}

	t.AdoptReq = &petproto.AdoptPetRequest{
		PetId:  t.Pet.Id,
		UserId: t.Pet.AdoptBy,
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
		Message:    constant.InvalidArgumentMessage,
		Data:       nil,
	}
}

func (t *PetServiceTest) TestFindAllSuccess() {
	protoReq := &petproto.FindAllPetRequest{}
	protoResp := &petproto.FindAllPetResponse{
		Pets: t.Pets,
	}

	expected := utils.ProtoToDtoList(t.Pets, t.ImagesList)

	client := petmock.PetClientMock{}
	client.On("FindAll", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindAll()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestFindAllUnavailableServiceError() {
	protoReq := &petproto.FindAllPetRequest{}

	expected := t.UnavailableServiceErr

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	client := petmock.PetClientMock{}
	client.On("FindAll", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindAll()

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestFindOneSuccess() {
	protoReq := &petproto.FindOnePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petproto.FindOnePetResponse{
		Pet: t.Pet,
	}

	expected := utils.ProtoToDto(t.Pet, t.Pet.Images)

	client := petmock.PetClientMock{}
	client.On("FindOne", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestFindOneNotFoundError() {
	protoReq := &petproto.FindOnePetRequest{
		Id: t.Pet.Id,
	}

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := petmock.PetClientMock{}
	client.On("FindOne", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestFindOneUnavailableServiceError() {
	protoReq := &petproto.FindOnePetRequest{
		Id: t.Pet.Id,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := petmock.PetClientMock{}
	client.On("FindOne", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateSuccess() {
	protoReq := t.CreatePetReq
	protoResp := &petproto.CreatePetResponse{
		Pet: t.Pet,
	}

	expected := utils.ProtoToDto(t.Pet, t.Pet.Images)

	client := &petmock.PetClientMock{}
	client.On("Create", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestCreateInvalidArgumentError() {
	protoReq := t.CreatePetReq

	expected := t.InvalidArgumentErr

	clientErr := status.Error(codes.InvalidArgument, constant.InvalidArgumentMessage)

	client := &petmock.PetClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateInternalError() {
	protoReq := t.CreatePetReq

	expected := t.InternalErr

	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)

	client := &petmock.PetClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateUnavailableServiceError() {
	protoReq := t.CreatePetReq

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petmock.PetClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestUpdateSuccess() {
	protoReq := t.UpdatePetReq
	protoResp := &petproto.UpdatePetResponse{
		Pet: t.Pet,
	}

	expected := utils.ProtoToDto(t.Pet, t.Pet.Images)

	client := &petmock.PetClientMock{}
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

	client := &petmock.PetClientMock{}
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

	client := &petmock.PetClientMock{}
	client.On("Update", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestDeleteSuccess() {
	protoReq := &petproto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petproto.DeletePetResponse{
		Success: true,
	}

	expected := &dto.DeleteResponse{Success: true}

	client := &petmock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Delete(t.Pet.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestDeleteNotFound() {
	protoReq := &petproto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petproto.DeletePetResponse{
		Success: false,
	}
	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petmock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Pet.Id)

	assert.Equal(t.T(), &dto.DeleteResponse{Success: false}, actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestDeleteServiceUnavailableError() {
	protoReq := &petproto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petproto.DeletePetResponse{
		Success: false,
	}
	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petmock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Pet.Id)

	assert.Equal(t.T(), &dto.DeleteResponse{Success: false}, actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestChangeViewSuccess() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petproto.ChangeViewPetResponse{
		Success: true,
	}

	client := &petmock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), actual, &dto.ChangeViewPetResponse{Success: true})
}

func (t *PetServiceTest) TestChangeViewNotFoundError() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petproto.ChangeViewPetResponse{
		Success: false,
	}

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petmock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.Equal(t.T(), &dto.ChangeViewPetResponse{Success: false}, actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestChangeViewUnavailableServiceError() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petproto.ChangeViewPetResponse{
		Success: false,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petmock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, clientErr)

	svc := NewService(client)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.Equal(t.T(), &dto.ChangeViewPetResponse{Success: false}, actual)
	assert.Equal(t.T(), expected, err)
}
