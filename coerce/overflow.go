package coerce

import "math"

// FloatOverflowsInt tests if float-overflows-int.
func FloatOverflowsInt(f float64, bitsize int) bool {
	switch bitsize {
	case 64:
		return f < math.MinInt64 || f > math.MaxInt64
	case 8:
		return f < math.MinInt8 || f > math.MaxInt8
	case 16:
		return f < math.MinInt16 || f > math.MaxInt16
	default:
		return f < math.MinInt32 || f > math.MaxInt32
	}
}

// FloatOverflowsUint tests if float-overflows-uint.
func FloatOverflowsUint(f float64, bitsize int) bool {
	switch bitsize {
	case 64:
		return f < 0 || f > math.MaxUint64
	case 8:
		return f < 0 || f > math.MaxUint8
	case 16:
		return f < 0 || f > math.MaxUint16
	default:
		return f < 0 || f > math.MaxUint32
	}
}

// IntOverflowsInt tests if int-overflows-int.
func IntOverflowsInt(i int64, bitsize int) bool {
	switch bitsize {
	case 64:
		return false
	case 8:
		return i < math.MinInt8 || i > math.MaxInt8
	case 16:
		return i < math.MinInt16 || i > math.MaxInt16
	default:
		return i < math.MinInt32 || i > math.MaxInt32
	}
}

// IntOverflowsUint tests if int-overflows-uint.
func IntOverflowsUint(i int64, bitsize int) bool {
	switch bitsize {
	case 64:
		return i < 0
	case 8:
		return i < 0 || i > math.MaxUint8
	case 16:
		return i < 0 || i > math.MaxUint16
	default:
		return i < 0 || i > math.MaxUint32
	}
}

// UintOverflowsInt tests if uint-overflows-int.
func UintOverflowsInt(u uint64, bitsize int) bool {
	switch bitsize {
	case 64:
		return u > math.MaxInt64
	case 8:
		return u > math.MaxInt8
	case 16:
		return u > math.MaxInt16
	default:
		return u > math.MaxInt32
	}
}

// UintOverflowsUint tests if uint-overflows-uint.
func UintOverflowsUint(u uint64, bitsize int) bool {
	switch bitsize {
	case 64:
		return false
	case 8:
		return u > math.MaxUint8
	case 16:
		return u > math.MaxUint16
	default:
		return u > math.MaxUint32
	}
}
