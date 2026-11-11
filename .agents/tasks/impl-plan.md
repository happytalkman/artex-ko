# Implementation Plan — 17 New Features for artex-ko

> Generated after reading: similar-projects-report.md, db/schema.sql, agent/finding_recorder.go,
> agent/prompt.go, server/server.go, server/finding_retests.go, report/report.go,
> llmrec/llmrec.go, skills/api-recon/SKILL.md, agent/retester.go, agent/finding_workflow.go,
> db/findings.go, db/llm_usage.go, server/scheduler.go, agent/promptcatalog.go, go.mod.

---

## 아키텍처 결정 사항

| 결정 | 근거 |
|------|------|
| CVSS 계산기를 외부 라이브러리 없이 순수 Go로 구현 | go.mod 의존성 최소화. CVSS 3.1 공식은 500줄 이내 |
| PoC 게이트를 server/finding_workflow.go에서 처리 | 서버 레이어에서 한국어 에러 메시지 반환 용이 |
| 연속 모니터링을 기존 scheduler.go에 통합 | 별도 서비스 불필요; scheduler.step()에 추가 |
| task_type 컬럼으로 CTF/pentest 구분 | 기존 status/queue_mode 패턴과 일관성 |
| pgvector 벡터를 `'[...]'::vector` 텍스트로 직접 삽입 | pgvector-go 라이브러리 불필요 |
| PDF 변환에 wkhtmltopdf 사용 | 없으면 501 반환 |
| Sigma 규칙을 파일시스템에 저장 | 기존 detections/sigma/ 활용 |
| 벤치마크 YAML 파싱에 gopkg.in/yaml.v3 추가 | Go 생태계 표준 |

---

## 기능 간 의존성

```
DB 스키마 변경 (schema.sql)
  ├── FEAT-001: findings.cvss_vector/cvss_score → AddFinding 시그니처
  ├── FEAT-002: tasks.schedule/task_type → scheduler 로직
  └── FEAT-004: resources 테이블, embeddings 테이블

agent/finding_workflow.go
  ├── FEAT-001: report_finding 도구에 cvss_vector 추가
  └── FEAT-004: 중복 감지 로직 주입

서버 라우트 (server/server.go)
  ├── FEAT-001: GET /api/llm/stats
  ├── FEAT-004: /api/resources CRUD, GET /api/tasks/{id}/report
  └── FEAT-005: GET /api/tasks/{id}/sarif, GET /api/benchmarks/targets

스킬 파일 (skills/)
  └── FEAT-003: 완전 독립 (Go 코드 변경 없음)
```

---

## 구현 계획

- [ ] 1. **[FEAT-001] DB: findings에 CVSS 컬럼 추가**
  - `db/schema.sql`: `ALTER TABLE findings ADD COLUMN IF NOT EXISTS cvss_vector TEXT NOT NULL DEFAULT '';` 및 `cvss_score NUMERIC(4,1) NOT NULL DEFAULT 0;` 추가. `idx_llm_usage_worker` 인덱스 추가.
  - Files: `db/schema.sql`
  - Verify: `go build ./...`

- [ ] 2. **[FEAT-001] CVSS v3.1 계산기**
  - `db/cvss.go` (new): `CalcCVSS31Score(vector string) (float64, error)`. `db/cvss_test.go`: AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H → 9.8.
  - Files: `db/cvss.go`, `db/cvss_test.go`
  - Verify: `go test ./db/... -run TestCalcCVSS31`

- [ ] 3. **[FEAT-001] DBFinding + AddFinding에 CVSS 필드 추가**
  - `db/findings.go`: DBFinding에 CvssVector/CvssScore 추가. `findingSelectCols`에 컬럼 추가. `scanFindings` 업데이트. `AddFinding` 시그니처 확장. `RecordFindingInput` 구조체도 동일 필드 추가.
  - Files: `db/findings.go`
  - Verify: `go build ./...`

- [ ] 4. **[FEAT-001] report_finding 도구 스키마에 CVSS + PoC 게이트**
  - `agent/toolcatalog.go`: report_finding 스키마에 `cvss_vector`, `cvss_score` 옵셔널 필드 추가.
  - `agent/finding_workflow.go`: `findingWorkflowTools()`에 `findingCVSSGuidance` 상수 주입.
  - `server/finding_workflow.go`: CalcCVSS31Score 호출 + PoC 게이트 (traffic_refs==[] && evidence=="" → severity='info').
  - Files: `agent/toolcatalog.go`, `agent/finding_workflow.go`, `server/finding_workflow.go`
  - Verify: `go test ./agent/... ./server/...`

- [ ] 5. **[FEAT-001] report/report.go에 CVSS 표시**
  - `findingView`에 CvssVector/CvssScore. `Markdown()` 헤더 형식 `[SEVERITY / CVSS: X.X]`.
  - Files: `report/report.go`
  - Verify: `go test ./report/...`

- [ ] 6. **[FEAT-001] LLM 통계 엔드포인트 /api/llm/stats**
  - `db/llm_usage.go`: `WorkerLatencyStats()` 추가 (avg/p95 레이턴시, 에러율).
  - `server/llm_stats.go` (new): `getLLMStats` 핸들러.
  - `server/server.go`: `mux.HandleFunc("GET /api/llm/stats", s.getLLMStats)` 등록.
  - Files: `db/llm_usage.go`, `server/llm_stats.go`, `server/server.go`
  - Verify: `go build ./...`

- [ ] 7. **[FEAT-002] DB: tasks에 schedule/task_type 컬럼 추가**
  - `db/schema.sql`: `schedule TEXT NOT NULL DEFAULT ''`, `schedule_last_run_at TIMESTAMPTZ`, `task_type TEXT NOT NULL DEFAULT 'pentest'`.
  - Files: `db/schema.sql`
  - Verify: `go build ./...`

- [ ] 8. **[FEAT-002] db/tasks.go + cron 파서**
  - Task 구조체에 Schedule/ScheduleLastRunAt/TaskType 추가. `ListScheduledTasks()`, `SetTaskScheduleLastRun()` 구현.
  - `server/scheduler.go`: `parseCronNextFire(expr string, since time.Time) (time.Time, error)` 헬퍼.
  - Files: `db/tasks.go`, `server/scheduler.go`
  - Verify: `go build ./...`

- [ ] 9. **[FEAT-002] 연속 모니터링: scheduler.go에 fireScheduledTasks**
  - `scheduler.step()`에 `sc.fireScheduledTasks()` 추가. 기한 도래 시 SetTaskScheduleLastRun + engine start. 신규 취약점 watermark 비교 후 notify 연동.
  - Files: `server/scheduler.go`
  - Verify: `go build ./...`

- [ ] 10. **[FEAT-002] task metadata API에 schedule/task_type 반영**
  - `server/task_metadata.go`: updateTaskMetadata에서 schedule/task_type 파싱·저장. 한국어 에러 상수 추가.
  - Files: `server/task_metadata.go`, `server/server.go`
  - Verify: `go build ./...`

- [ ] 11. **[FEAT-002] CTF 모드: 플래그 감지**
  - `agent/goals.go`: `detectCTFFlag(text string) (flag string, found bool)` — 정규식 `(?i)(flag|htb|ctf|picoctf|root)\{[^}]+\}`.
  - `agent/worker.go`: task_type=='ctf' 조건 하에 감지 → report_finding 자동 호출.
  - `agent/promptcatalog.go`: `ctfModeGuidance` 상수.
  - Files: `agent/goals.go`, `agent/worker.go`, `agent/promptcatalog.go`
  - Verify: `go test ./agent/... -run TestDetectCTFFlag`

- [ ] 12. **[FEAT-002] 복테스트 fixed 시 findings.status 자동 변경 보완**
  - `db/finding_retests.go`: RecordFindingRetestResult에서 verdict='fixed' 시 동일 트랜잭션 내 `UPDATE findings SET status='fixed'` 확인·추가.
  - Files: `db/finding_retests.go`
  - Verify: `go test ./db/... -run TestFindingRetestFixed`

- [ ] 13. **[FEAT-003] 스킬: Linux LPE**
  - `skills/linux-lpe/SKILL.md` (frontmatter + SUID/sudo/cron/PATH/kernel 체크리스트).
  - `skills/linux-lpe/scripts/linpeas_parse.py` (LinPEAS 출력 파싱, JSON 출력).
  - Files: `skills/linux-lpe/SKILL.md`, `skills/linux-lpe/scripts/linpeas_parse.py`
  - Verify: `python3 -m py_compile skills/linux-lpe/scripts/linpeas_parse.py`

- [ ] 14. **[FEAT-003] 스킬: Whitebox Analysis**
  - `skills/whitebox-analysis/SKILL.md`. `skills/whitebox-analysis/scripts/ast_scan.py` (Python ast, 위험 싱크 탐지).
  - Files: `skills/whitebox-analysis/SKILL.md`, `skills/whitebox-analysis/scripts/ast_scan.py`
  - Verify: `python3 -m py_compile skills/whitebox-analysis/scripts/ast_scan.py`

- [ ] 15. **[FEAT-003] 스킬: GraphQL Audit + WebSocket Fuzz**
  - `skills/graphql-audit/SKILL.md` (인트로스펙션→퍼징→권한 우회).
  - `skills/websocket-fuzz/SKILL.md` (핸드셰이크→메시지 열거→인증 우회).
  - Files: `skills/graphql-audit/SKILL.md`, `skills/websocket-fuzz/SKILL.md`
  - Verify: Both files have valid frontmatter

- [ ] 16. **[FEAT-003] 스킬: Cloud Audit + TOTP**
  - `skills/cloud-audit/SKILL.md` (AWS/GCP/Azure). `skills/cloud-audit/scripts/aws_check.sh`.
  - `skills/playwright-cli/references/totp-login.md`. `skills/playwright-cli/scripts/totp.py` (RFC 6238, stdlib).
  - Files: 위 4개 파일
  - Verify: `python3 skills/playwright-cli/scripts/totp.py JBSWY3DPEHPK3PXP` → 6자리

- [ ] 17. **[FEAT-004] DB: resources 테이블 추가**
  - `db/schema.sql`: resources 테이블 DDL.
  - Files: `db/schema.sql`
  - Verify: `go build ./...`

- [ ] 18. **[FEAT-004] 사용자 자원 라이브러리 — Go 구현**
  - `db/resources.go` (new): Resource 구조체 + CRUD.
  - `server/resources.go` (new): listResources/uploadResource/downloadResource/deleteResource. `data/resources/` 자동 생성. 최대 50MB.
  - `server/server.go`: 5개 라우트 등록.
  - Files: `db/resources.go`, `server/resources.go`, `server/server.go`
  - Verify: `go build ./...`

- [ ] 19. **[FEAT-004] 에이전트 도구: list_resources + read_resource**
  - `agent/toolcatalog.go` 또는 server SeedTool 패턴으로 등록. list_resources: DB 조회. read_resource: 파일 읽기 (1MB 상한, 바이너리 base64).
  - Files: `agent/toolcatalog.go` (or server/resources.go)
  - Verify: `go build ./...`

- [ ] 20. **[FEAT-004] 취약점 중복 감지 강화**
  - `db/findings.go`: `FindDuplicateFinding(taskID int64, vulnclass, urlKeyword string)` 추가.
  - `server/finding_workflow.go`: RecordFinding 전 FindDuplicateFinding 호출. 중복 시 `{duplicate:true, ...}` 반환.
  - `agent/finding_workflow.go`: findingWorkflowTools()에 중복 처리 가이드 추가.
  - Files: `db/findings.go`, `server/finding_workflow.go`, `agent/finding_workflow.go`
  - Verify: `go test ./db/... ./server/...`

- [ ] 21. **[FEAT-004] DB: embeddings 테이블 추가 (pgvector)**
  - `db/schema.sql`: `CREATE EXTENSION IF NOT EXISTS vector;` + embeddings 테이블 (task_id, content_type, content_id, text_summary, `embedding vector(1536)`). IVFFlat 인덱스.
  - 주의: pgvector 미설치 시 조용히 실패 — `CREATE EXTENSION IF NOT EXISTS` 패턴이 idempotent.
  - Files: `db/schema.sql`
  - Verify: `go build ./...`

- [ ] 22. **[FEAT-004] memory 패키지 구현**
  - `memory/memory.go` (new, package memory): Store 구조체, StoreExperience(), SearchSimilar(). 임베딩 API: `POST <base_url>/embeddings`. 벡터 저장: `INSERT ... $1::vector` (텍스트 `[0.1,...]` 형식, pgvector-go 불필요).
  - Files: `memory/memory.go`
  - Verify: `go build ./memory/...`

- [ ] 23. **[FEAT-004] HTML/PDF 보고서 export**
  - `report/html.go` (new): `HTMLReport(in Input) string` — Go html/template.
  - `server/task_report.go` (new): `taskReport` — format=md/html/pdf/sarif. pdf: wkhtmltopdf; 없으면 501.
  - `server/server.go`: `mux.HandleFunc("GET /api/tasks/{id}/report", s.taskReport)`.
  - Files: `report/html.go`, `server/task_report.go`, `server/server.go`
  - Verify: `go build ./...`; `go test ./report/... -run TestHTMLReport`

- [ ] 24. **[FEAT-004] Sigma 규칙 자동 생성**
  - `report/sigma.go` (new): `GenerateSigmaRule(finding *db.DBFinding) string` — text/template YAML.
  - `server/finding_workflow.go`: report_finding 성공 후 GenerateSigmaRule → `detections/sigma/<vulnclass>_<finding_id>.yml` 저장.
  - Files: `report/sigma.go`, `server/finding_workflow.go`
  - Verify: `go build ./...`; generated YAML valid

- [ ] 25. **[FEAT-005] SARIF 보고서 export**
  - `report/sarif.go` (new): `SARIF(in Input) ([]byte, error)` — SARIF 2.1.0. 각 finding → result (ruleId, level, message, locations).
  - `server/task_report.go`: format=sarif 처리 추가. `server/server.go`: `GET /api/tasks/{id}/sarif`.
  - `report/sarif_test.go`: TestSARIF validates schema/version.
  - Files: `report/sarif.go`, `report/sarif_test.go`, `server/task_report.go`, `server/server.go`
  - Verify: `go test ./report/... -run TestSARIF`

- [ ] 26. **[FEAT-005] GitHub Actions 액션 파일**
  - `.github/actions/artex-scan/action.yml` (composite action: artex_url/api_key/task_id/wait_seconds inputs).
  - `.github/actions/artex-scan/README.md`. `docs/github-actions-guide.md` (한국어).
  - Files: 위 3개 파일
  - Verify: `go build ./...` still passes; YAML valid

- [ ] 27. **[FEAT-005] 벤치마크 타겟 정의**
  - `benchmarks/targets/dvwa.yaml` (SQLi/XSS/FileUpload/CSRF/CommandInjection expected_findings).
  - `benchmarks/targets/juice-shop.yaml`.
  - Files: `benchmarks/targets/dvwa.yaml`, `benchmarks/targets/juice-shop.yaml`
  - Verify: YAML valid

- [ ] 28. **[FEAT-005] 벤치마크 CLI 러너 + go.mod 업데이트**
  - `benchmarks/run.go` (package main): --server/--token/--target/--timeout. 작업 생성→폴링→recall 계산. stdlib 전용.
  - `go.mod`: `gopkg.in/yaml.v3 v3.0.1` 추가. `go mod tidy` 실행.
  - Files: `benchmarks/run.go`, `go.mod`, `go.sum`
  - Verify: `go build ./benchmarks/...`

- [ ] 29. **[FEAT-005] 벤치마크 API 엔드포인트**
  - `server/benchmarks.go` (new): `listBenchmarkTargets` — `benchmarks/targets/*.yaml` 읽기, JSON 반환.
  - `server/server.go`: `mux.HandleFunc("GET /api/benchmarks/targets", s.listBenchmarkTargets)`.
  - Files: `server/benchmarks.go`, `server/server.go`
  - Verify: `go build ./...`

---

## 주의 사항

### DB 패턴
- 모든 DB 마이그레이션은 `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` 패턴 (기존 schema.sql 준수)
- `scanFindings()`은 `findingSelectCols` 상수와 정확히 동기화 — 컬럼 추가 시 두 곳 모두 수정

### 에러 메시지
- `writeErr()` 반환 HTTP 에러는 모두 한국어 (server/server.go의 const 블록 패턴)
- 에이전트 두뇌(시스템 프롬프트) 중국어는 보존

### 도구 등록 패턴
- 에이전트 도구: `agent/toolcatalog.go`의 BuiltinToolSeeds() 패턴
- 서버 관리 도구: `server/finding_retests.go`의 `seedFindingRetester()` 패턴 (SeedTool)

### skill 파일 형식
- 반드시 frontmatter `---\nname: <name>\ndescription: <desc>\n---` 포함
- `skills/api-recon/SKILL.md` 형식 참조

### memory 패키지
- pgvector 미설치 환경 고려 — StoreExperience/SearchSimilar 실패 시 로그 후 주 흐름 계속
- 벡터 텍스트 형식: `fmt.Sprintf("[%s]", strings.Join(floatStrs, ",")) + "::vector"`

### FEAT 실행 순서 준수
- FEAT-001 → FEAT-002 → FEAT-003 → FEAT-004 → FEAT-005 순으로 실행
- FEAT-003은 완전 독립 (Go 코드 변경 없음) — 다른 FEAT와 병행 가능하나 순서 유지
