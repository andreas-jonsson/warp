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

package play

import (
	"github.com/mode13/warp/game"
	_ "github.com/mode13/warp/game/entity/mothership"
	"github.com/mode13/warp/game/universe"
	"github.com/shibukawa/nanovgo"
)

type playState struct {
	uni  *universe.Universe
	gctl game.GameControl
}

func NewPlayState() *playState {
	return &playState{uni: universe.NewUniverse()}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from game.GameState, args ...interface{}) error {
	s.gctl = args[0].(game.GameControl)

	s.uni.SpawnEntity("mothership", 0)
	return nil
}

func (s *playState) Exit(to game.GameState) error {
	return nil
}

func (s *playState) Update(gctl game.GameControl) error {
	gctl.PollAll()
	return s.uni.Update()
}

func (s *playState) Render(ctx *nanovgo.Context) error {
	return s.uni.Render(ctx)
}
