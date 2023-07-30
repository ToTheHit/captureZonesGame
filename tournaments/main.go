package tournaments

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/internal/database"
	"main/internal/env"
	"net/http"
)

type application struct {
	db     *database.DB
	logger *zap.Logger
}

type config struct {
	db struct {
		dsn string
	}
}

func BuildRoutes() (http.Handler, error) {
	loggerConfig := zap.NewDevelopmentEncoderConfig()
	loggerConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(loggerConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	var cfg config
	cfg.db.dsn = env.GetString("DB_DSN", "admin:pass@localhost:5432/main?sslmode=disable")

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	app := &application{
		db:     db,
		logger: logger,
	}

	app.initDatabase()

	return app.routes(), nil
}
