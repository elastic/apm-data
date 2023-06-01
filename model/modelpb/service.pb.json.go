package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (s *Service) toModelJSON(out *modeljson.Service) {
	*out = modeljson.Service{
		Name:        s.Name,
		Version:     s.Version,
		Environment: s.Environment,
	}
	if s.Node != nil {
		out.Node = &modeljson.ServiceNode{
			Name: s.Node.Name,
		}
	}
	if s.Language != nil {
		out.Language = &modeljson.Language{
			Name:    s.Language.Name,
			Version: s.Language.Version,
		}
	}
	if s.Runtime != nil {
		out.Runtime = &modeljson.Runtime{
			Name:    s.Runtime.Name,
			Version: s.Runtime.Version,
		}
	}
	if s.Framework != nil {
		out.Framework = &modeljson.Framework{
			Name:    s.Framework.Name,
			Version: s.Framework.Version,
		}
	}
	if s.Origin != nil {
		out.Origin = &modeljson.ServiceOrigin{
			ID:      s.Origin.Id,
			Name:    s.Origin.Name,
			Version: s.Origin.Version,
		}
	}
	if s.Target != nil {
		out.Target = &modeljson.ServiceTarget{
			Name: s.Target.Name,
			Type: s.Target.Type,
		}
	}
}
