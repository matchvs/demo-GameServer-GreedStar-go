/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:55:45
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-05 19:08:08
 * @Description: file content
 */
package app

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

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
	foodNum  int
}

// 新建一个房间
func NewRoomItem(gameID uint32, roomID uint64) *RoomItem {
	room := new(RoomItem)
	room.roomID = roomID
	room.userList = NewGameUserSlice()
	room.gameTime = 54000
	// 初始化食物
	room.InitFoods()
	return room
}

// 初始化食物
func (self *RoomItem) InitFoods() {
	self.foodList = make([]*Food, 0, FOOD_INITIAL_NUM)
	for i := 0; i < FOOD_INITIAL_NUM; i++ {
		food := NewFood(i)
		self.foodList = append(self.foodList, food)
	}
	self.foodNum = FOOD_INITIAL_NUM
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
		uids = append(uids, self.userList[key].UserID)
	}
	return uids
}

// 查询用户ID
func (self *RoomItem) FindUser(userID uint32) (*GameUser, int) {
	for key := range self.userList {
		tmp := self.userList[key].UserID
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
		userMap[self.userList[key].UserID] = self.userList[key]
	}
	self.userList = self.userList[:0]
	for _, value := range userMap {
		self.userList = append(self.userList, value)
	}
	self.SortUser()
}

// 添加用户
func (self *RoomItem) AddUser(userID uint32, userProfile []byte) *GameUser {
	// 先判断是否重复
	if user, index := self.FindUser(userID); index >= 0 {
		return user
	}
	user := NewGameUser(userID, STATE_USER_PREPARED, DEFAULT_SCORE, USER_SIZE, DEFAULT_SPEED)
	self.userList = append(self.userList, user)
	// 添加新值就排序
	self.SortUser()
	return user
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

// 发送食物给指定用户
func (self *RoomItem) SendFoodMsg(userID uint32) {
	times := 3
	oneslen := len(self.foodList) / times
	for i := 0; i < 3; i++ {
		food := self.foodList[i*oneslen : (i+1)*oneslen]
		event := &RoomEventSend{
			Type: "addFood",
			Data: food,
		}
		msg, _ := json.Marshal(event)
		self.PushEventOther([]uint32{userID}, msg)
	}
}

// 发送其他用户给当 给指定用户
func (self *RoomItem) SendOtherUsers(userID uint32) {
	users := make([]uint32, len(self.userList))
	for key := range self.userList {
		if self.userList[key].UserID != userID {
			users = append(users, self.userList[key].UserID)
		}
	}
	event := &RoomEventSend{
		Type: "otherPlayer",
		Data: users,
	}
	msg, _ := json.Marshal(event)
	self.PushEventOther([]uint32{userID}, msg)
}

// 推送消息给指定部分用户
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

// 发送消息给所有用户
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

// 有人加入房间
func (self *RoomItem) UserJoinRoom(userID uint32, userProfile []byte) {
	user := self.AddUser(userID, userProfile)
	event := &RoomEventSend{
		Type: "addPlayer",
		Data: user,
	}
	msg, _ := json.Marshal(event)
	self.PushEventOther([]uint32{userID}, msg)

	event2 := &RoomEventSend{
		Type: "countDown",
		Data: self.gameTime,
	}
	msg, _ = json.Marshal(event2)
	self.PushEventOther([]uint32{userID}, msg)

	self.SendFoodMsg(userID)
	self.SendOtherUsers(userID)

}

// 房间游戏定时器
func (self *RoomItem) RoomTimer() {
	for {
		if self.gameTime <= 0 {
			event := &RoomEventRecv{
				Type: "GameOver",
				Data: "",
			}
			msg, _ := json.Marshal(event)
			// 发送给所有人
			self.PushEvent(msg)
			return
		}
		self.IsUserscollision()
		self.IsFoodCollision()
		self.IsBorderCollision()
		self.IsFoodListFull()
		self.RoomUserRank()
		if self.IsPlayerMove() {

		}
		time.Sleep(time.Microsecond * GAME_TIMER_INTERVAL)
	}
}

// 碰撞检测
func (self *RoomItem) IsUserscollision() {
	userNum := len(self.userList)
	for i := 0; i < userNum; i++ {
		for j := 0; j < userNum; j++ {
		}
	}
}

// 判断是否吃食物
func (self *RoomItem) IsFoodCollision() {

}

// 边界碰撞判断
func (self *RoomItem) IsBorderCollision() {

}

// 食物是否满
func (self *RoomItem) IsFoodListFull() {

}

// 玩家按照分数排行
func (self *RoomItem) RoomUserRank() {

}

// 玩家移动
func (self *RoomItem) IsPlayerMove() bool {
	return false
}
