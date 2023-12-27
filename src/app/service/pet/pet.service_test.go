package pet

import (
	"testing"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/suite"
)

type PetServiceTest struct {
	suite.Suite
	Pet            *proto.Pet
	PetReq         *proto.CreatePetRequest
	UpdatePetReq   *proto.UpdatePetRequest
	PetDto         *dto.PetDto
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {}

func (t *PetServiceTest) TestFindAllSuccess() {}

func (t *PetServiceTest) TestFindOneSuccess() {}

func (t *PetServiceTest) TestFindOneNotFound() {}

func (t *PetServiceTest) TestFindOneGrpcErr() {}

func (t *PetServiceTest) TestCreateSuccess() {}

func (t *PetServiceTest) TestCreateGrpcErr() {}

func (t *PetServiceTest) TestUpdateSuccess() {}

func (t *PetServiceTest) TestUpdateNotFound() {}

func (t *PetServiceTest) TestUpdateGrpcErr() {}

func (t *PetServiceTest) TestDeleteSuccess() {}

func (t *PetServiceTest) TestDeleteNotFound() {}

func (t *PetServiceTest) TestDeleteGrpcErr() {}
