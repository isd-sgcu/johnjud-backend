package pet

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	petConst "github.com/isd-sgcu/johnjud-backend/src/constant/pet"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	image_proto "github.com/isd-sgcu/johnjud-go-proto/johnjub/file/image/v1"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/suite"
)

type PetServiceTest struct {
	suite.Suite
	Pets           []*proto.Pet
	Pet            *proto.Pet
	PetReq         *proto.CreatePetRequest
	UpdatePetReq   *proto.UpdatePetRequest
	PetDto         *dto.PetDto
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
			Gender:       petConst.Gender(rand.Intn(1) + 1),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       petConst.Status(rand.Intn(1) + 1),
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			IsClubPet:    true,
			Background:   faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
		}
		var images []*image_proto.Image
		var imageUrls []string
		for i := 0; i < 3; i++ {
			url := faker.URL()
			images = append(images, &image_proto.Image{
				Id:       faker.UUIDDigit(),
				PetId:    pet.Id,
				ImageUrl: url,
			})
			imageUrls = append(imageUrls, url)
		}
		t.ImagesList = append(t.ImagesList, images)
		t.ImageUrlsList = append(t.ImageUrlsList, imageUrls)
		pets = append(pets, pet)
	}

	t.Pets = pets
	t.Pet = t.Pets[0]

	t.PetReq = &proto.CreatePetRequest{
		Pet: &proto.Pet{
			Type:         faker.Word(),
			Species:      faker.Word(),
			Name:         faker.Name(),
			Birthdate:    faker.Word(),
			Gender:       petConst.Gender(rand.Intn(1) + 1),
			Habit:        faker.Paragraph(),
			Caption:      faker.Paragraph(),
			Status:       petConst.Status(rand.Intn(1) + 1),
			IsSterile:    true,
			IsVaccinated: true,
			IsVisible:    true,
			IsClubPet:    true,
			Background:   faker.Paragraph(),
			Address:      faker.Paragraph(),
			Contact:      faker.Paragraph(),
		},
	}

	t.PetDto = &dto.PetDto{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       petConst.Gender(t.Pet.Gender),
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       petConst.Status(t.Pet.Status),
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}

	t.UpdatePetReq = &proto.UpdatePetRequest{
		Pet: &proto.Pet{
			Id:           t.Pet.Id,
			Type:         t.Pet.Type,
			Species:      t.Pet.Species,
			Name:         t.Pet.Name,
			Birthdate:    t.Pet.Birthdate,
			Gender:       petConst.Gender(t.Pet.Gender),
			Habit:        t.Pet.Habit,
			Caption:      t.Pet.Caption,
			Status:       petConst.Status(t.Pet.Status),
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
		StatusCode: http.StatusServiceUnavailable,
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

func (t *PetServiceTest) TestFindOneSuccess() {}

func (t *PetServiceTest) TestFindOneNotFound() {}

func (t *PetServiceTest) TestFindOneGrpcErr() {}

func (t *PetServiceTest) TestCreateSuccess() {}

func (t *PetServiceTest) TestCreateGrpcErr() {}

func (t *PetServiceTest) TestUpdateSuccess() {}

func (t *PetServiceTest) TestUpdateNotFound() {}

func (t *PetServiceTest) TestUpdateGrpcErr() {}

func (t *PetServiceTest) TestDeleteSuccess() {}

func (t *PetServiceTest) TestDeleteNotFound() {}

func (t *PetServiceTest) TestDeleteGrpcErr() {}
