package service

import (
	"fmt"
	"lawdrive/internal/audit"
	"lawdrive/internal/authz"
	"lawdrive/internal/model"
	"lawdrive/internal/storage"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) CreateCase(command model.CreateCaseCommand) (model.CaseMatter, error) {
	if err := command.Validate(); err != nil {
		return model.CaseMatter{}, err
	}
	value := model.CaseMatter{ID: command.ID, Number: command.Number, Title: command.Title, Client: command.Client, Status: model.CaseOpen, Members: []model.Member{{UserID: command.CreatorID, Role: command.CreatorRole}}, CreatedAt: command.At}
	if err := s.store.PutCase(value); err != nil {
		return model.CaseMatter{}, err
	}
	entry, _ := audit.New("audit-case-"+command.ID, value.ID, command.CreatorID, "case.created", value.ID, map[string]string{"number": value.Number}, command.At)
	if err := s.store.AppendAudit(entry); err != nil {
		return model.CaseMatter{}, err
	}
	return value, nil
}

func (s *Service) AddMember(caseID, actorID, userID string, role model.Role, at int64) (model.CaseMatter, error) {
	value, err := s.store.Case(caseID)
	if err != nil {
		return model.CaseMatter{}, err
	}
	if err := authz.Require(value, actorID, authz.Manage); err != nil {
		return model.CaseMatter{}, err
	}
	if _, exists := model.RoleFor(value, userID); exists {
		return model.CaseMatter{}, fmt.Errorf("user %s is already a member", userID)
	}
	value.Members = append(value.Members, model.Member{UserID: userID, Role: role})
	if err := s.store.UpdateCase(value); err != nil {
		return model.CaseMatter{}, err
	}
	entry, _ := audit.New(fmt.Sprintf("audit-member-%s-%s", caseID, userID), caseID, actorID, "member.added", userID, map[string]string{"role": string(role)}, at)
	return value, s.store.AppendAudit(entry)
}

func (s *Service) CreateFolder(command model.CreateFolderCommand) (model.Folder, error) {
	value, err := s.store.Case(command.CaseID)
	if err != nil {
		return model.Folder{}, err
	}
	if err := authz.Require(value, command.ActorID, authz.Edit); err != nil {
		return model.Folder{}, err
	}
	folder := model.Folder{ID: command.ID, CaseID: command.CaseID, ParentID: command.ParentID, Name: command.Name, CreatedAt: command.At}
	if err := s.store.PutFolder(folder); err != nil {
		return model.Folder{}, err
	}
	entry, _ := audit.New("audit-folder-"+command.ID, command.CaseID, command.ActorID, "folder.created", command.ID, map[string]string{"name": command.Name}, command.At)
	return folder, s.store.AppendAudit(entry)
}
