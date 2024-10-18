package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jenkins-x/logrus-stackdriver-formatter/pkg/stackdriver"
	"github.com/sirupsen/logrus"
)

var (
	// colorStatus returns a new function that returns status-colorized (cyan) strings for the
	// given arguments with fmt.Sprint().
	colorStatus = color.New(color.FgCyan).SprintFunc()

	// colorWarn returns a new function that returns status-colorized (yellow) strings for the
	// given arguments with fmt.Sprint().
	colorWarn = color.New(color.FgYellow).SprintFunc()

	// colorError returns a new function that returns error-colorized (red) strings for the
	// given arguments with fmt.Sprint().
	colorError = color.New(color.FgRed).SprintFunc()

	logger *logrus.Entry
)

// FormatLayoutType the layout kind
type FormatLayoutType string

const (
	// FormatLayoutJSON uses JSON layout
	FormatLayoutJSON FormatLayoutType = "json"

	// FormatLayoutText uses classic colorful Jenkins X layout
	FormatLayoutText FormatLayoutType = "text"

	// FormatLayoutStackdriver uses a custom formatter for stackdriver
	FormatLayoutStackdriver FormatLayoutType = "stackdriver"

	JxLogFormat = "JX_LOG_FORMAT"
	JxLogFile   = "JX_LOG_FILE"
	JxLogLevel  = "JX_LOG_LEVEL"
)

func initializeLogger() error {
	if logger == nil {
		_, err := forceInitLogger()
		if err != nil {
			return err
		}
	}
	return nil
}

func forceInitLogger() (*logrus.Entry, error) {
	format := os.Getenv(JxLogFormat)
	if format == "json" {
		setFormatter(FormatLayoutJSON)
	} else if format == "stackdriver" {
		setFormatter(FormatLayoutStackdriver)
	} else {
		setFormatter(FormatLayoutText)
	}

	level := os.Getenv(JxLogLevel)
	if level != "" {
		err := SetLevel(level)
		if err != nil {
			return nil, fmt.Errorf("unable to set level to %s: %w", level, err)
		}
	}

	debugFile := os.Getenv(JxLogFile)
	if debugFile != "" {
		hook := NewHook(debugFile, logrus.AllLevels)
		logrus.AddHook(hook)
	}

	logger = logrus.NewEntry(logrus.StandardLogger())
	return logger, nil
}

// Logger obtains the logger for use in the jx codebase
// This is the only way you should obtain a logger
func Logger() *logrus.Entry {
	err := initializeLogger()
	if err != nil {
		logrus.Warnf("error initializing logrus %v", err)
	}
	return logger
}

// SetLevel sets the logging level
func SetLevel(s string) error {
	level, err := logrus.ParseLevel(s)
	if err != nil {
		return fmt.Errorf("Invalid log level '%s'", s)
	}
	logrus.SetLevel(level)
	return nil
}

// GetLevel gets the current log level
func GetLevel() string {
	return logrus.GetLevel().String()
}

// GetLevels returns the list of valid log levels
func GetLevels() []string {
	var levels []string
	for _, level := range logrus.AllLevels {
		levels = append(levels, level.String())
	}
	return levels
}

// setFormatter sets the logrus format to use either text or JSON formatting
func setFormatter(layout FormatLayoutType) {
	switch layout {
	case FormatLayoutJSON:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case FormatLayoutStackdriver:
		var options []stackdriver.Option
		serviceName := os.Getenv("JX_LOG_SERVICE")
		if serviceName != "" {
			options = append(options, stackdriver.WithService(serviceName))
		}
		version := os.Getenv("JX_LOG_SERVICE_VERSION")
		if version != "" {
			options = append(options, stackdriver.WithVersion(version))
		}
		stackSkip := os.Getenv("JX_LOG_STACK_SKIP")
		if stackSkip != "" {
			values := strings.Split(stackSkip, ",")
			for _, v := range values {
				options = append(options, stackdriver.WithStackSkip(v))
			}
		}
		logrus.SetFormatter(stackdriver.NewFormatter(options...))
	default:
		logrus.SetFormatter(NewJenkinsXTextFormat())
	}
}

// CaptureOutput calls the specified function capturing and returning all logged messages.
func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	f()
	logrus.SetOutput(os.Stdout)
	return buf.String()
}

// SetOutput sets the outputs for the default logger.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}
