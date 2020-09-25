package catalog

func (c *Catalog) FindControlById(id string) *Control {
	for i, ctrl := range c.Controls {
		if ctrl.Id == id {
			return &c.Controls[i]
		}
	}
	for i, _ := range c.Groups {
		ctrl := c.Groups[i].FindControlById(id)
		if ctrl != nil {
			return ctrl
		}
	}
	return nil
}
