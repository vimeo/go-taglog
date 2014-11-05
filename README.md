# TagLog #

**Go Log Package**

Based on, and compatible with, the Go standard log package, but also provides
additional functionality and features.

## Installation ##

```
go get github.com/vimeo/go-taglog/taglog
```

## Features ##

- Can be used as a drop-in replacement for the Go standard log package
- Default values mimic those of the Go standard log package
- Basic extra features
    - Set the timestamp format using the same format as time.Time.Format() rather than using flags
    - Ability to change the output writer for the Logger type
    - Ability to get the current output writer
- Advanced extra features
    - Add tags to the log format to add context to log messages and allow for easier machine processing
    - Output log lines in JSON format
    - Provides a pre-defined timestamp format that is compatible with elasticsearch (TimestampFormatISO)

## Details ##

**Drop-in Replacement**

The easiest way to use TagLog as a drop-in replacement for the standard package
is to use:
```go
import (
    log "github.com/vimeo/go-taglog/taglog"
)
```

Conflicts:
- The flags Llongfile and Lshortfile are defined for compatibility, but they do not have any effect
- The flags Ldate, Ltime, and Lmicroseconds only apply when not using a custom timestamp format
    - If SetTimestampFormat() is called with an undefined value, the flags are subsequently ignored
    - If SetTimestampFormatType() is called, the timestamp format is reset and the flags will be used
    - With TimestampFormatTypeISO, the Lmicroseconds flag actually prints milliseconds

**Tags**

You can add, delete, get, set, push, and pop tags as key/value strings. A single
key can have multiple values.

Special-case tags:
- In JSON format, the "timestamp" and "msg" tags are overwritten when logging a line
    - When switching from JSON to plain format, the "timestamp" and "msg" tags are deleted
- Using "tags" or an empty string ("") as the key treats values as a global tags without a key
    - In plain format, global tags have the key omitted and are printed separately
    - In JSON format, global tags are exported in the "tags" field
    - Using GetTag() with either "" or "tags" will access global tags

**Defaults**

The default values will result in identical behavior to the Go standard log
package.
- format is FormatPlain
- timestamp format type is TimestampFormatTypeStd
- timestamp format is TimestampFormatStd
- flags are LstdFlags
- output is os.Stderr

**Examples**

Defaults:
```
import log "github.com/vimeo/go-taglog/taglog"
2014/07/24 22:08:56 Message String
```

Set Timestamp Flags:
```
log.SetFlags(log.Ldate)
2014/07/24 Message String
log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
2014/07/24 22:08:56.840626 Message String
```

Set Timestamp Format Type:
```
log.SetTimestampFormatType(log.TimestampFormatTypeISO)
2014-07-24T22:08:56.840Z Message String
```

Set Timestamp Format:
```
log.SetTimestampFormat("(Mon Jan 2 15:04:05 2006)")
(Thu Jul 24 22:08:56 2014) Message String
log.SetTimestampFormat(log.TimestampFormatISO)
2014-07-24T22:08:56Z Message String
```

Add Some Global Tags:
```
log.AddTag("", "jobserver", "dev", "123456")
2014-07-24T22:08:56Z [jobserver] [dev] [123456] Message String
```

Push/Pop a Tag:
```
log.DelTags()
log.AddTag("", "jobserver")
log.PushTag("", "go-vimeo-http")
2014-07-24T22:08:56Z [jobserver] [go-vimeo-http] Message String
log.PopTag("")
2014-07-24T22:08:56Z [jobserver] Message String
```

Add Some Key/Value Tags:
```
log.DelTags()
log.AddTag("job_id", "123456")
log.AddTag("filters", "vol", "delay")
2014-07-24T22:08:56Z [job_id=123456] [filters=vol,delay] Message String
```

Copy Global Logger and Add Tags:
```
logger := log.Copy()
logger.AddTag("req_tag", "some_tag_value")
log.DelTags()
log.Println(message)
2014-07-28T16:17:11Z Message String
logger.Println(message)
2014-07-28T16:17:11Z [filters=vol,delay] [job_id=123456] [req_tag=some_tag_value] Message String
```

Set Format to JSON:
```
log.AddTag("", "jobserver")
log.SetFormat(log.FormatJSON)
{"filters":["vol","delay"],"job_id":"123456","msg":"Message String","tags":"jobserver","timestamp":"2014-07-24T22:08:56Z"}
```

Parse and Aggregate Tags:
```
2014-07-29T18:34:23Z [clip_id=84009894] [job_id=123456] Something Happened
2014-07-29T18:34:24Z [clip_id=84009894] [job_id=123457] Something Happened
2014-07-29T18:34:25Z [clip_id=84009894] [job_id=123456] Something Else Happened
2014-07-29T18:34:26Z [clip_id=84009894] [job_id=123458] Something Happened
That
Takes
Multiple
Lines
2014-07-29T18:34:27Z [clip_id=84009894] [job_id=123459] Something Happened
2014-07-29T18:34:28Z [clip_id=84009999] [job_id=456789] Something Happened

{
    "clip_id": [
        "84009894",
        "84009999"
    ],
    "job_id": [
        "123456",
        "123457",
        "123458",
        "123459",
        "456789"
    ]
}
```

Convert Log from Plain to JSON:
```
2014/07/29 18:34:23 [clip_id=84009894] [job_id=123456] Something Happened
2014/07/29 18:34:24 [clip_id=84009894] [job_id=123457] Something Happened
2014/07/29 18:34:25 [clip_id=84009894] [job_id=123456] Something Else Happened
2014/07/29 18:34:26 [clip_id=84009894] [job_id=123458] Something Happened
That
Takes
Multiple
Lines
2014/07/29 18:34:27 [clip_id=84009894] [job_id=123459] Something Happened
2014/07/29 18:34:28 [clip_id=84009999] [job_id=456789] Something Happened

{"clip_id":"84009894","job_id":"123456","msg":"Something Happened","timestamp":"2014-07-29T18:34:23Z"}
{"clip_id":"84009894","job_id":"123457","msg":"Something Happened","timestamp":"2014-07-29T18:34:24Z"}
{"clip_id":"84009894","job_id":"123456","msg":"Something Else Happened","timestamp":"2014-07-29T18:34:25Z"}
{"clip_id":"84009894","job_id":"123458","msg":"Something Happened\nThat\nTakes\nMultiple\nLines","timestamp":"2014-07-29T18:34:26Z"}
{"clip_id":"84009894","job_id":"123459","msg":"Something Happened","timestamp":"2014-07-29T18:34:27Z"}
{"clip_id":"84009999","job_id":"456789","msg":"Something Happened","timestamp":"2014-07-29T18:34:28Z"}
```
