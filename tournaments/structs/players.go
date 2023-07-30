package structs

type Player struct {
	ID           int         `db:"id" json:"id"`
	Name         string      `db:"name" json:"name"`
	Phone        string      `db:"phone" json:"phone"`
	Points       int         `db:"points" json:"points"`
	TotalPoints  int         `db:"total_points" json:"totalPoints"`
	TournamentId *Tournament `db:"tournament_id" json:"tournamentId"`
	GameId       int         `db:"game_id" json:"gameId"` // TODO: gameId
	IsBot        bool        `db:"is_bot" json:"isBot"`
	IsBlocked    bool        `db:"is_blocked" json:"isBlocked"`
	CreatedBy    int         `db:"created_by" json:"createdBy"` // TODO: userId
	Place        int         `db:"place" json:"place"`
	IsDone       bool        `db:"is_done" json:"isDone"`
}

var PlayerSchema = `
	CREATE TABLE IF NOT EXISTS players (
		id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		name text,
		phone text,
		total_points int, 
		tournament_id int,
		game_id int,
		is_bot boolean,
		is_blocked boolean,
		created_by int,
		winners JSON[],
		place int,
		is_done boolean				
	);
`
