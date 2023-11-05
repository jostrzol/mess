package mess

import (
	"fmt"
)

type OptionTree interface {
	Accept(visitor OptionTreeVisitor)
	filterOptionSets([]Option, func([]Option) bool) OptionTree
	len() int
}

func MakeOptionTree(generators []ChoiceGenerator) OptionTree {
	return makeOptionTree([]Option{}, generators)
}

func makeOptionTree(options []Option, generators []ChoiceGenerator) OptionTree {
	if len(generators) == 0 {
		return nil
	}

	choice := generators[0](options)
	if choice == nil {
		// empty option tree - action won't be performed
		return &MessageNode{}
	}

	optionNodes := make(map[string]OptionNode)
	for _, option := range choice.GenerateOptions() {
		node, found := optionNodes[option.Message()]
		if !found {
			f := &factory{}
			option.Accept(f)
			node = f.node
			optionNodes[option.Message()] = node
		}
		msgNode := makeOptionTree(append(options, option), generators[1:])
		if msgNode == nil || msgNode.len() != 0 {
			err := node.addOption(option, msgNode)
			if err != nil {
				panic(err)
			}
		}
	}

	optionTreeNodes := make(map[string]OptionTree, len(optionNodes))
	for key, node := range optionNodes {
		optionTreeNodes[key] = node
	}

	return &MessageNode{optionNodes: optionTreeNodes}
}

type MessageNode struct {
	optionNodes map[string]OptionTree
}

func (n *MessageNode) Accept(visitor OptionTreeVisitor) {
	visitor.VisitMessageNode(n.optionNodes)
}

func (n *MessageNode) len() int {
	return len(n.optionNodes)
}

func (n *MessageNode) filterOptionSets(options []Option, predicate func([]Option) bool) OptionTree {
	if n.optionNodes == nil && predicate(nil) {
		return &MessageNode{}
	}
	newOptionNodes := make(map[string]OptionTree)
	for msg, optionNode := range n.optionNodes {
		optionNode.filterOptionSets(options, predicate)
		if optionNode.len() != 0 {
			newOptionNodes[msg] = optionNode
		}
	}
	return &MessageNode{optionNodes: newOptionNodes}
}

type OptionNode interface {
	OptionTree
	addOption(Option, OptionTree) error
	len() int
}

type factory struct{ node OptionNode }

func (f *factory) VisitPieceTypeOption(_ *PieceTypeOption) { f.node = &PieceTypeOptionNode{} }
func (f *factory) VisitSquareOption(_ *SquareOption)       { f.node = &SquareOptionNode{} }
func (f *factory) VisitMoveOption(_ *MoveOption)           { f.node = &MoveOptionNode{} }
func (f *factory) VisitUnitOption(_ *UnitOption)           { f.node = &UnitOptionNode{} }

type PieceTypeOptionNode struct {
	OptionNodeBase[*PieceTypeOption]
}

func (n *PieceTypeOptionNode) Accept(visitor OptionTreeVisitor) {
	visitor.VisitPieceTypeNode(n.Options)
}

func (n *PieceTypeOptionNode) filterOptionSets(options []Option, predicate func([]Option) bool) OptionTree {
	return &PieceTypeOptionNode{OptionNodeBase: n.OptionNodeBase.filterOptionSets(options, predicate)}
}

type SquareOptionNode struct{ OptionNodeBase[*SquareOption] }

func (n *SquareOptionNode) Accept(visitor OptionTreeVisitor) {
	visitor.VisitSquareNode(n.Options)
}

func (n *SquareOptionNode) filterOptionSets(options []Option, predicate func([]Option) bool) OptionTree {
	return &SquareOptionNode{OptionNodeBase: n.OptionNodeBase.filterOptionSets(options, predicate)}
}

type MoveOptionNode struct{ OptionNodeBase[*MoveOption] }

func (n *MoveOptionNode) Accept(visitor OptionTreeVisitor) {
	visitor.VisitMoveNode(n.Options)
}

func (n *MoveOptionNode) filterOptionSets(options []Option, predicate func([]Option) bool) OptionTree {
	return &MoveOptionNode{OptionNodeBase: n.OptionNodeBase.filterOptionSets(options, predicate)}
}

type UnitOptionNode struct{ OptionNodeBase[*UnitOption] }

func (n *UnitOptionNode) Accept(visitor OptionTreeVisitor) {
	visitor.VisitUnitNode(n.Options)
}

func (n *UnitOptionNode) filterOptionSets(options []Option, predicate func([]Option) bool) OptionTree {
	return &UnitOptionNode{OptionNodeBase: n.OptionNodeBase.filterOptionSets(options, predicate)}
}

type OptionNodeBase[T interface {
	Option
	comparable
}] struct {
	Options map[T]OptionTree
}

func (n *OptionNodeBase[T]) addOption(option Option, node OptionTree) error {
	opt, ok := option.(T)
	if !ok {
		return fmt.Errorf("tried to add option of type %T to %T", option, n)
	}
	if n.Options == nil {
		n.Options = make(map[T]OptionTree, 1)
	}
	n.Options[opt] = node
	return nil
}

func (n *OptionNodeBase[T]) len() int {
	return len(n.Options)
}

func (n *OptionNodeBase[T]) filterOptionSets(options []Option, predicate func([]Option) bool) OptionNodeBase[T] {
	newOptions := make(map[T]OptionTree)
	for option, msgNode := range n.Options {
		optionSet := append([]Option{option}, options...)
		if msgNode == nil {
			// last node
			if !predicate(optionSet) {
				newOptions[option] = msgNode
			}
		} else {
			// not last node
			newMsgNode := msgNode.filterOptionSets(optionSet, predicate)
			if newMsgNode.len() != 0 {
				newOptions[option] = newMsgNode
			}
		}
	}
	return OptionNodeBase[T]{Options: newOptions}
}

func FilterOptionTree(tree OptionTree, predicate func([]Option) bool) OptionTree {
	switch {
	case tree != nil:
		return tree.filterOptionSets([]Option{}, predicate)
	case predicate([]Option{}):
		return nil
	default:
		// empty message node with non-nil optionNodes -- move will be filtered
		return &MessageNode{optionNodes: map[string]OptionTree{}}
	}
}

func AllOptions(tree OptionTree) (result [][]Option) {
	if tree == nil {
		return [][]Option{{}}
	}
	tree.filterOptionSets([]Option{}, func(options []Option) bool {
		result = append(result, options)
		return false
	})
	return
}

type OptionTreeVisitor interface {
	VisitPieceTypeNode(options map[*PieceTypeOption]OptionTree)
	VisitSquareNode(options map[*SquareOption]OptionTree)
	VisitMoveNode(options map[*MoveOption]OptionTree)
	VisitUnitNode(options map[*UnitOption]OptionTree)
	VisitMessageNode(options map[string]OptionTree)
}
