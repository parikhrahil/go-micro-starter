package logger

import (
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type zLog struct {
	log    zerolog.Logger
	fields []Field
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

func (l *zLog) WithFields(fields ...Field) {
	if len(fields) > 0 {
		l.fields = fields
	}
}

func (l *zLog) Info(msg string, fields ...Field) {
	loglevel(l.log.Info(), msg, append(l.fields, fields...))
}

func (l *zLog) Warn(msg string, fields ...Field) {
	loglevel(l.log.Warn(), msg, append(l.fields, fields...))
}

func (l *zLog) Debug(msg string, fields ...Field) {
	loglevel(l.log.Debug(), msg, append(l.fields, fields...))
}

func (l *zLog) Error(msg string, fields ...Field) {
	loglevel(l.log.Error(), msg, append(l.fields, fields...))
}

func (l *zLog) Fatal(msg string, fields ...Field) {
	loglevel(l.log.Fatal(), msg, append(l.fields, fields...))
}

func (l *zLog) Panic(msg string, fields ...Field) {
	loglevel(l.log.Panic(), msg, append(l.fields, fields...))
}

func loglevel(e *zerolog.Event, msg string, fields []Field) {
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}
