package model

import (
	"errors"
	"fmt"
	"strings"
)

type Role string

const (
	RolePartner   Role = "partner"
	RoleLawyer    Role = "lawyer"
	RoleAssistant Role = "assistant"
	RoleGuest     Role = "guest"
)

type DocumentKind string

const (
	KindContract DocumentKind = "contract"
	KindEvidence DocumentKind = "evidence"
	KindMinutes  DocumentKind = "minutes"
	KindPleading DocumentKind = "pleading"
)

type CaseStatus string

const (
	CaseOpen   CaseStatus = "open"
	CaseStayed CaseStatus = "stayed"
	CaseClosed CaseStatus = "closed"
)

type Member struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
}
type CaseMatter struct {
	ID        string     `json:"id"`
	Number    string     `json:"number"`
	Title     string     `json:"title"`
	Client    string     `json:"client"`
	Status    CaseStatus `json:"status"`
	Members   []Member   `json:"members"`
	CreatedAt int64      `json:"created_at"`
}
type Document struct {
	ID             string       `json:"id"`
	CaseID         string       `json:"case_id"`
	FolderID       string       `json:"folder_id"`
	Name           string       `json:"name"`
	Kind           DocumentKind `json:"kind"`
	MediaType      string       `json:"media_type"`
	CurrentVersion int          `json:"current_version"`
	Status         string       `json:"status"`
	CreatedBy      string       `json:"created_by"`
	CreatedAt      int64        `json:"created_at"`
}
type DocumentVersion struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Sequence   int    `json:"sequence"`
	Checksum   string `json:"checksum"`
	Content    []byte `json:"content"`
	EditorID   string `json:"editor_id"`
	CreatedAt  int64  `json:"created_at"`
}
type PermissionGrant struct {
	ID         string   `json:"id"`
	CaseID     string   `json:"case_id"`
	DocumentID string   `json:"document_id"`
	SubjectID  string   `json:"subject_id"`
	Actions    []string `json:"actions"`
	ExpiresAt  int64    `json:"expires_at"`
	GrantedBy  string   `json:"granted_by"`
}
type ShareLink struct {
	ID            string `json:"id"`
	DocumentID    string `json:"document_id"`
	Token         string `json:"token"`
	AllowDownload bool   `json:"allow_download"`
	ExpiresAt     int64  `json:"expires_at"`
	Revoked       bool   `json:"revoked"`
}
type AuditEntry struct {
	ID         string            `json:"id"`
	CaseID     string            `json:"case_id"`
	ActorID    string            `json:"actor_id"`
	Action     string            `json:"action"`
	TargetID   string            `json:"target_id"`
	Detail     map[string]string `json:"detail"`
	OccurredAt int64             `json:"occurred_at"`
}
type Folder struct {
	ID        string `json:"id"`
	CaseID    string `json:"case_id"`
	ParentID  string `json:"parent_id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}
type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Active      bool   `json:"active"`
}

func (c CaseMatter) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("case id is required")
	}
	if strings.TrimSpace(c.Number) == "" {
		return errors.New("case number is required")
	}
	if strings.TrimSpace(c.Title) == "" {
		return errors.New("case title is required")
	}
	if c.Status != CaseOpen && c.Status != CaseStayed && c.Status != CaseClosed {
		return fmt.Errorf("unsupported case status %q", c.Status)
	}
	if len(c.Members) == 0 {
		return errors.New("case needs at least one member")
	}
	seen := map[string]bool{}
	for _, member := range c.Members {
		if member.UserID == "" {
			return errors.New("member user id is required")
		}
		if seen[member.UserID] {
			return fmt.Errorf("duplicate member %s", member.UserID)
		}
		seen[member.UserID] = true
	}
	return nil
}

func (d Document) Validate() error {
	if d.ID == "" || d.CaseID == "" || d.FolderID == "" {
		return errors.New("document identity is incomplete")
	}
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("document name is required")
	}
	switch d.Kind {
	case KindContract, KindEvidence, KindMinutes, KindPleading:
	default:
		return fmt.Errorf("unsupported document kind %q", d.Kind)
	}
	if d.CurrentVersion < 1 {
		return errors.New("document needs a version")
	}
	return nil
}

func (v DocumentVersion) Validate() error {
	if v.ID == "" || v.DocumentID == "" {
		return errors.New("version identity is incomplete")
	}
	if v.Sequence < 1 {
		return errors.New("version sequence must be positive")
	}
	if len(v.Content) == 0 {
		return errors.New("empty document content")
	}
	if v.Checksum == "" {
		return errors.New("checksum is required")
	}
	return nil
}

func RoleFor(c CaseMatter, userID string) (Role, bool) {
	for _, member := range c.Members {
		if member.UserID == userID {
			return member.Role, true
		}
	}
	return "", false
}
