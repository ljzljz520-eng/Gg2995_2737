package onlyoffice

import (
	"fmt"
	"lawdrive/internal/model"
)

type UserContext struct{ ID, Name string }
type Permissions struct{ Edit, Download, Print, Review, Comment bool }
type DocumentConfig struct {
	Key, Title, URL, FileType, DocumentType, CallbackURL string
	User                                                 UserContext
	Permissions                                          Permissions
}

func BuildConfig(document model.Document, version model.DocumentVersion, user model.User, role model.Role, baseURL string) (DocumentConfig, error) {
	if document.ID == "" || version.ID == "" || user.ID == "" {
		return DocumentConfig{}, fmt.Errorf("editor context is incomplete")
	}
	if baseURL == "" {
		return DocumentConfig{}, fmt.Errorf("document base url is required")
	}
	fileType := extension(document.Name)
	if fileType == "" {
		return DocumentConfig{}, fmt.Errorf("document extension is required")
	}
	permissions := Permissions{Comment: true}
	switch role {
	case model.RolePartner, model.RoleLawyer:
		permissions.Edit = true
		permissions.Download = true
		permissions.Print = true
		permissions.Review = true
	case model.RoleAssistant:
		permissions.Edit = true
		permissions.Review = true
	case model.RoleGuest:
	default:
		return DocumentConfig{}, fmt.Errorf("unsupported editor role %q", role)
	}
	return DocumentConfig{Key: version.ID, Title: document.Name, URL: baseURL + "/content/" + document.ID, CallbackURL: baseURL + "/callbacks/" + document.ID, FileType: fileType, DocumentType: documentType(fileType), User: UserContext{ID: user.ID, Name: user.DisplayName}, Permissions: permissions}, nil
}

func extension(name string) string {
	for index := len(name) - 1; index >= 0; index-- {
		if name[index] == '.' && index+1 < len(name) {
			return name[index+1:]
		}
	}
	return ""
}
func documentType(extension string) string {
	switch extension {
	case "doc", "docx", "odt", "txt":
		return "word"
	case "xls", "xlsx", "ods", "csv":
		return "cell"
	case "ppt", "pptx", "odp":
		return "slide"
	default:
		return "word"
	}
}
