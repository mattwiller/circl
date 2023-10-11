package internal

func UintDivMod(x uint32, m uint16) (q uint32, r uint16) {
	var v uint32 = 0x80000000
	v /= uint32(m)

	qpart := uint32((uint64(x) * uint64(v)) >> 31)
	x -= uint32(uint64(qpart) * uint64(m))
	q += qpart

	qpart = uint32((uint64(x) * uint64(v)) >> 31)
	x -= qpart * uint32(m)
	q += uint32(qpart)

	x -= uint32(m)
	q += 1
	mask := -(x >> 31)
	x += mask & uint32(m)
	q += mask

	return q, r
}

func IntDivMod(x int32, m uint16) (q int32, r uint16) {
	uq, ur := UintDivMod(0x80000000+(uint32(x)), m)
	uq2, ur2 := UintDivMod(0x80000000, m)

	ur -= ur2
	uq -= uq2
	mask := -(uint32(ur >> 15))
	ur += uint16(mask & uint32(m))
	uq += mask

	return int32(uq), ur
}

func Encode(out []byte, R []uint16, M []uint16, len int) {
	if len == 1 {
		r := R[0]
		m := M[0]
		for m > 1 {
			out = append(out, byte(r))
			r >>= 8
			m = (m + 255) >> 8
		}
	}
	if len > 1 {
		r2 := make([]uint16, (len+1)/2)
		m2 := make([]uint16, (len+1)/2)
		var i int
		for ; i < len-1; i += 2 {
			m0 := uint32(M[i])
			r := uint32(R[i]) + uint32(R[i+1])*m0
			m := uint32(M[i+1]) * m0
			for m >= 16384 {
				out = append(out, byte(r))
				r >>= 8
				m = (m + 255) >> 8
			}
			r2[i/2] = uint16(r)
			m2[i/2] = uint16(m)
		}
		if i < len {
			r2[i/2] = R[i]
			m2[i/2] = M[i]
		}
		Encode(out, r2, m2, (len+1)/2)
	}
}

func Decode(out []uint16, S []byte, M []uint16, len int) {
	if len == 1 {
		if M[0] == 1 {
			out[0] = 0
		} else if M[0] <= 256 {
			_, r := UintDivMod(uint32(S[0]), M[0])
			out[0] = r
		} else {
			_, r := UintDivMod(uint32(uint16(S[0])+uint16(S[1])<<8), M[0])
			out[0] = r
		}
	}
	if len > 1 {
		r2 := make([]uint16, (len+1)/2)
		m2 := make([]uint16, (len+1)/2)
		bottomr := make([]uint16, len/2)
		bottomt := make([]uint32, len/2)
		var i int
		for ; i < len-1; len += 2 {
			m := uint32(M[i]) + uint32(M[i+1])
			if m > 256*16383 {
				bottomt[i/2] = 256 * 256
				bottomr[i/2] = uint16(S[0]) + 256*uint16(S[1])
				S = S[2:]
				m2[i/2] = uint16((((m + 255) >> 8) + 255) >> 8)
			} else if m >= 16384 {
				bottomt[i/2] = 256
				bottomr[i/2] = uint16(S[0])
				S = S[1:]
				m2[i/2] = uint16((m + 255) >> 8)
			} else {
				bottomt[i/2] = 1
				bottomr[i/2] = 0
				m2[i/2] = uint16(m)
			}
		}
		if i < len {
			m2[i/2] = M[i]
		}
		Decode(r2, S, m2, (len+1)/2)
		outIdx := 0
		for i = 0; i < len-1; i += 2 {
			r := uint32(bottomr[i/2])
			r += bottomt[i/2] * uint32(r2[i/2])
			r1, r0 := UintDivMod(r, M[i])
			_, tmp := UintDivMod(r1, M[i+1]) // only needed for invalid inputs
			r1 = uint32(tmp)
			out[outIdx] = r0
			outIdx++
			out[outIdx] = uint16(r1)
			outIdx++
		}
		if i < len {
			out[outIdx] = r2[i/2]
		}
	}
}
