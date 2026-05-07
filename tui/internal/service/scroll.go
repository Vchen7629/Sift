package service

import (
	"charm.land/bubbles/v2/viewport"
)


func ScrollToFocused(vp *viewport.Model, focusedIndex, cardHeight int) {
	// this is to calculate the lines the curr focused card occupies
	itemTop    := focusedIndex * cardHeight                                                                                                       
	itemBottom := itemTop + cardHeight

	viewTop    := vp.YOffset()
	viewBottom := viewTop + vp.Height()

	if itemBottom > viewBottom {
		// item went below visible area, scroll down just enough
		vp.SetYOffset(itemBottom - vp.Height())
	} else if itemTop < viewTop {
		// item went above visible area, scroll up just enough
		vp.SetYOffset(itemTop)
	}
}