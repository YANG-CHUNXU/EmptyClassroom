# 空教室查询 V2

本项目基于 [Jraaay/EmptyClassroom](https://github.com/Jraaay/EmptyClassroom) fork 后继续开发，用于查询 BUPT 空教室信息。

当前版本在原项目基础上调整了前后端实现、部署方式和数据刷新流程，适配 Vercel Serverless 部署。

相比原项目主要变更：

- 用 go 重构了后端

- 用 React 重构了前端

## 环境变量

- `OUTBOUND_HTTP_TIMEOUT`：后端出站 HTTP 请求超时，使用 Go duration 格式，默认 `15s`
- `BLOB_READ_WRITE_TOKEN`：生产环境必填，用于从 Vercel Blob 读取/写入空教室快照。Vercel 环境中缺失时，普通数据接口会快速返回“数据暂未准备好”，不会退回 Serverless 本地文件缓存
- `CRON_SECRET`：生产环境建议配置，用于保护 `/api/cron/refresh`；Vercel Cron 会自动带上 `Authorization: Bearer <CRON_SECRET>` 请求头
- `JW_USERNAME`：实时教务集成测试使用的账号，仅在显式运行在线测试时需要
- `JW_PASSWORD`：实时教务集成测试使用的密码，仅在显式运行在线测试时需要

## 测试

- 默认执行 `go test ./...`，只运行可离线复现的测试
- 在线实时教务冒烟测试需要显式开启，并同时提供有效凭据：

```bash
RUN_REALTIME_INTEGRATION_TESTS=1 JW_USERNAME=your_username JW_PASSWORD=your_password go test ./service -run 'Test(Login|QueryOne|QueryAll)'
```

- 若未设置 `RUN_REALTIME_INTEGRATION_TESTS=1` 或缺少教务凭据，相关在线测试会被 `skip`，不会让默认测试失败
