import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import type { GroundResource } from "../types/GroundResource";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { DelayEvent } from "../types/DelayEvent";

// 与后端 config/seed.go 保持一致的本地兜底数据：
// TA1 门禁全绿；TA2 三类阻塞齐全；TA3 仅延误阻塞。
export const mockData: {
  flightTurnaround: FlightTurnaround[];
  groundTask: GroundTask[];
  groundResource: GroundResource[];
  resourceBooking: ResourceBooking[];
  delayEvent: DelayEvent[];
} = {
  flightTurnaround: [
    {
      id: 1, flight_no: "MU518", aircraft_reg: "B-2201", stand_no: "103",
      arrival_time: "2026-09-22T08:00:00Z", departure_time: "2026-09-22T08:55:00Z",
      turnaround_status: "IN_SERVICE", delay_reason: "", ready_at: null,
      progress: { total_tasks: 3, signed_tasks: 3, open_delays: 0, active_bookings: 3 },
      gate: { releasable: true, unsigned_tasks: [], open_delays: [], blocking_bookings: [] }
    },
    {
      id: 2, flight_no: "CA183", aircraft_reg: "B-5520", stand_no: "117",
      arrival_time: "2026-09-22T08:40:00Z", departure_time: "2026-09-22T09:40:00Z",
      turnaround_status: "IN_SERVICE", delay_reason: "货舱装卸等待", ready_at: null,
      progress: { total_tasks: 3, signed_tasks: 0, open_delays: 1, active_bookings: 1 },
      gate: {
        releasable: false,
        unsigned_tasks: [
          { id: 4, task_type: "BAGGAGE", team_id: 4, status: "IN_PROGRESS", deadline: "2026-09-22 09:10", blocker_note: "传送带车未到位", reason: "任务 BAGGAGE 当前状态 IN_PROGRESS，放行前必须完成签收（SIGNED）" },
          { id: 5, task_type: "WATER_SERVICE", team_id: 5, status: "BLOCKED", deadline: "2026-09-22 09:15", blocker_note: "清水车故障", reason: "任务 WATER_SERVICE 当前状态 BLOCKED，放行前必须完成签收（SIGNED）" },
          { id: 6, task_type: "PUSHBACK", team_id: 6, status: "DISPATCHED", deadline: "2026-09-22 09:35", blocker_note: "", reason: "任务 PUSHBACK 当前状态 DISPATCHED，放行前必须完成签收（SIGNED）" }
        ],
        open_delays: [
          { id: 1, delay_type: "BAGGAGE", minutes: 18, responsibility_team: "行李组", root_cause: "行李分拣系统卡包", reason: "延误事件 #1（BAGGAGE，+18 分钟）尚未关闭" }
        ],
        blocking_bookings: [
          { id: 5, resource_id: 4, resource_code: "BELT-12", task_id: 5, booking_status: "CONFLICT", conflict_reason: "与 CA166 的行李传送带占用重叠", reason: "资源 BELT-12 的预约存在冲突未解决：与 CA166 的行李传送带占用重叠" }
        ]
      }
    },
    {
      id: 3, flight_no: "CZ302", aircraft_reg: "B-8819", stand_no: "109",
      arrival_time: "2026-09-22T09:10:00Z", departure_time: "2026-09-22T10:15:00Z",
      turnaround_status: "IN_SERVICE", delay_reason: "航油车排队", ready_at: null,
      progress: { total_tasks: 2, signed_tasks: 2, open_delays: 1, active_bookings: 1 },
      gate: {
        releasable: false, unsigned_tasks: [],
        open_delays: [{ id: 2, delay_type: "REFUEL", minutes: 12, responsibility_team: "加油组", root_cause: "航油车调度排队", reason: "延误事件 #2（REFUEL，+12 分钟）尚未关闭" }],
        blocking_bookings: []
      }
    }
  ],
  groundTask: [
    { id: 1, turnaround_id: 1, task_type: "CLEANING", team_id: 1, planned_start: "2026-09-22T08:05:00Z", deadline: "2026-09-22T08:30:00Z", actual_finish: "2026-09-22T08:35:00Z", status: "SIGNED", blocker_note: "" },
    { id: 2, turnaround_id: 1, task_type: "CATERING", team_id: 2, planned_start: "2026-09-22T08:08:00Z", deadline: "2026-09-22T08:32:00Z", actual_finish: "2026-09-22T08:35:00Z", status: "SIGNED", blocker_note: "" },
    { id: 3, turnaround_id: 1, task_type: "REFUEL", team_id: 3, planned_start: "2026-09-22T08:10:00Z", deadline: "2026-09-22T08:40:00Z", actual_finish: "2026-09-22T08:35:00Z", status: "SIGNED", blocker_note: "" },
    { id: 4, turnaround_id: 2, task_type: "BAGGAGE", team_id: 4, planned_start: "2026-09-22T08:45:00Z", deadline: "2026-09-22T09:10:00Z", actual_finish: null, status: "IN_PROGRESS", blocker_note: "传送带车未到位" },
    { id: 5, turnaround_id: 2, task_type: "WATER_SERVICE", team_id: 5, planned_start: "2026-09-22T08:50:00Z", deadline: "2026-09-22T09:15:00Z", actual_finish: null, status: "BLOCKED", blocker_note: "清水车故障" },
    { id: 6, turnaround_id: 2, task_type: "PUSHBACK", team_id: 6, planned_start: "2026-09-22T09:20:00Z", deadline: "2026-09-22T09:35:00Z", actual_finish: null, status: "DISPATCHED", blocker_note: "" },
    { id: 7, turnaround_id: 3, task_type: "BAGGAGE", team_id: 4, planned_start: "2026-09-22T09:15:00Z", deadline: "2026-09-22T09:40:00Z", actual_finish: "2026-09-22T09:45:00Z", status: "SIGNED", blocker_note: "" },
    { id: 8, turnaround_id: 3, task_type: "REFUEL", team_id: 3, planned_start: "2026-09-22T09:20:00Z", deadline: "2026-09-22T09:50:00Z", actual_finish: "2026-09-22T09:45:00Z", status: "SIGNED", blocker_note: "" }
  ],
  groundResource: [
    { id: 1, resource_code: "GPU-01", resource_type: "GPU", location: "T2-103", availability_status: "BOOKED", maintenance_due_at: "2026-10-22T08:00:00Z", owner_team: "机务一组" },
    { id: 2, resource_code: "CART-CAT", resource_type: "CATERING", location: "T2-117", availability_status: "BOOKED", maintenance_due_at: "2026-10-22T08:00:00Z", owner_team: "配餐组" },
    { id: 3, resource_code: "FUEL-07", resource_type: "REFUEL", location: "T2-109", availability_status: "BOOKED", maintenance_due_at: "2026-10-22T08:00:00Z", owner_team: "加油组" },
    { id: 4, resource_code: "BELT-12", resource_type: "BAGGAGE", location: "T2-117", availability_status: "BOOKED", maintenance_due_at: "2026-10-22T08:00:00Z", owner_team: "行李组" },
    { id: 5, resource_code: "WATER-03", resource_type: "WATER", location: "T2-103", availability_status: "BOOKED", maintenance_due_at: "2026-10-22T08:00:00Z", owner_team: "清水组" },
    { id: 6, resource_code: "PB-02", resource_type: "PUSHBACK", location: "T2-103", availability_status: "BOOKED", maintenance_due_at: "2026-10-22T08:00:00Z", owner_team: "牵引车组" }
  ],
  resourceBooking: [
    { id: 1, resource_id: 1, resource_code: "GPU-01", turnaround_id: 1, task_id: 1, start_time: "2026-09-22T08:05:00Z", end_time: "2026-09-22T08:40:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null },
    { id: 2, resource_id: 5, resource_code: "WATER-03", turnaround_id: 1, task_id: 2, start_time: "2026-09-22T08:08:00Z", end_time: "2026-09-22T08:42:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null },
    { id: 3, resource_id: 6, resource_code: "PB-02", turnaround_id: 1, task_id: 3, start_time: "2026-09-22T08:10:00Z", end_time: "2026-09-22T08:48:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null },
    { id: 4, resource_id: 2, resource_code: "CART-CAT", turnaround_id: 2, task_id: 4, start_time: "2026-09-22T08:45:00Z", end_time: "2026-09-22T09:20:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null },
    { id: 5, resource_id: 4, resource_code: "BELT-12", turnaround_id: 2, task_id: 5, start_time: "2026-09-22T08:50:00Z", end_time: "2026-09-22T09:25:00Z", booking_status: "CONFLICT", conflict_reason: "与 CA166 的行李传送带占用重叠", released_at: null },
    { id: 6, resource_id: 3, resource_code: "FUEL-07", turnaround_id: 3, task_id: 8, start_time: "2026-09-22T09:20:00Z", end_time: "2026-09-22T09:55:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null }
  ],
  delayEvent: [
    { id: 1, turnaround_id: 2, delay_type: "BAGGAGE", minutes: 18, root_cause: "行李分拣系统卡包", responsibility_team: "行李组", resolved_at: null },
    { id: 2, turnaround_id: 3, delay_type: "REFUEL", minutes: 12, root_cause: "航油车调度排队", responsibility_team: "加油组", resolved_at: null },
    { id: 3, turnaround_id: 1, delay_type: "CATERING", minutes: 5, root_cause: "配餐晚到，已追回", responsibility_team: "配餐组", resolved_at: "2026-09-22T08:35:00Z" }
  ]
};
