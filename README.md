# 空教室查询 V2

> **⚠️ 本项目已停止维护** 
> 
> 该项目不再进行功能更新与 Bug 修复，仅作为学习参考保留，感谢大家的支持，对带来的不便感到十分抱歉🙇

这是一个提供给 BUPT 学生的空教室查询系统，方便同学们进行灵活游击战自习。

如果对你有帮助，欢迎 star！

如果你有更好的想法或者建议，欢迎与我交流！

ps: 终于抽出时间重构了代码，现在整体的体验和稳定性应该都有了一定提升！

相比 V1 变更的地方

- 用 go 重构了后端

- 用 React + Antd 重构了前端

## 环境变量

- `OUTBOUND_HTTP_TIMEOUT`：后端出站 HTTP 请求超时，使用 Go duration 格式，默认 `15s`
- `JW_USERNAME`：实时教务集成测试使用的账号，仅在显式运行在线测试时需要
- `JW_PASSWORD`：实时教务集成测试使用的密码，仅在显式运行在线测试时需要

## 测试

- 默认执行 `go test ./...`，只运行可离线复现的测试
- 在线实时教务冒烟测试需要显式开启，并同时提供有效凭据：

```bash
RUN_REALTIME_INTEGRATION_TESTS=1 JW_USERNAME=your_username JW_PASSWORD=your_password go test ./service -run 'Test(Login|QueryOne|QueryAll)'
```

- 若未设置 `RUN_REALTIME_INTEGRATION_TESTS=1` 或缺少教务凭据，相关在线测试会被 `skip`，不会让默认测试失败
