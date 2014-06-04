package client

type Message struct {
	Type    string
	Message string `json:",omitempty"`

	//Auth
	Login    string `json:",omitempty"`
	Password string `json:",omitempty"`

	//TankCommand
	LeftMotor  float64 `json:",omitempty"`
	RightMotor float64 `json:",omitempty"`
	Gun        *Gun    `json:",omitempty"`
	Fire       bool    `json:",omitempty"`

	//Tank
	Tank *Tank `json:",omitempty"`

	//World
	Tanks []Tank `json:",omitempty"`
	Map   *Map   `json:",omitempty"`

	//Client
	WorldFrequency int64 `json:",omitempty"`

	//Ping
	PingId int64 `json:",omitempty"`
}

type Map struct {
	Width  int64 `json:",omitempty"`
	Height int64 `json:",omitempty"`
}

type Tank struct {
	Id int64 `json:",omitempty"`

	Coords Coords `json:",omitempty"`
	Gun    Gun    `json:",omitempty"`

	Direction  float64 `json:",omitempty"`
	LeftMotor  float64 `json:",omitempty"`
	RightMotor float64 `json:",omitempty"`

	Radius int64 `json:",omitempty"`
	Health int64 `json:",omitempty"`
}

type Coords struct {
	X float64 `json:",omitempty"`
	Y float64 `json:",omitempty"`
}

type Gun struct {
	Direction      float64 `json:",omitempty"`
	ReloadProgress float64 `json:",omitempty"`
	Temperature    float64 `json:",omitempty"`
	TurnAngle      float64 `json:",omitempty"`
}

func NewMessage() Message {
	return Message{}
}

func (m Message) Motors(left, right float64) Message {
	m.LeftMotor = left
	m.RightMotor = right
	return m
}

func (m Message) SetFire() Message {
	m.Fire = true
	return m
}

func (m Message) TurnGun(angle float64) Message {
	m.Gun = &Gun{TurnAngle: angle}
	return m
}

func (m Message) IsWorld() bool {
	return m.Type == "World"
}
