package storage

import (
	"lawdrive/internal/model"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lawdrive.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	value := model.CaseMatter{ID: "case-reopen", Number: "2026-R", Title: "Reopen", Status: model.CaseOpen, Members: []model.Member{{UserID: "partner", Role: model.RolePartner}}}
	if err := store.PutCase(value); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	loaded, err := store.Case(value.ID)
	if err != nil || loaded.Number != value.Number {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
}
