package mess

type Choice struct {
	Message     string
	NextChoices []*Choice
	Generator   ChoiceGenerator
}

type ChoiceGenerator interface {
	GenerateOptions() IOptionNodeData
}

func (c *Choice) GenerateOptions() *OptionNode {
	if c == nil {
		return nil
	}

	optionData := c.Generator.GenerateOptions()
	var children []*OptionNode
	for _, choice := range c.NextChoices {
		children = append(children, choice.GenerateOptions())
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
	Data    IOptionNodeData
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

func (n *OptionNode) Accept(visitor OptionNodeDataVisitor) {
	n.Data.accept(n.Message, visitor)
}

func (n *OptionNode) AllRoutes() <-chan Route {
	result := make(chan Route)
	go func() {
		n.FilterRoutes(func(route Route) bool {
			result <- route
			return false
		})
		close(result)
	}()
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

type IOptionNodeData interface {
	accept(message string, visitor OptionNodeDataVisitor)
	setLeavesChildren(children []*OptionNode)
	filter(parentRoute Route, predicate func(Route) bool) IOptionNodeData
	len() int
}

type OptionNodeData[T Option] []*OptionNodeDatum[T]

func (d OptionNodeData[T]) setLeavesChildren(children []*OptionNode) {
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

func (d OptionNodeData[T]) filter(parentRoute Route, predicate func(Route) bool) (result OptionNodeData[T]) {
	for _, datum := range d {
		route := append(parentRoute, datum.Option)
		if datum.Children != nil {
			var newChildren []*OptionNode
			for _, child := range datum.Children {
				newChild := child.filterRoutes(route, predicate)
				newChildren = append(newChildren, newChild)
			}
			if newChildren != nil {
				result = append(result, &OptionNodeDatum[T]{Option: datum.Option, Children: newChildren})
			}
		} else if predicate(route) {
			result = append(result, datum)
		}
	}
	return
}

func (d OptionNodeData[T]) len() int {
	return len(d)
}

type IOptionNodeDatum interface {
	IOption() Option
	NonEmptyChildren() []*OptionNode
}

type OptionNodeDatum[T Option] struct {
	Option   T
	Children []*OptionNode
}

func (d *OptionNodeDatum[T]) IOption() Option {
	return d.Option
}

func (d *OptionNodeDatum[T]) NonEmptyChildren() (result []*OptionNode) {
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

func (d *OptionNodeDatum[T]) String() string {
	return d.Option.String()
}

type Option interface {
	String() string
}

type Route = []Option
