/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:55:57
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-04 19:29:25
 * @Description: file content
 */

package app

type Food struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Score int `json:"score"`
	ID    int `json:"id"`
	Size  int `json:"size"`
}

func NewFood() *Food {
	food := new(Food)
	return food
}
