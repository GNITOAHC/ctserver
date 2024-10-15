package helper

import "errors"

type UserOp struct {
	Mail      string // username
	Operation OpCode
	Path      string // Should start with /
	Ancestor  string // ancestor_id

	// For file
	File []byte

	// For url, text & dir name
	Content string
}

func (h *Helper) UserOperation(u UserOp) (string, error) {
	if u.Mail == "" || u.Operation == OpUndefined || u.Path == "" || u.Ancestor == "" {
		return "", errors.New("Mail, OpCode and Path are required")
	}

	switch u.Operation {
	case CreateFile:
	case CreateUrl:
	case CreateDir:
	case CreateText:
	default:
	}

	return "", nil
}
