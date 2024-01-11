package image

import (
	"fmt"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

func ProtoToDto(in *imageProto.Image) *dto.ImageResponse {
	return &dto.ImageResponse{
		Id:        in.Id,
		Url:       in.ImageUrl,
		ObjectKey: in.ObjectKey,
	}
}

func ProtoToDtoList(in []*imageProto.Image) []*dto.ImageResponse {
	var res []*dto.ImageResponse
	for _, i := range in {
		res = append(res, &dto.ImageResponse{
			Id:        i.Id,
			Url:       i.ImageUrl,
			ObjectKey: i.ObjectKey,
		})
	}
	return res
}

func CreateDtoToProto(in *dto.UploadImageRequest) *imageProto.UploadImageRequest {
	return &imageProto.UploadImageRequest{
		Filename: in.Filename,
		Data:     in.Data,
		PetId:    in.PetId,
	}
}

func MockImageList(n int) [][]*imageProto.Image {
	var imagesList [][]*imageProto.Image
	for i := 0; i <= n; i++ {
		var images []*imageProto.Image
		for j := 0; j <= 3; j++ {
			images = append(images, &imageProto.Image{
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
