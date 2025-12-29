package kfl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service/period"
)

var (
	ErrMalformedFilterQuery   = errors.New("Malformed filter query") // This is only a just-in-case fallback.
	ErrCannotMixAndOr         = errors.New("Cannot mix && and || operators on the same level. Please use parenthesis () for grouping.")
	errUnbalancedBrackets     = errors.New("Missing") // Internal “base” class
	ErrUnbalancedOpenBracket  = fmt.Errorf("%w opening parenthesis. Please make sure that the number of opening and closing parentheses matches.", errUnbalancedBrackets)
	ErrUnbalancedCloseBracket = fmt.Errorf("%w closing parenthesis. Please make sure that the number of opening and closing parentheses matches.", errUnbalancedBrackets)
	errOperatorOperand        = errors.New("Missing") // Internal “base” class
	ErrOperatorExpected       = fmt.Errorf("%w operator. Please put logical operators ('&&' or '||') between the search operands.", errOperatorOperand)
	ErrOperandExpected        = fmt.Errorf("%w filter term. Please remove redundant logical operators.", errOperatorOperand)
	ErrIllegalTokenValue      = errors.New("Illegal value. Please make sure to use only valid operand values.")
)

func Parse(filterQuery string) (Predicate, ParseError) {
	p, pErr := func() (Predicate, ParseError) {
		tokens, pos, pErr := tokenise(filterQuery)
		if pErr != nil {
			return nil, pErr
		}
		tp := newTokenParser(
			append(tokens, tokenCloseBracket{}),
			append(pos, len(filterQuery)),
		)
		p, pErr := parseGroup(&tp)
		if pErr != nil {
			return nil, pErr
		}
		// Check whether there are tokens left, which would indicate
		// unbalanced brackets.
		if nextToken, _ := tp.next(); nextToken != nil {
			return nil, parseError{
				err:      ErrUnbalancedOpenBracket,
				position: 0,
				length:   len(filterQuery),
			}
		}
		return p, nil
	}()
	if pErr != nil {
		if pErr, ok := pErr.(parseError); ok {
			pErr.query = filterQuery
			return nil, pErr
		}
	}
	return p, nil
}

func parseGroup(tp *tokenParser) (Predicate, ParseError) {
	g := predicateGroup{}

	if pErr := tp.checkNextIsOperand(); pErr != nil {
		return nil, pErr
	}

	for {
		nextToken, position := tp.next()
		if nextToken == nil {
			return nil, parseError{
				err:      ErrUnbalancedCloseBracket,
				position: 0,
			}
		}

		switch tk := nextToken.(type) {

		case tokenOpenBracket:
			if pErr := tp.checkNextIsOperand(); pErr != nil {
				return nil, pErr
			}
			p, pErr := parseGroup(tp)
			if pErr != nil {
				return nil, pErr
			}
			g.append(p)

		case tokenCloseBracket:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			p, pErr := g.make()
			return p, pErr

		case tokenDate:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			date, err := klog.NewDateFromString(tk.date)
			if err != nil {
				return nil, parseError{
					err:      err,
					position: position,
					length:   len(tk.date),
				}
			}
			g.append(IsInDateRange{date, date})

		case tokenDateRange:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			dateBoundaries := []klog.Date{nil, nil}
			for i, v := range tk.bounds {
				if v == "" {
					continue
				}
				date, err := klog.NewDateFromString(v)
				if err != nil {
					return nil, parseError{
						err:      err,
						position: position,
						length:   len(strings.Join(tk.bounds, "...")),
					}
				}
				dateBoundaries[i] = date
			}
			g.append(IsInDateRange{dateBoundaries[0], dateBoundaries[1]})

		case tokenPeriod:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			prd, err := period.NewPeriodFromPatternString(tk.period)
			if err != nil {
				return nil, parseError{
					err:      err,
					position: position,
					length:   len(tk.period),
				}
			}
			g.append(IsInDateRange{prd.Since(), prd.Until()})

		case tokenAnd, tokenOr:
			if pErr := tp.checkNextIsOperand(); pErr != nil {
				return nil, pErr
			}
			pErr := g.setOperator(tk, position)
			if pErr != nil {
				return nil, pErr
			}

		case tokenNot:
			if pErr := tp.checkNextIsOperand(); pErr != nil {
				return nil, pErr
			}
			g.negateNextOperand()

		case tokenEntryType:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			et, err := NewEntryTypeFromString(tk.entryType)
			if err != nil {
				return nil, parseError{
					err:      err,
					position: position,
					length:   len("type:") + len(tk.entryType),
				}
			}
			g.append(IsEntryType{et})

		case tokenTag:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			tag, err := klog.NewTagFromString(tk.tag)
			if err != nil {
				return nil, parseError{
					err:      err,
					position: position,
					length:   len(tk.tag),
				}
			}
			g.append(HasTag{tag})

		default:
			panic("Unrecognized token")
		}
	}
}
