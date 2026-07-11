package soap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const envelopeNS = "http://schemas.xmlsoap.org/soap/envelope/"

type HandlerFunc func(context.Context, map[string]string) (map[string]string, error)

type Action struct {
	Service string
	Name    string
	Args    map[string]string
}

type Service struct {
	Type    string
	Actions map[string]HandlerFunc
}

type Dispatcher struct {
	services map[string]Service
}

type Error struct {
	Code        int
	Description string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Description)
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{services: map[string]Service{}}
}

func (d *Dispatcher) Register(path string, serviceType string, actions map[string]HandlerFunc) {
	d.services[path] = Service{Type: serviceType, Actions: actions}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, ok := d.services[r.URL.Path]
	if !ok {
		WriteFault(w, Error{Code: 401, Description: "Invalid Action"})
		return
	}
	action, err := ParseRequest(r, service.Type)
	if err != nil {
		WriteFault(w, asSOAPError(err, 402, "Invalid Args"))
		return
	}
	slog.Debug("dlna soap action", "path", r.URL.Path, "service", action.Service, "action", action.Name)
	handler, ok := service.Actions[action.Name]
	if !ok {
		WriteFault(w, Error{Code: 401, Description: "Invalid Action"})
		return
	}
	values, err := handler(ContextWithRequest(r.Context(), r), action.Args)
	if err != nil {
		WriteFault(w, asSOAPError(err, 501, "Action Failed"))
		return
	}
	WriteResponse(w, service.Type, action.Name, values)
}

type requestContextKey struct{}

func ContextWithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey{}, r)
}

func RequestFromContext(ctx context.Context) (*http.Request, bool) {
	r, ok := ctx.Value(requestContextKey{}).(*http.Request)
	return r, ok
}

func ParseSOAPAction(value string) (string, string, bool) {
	value = strings.Trim(strings.TrimSpace(value), `"`)
	service, action, ok := strings.Cut(value, "#")
	if !ok || strings.TrimSpace(service) == "" || strings.TrimSpace(action) == "" {
		return "", "", false
	}
	return strings.TrimSpace(service), strings.TrimSpace(action), true
}

func asSOAPError(err error, code int, description string) Error {
	if soapErr, ok := err.(Error); ok {
		return soapErr
	}
	return Error{Code: code, Description: description}
}

func InvalidArgs(message string) Error {
	if strings.TrimSpace(message) == "" {
		message = "Invalid Args"
	}
	return Error{Code: 402, Description: message}
}

func RequiredArg(args map[string]string, name string) (string, error) {
	value := strings.TrimSpace(args[name])
	if value == "" {
		return "", InvalidArgs("Missing argument: " + name)
	}
	return value, nil
}
