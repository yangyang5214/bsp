package pkg

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/multierr"
	"os"
	"runtime"
	"strings"
)

type ChromePool struct {
	previousPIDs map[int32]struct{}
	log          *log.Helper
	browser      *rod.Browser
}

func NewChromePool() (*ChromePool, error) {
	dataStore, _ := os.MkdirTemp("", "bsp-*")
	chromeLauncher := launcher.New().
		Leakless(false).
		Set("disable-gpu", "true").
		Set("ignore-certificate-errors", "true").
		Set("ignore-certificate-errors", "1").
		Set("disable-crash-reporter", "true").
		Set("disable-notifications", "true").
		Set("hide-scrollbars", "true").
		Set("window-size", fmt.Sprintf("%d,%d", 1080, 1920)).
		Set("mute-audio", "true").
		Delete("use-mock-keychain").
		Env(append(os.Environ(), "TZ=Asia/Shanghai")...).
		UserDataDir(dataStore)

	if runtime.GOOS == "darwin" {
		chromeLauncher.Headless(false)
	}
	launcherURL, err := chromeLauncher.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(launcherURL)
	if browserErr := browser.Connect(); browserErr != nil {
		return nil, browserErr
	}

	previousPIDs := findChromeProcesses()

	return &ChromePool{
		previousPIDs: previousPIDs,
		log:          log.NewHelper(log.DefaultLogger),
		browser:      browser,
	}, nil
}

func (s *ChromePool) Clone() error {
	return s.browser.Close()
}

func (s *ChromePool) NavigateUrl(urlStr string, process func(*rod.Page) error) (string, error) {
	opts := proto.TargetCreateTarget{}
	p, err := s.browser.Page(opts)
	if err != nil {
		return "", err
	}
	defer p.Close()

	err = p.Navigate(urlStr)
	if err != nil {
		return "", err
	}

	err = p.WaitLoad()
	if err != nil {
		return "", err
	}

	err = process(p)
	if err != nil {
		return "", err
	}

	return p.HTML()
}

func (c *ChromePool) killChromeProcesses() error {
	var errs []error
	processes, _ := process.Processes()

	for _, p := range processes {
		// skip non-chrome processes
		if !isChromeProcess(p) {
			continue
		}

		// skip chrome processes that were already running
		if _, ok := c.previousPIDs[p.Pid]; ok {
			continue
		}

		if err := p.Kill(); err != nil {
			c.log.Infof("kill chrome process error %d, %s", p.Pid, err.Error())
			errs = append(errs, err)
		}
	}

	return multierr.Combine(errs...)
}

// findChromeProcesses finds chrome process running on host
func findChromeProcesses() map[int32]struct{} {
	processes, _ := process.Processes()
	list := make(map[int32]struct{})
	for _, p := range processes {
		if isChromeProcess(p) {
			list[p.Pid] = struct{}{}
			if ppid, err := p.Ppid(); err == nil {
				list[ppid] = struct{}{}
			}
		}
	}
	return list
}

// isChromeProcess checks if a process is chrome/chromium
func isChromeProcess(process *process.Process) bool {
	name, _ := process.Name()
	if name == "" {
		return false
	}
	return strings.HasPrefix(strings.ToLower(name), "chromium")
}
