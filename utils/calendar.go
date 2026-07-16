package utils

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"server/model/other"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var solarTerm = []string{
	"立春", "雨水", "惊蛰", "春分", "清明", "谷雨",
	"立夏", "小满", "芒种", "夏至", "小暑", "大暑",
	"立秋", "处暑", "白露", "秋分", "寒露", "霜降",
	"立冬", "小雪", "大雪", "冬至", "小寒", "大寒",
}

func GetCalendar(dateStr string) (other.Calendar, error) {
	resp, err := http.Get("https://www.rili.com.cn/rili/json/today/" + dateStr + ".js")
	if err != nil {
		return other.Calendar{}, err
	}
	//发送http请求后一定要关闭
	defer resp.Body.Close()
	//把数据一次性全部读进内存
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return other.Calendar{}, err
	}

	var jsonStr string

	//正则提取json
	reg := regexp.MustCompile(`(?s)\((.*?)\);`)
	//FindAllStringSubmatch 返回的是二维切片 [][]string,result = [[0] = 完整匹配的字符串,[1] = 第1个括号捕获的内容...]
	result := reg.FindAllStringSubmatch(string(body), -1)

	//至少匹配到一组结果并且第1组结果里至少有2个元素
	if len(result) > 0 && len(result[0]) > 1 {
		// 返回第1个括号捕获的内容
		jsonStr = result[0][1]
	} else {
		return other.Calendar{}, errors.New("failed to get data")
	}

	jieqi := gjson.Get(jsonStr, "jieqi.jieqi").Str

	calendar := other.Calendar{
		Date:         gjson.Get(jsonStr, "yangli.date").Str + " " + gjson.Get(jsonStr, "yangli.xingqi").Str,
		LunarDate:    gjson.Get(jsonStr, "nongli.yueri").Str,
		Ganzhi:       gjson.Get(jsonStr, "nongli.ganzhi").Str,
		Zodiac:       gjson.Get(jsonStr, "xingzuo.xingzuo").Str + "座",
		DayOfYear:    "今年第" + strconv.FormatInt(gjson.Get(jsonStr, "nian_index").Int(), 10) + "天",
		SolarTerm:    jieqi + "第" + strconv.FormatInt(gjson.Get(jsonStr, "jieqi.jieqi_index").Int(), 10) + "天　距离" + nextSolarTerm(jieqi) + "还有" + strconv.FormatInt(gjson.Get(jsonStr, "jieqi.jieqi_next").Int(), 10) + "天",
		Auspicious:   strings.ReplaceAll(gjson.Get(jsonStr, "yi").Str, ",", " "),
		Inauspicious: strings.ReplaceAll(gjson.Get(jsonStr, "ji").Str, ",", " "),
	}

	return calendar, nil
}

// 在 solarTerm 数组中查找当前节气，返回下一个节气。循环到末尾则回到开头（节气是循环的）。
func nextSolarTerm(currentTerm string) string {
	for i, term := range solarTerm {
		if term == currentTerm {
			if i == len(solarTerm)-1 {
				return solarTerm[0]
			}
			return solarTerm[i+1]
		}
	}
	return ""
}
