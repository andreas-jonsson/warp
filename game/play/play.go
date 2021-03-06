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
	"log"
	"time"

	"github.com/andreas-jonsson/nanovgo"
	"github.com/andreas-jonsson/svgo/svgo"
	"github.com/andreas-jonsson/warp/data"
	"github.com/andreas-jonsson/warp/game"
	_ "github.com/andreas-jonsson/warp/game/entity/mothership"
	"github.com/andreas-jonsson/warp/game/universe"
	"github.com/andreas-jonsson/warp/platform"
	"github.com/ungerik/go3d/vec3"
)

type playState struct {
	uni       *universe.Universe
	gctl      game.GameControl
	mouseGrab bool
	cameraPos vec3.T

	svg *svgo.Svg

	warping       bool
	warpPos       vec3.T
	warpStartTime time.Time
}

func NewPlayState() *playState {
	fp, err := data.FS.Open("tiger.svg")
	if err != nil {
		log.Panicln(err)
	}
	defer fp.Close()

	svg, err := svgo.ParseSvg(fp, 1)
	if err != nil {
		log.Panicln(err)
	}

	return &playState{uni: universe.NewUniverse(), svg: svg}
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

func (s *playState) startWarp(x, y int) {
	s.warping = true
	s.warpStartTime = time.Now()
	s.warpPos = vec3.T{float32(x), float32(y), 0}
}

func (s *playState) stopWarp() {
	s.warping = false
}

func (s *playState) Update(gctl game.GameControl) error {
	for event := gctl.PollEvent(); event != nil; event = gctl.PollEvent() {
		switch t := event.(type) {
		case *platform.MouseButtonEvent:
			if t.Button == 1 {
				if t.Type == platform.MouseButtonDown {
					s.startWarp(t.X, t.Y)
				} else if t.Type == platform.MouseButtonUp {
					s.stopWarp()
				}
			} else if t.Button == 3 {
				if t.Type == platform.MouseButtonDown {
					s.mouseGrab = true
				} else if t.Type == platform.MouseButtonUp {
					s.mouseGrab = false
				}
			}
		case *platform.MouseMotionEvent:
			if s.mouseGrab {
				s.cameraPos = vec3.T{s.cameraPos[0] + float32(t.XRel), s.cameraPos[1] + float32(t.YRel), s.cameraPos[2]}
			}
		case *platform.MouseWheelEvent:
			s.cameraPos[2] += float32(t.Y)
		}
	}

	dt, _ := gctl.Timing()
	return s.uni.Update(dt, s.cameraPos)
}

func (s *playState) Render(ctx *nanovgo.Context) error {
	if err := s.uni.Render(ctx); err != nil {
		return err
	}

	if s.warping {
		warpScreenPos := vec3.Sub(&s.warpPos, &s.cameraPos)
		warpTime := time.Since(s.warpStartTime).Seconds()

		ctx.BeginPath()
		ctx.SetFillColor(nanovgo.RGBA(0, 0, 255, 255))
		ctx.Circle(warpScreenPos[0], warpScreenPos[1], float32(warpTime*10+4))
		ctx.Fill()

		ctx.BeginPath()
		ctx.MoveTo(0, 0)
		ctx.LineTo(warpScreenPos[0], warpScreenPos[1])
		ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 255, 255))
		ctx.SetStrokeWidth(2)
		ctx.Stroke()
	}

	if err := svgo.Render(ctx, s.svg); err != nil {
		panic(err)
	}

	return nil
}
