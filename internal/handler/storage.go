package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/modarreszadeh/arvancloud-interview/internal/domain/quota"
	"github.com/modarreszadeh/arvancloud-interview/pkg/constants"
	"github.com/modarreszadeh/arvancloud-interview/pkg/helper"
	"github.com/modarreszadeh/arvancloud-interview/pkg/queue"
	"net/http"
	"strconv"
	"time"
)

var (
	storageQueue = queue.New(1, 100_000)
)

type ObjectStorage struct {
	UserId     int    `json:"userId"`
	ObjectId   string `json:"objectId"`
	FileName   string `json:"fileName"`
	BucketName string `json:"bucketName"`
	SizeInByte int    `json:"sizeInByte"`
}

func HandleStorageRequest(c echo.Context) error {
	var object ObjectStorage
	userId, _ := strconv.Atoi(helper.GetQueryParameters(c.QueryString())["userid"])
	err := c.Bind(&object)
	object.UserId = userId
	if err != nil {
		return err
	}
	storageQueue.Enqueue(object)

	go storageQueue.DispatchProcess(objectStoreProcess)

	return c.JSON(http.StatusOK, "ok")
}

func objectStoreProcess(process interface{}) {
	if object, ok := process.(ObjectStorage); ok {
		q, found := quota.Get(object.UserId)
		if found {
			if q.DataVolume < object.SizeInByte {
				fmt.Printf("Your volume quota has been reached\n")
				return
			}
			q.DataVolume -= object.SizeInByte
			quota.Update(object.UserId, q)
		} else {
			q := quota.New(object.UserId, 120, 50*constants.Megabyte)
			quota.Create(q)
		}

		time.Sleep(500 * time.Millisecond) // certain delay for data to be stored in the bucket
		fmt.Printf("object %.4s store in bucket with fileName: %s\n", object.ObjectId, object.FileName)
	}
}
