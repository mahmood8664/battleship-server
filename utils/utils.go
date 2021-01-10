package utils

import (
	"math/rand"
)

func MaskId(id string) string {
	if len(id) != 24 {
		return id
	}

	seed := 0
	bytes := []byte(id)
	for _, b := range bytes {
		seed = int(b) + seed
	}

	if seed < 0 {
		seed = seed * -1
	}

	rand.Seed(int64(seed))

	for i := 0; i < 100; i++ {
		r1 := rand.Int31n(24)
		r2 := 23 - r1
		temp := bytes[r1]
		bytes[r1] = bytes[r2]
		bytes[r2] = temp
	}

	return string(bytes)

}

////////////////////////////

func GetMapKeySlice(myMap map[int]bool) []int {
	keys := make([]int, len(myMap))
	i := 0
	for k := range myMap {
		keys[i] = k
		i++
	}
	return keys
}
