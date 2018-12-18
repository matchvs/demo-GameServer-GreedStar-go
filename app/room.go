/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:55:45
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-18 19:04:54
 * @Description: file content
 */
package app

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	matchvs "github.com/matchvs/gameServer-go"
	"github.com/matchvs/gameServer-go/src/defines"
	"github.com/matchvs/gameServer-go/src/log"
)

type RoomItem struct {
	roomID    uint64
	userList  GameUserSlice
	gameTime  int32
	push      matchvs.PushHandler
	gameID    uint32
	foodList  []*Food
	foodNum   int
	roomClose chan int
}

// 新建一个房间
func NewRoomItem(gameID uint32, roomID uint64) *RoomItem {
	room := new(RoomItem)
	room.roomID = roomID
	room.userList = NewGameUserSlice()
	room.gameTime = int32(GAME_TIME * GAME_FPS_INTERVAL)
	room.roomClose = make(chan int)
	room.foodNum = 0
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
		self.foodNum++
	}
	log.LogD("init food number %d", self.foodNum)
}

func (self *RoomItem) SetPush(p matchvs.PushHandler) {
	self.push = p
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
}

// 添加用户
func (self *RoomItem) AddUser(userID uint32, userProfile []byte) *GameUser {
	// 先判断是否重复
	if user, index := self.FindUser(userID); index >= 0 {
		return user
	}
	user := NewGameUser(userID)
	self.userList = append(self.userList, user)
	self.RoomUserRank()
	// log.LogD("当前用户：%v", self.userList)
	return user
}

// 删除用户
func (self *RoomItem) DelUser(userID uint32) {
	_, index := self.FindUser(userID)
	self.userList = append(self.userList[0:index], self.userList[index+1:]...)
	log.LogD("[%d] leave room，user number [%d]", userID, len(self.userList))
}

// 更新用户输入操作
func (self *RoomItem) UpateUserInput(userid uint32, data map[string]interface{}) error {
	input := &ClientInput{}
	input.Left = int(data["l"].(float64))
	input.Right = int(data["r"].(float64))
	input.Up = int(data["u"].(float64))
	input.Down = int(data["d"].(float64))
	input.Speed = int(data["p"].(float64))
	user, index := self.FindUser(userid)
	if index < 0 {
		return errors.New("no this user")
	}
	user.UpdateInput(input)
	return nil
}

// 接收客户端发送的消息，同步开始游戏
func (self *RoomItem) StartGame(gameID, userID uint32) {
	self.gameTime = int32(GAME_TIME * GAME_TIMER_INTERVAL)
	// 做消息中转
	room := &RoomEventSend{
		Type:    "startGame",
		Data:    self.userList,
		Profile: self.gameTime,
	}
	log.LogD("开始游戏：[%d]", userID)
	self.PushEventOther([]uint32{userID}, room)
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
			Type: OPTION_ADDFOOD,
			Data: food,
		}
		self.PushEventOther([]uint32{userID}, event)
	}
}

// 发送其他用户给当 给指定用户, 返回除自己以外的其他玩家列表
func (self *RoomItem) SendOtherUsers(userID uint32) {
	users := make([]*GameUser, 0, len(self.userList))
	uis := make([]uint32, 0, len(self.userList))
	for key := range self.userList {
		if self.userList[key].UserID != userID {
			users = append(users, self.userList[key])
			uis = append(uis, self.userList[key].UserID)
		}
	}
	if len(users) > 0 {
		event := &RoomEventSend{
			Type: "otherPlayer",
			Data: users,
		}
		log.LogD("Tell [%d] OtherUsers:[%v]：", userID, uis)
		self.PushEventOther([]uint32{userID}, event)
	}
}

// 推送消息给指定部分用户
func (self *RoomItem) PushEventOther(uids []uint32, event *RoomEventSend) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	req := &defines.MsPushEventReq{
		RoomID:    self.roomID,
		PushType:  1,
		GameID:    self.gameID,
		DestsList: uids[:],
		CpProto:   msg[:],
	}
	return self.push.PushEvent(req)
}

// 发送消息给所有用户
func (self *RoomItem) PushEvent(event *RoomEventSend) error {
	msg, _ := json.Marshal(event)
	req := &defines.MsPushEventReq{
		RoomID:    self.roomID,
		PushType:  3,
		GameID:    self.gameID,
		DestsList: self.GetUserIDList(),
		CpProto:   msg[:],
	}
	return self.push.PushEvent(req)
}

// 有人加入房间
func (self *RoomItem) UserJoinRoom(userID uint32, userProfile []byte) {
	self.AddUser(userID, userProfile)
	user := self.AddUser(userID, userProfile)
	event := &RoomEventSend{
		Type: "addPlayer",
		Data: user,
	}

	self.PushEvent(event)
	// 告诉当前玩家 游戏时间
	event2 := &RoomEventSend{
		Type: "countDown",
		Data: self.gameTime,
	}
	self.PushEventOther([]uint32{userID}, event2)
	// 给当前玩家发送食物列表
	self.SendFoodMsg(userID)
	// 告诉当前玩家，现在房间有哪些人
	self.SendOtherUsers(userID)

}

// 房间游戏定时器
func (self *RoomItem) StartTimer() {
	timerInterval := time.Duration(GAME_TIMER_INTERVAL)
	// log.LogD("timerInterval:%v", timerInterval)
	for {
		select {
		case _, ok := <-self.roomClose:
			if !ok {
				log.LogI("room[%d] is deleted", self.roomID)
				return
			}
		case <-time.After(time.Millisecond * timerInterval):
			self.gameTime--
			if self.gameTime <= 0 {
				event := &RoomEventSend{
					Type: "GameOver",
					Data: "",
				}
				// 发送给所有人
				self.PushEvent(event)
				return
			}
			self.IsUserscollision()
			self.IsFoodListFull()
			self.RoomUserRank()
			if self.IsPlayerMove() {
				event := &RoomEventSend{
					Type: OPTION_MOVE,
					Data: self.userList,
				}
				if err := self.PushEvent(event); err != nil {
					log.LogE("push player move error : %v", err)
					return
				}
			}
		}
	}
}
func (self *RoomItem) StopTimer() {
	close(self.roomClose)
}

// 碰撞检测
func (self *RoomItem) IsUserscollision() {
	userNum := len(self.userList)
	for i := 0; i < userNum; i++ {
		p1 := self.userList[i]
		for j := i + 1; j < userNum; j++ {
			p2 := self.userList[j]
			if IsCollisionWithCircle(p1.X, p1.Y, p1.Size, p2.X, p2.Y, p2.Size) {
				if p1.Score == p2.Score {
					break
				} else if p1.Score > p2.Score {
					p1.Score += p1.Score + p2.Score
					p2.ResetState()
				} else {
					p2.Score += p2.Score + p2.Score
					p1.ResetState()
				}
			}
		}
		// 检测食物碰撞
		self.IsFoodCollision(p1)
		// 检测边界碰撞
		self.IsBorderCollision(p1)
	}
}

// 判断是否吃食物
func (self *RoomItem) IsFoodCollision(user *GameUser) {
	foodNum := len(self.foodList)
	list := self.foodList[:]
	for i := 0; i < foodNum; i++ {
		food := self.foodList[i]
		if IsCollisionWithCircle(food.X, food.Y, food.Size, user.X, user.Y, user.Size) {
			user.Score += food.Score
			user.Size = USER_SIZE + user.Score/SIZE_MULTIPLE
			user.Speed = DEFAULT_SPEED - user.Score/SPEED_MULTIPLE
			if user.Speed < USER_MIN_SPEED {
				user.Speed = USER_MIN_SPEED
			}
			list = append(list[:i], list[i+1:]...)
			// log.LogD("吃掉了食物：[%v]", food.ID)
			event := &RoomEventSend{
				Type: "removeFood",
				Data: food.ID,
			}
			self.PushEvent(event)
		}
	}
	self.foodList = list
}

// 边界碰撞判断
func (self *RoomItem) IsBorderCollision(user *GameUser) {
	lAcme := user.X - user.Size
	rAcme := user.X + user.Size
	uAcme := user.Y + user.Size
	dAcme := user.Y - user.Size
	if lAcme <= 0 || rAcme >= GAME_MAP_WIDTH || uAcme >= GAME_MAP_HGITH || dAcme <= 0 {
		// log.LogD("出界死了: ")
		// log.LogD("当前左边位置：[x=%d, size=%d, lAcme=%v]", user.X, user.Size, lAcme)
		// log.LogD("当前右边位置：[x=%d, size=%d, rAcme=%v]", user.X, user.Size, rAcme)
		// log.LogD("当前上边位置：[x=%d, size=%d, lAcme=%v]", user.Y, user.Size, uAcme)
		// log.LogD("当前下边位置：[x=%d, size=%d, lAcme=%v]", user.Y, user.Size, dAcme)
		user.ResetState()
	}
}

// 食物是否满
func (self *RoomItem) IsFoodListFull() {
	list := make([]*Food, 0, FOOD_INITIAL_NUM-len(self.foodList))
	for i := 0; i < FOOD_INITIAL_NUM; i++ {
		number := len(self.foodList)
		if number < FOOD_INITIAL_NUM {
			// log.LogD("当前食物数量：[%d]，下一个食物ID[%d]", number, self.foodNum)
			food := NewFood(self.foodNum)
			self.foodList = append(self.foodList, food)
			self.foodNum++
			list = append(list, food)
		} else {
			if len(list) > 0 {
				// log.LogD("新增加了食物：[%v] 个", len(list))
				event := &RoomEventSend{
					Type: OPTION_ADDFOOD,
					Data: list,
				}
				self.PushEvent(event)
			}
			return
		}
	}
}

// 玩家按照分数排行
func (self *RoomItem) RoomUserRank() {
	sort.Sort(self.userList)
}

// 玩家移动
func (self *RoomItem) IsPlayerMove() bool {
	ismove := false
	for i := 0; i < len(self.userList); i++ {
		if self.userList[i].IsMove() {
			// log.LogD("用户[%d]在移动", self.userList[i].UserID)
			ismove = true
		}
	}
	return ismove
}
