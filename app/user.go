/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:34
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-05 16:56:49
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
func (s GameUserSlice) Less(i, j int) bool { return s[i].UserID > s[j].UserID }
func (s GameUserSlice) String() string {
	var str string
	for i := range s {
		str = str + fmt.Sprintf("{%v %v}", s[i].UserID, &s[i].input)
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
func NewGameUser(uid uint32, status, score, size, speed int) *GameUser {
	user := new(GameUser)
	user.input = &ClientInput{}
	user.UserID = uid
	user.X, user.Y = GetRandomPosition()
	user.Size = size
	user.Status = status
	user.Speed = speed
	user.Score = score
	return user
}

func (u *GameUser) UpdateInput(input *ClientInput) {
	u.input.down = input.down
	u.input.left = input.left
	u.input.right = input.right
	u.input.speed = input.speed
	u.input.up = input.up
}
