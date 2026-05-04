// pkg/parser/directive/validate.go
package directive

import (
	"fmt"
	"strings"

	"github.com/arkannsk/elval/pkg/errs"
	ann "github.com/arkannsk/elval/pkg/parser/annotations"
)

// Validate is the main entry point for static validation of a directive.
func Validate(ft FieldInfo, dir ann.Directive, loc errs.Location) []errs.Diagnostic {
	//  Custom directives (x-*) are always valid — skip validation
	if strings.HasPrefix(dir.Type, "x-") {
		return nil
	}

	dType := Type(dir.Type)

	validator, exists := Registry[dType]
	if !exists {
		available := make([]string, 0, len(Registry))
		for k := range Registry {
			available = append(available, string(k))
		}
		return []errs.Diagnostic{warn(loc, dir.Type, ft.StructName, ft.FieldName,
			fmt.Sprintf("unknown directive '%s'", dir.Type),
			fmt.Sprintf("available: %s", strings.Join(available, ", ")))}
	}

	// 3. Run the specific validator function — pass dirType explicitly
	diags := validator(dType, ft, dir.Params, loc)

	// 4. Add deprecation warning if applicable
	if meta := GetInfo(dType); meta.Deprecated {
		diags = append(diags, warn(loc, dir.Type, ft.StructName, ft.FieldName,
			fmt.Sprintf("directive '%s' is deprecated: %s", dir.Type, meta.Description),
			meta.Example))
	}

	return diags
}

// HasErrors checks if any diagnostic has SeverityError.
func HasErrors(diags []errs.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == errs.SeverityError || d.Severity == errs.SeverityWarning {
			return true
		}
	}
	return false
}

// FilterBySeverity returns only diagnostics matching the given severity.
func FilterBySeverity(diags []errs.Diagnostic, sev errs.Severity) []errs.Diagnostic {
	result := make([]errs.Diagnostic, 0)
	for _, d := range diags {
		if d.Severity == sev {
			result = append(result, d)
		}
	}
	return result
}
