package main

import "math"

type vector struct {
	X, Y, Z float64
}

func (v vector) Add(ov vector) vector { return vector{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z} }

func (v vector) Sub(ov vector) vector { return vector{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z} }

func (v vector) Mul(m float64) vector { return vector{m * v.X, m * v.Y, m * v.Z} }

func (v vector) Len() float64 { return math.Sqrt(v.Dot(v)) }

func (v vector) Dot(ov vector) float64 { return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z }

func (v vector) Unit() vector {
	vdot := v.Dot(v)
	if vdot == 0 {
		return vector{0, 0, 0}
	}
	return v.Mul(1 / math.Sqrt(vdot))
}
