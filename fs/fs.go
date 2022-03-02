package fs

import (
	"math"
)

type DirEntry struct {
	Path               string
	Size               int64
	Deleted            bool
	DeletionInProgress bool
}

func pow(n float64) int64 {
	return int64(math.Pow(1024, n))
}

var dirs = []DirEntry{
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\old\\build", Size: pow(3.1)},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\blog\\build", Size: pow(1.2)},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\practise\\build", Size: pow(2.76)},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\example\\build", Size: pow(3.001)},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\blog-old\\build", Size: pow(0.11)},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\website\\build", Size: pow(0.5)},
}

func Traverse(searchPath string, searchDirs, ignoreDirs []string, searchAll bool, ch chan DirEntry) {
	defer close(ch)
	for _, dir := range dirs {
		ch <- dir
	}
}

func Delete(path string) error {
	return nil
}
