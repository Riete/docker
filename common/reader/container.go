package reader

import (
	"bufio"
	"bytes"
	"io"

	"github.com/docker/docker/pkg/stdcopy"
)

func ParseToStdoutStderr(r io.Reader) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	_, err := stdcopy.StdCopy(&stdout, &stderr, r)
	return stdout.String(), stderr.String(), err
}

func ParseToCombinedOutput(r io.Reader) (string, error) {
	var b bytes.Buffer
	_, err := stdcopy.StdCopy(&b, &b, r)
	return b.String(), err
}

func ParseToCombinedStreamOutput(r io.Reader) <-chan string {
	ch := make(chan string)
	rb := bufio.NewReader(r)
	b := &bytes.Buffer{}

	go func() {
		defer close(ch)
		for {
			d, err := rb.ReadBytes('\n')
			if _, err := b.Write(d); err != nil {
				ch <- err.Error()
				return
			}
			if s, err := ParseToCombinedOutput(b); err != nil {
				ch <- err.Error()
				return
			} else {
				ch <- s
			}
			if err == io.EOF {
				return
			}
			if err != nil {
				ch <- err.Error()
				return
			}
		}
	}()
	return ch
}
