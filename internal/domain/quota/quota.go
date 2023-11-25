package quota

import (
	"context"
	"github.com/modarreszadeh/arvancloud-interview/internal/setting"
	r "github.com/modarreszadeh/arvancloud-interview/pkg/redis"
	"strconv"
	"time"
)

const (
	RequestRate = "RequestRate"
	DataVolume  = "DataVolume"
	LastUpdate  = "LastUpdate"
)

const KeyPattern = "quota_"

type Quota struct {
	UserId      int
	RequestRate int // request rate per minuet
	DataVolume  int // data volume as byte per month
	LastUpdate  time.Time
}

var (
	redisClient = r.NewRedisClient(&r.Config{Addr: setting.Redis.Addr, Password: setting.Redis.Password,
		DB: setting.Redis.DB})
	ctx = context.Background()
)

func New(userId, requestRate, dataVolume int) *Quota {
	return &Quota{
		DataVolume:  dataVolume,
		RequestRate: requestRate,
		UserId:      userId,
		LastUpdate:  time.Now(),
	}
}

func Get(userId int) (*Quota, bool) {
	var q = Quota{}
	QuotaMap, err := redisClient.HGetAll(ctx, KeyPattern+strconv.Itoa(userId)).Result()
	if err != nil {
		return nil, false
	}
	if len(QuotaMap) == 0 {
		return nil, false
	}
	q.RequestRate, err = strconv.Atoi(QuotaMap[RequestRate])
	q.DataVolume, _ = strconv.Atoi(QuotaMap[DataVolume])
	q.LastUpdate, _ = time.Parse(time.RFC3339, QuotaMap[LastUpdate])
	q.UserId = userId
	return &q, true
}

func Create(q *Quota) bool {
	_, err := redisClient.HSet(ctx, KeyPattern+strconv.Itoa(q.UserId),
		DataVolume, q.DataVolume,
		RequestRate, q.RequestRate,
		LastUpdate, q.LastUpdate,
	).Result()
	if err != nil {
		return false
	}
	return true
}

func Delete(userId int) bool {
	_, err := redisClient.Del(ctx, KeyPattern+strconv.Itoa(userId)).Result()
	if err != nil {
		return false
	}
	return true
}

func Update(userId int, q *Quota) bool {
	if Delete(userId) == true {
		q.UserId = userId
		q.LastUpdate = time.Now()
		isCreate := Create(q)
		if isCreate == true {
			return true
		}
	}
	return false
}
