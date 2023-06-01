package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (u *UserExperience) toModelJSON(out *modeljson.UserExperience) {
	*out = modeljson.UserExperience{
		CumulativeLayoutShift: u.CumulativeLayoutShift,
		FirstInputDelay:       u.FirstInputDelay,
		TotalBlockingTime:     u.TotalBlockingTime,
	}
	if u.LongTask != nil {
		out.Longtask = modeljson.LongtaskMetrics{
			Count: int(u.LongTask.Count),
			Sum:   u.LongTask.Sum,
			Max:   u.LongTask.Max,
		}
	}
}
