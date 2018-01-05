package adpcm

import (
	"math"
)

var indexTable = []int{
	-1, -1, -1, -1, 2, 4, 6, 8,
	-1, -1, -1, -1, 2, 4, 6, 8,
}

var stepTable = []int{
	256, 272, 304, 336, 368, 400, 448, 496, 544, 592, 656, 720, 800, 880, 960,
	1056, 1168, 1280, 1408, 1552, 1712, 1888, 2080, 2288, 2512, 2768, 3040,
	3344, 3680, 4048, 4464, 4912, 5392, 5936, 6528, 7184, 7904, 8704, 9568,
	10528, 11584, 12736, 14016, 15408, 16960, 18656, 20512, 22576, 24832,
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
