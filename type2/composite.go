/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http//www.gnu.org/licenses/>.
 *
 * composite.go
 *
 *  Created on May 07, 2017
 *      Author Massimiliano Ghilardi
 */

package type2

import (
	"go/token"
	"go/types"
	"reflect"
)

// ChanDir returns a channel type's direction.
// It panics if the type's Kind is not Chan.
func (t *timpl) ChanDir() reflect.ChanDir {
	if t.Kind() != reflect.Chan {
		errorf("ChanDir of non-chan type %v", t)
	}
	return t.rtype.ChanDir()
}

// Elem returns a type's element type.
// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
func (t *timpl) Elem() Type {
	gtype := t.underlying()
	rtype := t.rtype
	switch gtype := gtype.(type) {
	case *types.Array:
		return maketype(gtype.Elem(), rtype.Elem())
	case *types.Chan:
		return maketype(gtype.Elem(), rtype.Elem())
	case *types.Map:
		return maketype(gtype.Elem(), rtype.Elem())
	case *types.Pointer:
		return maketype(gtype.Elem(), rtype.Elem())
	case *types.Slice:
		return maketype(gtype.Elem(), rtype.Elem())
	default:
		errorf("Elem of invalid type %v", t)
		return Type{}
	}
}

// Key returns a map type's key type.
// It panics if the type's Kind is not Map.
func (t *timpl) Key() Type {
	if t.Kind() != reflect.Map {
		errorf("Key of non-map type %v", t)
	}
	gtype := t.underlying().(*types.Map)
	return maketype(gtype.Key(), t.rtype.Key())
}

// Len returns an array type's length.
// It panics if the type's Kind is not Array.
func (t *timpl) Len() int {
	if t.Kind() != reflect.Func {
		errorf("Len of non-array type %v", t)
	}
	return t.rtype.Len()
}

func ArrayOf(count int, elem Type) Type {
	return maketype(
		types.NewArray(elem.gtype, int64(count)),
		reflect.ArrayOf(count, elem.rtype))
}

func ChanOf(dir reflect.ChanDir, elem Type) Type {
	var gdir types.ChanDir
	switch dir {
	case reflect.RecvDir:
		gdir = types.RecvOnly
	case reflect.SendDir:
		gdir = types.SendOnly
	case reflect.BothDir:
		gdir = types.SendRecv
	}
	return maketype(
		types.NewChan(gdir, elem.gtype),
		reflect.ChanOf(dir, elem.rtype))
}

func MapOf(key, elem Type) Type {
	return maketype(
		types.NewMap(key.gtype, elem.gtype),
		reflect.MapOf(key.rtype, elem.rtype))
}

// NamedOf returns a new named type for the given type name.
// Initially, the underlying type is set to interface{} - use SetUnderlying to change it.
// These two steps are separate to allow creating self-referencing types,
// as for example type List struct { Elem int; Rest *List }
func NamedOf(name string, pkg *Package) Type {
	underlying := TypeOfInterface
	typename := types.NewTypeName(token.NoPos, pkg.impl, name, underlying.gtype)
	return maketype(
		types.NewNamed(typename, underlying.gtype, nil),
		underlying.rtype)
}

// SetUnderlying sets the underlying type of a named type and marks t as complete.
// It panics if the type is unnamed, or if the underlying type is named.
func (t *timpl) SetUnderlying(underlying Type) {
	switch gtype := t.gtype.(type) {
	case *types.Named:
		gtype.SetUnderlying(underlying.gtype)
	default:
		errorf("SetUnderlying of unnamed type %v", t)
	}
}

// AddMethod adds method 'name' to type, unless it is already in the method list.
// It panics if the type is unnamed, or if the signature is not a function type.
func (t *timpl) AddMethod(pkg Package, name string, signature Type) {
	gtype, ok := t.gtype.(*types.Named)
	if !ok {
		errorf("AddMethod on unnamed type %v", t)
	}
	if signature.kind != reflect.Func {
		errorf("AddMethod of non-func signature: %v", signature)
	}
	gsig := signature.underlying().(*types.Signature)
	gfun := types.NewFunc(token.NoPos, pkg.impl, name, gsig)
	gtype.AddMethod(gfun)
}

func PtrTo(elem Type) Type {
	return maketype(
		types.NewPointer(elem.gtype),
		reflect.PtrTo(elem.rtype))
}

func SliceOf(elem Type) Type {
	return maketype(
		types.NewSlice(elem.gtype),
		reflect.SliceOf(elem.rtype))
}