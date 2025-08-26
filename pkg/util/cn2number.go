package util

import "strings"

// Numbers
var chnNumChar = [10]string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}

// Weight positions
var chnUnitSection = [4]string{"", "万", "亿", "万亿"}

// Number weight positions
var chnUnitChar = [4]string{"", "十", "百", "千"}

type chnNameValue struct {
	name    string
	value   int
	secUnit bool
}

// Weight position to node relationships
var chnValuePair = []chnNameValue{{"十", 10, false}, {"百", 100, false}, {"千", 1000, false}, {"万", 10000, true}, {"亿", 100000000, true}}

//func main() {
//	for {
//		var typeStr string
//		var scanStr string
//
//		fmt.Println("1 阿拉伯转中文数字 2 中文数字转阿拉伯数字")
//		fmt.Println("请输入")
//
//		//fmt.Scanf("%s", &a)
//		fmt.Scan(&typeStr)
//
//		fmt.Println("请输入要转换的内容")
//		fmt.Scan(&scanStr)
//		if typeStr == "1" {
//			num, _ := strconv.ParseInt(scanStr, 10, 64)
//			var chnStr = numberToChinese(num)
//			fmt.Println(chnStr)
//		} else {
//			var numInt = chineseToNumber(scanStr)
//			fmt.Println(numInt)
//		}
//	}
//}

// Convert Arabic numbers to Chinese characters
func NumberToChinese(num int64) (numStr string) {
	var unitPos = 0
	var needZero = false

	for num > 0 { //小于零特殊处理
		section := num % 10000 // 已万为小结处理
		if needZero {
			numStr = chnNumChar[0] + numStr
		}
		strIns := sectionToChinese(section)
		if section != 0 {
			strIns += chnUnitSection[unitPos]
		} else {
			strIns += chnUnitSection[0]
		}
		numStr = strIns + numStr
		// When thousands digit is 0, need to add zero in next section
		needZero = (section < 1000) && (section > 0)
		num = num / 10000
		unitPos++
	}
	return
}
func sectionToChinese(section int64) (chnStr string) {
	var strIns string
	var unitPos = 0
	var zero = true
	for section > 0 {
		var v = section % 10
		if v == 0 {
			if !zero {
				zero = true // Need to add zero, ensures only one Chinese zero for consecutive zeros
				chnStr = chnNumChar[v] + chnStr
			}
		} else {
			zero = false                   // At least one digit is not zero
			strIns = chnNumChar[v]         // Chinese number for this position
			strIns += chnUnitChar[unitPos] // Chinese weight position for this digit
			chnStr = strIns + chnStr
		}
		unitPos++ // Shift position
		section = section / 10
	}
	return
}

// Convert Chinese characters to Arabic numbers
func ChineseToNumber(chnStr string) (rtnInt int) {
	var section = 0
	var number = 0
	// Handle special cases like 十一、十二、一百十一、一百十二 separately
	if len(chnStr) == 6 || strings.Contains(chnStr, "百十") {
		chnStr = strings.Replace(chnStr, "十", "一十", -1)
	}
	for index, value := range chnStr {
		var num = chineseToValue(string(value))
		if num > 0 {
			number = num
			if index == len(chnStr)-3 {
				section += number
				rtnInt += section
				break
			}
		} else {
			unit, secUnit := chineseToUnit(string(value))
			if secUnit {
				section = (section + number) * unit
				rtnInt += section
				section = 0

			} else {
				section += (number * unit)

			}
			number = 0
			if index == len(chnStr)-3 {
				rtnInt += section
				break
			}
		}
	}

	return
}
func chineseToUnit(chnStr string) (unit int, secUnit bool) {

	for i := 0; i < len(chnValuePair); i++ {
		if chnValuePair[i].name == chnStr {
			unit = chnValuePair[i].value
			secUnit = chnValuePair[i].secUnit
		}
	}
	return
}
func chineseToValue(chnStr string) (num int) {
	switch chnStr {
	case "零":
		num = 0
		break
	case "一":
		num = 1
		break
	case "二":
		num = 2
		break
	case "三":
		num = 3
		break
	case "四":
		num = 4
		break
	case "五":
		num = 5
		break
	case "六":
		num = 6
		break
	case "七":
		num = 7
		break
	case "八":
		num = 8
		break
	case "九":
		num = 9
		break
	default:
		num = -1
	}
	return
}
