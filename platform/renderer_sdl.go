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
	"image"
	"log"
	"unsafe"

	"github.com/goxjs/gl"
	"github.com/mode13/nanovgo"
	"github.com/veandco/go-sdl2/sdl"
)

const blurEffect = false

const fulscreenFlag = sdl.WINDOW_FULLSCREEN_DESKTOP //sdl.WINDOW_FULLSCREEN

type Config func(*sdlRenderer) error

func ConfigWithSize(w, h int) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.windowSize = image.Point{w, h}
		return nil
	}
}

func ConfigWithTitle(title string) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.windowTitle = title
		return nil
	}
}

func ConfigWithDiv(n int) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.resolutionDiv = n
		return nil
	}
}

func ConfigWithFulscreen(rnd *sdlRenderer) error {
	rnd.config.fulscreen = true
	return nil
}

func ConfigWithDebug(rnd *sdlRenderer) error {
	rnd.config.debug = true
	return nil
}

func ConfigWithNoVSync(rnd *sdlRenderer) error {
	rnd.config.novsync = true
	return nil
}

type sdlRenderer struct {
	window    *sdl.Window
	glContext sdl.GLContext
	vgContext *nanovgo.Context

	glBlurTexture gl.Texture
	glProgram     gl.Program
	glSquareBuffer,
	glSquareUVBuffer gl.Buffer

	config struct {
		windowTitle   string
		windowSize    image.Point
		resolutionDiv int
		debug, novsync,
		fulscreen bool
	}
}

func NewRenderer(configs ...Config) (*sdlRenderer, error) {
	var (
		err error
		rnd sdlRenderer
		dm  sdl.DisplayMode

		sdlFlags uint32 = sdl.WINDOW_SHOWN | sdl.WINDOW_OPENGL
		vgFlags         = nanovgo.StencilStrokes
	)

	for _, cfg := range configs {
		if err = cfg(&rnd); err != nil {
			return nil, err
		}
	}

	if err = sdl.GetDesktopDisplayMode(0, &dm); err != nil {
		return &rnd, err
	}

	cfg := &rnd.config
	if cfg.windowSize.X <= 0 {
		cfg.windowSize.X = int(dm.W)
	}
	if cfg.windowSize.Y <= 0 {
		cfg.windowSize.Y = int(dm.H)
	}

	if cfg.resolutionDiv > 0 {
		cfg.windowSize.X /= cfg.resolutionDiv
		cfg.windowSize.Y /= cfg.resolutionDiv
	}

	sdl.GL_SetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_BLUE_SIZE, 8)
	//sdl.GL_SetAttribute(sdl.GL_ALPHA_SIZE, 8)
	//sdl.GL_SetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GL_SetAttribute(sdl.GL_STENCIL_SIZE, 8)

	sdl.GL_SetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)

	sdl.GL_SetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)

	rnd.window, err = sdl.CreateWindow(cfg.windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, cfg.windowSize.X, cfg.windowSize.Y, sdlFlags)
	if err != nil {
		return &rnd, err
	}

	rnd.glContext, err = sdl.GL_CreateContext(rnd.window)
	if err != nil {
		return &rnd, err
	}

	gl.ContextWatcher.OnMakeCurrent(nil)
	if cfg.novsync {
		sdl.GL_SetSwapInterval(0)
	} else {
		sdl.GL_SetSwapInterval(1)
	}

	rnd.vgContext, err = nanovgo.NewContext(vgFlags)
	if err != nil {
		return &rnd, err
	}

	rnd.createBlurTexture()
	rnd.createGeometry()
	rnd.createShaders()

	rnd.window.SetGrab(true)
	//sdl.ShowCursor(0)
	return &rnd, nil
}

func (rnd *sdlRenderer) createBlurTexture() {
	rnd.glBlurTexture = gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, rnd.glBlurTexture)

	size := rnd.config.windowSize
	gl.TexImage2D(gl.TEXTURE_2D, 0, size.X, size.Y, gl.RGB, gl.UNSIGNED_BYTE, nil)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
}

func (rnd *sdlRenderer) createShaders() {
	prog := gl.CreateProgram()

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vs, vertexShader)
	gl.CompileShader(vs)

	if gl.GetShaderi(vs, gl.COMPILE_STATUS) == 0 {
		log.Panicln(gl.GetShaderInfoLog(vs))
	}

	ps := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(ps, pixelShader)
	gl.CompileShader(ps)

	if gl.GetShaderi(ps, gl.COMPILE_STATUS) == 0 {
		log.Panicln(gl.GetShaderInfoLog(ps))
	}

	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, ps)
	gl.LinkProgram(prog)

	if gl.GetProgrami(prog, gl.LINK_STATUS) <= 0 {
		log.Panicln("program linking error")
	}

	gl.DeleteShader(vs)
	gl.DeleteShader(ps)

	rnd.glProgram = prog
}

func (rnd *sdlRenderer) createGeometry() {
	squareVerticesData := []float32{
		-1, -1,
		1, -1,
		-1, 1,
		1, 1,
	}

	rnd.glSquareBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, rnd.glSquareBuffer)

	ptr := unsafe.Pointer(&squareVerticesData[0])
	gl.BufferData(gl.ARRAY_BUFFER, (*[1 << 30]byte)(ptr)[:len(squareVerticesData)*4], gl.STATIC_DRAW)

	textureUVData := []float32{
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	}

	rnd.glSquareUVBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, rnd.glSquareUVBuffer)

	ptr = unsafe.Pointer(&textureUVData[0])
	gl.BufferData(gl.ARRAY_BUFFER, (*[1 << 30]byte)(ptr)[:len(textureUVData)*4], gl.STATIC_DRAW)
}

func (rnd *sdlRenderer) renderBlurEffect() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, rnd.glBlurTexture)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.UseProgram(rnd.glProgram)

	pos := gl.GetAttribLocation(rnd.glProgram, "a_position")
	uv := gl.GetAttribLocation(rnd.glProgram, "a_uv")

	gl.BindBuffer(gl.ARRAY_BUFFER, rnd.glSquareBuffer)
	gl.VertexAttribPointer(pos, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(pos)

	gl.BindBuffer(gl.ARRAY_BUFFER, rnd.glSquareUVBuffer)
	gl.VertexAttribPointer(uv, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(uv)

	samp := gl.GetUniformLocation(rnd.glProgram, "s_texture")
	gl.Uniform1i(samp, 0)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
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
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.CULL_FACE)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	size := rnd.config.windowSize
	w, h := size.X, size.Y
	rnd.vgContext.BeginFrame(w, h, float32(w)/float32(h))

	return rnd.vgContext
}

func (rnd *sdlRenderer) Present() {
	rnd.vgContext.EndFrame()

	if blurEffect {
		gl.BindTexture(gl.TEXTURE_2D, rnd.glBlurTexture)

		size := rnd.config.windowSize
		w, h := size.X, size.Y
		gl.CopyTexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 0, 0, w, h, 0)

		gl.Disable(gl.BLEND)
		gl.Disable(gl.SCISSOR_TEST)

		rnd.renderBlurEffect()
	}

	sdl.GL_SwapWindow(rnd.window)
	if rnd.config.debug {
		checkGLError()
	}
}

func (rnd *sdlRenderer) Shutdown() {
	gl.DeleteTexture(rnd.glBlurTexture)
	gl.DeleteBuffer(rnd.glSquareBuffer)
	gl.DeleteBuffer(rnd.glSquareUVBuffer)
	gl.DeleteProgram(rnd.glProgram)

	rnd.vgContext.Delete()
	gl.ContextWatcher.OnDetach()
	sdl.GL_DeleteContext(rnd.glContext)
	rnd.window.Destroy()
}

func (rnd *sdlRenderer) SetWindowTitle(title string) {
	rnd.window.SetTitle(title)
}

func checkGLError() {
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Panicf("GL error: 0x%x\n", err)
	}
}

var vertexShader = `
	#version 120

	attribute vec4 a_position;
	attribute vec4 a_uv;
	varying vec2 v_uv;

	void main()
	{
	    gl_Position = a_position;
	    v_uv = a_uv.xy;
	}
`

var pixelShader = `
	#version 120

	#define N 4

	uniform sampler2D s_texture;
	varying vec2 v_uv;

	void main()
	{
		vec3 col = texture2D(s_texture, v_uv, 0).xyz;
		for (int i = 1; i < N; i++)
			col += texture2D(s_texture, v_uv, i).xyz;

		gl_FragColor = vec4(col / (N - 1), 1);
	}
`
