package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/esvarez/go-api/internal/pokemon"
	errs "github.com/esvarez/go-api/pkg/error"
	"github.com/esvarez/go-api/pkg/web"

	"github.com/gorilla/mux"
)

const pokemonID = "pokemon_id"

type PokemonHandler struct {
	PokemonService pokemonService
}

type pokemonService interface {
	FindByID(id string) (*pokemon.Pokemon, error)
}

func NewPokemonHandler(service pokemonService) *PokemonHandler {
	return &PokemonHandler{
		PokemonService: service,
	}
}

func (p PokemonHandler) findPokemon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		poke, err := p.PokemonService.FindByID(params[pokemonID])
		if err != nil {
			var status web.AppError
			log.Printf("error getting pokemon: %v", err)
			switch {
			case errors.Is(err, errs.ErrNotFound):
				status = web.ResourceNotFoundError
			default:
				status = web.InternalServerError
			}
			status.Send(w)
			return
		}

		web.Success(poke, http.StatusOK).Send(w)
	})
}

func MakePokemonHandler(r *mux.Router, handler *PokemonHandler) {
	r.Handle("/pokemon/{pokemon_id}", handler.findPokemon()).
		Methods(http.MethodGet)
}
