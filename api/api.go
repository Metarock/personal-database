package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Metarock/personal-database/vessel"
	"github.com/labstack/echo/v4"
)

type Server struct {
	db *vessel.Vessel
}

func NewServer(db *vessel.Vessel) *Server {
	return &Server{db: db}
}

func (server *Server) HandlePostInsert(context echo.Context) error {
	collname := context.Param("collname")
	var data vessel.Map
	if err := json.NewDecoder(context.Request().Body).Decode(&data); err != nil {
		return err
	}
	id, err := server.db.Coll(collname).Insert(data)
	if err != nil {
		return err
	}
	return context.JSON(http.StatusCreated, vessel.Map{"id": id})
}

func (server *Server) HandleGetQuery(context echo.Context) error {
	var (
		collname  = context.Param("collname")
		filterMap = NewFilterMap()
	)

	for key, value := range context.QueryParams() {
		filterParts := strings.Split(key, ".")
		if len(filterParts) != 2 {
			return fmt.Errorf("malformed query")
		}

		if len(value) == 0 {
			return fmt.Errorf("malformed query")
		}

		if value[0] == "" {
			return fmt.Errorf("malformed query")
		}

		var (
			filterType = filterParts[1]
			filterKey  = filterParts[0]
			filterVal  = value[0]
		)

		filterMap.Add(filterType, filterKey, filterVal)
	}
	records, err := server.db.Coll(collname).Eq(filterMap.Get(vessel.FilterTypeEQ)).Find()
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, records)
}
