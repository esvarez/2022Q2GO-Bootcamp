package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/esvarez/go-api/internal/pokemon"
	"github.com/esvarez/go-api/pkg/web"

	"github.com/gorilla/mux"
)

const (
	pokemonID           = "pokemon_id"
	items               = "items"
	tpe                 = "type"
	itemsPerWorker      = "items_per_workers"
	defaultItems        = 5
	defaultItemsWorkers = 1
)

var (
	validTypes = map[string]bool{"odd": true, "even": true}
)

type PokemonHandler struct {
	pokemonService pokemonService
}

type pokemonService interface {
	FindByID(id string) (*pokemon.Pokemon, error)
	GetPokemon(tpe string, items, itemsWorker int) ([]pokemon.Pokemon, error)
}

func NewPokemonHandler(service pokemonService) *PokemonHandler {
	return &PokemonHandler{
		pokemonService: service,
	}
}

func (p PokemonHandler) findPokemon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		poke, err := p.pokemonService.FindByID(params[pokemonID])
		if err != nil {
			log.Printf("error getting pokemon: %v", err)
			web.ErrorResponse(err).Send(w)
			return
		}

		web.Success(poke, http.StatusOK).Send(w)
	})
}

func (p PokemonHandler) getPokemon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if _, ok := params[tpe]; !ok {
			log.Println("no type")
			web.BadRequestError.Send(w)
			return

		}

		if val, ok := validTypes[params[tpe][0]]; !ok {
			log.Println("invalid type", val)
			web.BadRequestError.Send(w)
			return
		}

		t := params[tpe][0]
		var itms, itmsWorker int
		if val, ok := params[items]; ok {
			itms, _ = strconv.Atoi(val[0])
			if itms < 1 {
				log.Println("invalid items", val[0])
				web.BadRequestError.Send(w)
				return
			}
		} else {
			itms = 5
		}

		if val, ok := params[itemsPerWorker]; ok {
			itmsWorker, _ = strconv.Atoi(val[0])
			if itmsWorker < 1 {
				log.Println("invalid items", val[0])
				web.BadRequestError.Send(w)
				return
			}
		} else {
			itmsWorker = 1
		}

		response, err := p.pokemonService.GetPokemon(t, itms, itmsWorker)
		if err != nil {
			web.InternalServerError.Send(w)
			return
		}

		web.Success(response, http.StatusOK).Send(w)
	})
}

func MakePokemonHandler(r *mux.Router, handler *PokemonHandler) {
	r.Handle("/pokemon/{pokemon_id}", handler.findPokemon()).
		Methods(http.MethodGet)
	r.Handle("/pokemon", handler.getPokemon()).
		Methods(http.MethodGet)
}
