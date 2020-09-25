package catalog

func (c *Control) FindParamById(id string) *Param {
	for i, param := range c.Parameters {
		if param.Id == id {
			return &c.Parameters[i]
		}
	}
	return nil
}
