package test

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/constant"
	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"github.com/isd-sgcu/johnjud-backend/internal/pet"
	mock "github.com/isd-sgcu/johnjud-backend/mocks/repository/pet"
	img_mock "github.com/isd-sgcu/johnjud-backend/mocks/service/image"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PetServiceTest struct {
	suite.Suite
	Pet                  *model.Pet
	UpdatePet            *model.Pet
	ChangeViewPet        *model.Pet
	Pets                 []*model.Pet
	PetDto               *dto.PetResponse
	CreatePetReqMock     *dto.CreatePetRequest
	UpdatePetReqMock     *dto.UpdatePetRequest
	ChangeViewPetReqMock *dto.ChangeViewPetRequest
	Images               []*dto.ImageResponse
	ImageUrls            []string
	ImagesList           [][]*dto.ImageResponse
	ChangeAdoptBy        *model.Pet
	AdoptByReq           *dto.AdoptByRequest
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {
	var pets []*model.Pet
	t.ImageUrls = []string{}
	genders := []constant.Gender{constant.MALE, constant.FEMALE}
	statuses := []constant.Status{constant.ADOPTED, constant.FINDHOME}

	for i := 0; i <= 3; i++ {
		pet := &model.Pet{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			Type:         faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       genders[rand.Intn(2)],
			Color:        faker.Word(),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       statuses[rand.Intn(2)],
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			Origin:       faker.Paragraph(),
			Owner:        faker.Paragraph(),
			Contact:      faker.Paragraph(),
			Tel:          "",
		}
		var images []*dto.ImageResponse
		for i := 0; i < 3; i++ {
			url := faker.URL()
			images = append(images, &dto.ImageResponse{
				Id:    faker.UUIDDigit(),
				PetId: pet.ID.String(),
				Url:   url,
			})
			t.ImageUrls = append(t.ImageUrls, url)
		}
		t.ImagesList = append(t.ImagesList, images)
		pets = append(pets, pet)
	}

	t.Pets = pets
	t.Pet = pets[0]
	t.Images = t.ImagesList[0]

	t.PetDto = &dto.PetResponse{
		Id:           t.Pet.ID.String(),
		Type:         t.Pet.Type,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Color:        t.Pet.Color,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
		Images:       t.Images,
	}

	t.UpdatePet = &model.Pet{
		Base: model.Base{
			ID:        t.Pet.Base.ID,
			CreatedAt: t.Pet.Base.CreatedAt,
			UpdatedAt: t.Pet.Base.UpdatedAt,
			DeletedAt: t.Pet.Base.DeletedAt,
		},
		Type:         t.Pet.Type,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Color:        t.Pet.Color,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
	}

	t.ChangeViewPet = &model.Pet{
		Base: model.Base{
			ID:        t.Pet.Base.ID,
			CreatedAt: t.Pet.Base.CreatedAt,
			UpdatedAt: t.Pet.Base.UpdatedAt,
			DeletedAt: t.Pet.Base.DeletedAt,
		},
		Type:      t.Pet.Type,
		Name:      t.Pet.Name,
		Birthdate: t.Pet.Birthdate,
		Gender:    t.Pet.Gender,
		Color:     t.Pet.Color,

		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    false,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
	}

	t.CreatePetReqMock = &dto.CreatePetRequest{
		Type:      t.Pet.Type,
		Name:      t.Pet.Name,
		Birthdate: t.Pet.Birthdate,
		Gender:    t.Pet.Gender,
		Color:     t.Pet.Color,

		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		Images:       t.ImageUrls,
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVaccinated,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
	}

	t.UpdatePetReqMock = &dto.UpdatePetRequest{
		Type:      t.Pet.Type,
		Name:      t.Pet.Name,
		Birthdate: t.Pet.Birthdate,
		Gender:    t.Pet.Gender,
		Color:     t.Pet.Color,

		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		Images:       t.ImageUrls,
		IsSterile:    &t.Pet.IsSterile,
		IsVaccinated: &t.Pet.IsVaccinated,
		IsVisible:    &t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
	}

	t.ChangeViewPetReqMock = &dto.ChangeViewPetRequest{
		Visible: false,
	}

	t.ChangeAdoptBy = &model.Pet{
		Base: model.Base{
			ID:        t.Pet.Base.ID,
			CreatedAt: t.Pet.Base.CreatedAt,
			UpdatedAt: t.Pet.Base.UpdatedAt,
			DeletedAt: t.Pet.Base.DeletedAt,
		},
		Type:      t.Pet.Type,
		Name:      t.Pet.Name,
		Birthdate: t.Pet.Birthdate,
		Gender:    t.Pet.Gender,
		Color:     t.Pet.Color,

		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
		Tel:          t.Pet.Tel,
	}

	t.AdoptByReq = &dto.AdoptByRequest{
		UserID: t.ChangeAdoptBy.Owner,
	}

}
func (t *PetServiceTest) TestDeleteSuccess() {
	want := &dto.DeleteResponse{Success: true}

	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(nil)
	imgSrv := new(img_mock.ServiceMock)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.Delete(t.Pet.ID.String())

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteNotFound() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(gorm.ErrRecordNotFound)
	imgSrv := new(img_mock.ServiceMock)

	srv := pet.NewService(repo, imgSrv)
	_, err := srv.Delete(t.Pet.ID.String())

	assert.Equal(t.T(), http.StatusNotFound, err.StatusCode)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteWithDatabaseError() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(errors.New("internal server error"))
	imgSrv := new(img_mock.ServiceMock)

	srv := pet.NewService(repo, imgSrv)
	_, err := srv.Delete(t.Pet.ID.String())

	assert.Equal(t.T(), http.StatusInternalServerError, err.StatusCode)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteWithUnexpectedError() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(errors.New("unexpected error"))
	imgSrv := new(img_mock.ServiceMock)

	srv := pet.NewService(repo, imgSrv)
	_, err := srv.Delete(t.Pet.ID.String())

	assert.NotNil(t.T(), err)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestFindOneSuccess() {
	want := t.PetDto

	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(t.Pet, nil)
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.FindOne(t.Pet.ID.String())

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

// func (t *PetServiceTest) TestFindAllSuccess() {

// 	want := &proto.FindAllPetResponse{
// 		Pets: t.createPetsDto(t.Pets, t.ImagesList),
// 		Metadata: &proto.FindAllPetMetaData{
// 			Page:       1,
// 			TotalPages: 1,
// 			PageSize:   int32(len(t.Pets)),
// 			Total:      int32(len(t.Pets)),
// 		},
// 	}

// 	var petsIn []*model.Pet

// 	repo := &mock.RepositoryMock{}
// 	repo.On("FindAll", petsIn).Return(&t.Pets, nil)

// 	imgSrv := new(img_mock.ServiceMock)
// 	for i, pet := range t.Pets {
// 		imgSrv.On("FindByPetId", pet.ID.String()).Return(t.ImagesList[i], nil)
// 	}

// 	srv := pet.NewService(repo, imgSrv)

// 	actual, err := srv.FindAll()
// 	assert.Nil(t.T(), err)
// 	assert.Equal(t.T(), want, actual)
// }

func (t *PetServiceTest) TestFindOneNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(nil, errors.New("Not found pet"))
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(nil, nil)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.FindOne(t.Pet.ID.String())

	assert.Equal(t.T(), http.StatusNotFound, err.StatusCode)
	assert.Nil(t.T(), actual)
}

func createPets() []*model.Pet {
	var result []*model.Pet
	genders := []constant.Gender{constant.MALE, constant.FEMALE}
	statuses := []constant.Status{constant.ADOPTED, constant.FINDHOME}

	for i := 0; i < rand.Intn(4)+1; i++ {
		r := &model.Pet{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			Type:         faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       genders[rand.Intn(2)],
			Color:        faker.Word(),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       statuses[rand.Intn(2)],
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			Origin:       faker.Paragraph(),
			Owner:        faker.Paragraph(),
			Contact:      faker.Paragraph(),
		}
		result = append(result, r)
	}

	return result
}

func (t *PetServiceTest) createPetsDto(in []*model.Pet, imagesList [][]*dto.ImageResponse) []*dto.PetResponse {
	var result []*dto.PetResponse

	for i, p := range in {
		r := &dto.PetResponse{
			Id:        p.ID.String(),
			Type:      p.Type,
			Name:      p.Name,
			Birthdate: p.Birthdate,
			Gender:    p.Gender,
			Color:     p.Color,

			Habit:        p.Habit,
			Caption:      p.Caption,
			Status:       p.Status,
			Images:       imagesList[i],
			IsSterile:    &p.IsSterile,
			IsVaccinated: &p.IsVaccinated,
			IsVisible:    &p.IsVisible,
			Origin:       p.Origin,
			Owner:        p.Owner,
			Contact:      p.Contact,
		}

		result = append(result, r)
	}

	return result
}

func (t *PetServiceTest) TestCreateSuccess() {
	want := t.PetDto
	want.Images = t.Images

	repo := &mock.RepositoryMock{}

	in := &model.Pet{
		Type:      t.Pet.Type,
		Name:      t.Pet.Name,
		Birthdate: t.Pet.Birthdate,
		Gender:    t.Pet.Gender,
		Color:     t.Pet.Color,

		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
	}

	repo.On("Create", in).Return(t.Pet, nil)
	imgSrv := new(img_mock.ServiceMock)

	// imageIds := []string{t.CreatePetReqMock.Images[0], t.CreatePetReqMock.Images[1], t.CreatePetReqMock.Images[2]}
	imgSrv.On("AssignPet", &dto.AssignPetRequest{PetId: t.Pet.ID.String(), Ids: t.ImageUrls}).Return(&dto.AssignPetResponse{Success: true})

	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := pet.NewService(repo, imgSrv)

	actual, err := srv.Create(t.CreatePetReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestCreateInternalErr() {
	repo := &mock.RepositoryMock{}

	in := &model.Pet{
		Type:      t.Pet.Type,
		Name:      t.Pet.Name,
		Birthdate: t.Pet.Birthdate,
		Gender:    t.Pet.Gender,
		Color:     t.Pet.Color,

		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		Origin:       t.Pet.Origin,
		Owner:        t.Pet.Owner,
		Contact:      t.Pet.Contact,
	}

	repo.On("Create", in).Return(nil, errors.New("something wrong"))
	imgSrv := new(img_mock.ServiceMock)

	srv := pet.NewService(repo, imgSrv)

	actual, err := srv.Create(t.CreatePetReqMock)

	assert.Equal(t.T(), http.StatusInternalServerError, err.StatusCode)
	assert.Nil(t.T(), actual)
}

func (t *PetServiceTest) TestUpdateSuccess() {
	want := t.PetDto
	updatePet := t.UpdatePet
	updatePet.ID = uuid.Nil

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Pet.ID.String(), t.UpdatePet).Return(t.Pet, nil)
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.Update(t.Pet.ID.String(), t.UpdatePetReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestUpdateNotFound() {
	updatePet := t.UpdatePet
	updatePet.ID = uuid.Nil
	repo := &mock.RepositoryMock{}
	repo.On("Update", t.UpdatePet.ID.String(), t.UpdatePet).Return(nil, errors.New("Not found pet"))
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.UpdatePet.ID.String()).Return(t.Images, nil)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.Update(t.UpdatePet.ID.String(), t.UpdatePetReqMock)

	assert.Equal(t.T(), http.StatusNotFound, err.StatusCode)
	assert.Nil(t.T(), actual)
}

func (t *PetServiceTest) TestChangeViewSuccess() {
	want := &dto.ChangeViewPetResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(t.Pet, nil)
	repo.On("Update", t.Pet.ID.String(), t.ChangeViewPet).Return(t.ChangeViewPet, nil)
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.ChangeView(t.Pet.ID.String(), t.ChangeViewPetReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestChangeViewNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(nil, errors.New("Not found pet"))
	repo.On("Update", t.Pet.ID.String(), t.UpdatePet).Return(nil, errors.New("Not found pet"))
	imgSrv := new(img_mock.ServiceMock)

	srv := pet.NewService(repo, imgSrv)
	actual, err := srv.ChangeView(t.Pet.ID.String(), t.ChangeViewPetReqMock)

	assert.Equal(t.T(), http.StatusNotFound, err.StatusCode)
	assert.Nil(t.T(), actual)
}

func (t *PetServiceTest) TestAdoptBySuccess() {
	want := &dto.AdoptByResponse{Success: true}
	repo := &mock.RepositoryMock{}

	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(t.Pet, nil)
	repo.On("Update", t.Pet.ID.String(), t.ChangeAdoptBy).Return(t.ChangeAdoptBy, nil)

	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := pet.NewService(repo, imgSrv)

	actual, err := srv.Adopt(t.Pet.ID.String(), t.AdoptByReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestAdoptByPetNotFound() {
	wantError := status.Error(codes.NotFound, "pet not found")
	repo := &mock.RepositoryMock{}

	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(nil, wantError)

	imgSrv := new(img_mock.ServiceMock)
	srv := pet.NewService(repo, imgSrv)

	actual, err := srv.Adopt(t.Pet.ID.String(), t.AdoptByReq)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusNotFound, err.StatusCode)
	assert.Nil(t.T(), actual)

	repo.AssertNotCalled(t.T(), "Update", t.Pet.ID.String(), t.ChangeAdoptBy)
}

func (t *PetServiceTest) TestAdoptByUpdateError() {
	wantError := &dto.ResponseErr{StatusCode: http.StatusInternalServerError}
	repo := &mock.RepositoryMock{}

	repo.On("FindOne", t.Pet.ID.String(), &model.Pet{}).Return(t.Pet, nil)
	repo.On("Update", t.Pet.ID.String(), t.ChangeAdoptBy).Return(nil, errors.New("update error"))

	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(nil, wantError)

	srv := pet.NewService(repo, imgSrv)

	actual, err := srv.Adopt(t.Pet.ID.String(), t.AdoptByReq)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, err.StatusCode)
	assert.Nil(t.T(), actual)
}
