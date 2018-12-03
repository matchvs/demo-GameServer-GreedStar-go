/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:55:45
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-03 18:54:36
 * @Description: file content
 */
package app

import (
	"sync"
)

type RoomItem struct {
	roomID   uint64
	userMap  map[uint32]*GameUser
	userlock sync.RWMutex
}

func NewRoomItem(roomID uint64) *RoomItem {
	room := new(RoomItem)
	room.roomID = roomID
	room.userMap = make(map[uint32]*GameUser)
	return room
}

func (self *RoomItem) AddUser(userID uint32, userProfile []byte) {
	user := NewGameUser(userID)
	self.userlock.Lock()
	defer self.userlock.Unlock()
	self.userMap[userID] = user
}

func (self *RoomItem) DelUser(userID uint32) {
	self.userlock.Lock()
	defer self.userlock.Unlock()
	delete(self.userMap, userID)
}
