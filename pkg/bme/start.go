package bme

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudkucooland/SmartAlac/pkg/cdio"
	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
	"github.com/ebitengine/purego"
)

var debug bool
var ripdir string
var encodedir string
var tagdir string
var finaldir string
var tuiProgram *tea.Program

var paranoiaLevel cdio.ParanoiaMode = cdio.ParanoiaModeFull ^ cdio.ParanoiaModeNeverSkip

func ToggleParanoia() string {
	if paranoiaLevel == (cdio.ParanoiaModeFull ^ cdio.ParanoiaModeNeverSkip) {
		paranoiaLevel = cdio.ParanoiaModeOverlap
		return "Good Enough (Overlap)"
	}
	paranoiaLevel = cdio.ParanoiaModeFull ^ cdio.ParanoiaModeNeverSkip
	return "Full Paranoia"
}

func GetParanoiaName() string {
	if paranoiaLevel == (cdio.ParanoiaModeFull ^ cdio.ParanoiaModeNeverSkip) {
		return "Full"
	}
	return "Fast"
}

func PurgeDirectories() error {
	dirs := []string{ripdir, encodedir, tagdir}
	for _, d := range dirs {
		entries, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			os.RemoveAll(filepath.Join(d, entry.Name()))
		}
	}
	return nil
}

func Debug(d bool) {
	if d {
		slog.Debug("enabling debug")
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	debug = d
}

var AppConfig *Config

func cdioLogHandler(level int, message *byte) {
	if tuiProgram == nil {
		return
	}

	// Manual conversion of *byte (C string) to Go string
	var b []byte
	ptr := uintptr(unsafe.Pointer(message))
	for {
		c := *(*byte)(unsafe.Pointer(ptr))
		if c == 0 {
			break
		}
		b = append(b, c)
		ptr++
	}
	msg := string(b)

	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}

	levelStr := "info"
	switch cdio.LogLevel(level) {
	case cdio.LogDebug:
		levelStr = "debug"
	case cdio.LogWarn:
		levelStr = "warn"
	case cdio.LogError, cdio.LogAssert:
		levelStr = "error"
	}

	tuiProgram.Send(StatusMsg{
		Component: "cdio",
		Status:    fmt.Sprintf("[%s] %s", levelStr, msg),
	})
}

func Start(ctx context.Context, cfg *Config, p *tea.Program) error {
	AppConfig = cfg
	ripdir = cfg.RipDir
	encodedir = cfg.EncodeDir
	tagdir = cfg.TagDir
	finaldir = cfg.DoneDir
	tuiProgram = p

	if err := os.MkdirAll(ripdir, 0755); err != nil {
		return fmt.Errorf("failed to create rip directory %s: %w", ripdir, err)
	}
	if err := os.MkdirAll(encodedir, 0755); err != nil {
		return fmt.Errorf("failed to create encode directory %s: %w", encodedir, err)
	}
	if err := os.MkdirAll(tagdir, 0755); err != nil {
		return fmt.Errorf("failed to create tag directory %s: %w", tagdir, err)
	}
	if err := os.MkdirAll(finaldir, 0755); err != nil {
		return fmt.Errorf("failed to create final directory %s: %w", finaldir, err)
	}

	if err := cdio.Init(); err != nil {
		return fmt.Errorf("failed to initialize cdio: %w", err)
	}
	if err := mb5.Init(); err != nil {
		return fmt.Errorf("failed to initialize mb5: %w", err)
	}

	// Register custom libcdio log handler to redirect stderr logs to TUI
	callback := purego.NewCallback(cdioLogHandler)
	cdio.LogSetHandler(callback)

	var wg sync.WaitGroup

	// start batch ripper(s)
	devices := cfg.Devices
	if len(devices) == 0 {
		devices = []string{cdio.GetDefaultDevice(nil)}
	}

	for _, dev := range devices {
		if dev == "" {
			continue
		}
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			cdio_worker(ctx, d, p)
		}(dev)
	}

	// start batch encoder
	wg.Add(1)
	go func() {
		defer wg.Done()
		encoder(ctx, p)
	}()

	// start batch tagger
	wg.Add(1)
	go func() {
		defer wg.Done()
		tagger(ctx, p)
	}()

	<-ctx.Done()
	if p != nil {
		p.Send(StatusMsg{Component: "system", Status: "Shutdown requested"})
	}

	slog.Debug("waiting for background processes to finish")
	wg.Wait()
	slog.Debug("shutdown complete")
	return nil
}
