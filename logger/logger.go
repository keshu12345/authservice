package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/keshucs12345/authservice/config"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Fields = map[string]interface{}

type LogEntry interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugln(...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	Trace(...interface{})
	Tracef(string, ...interface{})
	Traceln(...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})
	WithField(string, interface{}) LogEntry
	WithFields(Fields) LogEntry
}

type Logger interface {
	LogEntry
	NewEntry() LogEntry
	AddHook(interface{}) error
	SetFormatter(string)
	SetLevel(string)
	SetOutput(io.Writer)
	SetReportCaller(bool)
}

var (
	instance Logger
	once     sync.Once
)

type Config struct {
	LogDir       string
	Level        string
	Format       string
	ReportCaller bool
}



func logrusLoggerImpl() Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		ForceColors:     true,
		DisableQuote:    true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	return &logrusLogger{logger: l}
}

func New(cfg *config.Configuration) Logger {
	once.Do(func() {
		logger := logrusLoggerImpl()

		if cfg.Logger.Level == "" {
			cfg.Logger.Level = "info"
		}
		if cfg.Logger.Format == "" {
			cfg.Logger.Format = "json"
		}

		logger.SetLevel(cfg.Logger.Level)
		logger.SetReportCaller(cfg.Logger.ReportCaller)

		date := time.Now().Format("2006-01-02")
		pid := os.Getpid()
		logPath := filepath.Join(cfg.Logger.LogDir, date)
		_ = os.MkdirAll(logPath, 0755)
		filename := filepath.Join(logPath, fmt.Sprintf("authservice-%d.log", pid))

		rotatedWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.Logger.MaxSize,
			MaxAge:     cfg.Logger.MaxAge,
			MaxBackups: cfg.Logger.MaxBackups,
			LocalTime:  cfg.Logger.LocalTime,
			Compress:   cfg.Logger.Compress,
		}

		// Split output between file and stdout
		multiWriter := io.MultiWriter(os.Stdout, rotatedWriter)
		logger.SetOutput(multiWriter)

		// Set different formatter for terminal and file (optional)
		logger.SetFormatter(cfg.Logger.Format)

		instance = logger
	})

	return instance
}

func (l *logrusEntry) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logrusEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusEntry) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l *logrusEntry) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logrusEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusEntry) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l *logrusEntry) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logrusEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusEntry) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *logrusEntry) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logrusEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusEntry) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *logrusEntry) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *logrusEntry) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logrusEntry) Panicln(args ...interface{}) {
	l.entry.Panicln(args...)
}

func (l *logrusEntry) Print(args ...interface{}) {
	l.entry.Print(args...)
}

func (l *logrusEntry) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logrusEntry) Println(args ...interface{}) {
	l.entry.Println(args...)
}

func (l *logrusEntry) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

func (l *logrusEntry) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args...)
}

func (l *logrusEntry) Traceln(args ...interface{}) {
	l.entry.Traceln(args...)
}

func (l *logrusEntry) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logrusEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusEntry) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

func (l *logrusEntry) WithField(key string, value interface{}) LogEntry {
	return &logrusEntry{entry: l.entry.WithField(key, value)}
}

func (l *logrusEntry) WithFields(fields Fields) LogEntry {
	return &logrusEntry{entry: l.entry.WithFields(fields)}
}

type logrusLogger struct {
	logger *logrus.Logger
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l *logrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *logrusLogger) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}

func (l *logrusLogger) Print(args ...interface{}) {
	l.logger.Print(args...)
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *logrusLogger) Println(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *logrusLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *logrusLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *logrusLogger) Traceln(args ...interface{}) {
	l.logger.Traceln(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

func (l *logrusLogger) WithField(key string, value interface{}) LogEntry {
	return l.NewEntry().WithField(key, value)
}

func (l *logrusLogger) WithFields(fields Fields) LogEntry {
	return l.NewEntry().WithFields(fields)
}

func (l *logrusLogger) SetLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		l.logger.SetLevel(logrus.InfoLevel)
	} else {
		l.logger.SetLevel(lvl)
	}
}

func (l *logrusLogger) SetFormatter(formatter string) {
	switch formatter {
	case "text":
		l.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
			ForceColors:     true,
			DisableQuote:    true,
		})
	case "prefixed-text":
		l.logger.SetFormatter(&prefixed.TextFormatter{FullTimestamp: true})
	default:
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}
}

func (l *logrusLogger) SetReportCaller(v bool)      { l.logger.SetReportCaller(v) }
func (l *logrusLogger) SetOutput(w io.Writer)       { l.logger.SetOutput(w) }
func (l *logrusLogger) AddHook(h interface{}) error { return nil }
func (l *logrusLogger) NewEntry() LogEntry          { return &logrusEntry{entry: logrus.NewEntry(l.logger)} }

type logrusEntry struct {
	entry *logrus.Entry
}
