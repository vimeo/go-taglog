/*
   Package taglog is based on, and compatible with, the Go standard log
   package, but it also provides additional functionality and features.

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

   Drop-in Replacement

   The easiest way to use TagLog as a drop-in replacement for the standard log
   package is to add it to your imports as:

       import (
           log "github.com/vimeo/go-taglog/taglog"
       )


   Conflicts

       - The flags Llongfile and Lshortfile are defined for compatibility, but they do not have any effect
       - The flags Ldate, Ltime, and Lmicroseconds only apply when not using a custom timestamp format
           - If SetTimestampFormat() is called with an undefined value, the flags are subsequently ignored
           - If SetTimestampFormatType() is called, the timestamp format is reset and the flags will be used
           - With TimestampFormatTypeISO, the Lmicroseconds flag actually prints milliseconds

   Tags

   You can add, delete, get, set, push, and pop tags as key/value strings. A single
   key can have multiple values.

   Special-case tags:

       - In JSON format, the "timestamp" and "msg" tags are overwritten when logging a line
           - When switching from JSON to plain format, the "timestamp" and "msg" tags are deleted
       - Using "tags" or an empty string ("") as the key treats values as a global tags without a key
           - In plain format, global tags have the key omitted and are printed separately
           - In JSON format, global tags are exported in the "tags" field
           - Using GetTag() with either "" or "tags" will access global tags

   Defaults

   The default values will result in identical behavior to the Go standard log package.
       - format is FormatPlain
       - timestamp format type is TimestampFormatTypeStd
       - timestamp format is TimestampFormatStd
       - flags are LstdFlags
       - output is os.Stderr
*/
package taglog
