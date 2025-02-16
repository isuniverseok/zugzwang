package notation

import "strconv"

func FileToInt(file byte) int {
	return int(file - 'a')
}

func RankToInt(rank byte) int {
	return 7 - int(rank-'1')
}

func IntToFile(file int) byte {
	return 'a' + byte(file)
}

func IntToRank(rank int) byte {
	return strconv.Itoa(7 - rank + 1)[0]
}