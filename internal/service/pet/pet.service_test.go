package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/internal/service/image"
	utils "github.com/isd-sgcu/johnjud-gateway/internal/utils/pet"
	imagemock "github.com/isd-sgcu/johnjud-gateway/mocks/client/image"
	petmock "github.com/isd-sgcu/johnjud-gateway/mocks/client/pet"
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
	MetadataDto           *dto.FindAllMetadata
	MetadataProto         *petproto.FindAllPetMetaData
	PetNotVisible         *petproto.Pet
	FindAllPetReq         *petproto.FindAllPetRequest
	FindAllImageReq       *imgproto.FindAllImageRequest
	UpdatePetReq          *petproto.UpdatePetRequest
	CreatePetReq          *petproto.CreatePetRequest
	ChangeViewPetReq      *petproto.ChangeViewPetRequest
	DeletePetReq          *petproto.DeletePetRequest
	AdoptReq              *petproto.AdoptPetRequest
	PetDto                *dto.PetResponse
	FindAllPetDto         *dto.FindAllPetRequest
	CreatePetDto          *dto.CreatePetRequest
	UpdatePetDto          *dto.UpdatePetRequest
	NotFoundErr           *dto.ResponseErr
	UnavailableServiceErr *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
	InternalErr           *dto.ResponseErr
	ChangeViewedPetDto    *dto.ChangeViewPetRequest
	AdoptDto              *dto.AdoptByRequest

	Images     []*imgproto.Image
	AllImages  []*imgproto.Image
	ImagesList map[string][]*imgproto.Image

	AssignPetReq *imgproto.AssignPetRequest
	AssignPetDto *dto.AssignPetRequest

	FindByPetIdReq *imgproto.FindImageByPetIdRequest
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {
	petIds := []string{faker.UUIDDigit(), faker.UUIDDigit(), faker.UUIDDigit(), faker.UUIDDigit()}
	t.AllImages = []*imgproto.Image{}
	imagesList := make(map[string][]*imgproto.Image)
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 3; j++ {
			img := &imgproto.Image{
				Id:        faker.UUIDDigit(),
				PetId:     petIds[i],
				ImageUrl:  faker.URL(),
				ObjectKey: faker.Word(),
			}
			imagesList[petIds[i]] = append(imagesList[petIds[i]], img)
			t.AllImages = append(t.AllImages, img)
		}
	}
	t.ImagesList = imagesList
	t.Images = imagesList[petIds[0]]
	genders := []constant.Gender{constant.MALE, constant.FEMALE}
	statuses := []constant.Status{constant.ADOPTED, constant.FINDHOME}

	var pets []*petproto.Pet
	for i := 0; i <= 3; i++ {
		pet := &petproto.Pet{
			Id:           petIds[i],
			Type:         faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       string(genders[rand.Intn(2)]),
			Color:        faker.Word(),
			Pattern:      faker.Word(),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Images:       imagesList[petIds[i]],
			Status:       string(statuses[rand.Intn(2)]),
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			Origin:       faker.Paragraph(),
			Owner:        faker.Paragraph(),
			Contact:      faker.Paragraph(),
			Tel:          faker.UUIDDigit(),
		}

		pets = append(pets, pet)
	}

	t.MetadataDto = &dto.FindAllMetadata{
		Page:       1,
		TotalPages: 1,
		PageSize:   len(t.Pets),
		Total:      len(t.Pets),
	}

	t.MetadataProto = &petproto.FindAllPetMetaData{
		Page:       1,
		TotalPages: 1,
		PageSize:   int32(len(t.Pets)),
		Total:      int32(len(t.Pets)),
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetNotVisible = &petproto.Pet{
		Id:           t.Pet.Id,
		Type:         t.Pet.Type,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Color:        t.Pet.Color,
		Pattern:      t.Pet.Pattern,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Images:       t.Pet.Images,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    false,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
		Tel:          t.Pet.Tel,
	}

	t.PetDto = utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Pet.Images))

	t.FindAllPetDto = &dto.FindAllPetRequest{
		Search:   "",
		Type:     "",
		Gender:   "",
		Color:    "",
		Pattern:  "",
		Age:      "",
		Origin:   "",
		PageSize: len(t.Pets),
		Page:     1,
	}

	t.CreatePetDto = &dto.CreatePetRequest{
		Type:         t.Pet.Type,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       constant.Gender(t.Pet.Gender),
		Color:        t.Pet.Color,
		Pattern:      t.Pet.Pattern,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Images:       []string{},
		Status:       constant.Status(t.Pet.Status),
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
		Tel:          t.Pet.Tel,
	}

	t.UpdatePetDto = &dto.UpdatePetRequest{
		Type:         t.Pet.Type,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       constant.Gender(t.Pet.Gender),
		Color:        t.Pet.Color,
		Pattern:      t.Pet.Pattern,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Images:       []string{},
		Status:       constant.Status(t.Pet.Status),
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
		Tel:          t.Pet.Tel,
	}

	t.FindAllPetReq = utils.FindAllDtoToProto(t.FindAllPetDto, true)
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
		UserId: t.Pet.Owner,
	}

	t.AdoptDto = &dto.AdoptByRequest{
		UserID: t.Pet.Owner,
	}

	t.FindAllImageReq = &imgproto.FindAllImageRequest{}

	t.AssignPetReq = &imgproto.AssignPetRequest{
		Ids:   []string{},
		PetId: t.Pet.Id,
	}

	t.AssignPetDto = &dto.AssignPetRequest{
		Ids:   []string{},
		PetId: t.Pet.Id,
	}

	t.FindByPetIdReq = &imgproto.FindImageByPetIdRequest{
		PetId: t.Pet.Id,
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
	protoResp := &petproto.FindAllPetResponse{
		Pets:     t.Pets,
		Metadata: t.MetadataProto,
	}

	findAllPPetsDto := utils.ProtoToDtoList(t.Pets, t.ImagesList, false)
	metadataDto := t.MetadataDto

	expected := &dto.FindAllPetResponse{
		Pets:     findAllPPetsDto,
		Metadata: metadataDto,
	}

	client := petmock.PetClientMock{}
	client.On("FindAll", t.FindAllPetReq).Return(protoResp, nil)

	findAllImageResp := &imgproto.FindAllImageResponse{
		Images: t.AllImages,
	}
	imageClient := imagemock.ImageClientMock{}
	imageClient.On("FindAll", t.FindAllImageReq).Return(findAllImageResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(&client, imageSvc)
	actual, err := svc.FindAll(t.FindAllPetDto, true)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *PetServiceTest) TestFindAllUnavailableServiceError() {
	expected := t.UnavailableServiceErr

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	client := petmock.PetClientMock{}
	client.On("FindAll", t.FindAllPetReq).Return(nil, clientErr)

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(&client, imageSvc)
	actual, err := svc.FindAll(t.FindAllPetDto, true)

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

	findByPetIdReq := t.FindByPetIdReq
	findByPetIdResp := &imgproto.FindImageByPetIdResponse{
		Images: t.Images,
	}

	expected := utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Pet.Images))

	client := petmock.PetClientMock{}
	client.On("FindOne", protoReq).Return(protoResp, nil)

	imageClient := imagemock.ImageClientMock{}
	imageClient.On("FindByPetId", findByPetIdReq).Return(findByPetIdResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(&client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(&client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)

	svc := NewService(&client, imageSvc)
	actual, err := svc.FindOne(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestCreateSuccess() {
	protoReq := t.CreatePetReq
	protoResp := &petproto.CreatePetResponse{
		Pet: t.Pet,
	}

	assignPetReq := t.AssignPetReq
	assignPetResp := &imgproto.AssignPetResponse{
		Success: true,
	}

	findByPetIdReq := t.FindByPetIdReq
	findByPetIdResp := &imgproto.FindImageByPetIdResponse{
		Images: t.Images,
	}

	expected := utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Pet.Images))

	client := &petmock.PetClientMock{}
	client.On("Create", protoReq).Return(protoResp, nil)

	imageClient := imagemock.ImageClientMock{}
	imageClient.On("AssignPet", assignPetReq).Return(assignPetResp, nil)
	imageClient.On("FindByPetId", findByPetIdReq).Return(findByPetIdResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Create(t.CreatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestUpdateSuccess() {
	protoReq := t.UpdatePetReq
	protoResp := &petproto.UpdatePetResponse{
		Pet: t.Pet,
	}

	expected := utils.ProtoToDto(t.Pet, utils.ImageProtoToDto(t.Pet.Images))

	client := &petmock.PetClientMock{}
	client.On("Update", protoReq).Return(protoResp, nil)

	findByPetIdResp := &imgproto.FindImageByPetIdResponse{
		Images: t.Images,
	}
	imageClient := imagemock.ImageClientMock{}
	imageClient.On("FindByPetId", t.FindByPetIdReq).Return(findByPetIdResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Update(t.Pet.Id, t.UpdatePetDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestDeleteSuccess() {
	petProtoReq := &petproto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	petProtoResp := &petproto.DeletePetResponse{
		Success: true,
	}
	imageProtoReq := &imgproto.DeleteByPetIdRequest{
		PetId: t.Pet.Id,
	}
	imageProtoResp := &imgproto.DeleteByPetIdResponse{
		Success: true,
	}

	expected := &dto.DeleteResponse{Success: true}

	client := &petmock.PetClientMock{}
	client.On("Delete", petProtoReq).Return(petProtoResp, nil)

	imageClient := imagemock.ImageClientMock{}
	imageClient.On("DeleteByPetId", imageProtoReq).Return(imageProtoResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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
	imageProtoReq := &imgproto.DeleteByPetIdRequest{
		PetId: t.Pet.Id,
	}
	imageProtoResp := &imgproto.DeleteByPetIdResponse{
		Success: true,
	}

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petmock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, clientErr)

	imageClient := imagemock.ImageClientMock{}
	imageClient.On("DeleteByPetId", imageProtoReq).Return(imageProtoResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Delete(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestDeleteServiceUnavailableError() {
	protoReq := &petproto.DeletePetRequest{
		Id: t.Pet.Id,
	}
	protoResp := &petproto.DeletePetResponse{
		Success: false,
	}
	imageProtoReq := &imgproto.DeleteByPetIdRequest{
		PetId: t.Pet.Id,
	}
	imageProtoResp := &imgproto.DeleteByPetIdResponse{
		Success: true,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petmock.PetClientMock{}
	client.On("Delete", protoReq).Return(protoResp, clientErr)

	imageClient := imagemock.ImageClientMock{}
	imageClient.On("DeleteByPetId", imageProtoReq).Return(imageProtoResp, nil)

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Delete(t.Pet.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestChangeViewSuccess() {
	protoReq := t.ChangeViewPetReq
	protoResp := &petproto.ChangeViewPetResponse{
		Success: true,
	}

	client := &petmock.PetClientMock{}
	client.On("ChangeView", protoReq).Return(protoResp, nil)

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
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

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.ChangeView(t.Pet.Id, t.ChangeViewedPetDto)

	assert.Equal(t.T(), &dto.ChangeViewPetResponse{Success: false}, actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestAdoptSuccess() {
	protoReq := t.AdoptReq
	protoResp := &petproto.AdoptPetResponse{
		Success: true,
	}

	client := &petmock.PetClientMock{}
	client.On("AdoptPet", protoReq).Return(protoResp, nil)

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Adopt(t.Pet.Id, t.AdoptDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), actual, &dto.AdoptByResponse{Success: true})
}

func (t *PetServiceTest) TestAdoptNotFoundError() {
	protoReq := t.AdoptReq

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := &petmock.PetClientMock{}
	client.On("AdoptPet", protoReq).Return(nil, clientErr)

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Adopt(t.Pet.Id, t.AdoptDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *PetServiceTest) TestAdoptUnavailableServiceError() {
	protoReq := t.AdoptReq

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &petmock.PetClientMock{}
	client.On("AdoptPet", protoReq).Return(nil, clientErr)

	imageClient := imagemock.ImageClientMock{}

	imageSvc := imageSvc.NewService(&imageClient)
	svc := NewService(client, imageSvc)
	actual, err := svc.Adopt(t.Pet.Id, t.AdoptDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}
