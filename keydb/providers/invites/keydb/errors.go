package keydb

import (
	"github.com/go-faster/errors"
)

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists. Use UpdateLink instead AddLink")
)
