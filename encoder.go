package adpcm

func Encode(in []int, out *[]byte) {
	leftStatus := NewStatus()
	rightStatus := NewStatus()

	for i := 0; i < len(in); i += 2 {
		sample := in[i]
		left := leftStatus.Encode(sample)

		sample = in[i+1]
		right := rightStatus.Encode(sample)

		b := left<<4 | right&0x0f
		*out = append(*out, b)
	}
}
