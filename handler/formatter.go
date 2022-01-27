package handler

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/forgoes/logging"
)

type Formatter func(e *logging.Event) ([]byte, error)

func JsonFormatter(e *logging.Event) ([]byte, error) {
	m := make(map[string]interface{})
	m["logger"] = e.GetLogger().Name()
	if e.GetTags() != nil {
		m["tag"] = e.GetTags()
	}
	if e.GetKvs() != nil {
		m["kvs"] = e.GetKvs()
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
	m["msg"] = e.GetMsg()
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
	switch level {
	case logging.DEBUG:
		buf.WriteString(White + level.String() + "  " + Reset)
	case logging.INFO:
		buf.WriteString(Green + level.String() + "   " + Reset)
	case logging.WARN:
		buf.WriteString(Purple + level.String() + "   " + Reset)
	case logging.ERROR:
		fallthrough
	case logging.PANIC:
		fallthrough
	case logging.FATAL:
		buf.WriteString(Red + level.String() + "  " + Reset)
	default:
		buf.WriteString(Gray + "UNKNOWN" + Reset)
	}

	buf.WriteString(Blue)
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buf.WriteString(Reset)

	buf.WriteString(" ")

	buf.WriteString(GrayBlack)
	if e.GetLogger().Name() == "" {
		buf.WriteString("root")
	} else {
		buf.WriteString(e.GetLogger().Name())
	}
	buf.WriteString(Reset)

	if e.GetCaller().GetOK() {
		buf.WriteString(" ")
		buf.WriteString(YellowBlue)
		buf.WriteString(path.Base(e.GetCaller().GetFile()))
		buf.WriteString(Reset)
		buf.WriteString(" ")
		f := strings.Split(e.GetCaller().GetFunc(), ".")
		if len(f) > 0 {
			buf.WriteString(BlackGray)
			buf.WriteString(f[len(f)-1])
			buf.WriteString(Reset)
			buf.WriteString(" ")
		}
		buf.WriteString(Red)
		buf.WriteString("[")
		buf.WriteString(strconv.Itoa(e.GetCaller().GetLine()))
		buf.WriteString("]")
		buf.WriteString(Reset)
	}

	buf.WriteString(" ")

	buf.WriteString(Cyan)
	buf.WriteString(e.GetMsg())
	buf.WriteString(Reset)

	if len(e.GetStack()) > 0 {
		buf.WriteString("\n")
		buf.WriteString(e.GetStack())
	}

	return buf.Bytes(), nil
}
