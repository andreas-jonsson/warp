// +build !js,!mobile

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

package platform

import (
	"log"

	"github.com/goxjs/gl"
	"github.com/shibukawa/nanovgo"
	"github.com/veandco/go-sdl2/sdl"
)

const fulscreenFlag = sdl.WINDOW_FULLSCREEN_DESKTOP //sdl.WINDOW_FULLSCREEN

type sdlRenderer struct {
	window         *sdl.Window
	glContext      sdl.GLContext
	vgContext      *nanovgo.Context
	debug, novsync bool
}

func NewRenderer(w, h int, data ...interface{}) (Renderer, error) {
	var (
		err error
		rnd sdlRenderer
		dm  sdl.DisplayMode

		title           = "Warp"
		sdlFlags uint32 = sdl.WINDOW_SHOWN | sdl.WINDOW_OPENGL
		vgFlags         = nanovgo.StencilStrokes
	)

	if err = sdl.GetDesktopDisplayMode(0, &dm); err != nil {
		return &rnd, err
	}

	if w <= 0 {
		w = int(dm.W)
	}
	if h <= 0 {
		h = int(dm.H)
	}

	for i := 0; i < len(data); i++ {
		handled := true
		p := data[i]

		ps, ok := p.(string)
		if ok {
			switch ps {
			case "fullscreen":
				sdlFlags |= fulscreenFlag
			case "debug":
				rnd.debug = true
				vgFlags |= nanovgo.Debug
			case "novsync":
				rnd.novsync = true
			case "title":
				i++
				title = data[i].(string)
			case "div":
				i++
				n := data[i].(int)
				w /= n
				h /= n
			default:
				handled = false
			}
		}

		if !handled {
			log.Println("invalid parameter passed to renderer:", p)
		}
	}

	sdl.GL_SetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_BLUE_SIZE, 8)
	//sdl.GL_SetAttribute(sdl.GL_ALPHA_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GL_SetAttribute(sdl.GL_STENCIL_SIZE, 8)

	sdl.GL_SetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)

	sdl.GL_SetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)

	rnd.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdlFlags)
	if err != nil {
		return &rnd, err
	}

	rnd.glContext, err = sdl.GL_CreateContext(rnd.window)
	if err != nil {
		return &rnd, err
	}

	gl.ContextWatcher.OnMakeCurrent(nil)
	if rnd.novsync {
		sdl.GL_SetSwapInterval(0)
	} else {
		sdl.GL_SetSwapInterval(1)
	}

	rnd.vgContext, err = nanovgo.NewContext(vgFlags)
	if err != nil {
		return &rnd, err
	}

	rnd.window.SetGrab(true)
	//sdl.ShowCursor(0)
	return &rnd, nil
}

func (rnd *sdlRenderer) ToggleFullscreen() {
	isFullscreen := (rnd.window.GetFlags() & fulscreenFlag) != 0
	if isFullscreen {
		rnd.window.SetFullscreen(0)
	} else {
		rnd.window.SetFullscreen(fulscreenFlag)
	}
}

func (rnd *sdlRenderer) Clear() *nanovgo.Context {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	w, h := rnd.window.GetSize()
	rnd.vgContext.BeginFrame(w, h, float32(w)/float32(h))
	return rnd.vgContext
}

func (rnd *sdlRenderer) Present() {
	rnd.vgContext.EndFrame()
	sdl.GL_SwapWindow(rnd.window)

	if rnd.debug {
		if err := gl.GetError(); err != gl.NO_ERROR {
			log.Panicln(err)
		}
	}
}

func (rnd *sdlRenderer) Shutdown() {
	rnd.vgContext.Delete()
	gl.ContextWatcher.OnDetach()
	sdl.GL_DeleteContext(rnd.glContext)
	rnd.window.Destroy()
}

func (rnd *sdlRenderer) SetWindowTitle(title string) {
	rnd.window.SetTitle(title)
}
