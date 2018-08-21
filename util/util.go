package util

import (
	"fmt"
	"crypto/md5"
	"hash"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"strconv"
	"encoding/json"
	"os"
	"github.com/muesli/cache2go"
	"time"
)

type P map[string]interface{}

func Md5(s ...interface{}) (r string) {
	return Hash("md5", s...)
}
func Hash(algorithm string, s ...interface{}) (r string) {
	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha2", "sha256":
		h = sha256.New()
	}
	for _, value := range s {
		switch value.(type) {
		case []byte:
			h.Write(value.([]byte))
		default:
			h.Write([]byte(ToString(value)))
		}
	}
	r = hex.EncodeToString(h.Sum(nil))
	return
}
func ToString(v interface{}) string {
	if v != nil {
		switch v.(type) {
		case bson.ObjectId:
			return v.(bson.ObjectId).Hex()
		case []byte:
			return string(v.([]byte))
		case *P, P:
			var p P
			switch v.(type) {
			case *P:
				if v.(*P) != nil {
					p = *v.(*P)
				}
			case P:
				p = v.(P)
			}
			var keys []string
			for k := range p {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			r := "P{"
			for _, k := range keys {
				r = JoinStr(r, k, ":", p[k], " ")
			}
			r = JoinStr(r, "}")
			return r
		case int64:
			return strconv.FormatInt(v.(int64), 10)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
func JoinStr(val ...interface{}) (r string) {
	for _, v := range val {
		r += ToString(v)
	}
	return
}

//string 转P
func JsonDecode(b []byte) (p *map[string]interface{}) {
	p = &map[string]interface{}{}
	err := json.Unmarshal(b, p)
	if err != nil {
		fmt.Print(err)
	}
	return
}

func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case P:
		return len(v.(P)) == 0
	}
	return ToString(v) == ""
}
func (p *P) ToInt(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		(*p)[k] = ToInt(v)
	}
}
func ToInt(s interface{}, default_v ...int) int {
	i, e := strconv.Atoi(ToString(s))
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return i
}

func WriteFile(url string, body []byte) bool {
	f, err := os.Create(url)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		_, err = f.Write(body)
		if err != nil {
			return false
		} else {
			return true
		}
	}
}


/****************************------------------以下方法为缓存-------------------------************************************/
type Cacha struct {
	openId   string
	moreData []byte
}

//添加缓存
func AddCache(token, openId string) bool {
	//创建缓存表,有则忽略，无则创建
	cache := cache2go.Cache("Cache")

	val := Cacha{openId, []byte{}}
	cache.Add(token, 120*time.Minute, &val)

	// 验证是否存在
	res, err := cache.Value(token)
	if err == nil {
		fmt.Print(res.Data().(Cacha).openId)
		return true
	} else {
		return false
	}

}
//获取缓存
func GetCache(token string) string {
	//创建缓存表,有则获取Cache表，无则创建
	cache := cache2go.Cache("Cache")
	res, err := cache.Value(token)
	if err == nil {
		return res.Data().(Cacha).openId
	} else {
		return ""
	}

}
