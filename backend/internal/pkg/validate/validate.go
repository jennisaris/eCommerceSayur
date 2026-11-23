package validate

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	slugRegex  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	phoneRegex = regexp.MustCompile(`^[0-9+\-\s()]{8,25}$`)
)

func New() *validator.Validate {
	v := validator.New()

	// Use JSON field names in validation errors instead of internal Go struct field names.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	_ = v.RegisterValidation("safe_url", validateSafeURL)
	_ = v.RegisterValidation("slug", validateSlug)
	_ = v.RegisterValidation("phone", validatePhone)

	return v
}

func validateSafeURL(fl validator.FieldLevel) bool {
	raw := strings.TrimSpace(fl.Field().String())
	if raw == "" {
		return true
	}
	for _, r := range raw {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return false
		}
	}
	// Allow site-relative asset paths like /images/... or /assets/..., but block protocol-relative //evil.com
	if strings.HasPrefix(raw, "/") {
		return !strings.HasPrefix(raw, "//")
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(parsed.Scheme)
	return (scheme == "http" || scheme == "https") && parsed.Host != ""
}

func validateSlug(fl validator.FieldLevel) bool {
	raw := fl.Field().String()
	return slugRegex.MatchString(raw)
}

func validatePhone(fl validator.FieldLevel) bool {
	raw := strings.TrimSpace(fl.Field().String())
	if !phoneRegex.MatchString(raw) {
		return false
	}
	digits := 0
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			digits++
		}
	}
	return digits >= 8 && digits <= 15
}

func FormatError(err error) string {
	if err == nil {
		return ""
	}
	var ve validator.ValidationErrors
	if ok := AsValidationErrors(err, &ve); ok && len(ve) > 0 {
		first := ve[0]
		field := first.Field()
		switch first.Tag() {
		case "required":
			return fmt.Sprintf("Field '%s' wajib diisi", field)
		case "email":
			return fmt.Sprintf("Format email pada '%s' tidak valid", field)
		case "min":
			return fmt.Sprintf("Field '%s' minimal %s karakter/nilai", field, first.Param())
		case "max":
			return fmt.Sprintf("Field '%s' maksimal %s karakter/nilai", field, first.Param())
		case "gt":
			return fmt.Sprintf("Field '%s' harus lebih besar dari %s", field, first.Param())
		case "gte":
			return fmt.Sprintf("Field '%s' harus minimal %s", field, first.Param())
		case "lte":
			return fmt.Sprintf("Field '%s' maksimal %s", field, first.Param())
		case "oneof":
			return fmt.Sprintf("Nilai '%s' tidak diizinkan", field)
		case "safe_url":
			return fmt.Sprintf("URL pada '%s' harus berupa http(s):// atau path lokal diawali '/'", field)
		case "slug":
			return fmt.Sprintf("Slug pada '%s' hanya boleh huruf kecil, angka, dan tanda hubung (-)", field)
		case "phone":
			return fmt.Sprintf("Format nomor telepon pada '%s' tidak valid", field)
		case "latitude", "longitude":
			return fmt.Sprintf("Koordinat pada '%s' tidak valid", field)
		default:
			return fmt.Sprintf("Nilai field '%s' tidak valid", field)
		}
	}
	return "Data input tidak valid"
}

func AsValidationErrors(err error, target *validator.ValidationErrors) bool {
	if ve, ok := err.(validator.ValidationErrors); ok {
		*target = ve
		return true
	}
	return false
}
