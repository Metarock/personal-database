package api

import (
	"strconv"

	"github.com/Metarock/personal-database/vessel"
)

type FilterMap struct {
	filters map[string]vessel.Map
}

func NewFilterMap() *FilterMap {
	filters := make(map[string]vessel.Map)
	filters[vessel.FilterTypeEQ] = vessel.Map{}
	return &FilterMap{filters: filters}
}

func (filter *FilterMap) Get(filterType string) vessel.Map {
	value, ok := filter.filters[filterType]
	if !ok {
		return vessel.Map{}
	}
	return value
}

func (filter *FilterMap) Add(filterType, key string, value string) {
	if _, ok := filter.filters[filterType]; !ok {
		return
	}

	filter.filters[filterType][key] = ensureCorrectTypeFromString(value)
}

func ensureCorrectTypeFromString(value string) any {
	switch {
	case value == "true":
		return true
	case value == "false":
		return false
	case isInteger(value):
		val, _ := strconv.Atoi(value)
		return val
	case isFloat(value):
		val, _ := strconv.ParseFloat(value, 64)
		return val
	default:
		return value
	}
}

func isFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
