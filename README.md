# 航空地勤周转保障平台

面向机场地勤团队的航班过站保障、资源调度、异常延误和任务签收系统。

## 快速启动

```bash
cp .env.example .env && docker compose up -d
```

## 访问地址或 CLI 示例

前端：<http://localhost:20108>

后端健康检查：<http://localhost:21108/health>


## 过站放行联动（核心业务规则）

放行接口：`POST /api/flight-turnaround/:id/release`

- **放行门禁（三项必须同时满足，缺一整次拒绝，HTTP 409）**
  1. 航班下所有地勤任务必须是 `SIGNED`
  2. 未关闭延误事件（`resolved_at` 为空）必须为零
  3. 资源预约不能存在无法自动释放的状态（`CONFLICT` 必须先人工解决；普通 `ACTIVE` 会在放行事务中被一次性释放，放行后预约不再 ACTIVE）
- 拒绝响应 `error.details` 分别返回 `unsigned_tasks` / `open_delays` / `blocking_bookings` 三类阻塞清单；失败时航班、任务、预约全部保持原样。
- 门禁通过后单事务完成：航班进入 `READY`（写 `ready_at`）→ 该航班下 `ACTIVE` 预约一次性置为 `RELEASED`、对应资源回到 `AVAILABLE` → 任务与过站时间线同步刷新。
- **重复放行 / 并发提交只生效一次**：按航班 ID 互斥 + 条件更新（仅非 READY/DEPARTED 可推进），后到者得到 `RELEASE_ALREADY_DONE` / `RELEASE_RACE_LOST`。
- 看板（`/dashboard`）逐航班展示阻塞明细，并可就地签收任务、关闭延误、解决预约冲突后再放行。
- 配套写接口：`POST /api/ground-task/:id/sign`、`POST /api/delay-event/:id/resolve`、`POST /api/resource-booking/:id/resolve-conflict`；门禁快照：`GET /api/flight-turnaround/:id/gate`。

## 本地开发方式

- 前端：`cd frontend && npm install && npm run dev`
- 后端：进入 `backend` 后 `go run .`（本地默认使用纯 Go SQLite 文件 `ground-turn.db`，无需起 MySQL，首启自动建表并灌入 3 个演示航班；容器内存在 `DB_HOST` 时切换为 MySQL）
- 接口统一挂在 `/api`。Vite dev server 已把 `/api`、`/health` 代理到 `http://localhost:3000`。



## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design + Redux Toolkit |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 部署 | Docker Compose |

## 项目目录结构

```text
frontend/src/api, stores, types, constants, constructors, components/common, hooks, pages, router, utils, mocks
backend/src/routes, controllers, services, models, repositories, middlewares, constants, constructors, utils, types, config
```

## 环境变量说明

- `COMPOSE_PROJECT_NAME`: Compose 项目名，默认 `ground-turn`
- `FRONTEND_PORT`: 前端端口，默认 `20108`
- `BACKEND_PORT`: 后端端口，默认 `21108`
- `DB_PORT`: 数据库宿主机端口
- `DB_USER/DB_PASSWORD/DB_NAME`: 本地数据库凭据

## Docker 部署说明

- 根 Compose 文件不写 `version`，顶层 `name: ground-turn`。
- 容器名均使用 `${COMPOSE_PROJECT_NAME:-ground-turn}` 前缀。
- 数据库使用命名卷，避免绑定中文路径。
- 常见问题：端口占用时修改 `.env` 中端口后重启；需要重置数据时执行 `docker compose down -v`。

## 枚举/常量出现位置清单

- GroundTaskType: constants/GroundTaskType（前后端）、types/GroundTaskType（前端重复定义）、constructors、logTemplates、errorMessages、TasksPage 筛选/展示、后端 seed/controller 均有引用。
- TurnaroundStatus: constants/TurnaroundStatus（前后端）、types/TurnaroundStatus（前端重复定义）、constructors、logTemplates、errorMessages、Dashboard/Turnarounds 筛选器、StatusBadge/TurnaroundTimeline 展示、后端 release service 条件更新。
- ResourceStatus: constants/ResourceStatus（前后端）、types/ResourceStatus（前端重复定义）、constructors、logTemplates、errorMessages、ResourcesPage 筛选展示、放行事务资源置 AVAILABLE。
- GroundTaskStatus（DISPATCHED/IN_PROGRESS/SIGNED/BLOCKED）: 后端 constants/GroundTaskStatus.go、models/GroundTask、repositories、constructors、logTemplates、放行门禁 service；前端 constants/GroundTaskStatus.ts、types/GroundTask、TasksPage、BlockerDetails、StatusBadge、seedData。
- BookingStatus（ACTIVE/RELEASED/CONFLICT）: 后端 constants/BookingStatus.go、models/ResourceBooking、repositories（批量释放/解冲突）、logTemplates、放行门禁 service；前端 constants/BookingStatus.ts、types/ResourceBooking、ResourcesPage、BlockerDetails、useResourceConflict、seedData。
- 放行错误码（RELEASE_BLOCKED/RELEASE_ALREADY_DONE/RELEASE_RACE_LOST/TASK_ALREADY_SIGNED 等）: 后端 constants/errorCodes.go、errorMessages.go、controllers/response.go；前端 constants/errorCodes.ts、errorMessages.ts、types/Api.ts、ReleaseStore、DashboardPage。

## 为什么会牵一发动全身

实体字段、枚举、日志模板、错误消息、构造器、筛选器和展示组件被刻意拆散到多个目录；修改一个状态值通常需要同步类型、构造器、服务、控制器、store、页面、README 与数据库种子。

## License

MIT
