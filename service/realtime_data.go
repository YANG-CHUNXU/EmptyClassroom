package service

import (
	"EmptyClassroom/config"
	"EmptyClassroom/logs"
	"EmptyClassroom/service/model"
	"EmptyClassroom/utils"
	"context"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LoginURL = "https://jwglweixin.bupt.edu.cn/bjyddx/login"
	QueryURL = "https://jwglweixin.bupt.edu.cn/bjyddx/todayClassrooms"

	LoginUsernameKey = "JW_USERNAME"
	LoginPasswordKey = "JW_PASSWORD"

	loginAESKey = "qzkj1kjghd=876&*"
)

var (
	Token                string
	ErrLoginRejected     = errors.New("login rejected")
	realtimeLoginRequest = utils.HttpPostForm
	realtimeQueryRequest = utils.HttpPostFormWithHeader
)

const currentConfigVersion = "2026-04-13-shahe-building-split-v2"

func Login(ctx context.Context) error {
	userNo := os.Getenv(LoginUsernameKey)
	pwd := os.Getenv(LoginPasswordKey)
	encryptedPwd, err := encodeLoginPassword(pwd)
	if err != nil {
		logs.CtxError(ctx, "encrypt login password failed: %v", err)
		return err
	}
	req := map[string]string{
		"userNo":      userNo,
		"pwd":         encryptedPwd,
		"encode":      "1",
		"captchaData": "",
		"codeVal":     "",
	}
	code, _, body, err := realtimeLoginRequest(ctx, LoginURL, req)
	if err != nil {
		logs.CtxError(ctx, "login failed: %v", err)
		return err
	}
	if code != 200 {
		logs.CtxError(ctx, "login failed - code not 200: %v", err)
		return errors.New("login failed")
	}
	var resp model.LoginResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.CtxError(ctx, "login failed - resp unmarshal failed: %v", err)
		return err
	}
	if resp.Code != "1" {
		msg := strings.TrimSpace(resp.Msg)
		logs.CtxError(ctx, "login failed - code not 1: %s", msg)
		if msg == "" {
			msg = "login failed"
		}
		return fmt.Errorf("%w: %s", ErrLoginRejected, msg)
	}
	Token = resp.Data.Token
	return nil
}

func QueryOne(ctx context.Context, id int) ([]model.JWClassInfo, error) {
	err := Login(ctx)
	if err != nil {
		logs.CtxError(ctx, "login failed: %v", err)
		return nil, err
	}
	header := map[string]string{
		"token": Token,
	}
	req := map[string]string{
		"campusId": resolveRealtimeCampusID(id),
	}
	code, _, body, err := realtimeQueryRequest(ctx, QueryURL, req, header)
	if err != nil {
		logs.CtxError(ctx, "query failed: %v", err)
		return nil, err
	}
	if code != 200 {
		logs.CtxError(ctx, "query failed - code not 200: %v", err)
		return nil, errors.New("query failed")
	}
	var resp model.QueryResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.CtxError(ctx, "query failed - resp unmarshal failed: %v", err)
		return nil, err
	}
	if resp.Code != "1" {
		logs.CtxError(ctx, "query failed - code not 1: %v", err)
		return nil, errors.New("query failed")
	}
	return resp.Data, nil
}

func encodeLoginPassword(pwd string) (string, error) {
	quotedPwd, err := json.Marshal(pwd)
	if err != nil {
		return "", err
	}
	encrypted, err := aesEncryptECBPKCS7([]byte(quotedPwd), []byte(loginAESKey))
	if err != nil {
		return "", err
	}
	firstBase64 := base64.StdEncoding.EncodeToString(encrypted)
	return base64.StdEncoding.EncodeToString([]byte(firstBase64)), nil
}

func aesEncryptECBPKCS7(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	padded := pkcs7Pad(plaintext, blockSize)
	encrypted := make([]byte, len(padded))
	for start := 0; start < len(padded); start += blockSize {
		block.Encrypt(encrypted[start:start+blockSize], padded[start:start+blockSize])
	}
	return encrypted, nil
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	if padding == 0 {
		padding = blockSize
	}
	padded := make([]byte, len(src)+padding)
	copy(padded, src)
	for i := len(src); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func resolveRealtimeCampusID(id int) string {
	switch id {
	case 1:
		return "01"
	case 2:
		return "04"
	default:
		return fmt.Sprintf("%02d", id)
	}
}

func QueryAll(ctx context.Context) (classInfo *model.ClassInfo, err error) {
	classInfo = &model.ClassInfo{
		UpdateAt:       time.Now(),
		ConfigVersion:  currentConfigVersion,
		IsFallback:     map[string]bool{},
		FallbackReason: map[string]string{},
	}
	sysConfig := config.GetConfig()
	for _, campus := range sysConfig.Campus {
		err = ProcessClassTableInfo(ctx, classInfo, campus.Name)
		if err != nil {
			logs.CtxError(ctx, "process failed: %v", err)
			return nil, err
		}
		if campus.HasRealtime {
			jwClassInfo, err := QueryOne(ctx, campus.Id)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil, err
				}
				logs.CtxError(ctx, "query failed: %v", err)
				classInfo.IsFallback[campus.Name] = true
				classInfo.FallbackReason[campus.Name] = describeRealtimeFailure(err)
			}
			// 即使查询报错也不返回，用课表数据进行兜底
			err = ProcessJWClassInfo(ctx, jwClassInfo, classInfo, campus)
			if err != nil {
				logs.CtxError(ctx, "process failed: %v", err)
				return nil, err
			}
		}
	}
	startTime, _ := time.Parse("2006-01-02 15:04:05", sysConfig.Notification.Start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", sysConfig.Notification.End)
	if time.Now().After(startTime) && time.Now().Before(endTime) {
		classInfo.Notification = &sysConfig.Notification
	} else {
		classInfo.Notification = nil
	}
	classTableStartWeek, _ := time.Parse("2006-01-02", sysConfig.ClassTable.StartWeek)
	classTableEndWeek, _ := time.Parse("2006-01-02", sysConfig.ClassTable.EndWeek)
	if time.Now().Before(classTableStartWeek) || time.Now().After(classTableEndWeek.AddDate(0, 0, 1)) {
		classInfo.ClassTable = nil
	} else {
		classInfo.ClassTable = &sysConfig.ClassTable
	}
	classInfo.EmptyReason = buildEmptyReason(classInfo, sysConfig.ClassTable)
	return classInfo, nil
}

func describeRealtimeFailure(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, ErrLoginRejected) {
		return "实时教务登录失败，请检查服务端教务账号配置"
	}
	return "实时教务查询失败，请稍后重试"
}

func buildEmptyReason(classInfo *model.ClassInfo, classTableConfig config.ClassTableConfig) string {
	if classInfo == nil || len(classInfo.CampusInfoMap) > 0 {
		return ""
	}
	if !classTableConfig.IsAvailable {
		if len(classInfo.IsFallback) > 0 {
			return "当前暂无可用教室数据：实时教务查询失败，且当前未启用课表数据兜底。"
		}
		return "当前暂无可用教室数据：当前仓库未启用课表数据兜底。"
	}
	if len(classInfo.IsFallback) > 0 {
		return "当前暂无可用教室数据：实时教务查询失败。"
	}
	return "当前暂无可用教室数据。"
}

func splitShaheTeachingBuilding(campusName string, buildingName string, classroomName string) string {
	if campusName != "沙河" || buildingName != "教学实验综合楼" || classroomName == "" {
		return buildingName
	}
	switch classroomName[0] {
	case 'N':
		return "N"
	case 'S':
		return "S"
	default:
		return buildingName
	}
}

func ProcessJWClassInfo(ctx context.Context, jwClassInfo []model.JWClassInfo, classInfo *model.ClassInfo, campusConfig config.CampusConfig) error {
	sysConfig := config.GetConfig()
	if jwClassInfo == nil {
		return nil
	}
	campusInfo := model.CampusInfo{
		Name:            campusConfig.Name,
		BuildingInfoMap: map[int]*model.BuildingInfo{},
		BuildingIdMap:   map[string]int{},
		MaxBuildingId:   0,
	}
	if classInfo.CampusInfoMap != nil && classInfo.CampusInfoMap[campusConfig.Name] != nil {
		campusInfo = *classInfo.CampusInfoMap[campusConfig.Name]
	}
	campusClassTableConfig := sysConfig.ClassTable.ClassTableMap[campusConfig.Name]
	for _, info := range jwClassInfo {
		classroomList := strings.Split(info.Classrooms, ",")
		for _, classroom := range classroomList {
			for _, replaceConfig := range campusConfig.ReplaceRegex {
				re, err := regexp.Compile(replaceConfig.Regex)
				if err != nil {
					logs.CtxError(ctx, "regex compile failed: %v", err)
					return err
				}
				classroom = re.ReplaceAllString(classroom, replaceConfig.Replace)
			}
			classroomInfo := model.ClassroomInfo{}
			var found bool
			var buildingName string
			buildingName, classroomInfo.Name, found = strings.Cut(strings.Split(classroom, "(")[0], "-")
			if !found {
				logs.CtxWarn(ctx, "classroom format error: %v", classroom)
				continue
			}
			buildingName = splitShaheTeachingBuilding(campusConfig.Name, buildingName, classroomInfo.Name)
			classroomInfo.Size, _ = strconv.ParseInt(strings.Split(strings.Split(classroom, "(")[1], ")")[0], 10, 32)
			classroomInfo.CanTrust = true
			if _, ok := campusInfo.BuildingIdMap[buildingName]; !ok {
				campusInfo.BuildingIdMap[buildingName] = campusInfo.MaxBuildingId
				campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId] = &model.BuildingInfo{
					Name:             buildingName,
					ClassroomInfoMap: map[int]*model.ClassroomInfo{},
					ClassroomIdMap:   map[string]int{},
					ClassMatrix:      [][]int{},
					MaxClassroomId:   0,
				}
				for i := 0; i < 14; i++ {
					campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix = append(campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix, []int{})
				}
				campusInfo.MaxBuildingId++
			}
			buildingId := campusInfo.BuildingIdMap[buildingName]
			if _, ok := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]; !ok {
				campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name] = campusInfo.BuildingInfoMap[buildingId].MaxClassroomId
				classroomInfo.BuildingId = buildingId
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId] = &classroomInfo
				campusInfo.BuildingInfoMap[buildingId].MaxClassroomId++
				for i := 0; i < 14; i++ {
					campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i] = append(campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i], 1)
				}
			} else if !campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]].CanTrust {
				// 覆盖
				classroomInfo.BuildingId = buildingId
				classroomType, typeOk := campusClassTableConfig.TypeMap[classroomInfo.Name]
				if typeOk {
					classroomInfo.Type = classroomType
				} else {
					classroomInfo.Type = campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]].Type
				}
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]] = &classroomInfo
				for i := 0; i < 14; i++ {
					campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i][campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]] = 1
				}
			}
			classroomId := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]
			nodeName, err := strconv.ParseInt(info.NodeName, 10, 32)
			if err != nil {
				logs.CtxWarn(ctx, "node name parse failed: %v", err)
				continue
			}
			campusInfo.BuildingInfoMap[buildingId].ClassMatrix[nodeName-1][classroomId] = 0
		}
	}
	if classInfo.CampusInfoMap == nil {
		classInfo.CampusInfoMap = map[string]*model.CampusInfo{}
	}
	classInfo.CampusInfoMap[campusConfig.Name] = &campusInfo
	return nil
}

func ProcessClassTableInfo(ctx context.Context, classInfo *model.ClassInfo, campusName string) error {
	sysConfig := config.GetConfig()
	classTableStartWeek, err := time.Parse("2006-01-02", sysConfig.ClassTable.StartWeek)
	if err != nil {
		logs.CtxError(ctx, "start week parse failed: %v", err)
		return err
	}
	classTableEndWeek, err := time.Parse("2006-01-02", sysConfig.ClassTable.EndWeek)
	if err != nil {
		logs.CtxError(ctx, "end week parse failed: %v", err)
		return err
	}
	if time.Now().Before(classTableStartWeek) || time.Now().After(classTableEndWeek.AddDate(0, 0, 1)) {
		return nil
	}
	nowWeek := int((time.Now().Unix() - classTableStartWeek.Unix()) / 604800)
	today := int(time.Now().Weekday())

	campusClassTableConfig := sysConfig.ClassTable.ClassTableMap[campusName]
	campusInfo := model.CampusInfo{
		Name:            campusName,
		BuildingInfoMap: map[int]*model.BuildingInfo{},
		BuildingIdMap:   map[string]int{},
		MaxBuildingId:   0,
	}
	if classInfo.CampusInfoMap != nil && classInfo.CampusInfoMap[campusName] != nil {
		campusInfo = *classInfo.CampusInfoMap[campusName]
	}
	for _, classItemInfo := range campusClassTableConfig.Class {
		buildingName, classroomName, found := strings.Cut(classItemInfo.Name, "-")
		if !found {
			logs.CtxWarn(ctx, "classroom format error: %v", classItemInfo.Name)
			continue
		}
		buildingName = splitShaheTeachingBuilding(campusName, buildingName, classroomName)
		if _, ok := campusInfo.BuildingIdMap[buildingName]; !ok {
			campusInfo.BuildingIdMap[buildingName] = campusInfo.MaxBuildingId
			campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId] = &model.BuildingInfo{
				Name:             buildingName,
				ClassroomInfoMap: map[int]*model.ClassroomInfo{},
				ClassroomIdMap:   map[string]int{},
				ClassMatrix:      [][]int{},
				MaxClassroomId:   0,
			}
			for i := 0; i < 14; i++ {
				campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix = append(campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix, []int{})
			}
			campusInfo.MaxBuildingId++
		}
		buildingId := campusInfo.BuildingIdMap[buildingName]
		if _, ok := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomName]; !ok {
			campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomName] = campusInfo.BuildingInfoMap[buildingId].MaxClassroomId
			campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId] = &model.ClassroomInfo{
				Name:       classroomName,
				Size:       0,
				CanTrust:   false,
				BuildingId: buildingId,
			}
			classroomSize, err := strconv.ParseInt(classItemInfo.Seat, 10, 32)
			if err == nil {
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId].Size = classroomSize
			}
			classroomType, typeOk := campusClassTableConfig.TypeMap[classItemInfo.Name]
			if typeOk {
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId].Type = classroomType
			}
			campusInfo.BuildingInfoMap[buildingId].MaxClassroomId++
			for i := 0; i < 14; i++ {
				campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i] = append(campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i], 0)
			}
		} else {
			// 跳过
			continue
		}
		classroomId := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomName]
		for i := 0; i < 14; i++ {
			for _, week := range classItemInfo.Classes[i][today] {
				if week == nowWeek {
					campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i][classroomId] = 1
				}
			}
		}
	}
	if classInfo.CampusInfoMap == nil {
		classInfo.CampusInfoMap = map[string]*model.CampusInfo{}
	}
	classInfo.CampusInfoMap[campusName] = &campusInfo
	return nil
}
