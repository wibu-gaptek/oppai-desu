package utils

import (
	"io"

	"github.com/docker/docker/pkg/archive"
)

func TarWithOpt(src string) (io.ReadCloser, error) {

	tar, err := archive.TarWithOptions(src, &archive.TarOptions{})
	if err != nil {
		return nil, err
	}
	return tar, nil
}
