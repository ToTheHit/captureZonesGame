package categories

import (
	"go.uber.org/zap"
	"main/internal/request"
	"main/internal/response"
	"net/http"
	"strconv"
	"time"
)

func (app *application) stst(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"Status": "FROM CATEGORIES",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		return
	}
}

type Ani interface {
}

func (app *application) getCategoriesListHandler(w http.ResponseWriter, r *http.Request) {
	categories := []Category{}
	err := app.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		panic(err)
	}

	//response.JSON(w, http.StatusOK, categories)
	data := map[string]Ani{"categories": categories}

	errResponse := response.JSON(w, http.StatusOK, data)
	if errResponse != nil {
		panic(errResponse)
	}
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	//var reqData struct {
	//	Title        string `json:"title"`
	//	IsSingleGame bool   `json:"isSingleGame"`
	//}

	var reqData Category
	err := request.DecodeJSON(w, r, &reqData)
	if err != nil {
		//app.badRequest(w, r, err)
		panic(err)

		return
	}

	//result, err := app.db.NamedExec(`INSERT INTO categories (title, issinglegame, _createdon) VALUES (:title,:isSingleGame,:_createdOn) RETURNING ID`,
	//	map[string]interface{}{
	//		"title":        reqData.Title,
	//		"isSingleGame": reqData.IsSingleGame,
	//		"_createdOn":   time.Now().UnixMilli(),
	//	})
	//

	sqlStr := `insert into categories (title, issinglegame, _createdon) values ($1, $2, $3) returning ID`

	reqData.CreatedOn = time.Now().UnixMilli()
	var id int
	err = app.db.QueryRow(sqlStr, reqData.Title, reqData.IsSingleGame, reqData.CreatedOn).Scan(&id)
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

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {

}
