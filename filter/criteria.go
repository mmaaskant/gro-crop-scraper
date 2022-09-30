package filter

// CriteriaBuilder simplifies the adding of children to a Criteria.
type CriteriaBuilder struct {
	criteria *Criteria
}

func NewCriteriaBuilder(c *Criteria) *CriteriaBuilder {
	return &CriteriaBuilder{
		c,
	}
}

func (cb *CriteriaBuilder) AddChild(child *Criteria) *CriteriaBuilder {
	child.Parent = cb.criteria
	cb.criteria.Child = child
	cb.criteria = child
	return cb
}

func (cb *CriteriaBuilder) Build() *Criteria {
	criteria := cb.criteria
	for parent := criteria.Parent; parent != nil; parent = parent.Parent {
		criteria = parent
	}
	return criteria
}

// Criteria defines if a set of data passes its requirements or not, and can optionally extract the matched data
// using Extractor.
type Criteria struct {
	Extractor    Extractor
	interpreters []ConditionInterpreter
	Depth        int
	Parent       *Criteria
	Child        *Criteria
}

func NewCriteria(extractor Extractor, interpreters ...ConditionInterpreter) *Criteria {
	return &Criteria{
		extractor,
		interpreters,
		0,
		nil,
		nil,
	}
}

func (c *Criteria) Match(data any) bool {
	for _, i := range c.interpreters {
		if !i.Interpret(data) {
			return false
		}
	}
	return true
}

func (c *Criteria) Previous() *Criteria {
	return c.Parent
}

func (c *Criteria) Next() *Criteria {
	return c.Child
}

func (c *Criteria) Clone() *Criteria {
	cCopy := *c
	for child := cCopy.Next(); child != nil; child = child.Next() {
		parentCopy := *child.Parent
		child.Parent = &parentCopy
		if child.Child != nil {
			childCopy := *child.Child
			child.Child = &childCopy
		}
	}
	return &cCopy
}
