package vessel

import (
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

const (
	defaultDBName = "deathstar"
	ext           = "vessel"
)

type Map map[string]any

// Vessel represents an entity
type Vessel struct {
	// Define vessel properties here
	currentDB string
	*Options
	db *bbolt.DB
}

// stand alone binary
func New(options ...OptFunc) (*Vessel, error) {
	// 	ar opt Options
	// p := &opt // p has type *Options

	// // Example 2: pointer to a composite literal (common idiom)
	// p2 := &Options{Field: 1} // p2 is *Options

	// // Example 3: dereference
	// fmt.Println((*p2).Field) // or p2.Field, Go auto-dereferences for selectors
	opts := &Options{
		DBName:  defaultDBName,
		Encoder: JSONEncoder{},
		Decoder: JSONDecoder{},
	}

	for _, function := range options {
		function(opts)
	}

	dbname := fmt.Sprintf("%s.%s", opts.DBName, ext)
	db, err := bbolt.Open(dbname, 0o666, nil)
	if err != nil {
		return nil, err
	}

	return &Vessel{
		db:        db,
		currentDB: dbname,
		Options:   opts,
	}, nil
}

// create collection
// For myself:
// CreateCollection is type *Vessel (receiver is vessel)
//   - pointer receiver (*Vessel) means the method can modify the original Vessel instance
//
// takes a string argument name
// returns two values : a pointer to a bbolt.Bucket and an error
func (vessel *Vessel) CreateCollection(name string) (*bbolt.Bucket, error) {
	tx, err := vessel.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return bucket, err
}

func (vessel *Vessel) Close() error {
	return vessel.db.Close()
}

// clear collection
func (vessel *Vessel) DropDatabase(name string) error {
	dbname := fmt.Sprintf("%s.%s", name, ext)
	return os.Remove(dbname)
}

func (vessel *Vessel) Coll(name string) *Filter {
	return NewFilter(vessel, name)
}
