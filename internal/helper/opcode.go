package helper

type OpCode int

const (
	OpUndefined OpCode = iota
	CreateFile
	CreateUrl
	CreateDir
	CreateText
	ReadFile
	ReadUrl
	ReadDir
	ReadText
	UpdateFile
	UpdateUrl
	UpdateDir
	UpdateText
	DeleteFile
	DeleteUrl
	DeleteDir
	DeleteText
)
