package account

type NewAccount struct {
	Owner    string
	Currency string
}

type QueryFilter struct {
	Owner *string
}

type QueryLimiter struct {
	PageId   int32
	PageSize int32
}
