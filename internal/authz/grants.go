package authz

import "lawdrive/internal/model"

func GrantAllows(grant model.PermissionGrant, action Action, now int64) bool {
	if grant.ExpiresAt > 0 && now >= grant.ExpiresAt {
		return false
	}
	for _, granted := range grant.Actions {
		if granted == string(action) {
			return true
		}
	}
	return false
}

func LinkAllows(link model.ShareLink, action Action, now int64) bool {
	if link.Revoked || now >= link.ExpiresAt {
		return false
	}
	if action == View {
		return true
	}
	return action == Download && link.AllowDownload
}
