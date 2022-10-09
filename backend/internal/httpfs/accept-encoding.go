package httpfs

import (
	"strconv"
)

const (
	// defaultQValue default quality value
	defaultQValue = 1.0
)

func findStrInSlice(s string, ss []string) (findIdx int) {
	ssLen := len(ss)
	for ; findIdx < ssLen; findIdx++ {
		if s == ss[findIdx] {
			return
		}
	}

	findIdx = -1
	return
}

// ParseAcceptHeaderWithPoo 解析以 "," 和 ";" 為分隔的 Accept-Encoding 、Accept-Language Header 格式
func ParseAcceptHeader(s string, serverAcceptList []string) (finalV string) {
	strLen := len(s)

	var finalQuality float32
	var finalServerPriority int

	var head, tail int
	for i := 0; i < strLen; i++ {
		dealSegments := false
		if s[i] == ',' {
			if s[head] == ',' {
				// 當 s[0] == ',' 或者連續 ',' 時(ex:",,,...,,,")
				head++
			} else {
				tail = i
				dealSegments = true
			}
		} else if i == strLen-1 {
			// 最尾段可能不會有 ',' 所以需要特別判斷
			tail = strLen
			dealSegments = true
		}

		if dealSegments {
			v, q := ParseAcceptValue(s, head, tail)
			if v != "" && q > 0 {
				if v == "*" {
					if q > finalQuality || (q == finalQuality && finalV != serverAcceptList[0]) {
						finalV = serverAcceptList[0]
						finalQuality = q
						finalServerPriority = 0
					}
				} else {
					fIdx := findStrInSlice(v, serverAcceptList)
					if fIdx != -1 {
						if q > finalQuality || (q == finalQuality && fIdx < finalServerPriority) {
							finalV = v
							finalQuality = q
							finalServerPriority = fIdx
						}
					}
				}
			}
			head = i + 1
		}
	}

	return finalV
}

// ParseAcceptValue 將 header 分割為 value 和 quality
//
// ex: en-US;q=0.6 => value:"en-US" quality=0.6
func ParseAcceptValue(s string, head, tail int) (value string, quality float32) {
	if head > tail {
		return
	}

	curIdx := head

	// find semicolon and trim space
	const maxSwitch = 2
	// sRecord => space switch record (非空格和空格的交換位置)
	var sRecord [maxSwitch]int
	// sNum = > switch amount
	var sAmount int
	isSpace := true
	for ; curIdx < tail; curIdx++ {
		if s[curIdx] == ';' {
			break
		}
		if s[curIdx] == ' ' && !isSpace || s[curIdx] != ' ' && isSpace {
			if sAmount == maxSwitch {
				return "", 0
			}
			isSpace = !isSpace
			sRecord[sAmount] = curIdx
			sAmount++
		}
	}
	if !isSpace {
		if curIdx >= tail {
			sRecord[1] = tail
		} else {
			sRecord[1] = curIdx
		}
	}

	// find "q=" string
	sIdx := curIdx
	tailMin2 := tail - 2
	findNotSpace := false
	for curIdx++; curIdx < tailMin2; curIdx++ {
		if s[curIdx] != ' ' {
			findNotSpace = true
			break
		}
	}

	qHead := curIdx + 2
	if findNotSpace {
		if s[curIdx] != 'q' || s[curIdx+1] != '=' || s[curIdx+2] == ' ' {
			return "", 0
		}

		// find last quality value index
		for curIdx = qHead; curIdx < tail; curIdx++ {
			if s[curIdx] == ' ' {
				break
			}
		}
	} else {
		if curIdx > tailMin2 && curIdx > sIdx+1 {
			curIdx--
		}
		for ; curIdx < tail; curIdx++ {
			if s[curIdx] != ' ' {
				return "", 0
			}
		}
	}

	// parse quality value
	var qValue float32 = defaultQValue
	if curIdx <= tail && qHead < curIdx {
		const parseBit = 32
		var parseErr error
		var parseValue float64
		parseValue, parseErr = strconv.ParseFloat(s[qHead:curIdx], parseBit)
		if parseErr != nil {
			return
		}
		if parseValue < 0 {
			return
		}
		qValue = float32(parseValue)
	}

	value = s[sRecord[0]:sRecord[1]]
	quality = qValue
	return value, quality
}
