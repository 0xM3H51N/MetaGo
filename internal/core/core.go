package core

type FileMeta struct {
	Name    string
	Size    int64
	Hash    string
	ModTime string
}

type Config struct {
	FilePath  string
	DirPath   string
	Json      bool
	HashType  string
	Recursive bool
}
