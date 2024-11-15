
# Configuration Documentation for `config.yml`

This document provides a comprehensive guide to configuring the `config.yml` file for filtering events based on various criteria, such as `camera`, `zone`, `label`, `sublabel`, `score`, `length`, `severity`, `media`, and `time`. The configuration structure supports both global and camera-specific settings, with options for flexible "and" and "or" filtering logic.

## Table of Contents

1. [Overview](#overview)
2. [Environment Variables](#environment-variables)
3. [Global and Camera-Specific Filters](#global-and-camera-specific-filters)
4. [Filter Types](#filter-types)
5. [Examples](#examples)
6. [Logic and Behavior](#logic-and-behavior)
7. [FAQ](#faq)

---

## Overview

This configuration file allows users to define global and camera-specific event filters to customize which events are included in notifications. You can create both broad (global) filters that apply across all cameras and specific filters that target individual cameras, zones, labels, or sublabels.

The configuration file uses both **"and"** and **"or"** filtering logic:
- **"And" Logic**: Achieved by nesting fields, meaning all conditions within a nested structure must be met.
- **"Or" Logic**: Configured by setting fields at the same level, meaning any of the specified conditions at that level will match.

This flexible approach allows for detailed filtering, where only the events that meet the specified criteria will be included.

---

## Environment Variables

The `environment` section contains general configuration options related to system setup and integration settings. These variables control behavior that is not related to event filtering but is essential for initializing and connecting the notification system.

### Example

```yaml
environment:
  TELEGRAM_BOT_TOKEN: "your-telegram-bot-token"       # Token for Telegram bot integration
  FRIGATE_URL: "http://localhost:5000"                # Internal URL for Frigate
  DEBUG: False                                        # Enable or disable debug mode
  TELEGRAM_CHAT_ID: 123456789                         # Telegram chat ID for notifications
  SLEEP_TIME: 5                                       # Interval between event checks, in seconds
  FRIGATE_EVENT_LIMIT: 20                             # Max number of events to retrieve per cycle
  FRIGATE_EXTERNAL_URL: "http://localhost:5000"       # External Frigate URL for notifications
  TZ: "America/New_York"                              # Set timezone for event timestamps
  REDIS_ADDR: "localhost:6379"                        # Redis server address
  REDIS_DB: 0                                         # Redis database index
  TIME_WAIT_SAVE: 30                                  # Wait time before saving event data, in seconds
```


## Global and Camera-Specific Filters

The main filtering configurations are set in the `cameras` section. This section supports both **global filters** that apply to all cameras and **camera-specific filters** for fine-grained control.

### Structure Overview

```yaml
cameras:
  # Global filters applied to all cameras unless overridden
  camera:
    - kitchen
    - frontyard
  zone:
    - door
    - yard
  label:
    - person
    - vehicle
  sublabel:
    - Max
    - Visitor
  score:
    min_score: 0.7
    max_score: 1.0
  length:
    min_length: 5.0
    max_length: 60.0
  severity:
    reviewed: false
    severity: alert
  media:
    has_snapshot: true
    has_clip: true
    include_thumbnails: true
  time:
    before: 1627804800
    after: 1627752000
    timezone: "America/New_York"
    time_range: "08:00,22:00"
  
  # Camera-specific filters
  frontyard:
    zone:
      porch:
        label:
          person:
            sublabel:
              - visitor
            score:
              min_score: 0.9
    time:
      after: 1627752000
```

---

## Filter Types

Each filter type controls a specific aspect of event selection. The following explains the available filters and their usage.

### Camera Filter (`camera`)

- **Purpose**: Specify which cameras are included in the filter.
- **Options**: List camera names to include.
- **Example**:
  ```yaml
  camera:
    - kitchen
    - frontyard
  ```

### Zone Filter (`zone`)

- **Purpose**: Include only events occurring in specified zones.
- **Options**: List of zones.
- **Example**:
  ```yaml
  zone:
    - yard
    - driveway
  ```

### Label and Sublabel Filters (`label`, `sublabel`)

- **Purpose**: Filter events by label (e.g., "person") and sublabel (e.g., "Max").
- **Usage**: Set `label` and/or `sublabel` as top-level filters for "or" behavior, or nest `sublabel` under `label` for "and" behavior.
- **Example**:
  ```yaml
  label:
    - person
  sublabel:
    - Max


### Score Filter (`score`)

- **Purpose**: Include events only within a certain score range.
- **Usage**: Define `min_score` and `max_score` values.
- **Example**:
  ```yaml
  score:
    min_score: 0.7
    max_score: 1.0
  ```

### Length Filter (`length`)

- **Purpose**: Filter events by length (duration in seconds).
- **Example**:
  ```yaml
  length:
    min_length: 5.0
    max_length: 60.0
  ```

### Severity Filter (`severity`)

- **Purpose**: Filter events based on severity (e.g., `alert`).
- **Example**:
  ```yaml
  severity:
    reviewed: false
    severity: alert
  ```

### Media Filter (`media`)

- **Purpose**: Include events based on available media types (snapshot, clip, etc.).
- **Example**:
  ```yaml
  media:
    has_snapshot: true
    has_clip: true
    include_thumbnails: true
  ```

### Time Filter (`time`)

- **Purpose**: Restrict events to a specific time range.
- **Options**: `before`, `after`, `timezone`, and `time_range`.
- **Example**:
  ```yaml
  time:
    before: 1627804800
    after: 1627752000
    timezone: "America/New_York"
    time_range: "08:00,22:00"
  ```

---

## Examples

### Basic Global Filter

This includes all events from the `kitchen` camera with a score of at least 0.7:
```yaml
cameras:
  camera:
    - kitchen
  score:
    min_score: 0.7
```

### Camera-Specific Filter with Nested Conditions

This example only includes events in the `driveway` zone of the `frontyard` camera with `label: dog` and `sublabel: max`:
```yaml
cameras:
  frontyard:
    zone:
      driveway:
        label:
          dog:
            sublabel:
              - max
```

### Mixed Global and Camera-Specific Filter

This includes any event with `label: person` or `sublabel: visitor` across all cameras, with a nested score filter specific to `frontyard`:
```yaml
cameras:
  label:
    - person
  sublabel:
    - visitor
  frontyard:
    score:
      min_score: 0.8

---

## Logic and Behavior

- **"And" Logic**: Achieved by nesting fields (e.g., `label` and `sublabel` under `zone`). All nested conditions must match.
- **"Or" Logic**: Fields at the same level (e.g., multiple `label` or `sublabel` entries) act independently, meaning any of them can match.

### Combining "And" and "Or" Logic

The configuration supports combining "and" and "or" logic through a mix of nested and non-nested structures.

---

## FAQ

### 1. **What if I only want to filter by sublabel without specifying a label?**
   - Specify the `sublabel` at the global level or directly under the camera without nesting.

### 2. **Can I override global settings for a specific camera?**
   - Yes. Any setting defined within a specific camera overrides the global setting for that camera.

### 3. **How does the time filter work?**
   - Set `before` and `after` timestamps to limit events to a specific date/time range. You can also use `time_range` to specify daily hours.

### 4. **What’s the difference between nesting and non-nesting?**
   - Nesting enforces "and" conditions, while non-nesting allows independent "or" conditions. 

---

This `README.md` provides a detailed reference for configuring and using the `config.yml` file for event filtering. You can use this guide to create complex filter configurations based on your requirements. Let me know if you need more adjustments!
