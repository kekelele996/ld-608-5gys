CREATE TABLE IF NOT EXISTS flight_turnaround (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  flight_no VARCHAR(32) NOT NULL,
  aircraft_reg VARCHAR(32) NOT NULL,
  stand_no VARCHAR(32) NOT NULL,
  arrival_time DATETIME NOT NULL,
  departure_time DATETIME NOT NULL,
  turnaround_status VARCHAR(32) NOT NULL,
  delay_reason VARCHAR(255) NOT NULL DEFAULT '',
  ready_at DATETIME NULL,
  timeline_synced_at DATETIME NULL,
  release_in_progress TINYINT(1) NOT NULL DEFAULT 0,
  INDEX idx_flight_status (turnaround_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ground_resource (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  resource_code VARCHAR(32) NOT NULL UNIQUE,
  resource_type VARCHAR(32) NOT NULL,
  location VARCHAR(32) NOT NULL,
  availability_status VARCHAR(32) NOT NULL,
  maintenance_due_at DATETIME NULL,
  owner_team VARCHAR(64) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ground_task (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  turnaround_id BIGINT NOT NULL,
  task_type VARCHAR(32) NOT NULL,
  team_id BIGINT NOT NULL,
  planned_start DATETIME NOT NULL,
  deadline DATETIME NOT NULL,
  actual_finish DATETIME NULL,
  status VARCHAR(32) NOT NULL,
  blocker_note VARCHAR(255) NOT NULL DEFAULT '',
  timeline_synced_at DATETIME NULL,
  CONSTRAINT fk_ground_task_turnaround FOREIGN KEY (turnaround_id) REFERENCES flight_turnaround(id),
  INDEX idx_task_turnaround_status (turnaround_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS resource_booking (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  resource_id BIGINT NOT NULL,
  turnaround_id BIGINT NOT NULL,
  task_id BIGINT NULL,
  start_time DATETIME NOT NULL,
  end_time DATETIME NOT NULL,
  booking_status VARCHAR(32) NOT NULL,
  conflict_reason VARCHAR(255) NOT NULL DEFAULT '',
  released_at DATETIME NULL,
  CONSTRAINT fk_booking_resource FOREIGN KEY (resource_id) REFERENCES ground_resource(id),
  CONSTRAINT fk_booking_turnaround FOREIGN KEY (turnaround_id) REFERENCES flight_turnaround(id),
  CONSTRAINT fk_booking_task FOREIGN KEY (task_id) REFERENCES ground_task(id),
  INDEX idx_booking_turnaround_status (turnaround_id, booking_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS delay_event (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  turnaround_id BIGINT NOT NULL,
  delay_type VARCHAR(64) NOT NULL,
  minutes INT NOT NULL DEFAULT 0,
  root_cause VARCHAR(255) NOT NULL DEFAULT '',
  responsibility_team VARCHAR(64) NOT NULL DEFAULT '',
  resolved_at DATETIME NULL,
  CONSTRAINT fk_delay_turnaround FOREIGN KEY (turnaround_id) REFERENCES flight_turnaround(id),
  INDEX idx_delay_turnaround_resolved (turnaround_id, resolved_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS audit_log (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  actor VARCHAR(64) NOT NULL,
  action VARCHAR(64) NOT NULL,
  target_type VARCHAR(64) NOT NULL,
  target_id VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL,
  INDEX idx_audit_target (target_type, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO flight_turnaround (id, flight_no, aircraft_reg, stand_no, arrival_time, departure_time, turnaround_status, delay_reason)
VALUES
  (1, 'CA1801', 'B-2026', 'A12', '2026-09-22 08:00:00', '2026-09-22 09:10:00', 'IN_SERVICE', ''),
  (2, 'MU5102', 'B-5088', 'B07', '2026-09-22 10:00:00', '2026-09-22 11:10:00', 'IN_SERVICE', '行李等待'),
  (3, 'CZ3120', 'B-6310', 'C03', '2026-09-22 12:00:00', '2026-09-22 13:10:00', 'ON_STAND', '');

INSERT IGNORE INTO ground_resource (id, resource_code, resource_type, location, availability_status, maintenance_due_at, owner_team)
VALUES
  (1, 'CLN-A12', 'CLEANING', 'A12', 'BOOKED', '2026-10-01 00:00:00', '清洁一组'),
  (2, 'BAG-07', 'BAGGAGE', 'B07', 'BOOKED', '2026-10-02 00:00:00', '行李班组'),
  (3, 'FUEL-B07', 'REFUEL', 'B07', 'BOOKED', '2026-10-03 00:00:00', '加油班组'),
  (4, 'CAT-C03', 'CATERING', 'C03', 'BOOKED', '2026-10-04 00:00:00', '配餐班组');

INSERT IGNORE INTO ground_task (id, turnaround_id, task_type, team_id, planned_start, deadline, actual_finish, status, blocker_note)
VALUES
  (101, 1, 'CLEANING', 1, '2026-09-22 08:05:00', '2026-09-22 08:40:00', '2026-09-22 09:00:00', 'SIGNED', ''),
  (102, 1, 'BAGGAGE', 2, '2026-09-22 08:05:00', '2026-09-22 08:50:00', '2026-09-22 09:00:00', 'SIGNED', ''),
  (201, 2, 'CLEANING', 1, '2026-09-22 10:05:00', '2026-09-22 10:40:00', NULL, 'SIGNED', ''),
  (202, 2, 'REFUEL', 3, '2026-09-22 10:10:00', '2026-09-22 10:50:00', NULL, 'BLOCKED', '加油栓井占用'),
  (301, 3, 'CATERING', 4, '2026-09-22 12:05:00', '2026-09-22 12:40:00', NULL, 'SIGNED', '');

INSERT IGNORE INTO resource_booking (id, resource_id, turnaround_id, task_id, start_time, end_time, booking_status, conflict_reason)
VALUES
  (1001, 1, 1, 101, '2026-09-22 08:05:00', '2026-09-22 08:40:00', 'ACTIVE', ''),
  (1002, 2, 1, 102, '2026-09-22 08:05:00', '2026-09-22 08:50:00', 'ACTIVE', ''),
  (2001, 3, 2, 202, '2026-09-22 10:10:00', '2026-09-22 10:50:00', 'CONFLICT', '加油栓井与邻机位冲突'),
  (3001, 4, 3, 301, '2026-09-22 12:05:00', '2026-09-22 12:40:00', 'ACTIVE', '');

INSERT IGNORE INTO delay_event (id, turnaround_id, delay_type, minutes, root_cause, responsibility_team, resolved_at)
VALUES (1, 2, 'GROUND_HANDLING', 25, '加油资源冲突', '加油班组', NULL);
