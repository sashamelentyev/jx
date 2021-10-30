package jx

import (
	"golang.org/x/xerrors"
)

// Null reads a json object as null and
// returns whether it's a null or not.
func (d *Decoder) Null() error {
	if err := d.expectNext('n'); err != nil {
		return err
	}
	return d.skipThreeBytes('u', 'l', 'l') // null
}

// Bool reads a json object as Bool
func (d *Decoder) Bool() (bool, error) {
	c, err := d.next()
	if err != nil {
		return false, err
	}
	switch c {
	case 't':
		if err := d.skipThreeBytes('r', 'u', 'e'); err != nil {
			return false, err
		}
		return true, nil
	case 'f':
		return false, d.skipFourBytes('a', 'l', 's', 'e')
	default:
		return false, badToken(c)
	}
}

// Skip skips a json object and positions to relatively the next json object.
func (d *Decoder) Skip() error {
	c, err := d.next()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		if err := d.strSkip(); err != nil {
			return xerrors.Errorf("str: %w", err)
		}
		return nil
	case 'n':
		return d.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		return d.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		return d.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		d.unread()
		_, err := d.Float32()
		return err
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.skipNumber()
	case '[':
		if err := d.skipArray(); err != nil {
			return xerrors.Errorf("array: %w", err)
		}
		return nil
	case '{':
		if err := d.skipObject(); err != nil {
			return xerrors.Errorf("object: %w", err)
		}
		return nil
	default:
		return badToken(c)
	}
}

func (d *Decoder) skipFourBytes(b1, b2, b3, b4 byte) error {
	for _, b := range [...]byte{b1, b2, b3, b4} {
		c, err := d.byte()
		if err != nil {
			return err
		}
		if c != b {
			return badToken(c)
		}
	}
	return nil
}

func (d *Decoder) skipThreeBytes(b1, b2, b3 byte) error {
	for _, b := range [...]byte{b1, b2, b3} {
		c, err := d.byte()
		if err != nil {
			return err
		}
		if c != b {
			return badToken(c)
		}
	}
	return nil
}

func (d *Decoder) skipNumber() error {
	ok, err := d.skipNumberFast()
	if err != nil || ok {
		return err
	}
	d.unread()
	if _, err := d.Float64(); err != nil {
		return err
	}
	return nil
}

func (d *Decoder) skipNumberFast() (ok bool, err error) {
	dotFound := false
	for i := d.head; i < d.tail; i++ {
		c := d.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case tDot:
			if dotFound {
				return false, xerrors.New("more than one dot")
			}
			if i+1 == d.tail {
				return false, nil
			}
			c = d.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				return false, xerrors.New("no digit after dot")
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if d.head == i {
					return false, nil // if - without following digits
				}
				d.head = i
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

func (d *Decoder) strSkip() error {
	ok, err := d.strFastSkip()
	if err != nil || ok {
		return err
	}

	d.unread()
	if _, err := d.str(value{ignore: true}); err != nil {
		return err
	}
	return nil
}

func (d *Decoder) strFastSkip() (ok bool, err error) {
	for i := d.head; i < d.tail; i++ {
		c := d.buf[i]
		switch {
		case c == '"':
			d.head = i + 1
			return true, nil
		case c == '\\':
			return false, nil
		case c < ' ':
			return false, badToken(c)
		}
	}
	return false, nil
}

func (d *Decoder) skipObject() error {
	d.unread()
	return d.ObjectBytes(func(iter *Decoder, _ []byte) error {
		return iter.Skip()
	})
}

func (d *Decoder) skipArray() error {
	d.unread()
	return d.Array(func(iter *Decoder) error {
		return iter.Skip()
	})
}