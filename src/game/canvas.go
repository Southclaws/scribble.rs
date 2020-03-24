package game

type Pixel struct{ R, G, B byte } // single pixel
type Pixels []Pixel               // collection of dimensionally contiguous pixels
type Coords [2]uint16             // two-dimensional coordinate system
type Draw struct {                // a draw instruction
	Pixel
	Coords
}
type Draws []Draw    // a sequence of draw instructions
type Canvas struct { // a canvas object
	Coords
	Pixels
}

// executes the given draw instructions against the canvas in a mutable way
func (c *Canvas) draw(d Draws) {
	w := c.Coords[0]
	for _, p := range d {
		x, y := p.Coords[0], c.Coords[1]
		idx := x + y*w
		c.Pixels[idx] = p.Pixel
	}
}
