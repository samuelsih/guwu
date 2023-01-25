package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/samuelsih/guwu/pkg/errs"
)

var (
	logger zerolog.Logger
	once   sync.Once
	debug  bool
)

type P struct {
	Key   string
	Value any
}

func SetMode(debugMode bool) {
	once.Do(func() {
		if debugMode {
			debug = debugMode

			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			output := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.UnixDate,
				FormatLevel: func(i any) string {
					return strings.ToUpper(fmt.Sprintf("[%s]", i))
				},
				FormatMessage: func(i any) string {
					return fmt.Sprintf("<%s>", i)
				},
			}

			logger = zerolog.New(output).With().Timestamp().Logger()
		}
	})
}

func Err(err error) {
	ops := errs.Ops(err)

	logger.Error().Int("status", errs.GetKind(err)).Str("trace", strings.Join(ops, "->")).Msg(err.Error())
}

func Debug(msg string) {
	if debug {
		logger.Debug().Msg(msg)
	}
}

func Debugs(pairs ...P) {
	if debug {
		dict := zerolog.Dict()
		for _, pair := range pairs {
			dict.Interface(pair.Key, pair.Value)
		}

		logger.Debug().Dict("data", dict)
	}
}

func SysInfo(msg string) {
	logger.Info().Msgf(msg)
}

func SysErr(err error) {
	logger.Error().Stack().Err(err).Msg("")
}

func SysFatal(msg string, args ...any) {
	info := fmt.Errorf(msg, args...)
	logger.Fatal().Err(info)
}

func Errorf(msg string, args ...any) {
	result := fmt.Errorf(msg, args...)
	logger.Error().Err(result)
}
