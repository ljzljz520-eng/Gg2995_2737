package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lawdrive/internal/audit"
	"lawdrive/internal/authz"
	"lawdrive/internal/model"
)

func checksum(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func (s *Service) Upload(command model.UploadCommand) (model.Document, error) {
	if err := command.Validate(); err != nil {
		return model.Document{}, err
	}
	caseMatter, err := s.store.Case(command.CaseID)
	if err != nil {
		return model.Document{}, err
	}
	if err := authz.Require(caseMatter, command.ActorID, authz.Edit); err != nil {
		return model.Document{}, err
	}
	document := model.Document{ID: command.ID, CaseID: command.CaseID, FolderID: command.FolderID, Name: command.Name, Kind: command.Kind, MediaType: command.MediaType, CurrentVersion: 1, Status: "active", CreatedBy: command.ActorID, CreatedAt: command.At}
	version := model.DocumentVersion{ID: command.VersionID, DocumentID: command.ID, Sequence: 1, Checksum: checksum(command.Content), Content: append([]byte(nil), command.Content...), EditorID: command.ActorID, CreatedAt: command.At}
	if err := s.store.CreateDocument(document, version); err != nil {
		return model.Document{}, err
	}
	entry, _ := audit.New("audit-upload-"+command.ID, command.CaseID, command.ActorID, "document.uploaded", command.ID, map[string]string{"kind": string(command.Kind), "checksum": version.Checksum}, command.At)
	return document, s.store.AppendAudit(entry)
}

func (s *Service) Edit(command model.EditCommand) (model.Document, error) {
	if err := command.Validate(); err != nil {
		return model.Document{}, err
	}
	document, err := s.store.Document(command.DocumentID)
	if err != nil {
		return model.Document{}, err
	}
	caseMatter, err := s.store.Case(document.CaseID)
	if err != nil {
		return model.Document{}, err
	}
	if err := authz.Require(caseMatter, command.ActorID, authz.Edit); err != nil {
		return model.Document{}, err
	}
	version := model.DocumentVersion{ID: command.VersionID, DocumentID: document.ID, Sequence: document.CurrentVersion + 1, Checksum: checksum(command.Content), Content: append([]byte(nil), command.Content...), EditorID: command.ActorID, CreatedAt: command.At}
	updated, err := s.store.AddVersion(version)
	if err != nil {
		return model.Document{}, err
	}
	entry, _ := audit.New("audit-edit-"+command.VersionID, document.CaseID, command.ActorID, "document.edited", document.ID, map[string]string{"version": fmt.Sprint(version.Sequence), "checksum": version.Checksum}, command.At)
	return updated, s.store.AppendAudit(entry)
}
