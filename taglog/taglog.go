package taglog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	FormatPlain = iota // Plain log format. Simple format for easy human reading.
	FormatJSON         // JSON log format. Each log line is a JSON blob for easy machine reading.
)

// Get a log format from a string.
func ParseFormat(fmt string) int {
	switch strings.ToLower(fmt) {
	case "formatplain":
		return FormatPlain
	case "formatjson":
		return FormatJSON
	}
	return -1
}

const (
	// Format type used in conjunction with flags to set the timestamp format.
	TimestampFormatTypeISO     = iota // ISO 8601 timestamp format in UTC
	TimestampFormatTypeStd            // Standard log package timestamp format
	TimestampFormatTypeUnknown        // Custom timestamp format
)

// Get a timestamp format type from a string.
func ParseTimestampFormatType(fmt string) int {
	switch strings.ToLower(fmt) {
	case "timestampformattypeiso":
		return TimestampFormatTypeISO
	case "timestampformatstd":
		return TimestampFormatTypeStd
	}
	return -1
}

const (
	// Format used for printing the timestamp. See time.Time.Format()
	TimestampFormatISO     = time.RFC3339          // ISO 8601 format with date and time
	TimestampFormatISOTime = "15:04:05Z07:00"      // ISO 8601 format with time only
	TimestampFormatISODate = "2006-01-02"          // ISO 8601 format with date only
	TimestampFormatStd     = "2006/01/02 15:04:05" // Standard format with date and time
	TimestampFormatStdTime = "15:04:05"            // Standard format with time only
	TimestampFormatStdDate = "2006/01/02"          // Standard format with date only
)

// Get a timestamp format from a string. If it does not match a predefined
// format name, it is assumed to be a custom format and the string is returned
// as-is.
func ParseTimestampFormat(fmt string) string {
	switch strings.ToLower(fmt) {
	case "timestampformatiso":
		return TimestampFormatISO
	case "timestampformatisotime":
		return TimestampFormatISOTime
	case "timestampformatisodate":
		return TimestampFormatISODate
	case "timestampformatstd":
		return TimestampFormatStd
	case "timestampformatstdtime":
		return TimestampFormatStdTime
	case "timestampformatstddate":
		return TimestampFormatStdDate
	}
	return fmt
}

const (
	// Bits or'ed together to control what's printed.
	Ldate         = 1 << iota     // the date
	Ltime                         // the time
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // IGNORED
	Lshortfile                    // IGNORED
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger

	// Custom taglog flags (hoping they don't add 12 more standard tags)
	Lmilliseconds = 1 << 16 // millisecond resolution: 01:23:23.123.  assumes Ltime.
)

// Get flags from a string.
func ParseFlags(flags string) int {
	var out int
	flagss := strings.Split(flags, "|")
	for _, f := range flagss {
		f = strings.ToLower(strings.TrimSpace(f))
		switch f {
		case "ldate":
			out |= Ldate
		case "ltime":
			out |= Ltime
		case "lmicroseconds":
			out |= Lmicroseconds
		case "llongfile":
			out |= Llongfile
		case "lshortfile":
			out |= Lshortfile
		case "lutc":
			out |= LUTC
		case "lstdflags":
			out |= LstdFlags
		case "lmilliseconds":
			out |= Lmilliseconds
		}
	}
	return out
}

// Parameters which control how the logs are formatted
type Params struct {
	Format              int
	TimestampFormat     string
	TimestampFormatType int
	Prefix              string
	Flag                int
}

// taglog counterpart to the log.Logger type
type Logger struct {
	mu            sync.Mutex
	tags          Tags
	levelset      *LevelSet
	level         string
	levelTag      string
	standardLevel string
	out           io.Writer
	params        Params
}

// See log.New
func New(out io.Writer, prefix string, flag int) *Logger {
	tl := new(Logger)
	tl.tags = make(Tags)
	tl.out = out
	tl.params = DefaultParams
	tl.params.Prefix = prefix
	tl.params.Flag = flag
	tl.DefineLevels(DefaultLevelSet)
	tl.levelTag = "level"
	return tl
}

var std = New(os.Stderr, "", LstdFlags)

// Default parameters
var DefaultParams = Params{
	Format:              FormatPlain,
	TimestampFormat:     TimestampFormatStd,
	TimestampFormatType: TimestampFormatTypeStd,
	Prefix:              "",
	Flag:                LstdFlags,
}

// Directly access the standard global logger
func Global() *Logger {
	return std
}

// Create a new Logger by copying the formatting and tags from another Logger.
func (this *Logger) Copy() *Logger {
	this.mu.Lock()
	tl := *this

	// deep copy tags
	tl.tags = make(Tags)
	for k, v := range this.tags {
		tl.tags[k] = v
	}

	tl.mu.Unlock()
	this.mu.Unlock()
	return &tl
}

// Generate the timestamp format from the type and flags
func GenTimestampFormat(tsFormatType int, flag int) string {
	switch tsFormatType {
	case TimestampFormatTypeISO:
		if flag&(Ldate|Ltime) != 0 {
			if flag&Ldate != 0 && flag&Ltime != 0 {
				if flag&(Lmicroseconds) != 0 {
					return "2006-01-02T15:04:05.000000Z07:00"
				} else if flag&(Lmilliseconds) != 0 {
					return "2006-01-02T15:04:05.000Z07:00"
				} else {
					return TimestampFormatISO
				}
			} else if flag&Ldate != 0 {
				return TimestampFormatISODate
			} else {
				if flag&(Lmicroseconds) != 0 {
					return "15:04:05.000000Z07:00"
				} else if flag&(Lmilliseconds) != 0 {
					return "15:04:05.000Z07:00"
				} else {
					return TimestampFormatISOTime
				}
			}
		}
	case TimestampFormatTypeStd:
		if flag&(Ldate|Ltime) != 0 {
			if flag&Ldate != 0 && flag&Ltime != 0 {
				if flag&(Lmicroseconds) != 0 {
					return "2006/01/02 15:04:05.000000"
				} else if flag&(Lmilliseconds) != 0 {
					return "2006/01/02 15:04:05.000"
				} else {
					return TimestampFormatStd
				}
			} else if flag&Ldate != 0 {
				return TimestampFormatStdDate
			} else {
				if flag&(Lmicroseconds) != 0 {
					return "15:04:05.000000"
				} else if flag&(Lmilliseconds) != 0 {
					return "15:04:05.000"
				} else {
					return TimestampFormatStdTime
				}
			}
		}
	}
	return ""
}

func calcTsFormat(params *Params) string {
	if params.TimestampFormat != "" {
		return params.TimestampFormat
	}
	return GenTimestampFormat(params.TimestampFormatType, params.Flag)
}

// See log.Logger.Output
func (this *Logger) Output(s string) error {
	return this.Loutput(this.standardLevel, s)
}

// See log.Logger.Output
func (this *Logger) Loutput(level string, s string) error {
	var err error
	var b []byte

	now := time.Now()
	if this.params.Flag&(LUTC) != 0 {
		now = now.UTC()
	}
	this.mu.Lock()
	defer this.mu.Unlock()

	tsFormat := calcTsFormat(&this.params)
	nowStr := now.Format(tsFormat)

	if level != "" && this.levelset != nil && this.level != "" {
		// discard messages lower than the current log level
		if this.levelset.Less(level, this.level) {
			return nil
		}

		// set level tag
		if this.levelTag != "" {
			if this.levelset.Contains(level) {
				this.tags.Set(this.levelTag, strings.ToUpper(level))
				defer this.tags.Del(this.levelTag)
			}
		}
	}

	if this.params.Format == FormatJSON {
		if nowStr != "" {
			this.tags.Set("timestamp", nowStr)
		}
		this.tags.Set("msg", s)

		b, err = json.Marshal(&this.tags)
		if err != nil {
			return err
		}
	} else if this.params.Format == FormatPlain {
		line := []string{}
		if nowStr != "" {
			line = append(line, nowStr)
		}
		lineTags := []string{}
		for k, v := range this.tags {
			switch vs := v.(type) {
			case string:
				if k == "tags" {
					lineTags = append(lineTags, fmt.Sprintf("[%s]", vs))
				} else {
					lineTags = append(lineTags, fmt.Sprintf("[%s=%s]", k, vs))
				}
			case []string:
				if k == "tags" {
					for _, v0 := range vs {
						lineTags = append(lineTags, fmt.Sprintf("[%s]", v0))
					}
				} else {
					lineTags = append(lineTags, fmt.Sprintf("[%s=%s]", k, strings.Join(vs, ",")))
				}
			}
		}
		sort.Strings(lineTags)
		line = append(line, lineTags...)
		if s != "" {
			line = append(line, s)
		}
		b = []byte(this.params.Prefix + strings.Join(line, " "))
	}

	b = append(b, '\n')
	_, err = this.out.Write(b)
	return err
}

// Get the formatting parameters.
func (this *Logger) Params() Params {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.params
}

// See log.Logger.SetFlags
func (this *Logger) SetFlags(flag int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.params.Flag = flag
}

// See log.Logger.Flags
func (this *Logger) Flags() int {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.params.Flag
}

// See log.Logger.SetPrefix
func (this *Logger) SetPrefix(prefix string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.params.Prefix = prefix
}

// See log.Logger.Prefix
func (this *Logger) Prefix() string {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.params.Prefix
}

// Set the timestamp format type.
func (this *Logger) SetTimestampFormatType(tsFormatType int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.params.TimestampFormatType = tsFormatType
	this.params.TimestampFormat = ""
}

// Get the timestamp format type.
func (this *Logger) TimestampFormatType() int {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.params.TimestampFormatType
}

// Set the timestamp format.
func (this *Logger) SetTimestampFormat(tsFormat string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	tsFormat = ParseTimestampFormat(tsFormat)
	this.params.TimestampFormat = tsFormat
	switch tsFormat {
	case TimestampFormatISO, TimestampFormatISOTime, TimestampFormatISODate:
		this.params.TimestampFormatType = TimestampFormatTypeISO
	case TimestampFormatStd, TimestampFormatStdTime, TimestampFormatStdDate:
		this.params.TimestampFormatType = TimestampFormatTypeStd
	default:
		this.params.TimestampFormatType = TimestampFormatTypeUnknown
	}
}

// Get the timestamp format.
func (this *Logger) TimestampFormat() string {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.params.TimestampFormat
}

// Set the log format.
func (this *Logger) SetFormat(format int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.params.Format == FormatJSON && format == FormatPlain {
		this.tags.Del("timestamp")
		this.tags.Del("msg")
	}
	this.params.Format = format
}

// Get the log format.
func (this *Logger) Format() int {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.params.Format
}

// Add one or more values to a key.
func (this *Logger) AddTag(key string, value ...string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	this.tags.Add(key, value...)
}

// Add one or more values to a key, merging any duplicate values.
func (this *Logger) MergeTag(key string, value ...string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	this.tags.Merge(key, value...)
}

// Append one or more values to a key. This the same as AddTag and is only
// provided to couple with Pop() for code clarity.
func (this *Logger) PushTag(key string, value ...string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	this.tags.Push(key, value...)
}

// Remove the last value for a key
func (this *Logger) PopTag(key string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	this.tags.Pop(key)
}

// Set one or more values for a key. Any existing values are discarded.
func (this *Logger) SetTag(key string, value ...string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	this.tags.Set(key, value...)
}

// Get the first value for a key. If the key does not exist, an empty string is
// returned.
func (this *Logger) GetTag(key string) string {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	return this.tags.Get(key)
}

// Get all the values for a key. If the key does not exist, a nil slice is
// returned.
func (this *Logger) GetTags(key string) []string {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		key = "tags"
	}
	return this.tags.GetAll(key)
}

// Delete a key.
func (this *Logger) DelTag(key string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if key == "" {
		this.tags.Del("tags")
	}
	this.tags.Del(key)
}

// Delete all keys.
func (this *Logger) DelTags() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.tags.DelAll()
}

// Export all tags as a map of string slices.
func (this *Logger) ExportTags() map[string][]string {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.tags.Export()
}

// Import tags from a map of string slices.
func (this *Logger) ImportTags(tags map[string][]string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.tags.Import(tags)
}

// Set the output Writer.
func (this *Logger) SetOutput(w io.Writer) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.out = w
}

// Get the output Writer.
func (this *Logger) GetOutput() io.Writer {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.out
}

// Parse tags from a list of "key=value" strings.
func (this *Logger) ParseTags(tags []string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, s := range tags {
		ss := strings.Split(s, "=")
		if len(ss) == 1 {
			if len(ss[0]) > 0 {
				this.tags.Add("tags", ss[0])
			}
		} else if len(ss) > 1 {
			key := ss[0]
			if key == "" {
				key = "tags"
			}
			this.tags.Merge(key, ss[1])
		}
	}
}

// See log.Logger.Printf
func (this *Logger) Printf(format string, v ...interface{}) {
	this.Output(fmt.Sprintf(format, v...))
}

// See log.Logger.Print
func (this *Logger) Print(v ...interface{}) {
	this.Output(fmt.Sprint(v...))
}

// See log.Logger.Println
func (this *Logger) Println(v ...interface{}) {
	this.Output(fmt.Sprint(v...))
}

func (this *Logger) Lprintf(level string, format string, v ...interface{}) {
	this.Loutput(level, fmt.Sprintf(format, v...))
}

func (this *Logger) Lprint(level string, v ...interface{}) {
	this.Loutput(level, fmt.Sprint(v...))
}

func (this *Logger) Lprintln(level string, v ...interface{}) {
	this.Loutput(level, fmt.Sprint(v...))
}

// See log.Logger.Fatal
func (this *Logger) Fatal(v ...interface{}) {
	this.Output(fmt.Sprint(v...))
	os.Exit(1)
}

// See log.Logger.Fatalf
func (this *Logger) Fatalf(format string, v ...interface{}) {
	this.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

// See log.Logger.Fatalln
func (this *Logger) Fatalln(v ...interface{}) {
	this.Output(fmt.Sprintln(v...))
	os.Exit(1)
}

func (this *Logger) Lfatal(level string, v ...interface{}) {
	this.Loutput(level, fmt.Sprint(v...))
	os.Exit(1)
}

func (this *Logger) Lfatalf(level string, format string, v ...interface{}) {
	this.Loutput(level, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (this *Logger) Lfatalln(level string, v ...interface{}) {
	this.Loutput(level, fmt.Sprintln(v...))
	os.Exit(1)
}

// See log.Logger.Panic
func (this *Logger) Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	this.Output(s)
	panic(s)
}

// See log.Logger.Panicf
func (this *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	this.Output(s)
	panic(s)
}

// See log.Logger.Panicln
func (this *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	this.Output(s)
	panic(s)
}

// Create a new Logger by copying the formatting and tags from the Standard Logger.
func Copy() *Logger {
	return std.Copy()
}

// See log.SetFlags
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// See log.Flags
func Flags() int {
	return std.Flags()
}

// See log.SetPrefix
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// See log.Prefix
func Prefix() string {
	return std.Prefix()
}

// Set the timestamp format type for the Standard Logger.
func SetTimestampFormatType(tsFormatType int) {
	std.SetTimestampFormatType(tsFormatType)
}

// Get the timestamp format type for the Standard Logger.
func TimestampFormatType() int {
	return std.TimestampFormatType()
}

// Set the timestamp format for the Standard Logger.
func SetTimestampFormat(tsFormat string) {
	std.SetTimestampFormat(tsFormat)
}

// Get the timestamp format for the Standard Logger.
func TimestampFormat() string {
	return std.TimestampFormat()
}

// Set the log format for the Standard Logger.
func SetFormat(format int) {
	std.SetFormat(format)
}

// Get the log format for the Standard Logger.
func Format() int {
	return std.Format()
}

// Add one or more values to a key.
func AddTag(key string, value ...string) {
	std.AddTag(key, value...)
}

// Add one or more values to a key, merging any duplicate values.
func MergeTag(key string, value ...string) {
	std.MergeTag(key, value...)
}

// Append one or more values to a key. This the same as AddTag and is only
// provided to couple with Pop() for code clarity.
func PushTag(key string, value ...string) {
	std.PushTag(key, value...)
}

// Remove the last value for a key
func PopTag(key string) {
	std.PopTag(key)
}

// Set one or more values for a key. Any existing values are discarded.
func SetTag(key string, value ...string) {
	std.SetTag(key, value...)
}

// Get the first value for a key. If the key does not exist, an empty string is
// returned.
func GetTag(key string) string {
	return std.GetTag(key)
}

// Get all the values for a key. If the key does not exist, a nil slice is
// returned.
func GetTags(key string) []string {
	return std.GetTags(key)
}

// Delete a key.
func DelTag(key string) {
	std.DelTag(key)
}

// Delete all keys.
func DelTags() {
	std.DelTags()
}

// Export all tags as a map of string slices.
func ExportTags() map[string][]string {
	return std.ExportTags()
}

// Import tags from a map of string slices.
func ImportTags(tags map[string][]string) {
	std.ImportTags(tags)
}

// See log.SetOutput
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// Get the output Writer.
func GetOutput() io.Writer {
	return std.GetOutput()
}

// Parse tags from a list of "key=value" strings.
func ParseTags(tags []string) {
	std.ParseTags(tags)
}

// See log.Printf
func Printf(format string, v ...interface{}) {
	std.Output(fmt.Sprintf(format, v...))
}

// See log.Print
func Print(v ...interface{}) {
	std.Output(fmt.Sprint(v...))
}

// See log.Println
func Println(v ...interface{}) {
	std.Output(fmt.Sprint(v...))
}

func Lprintf(level string, format string, v ...interface{}) {
	std.Loutput(level, fmt.Sprintf(format, v...))
}

func Lprint(level string, v ...interface{}) {
	std.Loutput(level, fmt.Sprint(v...))
}

func Lprintln(level string, v ...interface{}) {
	std.Loutput(level, fmt.Sprint(v...))
}

// See log.Fatal
func Fatal(v ...interface{}) {
	std.Output(fmt.Sprint(v...))
	os.Exit(1)
}

// See log.Fatalf
func Fatalf(format string, v ...interface{}) {
	std.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

// See log.Fatalln
func Fatalln(v ...interface{}) {
	std.Output(fmt.Sprintln(v...))
	os.Exit(1)
}

func Lfatal(level string, v ...interface{}) {
	std.Loutput(level, fmt.Sprint(v...))
	os.Exit(1)
}

func Lfatalf(level string, format string, v ...interface{}) {
	std.Loutput(level, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Lfatalln(level string, v ...interface{}) {
	std.Loutput(level, fmt.Sprintln(v...))
	os.Exit(1)
}

// See log.Panic
func Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(s)
	panic(s)
}

// See log.Panicf
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(s)
	panic(s)
}

// See log.Panicln
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(s)
	panic(s)
}
