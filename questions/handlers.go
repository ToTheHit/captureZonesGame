package questions

import (
	"go.uber.org/zap"
	"main/internal/request"
	"main/internal/response"
	"net/http"
	"strconv"
	"time"
)

func (app *application) getQuestionsListHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getQuestionHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var reqData Question
	err := request.DecodeJSON(w, r, &reqData)
	if err != nil {
		//app.badRequest(w, r, err)
		panic(err)

		return
	}

	reqData.CreatedOn = time.Now().UnixMilli()

	sqlStr := `insert into questions (title, intro, text, answers, _createdon) values ($1, $2, $3, $4, $5) returning ID`

	var id int
	err = app.db.QueryRow(sqlStr, reqData.Title, reqData.Intro, reqData.Text, reqData.Answers, reqData.CreatedOn).Scan(&id)
	if err != nil {
		app.logger.Error("insert failed, err: " + err.Error())
		panic(err)
	}

	app.logger.Debug("insert success: ",
		zap.String("id", strconv.Itoa(id)),
	)

	errResponse := response.JSON(w, http.StatusOK, reqData)
	if errResponse != nil {
		panic(errResponse)
	}
}

func (app *application) updateQuestionHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteQuestionHandler(w http.ResponseWriter, r *http.Request) {

}
