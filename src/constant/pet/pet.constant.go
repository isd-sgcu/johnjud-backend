package pet

import (
	"encoding/json"
	"errors"
	"strings"
)

type Gender int

const (
	MALE   = 1
	FEMALE = 2
)

type Status int

const (
	ADOPTED  = 1
	FINDHOME = 2
)

func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToUpper(s)
	switch s {
	case "MALE":
		*g = Gender(1)
		return nil
	case "FEMALE":
		*g = Gender(2)
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
		*st = Status(1)
		return nil
	case "FINDHOME":
		*st = Status(2)
		return nil

	default:
		return errors.New("invalid gender")
	}
}

const FindAllPetSuccessMessage = "find all pets success"
const FindOnePetSuccessMessage = "find one pet success"
const CreatePetSuccessMessage = "create pet success"
const UpdatePetSuccessMessage = "update pet success"
const ChangeViewPetSuccessMessage = "change view pet success"
const DeletePetSuccessMessage = "delete pet success"
