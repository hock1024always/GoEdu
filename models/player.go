package models

import (
	"github.com/hock1024always/GoEdu/dao"
	"gorm.io/gorm"
	"sort"
	"strconv"
)

type Player struct {
	Id          int    `json:"id"`
	Nickname    string `json:"nickname"`
	Aid         int    `json:"aid"`         // 活动序号
	Ref         int    `json:"ref"`         // 注销号码
	Avatar      int    `json:"avatar"`      // 头像序号
	Score       int    `json:"score"`       // 积分
	Declaration string `json:"declaration"` // 将类型改为 string
	Password    string `json:"password"`
	// UpdateTime  int    `json:"updateTime"`
}

// 用来通过便利法得到排名
type PlayerScore struct {
	PlayerID int
	Score    int
}

func (Player) TableName() string {
	return "player"
}

func AddPlayer(nickname, password string) (Player, error) {
	player := Player{Nickname: nickname, Password: password}
	err := dao.Db.Create(&player).Error
	return player, err
}

// 获取某种顺序排列的某一活动的玩家列表 DESC降序 ASC升序
func GetPlayers(aid int, sort string) ([]Player, error) {
	var players []Player
	err := dao.Db.Where("aid =?", aid).Order(sort).Find(&players).Error
	return players, err
}

func GetPlayerById(id int) (Player, error) {
	var player Player
	err := dao.Db.Where("id =?", id).First(&player).Error
	return player, err
}

// 通过投票来更新得分 加分
func UpdateScoreByVote(id int) error {
	var player Player
	err := dao.Db.Model(&player).Where("id =?", id).UpdateColumn("score", gorm.Expr("score + ?", 1)).Error
	return err
}

// 通过活动来更新得分 减分
func UpdateScoreByActivity(playerId int) error {
	var player Player
	err := dao.Db.Model(&player).Where("id =?", playerId).UpdateColumn("score", gorm.Expr("score - ?", 1)).Error
	return err
}

// 通过管理员修改来更新得分
func UpdateScoreByAdmin(nickname string, score int) error {
	var player Player
	err := dao.Db.Model(&player).Where("nickname =?", nickname).UpdateColumn("score", score).Error
	return err
}

func CheckPlayerExistsByNickname(nickname string) (Player, error) {
	var player Player
	err := dao.Db.Where("nickname =?", nickname).First(&player).Error
	return player, err
}

// 更改宣言
func AddDeclaration(id int, declaration string) (Player, error) {
	var player Player
	err := dao.Db.Model(&player).Where("id =?", id).Update("declaration", declaration).Error
	return player, err
}

// 删除涉及玩家和活动的投票记录 并返回得分
func DeletePlayerByActivityId(playerId, activityId int) (int, error) {
	// 定义一个变量来存储删除的记录数量
	result := dao.Db.Where("player_id = ? AND activity_id = ?", playerId, activityId).Delete(&Vote{})
	var voteNum int = int(result.RowsAffected)
	// 返回删除的条目数量和可能发生的错误
	return voteNum, result.Error
}

// 通过删除记录来扣分
func ReduceScore(playerId int, voteNum int) error {
	var player Player
	err := dao.Db.Model(&player).Where("id =?", playerId).UpdateColumn("score", gorm.Expr("score - ?", voteNum)).Error
	return err
}

func UpdatePlayerAid(id int, aid int) error {
	var player Player
	err := dao.Db.Model(&player).Where("id =?", id).Update("aid", aid).Error
	return err
}

// 查看给这个参赛者投票的投票记录
func GetVoteListForPlayer(playerId int, sort string) ([]Vote, error) {
	var votes []Vote
	err := dao.Db.Where("player_id =?", playerId).Order(sort).Find(&votes).Error
	return votes, err
}

// 查看给这个活动投票给某个用户的投票记录
func GetVoteListForPlayerByActivityId(playerId int, activityId int, sort string) ([]Vote, error) {
	var votes []Vote
	err := dao.Db.Where("player_id =? AND activity_id =?", playerId, activityId).Order(sort).Find(&votes).Error
	return votes, err
}

// 查找多有参赛者 按分数排序
func GetAllPlayers(sort string) ([]Player, error) {
	var players []Player
	err := dao.Db.Order(sort).Find(&players).Error
	return players, err
}

// 根据用户删除投票记录
func DeleteVoteByUserIdAndActivityId(userId, activityId int) ([]int, error) {
	var votes []Vote

	// 查询符合条件的投票记录
	err := dao.Db.Where("user_id = ? AND activity_id = ?", userId, activityId).Find(&votes).Error
	if err != nil {
		return nil, err // 如果没有找到记录，返回错误
	}

	// 如果没有找到投票记录，直接返回一个空切片
	if len(votes) == 0 {
		return []int{}, nil
	}

	// 删除符合条件的投票记录
	result := dao.Db.Where("user_id = ? AND activity_id = ?", userId, activityId).Delete(&Vote{})
	if result.Error != nil {
		return nil, result.Error // 返回删除时的错误
	}

	// 提取所有的 player_id，返回
	playerIds := make([]int, len(votes))
	for i, vote := range votes {
		playerIds[i] = vote.PlayerId
	}

	return playerIds, nil
}

func DeleteVoteScore(playerIds []int) error {
	for _, playerId := range playerIds {
		err := dao.Db.Model(&Player{}).Where("id =?", playerId).UpdateColumn("score", gorm.Expr("score - ?", 1)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// IntArrayToStringArray 将整型数组转换为字符串型数组
func IntArrayToStringArray(intArray []int) []string {
	// 创建一个字符串型数组，长度与整型数组相同
	stringArray := make([]string, len(intArray))

	// 遍历整型数组，并将每个元素转换为字符串
	for i, v := range intArray {
		stringArray[i] = strconv.Itoa(v) // 使用 strconv.Itoa 将整型转换为字符串
	}

	return stringArray
}

// 输入活动的编号 返回参赛者数组
func CheckVoteByActivityId(activityId int) ([]int, error) {
	var votes []Vote

	// 查询符合条件的投票记录
	err := dao.Db.Where("activity_id = ?", activityId).Find(&votes).Error
	if err != nil {
		return nil, err // 如果没有找到记录，返回错误
	}

	// 如果没有找到投票记录，直接返回一个空切片
	if len(votes) == 0 {
		return []int{}, nil
	}

	// 提取所有的 player_id，返回
	playerIds := make([]int, len(votes))
	for i, vote := range votes {
		playerIds[i] = vote.PlayerId
	}

	return playerIds, nil
}

// 根据参赛者编号计算得分
func GetScoresFromVotes(playerIds []int) []PlayerScore {
	scoreMap := make(map[int]int)

	// 统计每个参赛者的得分
	for _, playerId := range playerIds {
		scoreMap[playerId]++
	}

	// 将得分数据转换为 PlayerScore 数组
	playerScores := make([]PlayerScore, 0, len(scoreMap))
	for playerId, score := range scoreMap {
		playerScores = append(playerScores, PlayerScore{PlayerID: playerId, Score: score})
	}

	// 根据得分从高到低排序
	sort.Slice(playerScores, func(i, j int) bool {
		return playerScores[i].Score > playerScores[j].Score
	})

	return playerScores
}
