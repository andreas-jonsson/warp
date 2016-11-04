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
	"os"
	"os/user"
	"path"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var keyMapping = map[sdl.Keycode]int{
	sdl.K_UP:     KeyUp,
	sdl.K_DOWN:   KeyDown,
	sdl.K_LEFT:   KeyLeft,
	sdl.K_RIGHT:  KeyRight,
	sdl.K_ESCAPE: KeyEsc,
	sdl.K_RETURN: KeyReturn,
}

var mouseMapping = map[int]int{
	sdl.MOUSEBUTTONDOWN: MouseButtonDown,
	sdl.MOUSEBUTTONUP:   MouseButtonUp,
	sdl.MOUSEWHEEL:      MouseWheel,
}

func init() {
	runtime.LockOSThread()

	if runtime.GOOS == "windows" {
		ConfigPath = path.Join(os.Getenv("LOCALAPPDATA"), "Warp")
	} else {
		if usr, err := user.Current(); err == nil {
			ConfigPath = path.Join(usr.HomeDir, ".config", "warp")
		}
	}

	ConfigPath = path.Clean(ConfigPath)
	os.MkdirAll(ConfigPath, 0755)
}

func Init() error {
	idCounter = 0
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	return nil
}

func Shutdown() {
	sdl.Quit()
}

func Mouse() MouseState {
	x, y, buttons := sdl.GetMouseState()

	window := sdl.GetMouseFocus()
	if window == nil {
		return MouseState{}
	}

	left := (buttons & sdl.ButtonLMask()) != 0
	middle := (buttons & sdl.ButtonMMask()) != 0
	right := (buttons & sdl.ButtonRMask()) != 0

	return MouseState{X: x, Y: y, Buttons: [3]bool{left, middle, right}}
}

func PollEvent() Event {
	event := sdl.PollEvent()
	if event == nil {
		return nil
	}

	switch t := event.(type) {
	case *sdl.QuitEvent:
		return &QuitEvent{}
	case *sdl.KeyUpEvent:
		ev := &KeyUpEvent{}
		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
			ev.Rune = rune(t.Keysym.Unicode)
		} else {
			ev.Key = KeyUnknown
		}
		return ev
	case *sdl.KeyDownEvent:
		ev := &KeyDownEvent{}
		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
			ev.Rune = rune(t.Keysym.Unicode)
		} else {
			ev.Key = KeyUnknown
		}
		return ev
	case *sdl.MouseButtonEvent:
		ev := &MouseButtonEvent{
			Button: int(t.Button),
			X:      int(t.X),
			Y:      int(t.Y),
		}

		switch t.Type {
		case sdl.MOUSEBUTTONDOWN:
			ev.Type = MouseButtonDown
		case sdl.MOUSEBUTTONUP:
			ev.Type = MouseButtonUp
		case sdl.MOUSEWHEEL:
			ev.Type = MouseWheel
		}
		return ev
	case *sdl.MouseMotionEvent:
		return &MouseMotionEvent{
			X:    int(t.X),
			Y:    int(t.Y),
			XRel: int(t.XRel),
			YRel: int(t.YRel),
		}
	}

	return nil
}
