package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/isd-sgcu/johnjud-gateway/internal/pet"
	routerMock "github.com/isd-sgcu/johnjud-gateway/mocks/router"
	imageMock "github.com/isd-sgcu/johnjud-gateway/mocks/service/image"
	petMock "github.com/isd-sgcu/johnjud-gateway/mocks/service/pet"
	validatorMock "github.com/isd-sgcu/johnjud-gateway/mocks/validator"

	utils "github.com/isd-sgcu/johnjud-gateway/internal/utils/pet"
	petProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imgProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"

	"github.com/stretchr/testify/suite"
)

type PetHandlerTest struct {
	suite.Suite
	Pet                  *petProto.Pet
	Pets                 []*petProto.Pet
	PetDto               *dto.PetResponse
	QueriesMock          map[string]string
	Metadata             *dto.FindAllMetadata
	FindAllPetRequest    *dto.FindAllPetRequest
	CreatePetRequest     *dto.CreatePetRequest
	ChangeViewPetRequest *dto.ChangeViewPetRequest
	UpdatePetRequest     *dto.UpdatePetRequest
	AdoptByRequest       *dto.AdoptByRequest
	BindErr              *dto.ResponseErr
	NotFoundErr          *dto.ResponseErr
	ServiceDownErr       *dto.ResponseErr
	InternalErr          *dto.ResponseErr
	Images               []*imgProto.Image
	ImagesList           map[string][]*imgProto.Image
}

func TestPetHandler(t *testing.T) {
	suite.Run(t, new(PetHandlerTest))
}

func (t *PetHandlerTest) SetupTest() {
	petIds := []string{faker.UUIDDigit(), faker.UUIDDigit(), faker.UUIDDigit(), faker.UUIDDigit()}
	imagesList := make(map[string][]*imgProto.Image)
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 3; j++ {
			img := &imgProto.Image{
				Id:        faker.UUIDDigit(),
				PetId:     petIds[i],
				ImageUrl:  faker.URL(),
				ObjectKey: faker.Word(),
			}
			imagesList[petIds[i]] = append(imagesList[petIds[i]], img)
		}
	}
	t.ImagesList = imagesList
	t.Images = imagesList[petIds[0]]
	var pets []*petProto.Pet
	genders := []constant.Gender{constant.MALE, constant.FEMALE}
	statuses := []constant.Status{constant.ADOPTED, constant.FINDHOME}

	for i := 0; i <= 3; i++ {
		pet := &petProto.Pet{
			Id:           faker.UUIDDigit(),
			Type:         faker.Word(),
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
			Origin:       faker.Paragraph(),
			Owner:        faker.Paragraph(),
			Contact:      faker.Paragraph(),
			Tel:          "",
		}

		pets = append(pets, pet)
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.Metadata = &dto.FindAllMetadata{
		Page:       1,
		TotalPages: 1,
		PageSize:   len(t.Pets),
		Total:      len(t.Pets),
	}

	t.PetDto = &dto.PetResponse{
		Id:           t.Pet.Id,
		Type:         t.Pet.Type,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       constant.Gender(t.Pet.Gender),
		Color:        t.Pet.Color,
		Pattern:      t.Pet.Pattern,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       constant.Status(t.Pet.Status),
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
		Tel:          t.Pet.Tel,
	}

	t.QueriesMock = map[string]string{
		"search":   "",
		"type":     "",
		"gender":   "",
		"color":    "",
		"pattern":  "",
		"age":      "",
		"origin":   "",
		"pageSize": "0",
		"page":     "0",
	}

	t.FindAllPetRequest = &dto.FindAllPetRequest{}

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
		Message:    constant.PetNotFoundMessage,
		Data:       nil,
	}

	t.BindErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidIDMessage,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}
}

func (t *PetHandlerTest) TestFindAllSuccess() {
	findAllResponse := utils.ProtoToDtoList(t.Pets, t.ImagesList, false)
	metadataResponse := t.Metadata
	expectedResponse := &dto.FindAllPetResponse{
		Pets:     findAllResponse,
		Metadata: metadataResponse,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Queries().Return(t.QueriesMock)
	petSvc.EXPECT().FindAll(t.FindAllPetRequest, false).Return(expectedResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.FindAll(context)
}

func (t *PetHandlerTest) TestFindOneSuccess() {
	findOneResponse := utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Images))
	expectedResponse := findOneResponse

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(findOneResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestFindOneNotFoundErr() {
	findOneResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(nil, findOneResponse)
	context.EXPECT().JSON(http.StatusNotFound, findOneResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestFindOneGrpcErr() {
	findOneResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().FindOne(t.Pet.Id).Return(nil, findOneResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, findOneResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.FindOne(context)
}

func (t *PetHandlerTest) TestCreateSuccess() {
	createResponse := utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Images))
	expectedResponse := createResponse

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.CreatePetRequest).Return(nil)
	petSvc.EXPECT().Create(t.CreatePetRequest).Return(createResponse, nil)
	context.EXPECT().JSON(http.StatusCreated, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Create(context)
}

func (t *PetHandlerTest) TestCreateGrpcErr() {
	createErrorResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.CreatePetRequest).Return(nil)
	petSvc.EXPECT().Create(t.CreatePetRequest).Return(nil, createErrorResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, createErrorResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Create(context)
}

func (t *PetHandlerTest) TestUpdateSuccess() {
	updateResponse := utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Images))
	expectedResponse := updateResponse

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(updateResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestUpdateNotFound() {
	updateResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(nil, updateResponse)
	context.EXPECT().JSON(http.StatusNotFound, updateResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Update(context)
}

func (t *PetHandlerTest) TestUpdateGrpcErr() {
	updateResponse := t.ServiceDownErr

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.UpdatePetRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdatePetRequest).Return(nil)
	petSvc.EXPECT().Update(t.Pet.Id, t.UpdatePetRequest).Return(nil, updateResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, updateResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
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
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}
func (t *PetHandlerTest) TestDeleteNotFound() {
	deleteResponse := &dto.DeleteResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Delete(context)
}

func (t *PetHandlerTest) TestDeleteGrpcErr() {
	deleteResponse := &dto.DeleteResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	petSvc.EXPECT().Delete(t.Pet.Id).Return(deleteResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
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
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.ChangeViewPetRequest).Return(nil)
	validator.EXPECT().Validate(t.ChangeViewPetRequest).Return(nil)
	petSvc.EXPECT().ChangeView(t.Pet.Id, t.ChangeViewPetRequest).Return(changeViewResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.ChangeView(context)
}

func (t *PetHandlerTest) TestChangeViewNotFound() {
	changeViewResponse := &dto.ChangeViewPetResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.ChangeViewPetRequest).Return(nil)
	validator.EXPECT().Validate(t.ChangeViewPetRequest).Return(nil)
	petSvc.EXPECT().ChangeView(t.Pet.Id, t.ChangeViewPetRequest).Return(changeViewResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.ChangeView(context)
}

func (t *PetHandlerTest) TestChangeViewGrpcErr() {
	changeViewResponse := &dto.ChangeViewPetResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.ChangeViewPetRequest).Return(nil)
	validator.EXPECT().Validate(t.ChangeViewPetRequest).Return(nil)
	petSvc.EXPECT().ChangeView(t.Pet.Id, t.ChangeViewPetRequest).Return(changeViewResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
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
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.AdoptByRequest).Return(nil)
	validator.EXPECT().Validate(t.AdoptByRequest).Return(nil)
	petSvc.EXPECT().Adopt(t.Pet.Id, t.AdoptByRequest).Return(adoptByResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Adopt(context)
}

func (t *PetHandlerTest) TestAdoptNotFound() {
	adoptByResponse := &dto.AdoptByResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.AdoptByRequest).Return(nil)
	validator.EXPECT().Validate(t.AdoptByRequest).Return(nil)
	petSvc.EXPECT().Adopt(t.Pet.Id, t.AdoptByRequest).Return(adoptByResponse, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Adopt(context)
}

func (t *PetHandlerTest) TestAdoptGrpcErr() {
	adoptByResponse := &dto.AdoptByResponse{
		Success: false,
	}

	controller := gomock.NewController(t.T())

	petSvc := petMock.NewMockService(controller)
	imageSvc := imageMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Pet.Id, nil)
	context.EXPECT().Bind(t.AdoptByRequest).Return(nil)
	validator.EXPECT().Validate(t.AdoptByRequest).Return(nil)
	petSvc.EXPECT().Adopt(t.Pet.Id, t.AdoptByRequest).Return(adoptByResponse, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := pet.NewHandler(petSvc, imageSvc, validator)
	handler.Adopt(context)
}
