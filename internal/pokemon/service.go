package pokemon

import (
	"context"
	"encoding/json"
	"fmt"
	errs "github.com/esvarez/go-api/pkg/error"
	"log"
	"net/http"
	"sync"
)

const (
	endpoint    = "https://pokeapi.co/api/v2/"
	pokemonInfo = "pokemon/"

	numWorkers = 8
)

type writer interface {
	AddPokemon(pokemon *Pokemon) error
}

type repository interface {
	writer
}

type Service struct {
	endpoint string
	workerPool
	client http.Client
	repo   repository
}

type workFunc func(ctx context.Context, pokemon []string) ([]string, error)

type workerPool struct {
	queue  chan Work
	result chan Result
	done   chan any
}

type Work struct {
	fn        workFunc
	pokemon   []string
	condition int
}

type Result struct {
	Value []string
	err   error
}

func (w Work) execute(ctx context.Context) Result {
	val, err := w.fn(ctx, w.pokemon)
	if err != nil {
		return Result{err: err}
	}
	return Result{Value: val}
}

func worker(ctx context.Context, id int, wg *sync.WaitGroup, works <-chan Work, result chan<- Result) {
	defer wg.Done()
	fmt.Println("worker", id, "started")
	count := 0
	for {
		select {
		case work, ok := <-works:
			if !ok {
				return
			}
			//result <- work.execute(ctx)
			// TODO handle error
			//i, _ := strconv.Atoi(work.pokemon[0])
			count++
			//if i%2 == work.condition {
			result <- Result{Value: work.pokemon, err: nil}

			//}
			fmt.Println("worker", id, "started", "processed", count)
			//if count == 10 {
			//	return
			//}
		case <-ctx.Done():
			fmt.Println("workers done")
			result <- Result{err: ctx.Err()}
			return
		}
	}
}

func NewWorkerPool() workerPool {
	// TODO move worker pool to other package
	return workerPool{
		queue:  make(chan Work),
		result: make(chan Result),
		done:   make(chan any),
	}
}

func (w workerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg, w.queue, w.result)
	}
	wg.Wait()
	w.Close()
}

func NewService(repo repository) *Service {
	return &Service{
		endpoint:   endpoint,
		repo:       repo,
		client:     http.Client{},
		workerPool: NewWorkerPool(),
	}
}

func (w workerPool) Close() {
	close(w.done)
	close(w.result)
}

func (w workerPool) AddWork(works []Work) {
	for i := range works {
		w.queue <- works[i]
	}
	close(w.queue)
}

/*
func (w workerPool) AddWorkers(workers int) {
	w.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(workerId int) {
			count := 0
			for job := range w.queue {
				// Logic here
				job.fn.Run()
				fmt.Printf(" worker %d prcesing", workerId)
				count++
				if count == 2 {
					break
				}

			}
			fmt.Printf("worker %d finished %d jobs\n", workerId, count)
			w.wg.Done()
		}(i)
	}
}
*/

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

func (s Service) GetPokemon(tpe string, items, itemsWorker int) (*Pokemon, error) {
	//TODO implement me
	panic("implement me")
}
