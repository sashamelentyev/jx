package jir

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_read_uint64_invalid(t *testing.T) {
	should := require.New(t)
	iter := ParseString(Default, ",")
	iter.Uint64()
	should.NotNil(iter.Error)
}

func Test_int8(t *testing.T) {
	inputs := []string{`127`, `-128`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input)
			expected, err := strconv.ParseInt(input, 10, 8)
			should.Nil(err)
			should.Equal(int8(expected), iter.Int8())
		})
	}
}

func Test_read_int16(t *testing.T) {
	inputs := []string{`32767`, `-32768`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input)
			expected, err := strconv.ParseInt(input, 10, 16)
			should.Nil(err)
			should.Equal(int16(expected), iter.Int16())
		})
	}
}

func Test_read_int32(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `2147483647`, `-2147483648`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.Nil(err)
			should.Equal(int32(expected), iter.Int32())
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(Default, bytes.NewBufferString(input), 2)
			expected, err := strconv.ParseInt(input, 10, 32)
			should.Nil(err)
			should.Equal(int32(expected), iter.Int32())
		})
	}
}

func Test_read_int_overflow(t *testing.T) {
	should := require.New(t)
	inputArr := []string{"123451", "-123451"}
	for _, s := range inputArr {
		iter := ParseString(Default, s)
		iter.Int8()
		should.NotNil(iter.Error)

		iterU := ParseString(Default, s)
		iterU.Uint8()
		should.NotNil(iterU.Error)

	}

	inputArr = []string{"12345678912", "-12345678912"}
	for _, s := range inputArr {
		iter := ParseString(Default, s)
		iter.Int16()
		should.NotNil(iter.Error)

		iterUint := ParseString(Default, s)
		iterUint.Uint16()
		should.NotNil(iterUint.Error)
	}

	inputArr = []string{"3111111111", "-3111111111", "1234232323232323235678912", "-1234567892323232323212"}
	for _, s := range inputArr {
		iter := ParseString(Default, s)
		iter.Int32()
		should.NotNil(iter.Error)

		iterUint := ParseString(Default, s)
		iterUint.Uint32()
		should.NotNil(iterUint.Error)
	}

	inputArr = []string{"9223372036854775811", "-9523372036854775807", "1234232323232323235678912", "-1234567892323232323212"}
	for _, s := range inputArr {
		iter := ParseString(Default, s)
		iter.Int64()
		should.NotNil(iter.Error)

		iterUint := ParseString(Default, s)
		iterUint.Uint64()
		should.NotNil(iterUint.Error)
	}
}

func Test_read_int64(t *testing.T) {
	inputs := []string{`1`, `12`, `123`, `1234`, `12345`, `123456`, `9223372036854775807`, `-9223372036854775808`}
	for _, input := range inputs {
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(Default, input)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.Nil(err)
			should.Equal(expected, iter.Int64())
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(Default, bytes.NewBufferString(input), 2)
			expected, err := strconv.ParseInt(input, 10, 64)
			should.Nil(err)
			should.Equal(expected, iter.Int64())
		})
	}
}

func Test_write_uint8(t *testing.T) {
	vals := []uint8{0, 1, 11, 111, 255}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint8(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 3)
	stream.WriteRaw("a")
	stream.WriteUint8(100) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a100", buf.String())
}

func Test_write_int8(t *testing.T) {
	vals := []int8{0, 1, -1, 99, 0x7f, -0x80}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt8(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 4)
	stream.WriteRaw("a")
	stream.WriteInt8(-100) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-100", buf.String())
}

func Test_write_uint16(t *testing.T) {
	vals := []uint16{0, 1, 11, 111, 255, 0xfff, 0xffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint16(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 5)
	stream.WriteRaw("a")
	stream.WriteUint16(10000) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a10000", buf.String())
}

func Test_write_int16(t *testing.T) {
	vals := []int16{0, 1, 11, 111, 255, 0xfff, 0x7fff, -0x8000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt16(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 6)
	stream.WriteRaw("a")
	stream.WriteInt16(-10000) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-10000", buf.String())
}

func Test_write_uint32(t *testing.T) {
	vals := []uint32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint32(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("a")
	stream.WriteUint32(0xffffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}

func Test_write_int32(t *testing.T) {
	vals := []int32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0x7fffffff, -0x80000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt32(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 11)
	stream.WriteRaw("a")
	stream.WriteInt32(-0x7fffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-2147483647", buf.String())
}

func Test_write_uint64(t *testing.T) {
	vals := []uint64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0xffffffffffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteUint64(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("a")
	stream.WriteUint64(0xffffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}

func Test_write_int64(t *testing.T) {
	vals := []int64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff, 0xffffffffff, 0xfffffffffff, 0xffffffffffff, 0xfffffffffffff, 0xffffffffffffff,
		0xfffffffffffffff, 0x7fffffffffffffff, -0x8000000000000000}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(Default, buf, 4096)
			stream.WriteInt64(val)
			_ = stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(val, 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Default, buf, 10)
	stream.WriteRaw("a")
	stream.WriteInt64(0xffffffff) // should clear buffer
	_ = stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}