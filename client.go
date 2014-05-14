package client

import (
	"io"
	"net"
)

type Client interface {
	ReadMessage() (*Message, error)

	Auth(login, pass string) error
	Fire() error
	Motors(left, right float64) error
	TurnGun(angle float64) error

	WorldFrequency(n int64) error
	Ping(id int64) error

	Disconnect() error
}

type client struct {
	conn *connection
}

func ConnectTo(addrStr string) (Client, error) {

	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

func NewClient(rw io.ReadWriteCloser) Client {
	return &client{newConnection(rw)}
}

func (c *client) ReadMessage() (*Message, error) {
	return c.conn.readMessage()
}

func (c *client) WorldFrequency(n int64) error {
	return c.SendClientCommand(Message{WorldFrequency: n})
}

func (c *client) Ping(id int64) error {
	return c.conn.sendMessage(&Message{
		Type:   "Ping",
		PingId: id,
	})
}

func (c *client) TurnGun(angle float64) error {
	return c.SendTankCommand(Message{Gun: Gun{TurnAngle: angle}})
}

func (c *client) Motors(left, right float64) error {
	return c.SendTankCommand(Message{LeftMotor: left, RightMotor: right})
}

func (c *client) Fire() error {
	return c.SendTankCommand(Message{Fire: true})
}

func (c *client) SendClientCommand(command Message) error {
	command.Type = "Client"
	return c.conn.sendMessage(&command)
}

func (c *client) SendTankCommand(command Message) error {
	command.Type = "TankCommand"
	return c.conn.sendMessage(&command)
}

func (c *client) Auth(login, pass string) error {
	return c.conn.sendMessage(&Message{
		Type:     "Auth",
		Login:    login,
		Password: pass,
	})
}

func (c *client) Disconnect() error {
	return c.conn.Close()
}
