package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

// UpdateRequest contains all update operations for a item and entities in its hierarchy
type UpdateRequest struct {
	Items []*ItemUpdate
}

type UpdateResponse struct {
	Changed int
}

// ItemUpdate updates an item for the given package name
type ItemUpdate struct {
	Name     string
	State    *StateUpdate
	Template *string
	Package  *PackageUpdate
}

// StateUpdate updates the state for a given item
type StateUpdate struct {
	LatestVersion *string
}

// PackageUpdate updates the given package
type PackageUpdate struct {
	Name     string
	NewName  *string // for rename
	Versions *VersionUpdate
}

// VersionUpdate updates the package version for a given package
type VersionUpdate struct {
	Add    []*VersionAdd
	Remove []string
	Modify []*VersionModify
}

type VersionModify struct {
	Version    string
	NewVersion *string
	Targets    *TargetUpdate
}

type VersionAdd struct {
	Version string
	Targets []*TargetAdd
}

type TargetUpdate struct {
	Add    []*TargetAdd
	Remove []*PlatformArchitectureCriteria
	Modify []*TargetModify
}

type TargetModify struct {
	Criteria        *PlatformArchitectureCriteria
	NewPlatform     *string
	NewArchitecture *string
	Tasks           []*TaskPatch
}

type PlatformArchitectureCriteria struct {
	Platform     string
	Architecture string
}

func (c *PlatformArchitectureCriteria) IsMatch(target *packages.Target) bool {
	return c.Architecture == target.Architecture && c.Platform == target.Architecture
}

type TargetAdd struct {
	Platform     string
	Architecture string
	Tasks        []*TaskAdd
}

type PatchOperation int

const (
	PatchAdd     PatchOperation = 0
	PatchRemove  PatchOperation = 1
	PatchReplace PatchOperation = 2
)

type TaskPatch struct {
	Index     int
	Operation PatchOperation
	Value     *TaskAdd
}

type TaskAdd struct {
	Name       string
	Properties map[string]string
}
