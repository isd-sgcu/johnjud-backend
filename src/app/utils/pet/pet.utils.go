package pet

import (
	"fmt"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
	petproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imgproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

func MockImageList(n int) [][]*imgproto.Image {
	var imagesList [][]*imgproto.Image
	for i := 0; i <= n; i++ {
		var images []*imgproto.Image
		for j := 0; j <= 3; j++ {
			images = append(images, &imgproto.Image{
				Id:        fmt.Sprintf("%v%v", i, j),
				PetId:     fmt.Sprintf("%v%v", i, j),
				ImageUrl:  fmt.Sprintf("%v%v", i, j),
				ObjectKey: fmt.Sprintf("%v%v", i, j),
			})
		}
		imagesList = append(imagesList, images)
	}

	return imagesList
}

func ProtoToDto(in *petproto.Pet, images []*imgproto.Image) *dto.PetResponse {
	pet := &dto.PetResponse{
		Id:           in.Id,
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       pet.Gender(in.Gender),
		Color:        in.Color,
		Pattern:      in.Pattern,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       pet.Status(in.Status),
		IsSterile:    &in.IsSterile,
		IsVaccinated: &in.IsVaccinated,
		IsVisible:    &in.IsVisible,
		Origin:       in.Origin,
		Address:      in.Address,
		Contact:      in.Contact,
		AdoptBy:      in.AdoptBy,
		Images:       extractImages(images),
	}
	return pet
}

func CreateDtoToProto(in *dto.CreatePetRequest) *petproto.CreatePetRequest {
	return &petproto.CreatePetRequest{
		Pet: &petproto.Pet{
			Type:         in.Type,
			Name:         in.Name,
			Birthdate:    in.Birthdate,
			Gender:       string(in.Gender),
			Color:        in.Color,
			Pattern:      in.Pattern,
			Habit:        in.Habit,
			Caption:      in.Caption,
			Images:       []*imgproto.Image{},
			Status:       string(in.Status),
			IsSterile:    *in.IsSterile,
			IsVaccinated: *in.IsVaccinated,
			IsVisible:    *in.IsVisible,
			Origin:       in.Origin,
			Address:      in.Address,
			Contact:      in.Contact,
			AdoptBy:      in.AdoptBy,
		},
	}
}

func UpdateDtoToProto(id string, in *dto.UpdatePetRequest) *petproto.UpdatePetRequest {
	req := &petproto.UpdatePetRequest{
		Pet: &petproto.Pet{
			Id:        id,
			Type:      in.Type,
			Name:      in.Name,
			Birthdate: in.Birthdate,
			Gender:    string(in.Gender),
			Color:     in.Color,
			Pattern:   in.Pattern,
			Habit:     in.Habit,
			Caption:   in.Caption,
			Images:    []*imgproto.Image{},
			Status:    string(in.Status),
			Origin:    in.Origin,
			Address:   in.Address,
			Contact:   in.Contact,
			AdoptBy:   in.AdoptBy,
		},
	}

	return req
}

func ProtoToDtoList(in []*petproto.Pet, imagesList [][]*imgproto.Image) []*dto.PetResponse {
	var resp []*dto.PetResponse
	for i, p := range in {
		pet := &dto.PetResponse{
			Id:           p.Id,
			Type:         p.Type,
			Name:         p.Name,
			Birthdate:    p.Birthdate,
			Gender:       pet.Gender(p.Gender),
			Color:        p.Color,
			Pattern:      p.Pattern,
			Habit:        p.Habit,
			Caption:      p.Caption,
			Status:       pet.Status(p.Status),
			IsSterile:    &p.IsSterile,
			IsVaccinated: &p.IsVaccinated,
			IsVisible:    &p.IsVisible,
			Origin:       p.Origin,
			Address:      p.Address,
			Contact:      p.Contact,
			AdoptBy:      p.AdoptBy,
			Images:       extractImages(imagesList[i]),
		}
		resp = append(resp, pet)
	}
	return resp
}

func extractImages(images []*imgproto.Image) []dto.ImageResponse {
	var result []dto.ImageResponse
	for _, img := range images {
		result = append(result, dto.ImageResponse{
			Id:  img.Id,
			Url: img.ImageUrl,
		})
	}
	return result
}
