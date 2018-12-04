/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:59:24
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-04 19:11:29
 * @Description: file content
 */
package app

import (
	"encoding/json"
	"sync"

	"github.com/matchvs/gameServer-go/src/log"

	matchvs "github.com/matchvs/gameServer-go"
)

type GreedyStar struct {
	push    *matchvs.PushManager
	roomMap sync.Map
	gameID  uint32
}

func NewGreedyStar() *GreedyStar {
	gst := new(GreedyStar)
	return gst
}

func (self *GreedyStar) SetPush(p *matchvs.PushManager) {
	self.push = p
}

func (self *GreedyStar) CreateRoom(gameID uint32, roomID uint64, userProperty []byte) {
	self.gameID = gameID
	room := NewRoomItem(gameID, roomID)
	room.SetPush(self.push)
	self.roomMap.Store(roomID, room)
}

// 有人加入房间
func (self *GreedyStar) JoinRoom(userID uint32, roomID uint64, userProfile []byte) {
	item, ok := self.roomMap.Load(roomID)
	if ok {
		room := item.(*RoomItem)
		room.AddUser(userID, userProfile)
	} else {
		log.LogD("no this room [%d]", roomID)
	}
}

// 有人离开房间
func (self *GreedyStar) LeaveRoom(userID uint32, roomID uint64) {
	itme, ok := self.roomMap.Load(roomID)
	if ok {
		room := itme.(*RoomItem)
		room.DelUser(userID)
	} else {
		log.LogD("no this room [%v]", roomID)
	}
}

// 房间销毁
func (self *GreedyStar) DeleteRoom(roomID uint64) {
	self.roomMap.Delete(roomID)
}

// 处理来自客户端的消息
func (self *GreedyStar) ClientEvent(userID uint32, roomID uint64, datas []byte) {
	event := &RoomEventRecv{}
	json.Unmarshal(datas, event)
	item, ok := self.roomMap.Load(roomID)
	if !ok {
		return
	}
	room := item.(*RoomItem)
	switch event.Type {
	case "input":
		room.UpateUserInput(userID, event.Data.([]byte))
	case "startGame":
		room.StartGame(self.gameID, userID)
	}
}
