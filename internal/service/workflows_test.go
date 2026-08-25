package service

import (
	"lawdrive/internal/model"
	"lawdrive/internal/storage"
	"path/filepath"
	"testing"
)

func testService(t *testing.T) (*Service, *storage.Store) {
	t.Helper()
	store, err := storage.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	return New(store), store
}

func TestWorkflowCaseSetup(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	created, err := service.CreateCase(model.CreateCaseCommand{ID: "case-1", Number: "2026-001", Title: "Supply dispute", Client: "Acme", CreatorID: "partner", CreatorRole: model.RolePartner, At: 100})
	if err != nil {
		t.Fatal(err)
	}
	created, err = service.AddMember(created.ID, "partner", "lawyer", model.RoleLawyer, 110)
	if err != nil {
		t.Fatal(err)
	}
	folder, err := service.CreateFolder(model.CreateFolderCommand{ID: "folder-1", CaseID: created.ID, Name: "Evidence", ActorID: "lawyer", At: 120})
	if err != nil || folder.Name != "Evidence" {
		t.Fatalf("folder=%+v err=%v", folder, err)
	}
}

func TestWorkflowUploadEditPreview(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	_, err := service.CreateCase(model.CreateCaseCommand{ID: "case-2", Number: "2026-002", Title: "Lease dispute", CreatorID: "lawyer", CreatorRole: model.RoleLawyer, At: 200})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.PutFolder(model.Folder{ID: "folder-2", CaseID: "case-2", Name: "Contracts", CreatedAt: 201}); err != nil {
		t.Fatal(err)
	}
	document, err := service.Upload(model.UploadCommand{ID: "doc-2", VersionID: "v2-1", CaseID: "case-2", FolderID: "folder-2", Name: "lease.docx", MediaType: "application/docx", ActorID: "lawyer", Kind: model.KindContract, Content: []byte("first"), At: 210})
	if err != nil {
		t.Fatal(err)
	}
	document, err = service.Edit(model.EditCommand{VersionID: "v2-2", DocumentID: document.ID, ActorID: "lawyer", Content: []byte("second"), At: 220})
	if err != nil || document.CurrentVersion != 2 {
		t.Fatalf("document=%+v err=%v", document, err)
	}
	entries, err := store.AuditEntries("case-2")
	if err != nil || len(entries) != 3 {
		t.Fatalf("entries=%d err=%v", len(entries), err)
	}
}

func TestWorkflowShareDownload(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	_, err := service.CreateCase(model.CreateCaseCommand{ID: "case-3", Number: "2026-003", Title: "Evidence review", CreatorID: "partner", CreatorRole: model.RolePartner, At: 300})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.PutFolder(model.Folder{ID: "folder-3", CaseID: "case-3", Name: "Evidence", CreatedAt: 301}); err != nil {
		t.Fatal(err)
	}
	_, err = service.Upload(model.UploadCommand{ID: "doc-3", VersionID: "v3-1", CaseID: "case-3", FolderID: "folder-3", Name: "photo.txt", MediaType: "text/plain", ActorID: "partner", Kind: model.KindEvidence, Content: []byte("evidence"), At: 310})
	if err != nil {
		t.Fatal(err)
	}
	link, err := service.Share(model.ShareCommand{ID: "share-3", DocumentID: "doc-3", ActorID: "partner", Token: "fixed-token", AllowDownload: true, ExpiresAt: 500, At: 320})
	if err != nil {
		t.Fatal(err)
	}
	content, err := service.Download(model.DownloadCommand{DocumentID: "doc-3", ActorID: "external", ShareToken: link.Token, At: 330})
	if err != nil || string(content) != "evidence" {
		t.Fatalf("content=%s err=%v", content, err)
	}
}
