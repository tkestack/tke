package storage

import "tkestack.io/tke/pkg/audit/storage/types"

type QueryParameter struct {
	Offset      int
	Size        int
	ClusterName string
	Namespace   string
	Resource    string
	Name        string
	StartTime   int64
	EndTime     int64
	UserName    string
	Query       string
}

type AuditStorage interface {
	Query(param *QueryParameter) ([]*types.Event, int, error)
	Save([]*types.Event) error
	// list option values for field
	FieldValues() map[string][]string

	Start()
	Stop()
}
