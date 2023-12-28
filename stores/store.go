package stores

type ReadStore interface {
	Get(int64, OrderBy, interface{}) error
	All(interface{}) error
	Count() (int64, error)
	Close() error
}

type ReadWriteStore interface {
	ReadStore
	Save(interface{}) error
	SaveAll([]interface{}) error
	UpdateAllById(map[interface{}]interface{}) error
}

type OrderBy struct {
	Column string
	Order  int
}
