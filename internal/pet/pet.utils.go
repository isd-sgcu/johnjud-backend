package pet

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/constant"
	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func FilterPet(in *[]*model.Pet, query *dto.FindAllPetRequest) error {
	if query.MaxAge == 0 {
		query.MaxAge = math.MaxInt32
	}

	var results []*model.Pet
	for _, p := range *in {
		res, err := filterAge(p, query.MinAge, query.MaxAge)
		if err != nil {
			return err
		}
		if !res {
			continue
		}
		if query.Search != "" && !strings.Contains(p.Name, query.Search) {
			continue
		}
		if query.Type != "" && p.Type != query.Type {
			continue
		}
		if query.Gender != "" && p.Gender != constant.Gender(query.Gender) {
			continue
		}
		if query.Color != "" && p.Color != query.Color {
			continue
		}
		if query.Origin != "" && p.Origin != query.Origin {
			continue
		}

		results = append(results, p)
	}
	*in = results
	return nil
}

func PaginatePets(pets *[]*model.Pet, page int, pageSize int, metadata *dto.FindAllMetadata) error {
	totalsPets := int(len(*pets))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = totalsPets
	}
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalsPets {
		*pets = []*model.Pet{}
		return nil
	}
	if end > totalsPets {
		end = totalsPets
	}
	*pets = (*pets)[start:end]

	totalPages := int(math.Ceil(float64(totalsPets) / float64(pageSize)))

	metadata.Page = page
	metadata.PageSize = pageSize
	metadata.Total = totalsPets
	metadata.TotalPages = totalPages
	return nil
}

func RawToDtoList(in *[]*model.Pet, images map[string][]*dto.ImageResponse, query *dto.FindAllPetRequest) ([]*dto.PetResponse, error) {
	var result []*dto.PetResponse
	if len(*in) != len(images) {
		return nil, errors.New("length of in and imageUrls have to be the same")
	}

	for _, p := range *in {
		// TODO: create new filter image function this wont work
		result = append(result, RawToDto(p, images[p.ID.String()]))
	}
	return result, nil
}

func RawToDto(in *model.Pet, images []*dto.ImageResponse) *dto.PetResponse {
	return &dto.PetResponse{
		Id:           in.ID.String(),
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       in.Gender,
		Color:        in.Color,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       in.Status,
		Images:       images,
		IsSterile:    &in.IsSterile,
		IsVaccinated: &in.IsVaccinated,
		IsVisible:    &in.IsVisible,
		Origin:       in.Origin,
		Owner:        in.Owner,
		Contact:      in.Contact,
		Tel:          in.Tel,
	}
}

func DtoToRaw(in *dto.PetResponse) (res *model.Pet, err error) {
	var id uuid.UUID

	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	return &model.Pet{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       in.Gender,
		Color:        in.Color,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       in.Status,
		IsSterile:    *in.IsSterile,
		IsVaccinated: *in.IsVaccinated,
		IsVisible:    *in.IsVisible,
		Origin:       in.Origin,
		Owner:        in.Owner,
		Contact:      in.Contact,
		Tel:          in.Tel,
	}, nil
}

func ExtractImageUrls(in []*dto.ImageResponse) []string {
	var result []string
	for _, e := range in {
		result = append(result, e.Url)
	}
	return result
}

func ExtractImageIDs(in []*dto.ImageResponse) []string {
	var result []string
	for _, e := range in {
		result = append(result, e.Id)
	}
	return result
}

func UpdateMap(in *model.Pet) map[string]interface{} {
	updateMap := make(map[string]interface{})
	t := reflect.TypeOf(*in)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeName := field.Type.Name()
		value := reflect.ValueOf(*in).Field(i).Interface()
		if (typeName == "string" || typeName == "Gender" || typeName == "Status") && value != "" {
			updateMap[field.Name] = value
		}
		if typeName == "bool" {
			updateMap[field.Name] = value
		}
	}
	return updateMap
}

func parseDate(dateStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func filterAge(pet *model.Pet, minAge, maxAge int) (bool, error) {
	log.Info().
		Str("service", "filterAge").
		Str("module", "birth").Msgf("birthdate: %s", pet.Birthdate)
	birthdate, err := parseDate(pet.Birthdate)
	if err != nil {
		return false, err
	}

	currYear := time.Now()
	birthYear := birthdate
	diff := currYear.Sub(birthYear).Hours() / constant.DAY / constant.YEAR
	log.Info().
		Str("service", "filterAge").
		Str("module", "FilterAge").Msgf("diff: %f", diff)

	return diff >= float64(minAge) && diff <= float64(maxAge), nil
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

func UpdateDtoToModel(in *dto.UpdatePetRequest) *model.Pet {
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

	req := &model.Pet{
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       in.Gender,
		Color:        in.Color,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       in.Status,
		IsSterile:    isSterile,
		IsVaccinated: isVaccinated,
		IsVisible:    isVisible,
		Origin:       in.Origin,
		Owner:        in.Owner,
		Contact:      in.Contact,
		Tel:          in.Tel,
	}

	return req
}

func CreateDtoToModel(in *dto.CreatePetRequest) *model.Pet {
	return &model.Pet{
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       in.Gender,
		Color:        in.Color,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       in.Status,
		IsSterile:    *in.IsSterile,
		IsVaccinated: *in.IsVaccinated,
		IsVisible:    *in.IsVisible,
		Origin:       in.Origin,
		Owner:        in.Owner,
		Contact:      in.Contact,
		Tel:          in.Tel,
	}
}

func ImageList(in []*dto.ImageResponse) map[string][]*dto.ImageResponse {
	imagesList := make(map[string][]*dto.ImageResponse)
	for _, image := range in {
		img := &dto.ImageResponse{
			Id:        image.Id,
			PetId:     image.PetId,
			Url:       image.Url,
			ObjectKey: image.ObjectKey,
		}
		imagesList[image.PetId] = append(imagesList[image.PetId], img)
	}

	return imagesList
}

// func ProtoToDtoList(in []*petproto.Pet, imagesList map[string][]*imgproto.Image, isAdmin bool) []*dto.PetResponse {
// 	var resp []*dto.PetResponse
// 	for _, p := range in {
// 		if !isAdmin && !p.IsVisible {
// 			continue
// 		}
// 		pet := &dto.PetResponse{
// 			Id:           p.Id,
// 			Type:         p.Type,
// 			Name:         p.Name,
// 			Birthdate:    p.Birthdate,
// 			Gender:       constant.Gender(p.Gender),
// 			Color:        p.Color,
// 			Pattern:      p.Pattern,
// 			Habit:        p.Habit,
// 			Caption:      p.Caption,
// 			Status:       constant.Status(p.Status),
// 			IsSterile:    &p.IsSterile,
// 			IsVaccinated: &p.IsVaccinated,
// 			IsVisible:    &p.IsVisible,
// 			Origin:       p.Origin,
// 			Owner:        p.Owner,
// 			Contact:      p.Contact,
// 			Tel:          p.Tel,
// 			Images:       ImageProtoToDto(imagesList[p.Id]),
// 		}
// 		resp = append(resp, pet)
// 	}
// 	return resp
// }
