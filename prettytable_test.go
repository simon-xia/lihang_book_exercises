package prettytable

import (
	"fmt"
	"testing"
)

const unit_format = "%.3f %s"

type testStruct struct {
	Field1 string
	Field2 int     `table:"fuck"`
	Field3 float64 `table:"f"`
}

func float64Fmt(v interface{}) string {
	return fmt.Sprintf("%.3f", v)
}

func TestNew(t *testing.T) {

	data := []testStruct{
		{"haha", 12314124, 3.143241123412},
		{"wwww", 2, 8.144234525},
		{"oooo", 3, 9.14},
	}

	table := NewTableFromStructs(data)

	tmp := testStruct{"new", 0, 9.99}
	tmp2 := testStruct{"new2", 0, 9.99}
	table.InsertData(2, tmp, tmp2)
	//tt := []testStruct{tmp, tmp2}
	table.AppendData(tmp, tmp2)

	table.String()

	table.AddFormat(map[string]ColFmt{
		"field3": ColFmt{
			Func: float64Fmt,
		},
	})
	table.AddTitle("title")
	fmt.Println(table.HtmlTable())
	fmt.Printf("%+v", table)
}
