package categorize

import (
	"messenger/jobs/scheduler"
	"time"
)

func init() {
	scheduler.Instance.Schedule(time.Second, func() {
		newJob().Run()
	})
}
