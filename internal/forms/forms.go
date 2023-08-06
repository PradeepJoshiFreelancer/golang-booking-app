package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// new strut to save the form data
type Form struct {
	url.Values
	Errors errors
}

// Creates a new instance of form struct
func NewForm(formData url.Values) *Form {
	return &Form{
		formData,
		errors(map[string][]string{}),
	}
}

// checks if the field exists in the form
func (f *Form) HasField(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// validates if the list of feilds are not spaces
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.AddError(field, "This field cannot be blank")
		}
	}
}

// validates if the fields min length.
func (f *Form) MinLength(field string, minLength int, r *http.Request) bool {
	value := f.Get(field)
	if len(strings.TrimSpace(value)) < minLength {
		f.Errors.AddError(field, fmt.Sprintf("This field length cannot be less than %d charecters", minLength))
		return false
	}
	return true
}

// validates email address field is valid uses govalidator package
func (f *Form) IsEmailValid(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.AddError(field, "Not a valid email address")

	}
}

// cheks if the form instance is valid
func (f *Form) IsFormValid() bool {
	return len(f.Errors) == 0
}
