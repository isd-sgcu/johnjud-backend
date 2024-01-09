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
		Species:      in.Species,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       pet.Gender(in.Gender),
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       pet.Status(in.Status),
		IsSterile:    &in.IsSterile,
		IsVaccinated: &in.IsVaccinated,
		IsVisible:    &in.IsVisible,
		IsClubPet:    &in.IsClubPet,
		Background:   in.Background,
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
			Species:      in.Species,
			Name:         in.Name,
			Birthdate:    in.Birthdate,
			Gender:       petproto.Gender(in.Gender),
			Habit:        in.Habit,
			Caption:      in.Caption,
			Images:       []*imgproto.Image{},
			Status:       petproto.PetStatus(in.Status),
			IsSterile:    *in.IsSterile,
			IsVaccinated: *in.IsVaccinated,
			IsVisible:    *in.IsVisible,
			IsClubPet:    *in.IsClubPet,
			Background:   in.Background,
			Address:      in.Address,
			Contact:      in.Contact,
			AdoptBy:      in.AdoptBy,
		},
	}
}

func UpdateDtoToProto(id string, in *dto.UpdatePetRequest) *petproto.UpdatePetRequest {
	req := &petproto.UpdatePetRequest{
		Pet: &petproto.Pet{
			Id:         id,
			Type:       in.Type,
			Species:    in.Species,
			Name:       in.Name,
			Birthdate:  in.Birthdate,
			Gender:     petproto.Gender(in.Gender),
			Habit:      in.Habit,
			Caption:    in.Caption,
			Images:     []*imgproto.Image{},
			Status:     petproto.PetStatus(in.Status),
			Background: in.Background,
			Address:    in.Address,
			Contact:    in.Contact,
			AdoptBy:    in.AdoptBy,
		},
	}

	if in.IsClubPet == nil {
		req.Pet.IsClubPet = false
	} else {
		req.Pet.IsClubPet = *in.IsClubPet
	}

	if in.IsSterile == nil {
		req.Pet.IsSterile = false
	} else {
		req.Pet.IsSterile = *in.IsSterile
	}

	if in.IsVaccinated == nil {
		req.Pet.IsVaccinated = false
	} else {
		req.Pet.IsVaccinated = *in.IsVaccinated
	}

	if in.IsVisible == nil {
		req.Pet.IsVisible = false
	} else {
		req.Pet.IsVisible = *in.IsVisible
	}

	return req
}

func ProtoToDtoList(in []*petproto.Pet, imagesList [][]*imgproto.Image) []*dto.PetResponse {
	var resp []*dto.PetResponse
	for i, p := range in {
		pet := &dto.PetResponse{
			Id:           p.Id,
			Type:         p.Type,
			Species:      p.Species,
			Name:         p.Name,
			Birthdate:    p.Birthdate,
			Gender:       pet.Gender(p.Gender),
			Habit:        p.Habit,
			Caption:      p.Caption,
			Status:       pet.Status(p.Status),
			IsSterile:    &p.IsSterile,
			IsVaccinated: &p.IsVaccinated,
			IsVisible:    &p.IsVisible,
			IsClubPet:    &p.IsClubPet,
			Background:   p.Background,
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
func IsLike(petId string, likes []*dto.LikeResponse) *bool {
	for _, like := range likes {
		if like.PetID == petId {
			return BoolAddr(true)
		}
	}
	return BoolAddr(false)
}

func MapIsLikeToPets(likes []*dto.LikeResponse, pets []*dto.PetResponse) []*dto.PetResponse {
	for _, pet := range pets {
		pet.IsLike = IsLike(pet.Id, likes)
	}
	return pets
}

func BoolAddr(b bool) *bool {
	return &b
}
