package logger

import (
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type zLog struct {
	log zerolog.Logger
}

type Opts struct {
	LogLevel string    `validate:"required,oneof=trace debug info warn error fatal panic"`
	Writer   io.Writer `validate:"omitempty"`
}

func New(opts *Opts) (Logger, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil {
		return nil, err
	}

	zerolog.SetGlobalLevel(level)

	out := opts.Writer
	if out == nil {
		out = zerolog.ConsoleWriter{Out: os.Stderr}
	}

	log := zerolog.New(out).With().Timestamp().Logger()
	return &zLog{log: log}, nil
}

func (l *zLog) Info(msg string, fields ...Field) {
	loglevel(l.log.Info(), msg, fields...)
}

func (l *zLog) Warn(msg string, fields ...Field) {
	loglevel(l.log.Warn(), msg, fields...)
}

func (l *zLog) Debug(msg string, fields ...Field) {
	loglevel(l.log.Debug(), msg, fields...)
}

func (l *zLog) Error(msg string, fields ...Field) {
	loglevel(l.log.Error(), msg, fields...)
}

func (l *zLog) Fatal(msg string, fields ...Field) {
	loglevel(l.log.Fatal(), msg, fields...)
}

func (l *zLog) Panic(msg string, fields ...Field) {
	loglevel(l.log.Panic(), msg, fields...)
}

func loglevel(e *zerolog.Event, msg string, fields ...Field) {
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}
