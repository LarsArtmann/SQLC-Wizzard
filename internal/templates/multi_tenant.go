package templates

// MultiTenantTemplate generates sqlc config for multi-tenant SaaS applications.
type MultiTenantTemplate struct {
	ConfiguredTemplate
}

// NewMultiTenantTemplate creates a new multi-tenant template.
func NewMultiTenantTemplate() *MultiTenantTemplate {
	t := &MultiTenantTemplate{}
	t.ConfiguredTemplate = ConfiguredTemplate{
		TemplateName:        "multi-tenant",
		TemplateDescription: "Optimized for SaaS multi-tenant architecture with tenant isolation and strict safety rules",
		DefaultPackageName:  "multi-tenant",
		DefaultProjectName:  "multi-tenant-app",
		StrictMode:          true,
		ProjectType:         "multi-tenant",
		DbEngine:            "postgresql",
		PackagePath:         "internal/db",
		BaseOutput:          "internal/db",
		UseManaged:          true,
		UseUUIDs:            true,
		UseJSON:             true,
		UseArrays:           true,
		UseFullText:         false,
		EmitJSONTags:        true,
		EmitPreparedQueries: true,
		EmitInterface:       true,
		EmitEmptySlices:     false,
		EmitResultStructPointers: true,
		EmitParamsStructPointers: true,
		EmitEnumValidMethod:      false,
		EmitAllEnumValues:        false,
		JSONTagsCaseStyle:        "camel",
		StrictFunctions:         true,
		StrictOrderBy:           true,
		NoSelectStar:             true,
		RequireWhere:             true,
		NoDropTable:             true,
		NoTruncate:              false,
		RequireLimit:            true,
		Features:                []string{"emit_interface", "prepared_queries", "json_tags", "tenant_isolation", "strict_checks"},
		CustomRenameRules: map[string]string{
			"id":     "ID",
			"uuid":   "UUID",
			"tenant": "Tenant",
			"url":    "URL",
			"uri":    "URI",
			"json":   "JSON",
			"api":    "API",
			"http":   "HTTP",
			"db":     "DB",
			"sql":    "SQL",
		},
	}
	return t
}

// RequiredFeatures returns which features this template requires.
func (t *MultiTenantTemplate) RequiredFeatures() []string {
	return []string{"emit_interface", "prepared_queries", "json_tags", "tenant_isolation", "strict_checks"}
}