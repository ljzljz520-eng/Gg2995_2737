package authz

import (
	"lawdrive/internal/model"
	"testing"
)

func TestRolePolicy(t *testing.T) {
	if !Decide(model.RoleLawyer, Download, model.CaseOpen).Allowed {
		t.Fatal("lawyer should download")
	}
	if Decide(model.RoleAssistant, Share, model.CaseOpen).Allowed {
		t.Fatal("assistant should not share")
	}
	if Decide(model.RolePartner, Edit, model.CaseClosed).Allowed {
		t.Fatal("closed case should be read only")
	}
}
