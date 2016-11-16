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
	"image"
	_ "image/jpeg"
	"log"

	"github.com/mode13/nanovgo"
	"github.com/mode13/warp/data"
	"github.com/mode13/warp/game/entity"
	"github.com/mode13/warp/platform"
	"github.com/ungerik/go3d/vec2"
	"github.com/ungerik/go3d/vec3"
)

type Universe struct {
	entities          map[uint64]entity.Entity
	backgroundImage   image.Image
	backgroundImageID int
	cameraPos         vec3.T
	tick              float64
}

func NewUniverse() *Universe {
	var spaceImage image.Image

	fp, err := data.FS.Open("space.jpg")
	if err == nil {
		defer fp.Close()
		spaceImage, _, err = image.Decode(fp)
	}

	if err != nil {
		log.Println(err)
		spaceImage = image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	}

	return &Universe{
		entities:          make(map[uint64]entity.Entity),
		backgroundImage:   spaceImage,
		backgroundImageID: -1,
	}
}

func (uni *Universe) SpawnEntity(ty string, owner int) entity.Entity {
	entity := entity.NewEntity(ty, platform.NewId64(), owner)
	uni.entities[entity.Id()] = entity
	return entity
}

func (uni *Universe) FindAll(pos vec2.T, rad float32, filter uint32) []entity.Entity {
	return nil
}

func (uni *Universe) CameraPosition() vec3.T {
	return uni.cameraPos
}

func (uni *Universe) Bounds() image.Rectangle {
	return uni.backgroundImage.Bounds()
}

func (uni *Universe) Update(dt float64, cameraPos vec3.T) error {
	uni.cameraPos = cameraPos
	uni.tick += dt

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
	if uni.backgroundImageID < 0 {
		uni.backgroundImageID = ctx.CreateImageFromGoImage(0, uni.backgroundImage)
	}

	imgSize := uni.backgroundImage.Bounds().Size()
	imgSizeX := float32(imgSize.X)
	imgSizeY := float32(imgSize.Y)

	imgPaint := nanovgo.ImagePattern(0, 0, imgSizeX, imgSizeY, 0, uni.backgroundImageID, 1)

	ctx.Scissor(uni.cameraPos[0]-1, uni.cameraPos[1]-1, imgSizeX+2, imgSizeY+2)

	const parallaxEffect = 0.8

	ctx.Save()
	ctx.Translate(uni.cameraPos[0]*parallaxEffect, uni.cameraPos[1]*parallaxEffect)

	ctx.BeginPath()
	ctx.SetFillPaint(imgPaint)
	ctx.Rect(0, 0, imgSizeX, imgSizeY)
	ctx.Fill()

	ctx.Restore()
	ctx.Translate(uni.cameraPos[0], uni.cameraPos[1])

	// Draw border

	ctx.BeginPath()
	ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 255, 255))
	ctx.SetStrokeWidth(1)
	ctx.Rect(0, 0, imgSizeX, imgSizeY)
	ctx.Stroke()

	// Draw grid

	const gridStep = 75
	for x := 0; x <= imgSize.X; x += gridStep {
		if x%(5*gridStep) == 0 {
			ctx.SetStrokeWidth(2)
			ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 255, 255))
		} else {
			ctx.SetStrokeWidth(1)
			ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 200, 75))
		}

		ctx.BeginPath()
		ctx.MoveTo(float32(x), 0)
		ctx.LineTo(float32(x), imgSizeY)
		ctx.Stroke()
	}

	for y := 0; y <= imgSize.Y; y += gridStep {
		if y%(5*gridStep) == 0 {
			ctx.SetStrokeWidth(2)
			ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 255, 255))
		} else {
			ctx.SetStrokeWidth(1)
			ctx.SetStrokeColor(nanovgo.RGBA(0, 0, 200, 75))
		}

		ctx.BeginPath()
		ctx.MoveTo(0, float32(y))
		ctx.LineTo(imgSizeX, float32(y))
		ctx.Stroke()
	}

	for _, entity := range uni.entities {
		if err := entity.Render(ctx); err != nil {
			return err
		}
	}
	return nil
}
