/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:55:57
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-05 11:32:27
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

func NewFood(id int) *Food {
	food := new(Food)
	food.X, food.Y = GetRandomPosition()
	food.ID = id
	food.Score = GenerateScore()
	food.Size = food.Score
	return food
}

func GenerateScore() int {
	return ScoreList[GetRandNum(0, len(ScoreList)-1)]
}
