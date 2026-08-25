package model

import "fmt"

type CreateCaseCommand struct {
	ID, Number, Title, Client, CreatorID string
	CreatorRole                          Role
	At                                   int64
}
type CreateFolderCommand struct {
	ID, CaseID, ParentID, Name, ActorID string
	At                                  int64
}
type UploadCommand struct {
	ID, VersionID, CaseID, FolderID, Name, MediaType, ActorID string
	Kind                                                      DocumentKind
	Content                                                   []byte
	At                                                        int64
}
type EditCommand struct {
	VersionID, DocumentID, ActorID string
	Content                        []byte
	At                             int64
}
type ShareCommand struct {
	ID, DocumentID, ActorID, Token string
	AllowDownload                  bool
	ExpiresAt, At                  int64
}
type DownloadCommand struct {
	DocumentID, ActorID, ShareToken string
	At                              int64
}

func (c CreateCaseCommand) Validate() error {
	if c.ID == "" || c.Number == "" || c.CreatorID == "" {
		return fmt.Errorf("create case identity is incomplete")
	}
	if c.CreatorRole != RolePartner && c.CreatorRole != RoleLawyer {
		return fmt.Errorf("role %s cannot create cases", c.CreatorRole)
	}
	if c.At <= 0 {
		return fmt.Errorf("creation time is required")
	}
	return nil
}
func (c UploadCommand) Validate() error {
	if c.ID == "" || c.VersionID == "" || c.CaseID == "" || c.ActorID == "" {
		return fmt.Errorf("upload identity is incomplete")
	}
	if len(c.Content) == 0 {
		return fmt.Errorf("upload content is empty")
	}
	if c.At <= 0 {
		return fmt.Errorf("upload time is required")
	}
	return nil
}
func (c EditCommand) Validate() error {
	if c.DocumentID == "" || c.VersionID == "" || c.ActorID == "" {
		return fmt.Errorf("edit identity is incomplete")
	}
	if len(c.Content) == 0 {
		return fmt.Errorf("edit content is empty")
	}
	return nil
}
