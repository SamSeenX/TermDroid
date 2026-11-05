# TermDroid Project Rules

## 1. 200-Line Limit
* Every Go source file must strictly stay under 200 lines of code.
* Files nearing 180 lines must be split into dedicated submodules.
* Project markdown files (BRAIN.md, RULES.md, CHANGELOG.md, README.md, SCOPE.md) are exempt.

## 2. No AI Writing Markers
* Never use semicolons in documentation or conversational messages.
* Never use em dashes or double hyphens as dashes in text.

## 3. Subprocess Safety
* Never invoke blocking system commands directly on the Bubble Tea UI thread.
* Always manage subprocess execution using Go contexts with cancellation.
* Clean up all spawned Scrcpy windows and ADB logcat streams upon exit.

## 4. UI Standards
* Follow the Elm architecture (Model, Update, View) strictly.
* Design with responsive terminal dimensions (minimum 80x24 terminal window).
* Keep styles centralized within internal/app/styles.go using Lip Gloss.
