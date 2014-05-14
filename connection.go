package client

import (
	"encoding/json"
	"io"
)

const (
	EOL = "\n"
)

type connection struct {
	rw      io.ReadWriteCloser
	decoder *json.Decoder
}

func newConnection(rw io.ReadWriteCloser) *connection {
	return &connection{rw, json.NewDecoder(rw)}
}

func (c *connection) sendMessage(message interface{}) error {
	jsonStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	jsonStr = append(jsonStr, []byte(EOL)...)

	for len(jsonStr) > 0 {
		n, err := c.rw.Write(jsonStr)
		if err != nil {
			return err
		}
		jsonStr = jsonStr[n:]
	}

	return nil
}

func (c *connection) readMessage() (*Message, error) {
	var m Message

	if err := c.decoder.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *connection) Close() error {
	return c.rw.Close()
}
