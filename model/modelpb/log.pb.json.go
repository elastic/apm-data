package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (l *Log) toModelJSON(out *modeljson.Log) {
	*out = modeljson.Log{
		Level:  l.Level,
		Logger: l.Logger,
	}
	if l.Origin != nil {
		out.Origin = modeljson.LogOrigin{
			Function: l.Origin.FunctionName,
		}
		if l.Origin.File != nil {
			out.Origin.File = modeljson.LogOriginFile{
				Name: l.Origin.File.Name,
				Line: int(l.Origin.File.Line),
			}
		}
	}
}
