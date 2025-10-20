package queue

import (
	"fmt"
	"sync"

	"github.com/Cere6rum/MicroBlog2/internal/models"
)

// LikeQueue - очередь для асинхронной обработки лайков
type LikeQueue struct {
	queue   chan models.LikeEvent
	workers int
	wg      sync.WaitGroup
	done    chan struct{}
}

// NewLikeQueue создает новую очередь лайков
func NewLikeQueue(bufferSize, workers int) *LikeQueue {
	return &LikeQueue{
		queue:   make(chan models.LikeEvent, bufferSize),
		workers: workers,
		done:    make(chan struct{}),
	}
}

// Start запускает обработчики (воркеры) очереди
func (lq *LikeQueue) Start(processFunc func(models.LikeEvent) error) {
	for i := 0; i < lq.workers; i++ {
		lq.wg.Add(1)
		go lq.worker(i, processFunc)
	}
}

// worker - горутина-обработчик событий лайков
func (lq *LikeQueue) worker(id int, processFunc func(models.LikeEvent) error) {
	defer lq.wg.Done()

	for {
		select {
		case event := <-lq.queue:
			// Обрабатываем событие лайка
			if err := processFunc(event); err != nil {
				fmt.Printf("Worker %d: ошибка обработки лайка: %v\n", id, err)
			}

		case <-lq.done:
			// Завершаем работу воркера
			return
		}
	}
}

// Enqueue добавляет событие лайка в очередь
func (lq *LikeQueue) Enqueue(event models.LikeEvent) {
	lq.queue <- event
}

// Stop останавливает обработку очереди
func (lq *LikeQueue) Stop() {
	close(lq.done)
	lq.wg.Wait()
	close(lq.queue)
}
