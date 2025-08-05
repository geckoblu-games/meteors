package game

func CheckCollision(a, b *Sprite) bool {
	aMask := a.GetAlphaMask()
	bMask := b.GetAlphaMask()

	ax1, ay1 := int(a.position.X-float64(aMask.Width)/2), int(a.position.Y-float64(aMask.Height)/2)
	ax2, ay2 := int(a.position.X+float64(aMask.Width)/2), int(a.position.Y+float64(aMask.Height)/2)

	bx1, by1 := int(b.position.X-float64(bMask.Width)/2), int(b.position.Y-float64(bMask.Height)/2)
	bx2, by2 := int(b.position.X+float64(bMask.Width)/2), int(b.position.Y+float64(bMask.Height)/2)

	// AABB rejection
	if ax1 >= bx2 || bx1 >= ax2 || ay1 >= by2 || by1 >= ay2 {
		// fmt.Printf("AABB rejection [%v %v %v %v] [%v %v %v %v]\n", ax1, ay1, ax2, ay2, bx1, by1, bx2, by2)
		return false
	}

	//fmt.Printf("AABB collision [%v %v %v %v] [%v %v %v %v]\n", ax1, ay1, ax2, ay2, bx1, by1, bx2, by2)

	// Overlapping region
	ix1 := max(ax1, bx1)
	iy1 := max(ay1, by1)
	ix2 := min(ax2, bx2)
	iy2 := min(ay2, by2)

	for y := iy1; y < iy2; y++ {
		for x := ix1; x < ix2; x++ {
			if aMask.At(x-ax1, y-ay1) && bMask.At(x-bx1, y-by1) {
				return true
			}
		}
	}

	return false
}
