package model

import "testing"

func TestClassificationCatalog(t *testing.T) {
	items := ClassificationsByPhase("filing")
	if len(items) != 30 {
		t.Fatalf("got %d classifications", len(items))
	}
	deadline, err := RetentionDeadline(100, "filing-pleading-filed")
	if err != nil || deadline <= 100 {
		t.Fatalf("deadline=%d err=%v", deadline, err)
	}
}

func TestCaseValidation(t *testing.T) {
	value := CaseMatter{ID: "case-1", Number: "2026-001", Title: "Contract dispute", Status: CaseOpen, Members: []Member{{UserID: "partner", Role: RolePartner}}}
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
}
