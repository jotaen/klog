package klog

import (
	"errors"
	"github.com/alecthomas/kong"
	klog "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service"
	"github.com/jotaen/klog/src/service/period"
	"reflect"
	"strings"
)

func dateDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("date", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid date")
		}
		d, err := klog.NewDateFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid date")
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}

func timeDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("time", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid time")
		}
		t, err := klog.NewTimeFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid time")
		}
		target.Set(reflect.ValueOf(t))
		return nil
	}
}

func shouldTotalDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("should", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid should-total duration")
		}
		valueAsDuration := strings.TrimSuffix(value, "!")
		duration, err := klog.NewDurationFromString(valueAsDuration)
		if err != nil {
			return errors.New("`" + value + "` is not a valid should total")
		}
		shouldTotal := klog.NewShouldTotal(0, duration.InMinutes())
		target.Set(reflect.ValueOf(shouldTotal))
		return nil
	}
}

func periodDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("period", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid period")
		}
		p, err := period.NewPeriodFromPatternString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid period")
		}
		target.Set(reflect.ValueOf(p))
		return nil
	}
}

func roundingDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("rounder", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid rounding value")
		}
		r, err := service.NewRoundingFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid rounding value")
		}
		target.Set(reflect.ValueOf(r))
		return nil
	}
}

func tagDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("tag", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid tag")
		}
		t, err := klog.NewTagFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid tag")
		}
		target.Set(reflect.ValueOf(t))
		return nil
	}
}

func recordSummaryDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("recordSummary", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid record summary")
		}
		// Normalize potential double-escaping (from CLI)
		value = strings.ReplaceAll(value, "\\n", "\n")
		summary, sErr := klog.NewRecordSummary(strings.Split(value, "\n")...)
		if sErr != nil {
			return errors.New("A record summary cannot contain blank lines, and none of its lines can start with whitespace characters")
		}
		target.Set(reflect.ValueOf(summary))
		return nil
	}
}

func entrySummaryDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("entrySummary", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid record summary")
		}
		// Normalize potential double-escaping (from CLI)
		value = strings.ReplaceAll(value, "\\n", "\n")
		// When passing entries like `-45m` the leading dash must be escaped
		// by the user, otherwise it would be treated like a flag. Therefore, we
		// have to remove the potential escaping backslash. Examples:
		// - `\\-45m` (unquoted, with double-escape)
		// - `'\-45m'` (quoted, with single-escape)
		if strings.HasPrefix(value, "\\-") {
			value = strings.TrimPrefix(value, "\\")
		} else if strings.HasPrefix(value, "\\\\-") {
			value = strings.TrimPrefix(value, "\\\\")
		}
		summary, sErr := klog.NewEntrySummary(strings.Split(value, "\n")...)
		if sErr != nil {
			return errors.New("An entry summary cannot contain blank lines")
		}
		target.Set(reflect.ValueOf(summary))
		return nil
	}
}
