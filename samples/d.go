package samples

import (
	"bytes"
)

type D struct {
	b  B
	i2 int
	s2 string

	dirtyFlags uint32
}

func (d *D) B() *B {
	return &d.b
}

func (d *D) Dirty(buf *bytes.Buffer) error {
	if d.dirtyFlags == 0 {
		return ErrNoDirtyData
	}

	if err := writeBinary(buf, d.dirtyFlags); err != nil {
		return err
	}

	if d.dirtyFlags&Bit0 != 0 {
		if err := d.b.Dirty(buf); err != nil {
			return err
		}
	}

	if d.dirtyFlags&Bit1 != 0 {
		if _, err := buf.WriteString(d.s2); err != nil {
			return err
		}
		if err := buf.WriteByte(0); err != nil {
			return err
		}
	}

	return nil
}

func (d *D) CleanDirty() {
	d.dirtyFlags = 0
}

func (d *D) MergeFrom(diff []byte) error {
	var buf bytes.Buffer
	buf.Write(diff)

	var dirtyFlags uint32
	if err := readBinary(&buf, &dirtyFlags); err != nil {
		return err
	}

	if dirtyFlags&Bit0 != 0 {
		if err := d.B().MergeFrom(buf.Bytes()); err != nil {
			return err
		}
	}

	if dirtyFlags&Bit1 != 0 {
		by, err := buf.ReadBytes(0)
		if err != nil {
			return err
		}
		d.s2 = string(by[:len(by)-1])
	}

	return nil
}
