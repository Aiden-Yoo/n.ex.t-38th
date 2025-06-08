package util

import (
	"fmt"
	"io"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

// ResultFormatter는 결과를 문자열로 포맷팅하는 인터페이스
type ResultFormatter interface {
	Format() string
}

// CustomRenderResults는 결과를 커스텀 포맷으로 출력
func CustomRenderResults(w io.Writer, results []*gornir.JobResult, title string) {
	fmt.Fprintf(w, "\n====== %s ======\n", title)

	var success, failed int
	for _, r := range results {
		if r.Err() != nil {
			failed++
		} else {
			success++
		}
	}

	fmt.Fprintf(w, "총 장비: %d, 성공: %d, 실패: %d\n\n", len(results), success, failed)

	for _, r := range results {
		hostname := r.Host().Hostname

		if r.Err() != nil {
			fmt.Fprintf(w, "❌ %s: 오류 발생 - %v\n", hostname, r.Err())
			continue
		}

		data := r.Data()

		// TaskAll: []interface{} 처리
		if list, ok := data.([]interface{}); ok {
			fmt.Fprintf(w, "✅ %s:\n", hostname)
			for _, item := range list {
				if formatter, ok := item.(ResultFormatter); ok {
					fmt.Fprintf(w, "%s\n", formatter.Format())
				} else {
					fmt.Fprintf(w, "  - (알 수 없는 결과 타입: %T)\n", item)
				}
			}
			continue
		}

		// 단일 Task 결과
		if formatter, ok := data.(ResultFormatter); ok {
			fmt.Fprintf(w, "✅ %s:\n%s\n", hostname, formatter.Format())
		} else {
			fmt.Fprintf(w, "❌ %s: 지원하지 않는 형식 (%T)\n", hostname, data)
		}
	}
}
