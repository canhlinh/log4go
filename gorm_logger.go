package log4go

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

type GormLogger struct {
}

func (this GormLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		source := fmt.Sprintf("%v", values[1])
		lastIndex := strings.LastIndex(source, "/") + 1
		writeGormLog(source[lastIndex:])

		if level == "sql" {
			// duration

			messages := fmt.Sprintf("%s[ %.2fms] ", COLOR_OFF, float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)
			// sql
			var sql string
			var formattedValues []string

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}

			var formattedValuesLength = len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
			messages = messages + sql
			writeGojiLog(messages)
		}
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func writeGormLog(msg string) {
	rec := &LogRecord{
		Level:   DEBUG,
		Created: time.Now(),
		Source:  "GORM",
		Message: msg,
	}
	for _, filt := range Global {
		if filt.Level > rec.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}
