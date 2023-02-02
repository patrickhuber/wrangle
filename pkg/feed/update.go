package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

// For an Array:
// create an {Entity}Update struct
//
//    type EntityUpdate struct{
//       Add []*Entity
//       Modify []*EntityModify
//       Remove []int // this is the index of the array
//    }
//
// For a Map:
// create an {Entity}Update struct
//
//    type EntityUpdate struct{
//       Add []*EntityAdd
//       Modify []*EntityModify
//       Remove []string // this is the key of the map
//    }
//
// For a Single Entity:
// create an {Entity}Modify struct
//
//    type EntityModify struct{
//	    Name string       // this is a key and is immutable
//      Action *string    // nil == unchanged, value == change
//    }
//
// For a Single Entity with Mutable key:
// create an {Entity}Modify struct
//
//    type EntityModify struct{
//      Name string     // this is used to look up the entity
//      NewName *string // this is used to rename the key
//    }

// UpdateRequest contains all update operations for a item and entities in its hierarchy
type UpdateRequest struct {
	Items *ItemUpdate
}

type UpdateResponse struct {
	Changed int
}

// ItemUpdate updates an item for the given package name
type ItemUpdate struct {
	Add    []*Item
	Modify []*ItemModify
}

type ItemModify struct {
	Name     string
	State    *StateModify
	Template *string
	Package  *PackageModify
}

// StateModify updates the state for a given item
type StateModify struct {
	LatestVersion *string
}

// PackageModify updates the given package
type PackageModify struct {
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
	Manifest   *ManifestModify
}

type VersionAdd struct {
	Version  string
	Manifest *ManifestAdd
}

type ManifestAdd struct {
	Package *ManifestPackageAdd
}

type ManifestModify struct {
	Package *ManifestPackageModify
}

type ManifestPackageAdd struct {
	Name    string
	Version string
	Targets []*ManifestTargetAdd
}

type ManifestPackageModify struct {
	NewName    *string
	NewVersion *string
	Targets    *ManifestTargetUpdate
}

type ManifestTargetUpdate struct {
	Add    []*ManifestTargetAdd
	Remove []*PlatformArchitectureCriteria
	Modify []*ManifestTargetModify
}

type ManifestTargetModify struct {
	Criteria        *PlatformArchitectureCriteria
	NewPlatform     *string
	NewArchitecture *string
	Steps           []*ManifestStepPatch
}

type PlatformArchitectureCriteria struct {
	Platform     string
	Architecture string
}

func (c *PlatformArchitectureCriteria) IsMatch(target *packages.ManifestTarget) bool {
	return c.Architecture == target.Architecture && c.Platform == target.Architecture
}

type ManifestTargetAdd struct {
	Platform     string
	Architecture string
	Steps        []*ManifestStepAdd
}

type PatchOperation int

const (
	PatchAdd     PatchOperation = 0
	PatchRemove  PatchOperation = 1
	PatchReplace PatchOperation = 2
)

type ManifestStepPatch struct {
	Index     int
	Operation PatchOperation
	Value     *ManifestStepAdd
}

type ManifestStepAdd struct {
	Action string
	With   map[string]string
}
