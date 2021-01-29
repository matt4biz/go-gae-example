package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type entry uint8

const (
	gray entry = iota
	maroon
	red
	orange
	yellow
	olive
	chartreuse
	green
	lime
	teal
	aqua
	sky
	blue
	navy
	purple
	violet
)

type paletteColor struct {
	name entry
	rgb  color.RGBA
}

func (p paletteColor) RGBA() (r, g, b, a uint32) {
	return p.rgb.RGBA()
}

var rainbow = color.Palette{
	&paletteColor{gray, color.RGBA{120, 120, 120, 255}},
	&paletteColor{maroon, color.RGBA{128, 0, 0, 255}},
	&paletteColor{red, color.RGBA{255, 0, 0, 255}},
	&paletteColor{orange, color.RGBA{255, 128, 0, 255}},
	&paletteColor{yellow, color.RGBA{255, 255, 0, 255}},
	&paletteColor{olive, color.RGBA{128, 128, 0, 255}},
	&paletteColor{chartreuse, color.RGBA{128, 255, 0, 255}},
	&paletteColor{green, color.RGBA{0, 128, 0, 255}},
	&paletteColor{lime, color.RGBA{0, 255, 0, 255}},
	&paletteColor{teal, color.RGBA{0, 128, 128, 255}},
	&paletteColor{aqua, color.RGBA{0, 255, 255, 255}},
	&paletteColor{sky, color.RGBA{0, 128, 255, 255}},
	&paletteColor{blue, color.RGBA{0, 0, 255, 255}},
	&paletteColor{navy, color.RGBA{0, 0, 128, 255}},
	&paletteColor{purple, color.RGBA{128, 0, 128, 255}},
	&paletteColor{violet, color.RGBA{128, 0, 255, 255}},
}

const (
	scale = 8
	size  = 1024
	n     = size / scale
)

var source = rand.New(rand.NewSource(time.Now().UnixNano()))

func makeRandSlice(n int) []entry {
	length := len(rainbow) - 1
	result := make([]entry, n)

	for i := range result {
		// never the default color (here gray)
		result[i] = entry(source.Intn(length) + 1)
	}

	return result
}

func paintSquare(i, k int, src []entry, img *image.Paletted) {
	// lay down a square with an outline using the default
	// color (gray; we deliberately excluded it from the data)

	for x := 0; x < scale; x++ {
		for y := 0; y < scale; y++ {
			idx := uint8(src[i])

			if x == 0 || y == 0 || x == scale-1 || y == scale-1 {
				idx = 0
			}

			img.SetColorIndex(i*scale+x, k*scale+y, idx)
		}
	}
}

func animate(out io.Writer, loop, delay int, sort func(int, []entry) int) {
	array := makeRandSlice(n)
	step := make([][]entry, n)
	anim := gif.GIF{LoopCount: loop}

	for i := 0; i < n; i++ {
		rect := image.Rect(0, 0, size, size)
		img := image.NewPaletted(rect, rainbow)

		// at each step we'll copy the array after
		// sorting it so we can draw the right view

		step[i] = make([]entry, n)

		c := sort(i, array)

		if c < 0 {
			break
		}

		copy(step[i], array)

		// now we're going to color squares based on the
		// color of paletteColors for each entry in the last step

		for k := 0; k < n; k++ {
			for id := 0; id < n; id++ {
				// we use the current step unless we're at a row and
				// column that should show the previous state

				var src []entry = step[i]

				if k < i && id <= c {
					src = step[k]
				}

				paintSquare(id, k, src, img)
			}
		}

		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim)
}

type qsort struct {
	part  func(int, int, []entry) (int, int)
	stack []int
}

// Lomuto's partition (see also Programming Pearls)
func partHigh(l, h int, A []entry) (int, int) {
	pivot := A[h]
	i := l - 1

	for j := l; j < h; j++ {
		if A[j] <= pivot {
			i++
			A[i], A[j] = A[j], A[i]
		}
	}

	i++
	A[i], A[h] = A[h], A[i]

	return i, 1
}

// Hoare's original partition
func partMiddle(l, h int, array []entry) (int, int) {
	i := (h + l) / 2
	pivot := array[i]

	for l <= h {
		for array[l] < pivot {
			l++
		}

		for array[h] > pivot {
			h--
		}

		if l <= h {
			array[l], array[h] = array[h], array[l]
			l++
			h--
		}
	}

	return l, 0
}

// Lomuto & median-of-three, but we use insertion sort
// for short subarrays (can't really see this animated)
func partInsert(l, h int, array []entry) (int, int) {
	if j := h - l + 1; j < 7 {
		for i := 0; i < j; i++ {
			insertionStep(i, array[l:h+1])
		}

		return l, 1
	}

	return partMedian(l, h, array)
}

// Lomuto using median-of-three as a pivot choice
func partMedian(l, h int, array []entry) (int, int) {
	m := (l + h) / 2

	if array[m] < array[l] {
		array[m], array[l] = array[l], array[m]
	}
	if array[l] < array[h] {
		array[h], array[l] = array[l], array[h]
	}
	if array[m] < array[h] {
		array[h], array[m] = array[m], array[h]
	}

	return partHigh(l, h, array)
}

func (q *qsort) push(l, h int) {
	q.stack = append(q.stack, l, h)
}

func (q *qsort) pop() (l, h int) {
	top := len(q.stack)

	h = q.stack[top-1]
	l = q.stack[top-2]

	q.stack = q.stack[:top-2]
	return
}

func (q *qsort) qsStep(i int, array []entry) int {
	if i == 0 {
		// we do this so the first frame is always untouched

		q.stack = make([]int, 0, len(array))
		q.stack = append(q.stack, 0, len(array)-1)
		return 0
	}

	if len(q.stack) > 1 {
		low, high := q.pop()
		pivot, off := q.part(low, high, array)

		if pivot-1 > low {
			q.push(low, pivot-1)
		}

		if pivot+off < high {
			q.push(pivot+off, high)
		}
	}

	// we do this so we can stop animation early

	if q.stack == nil {
		return -1
	} else if len(q.stack) == 0 {
		q.stack = nil
	}

	return len(array)
}

// This is the dutch-flag three-way partition based
// on picking the middle entry as the pivot; it needs
// a slightly different quicksort step (below)
func partFlag(l, h int, A []entry) (int, int) {
	p := (h + l) / 2
	pivot := A[p]

	for j := l; j <= h; {
		if A[j] < pivot {
			A[j], A[l] = A[l], A[j]
			l++
			j++
		} else if A[j] > pivot {
			A[j], A[h] = A[h], A[j]
			h--
		} else {
			j++
		}
	}

	return l, h
}

func (q *qsort) qsStepFlag(i int, array []entry) int {
	if i == 0 {
		// we do this so the first frame is always untouched

		q.stack = make([]int, 0, len(array))
		q.stack = append(q.stack, 0, len(array)-1)
		return 0
	}

	if len(q.stack) > 1 {
		low, high := q.pop()
		l, h := q.part(low, high, array)

		if l-1 > low {
			q.push(low, l-1)
		}

		if h+1 < high {
			q.push(h+1, high)
		}
	}

	// we do this so we can stop animation early

	if q.stack == nil {
		return -1
	} else if len(q.stack) == 0 {
		q.stack = nil
	}

	return len(array)
}

func qsortHigh(w http.ResponseWriter, r *http.Request) {
	q := qsort{part: partHigh}
	animate(w, getLoop(r), getDelay(r), q.qsStep)
}

func qsortMiddle(w http.ResponseWriter, r *http.Request) {
	q := qsort{part: partMiddle}
	animate(w, getLoop(r), getDelay(r), q.qsStep)
}

func qsortMedian(w http.ResponseWriter, r *http.Request) {
	q := qsort{part: partMedian}
	animate(w, getLoop(r), getDelay(r), q.qsStep)
}

func qsortInsert(w http.ResponseWriter, r *http.Request) {
	q := qsort{part: partInsert}
	animate(w, getLoop(r), getDelay(r), q.qsStep)
}

func qsortFlag(w http.ResponseWriter, r *http.Request) {
	q := qsort{part: partFlag}
	animate(w, getLoop(r), getDelay(r), q.qsStepFlag)
}

func insertionStep(i int, array []entry) int {
	for j := i; j > 0 && array[j] < array[j-1]; j-- {
		array[j], array[j-1] = array[j-1], array[j]
	}

	return i
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	animate(w, getLoop(r), getDelay(r), insertionStep)
}
