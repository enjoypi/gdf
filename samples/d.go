package samples

import (
	"bytes"
)

type D struct {
	b  B
	s2 string
	i2 int64

	dirtyFlags uint32
	MarkDirty
}

func (d *D) Init() {
	d.b.SetMark(d.MarkB)
}

func (d *D) B() *B {
	return &d.b
}

func (d *D) SetB(b *B) {
	d.b = *b
	b.MarkAll()
	d.MarkB()
}

func (d *D) MarkB() {
	d.dirtyFlags |= Bit0
}

func (d *D) S2() string {
	return d.s2
}

func (d *D) SetS2(s2 string) {
	d.s2 = s2
	d.dirtyFlags |= Bit1
}

func (d *D) I2() int64 {
	return d.i2
}

func (d *D) SetI2(i2 int64) {
	d.i2 = i2
	d.dirtyFlags |= Bit2
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

	if d.dirtyFlags&Bit2 != 0 {
		if err := writeBinary(buf, d.i2); err != nil {
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

	if dirtyFlags&Bit2 != 0 {
		if err := readBinary(&buf, &d.i2); err != nil {
			return err
		}
	}

	return nil
}
