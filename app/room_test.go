package app

import (
	"fmt"
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

func Test_RandomPosition(t *testing.T) {
	for i := 0; i < 10; i++ {
		x, y := GetRandomPosition()
		fmt.Printf("x = %d, y = %d \n", x, y)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("list %v\n", GetRandNums(123, 4561441, 10))
	}
}
