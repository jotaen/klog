package main

import (
	"errors"
	"github.com/alecthomas/kong"
	"klog"
	"klog/app/cli/lib"
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

func durationDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("duration", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid duration")
		}
		value = strings.TrimSuffix(value, "!")
		d, err := klog.NewDurationFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid duration")
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}

func periodDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("period", &value); err != nil {
			return err
		}
		p, err := lib.NewPeriodFromString(value)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(p))
		return nil
	}
}
