package repository

type Querier interface{}

var _ Querier = (*Queries)(nil)
