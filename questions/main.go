package questions

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/categories"
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

type Answer struct {
	IsCorrect bool   `json:"isCorrect"`
	Text      string `json:"text"`
}

type Question struct {
	ID         int                  `db:"id"`
	CreatedOn  int64                `db:"_createdon"`
	UpdatedOn  *int64               `db:"_updatedon"`
	CategoryID *categories.Category `db:"category_id" json:"categoryId"`
	Title      string               `db:"title" json:"title"`
	Intro      string               `db:"intro" json:"intro"`
	Text       string               `db:"text" json:"text"`
	Answers    []*Answer            `db:"answers" json:"answers"`
	//Image      string               `db:"image" json:"image"`
}

var questionSchema = `
	CREATE TABLE IF NOT EXISTS questions (
		id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		_createdOn TIMESTAMP,
		_updatedOn TIMESTAMP,
		category_id int, 
		title text,
		intro text,
		text text,
		answers JSON[]
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
	var cfg config
	cfg.db.dsn = env.GetString("DB_DSN", "admin:pass@localhost:5432/main?sslmode=disable")

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// TODO: ==== move to file =====
	db.MustExec(questionSchema)

	app := &application{
		db:     db,
		logger: logger,
	}

	return app.routes(), nil
}
