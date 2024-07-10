package slogfluentd

import (
	"context"

	"log/slog"

	"github.com/fluent/fluent-logger-golang/fluent"
	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// connection to Fluentd
	Client *fluent.Fluent
	Tag    string

	// optional: customize json payload builder
	Converter Converter
	// optional: fetch attributes from context
	AttrFromContext []func(ctx context.Context) []slog.Attr

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func (o Option) NewFluentdHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Client == nil {
		panic("missing Fuentd client")
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	if o.AttrFromContext == nil {
		o.AttrFromContext = []func(ctx context.Context) []slog.Attr{}
	}

	return &FluentdHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*FluentdHandler)(nil)

type FluentdHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *FluentdHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *FluentdHandler) Handle(ctx context.Context, record slog.Record) error {
	tag := h.getTag(&record)
	fromContext := slogcommon.ContextExtractor(ctx, h.option.AttrFromContext)
	message := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, append(h.attrs, fromContext...), h.groups, &record, tag)

	return h.option.Client.PostWithTime(tag, record.Time, message)
}

func (h *FluentdHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &FluentdHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *FluentdHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &FluentdHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

func (h *FluentdHandler) getTag(record *slog.Record) string {
	tag := h.option.Tag

	for i := range h.attrs {
		if h.attrs[i].Key == "tag" && h.attrs[i].Value.Kind() == slog.KindString {
			tag = h.attrs[i].Value.String()
			break
		}
	}

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "tag" && attr.Value.Kind() == slog.KindString {
			tag = attr.Value.String()
			return false
		}
		return true
	})

	return tag
}
