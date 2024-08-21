package pet

import (
	"errors"
	"strconv"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	petproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imgproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

func ImageList(in []*dto.ImageResponse) map[string][]*imgproto.Image {
	imagesList := make(map[string][]*imgproto.Image)
	for _, image := range in {
		img := &imgproto.Image{
			Id:        image.Id,
			PetId:     image.PetId,
			ImageUrl:  image.Url,
			ObjectKey: image.ObjectKey,
		}
		imagesList[image.PetId] = append(imagesList[image.PetId], img)
	}

	return imagesList
}

func ImageProtoToDto(images []*imgproto.Image) []*dto.ImageResponse {
	var imageDto []*dto.ImageResponse
	for _, image := range images {
		imageDto = append(imageDto, &dto.ImageResponse{
			Id:        image.Id,
			PetId:     image.PetId,
			Url:       image.ImageUrl,
			ObjectKey: image.ObjectKey,
		})
	}
	return imageDto
}

func ImageListProtoToDto(imagesList [][]*imgproto.Image) [][]*dto.ImageResponse {
	var imageListDto [][]*dto.ImageResponse
	for _, images := range imagesList {
		var imageDto []*dto.ImageResponse
		for _, image := range images {
			imageDto = append(imageDto, &dto.ImageResponse{
				Id:        image.Id,
				PetId:     image.PetId,
				Url:       image.ImageUrl,
				ObjectKey: image.ObjectKey,
			})
		}
		imageListDto = append(imageListDto, imageDto)
	}
	return imageListDto
}

func ProtoToDto(in *petproto.Pet, images []*dto.ImageResponse) *dto.PetResponse {
	pet := &dto.PetResponse{
		Id:           in.Id,
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       constant.Gender(in.Gender),
		Color:        in.Color,
		Pattern:      in.Pattern,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       constant.Status(in.Status),
		IsSterile:    &in.IsSterile,
		IsVaccinated: &in.IsVaccinated,
		IsVisible:    &in.IsVisible,
		Origin:       in.Origin,
		Owner:        in.Owner,
		Contact:      in.Contact,
		Tel:          in.Tel,
		Images:       images,
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
			Owner:        in.Owner,
			Contact:      in.Contact,
			Tel:          in.Tel,
		},
	}
}

func UpdateDtoToProto(id string, in *dto.UpdatePetRequest) *petproto.UpdatePetRequest {
	isSterile := false
	if in.IsSterile != nil {
		isSterile = *in.IsSterile
	}
	isVaccinated := false
	if in.IsVaccinated != nil {
		isVaccinated = *in.IsVaccinated
	}
	isVisible := false
	if in.IsVisible != nil {
		isVisible = *in.IsVisible
	}

	req := &petproto.UpdatePetRequest{
		Pet: &petproto.Pet{
			Id:           id,
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
			IsSterile:    isSterile,
			IsVaccinated: isVaccinated,
			IsVisible:    isVisible,
			Origin:       in.Origin,
			Owner:        in.Owner,
			Contact:      in.Contact,
			Tel:          in.Tel,
		},
	}

	return req
}

func ProtoToDtoList(in []*petproto.Pet, imagesList map[string][]*imgproto.Image, isAdmin bool) []*dto.PetResponse {
	var resp []*dto.PetResponse
	for _, p := range in {
		if !isAdmin && !p.IsVisible {
			continue
		}
		pet := &dto.PetResponse{
			Id:           p.Id,
			Type:         p.Type,
			Name:         p.Name,
			Birthdate:    p.Birthdate,
			Gender:       constant.Gender(p.Gender),
			Color:        p.Color,
			Pattern:      p.Pattern,
			Habit:        p.Habit,
			Caption:      p.Caption,
			Status:       constant.Status(p.Status),
			IsSterile:    &p.IsSterile,
			IsVaccinated: &p.IsVaccinated,
			IsVisible:    &p.IsVisible,
			Origin:       p.Origin,
			Owner:        p.Owner,
			Contact:      p.Contact,
			Tel:          p.Tel,
			Images:       ImageProtoToDto(imagesList[p.Id]),
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

func FindAllDtoToProto(in *dto.FindAllPetRequest, isAdmin bool) *petproto.FindAllPetRequest {
	return &petproto.FindAllPetRequest{
		Search:   in.Search,
		Type:     in.Type,
		Gender:   in.Gender,
		Color:    in.Color,
		Pattern:  in.Pattern,
		Origin:   in.Origin,
		PageSize: int32(in.PageSize),
		Page:     int32(in.Page),
		MaxAge:   int32(in.MaxAge),
		MinAge:   int32(in.MinAge),
		IsAdmin:  isAdmin,
	}
}

func MetadataProtoToDto(in *petproto.FindAllPetMetaData) *dto.FindAllMetadata {
	return &dto.FindAllMetadata{
		Page:       int(in.Page),
		TotalPages: int(in.TotalPages),
		PageSize:   int(in.PageSize),
		Total:      int(in.Total),
	}
}

func QueriesToFindAllDto(queries map[string]string) (*dto.FindAllPetRequest, error) {
	request := &dto.FindAllPetRequest{
		Search:   "",
		Type:     "",
		Gender:   "",
		Color:    "",
		Pattern:  "",
		MinAge:   0,
		MaxAge:   0,
		Origin:   "",
		PageSize: 0,
		Page:     0,
	}

	for q, v := range queries {
		switch q {
		case "search":
			request.Search = v
		case "type":
			request.Type = v
		case "gender":
			request.Gender = v
		case "color":
			request.Color = v
		case "pattern":
			request.Pattern = v
		case "minAge":
			minAge, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("error parsing minAge")
			}
			request.MinAge = minAge
		case "maxAge":
			maxAge, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("error parsing maxAge")
			}
			request.MaxAge = maxAge
		case "origin":
			request.Origin = v
		case "pageSize":
			pageSize, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("error pasring pageSize")
			}
			request.PageSize = pageSize
		case "page":
			page, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("error pasring page")
			}
			request.Page = page
		}
	}

	return request, nil
}
