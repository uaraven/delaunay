package delaunay

import (
	"math"
	"testing"
)

func TestDelaunay_circumcenter(t *testing.T) {
	p1 := Point{X: 5, Y: 7}
	p2 := Point{X: 6, Y: 6}
	p3 := Point{X: 2, Y: -2}

	tr := NewTriangle(p1, p2, p3)
	cpt := tr.Circumcenter

	if math.Abs(cpt.X-2) >= eps || math.Abs(cpt.Y-3) >= eps {
		t.Fatalf("Expected circumcenter to be (2,3), but got %v", cpt)
	}
	if math.Abs(tr.Circumradius-5) > eps {
		t.Fatalf("Expected radius to be 5, but got %3.3f", tr.Circumradius)
	}
}

func TestDelaunay_Triangulate(t *testing.T) {
	p1 := Point{X: 2, Y: 4}
	p2 := Point{X: 4, Y: 2}
	p3 := Point{X: 5, Y: 8}
	p4 := Point{X: 8, Y: 3}

	delaunay := InitDelaunay([]Point{p1, p2, p3, p4})

	triangles := delaunay.Triangulate()

	if len(triangles) != 2 {
		t.Fatalf("Expected two triangles, but got %v", triangles)
	}
	expected0 := NewTriangle(p1, p2, p3)
	expected1 := NewTriangle(p3, p4, p2)

	if triangles[0] != expected0 {
		t.Errorf("Expected triangle %v, but got %v", expected0, triangles[0])
	}

	if triangles[1] != expected1 {
		t.Errorf("Expected triangle %v, but got %v", expected1, triangles[1])
	}
}
