package object

import (
	"fmt"
	"strings"
)

// ObjectType represents the type of an object
type ObjectType string

// Object type constants
const (
	INTEGER_OBJ  ObjectType = "INTEGER"
	FLOAT_OBJ    ObjectType = "FLOAT"
	BOOLEAN_OBJ  ObjectType = "BOOLEAN"
	NULL_OBJ     ObjectType = "NULL"
	STRING_OBJ   ObjectType = "STRING"
	ARRAY_OBJ    ObjectType = "ARRAY"
	MAP_OBJ      ObjectType = "MAP"
	RESOURCE_OBJ ObjectType = "RESOURCE"
	VARIABLE_OBJ ObjectType = "VARIABLE"
	OUTPUT_OBJ   ObjectType = "OUTPUT"
	ERROR_OBJ    ObjectType = "ERROR"
)

// Object is the interface that all runtime objects must implement
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer value
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Float represents a floating-point value
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null represents a null value
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// String represents a string value
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Array represents an array of objects
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out strings.Builder
	elements := make([]string, len(a.Elements))
	for i, el := range a.Elements {
		elements[i] = el.Inspect()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Map represents a map/dictionary of string keys to objects
type Map struct {
	Pairs map[string]Object
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out strings.Builder
	pairs := make([]string, 0, len(m.Pairs))
	for k, v := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Resource represents a cloud resource
type Resource struct {
	ResourceType string
	Name         string
	Properties   map[string]Object
}

func (r *Resource) Type() ObjectType { return RESOURCE_OBJ }
func (r *Resource) Inspect() string {
	return fmt.Sprintf("resource %s \"%s\"", r.ResourceType, r.Name)
}

// Variable represents a variable declaration
type Variable struct {
	Name        string
	VarType     string
	Default     Object
	Description string
}

func (v *Variable) Type() ObjectType { return VARIABLE_OBJ }
func (v *Variable) Inspect() string {
	return fmt.Sprintf("variable \"%s\"", v.Name)
}

// Output represents an output declaration
type Output struct {
	Name  string
	Value Object
}

func (o *Output) Type() ObjectType { return OUTPUT_OBJ }
func (o *Output) Inspect() string {
	return fmt.Sprintf("output \"%s\" = %s", o.Name, o.Value.Inspect())
}

// Error represents an error object
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Environment stores variable bindings
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

// NewEnclosedEnvironment creates a new environment enclosed by an outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores a value in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// All returns all bindings in the current environment (not including outer)
func (e *Environment) All() map[string]Object {
	result := make(map[string]Object)
	for k, v := range e.store {
		result[k] = v
	}
	return result
}
