package fsnapshot

import "regexp"

// FileStatus constants
const (
	fNew        = 1
	fChanged    = 2
	fDeleted    = 3
	fUnmodified = 4
)

// FileData contains file information that needs to be saved
type FileData struct {
	Filename         string
	FilenameChecksum []byte
	Subdirectory     string
	Checksum         []byte
	FileStatus       int
	seen             bool
}

func filenameChecksum(filename []byte) []byte {
	t := filename
	r1 := regexp.MustCompile("/")
	r2 := regexp.MustCompile(`\\`)
	r1.ReplaceAll(t, []byte("_"))
	r2.ReplaceAll(t, []byte("_"))
	return checksum(t)
}
