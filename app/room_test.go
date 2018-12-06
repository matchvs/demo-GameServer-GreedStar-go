package app

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/matchvs/gameServer-go/src/log"
)

func Test_CreateRoom(t *testing.T) {
	greed := NewGreedyStar()
	greed.CreateRoom(123456, 1234567897894, []byte("sdfsd"))
	greed.CreateRoom(123456, 1234567897895, []byte("sdfsd"))
	greed.CreateRoom(123456, 1234567897896, []byte("sdfsd"))

	greed.JoinRoom(123456, 1234567897894, []byte("sdfsd"))
	greed.JoinRoom(45644, 1234567897894, []byte("sdfsd"))
	greed.JoinRoom(64654, 1234567897894, []byte("sdfsd"))
	greed.JoinRoom(646589, 1234567897894, []byte("sdfsd"))

	greed.JoinRoom(86444, 1234567897895, []byte("sdfsd"))
	greed.JoinRoom(34676, 1234567897895, []byte("sdfsd"))
	greed.JoinRoom(36541, 1234567897896, []byte("sdfsd"))
	greed.JoinRoom(986454, 1234567897896, []byte("sdfsd"))

	time.Sleep(time.Second * 3)
	greed.LeaveRoom(64654, 1234567897894)
	greed.LeaveRoom(646589, 1234567897894)
	greed.LeaveRoom(123456, 1234567897894)
	greed.LeaveRoom(45644, 1234567897894)
	greed.LeaveRoom(86444, 1234567897895)
	greed.LeaveRoom(34676, 1234567897895)
	greed.LeaveRoom(36541, 1234567897896)
	greed.LeaveRoom(986454, 1234567897896)
	time.Sleep(time.Second * 3)
	greed.DeleteRoom(1234567897894)
	time.Sleep(time.Second * 3)
	greed.DeleteRoom(1234567897895)
	time.Sleep(time.Second * 3)
	greed.DeleteRoom(1234567897896)
}

func Test_RoomUser(t *testing.T) {
	room := NewRoomItem(20000, 123456789)
	room.AddUser(54321, []byte("hello"))
	room.AddUser(12345, []byte("hello"))
	room.AddUser(6546514, []byte("hello"))
	room.AddUser(6546514, []byte("hello"))
	room.AddUser(6546514, []byte("hello"))
	room.AddUser(6546514, []byte("hello"))
	room.AddUser(8761645, []byte("hello"))
	room.DelUser(54321)
}

func Test_RoomEventSend(t *testing.T) {
	userList := make([]*GameUser, 10)
	for i := 0; i < 10; i++ {
		user := NewGameUser(uint32(i + 1))
		userList[i] = user
	}
	List := &RoomEventSend{
		Type: "list",
		Data: userList,
	}

	if msg, err := json.Marshal(List); err != nil {
		log.LogW("json list error %v", err)
	} else {
		log.LogD("%v", string(msg))
	}

	str := &RoomEventSend{
		Type: "list",
		Data: "sdf",
	}
	if msg, err := json.Marshal(str); err != nil {
		log.LogW("json string error %v", err)
	} else {
		log.LogD("%v", string(msg))
	}

	str.Data = 124
	if msg, err := json.Marshal(str); err != nil {
		log.LogW("json number error %v", err)
	} else {
		log.LogD("%v", string(msg))
	}
}
