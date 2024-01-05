package pet

import (
	"fmt"

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
