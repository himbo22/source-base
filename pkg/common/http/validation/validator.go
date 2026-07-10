package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Use json tag instead of struct field name.
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return field.Name
		}
		return name
	})
}

// Validate validates any request struct. Safe to call with non-struct input;
// returns an error instead of panicking.
func Validate(obj any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid validation target: %v", r)
		}
	}()

	if rv := reflect.ValueOf(obj); !isValidatable(rv) {
		return fmt.Errorf("validation target must be a struct or pointer to struct, got %T", obj)
	}

	if verr := validate.Struct(obj); verr != nil {
		if errs, ok := verr.(validator.ValidationErrors); ok {
			return buildValidationError(errs)
		}
		return verr
	}
	return nil
}

func isValidatable(rv reflect.Value) bool {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	return rv.Kind() == reflect.Struct
}

func buildValidationError(errs validator.ValidationErrors) error {
	// Return first error only. Can easily change to multiple errors later.
	fe := errs[0]
	field := fieldPath(fe)

	switch fe.Tag() {
	case "required":
		return fmt.Errorf("%s is required", field)
	case "email":
		return fmt.Errorf("%s must be a valid email", field)
	case "min":
		return fmt.Errorf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Errorf("%s must not exceed %s characters", field, fe.Param())
	case "len":
		return fmt.Errorf("%s must be exactly %s characters", field, fe.Param())
	case "oneof":
		return fmt.Errorf("%s must be one of [%s]", field, fe.Param())
	case "gte":
		return fmt.Errorf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Errorf("%s must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Errorf("%s must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Errorf("%s must be less than %s", field, fe.Param())
	case "uuid":
		return fmt.Errorf("%s must be a valid uuid", field)
	case "url":
		return fmt.Errorf("%s must be a valid url", field)
	case "datetime":
		return fmt.Errorf("%s must match datetime format %s", field, fe.Param())
	case "alpha":
		return fmt.Errorf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Errorf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Errorf("%s must be numeric", field)
	case "unique":
		return fmt.Errorf("%s must contain unique values", field)
	default:
		return fmt.Errorf("%s is invalid (%s)", field, fe.Tag())
	}
}

// fieldPath strips the root struct name from Namespace(), leaving a dotted
// path like "address.city" for nested fields, or just "email" for top-level ones.
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	parts := strings.SplitN(ns, ".", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return fe.Field()
}
