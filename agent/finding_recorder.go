package agent

import (
	"context"

	"github.com/Autumn-27/artex/db"
)

// FindingRecorder is injected by the host; agents never synthesize or copy
// evidence bodies themselves. Its implementation owns the atomic write.
type FindingRecorder interface {
	Record(context.Context, db.RecordFindingInput, []db.TrafficRef) (*db.RecordedFinding, error)
}

// Tool-use guidance is appended without replacing the user's editable prompt.
// It does not require capture or claim that unavailable traffic tools exist.
const findingTrafficGuidance = "\n\n**漏洞流量证据（可选）**：调用 report_finding 上报漏洞时，如有已查看并确认支持漏洞结论的 HTTP 请求/响应，可用 traffic_refs 按复现顺序绑定真实 ID；域名和时间只作候选筛选，不推定关联。TCP 等非 HTTP 漏洞、未采集或无确切匹配时省略或传 []，在 evidence 保留命令输出、日志等其他可验证证据，建议说明未绑定原因。不要猜测 ID，也不要仅为补包重复探测。"

// findingCVSSGuidance instructs the worker to determine and supply CVSS v3.1
// base vectors for each finding via the report_finding tool.
// PoC gate: a finding with no traffic_refs and no evidence is downgraded to 'info'
// severity by the recording layer, since it cannot be confirmed without proof.
const findingCVSSGuidance = "\n\n**CVSS v3.1 자동 채점 및 PoC 게이트**：每条漏洞上报时，请在 cvss_vector 字段中提供 CVSS v3.1 基础评分向量（格式：AV:_/AC:_/PR:_/UI:_/S:_/C:_/I:_/A:_，例如 AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H）。评分会自动从向量计算（也可在 cvss_score 手动覆盖）。若无 traffic_refs 且 evidence 为空，说明无可验证证据，系统会将严重级别自动降为 info（PoC 门控）——请务必在 evidence 保留命令输出、日志或 PoC 片段，或绑定真实流量 ID。"

func (t *ToolSet) SetFindingRecorder(r FindingRecorder)   { t.findingRecorder = r }
func (w *Worker) SetFindingRecorder(r FindingRecorder)    { w.findingRecorder = r }
func (p *Planner) SetFindingRecorder(r FindingRecorder)   { p.findingRecorder = r }
func (m *MainAgent) SetFindingRecorder(r FindingRecorder) { m.findingRecorder = r }
