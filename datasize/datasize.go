package datasize

import (
	"fmt"
	"strings"
)

type ByteSize int64

const (
	B ByteSize = 1

	// for computer storage
	StorageSpan ByteSize = 1000
	KB          ByteSize = StorageSpan * B  // kilobyte
	MB          ByteSize = StorageSpan * KB // megabyte
	GB          ByteSize = StorageSpan * MB // gigabyte
	TB          ByteSize = StorageSpan * GB // terabyte
	PB          ByteSize = StorageSpan * TB // petabyte
	EB          ByteSize = StorageSpan * PB // exabyte

	// for computer memory
	MemorySpan ByteSize = 1024
	KiB        ByteSize = MemorySpan * B   // kibibyte
	MiB        ByteSize = MemorySpan * KiB // mebibyte
	GiB        ByteSize = MemorySpan * MiB // gibibyte
	TiB        ByteSize = MemorySpan * GiB // tebibyte
	PiB        ByteSize = MemorySpan * TiB // pebibyte
	EiB        ByteSize = MemorySpan * PiB // exbibyte
)

// New new a ByteSize from byte size
func New(n int64) ByteSize {
	return ByteSize(n)
}

func (bs ByteSize) KB() float64 {
	return float64(bs) / float64(KB)
}

func (bs ByteSize) MB() float64 {
	return float64(bs) / float64(MB)
}

func (bs ByteSize) GB() float64 {
	return float64(bs) / float64(GB)
}

func (bs ByteSize) TB() float64 {
	return float64(bs) / float64(TB)
}

func (bs ByteSize) PB() float64 {
	return float64(bs) / float64(PB)
}

func (bs ByteSize) EB() float64 {
	return float64(bs) / float64(EB)
}

// float2String float64 to string,
func float2String(value float64) string {
	s := fmt.Sprintf("%.3f", value)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

// String1 show unit by StorageSpan(1000)
func (bs ByteSize) String1() string {
	if bs >= EB {
		return fmt.Sprintf("%sEB", float2String(bs.EB()))
	}
	if bs >= PB {
		return fmt.Sprintf("%sPB", float2String(bs.PB()))
	}
	if bs >= TB {
		return fmt.Sprintf("%sTB", float2String(bs.TB()))
	}
	if bs >= GB {
		return fmt.Sprintf("%sGB", float2String(bs.GB()))
	}
	if bs >= MB {
		return fmt.Sprintf("%sMB", float2String(bs.MB()))
	}
	if bs >= KB {
		return fmt.Sprintf("%sKB", float2String(bs.KB()))
	}
	return fmt.Sprintf("%d", bs)
}

func (bs ByteSize) KiB() float64 {
	return float64(bs) / float64(KiB)
}

func (bs ByteSize) MiB() float64 {
	return float64(bs) / float64(MiB)
}

func (bs ByteSize) GiB() float64 {
	return float64(bs) / float64(GiB)
}

func (bs ByteSize) TiB() float64 {
	return float64(bs) / float64(TiB)
}

func (bs ByteSize) PiB() float64 {
	return float64(bs) / float64(PiB)
}

func (bs ByteSize) EiB() float64 {
	return float64(bs) / float64(EiB)
}

// String1 show unit by MemorySpan(1024)
func (bs ByteSize) String2() string {
	if bs > EiB {
		return fmt.Sprintf("%sEiB", float2String(bs.EiB()))
	}
	if bs > PiB {
		return fmt.Sprintf("%sPiB", float2String(bs.PiB()))
	}
	if bs > TiB {
		return fmt.Sprintf("%sTiB", float2String(bs.TiB()))
	}
	if bs > GiB {
		return fmt.Sprintf("%sGiB", float2String(bs.GiB()))
	}
	if bs > MiB {
		return fmt.Sprintf("%sMiB", float2String(bs.MiB()))
	}
	if bs > KiB {
		return fmt.Sprintf("%sKiB", float2String(bs.KiB()))
	}
	return fmt.Sprintf("%d", bs)
}
