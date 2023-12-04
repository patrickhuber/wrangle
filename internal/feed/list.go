package feed

type ListRequest struct {
	Where  []*ItemReadAnyOf
	Expand *ItemReadExpand
}

type ItemReadAnyOf struct {
	AnyOf []*ItemReadAllOf
}

type ItemReadAllOf struct {
	AllOf []*ItemReadPredicate
}

type ItemReadPredicate struct {
	Name string
}

type ItemReadExpand struct {
	Package *ItemReadExpandPackage
}

type ItemReadExpandPackage struct {
	Where []*ItemReadExpandPackageAnyOf
}

type ItemReadExpandPackageAnyOf struct {
	AnyOf []*ItemReadExpandPackageAllOf
}

type ItemReadExpandPackageAllOf struct {
	AllOf []*ItemReadExpandPackagePredicate
}

type ItemReadExpandPackagePredicate struct {
	Latest  bool
	Version string
}

type ListResponse struct {
	Items []*Item
}

func (q *ItemReadExpandPackage) IsMatch(version, latest string) bool {
	if q == nil || len(q.Where) == 0 {
		return true
	}

	for _, any := range q.Where {
		for _, all := range any.AnyOf {
			allTrue := true
			for _, p := range all.AllOf {
				if p.Latest && version == latest {
					continue
				}
				if p.Version == version {
					continue
				}
				allTrue = false
				break
			}
			if allTrue {
				return true
			}
		}
	}
	return false
}

func IsMatch(q []*ItemReadAnyOf, name string) bool {
	if len(q) == 0 {
		return true
	}

	// filter out any folders that don't match the search criteria
	for _, any := range q {
		for _, all := range any.AnyOf {
			allTrue := true
			for _, p := range all.AllOf {
				if p.Name != name {
					allTrue = false
					break
				}
			}
			if allTrue {
				return true
			}
		}
	}
	return false
}
