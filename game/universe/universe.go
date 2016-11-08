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
	"time"

	"github.com/shibukawa/nanovgo"
	"github.com/ungerik/go3d/vec2"
)

const LightSpeed = 1.0

type Entity interface {
	Id() uint64
	Position() vec2.T
	Render(ctx *nanovgo.Context)
	Clone() Entity
}

type (
	universeState struct {
		timeStamp time.Time
		entities  map[uint64]Entity
	}

	universe struct {
		frames   []universeState
		entities []Entity

		frameIdx,
		maxFrames int64
		fps int
	}
)

func NewUniverse(rad float32, fps int) *universe {
	lightTravelTime := rad * 2.0 / LightSpeed
	maxFrames := int64(lightTravelTime * float32(fps))

	return &universe{fps: fps, maxFrames: maxFrames, frames: make([]universeState, maxFrames)}
}

func (uni *universe) NewFrame() {
	state := universeState{timeStamp: time.Now(), entities: make(map[uint64]Entity)}
	for _, entity := range uni.entities {
		state.entities[entity.Id()] = entity.Clone()
	}

	uni.frames[uni.frameIdx] = state
	uni.frameIdx++
}

func (uni *universe) StateFromObserver(pos *vec2.T) []Entity {
	entities := make([]Entity, 1)

	for _, entity := range uni.entities {
		p := entity.Position()
		dist := pos.Sub(&p).Length()
		frameDist := int64((dist / LightSpeed) * float32(uni.fps))
		idx := uni.frameIdx - frameDist

		if idx < 0 {
			idx = 0
		}

		frame := uni.frames[idx]
		entity = frame.entities[entity.Id()]
		entities = append(entities, entity)
	}

	return entities
}
