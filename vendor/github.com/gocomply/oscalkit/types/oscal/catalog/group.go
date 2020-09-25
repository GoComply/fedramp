package catalog

func (g *Group) FindControlById(id string) *Control {
	for i, ctrl := range g.Controls {
		if ctrl.Id == id {
			return &g.Controls[i]
		}
	}
	for i, _ := range g.Groups {
		ctrl := g.Groups[i].FindControlById(id)
		if ctrl != nil {
			return ctrl
		}
	}
	return nil
}
