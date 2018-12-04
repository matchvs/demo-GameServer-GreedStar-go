/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:34
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-04 17:22:10
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
func (s GameUserSlice) Less(i, j int) bool { return s[i].userID > s[j].userID }
func (s GameUserSlice) String() string {
	var str string
	for i := range s {
		str = str + fmt.Sprintf("{%v %v}", s[i].userID, &s[i].input)
	}
	return str
}

//
type GameUser struct {
	userID uint32
	input  *ClientInput
}

//
func NewGameUser(uid uint32) *GameUser {
	user := new(GameUser)
	user.input = &ClientInput{}
	user.userID = uid
	return user
}

func (u *GameUser) UpdateInput(input *ClientInput) {
	u.input.down = input.down
	u.input.left = input.left
	u.input.right = input.right
	u.input.speed = input.speed
	u.input.up = input.up
}
