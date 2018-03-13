package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"GochannelQueue"

	log "github.com/sirupsen/logrus"
)

type ChannelInfo struct {
	GetTime int64
}

type ChannelInfoFIFO struct {
	InfoMap   map[string]bool
	InfoFIFO  []ChannelInfo
	RwMutex   sync.RWMutex
	MaxLength int
}

var (
	logProducer = log.WithFields(log.Fields{
		"Modul": "ChannelInfoProducer",
	})
	logConsumer = log.WithFields(log.Fields{
		"Modul": "ChannelInfoConsumer",
	})

	channelInfoFIFO ChannelInfoFIFO
)

func initLog() {
	//log format
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02T15:04:05.000000"
	log.SetFormatter(formatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func init() {
	initLog()
	channelInfoFIFO.InfoMap = make(map[string]bool)

}

func ChannelInfoProducer(q *queue.Queue) {
	logProducer.Debug("ChannelInfoProducer1 start")
	defer logProducer.Debug("ChannelInfoProducer1 end")

	for {

		now := time.Now().Unix()

		logProducer.Info(fmt.Sprintf("time now:%v\n", now))

		queueNode := &ChannelInfo{
			GetTime: time.Now().Unix(),
		}
		ok := q.Push(queueNode)
		if ok {
			logProducer.Debug(fmt.Sprintf("PUSH queue success! queue len:%v, push node:%+v, \n", q.Len(), queueNode))
		} else {
			logProducer.Debug(fmt.Sprintf("PUSH queue FUll! queue len:%v, push node:%+v, \n", q.Len(), queueNode))
		}
		time.Sleep(time.Second * 5)
	}

}

func ChannelInfoConsumer(queue *queue.Queue, afterSeconds int64) {
	logProducer.Debug(fmt.Sprintf("consume start"))
	for {

		//get head elem
		qNode := queue.Pull()
		if qNode == nil {
			logConsumer.Warn("empty queue!")
			time.Sleep(time.Second)
			continue
		}

		cinfo, ok := qNode.(*ChannelInfo)
		if !ok {
			logConsumer.Warn(fmt.Sprintf("qNode Type:%T, qNode:%+v\n", qNode, qNode))
			time.Sleep(time.Second)
			continue
		}

		logConsumer.Debug(fmt.Sprintf("consume a elem, get_time:%+v, info:%+v\n",
			cinfo.GetTime, cinfo))

		time.Sleep(time.Second * 2)
	}

}

func main() {
	var afterSeconds int64
	afterSeconds = 10

	var channelQueue *queue.Queue
	channelQueue = queue.NewQueue(int(afterSeconds))

	// produce 1 every 2s
	go ChannelInfoProducer(channelQueue)

	time.Sleep(time.Second * 5)
	//consume 1 every 1s
	go ChannelInfoConsumer(channelQueue, afterSeconds)
	for {
		time.Sleep(time.Minute * 5)
	}
}
