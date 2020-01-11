package samples

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Structs", func() {
	var ()

	BeforeEach(func() {

	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				b := B{
					i: 1,
					s: "s",
				}
				// default
				Expect(b.I()).To(Equal(int64(1)))
				Expect(b.S()).To(Equal("s"))
				var buf bytes.Buffer
				err := b.Dirty(&buf)
				Expect(buf.Bytes()).To(BeEmpty())
				Expect(err).To(MatchError(ErrNoDirtyData))

				// change i
				b.SetI(2)
				Expect(b.I()).To(Equal(int64(2)))

				buf.Reset()
				err = b.Dirty(&buf)
				Expect(buf.Bytes()).NotTo(BeEmpty())
				Expect(err).To(BeNil(), err)

				var b2 B
				Expect(b2.MergeFrom(buf.Bytes())).To(BeNil())
				Expect(b2.I()).To(Equal(b.I()))
				Expect(b2.S()).To(BeZero())

				// change s
				b.SetS("ss")
				Expect(b.S()).To(Equal("ss"))
				buf.Reset()
				err = b.Dirty(&buf)
				Expect(buf.Bytes()).NotTo(BeEmpty())
				Expect(err).To(BeNil(), err)

				var b3 B
				Expect(b3.MergeFrom(buf.Bytes())).To(BeNil())
				Expect(b3.I()).To(Equal(b.I()))
				Expect(b3.S()).To(Equal(b.S()))

				// clean dirty flags
				b.CleanDirty()
				buf.Reset()
				err = b.Dirty(&buf)
				Expect(err).To(MatchError(ErrNoDirtyData))
				Expect(buf.Bytes()).To(BeEmpty())
			})
		})

		Context("Nesting struct", func() {
			It("Nesting struct", func() {
				d := D{
					b:  B{i: 1, s: "s"},
					i2: 2,
					s2: "s2",
				}
				d.Init()
				// default
				Expect(d.B().I()).To(Equal(int64(1)))
				Expect(d.B().S()).To(Equal("s"))
				Expect(d.I2()).To(Equal(int64(2)))
				Expect(d.S2()).To(Equal("s2"))

				var buf bytes.Buffer
				err := d.Dirty(&buf)
				Expect(buf.Bytes()).To(BeEmpty())
				Expect(err).To(MatchError(ErrNoDirtyData))

				// change child
				buf.Reset()
				b := d.B()
				b.SetI(11)
				fmt.Println("change child", d)
				err = d.Dirty(&buf)
				Expect(err).To(BeNil())

				var d2 D
				Expect(d2.MergeFrom(buf.Bytes())).To(BeNil())
				//Expect(d.B().S()).To(Equal("s"))
			})
		})
	})
})
