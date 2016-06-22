package zipext

import (
	"archive/zip"
	"errors"
	"fmt"
	"time"
)

// Extra is accessor for zip.File#Extra information.
type Extra struct {
	zf   *zip.File
	perr []error

	// extended timestamp (0x5455)
	extime struct {
		m time.Time
		a time.Time
		c time.Time
	}
}

// ModTime returns the modication time in UTC. The resolution is 1s.
func (ex *Extra) ModTime() time.Time {
	if !ex.extime.m.IsZero() {
		return ex.extime.m
	}
	return ex.zf.ModTime()
}

// AcTime returns the access time in UTC. The resolution is 1s.
func (ex *Extra) AcTime() time.Time {
	return ex.extime.a
}

// CrTime returns the creation time in UTC. The resolution is 1s.
func (ex *Extra) CrTime() time.Time {
	return ex.extime.c
}

// Parse parses zip.File#Extra decodes extra information.
func Parse(zf *zip.File) *Extra {
	ex := &Extra{
		zf: zf,
	}
	r := NewReader(zf)
	for {
		f, err := r.Read()
		switch {
		case err != nil:
			ex.perr = append(ex.perr, err)
			fallthrough
		case f == nil:
			return ex
		}
		if err := ex.procField(f); err != nil {
			ex.perr = append(ex.perr, err)
		}
	}
}

func (ex *Extra) procField(f *Field) error {
	r := f.readBuf()
	switch f.Tag {

	case 0x5455: // extended timestamp.
		if len(r) < 1 {
			return errors.New("too short extended timestamp")
		}
		flag := r.uint8()
		if flag&0x01 != 0 && len(r) >= 4 {
			ex.extime.m = time.Unix(int64(r.uint32()), 0)
		}
		if flag&0x02 != 0 && len(r) >= 4 {
			ex.extime.a = time.Unix(int64(r.uint32()), 0)
		}
		if flag&0x04 != 0 && len(r) >= 4 {
			ex.extime.c = time.Unix(int64(r.uint32()), 0)
		}
		return nil

	// TODO: parse other tag.

	default:
		return fmt.Errorf("unsupported tag: 0x%04x", f.Tag)
	}
}
