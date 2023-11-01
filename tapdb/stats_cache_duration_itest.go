//go:build itest

package tapdb

import "time"

// StatsCacheDuration is the duration for which the stats cache is valid. For
// itests, we reduce this to pretty much nothing.
const StatsCacheDuration = time.Microsecond * 1
