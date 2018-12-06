/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:20
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-06 18:17:15
 * @Description: file content
 */

package app

type RoomEventRecv struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type RoomEventSend struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Profile interface{} `json:"profile"`
}
