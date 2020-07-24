package hero

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/kataras/iris/v12/context"
)

// binding contains the Dependency and the Input, it's the result of a function or struct + dependencies.
type binding struct {
	Dependency *Dependency
	Input      *Input
}

// Input contains the input reference of which a dependency is binded to.
type Input struct {
	Index            int   // for func inputs
	StructFieldIndex []int // for struct fields in order to support embedded ones.
	Type             reflect.Type

	selfValue reflect.Value // reflect.ValueOf(*Input) cache.
}

func newInput(typ reflect.Type, index int, structFieldIndex []int) *Input {
	in := &Input{
		Index:            index,
		StructFieldIndex: structFieldIndex,
		Type:             typ,
	}

	in.selfValue = reflect.ValueOf(in)
	return in
}

// String returns the string representation of a binding.
func (b *binding) String() string {
	index := fmt.Sprintf("%d", b.Input.Index)
	if len(b.Input.StructFieldIndex) > 0 {
		for j, i := range b.Input.StructFieldIndex {
			if j == 0 {
				index = fmt.Sprintf("%d", i)
				continue
			}
			index += fmt.Sprintf(".%d", i)
		}
	}

	return fmt.Sprintf("[%s:%s] maps to [%s]", index, b.Input.Type.String(), b.Dependency)
}

// Equal compares "b" and "other" bindings and reports whether they are referring to the same values.
func (b *binding) Equal(other *binding) bool {
	if b == nil {
		return other == nil
	}

	if other == nil {
		return false
	}

	// if b.String() != other.String() {
	// 	return false
	// }

	if expected, got := b.Dependency != nil, other.Dependency != nil; expected != got {
		return false
	}

	if expected, got := fmt.Sprintf("%v", b.Dependency.OriginalValue), fmt.Sprintf("%v", other.Dependency.OriginalValue); expected != got {
		return false
	}

	if expected, got := b.Dependency.DestType != nil, other.Dependency.DestType != nil; expected != got {
		return false
	}

	if b.Dependency.DestType != nil {
		if expected, got := b.Dependency.DestType.String(), other.Dependency.DestType.String(); expected != got {
			return false
		}
	}

	if expected, got := b.Input != nil, other.Input != nil; expected != got {
		return false
	}

	if b.Input != nil {
		if expected, got := b.Input.Index, other.Input.Index; expected != got {
			return false
		}

		if expected, got := b.Input.Type.String(), other.Input.Type.String(); expected != got {
			return false
		}

		if expected, got := b.Input.StructFieldIndex, other.Input.StructFieldIndex; !reflect.DeepEqual(expected, got) {
			return false
		}
	}

	return true
}

func matchDependency(dep *Dependency, in reflect.Type) bool {
	if dep.Explicit {
		return dep.DestType == in
	}

	return dep.DestType == nil || equalTypes(dep.DestType, in)
}

func getBindingsFor(inputs []reflect.Type, deps []*Dependency, paramsCount int) (bindings []*binding) {
	// Path parameter start index is the result of [total path parameters] - [total func path parameters inputs],
	// moving from last to first path parameters and first to last (probably) available input args.
	//
	// That way the above will work as expected:
	// 1. mvc.New(app.Party("/path/{firstparam}")).Handle(....Controller.GetBy(secondparam string))
	// 2. mvc.New(app.Party("/path/{firstparam}/{secondparam}")).Handle(...Controller.GetBy(firstparam, secondparam string))
	// 3. usersRouter := app.Party("/users/{id:uint64}"); usersRouter.ConfigureContainer().Handle(method, "/", handler(id uint64))
	// 4. usersRouter.Party("/friends").ConfigureContainer().Handle(method, "/{friendID:uint64}", handler(friendID uint64))
	//
	// Therefore, count the inputs that can be path parameters first.
	shouldBindParams := make(map[int]struct{})
	totalParamsExpected := 0
	if paramsCount != -1 {
		for i, in := range inputs {
			if _, canBePathParameter := context.ParamResolvers[in]; !canBePathParameter {
				continue
			}
			shouldBindParams[i] = struct{}{}

			totalParamsExpected++
		}
	}

	startParamIndex := paramsCount - totalParamsExpected
	if startParamIndex < 0 {
		startParamIndex = 0
	}

	lastParamIndex := startParamIndex

	getParamIndex := func() int {
		paramIndex := lastParamIndex
		lastParamIndex++
		return paramIndex
	}

	bindedInput := make(map[int]struct{})

	for i, in := range inputs { //order matters.
		_, canBePathParameter := shouldBindParams[i]

		prevN := len(bindings) // to check if a new binding is attached; a dependency was matched (see below).

		for j := len(deps) - 1; j >= 0; j-- {
			d := deps[j]
			// Note: we could use the same slice to return.
			//
			// Add all dynamic dependencies (caller-selecting) and the exact typed dependencies.
			//
			// A dependency can only be matched to 1 value, and 1 value has a single dependency
			// (e.g. to avoid conflicting path parameters of the same type).
			if _, alreadyBinded := bindedInput[j]; alreadyBinded {
				continue
			}

			match := matchDependency(d, in)
			if !match {
				continue
			}

			if canBePathParameter {
				// wrap the existing dependency handler.
				paramHandler := paramDependencyHandler(getParamIndex())
				prevHandler := d.Handle
				d.Handle = func(ctx *context.Context, input *Input) (reflect.Value, error) {
					v, err := paramHandler(ctx, input)
					if err != nil {
						v, err = prevHandler(ctx, input)
					}

					return v, err
				}
				d.Static = false
				d.OriginalValue = nil
			}

			bindings = append(bindings, &binding{
				Dependency: d,
				Input:      newInput(in, i, nil),
			})

			if !d.Explicit { // if explicit then it can be binded to more than one input
				bindedInput[j] = struct{}{}
			}

			break
		}

		if prevN == len(bindings) {
			if canBePathParameter {
				// no new dependency added for this input,
				// let's check for path parameters.
				bindings = append(bindings, paramBinding(i, getParamIndex(), in))
				continue
			}

			// else add builtin bindings that may be registered by user too, but they didn't.
			if isPayloadType(in) {
				bindings = append(bindings, payloadBinding(i, in))
				continue
			}
		}
	}

	return
}

func isPayloadType(in reflect.Type) bool {
	switch indirectType(in).Kind() {
	case reflect.Struct, reflect.Slice, reflect.Ptr:
		return true
	default:
		return false
	}
}

func getBindingsForFunc(fn reflect.Value, dependencies []*Dependency, paramsCount int) []*binding {
	fnTyp := fn.Type()
	if !isFunc(fnTyp) {
		panic("bindings: unresolved: not a func type")
	}

	n := fnTyp.NumIn()
	inputs := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		inputs[i] = fnTyp.In(i)
	}

	bindings := getBindingsFor(inputs, dependencies, paramsCount)
	if expected, got := n, len(bindings); expected > got {
		panic(fmt.Sprintf("expected [%d] bindings (input parameters) but got [%d]", expected, got))
	}

	return bindings
}

func getBindingsForStruct(v reflect.Value, dependencies []*Dependency, paramsCount int, sorter Sorter) (bindings []*binding) {
	typ := indirectType(v.Type())
	if typ.Kind() != reflect.Struct {
		panic("bindings: unresolved: no struct type")
	}

	// get bindings from any struct's non zero values first, including unexported.
	elem := reflect.Indirect(v)
	nonZero := lookupNonZeroFieldValues(elem)
	for _, f := range nonZero {
		// fmt.Printf("Controller [%s] | NonZero | Field Index: %v | Field Type: %s\n", typ, f.Index, f.Type)
		bindings = append(bindings, &binding{
			Dependency: NewDependency(elem.FieldByIndex(f.Index).Interface()),
			Input:      newInput(f.Type, f.Index[0], f.Index),
		})
	}

	fields, stateless := lookupFields(elem, true, true, nil)
	n := len(fields)

	if n > 1 && sorter != nil {
		sort.Slice(fields, func(i, j int) bool {
			return sorter(fields[i].Type, fields[j].Type)
		})
	}

	inputs := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		// fmt.Printf("Controller [%s] | Field Index: %v | Field Type: %s\n", typ, fields[i].Index, fields[i].Type)
		inputs[i] = fields[i].Type
	}
	exportedBindings := getBindingsFor(inputs, dependencies, paramsCount)
	// fmt.Printf("Controller [%s] | Inputs length: %d vs Bindings length: %d\n", typ, n, len(exportedBindings))

	if stateless == 0 && len(nonZero) >= len(exportedBindings) {
		// if we have not a single stateless and fields are defined then just return.
		// Note(@kataras): this can accept further improvements.
		return
	}

	// get declared bindings from deps.
	bindings = append(bindings, exportedBindings...)
	for _, binding := range bindings {
		// fmt.Printf(""Controller [%s] | Binding: %s\n", typ, binding.String())

		if len(binding.Input.StructFieldIndex) == 0 {
			// set correctly the input's field index.
			structFieldIndex := fields[binding.Input.Index].Index
			binding.Input.StructFieldIndex = structFieldIndex
		}

		// fmt.Printf("Controller [%s] | binding Index: %v | binding Type: %s\n", typ, binding.Input.StructFieldIndex, binding.Input.Type)
		// fmt.Printf("Controller [%s] Set [%s] to struct field index: %v\n", typ.String(), binding.Input.Type.String(), binding.Input.StructFieldIndex)
	}

	return
}

/*
	Builtin dynamic bindings.
*/

func paramBinding(index, paramIndex int, typ reflect.Type) *binding {
	return &binding{
		Dependency: &Dependency{Handle: paramDependencyHandler(paramIndex), DestType: typ, Source: getSource()},
		Input:      newInput(typ, index, nil),
	}
}

func paramDependencyHandler(paramIndex int) DependencyHandler {
	return func(ctx *context.Context, input *Input) (reflect.Value, error) {
		if ctx.Params().Len() <= paramIndex {
			return emptyValue, ErrSeeOther
		}

		return reflect.ValueOf(ctx.Params().Store[paramIndex].ValueRaw), nil
	}
}

// registered if input parameters are more than matched dependencies.
// It binds an input to a request body based on the request content-type header (JSON, XML, YAML, Query, Form).
func payloadBinding(index int, typ reflect.Type) *binding {
	return &binding{
		Dependency: &Dependency{
			Handle: func(ctx *context.Context, input *Input) (newValue reflect.Value, err error) {
				wasPtr := input.Type.Kind() == reflect.Ptr

				if serveDepsV := ctx.Values().Get(context.DependenciesContextKey); serveDepsV != nil {
					if serveDeps, ok := serveDepsV.(context.DependenciesMap); ok {
						if newValue, ok = serveDeps[typ]; ok {
							return
						}
					}
				}

				if input.Type.Kind() == reflect.Slice {
					newValue = reflect.New(reflect.SliceOf(indirectType(input.Type)))
				} else {
					newValue = reflect.New(indirectType(input.Type))
				}

				ptr := newValue.Interface()
				err = ctx.ReadBody(ptr)

				if !wasPtr {
					newValue = newValue.Elem()
				}

				return
			},
			Source: getSource(),
		},
		Input: newInput(typ, index, nil),
	}

}