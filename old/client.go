package client

import (
	"encoding/json"
	"net"
)

const (
	EOL = "\n"
)

type Client interface {
	Connect() error
	Disconnect() error

	Authorizate() error
	TankCommand(message Message) error
	SendMessage(message Message) error

	ReadMessage() (Message, error)
	ReadType(string) (Message, error)
	ReadWorld() (Message, error)
	ReadTank() (TankMessage, error)
}

type client struct {
	ServerAddr string

	Login    string
	Password string

	connection *net.TCPConn
	jsonDec    *json.Decoder
}

func NewClient(login, pass, serverAddr string) Client {
	return &client{
		ServerAddr: serverAddr,
		Login:      login,
		Password:   pass,
	}
}

func (c *client) Connect() error {

	addr, err := net.ResolveTCPAddr("tcp", c.ServerAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	c.connection = conn
	c.jsonDec = json.NewDecoder(c.connection)

	return nil
}

func (c *client) Disconnect() error {
	return c.connection.Close()
}

func (c *client) Authorizate() error {
	message := NewAuthMessage(c.Login, c.Password)
	err := c.SendMessage(message)
	return err
}

func (c *client) TankCommand(message Message) error {
	message["Type"] = "TankCommand"
	err := c.SendMessage(message)
	return err
}

func (c *client) ReadMessage() (Message, error) {
	var m Message

	if err := c.jsonDec.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *client) ReadType(typeStr string) (Message, error) {
	for {
		m, err := c.ReadMessage()
		if err != nil {
			return nil, err
		}
		if m["Type"] == typeStr {
			return m, nil
		}
	}
}

func (c *client) ReadWorld() (Message, error) {
	message, err := c.ReadType("World")

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (c *client) ReadTank() (TankMessage, error) {
	message, err := c.ReadType("Tank")
	if err != nil {
		return nil, err
	}

	return TankMessage(message), nil
}

func (c *client) SendMessage(message Message) error {
	jsonStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	jsonStr = append(jsonStr, []byte(EOL)...)

	for len(jsonStr) > 0 {
		n, err := c.connection.Write(jsonStr)
		if err != nil {
			return err
		}
		jsonStr = jsonStr[n:]
	}

	return nil
}
