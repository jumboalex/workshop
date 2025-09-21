package main

func canPlaceFlowers(flowerbed []int, n int) bool {
	/*
		if len(flowerbed) < 2 {
			if flowerbed[0] == 0 {
				return true
			} else {
				if n == 0 {
					return true
				} else {
					return false
				}

			}
		} else {
			for i := 0; i < len(flowerbed); i++ {
				if n == 0 {
					return true
				}
				if i == 0 {
					if flowerbed[i] == 0 && flowerbed[i+1] == 0 {
						flowerbed[i] = 1
						n--
					}

					continue
				}
				if i < len(flowerbed)-1 {
					if flowerbed[i-1] == 0 && flowerbed[i] == 0 && flowerbed[i+1] == 0 {
						flowerbed[i] = 1
						n--
					}

					continue
				}
				if flowerbed[i-1] == 0 && flowerbed[i] == 0 {
					flowerbed[i] = 1
					n--
				}

			}
			return n == 0
		}
	*/
	if n == 0 {
		return true
	}
	for i := 0; i < len(flowerbed); i++ {
		if flowerbed[i] == 0 {
			prev := (i == 0) || (flowerbed[i-1] == 0)
			next := (i == len(flowerbed)-1) || (flowerbed[i+1] == 0)

			if prev && next {
				flowerbed[i] = 1
				n--
				if n == 0 {
					return true
				}
			}
		}
	}
	return n == 0
}
