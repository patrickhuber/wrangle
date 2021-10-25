package feed

// UpdateRequest contains all update operations for a item and entities in its hierarchy
type UpdateRequest struct {
	Items []*ItemUpdate
}

type UpdateResponse struct {
	Items []*Item
}

// ItemUpdate updates an item for the given package name
type ItemUpdate struct {
	Name      string
	State     *StateUpdate
	Template  *string
	Platforms *PlatformUpdate
	Package   *PackageUpdate
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
	Platform        string
	NewPlatform     *string
	Architecture    string
	NewArchitecture *string
	Tasks           *TaskUpdate
}

type PlatformArchitectureCriteria struct {
	Platform     string
	Architecture string
}

type TargetAdd struct {
	Platform     string
	Architecture string
	Tasks        []*TaskAdd
}

type TaskUpdate struct {
	Add    []*TaskAdd
	Remove []string
	Modify []*TaskModify
}

type TaskModify struct {
	Name       string
	NewName    *string
	Properties *StringMapUpdate
}

type TaskAdd struct {
	Name       string
	Properties map[string]string
}

type PlatformUpdate struct {
	Add    []*PlatformAdd
	Remove []string
	Modify []*PlatformModify
}

type PlatformModify struct {
	Name          string
	Architectures *UniqueStringListUpdate
}

type PlatformAdd struct {
	Name          string
	Architectures []string
}

type StringMapUpdate struct {
	Modify map[string]string
	Remove []string
}

type UniqueStringListUpdate struct {
	Add    []string
	Remove []string
}
