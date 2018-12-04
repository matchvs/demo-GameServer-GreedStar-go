package app

import (
	"testing"
)

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
