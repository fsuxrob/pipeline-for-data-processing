package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const bufferSize int = 10
const bufferDrainInterval time.Duration = 5 * time.Second

type RingBuffer struct {
	size int
	data []int
	pos  int
	m    sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{size, make([]int, size), -1, sync.Mutex{}}
}

func (r *RingBuffer) Push(el int) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		for i := 1; i < r.size-1; i++ {
			r.data[i-1] = r.data[i]
		}
		r.data[r.pos] = el
	} else {
		r.pos++
		r.data[r.pos] = el
	}
}

func (r *RingBuffer) Get() []int {
	if r.pos < 0 {
		return nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	var output []int = r.data[:r.pos+1]
	r.pos = -1
	return output
}

type StageInt func(<-chan bool, <-chan int) <-chan int

type PipeLineInt struct {
	stages []StageInt
	done   <-chan bool
}

func NewPipelineInt(done <-chan bool, stages ...StageInt) *PipeLineInt {
	return &PipeLineInt{done: done, stages: stages}
}

func (p *PipeLineInt) Run(source <-chan int) <-chan int {
	var c <-chan int = source
	for index := range p.stages {
		c = p.runStageInt(p.stages[index], c)
	}
	return c
}

func (p *PipeLineInt) runStageInt(stage StageInt, sourceChan <-chan int) <-chan int {
	return stage(p.done, sourceChan)
}

func dataSource() (<-chan int, <-chan bool) {
	c := make(chan int)
	done := make(chan bool)

	go func() {
		defer close(done)
		scanner := bufio.NewScanner(os.Stdin)
		var data string
		for {
			scanner.Scan()
			data = scanner.Text()
			if strings.EqualFold(data, "exit") {
				log.Println("Программа завершила работу!")
				return
			}
			i, err := strconv.Atoi(data)
			if err != nil {
				log.Println("Программа обрабатывает только целые числа!")
				continue
			}
			c <- i
		}
	}()
	return c, done
}

func negativeFilterStageInt(done <-chan bool, input <-chan int) <-chan int {
	convertedIntChan := make(chan int)
	go func() {
		for {
			select {
			case data := <-input:
				if data > 0 {
					select {
					case convertedIntChan <- data:
					case <-done:
						return
					}
				}
			case <-done:
				return
			}
		}
	}()
	return convertedIntChan
}

func specialFilterStageInt(done <-chan bool, input <-chan int) <-chan int {
	filteredIntChan := make(chan int)
	go func() {
		for {
			select {
			case data := <-input:
				if data != 0 && data%3 == 0 {
					select {
					case filteredIntChan <- data:
					case <-done:
						return
					}
				}
			case <-done:
				return
			}
		}
	}()
	return filteredIntChan
}

func bufferStageInt(done <-chan bool, input <-chan int) <-chan int {
	bufferedIntChan := make(chan int)
	buffer := NewRingBuffer(bufferSize)
	go func() {
		for {
			select {
			case data := <-input:
				buffer.Push(data)
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(bufferDrainInterval):
				bufferData := buffer.Get()
				for _, data := range bufferData {
					select {
					case bufferedIntChan <- data:
					case <-done:
						return
					}
				}
			case <-done:
				return
			}
		}
	}()
	return bufferedIntChan
}

func consumer(done <-chan bool, input <-chan int) {
	for {
		select {
		case data := <-input:
			log.Printf("Обработаны данные: %d\n", data)
		case <-done:
			return
		}
	}
}

func main() {
	source, done := dataSource()
	pipline := NewPipelineInt(done, negativeFilterStageInt, specialFilterStageInt, bufferStageInt)
	log.Println("Пайплайн запущен!")
	consumer(done, pipline.Run(source))
}
