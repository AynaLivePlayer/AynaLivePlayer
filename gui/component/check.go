package component

import "fyne.io/fyne/v2/widget"

func NewCheckOneWayBinding(name string, val *bool, checked bool) *widget.Check {
	check := widget.NewCheck(name, func(b bool) {
		*val = b
	})
	check.SetChecked(checked)
	return check
}
