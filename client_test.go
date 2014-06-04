package client

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var host = "nox73.ru:9292"
var login = "1"
var pass = "1"

func AuthReader() io.Reader {
	return strings.NewReader(`{"Message":"Hello! You should authorize befor join the game!","Type":"Auth"}` + "\n")
}

func TankReader() io.Reader {
	return strings.NewReader(`{"Message":"Your tank id = 4035","Tank":{"Id":4035,"Coords":{"X":604,"Y":188},"Direction":0,"LeftMotor":0,"RightMotor":0,"Gun":{"Direction":0,"ReloadProgress":0,"Temperature":0},"Radius":10,"Health":1},"Type":"Tank"}` + "\n")
}

func WorldReader() io.Reader {
	return strings.NewReader(`{"Bullets":[],"Id":558528599,"Map":{"Width":1024,"Height":768},"Tanks":[{"Id":4034,"Coords":{"X":860,"Y":449},"Direction":0,"LeftMotor":0,"RightMotor":0,"Gun":{"Direction":0,"ReloadProgress":0,"Temperature":0},"Radius":10,"Health":1}],"Type":"World"}` + "\n")
}

type MockConn struct {
	r io.Reader
	mock.Mock
}

func (m *MockConn) Read(p []byte) (int, error) {
	return m.r.Read(p)
}

func (m *MockConn) Write(p []byte) (int, error) {
	return len(p), nil
}

func (m *MockConn) Close() error {
	return nil
}

func TestConnect(t *testing.T) {
	c, err := ConnectTo(host)
	assert.NoError(t, err)

	err = c.Disconnect()
	assert.NoError(t, err)
}

func TestAuth(t *testing.T) {
	m := &MockConn{r: io.MultiReader(AuthReader(), TankReader(), WorldReader())}

	c := NewClient(m)

	message, err := c.ReadMessage()

	assert.NoError(t, err)
	assert.Equal(t, "Auth", message.Type, "type should be 'Auth'")

	err = c.Auth(login, pass)
	assert.NoError(t, err)

	message, err = c.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "Tank", message.Type, "type should be 'Tank'")
	assert.NotNil(t, message.Tank, "Tank in TankMessage should not be nil")
	assert.Equal(t, message.Tank.Id, 4035, "Tank.Id in TankMessage wrong")

	message, err = c.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "World", message.Type, "type should be 'World'")
	assert.Equal(t, true, message.IsWorld(), "IsWorld should be true")
	assert.NotNil(t, message.Map, "Map in WorldMessage should not be nil")
	assert.Equal(t, 1024, message.Map.Width, "Map.Width should be 1024")
	assert.Equal(t, 768, message.Map.Height, "Map.Height should be 768")

	assert.NotNil(t, message.Tanks, "Tanks in WorldMessage should be array of Tanks")
	assert.Equal(t, 1, len(message.Tanks), "Tanks count should be equal 1")
	assert.Equal(t, 0, message.Tanks[0].Gun.Direction, "First Tank direction should be 0")
	assert.Equal(t, 860, message.Tanks[0].Coords.X)

	c.Disconnect()
}

func TestMessage(t *testing.T) {
	m := NewMessage()

	left := 21.0
	right := 31.0

	angle := 324.0

	m = m.Motors(left, right).SetFire().TurnGun(angle)

	assert.Equal(t, left, m.LeftMotor)
	assert.Equal(t, right, m.RightMotor)
	assert.Equal(t, true, m.Fire)
	assert.Equal(t, angle, m.Gun.TurnAngle)
}
