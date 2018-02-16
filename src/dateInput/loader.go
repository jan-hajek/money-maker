package dateInput

type Loader interface {
	Load() ([]DateInput, error)
}
