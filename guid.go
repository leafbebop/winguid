// Package guid provides data structure to store a windows guid

package guid

type GUID struct {
	V1 uint32
	V2 uint16
	V3 uint16
	V4 [8]uint8
}

type ErrBadFormat struct {
	ch  byte
	pos int
}

func (e ErrBadFormat) Error() string {
	return string(e.ch) + " ( at " + itoa(e.pos) + " ) is not valid in hex"
}

func itoa(x int) string {
	if x < 10 {
		return string(byte(x) + '0')
	}
	return string([]byte{byte(x/10) + '0', byte(x%10) + '0'})
}

type ErrIncomplete struct {
	Where string
}

func (e ErrIncomplete) Error() string {
	return "GUID sequence is not complete when Parsing " + e.Where
}

func dehex(x byte, pos int) (uint8, error) {
	switch {
	case x >= '0' && x <= '9':
		return x - uint8('0'), nil
	case x >= 'A' && x <= 'F':
		return x - uint8('A') + 10, nil
	case x >= 'a' && x <= 'f':
		return x - uint8('a') + 10, nil
	default:
		return 0, ErrBadFormat{x, pos}
	}

}

// MustFromString is like FromString but panics when the String cannot be parsed.
func MustFromString(s string) GUID {
	g, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return g
}

// FromString parses a string representing a windows guid.
// Supported formats include:
//
// XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
//
// {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
//
// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
//
// {XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX}
func FromString(s string) (GUID, error) {
	var g GUID

	i, j := 0, 0

	l := len(s)

	if s[i] == '{' {
		i++
	}

	j = i + 8
	if j >= l {
		return GUID{}, ErrIncomplete{"V1"}
	}
	for ; i < j; i++ {
		v, err := dehex(s[i], i)
		if err != nil {
			return GUID{}, err
		}
		g.V1 = g.V1<<4 | uint32(v)
	}

	if s[i] == '-' {
		i++
	}

	j = i + 4
	if j >= l {
		return GUID{}, ErrIncomplete{"V2"}
	}

	for ; i < j; i++ {
		v, err := dehex(s[i], i)
		if err != nil {
			return GUID{}, err
		}
		g.V2 = g.V2<<4 | uint16(v)
	}

	if s[i] == '-' {
		i++
	}

	j = i + 4
	if j >= l {
		return GUID{}, ErrIncomplete{"V3"}
	}
	for ; i < j; i++ {
		v, err := dehex(s[i], i)
		if err != nil {
			return GUID{}, err
		}

		g.V3 = g.V3<<4 | uint16(v)
	}

	if s[i] == '-' {
		i++
	}

	j = i + 2
	if j >= l {
		return GUID{}, ErrIncomplete{"V4[0]"}
	}
	for ; i < j; i++ {
		v, err := dehex(s[i], i)
		if err != nil {
			return GUID{}, err
		}

		g.V4[0] = g.V4[0]<<4 | v
	}
	j = i + 2
	if j >= l {
		return GUID{}, ErrIncomplete{"V4[1]"}
	}
	for ; i < j; i++ {
		v, err := dehex(s[i], i)
		if err != nil {
			return GUID{}, err
		}

		g.V4[1] = g.V4[1]<<4 | v
	}

	if s[i] == '-' {
		i++
	}

	j = i + 12
	if j > l {
		return GUID{}, ErrIncomplete{"V4[2...]"}
	}
	for k := 2; k < 8; k++ {
		v1, err := dehex(s[i], i)
		if err != nil {
			return GUID{}, err
		}

		v2, err := dehex(s[i+1], i+1)

		g.V4[k] = v1<<4 | v2
		i += 2
	}

	if l > i && s[i] != '}' {
		return GUID{}, ErrBadFormat{s[i], i}
	}

	return g, nil
}

// String formats a GUID into a string. It uses format
// {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX} .
func (g GUID) String() string {
	b := make([]byte, 38)
	i := 0
	b[i] = byte('{')
	i++

	hex32(g.V1, b, &i)
	b[i] = byte('-')
	i++

	hex16(g.V2, b, &i)
	b[i] = byte('-')
	i++

	hex16(g.V3, b, &i)
	b[i] = byte('-')
	i++

	hex8(g.V4[0], b, &i)
	hex8(g.V4[1], b, &i)
	b[i] = byte('-')
	i++

	for _, x := range g.V4[2:] {
		hex8(x, b, &i)
	}

	b[i] = byte('}')

	return string(b)

}

const (
	mask8 = 255 << (8 * iota)
	mask16
	mask24
	mask32
)

func hex32(x uint32, b []byte, i *int) {
	hex8(uint8(x&mask32>>24), b, i)
	hex8(uint8(x&mask24>>16), b, i)
	hex8(uint8(x&mask16>>8), b, i)
	hex8(uint8(x&mask8), b, i)
}

func hex16(x uint16, b []byte, i *int) {
	hex8(uint8(x&mask16>>8), b, i)
	hex8(uint8(x&mask8), b, i)
}

func hex8(x uint8, b []byte, i *int) {
	h, l := x>>4, x&15

	if h >= 10 {
		b[*i] = h - 10 + byte('a')
	} else {
		b[*i] = h + byte('0')
	}

	(*i)++

	if l >= 10 {
		b[*i] = l - 10 + byte('a')
	} else {
		b[*i] = l + byte('0')
	}

	(*i)++
}
