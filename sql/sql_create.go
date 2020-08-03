package sql

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sunary/kitchen/str"
)

const (
	typePrefix          = "type:"
	defaultPrefix       = "default:"
	indexPrefix         = "index:"
	foreignKeyPrefix    = "foreignkey:"
	associationFkPrefix = "association_foreignkey:"
	isPrimaryKey        = "primary_key"
	isUnique            = "unique"
	isAutoIncrement     = "auto_increment"
	funcTableName       = "TableName"
)

func SqlCreateTable(tb interface{}) string {
	tableName := getTableName(tb)
	maxLen := 0

	fields := [][]string{}
	indexes := []string{}
	v := reflect.ValueOf(tb)
	t := reflect.TypeOf(tb)
	for j := 0; j < t.NumField(); j++ {
		field := t.Field(j)
		gtag := field.Tag.Get(gormTag)
		if gtag == "-" {
			continue
		}

		gts := strings.Split(gtag, ";")
		columnDeclare := str.ToSnakeCase(field.Name)
		defaultDeclare := ""
		typeDeclare := ""
		isFkDeclare := false
		isPkDeclare := false
		isUniqueDeclare := false
		indexDeclare := ""
		isAutoDeclare := false
		for _, gt := range gts {
			gtLower := strings.ToLower(gt)
			if strings.HasPrefix(gtLower, foreignKeyPrefix) || strings.HasPrefix(gtLower, associationFkPrefix) {
				isFkDeclare = true
				break
			} else if strings.HasPrefix(gtLower, columnPrefix) {
				columnDeclare = gt[len(columnPrefix):]
			} else if strings.HasPrefix(gtLower, typePrefix) {
				typeDeclare = gt[len(typePrefix):]
			} else if strings.HasPrefix(gtLower, defaultPrefix) {
				defaultDeclare = "DEFAULT " + gt[len(defaultPrefix):]
			} else if strings.HasPrefix(gtLower, indexPrefix) {
				indexDeclare = gt[len(indexPrefix):]
				if idxFields := strings.Split(indexDeclare, ","); len(idxFields) > 1 {
					idxQuotes := make([]string, len(idxFields))
					for i := range idxFields {
						idxQuotes[i] = fmt.Sprintf("`%s`", idxFields[i])
					}
					indexDeclare = fmt.Sprintf("`idx_%s` ON `%s`(%s);", strings.Join(idxFields, "_"), tableName, strings.Join(idxQuotes, ", "))
				} else {
					indexDeclare = fmt.Sprintf("`%s` ON `%s`(`%s`);", indexDeclare, tableName, columnDeclare)
				}
			} else if gtLower == isPrimaryKey {
				isPkDeclare = true
			} else if gtLower == isUnique {
				isUniqueDeclare = true
			} else if gtLower == isAutoIncrement {
				isAutoDeclare = true
			}
		}

		if isFkDeclare {
			continue
		}

		if indexDeclare != "" {
			if isUniqueDeclare {
				indexes = append(indexes, "CREATE UNIQUE INDEX "+indexDeclare)
			} else {
				indexes = append(indexes, "CREATE INDEX "+indexDeclare)
			}
		}

		if len(columnDeclare) > maxLen {
			maxLen = len(columnDeclare)
		}

		fs := []string{columnDeclare}
		if typeDeclare != "" {
			fs = append(fs, typeDeclare)
		} else {
			fs = append(fs, sqlType(v.Field(j).Interface(), ""))
		}
		if defaultDeclare != "" {
			fs = append(fs, defaultDeclare)
		}
		if isAutoDeclare {
			fs = append(fs, "AUTO_INCREMENT")
		}
		if isPkDeclare {
			fs = append(fs, "PRIMARY KEY")
		}

		fields = append(fields, fs)
	}

	fs := []string{}
	for _, f := range fields {
		fs = append(fs, fmt.Sprintf("  `%s`%s%s", f[0], strings.Repeat(" ", maxLen-len(f[0])+1), strings.Join(f[1:], " ")))
	}

	sql := []string{fmt.Sprintf("CREATE TABLE `%s`(\n%s\n);", tableName, strings.Join(fs, ",\n"))}
	sql = append(sql, indexes...)
	return strings.Join(sql, "\n")
}

func SqlDropTable(tb interface{}) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", getTableName(tb))
}

func sqlType(v interface{}, suffix string) string {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		vv := reflect.Indirect(reflect.ValueOf(v))
		if vv.IsZero() {
			return "UNSPECIFIED"
		}
		return sqlType(vv.Interface(), "NULL")
	}

	if suffix != "" {
		suffix = " " + suffix
	}
	switch v.(type) {
	case bool:
		return "BOOLEAN" + suffix
	case int8, uint8:
		return "TINYINT" + suffix
	case int16, uint16:
		return "SMALLINT" + suffix
	case int, int32, uint32:
		return "INT" + suffix
	case int64, uint64:
		return "BIGINT" + suffix
	case float32:
		return "FLOAT" + suffix
	case float64:
		return "DOUBLE" + suffix
	case string:
		return "TEXT" + suffix
	case time.Time:
		return "DATETIME" + suffix
	default:
		return "UNSPECIFIED"
	}
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
