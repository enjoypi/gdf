package samples

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	BitAll = 0xFFFFFFFF
	Bit00  = 1
	Bit01  = 1 << 1
	Bit02  = 1 << 2
	Bit03  = 1 << 3
	Bit04  = 1 << 4
	Bit05  = 1 << 5
	Bit06  = 1 << 6
	Bit07  = 1 << 7
	Bit08  = 1 << 8
	Bit09  = 1 << 9
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
	i     int64
	str     string
	m     map[string]string
	slice []string

	dirtyFlags uint32

	dirtyM map[string]bool		// true表示删除，false表示变更
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
	b.dirtyFlags |= Bit00
	b.MarkParent()
}

func (b *B) S() string {
	return b.str
}

func (b *B) SetS(s string) {
	b.str = s
	b.dirtyFlags |= Bit01
	b.MarkParent()
}

func (b *B) LoadFromM(key string) (string, bool) {
	value, ok := b.m[key]
	return value, ok
}

func (b *B) StoreIntoM(key, value string) {
	b.m[key] = value
	b.dirtyFlags |= Bit02
}

func (b *B) Slice() []string {
	return b.slice
}

func (b *B) Dirty(buf *bytes.Buffer) error {
	if b.dirtyFlags == 0 {
		return ErrNoDirtyData
	}

	if err := writeBinary(buf, b.dirtyFlags); err != nil {
		return err
	}

	if b.dirtyFlags&Bit00 != 0 {
		if err := writeBinary(buf, b.i); err != nil {
			return err
		}
	}

	if b.dirtyFlags&Bit01 != 0 {
		if _, err := buf.WriteString(b.str); err != nil {
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

	if dirtyFlags&Bit00 != 0 {
		i := b.i
		if err := readBinary(&buf, &b.i); err != nil {
			b.i = i
			return err
		}
	}

	if dirtyFlags&Bit01 != 0 {
		by, err := buf.ReadBytes(0)
		if err != nil {
			return err
		}
		b.str = string(by[:len(by)-1])
	}

	return nil
}
