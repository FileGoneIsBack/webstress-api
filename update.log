# Changelog

## 2025-04-18
- Reworked attack handler to centralize parameter parsing from panel/API into a single `Attack` struct
- Combined validation of all attack parameters in one place instead of individually
- Implemented attack splitting and recorded each attack in the database for tracking

## 2025-04-21
- Enhanced attack handler to log which servers and methods receive each connection
- Added retry logic: if a server or API fails, automatically retry with alternate endpoints
- Display detailed error reporting when any connection drops
- Replaced fixed thread count with a configurable slider-controlled concurrency value

## 2025-04-22
- Migrated site settings storage to SQL
- Built a dedicated site settings page with new CSS and JS
- Created backend endpoint for site settings management
- Replaced most in-code config options with database-driven settings

## 2025-04-23
- Migrated flood methods management to SQL, removing `methods.json`
- Created SQL table `floods_methods` and updated Go models
- Added methods management UI and CRUD backend
- Integrated workflow functions connecting site settings and methods modules

## Coming Soon
- Migrate `api.json` configuration to SQL-backed storage
- Extend methods functions to integrate directly with external APIs
- Add full audit logging for settings and method changes

