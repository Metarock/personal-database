package vessel

import (
	"fmt"

	"go.etcd.io/bbolt"
)

const (
	FilterTypeEQ = "eq"
)

func eq(a, b any) bool {
	return a == b
}

type comparison func(a, b any) bool

type compFilter struct {
	kvs  Map
	comp comparison
}

// func — starts a function declaration.
// (filter compFilter) — the receiver: this makes apply a method on type compFilter. It’s like a class method where filter is the method’s this. Because it’s a value receiver (not *compFilter), the method gets a copy of the compFilter.
// apply — the method name.
// (record Map) — parameters: a single parameter named record of type Map (Map is defined as map[string]any).
// bool (if present after the parameters) — the result type: apply returns a boolean.
func (filter compFilter) apply(record Map) bool {
	for key, v := range filter.kvs {
		value, ok := record[key]
		if !ok {
			return false
		}

		if key == "id" {
			return filter.comp(value, uint64(v.(int)))
		}
		return filter.comp(value, v)
	}
	return true
}

type Filter struct {
	vessel      *Vessel
	coll        string
	compFilters []compFilter
	slct        []string
	limit       int
}

func NewFilter(db *Vessel, coll string) *Filter {
	return &Filter{
		vessel:      db,
		coll:        coll,
		compFilters: make([]compFilter, 0),
	}
}

func (filter *Filter) Eq(values Map) *Filter {
	filt := compFilter{
		comp: eq,
		kvs:  values,
	}
	filter.compFilters = append(filter.compFilters, filt)

	return filter
}

// Insert insert the given values.
func (filter *Filter) Insert(values Map) (uint64, error) {
	tx, err := filter.vessel.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	collBucket, err := tx.CreateBucketIfNotExists([]byte(filter.coll))
	if err != nil {
		return 0, err
	}
	id, err := collBucket.NextSequence()
	if err != nil {
		return 0, err
	}
	b, err := filter.vessel.Encoder.Encode(values)
	if err != nil {
		return 0, err
	}
	if err := collBucket.Put(uint64Bytes(id), b); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (filter *Filter) Find() ([]Map, error) {
	tx, err := filter.vessel.db.Begin(true)
	if err != nil {
		return nil, err
	}
	bucket := tx.Bucket([]byte(filter.coll))
	if bucket == nil {
		return nil, fmt.Errorf("bucket (%s) not found", filter.coll)
	}
	records, err := filter.findFiltered(bucket)
	fmt.Println("records", records)
	if err != nil {
		return nil, err
	}
	return records, tx.Commit()
}

func (filter *Filter) Update(values Map) ([]Map, error) {
	tx, err := filter.vessel.db.Begin(true)
	if err != nil {
		return nil, err
	}
	bucket := tx.Bucket([]byte(filter.coll))
	if bucket == nil {
		return nil, fmt.Errorf("bucket (%s) not found", filter.coll)
	}
	records, err := filter.findFiltered(bucket)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		for k, v := range values {
			if _, ok := record[k]; ok {
				record[k] = v
			}
		}
		b, err := filter.vessel.Encoder.Encode(record)
		if err != nil {
			return nil, err
		}
		if err := bucket.Put(uint64Bytes(record["id"].(uint64)), b); err != nil {
			return nil, err
		}
	}
	return records, tx.Commit()
}

func (filter *Filter) Delete() error {
	tx, err := filter.vessel.db.Begin(true)
	if err != nil {
		return err
	}
	bucket := tx.Bucket([]byte(filter.coll))
	if bucket == nil {
		return fmt.Errorf("bucket (%s) not found", filter.coll)
	}
	records, err := filter.findFiltered(bucket)
	if err != nil {
		return err
	}
	for _, r := range records {
		idbytes := uint64Bytes(r["id"].(uint64))
		if err := bucket.Delete(idbytes); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (filter *Filter) Limit(n int) *Filter {
	filter.limit = n
	return filter
}

func (filter *Filter) Select(values ...string) *Filter {
	filter.slct = append(filter.slct, values...)
	return filter
}

func (filter *Filter) findFiltered(bucket *bbolt.Bucket) ([]Map, error) {
	results := []Map{}
	bucket.ForEach(func(k, v []byte) error {
		record := Map{
			"id": uint64FromBytes(k),
		}
		if err := filter.vessel.Decoder.Decode(v, &record); err != nil {
			return err
		}
		include := true
		for _, filter := range filter.compFilters {
			if !filter.apply(record) {
				include = false
				break
			}
		}
		if !include {
			return nil
		}
		record = filter.applySelect(record)
		results = append(results, record)
		return nil
	})
	return results, nil
}

func (filter *Filter) applySelect(record Map) Map {
	if len(filter.slct) == 0 {
		return record
	}
	data := Map{}
	for _, key := range filter.slct {
		if _, ok := record[key]; ok {
			data[key] = record[key]
		}
	}
	return data
}
