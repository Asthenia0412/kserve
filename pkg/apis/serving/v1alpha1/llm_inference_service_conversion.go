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

package v1alpha1

import (
	"encoding/json"
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	igwapiv1 "sigs.k8s.io/gateway-api-inference-extension/api/v1"

	igwapiv1alpha2 "github.com/kserve/kserve/pkg/apis/gie/v1alpha2pool"

	"github.com/kserve/kserve/pkg/apis/serving/v1alpha2"
)

const (
	// ModelCriticalityAnnotationKey stores the model criticality when converting to the Hub version.
	ModelCriticalityAnnotationKey = "internal.serving.kserve.io/model-criticality"
	// LoRACriticalitiesAnnotationKey stores LoRA adapter criticalities as JSON when converting to the Hub version.
	LoRACriticalitiesAnnotationKey = "internal.serving.kserve.io/lora-criticalities"
)

// Compile-time interface compliance checks.
var (
	_ conversion.Convertible = &LLMInferenceService{}
	_ conversion.Convertible = &LLMInferenceServiceConfig{}
)

// ConvertTo converts this LLMInferenceService (v1alpha1) to the Hub version (v1beta1).
// Uses reflection-based field copying for identity bulk copy (avoids JSON round-trip
// issues with apis.URL), then handles known type-mismatch edge cases and preserves
// v1alpha1-only Criticality fields in annotations. The Hub concrete type lives in
// pkg/apis/serving/v1beta1; we avoid a direct Go import to prevent an import cycle
// with v1beta1/configmap.go → v1alpha1.KernelCachePath.
func (src *LLMInferenceService) ConvertTo(dstRaw conversion.Hub) error {
	if err := reflectFieldCopy(src, dstRaw); err != nil {
		return err
	}
	// Edge case: InferencePoolSpec uses different GIE package types in v1alpha1 (fork)
	// vs Hub (official). Reflection copy skips incompatible types, so we set it explicitly.
	if src.Spec.Router != nil && src.Spec.Router.Scheduler != nil && src.Spec.Router.Scheduler.Pool != nil {
		if err := setHubPoolSpec(dstRaw, convertPoolSpecForHub(src.Spec.Router.Scheduler.Pool)); err != nil {
			return fmt.Errorf("setting InferencePoolSpec on Hub: %w", err)
		}
	}
	saveCriticalityToAnnotationsOnHub(dstRaw, &src.Spec.Model)
	return nil
}

// ConvertFrom converts the Hub version (v1beta1) to this LLMInferenceService (v1alpha1).
func (dst *LLMInferenceService) ConvertFrom(srcRaw conversion.Hub) error {
	if err := reflectFieldCopy(srcRaw, dst); err != nil {
		return err
	}
	// Edge case: InferencePoolSpec type mismatch (v1alpha1 fork vs Hub official GIE).
	// Reflection copy handles Pool.Ref correctly (same *corev1.LocalObjectReference type)
	// but cannot copy Pool.Spec because the GIE package types differ. We explicitly
	// convert and set Spec only when Hub has a non-nil Spec. If Hub only has Pool.Ref
	// (no Spec), the reflection copy already handled it and we must NOT overwrite
	// Pool with nil.
	if dst.Spec.Router != nil && dst.Spec.Router.Scheduler != nil && dst.Spec.Router.Scheduler.Pool != nil {
		hubPoolSpec := getHubPoolSpec(srcRaw)
		if hubPoolSpec != nil {
			dst.Spec.Router.Scheduler.Pool.Spec = convertInferencePoolSpecFromV1(hubPoolSpec)
		}
	}
	restoreCriticalityFromHubAnnotations(srcRaw, &dst.Spec.Model)
	cleanupCriticalityAnnotations(&dst.ObjectMeta)
	return nil
}

// ConvertTo converts this LLMInferenceServiceConfig (v1alpha1) to the Hub version (v1beta1).
func (src *LLMInferenceServiceConfig) ConvertTo(dstRaw conversion.Hub) error {
	if err := reflectFieldCopy(src, dstRaw); err != nil {
		return err
	}
	saveCriticalityToAnnotationsOnHub(dstRaw, &src.Spec.Model)
	return nil
}

// ConvertFrom converts the Hub version (v1beta1) to this LLMInferenceServiceConfig (v1alpha1).
func (dst *LLMInferenceServiceConfig) ConvertFrom(srcRaw conversion.Hub) error {
	if err := reflectFieldCopy(srcRaw, dst); err != nil {
		return err
	}
	restoreCriticalityFromHubAnnotations(srcRaw, &dst.Spec.Model)
	cleanupCriticalityAnnotations(&dst.ObjectMeta)
	return nil
}

// setHubPoolSpec navigates the Hub struct via reflection to set the Spec.Router.Scheduler.Pool.Spec
// field. We use reflection because we can't import the Hub's concrete type (v1beta1) here.
func setHubPoolSpec(hub conversion.Hub, poolSpec *igwapiv1.InferencePoolSpec) error {
	if poolSpec == nil {
		return nil
	}
	hubVal := reflect.ValueOf(hub)
	// Dereference interface and pointer
	for hubVal.Kind() == reflect.Interface || hubVal.Kind() == reflect.Ptr {
		if hubVal.IsNil() {
			return nil
		}
		hubVal = hubVal.Elem()
	}
	if hubVal.Kind() != reflect.Struct {
		return fmt.Errorf("Hub is not a struct: %s", hubVal.Type())
	}

	// Navigate: Spec.Router.Scheduler.Pool.Spec
	path := []string{"Spec", "Router", "Scheduler", "Pool", "Spec"}
	for i, fieldName := range path {
		isLast := i == len(path)-1
		found := findFieldInStruct(hubVal, fieldName)
		if !found.IsValid() {
			return fmt.Errorf("field %s not found in %s", fieldName, hubVal.Type())
		}
		hubVal = found
		// If it's a pointer, allocate if nil, then deref (except for the last one)
		if hubVal.Kind() == reflect.Ptr {
			if hubVal.IsNil() {
				hubVal.Set(reflect.New(hubVal.Type().Elem()))
			}
			if !isLast {
				hubVal = hubVal.Elem()
			}
		}
		// Non-pointer non-final kinds (struct, etc.) are already dereferenced — no-op needed.
	}
	// hubVal is now the Spec pointer field (already allocated above if nil)
	dstPtr := reflect.New(hubVal.Type().Elem())
	dstPtr.Elem().Set(reflect.ValueOf(poolSpec).Elem())
	hubVal.Set(dstPtr)
	return nil
}

// getHubPoolSpec extracts Spec.Router.Scheduler.Pool.Spec from Hub via reflection.
// Returns nil if not present or nil.
func getHubPoolSpec(hub conversion.Hub) *igwapiv1.InferencePoolSpec {
	hubVal := reflect.ValueOf(hub)
	for hubVal.Kind() == reflect.Interface || hubVal.Kind() == reflect.Ptr {
		if hubVal.IsNil() {
			return nil
		}
		hubVal = hubVal.Elem()
	}
	if hubVal.Kind() != reflect.Struct {
		return nil
	}
	for _, fieldName := range []string{"Spec", "Router", "Scheduler", "Pool", "Spec"} {
		found := findFieldInStruct(hubVal, fieldName)
		if !found.IsValid() {
			return nil
		}
		hubVal = found
		if hubVal.Kind() == reflect.Ptr {
			if hubVal.IsNil() {
				return nil
			}
			hubVal = hubVal.Elem()
		}
	}
	// hubVal is the dereferenced InferencePoolSpec from GIE v1
	if !hubVal.CanAddr() {
		return nil
	}
	spec, ok := hubVal.Addr().Interface().(*igwapiv1.InferencePoolSpec)
	if !ok {
		return nil
	}
	return spec
}

// findFieldInStruct finds a named field in a struct, including promoted fields from
// embedded (anonymous) structs. Returns zero Value if not found.
func findFieldInStruct(val reflect.Value, name string) reflect.Value {
	// First try direct lookup (handles non-embedded fields and top-level embedded)
	if f := val.FieldByName(name); f.IsValid() {
		return f
	}
	// Walk embedded structs
	t := val.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			nested := findFieldInStruct(val.Field(i), name)
			if nested.IsValid() {
				return nested
			}
		}
	}
	return reflect.Value{}
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

// cleanupCriticalityAnnotations removes criticality-related annotations from an ObjectMeta.
func cleanupCriticalityAnnotations(meta *metav1.ObjectMeta) {
	if meta == nil || meta.Annotations == nil {
		return
	}
	delete(meta.Annotations, ModelCriticalityAnnotationKey)
	delete(meta.Annotations, LoRACriticalitiesAnnotationKey)
	if len(meta.Annotations) == 0 {
		meta.Annotations = nil
	}
}

// convertPoolSpecForHub converts v1alpha1's GIE pool spec type to v1beta1's (official GIE v1).
// v1alpha1 uses github.com/kserve/kserve/pkg/apis/gie/v1alpha2pool while v1beta1 uses
// sigs.k8s.io/gateway-api-inference-extension/api/v1 — different packages, different structs.
// The old field-by-field reflection copy cannot handle this, so we delegate to the
// built-in ConvertTo/ConvertFrom on the GIE types.
func convertPoolSpecForHub(pool *InferencePoolSpec) *igwapiv1.InferencePoolSpec {
	if pool == nil || pool.Spec == nil {
		return nil
	}
	return convertInferencePoolSpecToV1(pool.Spec)
}

// convertPoolSpecFromHub converts Hub's GIE pool spec type (official GIE v1) back to
// v1alpha1's fork type (gie/v1alpha2pool).
func convertPoolSpecFromHub(poolSpec *igwapiv1.InferencePoolSpec) *InferencePoolSpec {
	if poolSpec == nil {
		return nil
	}
	return &InferencePoolSpec{
		Spec: convertInferencePoolSpecFromV1(poolSpec),
	}
}

// saveCriticalityToAnnotationsOnHub serialises v1alpha1-only Criticality fields into
// annotations on the Hub object. Operates on conversion.Hub via JSON indirection so we
// don't need a direct import on the Hub's concrete package.
func saveCriticalityToAnnotationsOnHub(hub conversion.Hub, model *LLMModelSpec) {
	if model.Criticality != nil && *model.Criticality != "" {
		setHubAnnotation(hub, ModelCriticalityAnnotationKey, string(*model.Criticality))
	}
	if model.LoRA != nil && len(model.LoRA.Adapters) > 0 {
		loraCriticalities := make(map[int]string)
		hasCriticality := false
		for i, adapter := range model.LoRA.Adapters {
			if adapter.Criticality != nil && *adapter.Criticality != "" {
				loraCriticalities[i] = string(*adapter.Criticality)
				hasCriticality = true
			}
		}
		if hasCriticality {
			if data, err := json.Marshal(loraCriticalities); err == nil {
				setHubAnnotation(hub, LoRACriticalitiesAnnotationKey, string(data))
			}
		}
	}
}

// restoreCriticalityFromHubAnnotations reads criticality annotations from the Hub object
// and restores them into the v1alpha1 model spec.
func restoreCriticalityFromHubAnnotations(hub conversion.Hub, model *LLMModelSpec) {
	ann := getHubAnnotations(hub)
	if ann == nil {
		return
	}
	if criticality, ok := ann[ModelCriticalityAnnotationKey]; ok && criticality != "" {
		c := Criticality(criticality)
		model.Criticality = &c
		delete(ann, ModelCriticalityAnnotationKey)
	}
	if loraData, ok := ann[LoRACriticalitiesAnnotationKey]; ok && loraData != "" {
		var loraCriticalities map[int]string
		if err := json.Unmarshal([]byte(loraData), &loraCriticalities); err == nil {
			if model.LoRA != nil {
				for i := range model.LoRA.Adapters {
					if criticality, exists := loraCriticalities[i]; exists {
						c := Criticality(criticality)
						model.LoRA.Adapters[i].Criticality = &c
					}
				}
			}
		}
		delete(ann, LoRACriticalitiesAnnotationKey)
	}
	if len(ann) == 0 {
		setHubAnnotations(hub, nil)
	} else {
		setHubAnnotations(hub, ann)
	}
}

// getHubAnnotations reads the metadata.annotations map from a conversion.Hub via JSON.
func getHubAnnotations(hub conversion.Hub) map[string]string {
	data, err := json.Marshal(hub)
	if err != nil {
		return nil
	}
	var wrap struct {
		Metadata struct {
			Annotations map[string]string `json:"annotations"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(data, &wrap); err != nil {
		return nil
	}
	return wrap.Metadata.Annotations
}

// setHubAnnotation sets a single annotation key/value on the Hub object via JSON.
func setHubAnnotation(hub conversion.Hub, key, value string) {
	ann := getHubAnnotations(hub)
	if ann == nil {
		ann = make(map[string]string)
	}
	ann[key] = value
	setHubAnnotations(hub, ann)
}

// setHubAnnotations writes the annotations map back into the Hub object via JSON round-trip.
func setHubAnnotations(hub conversion.Hub, ann map[string]string) {
	data, err := json.Marshal(hub)
	if err != nil {
		return
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return
	}
	var metaRaw map[string]interface{}
	if m, ok := raw["metadata"].(map[string]interface{}); ok {
		metaRaw = m
	} else {
		metaRaw = make(map[string]interface{})
	}
	metaBytes, _ := json.Marshal(ann)
	var annRaw interface{}
	_ = json.Unmarshal(metaBytes, &annRaw)
	metaRaw["annotations"] = annRaw
	raw["metadata"] = metaRaw
	out, err := json.Marshal(raw)
	if err != nil {
		return
	}
	_ = json.Unmarshal(out, hub)
}

func convertSpecToV1Alpha2(src *LLMInferenceServiceSpec) v1alpha2.LLMInferenceServiceSpec {
	dst := v1alpha2.LLMInferenceServiceSpec{
		Model:    convertModelSpecToV1Alpha2(&src.Model),
		Runtime:  src.Runtime,
		BaseRefs: src.BaseRefs,
	}

	// StorageInitializer - direct copy since structs are identical
	if src.StorageInitializer != nil {
		dst.StorageInitializer = &v1alpha2.StorageInitializerSpec{
			Enabled: src.StorageInitializer.Enabled,
		}
	}

	// WorkloadSpec (inline)
	dst.WorkloadSpec = convertWorkloadSpecToV1Alpha2(&src.WorkloadSpec)

	// Router
	if src.Router != nil {
		dst.Router = convertRouterSpecToV1Alpha2(src.Router)
	}

	// Prefill
	if src.Prefill != nil {
		prefill := convertWorkloadSpecToV1Alpha2(src.Prefill)
		dst.Prefill = &prefill
	}

	// Tracing
	if src.Tracing != nil {
		dst.Tracing = &v1alpha2.TracingSpec{
			ExporterEndpoint: src.Tracing.ExporterEndpoint,
			Sampler:          src.Tracing.Sampler,
			SamplerArg:       src.Tracing.SamplerArg,
			Exporter:         src.Tracing.Exporter,
		}
	}

	return dst
}

func convertSpecFromV1Alpha2(src *v1alpha2.LLMInferenceServiceSpec) LLMInferenceServiceSpec {
	dst := LLMInferenceServiceSpec{
		Model:    convertModelSpecFromV1Alpha2(&src.Model),
		Runtime:  src.Runtime,
		BaseRefs: src.BaseRefs,
	}

	// StorageInitializer - direct copy since structs are identical
	if src.StorageInitializer != nil {
		dst.StorageInitializer = &StorageInitializerSpec{
			Enabled: src.StorageInitializer.Enabled,
		}
	}

	// WorkloadSpec (inline)
	dst.WorkloadSpec = convertWorkloadSpecFromV1Alpha2(&src.WorkloadSpec)

	// Router
	if src.Router != nil {
		dst.Router = convertRouterSpecFromV1Alpha2(src.Router)
	}

	// Prefill
	if src.Prefill != nil {
		prefill := convertWorkloadSpecFromV1Alpha2(src.Prefill)
		dst.Prefill = &prefill
	}

	// Tracing
	if src.Tracing != nil {
		dst.Tracing = &TracingSpec{
			ExporterEndpoint: src.Tracing.ExporterEndpoint,
			Sampler:          src.Tracing.Sampler,
			SamplerArg:       src.Tracing.SamplerArg,
			Exporter:         src.Tracing.Exporter,
		}
	}

	return dst
}

func convertScalingSpecToV1Alpha2(src *ScalingSpec) *v1alpha2.ScalingSpec {
	if src == nil {
		return nil
	}

	dst := &v1alpha2.ScalingSpec{
		MinReplicas: src.MinReplicas,
		MaxReplicas: src.MaxReplicas,
	}

	if src.WVA != nil {
		dst.WVA = &v1alpha2.WVASpec{
			VariantCost: src.WVA.VariantCost,
		}
		if src.WVA.HPA != nil {
			dst.WVA.HPA = &v1alpha2.HPAScalingSpec{
				Behavior: src.WVA.HPA.Behavior,
			}
		}
		if src.WVA.KEDA != nil {
			dst.WVA.KEDA = convertKEDAScalingSpecToV1Alpha2(src.WVA.KEDA)
		}
	}

	if src.KEDA != nil {
		dst.KEDA = convertDirectKEDAScalingSpecToV1Alpha2(src.KEDA)
	}

	return dst
}

func convertScalingSpecFromV1Alpha2(src *v1alpha2.ScalingSpec) *ScalingSpec {
	if src == nil {
		return nil
	}

	dst := &ScalingSpec{
		MinReplicas: src.MinReplicas,
		MaxReplicas: src.MaxReplicas,
	}

	if src.WVA != nil {
		dst.WVA = &WVASpec{
			VariantCost: src.WVA.VariantCost,
		}
		if src.WVA.HPA != nil {
			dst.WVA.HPA = &HPAScalingSpec{
				Behavior: src.WVA.HPA.Behavior,
			}
		}
		if src.WVA.KEDA != nil {
			dst.WVA.KEDA = convertKEDAScalingSpecFromV1Alpha2(src.WVA.KEDA)
		}
	}

	if src.KEDA != nil {
		dst.KEDA = convertDirectKEDAScalingSpecFromV1Alpha2(src.KEDA)
	}

	return dst
}

func convertKEDAScalingSpecToV1Alpha2(src *KEDAScalingSpec) *v1alpha2.KEDAScalingSpec {
	if src == nil {
		return nil
	}
	return &v1alpha2.KEDAScalingSpec{
		PollingInterval:       src.PollingInterval,
		CooldownPeriod:        src.CooldownPeriod,
		InitialCooldownPeriod: src.InitialCooldownPeriod,
		IdleReplicaCount:      src.IdleReplicaCount,
		Fallback:              src.Fallback,
		Advanced:              src.Advanced,
	}
}

func convertKEDAScalingSpecFromV1Alpha2(src *v1alpha2.KEDAScalingSpec) *KEDAScalingSpec {
	if src == nil {
		return nil
	}
	return &KEDAScalingSpec{
		PollingInterval:       src.PollingInterval,
		CooldownPeriod:        src.CooldownPeriod,
		InitialCooldownPeriod: src.InitialCooldownPeriod,
		IdleReplicaCount:      src.IdleReplicaCount,
		Fallback:              src.Fallback,
		Advanced:              src.Advanced,
	}
}

func convertDirectKEDAScalingSpecToV1Alpha2(src *DirectKEDAScalingSpec) *v1alpha2.DirectKEDAScalingSpec {
	if src == nil {
		return nil
	}
	return &v1alpha2.DirectKEDAScalingSpec{
		KEDAScalingSpec: *convertKEDAScalingSpecToV1Alpha2(&src.KEDAScalingSpec),
		Triggers:        src.Triggers,
	}
}

func convertDirectKEDAScalingSpecFromV1Alpha2(src *v1alpha2.DirectKEDAScalingSpec) *DirectKEDAScalingSpec {
	if src == nil {
		return nil
	}
	return &DirectKEDAScalingSpec{
		KEDAScalingSpec: *convertKEDAScalingSpecFromV1Alpha2(&src.KEDAScalingSpec),
		Triggers:        src.Triggers,
	}
}

func convertModelSpecToV1Alpha2(src *LLMModelSpec) v1alpha2.LLMModelSpec {
	dst := v1alpha2.LLMModelSpec{
		URI:  src.URI,
		Name: src.Name,
		// Note: Criticality is preserved via annotations in ConvertTo/ConvertFrom
	}

	if src.LoRA != nil {
		dst.LoRA = convertLoRASpecToV1Alpha2(src.LoRA)
	}

	return dst
}

func convertModelSpecFromV1Alpha2(src *v1alpha2.LLMModelSpec) LLMModelSpec {
	dst := LLMModelSpec{
		URI:  src.URI,
		Name: src.Name,
		// Note: Criticality is restored from annotations in ConvertTo/ConvertFrom
	}

	if src.LoRA != nil {
		dst.LoRA = convertLoRASpecFromV1Alpha2(src.LoRA)
	}

	return dst
}

func convertLoRASpecToV1Alpha2(src *LoRASpec) *v1alpha2.LoRASpec {
	if src == nil {
		return nil
	}

	dst := &v1alpha2.LoRASpec{
		MaxRank:        src.MaxRank,
		MaxAdapters:    src.MaxAdapters,
		MaxCpuAdapters: src.MaxCpuAdapters,
	}
	for _, adapter := range src.Adapters {
		dst.Adapters = append(dst.Adapters, convertModelSpecToV1Alpha2(&adapter))
	}

	return dst
}

func convertLoRASpecFromV1Alpha2(src *v1alpha2.LoRASpec) *LoRASpec {
	if src == nil {
		return nil
	}

	dst := &LoRASpec{
		MaxRank:        src.MaxRank,
		MaxAdapters:    src.MaxAdapters,
		MaxCpuAdapters: src.MaxCpuAdapters,
	}
	for _, adapter := range src.Adapters {
		dst.Adapters = append(dst.Adapters, convertModelSpecFromV1Alpha2(&adapter))
	}

	return dst
}

func convertWorkloadSpecToV1Alpha2(src *WorkloadSpec) v1alpha2.WorkloadSpec {
	dst := v1alpha2.WorkloadSpec{
		Replicas:    src.Replicas,
		Labels:      src.Labels,
		Annotations: src.Annotations,
		Template:    src.Template,
		Worker:      src.Worker,
	}

	if src.Parallelism != nil {
		dst.Parallelism = &v1alpha2.ParallelismSpec{
			Tensor:      src.Parallelism.Tensor,
			Pipeline:    src.Parallelism.Pipeline,
			Data:        src.Parallelism.Data,
			DataLocal:   src.Parallelism.DataLocal,
			DataRPCPort: src.Parallelism.DataRPCPort,
			Expert:      src.Parallelism.Expert,
		}
	}

	if src.Scaling != nil {
		dst.Scaling = convertScalingSpecToV1Alpha2(src.Scaling)
	}

	if src.KVCacheOffloading != nil {
		dst.KVCacheOffloading = convertKVCacheOffloadingToV1Alpha2(src.KVCacheOffloading)
	}

	if src.RolloutStrategy != nil {
		dst.RolloutStrategy = &v1alpha2.RolloutStrategy{
			MaxUnavailable: src.RolloutStrategy.MaxUnavailable,
			MaxSurge:       src.RolloutStrategy.MaxSurge,
		}
	}

	return dst
}

func convertWorkloadSpecFromV1Alpha2(src *v1alpha2.WorkloadSpec) WorkloadSpec {
	dst := WorkloadSpec{
		Replicas:    src.Replicas,
		Labels:      src.Labels,
		Annotations: src.Annotations,
		Template:    src.Template,
		Worker:      src.Worker,
	}

	if src.Parallelism != nil {
		dst.Parallelism = &ParallelismSpec{
			Tensor:      src.Parallelism.Tensor,
			Pipeline:    src.Parallelism.Pipeline,
			Data:        src.Parallelism.Data,
			DataLocal:   src.Parallelism.DataLocal,
			DataRPCPort: src.Parallelism.DataRPCPort,
			Expert:      src.Parallelism.Expert,
		}
	}

	if src.Scaling != nil {
		dst.Scaling = convertScalingSpecFromV1Alpha2(src.Scaling)
	}

	if src.KVCacheOffloading != nil {
		dst.KVCacheOffloading = convertKVCacheOffloadingFromV1Alpha2(src.KVCacheOffloading)
	}

	if src.RolloutStrategy != nil {
		dst.RolloutStrategy = &RolloutStrategy{
			MaxUnavailable: src.RolloutStrategy.MaxUnavailable,
			MaxSurge:       src.RolloutStrategy.MaxSurge,
		}
	}

	return dst
}

func convertKVCacheOffloadingToV1Alpha2(src *KVCacheOffloadingSpec) *v1alpha2.KVCacheOffloadingSpec {
	dst := &v1alpha2.KVCacheOffloadingSpec{
		CPU:            src.CPU,
		EvictionPolicy: src.EvictionPolicy,
	}
	for _, s := range src.Secondary {
		dst.Secondary = append(dst.Secondary, convertSecondaryTierToV1Alpha2(&s))
	}
	return dst
}

func convertKVCacheOffloadingFromV1Alpha2(src *v1alpha2.KVCacheOffloadingSpec) *KVCacheOffloadingSpec {
	dst := &KVCacheOffloadingSpec{
		CPU:            src.CPU,
		EvictionPolicy: src.EvictionPolicy,
	}
	for _, s := range src.Secondary {
		dst.Secondary = append(dst.Secondary, convertSecondaryTierFromV1Alpha2(&s))
	}
	return dst
}

func convertSecondaryTierToV1Alpha2(src *SecondaryTierSpec) v1alpha2.SecondaryTierSpec {
	dst := v1alpha2.SecondaryTierSpec{}
	if src.FileSystem != nil {
		dst.FileSystem = &v1alpha2.FileSystemTierSpec{}
		if src.FileSystem.EmptyDir != nil {
			dst.FileSystem.EmptyDir = &v1alpha2.EmptyDirTierSpec{
				Size: src.FileSystem.EmptyDir.Size,
			}
		}
		if src.FileSystem.PVC != nil {
			dst.FileSystem.PVC = &v1alpha2.PVCTierSpec{
				Spec: src.FileSystem.PVC.Spec,
			}
			if src.FileSystem.PVC.Ref != nil {
				dst.FileSystem.PVC.Ref = &v1alpha2.PVCRefTierSpec{
					Name: src.FileSystem.PVC.Ref.Name,
					Path: src.FileSystem.PVC.Ref.Path,
				}
			}
		}
	}
	return dst
}

func convertSecondaryTierFromV1Alpha2(src *v1alpha2.SecondaryTierSpec) SecondaryTierSpec {
	dst := SecondaryTierSpec{}
	if src.FileSystem != nil {
		dst.FileSystem = &FileSystemTierSpec{}
		if src.FileSystem.EmptyDir != nil {
			dst.FileSystem.EmptyDir = &EmptyDirTierSpec{
				Size: src.FileSystem.EmptyDir.Size,
			}
		}
		if src.FileSystem.PVC != nil {
			dst.FileSystem.PVC = &PVCTierSpec{
				Spec: src.FileSystem.PVC.Spec,
			}
			if src.FileSystem.PVC.Ref != nil {
				dst.FileSystem.PVC.Ref = &PVCRefTierSpec{
					Name: src.FileSystem.PVC.Ref.Name,
					Path: src.FileSystem.PVC.Ref.Path,
				}
			}
		}
	}
	return dst
}

func convertRouterSpecToV1Alpha2(src *RouterSpec) *v1alpha2.RouterSpec {
	if src == nil {
		return nil
	}

	dst := &v1alpha2.RouterSpec{}

	if src.Route != nil {
		dst.Route = &v1alpha2.GatewayRoutesSpec{}
		if src.Route.HTTP != nil {
			dst.Route.HTTP = &v1alpha2.HTTPRouteSpec{
				Refs: src.Route.HTTP.Refs,
				Spec: src.Route.HTTP.Spec,
			}
		}
		dst.Route.Group = src.Route.Group
		dst.Route.Weight = src.Route.Weight
	}

	if src.Gateway != nil {
		dst.Gateway = &v1alpha2.GatewaySpec{}
		for _, ref := range src.Gateway.Refs {
			dst.Gateway.Refs = append(dst.Gateway.Refs, v1alpha2.GatewayObjectReference{
				UntypedObjectReference: v1alpha2.UntypedObjectReference{
					Name:      ref.Name,
					Namespace: ref.Namespace,
				},
				SectionName: ref.SectionName,
			})
		}
	}

	if src.Ingress != nil {
		dst.Ingress = &v1alpha2.IngressSpec{}
		for _, ref := range src.Ingress.Refs {
			dst.Ingress.Refs = append(dst.Ingress.Refs, v1alpha2.UntypedObjectReference{
				Name:      ref.Name,
				Namespace: ref.Namespace,
			})
		}
	}

	if src.Scheduler != nil {
		dst.Scheduler = &v1alpha2.SchedulerSpec{
			Labels:      src.Scheduler.Labels,
			Annotations: src.Scheduler.Annotations,
			Template:    src.Scheduler.Template,
			Replicas:    src.Scheduler.Replicas,
		}
		if src.Scheduler.Pool != nil {
			dst.Scheduler.Pool = &v1alpha2.InferencePoolSpec{
				Ref:  src.Scheduler.Pool.Ref,
				Spec: convertInferencePoolSpecToV1(src.Scheduler.Pool.Spec),
			}
		}
		if src.Scheduler.Config != nil {
			dst.Scheduler.Config = &v1alpha2.SchedulerConfigSpec{
				Inline: src.Scheduler.Config.Inline,
				Ref:    src.Scheduler.Config.Ref,
			}
		}
		if src.Scheduler.Tokenizer != nil {
			dst.Scheduler.Tokenizer = &v1alpha2.TokenizerSpec{
				Template: src.Scheduler.Tokenizer.Template,
			}
		}
	}

	return dst
}

// convertInferencePoolSpecToV1 converts igwapi v1alpha2 InferencePoolSpec to igwapi v1 InferencePoolSpec
// using the built-in conversion from the gateway-api-inference-extension library.
func convertInferencePoolSpecToV1(src *igwapiv1alpha2.InferencePoolSpec) *igwapiv1.InferencePoolSpec {
	if src == nil {
		return nil
	}

	// Use the built-in conversion by wrapping in temporary InferencePool objects
	srcPool := &igwapiv1alpha2.InferencePool{Spec: *src}
	dstPool := &igwapiv1.InferencePool{}

	if err := srcPool.ConvertTo(dstPool); err != nil {
		// Return nil rather than an empty spec — callers nil-check Pool.Spec, and an
		// empty spec would bypass those guards with invalid zero-value fields.
		return nil
	}

	return &dstPool.Spec
}

func convertRouterSpecFromV1Alpha2(src *v1alpha2.RouterSpec) *RouterSpec {
	if src == nil {
		return nil
	}

	dst := &RouterSpec{}

	if src.Route != nil {
		dst.Route = &GatewayRoutesSpec{}
		if src.Route.HTTP != nil {
			dst.Route.HTTP = &HTTPRouteSpec{
				Refs: src.Route.HTTP.Refs,
				Spec: src.Route.HTTP.Spec,
			}
		}
		dst.Route.Group = src.Route.Group
		dst.Route.Weight = src.Route.Weight
	}

	if src.Gateway != nil {
		dst.Gateway = &GatewaySpec{}
		for _, ref := range src.Gateway.Refs {
			dst.Gateway.Refs = append(dst.Gateway.Refs, GatewayObjectReference{
				UntypedObjectReference: UntypedObjectReference{
					Name:      ref.Name,
					Namespace: ref.Namespace,
				},
				SectionName: ref.SectionName,
			})
		}
	}

	if src.Ingress != nil {
		dst.Ingress = &IngressSpec{}
		for _, ref := range src.Ingress.Refs {
			dst.Ingress.Refs = append(dst.Ingress.Refs, UntypedObjectReference{
				Name:      ref.Name,
				Namespace: ref.Namespace,
			})
		}
	}

	if src.Scheduler != nil {
		dst.Scheduler = &SchedulerSpec{
			Labels:      src.Scheduler.Labels,
			Annotations: src.Scheduler.Annotations,
			Template:    src.Scheduler.Template,
			Replicas:    src.Scheduler.Replicas,
		}
		if src.Scheduler.Pool != nil {
			dst.Scheduler.Pool = &InferencePoolSpec{
				Ref:  src.Scheduler.Pool.Ref,
				Spec: convertInferencePoolSpecFromV1(src.Scheduler.Pool.Spec),
			}
		}
		if src.Scheduler.Config != nil {
			dst.Scheduler.Config = &SchedulerConfigSpec{
				Inline: src.Scheduler.Config.Inline,
				Ref:    src.Scheduler.Config.Ref,
			}
		}
		if src.Scheduler.Tokenizer != nil {
			dst.Scheduler.Tokenizer = &TokenizerSpec{
				Template: src.Scheduler.Tokenizer.Template,
			}
		}
	}

	return dst
}

// convertInferencePoolSpecFromV1 converts igwapi v1 InferencePoolSpec to igwapi v1alpha2 InferencePoolSpec
// using the built-in conversion from the gateway-api-inference-extension library.
func convertInferencePoolSpecFromV1(src *igwapiv1.InferencePoolSpec) *igwapiv1alpha2.InferencePoolSpec {
	if src == nil {
		return nil
	}

	// Use the built-in conversion by wrapping in temporary InferencePool objects
	srcPool := &igwapiv1.InferencePool{Spec: *src}
	dstPool := &igwapiv1alpha2.InferencePool{}

	if err := dstPool.ConvertFrom(srcPool); err != nil {
		// Return nil rather than an empty spec — callers nil-check Pool.Spec, and an
		// empty spec would bypass those guards with invalid zero-value fields.
		return nil
	}

	return &dstPool.Spec
}

// saveCriticalityToAnnotations stores criticality values in annotations to prevent data loss
// when converting from v1alpha1 to v1alpha2.
func saveCriticalityToAnnotations(meta *metav1.ObjectMeta, model *LLMModelSpec) {
	// Save model criticality
	if model.Criticality != nil && *model.Criticality != "" {
		if meta.Annotations == nil {
			meta.Annotations = make(map[string]string)
		}
		meta.Annotations[ModelCriticalityAnnotationKey] = string(*model.Criticality)
	}

	// Save LoRA adapter criticalities
	if model.LoRA != nil && len(model.LoRA.Adapters) > 0 {
		loraCriticalities := make(map[int]string)
		hasCriticality := false
		for i, adapter := range model.LoRA.Adapters {
			if adapter.Criticality != nil && *adapter.Criticality != "" {
				loraCriticalities[i] = string(*adapter.Criticality)
				hasCriticality = true
			}
		}
		if hasCriticality {
			if meta.Annotations == nil {
				meta.Annotations = make(map[string]string)
			}
			if data, err := json.Marshal(loraCriticalities); err == nil {
				meta.Annotations[LoRACriticalitiesAnnotationKey] = string(data)
			}
		}
	}
}

// restoreCriticalityFromAnnotations restores criticality values from annotations
// when converting from v1alpha2 to v1alpha1.
func restoreCriticalityFromAnnotations(meta *metav1.ObjectMeta, model *LLMModelSpec) {
	if meta.Annotations == nil {
		return
	}

	// Restore model criticality
	if criticality, ok := meta.Annotations[ModelCriticalityAnnotationKey]; ok && criticality != "" {
		c := Criticality(criticality)
		model.Criticality = &c
		delete(meta.Annotations, ModelCriticalityAnnotationKey)
	}

	// Restore LoRA adapter criticalities
	if loraData, ok := meta.Annotations[LoRACriticalitiesAnnotationKey]; ok && loraData != "" {
		var loraCriticalities map[int]string
		if err := json.Unmarshal([]byte(loraData), &loraCriticalities); err == nil {
			if model.LoRA != nil {
				for i := range model.LoRA.Adapters {
					if criticality, exists := loraCriticalities[i]; exists {
						c := Criticality(criticality)
						model.LoRA.Adapters[i].Criticality = &c
					}
				}
			}
		}
		delete(meta.Annotations, LoRACriticalitiesAnnotationKey)
	}

	// Clean up empty annotations map
	if len(meta.Annotations) == 0 {
		meta.Annotations = nil
	}
}
