package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type EchoConfig struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

func NewEchoServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	return e
}

func RunHttpServer(ctx context.Context, echo *echo.Echo, cfg *EchoConfig) {
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout
	echo.Server.MaxHeaderBytes = MaxHeaderBytes

	go func() {
		err := echo.Start(cfg.Host + ":" + cfg.Port)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("port", cfg.Port).Msg("shutting down Http server")
				err := echo.Shutdown(ctx)
				if err != nil {
					log.Error().Err(err).Msg("error shutting down Http server")
					return
				}
				log.Info().Msg("server shut down")
				return
			}
		}
	}()
}
