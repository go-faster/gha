package lang

import (
	_ "embed"
	"strconv"
	"strings"

	"github.com/go-faster/jx"
)

//go:embed languages.enum.json
var Raw []byte

// Enum value of
func Enum(v string) int {
	return all[v]
}

func Ok(v string) bool {
	if strings.Contains(v, `'`) {
		return false
	}
	_, ok := all[v]
	return ok
}

var ddl string

func DDL() string {
	return ddl
}

var all map[string]int

func init() {
	all = map[string]int{}
	var b strings.Builder
	b.WriteString(`Enum16(`)
	_ = jx.DecodeBytes(Raw).Obj(func(d *jx.Decoder, key string) error {
		v, _ := d.Int()
		all[key] = v

		name := strings.ReplaceAll(key, "'", "\\'")

		if v != 0 {
			b.WriteString(`, `)
		}

		b.WriteString(`'`)
		b.WriteString(name)
		b.WriteString(`'`)

		b.WriteString(" = ")
		b.WriteString(strconv.Itoa(v))

		return nil
	})
	b.WriteString(`)`)

	ddl = b.String()
}
