package client

type TankMessage map[string]interface{}

func (m TankMessage) Id() int64 {
	return int64(m.Tank()["Id"].(float64))
}

func (m TankMessage) Tank() Message {
	return Message(m["Tank"].(map[string]interface{}))
}
