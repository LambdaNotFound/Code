package bfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_minKnightMoves(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{name: "origin", x: 0, y: 0, expected: 0},
		{name: "one_move", x: 2, y: 1, expected: 1},
		{name: "two_moves", x: 5, y: 5, expected: 4},
		{name: "leetcode_example1", x: 2, y: 1, expected: 1},
		{name: "leetcode_example2", x: 5, y: 5, expected: 4},
		{name: "negative_quadrant", x: -2, y: -1, expected: 1},
		{name: "far_target", x: 100, y: 100, expected: 68},
		{name: "one_one", x: 1, y: 1, expected: 2}, // (0,0)→(2,-1)→(1,1)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, minKnightMoves(tt.x, tt.y))
		})
	}
}

func Test_ladderLengthBidirectionalBFS(t *testing.T) {
	testCases := []struct {
		name      string
		beginWord string
		endWord   string
		wordList  []string
		expected  int
	}{
		{
			"case 1",
			"hit",
			"cog",
			[]string{"hot", "dot", "dog", "lot", "log", "cog"},
			5,
		},
		{
			"case 2",
			"hit",
			"cog",
			[]string{"hot", "dot", "dog", "lot", "log"},
			0,
		},
		{
			"case 3",
			"ymain",
			"oecij",
			[]string{"ymann", "yycrj", "oecij", "ymcnj", "yzcrj", "yycij", "xecij", "yecij", "ymanj", "yzcnj", "ymain"},
			10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ladderLengthBidirectionalBFS(tc.beginWord, tc.endWord, tc.wordList)
			assert.Equal(t, tc.expected, result)
		})
	}
}
