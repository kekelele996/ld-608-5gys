export const GroundTaskType = ["CLEANING","CATERING","BAGGAGE","REFUEL","WATER_SERVICE","PUSHBACK"] as const;
export type GroundTaskType = (typeof GroundTaskType)[number];
export const GroundTaskTypeText: Record<GroundTaskType, string> = {
  CLEANING: "客舱清洁",
  CATERING: "配餐保障",
  BAGGAGE: "行李装卸",
  REFUEL: "航油加注",
  WATER_SERVICE: "清水排污",
  PUSHBACK: "牵引车推出"
};
