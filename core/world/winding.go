package world

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

type Winding struct {
	Points []*math32.Vector3
}

func NewWinding(points int) *Winding {
	w := &Winding{
		Points: make([]*math32.Vector3, points),
	}

	for i := range w.Points {
		w.Points[i] = &math32.Vector3{}
	}

	return w
}

const splitEpsilon = 0.01

const (
	splitFront = 0
	splitBack  = 1
	splitOn    = 2
)

func (w *Winding) Clip(split *Plane) {
	// Figure out which side of the split
	// each point of this winding is on

	// This is taken completely wholesale from the valve sdk
	// look in brushops.cpp

	// Counts of how many points are on which side
	windingLength := len(w.Points) + 1

	counts := make([]int, 3)

	// Which side the point is on
	sides := make([]int, windingLength)

	// Distance from the split for each point
	dists := make([]float32, windingLength)

	for i, point := range w.Points {
		dot := point.Dot(&split.Normal)
		dot -= split.Dist

		dists[i] = dot

		if dot > splitEpsilon {
			sides[i] = splitFront
		} else if dot < -splitEpsilon {
			sides[i] = splitBack
		} else {
			sides[i] = splitOn
		}

		counts[sides[i]]++
	}

	sides[len(w.Points)] = sides[0]
	dists[len(w.Points)] = dists[0]

	if counts[splitFront] == 0 && counts[splitBack] == 0 {
		// Nothing to split (everything was on the plane)
		fmt.Println("All on...")
		return
	}

	if counts[splitFront] == 0 {
		// Everything was behind this plane
		// so we no longer have any points
		fmt.Println("All behind...")
		*w = *NewWinding(0)
		return
	}

	if counts[splitBack] == 0 {
		// Nothing was behind the split
		// so nothing to change
		fmt.Println("All in front...")
		return
	}

	maxPoints := len(w.Points) + 4
	numPoints := 0

	newWinding := NewWinding(maxPoints)

	for i, point := range w.Points {
		mid := newWinding.Points[numPoints]

		if sides[i] == splitOn {
			*mid = *point
			numPoints++
			continue
		}

		if sides[i] == splitFront {
			*mid = *point
			numPoints++
			mid = newWinding.Points[numPoints]
		}

		if sides[i+1] == splitOn || sides[i+1] == sides[i] {
			continue
		}

		// Generate a split point
		p2 := w.Points[0]
		if i != len(w.Points)-1 {
			p2 = w.Points[i+1]
		}

		numPoints++

		dot := dists[i] / (dists[i] - dists[i+1])

		for j := 0; j < 3; j++ {
			// avoid round off error when possible
			comp := point.Component(j)

			if split.Normal.Component(j) == 1 {
				mid.SetComponent(j, split.Dist)
			} else if split.Normal.Component(j) == -1 {
				mid.SetComponent(j, -split.Dist)
			} else {
				mid.SetComponent(j, comp+dot*(p2.Component(j)-comp))
			}
		}
	}

	w.Points = newWinding.Points[:numPoints]
}

// Source defines this constant as sqrt(3) * 2 * 16584
// its the max diagonal length a map could possibly be
const maxTrace = 56755.8408624

func VectorMAInline(start, direction, dest *math32.Vector3, scale float32) {
	dest.X = start.X + direction.X*scale
	dest.Y = start.Y + direction.Y*scale
	dest.Z = start.Z + direction.Z*scale
}

func CreateWindingFromPlane(p *Plane) *Winding {
	// https://github.com/emily33901/HammerFromScratch/blob/a0f669718a70632138545fd1a5a493b8299221a0/hammer/brushops.cpp

	// Find the major axis
	up := &math32.Vector3{}

	normalArray := []float32{p.Normal.X, p.Normal.Y, p.Normal.Z}
	max := math32.Abs(normalArray[0])
	idx := 0

	for i := 1; i < 3; i++ {
		v := math32.Abs(normalArray[i])
		if v > max {
			idx = i
		}
	}

	if idx == -1 {
		panic("No major axis found...")
	}

	if idx == 0 || idx == 1 {
		up.Z = 1
	} else {
		up.X = 1
	}

	// If X or Y are greater than Z
	if math32.Abs(p.Normal.X) < math32.Abs(p.Normal.Z) || math32.Abs(p.Normal.Y) < math32.Abs(p.Normal.Z) {
		up.Z = 1
	} else {
		// Z must be the largest
		up.X = 1
	}

	v := up.Dot(&p.Normal)

	VectorMAInline(up, &p.Normal, up, -v)
	up.Normalize()

	org := p.Normal.Clone().MultiplyScalar(p.Dist)

	right := up.Clone().Cross(&p.Normal).MultiplyScalar(maxTrace)
	up = up.MultiplyScalar(maxTrace)

	w := NewWinding(4)

	w.Points[0] = org.Clone().Sub(right).Add(up)
	w.Points[1] = org.Clone().Add(right).Add(up)
	w.Points[2] = org.Clone().Add(right).Sub(up)
	w.Points[3] = org.Clone().Sub(right).Sub(up)

	return w
}
