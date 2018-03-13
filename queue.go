package queue

import (
	"log"
	"sync"
)

type Queue struct {
	PoolSize int
	PoolChan chan interface{}
	Count    int
	Mutex    sync.Mutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		PoolSize: size,
		PoolChan: make(chan interface{}, size),
		Count:    0,
	}
}

func (this *Queue) Init(size int) *Queue {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.PoolSize = size
	this.Count = 0
	this.PoolChan = make(chan interface{}, size)
	return this
}

func (this *Queue) Push(i interface{}) (ret bool) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if len(this.PoolChan) >= this.PoolSize || this.Count > this.PoolSize {
		return false
	}

	this.PoolChan <- i
	// log.Printf("Push count1 %+v\n", this.Count)
	this.Count++
	log.Printf("Push count:%+v\n", this.Count)

	return true
}

func (this *Queue) PushSlice(s []interface{}) {
	for _, i := range s {
		this.Push(i)
	}
}

func (this *Queue) Pull() (node interface{}) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if 0 == this.Count {
		node = nil
	} else {
		node = <-this.PoolChan
		this.Count--

		log.Printf("Pull count %+v\n", this.Count)
	}

	return node
}

// 二次使用Queue实例时，根据容量需求进行高效转换
func (this *Queue) Exchange(num int) (add int) {
	last := len(this.PoolChan)

	if last >= num {
		add = int(0)
		return
	}

	if this.PoolSize < num {
		pool := []interface{}{}
		for i := 0; i < last; i++ {
			pool = append(pool, <-this.PoolChan)
		}
		// 重新定义、赋值
		this.Init(num).PushSlice(pool)
	}

	add = num - last
	return
}

func (this *Queue) Len() (len int) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	len = this.Count

	return len
}
