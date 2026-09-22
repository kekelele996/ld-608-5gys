-- ground-turn 航空地勤周转保障平台 · MySQL 8.0 初始化脚本
-- 实际部署由后端 GORM AutoMigrate 自动建表，本脚本作为显式 DDL 与评审参考。

CREATE TABLE IF NOT EXISTS flight_turnaround (
  id                 BIGINT AUTO_INCREMENT PRIMARY KEY,
  flight_no          VARCHAR(32)  NOT NULL,
  aircraft_reg       VARCHAR(32)  NOT NULL,
  stand_no           VARCHAR(32),
  arrival_time       DATETIME     NOT NULL,
  departure_time     DATETIME     NOT NULL,
  turnaround_status  VARCHAR(32)  NOT NULL DEFAULT 'ARRIVING',
  delay_reason       VARCHAR(255) DEFAULT '',
  ready_at           DATETIME     NULL,
  created_at         DATETIME(3)  NULL,
  updated_at         DATETIME(3)  NULL,
  INDEX idx_turnaround_status (turnaround_status),
  INDEX idx_turnaround_ready_at (ready_at)
);

CREATE TABLE IF NOT EXISTS ground_task (
  id            BIGINT AUTO_INCREMENT PRIMARY KEY,
  turnaround_id BIGINT      NOT NULL,
  task_type     VARCHAR(32) NOT NULL,
  team_id       BIGINT      NOT NULL,
  planned_start DATETIME    NOT NULL,
  deadline      DATETIME    NOT NULL,
  actual_finish DATETIME    NULL,
  status        VARCHAR(32) NOT NULL DEFAULT 'DISPATCHED',
  blocker_note  VARCHAR(255) DEFAULT '',
  created_at    DATETIME(3) NULL,
  updated_at    DATETIME(3) NULL,
  INDEX idx_task_turnaround (turnaround_id),
  INDEX idx_task_status (status)
);

CREATE TABLE IF NOT EXISTS ground_resource (
  id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
  resource_code       VARCHAR(32) NOT NULL,
  resource_type       VARCHAR(32) NOT NULL,
  location            VARCHAR(64),
  availability_status VARCHAR(32) NOT NULL DEFAULT 'AVAILABLE',
  maintenance_due_at  DATETIME    NULL,
  owner_team          VARCHAR(64),
  created_at          DATETIME(3) NULL,
  updated_at          DATETIME(3) NULL,
  INDEX idx_resource_status (availability_status)
);

CREATE TABLE IF NOT EXISTS resource_booking (
  id              BIGINT AUTO_INCREMENT PRIMARY KEY,
  resource_id     BIGINT      NOT NULL,
  turnaround_id   BIGINT      NOT NULL,
  task_id         BIGINT      NOT NULL,
  start_time      DATETIME    NOT NULL,
  end_time        DATETIME    NOT NULL,
  booking_status  VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  conflict_reason VARCHAR(255) DEFAULT '',
  released_at     DATETIME    NULL,
  created_at      DATETIME(3) NULL,
  updated_at      DATETIME(3) NULL,
  INDEX idx_booking_turnaround (turnaround_id),
  INDEX idx_booking_resource (resource_id),
  INDEX idx_booking_status (booking_status)
);

CREATE TABLE IF NOT EXISTS delay_event (
  id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
  turnaround_id       BIGINT      NOT NULL,
  delay_type          VARCHAR(32) NOT NULL,
  minutes             INT         NOT NULL DEFAULT 0,
  root_cause          VARCHAR(255),
  responsibility_team VARCHAR(64),
  resolved_at         DATETIME    NULL,
  created_at          DATETIME(3) NULL,
  updated_at          DATETIME(3) NULL,
  INDEX idx_delay_turnaround (turnaround_id),
  INDEX idx_delay_resolved (resolved_at)
);

CREATE TABLE IF NOT EXISTS audit_log (
  id          BIGINT AUTO_INCREMENT PRIMARY KEY,
  actor       VARCHAR(64),
  action      VARCHAR(64),
  target_type VARCHAR(64),
  target_id   VARCHAR(64),
  detail      VARCHAR(512),
  created_at  DATETIME(3) NULL,
  INDEX idx_audit_target (target_type, target_id),
  INDEX idx_audit_created (created_at)
);
