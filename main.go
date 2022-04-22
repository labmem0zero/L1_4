package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartWorker(id int, terminate <-chan int,jobs <-chan int){
	for {
		select {
		case <-terminate:
			fmt.Printf("Остановлен воркер №%v\n",id)
			return
		case <-jobs:
			fmt.Printf("Воркер №%v:%v\n",id,<-jobs)
			time.Sleep(time.Second/20)
		}
	}
}

func WriteChannelINFINITE(jobs chan<-int){
	for {
		jobs<-rand.Intn(100)
		time.Sleep(time.Duration(rand.Float64()/10)*time.Second)
	}
}

func main(){
	//Решил останавливать сигнальным канал, потому что в условиях текущей задачи мне не требуется контекст
	//так что чем проще, тем лучше.
	//Для остановки горутины в серьезном проекте я бы воспользовался
	//контекстом с таймаутом. С контекстом проще систематизировать данные
	//когда уже есть много кода
	workersAmount:=2
	jobs:=make(chan int)
	terminateChanel:=make(chan int)
	go WriteChannelINFINITE(jobs)
	for i:=1;i<workersAmount+1;i++{
		go StartWorker(i,terminateChanel,jobs)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,syscall.SIGINT)
	<-sigChan
	close(terminateChanel)
	<-sigChan
}
