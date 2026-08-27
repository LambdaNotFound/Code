---
description: Summarize uncommitted working-tree changes in a few bullets and flag risks such as missing error handling, hardcoded values, or tests needing updates. Use when the user asks what changed, wants a commit message, or wants a quick read on their uncommitted diff. Not for a defect-hunting review of that diff (use review-code, or pr-review for a PR). Not for committing or pushing.
---

## Current changes

!`git diff HEAD`

## Instructions

Summarize the changes above in two or three bullet points, then list any risks you notice such as missing error handling, hardcoded values, or tests that need updating. If the diff is empty, say there are no uncommitted changes.