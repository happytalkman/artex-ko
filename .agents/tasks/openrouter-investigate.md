# OpenRouter API 통합 조사 보고서

## 요약 결론

**OpenRouter는 ARTEX에서 즉시 사용 가능하다.** OpenRouter는 OpenAI 호환 API를 제공하므로, UI에서 LLM 프로필을 생성할 때 `format=openai`, `base_url=https://openrouter.ai/api/v1`, `api_key=<openrouter키>` 를 설정하면 추가 코드 수정 없이 동작한다.

**유일한 제약:** OpenRouter가 권장하는 `HTTP-Referer` 및 `X-Title` 헤더는 현재 코드에서 지원하지 않는다. 그러나 이 헤더들은 **선택사항(Optional)** 이며 API 동작에 필수가 아니다 — 이것이 없어도 요청은 정상 처리된다.

---

## 1. LLM 공급자 설정 방식

### 1.1 두 가지 설정 경로

ARTEX는 LLM 공급자를 두 가지 방식으로 설정한다:

**경로 A: 환경변수 기반 (레거시/단순)**
- `agent/provider.go` → `FromEnv()` 함수
- 관련 환경변수:
  - `ARTEX_LLM_PROVIDER` = `anthropic` | `openai` (기본값: 키로 추론)
  - `ARTEX_LLM_MODEL` = 모델 ID (기본값: 공급자별 기본 모델)
  - `ARTEX_LLM_BASE_URL` = 커스텀 엔드포인트 (선택)
  - `ARTEX_LLM_PROXY` = 프록시 URL (선택)
  - `ANTHROPIC_API_KEY` / `OPENAI_API_KEY` = 인증 키
- OpenRouter 사용 시: `ARTEX_LLM_PROVIDER=openai`, `OPENAI_API_KEY=sk-or-...`, `ARTEX_LLM_BASE_URL=https://openrouter.ai/api/v1`

**경로 B: UI/DB 기반 LLM 프로필 (권장, 다중 프로필 지원)**
- `db/config.go` → `LLMProfile` 구조체 (PostgreSQL `llm_profiles` 테이블)
- API 엔드포인트 (`server/server.go` 라우트 등록 확인):
  - `GET  /api/llm/profiles` — 프로필 목록
  - `POST /api/llm/profiles` — 프로필 생성/수정
  - `DELETE /api/llm/profiles/{id}` — 프로필 삭제
  - `POST /api/llm/profiles/active` — 활성 프로필 지정
  - `POST /api/llm/test` — 연결 테스트
- UI 경로: 시스템 → LLM 설정

### 1.2 `LLMProfile` 구조체 주요 필드 (`db/config.go`)

```go
type LLMProfile struct {
    ID       int64   `json:"id"`
    Name     string  `json:"name"`
    Format   string  `json:"format"`    // "anthropic" | "openai" | "openai-responses"
    BaseURL  string  `json:"base_url"`  // 커스텀 엔드포인트 (선택)
    Proxy    string  `json:"proxy"`     // 프록시 URL (선택)
    Model    string  `json:"model"`
    APIKey   string  `json:"-"`         // UI에 노출 안 됨; api_key 필드로 POST
    // ... 기타 설정들
    Streaming bool   `json:"streaming"` // 기본값: true (SSE)
    MaxTokens int    `json:"max_tokens"`
    SessionHeaderKey string `json:"session_header_key"` // 커스텀 헤더명(값은 session id)
}
```

---

## 2. OpenRouter base_url 설정 동작 여부

### 2.1 결론: 동작한다

`agent/provider.go` → `ConfigFrom()` 함수 분석:

```go
func ConfigFrom(provider, model, baseURL, apiKey, proxy string) Config {
    // ...
    switch strings.TrimSpace(provider) {
    case "openai":
        c.Format = llm.FormatOpenAI
        // provider appends "/chat/completions"; tolerate a full endpoint URL.
        c.BaseURL = strings.TrimRight(strings.TrimSuffix(c.BaseURL, "/chat/completions"), "/")
        // ...
    }
    return c
}
```

`agent.Config.NewProvider()` 호출 시:
```go
lc := llm.Config{
    Format:  c.Format,
    BaseURL: c.BaseURL,  // = https://openrouter.ai/api/v1
    APIKey:  c.APIKey,
    Model:   c.Model,
    // ...
}
return llm.NewProvider(lc)
```

`norma@v0.4.3/llm/openai.go` → `openaiProvider` 는 `BaseURL + "/chat/completions"` 로 요청을 보낸다:
```go
resp, err := doStream(ctx, p.cfg, p.cfg.BaseURL+"/chat/completions", body, func(r *http.Request) {
    r.Header.Set("content-type", "application/json")
    r.Header.Set("authorization", "Bearer "+p.cfg.APIKey)
    r.Header.Set("accept", "text/event-stream")
}, "openai")
```

따라서 `base_url = https://openrouter.ai/api/v1` 로 설정하면 실제 요청은 `https://openrouter.ai/api/v1/chat/completions` 로 전송된다. **OpenRouter가 이 경로를 정확히 제공하므로 정상 동작한다.**

---

## 3. 추가 헤더 지원 여부 (HTTP-Referer, X-Title)

### 3.1 현재 코드 분석

`norma@v0.4.3/llm/openai.go` 에서 전송하는 HTTP 헤더:
- `content-type: application/json`
- `authorization: Bearer <api_key>`
- `accept: text/event-stream` (스트리밍) / `application/json` (논스트리밍)

`ExtraHeaders`, `AdditionalHeaders` 등의 커스텀 헤더 필드는 **코드베이스 전체에 존재하지 않는다** (grep 결과 zero matches).

`SessionHeaderKey` 필드가 존재하나, 이것은 헤더 이름을 설정하고 **값은 현재 세션 ID가 자동으로 주입되는** 특수 목적 필드이다 — 임의 값을 가진 임의 헤더를 추가하는 범용 메커니즘이 아니다.

### 3.2 OpenRouter 헤더 요구사항

[OpenRouter 공식 문서](https://openrouter.ai/docs/api/reference/authentication)에 따르면:

| 헤더 | 필수 여부 | 용도 |
|------|----------|------|
| `Authorization: Bearer <key>` | **필수** | 인증 (norma가 자동 설정) |
| `HTTP-Referer: <site_url>` | **선택** | OpenRouter 리더보드 랭킹용 |
| `X-OpenRouter-Title: <name>` | **선택** | OpenRouter 리더보드 랭킹용 |

**`HTTP-Referer`와 `X-Title`은 선택사항**이다. 없어도 API 요청은 정상적으로 처리된다. 단지 OpenRouter 통계 페이지에서 앱이 표시되지 않을 뿐이다.

---

## 4. LLM 프로필 생성 API 구조 (`POST /api/llm/profiles`)

`server/server_mgmt.go` → `pgSaveProfile` 핸들러 분석:

### 요청 본문 (JSON)

```json
{
  "name": "OpenRouter Free",
  "format": "openai",
  "base_url": "https://openrouter.ai/api/v1",
  "model": "google/gemma-4-26b-a4b-it:free",
  "api_key": "sk-or-v1-xxxxxxxxxxxx",
  "proxy": "",
  "streaming": true,
  "max_tokens": 0,
  "rate_per_second": 0,
  "rate_per_minute": 0,
  "context_window_k": 256,
  "reasoning_effort": "",
  "thinking_type": "",
  "session_header_key": "",
  "priority": 0,
  "pool_exclude": false
}
```

### 응답

```json
{"id": 123}
```

이후 `POST /api/llm/profiles/active` 로 `{"id": 123}` 를 전송해 해당 프로필을 활성화해야 한다.

### UI 경로

시스템 → LLM 설정 → [+ 추가] 버튼 → 공급자: **openai**, 모델: 모델 ID, Base URL: `https://openrouter.ai/api/v1`, API Key: OpenRouter 키 입력.

---

## 5. ARTEX에 적합한 OpenRouter 무료 모델 추천

> ⚠️ 무료 모델 목록은 빠르게 변경됩니다. 사용 전 [https://openrouter.ai/collections/free-models](https://openrouter.ai/collections/free-models) 에서 현재 사용 가능한 모델을 확인하세요.

### 추천 모델 (2026년 10월 기준, [출처](https://www.teamday.ai/blog/best-free-ai-models-openrouter-2026))

| 모델 ID | 컨텍스트 | 입력 타입 | ARTEX 적합성 |
|---------|---------|----------|------------|
| `nvidia/nemotron-3-ultra-550b-a55b:free` | 1,000,000 토큰 | 텍스트 | ★★★★★ 대규모 컨텍스트, 최고 품질 |
| `nvidia/nemotron-3.5-lightning:free` | 1,000,000 토큰 | 텍스트 | ★★★★★ 대규모 컨텍스트, 빠름 |
| `nvidia/nemotron-3-super-120b-a12b:free` | 262,144 토큰 | 텍스트 | ★★★★☆ 대형 모델, 범용 |
| `google/gemma-4-26b-a4b-it:free` | 262,144 토큰 | 이미지+텍스트+비디오 | ★★★★☆ 멀티모달, Google |
| `qwen/qwen3.8-27b:free` | 262,144 토큰 | 이미지+텍스트+비디오 | ★★★★☆ 멀티모달, 강력한 코딩 |
| `poolside/laguna-s-2.1:free` | 262,144 토큰 | 텍스트 | ★★★☆☆ 코딩 특화 |

**ARTEX용 최우선 추천:** `nvidia/nemotron-3-ultra-550b-a55b:free` — 1M 토큰 컨텍스트 윈도우가 보안 탐색/분석 작업에 유리하며 텍스트 품질이 높다. ARTEX 설정 시 `context_window_k: 1000` 으로 설정.

**2순위:** `nvidia/nemotron-3.5-lightning:free` — 동일한 1M 컨텍스트, 더 빠른 응답속도.

> **참고:** 과거 자주 언급되던 `google/gemini-flash-1.5-8b:free`, `meta-llama/llama-3.1-8b-instruct:free` 는 현재 목록에서 제거된 상태이다. `:free` suffix 모델 ID는 변경이 잦으므로 최신 목록을 반드시 확인할 것.

---

## 6. OpenRouter 무료 티어 제한

| 조건 | 분당 요청 | 일일 요청 |
|------|----------|----------|
| 평생 크레딧 구매액 $10 미만 | 20 RPM | 50 RPD |
| 평생 크레딧 구매액 $10 이상 | 20 RPM | 1,000 RPD |

무료 API 키 자체는 신용카드 없이 발급 가능하며 무료 모델에는 크레딧이 차감되지 않는다. 단, 일일 50회 제한이 있으므로 프로덕션 사용보다는 테스트/소규모 탐색 목적에 적합하다.

---

## 7. 단계별 설정 방법

### UI를 통한 설정

1. ARTEX UI → **시스템** → **LLM 설정**
2. **[+ 추가]** 클릭
3. 다음 값 입력:
   - **이름:** `OpenRouter Free` (임의)
   - **공급자 형식:** `openai`
   - **Base URL:** `https://openrouter.ai/api/v1`
   - **모델:** `nvidia/nemotron-3-ultra-550b-a55b:free`
   - **API Key:** `sk-or-v1-...` (OpenRouter 대시보드에서 발급)
   - **컨텍스트 윈도우(K):** `1000` (Nemotron 3 Ultra의 경우)
4. **[저장]** → **[테스트]** 로 연결 확인
5. **[활성화]** 클릭

### 환경변수를 통한 설정 (대안)

```sh
export ARTEX_LLM_PROVIDER=openai
export OPENAI_API_KEY=sk-or-v1-xxxxxxxxxxxx
export ARTEX_LLM_BASE_URL=https://openrouter.ai/api/v1
export ARTEX_LLM_MODEL=nvidia/nemotron-3-ultra-550b-a55b:free
```

---

## 8. 향후 개선 권고사항 (구현 불필요)

현재 코드는 OpenRouter에서 정상 동작하지만, 아래 개선이 이루어지면 완성도가 높아진다:

1. **커스텀 헤더 지원 추가** (`db.LLMProfile`에 `extra_headers map[string]string` 필드 추가, `norma` 레벨에서 `llm.Config`에 `ExtraHeaders` 추가): OpenRouter `HTTP-Referer`, `X-Title` 헤더, 기타 게이트웨이 전용 헤더를 지원할 수 있게 된다.
   - 변경 필요 파일: `norma/llm/provider.go` (Config 구조체), `norma/llm/openai.go` (헤더 주입), `db/config.go` (LLMProfile 필드), `server/server_mgmt.go` (pgSaveProfile 파싱)
   - 단, `norma`는 외부 모듈(`v0.4.3`)이므로 해당 레포 변경이 선행되어야 함

2. **`SessionHeaderKey`를 OpenRouter `HTTP-Referer`로 우회 활용** (현재도 가능): `session_header_key = HTTP-Referer` 로 설정하면 모든 요청에 `HTTP-Referer: <session-id>` 헤더가 붙는다. 의미상 맞지 않으나 OpenRouter가 선택 헤더이므로 큰 문제는 없다.

---

## 참고 파일 목록

| 파일 | 역할 |
|------|------|
| `agent/provider.go` | LLM Config 정의, `ConfigFrom()`, `FromEnv()`, `NewProvider()` |
| `db/config.go` | `LLMProfile` DB 구조체, CRUD |
| `server/server.go` | HTTP 라우트 등록 (`/api/llm/*`) |
| `server/server_mgmt.go` | `pgSaveProfile`, `pgListProfiles`, `testLLM` 핸들러 |
| `server/llmpool.go` | 다중 프로필 장애조치(failover) 풀 |
| `norma@v0.4.3/llm/provider.go` | `llm.Config`, `NewProvider()` |
| `norma@v0.4.3/llm/openai.go` | OpenAI 형식 HTTP 요청 빌더 |

*Content was rephrased for compliance with licensing restrictions. Sources: [OpenRouter Auth Docs](https://openrouter.ai/docs/api/reference/authentication), [OpenRouter Free Models](https://openrouter.ai/collections/free-models), [Free Model Catalog](https://www.teamday.ai/blog/best-free-ai-models-openrouter-2026)*
