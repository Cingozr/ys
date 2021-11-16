package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var lastRecordTimeStamp int64

type IJsonStorage interface {
	Save(storage *MemoryStorage) error
}

type JsonStorage struct {
	wg                 sync.WaitGroup
	ctx                context.Context
	fromTimerToJsonChn chan *struct{}
	storage            *MemoryStorage
}

type JsonStorageModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewJsonStorage(ctx context.Context, fromTimerToJsonChn chan *struct{}) *JsonStorage {

	var res = InsertData()
	return &JsonStorage{
		ctx:                ctx,
		fromTimerToJsonChn: fromTimerToJsonChn,
		storage:            NewMemoryStorage(res),
	}
}

func (j *JsonStorage) IntervalSave() {
	j.wg.Add(1)
	go func() {
		for {
			select {
			case <-j.fromTimerToJsonChn:
				lastRecordTimeStamp = time.Now().UnixNano()
				fileName := fmt.Sprintf("./tmp/%v-data.json", lastRecordTimeStamp)
				if responseMemoryModel, err := j.storage.GetAll(); err != nil {
					log.Println(err)
				} else {
					saveModels, _ := json.MarshalIndent(responseMemoryModel, "", " ")
					if err := ioutil.WriteFile(fileName, saveModels, 0644); err != nil {
						j.wg.Done()
						break
					}
				}
			case <-j.ctx.Done():
				j.wg.Done()
				break

			}
		}
	}()
	j.wg.Wait()
}

func InsertData() *map[string]string {
	const filePath = "tmp"
	files, _ := ioutil.ReadDir(filePath)
	var jsonDataToMemoryStorage = new(map[string]string)
	for _, f := range files {
		jsonFile, _ := os.Open(path.Join(filePath, f.Name()))
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &jsonDataToMemoryStorage)
	}
	return jsonDataToMemoryStorage
}
