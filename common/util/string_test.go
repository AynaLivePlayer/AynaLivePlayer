package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	assert.Equal(t, 3, LevenshteinDistance("kitten", "sitting"))
	assert.Equal(t, 0, LevenshteinDistance("kitten", "kitten"))
	assert.Equal(t, 1, LevenshteinDistance("kitten", "kittens"))
	assert.Equal(t, 2, LevenshteinDistance("kitten", "kitt"))
	assert.Greater(t, LevenshteinDistance("夜曲 周杰伦/方文山", "夜曲 周杰伦"), LevenshteinDistance("夜曲 翻唱A", "夜曲 周杰伦"))
	assert.Greater(t,
		WeightedLevenshteinDistance("Mojito Tommy Hong", "Mojito 周杰伦", 1, 1, 3),
		WeightedLevenshteinDistance("Mojito 周杰伦", "Mojito 周杰伦", 1, 1, 3))
	assert.Greater(t,
		WeightedLevenshteinDistance("默 (Live) 李荣浩/周杰伦", "Mojito 周杰伦", 1, 1, 3),
		WeightedLevenshteinDistance("Mojito 周杰伦", "Mojito 周杰伦", 1, 1, 3))
	assert.Greater(t,
		WeightedLevenshteinDistance("布拉格广场 周杰伦", "Mojito 周杰伦", 1, 1, 3),
		WeightedLevenshteinDistance("Mojito 周杰伦", "Mojito 周杰伦", 1, 1, 3))
	assert.Greater(t,
		WeightedLevenshteinDistance("Mojito（翻自 cover 周杰伦）野猪佩奇", "Mojito 周杰伦", 1, 1, 3),
		WeightedLevenshteinDistance("Mojito 周杰伦", "Mojito 周杰伦", 1, 1, 3))
	//assert.Less(t, WeightedLevenshteinDistance("夜曲 周杰伦/方文山", "夜曲 周杰伦",1,1,3), WeightedLevenshteinDistance("夜曲 翻唱A", "夜曲 周杰伦",1,1,3))
}

func TestLongestCommonString(t *testing.T) {
	assert.Equal(t, "itt", LongestCommonString("kitten", "sitting"))
	assert.Equal(t, "布拉格广场", LongestCommonString("布拉格广场 周杰伦", "布拉格广场"))
}
