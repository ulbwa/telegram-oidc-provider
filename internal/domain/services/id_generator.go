package services

// IdGenerator is an interface for generating unique identifiers.
type IdGenerator interface {
	Generate() string
}
