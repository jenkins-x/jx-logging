package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_CanSetLevel(t *testing.T) {
	l, err := forceInitLogger()
	assert.NoError(t, err)

	err = SetLevel("info")
	assert.NoError(t, err)

	out := CaptureOutput(func() { l.Debug("hello") })
	assert.Empty(t, out)

	out = CaptureOutput(func() { l.Warn("hello") })
	assert.Equal(t, "WARNING: hello\n", out)

	out = CaptureOutput(func() { l.Error("hello") })
	assert.Equal(t, "ERROR: hello\n", out)

	err = SetLevel("debug")
	assert.NoError(t, err)

	out = CaptureOutput(func() { l.Debug("hello") })
	assert.Equal(t, "DEBUG: hello\n", out)
}

func TestLogger_CannotSetInvalidLevel(t *testing.T) {
	_, err := forceInitLogger()
	assert.NoError(t, err)

	err = SetLevel("foo")
	assert.Error(t, err)
	assert.Equal(t, "Invalid log level 'foo'", err.Error())
}

func TestLogger_CanHandleLogLevelEnvVar(t *testing.T) {
	l, err := forceInitLogger()
	assert.NoError(t, err)

	err = SetLevel("info")
	assert.NoError(t, err)

	assert.Equal(t, "info", GetLevel())
	out := CaptureOutput(func() { l.Debug("hello") })
	assert.Empty(t, out)

	out = CaptureOutput(func() { l.Info("hello") })
	assert.Equal(t, "hello\n", out)

	os.Setenv("JX_LOG_LEVEL", "debug")
	defer os.Unsetenv("JX_LOG_LEVEL")

	l, err = forceInitLogger()
	assert.NoError(t, err)

	assert.Equal(t, "debug", GetLevel())
	out = CaptureOutput(func() { l.Debug("hello") })
	assert.Equal(t, "DEBUG: hello\n", out)
}

func TestLogger_JsonLogFormatter(t *testing.T) {
	os.Setenv("JX_LOG_FORMAT", "json")
	defer os.Unsetenv("JX_LOG_FORMAT")

	l, err := forceInitLogger()
	assert.NoError(t, err)

	out := CaptureOutput(func() { l.Infof("hello") })

	assert.Equal(t, strings.HasPrefix(out, "{"), true)
	assert.Equal(t, strings.HasSuffix(out, "}\n"), true)
	assert.Contains(t, out, `"level":"info"`)
}

func Test_Stackdriver_log_formatter(t *testing.T) {
	os.Setenv("JX_LOG_FORMAT", "stackdriver")
	defer os.Unsetenv("JX_LOG_FORMAT")

	l, err := forceInitLogger()
	assert.NoError(t, err)

	out := CaptureOutput(func() { l.Infof("hello") })
	t.Logf(out)
	assert.Equal(t, strings.HasPrefix(out, "{"), true)
	assert.Equal(t, strings.HasSuffix(out, "}\n"), true)
	assert.Contains(t, out, `"severity":"INFO"`)
	assert.Contains(t, out, `"context":{}`)
}

func Test_CanLogToFile(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "test1.*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	t.Logf("got file %s", file.Name())

	os.Setenv("JX_LOG_FORMAT", "text")
	defer os.Unsetenv("JX_LOG_FORMAT")
	os.Setenv("JX_LOG_FILE", file.Name())
	defer os.Unsetenv("JX_LOG_FILE")

	l, err := forceInitLogger()
	assert.NoError(t, err)

	err = SetLevel("debug")

	assert.NoError(t, err)
	assert.Equal(t, "debug", GetLevel())

	out := CaptureOutput(func() { l.Debug("debug-hello") })
	t.Logf("Out>: '%s'", out)
	assert.Equal(t, "DEBUG: debug-hello\n", out)

	b, err := ioutil.ReadFile(file.Name())
	assert.NoError(t, err)
	t.Logf("File>: '%s'", b)
	assert.Contains(t, string(b), `"msg":"debug-hello"`)

	out = CaptureOutput(func() { l.Info("info-hello") })
	assert.Equal(t, "info-hello\n", out)
	t.Logf("Out>: '%s'", out)

	b, err = ioutil.ReadFile(file.Name())
	assert.NoError(t, err)
	t.Logf("File>: '%s'", b)
	assert.Contains(t, string(b), `"msg":"debug-hello"`)
	assert.Contains(t, string(b), `"msg":"info-hello"`)

	out = CaptureOutput(func() { l.Warn("warning-hello") })
	assert.Equal(t, "WARNING: warning-hello\n", out)
	t.Logf("Out>: '%s'", out)

	b, err = ioutil.ReadFile(file.Name())
	assert.NoError(t, err)
	t.Logf("File>: '%s'", b)
	assert.Contains(t, string(b), `"msg":"debug-hello"`)
	assert.Contains(t, string(b), `"msg":"info-hello"`)
	assert.Contains(t, string(b), `"msg":"warning-hello"`)

	out = CaptureOutput(func() { l.Error("error-hello") })
	assert.Equal(t, "ERROR: error-hello\n", out)
	t.Logf("Out>: '%s'", out)

	b, err = ioutil.ReadFile(file.Name())
	assert.NoError(t, err)
	t.Logf("File>: '%s'", b)
	assert.Contains(t, string(b), `"msg":"debug-hello"`)
	assert.Contains(t, string(b), `"msg":"info-hello"`)
	assert.Contains(t, string(b), `"msg":"warning-hello"`)
	assert.Contains(t, string(b), `"msg":"error-hello"`)
}
