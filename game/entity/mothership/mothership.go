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
	"math"

	"github.com/mode13/nanovgo"
	"github.com/mode13/warp/game/entity"
	"github.com/ungerik/go3d/vec2"
)

const (
	hp         = 100
	circleSize = 150
)

type ship struct {
	id         uint64
	hp         float32
	numCircles int
	pos        vec2.T
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
	if e.numCircles == 0 {
		size := uni.Bounds().Size()
		diagonal := math.Sqrt(float64(size.X*size.X + size.Y*size.Y))
		e.numCircles = int(diagonal/circleSize) + 1
	}
	return nil
}

func (e *ship) Render(ctx *nanovgo.Context) error {
	ctx.BeginPath()
	ctx.SetFillColor(nanovgo.RGBA(25, 25, 200, 200))
	ctx.Circle(e.pos[0], e.pos[1], 10)
	ctx.Fill()

	for i := 1; i < e.numCircles; i++ {
		ctx.BeginPath()
		ctx.Circle(e.pos[0], e.pos[1], float32(i)*150)
		ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 200, 175))
		ctx.SetStrokeWidth(1)
		ctx.Stroke()
	}

	return nil
}
