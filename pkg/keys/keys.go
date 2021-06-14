package keys

import "github.com/hajimehoshi/ebiten/v2/inpututil"

func IsAnyKeyJustPressed() bool {
	for _, k := range inpututil.PressedKeys() {
		if inpututil.IsKeyJustPressed(k) {
			return true
		}
	}
	return false
}
