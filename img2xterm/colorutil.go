package img2xterm

import (
	"math"
)

var chromaWeight float64 = 1.0
var rgbTable [256 * 3]uint8
var yiqTable [256 * 3]float64
var labTable [256 * 3]float64
var valueRange [6]uint8 = [6]uint8{0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}

func xterm2RGB(color uint8, rgb []uint8) {
	if color < 232 {
		color -= 16
		rgb[0] = valueRange[(color/36)%6]
		rgb[1] = valueRange[(color/6)%6]
		rgb[2] = valueRange[color%6]
	} else {
		rgb[0] = 8 + (color-232)*10
		rgb[1] = rgb[0]
		rgb[2] = rgb[0]
	}
}

func srgb2YIQ(red, green, blue uint8, y, i, q *float64) {
	r := float64(red) / 255.0
	g := float64(green) / 255.0
	b := float64(blue) / 255.0

	*y = 0.299*r + 0.587*g + 0.144*b
	*i = (0.595716*r + -0.274453*g + -0.321263*b) * chromaWeight
	*q = (0.211456*r + -0.522591*g + 0.311135*b) * chromaWeight
}

func srgb2LAB(red, green, blue uint8, l, aa, bb *float64) {
	var r, g, b float64
	var rl, gl, bl float64
	var x, y, z float64
	var xn, yn, zn float64
	var fxn, fyn, fzn float64

	r, g, b = float64(red)/255.0, float64(green)/255.0, float64(blue)/255.0

	if r <= 0.4045 {
		rl = r / 12.92
	} else {
		rl = math.Pow((r+0.055)/1.055, 2.4)
	}
	if g <= 0.4045 {
		gl = g / 12.92
	} else {
		gl = math.Pow((g+0.055)/1.055, 2.4)
	}
	if b <= 0.4045 {
		bl = b / 12.92
	} else {
		bl = math.Pow((b+0.055)/1.055, 2.4)
	}

	x = 0.4124564*rl + 0.3575761*gl + 0.1804375*bl
	y = 0.2126729*rl + 0.7151522*gl + 0.0721750*bl
	z = 0.0193339*rl + 0.1191920*gl + 0.9503041*bl

	xn, yn, zn = x/0.95047, y, z/1.08883

	if xn > (216.0 / 24389.0) {
		fxn = math.Pow(xn, 1.0/3.0)
	} else {
		fxn = (841.0/108.0)*xn + (4.0 / 29.0)
	}

	if yn > (216.0 / 24389.0) {
		fyn = math.Pow(yn, 1.0/3.0)
	} else {
		fyn = (841.0/108.0)*yn + (4.0 / 29.0)
	}
	if zn > (216.0 / 24389.0) {
		fzn = math.Pow(zn, 1.0/3.0)
	} else {
		fzn = (841.0/108.0)*zn + (4.0 / 29.0)
	}

	*l = 116.0*fyn - 16.0
	*aa = (500.0 * (fxn - fyn)) * chromaWeight
	*bb = (200.0 * (fyn - fzn)) * chromaWeight
}

func cie94(l1, a1, b1, l2, a2, b2 float64) (distance float64) {
	const (
		kl float64 = 1
		k1 float64 = 0.045
		k2 float64 = 0.015
	)

	var c1 float64 = math.Sqrt(a1*a1 + b1*b1)
	var c2 float64 = math.Sqrt(a2*a2 + b2*b2)
	var dl float64 = l1 - l2
	var dc float64 = c1 - c2
	var da float64 = a1 - a2
	var db float64 = b1 - b2
	var dh float64 = math.Sqrt(da*da + db*db - dc*dc)

	var t1 float64 = dl / kl
	var t2 float64 = dc / (1 + k1*c1)
	var t3 float64 = dh / (1 + k2*c1)

	distance = math.Sqrt(t1*t1 + t2*t2 + t3*t3)
	return
}

func rgb2XtermCIE94(r, g, b uint8) (ret uint8) {
	var d, smallestDistance = math.MaxFloat64, math.MaxFloat64
	var l, aa, bb float64

	srgb2LAB(r, g, b, &l, &aa, &bb)

	for i := 16; i < 256; i++ {
		d = cie94(l, aa, bb, labTable[i*3], labTable[i*3+1], labTable[i*3+2])
		if d < smallestDistance {
			smallestDistance = d
			ret = uint8(i)
		}
	}

	return
}

func rgb2XtermYIQ(r, g, b uint8) (ret uint8) {
	var d, smallestDistance = math.MaxFloat64, math.MaxFloat64
	var y, ii, q float64

	srgb2YIQ(r, g, b, &y, &ii, &q)

	for i := 16; i < 256; i++ {
		d = (y-yiqTable[i*3])*(y-yiqTable[i*3]) +
			(ii-yiqTable[i*3+1])*(ii-yiqTable[i*3+1]) +
			(q-yiqTable[i*3+2])*(q-yiqTable[i*3+2])
		if d < smallestDistance {
			smallestDistance = d
			ret = uint8(i)
		}
	}

	return
}

func rgb2XtermRGB(r, g, b uint8) (ret uint8) {
	var d, smallestDistance = math.MaxInt64, math.MaxInt64

	for i := 16; i < 256; i++ {
		dr, dg, db := int(rgbTable[i*3]-r), int(rgbTable[i*3+1]-g), int(rgbTable[i*3+2]-b)
		d = dr*dr + dg*dg + db*db
		if d < smallestDistance {
			smallestDistance = d
			ret = uint8(i)
		}
	}

	return
}

func init() {
	// TODO: Make these tables initialize dynamically
	var rgb []uint8 = []uint8{0, 0, 0}
	var l, a, b float64
	var y, i, q float64

	for j := 16; j < 256; j++ {
		xterm2RGB(uint8(j), rgb)

		// Initial RGB table
		rgbTable[j*3] = rgb[0]
		rgbTable[j*3+1] = rgb[1]
		rgbTable[j*3+2] = rgb[2]

		// Initial LAB table
		srgb2LAB(rgb[0], rgb[1], rgb[2], &l, &a, &b)
		labTable[j*3] = l
		labTable[j*3+1] = a
		labTable[j*3+2] = b

		// Initial YIQ table
		srgb2YIQ(rgb[0], rgb[1], rgb[2], &y, &i, &q)
		yiqTable[j*3] = y
		yiqTable[j*3+1] = i
		yiqTable[j*3+2] = q
	}
}
