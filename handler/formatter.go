package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/forgoes/logging"
)

type Formatter func(e *logging.Event) ([]byte, error)

type levelStyle struct {
	badge string
	text  string
}

type contextEntry struct {
	key   string
	value interface{}
}

func normalizeJSONContext(values map[interface{}]interface{}) (map[string]interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	normalized := make(map[string]interface{}, len(values))
	for k, v := range values {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("json formatter only supports string context keys, got %T", k)
		}
		normalized[key] = v
	}

	return normalized, nil
}

func sortedContextEntries(values map[interface{}]interface{}) []contextEntry {
	if len(values) == 0 {
		return nil
	}

	entries := make([]contextEntry, 0, len(values))
	for k, v := range values {
		entries = append(entries, contextEntry{
			key:   fmt.Sprint(k),
			value: v,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	return entries
}

func writeContextBlock(buf *bytes.Buffer, label string, labelColor string, keyColor string, valueColor string, values map[interface{}]interface{}) {
	entries := sortedContextEntries(values)
	if len(entries) == 0 {
		return
	}

	buf.WriteString(" ")
	buf.WriteString(labelColor)
	buf.WriteString(" ")
	buf.WriteString(label)
	buf.WriteString(" ")
	buf.WriteString(Reset)

	for _, entry := range entries {
		buf.WriteString(" ")
		buf.WriteString(keyColor)
		buf.WriteString(entry.key)
		buf.WriteString(Reset)
		buf.WriteString("=")
		buf.WriteString(valueColor)
		buf.WriteString(fmt.Sprint(entry.value))
		buf.WriteString(Reset)
	}
}

func stdLevelStyle(level logging.Level) levelStyle {
	switch level {
	case logging.DEBUG:
		return levelStyle{
			badge: BlackWhite + " DEBUG " + Reset,
			text:  White,
		}
	case logging.INFO:
		return levelStyle{
			badge: GreenPurple + " INFO " + Reset,
			text:  Green,
		}
	case logging.WARN:
		return levelStyle{
			badge: BlueYellow + " WARN " + Reset,
			text:  Yellow,
		}
	case logging.ERROR:
		return levelStyle{
			badge: RedCyan + " ERROR " + Reset,
			text:  Red,
		}
	case logging.PANIC:
		return levelStyle{
			badge: CyanRed + " PANIC " + Reset,
			text:  Red,
		}
	case logging.FATAL:
		return levelStyle{
			badge: RedCyan + " FATAL " + Reset,
			text:  Red,
		}
	default:
		return levelStyle{
			badge: GrayBlack + " UNKNOWN " + Reset,
			text:  Gray,
		}
	}
}

func writeSeparator(buf *bytes.Buffer) {
	buf.WriteString(" ")
	buf.WriteString(Faint)
	buf.WriteString("|")
	buf.WriteString(Reset)
	buf.WriteString(" ")
}

func writeMessageBlock(buf *bytes.Buffer, level logging.Level, msg string) {
	style := stdLevelStyle(level)
	buf.WriteString(style.text)
	buf.WriteString(Bold)
	buf.WriteString(msg)
	buf.WriteString(Reset)
}

func writeErrorBlock(buf *bytes.Buffer, err error) {
	if err == nil {
		return
	}

	writeSeparator(buf)
	buf.WriteString(CyanRed)
	buf.WriteString(" err ")
	buf.WriteString(Reset)
	buf.WriteString(" ")
	buf.WriteString(Red)
	buf.WriteString(err.Error())
	buf.WriteString(Reset)
}

func formatCaller(e *logging.Event) string {
	if !e.GetCaller().GetOK() {
		return ""
	}

	parts := make([]string, 0, 3)
	parts = append(parts, path.Base(e.GetCaller().GetFile()))
	f := strings.Split(e.GetCaller().GetFunc(), ".")
	if len(f) > 0 {
		parts = append(parts, f[len(f)-1])
	}
	parts = append(parts, strconv.Itoa(e.GetCaller().GetLine()))

	return strings.Join(parts, ":")
}

func writeStackBlock(buf *bytes.Buffer, stack string) {
	if len(stack) == 0 {
		return
	}

	lines := strings.Split(strings.TrimRight(stack, "\n"), "\n")
	buf.WriteString("\n")
	buf.WriteString(Faint)
	buf.WriteString("stack:")
	buf.WriteString(Reset)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		buf.WriteString("\n")
		buf.WriteString(Faint)
		buf.WriteString("  ")
		buf.WriteString(strings.TrimSpace(line))
		buf.WriteString(Reset)
	}
}

func JsonFormatter(e *logging.Event) ([]byte, error) {
	m := make(map[string]interface{})
	m["logger"] = e.GetLogger().Name()
	if e.GetTags() != nil {
		tags, err := normalizeJSONContext(e.GetTags())
		if err != nil {
			return nil, err
		}
		m["tag"] = tags
	}
	if e.GetKvs() != nil {
		kvs, err := normalizeJSONContext(e.GetKvs())
		if err != nil {
			return nil, err
		}
		m["kvs"] = kvs
	}
	m["time"] = e.GetTime().Format("2006-01-02 15:04:05")
	m["level"] = e.GetLevel().String()
	m["caller"] = map[string]interface{}{
		"ok":   e.GetCaller().GetOK(),
		"pc":   e.GetCaller().GetPC(),
		"file": e.GetCaller().GetFile(),
		"func": e.GetCaller().GetFunc(),
		"line": e.GetCaller().GetLine(),
	}
	msg := e.GetMsg()
	args := e.GetArgs()
	if len(args) > 0 {
		m["msg"] = fmt.Sprintf(msg, args...)
	} else {
		m["msg"] = msg
	}
	m["stack"] = e.GetStack()
	m["extra"] = e.GetExtra()

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func StdFormatter(e *logging.Event) ([]byte, error) {
	b := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(b)

	level := e.GetLevel()
	style := stdLevelStyle(level)
	buf.WriteString(style.badge)
	buf.WriteString(" ")
	buf.WriteString(Faint)
	buf.WriteString(e.GetTime().Format("2006-01-02 15:04:05"))
	buf.WriteString(Reset)

	writeSeparator(buf)
	buf.WriteString(BlackGray)
	if e.GetLogger().Name() == "" {
		buf.WriteString("root")
	} else {
		buf.WriteString(e.GetLogger().Name())
	}
	buf.WriteString(Reset)

	if caller := formatCaller(e); caller != "" {
		writeSeparator(buf)
		buf.WriteString(YellowBlue)
		buf.WriteString(" caller ")
		buf.WriteString(Reset)
		buf.WriteString(" ")
		buf.WriteString(Yellow)
		buf.WriteString(caller)
		buf.WriteString(Reset)
	}

	msg := e.GetMsg()
	args := e.GetArgs()
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	writeSeparator(buf)
	writeMessageBlock(buf, level, msg)

	writeErrorBlock(buf, e.GetError())

	writeContextBlock(buf, "TAG", BlueYellow, Yellow, White, e.GetTags())
	writeContextBlock(buf, "KVS", GreenPurple, Cyan, White, e.GetKvs())

	writeStackBlock(buf, e.GetStack())

	return buf.Bytes(), nil
}
