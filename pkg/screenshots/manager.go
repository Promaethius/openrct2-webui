package screenshots

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/promaethius/openrct2-webui/pkg/plugin"
)

type Screenshot struct {
	Name      string
	Timestamp time.Time
	Image     image.Image
}

type Manager struct {
	m sync.Mutex

	c    *plugin.Client
	path string

	interval time.Duration

	retain      uint32
	screenshots []Screenshot
}

func (m *Manager) addScreenshot(name string, t time.Time, i image.Image) {
	m.m.Lock()
	defer m.m.Unlock()

	m.screenshots = append(m.screenshots, Screenshot{name, t, i})

	if len(m.screenshots) > int(m.retain) {
		os.Remove(filepath.Join(m.path, m.screenshots[0].Name))
		m.screenshots = m.screenshots[1:]
	}
}

func (m *Manager) ScanDirectory() error {
	slog.Debug("building list of existing screenshots")

	images, err := os.ReadDir(m.path)
	if err != nil {
		return err
	}

	sort.Slice(images, func(i, j int) bool {
		iInfo, err := images[i].Info()
		if err != nil {
			return false
		}
		jInfo, err := images[j].Info()
		if err != nil {
			return false
		}

		return iInfo.ModTime().Before(jInfo.ModTime())
	})

	for _, i := range images {
		if i.IsDir() {
			continue
		}
		if filepath.Ext(i.Name()) != ".png" {
			continue
		}

		info, err := i.Info()
		if err != nil {
			return err
		}

		file, err := os.OpenFile(filepath.Join(m.path, i.Name()), os.O_RDONLY, 0644)
		if err != nil {
			return err
		}

		image, err := png.Decode(file)
		file.Close()
		if err != nil {
			return err
		}

		slog.Debug("found screenshot", slog.String("name", i.Name()), slog.Time("time", info.ModTime()))

		m.addScreenshot(i.Name(), info.ModTime(), image)
	}

	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	t := time.NewTicker(m.interval)
	defer t.Stop()

	slog.Debug("initial scanning screenshot directory")

	err := m.ScanDirectory()
	if err != nil {
		return err
	}

	for {
		var now time.Time
		select {
		case <-ctx.Done():
			return nil
		case now = <-t.C:
		}

		filename := fmt.Sprintf("%s.png", uuid.New())

		slog.Debug("generating screenshot", slog.String("name", filename), slog.Time("time", now))

		_, err := m.c.Command(fmt.Sprintf("context.captureImage({filename: \"%s\", zoom: 0, rotation: 0})", filename))
		if err != nil {
			return err
		}

		slog.Debug("scanning screenshot directory")

		err = m.ScanDirectory()
		if err != nil {
			return err
		}
	}
}

func (m *Manager) GetScreenshots() []Screenshot {
	m.m.Lock()
	defer m.m.Unlock()

	return m.screenshots
}

func NewManager(client *plugin.Client, path string, interval time.Duration, retain uint32) (*Manager, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if interval < time.Second {
		return nil, errors.New("interval cannot be less than one second")
	}

	return &Manager{
		c:        client,
		path:     path,
		interval: interval,
		retain:   retain,
	}, nil
}
