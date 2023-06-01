package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (f *Faas) toModelJSON(out *modeljson.FAAS) {
	*out = modeljson.FAAS{
		ID:        f.Id,
		Name:      f.Name,
		Version:   f.Version,
		Execution: f.Execution,
		Coldstart: f.ColdStart,
		Trigger: modeljson.FAASTrigger{
			Type:      f.TriggerType,
			RequestID: f.TriggerRequestId,
		},
	}

}
