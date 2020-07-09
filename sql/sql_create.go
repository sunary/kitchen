package sql

import (
	"reflect"
	"strings"
	"time"

	"github.com/sunary/kitchen/str"
)

const (
	typePrefix      = "type:"
	defaultPrefix   = "default:"
	isPrimaryKey    = "primary_key"
	isAutoIncrement = "auto_increment"
	funcTableName   = "TableName"
)

func SqlCreateTable(tb interface{}) string {
	fields := []string{}
	v := reflect.ValueOf(tb)
	t := reflect.TypeOf(tb)
	for j := 0; j < t.NumField(); j++ {
		field := t.Field(j)
		gtag := field.Tag.Get(gormTag)
		if gtag == "-" {
			continue
		}

		gts := strings.Split(gtag, ";")
		fs := []string{""}
		columnDeclare := str.ToSnakeCase(field.Name)
		defaultDeclare := ""
		typeDeclare := ""
		isPkDeclare := false
		isAutoDeclare := false
		for _, gt := range gts {
			if strings.HasPrefix(gt, columnPrefix) {
				columnDeclare = strings.TrimPrefix(gt, columnPrefix)
			} else if strings.HasPrefix(gt, typePrefix) {
				typeDeclare = strings.TrimPrefix(gt, typePrefix)
			} else if strings.HasPrefix(gt, defaultPrefix) {
				defaultDeclare = "DEFAULT " + strings.TrimPrefix(gt, defaultPrefix)
			} else if strings.ToLower(gt) == isPrimaryKey {
				isPkDeclare = true
			} else if strings.ToLower(gt) == isAutoIncrement {
				isAutoDeclare = true
			}
		}

		fs = append(fs, "`"+columnDeclare+"`")
		if typeDeclare != "" {
			fs = append(fs, typeDeclare)
		} else {
			fs = append(fs, sqlType(v.Field(j).Interface()))
		}
		if defaultDeclare != "" {
			fs = append(fs, defaultDeclare)
		}
		if isAutoDeclare {
			fs = append(fs, "AUTO INCREMENT")
		}
		if isPkDeclare {
			fs = append(fs, "PRIMARY KEY")
		}

		fields = append(fields, strings.Join(fs, " "))
	}

	return "CREATE TABLE " + getTableName(tb) + "(" + strings.Join(fields, ",\n") + ");"
}

func sqlType(v interface{}) string {
	switch v.(type) {
	case bool:
		return "BOOLEAN"
	case int8, uint8:
		return "TINYINT"
	case int16, uint16:
		return "SMALLINT"
	case int, int32, uint32:
		return "INT"
	case int64, uint64:
		return "BIGINT"
	case float32:
		return "FLOAT"
	case float64:
		return "DOUBLE"
	case string:
		return "TEXT"
	case time.Time:
		return "TIMESTAMP"
	}
	return "UNSPECIFIED"
}

func getTableName(t interface{}) string {
	st := reflect.TypeOf(t)
	if _, ok := st.MethodByName(funcTableName); ok {
		v := reflect.ValueOf(t).MethodByName(funcTableName).Call(nil)
		if len(v) > 0 {
			return v[0].String()
		}
	}

	name := ""
	if t := reflect.TypeOf(t); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return str.ToSnakeCase(name)
}
