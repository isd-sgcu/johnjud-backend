package constant

const FindOneUserSuccessMessage = "find one user success"
const UpdateUserSuccessMessage = "update user success"
const DeleteUserSuccessMessage = "delete user success"

type Role string

const (
	USER  Role = "user"
	ADMIN Role = "admin"
)
