package samples

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	BitAll = 0xFFFFFFFF
	Bit0   = 1
	Bit1   = 1 << 1
	Bit2   = 1 << 2
	Bit3   = 1 << 3
	Bit4   = 1 << 4
	Bit5   = 1 << 5
	Bit6   = 1 << 6
	Bit7   = 1 << 7
	Bit8   = 1 << 8
	Bit9   = 1 << 9
	Bit10  = 1 << 10
	Bit11  = 1 << 11
)

var (
	ErrNoDirtyData = errors.New("dirty: no dirty data")
)

func writeBinary(writer io.Writer, data interface{}) error {
	return binary.Write(writer, binary.LittleEndian, data)
}

func readBinary(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.LittleEndian, data)
}

type MarkDirty func()

type B struct {
	i int64
	s string

	dirtyFlags uint32
	MarkDirty
}

func (b *B) SetMark(mark MarkDirty) {
	b.MarkDirty = mark
}

func (b *B) MarkParent() {
	if b.MarkDirty != nil {
		b.MarkDirty()
	}
}

func (b *B) MarkAll() {
	b.dirtyFlags &= BitAll
}

func (b *B) I() int64 {
	return b.i
}

func (b *B) SetI(i int64) {
	b.i = i
	b.dirtyFlags |= Bit0
	b.MarkParent()
}

func (b *B) S() string {
	return b.s
}

func (b *B) SetS(s string) {
	b.s = s
	b.dirtyFlags |= Bit1
	b.MarkParent()
}

func (b *B) Dirty(buf *bytes.Buffer) error {
	if b.dirtyFlags == 0 {
		return ErrNoDirtyData
	}

	if err := writeBinary(buf, b.dirtyFlags); err != nil {
		return err
	}

	if b.dirtyFlags&Bit0 != 0 {
		if err := writeBinary(buf, b.i); err != nil {
			return err
		}
	}

	if b.dirtyFlags&Bit1 != 0 {
		if _, err := buf.WriteString(b.s); err != nil {
			return err
		}
		if err := buf.WriteByte(0); err != nil {
			return err
		}
	}

	return nil
}

func (b *B) CleanDirty() {
	b.dirtyFlags = 0
}

func (b *B) MergeFrom(diff []byte) error {
	var buf bytes.Buffer
	buf.Write(diff)

	var dirtyFlags uint32
	if err := readBinary(&buf, &dirtyFlags); err != nil {
		return err
	}

	if dirtyFlags&Bit0 != 0 {
		i := b.i
		if err := readBinary(&buf, &b.i); err != nil {
			b.i = i
			return err
		}
	}

	if dirtyFlags&Bit1 != 0 {
		by, err := buf.ReadBytes(0)
		if err != nil {
			return err
		}
		b.s = string(by[:len(by)-1])
	}

	return nil
}
