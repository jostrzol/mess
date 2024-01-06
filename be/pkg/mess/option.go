package mess

type Choice struct {
	Message         string
	NextChoices     []*Choice
	OptionGenerator OptionGenerator
}

type OptionGenerator interface {
	GenerateOptionData() IOptionData
}

func (c *Choice) GenerateOptions() *OptionNode {
	if c == nil {
		return nil
	}

	optionData := c.OptionGenerator.GenerateOptionData()
	var children []*OptionNode
	for _, choice := range c.NextChoices {
		child := choice.GenerateOptions()
		if child.Len() != 0 {
			children = append(children, child)
		}
	}
	if children != nil {
		optionData.setLeavesChildren(children)
	}
	return &OptionNode{
		Message: c.Message,
		Data:    optionData,
	}
}

type OptionNode struct {
	Message string
	Data    IOptionData
}

func EmptyOptionNode() *OptionNode {
	return &OptionNode{
		Message: "<empty>",
		Data:    nil,
	}
}

func (n *OptionNode) Len() int {
	if n == nil || n.Data == nil {
		return 0
	}
	return n.Data.len()
}

func (n *OptionNode) Accept(visitor OptionDataVisitor) {
	n.Data.accept(n.Message, visitor)
}

func (n *OptionNode) AllRoutes() (result []Route) {
	n.FilterRoutes(func(route Route) bool {
		result = append(result, route)
		return false
	})
	return result
}

func (n *OptionNode) FilterRoutes(predicate func(Route) bool) *OptionNode {
	if n == nil {
		if predicate(Route{}) {
			return nil
		}
		return EmptyOptionNode()
	}
	return n.filterRoutes(nil, predicate)
}

func (n *OptionNode) filterRoutes(parentRoute Route, predicate func(Route) bool) *OptionNode {
	if n.Data == nil {
		return EmptyOptionNode()
	}
	newData := n.Data.filter(parentRoute, predicate)

	return &OptionNode{
		Message: n.Message,
		Data:    newData,
	}
}

type IOptionData interface {
	accept(message string, visitor OptionDataVisitor)
	setLeavesChildren(children []*OptionNode)
	filter(parentRoute Route, predicate func(Route) bool) IOptionData
	len() int
}

type OptionData[T Option] []*OptionDatum[T]

func (d OptionData[T]) setLeavesChildren(children []*OptionNode) {
	for _, datum := range d {
		if datum.Children == nil {
			datum.Children = children
		} else {
			for _, child := range datum.Children {
				child.Data.setLeavesChildren(children)
			}
		}
	}
}

func (d OptionData[T]) filter(parentRoute Route, predicate func(Route) bool) (result OptionData[T]) {
	for _, datum := range d {
		route := append(parentRoute, datum.Option)
		if datum.Children != nil {
			var newChildren []*OptionNode
			for _, child := range datum.Children {
				newChild := child.filterRoutes(route, predicate)
				if newChild.Len() != 0 {
					newChildren = append(newChildren, newChild)
				}
			}
			if newChildren != nil {
				result = append(result, &OptionDatum[T]{Option: datum.Option, Children: newChildren})
			}
		} else if predicate(route) {
			result = append(result, datum)
		}
	}
	return
}

func (d OptionData[T]) len() int {
	return len(d)
}

type IOptionDatum interface {
	IOption() Option
	NonEmptyChildren() []*OptionNode
}

type OptionDatum[T Option] struct {
	Option   T
	Children []*OptionNode
}

func (d *OptionDatum[T]) IOption() Option {
	return d.Option
}

func (d *OptionDatum[T]) NonEmptyChildren() (result []*OptionNode) {
	if d == nil {
		return
	}
	for _, child := range d.Children {
		if child != nil && child.Data.len() > 0 {
			result = append(result, child)
		}
	}
	return
}

func (d *OptionDatum[T]) String() string {
	return d.Option.String()
}

type Option interface {
	String() string
}

type Route = []Option
