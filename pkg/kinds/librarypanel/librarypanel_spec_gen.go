// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     GoResourceTypes
//
// Run 'make gen-cue' from repository root to regenerate.

// Code generated - EDITING IS FUTILE. DO NOT EDIT.

package librarypanel

import (
	time "time"
)

type Spec struct {
	// Folder UID
	FolderUid *string `json:"folderUid,omitempty"`
	// Library element UID
	Uid string `json:"uid"`
	// Panel name (also saved in the model)
	Name string `json:"name"`
	// Panel description
	Description *string `json:"description,omitempty"`
	// The panel type (from inside the model)
	Type string `json:"type"`
	// Dashboard version when this was saved (zero if unknown)
	SchemaVersion *uint16 `json:"schemaVersion,omitempty"`
	// panel version, incremented each time the dashboard is updated.
	Version int64 `json:"version"`
	// TODO: should be the same panel schema defined in dashboard
	// Typescript: Omit<Panel, 'gridPos' | 'id' | 'libraryPanel'>;
	Model map[string]any `json:"model"`
	// Object storage metadata
	Meta *LibraryElementDTOMeta `json:"meta,omitempty"`
}

// NewSpec creates a new Spec object.
func NewSpec() *Spec {
	return &Spec{}
}

type LibraryElementDTOMeta struct {
	FolderName          string                    `json:"folderName"`
	FolderUid           string                    `json:"folderUid"`
	ConnectedDashboards int64                     `json:"connectedDashboards"`
	Created             time.Time                 `json:"created"`
	Updated             time.Time                 `json:"updated"`
	CreatedBy           LibraryElementDTOMetaUser `json:"createdBy"`
	UpdatedBy           LibraryElementDTOMetaUser `json:"updatedBy"`
}

// NewLibraryElementDTOMeta creates a new LibraryElementDTOMeta object.
func NewLibraryElementDTOMeta() *LibraryElementDTOMeta {
	return &LibraryElementDTOMeta{
		CreatedBy: *NewLibraryElementDTOMetaUser(),
		UpdatedBy: *NewLibraryElementDTOMetaUser(),
	}
}

type LibraryElementDTOMetaUser struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}

// NewLibraryElementDTOMetaUser creates a new LibraryElementDTOMetaUser object.
func NewLibraryElementDTOMetaUser() *LibraryElementDTOMetaUser {
	return &LibraryElementDTOMetaUser{}
}
