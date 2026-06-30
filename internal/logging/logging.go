package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"
)

const maxBufferedEntries = 300

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type Entry struct {
	ID         string         `json:"id"`
	Time       time.Time      `json:"time"`
	Level      Level          `json:"level"`
	Message    string         `json:"message"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type Manager struct {
	level       slog.LevelVar
	fileSink    *fileSink
	mu          sync.Mutex
	nextID      uint64
	entries     []Entry
	subscribers map[chan Entry]struct{}
}

var Default = NewManager()

func NewManager() *Manager {
	manager := &Manager{subscribers: make(map[chan Entry]struct{}), fileSink: newFileSink()}
	manager.level.Set(slog.LevelInfo)
	return manager
}

func ConfigureDefault() {
	slog.SetDefault(slog.New(Default.Handler(os.Stdout)))
}

func (m *Manager) Handler(writer io.Writer) slog.Handler {
	return &streamHandler{
		base: slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: &m.level}),
		m:    m,
	}
}

func (m *Manager) Level() Level {
	switch m.level.Level() {
	case slog.LevelDebug:
		return LevelDebug
	case slog.LevelWarn:
		return LevelWarn
	case slog.LevelError:
		return LevelError
	default:
		return LevelInfo
	}
}

func (m *Manager) SetLevel(level Level) error {
	switch level {
	case LevelDebug:
		m.level.Set(slog.LevelDebug)
	case LevelInfo:
		m.level.Set(slog.LevelInfo)
	case LevelWarn:
		m.level.Set(slog.LevelWarn)
	case LevelError:
		m.level.Set(slog.LevelError)
	default:
		return errors.New("unsupported log level")
	}
	slog.Info("log level changed", "level", level)
	return nil
}

func (m *Manager) Subscribe() (<-chan Entry, func()) {
	ch := make(chan Entry, maxBufferedEntries)

	m.mu.Lock()
	for _, entry := range m.entries {
		ch <- entry
	}
	m.subscribers[ch] = struct{}{}
	m.mu.Unlock()

	return ch, func() {
		m.mu.Lock()
		delete(m.subscribers, ch)
		close(ch)
		m.mu.Unlock()
	}
}

func (m *Manager) publish(record slog.Record) {
	m.mu.Lock()
	m.nextID++
	entry := Entry{
		ID:         strconv.FormatUint(m.nextID, 10),
		Time:       record.Time.UTC(),
		Level:      levelFromSlog(record.Level),
		Message:    record.Message,
		Attributes: attrsFromRecord(record),
	}
	m.entries = append(m.entries, entry)
	if len(m.entries) > maxBufferedEntries {
		m.entries = m.entries[len(m.entries)-maxBufferedEntries:]
	}
	for ch := range m.subscribers {
		select {
		case ch <- entry:
		default:
		}
	}
	m.mu.Unlock()
}

type streamHandler struct {
	base slog.Handler
	m    *Manager
}

func (h *streamHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

func (h *streamHandler) Handle(ctx context.Context, record slog.Record) error {
	err := h.base.Handle(ctx, record)
	h.m.writeFileRecord(record)
	h.m.publish(record)
	return err
}

func (h *streamHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &streamHandler{base: h.base.WithAttrs(attrs), m: h.m}
}

func (h *streamHandler) WithGroup(name string) slog.Handler {
	return &streamHandler{base: h.base.WithGroup(name), m: h.m}
}

func levelFromSlog(level slog.Level) Level {
	switch {
	case level <= slog.LevelDebug:
		return LevelDebug
	case level >= slog.LevelError:
		return LevelError
	case level >= slog.LevelWarn:
		return LevelWarn
	default:
		return LevelInfo
	}
}

func attrsFromRecord(record slog.Record) map[string]any {
	attrs := map[string]any{}
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = slogValue(attr.Value)
		return true
	})
	if len(attrs) == 0 {
		return nil
	}
	return attrs
}

func slogValue(value slog.Value) any {
	switch value.Kind() {
	case slog.KindString:
		return value.String()
	case slog.KindBool:
		return value.Bool()
	case slog.KindInt64:
		return value.Int64()
	case slog.KindUint64:
		return value.Uint64()
	case slog.KindFloat64:
		return value.Float64()
	case slog.KindDuration:
		return value.Duration().String()
	case slog.KindTime:
		return value.Time().UTC()
	case slog.KindGroup:
		group := map[string]any{}
		for _, attr := range value.Group() {
			group[attr.Key] = slogValue(attr.Value)
		}
		return group
	case slog.KindLogValuer:
		return slogValue(value.Resolve())
	default:
		return fmt.Sprint(value.Any())
	}
}
