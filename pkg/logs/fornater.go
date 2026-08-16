package logs

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type PrettyFormatter struct {
	ShowTimestamp bool
}

func InitLogger() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&PrettyFormatter{
		ShowTimestamp: true,
	})

	log.SetLevel(log.DebugLevel)
}

func (f *PrettyFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b strings.Builder

	// Время
	if f.ShowTimestamp {
		b.WriteString(entry.Time.Format("15:04:05"))
		b.WriteString("  ")
	}

	// Уровень
	level := strings.ToUpper(entry.Level.String())

	var levelColor string

	switch entry.Level {
	case log.DebugLevel:
		levelColor = "\033[36m" // cyan
	case log.InfoLevel:
		levelColor = "\033[32m" // green
	case log.WarnLevel:
		levelColor = "\033[33m" // yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = "\033[31m" // red
	default:
		levelColor = "\033[37m" // white
	}

	reset := "\033[0m"

	b.WriteString(levelColor)
	b.WriteString(fmt.Sprintf("%-5s", level))
	b.WriteString(reset)
	b.WriteString("  ")

	// Сообщение
	b.WriteString(entry.Message)

	// Поля
	if len(entry.Data) > 0 {
		b.WriteString("  ")

		keys := make([]string, 0, len(entry.Data))
		for key := range entry.Data {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for i, key := range keys {
			if i > 0 {
				b.WriteString(" ")
			}

			b.WriteString(key)
			b.WriteString("=")
			b.WriteString(formatValue(entry.Data[key]))
		}
	}

	b.WriteString("\n")

	return []byte(b.String()), nil
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		if strings.ContainsAny(v, " \t\n\"") {
			return fmt.Sprintf("%q", v)
		}

		return v

	case time.Duration:
		return v.String()

	default:
		return fmt.Sprintf("%v", v)
	}
}
