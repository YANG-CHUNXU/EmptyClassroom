package service

import (
	"EmptyClassroom/config"
	"EmptyClassroom/service/model"
	"context"
	"testing"
)

func TestProcessJWClassInfoSplitsShaheTeachingBuildingByWing(t *testing.T) {
	var shaheConfig config.CampusConfig
	for _, campus := range config.GetConfig().Campus {
		if campus.Name == "沙河" {
			shaheConfig = campus
			break
		}
	}
	if shaheConfig.Name == "" {
		t.Fatal("沙河校区配置不存在")
	}

	classInfo := &model.ClassInfo{}
	jwClassInfo := []model.JWClassInfo{
		{
			Classrooms: "教学实验综合楼-N101(120)",
			NodeName:   "1",
		},
		{
			Classrooms: "教学实验综合楼-S201(80)",
			NodeName:   "2",
		},
	}

	if err := ProcessJWClassInfo(context.Background(), jwClassInfo, classInfo, shaheConfig); err != nil {
		t.Fatalf("ProcessJWClassInfo() error = %v", err)
	}

	campusInfo := classInfo.CampusInfoMap["沙河"]
	if campusInfo == nil {
		t.Fatal("未生成沙河校区数据")
	}

	if _, ok := campusInfo.BuildingIdMap["N"]; !ok {
		t.Fatalf("building_id_map 缺少 N: %#v", campusInfo.BuildingIdMap)
	}
	if _, ok := campusInfo.BuildingIdMap["S"]; !ok {
		t.Fatalf("building_id_map 缺少 S: %#v", campusInfo.BuildingIdMap)
	}
	if _, ok := campusInfo.BuildingIdMap["教学实验综合楼"]; ok {
		t.Fatalf("building_id_map 不应再包含未拆分的教学实验综合楼: %#v", campusInfo.BuildingIdMap)
	}
}
