package adpcm

import (
	"bytes"
	"io"
)

type Decoder struct {
	NumChannels int

	left  *Status
	right *Status
}

func NewDecoder(numChannels int) *Decoder {
	return &Decoder{
		NumChannels: numChannels,

		left:  NewStatus(),
		right: NewStatus(),
	}
}

func (decoder *Decoder) Decode(in []byte, out *[]int) {
	reader := bytes.NewReader(in)

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		sample := decoder.left.Decode(b >> 4)
		*out = append(*out, sample)

		right := b & 0x0f
		if decoder.NumChannels == 1 {
			sample = decoder.left.Decode(right)
		} else {
			sample = decoder.right.Decode(right)
		}
		*out = append(*out, sample)
	}
}
