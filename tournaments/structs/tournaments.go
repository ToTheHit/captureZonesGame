package structs

import (
	"main/categories"
	"main/questions"
)

type Winner struct {
	ID     int `db:"id" json:"id"`
	Points int `db:"points" json:"points"`
	Zones  int `db:"zones" json:"zones"`
}

type Tournament struct {
	ID         int                  `db:"id"`
	CreatedOn  int64                `db:"_createdon"`
	Name       string               `db:"name" json:"name"` // title
	CategoryID *categories.Category `db:"category_id" json:"categoryId"`
	//MapTemplateId
	MaxPlayers   int                   `db:"max_players" json:"maxPlayers"`
	Status       string                `db:"status" json:"status"`
	Stage        int                   `db:"stage" json:"stage"`
	Questions    []*questions.Question `db:"questions" json:"questions"`
	QuestionId   int                   `db:"question_id" json:"questionId"`
	StartsAt     int64                 `db:"starts_at" json:"startsAt"`
	Winners      []*Winner             `db:"winners" json:"winners"`
	IsPublished  bool                  `db:"is_published" json:"isPublished"`
	IsSingleGame bool                  `db:"is_single_game" json:"isSingleGame"`
}

var TournamentSchema = `
	CREATE TABLE IF NOT EXISTS tournaments (
		id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		_createdOn TIMESTAMP,
		name text,
		category_id int,
		max_players int, 
		status text,
		stage int,
		questions JSON[],
		question_id int,
		starts_at TIMESTAMP,
		winners JSON[],
		is_published boolean,
		is_single_game boolean				
	);
`
