package pokemon

import (
	"encoding/json"
	"log"
	"net/http"

	errs "github.com/esvarez/go-api/pkg/error"
)

const (
	endpoint    = "https://pokeapi.co/api/v2/"
	pokemonInfo = "pokemon/"
)

type writer interface {
	AddPokemon(pokemon *Pokemon) error
}

type repository interface {
	writer
}

type Service struct {
	endpoint string
	client   http.Client
	repo     repository
}

func NewService(repo repository) *Service {
	return &Service{
		endpoint: endpoint,
		repo:     repo,
		client:   http.Client{},
	}
}

func (s Service) FindByID(id string) (*Pokemon, error) {
	resp, err := s.client.Get(s.endpoint + pokemonInfo + id)
	if err != nil {
		log.Println("error getting pokemon: ", err)
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, errs.ErrNotFound
	}

	pokemon := &Pokemon{}
	json.NewDecoder(resp.Body).Decode(pokemon)

	if err := s.repo.AddPokemon(pokemon); err != nil {
		log.Println("error adding pokemon: ", err)
		return nil, err
	}

	return pokemon, nil
}
