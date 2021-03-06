/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-12-03 16:56:50
 * @LastEditors: Ville
 * @LastEditTime: 2018-12-11 18:20:51
 * @Description: file content
 */

package app

const (
	STATE_USER_PREPARED = 1
	STATE_USER_PLAYING  = 2
	STATE_USER_DIEAD    = 3
)

var (
	USER_SIZE               = 20
	GAME_MAP_WIDTH          = 2560
	GAME_MAP_HGITH          = 1440
	FOOD_SHOW               = 1
	USER_MIN_SPEED          = 1
	FOOD_HIDE               = 2
	SIZE_MULTIPLE           = 200
	SPEED_MULTIPLE          = 3000
	DEFAULT_SPEED           = 5
	DEFAULT_SCORE           = 0
	SPEED_UP                = 2
	SPEED_DISSIPATION_SCORE = 3
	GAME_TIME               = 120 // 游戏时间
	GAME_FPS_INTERVAL       = 20  // 每秒执行多少次检测
	GAME_TIMER_INTERVAL     = 1000 / GAME_FPS_INTERVAL
	ScoreList               = []int{20, 40, 60}
	FOOD_INITIAL_SIZE       = []int{2, 3, 8}
	FOOD_INITIAL_NUM        = 60
)

const (
	OPTION_ADDFOOD = "addFood"
	OPTION_MOVE    = "move"
)
