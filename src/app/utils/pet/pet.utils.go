package pet

import (
	"github.com/go-faker/faker/v4"
	imgproto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

func MockImageList(n int) [][]*imgproto.Image {
	var imagesList [][]*imgproto.Image
	for i := 0; i <= n; i++ {
		var images []*imgproto.Image
		for j := 0; j <= 3; j++ {
			images = append(images, &imgproto.Image{
				Id:        faker.UUIDDigit(),
				PetId:     faker.UUIDDigit(),
				ImageUrl:  faker.URL(),
				ObjectKey: faker.Word(),
			})
		}
		imagesList = append(imagesList, images)
	}

	return imagesList
}

func MockImages() []*imgproto.Image {
	var images []*imgproto.Image
	for j := 0; j <= 3; j++ {
		images = append(images, &imgproto.Image{
			Id:        faker.UUIDDigit(),
			PetId:     faker.UUIDDigit(),
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Word(),
		})
	}
	return images
}
