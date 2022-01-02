package memory

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type versionRepository struct {
	items map[string]*feed.Item
}

func NewVersionRepository(items map[string]*feed.Item) feed.VersionRepository {
	return &versionRepository{
		items: items,
	}
}

func (r *versionRepository) Get(packageName string, version string) (*packages.Version, error) {
	item, ok := r.items[packageName]
	if !ok {
		return nil, nil
	}
	for _, v := range item.Package.Versions {
		if v.Version == version {
			return v, nil
		}
	}
	return nil, nil
}

func (r *versionRepository) List(packageName string, query *feed.ItemReadExpandPackage) ([]*packages.Version, error) {
	item, ok := r.items[packageName]
	versions := []*packages.Version{}

	if !ok {
		return versions, nil
	}
	latestVersion := ""
	if item.State != nil {
		latestVersion = item.State.LatestVersion
	}
	for _, v := range item.Package.Versions {
		isMatch := query.IsMatch(v.Version, latestVersion)
		if isMatch {
			versions = append(versions, v)
		}
	}
	return versions, nil
}

type state struct {
	modified map[string]*packages.Version
	removed  map[string]bool
	added    map[string]*packages.Version
}

func newState() *state {
	return &state{
		modified: make(map[string]*packages.Version),
		removed:  make(map[string]bool),
		added:    make(map[string]*packages.Version),
	}
}

func (r *versionRepository) Update(packageName string, command *feed.VersionUpdate) ([]*packages.Version, error) {
	updated := []*packages.Version{}
	item, ok := r.items[packageName]
	if !ok || item.Package == nil || command == nil {
		return updated, nil
	}
	state := newState()
	item.Package.Versions = updateVersions(command, state, item.Package.Versions)
	for _, m := range state.modified {
		updated = append(updated, m)
	}
	for _, a := range state.added {
		updated = append(updated, a)
	}
	return updated, nil
}

func ToVersion(add *feed.VersionAdd) *packages.Version {
	version := &packages.Version{
		Version: add.Version,
		Targets: []*packages.Target{},
	}
	for _, target := range add.Targets {
		version.Targets = append(version.Targets, ToTarget(target))
	}
	return version
}

func updateVersions(update *feed.VersionUpdate, state *state, versions []*packages.Version) []*packages.Version {
	versions = FilterVersions(update.Remove, state, versions)
	versions = ModifyVersions(update.Modify, state, versions)
	versions = AddVersions(update.Add, state, versions)
	return versions
}

func FilterVersions(remove []string, state *state, versions []*packages.Version) []*packages.Version {
	if len(remove) == 0 {
		return versions
	}
	n := 0
	for _, v := range versions {
		keep := true
		for _, r := range remove {
			if r == v.Version {
				keep = false
				state.removed[v.Version] = true
				break
			}
		}
		if keep {
			versions[n] = v
			n++
		}
	}
	return versions[:n]
}

func ModifyVersions(modify []*feed.VersionModify, state *state, versions []*packages.Version) []*packages.Version {
	if len(modify) == 0 {
		return versions
	}
	for _, v := range versions {
		for _, m := range modify {
			if v.Version != m.Version {
				continue
			}
			ModifyVersion(m, state, v)
		}
	}
	return versions
}

func AddVersions(add []*feed.VersionAdd, state *state, versions []*packages.Version) []*packages.Version {
	for _, a := range add {
		version := ToVersion(a)
		versions = append(versions, version)
		state.added[version.Version] = version
	}
	return versions
}

func ModifyVersion(modify *feed.VersionModify, state *state, version *packages.Version) {
	if modify.NewVersion != nil {
		version.Version = *modify.NewVersion
		state.modified[version.Version] = version
	}
	modified := false
	version.Targets, modified = UpdateTargets(modify.Targets, version.Targets)
	if modified {
		state.modified[version.Version] = version
	}
}

func ToTarget(add *feed.TargetAdd) *packages.Target {
	target := &packages.Target{
		Platform:     add.Platform,
		Architecture: add.Architecture,
		Tasks:        []*packages.Task{},
	}
	for _, task := range add.Tasks {
		target.Tasks = append(target.Tasks, ToTask(task))
	}
	return target
}

func ToTask(add *feed.TaskAdd) *packages.Task {
	return &packages.Task{
		Name:       add.Name,
		Properties: add.Properties,
	}
}

func UpdateTargets(update *feed.TargetUpdate, targets []*packages.Target) ([]*packages.Target, bool) {
	if update == nil {
		return targets, false
	}
	var removed bool
	var modified bool
	targets, removed = FilterTargets(update.Remove, targets)
	targets, modified = ModifyTargets(update.Modify, targets)
	targets = AddTargets(update.Add, targets)
	return targets, removed || modified || len(update.Add) > 0
}

func TargetIsMatch(criteria *feed.PlatformArchitectureCriteria, target *packages.Target) bool {
	return criteria.Platform == target.Platform &&
		(criteria.Architecture == target.Architecture ||
			criteria.Architecture == "")
}

func FilterTargets(remove []*feed.PlatformArchitectureCriteria, targets []*packages.Target) ([]*packages.Target, bool) {
	if len(remove) == 0 {
		return targets, false
	}

	n := 0
	modified := false
	for _, x := range targets {
		keep := true
		for _, r := range remove {
			if TargetIsMatch(r, x) {
				keep = false
				break
			}
		}
		if keep {
			targets[n] = x
			n++
		} else {
			modified = true
		}
	}
	return targets[:n], modified
}

func ModifyTargets(modify []*feed.TargetModify, targets []*packages.Target) ([]*packages.Target, bool) {
	if len(modify) == 0 {
		return targets, false
	}
	modified := false
	for _, t := range targets {
		for _, m := range modify {
			if !TargetIsMatch(m.Criteria, t) {
				continue
			}
			modified = modified || ModifyTarget(m, t)
		}
	}
	return targets, modified
}

func AddTargets(add []*feed.TargetAdd, targets []*packages.Target) []*packages.Target {
	for _, a := range add {
		target := ToTarget(a)
		targets = append(targets, target)
	}
	return targets
}

func ModifyTarget(modify *feed.TargetModify, target *packages.Target) bool {
	modified := false
	if modify.NewArchitecture != nil {
		target.Architecture = *modify.NewArchitecture
		modified = true
	}
	if modify.NewPlatform != nil {
		target.Platform = *modify.NewArchitecture
		modified = true
	}
	return modified || PatchTasks(modify.Tasks, target.Tasks)
}

func PatchTasks(patches []*feed.TaskPatch, tasks []*packages.Task) bool {
	modified := false
	for _, p := range patches {
		switch p.Operation {
		case feed.PatchAdd:
			tasks = AddTask(p, tasks)
			modified = true

		case feed.PatchRemove:
			var removed = false
			tasks, removed = RemoveTask(p, tasks)
			modified = modified || removed

		case feed.PatchReplace:
			modified = modified || ReplaceTask(p, tasks)
		}
	}
	return modified
}

func AddTask(patch *feed.TaskPatch, tasks []*packages.Task) []*packages.Task {
	if patch.Operation != feed.PatchAdd {
		return tasks
	}
	task := ToTask(patch.Value)
	if patch.Index == nil {
		// add to end
		tasks = append(tasks, task)
		return tasks
	}

	// insert before
	index := *patch.Index
	tasks = append(tasks[:index+1], tasks[index:]...)
	tasks[index] = task
	return tasks
}

func RemoveTask(patch *feed.TaskPatch, tasks []*packages.Task) ([]*packages.Task, bool) {
	if patch.Operation != feed.PatchRemove || patch.Index == nil {
		return tasks, false
	}
	index := *patch.Index
	tasks = append(tasks[:index], tasks[index+1:]...)
	return tasks, true
}

func ReplaceTask(patch *feed.TaskPatch, tasks []*packages.Task) bool {
	if patch.Operation != feed.PatchReplace {
		return false
	}
	if patch.Index == nil || patch.Value == nil {
		return false
	}
	index := *patch.Index
	tasks[index] = ToTask(patch.Value)
	return true
}
