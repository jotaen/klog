# Specification (v1)

## Directory structure

```
2019/
2020/
	01/
	02/
	03/
		12.yml
		13.yml
```

## File format

```yaml
date: 2020-01-31
summary: Some text
tags: [custom, keywords]
worktime: 8h
log:
- start: 09:00
	end: 1:00pm
- time: 1:30
- time: 2h 15m
```

### Field `date`
- Must be present
- Is of type `date`

### Field `summary`
- Is optional
- Must be a string (single or multiline)

### Field `tags`
- Is optional
- Must be a list of string values

### Field `worktime`
- Is optional
- Is of type “period of time”

### Field `log`
- Is optional
- Must be a list of times, either denotes as range or duration
	- Range:
		- Must have a `start` and an `end` field
		- Is of type “time of day”
	- Duration (period of time):
		- Must have a `time` field
		- Is of type “period of time”

### Type “date”
- Denotes a day in the gregorian calendar
- Must be formatted either `YYYY-MM-DD` or `DD.MM.YYYY`

### Type “time of day”
- Denotes a point-in-time of the day (as represented by a wall clock)
- 1–2 *hour* digits (the leading `0` is optional), followed by a colon, followed by 2 *minute* digits. If `am` or `pm` suffix is present, the hour range is 0–12, otherwise it’s 0–24.

### Type “period of time”
- Denotes a duration (as represented by a stop watch)
- Either *hour* digits (leading `0` is optional), followed by a colon, followed by *minute* digits.
- Or *hour* digits with a `h` suffix, followed by a space, optionally followed by two *minute* digits with `m` suffix.

### Version marker
- Is optional
- Must be the first line in the file (prefixed by `# `)
- Contains the version number, which is an unsigned integer prefixed by a `v`
- When absent the file is assumed to follow the latest version of the specification.

### Additional properties
- No other property is allowed to appear
