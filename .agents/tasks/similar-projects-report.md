# ARTEX-KO 유사 프로젝트 조사 보고서

> 작성일: 2025년 조사  
> 목적: GitHub 상위 AI 자율 침투 테스트 오픈소스 프로젝트를 조사하고, artex-ko에 추가하면 좋을 기능을 우선순위별로 추천한다.

---

## 1. 요약

GitHub에서 AI 자율 침투 테스트 분야의 상위 프로젝트들(PentAGI ★27K+, PentestGPT ★8K+, HackingBuddyGPT, Shannon, CyberStrike 등)과 artex-ko를 비교한 결과, artex-ko는 이미 이중 그래프 구조·MITM 프록시·skill 시스템·MCP 지원 등 강력한 차별점을 보유하고 있다. 그러나 경쟁 프로젝트들은 **벡터 메모리/지식 그래프**, **화이트박스 소스코드 분석**, **CVSS 자동 채점·PoC 검증 게이트**, **LLM 옵저버빌리티**, **CI/CD 통합** 등 여러 기능을 갖추고 있어 artex-ko에 추가하면 큰 임팩트를 기대할 수 있다.

---

## 2. GitHub 상위 유사 프로젝트 목록

| 프로젝트 | Stars | 언어 | 핵심 차별점 |
|---|---|---|---|
| [vxcontrol/pentagi](https://github.com/vxcontrol/pentagi) | ★27K+ | Go+React | pgvector 장기 메모리, Neo4j 지식 그래프(Graphiti), Langfuse 옵저버빌리티 |
| [GreyDGL/PentestGPT](https://github.com/GreyDGL/PentestGPT) | ★8K+ | Python | 자율 CTF 파이프라인(PTT), Claude Code/Codex 백엔드 플러그인, 2FA/TOTP 자동 처리 |
| [ipa-lab/hackingBuddyGPT](https://github.com/ipa-lab/hackingBuddyGPT) | ★3K+ | Python | Linux 권한 상승 특화, 재현 가능 벤치마크 VM 세트 |
| [KeygraphHQ/shannon](https://github.com/KeygraphHQ/shannon) | ★2.5K+ | TypeScript | 화이트박스 소스코드 분석+실제 익스플로잇 증명, XBOW 96.15%, GitHub Actions 통합 |
| [CyberStrikeus/CyberStrike](https://github.com/CyberStrikeus/CyberStrike) | ★2K+ | Rust | 7,600+ 보안 스킬, 120+ OWASP 테스트 케이스, 13개 특화 에이전트 |
| [Usta0x001/Phantom](https://github.com/Usta0x001/Phantom) | 신규 | Python | ReAct 루프, 30+ 보안 도구, Docker 샌드박스, 완전 자율 보고서 |
| [FORGE (arXiv)](https://arxiv.org/pdf/2606.03453) | - | Python | CVE 익스플로잇 67.8%, Sigma/Snort 탐지 규칙 자동 생성, CVSS 4단계 등급 |
| [lordx64/pentestkit](https://github.com/lordx64/pentestkit) | 신규 | Python | PoC 필수 게이트, CVSS v3 자동 채점, XBOW 100% |
| [protectai/vulnhuntr](https://github.com/protectai/vulnhuntr) | ★2K+ | Python | 소스코드 콜 체인 추적, 제로샷 취약점 발견 |
| [Armur-Ai/Pentest-Swarm-AI](https://github.com/Armur-Ai/Pentest-Swarm-AI) | 신규 | Go | 스웜 에이전트, 버그바운티·CTF·연속 모니터링 모드 |

*Content was rephrased for compliance with licensing restrictions*

---

## 3. artex-ko 현재 구현 기능 범위

`README.md`, `server/server.go`, `agent/*.go`, `skills/` 디렉터리를 직접 확인한 결과:

### 확인된 기능

**아키텍처·핵심**
- 이중 그래프(자산 그래프 + 탐색 그래프) + 앵커 연결 구조
- 플래너-워커 멀티에이전트 루프 (이벤트 구동 폐곡선)
- 작업 간 worker trace 공유(`search_all_worker_traces`, `get_worker_trace`)
- planner 다중 라운드 공유 todolist (직렬 공격 체인 보존)
- `kill_work` / `steer_work` — 실행 중 의도 중단·방향 수정

**탐지·리포팅**
- finding recorder (`agent/finding_recorder.go`) — HTTP 트래픽 refs 연동
- 취약점 복테스트 에이전트 (`agent/retester.go`) — reproduced/fixed/inconclusive 판정
- 보고서 생성 (`report/` 패키지)
- CSV 내보내기 (UI)

**인프라·통합**
- MITM 기록 프록시 (`traffic/`) — HTTPS CA 포함
- MCP(Model Context Protocol) 지원 (`mcphttp/`)
- skill 시스템 디렉터리 로딩 (`skills/`)
  - `api-recon` skill: SPA/MPA API 열거, 정적+동적 분석
  - `playwright-cli` skill: 브라우저 자동화
  - `scopesentry` skill: ScopeSentry MCP 연동
- LLM 풀/레코더/페일오버 (`llmpool/`, `llmrec/`)
- 알림 시스템 (`notify/`)
- 도구 승인 게이트 (`guard/`, `intercept/`)
- 비동기 보강 (`enrich/`) — DNS/HTTP enrichment

**UX**
- 다국어 출력(한국어·영어·중국어·스페인어)
- SSE 실시간 스트리밍
- 사람 개입 대화(main agent)
- 커버리지 그래프 시각화

### artex-ko에 없는 기능 (경쟁사 대비)

1. 벡터 메모리 / 임베딩 검색 (PentAGI: pgvector)
2. 지식 그래프 (PentAGI: Neo4j+Graphiti)
3. LLM 옵저버빌리티 대시보드 (PentAGI: Langfuse/Grafana)
4. 화이트박스 소스코드 분석 (Shannon)
5. CVSS 자동 채점 (pentestkit, FORGE, ai-pentest-agent)
6. PoC 필수 검증 게이트 (pentestkit)
7. CVE/NVD 연동 (FORGE)
8. 자동 탐지 규칙 생성 (FORGE: Sigma/Snort)
9. CI/CD 통합 (Shannon: GitHub Actions)
10. Linux 권한 상승 특화 에이전트 (hackingBuddyGPT)
11. CTF 자율 모드 (PentestGPT)
12. 연속(24/7) 모니터링 모드
13. 스웜(병렬 분기) 에이전트 (Pentest-Swarm-AI)
14. 사용자 자원 라이브러리 (PentAGI: 파일 업로드/가상 경로 파일시스템)
15. 캠페인(복수 작업) 브랜치 관리
16. GraphQL API (PentAGI)
17. 2FA/TOTP 자동 처리 (PentestGPT)
18. 재현 가능한 벤치마크 평가 모드 (hackingBuddyGPT, XBOW)
19. 다중 도구 Docker 샌드박스 격리 (PentAGI: 200+ Kali 도구)
20. 웹 인텔리전스 내장 브라우저 에이전트 (PentAGI)

---

## 4. 추가 추천 기능 (우선순위별)

### 🔴 P1 — 높은 임팩트 / 구현 가능 (즉시 추천)

#### 4-1. CVSS v3.1 자동 채점 + PoC 검증 게이트

**근거**: pentestkit(XBOW 100%), ai-pentest-agent, FORGE 모두 CVSS 자동 채점을 핵심으로 삼는다. artex-ko는 severity를 에이전트가 텍스트로 판단하는데, 구조화된 CVSS 벡터 문자열을 자동 계산하면 보고서 품질이 대폭 향상된다.

**구현 방향**:
- `report_finding` 도구에 CVSS 벡터 필드 추가 (AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H 형태)
- Worker 프롬프트에 CVSS 채점 가이드 추가 (`agent/prompt.go`)
- DB 스키마에 `cvss_vector`, `cvss_score` 컬럼 추가
- PoC 증거(`traffic_refs` 또는 `evidence`)가 없으면 severity를 자동으로 `info`로 강등하는 게이트 로직

**예상 파일**: `agent/finding_recorder.go`, `db/`, `agent/prompt.go`, `report/`

---

#### 4-2. LLM 옵저버빌리티 — 비용·레이턴시·에러율 대시보드

**근거**: PentAGI는 Langfuse + Grafana 스택을 별도 compose 파일로 제공한다. artex-ko는 `llmrec/` 패키지로 LLM 호출을 이미 기록하고 있지만, UI에서 토큰 소비·레이턴시·에러율을 시각화하는 뷰가 없다. 프로 사용자와 팀 환경에서 비용 통제를 위해 필수적이다.

**구현 방향**:
- `llmrec/`의 기록 데이터를 집계해 `/api/llm/stats` 엔드포인트 추가
- 대시보드 카드에 "에이전트별 토큰·비용·레이턴시·실패율" 추가
- (선택) Langfuse OpenTelemetry SDK를 선택적 통합으로 추가

**예상 파일**: `server/server.go`, `llmrec/`, `web/` 대시보드 컴포넌트

---

#### 4-3. Linux 권한 상승(LPE) 전용 skill

**근거**: hackingBuddyGPT는 Linux LPE만 전문으로 연구하며 재현 가능한 벤치마크 VM 세트를 운영한다. artex-ko의 skill 시스템은 디렉터리 추가만으로 확장 가능하고, 현재 `skills/`에는 웹 중심 skill만 있다.

**구현 방향**:
- `skills/linux-lpe/SKILL.md` 신규 작성 — SUID 바이너리, sudo 권한, cron job, PATH 하이재킹, 커널 익스플로잇 체크리스트
- `skills/linux-lpe/scripts/` — LinPEAS/LinEnum 출력 자동 파싱 스크립트
- 워커가 호스트 접근 권한을 얻으면 LPE skill을 자동 로드하는 플래너 프롬프트 힌트

**예상 파일**: `skills/linux-lpe/SKILL.md`, `skills/linux-lpe/scripts/`

---

#### 4-4. 취약점 복테스트 자동화 강화 — retester 에이전트 UI 통합

**근거**: `agent/retester.go`가 이미 복테스트 로직을 구현하고 있으나, 취약점 목록에서 "복테스트 실행" 버튼과 결과(reproduced/fixed/inconclusive) 표시가 UI에 명확히 노출되어 있지 않다. 경쟁사 대비 이미 보유한 기능을 UI에서 더 잘 드러내면 된다.

**구현 방향**:
- 취약점 상세 화면에 "복테스트 실행" 버튼 추가
- 복테스트 이력 테이블 (날짜, 판정, 증거 요약)
- 복테스트 결과 `fixed`일 때 취약점 상태 자동 변경

**예상 파일**: `web/` 취약점 컴포넌트, `server/server.go` 복테스트 엔드포인트

---

#### 4-5. 연속 모니터링 모드 (24/7 자율 재실행)

**근거**: artex-deep-hunter fork와 Pentest-Swarm-AI는 연속 모니터링을 차별점으로 내세운다. 현재 artex-ko는 작업을 수동으로 시작하는 구조인데, 대상을 등록하면 주기적으로 자동 재스캔하고 새 취약점만 알림으로 보내는 모드를 추가하면 버그바운티·정기 감사 사용 사례를 커버할 수 있다.

**구현 방향**:
- 작업에 `schedule` 필드(cron 표현식 또는 간격) 추가
- `server/`의 스케줄러 루프에 주기적 재시작 로직 추가
- 이전 실행 대비 신규 취약점만 알림으로 전송 (`notify/` 연동)
- UI에 "예약된 작업" 탭 추가

**예상 파일**: `server/server.go`, `db/`, `notify/`, `web/`

---

### 🟡 P2 — 중간 임팩트 / 중간~높은 난이도

#### 4-6. 벡터 메모리 — 작업 간 경험 재활용 (pgvector)

**근거**: PentAGI의 가장 큰 차별점 중 하나다. 현재 artex-ko의 worker trace 공유는 같은 작업 내에서만 작동한다. pgvector를 활용해 과거 작업의 성공 접근법, 발견된 패턴, 익스플로잇 snippets를 임베딩으로 저장하고 새 작업 시 유사 경험을 검색해 주입하면 반복 타겟이나 유사 환경에서 성능이 대폭 향상된다.

**구현 방향**:
- PostgreSQL에 `pgvector` 확장 추가
- `memory/` 패키지 신설: 경험 임베딩 저장·검색
- Worker 시작 시 유사 작업 경험 검색 → 초기 컨텍스트에 주입
- 임베딩 모델은 LLM 설정에서 선택 가능하게

**예상 파일**: `db/schema.go`, `memory/` 신설, `agent/worker.go`, `agent/prompt.go`

---

#### 4-7. 화이트박스 소스코드 분석 skill

**근거**: Shannon은 소스코드를 읽어 공격 벡터를 먼저 파악한 뒤 익스플로잇을 실행한다. vulnhuntr는 콜 체인 추적으로 복잡한 다단계 취약점을 발견한다. artex-ko의 skill 시스템에 소스코드 분석 skill을 추가하면 화이트박스/그레이박스 테스트 사용 사례를 커버할 수 있다.

**구현 방향**:
- `skills/whitebox-analysis/SKILL.md` 신규 작성
- 소스코드 디렉터리를 입력받아 위험 함수(`eval`, `exec`, `query`, `innerHTML` 등) 추적
- 사용자 입력 → 위험 싱크 콜 체인 맵핑 (Python ast / Go ast / JS babel 파서 활용)
- 발견된 콜 체인을 `report_finding` 후보로 플래너에게 전달

**예상 파일**: `skills/whitebox-analysis/SKILL.md`, `skills/whitebox-analysis/scripts/`

---

#### 4-8. 자동 탐지 규칙 생성 (Sigma/Suricata)

**근거**: artex-ko는 이미 `detections/` 디렉터리에 Sigma·Suricata 규칙을 수동으로 관리한다. FORGE 논문은 익스플로잇 트레이스에서 Sigma/Snort 규칙을 자동으로 생성하는 멀티에이전트 파이프라인을 구현했다. artex-ko의 MITM 프록시 트래픽 데이터와 worker activity 로그를 활용하면 취약점 발견 → 탐지 규칙 자동 생성이 가능하다.

**구현 방향**:
- finding reporter에 "탐지 규칙 생성" 단계 추가
- 익스플로잇 트래픽 패턴(User-Agent, URI 패턴, 페이로드 특징)을 추출해 Sigma 템플릿으로 변환
- 생성된 규칙을 `detections/` 자동 추가, 기존 규칙과 중복 제거

**예상 파일**: `report/`, `agent/finding_workflow.go`, `detections/`

---

#### 4-9. CTF 자율 모드 (Supervisor-Executor 루프)

**근거**: PentestGPT는 SQLite 상태 지속과 반복 루프로 CTF 문제를 완전 자율로 푼다. artex-ko의 플래너-워커 구조는 이미 유사하지만, CTF 특화 목표 판정(플래그 발견)과 반복 한도 설정, 이전 시도 컨텍스트 복원 기능이 없다.

**구현 방향**:
- 작업 유형에 `ctf` 타입 추가 (목표: 플래그 패턴 `flag{...}` 또는 `HTB{...}` 발견)
- 플래그 발견 시 작업 자동 완료 처리
- CTF 모드에서 최대 반복 횟수·시간 제한 설정 옵션

**예상 파일**: `agent/goals.go`, `agent/planner.go`, `server/server.go` (작업 생성 UI)

---

#### 4-10. CI/CD 통합 — GitHub Actions / GitLab CI 지원

**근거**: Shannon은 GitHub Actions marketplace를 통해 코드 Push 시 자동 보안 테스트를 실행한다. artex-ko에 HTTP API가 이미 있으므로, 공식 GitHub Actions action과 간단한 CLI 래퍼를 제공하면 DevSecOps 파이프라인에 쉽게 통합할 수 있다.

**구현 방향**:
- `artex-action` GitHub Actions marketplace action 패키지 (Docker 기반)
- API key 기반 원격 작업 실행 엔드포인트 강화
- 작업 완료 시 SARIF 형식으로 결과 출력 (GitHub Security tab 연동)
- `.github/workflows/` 예시 파일 작성

**예상 파일**: `.github/actions/artex-scan/`, `server/server.go` (SARIF 내보내기), `docs/`

---

### 🟢 P3 — 낮은~중간 임팩트 / 낮은 난이도 (빠르게 추가 가능)

#### 4-11. skill 추가 — API 보안 테스트 (GraphQL, gRPC, WebSocket)

**근거**: 현재 `api-recon` skill은 REST/JSON 중심이다. GraphQL 인트로스펙션, gRPC reflection, WebSocket 메시지 퍼징을 별도 skill로 분리하면 현대적 API 스택 커버리지가 늘어난다.

**구현 방향**:
- `skills/graphql-audit/SKILL.md` — 인트로스펙션 → 뮤테이션 퍼징 → 권한 우회
- `skills/websocket-fuzz/SKILL.md` — 메시지 타입 열거, 인증 토큰 주입 테스트

---

#### 4-12. skill 추가 — 클라우드 환경(AWS/GCP/Azure) 오구성 감사

**근거**: 현대 버그바운티와 레드팀 작업의 큰 비중이 클라우드 오구성이다. SSRF → AWS IMDSv1, 공개 S3 버킷, 과도한 IAM 권한 등을 체계적으로 테스트하는 skill이 없다.

**구현 방향**:
- `skills/cloud-audit/SKILL.md` — AWS IMDSv2 우회, S3 ACL, GCP metadata server
- `skills/cloud-audit/scripts/` — AWS/GCP/Azure CLI 체크 스크립트

---

#### 4-13. 취약점 중복 제거 및 병합 로직 강화

**근거**: 여러 worker가 동시에 같은 취약점을 발견하면 중복 등록이 발생할 수 있다. 현재 `finding_recorder.go`에 TrafficRef 연동이 있지만, URL+취약점 유형+파라미터 기반 중복 감지 로직이 명시적으로 없다.

**구현 방향**:
- `report_finding` 호출 시 동일 대상·유형 기존 취약점 조회
- 중복으로 판단되면 에이전트에게 기존 ID를 반환하고 증거만 보강

**예상 파일**: `agent/finding_recorder.go`, `db/`

---

#### 4-14. 보고서 형식 다양화 — PDF/HTML 내보내기

**근거**: 현재 CSV 내보내기만 있다. 클라이언트·팀에 전달 가능한 전문 포맷(PDF, HTML)이 없어 별도 작업이 필요하다. PentestGPT와 CyberStrike는 PDF 보고서를 기본 제공한다.

**구현 방향**:
- `report/` 패키지에 HTML 템플릿 렌더러 추가
- headless Chrome(playwright)으로 HTML → PDF 변환
- `/api/tasks/{id}/report` 엔드포인트에 `format=pdf|html|json|csv` 파라미터 추가

---

#### 4-15. 사용자 자원 라이브러리 — 파일 업로드 및 에이전트 접근

**근거**: PentAGI는 사용자가 사전 수집한 wordlist, 인증서, 설정 파일을 업로드하고 에이전트가 작업 중에 참조할 수 있는 가상 경로 파일시스템을 제공한다. artex-ko에서는 외부 파일을 에이전트에게 제공할 방법이 없어 사전 준비된 자원을 활용하기 어렵다.

**구현 방향**:
- `/api/resources` CRUD 엔드포인트 (업로드, 목록, 삭제)
- 작업 시작 시 선택한 자원을 worker 샌드박스 디렉터리에 마운트
- 에이전트 도구에 `list_resources`, `read_resource` 추가

---

#### 4-16. 2FA/TOTP 자동 처리 — playwright-cli skill 확장

**근거**: PentestGPT는 2FA/TOTP 로그인 자동화를 차별점으로 내세운다. artex-ko는 이미 `playwright-cli` skill을 가지고 있어, TOTP 생성(`otpauth://` URL 파싱 + TOTP 계산)을 skill 스크립트로 추가하면 된다.

**구현 방향**:
- `skills/playwright-cli/references/totp-login.md` 추가 — TOTP 시드 등록, playwright로 OTP 필드 자동 입력
- Python `pyotp` 또는 순수 stdlib TOTP 계산 스크립트

---

#### 4-17. 에이전트 성능 벤치마크 평가 모드

**근거**: hackingBuddyGPT와 PentestGPT는 재현 가능한 벤치마크 VM/CTF 세트로 에이전트 성능을 측정한다. artex-ko에 "벤치마크 모드"를 추가하면 모델 교체·프롬프트 개선의 효과를 객관적으로 측정할 수 있다.

**구현 방향**:
- DVWA, OWASP Juice Shop, HackTheBox 머신 정의 파일 (`benchmarks/` 디렉터리)
- 작업 완료 후 목표 취약점 목록 대비 발견율(recall) 자동 계산
- 여러 LLM 모델/설정 비교 실행 스크립트

---

## 5. 결론 및 종합 권고

artex-ko는 이중 그래프 구조, MITM 프록시, skill 시스템, MCP 지원, 취약점 복테스트라는 경쟁사 대비 고유한 강점을 가진다. 아래 순서로 기능을 추가하면 단기간에 가장 큰 차별화 효과를 얻을 수 있다.

| 우선순위 | 기능 | 예상 임팩트 |
|---|---|---|
| ① | CVSS v3.1 자동 채점 + PoC 검증 게이트 | 보고서 신뢰도·전문성 대폭 향상 |
| ② | LLM 옵저버빌리티 대시보드 (토큰·비용·레이턴시) | 팀/엔터프라이즈 비용 관리 필수 |
| ③ | 연속 모니터링 모드 (주기적 재스캔 + 신규 취약점 알림) | 버그바운티·정기 감사 사용 사례 커버 |
| ④ | Linux LPE skill 추가 | 웹 너머 호스트 침투 커버리지 확장 |
| ⑤ | 화이트박스 소스코드 분석 skill | Shannon/vulnhuntr와 차별화 |
| ⑥ | 벡터 메모리 (pgvector) — 작업 간 경험 재활용 | 반복 타겟 성능 대폭 향상 |
| ⑦ | PDF/HTML 보고서 내보내기 | 전달 가능한 전문 보고서 |
| ⑧ | CI/CD 통합 (GitHub Actions) | DevSecOps 파이프라인 편입 |
| ⑨ | 자동 탐지 규칙 생성 (Sigma) | 기존 detections/ 자원과 시너지 |
| ⑩ | 에이전트 성능 벤치마크 모드 | 모델/프롬프트 개선 객관적 측정 |

---

## 참고 출처

- [vxcontrol/pentagi](https://github.com/vxcontrol/pentagi) — Go 백엔드, pgvector, Neo4j 지식 그래프
- [GreyDGL/PentestGPT](https://github.com/GreyDGL/PentestGPT) — USENIX Security 2024, PTT, 자율 CTF 파이프라인
- [ipa-lab/hackingBuddyGPT](https://github.com/ipa-lab/hackingBuddyGPT) — Linux LPE 벤치마크, 50줄 에이전트
- [KeygraphHQ/shannon](https://github.com/KeygraphHQ/shannon) — 화이트박스 소스코드 분석, XBOW 96.15%
- [CyberStrikeus/CyberStrike](https://github.com/CyberStrikeus/CyberStrike) — 7,600+ 스킬, 120+ OWASP 케이스
- [insidetrust/awesome-ai-pentest](https://github.com/insidetrust/awesome-ai-pentest) — AI 침투 테스트 종합 목록
- [FORGE arXiv:2606.03453](https://arxiv.org/pdf/2606.03453) — CVE 익스플로잇 + Sigma 자동 생성
- [lordx64/pentestkit](https://github.com/lordx64/pentestkit) — PoC 필수 게이트, CVSS v3
- [protectai/vulnhuntr](https://github.com/protectai/vulnhuntr) — 소스코드 콜 체인 취약점 발견
- [logkiller/artex-deep-hunter](https://github.com/logkiller/artex-deep-hunter) — ARTEX 포크, 캠페인 브랜치
- artex-ko 소스: `README.md`, `server/server.go`, `agent/*.go`, `skills/*/SKILL.md` (직접 확인)
