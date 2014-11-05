package taglog

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    "strings"
    "time"
)

// Used to parse logs created by taglog.
type Parser struct {
    params Params
    tags Tags
}

// Create a new Parser instance.
func NewParser(params Params) *Parser {
    p := new(Parser)
    p.params = params
    p.tags = make(Tags)
    return p
}

// Get all currently parsed tags as a map of string slices.
func (this *Parser) Tags() map[string][]string {
    return this.tags.Export()
}

// Clear all parsed tags.
func (this *Parser) Reset() {
    this.tags = make(Tags)
}

// Merge tags from a map of string slices.
func (this *Parser) MergeTags(newTags map[string][]string) {
    this.tags.Import(newTags)
}

func (this *Parser) parseLinePlain(line string, timestampFormat string) (Tags, error) {
    tags := make(Tags)

    if this.params.Prefix != "" {
        if !strings.HasPrefix(line, this.params.Prefix) {
            return nil, fmt.Errorf("Log format mismatch: prefix")
        }
        line = strings.TrimPrefix(line, this.params.Prefix)
    }

    tsFormat := calcTsFormat(&this.params)
    if tsFormat != "" {
        fmtTokens := len(strings.Split(tsFormat, " "))
        lineSplit := strings.Split(line, " ")
        if len(lineSplit) < fmtTokens {
            return nil, fmt.Errorf("Log format mismatch: timestamp")
        }
        lineSplit  = lineSplit[:fmtTokens]
        tsStr     := strings.Join(lineSplit, " ")

        ts, err := time.Parse(tsFormat, tsStr)
        if err != nil {
            return nil, fmt.Errorf("Log format mismatch: timestamp")
        }
        if timestampFormat == "" {
            tags.Add("timestamp", tsStr)
        } else {
            tags.Add("timestamp", ts.Format(timestampFormat))
        }

        line = strings.TrimPrefix(line, tsStr)
        line = strings.TrimLeft(line, " ")
    }

    tokenStart := -1
    msgStart   := 0
    var tagTokens []string
    for i, c := range line {
        if tokenStart < 0 {
            if c == ']' {
                break
            } else if c == '[' {
                tokenStart = i
            }
        } else {
            if c == ']' {
                tagTokens = append(tagTokens, line[tokenStart+1:i])
                tokenStart = -1
                msgStart = i + 2
            }
        }
    }
    if msgStart >= len(line) {
        return nil, fmt.Errorf("Empty line")
    }
    tags.Add("msg", line[msgStart:])

    for _, t := range tagTokens {
        pairs := strings.Split(t, "=")
        if len(pairs) == 1 {
            tags.Add("tags", pairs[0])
        } else {
            key    := pairs[0]
            values := strings.Split(pairs[1], ",")
            tags.Add(key, values...)
        }
    }

    return tags, nil
}

func (this *Parser) parseLineJSON(line string) (Tags, error) {
    tags := make(Tags)

    err := json.Unmarshal([]byte(line), &tags)
    if err != nil {
        return nil, err
    }

    return tags, nil
}

func (this *Parser) mergeLinePlain(line string) error {
    tags, err := this.parseLinePlain(line, "")
    if err != nil {
        return err
    }
    tags.Del("msg")
    tags.Del("timestamp")

    this.MergeTags(tags.Export())
    return nil
}

func (this *Parser) mergeLineJSON(line string) error {
    tags := make(Tags)

    err := json.Unmarshal([]byte(line), &tags)
    if err != nil {
        return err
    }
    tags.Del("msg")
    tags.Del("timestamp")

    this.MergeTags(tags.Export())
    return nil
}

// Parse a single log line.
func (this *Parser) ParseLine(line string) error {
    switch this.params.Format {
    case FormatPlain:
        return this.mergeLinePlain(line)
    case FormatJSON:
        return this.mergeLineJSON(line)
    }
    return fmt.Errorf("Invalid format")
}

func (this *Parser) parseInputPlain(input io.Reader) error {
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        this.ParseLine(scanner.Text())
    }
    return scanner.Err()
}

func (this *Parser) parseInputJSON(input io.Reader) error {
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        s := scanner.Text()
        err := this.ParseLine(s)
        if err != nil {
            return err
        }
    }
    return scanner.Err()
}

// Parse all lines from an io.Reader
func (this *Parser) ParseInput(input io.Reader) error {
    switch this.params.Format {
    case FormatPlain:
        return this.parseInputPlain(input)
    case FormatJSON:
        return this.parseInputJSON(input)
    }
    return fmt.Errorf("Invalid format")
}

// Convert Plain format input to JSON format output. timestampFormat specifies
// the output timestamp format. An empty string retains the timestamp format
// from the input.
func (this *Parser) PlainToJSON(input io.Reader, output io.Writer, timestampFormat string) error {
    var s string
    var lineTags Tags

    scanner := bufio.NewScanner(input)

    // get line 1
    // skip any non-starting lines at the beginning
    for scanner.Scan() {
        s = scanner.Text()
        tags, err := this.parseLinePlain(s, timestampFormat)
        if err == nil {
            lineTags = tags
            break
        }
    }
    if s == "" {
        return scanner.Err()
    }

    // get line 2
    for scanner.Scan() {
        s = scanner.Text()
        tags, err := this.parseLinePlain(s, timestampFormat)
        if err == nil {
            b, err := json.Marshal(&lineTags)
            if err == nil {
                fmt.Fprintln(output, string(b))
            }
            lineTags = tags
        } else {
            lineTags.Set("msg", lineTags.Get("msg") + "\n" + s)
        }
    }
    b, err := json.Marshal(&lineTags)
    if err == nil {
        fmt.Fprintln(output, string(b))
    }
    return nil
}
