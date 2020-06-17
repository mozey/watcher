package watcher

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const AppDir = "APP_DIR"

type MultiFlag []string

func (f *MultiFlag) String() string {
	return strings.Join(*f, " ")
}

func (f *MultiFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// Timeout executes the callback after the specified duration.
// Abort execution by sending a message on the cancel chan
func Timeout(cancel chan bool, duration time.Duration, cb func()) {
	if cancel == nil {
		return
	}
	select {
	case <-time.After(duration):
		cb()
	case <-cancel:
		// Abort
	}
}

// CmdIn for use with command functions
type CmdIn struct {
	// BaseDir for relative paths
	BaseDir string
	// Version
	Version bool
	// WatchDirs is the dirs to watch
	WatchDirs MultiFlag
	// Recursive can be set to watch sub dirs
	Recursive bool
	// Delay in milliseconds before printing changes
	Delay int
	// Limit sub dirs to watch
	Limit int
	// IncludeFiles matching patterns
	IncludeFiles MultiFlag
	// ExcludeDirs matching patterns
	ExcludeDirs MultiFlag
}

const CmdVersion = "version"
const CmdWatch = "watch"

// CmdOut for use with Cmd function
type CmdOut struct {
	// Cmd
	Cmd string
	// Watcher
	Watcher *fsnotify.Watcher
}

func ParseFlags() *CmdIn {
	in := CmdIn{}

	flag.BoolVar(&in.Version, "version", false, "Print version")
	flag.BoolVar(&in.Recursive, "r", false, "Recursively watch sub dirs")
	flag.IntVar(&in.Limit, "l", 100, "Limit dirs to include recursively")
	flag.IntVar(&in.Delay, "d", 1000,
		"Delay in milliseconds before printing changes")
	flag.StringVar(&in.BaseDir, "b", "", "Base dir for relative paths")
	flag.Var(&in.WatchDirs, "dir", "Dirs to watch")
	flag.Var(&in.IncludeFiles, "include", "Only include matching files")
	flag.Var(&in.ExcludeDirs, "excludeDir", "Exclude matching dirs")
	flag.Parse()

	return &in
}

func (in *CmdIn) FileIncluded(p string) (included bool, err error) {
	if len(in.IncludeFiles) == 0 {
		// All files are included by default
		return true, nil
	}
	// TODO Compile patterns once and cache
	for _, includeFile := range in.IncludeFiles {
		match, err := regexp.MatchString(includeFile, p)
		if err != nil {
			return included, errors.WithStack(err)
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (in *CmdIn) DirExcluded(p string) (excluded bool, err error) {
	if len(in.ExcludeDirs) == 0 {
		// No dirs are excluded by default
		return false, nil
	}
	// TODO Compile patterns once and cache
	for _, excludeDir := range in.ExcludeDirs {
		match, err := regexp.MatchString(excludeDir, p)
		if err != nil {
			return excluded, errors.WithStack(err)
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (in *CmdIn) Watch(watcher *fsnotify.Watcher) {
	var cancel chan bool
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Check if file must be included
			included, err := in.FileIncluded(event.Name)
			if err != nil {
				watcher.Errors <- err
				return
			}

			if included {
				// Cancel previous timeout if set
				if cancel != nil {
					close(cancel)
				}
				// Reset cancel chan
				cancel = make(chan bool)
				// Use a timeout in case multiple files were changed
				go Timeout(cancel, time.Duration(in.Delay)*time.Millisecond,
					func() {
						// Print path to last file that was modified
						fmt.Printf("%v\n", event.Name)
					})
			}
		}
	}
}

func Cmd(in *CmdIn) (out *CmdOut, err error) {
	out = &CmdOut{}

	if in.Version {
		out.Cmd = CmdVersion
		return out, nil
	}
	out.Cmd = CmdWatch

	out.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return out, errors.WithStack(err)
	}

	go in.Watch(out.Watcher)

	for _, relativePath := range in.WatchDirs {

		// Use absolute paths
		absolutePath := path.Join(in.BaseDir, relativePath)

		// Check dir exclusion filter
		excluded, err := in.DirExcluded(absolutePath)
		if err != nil {
			return out, errors.WithStack(err)
		}

		if !excluded {
			// Watch the specified dir
			log.Debug().Str("path", absolutePath).Msg("Add path")
			err = out.Watcher.Add(absolutePath)
			if err != nil {
				return out, errors.WithStack(err)
			}

			// Watch sub dirs recursively
			r := 0
			if in.Recursive {
				err = filepath.Walk(absolutePath,
					func(path string, info os.FileInfo, err error) error {
						// Check limit
						if r < in.Limit {
							if info.IsDir() {
								if strings.HasPrefix(info.Name(), ".") {
									// Skip hidden dirs
									return filepath.SkipDir
								} else {
									// Check dir exclusion filter
									excluded, err := in.DirExcluded(path)
									if err != nil {
										return errors.WithStack(err)
									}
									if excluded {
										// Skip excluded dirs
										return filepath.SkipDir
									} else {
										// Watch sub dir
										log.Debug().Str("path", path).
											Msg("Add sub path")
										err = out.Watcher.Add(path)
										if err != nil {
											return errors.WithStack(err)
										}
										r++
									}
								}
							}
						}
						return nil
					})
				if err != nil {
					return out, errors.WithStack(err)
				}
			}
		}
	}

	return out, nil
}

func Main() (out *CmdOut, err error) {
	// Parse flags
	in := ParseFlags()

	// Base dir is required, resolve in this order (flag, env, working dir)
	if in.BaseDir == "" {
		in.BaseDir = os.Getenv(AppDir)
		if in.BaseDir == "" {
			in.BaseDir, err = os.Getwd()
			if err != nil {
				return out, errors.WithStack(err)
			}
			_, err := os.Stat(in.BaseDir)
			if err != nil {
				if os.IsNotExist(err) {
					return out, errors.WithStack(
						fmt.Errorf("invalid base dir"))
				} else {
					return out, errors.WithStack(err)
				}
			}
		}
	}

	// Run cmd
	out, err = Cmd(in)
	if err != nil {
		return out, errors.WithStack(err)
	}

	return out, nil
}
