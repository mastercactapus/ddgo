package main

import (
	"github.com/cheggaaa/pb"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		logrus.Fatalln("Usage: ddgo <src> <dest>")
	}
	src := os.Args[1]
	dst := os.Args[2]
	fd, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		logrus.Fatalln("failed to open src:", err)
	}
	defer fd.Close()
	srcLen, err := fd.Seek(0, 2)
	if err != nil {
		logrus.Fatalln(err)
	}
	_, err = fd.Seek(0, 0)
	if err != nil {
		logrus.Fatalln(err)
	}

	bar := pb.New64(srcLen)
	bar.SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true

	out, err := os.OpenFile(dst, os.O_WRONLY|syscall.O_DIRECT, 0777)
	if err != nil {
		logrus.Fatalln("failed to open dst:", err)
	}
	defer out.Close()

	dstLen, err := out.Seek(0, 2)
	if err != nil {
		logrus.Fatalln(err)
	}
	_, err = out.Seek(0, 0)
	if err != nil {
		logrus.Fatalln(err)
	}

	if dstLen < srcLen {
		bar.Total = dstLen
	}

	bar.Start()
	if dstLen < srcLen {
		logrus.Warnln("destination device too small, not all bytes will be copied")
		_, err = io.Copy(io.MultiWriter(bar, out), io.LimitReader(fd, dstLen))
	} else {
		_, err = io.Copy(io.MultiWriter(bar, out), fd)
	}
	if err != nil {
		logrus.Fatalln("copy failed:", err)
	}
	bar.Finish()
	fd.Close()
	out.Close()
}
