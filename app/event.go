/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:20
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-04 18:21:36
 * @Description: file content
 */

package app

type RoomEventRecv struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type RoomEventSend struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Profile interface{} `json:"profile"`
}
