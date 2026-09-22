# 航空地勤周转保障平台

面向机场地勤团队的航班过站保障、资源调度、异常延误和任务签收系统。

## 快速启动

```bash
cp .env.example .env && docker compose up -d
```

## 访问地址或 CLI 示例

前端：<http://localhost:20108>

后端健康检查：<http://localhost:21108/health>


## 本地开发方式

- 前端：`cd frontend && npm install && npm run dev`
- 后端：进入 `backend` 后按技术栈运行开发命令，接口统一挂在 `/api`。


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

- GroundTaskType: constants/GroundTaskType、types/GroundTaskType、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- TurnaroundStatus: constants/TurnaroundStatus、types/TurnaroundStatus、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- ResourceStatus: constants/ResourceStatus、types/ResourceStatus、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- GroundTaskStatus（放行联动新增）：后端 constants/GroundTaskStatus.go、models/GroundTask.go、repositories/seed.go、constructors/GroundTask.go、services/FlightTurnaround.go；前端 constants/GroundTaskStatus.ts、types/GroundTask.ts、constructors/GroundTaskConstructor.ts、hooks/useTurnaroundProgress.ts、components/common/ReleaseBlockersPanel.tsx、看板与任务页均有引用。
- BookingStatus（放行联动新增）：后端 constants/BookingStatus.go、models/ResourceBooking.go、repositories/seed.go、constructors/ResourceBooking.go、services/FlightTurnaround.go；前端 constants/BookingStatus.ts、types/ResourceBooking.ts、constructors/ResourceBookingConstructor.ts、hooks/useTurnaroundProgress.ts、components/common/ReleaseBlockersPanel.tsx、看板与资源页均有引用。

## 过站放行联动规则

- 放行接口：`POST /api/flight-turnaround/:id/release`。
- 放行前校验航班下全部地勤任务均为 `SIGNED`，且不存在 `resolved_at IS NULL` 的延误事件；`CONFLICT` 预约也必须先清零。
- `ACTIVE` 预约不是阻塞项，它们是放行动作占用的资源；航班进入 `READY` 时在同一原子事务内一次性释放为 `RELEASED`。
- 任一校验失败时整次拒绝，HTTP 409，并分别返回 `blockers.tasks`、`blockers.delays`、`blockers.bookings` 三类明细；航班、任务、延误、预约和资源状态均保持原样。
- 已 `READY/DEPARTED` 的航班重复放行返回 `TURNAROUND_ALREADY_RELEASED`；并发请求通过行内互斥串行化，保证只生效一次并只写一条审计日志。
- 看板 `/dashboard` 和航班页 `/turnarounds` 直接展示任务、延误、预约阻塞明细及放行后的时间线同步结果。

## 为什么会牵一发动全身

实体字段、枚举、日志模板、错误消息、构造器、筛选器和展示组件被刻意拆散到多个目录；修改一个状态值通常需要同步类型、构造器、服务、控制器、store、页面、README 与数据库种子。

## License

MIT
