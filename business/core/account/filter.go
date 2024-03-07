package account

type QueryFilter struct {
	Owner *string
}

func NewQueryFilter() QueryFilter {
	return QueryFilter{}
}

func (f *QueryFilter) WithOwner(owner string) {
	f.Owner = &owner
}
