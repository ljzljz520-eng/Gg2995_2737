package authz

import (
	"fmt"
	"lawdrive/internal/model"
)

type Action string

const (
	View     Action = "view"
	Edit     Action = "edit"
	Download Action = "download"
	Share    Action = "share"
	Manage   Action = "manage"
)

type Decision struct {
	Allowed bool
	Reason  string
}

func Decide(role model.Role, action Action, status model.CaseStatus) Decision {
	if status == model.CaseClosed && (action == Edit || action == Share) {
		return Decision{Reason: "closed cases are read only"}
	}
	switch role {
	case model.RolePartner:
		return Decision{Allowed: true, Reason: "partner policy"}
	case model.RoleLawyer:
		if action == Manage {
			return Decision{Reason: "case management requires partner"}
		}
		return Decision{Allowed: true, Reason: "lawyer policy"}
	case model.RoleAssistant:
		if action == View || action == Edit {
			return Decision{Allowed: true, Reason: "assistant collaboration policy"}
		}
		return Decision{Reason: "assistant cannot download or share"}
	case model.RoleGuest:
		if action == View {
			return Decision{Allowed: true, Reason: "guest preview policy"}
		}
		return Decision{Reason: "guest is preview only"}
	default:
		return Decision{Reason: fmt.Sprintf("unknown role %q", role)}
	}
}

func Require(c model.CaseMatter, userID string, action Action) error {
	role, ok := model.RoleFor(c, userID)
	if !ok {
		return fmt.Errorf("user %s is not a case member", userID)
	}
	decision := Decide(role, action, c.Status)
	if !decision.Allowed {
		return fmt.Errorf("%s: %s", action, decision.Reason)
	}
	return nil
}
