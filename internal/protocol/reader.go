// RESP stream parser (frame-by-frame)
package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const maxBulk = 512 << 20

var ErrProtocol = errors.New("Protocol error")

type Reader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		br: bufio.NewReader(r),
	}
}

func (r *Reader) ReadCommand() (Command, error) {
	for {
		b, err := r.br.Peek(1)
		if err != nil {
			return Command{}, err
		}

		if b[0] == '*' {
			v, err := r.ReadValue()
			if err != nil {
				return Command{}, err
			}
			if v.Kind != KindArray {
				return Command{}, fmt.Errorf("%w: command must be an array", ErrProtocol)
			}
			args := make([]string, len(v.Array))
			for i, el := range v.Array {
				if el.Kind != KindBulk {
					return Command{}, fmt.Errorf("%w: command arguments must be bulk strings", ErrProtocol)
				}
				args[i] = el.Str
			}
			if len(args) == 0 {
				continue
			}
			return Command{Name: strings.ToUpper(args[0]), Args: args[1:]}, nil
		}
		return Command{}, fmt.Errorf("%w: invalid command prefix %c", ErrProtocol, b[0])
	}
}

// ReadValue reads any single RESP value. The CLI uses this to read replies.
func (r *Reader) ReadValue() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	if line == "" {
		return Value{}, fmt.Errorf("%w: empty line", ErrProtocol)
	}

	switch line[0] {
	case '+':
		return SimpleString(line[1:]), nil
	case '-':
		return NewError(line[1:]), nil
	case ':':
		n, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return Value{}, fmt.Errorf("%w: Bad Integer %q", ErrProtocol, line)
		}
		return Integer(n), nil
	case '$':
		n, err := parseLen(line[1:])
		if err != nil {
			return Value{}, err
		}
		if n < 0 {
			return Nil(), nil
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(r.br, buf); err != nil {
			return Value{}, err
		}
		if err := r.expectCRLF(); err != nil {
			return Value{}, err
		}
		return Bulk(string(buf)), nil
	case '*':
		n, err := parseLen(line[1:])
		if err != nil {
			return Value{}, err
		}
		if n < 0 {
			return Nil(), nil
		}
		items := make([]Value, 0, min(n, 1024))
		for i := 0; i < n; i++ {
			el, err := r.ReadValue()
			if err != nil {
				return Value{}, err
			}
			items = append(items, el)
		}
		return Value{Kind: KindArray, Array: items}, nil
	default:
		return Value{}, fmt.Errorf("%w: unexpected byte %q", ErrProtocol, line[0])
	}
}

func parseLen(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < -1 || n > maxBulk {
		return 0, fmt.Errorf("%w: invalid length %q", ErrProtocol, s)
	}
	return n, nil
}

// readLine reads one \r\n-terminated line and returns it without the terminator.
func (r *Reader) readLine() (string, error) {
	line, err := r.br.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return line, nil
}

func (r *Reader) expectCRLF() error {
	b1, err := r.br.ReadByte()
	if err != nil {
		return err
	}
	b2, err := r.br.ReadByte()
	if err != nil {
		return err
	}
	if b1 != '\r' || b2 != '\n' {
		return fmt.Errorf("%w: expected CRLF after bulk payload", ErrProtocol)
	}
	return nil
}
