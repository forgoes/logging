package handler

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgoes/logging"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestJsonFormatter_MergeLoggerAndEventContext(t *testing.T) {
	logger := logging.GetLogger("formatter.json").Clone()
	logger.SetTag("service", "billing")
	logger.SetTag("request_id", "logger")
	logger.SetKv("region", "ca")
	logger.SetKv("request_id", "logger")

	e := logger.Info()
	e.Tag("request_id", "event")
	e.Tag("step", "db")
	e.Kv("request_id", "event")
	e.Kv("user_id", 42)

	b, err := JsonFormatter(e)
	assert.Nil(t, err)

	var payload map[string]interface{}
	err = json.Unmarshal(b, &payload)
	assert.Nil(t, err)

	tags := payload["tag"].(map[string]interface{})
	assert.Equal(t, "billing", tags["service"])
	assert.Equal(t, "event", tags["request_id"])
	assert.Equal(t, "db", tags["step"])

	kvs := payload["kvs"].(map[string]interface{})
	assert.Equal(t, "ca", kvs["region"])
	assert.Equal(t, "event", kvs["request_id"])
	assert.Equal(t, float64(42), kvs["user_id"])
}

func TestStdFormatter_MergeLoggerAndEventContext(t *testing.T) {
	logger := logging.GetLogger("formatter.std").Clone()
	logger.SetTag("service", "billing")
	logger.SetTag("request_id", "logger")
	logger.SetKv("region", "ca")
	logger.SetKv("request_id", "logger")

	e := logger.Info()
	e.Tag("request_id", "event")
	e.Kv("request_id", "event")
	e.Kv("user_id", 42)

	b, err := StdFormatter(e)
	assert.Nil(t, err)

	s := ansiPattern.ReplaceAllString(string(b), "")
	assert.True(t, strings.Contains(s, "TAG"))
	assert.True(t, strings.Contains(s, "request_id=event"))
	assert.True(t, strings.Contains(s, "service=billing"))
	assert.True(t, strings.Contains(s, "KVS"))
	assert.True(t, strings.Contains(s, "region=ca"))
	assert.True(t, strings.Contains(s, "user_id=42"))
}

func TestStdFormatter_RendersCliStyleSections(t *testing.T) {
	logger := logging.GetLogger("formatter.cli").Clone()
	logger.SetTag("service", "billing")

	e := logger.Error()
	e.Kv("request_id", "req-1")
	e.E(assert.AnError)
	e.Attach("ignored")

	b, err := StdFormatter(e)
	assert.Nil(t, err)

	s := ansiPattern.ReplaceAllString(string(b), "")
	assert.True(t, strings.Contains(s, "ERROR"))
	assert.True(t, strings.Contains(s, "formatter.cli"))
	assert.True(t, strings.Contains(s, "err"))
	assert.True(t, strings.Contains(s, assert.AnError.Error()))
	assert.True(t, strings.Contains(s, "TAG"))
	assert.True(t, strings.Contains(s, "service=billing"))
	assert.True(t, strings.Contains(s, "KVS"))
	assert.True(t, strings.Contains(s, "request_id=req-1"))
}

func TestJsonFormatter_RejectsNonStringContextKeys(t *testing.T) {
	logger := logging.GetLogger("formatter.json.invalid").Clone()

	e := logger.Info()
	e.Tag(1, "bad-key")

	_, err := JsonFormatter(e)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "only supports string context keys"))
}
