// Copyright Contributors to the Open Cluster Management project
package metrics

import (
	"time"

	"github.com/stolostron/search-indexer/pkg/config"
	"k8s.io/klog/v2"
)

var DEFAULT_SLOW_LOG = time.Duration(config.Cfg.SlowLog) * time.Millisecond

// Record the time when a function starts and logs if the function takes more than the expected duration.
// The returned function should be invoked with defer.
func SlowLog(msg string, logAfter time.Duration) func() {
	start := time.Now()

	// This function should be invoked with defer to execute at the end of the caller function.
	return func() {
		if (logAfter > 0 && time.Since(start) > logAfter) || (time.Since(start) > DEFAULT_SLOW_LOG) {
			klog.Warningf("%s - %s", time.Since(start), msg)
		}

		// We could emit metric here, but it could emit too much data.
	}
}
