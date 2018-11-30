package taglog

import (
	"strings"
)

type LevelSet struct {
	levels       map[string]int // use a map to avoid searches
	defaultLevel string
}

// Some convenience constants. You can use whatever values you want.
const (
	LevelTrace     = "TRACE"
	LevelDebug     = "DEBUG"
	LevelInfo      = "INFO"
	LevelDefault   = "DEFAULT"
	LevelNotice    = "NOTICE"
	LevelVerbose   = "VERBOSE"
	LevelWarning   = "WARNING"
	LevelWarn      = "WARN"
	LevelError     = "ERROR"
	LevelErr       = "ERR"
	LevelFatal     = "FATAL"
	LevelCritical  = "CRITICAL"
	LevelAlert     = "ALERT"
	LevelEmergency = "EMERGENCY"
	LevelAll       = "ALL"
	LevelOff       = "OFF"
	LevelFine      = "FINE"
	LevelFiner     = "FINER"
	LevelFinest    = "FINEST"
)

func NewLevelSet(levels []string, defaultLevel string) *LevelSet {
	ls := new(LevelSet)
	ls.levels = make(map[string]int)

	for i, lvl := range levels {
		ls.levels[strings.ToUpper(lvl)] = i
	}

	defaultNom := strings.ToUpper(defaultLevel)
	_, found := ls.levels[defaultNom]
	if !found {
		if len(levels) > 0 {
			defaultLevel = strings.ToUpper(levels[0])
		}
	} else {
		ls.defaultLevel = defaultLevel
	}

	return ls
}

func (ls *LevelSet) Default() string {
	return ls.defaultLevel
}

func (ls *LevelSet) Less(a, b string) bool {
	return ls.levels[strings.ToUpper(a)] < ls.levels[strings.ToUpper(b)]
}

func (ls *LevelSet) Contains(lvl string) bool {
	_, found := ls.levels[strings.ToUpper(lvl)]
	return found
}

var DefaultLevelSet = NewLevelSet([]string{
	LevelDebug,
	LevelInfo,
	LevelNotice,
	LevelWarning,
	LevelError,
	LevelCritical,
	LevelAlert,
	LevelEmergency,
}, LevelInfo)

func (this *Logger) DefineLevels(ls *LevelSet) {
	this.levelset = ls
	this.level = ls.Default()
}

func (this *Logger) SetLevel(lvl string) {
	if lvl == "" || !this.levelset.Contains(lvl) {
		this.level = ""
		return
	}
	this.level = strings.ToUpper(lvl)
}

func (this *Logger) GetLevel() string {
	return this.level
}

func (this *Logger) SetLevelTag(tag string) {
	this.levelTag = tag
}

func (this *Logger) SetStandardLevel(lvl string) {
	this.standardLevel = lvl
}

func DefineLevels(ls *LevelSet) {
	std.DefineLevels(ls)
}

func SetLevel(lvl string) {
	std.SetLevel(lvl)
}

func GetLevel() string {
	return std.GetLevel()
}

func SetLevelTag(tag string) {
	std.SetLevelTag(tag)
}

func SetStandardLevel(lvl string) {
	std.SetStandardLevel(lvl)
}
