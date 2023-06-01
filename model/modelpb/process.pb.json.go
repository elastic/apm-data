package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (p *Process) toModelJSON(out *modeljson.Process) {
	*out = modeljson.Process{
		Pid:         int(p.Pid),
		Title:       p.Title,
		CommandLine: p.CommandLine,
		Executable:  p.Executable,
		Args:        p.Argv,
		Parent: modeljson.ProcessParent{
			Pid: p.Ppid,
		},
	}
	if p.Thread != nil {
		out.Thread = modeljson.ProcessThread{
			Name: p.Thread.Name,
			ID:   int(p.Thread.Id),
		}
	}
}
