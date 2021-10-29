package jx

import (
	"io"
	"math"
	"strconv"

	"golang.org/x/xerrors"
)

var intDigits []int8

const uint32SafeToMultiply10 = uint32(0xffffffff)/10 - 1
const uint64SafeToMultiple10 = uint64(0xffffffffffffffff)/10 - 1
const maxFloat64 = 1<<53 - 1

func init() {
	intDigits = make([]int8, 256)
	for i := 0; i < len(intDigits); i++ {
		intDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		intDigits[i] = i - int8('0')
	}
}

// Uint read uint.
func (it *Iter) Uint() (uint, error) {
	if strconv.IntSize == 32 {
		v, err := it.Uint32()
		if err != nil {
			return 0, err
		}
		return uint(v), nil
	}
	v, err := it.Uint64()
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// Int reads integer.
func (it *Iter) Int() (int, error) {
	if strconv.IntSize == 32 {
		v, err := it.Int32()
		return int(v), err
	}
	v, err := it.Int64()
	return int(v), err
}

// Int32 reads int32 value.
func (it *Iter) Int32() (int32, error) {
	c, err := it.next()
	if err != nil {
		return 0, err
	}
	if c == '-' {
		val, err := it.readUint32()
		if err != nil {
			return 0, err
		}
		if val > math.MaxInt32+1 {
			return 0, xerrors.New("overflow")
		}
		return -int32(val), nil
	}
	it.unread()
	val, err := it.readUint32()
	if err != nil {
		return 0, err
	}
	if val > math.MaxInt32 {
		return 0, xerrors.New("overflow")
	}
	return int32(val), nil
}

// Uint32 read uint32
func (it *Iter) Uint32() (uint32, error) {
	return it.readUint32()
}

func (it *Iter) readUint32() (uint32, error) {
	c, err := it.next()
	if err != nil {
		return 0, err
	}
	ind := intDigits[c]
	if ind == 0 {
		return 0, it.assertInt() // single zero
	}
	if ind == invalidCharForNumber {
		return 0, xerrors.Errorf("bad token: %w", err)
	}
	value := uint32(ind)
	if it.tail-it.head > 10 {
		i := it.head
		ind2 := intDigits[it.buf[i]]
		if ind2 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value, nil
		}
		i++
		ind3 := intDigits[it.buf[i]]
		if ind3 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*10 + uint32(ind2), nil
		}
		i++
		ind4 := intDigits[it.buf[i]]
		if ind4 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*100 + uint32(ind2)*10 + uint32(ind3), nil
		}
		i++
		ind5 := intDigits[it.buf[i]]
		if ind5 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*1000 + uint32(ind2)*100 + uint32(ind3)*10 + uint32(ind4), nil
		}
		i++
		ind6 := intDigits[it.buf[i]]
		if ind6 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*10000 + uint32(ind2)*1000 + uint32(ind3)*100 + uint32(ind4)*10 + uint32(ind5), nil
		}
		i++
		ind7 := intDigits[it.buf[i]]
		if ind7 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*100000 + uint32(ind2)*10000 + uint32(ind3)*1000 + uint32(ind4)*100 + uint32(ind5)*10 + uint32(ind6), nil
		}
		i++
		ind8 := intDigits[it.buf[i]]
		if ind8 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*1000000 + uint32(ind2)*100000 + uint32(ind3)*10000 + uint32(ind4)*1000 + uint32(ind5)*100 + uint32(ind6)*10 + uint32(ind7), nil
		}
		i++
		ind9 := intDigits[it.buf[i]]
		value = value*10000000 + uint32(ind2)*1000000 + uint32(ind3)*100000 + uint32(ind4)*10000 + uint32(ind5)*1000 + uint32(ind6)*100 + uint32(ind7)*10 + uint32(ind8)
		it.head = i
		if ind9 == invalidCharForNumber {
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value, nil
		}
	}
	for {
		for i := it.head; i < it.tail; i++ {
			ind = intDigits[it.buf[i]]
			if ind == invalidCharForNumber {
				it.head = i
				if err := it.assertInt(); err != nil {
					return 0, err
				}
				return value, nil
			}
			if value > uint32SafeToMultiply10 {
				value2 := (value << 3) + (value << 1) + uint32(ind)
				if value2 < value {
					return 0, xerrors.New("overflow")
				}
				value = value2
				continue
			}
			value = (value << 3) + (value << 1) + uint32(ind)
		}
		err := it.read()
		if err == io.EOF {
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value, nil
		}
		if err != nil {
			return 0, err
		}
	}
}

// Int64 read int64
func (it *Iter) Int64() (int64, error) {
	c, err := it.next()
	if err != nil {
		return 0, err
	}
	if c == '-' {
		c, err := it.next()
		if err != nil {
			return 0, err
		}
		val, err := it.readUint64(c)
		if err != nil {
			return 0, err
		}
		if val > math.MaxInt64+1 {
			return 0, xerrors.Errorf("%d overflows", val)
		}
		return -int64(val), nil
	}
	val, err := it.readUint64(c)
	if err != nil {
		return 0, err
	}
	if val > math.MaxInt64 {
		return 0, xerrors.Errorf("%d overflows", val)
	}
	return int64(val), nil
}

// Uint64 read uint64
func (it *Iter) Uint64() (uint64, error) {
	c, err := it.next()
	if err != nil {
		return 0, err
	}
	return it.readUint64(c)
}

func (it *Iter) readUint64(c byte) (uint64, error) {
	ind := intDigits[c]
	if ind == 0 {
		if err := it.assertInt(); err != nil {
			return 0, err
		}
		return 0, nil // single zero
	}
	if ind == invalidCharForNumber {
		return 0, xerrors.Errorf("invalid number: %w", badToken(c))
	}
	value := uint64(ind)
	if it.tail-it.head > 10 {
		i := it.head
		ind2 := intDigits[it.buf[i]]
		if ind2 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return 0, nil
		}
		i++
		ind3 := intDigits[it.buf[i]]
		if ind3 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*10 + uint64(ind2), nil
		}
		i++
		ind4 := intDigits[it.buf[i]]
		if ind4 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*100 + uint64(ind2)*10 + uint64(ind3), nil
		}
		i++
		ind5 := intDigits[it.buf[i]]
		if ind5 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*1000 + uint64(ind2)*100 + uint64(ind3)*10 + uint64(ind4), nil
		}
		i++
		ind6 := intDigits[it.buf[i]]
		if ind6 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*10000 + uint64(ind2)*1000 + uint64(ind3)*100 + uint64(ind4)*10 + uint64(ind5), nil
		}
		i++
		ind7 := intDigits[it.buf[i]]
		if ind7 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*100000 + uint64(ind2)*10000 + uint64(ind3)*1000 + uint64(ind4)*100 + uint64(ind5)*10 + uint64(ind6), nil
		}
		i++
		ind8 := intDigits[it.buf[i]]
		if ind8 == invalidCharForNumber {
			it.head = i
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value*1000000 + uint64(ind2)*100000 + uint64(ind3)*10000 + uint64(ind4)*1000 + uint64(ind5)*100 + uint64(ind6)*10 + uint64(ind7), nil
		}
		i++
		ind9 := intDigits[it.buf[i]]
		value = value*10000000 + uint64(ind2)*1000000 + uint64(ind3)*100000 + uint64(ind4)*10000 + uint64(ind5)*1000 + uint64(ind6)*100 + uint64(ind7)*10 + uint64(ind8)
		it.head = i
		if ind9 == invalidCharForNumber {
			if err := it.assertInt(); err != nil {
				return 0, err
			}
			return value, nil
		}
	}
	for {
		for i := it.head; i < it.tail; i++ {
			ind = intDigits[it.buf[i]]
			if ind == invalidCharForNumber {
				it.head = i
				if err := it.assertInt(); err != nil {
					return 0, err
				}
				return value, nil
			}
			if value > uint64SafeToMultiple10 {
				value2 := (value << 3) + (value << 1) + uint64(ind)
				if value2 < value {
					return 0, xerrors.New("overflow")
				}
				value = value2
				continue
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
		err := it.read()
		if err == io.EOF {
			if err := it.assertInt(); err != nil {
				return 0, xerrors.Errorf("assert: %w", err)
			}
			return value, nil
		}
		if err != nil {
			return 0, xerrors.Errorf("read: %w", err)
		}
	}
}

func (it *Iter) assertInt() error {
	if it.head < it.tail && it.buf[it.head] == '.' {
		return xerrors.New("got float instead of int")
	}
	return nil
}
