package prettytable

import (
	"fmt"
	"testing"
)

type testStruct struct {
	Field1 string
	Field2 int `table:"fuck"`
	field3 float64
}

func TestNew(t *testing.T) {

	data := []testStruct{
		{"haha", 1, 3.14},
		{"wwww", 2, 8.14},
		{"oooo", 3, 9.14},
	}

	table := NewTableFromStructs(data)

	tmp := testStruct{"new", 0, 9.99}
	tmp2 := testStruct{"new2", 0, 9.99}
	table.InsertData(2, tmp, tmp2)
	//tt := []testStruct{tmp, tmp2}
	table.AppendData(tmp, tmp2)

	table.String()

	table.AddTitle("title")
	fmt.Println(table.HtmlTable())
}
