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

package universe

import (
	"github.com/mode13/warp/game/entity"
	"github.com/mode13/warp/platform"
	"github.com/shibukawa/nanovgo"
	"github.com/ungerik/go3d/vec2"
)

type Universe struct {
	entities map[uint64]entity.Entity
}

func NewUniverse() *Universe {
	return &Universe{entities: make(map[uint64]entity.Entity)}
}

func (uni *Universe) SpawnEntity(ty string, owner int) entity.Entity {
	entity := entity.NewEntity(ty, platform.NewId64(), owner)
	uni.entities[entity.Id()] = entity
	return entity
}

func (uni *Universe) FindAll(pos vec2.T, rad float32, filter uint32) []entity.Entity {
	return nil
}

func (uni *Universe) Update() error {
	for _, entity := range uni.entities {
		if err := entity.Update(uni); err != nil {
			return err
		}
	}

	for id, entity := range uni.entities {
		if !entity.Alive() {
			delete(uni.entities, id)
		}
	}
	return nil
}

func (uni *Universe) Render(ctx *nanovgo.Context) error {
	for _, entity := range uni.entities {
		if err := entity.Render(ctx); err != nil {
			return err
		}
	}
	return nil
}
