package server

import (
	"net/http"
)

// getLLMStats handles GET /api/llm/stats and returns aggregated LLM usage
// statistics: per-profile totals, per-day buckets, and per-worker latency/error stats.
func (s *Server) getLLMStats(w http.ResponseWriter, r *http.Request) {
	profiles, err := s.m.pg.UsageByProfile()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "LLM 통계 조회 실패: 프로필별")
		return
	}
	daily, err := s.m.pg.UsageDaily(30)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "LLM 통계 조회 실패: 일별")
		return
	}
	workers, err := s.m.pg.WorkerLatencyStats()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "LLM 통계 조회 실패: 워커별")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"profiles": profiles,
		"daily":    daily,
		"workers":  workers,
	})
}
