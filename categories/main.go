package categories

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/internal/database"
	"main/internal/env"
	"net/http"
)

type config struct {
	db struct {
		dsn string
	}
}

type application struct {
	db     *database.DB
	logger *zap.Logger
}

type Category struct {
	ID           int    `db:"id"`
	CreatedOn    int64  `db:"_createdon"`
	UpdatedOn    *int64 `db:"_updatedon"`
	Title        string `db:"title" json:"title"`
	IsSingleGame bool   `db:"issinglegame" json:"isSingleGame"`
}

var categorySchema = `
	CREATE TABLE IF NOT EXISTS categories (
		id int,
		_createdOn TIMESTAMP,
		_updatedOn TIMESTAMP,
		title text,
		isSingleGame bool
	);
`

func BuildRoutes() (http.Handler, error) {
	loggerConfig := zap.NewDevelopmentEncoderConfig()
	loggerConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(loggerConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	//logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	var cfg config
	cfg.db.dsn = env.GetString("DB_DSN", "admin:pass@localhost:5432/main?sslmode=disable")

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	// TODO: ==== move to file =====
	db.MustExec(categorySchema)

	app := &application{
		db:     db,
		logger: logger,
	}

	return app.routes(), nil
}
