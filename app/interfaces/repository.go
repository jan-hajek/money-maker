package interfaces

type Scanable interface {
	Scan(dest ...interface{}) error
}
