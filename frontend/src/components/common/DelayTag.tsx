// DelayTag 延误状态标签：OPEN 未关闭（阻塞放行）/ RESOLVED 已关闭。
export function DelayTag({ value }: { value: "OPEN" | "RESOLVED" | string }) {
  const text = value === "OPEN" ? "未关闭" : value === "RESOLVED" ? "已关闭" : value;
  return <span className={"badge delay-" + String(value).toLowerCase()}>{text}</span>;
}
