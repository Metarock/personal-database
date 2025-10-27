package vessel

import (
	"log"
	"testing"
)

func TestDelete(test *testing.T) {
	db, err := New(WithDBName("test"))
	if err != nil {
		test.Fatal(err)
	}
	defer db.DropDatabase("test")
	id, err := db.Coll("users").Insert(Map{"name": "foo"})
	if err != nil {
		test.Fatal(err)
	}
	delete := Map{"id": id}
	if err := db.Coll("users").Eq(delete).Delete(); err != nil {
		test.Fatal(err)
	}
	records, err := db.Coll("users").Find()
	if err != nil {
		test.Fatal(err)
	}
	if len(records) != 0 {
		test.Fatalf("expected to have 0 records got %d", len(records))
	}
}

func TestUpdate(test *testing.T) {
	db, err := New(WithDBName("test"))
	if err != nil {
		test.Fatal(err)
	}
	defer db.DropDatabase("test")
	_, err = db.Coll("users").Insert(Map{"name": "foo"})
	if err != nil {
		test.Fatal(err)
	}
	values := Map{"name": "bar"}
	results, err := db.Coll("users").Update(values)
	if err != nil {
		test.Fatal(err)
	}
	if len(results) != 1 {
		log.Fatalf("expected to have 1 result got %d", len(results))
	}
	records, err := db.Coll("users").Find()
	if err != nil {
		test.Fatal(err)
	}
	if len(records) != 1 {
		log.Fatalf("expected to have 1 result got %d", len(results))
	}
	if records[0]["name"] != values["name"] {
		test.Fatalf("expected name to be %s got %s", values["name"], records[0]["name"])
	}
}

func TestInsert(test *testing.T) {
	values := []Map{
		{
			"name": "Foo",
			"age":  10,
		},
		{
			"name": "Bar",
			"age":  88.3,
		},
		{
			"name": "Baz",
			"age":  10,
		},
	}

	db, err := New(WithDBName("test"))
	if err != nil {
		test.Fatal(err)
	}
	defer db.DropDatabase("test")
	for i, data := range values {
		id, err := db.Coll("users").Insert(data)
		if err != nil {
			test.Fatal(err)
		}
		if id != uint64(i+1) {
			test.Fatalf("expect ID %d got %d", i, id)
		}
	}
	users, err := db.Coll("users").Find()
	if err != nil {
		test.Fatal(err)
	}
	if len(users) != len(values) {
		test.Fatalf("expecting %d result got %d", len(values), len(users))
	}
}

func TestFind(test *testing.T) {
	db, err := New(WithDBName("test"))
	if err != nil {
		test.Fatal(err)
	}
	defer db.DropDatabase("test")

	coll := "users"
	db.Coll(coll).Insert(Map{"username": "James007"})
	db.Coll(coll).Insert(Map{"username": "Acice"})
	db.Coll(coll).Insert(Map{"username": "Bob"})
	db.Coll(coll).Insert(Map{"username": "Mike"})

	results, err := db.Coll("users").Eq(Map{"username": "James007"}).Find()
	if err != nil {
		test.Fatal(err)
	}
	if len(results) != 1 {
		test.Fatalf("expected 1 result got %d", len(results))
	}
}
