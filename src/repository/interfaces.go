package repository

type Scanable interface {
	Scan(dest ...interface{}) error
}
