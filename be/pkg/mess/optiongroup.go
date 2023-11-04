package mess

import (
	"fmt"
)

type OptionGroups = map[string]OptionGroup

func FilterOptions(
	optionSets [][]Option,
	predicate func(map[string]OptionGroup) (Option, error),
) ([][]Option, error) {
	if len(optionSets) == 0 {
		return [][]Option{}, nil
	} else if len(optionSets) == 1 && len(optionSets[0]) == 0 {
		return [][]Option{optionSets[0]}, nil
	}

	for i := 0; i < len(optionSets[0]); i++ {
		groups := groupOptions(optionSets, i)

		option, err := predicate(groups)
		if err != nil {
			return nil, err
		}

		newOptionSets := make([][]Option, 0)
		for _, options := range optionSets {
			if options[i] == option {
				newOptionSets = append(newOptionSets, options)
			}
		}
		optionSets = newOptionSets
	}

	return optionSets, nil
}

func groupOptions(optionSets [][]Option, i int) OptionGroups {
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
