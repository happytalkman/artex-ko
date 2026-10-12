# ARTEX-KO 설정 보고서

생성일시: 2026-10-12

---

## 프로젝트 개요

**artex-ko**는 [Autumn-27/ARTEX](https://github.com/Autumn-27/ARTEX)의 한국어판 포크로,
LLM 멀티 에이전트가 자율적으로 침투 테스트를 수행하는 공격 보안 도구입니다.

- Go 백엔드(단일 바이너리) + Next.js 프런트엔드(선택적 내장)
- SQLite 기반 듀얼 그래프 스토어, 이벤트 기반 탐색 엔진
- JSON HTTP API (기본 포트 8787)
- MITM 프록시를 통한 트래픽 기록 (기본 포트 8788)
- 자동 업데이트(셀프업데이트) 기능 내장
- 출력 언어 기본값: 한국어 (ko)

---

## 의존성 설치 결과

`go mod download` 명령 실행 완료 — 성공.

- Go 버전: go1.27.2 windows/amd64 (요구: 1.26 이상)
- 주요 의존성:
  - `modernc.org/sqlite` v1.52.0 — CGo 없는 순수 Go SQLite
  - `github.com/lqqyt2423/go-mitmproxy` v1.9.2 — 트래픽 기록 프록시
  - `github.com/Autumn-27/norma` v0.4.3 — 스캔/탐지 엔진
  - `github.com/projectdiscovery/dnsx` v1.2.3 — DNS 조회
  - `github.com/jackc/pgx/v5` v5.10.0 — PostgreSQL 클라이언트(선택적)
  - `github.com/golang-jwt/jwt/v5` v5.3.1 — 인증

---

## 빌드 결과

**빌드 성공.**

```
go build -o artex.exe ./cmd/artex/...
```

| 항목 | 값 |
|------|-----|
| 빌드 태그 | (없음, 개발 빌드) |
| 바이너리 위치 | `c:\apex\artex-ko\artex.exe` |
| 파일 크기 | ~44 MB |
| 프런트엔드 내장 | 아니오 (개발 빌드 — API 서버만 동작) |

> 프런트엔드를 내장하려면 아래 "프런트엔드 내장 빌드" 절을 참조하세요.

---

## 실행 명령

### 즉시 실행 (개발 / API 서버)

```powershell
# 기본 실행 (127.0.0.1:8787 에서 수신)
cd c:\apex\artex-ko
.\artex.exe

# 포트·데이터 디렉터리 지정
.\artex.exe -addr 127.0.0.1:8787 -data c:\apex\artex-ko\data

# 트래픽 프록시 비활성화
.\artex.exe -proxy ""

# go run으로 바로 실행
go run ./cmd/artex
```

### 플래그 전체 목록

| 플래그 | 기본값 | 설명 |
|--------|--------|------|
| `-addr` | `127.0.0.1:8787` | HTTP 수신 주소 (다른 호스트 허용: `:8787`) |
| `-data` | `.\data` | SQLite 스토어용 데이터 디렉터리 |
| `-proxy` | `127.0.0.1:8788` | 트래픽 기록 프록시 주소 (비우면 비활성화) |

### 처음 실행 후

1. 브라우저에서 `http://127.0.0.1:8787` 접속
2. `/setup` 화면에서 관리자 비밀번호 설정 (로그에 찍힌 토큰 필요)
3. 설정 화면에서 LLM API 키 입력

---

## 필수 환경 변수

아래 변수는 `.env` 파일 또는 시스템 환경 변수로 설정합니다.
**LLM API 키 중 하나 이상은 반드시 입력해야 에이전트가 동작합니다.**

| 변수 | 필수 여부 | 설명 |
|------|-----------|------|
| `ANTHROPIC_API_KEY` | LLM 사용 시 필수 | Anthropic Claude API 키 |
| `OPENAI_API_KEY` | LLM 사용 시 필수 | OpenAI API 키 |
| `ARTEX_LLM_PROVIDER` | 선택 | LLM 공급자 (anthropic / openai) |
| `ARTEX_LLM_MODEL` | 선택 | 사용할 모델명 |
| `ARTEX_LLM_BASE_URL` | 선택 | 커스텀 LLM 엔드포인트 URL |
| `ARTEX_LLM_PROXY` | 선택 | LLM 호출용 프록시 |
| `ARTEX_BIND` | 선택 | 수신 IP (기본 `127.0.0.1`, 외부 허용: `0.0.0.0`) |
| `ARTEX_SETUP_TOKEN` | 선택 | 초기 설정 토큰 (미설정 시 랜덤 생성 후 로그에 출력) |
| `ARTEX_ADMIN_PASSWORD` | 선택 | 초기 관리자 비밀번호 (8자 이상, 설정 시 /setup 건너뜀) |
| `ARTEX_CORS_ORIGINS` | 선택 | CORS 허용 출처 목록 (쉼표 구분) |
| `ARTEX_LANG` | 선택 | 기본 출력 언어 (`ko`·`en`·`zh`·`es`, 기본: `ko`) |
| `ARTEX_PG_DSN` | 선택 | PostgreSQL DSN (config.json의 database 항목보다 우선) |
| `ARTEX_SKILL_DIR` | 선택 | 스킬 디렉터리 경로 (기본: 실행 파일 옆 `skills/`) |
| `ARTEX_CONFIG` | 선택 | config.json 대체 경로 |

> Docker Compose 배포 시 추가로 필요한 변수: `ARTEX_TAG`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`

---

## 생성된 설정 파일

### `.env` (c:\apex\artex-ko\.env)
`.env.example`에서 복사, `ARTEX_LANG=ko` 기본값 적용.
Docker Compose 실행 시 사용. **LLM API 키를 직접 채워 넣어야 합니다.**

### `config.json` (c:\apex\artex-ko\config.json)
`config.example.json`에서 복사.
- DB 접속: `127.0.0.1:5433` (PostgreSQL, 선택 사항)
- 스킬 디렉터리: `skills` (현재 작업 디렉터리 기준 상대 경로)
- 출력 언어: `ko`

> **참고**: artex는 SQLite를 기본 스토어로 사용합니다. PostgreSQL(`config.json`의 `database`)은 선택 사항이며, DSN이 비어 있으면 SQLite만 사용합니다.

---

## 프런트엔드 내장 빌드

현재 바이너리는 프런트엔드 없이 빌드됐습니다. API는 동작하지만
`http://127.0.0.1:8787`에서 UI가 표시되지 않습니다.

UI 포함 바이너리를 빌드하려면 Node.js 22 이상이 필요합니다:

```powershell
# 프런트엔드 빌드
cd c:\apex\artex-ko\web
npm ci
npm run build:static

# 빌드 결과를 server/webui/dist로 복사
New-Item -ItemType Directory -Force c:\apex\artex-ko\server\webui\dist
Copy-Item -Recurse -Force c:\apex\artex-ko\web\out\* c:\apex\artex-ko\server\webui\dist\

# embedui 태그로 백엔드 재빌드
cd c:\apex\artex-ko
go build -tags embedui -o artex.exe ./cmd/artex/...
```

---

## 발생한 문제 및 해결 방법

| 문제 | 원인 | 해결 |
|------|------|------|
| 빌드 시간 지연 (~5분) | `modernc.org/sqlite` 등 대용량 패키지 최초 컴파일 | 정상 동작, 두 번째 빌드부터 캐시 적용되어 빠름 |
| `go.mod`의 `go 1.26.3` | 설치된 Go 1.27.2로 빌드 가능 | 문제 없음 |
| 프런트엔드 없음 | 개발 빌드 (embedui 태그 미사용) | API 서버로 동작, UI 필요 시 위의 "프런트엔드 내장 빌드" 절 참조 |

---

## 사용자가 직접 해야 할 작업

1. **LLM API 키 입력 (필수)**
   - `.env` 파일의 `ANTHROPIC_API_KEY` 또는 `OPENAI_API_KEY` 채우기
   - 또는 `.\artex.exe` 실행 후 웹 UI 설정 화면에서 입력

2. **초기 관리자 계정 설정**
   - `.\artex.exe` 첫 실행 → 로그에서 설정 토큰 확인
   - `http://127.0.0.1:8787/setup` 접속 → 토큰 입력 → 관리자 비밀번호 설정
   - 또는 `.env`에 `ARTEX_ADMIN_PASSWORD=<8자이상>` 설정하면 자동 초기화

3. **PostgreSQL 연결 (선택)**
   - PostgreSQL 사용 시 `config.json`의 `database.password` 채우기
   - 또는 `ARTEX_PG_DSN` 환경 변수로 DSN 직접 지정
   - 미설정 시 SQLite로 동작 (기본)

4. **프런트엔드 내장 빌드 (UI 필요 시)**
   - Node.js 22 이상 설치 후 위의 "프런트엔드 내장 빌드" 절 실행

5. **보안 경고**: ARTEX는 허가받은 대상에만 사용해야 합니다. 외부 노출 시(`-addr :8787` 또는 `ARTEX_BIND=0.0.0.0`) 초기 설정을 반드시 먼저 완료하세요.

---

## 빠른 시작 요약

```powershell
# 1. API 키 설정 (편집기에서 .env 열기)
notepad c:\apex\artex-ko\.env

# 2. 실행
cd c:\apex\artex-ko
.\artex.exe

# 3. 브라우저 접속
# http://127.0.0.1:8787
```
