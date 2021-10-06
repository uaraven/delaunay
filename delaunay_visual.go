package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	d "github.com/uaraven/delaunay/delaunay"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

var tr *d.Delaunay
var index int

func update() {
	if rl.IsKeyPressed(rl.KeySpace) {
		if index < tr.PointCount() {
			tr.TriangulationStep(tr.Point(index))
			index += 1
		} else if index == tr.PointCount() {
			tr.Finalize()
		}
	}
}

func draw() {
	// draw points
	for i, p := range tr.Points() {
		if i < index {
			rl.DrawCircle(int32(p.X), int32(p.Y), 3, rl.Red)
		} else {
			rl.DrawCircle(int32(p.X), int32(p.Y), 5, rl.SkyBlue)
		}
		rl.DrawText(p.String(), int32(p.X)+5, int32(p.Y), 12, rl.DarkBlue)
	}
	y := -200
	trg := tr.Triangles()
	for _, t := range trg {
		var color rl.Color
		if t.UsesAnyOfVertices(tr.SupertriangleVertices()) {
			color = rl.Blue
		} else {
			color = rl.Green
		}
		rl.DrawCircleLines(int32(t.Circumcenter.X), int32(t.Circumcenter.Y), float32(t.Circumradius), rl.ColorAlpha(rl.Beige, 0.5))
		for _, e := range t.Edges() {
			rl.DrawLine(int32(e.P1.X), int32(e.P1.Y), int32(e.P2.X), int32(e.P2.Y), color)
		}
		rl.DrawText(fmt.Sprintf("%v", t), -180, int32(y), 12, rl.Black)
		y += 15
	}

}

func init() {
	// p1 := d.Point{X: 200 / 2, Y: 400 / 2}
	// p2 := d.Point{X: 400 / 2, Y: 200 / 2}
	// p3 := d.Point{X: 500 / 2, Y: 800 / 2}
	// p4 := d.Point{X: 800 / 2, Y: 300 / 2}

	// tr = d.InitDelaunay([]d.Point{p1, p2, p3, p4})
	tr = d.InitDelaunay([]d.Point{
		d.NewPoint(0, 40),
		d.NewPoint(40, 0),
		d.NewPoint(50, 100),
		d.NewPoint(90, 10),
		d.NewPoint(100, 60),
	})
	index = 0
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Delaunay")

	cam := rl.Camera2D{}
	cam.Zoom = 1
	cam.Target = rl.Vector2{X: -300, Y: -400}

	for !rl.WindowShouldClose() {
		update()
		rl.BeginDrawing()
		rl.BeginMode2D(cam)
		rl.ClearBackground(rl.White)
		draw()
		rl.EndMode2D()
		if index < tr.PointCount() {
			rl.DrawText("Press space for the next step", 20, screenHeight-22, 20, rl.DarkBlue)
		}
		rl.EndDrawing()
	}
}
