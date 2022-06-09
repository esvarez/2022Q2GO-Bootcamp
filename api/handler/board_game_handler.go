package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/esvarez/go-api/internal/boardgame"
	"github.com/esvarez/go-api/pkg/web"

	"github.com/gorilla/mux"
)

const boardGameID = "board_game_id"

type BoardGameHandler struct {
	BoardGameService boardGameService
}

type boardGameService interface {
	FindByID(id int) (*boardgame.BoardGame, error)
}

func NewBoardGameHandler(service boardGameService) *BoardGameHandler {
	return &BoardGameHandler{
		BoardGameService: service,
	}
}

func (b BoardGameHandler) findBoardGameByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params[boardGameID])
		if err != nil {
			log.Printf("error gettin board game id: %v", err)
			web.InternalServerError.Send(w)
			return
		}
		bGame, err := b.BoardGameService.FindByID(id)
		if err != nil {
			log.Printf("error getting boardgame: %v", err)
			web.ErrorResponse(err).Send(w)
			return
		}

		web.Success(bGame, http.StatusOK).Send(w)
	})
}

func MakeBoardGameHandler(router *mux.Router, handler *BoardGameHandler) {
	router.Handle("/boardgame/{board_game_id}", handler.findBoardGameByID()).
		Methods(http.MethodGet)
}
