package samples

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Structs", func() {
	var (
		b B
		//d D
	)

	BeforeEach(func() {
		b = B{
			i: 1,
			s: "s",
		}

		//d = D{
		//	i2: 2,
		//	s2: "s2",
		//}
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				// default
				Expect(b.I()).To(Equal(int64(1)))
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

		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				//Expect(d.B().S()).To(Equal("s"))
			})
		})
	})
})
