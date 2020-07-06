package types

//Error type of storage
type Error string

func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	ErrDateBusy      Error = "this time is already in use"
	ErrNotFound      Error = "event not found"
	ErrEventDeleted  Error = "event was deleted"
	ErrEventIdExists Error = "event with this ID already exists"
)
