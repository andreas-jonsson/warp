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

//go:generate go run data/generate.go

package main

import (
	"fmt"
	"log"

	"github.com/mode13/warp/game"
	"github.com/mode13/warp/game/menu"
	"github.com/mode13/warp/game/play"
	"github.com/mode13/warp/platform"
)

func main() {
	if err := platform.Init(); err != nil {
		log.Panicln(err)
	}
	defer platform.Shutdown()

	rnd, err := platform.NewRenderer(platform.ConfigWithDiv(2), platform.ConfigWithNoVSync)
	if err != nil {
		log.Panicln(err)
	}
	defer rnd.Shutdown()

	states := map[string]game.GameState{
		"menu": menu.NewMenuState(),
		"play": play.NewPlayState(),
	}

	g, err := game.NewGame(states)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Shutdown()

	var gctl game.GameControl = g
	if err := g.SwitchState("menu", gctl); err != nil {
		log.Panicln(err)
	}

	for g.Running() {
		ctx := rnd.Clear()

		if err := g.Update(); err != nil {
			log.Panicln(err)
		}

		_, fps := g.Timing()
		rnd.SetWindowTitle(fmt.Sprintf("Warp - %d fps", fps))

		if err := g.Render(ctx); err != nil {
			log.Panicln(err)
		}

		rnd.Present()
	}
}
