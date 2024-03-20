package redis

import (
	"encoding/json"
	"reflect"
	"time"

	red "github.com/go-redis/redis"
	"github.com/tidwall/gjson"

	"github.com/130-133/go-common/utils/help"
)

type SetFunc func() ([]byte, error)
type ZSetFunc func() ([]red.Z, error)
type HSetFunc func() (map[string]interface{}, error)

type LoadResult struct {
	cmd      cmdType
	method   string
	data     []byte
	list     []string
	zSet     []red.Z
	hash     map[string]interface{}
	cacheHit bool
}
type cmdType string

const (
	Str  cmdType = "string"
	Hash cmdType = "hash"
	List cmdType = "list"
	ZSet cmdType = "zset"
)

// LoadSetEx 读取或写入带过期时间的map数据进str
func (r *MRedis) LoadSetEx(key string, ttl time.Duration, fc SetFunc) (LoadResult, error) {
	var (
		res = LoadResult{cmd: Str, method: "LoadSetEx"}
		err error
	)
	if r.Exists(key).Val() > 0 {
		cli := r.Get(key)
		if cli.Err() != nil {
			goto Func
		}
		res.data = []byte(cli.Val())
		res.cacheHit = true
		return res, nil
	}
Func:
	res.data, err = fc()
	if err != nil {
		return res, err
	}
	r.SetNX(key, string(res.data), ttl)
	return res, nil
}

// LoadHMSetEx 读取或写入带过期时间的map数据进hash
func (r *MRedis) LoadHMSetEx(key string, ttl time.Duration, fc SetFunc) (LoadResult, error) {
	var (
		res = LoadResult{cmd: Hash, method: "LoadHMSetEx"}
		err error
	)
	if r.Exists(key).Val() > 0 {
		cli := r.HGetAll(key)
		if cli.Err() != nil {
			goto Func
		}
		res.data, _ = json.Marshal(cli.Val())
		res.cacheHit = true
		return res, nil
	}

Func:
	res.data, err = fc()
	if err != nil {
		return res, err
	}
	mapData := make(map[string]interface{})
	_ = json.Unmarshal(res.data, &mapData)
	gjson.ParseBytes(res.data).ForEach(func(key, value gjson.Result) bool {
		vType := reflect.TypeOf(value)
		switch vType.Kind() {
		case reflect.Map, reflect.Struct:
			mapData[key.String()] = help.ToString(value.Value())
		default:
			mapData[key.String()] = value.Value()
		}
		return true
	})
	r.Pipelined(func(pp red.Pipeliner) error {
		pp.HMSet(key, mapData)
		pp.Expire(key, ttl)
		return nil
	})
	return res, nil
}

// Deprecated
func (r *MRedis) LoadMulitHMSetEx(key string, ttl time.Duration, fc SetFunc) (LoadResult, error) {
	return r.LoadMultiHMSetEx(key, ttl, fc)
}

// LoadMultiHMSetEx 读取或写入带过期时间的map数据进hash
func (r *MRedis) LoadMultiHMSetEx(key string, ttl time.Duration, fc SetFunc) (LoadResult, error) {
	var (
		res = LoadResult{cmd: Hash, method: "LoadMultiHMSetEx"}
		rl  LoadResult
		err error
	)
	rl, err = r.LoadHMSetEx(key, ttl, fc)
	rl.method = res.method
	return rl, err
}

// LoadHSetEx 读取或写入带过期时间的map数据进hash
func (r *MRedis) LoadHSetEx(key string, field string, ttl time.Duration, fc SetFunc) (LoadResult, error) {
	var (
		res = LoadResult{cmd: Hash, method: "LoadHSetEx"}
		err error
	)
	if r.HExists(key, field).Val() {
		cli := r.HGet(key, field)
		if cli.Err() != nil {
			goto Func
		}
		res.data, _ = cli.Bytes()
		res.cacheHit = true
		return res, nil
	}

Func:
	res.data, err = fc()
	if err != nil {
		return res, err
	}
	existsTTL := r.TTL(key).Val()
	r.Pipelined(func(pp red.Pipeliner) error {
		pp.HSetNX(key, field, res.String())
		if existsTTL.Seconds() <= 0 {
			pp.Expire(key, ttl)
		}
		return nil
	})
	return res, nil
}

func (r *MRedis) LoadListPushEx(key string, ttl time.Duration, fc SetFunc) (LoadResult, error) {
	var (
		res = LoadResult{cmd: List, method: "LoadListPushEx"}
		err error
	)
	if r.Exists(key).Val() > 0 {
		length := r.LLen(key).Val()
		cli := r.LRange(key, 0, length)
		if cli.Err() != nil {
			goto Func
		}
		res.list = cli.Val()
		res.cacheHit = true
		return res, nil
	}
Func:
	res.data, err = fc()
	if err != nil {
		return res, err
	}
	var listData []interface{}
	listStr := make([]string, 0)
	listInterface := make([]interface{}, 0)
	_ = json.Unmarshal(res.data, &listData)
	for _, v := range listData {
		bytes, _ := json.Marshal(v)
		listStr = append(listStr, string(bytes))
		listInterface = append(listInterface, string(bytes))
	}
	if len(listStr) == 0 {
		return res, nil
	}
	res.list = listStr
	_, _ = r.Pipelined(func(pp red.Pipeliner) error {
		pp.RPush(key, listInterface...)
		pp.Expire(key, ttl)
		return nil
	})
	return res, nil
}

func (r *MRedis) LoadZSetEx(key string, ttl time.Duration, fc ZSetFunc) (LoadResult, error) {
	var (
		res = LoadResult{cmd: ZSet, method: "LoadZSetEx"}
		err error
	)
	if r.Exists(key).Val() > 0 {
		cli := r.ZRevRangeWithScores(key, 0, -1)
		if cli.Err() != nil {
			goto Func
		}
		res.zSet = cli.Val()
		res.data, err = json.Marshal(res.zSet)
		res.cacheHit = true
		return res, err
	}
Func:
	res.zSet, err = fc()
	if err != nil {
		return res, err
	}
	_, _ = r.Pipelined(func(pp red.Pipeliner) error {
		pp.ZAdd(key, res.zSet...)
		pp.Expire(key, ttl)
		return nil
	})
	return res, err
}

func (r LoadResult) Unmarshal(data interface{}) error {
	switch r.cmd {
	case List:
		if r.list == nil {
			return nil
		}
		var tmp []interface{}
		for _, v := range r.list {
			var val interface{}
			if err := json.Unmarshal([]byte(v), &val); err != nil {
				val = v
			}
			tmp = append(tmp, val)
		}
		res, _ := json.Marshal(tmp)
		return json.Unmarshal(res, data)
	case Hash:
		switch r.method {
		case "LoadMultiHMSetEx", "LoadMulitHMSetEx":
			mapData := make(map[string]interface{})
			json.Unmarshal(r.data, &mapData)
			for k, v := range mapData {
				tmp := make(map[string]interface{})
				if err := json.Unmarshal([]byte(v.(string)), &tmp); err == nil {
					mapData[k] = tmp
				} else {
					mapData[k] = v
				}
			}
			allTmp, _ := json.Marshal(mapData)
			return help.ParseJsonToStruct(allTmp, data)
		}
	}
	if r.data == nil {
		return nil
	}
	return help.ParseJsonToStruct(r.data, data)
}

func (r LoadResult) String() string {
	if r.cmd != Str && r.cmd != Hash {
		return ""
	}
	return string(r.data)
}

func (r LoadResult) List() []string {
	if r.cmd != List {
		return nil
	}
	return r.list
}

func (r LoadResult) ZSet() []red.Z {
	if r.cmd != ZSet {
		return nil
	}
	return r.zSet
}

func (r LoadResult) Byte() []byte {
	return r.data
}

func (r LoadResult) CacheHit() bool {
	return r.cacheHit
}
