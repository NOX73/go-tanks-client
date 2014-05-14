package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var host = "nox73.ru:9292"
var login = "1"
var pass = "1"

func TestConnect(t *testing.T) {
	var err error
	c := NewClient(login, pass, host)

	err = c.Connect()
	assert.Nil(t, err)

	err = c.Authorizate()
	assert.Nil(t, err)

	for i := 0; i < 10; i++ {
		message, err := c.ReadMessage()
		if assert.NoError(t, err) {
			assert.NotNil(t, message["Type"])
		}
	}

	c.Disconnect()
}

func TestReadMethods(t *testing.T) {
	c := NewClient(login, pass, host)

	c.Connect()
	c.Authorizate()

	tank, err := c.ReadTank()
	if assert.NoError(t, err) {
		assert.Equal(t, tank["Type"], "Tank", "Type of tank message is not 'Tank'")
		assert.NotNil(t, tank["Tank"], "Tank message hasn't Tank field")
	}

	message, err := c.ReadWorld()
	if assert.NoError(t, err) {
		assert.Equal(t, message["Type"], "World", "Type of world message is not 'World'")
		assert.NotNil(t, message["Tanks"], "World message hasn't Tanks field")
	}

	c.Disconnect()
}

func TestTankCommand(t *testing.T) {
	var left = 0.4
	var right = 0.5

	c := NewClient(login, pass, host)

	c.Connect()
	c.Authorizate()

	selfTank, err := c.ReadTank()
	assert.NoError(t, err, "Error on read tank.")

	m := NewMessage().Motors(left, right)
	c.TankCommand(m)

	for i := 0; i < 30; i++ {
		c.ReadMessage()
	}

	message, err := c.ReadWorld()
	tanksArr := message["Tanks"].([]interface{})

	if assert.NoError(t, err) {
		tanks := make([]Message, len(tanksArr))
		for i, t := range tanksArr {
			tanks[i] = Message(t.(map[string]interface{}))
		}

		var tank Message
		for _, tank = range tanks {
			if selfTank.Id() == tank["Id"] {
				break
			}
		}

		assert.Equal(t, tank["LeftMotor"], left, "Left motor is wrong")
		assert.Equal(t, tank["RightMotor"], right, "Right motor is wrong")
	}

	c.Disconnect()
}
