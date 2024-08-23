package image

import (
	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

func ProtoToDto(in *imageProto.Image) *dto.ImageResponse {
	return &dto.ImageResponse{
		Id:        in.Id,
		PetId:     in.PetId,
		Url:       in.ImageUrl,
		ObjectKey: in.ObjectKey,
	}
}

func ProtoToDtoList(in []*imageProto.Image) []*dto.ImageResponse {
	var res []*dto.ImageResponse
	for _, i := range in {
		res = append(res, &dto.ImageResponse{
			Id:        i.Id,
			PetId:     i.PetId,
			Url:       i.ImageUrl,
			ObjectKey: i.ObjectKey,
		})
	}
	return res
}

func CreateDtoToProto(in *dto.UploadImageRequest) *imageProto.UploadImageRequest {
	return &imageProto.UploadImageRequest{
		Filename: in.Filename,
		Data:     in.File,
		PetId:    in.PetId,
	}
}

func ImageList(in []*dto.ImageResponse) map[string][]*imageProto.Image {
	imagesList := make(map[string][]*imageProto.Image)
	for _, image := range in {
		img := &imageProto.Image{
			Id:        image.Id,
			PetId:     image.PetId,
			ImageUrl:  image.Url,
			ObjectKey: image.ObjectKey,
		}
		imagesList[image.PetId] = append(imagesList[image.PetId], img)
	}

	return imagesList
}
