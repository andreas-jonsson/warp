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

package main

import (
	"log"
	"math"
	"time"

	"github.com/andreas-jonsson/warp/platform"
	"github.com/shibukawa/nanovgo"
)

func main() {
	if err := platform.Init(); err != nil {
		log.Panicln(err)
	}
	defer platform.Shutdown()

	rnd, err := platform.NewRenderer(640, 360)
	if err != nil {
		log.Panicln(err)
	}
	defer rnd.Shutdown()

	for {
		for ev := platform.PollEvent(); ev != nil; ev = platform.PollEvent() {
			switch ev.(type) {
			case *platform.QuitEvent, *platform.KeyUpEvent:
				return
			}
		}

		rnd.Clear()
		drawVG(rnd.VG())
		rnd.Present()
	}
}

var startupTime = time.Now()

func drawVG(ctx *nanovgo.Context) {
	mouse := platform.Mouse()
	drawEyes(ctx, 100, 100, 50, 50, float32(mouse.X), float32(mouse.Y), float32(time.Since(startupTime)*time.Second)*0.001)
}

func drawEyes(ctx *nanovgo.Context, x, y, w, h, mx, my, t float32) {
	ex := w * 0.23
	ey := h * 0.5
	lx := x + ex
	ly := y + ey
	rx := x + w - ex
	ry := y + ey
	var dx, dy, d, br float32
	if ex < ey {
		br = ex * 0.5
	} else {
		br = ey * 0.5
	}
	blink := float32(1.0 - math.Pow(float64(sinF(t*0.5)), 200)*0.8)

	bg1 := nanovgo.LinearGradient(x, y+h*0.5, x+w*0.1, y+h, nanovgo.RGBA(0, 0, 0, 32), nanovgo.RGBA(0, 0, 0, 16))
	ctx.BeginPath()
	ctx.Ellipse(lx+3.0, ly+16.0, ex, ey)
	ctx.Ellipse(rx+3.0, ry+16.0, ex, ey)
	ctx.SetFillPaint(bg1)
	ctx.Fill()

	bg2 := nanovgo.LinearGradient(x, y+h*0.25, x+w*0.1, y+h, nanovgo.RGBA(220, 220, 220, 255), nanovgo.RGBA(128, 128, 128, 255))
	ctx.BeginPath()
	ctx.Ellipse(lx, ly, ex, ey)
	ctx.Ellipse(rx, ry, ex, ey)
	ctx.SetFillPaint(bg2)
	ctx.Fill()

	dx = (mx - rx) / (ex * 10)
	dy = (my - ry) / (ey * 10)
	d = sqrtF(dx*dx + dy*dy)
	if d > 1.0 {
		dx /= d
		dy /= d
	}
	dx *= ex * 0.4
	dy *= ey * 0.5
	ctx.BeginPath()
	ctx.Ellipse(lx+dx, ly+dy+ey*0.25*(1.0-blink), br, br*blink)
	ctx.SetFillColor(nanovgo.RGBA(32, 32, 32, 255))
	ctx.Fill()

	dx = (mx - rx) / (ex * 10)
	dy = (my - ry) / (ey * 10)
	d = sqrtF(dx*dx + dy*dy)
	if d > 1.0 {
		dx /= d
		dy /= d
	}
	dx *= ex * 0.4
	dy *= ey * 0.5
	ctx.BeginPath()
	ctx.Ellipse(rx+dx, ry+dy+ey*0.25*(1.0-blink), br, br*blink)
	ctx.SetFillColor(nanovgo.RGBA(32, 32, 32, 255))
	ctx.Fill()

	dx = (mx - rx) / (ex * 10)
	dy = (my - ry) / (ey * 10)
	d = sqrtF(dx*dx + dy*dy)
	if d > 1.0 {
		dx /= d
		dy /= d
	}
	dx *= ex * 0.4
	dy *= ey * 0.5
	ctx.BeginPath()
	ctx.Ellipse(rx+dx, ry+dy+ey*0.25*(1.0-blink), br, br*blink)
	ctx.SetFillColor(nanovgo.RGBA(32, 32, 32, 255))
	ctx.Fill()

	gloss1 := nanovgo.RadialGradient(lx-ex*0.25, ly-ey*0.5, ex*0.1, ex*0.75, nanovgo.RGBA(255, 255, 255, 128), nanovgo.RGBA(255, 255, 255, 0))
	ctx.BeginPath()
	ctx.Ellipse(lx, ly, ex, ey)
	ctx.SetFillPaint(gloss1)
	ctx.Fill()

	gloss2 := nanovgo.RadialGradient(rx-ex*0.25, ry-ey*0.5, ex*0.1, ex*0.75, nanovgo.RGBA(255, 255, 255, 128), nanovgo.RGBA(255, 255, 255, 0))
	ctx.BeginPath()
	ctx.Ellipse(rx, ry, ex, ey)
	ctx.SetFillPaint(gloss2)
	ctx.Fill()
}

func cosF(a float32) float32 {
	return float32(math.Cos(float64(a)))
}

func sinF(a float32) float32 {
	return float32(math.Sin(float64(a)))
}

func sqrtF(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}

func clampF(a, min, max float32) float32 {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func absF(a float32) float32 {
	if a > 0.0 {
		return a
	}
	return -a
}

func maxF(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
