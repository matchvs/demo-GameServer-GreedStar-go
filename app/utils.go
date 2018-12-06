/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:58:23
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-06 19:14:07
 * @Description: file content
 */
package app

import (
	"math"
	"math/rand"
	"time"
)

func GetRandomPosition() (x, y int) {
	pad := 40
	minX := pad
	minY := pad
	maxX := GAME_MAP_WIDTH - pad
	maxY := GAME_MAP_HGITH - pad
	x = GetRandNums(minX, maxX, 3)[1]
	y = GetRandNums(minY, maxY, 4)[3]
	return
}

func GetRandNum(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := r.Intn(max-min) + min
	return randNum
}

//生成count个[min,max)结束的不重复的随机数
func GetRandNums(min int, max int, count int) []int {
	//范围检查
	if max < min {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	for len(nums) < count {
		//生成随机数
		num := GetRandNum(min, max)

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

func IsCollisionWithCircle(x1, y1, r1, x2, y2, r2 int) bool {
	if math.Sqrt(math.Pow(float64(x1-x2), 2)+math.Pow(float64((y1-y2)), 2)) <= float64(r1+r2) {
		return true
	}
	return false
}
