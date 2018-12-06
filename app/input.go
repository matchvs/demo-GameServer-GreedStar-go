package app

type ClientInput struct {
	Left  int `json:"l"`
	Right int `json:"r"`
	Up    int `json:"u"`
	Down  int `json:"d"`
	Speed int `json:"p"`
}
