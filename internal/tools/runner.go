package tools

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var ErrOutputLimit = errors.New("media tool output exceeded limit")

type CommandSpec struct {
	Name           string
	Args           []string
	Timeout        time.Duration
	MaxOutputBytes int64
	MaxStderrBytes int64
}

type ProgressCommandSpec struct {
	CommandSpec
	Progress func(string)
}

func LookPath(name string) (string, error) {
	if err := validateToolName(name); err != nil {
		return "", err
	}
	return exec.LookPath(name)
}

func RunOutput(ctx context.Context, spec CommandSpec) ([]byte, error) {
	if err := validateToolName(spec.Name); err != nil {
		return nil, err
	}
	if spec.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, spec.Timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Env = mediaToolEnv()
	stdout := newLimitedBuffer(spec.MaxOutputBytes)
	stderr := newLimitedBuffer(spec.MaxStderrBytes)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return stdout.Bytes(), toolRunError(ctx, err, stderr)
	}
	if stdout.Limited() || stderr.Limited() {
		return stdout.Bytes(), ErrOutputLimit
	}
	return stdout.Bytes(), nil
}

func RunOutputProgress(ctx context.Context, spec ProgressCommandSpec) ([]byte, error) {
	if spec.Progress == nil {
		return RunOutput(ctx, spec.CommandSpec)
	}
	if err := validateToolName(spec.Name); err != nil {
		return nil, err
	}
	if spec.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, spec.Timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Env = mediaToolEnv()
	stderr := newLimitedBuffer(spec.MaxStderrBytes)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	output, readErr := readProgressOutput(stdout, spec.MaxOutputBytes, spec.Progress)
	waitErr := cmd.Wait()
	if readErr != nil {
		return output, readErr
	}
	if waitErr != nil {
		return output, toolRunError(ctx, waitErr, stderr)
	}
	if stderr.Limited() {
		return output, ErrOutputLimit
	}
	return output, nil
}

func readProgressOutput(stdout io.Reader, maxOutputBytes int64, progress func(string)) ([]byte, error) {
	output := newLimitedBuffer(maxOutputBytes)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		_, _ = output.Write(append([]byte(line), '\n'))
		progress(line)
	}
	if err := scanner.Err(); err != nil {
		return output.Bytes(), err
	}
	if output.Limited() {
		return output.Bytes(), ErrOutputLimit
	}
	return output.Bytes(), nil
}

func RunStream(ctx context.Context, name string, args []string, stdout io.Writer, maxStderrBytes int64) error {
	if err := validateToolName(name); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = mediaToolEnv()
	stderr := newLimitedBuffer(maxStderrBytes)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return toolRunError(ctx, err, stderr)
	}
	if stderr.Limited() {
		return ErrOutputLimit
	}
	return nil
}

func SafePathArg(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("media path is empty")
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("media path must be absolute: %s", path)
	}
	if strings.HasPrefix(filepath.Base(path), "-") {
		return fmt.Errorf("media path must not look like an option: %s", path)
	}
	return nil
}

func validateToolName(name string) error {
	if name == "" {
		return errors.New("media tool name is empty")
	}
	if strings.HasPrefix(name, "-") || strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid media tool name: %s", name)
	}
	return nil
}

func mediaToolEnv() []string {
	env := []string{"LANG=C"}
	if path := os.Getenv("PATH"); path != "" {
		env = append(env, "PATH="+path)
	}
	return env
}

func toolRunError(ctx context.Context, err error, stderr *limitedBuffer) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if stderr.Limited() {
		return ErrOutputLimit
	}
	message := strings.TrimSpace(stderr.String())
	if message == "" {
		return err
	}
	return errors.New(message)
}

type limitedBuffer struct {
	max     int64
	limited bool
	buf     bytes.Buffer
}

func newLimitedBuffer(max int64) *limitedBuffer {
	return &limitedBuffer{max: max}
}

func (b *limitedBuffer) Write(payload []byte) (int, error) {
	if b.max <= 0 {
		_, _ = b.buf.Write(payload)
		return len(payload), nil
	}
	remaining := b.max - int64(b.buf.Len())
	if remaining <= 0 {
		b.limited = true
		return len(payload), nil
	}
	if int64(len(payload)) > remaining {
		_, _ = b.buf.Write(payload[:int(remaining)])
		b.limited = true
		return len(payload), nil
	}
	_, _ = b.buf.Write(payload)
	return len(payload), nil
}

func (b *limitedBuffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *limitedBuffer) String() string {
	return b.buf.String()
}

func (b *limitedBuffer) Limited() bool {
	return b.limited
}
