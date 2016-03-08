package taglog

import (
	"fmt"
	"io"
	"os"
)

type MultiLogger struct {
	loggers []*Logger
}

func NewMultiLogger(loggers ...*Logger) *MultiLogger {
	mlog := new(MultiLogger)
	mlog.loggers = loggers
	return mlog
}

func (mlog *MultiLogger) Copy() *MultiLogger {
	newLoggers := make([]*Logger, len(mlog.loggers))
	for i, logger := range mlog.loggers {
		newLoggers[i] = logger.Copy()
	}
	return NewMultiLogger(newLoggers...)
}

func (mlog *MultiLogger) Output(s string) error {
	var anyErr error

	for _, logger := range mlog.loggers {
		err := logger.Output(s)
		if err != nil {
			anyErr = err
		}
	}

	return anyErr
}

func (mlog *MultiLogger) Loutput(level string, s string) error {
	var anyErr error

	for _, logger := range mlog.loggers {
		err := logger.Loutput(level, s)
		if err != nil {
			anyErr = err
		}
	}

	return anyErr
}

func (mlog *MultiLogger) Params() Params {
	if len(mlog.loggers) == 0 {
		return Params{}
	}
	return mlog.loggers[0].Params()
}

func (mlog *MultiLogger) SetFlags(flag int) {
	for _, logger := range mlog.loggers {
		logger.SetFlags(flag)
	}
}

func (mlog *MultiLogger) Flags() int {
	if len(mlog.loggers) == 0 {
		return 0
	}
	return mlog.loggers[0].Flags()
}

func (mlog *MultiLogger) SetPrefix(prefix string) {
	for _, logger := range mlog.loggers {
		logger.SetPrefix(prefix)
	}
}

func (mlog *MultiLogger) Prefix() string {
	if len(mlog.loggers) == 0 {
		return ""
	}
	return mlog.loggers[0].Prefix()
}

func (mlog *MultiLogger) SetTimestampFormatType(tsFormatType int) {
	for _, logger := range mlog.loggers {
		logger.SetTimestampFormatType(tsFormatType)
	}
}

func (mlog *MultiLogger) TimestampFormatType() int {
	if len(mlog.loggers) == 0 {
		return TimestampFormatTypeUnknown
	}
	return mlog.loggers[0].TimestampFormatType()
}

func (mlog *MultiLogger) SetTimestampFormat(tsFormat string) {
	for _, logger := range mlog.loggers {
		logger.SetTimestampFormat(tsFormat)
	}
}

func (mlog *MultiLogger) TimestampFormat() string {
	if len(mlog.loggers) == 0 {
		return ""
	}
	return mlog.loggers[0].TimestampFormat()
}

func (mlog *MultiLogger) SetFormat(format int) {
	for _, logger := range mlog.loggers {
		logger.SetFormat(format)
	}
}

func (mlog *MultiLogger) Format() int {
	if len(mlog.loggers) == 0 {
		return FormatPlain
	}
	return mlog.loggers[0].Format()
}

func (mlog *MultiLogger) AddTag(key string, value ...string) {
	for _, logger := range mlog.loggers {
		logger.AddTag(key, value...)
	}
}

func (mlog *MultiLogger) MergeTag(key string, value ...string) {
	for _, logger := range mlog.loggers {
		logger.MergeTag(key, value...)
	}
}

func (mlog *MultiLogger) PushTag(key string, value ...string) {
	for _, logger := range mlog.loggers {
		logger.PushTag(key, value...)
	}
}

func (mlog *MultiLogger) PopTag(key string) {
	for _, logger := range mlog.loggers {
		logger.PopTag(key)
	}
}

func (mlog *MultiLogger) SetTag(key string, value ...string) {
	for _, logger := range mlog.loggers {
		logger.SetTag(key, value...)
	}
}

func (mlog *MultiLogger) GetTag(key string) string {
	if len(mlog.loggers) == 0 {
		return ""
	}
	return mlog.loggers[0].GetTag(key)
}

func (mlog *MultiLogger) GetTags(key string) []string {
	if len(mlog.loggers) == 0 {
		return nil
	}
	return mlog.loggers[0].GetTags(key)
}

func (mlog *MultiLogger) DelTag(key string) {
	for _, logger := range mlog.loggers {
		logger.DelTag(key)
	}
}

func (mlog *MultiLogger) DelTags() {
	for _, logger := range mlog.loggers {
		logger.DelTags()
	}
}

func (mlog *MultiLogger) ExportTags() map[string][]string {
	if len(mlog.loggers) == 0 {
		return nil
	}
	return mlog.loggers[0].ExportTags()
}

func (mlog *MultiLogger) ImportTags(tags map[string][]string) {
	for _, logger := range mlog.loggers {
		logger.ImportTags(tags)
	}
}

func (mlog *MultiLogger) SetOutput(w io.Writer) {
	for _, logger := range mlog.loggers {
		logger.SetOutput(w)
	}
}

func (mlog *MultiLogger) GetOutput() io.Writer {
	if len(mlog.loggers) == 0 {
		return nil
	}
	return mlog.loggers[0].GetOutput()
}

func (mlog *MultiLogger) ParseTags(tags []string) {
	for _, logger := range mlog.loggers {
		logger.ParseTags(tags)
	}
}

func (mlog *MultiLogger) Printf(format string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Printf(format, v...)
	}
}

func (mlog *MultiLogger) Print(v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Print(v...)
	}
}

func (mlog *MultiLogger) Println(v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Println(v...)
	}
}

func (mlog *MultiLogger) Lprintf(level string, format string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Lprintf(level, format, v...)
	}
}

func (mlog *MultiLogger) Lprint(level string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Lprint(level, v...)
	}
}

func (mlog *MultiLogger) Lprintln(level string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Lprintln(level, v...)
	}
}

func (mlog *MultiLogger) Fatal(v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Print(v...)
	}
	os.Exit(1)
}

func (mlog *MultiLogger) Fatalf(format string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Printf(format, v...)
	}
	os.Exit(1)
}

func (mlog *MultiLogger) Fatalln(v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Println(v...)
	}
	os.Exit(1)
}

func (mlog *MultiLogger) Lfatal(level string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Lprint(level, v...)
	}
	os.Exit(1)
}

func (mlog *MultiLogger) Lfatalf(level string, format string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Lprintf(level, format, v...)
	}
	os.Exit(1)
}

func (mlog *MultiLogger) Lfatalln(level string, v ...interface{}) {
	for _, logger := range mlog.loggers {
		logger.Lprintln(level, v...)
	}
	os.Exit(1)
}

func (mlog *MultiLogger) Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	mlog.Output(s)
	panic(s)
}

func (mlog *MultiLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	mlog.Output(s)
	panic(s)
}

func (mlog *MultiLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	mlog.Output(s)
	panic(s)
}
