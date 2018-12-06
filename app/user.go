/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:34
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-06 19:17:39
 * @Description: file content
 */
package app

import (
	"fmt"
)

type GameUserSlice []*GameUser

func NewGameUserSlice() GameUserSlice {
	u := make(GameUserSlice, 0, 100)
	return u
}
func (s GameUserSlice) Len() int           { return len(s) }
func (s GameUserSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s GameUserSlice) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s GameUserSlice) String() string {
	var str string
	for i := range s {
		str = str + fmt.Sprintf("{UserID:%d Score:%d}", s[i].UserID, s[i].Score)
	}
	return str
}

//
type GameUser struct {
	UserID uint32       `json:"userID"`
	input  *ClientInput ``
	Status int          `json:"status"`
	Score  int          `json:"score"`
	X      int          `json:"x"`
	Y      int          `json:"y"`
	Speed  int          `json:"speed"`
	Size   int          `json:"size"`
}

//
func NewGameUser(uid uint32) *GameUser {
	user := new(GameUser)
	user.input = &ClientInput{}
	user.UserID = uid
	user.ResetState()
	return user
}

func (u *GameUser) UpdateInput(input *ClientInput) {
	u.input.Down = input.Down
	u.input.Left = input.Left
	u.input.Right = input.Right
	u.input.Speed = input.Speed
	u.input.Up = input.Up
}

func (u *GameUser) ResetState() {
	u.X, u.Y = GetRandomPosition()
	u.Score = DEFAULT_SCORE
	u.Size = USER_SIZE
	u.Speed = DEFAULT_SPEED
	u.Status = STATE_USER_PREPARED
}

func (u *GameUser) IsMove() bool {
	userSpeed := 0
	moveOk := false
	if u.Score >= SPEED_DISSIPATION_SCORE && u.input.Speed == 1 {
		u.Score -= SPEED_DISSIPATION_SCORE
		userSpeed = u.Speed + SPEED_UP
	} else {
		userSpeed = u.Speed
	}
	if u.input.Left == 1 {
		u.X -= userSpeed
		moveOk = true
	}
	if u.input.Right == 1 {
		u.X += userSpeed
		moveOk = true
	}
	if u.input.Up == 1 {
		u.Y += userSpeed
		moveOk = true
	}
	if u.input.Down == 1 {
		u.Y -= userSpeed
		moveOk = true
	}
	return moveOk
}
