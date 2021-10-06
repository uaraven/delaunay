package delaunay

import (
	"fmt"
	"math"
	"sort"
)

const eps = 1e-6

type Point struct {
	X float64
	Y float64
}

func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%3.1f, %3.1f)", p.X, p.Y)
}

func (p Point) CompareTo(other Point) int {
	if p == other {
		return 0
	}
	if p.X > other.X || (math.Abs(p.X-other.X) < eps && p.Y > other.Y) {
		return 1
	} else {
		return 0
	}
}

func (p Point) Distance(other Point) float64 {
	dx := other.X - p.X
	dy := other.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (p Point) AngleTo(other Point) int {
	return int(360.0*math.Atan2(other.X-p.X, other.Y-p.Y)/math.Pi+360) % 360
}

type Edge struct {
	P1 Point
	P2 Point
}

func NewEdge(p1, p2 Point) Edge {
	if p2.CompareTo(p1) > 0 {
		return Edge{P1: p2, P2: p1}
	} else {
		return Edge{P1: p1, P2: p2}
	}
}

func (e Edge) String() string {
	return fmt.Sprintf("%v - %v", e.P1, e.P2)
}

type Triangle struct {
	P1           Point
	P2           Point
	P3           Point
	Circumcenter Point
	Circumradius float64
}

func NewTriangle(p1, p2, p3 Point) Triangle {
	center := NewPoint((p1.X+p2.X+p3.X)/3, (p1.Y+p2.Y+p3.Y)/3)
	vertices := []Point{p1, p2, p3}

	sort.Slice(vertices, func(i, j int) bool {
		return center.AngleTo(vertices[i]) > center.AngleTo(vertices[j])
	})

	t := Triangle{P1: vertices[0], P2: vertices[1], P3: vertices[2]}
	t.Circumcenter, t.Circumradius = t.circumcenter()
	return t
}

func (t Triangle) circumcenter() (Point, float64) {
	if math.Abs(t.P1.Y-t.P2.Y) < eps && math.Abs(t.P2.Y-t.P3.Y) < eps {
		return NewPoint(0, 0), 0
	}
	e1 := t.Edge(0)
	e2 := t.Edge(1)
	if e1.IsHorizontal() {
		e1 = t.Edge(2)
	} else if e2.IsHorizontal() {
		e2 = t.Edge(2)
	}

	c1 := e1.Center()
	c2 := e2.Center()

	m1 := -1 / ((e1.P2.Y - e1.P1.Y) / (e1.P2.X - e1.P1.X))
	m2 := -1 / ((e2.P2.Y - e2.P1.Y) / (e2.P2.X - e2.P1.X))

	b1 := c1.Y - m1*c1.X
	b2 := c2.Y - m2*c2.X

	cx := (b2 - b1) / (m1 - m2)
	cy := m1*cx + b1

	center := NewPoint(cx, cy)

	r := center.Distance(e1.P1)

	return center, r
}

func (t Triangle) String() string {
	return fmt.Sprintf("[%v - %v - %v]", t.P1, t.P2, t.P3)
}

func (e Edge) Center() Point {
	return NewPoint(
		e.P1.X+(e.P2.X-e.P1.X)/2,
		e.P1.Y+(e.P2.Y-e.P1.Y)/2,
	)
}

func (e Edge) IsHorizontal() bool {
	return math.Abs(e.P1.Y-e.P2.Y) < eps
}

func (t Triangle) Edge(index int) Edge {
	switch index % 3 {
	case 0:
		return NewEdge(t.P1, t.P2)
	case 1:
		return NewEdge(t.P2, t.P3)
	case 2:
		return NewEdge(t.P3, t.P1)
	}
	panic(fmt.Errorf("invalid index: %d", index))
}

func (t Triangle) Edges() []Edge {
	return []Edge{t.Edge(0), t.Edge(1), t.Edge(2)}
}

func (t Triangle) Vertices() []Point {
	return []Point{t.P1, t.P2, t.P3}
}

func (t Triangle) UsesAnyOf(vertices []Point) bool {
	verts := make(map[Point]bool, len(vertices))
	for _, v := range vertices {
		verts[v] = true
	}
	for _, tv := range t.Vertices() {
		if _, ok := verts[tv]; ok {
			return true
		}
	}
	return false
}

func (t Triangle) UsesAnyOfVertices(vertices map[Point]bool) bool {
	for _, tv := range t.Vertices() {
		if _, ok := vertices[tv]; ok {
			return true
		}
	}
	return false
}

func reduce(points []Point, initial float64, reducer func(float64, Point) float64) float64 {
	result := initial
	for _, pt := range points {
		result = reducer(result, pt)
	}
	return result
}

func reduceMinX(prev float64, pt Point) float64 {
	if pt.X < prev {
		return pt.X
	} else {
		return prev
	}
}

func reduceMinY(prev float64, pt Point) float64 {
	if pt.Y < prev {
		return pt.Y
	} else {
		return prev
	}
}

func reduceMaxX(prev float64, pt Point) float64 {
	if pt.X > prev {
		return pt.X
	} else {
		return prev
	}
}

func reduceMaxY(prev float64, pt Point) float64 {
	if pt.Y > prev {
		return pt.Y
	} else {
		return prev
	}
}

const offset = 2

func boundingTriangle(points []Point) Triangle {
	minX := reduce(points, math.MaxInt64, reduceMinX)
	maxX := reduce(points, math.MinInt64, reduceMaxX)
	minY := reduce(points, math.MaxInt64, reduceMinY)
	maxY := reduce(points, math.MinInt64, reduceMaxY)

	bottomLeft := NewPoint(minX-offset, maxY+offset)
	topRight := NewPoint(maxX+offset, minY-offset)

	m := -1 / ((topRight.Y - bottomLeft.Y) / (topRight.X - bottomLeft.X - offset))
	b := topRight.Y - m*topRight.X

	topLeft := NewPoint(bottomLeft.X, m*bottomLeft.X+b)
	bottomRight := NewPoint((bottomLeft.Y-b)/m, bottomLeft.Y)

	return NewTriangle(topLeft, bottomRight, NewPoint(minX-offset*5, maxY+offset*5))
}

type Delaunay struct {
	triangles             []Triangle
	supertriangle         Triangle
	points                []Point
	superTriangleVertices map[Point]bool
}

func (d *Delaunay) Finalize() {
	result := make([]Triangle, 0)
	for _, t := range d.triangles {
		if !t.UsesAnyOfVertices(d.superTriangleVertices) {
			result = append(result, t)
		}
	}
	d.triangles = result
}

func (d *Delaunay) Point(index int) Point {
	return d.points[index]
}

func (d *Delaunay) PointCount() int {
	return len(d.points)
}

func (d *Delaunay) Points() []Point {
	return d.points
}

func (d *Delaunay) Triangles() []Triangle {
	return d.triangles
}

func (d *Delaunay) SupertriangleVertices() map[Point]bool {
	return d.superTriangleVertices
}

// TriangulationStep performs one step of the triangulation algorithm for a given point
func (d *Delaunay) TriangulationStep(v Point) {
	tidx := 0
	edges := make(map[Edge]int)
	for tidx < len(d.triangles) {
		t := d.triangles[tidx]
		if v.X-t.Circumcenter.X <= t.Circumradius { // x component of the distance from the current point to the circumcircle center is greater than the circumcircle radius, that triangle need never be considered for the later points
			if t.Circumradius > 0 && v.Distance(t.Circumcenter) <= t.Circumradius { // point is inside the circumcircle of the triangle
				for _, e := range t.Edges() {
					edges[e] += 1
				}
				d.triangles = append(d.triangles[:tidx], d.triangles[tidx+1:]...)
			} else {
				tidx++
			}
		} else {
			tidx++
		}
	}
	trid := make([]Triangle, 0, len(edges))
	for e, count := range edges {
		if count == 1 { // only process edges that were not shared by removed triangles
			t := NewTriangle(v, e.P1, e.P2)
			trid = append(trid, t)
		}
	}
	d.triangles = append(d.triangles, trid...)
}

// InitDelaunay Initializes data structures for Delaunay triangulation of the set of points
// Call Triangulate() on the returned struct to perform triangulation
func InitDelaunay(points []Point) *Delaunay {
	triangles := make([]Triangle, 0)
	supertriangle := boundingTriangle(points)
	points = append(points, supertriangle.Vertices()...)
	triangles = append(triangles, supertriangle)
	verts := make(map[Point]bool, len(supertriangle.Vertices()))
	for _, v := range supertriangle.Vertices() {
		verts[v] = true
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].X < points[j].X
	})
	return &Delaunay{
		triangles:             triangles,
		supertriangle:         supertriangle,
		points:                points,
		superTriangleVertices: verts,
	}
}

// Triangulate performs Delaunay triangulation using Bowyer–Watson algorithm as described by Paul Bourke
// See http://paulbourke.net/papers/triangulate/
func (d *Delaunay) Triangulate() []Triangle {
	for _, v := range d.points {
		d.TriangulationStep(v)
	}
	d.Finalize()
	return d.triangles
}
