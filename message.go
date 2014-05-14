package client

type Message struct {
	Type    string
	Message string

	//Auth
	Login    string
	Password string

	//TankCommand
	LeftMotor  float64
	RightMotor float64
	Gun        Gun
	Fire       bool

	//Tank
	Tank Tank

	//World
	Tanks []Tank
	Map   Map

	//Client
	WorldFrequency int64

	//Ping
	PingId int64
}

type Map struct {
	Width  int64
	Height int64
}

type Tank struct {
	Id int64

	Coords Coords
	Gun    Gun

	Direction  float64
	LeftMotor  float64
	RightMotor float64

	Radius int64
	Health int64
}

type Coords struct {
	X int64
	Y int64
}

type Gun struct {
	Direction      float64
	ReloadProgress float64
	Temperature    float64
	TurnAngle      float64
}
