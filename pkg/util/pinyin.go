package util

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// 中文转拼音
//
// sep 多个汉字间的分隔符，如设置 _ 则 中国 => zhong_guo
func Zh2Pinyin(s string, sep ...string) string {
	dSep := ``
	if len(sep) > 0 {
		dSep = sep[0]
	}
	return strings.Join(pinyin.LazyConvert(s, nil), dSep)
}
