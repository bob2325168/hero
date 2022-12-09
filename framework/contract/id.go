package contract

const IDKey = "hero:id"

type IDService interface {
	NewID() string
}
