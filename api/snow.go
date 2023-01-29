package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"snow-consensus-algo/snow"
	"snow-consensus-algo/transaction"
)

type Tx struct {
	ID string `json:"id"`
}

type QueryResp struct {
	IsOK bool `json:"is_ok"`
}

type SnowAPI struct {
	SnowConsensus *snow.SnowConsensus
}

func (a *SnowAPI) OnQuery(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var tx Tx
	err = json.Unmarshal(body, &tx)
	if err != nil {
		panic(err)
	}

	snowTx := &transaction.Tx{ID: tx.ID}

	b := a.SnowConsensus.OnQuery(snowTx)
	resp := QueryResp{IsOK: b}
	fmt.Println("aaa", resp)
	bytess, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Write(bytess)
}
