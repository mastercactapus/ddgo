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
	stat, err := os.Stat(src)
	if err != nil {
		logrus.Fatalln("stat failed:", err)
	}
	fd, err := os.Open(src)
	if err != nil {
		logrus.Fatalln("failed to open src:", err)
	}
	bar := pb.New64(stat.Size())
	bar.SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true

	out, err := os.OpenFile(dst, os.O_WRONLY|syscall.O_DIRECT, 0777)
	if err != nil {
		logrus.Fatalln("failed to open dst:", err)
	}
	bar.Start()
	_, err = io.Copy(io.MultiWriter(bar, out), fd)
	if err != nil {
		logrus.Fatalln("copy failed:", err)
	}
	bar.Finish()
	fd.Close()
	out.Close()
}
