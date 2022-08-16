package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := os.Args[1]
	srcStat, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	models, err := ParseBigmodelInterface(filename)
	if err != nil {
		panic(err)
	}
	for i := range models {
		outname := strings.Replace(filename, ".go", "_genimpl.go", 1)
		if err := renderModelImpl(&models[i], false, outname, srcStat.Mode().Perm()); err != nil {
			panic(err)
		}
	}
}

func x(bs []byte) []byte {
	lines := strings.Split(string(bs), "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%3d|%s", i+1, line)
	}
	return []byte(strings.Join(lines, "\n"))
}
