package jx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Num(t *testing.T) {
	var e Encoder
	e.Num(Num{'1', '2', '3'})
	require.Equal(t, e.String(), "123")
}

func TestNum(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		for _, cc := range []struct {
			Name   string
			String string
			Value  Num
		}{
			{
				Name:   "Int",
				String: "-12",
				Value:  Num("-12"),
			},
			{
				Name:   "IntStr",
				String: `"-12"`,
				Value:  Num(`"-12"`),
			},
		} {
			t.Run(cc.Name, func(t *testing.T) {
				require.Equal(t, cc.String, cc.Value.String())
			})
		}
	})
	t.Run("ZeroValue", func(t *testing.T) {
		// Zero value is invalid because there is no Num.Value.
		var v Num
		require.False(t, v.Zero())
		require.False(t, v.Positive())
		require.False(t, v.Negative())
		require.Equal(t, "<invalid>", v.String())
	})
	t.Run("Integer", func(t *testing.T) {
		t.Run("Int", func(t *testing.T) {
			v := Num{'1', '2', '3'}
			t.Run("Methods", func(t *testing.T) {
				assert.True(t, v.Positive())
				assert.False(t, v.Negative())
				assert.False(t, v.Zero())
				assert.Equal(t, 1, v.Sign())
				assert.Equal(t, "123", v.String())
				assert.True(t, v.Equal(v))
				assert.False(t, v.Equal(Num{}))
			})
			t.Run("Encode", func(t *testing.T) {
				var e Encoder
				e.Num(v)
				require.Equal(t, e.String(), "123")

				n, err := DecodeBytes(e.Bytes()).Int()
				require.NoError(t, err)
				require.Equal(t, 123, n)
			})
		})
		t.Run("FloatAsInt", func(t *testing.T) {
			t.Run("Positive", func(t *testing.T) {
				v := Num{'1', '2', '3', '.', '0'}
				n, err := v.Int64()
				require.NoError(t, err)
				require.Equal(t, int64(123), n)

				un, err := v.Uint64()
				require.NoError(t, err)
				require.Equal(t, uint64(123), un)

				f, err := v.Float64()
				require.NoError(t, err)
				require.InEpsilon(t, 123, f, epsilon)
			})
			t.Run("Negative", func(t *testing.T) {
				v := Num{'1', '2', '3', '.', '0', '0', '1'}
				_, err := v.Int64()
				require.Error(t, err)
			})
		})
		t.Run("Decode", func(t *testing.T) {
			n, err := DecodeStr("12345").Num()
			require.NoError(t, err)
			require.Equal(t, "12345", n.String())
		})
	})
	t.Run("Float", func(t *testing.T) {
		const (
			s = `1.23`
			f = 1.23
		)
		v := Num(s)
		t.Run("Encode", func(t *testing.T) {
			var e Encoder
			e.Num(v)
			require.Equal(t, e.String(), s)

			n, err := DecodeBytes(e.Bytes()).Float64()
			require.NoError(t, err)
			require.InEpsilon(t, f, n, epsilon)
		})
		t.Run("Decode", func(t *testing.T) {
			n, err := DecodeStr(s).Num()
			require.NoError(t, err)
			require.Equal(t, s, n.String())
		})
		t.Run("Methods", func(t *testing.T) {
			assert.True(t, v.Positive())
			assert.False(t, v.Negative())
			assert.False(t, v.Zero())
			assert.Equal(t, 1, v.Sign())
			assert.Equal(t, s, v.String())
		})
	})
}

func BenchmarkNum(b *testing.B) {
	b.Run("FloatAsInt", func(b *testing.B) {
		b.Run("Integer", func(b *testing.B) {
			v := Num{'1', '2', '3', '5', '7', '.', '0'}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := v.Int64(); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("Float30Chars", func(b *testing.B) {
			var v Num
			for i := 0; i < 30; i++ {
				v = append(v, '1')
			}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if err := v.floatAsInt(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("Integer", func(b *testing.B) {
		v := Num{'1', '2', '3', '5', '7'}
		b.Run("Positive", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !v.Positive() {
					b.Fatal("should be positive")
				}
			}
		})
		b.Run("Zero", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if v.Zero() {
					b.Fatal("should be not zero")
				}
			}
		})
		b.Run("Encode", func(b *testing.B) {
			b.ReportAllocs()
			var e Encoder
			for i := 0; i < b.N; i++ {
				e.Num(v)
				e.Reset()
			}
		})
	})
}
