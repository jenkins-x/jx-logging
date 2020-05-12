// +build unit

package log

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_debug_log_is_written_to_output_when_corresponding_level_is_set(t *testing.T) {
	err := SetLevel("info")
	assert.NoError(t, err)

	out := CaptureOutput(func() { Logger().Debug("hello") })
	assert.Empty(t, out)

	err = SetLevel("debug")
	assert.NoError(t, err)

	out = CaptureOutput(func() { Logger().Debug("hello") })
	assert.Equal(t, "DEBUG: hello\n", out)
}

func Test_setting_unknown_log_level_returns_error(t *testing.T) {
	err := SetLevel("foo")
	assert.Error(t, err)
	assert.Equal(t, "Invalid log level 'foo'", err.Error())
}

func Test_debug_log_is_written_to_output_when_env_var_is_set(t *testing.T) {
	err := SetLevel("info")
	assert.NoError(t, err)

	assert.Equal(t, "info", GetLevel())
	out := CaptureOutput(func() { Logger().Debug("hello") })
	assert.Empty(t, out)

	out = CaptureOutput(func() { Logger().Info("hello") })
	assert.Equal(t, "hello\n", out)

	out = CaptureOutput(func() { Logger().Warn("hello") })
	assert.Equal(t, "WARNING: hello\n", out)

	out = CaptureOutput(func() { Logger().Error("hello") })
	assert.Equal(t, "ERROR: hello\n", out)

	os.Setenv("JX_LOG_LEVEL", "debug")
	err = forceInitLogger()
	assert.NoError(t, err)

	assert.Equal(t, "debug", GetLevel())
	out = CaptureOutput(func() { Logger().Debug("hello") })
	assert.Equal(t, "DEBUG: hello\n", out)
}

func Test_Json_log_formatter(t *testing.T) {
	os.Setenv("JX_LOG_FORMAT", "json")
	err := forceInitLogger()
	assert.NoError(t, err)

	out := CaptureOutput(func() { Logger().Infof("hello") })
	t.Logf(out)
	assert.Equal(t, strings.HasPrefix(out, "{"), true)
	assert.Equal(t, strings.HasSuffix(out, "}\n"), true)
	assert.Equal(t, strings.Contains(out, `"level":"info"`), true)
}

func Test_Stackdriver_log_formatter(t *testing.T) {
	os.Setenv("JX_LOG_FORMAT", "stackdriver")
	err := forceInitLogger()
	assert.NoError(t, err)

	out := CaptureOutput(func() { Logger().Infof("hello") })
	t.Logf(out)
	assert.Equal(t, strings.HasPrefix(out, "{"), true)
	assert.Equal(t, strings.HasSuffix(out, "}\n"), true)
	assert.Equal(t, strings.Contains(out, `"severity":"INFO"`), true)
	assert.Equal(t, strings.Contains(out, `"context":{}`), true)
}

func Test_GetLevels(t *testing.T) {
	Logger()
	levels := GetLevels()
	assert.Equal(t, "panic fatal error warning info debug trace", strings.Join(levels, " "))
}


