package onlyoffice

import (
	"lawdrive/internal/model"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	document := model.Document{ID: "doc", Name: "memo.docx"}
	version := model.DocumentVersion{ID: "version"}
	user := model.User{ID: "lawyer", DisplayName: "Counsel"}
	config, err := BuildConfig(document, version, user, model.RoleLawyer, "http://office.local")
	if err != nil {
		t.Fatal(err)
	}
	if !config.Permissions.Edit || config.DocumentType != "word" {
		t.Fatalf("config=%+v", config)
	}
}

func TestHandleCallback(t *testing.T) {
	result, err := HandleCallback(Callback{Status: 2, Key: "v2", UserID: "lawyer", Content: []byte("saved")})
	if err != nil || !result.Save || !result.Close {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}
