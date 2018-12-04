/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:55:45
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-04 19:43:05
 * @Description: file content
 */
package app

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/matchvs/gameServer-go/src/defines"

	matchvs "github.com/matchvs/gameServer-go"
)

type RoomItem struct {
	roomID   uint64
	userList GameUserSlice
	gameTime int32
	push     *matchvs.PushManager
	gameID   uint32
	foodList []*Food
}

func NewRoomItem(gameID uint32, roomID uint64) *RoomItem {
	room := new(RoomItem)
	room.roomID = roomID
	room.userList = NewGameUserSlice()
	room.foodList = make([]*Food, 0, 100)
	room.gameTime = 54000
	return room
}

func (self *RoomItem) SetPush(p *matchvs.PushManager) {
	self.push = p
}

// 排序
func (self *RoomItem) SortUser() {
	sort.Sort(self.userList)
}

// 获取房间用户列ID 列表
func (self *RoomItem) GetUserIDList() []uint32 {
	uids := make([]uint32, len(self.userList))
	for key := range self.userList {
		uids = append(uids, self.userList[key].userID)
	}
	return uids
}

// 查询用户ID
func (self *RoomItem) FindUser(userID uint32) (*GameUser, int) {
	// var (
	// 	hight = len(self.userList)
	// 	low   = 0
	// 	mid   = 0
	// )
	// for i := 0; i < len(self.userList); i++ {
	// 	mid = (low + hight) / 2
	// 	tmp := self.userList[mid].userID
	// 	if tmp == userID {
	// 		return self.userList[mid], mid
	// 	} else if tmp > userID {
	// 		low = mid + 1
	// 	} else {
	// 		hight = mid - 1
	// 	}
	// 	if low > hight {
	// 		return nil, -1
	// 	}
	// }
	// return nil, -1
	for key := range self.userList {
		tmp := self.userList[key].userID
		if tmp == userID {
			return self.userList[key], key
		}
	}
	return nil, -1
}

// 去重
func (self *RoomItem) DelRepeat() {
	userMap := make(map[uint32]*GameUser)
	for key := range self.userList {
		userMap[self.userList[key].userID] = self.userList[key]
	}
	self.userList = self.userList[:0]
	for _, value := range userMap {
		self.userList = append(self.userList, value)
	}
	self.SortUser()
}

// 添加用户
func (self *RoomItem) AddUser(userID uint32, userProfile []byte) {
	// 先判断是否重复
	if _, index := self.FindUser(userID); index >= 0 {
		return
	}
	user := NewGameUser(userID)
	self.userList = append(self.userList, user)
	// 添加新值就排序
	self.SortUser()
}

// 删除用户
func (self *RoomItem) DelUser(userID uint32) {
	_, index := self.FindUser(userID)
	self.userList = append(self.userList[0:index], self.userList[index+1:]...)
}

// 更新用户输入操作
func (self *RoomItem) UpateUserInput(userid uint32, data []byte) error {
	input := &ClientInput{}
	if err := json.Unmarshal(data, input); err != nil {
		return err
	}
	user, index := self.FindUser(userid)
	if index < 0 {
		return errors.New("no this user")
	}
	user.UpdateInput(input)
	return nil
}

// 接收客户端发送的消息，同步开始游戏
func (self *RoomItem) StartGame(gameID, userID uint32) {
	// 做消息中转
	room := &RoomEventSend{
		Type:    "startGame",
		Data:    self.userList,
		Profile: self.gameTime,
	}
	msg, _ := json.Marshal(room)
	self.PushEventOther([]uint32{userID}, msg)
	self.SendFoodMsg(userID)
	return
}

// 发送食物
func (self *RoomItem) SendFoodMsg(userID uint32) {
	times := 3
	oneslen := len(self.foodList) / times
	for i := 0; i < 3; i++ {
		food := self.foodList[i*oneslen : (i+1)*oneslen]
		msg, _ := json.Marshal(food)
		self.PushEventOther([]uint32{userID}, msg)
	}
}

func (self *RoomItem) PushEventOther(uids []uint32, msg []byte) {
	req := &defines.MsPushEventReq{
		RoomID:    self.roomID,
		PushType:  1,
		GameID:    self.gameID,
		DestsList: uids[:],
		CpProto:   msg[:],
	}
	self.push.PushEvent(req)
}

func (self *RoomItem) PushEvent(msg []byte) {
	req := &defines.MsPushEventReq{
		RoomID:    self.roomID,
		PushType:  1,
		GameID:    self.gameID,
		DestsList: self.GetUserIDList(),
		CpProto:   msg[:],
	}
	self.push.PushEvent(req)
}
