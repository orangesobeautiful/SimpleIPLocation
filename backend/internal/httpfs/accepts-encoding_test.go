package httpfs

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

var acptSegFactor = [][]string{
	// value 前面的空格數
	{"", " ", "  ", "   "},
	// value 前半
	{"gz", "def", "b", "zh-", "e", "en-"},
	// value 中間的空格數
	{"", "  ", "   "},
	// value 後半
	{"ip", "late", "r", "TW", "n", "US"},
	// value 後面的空格數
	{"", " ", "  ", "   "},
	// 有無分號
	{"", ";"},
	// 分號後面的空格數
	{"", " ", "  ", "   "},
	// q=、沒東西、不合標準的其他字串
	{"q=", "", "q", "o=", "abc"},
	// q= 和 quality 中間的空格
	{"", " ", "  ", "   ", "    "},
	// 不同的 quality
	{"0.0", "0", "0.2", "0.3", "1", "", "notfloat"},
	// 最後面的空格數
	{"", " ", "  ", "   "},
}

// combineAcpSegFactor
//
// sIdx: 每個遞迴層級選擇的 index
func combineAcpSegFactor(t *testing.T, level int, sIdx []int) {
	if level == len(acptSegFactor) {
		var expectV string
		var expectQ float32

		var testStr string
		for i := 0; i < len(acptSegFactor); i++ {
			testStr += acptSegFactor[i][sIdx[i]]
		}

		splitSemicolon := strings.Split(testStr, ";")
		if len(splitSemicolon) <= 2 {
			valid := true
			if len(splitSemicolon) == 2 {
				qSeg := strings.TrimSpace(splitSemicolon[1])
				if qSeg == "" {
					expectQ = defaultQValue
				} else {
					if len(strings.Split(qSeg, " ")) == 1 {
						qEqualIdx := strings.Index(qSeg, "q=")
						if qEqualIdx >= 0 {
							parseFloat64, pErr := strconv.ParseFloat(qSeg[qEqualIdx+2:], 32)
							if pErr == nil {
								expectQ = float32(parseFloat64)
							} else {
								valid = false
							}
						} else {
							valid = false
						}
					} else {
						valid = false
					}
				}
			} else {
				expectQ = defaultQValue
			}

			if valid {
				expectV = strings.TrimSpace(splitSemicolon[0])
				if len(strings.Split(expectV, " ")) != 1 {
					expectV = ""
					expectQ = 0
				}
			}
		}

		resV, resQ := ParseAcceptValue(testStr, 0, len(testStr))
		if resV != expectV || resQ != expectQ {
			t.Fatalf("(%s) => (%v, %v) , want (%v, %v)", testStr, resV, resQ, expectV, expectQ)
		}
	} else {
		if level == 3 {
			sIdx[level] = sIdx[1]
			combineAcpSegFactor(t, level+1, sIdx)
		} else {
			for i := 0; i < len(acptSegFactor[level]); i++ {
				if level == 0 {
					sIdx = make([]int, len(acptSegFactor))
				}
				sIdx[level] = i
				combineAcpSegFactor(t, level+1, sIdx)
			}
		}
	}
}

func TestParseAcceptValue(t *testing.T) {
	combineAcpSegFactor(t, 0, nil)
}

func permDeal(totalNum int, dealFunc func([]int)) {
	idxAry := make([]int, totalNum)
	for i := 0; i < totalNum; i++ {
		idxAry[i] = i
	}

	var perm func(int)
	perm = func(level int) {
		if level == 1 {
			dealFunc(idxAry)
		} else {
			for i := 0; i < level; i++ {
				perm(level - 1)
				if level%2 == 1 {
					idxAry[i], idxAry[level-1] = idxAry[level-1], idxAry[i]
				} else {
					idxAry[0], idxAry[level-1] = idxAry[level-1], idxAry[0]
				}
			}
		}
	}
	perm(totalNum)
}

// TestParseAcceptHeaderSwitchClientOrder
// 測試 AcceptHeader 交換順序
func TestParseAcceptHeaderSwitchClientOrder(t *testing.T) {
	testCodingList := []string{"gzip", "deflate", "br"}
	testNum := len(testCodingList)

	const expectRes string = "br"
	serverAcptList := []string{"br", "gzip"}
	permDeal(testNum, func(idxAry []int) {
		codingBuffer := bytes.NewBufferString(testCodingList[idxAry[0]])
		for i := 1; i < testNum; i++ {
			codingBuffer.WriteByte(',')
			codingBuffer.WriteString(testCodingList[idxAry[i]])
		}

		coding := codingBuffer.String()
		res := ParseAcceptHeader(coding, serverAcptList)
		if res != expectRes {
			t.Fatalf("\"%s\", server acpt=%+v => %s, want %s", coding, serverAcptList, res, expectRes)
		}
	})
}

// TestParseAcceptHeaderQuality
// 測試 Quality
func TestParseAcceptHeaderQuality(t *testing.T) {
	coding := "gzip;q=0.6, deflate, br;q=0.1"

	const expectRes string = "deflate"
	serverAcptList := []string{"br", "gzip", "deflate"}

	res := ParseAcceptHeader(coding, serverAcptList)
	if res != expectRes {
		t.Fatalf("\"%s\", server acpt=%+v => %s, want %s", coding, serverAcptList, res, expectRes)
	}
}
