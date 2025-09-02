package solid_coding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trie(t *testing.T) {
	trie := Constructor()
	trie.Insert("apple")
    assert.Equal(t, true, trie.Search("apple"))

    assert.Equal(t, false, trie.Search("app"))
    assert.Equal(t, true, trie.StartsWith("app"))
    
    assert.Equal(t, false, trie.Search("banana"))

}