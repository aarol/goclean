package fs

import "time"

type DirEntry struct {
	Path               string
	Size               int64
	Deleted            bool
	DeletionInProgress bool
}

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
	tb = gb * 1024
)

var dirs = []DirEntry{
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\old\\build", Size: 3 * kb},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\blog\\build", Size: 19 * mb},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\practise\\build", Size: 1 * gb},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\example\\build", Size: 3 * gb},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\blog-old\\build", Size: 0},
	{Path: "C:\\Users\\aarol\\Documents\\Code\\projects\\website\\build", Size: 600 * kb},
}

func Traverse(searchPath string, searchDirs, ignoreDirs []string, searchAll bool, ch chan DirEntry) {
	defer close(ch)
	for _, dir := range dirs {
		ch <- dir
	}
}

func Delete(path string) error {
	time.Sleep(2 * time.Second)
	return nil
}
