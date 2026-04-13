package service

import (
	"EmptyClassroom/config"
	"EmptyClassroom/service/model"
	"errors"
	"fmt"
	"testing"
)

func TestDescribeRealtimeFailure(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "nil",
			err:  nil,
			want: "",
		},
		{
			name: "login rejected",
			err:  fmt.Errorf("%w: invalid credentials", ErrLoginRejected),
			want: "实时教务登录失败，请检查服务端教务账号配置",
		},
		{
			name: "temporary failure",
			err:  errors.New("timeout"),
			want: "实时教务查询失败，请稍后重试",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := describeRealtimeFailure(tc.err)
			if got != tc.want {
				t.Fatalf("describeRealtimeFailure() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBuildEmptyReason(t *testing.T) {
	tests := []struct {
		name             string
		classInfo        *model.ClassInfo
		classTableConfig config.ClassTableConfig
		want             string
	}{
		{
			name: "has campus data",
			classInfo: &model.ClassInfo{
				CampusInfoMap: map[string]*model.CampusInfo{
					"西土城": {},
				},
			},
			classTableConfig: config.ClassTableConfig{IsAvailable: false},
			want:             "",
		},
		{
			name: "no class table and fallback",
			classInfo: &model.ClassInfo{
				IsFallback: map[string]bool{
					"西土城": true,
				},
			},
			classTableConfig: config.ClassTableConfig{IsAvailable: false},
			want:             "当前暂无可用教室数据：实时教务查询失败，且当前未启用课表数据兜底。",
		},
		{
			name:             "no class table only",
			classInfo:        &model.ClassInfo{},
			classTableConfig: config.ClassTableConfig{IsAvailable: false},
			want:             "当前暂无可用教室数据：当前仓库未启用课表数据兜底。",
		},
		{
			name: "class table available but fallback",
			classInfo: &model.ClassInfo{
				IsFallback: map[string]bool{
					"沙河": true,
				},
			},
			classTableConfig: config.ClassTableConfig{IsAvailable: true},
			want:             "当前暂无可用教室数据：实时教务查询失败。",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildEmptyReason(tc.classInfo, tc.classTableConfig)
			if got != tc.want {
				t.Fatalf("buildEmptyReason() = %q, want %q", got, tc.want)
			}
		})
	}
}
