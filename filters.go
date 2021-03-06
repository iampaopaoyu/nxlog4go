// Copyright (C) 2017, ccpaging <ccpaging@gmail.com>.  All rights reserved.

package nxlog4go

import ()

/****** Filters map ******/

// Filters represents a collection of Appenders through which log messages are
// written.
type Filters map[string]*Filter

// NewFilters creates a new filters
func NewFilters() Filters {
	return Filters{}
}

// Add a new filter to the filters map which will only log messages at lvl or
// higher. And call Appender.Init() to allow the appender protocol to perform
// any initialization steps it needs.
// This function should be called before install filters to logger by Logger.SetFilters(fs)
// Returns the Filters for chaining.
func (fs Filters) Add(tag string, level Level, writer Appender) Filters {
	if filt, ok := fs[tag]; ok {
		filt.Close()
		delete(fs, tag)
	}
	writer.Init()
	fs[tag] = NewFilter(level, writer)
	return fs
}

// Close and remove all filters in preparation for exiting the program or a
// reconfiguration of logging.  Calling this is not really imperative, unless
// you want to guarantee that all log messages are written.
// This function should be called after release filters by Logger.SetFilters(nil)
func (fs Filters) Close() {
	// Close all filters
	for tag, filt := range fs {
		filt.Close()
		delete(fs, tag)
	}
}

// Skip check log level and return whether skip or not
func (fs Filters) Skip(lvl Level) bool {
	for _, filt := range fs {
		if lvl >= filt.Level {
			return false
		}
	}
	return true
}

// Dispatch the logs
func (fs Filters) Dispatch(rec *LogRecord) {
	for _, filt := range fs {
		if rec.Level < filt.Level {
			continue
		}
		filt.writeToChan(rec)
	}
}
