package filter

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
	ErrOperatorExpected       = fmt.Errorf("%w operator. Please put a logical operator ('&&' or '||') before this search operand.", errOperatorOperand)
	ErrOperandExpected        = fmt.Errorf("%w filter term. Please remove redundant logical operators.", errOperatorOperand)
	ErrIllegalTokenValue      = errors.New("Illegal value. Please make sure to use only valid operand values.")
)

func Parse(filterQuery string) (Predicate, ParseError) {
	p, pErr := func() (Predicate, ParseError) {
		tokens, pErr := tokenise(filterQuery)
		if pErr != nil {
			return nil, pErr
		}
		tp := newTokenParser(
			append(tokens, token{tokenCloseBracket, ")", len(filterQuery) - 1}),
		)
		p, pErr := parseGroup(&tp, filterQuery)
		if pErr != nil {
			return nil, pErr
		}
		// Check whether there are tokens left, which would indicate
		// unbalanced brackets.
		if tp.next() != (token{}) {
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

func parseGroup(tp *tokenParser, filterQuery string) (Predicate, ParseError) {
	g := newPredicateGroup()

	if pErr := tp.checkNextIsOperand(); pErr != nil {
		return nil, pErr
	}

	for {
		tk := tp.next()
		if tk == (token{}) {
			return nil, parseError{
				err:      ErrUnbalancedCloseBracket,
				position: 0,
				length:   len(filterQuery),
			}
		}

		switch tk.kind {

		case tokenOpenBracket:
			if pErr := tp.checkNextIsOperand(); pErr != nil {
				return nil, pErr
			}
			p, pErr := parseGroup(tp, filterQuery)
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
			date, err := klog.NewDateFromString(tk.value)
			if err != nil {
				return nil, parseError{
					err:      ErrIllegalTokenValue,
					position: tk.position,
					length:   len(tk.value),
				}
			}
			g.append(IsInDateRange{date, date})

		case tokenDateRange:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			dateBoundaries := []klog.Date{nil, nil}
			bounds := strings.Split(tk.value, "...")
			for i, v := range bounds {
				if v == "" {
					continue
				}
				// Try whether bound is period:
				prd, err := period.NewPeriodFromPatternString(v)
				if err == nil {
					if i == 0 {
						dateBoundaries[i] = prd.Since()
					} else {
						dateBoundaries[i] = prd.Until()
					}
					continue
				}
				// Try whether bound is date:
				date, err := klog.NewDateFromString(v)
				if err == nil {
					dateBoundaries[i] = date
					continue
				}
				// Otherwise, yield error:
				return nil, parseError{
					err:      ErrIllegalTokenValue,
					position: tk.position,
					length:   len(tk.value),
				}
			}
			if dateBoundaries[0] != nil && dateBoundaries[1] != nil {
				if !dateBoundaries[1].IsAfterOrEqual(dateBoundaries[0]) {
					return nil, parseError{
						err:      ErrIllegalTokenValue,
						position: tk.position,
						length:   len(tk.value),
					}
				}
			}
			g.append(IsInDateRange{dateBoundaries[0], dateBoundaries[1]})

		case tokenPeriod:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			prd, err := period.NewPeriodFromPatternString(tk.value)
			if err != nil {
				return nil, parseError{
					err:      ErrIllegalTokenValue,
					position: tk.position,
					length:   len(tk.value),
				}
			}
			g.append(IsInDateRange{prd.Since(), prd.Until()})

		case tokenAnd, tokenOr:
			if pErr := tp.checkNextIsOperand(); pErr != nil {
				return nil, pErr
			}
			pErr := g.setOperator(tk, tk.position)
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
			et, err := NewEntryTypeFromString(strings.TrimLeft(tk.value, "type:"))
			if err != nil {
				return nil, parseError{
					err:      ErrIllegalTokenValue,
					position: tk.position,
					length:   len(tk.value),
				}
			}
			g.append(IsEntryType{et})

		case tokenTag:
			if pErr := tp.checkNextIsOperatorOrEnd(); pErr != nil {
				return nil, pErr
			}
			tag, err := klog.NewTagFromString(tk.value)
			if err != nil {
				return nil, parseError{
					err:      ErrIllegalTokenValue,
					position: tk.position,
					length:   len(tk.value),
				}
			}
			g.append(HasTag{tag})

		default:
			// This should never happen.
			panic("Unrecognized token")
		}
	}
}
