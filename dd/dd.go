// Package dd is for distributed database using trie's and posets
package dd

import (
	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/trie"
)

type Block struct {
	ID   string
	Trie *trie.Trie
}

func Run(cnfg.Config) error {
	const chains = 3

	return nil
}
