package client

type Message map[string]interface{}

func NewMessage() Message {
	return Message{}
}

func NewAuthMessage(login, pass string) Message {
	return Message{
		"Type":     "Auth",
		"Login":    login,
		"Password": pass,
	}
}

func (m Message) Motors(left, right float64) Message {
	m["LeftMotor"] = left
	m["RightMotor"] = right
	return m
}

func (m Message) Fire() Message {
	m["Fire"] = true
	return m
}

func (m Message) TurnGun(angle float64) Message {
	m.Gun()["TurnAngle"] = angle
	return m
}

func (m Message) Gun() Message {
	var gun Message
	gun, ok := m["Gun"].(Message)
	if !ok {
		gun = Message{}
		m["Gun"] = gun
	}
	return gun
}
