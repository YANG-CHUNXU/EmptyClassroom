//go:build localserver

package main

import (
	"EmptyClassroom/bootstrap"
	"EmptyClassroom/service/model"
	"encoding/gob"
)

func Init() {
	gob.Register(&model.ClassInfo{})
	bootstrap.Init(true)
}
