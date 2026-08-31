# fixtures/

Static test data consumed by solution code at runtime — not Go source, not compiled, not test-framework fixtures in the xUnit sense.

- `file.txt` — whitespace-column table of `color`, `date`, `number` rows. Loaded by `golang/interview/affirm_spreadsheet.go` (`NewSpreadSheet("../../fixtures/file.txt")`) as sample input for a mock-interview spreadsheet problem. The relative path is hardcoded, so this file must stay at `fixtures/file.txt` relative to the repo root.

Add new fixture files here only when a solution needs an external data file it can't reasonably inline as a Go literal.
