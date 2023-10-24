package mess

import "fmt"

type OptionGroups = map[string]OptionGroup

func GroupOptions(optionSets [][]Option, i int) OptionGroups {
	result := make(OptionGroups)
	for _, options := range optionSets {
		option := options[i]
		group, found := result[option.Message()]
		if !found {
			f := &factory{}
			option.Accept(f)
			group = f.group
			result[option.Message()] = group
		}
		err := group.addOption(option)
		if err != nil {
			panic(err)
		}
	}
	return result
}

type factory struct{ group OptionGroup }

func (f *factory) VisitPieceTypeOption(_ *PieceTypeOption) { f.group = &PieceTypeOptionGroup{} }
func (f *factory) VisitSquareOption(_ *SquareOption)       { f.group = &SquareOptionGroup{} }
func (f *factory) VisitMoveOption(_ *MoveOption)           { f.group = &MoveOptionGroup{} }
func (f *factory) VisitUnitOption(_ *UnitOption)           { f.group = &UnitOptionGroup{} }

// Option groups

type PieceTypeOptionGroup struct {
	OptionGroupBase[*PieceTypeOption]
}

func (ptog *PieceTypeOptionGroup) Accept(visitor OptionGroupVisitor) {
	visitor.VisitPieceTypeOptions(ptog.Options)
}

type SquareOptionGroup struct{ OptionGroupBase[*SquareOption] }

func (sog *SquareOptionGroup) Accept(visitor OptionGroupVisitor) {
	visitor.VisitSquareOptions(sog.Options)
}

type MoveOptionGroup struct{ OptionGroupBase[*MoveOption] }

func (mog *MoveOptionGroup) Accept(visitor OptionGroupVisitor) {
	visitor.VisitMoveOptions(mog.Options)
}

type UnitOptionGroup struct{ OptionGroupBase[*UnitOption] }

func (uog *UnitOptionGroup) Accept(visitor OptionGroupVisitor) {
	visitor.VisitUnitOptions(uog.Options)
}

type OptionGroup interface {
	addOption(Option) error
	Accept(visitor OptionGroupVisitor)
}

type OptionGroupBase[T Option] struct {
	Options []T
}

func (ogb *OptionGroupBase[T]) addOption(option Option) error {
	opt, ok := option.(T)
	if !ok {
		return fmt.Errorf("tried to add option of type %T to %T", option, ogb)
	}
	ogb.Options = append(ogb.Options, opt)
	return nil
}

// Option group visitor

type OptionGroupVisitor interface {
	VisitPieceTypeOptions(options []*PieceTypeOption)
	VisitSquareOptions(options []*SquareOption)
	VisitMoveOptions(options []*MoveOption)
	VisitUnitOptions(options []*UnitOption)
}
