package lawdrive

import (
	"fmt"
	"lawdrive/internal/model"
	"lawdrive/internal/service"
	"lawdrive/internal/storage"
	"path/filepath"
	"testing"
)

type quotaPool struct{ active, limit int }
type quotaResource struct {
	pool     *quotaPool
	released bool
}

func (p *quotaPool) Acquire(string) (service.Resource, error) {
	if p.active >= p.limit {
		return nil, fmt.Errorf("quota exhausted")
	}
	p.active++
	return &quotaResource{pool: p}, nil
}
func (r *quotaResource) Release() {
	if !r.released {
		r.released = true
		r.pool.active--
	}
}

func TestBusinessChain06(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "batch.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	svc := service.New(store)
	_, err = svc.CreateCase(model.CreateCaseCommand{ID: "case-b", Number: "B-1", Title: "Batch", CreatorID: "partner", CreatorRole: model.RolePartner, At: 10})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.PutFolder(model.Folder{ID: "folder-b", CaseID: "case-b", Name: "Export", CreatedAt: 11}); err != nil {
		t.Fatal(err)
	}
	for index := 1; index <= 3; index++ {
		id := fmt.Sprintf("doc-b-%d", index)
		_, err = svc.Upload(model.UploadCommand{ID: id, VersionID: fmt.Sprintf("v-b-%d", index), CaseID: "case-b", FolderID: "folder-b", Name: id + ".docx", MediaType: "application/docx", ActorID: "partner", Kind: model.KindEvidence, Content: []byte(id), At: int64(20 + index)})
		if err != nil {
			t.Fatal(err)
		}
	}
	result, err := svc.ExportBatch([]string{"doc-b-1", "doc-b-2", "doc-b-3"}, &quotaPool{limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("got %d exports", len(result))
	}
}
