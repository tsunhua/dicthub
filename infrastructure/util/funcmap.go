package util

import (
	"app/service/dict/model"
	"encoding/json"
	"time"

	"github.com/Masterminds/sprig/v3"
)

type Title struct {
	Dict          model.DictBO
	SelectedCatId string
}

func DictHubFuncMap() map[string]interface{} {
	m := sprig.GenericFuncMap()
	// show the time before now's duration user-friendly
	m["durf"] = func(t time.Time) string {
		return GetDurationFriendly(GetCurrentShanghaiTime(), ToShanghaiTime(t)) + "前"
	}
	// reserve specific count of s, and end with "...".
	m["abbrev"] = func(count int, s string) string {
		rs := []rune(s)
		if count <= 0 || count >= len(rs) {
			return s
		}
		return string(rs[0:count+1]) + "..."
	}
	m["json"] = func(val interface{}) string {
		bs, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(bs)
	}
	// pretty json
	m["jsonp"] = func(val interface{}) string {
		bs, err := json.MarshalIndent(val, "", "\t")
		if err != nil {
			return ""
		}
		return string(bs)
	}

	m["title"] = func(dict model.DictBO, selectedCatId string) Title {
		return Title{
			Dict:          dict,
			SelectedCatId: selectedCatId,
		}
	}
	return m
}
