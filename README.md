# 日志流聚合服务

这是一个只使用 Go 标准库实现的日志处理 API 服务。它接收不同应用和主机产生的日志，维护日志源与采集状态，支持条件检索、告警规则、通知、标签、保留策略、批量处理和运行健康检查。

## 运行

在项目目录执行：

```bash
go run ./cmd/server
```

服务默认监听 `:8080`。

环境变量：

- `PORT`：服务端口，默认 8080
- `ADDR`：完整监听地址，设置后覆盖 PORT
- `MAX_PAGE_SIZE`：最大分页大小，默认 100
- `AUTH_TOKEN`：Bearer Token，默认 `demo-token`
- `LOG_LEVEL`：日志级别，可选 debug、info、warn、error

## 目录

```text
cmd/server/main.go          服务入口
internal/app/               依赖装配
internal/config/            环境配置
internal/model/             日志源、日志条目、告警等模型
internal/store/             数据访问接口和内存实现
internal/service/           日志处理与状态流转
internal/handler/           HTTP 路由和中间件
pkg/httpx/                  HTTP 响应与分页辅助
pkg/idgen/                  ID 生成
pkg/logger/                 分级日志
```

## 主要接口

所有 `/api/` 请求都需要 `Authorization: Bearer demo-token`。

- `/api/sources`：维护日志源及其 active、paused 状态
- `/api/entries`：写入、查询和清理日志条目
- `/api/indexes`：维护日志检索索引
- `/api/queries`：保存日志检索条件
- `/api/alert-rules`：维护告警规则和启停状态
- `/api/alerts`：处理告警事件的确认与解决流程
- `/api/tags`：维护日志标签
- `/api/notifications`：读取告警通知并标记已读
- `/api/retention-policies`：按日志源应用保留策略
- `/api/health`：检查日志源和服务运行状态
- `/api/stats`：查询日志采集计数、级别分布和错误率
- `/api/sources/{id}/export`：导出单个日志源的记录

## 测试

```bash
go test ./...
```
