package main

import (
	"net/http"
	"snow-consensus-algo/api"
	"snow-consensus-algo/snow"
)

type FakeSnowUtil struct {
}

func (u *FakeSnowUtil) RandomSample(txID string, k int) int {
	return k
}

func main() {
	snowUtil := &FakeSnowUtil{}
	SnowConsensus := snow.NewSnowConsensus(snowUtil, 5, 5, "123", "456", "789")
	go SnowConsensus.Loop()
	server := api.SnowAPI{SnowConsensus: SnowConsensus}
	http.HandleFunc("/query", server.OnQuery)
	http.ListenAndServe(":9000", nil)
}
