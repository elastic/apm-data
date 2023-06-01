package modelpb

import "github.com/elastic/apm-data/model/internal/modeljson"

func (u *User) toModelJSON(out *modeljson.User) {
	*out = modeljson.User{
		Domain: u.Domain,
		ID:     u.Id,
		Email:  u.Email,
		Name:   u.Name,
	}
}
