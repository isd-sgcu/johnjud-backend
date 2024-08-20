package pet

import (
	"encoding/json"
	"errors"
	"strings"
)

type Gender string

const (
	MALE   Gender = "male"
	FEMALE Gender = "female"
)

type Status string

const (
	ADOPTED  Status = "adopted"
	FINDHOME Status = "findhome"
)

func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToUpper(s)
	switch s {
	case "MALE":
		*g = MALE
		return nil
	case "FEMALE":
		*g = FEMALE
		return nil

	default:
		return errors.New("invalid gender")
	}
}

func (st *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToUpper(s)
	switch s {
	case "ADOPTED":
		*st = ADOPTED
		return nil
	case "FINDHOME":
		*st = FINDHOME
		return nil

	default:
		return errors.New("invalid status")
	}
}

const FindAllPetSuccessMessage = "find all pets success"
const FindOnePetSuccessMessage = "find one pet success"
const CreatePetSuccessMessage = "create pet success"
const UpdatePetSuccessMessage = "update pet success"
const ChangeViewPetSuccessMessage = "change view pet success"
const DeletePetSuccessMessage = "delete pet success"
const AdoptPetSuccessMessage = "adopt pet success"
