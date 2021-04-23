package bsnChainCode

import (
	"awesomeProject/utils"
)

var owner string

/*
respondOracles代表所有有回复的预言机
dataHashes代表每个respondOracles的回复
*/

func ReportScore(respondOracles []string, dataHashes []string, nRequestOracles int) string {
	var nMatches int = nRequestOracles / 2
	var correctHash string
	var points []int
	var Aggrement map[string]int = make(map[string]int)                  //查看每个回复的重复次数
	var DuplicateOracles map[string][]string = make(map[string][]string) //这个回复的所有预言机
	var TrueOracles []string = make([]string, 0)                         //所有这个回复的预言机

	for i := 0; i < len(dataHashes); i++ {
		if Aggrement[dataHashes[i]] != 0 {
			Aggrement[dataHashes[i]]++
			DuplicateOracles[dataHashes[i]] = append(DuplicateOracles[dataHashes[i]], respondOracles[i])
		} else {
			Aggrement[dataHashes[i]] = 1
			DuplicateOracles[dataHashes[i]] = []string{respondOracles[i]}
		}
	}
	/*
		找出正确的结果
		返回该结果的所有预言机
	*/
	for i := 0; i < len(dataHashes); i++ {
		if Aggrement[dataHashes[i]] >= nMatches {
			correctHash = dataHashes[i]
			for j := 0; j < len(DuplicateOracles[correctHash]); j++ {
				TrueOracles = append(TrueOracles, DuplicateOracles[correctHash][j])
			}
			break
		}
	}
	/*
		计算这次所有有返回结果的预言机的信誉值平均分
	*/
	for i := 0; i < len(dataHashes); i++ {
		for j := 0; j < len(TrueOracles); j++ {
			if respondOracles[i] == TrueOracles[j] {
				points = append(points, 100)
				tempOracle := utils.Oracles[respondOracles[i]]
				c := tempOracle.Times
				tempOracle.Times++
				tempOracle.AveScore = (tempOracle.AveScore*c + 100) / (c + 1)
				break
			} else {
				points = append(points, 0)
				tempOracle := utils.Oracles[respondOracles[i]]
				c := tempOracle.Times
				tempOracle.Times++
				tempOracle.AveScore = (tempOracle.AveScore * c) / (c + 1)
			}
		}
	}
	index := selectWinningOracles(TrueOracles)
	return utils.Oracles[TrueOracles[index]].Address
}

func selectWinningOracles(trueOraclesAddr []string) int {
	maxPoint := utils.Oracles[trueOraclesAddr[0]].AveScore
	maxIndex := 0
	for i := 1; i < len(trueOraclesAddr); i++ {
		if utils.Oracles[trueOraclesAddr[i]].AveScore > maxPoint {
			maxPoint = utils.Oracles[trueOraclesAddr[i]].AveScore
			maxIndex = i

		}
	}
	return maxIndex
}
