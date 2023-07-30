package tournaments

import "main/tournaments/structs"

func (app *application) initDatabase() {
	app.db.MustExec(structs.TournamentSchema)
	app.db.MustExec(structs.PlayerSchema)
}
