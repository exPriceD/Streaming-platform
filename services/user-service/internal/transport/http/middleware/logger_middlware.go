package middleware

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Skipper func(c echo.Context) bool

type LoggerMiddlewareConfig struct {
	Logger       *slog.Logger // Исходный логгер (опционально, для наследования атрибутов)
	ConsoleLevel slog.Level   // Уровень логирования для консоли
	FileLevel    slog.Level   // Уровень логирования для файла
	LogFilePath  string       // Путь к файлу логов
	Skipper      Skipper      // Пропускать определённые запросы
}

var DefaultLoggerConfig = LoggerMiddlewareConfig{
	ConsoleLevel: slog.LevelInfo,
	FileLevel:    slog.LevelDebug,
	LogFilePath:  "logs/requests.log",
	Skipper:      func(c echo.Context) bool { return false },
}

func NewLoggerMiddleware(cfg LoggerMiddlewareConfig) echo.MiddlewareFunc {
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultLoggerConfig.Skipper
	}

	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.ConsoleLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				level := a.Value.Any().(slog.Level)
				switch {
				case level >= slog.LevelError:
					return slog.String(a.Key, color.RedString("ERROR"))
				case level >= slog.LevelWarn:
					return slog.String(a.Key, color.YellowString("WARN"))
				case level >= slog.LevelInfo:
					return slog.String(a.Key, color.GreenString("INFO"))
				default:
					return slog.String(a.Key, color.CyanString("DEBUG"))
				}
			case slog.TimeKey:
				return slog.String(a.Key, a.Value.Time().Format("15:04:05.000"))
			}
			return a
		},
	})

	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: cfg.FileLevel,
	})

	multiHandler := NewMultiHandler(consoleHandler, fileHandler)
	logger := slog.New(multiHandler).WithGroup("http")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			err := next(c)

			req := c.Request()
			res := c.Response()
			status := res.Status
			method := req.Method
			uri := req.URL.RequestURI()
			duration := time.Since(start)
			remoteIP := c.RealIP()
			logErr := ""
			if err != nil {
				logErr = err.Error()
			}

			if errMsg, ok := c.Get("error_message").(string); ok && logErr == "" {
				logErr = errMsg
			}

			msg, _ := c.Get("log_message").(string)
			if msg == "" {
				msg = "Request processed"
			}

			var level slog.Level
			switch {
			case status >= 500:
				level = slog.LevelError
			case status >= 400:
				level = slog.LevelWarn
			case status >= 200:
				level = slog.LevelInfo
			default:
				level = slog.LevelDebug
			}

			logger.LogAttrs(context.Background(), level, msg,
				slog.String("method", method),
				slog.String("uri", uri),
				slog.Int("status", status),
				slog.String("error", logErr),
				slog.Duration("duration", duration),
				slog.String("remote_ip", remoteIP),
			)

			return err
		}
	}
}

func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return multiHandler(handlers)
}

type multiHandler []slog.Handler

func (h multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h))
	for i, handler := range h {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return multiHandler(newHandlers)
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h))
	for i, handler := range h {
		newHandlers[i] = handler.WithGroup(name)
	}
	return multiHandler(newHandlers)
}
