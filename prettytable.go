package prettytable

/*

* build table with []struct
	- init header with struct tag


* free style


* sort
* insert
* append
* delete

*/

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	tagname        = "table"
	field_spliter  = "|"
	header_spliter = "-"
	corner_char    = "+"
)

var (
	ErrNotTheSameTypeAsBefore = errors.New("not the same type as before")
	ErrOutOfBound             = errors.New("out of bound")
	ErrInvalidRowLength       = errors.New("invalid row length")
	ErrInvalidColumnLength    = errors.New("invalid column length")
)

type Align uint8

const (
	LEFT_ALIGN Align = iota
	RIGHT_ALIGN
	CENTRE_ALIGN
)

type row []interface{}

func headerGen(start, end, spliter string, strs []string) string {
	return start + strings.Join(strs, spliter) + end
}

func csvHeder(strs []string) string {
	return headerGen("", "\n", ",", strs)
}

// border etc to be process
func htmlHeader(strs []string) string {
	return headerGen("<tr><th>", "</tr>\n", "</th><th>", strs)
}

/*
args:  format
*/

func (t *Table) HtmlTable() string {
	var buf bytes.Buffer
	buf.WriteString("<table border=\"1\">")
	if len(t.title) > 0 {
		buf.WriteString(fmt.Sprintf("<tr><td colspan=\"%d\" align=\"center\">%s</td></tr>\n", len(t.header), t.title))
	}
	buf.WriteString(htmlHeader(t.header))
	buf.WriteString(t.renderHtml())
	buf.WriteString("</table>")
	return buf.String()
}

func (t *Table) renderHtml() string {
	var buf bytes.Buffer
	for _, row := range t.Data {
		buf.WriteString(t.readerHtmlOneRow(row))
	}
	return buf.String()
}

func (t *Table) readerHtmlOneRow(row []interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("<tr>")
	for i, r := range row {
		buf.WriteString("<td")
		cf, ok := t.colfmt[t.structMeta[i]]
		if ok {
			// TODO align
			buf.WriteString(fmt.Sprintf(">%s", cf.f(r)))
		} else {
			buf.WriteString(fmt.Sprintf(">%v", r))
		}
		buf.WriteString("</td>")
	}
	buf.WriteString("</tr>\n")
	return buf.String()
}

//TODO
func (t *Table) renderCSV() string {
	var buf bytes.Buffer
	buf.WriteString(csvHeder(t.header))
	for range t.Data {
		buf.WriteString("")
	}
	return buf.String()
}

func (t *Table) PrettyTable() {
}

type ColFmtFunc func(v interface{}) string

type TabFmt struct {
	Align Align
	// font  bold, xie
	// color
}

type ColFmt struct {
	TabFmt
	// %.3f => %s, fmt.Sprintf("%.3f", f)
	f ColFmtFunc // default use %+v
}

type Table struct {
	title      string
	header     []string
	tabfmt     TabFmt
	colfmt     map[string]ColFmt
	structMeta map[int]string
	//TODO header fmt

	Data       [][]interface{}
	table      []row
	stringData [][]string
	structType reflect.Type
	// mode
	// buildFromData bool
}

func (t *Table) AddHeader(header []string) {
	t.header = header
}

func (t *Table) AddTitle(title string) {
	t.title = title
}

func NewTableFromStructs(d interface{}) *Table {

	kind := reflect.TypeOf(d).Kind()

	if kind != reflect.Array && kind != reflect.Slice {
		return nil
	}

	v := reflect.ValueOf(d)
	if v.Len() < 1 {
		return nil
	}

	vv := v.Index(0)
	structType := vv.Type()

	// TODO: public and private field
	header := make([]string, 0, vv.NumField())
	structMeta := make(map[int]string, vv.NumField())

	j := 0
	for i := 0; i < vv.NumField(); i++ {
		field := structType.Field(i)
		n := field.Name
		tag := field.Tag.Get(tagname)
		if tag == "-" {
			continue
		}
		if tag != "" {
			n = tag
		}

		header = append(header, n)
		structMeta[j] = field.Name
		j++
	}

	table := make([][]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		vv := v.Index(i)
		row := make([]interface{}, vv.NumField())
		for i, name := range structMeta {
			row[i] = vv.FieldByName(name)
		}
		table = append(table, row)
	}

	return &Table{
		structType: structType,
		header:     header,
		Data:       table,
		structMeta: structMeta,
	}
}

func (t *Table) Print() {
	//fmt.Println(t.String())
}

func (t *Table) String() {
	for _, h := range t.header {
		fmt.Printf("%v ", h)
	}
	fmt.Println()

	for _, v := range t.Data {
		for _, vv := range v {
			fmt.Printf("%v ", vv)
		}
		fmt.Println()
	}
}

func (t *Table) AppendData(data ...interface{}) (err error) {
	return t.InsertData(len(t.Data), data...)
}

//TODO ...
func (t *Table) InsertData(idx int, data ...interface{}) (err error) {
	//TODO
	// if t.structType == nil {
	// forbidden
	//}

	if idx > len(t.Data) || idx < 0 {
		err = ErrOutOfBound
		return
	}

	rows := make([][]interface{}, len(data))
	for i, d := range data {
		tp := reflect.TypeOf(d)
		if tp != t.structType {
			err = ErrNotTheSameTypeAsBefore
			return
		}

		v := reflect.ValueOf(d)
		row := make([]interface{}, v.NumField())
		for j := 0; j < v.NumField(); j++ {
			row[j] = v.Field(j)
		}
		rows[i] = row
	}

	t.Data = append(t.Data[:idx], append(rows, t.Data[idx:]...)...)

	return
}

func (t *Table) AppendRow(data []interface{}) {
}

func (t *Table) InsertRow(idx int, data []interface{}) {
}

func (t *Table) SetAlign(fieldnames []string, align Align) (err error) {
	return
}

func (t *Table) AddColumn(fieldname string, data []interface{}) {
}
