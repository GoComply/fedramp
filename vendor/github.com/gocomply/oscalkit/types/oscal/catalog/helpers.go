package catalog

// ControlOpts to generate controls
type ControlOpts struct {
	Params   []Param
	Parts    []Part
	Controls []Control
}

// NewPart creates a new part
func NewPart(id, title, narrative string) Part {
	return Part{
		Id:    id,
		Title: Title(title),
		Prose: &Prose{Raw: narrative},
	}
}

// NewControl creates a new control
func NewControl(id, title string, opts *ControlOpts) Control {
	ctrl := Control{
		Id:    id,
		Title: Title(title),
	}
	if opts != nil {
		ctrl.Controls = opts.Controls
		ctrl.Parts = opts.Parts
		ctrl.Parameters = opts.Params
	}
	return ctrl
}
