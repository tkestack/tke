package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/pkg/audit/apis/config"
)

func ValidateAuditConfiguration(ac *config.AuditConfiguration) error {
	if ac.Storage.ElasticSearch == nil {
		fld := field.NewPath("components").Child("elasticSearch")
		return field.Required(fld, "must be specified")
	}
	if ac.Storage.ElasticSearch.Address == "" {
		fld := field.NewPath("components").Child("elasticSearch").Child("address")
		return field.Required(fld, "must be specified")
	}
	return nil
}
