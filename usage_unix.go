//go:build unix

package main

import (
	"github.com/pkg/errors"
	"syscall"
	"time"
)

type Usage struct {
	UserTime   time.Duration
	SystemTime time.Duration
	MaxRSS     int64
}

func GetChildrenUsage() (_ Usage, implemented bool, err error) {
	usage := syscall.Rusage{}
	err = syscall.Getrusage(syscall.RUSAGE_CHILDREN, &usage)
	return Usage{
		UserTime:   time.Duration(usage.Utime.Nano()) * time.Nanosecond,
		SystemTime: time.Duration(usage.Stime.Nano()) * time.Nanosecond,
		MaxRSS:     int64(usage.Maxrss),
	}, true, errors.WithStack(err)
}
