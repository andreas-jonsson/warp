/*
Copyright (C) 2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package entity

import (
	"image"
	"log"

	"github.com/andreas-jonsson/nanovgo"
	"github.com/ungerik/go3d/vec2"
	"github.com/ungerik/go3d/vec3"
)

const (
	EnemyUnits = iota
	PlayerUnits
	TeamUnits
)

const (
	LaserDamage DamageType = iota
)

type (
	Constructor func(uint64, int) Entity
	DamageType  int

	Universe interface {
		SpawnEntity(ty string, owner int) Entity
		FindAll(pos vec2.T, rad float32, filter uint32) []Entity
		CameraPosition() vec3.T
		Bounds() image.Rectangle
	}

	Entity interface {
		Id() uint64
		Alive() bool
		TakeFire(damage float32, ty DamageType) bool
		Update(uni Universe) error
		Render(ctx *nanovgo.Context) error
	}
)

var entityConstructors = make(map[string]Constructor)

func RegisterConstructor(ty string, init Constructor) {
	entityConstructors[ty] = init
}

func NewEntity(ty string, id uint64, owner int) Entity {
	c, ok := entityConstructors[ty]
	if !ok {
		log.Panicln("invalid entity type:", ty)
	}
	return c(id, owner)
}
