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

package mothership

import (
	"github.com/mode13/warp/game/entity"
	"github.com/shibukawa/nanovgo"
	"github.com/ungerik/go3d/vec2"
)

const hp = 100

type ship struct {
	id  uint64
	hp  float32
	pos vec2.T
}

func init() {
	entity.RegisterConstructor("mothership", newMothership)
}

func newMothership(id uint64, owner int) entity.Entity {
	return &ship{id: id, hp: hp}
}

func (e *ship) Id() uint64 {
	return e.id
}

func (e *ship) Alive() bool {
	return e.hp > 0
}

func (e *ship) TakeFire(damage float32, ty entity.DamageType) bool {
	e.hp -= damage
	return e.Alive()
}

func (e *ship) Update(uni entity.Universe) error {
	return nil
}

func (e *ship) Render(ctx *nanovgo.Context) error {
	ctx.BeginPath()
	ctx.SetFillColor(nanovgo.RGBA(255, 0, 0, 255))
	ctx.Circle(0, 0, 1000)
	ctx.Fill()
	return nil
}
