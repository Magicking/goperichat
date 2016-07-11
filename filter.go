package main

import (
	"fmt"
	"time"
	"log"
	"sync"
)
type StatsLive struct {
	UserCount map[string]int
	UserCountM sync.Mutex
	WordFreq map[string]int
	WordFreqCount int
	WordFreqM sync.Mutex
	Filter FilterFunc
	C chan LineEntry
}

type LineEntry struct {
	Id string
	Timestamp int
	Type int	// Death pile is default value (0)
	Data string
}

type FilterFunc func(sl *StatsLive, le LineEntry)

func WordFreq(sl *StatsLive, le LineEntry){
	switch le.Type {
		case 1:
			s, wc := TextReduce(le.Data)
			sl.WordFreqM.Lock()
			sl.WordFreqCount += wc
			for _, word := range s {
				if word == "" {
					continue
				}
				sl.WordFreq[word]++
			}
			sl.WordFreqM.Unlock()
		return
		case 2:
			//Lock Mutex
			//Set UserCount for this channel Id
			//Unlock Mutex
		return
	}
}

func NewStatsLive(f FilterFunc) *StatsLive {
	var ret StatsLive

	ret.UserCount = make(map[string]int)
	ret.WordFreq = make(map[string]int)
	ret.Filter = f
	ret.C = make(chan LineEntry)

	go func(){
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
				case req := <-ret.C:
					if req.Type == 0 {
						log.Println("Exiting", req.Id)
						return
					}
					go ret.Filter(&ret, req)
				case t := <-ticker.C:
					ret.WordFreqM.Lock()
					fmt.Printf("%s: word count: %d\n", t, ret.WordFreqCount)
					for s := range ret.WordFreq {
						fmt.Printf("%s: %d\n", s, ret.WordFreq[s])
					}
					fmt.Printf("\n\n")
					ret.WordFreqM.Unlock()
			}
		}
	}()

	return &ret
}
