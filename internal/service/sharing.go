package service

import (
	"fmt"
	"lawdrive/internal/audit"
	"lawdrive/internal/authz"
	"lawdrive/internal/model"
)

func (s *Service) Share(command model.ShareCommand) (model.ShareLink, error) {
	document, err := s.store.Document(command.DocumentID)
	if err != nil {
		return model.ShareLink{}, err
	}
	caseMatter, err := s.store.Case(document.CaseID)
	if err != nil {
		return model.ShareLink{}, err
	}
	if err := authz.Require(caseMatter, command.ActorID, authz.Share); err != nil {
		return model.ShareLink{}, err
	}
	if command.ExpiresAt <= command.At {
		return model.ShareLink{}, fmt.Errorf("share expiry must be in the future")
	}
	link := model.ShareLink{ID: command.ID, DocumentID: document.ID, Token: command.Token, AllowDownload: command.AllowDownload, ExpiresAt: command.ExpiresAt}
	if err := s.store.PutShare(link); err != nil {
		return model.ShareLink{}, err
	}
	entry, _ := audit.New("audit-share-"+command.ID, document.CaseID, command.ActorID, "document.shared", document.ID, map[string]string{"download": fmt.Sprint(command.AllowDownload)}, command.At)
	return link, s.store.AppendAudit(entry)
}

func (s *Service) Download(command model.DownloadCommand) ([]byte, error) {
	document, err := s.store.Document(command.DocumentID)
	if err != nil {
		return nil, err
	}
	caseMatter, err := s.store.Case(document.CaseID)
	if err != nil {
		return nil, err
	}
	allowed := authz.Require(caseMatter, command.ActorID, authz.Download) == nil
	if !allowed && command.ShareToken != "" {
		link, linkErr := s.store.ShareByToken(command.ShareToken)
		allowed = linkErr == nil && link.DocumentID == document.ID && authz.LinkAllows(link, authz.Download, command.At)
	}
	if !allowed {
		return nil, fmt.Errorf("download denied")
	}
	version, err := s.store.CurrentVersion(document.ID)
	if err != nil {
		return nil, err
	}
	entry, _ := audit.New(fmt.Sprintf("audit-download-%s-%d", document.ID, command.At), document.CaseID, command.ActorID, "document.downloaded", document.ID, map[string]string{"version": fmt.Sprint(version.Sequence)}, command.At)
	if err := s.store.AppendAudit(entry); err != nil {
		return nil, err
	}
	return append([]byte(nil), version.Content...), nil
}
