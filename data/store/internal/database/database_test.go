package database

import (
	"messenger/data/common"
	"testing"
)

type Type struct {
	Id    int `gorm:"primaryKey"`
	Array common.JsonArray[string]
}

func TestDatabase(t *testing.T) {
	var table = ChatDatabase.Table("test")
	var e = table.AutoMigrate(&Type{})
	if e != nil {
		t.Fatal(e)
	}

	// Insert a few elements
	{
		e = table.Create(&Type{
			Id:    1,
			Array: common.JsonArray[string]{"a", "b", "c"},
		}).Error
		if e != nil {
			t.Fatal(e)
		}
		e = table.Create(&Type{
			Id:    2,
			Array: common.JsonArray[string]{"c", "d", "e", "f"},
		}).Error
		if e != nil {
			t.Fatal(e)
		}
	}
	// Find the elements where the array contains "c"
	{
		var res []Type
		e = table.Raw("SELECT * FROM test, json_each(test.array) WHERE json_each.value = ?", "c").Scan(&res).Error
		if e != nil {
			t.Fatal(e)
		}
		if len(res) != 2 {
			t.Fatal("invalid result")
		}
	}
	// Find the elements where the array contains "d"
	{
		var res []Type
		e = table.Raw("SELECT * FROM test, json_each(test.array) WHERE json_each.value = ?", "d").Scan(&res).Error
		if e != nil {
			t.Fatal(e)
		}
		if len(res) != 1 {
			t.Fatal("invalid result")
		}
	}
}
