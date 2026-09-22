import { useFlightTurnaroundStore } from "../../stores/FlightTurnaroundStore";

type Props = {
  flightId: number;
  status: string;
  ready: boolean;
};

export function ReleaseButton({ flightId, status, ready }: Props) {
  const releasingId = useFlightTurnaroundStore((state) => state.releasingId);
  const release = useFlightTurnaroundStore((state) => state.release);
  const released = status === "READY" || status === "DEPARTED";
  const loading = releasingId === flightId;

  return <button
    className="release-button"
    disabled={loading || released || !ready}
    onClick={() => void release(flightId)}
    title={released ? "已经放行，重复提交只会返回冲突结果" : ready ? "校验通过后进入 READY 并释放 ACTIVE 预约" : "请先清空任务、延误和预约阻塞"}
  >
    {released ? "已放行" : loading ? "放行中…" : ready ? "放行联动" : "不可放行"}
  </button>;
}
