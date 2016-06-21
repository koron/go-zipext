package zipext

import (
	"archive/zip"
	"encoding/binary"
	"errors"
)

var (
	// ErrTooShortField is there are no enough data for an extra field.
	ErrTooShortField = errors.New("too short field")

	// ErrSizeMismatch is that detect size mismatch for an extra field.
	ErrSizeMismatch = errors.New("size mismatch")
)

// Field represents extra field of zip file.
type Field struct {
	Tag  uint16
	Data []byte
}

type readBuf []byte

func (b *readBuf) uint16() uint16 {
	v := binary.LittleEndian.Uint16(*b)
	*b = (*b)[2:]
	return v
}

func (b *readBuf) uint32() uint32 {
	v := binary.LittleEndian.Uint32(*b)
	*b = (*b)[4:]
	return v
}

func (b *readBuf) uint64() uint64 {
	v := binary.LittleEndian.Uint64(*b)
	*b = (*b)[8:]
	return v
}

// Reader reads extend fields in zip file.
type Reader struct {
	err error

	buf readBuf
}

// NewReader creates a reader.
func NewReader(zf *zip.File) *Reader {
	return &Reader{
		buf: readBuf(zf.Extra),
	}
}

func (r *Reader) Read() (*Field, error) {
	if r.err != nil {
		return nil, r.err
	}
	remain := len(r.buf)
	if remain == 0 {
		return nil, nil
	}
	if remain < 4 {
		r.err = ErrTooShortField
		return nil, r.err
	}
	tag := r.buf.uint16()
	size := r.buf.uint16()
	if int(size) > len(r.buf) {
		r.err = ErrSizeMismatch
		return nil, r.err
	}
	return &Field{
		Tag:  tag,
		Data: r.readBytes(int(size)),
	}, nil
}

func (r *Reader) readUint16() uint16 {
	v := binary.LittleEndian.Uint16(r.buf)
	r.buf = r.buf[2:]
	return v
}

func (r *Reader) readBytes(size int) []byte {
	v := r.buf[:size]
	r.buf = r.buf[size:]
	return v
}
