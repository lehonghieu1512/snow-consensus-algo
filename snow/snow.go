package snow

import (
	"fmt"
	"snow-consensus-algo/transaction"
	"snow-consensus-algo/utils"
	"sync"
	"time"
)

type SnowUtil interface {
	RandomSample(txID string, k int) int
}

type SnowConsensus struct {
	k                  int
	queryTxSet         []string
	knownTxSet         []string
	conflictSet        map[string][]string
	lastPre            map[string]string
	preference         map[string]string
	consecutiveSuccess map[string]int
	chits              map[string]int
	ancestry           map[string][]string
	snowUtil           SnowUtil
	mu                 sync.Mutex
	alpha              int
}

func (c *SnowConsensus) getParents(txID string) []string {
	var parents []string
	// find
	// leaf with chit = 1

	var allAncestors []string
	for _, v := range c.ancestry {
		allAncestors = append(allAncestors, v...)
	}

	for _, txID := range c.knownTxSet {
		if !utils.IsInSlice(txID, allAncestors) && c.chits[txID] == 1 {
			parents = append(parents, txID)
		}
	}

	return parents
}

func (c *SnowConsensus) Loop() {
	for {
		fmt.Println(c.ancestry)
		time.Sleep(3 * time.Second)
		var txID string
		for _, ID := range c.knownTxSet {
			if !utils.IsInSlice(ID, c.queryTxSet) {
				txID = ID
			}
		}
		if txID == "" {
			continue
		}

		okReps := c.snowUtil.RandomSample(txID, c.k)

		if okReps >= c.alpha {
			if _, ok := c.ancestry[txID]; !ok {
				parents := c.getParents(txID)
				c.ancestry[txID] = parents
			}
			c.chits[txID] = 1
			for _, ancestor := range c.ancestry[txID] {
				if c.confidence(ancestor) > c.confidence(c.preference[ancestor]) {
					c.preference[ancestor] = ancestor
				}
				if ancestor != c.lastPre[ancestor] {
					c.lastPre[ancestor] = ancestor
					c.consecutiveSuccess[ancestor] = 1
				} else {
					c.consecutiveSuccess[ancestor] += 1
				}
			}
		} else {
			for _, ancestor := range c.ancestry[txID] {
				c.consecutiveSuccess[ancestor] = 0
			}
			c.chits[txID] = 0
		}
		c.mu.Lock()
		c.queryTxSet = append(c.queryTxSet, txID)
		c.mu.Unlock()
	}
}

func (c *SnowConsensus) confidence(txID string) int {
	confidence := c.chits[txID]
	for k, v := range c.ancestry {
		if utils.IsInSlice(txID, v) && c.chits[k] == 1 {
			confidence += 1
		}
	}
	return confidence
}

func (c *SnowConsensus) OnReceiveTx(tx *transaction.Tx) {
	if !utils.IsInSlice(tx.ID, c.knownTxSet) {
		parents := c.getParents(tx.ID)
		c.ancestry[tx.ID] = parents
		c.registerTx(tx)
		c.knownTxSet = append(c.knownTxSet, tx.ID)
	}
}

func (c *SnowConsensus) registerTx(tx *transaction.Tx) {
	if _, ok := c.conflictSet[tx.ID]; !ok {
		c.preference[tx.ID] = tx.ID
		c.lastPre[tx.ID] = tx.ID
		c.consecutiveSuccess[tx.ID] = 0
	} else {
		c.conflictSet[tx.ID] = append(c.conflictSet[tx.ID], tx.ID)
	}
	c.chits[tx.ID] = 1
}

func (c *SnowConsensus) isPreferred(txID string) bool {
	fmt.Println(txID, c.preference[txID], c.preference, "ssss")
	return txID == c.preference[txID]
}

func (c *SnowConsensus) isStronglyPreferred(tx *transaction.Tx) bool {
	fmt.Print(c.ancestry[tx.ID], c.preference, "\n")
	for _, ancestor := range c.ancestry[tx.ID] {
		if !c.isPreferred(ancestor) {
			return false
		}
	}
	return true
}

func (c *SnowConsensus) OnQuery(tx *transaction.Tx) bool {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.OnReceiveTx(tx)
	return c.isStronglyPreferred(tx)
}

func NewSnowConsensus(snowUtil SnowUtil, k int, alpha int, initialTx ...string) *SnowConsensus {
	return &SnowConsensus{
		knownTxSet:         initialTx,
		preference:         make(map[string]string),
		conflictSet:        make(map[string][]string),
		lastPre:            make(map[string]string),
		consecutiveSuccess: make(map[string]int),
		chits:              make(map[string]int),
		ancestry:           make(map[string][]string),
		snowUtil:           snowUtil,
		k:                  k,
		alpha:              alpha,
	}
}
