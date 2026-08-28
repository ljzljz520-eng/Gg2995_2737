package query

import (
	"lawdrive/internal/model"
	"lawdrive/internal/service"
	"lawdrive/internal/storage"
	"path/filepath"
	"testing"
)

func TestSearchAndDashboard(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	svc := service.New(store)
	_, err = svc.CreateCase(model.CreateCaseCommand{ID: "case-q", Number: "Q-1", Title: "Query", CreatorID: "lawyer", CreatorRole: model.RoleLawyer, At: 10})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.PutFolder(model.Folder{ID: "folder-q", CaseID: "case-q", Name: "Minutes", CreatedAt: 11}); err != nil {
		t.Fatal(err)
	}
	_, err = svc.Upload(model.UploadCommand{ID: "doc-q", VersionID: "vq-1", CaseID: "case-q", FolderID: "folder-q", Name: "meeting.docx", MediaType: "application/docx", ActorID: "lawyer", Kind: model.KindMinutes, Content: []byte("minutes"), At: 12})
	if err != nil {
		t.Fatal(err)
	}
	queries := New(store)
	items, err := queries.SearchDocuments("lawyer", model.SearchFilter{CaseID: "case-q", Query: "meeting"})
	if err != nil || len(items) != 1 || !items[0].CanDownload {
		t.Fatalf("items=%v err=%v", items, err)
	}
	dashboard, err := queries.Dashboard("lawyer", 20)
	if err != nil || dashboard.OpenCases != 1 || dashboard.Documents != 1 {
		t.Fatalf("dashboard=%+v err=%v", dashboard, err)
	}
}
