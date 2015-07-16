package fos

import (
	"errors"
	"os"
)

var (
	ErrTruncToSmall       = errors.New("can't truncate to a smaller file")
	ErrReadAtInvalidParam = errors.New("invalid param to ReadAt")
)

type File struct {
	b      []byte
	offset int64
}

func (f *File) Fd() int {
	return 0
}

func (f *File) Truncate(size int64) error {
	if int64(f.length()) > size {
		return ErrTruncToSmall
	}

	tmp := make([]byte, size)
	copy(tmp, f.b)
	f.b = tmp

	return nil
}

func (f *File) length() int {
	return len(f.b)
}

func (f *File) Sync() error {
	return nil
}

func OpenFile(path string, flag int, mode os.FileMode) (*File, error) {
	return &File{}, nil
}

type Buffer struct {
	b      []byte
	offset int64
}

func (p *File) WriteAt(buf []byte, off int64) (n int, err error) {
	ioff := int(off)
	iend := ioff + len(buf)
	if len(p.b) < iend {
		if len(p.b) == ioff {
			p.b = append(p.b, buf...)
			return len(buf), nil
		}
		zero := make([]byte, iend-len(p.b))
		p.b = append(p.b, zero...)
	}
	copy(p.b[ioff:], buf)
	return len(buf), nil
}

type Stat struct {
	file *File
}

func (f *File) Stat() (*Stat, error) {
	return &Stat{f}, nil
}

func (s *Stat) Size() int {
	return s.file.length()
}

func (f *File) ReadAt(buf []byte, offset int) (int, error) {
	if offset != 0 {
		return 0, ErrReadAtInvalidParam
	}
	min := func(a, b int) int {
		if a < b {
			return a
		} else {
			return b
		}
	}
	n := min(f.length(), len(buf))
	readSize := copy(buf, f.b[:n])
	return readSize, nil
}

func (f *File) Close() error {
	return nil
}

func (f *File) Buffer() []byte {
	return f.b
}
