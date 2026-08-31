# notes/project_rules/

Prompt files, not code. Each is a "project rules" doc pasted into a separate Claude Code project to configure a session for one kind of interview-prep practice. None of this is read by the Go module or its tests.

Three topics, each with a `_v0` and a newer `_v1`:
- `coding_interview_v{0,1}.md` — Go/Python DSA and LeetCode review rules: submit-then-review workflow, bug-first feedback, complexity checks.
- `system_design_v{0,1}.md` — mock system-design interviewer rules: phase structure (scope → high-level → deep dive → scale/failure → debrief), trade-off-reasoning focus.
- `behavioral_quesiton_v{0,1}.md` — behavioral (BQ) interview rules: building and maintaining a reusable story catalog, drilling story selection/framing per question. (Filename typo — "quesiton" — kept as-is since it's an existing checked-in name.)

The `v1` file in each pair is the current, refined version; `v0` is kept as the earlier iteration rather than deleted. When editing prep rules, prefer updating `v1` unless asked to fork a new version.
