/*
Copyright 2025 The KServe Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"fmt"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// Compile-time interface compliance checks.
var (
	_ conversion.Convertible = &LLMInferenceService{}
	_ conversion.Convertible = &LLMInferenceServiceConfig{}
)

// ConvertTo converts this LLMInferenceService (v1alpha2) to the Hub version (v1beta1).
// Uses reflection-based field copying to avoid JSON round-trip issues with apis.URL
// (knative.dev/pkg/apis) which URL-encodes '/' in Host to '%2F' during Marshal and
// fails to parse back during Unmarshal (Go url.Parse rejects %2F in hostnames).
// The Hub type lives in pkg/apis/serving/v1beta1; we avoid a direct Go-type import
// here to break the v1alpha2→v1beta1→v1alpha1 import cycle caused by
// v1beta1/configmap.go pulling in v1alpha1.KernelCachePath.
func (src *LLMInferenceService) ConvertTo(dstRaw conversion.Hub) error {
	return reflectFieldCopy(src, dstRaw)
}

// ConvertFrom converts the Hub version (v1beta1) to this LLMInferenceService (v1alpha2).
func (dst *LLMInferenceService) ConvertFrom(srcRaw conversion.Hub) error {
	return reflectFieldCopy(srcRaw, dst)
}

// ConvertTo converts this LLMInferenceServiceConfig (v1alpha2) to the Hub version (v1beta1).
func (src *LLMInferenceServiceConfig) ConvertTo(dstRaw conversion.Hub) error {
	return reflectFieldCopy(src, dstRaw)
}

// ConvertFrom converts the Hub version (v1beta1) to this LLMInferenceServiceConfig (v1alpha2).
func (dst *LLMInferenceServiceConfig) ConvertFrom(srcRaw conversion.Hub) error {
	return reflectFieldCopy(srcRaw, dst)
}

// reflectFieldCopy copies all exported fields from src to dst by matching field names.
// Both src and dst must be pointers to structs. Structurally compatible fields (same name,
// same or assignable type) are copied recursively. Fields that don't exist in dst are
// silently skipped. Embedded struct fields are flattened so wrapper types (like test
// stubs that embed the Hub type) work correctly. This preserves values like apis.URL
// exactly as-is, avoiding the JSON marshal/unmarshal round-trip bug in
// knative.dev/pkg/apis.URL where '/' in Host gets URL-encoded to '%2F' and fails
// Go's url.Parse during UnmarshalJSON.
func reflectFieldCopy(src, dst interface{}) error {
	return reflectFieldCopyValues(reflect.ValueOf(src), reflect.ValueOf(dst))
}

// structField represents a single field in a struct, including those promoted from
// embedded structs. value is settable and represents the actual field location.
type structField struct {
	name  string
	value reflect.Value
	typ   reflect.Type
}

// flattenStructFields returns all exported fields of a struct, flattening embedded
// (anonymous) struct fields so promoted fields are accessible by name.
func flattenStructFields(val reflect.Value, depth int) []structField {
	if depth > 4 {
		return nil
	}
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	var fields []structField
	t := val.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		fv := val.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			fields = append(fields, flattenStructFields(fv, depth+1)...)
		} else {
			fields = append(fields, structField{name: f.Name, value: fv, typ: f.Type})
		}
	}
	return fields
}

func reflectFieldCopyValues(srcVal, dstVal reflect.Value) error {
	// Dereference pointers
	for srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return nil
		}
		srcVal = srcVal.Elem()
	}
	for dstVal.Kind() == reflect.Ptr {
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstVal.Type().Elem()))
		}
		dstVal = dstVal.Elem()
	}

	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		if srcVal.Type().AssignableTo(dstVal.Type()) {
			dstVal.Set(srcVal)
			return nil
		}
		return fmt.Errorf("cannot copy %s to %s", srcVal.Type(), dstVal.Type())
	}

	// Flatten both src and dst to handle embedded struct wrappers
	srcFields := flattenStructFields(srcVal, 0)
	dstFields := flattenStructFields(dstVal, 0)

	// Build dst lookup map
	dstByName := make(map[string]structField, len(dstFields))
	for _, df := range dstFields {
		if _, exists := dstByName[df.name]; !exists {
			dstByName[df.name] = df
		}
	}

	for _, sf := range srcFields {
		df, ok := dstByName[sf.name]
		if !ok {
			continue
		}
		if err := copySingleField(sf.value, df.value, sf.typ, df.typ); err != nil {
			return fmt.Errorf("field %s: %w", sf.name, err)
		}
	}
	return nil
}

func copySingleField(src, dst reflect.Value, srcType, dstType reflect.Type) error {
	// Identical type — direct assignment
	if srcType == dstType {
		dst.Set(src)
		return nil
	}

	// Pointer types
	if srcType.Kind() == reflect.Ptr && dstType.Kind() == reflect.Ptr {
		if src.IsNil() {
			dst.Set(reflect.Zero(dstType))
			return nil
		}
		if dst.IsNil() {
			dst.Set(reflect.New(dstType.Elem()))
		}
		return copySingleField(src.Elem(), dst.Elem(), srcType.Elem(), dstType.Elem())
	}

	// Slices
	if srcType.Kind() == reflect.Slice && dstType.Kind() == reflect.Slice {
		if src.IsNil() {
			dst.Set(reflect.Zero(dstType))
			return nil
		}
		newSlice := reflect.MakeSlice(dstType, src.Len(), src.Cap())
		for i := range src.Len() {
			if err := copySingleField(src.Index(i), newSlice.Index(i), srcType.Elem(), dstType.Elem()); err != nil {
				return err
			}
		}
		dst.Set(newSlice)
		return nil
	}

	// Maps
	if srcType.Kind() == reflect.Map && dstType.Kind() == reflect.Map {
		if src.IsNil() {
			dst.Set(reflect.Zero(dstType))
			return nil
		}
		newMap := reflect.MakeMapWithSize(dstType, src.Len())
		iter := src.MapRange()
		for iter.Next() {
			newKey := reflect.New(dstType.Key()).Elem()
			if err := copySingleField(iter.Key(), newKey, srcType.Key(), dstType.Key()); err != nil {
				return err
			}
			newVal := reflect.New(dstType.Elem()).Elem()
			if err := copySingleField(iter.Value(), newVal, srcType.Elem(), dstType.Elem()); err != nil {
				return err
			}
			newMap.SetMapIndex(newKey, newVal)
		}
		dst.Set(newMap)
		return nil
	}

	// Nested structs
	if srcType.Kind() == reflect.Struct && dstType.Kind() == reflect.Struct {
		return reflectFieldCopyValues(src, dst)
	}

	// Fallback: assignable types
	if srcType.AssignableTo(dstType) {
		dst.Set(src)
		return nil
	}

	// Incompatible types (e.g. different packages): skip silently
	// Special cases like InferencePoolSpec must be handled explicitly by callers
	return nil
}
