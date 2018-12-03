/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:34
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-03 19:03:05
 * @Description: file content
 */
package app

type GameUser struct {
	userID uint32
	input  *ClientInput
}

func NewGameUser(uid uint32) *GameUser {
	user := new(GameUser)
	user.input = &ClientInput{}
	user.userID = uid
	return user
}
