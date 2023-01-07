// Copyright 2019-2023 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

package validation

import (
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterAttributeValueValidation(v *validator.Validate, trans ut.Translator) error {
	return registerValidation(v, trans, "attributeValue", "{0} must be a string, integer, boolean or null value", func(fl validator.FieldLevel) bool {
		switch fl.Field().Kind() {
		// json.Number is internally represented as a string.
		case reflect.String:
			return true
		case reflect.Bool:
			return true
		case reflect.Interface:
			return fl.Field().IsNil()
		case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64,
			reflect.Float32, reflect.Float64, reflect.Func, reflect.Int, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Int8, reflect.Invalid,
			reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Uint,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr,
			reflect.UnsafePointer:

			return false
		default:
			return false
		}
	})
}
