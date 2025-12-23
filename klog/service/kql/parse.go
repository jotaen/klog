package kql

import (
	"errors"
	"fmt"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service/period"
)

var (
	ErrMalformedQuery     = errors.New("Malformed query") // This is only a just-in-case fallback.
	ErrCannotMixAndOr     = errors.New("Cannot mix && and || operators on the same level. Please use parenthesis () for grouping.")
	ErrUnbalancedBrackets = errors.New("Unbalanced parenthesis. Please make sure that the number of opening and closing parentheses matches.")
	errOperatorOperand    = errors.New("Missing expected") // Internal “base” class
	ErrOperatorExpected   = fmt.Errorf("%w operator. Please put logical operators ('&&' or '||') between the search operands.", errOperatorOperand)
	ErrOperandExpected    = fmt.Errorf("%w operand. Please remove redundant logical operators.", errOperatorOperand)
	ErrIllegalTokenValue  = errors.New("Illegal value. Please make sure to use only valid operand values.")
)

func Parse(query string) (Predicate, error) {
	tokens, err := tokenise(query)
	if err != nil {
		return nil, err
	}
	tp := newTokenParser(append(tokens, tokenCloseBracket{}))
	p, err := parseGroup(&tp)
	if err != nil {
		return nil, err
	}
	// Check whether there are tokens left, which would indicate
	// unbalanced brackets.
	if tp.next() != nil {
		return nil, ErrUnbalancedBrackets
	}
	return p, nil
}

func parseGroup(tp *tokenParser) (Predicate, error) {
	g := predicateGroup{}

	for {
		nextToken := tp.next()
		if nextToken == nil {
			break
		}

		switch tk := nextToken.(type) {

		case tokenOpenBracket:
			if err := tp.checkNextIsOperand(); err != nil {
				return nil, err
			}
			p, err := parseGroup(tp)
			if err != nil {
				return nil, err
			}
			g.append(p)

		case tokenCloseBracket:
			if err := tp.checkNextIsOperatorOrEnd(); err != nil {
				return nil, err
			}
			p, err := g.make()
			return p, err

		case tokenDate:
			if err := tp.checkNextIsOperatorOrEnd(); err != nil {
				return nil, err
			}
			date, err := klog.NewDateFromString(tk.date)
			if err != nil {
				return nil, err
			}
			g.append(IsInDateRange{date, date})

		case tokenDateRange:
			if err := tp.checkNextIsOperatorOrEnd(); err != nil {
				return nil, err
			}
			dateBoundaries := []klog.Date{nil, nil}
			for i, v := range tk.bounds {
				if v == "" {
					continue
				}
				date, err := klog.NewDateFromString(v)
				if err != nil {
					return nil, err
				}
				dateBoundaries[i] = date
			}
			g.append(IsInDateRange{dateBoundaries[0], dateBoundaries[1]})

		case tokenPeriod:
			if err := tp.checkNextIsOperatorOrEnd(); err != nil {
				return nil, err
			}
			prd, err := period.NewPeriodFromPatternString(tk.period)
			if err != nil {
				return nil, err
			}
			g.append(IsInDateRange{prd.Since(), prd.Until()})

		case tokenAnd, tokenOr:
			if err := tp.checkNextIsOperand(); err != nil {
				return nil, err
			}
			err := g.setOperator(tk)
			if err != nil {
				return nil, err
			}

		case tokenNot:
			if err := tp.checkNextIsOperand(); err != nil {
				return nil, err
			}
			g.negateNextOperand()

		case tokenTag:
			if err := tp.checkNextIsOperatorOrEnd(); err != nil {
				return nil, err
			}
			tag, err := klog.NewTagFromString(tk.tag)
			if err != nil {
				return nil, err
			}
			g.append(HasTag{tag})

		default:
			panic("Unrecognized token")
		}
	}

	return nil, ErrUnbalancedBrackets
}
