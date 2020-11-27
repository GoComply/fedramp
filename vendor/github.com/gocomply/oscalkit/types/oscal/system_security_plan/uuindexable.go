package system_security_plan

func (bc *ByComponent) SetUuid(uuid string) {
	bc.Uuid = uuid
}
func (c *Component) SetUuid(uuid string) {
	c.Uuid = uuid
}

func (ir *ImplementedRequirement) SetUuid(uuid string) {
	ir.Uuid = uuid
}

func (st *Statement) SetUuid(uuid string) {
	st.Uuid = uuid
}

func (ssp *SystemSecurityPlan) SetUuid(uuid string) {
	ssp.Uuid = uuid
}

func (u *User) SetUuid(uuid string) {
	u.Uuid = uuid
}
