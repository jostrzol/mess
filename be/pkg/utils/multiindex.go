package utils

type Multiindex struct {
	current []int
	end     []int
}

func NewMultiindexLike[T any](xs [][]T) Multiindex {
	end := make([]int, len(xs))
	for i, x := range xs {
		end[i] = len(x)
	}
	return NewMultiindex(end)
}

func NewMultiindex(end []int) Multiindex {
	current := make([]int, len(end))
	return Multiindex{current: current, end: end}
}

func (m Multiindex) Next() {
	carry := true
	for i := len(m.end) - 1; carry && i >= 0; i-- {
		m.current[i]++
		if m.current[i] >= m.end[i] {
			m.current[i] = 0
		} else {
			carry = false
		}
	}
	if carry {
		// reached the end
		m.current[0] = m.end[0]
	}
}

func (m Multiindex) Current() []int {
	return m.current
}

func (m Multiindex) IsEnd() bool {
	return m.current[0] >= m.end[0]
}
