package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	routerMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/router"
	imageMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/image"
	likeMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/like"
	petMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/pet"
	validatorMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/validator"

	errConst "github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	likeUtils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/like"
	petUtils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/pet"
	petConst "github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	userProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	likeProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
	petProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imgProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"

	"github.com/stretchr/testify/suite"
)

type PetHandlerTest struct {
	suite.Suite
	Pet                  *petProto.Pet
	User                 *userProto.User
	Likes                []*likeProto.Like
	Pets                 []*petProto.Pet
	PetDto               *dto.PetResponse
	CreatePetRequest     *dto.CreatePetRequest
	ChangeViewPetRequest *dto.ChangeViewPetRequest
	UpdatePetRequest     *dto.UpdatePetRequest
	AdoptByRequest       *dto.AdoptByRequest
	BindErr              *dto.ResponseErr
	NotFoundErr          *dto.ResponseErr
	ServiceDownErr       *dto.ResponseErr
	InternalErr          *dto.ResponseErr
	Images               []*imgProto.Image
	ImagesList           [][]*imgProto.Image
}

func TestPetHandler(t *testing.T) {
	suite.Run(t, new(PetHandlerTest))
}

func (t *PetHandlerTest) SetupTest() {
	imagesList := petUtils.MockImageList(3)
	t.ImagesList = imagesList
	t.Images = imagesList[0]
	var pets []*petProto.Pet
	genders := []petConst.Gender{petConst.MALE, petConst.FEMALE}
	statuses := []petConst.Status{petConst.ADOPTED, petConst.FINDHOME}

	for i := 0; i <= 3; i++ {
		pet := &petProto.Pet{
			Id:           faker.UUIDDigit(),
			Type:         faker.Word(),
			Species:      faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       string(genders[rand.Intn(2)]),
			Color:        faker.Word(),
			Pattern:      faker.Word(),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Images:       []*imgProto.Image{},
			Status:       string(statuses[rand.Intn(2)]),
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			IsClubPet:    true,
			Origin:       faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
			AdoptBy:      "",
		}

		pets = append(pets, pet)
	}

	var likes []*likeProto.Like
	for i := 0; i <= 3; i++ {
		like := &likeProto.Like{
			Id:     faker.UUIDDigit(),
			UserId: faker.UUIDDigit(),
			PetId:  faker.UUIDDigit(),
		}
		likes = append(likes, like)
	}
	t.Likes = likes

	t.User = &userProto.User{
		Id:        faker.UUIDDigit(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
		Role:      "USER",
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetDto = &dto.PetResponse{
		Id:           t.Pet.Id,
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       pet.Gender(t.Pet.Gender),
		Color:        t.Pet.Color,
		Pattern:      t.Pet.Pattern,
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
		AdoptBy:      t.Pet.AdoptBy,
		IsLike:       false,
	}

	t.CreatePetRequest = &dto.CreatePetRequest{}

	t.UpdatePetRequest = &dto.UpdatePetRequest{}

	t.ChangeViewPetRequest = &dto.ChangeViewPetRequest{}

	t.AdoptByRequest = &dto.AdoptByRequest{}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    errConst.PetNotFoundMessage,
		Data:       nil,
	}

	t.BindErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    errConst.InvalidIDMessage,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    errConst.InternalErrorMessage,
		Data:       nil,
	}
}

func (t *PetHandlerTest) TestFindAllSuccess() {
	findAllResponse := petUtils.ProtoToDtoList(t.Pets, t.ImagesList)
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petConst.FindAllPetSuccessMessage,
		Data:       findAllResponse,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().IsAuth().Return(true)
	context.EXPECT().UserID().Return(t.User.Id)
	petSvc.EXPECT().FindAll().Return(findAllResponse, nil)
	likeSvc.EXPECT().FindByUserId(t.User.Id).Return(likeUtils.ProtoToDtoList(t.Likes), nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.FindAll(context)
}

func (t *PetHandlerTest) TestFindOneSuccess() {
	findOneResponse := petUtils.ProtoToDto(t.Pet, t.Images)
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petConst.FindOnePetSuccessMessage,
		Data:       findOneResponse,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().IsAuth().Return(true)
	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().UserID().Return(t.User.Id)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(findOneResponse, nil)
	likeSvc.EXPECT().FindByUserId(t.User.Id).Return(likeUtils.ProtoToDtoList(t.Likes), nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestFindOneNotFoundErr() {
	findOneResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().IsAuth().Return(true)
	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(nil, findOneResponse)
	context.EXPECT().JSON(http.StatusNotFound, findOneResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestFindOneGrpcErr() {
	findOneResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().IsAuth().Return(true)
	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(nil, findOneResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, findOneResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestCreateSuccess() {
	createResponse := petUtils.ProtoToDto(t.Pet, t.Images)
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusCreated,
		Message:    petConst.CreatePetSuccessMessage,
		Data:       createResponse,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.CreatePetRequest).Return(nil)
	petSvc.EXPECT().Create(t.CreatePetRequest).Return(createResponse, nil)
	context.EXPECT().JSON(http.StatusCreated, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Create(context)
}

func (t *PetHandlerTest) TestCreateGrpcErr() {
	createErrorResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.CreatePetRequest).Return(nil)
	petSvc.EXPECT().Create(t.CreatePetRequest).Return(nil, createErrorResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, createErrorResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Create(context)
}

func (t *PetHandlerTest) TestUpdateSuccess() {
	updateResponse := petUtils.ProtoToDto(t.Pet, t.Images)
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petConst.UpdatePetSuccessMessage,
		Data:       updateResponse,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(updateResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestUpdateNotFound() {
	updateResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(nil, updateResponse)
	context.EXPECT().JSON(http.StatusNotFound, updateResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestUpdateGrpcErr() {
	updateResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(nil, updateResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, updateResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestDeleteSuccess() {
	deleteResponse := &dto.DeleteResponse{
		Success: true,
	}
	expectedResponse := deleteResponse

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Delete(context)
}
func (t *PetHandlerTest) TestDeleteNotFound() {
	deleteResponse := &dto.DeleteResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestDeleteGrpcErr() {
	deleteResponse := &dto.DeleteResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestChangeViewSuccess() {
	changeViewResponse := &dto.ChangeViewPetResponse{
		Success: true,
	}
	expectedResponse := changeViewResponse

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.ChangeViewPetRequest).Return(nil)
	validator.EXPECT().Validate(t.ChangeViewPetRequest).Return(nil)
	petSvc.EXPECT().ChangeView(t.Pet.Id, t.ChangeViewPetRequest).Return(changeViewResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.ChangeView(context)
}

func (t *PetHandlerTest) TestChangeViewNotFound() {
	changeViewResponse := &dto.ChangeViewPetResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.ChangeViewPetRequest).Return(nil)
	validator.EXPECT().Validate(t.ChangeViewPetRequest).Return(nil)
	petSvc.EXPECT().ChangeView(t.Pet.Id, t.ChangeViewPetRequest).Return(changeViewResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.ChangeView(context)
}

func (t *PetHandlerTest) TestChangeViewGrpcErr() {
	changeViewResponse := &dto.ChangeViewPetResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.ChangeViewPetRequest).Return(nil)
	validator.EXPECT().Validate(t.ChangeViewPetRequest).Return(nil)
	petSvc.EXPECT().ChangeView(t.Pet.Id, t.ChangeViewPetRequest).Return(changeViewResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.ChangeView(context)
}

func (t *PetHandlerTest) TestAdoptSuccess() {
	adoptByResponse := &dto.AdoptByResponse{
		Success: true,
	}
	expectedResponse := adoptByResponse

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.AdoptByRequest).Return(nil)
	validator.EXPECT().Validate(t.AdoptByRequest).Return(nil)
	petSvc.EXPECT().Adopt(t.Pet.Id, t.AdoptByRequest).Return(adoptByResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Adopt(context)
}

func (t *PetHandlerTest) TestAdoptNotFound() {
	adoptByResponse := &dto.AdoptByResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.AdoptByRequest).Return(nil)
	validator.EXPECT().Validate(t.AdoptByRequest).Return(nil)
	petSvc.EXPECT().Adopt(t.Pet.Id, t.AdoptByRequest).Return(adoptByResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Adopt(context)
}

func (t *PetHandlerTest) TestAdoptGrpcErr() {
	adoptByResponse := &dto.AdoptByResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.AdoptByRequest).Return(nil)
	validator.EXPECT().Validate(t.AdoptByRequest).Return(nil)
	petSvc.EXPECT().Adopt(t.Pet.Id, t.AdoptByRequest).Return(adoptByResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(petSvc, imageSvc, likeSvc, validator)
	handler.Adopt(context)
}
