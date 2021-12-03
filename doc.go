package harfbuzzgoperf

import "github.com/npillmayer/schuko/tracing"

// tracer traces to tracing key 'hbperf.base'.
func tracer() tracing.Trace {
	return tracing.Select("hbperf.base")
}
