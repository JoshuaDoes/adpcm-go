package adpcm

import (
	"math"
)

var indexTable = []int{
	-1, -1, -1, -1, 2, 4, 6, 8,
	-1, -1, -1, -1, 2, 4, 6, 8,
}

var stepTable = []int{
	7, 8, 9, 10, 11, 12, 13, 14, 16, 17,
	19, 21, 23, 25, 28, 31, 34, 37, 41, 45,
	50, 55, 60, 66, 73, 80, 88, 97, 107, 118,
	130, 143, 157, 173, 190, 209, 230, 253, 279, 307,
	337, 371, 408, 449, 494, 544, 598, 658, 724, 796,
	876, 963, 1060, 1166, 1282, 1411, 1552, 1707, 1878, 2066,
	2272, 2499, 2749, 3024, 3327, 3660, 4026, 4428, 4871, 5358,
	5894, 6484, 7132, 7845, 8630, 9493, 10442, 11487, 12635, 13899,
	15289, 16818, 18500, 20350, 22385, 24623, 27086, 29794, 32767,
}

type Status struct {
	sample int
	index  int
}

func NewStatus() *Status {
	return &Status{
		sample: 0,
		index:  0,
	}
}

func (status *Status) Decode(nibble byte) int {
	step := stepTable[status.index]

	diff := 0
	if nibble&4 != 0 {
		diff += step
	}
	if nibble&2 != 0 {
		diff += step >> 1
	}
	if nibble&1 != 0 {
		diff += step >> 2
	}
	diff += step >> 3

	if nibble&8 != 0 {
		diff = -diff
	}

	newSample := status.sample + diff
	if newSample > math.MaxInt16 {
		newSample = math.MaxInt16
	} else if newSample < math.MinInt16 {
		newSample = math.MinInt16
	}
	status.sample = newSample

	index := status.index + indexTable[nibble]
	if index < 0 {
		index = 0
	} else if index >= len(stepTable) {
		index = len(stepTable) - 1
	}
	status.index = index

	return newSample
}

func (status *Status) Encode(sample int) byte {
	diff := sample - status.sample
	var nibble byte = 0

	if diff < 0 {
		nibble = 8
		diff = -diff
	}

	var mask byte = 4
	tempStep := stepTable[status.index]
	for i := 0; i < 3; i++ {
		if diff > tempStep {
			nibble |= mask
			diff -= tempStep
		}
		mask >>= 1
		tempStep >>= 1
	}

	// XXX
	// Update the status
	status.Decode(nibble)

	return nibble
}
