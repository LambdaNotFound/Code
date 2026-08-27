#!/usr/bin/env python3
"""Validate this repo's agent and skill definitions.

Two suites, no arguments, exits non-zero on any hard failure:

  structural  frontmatter parses, agent names unique, hook targets resolve
              and are executable, no dangling repo paths
  semantic    router contract (Use when / Not for), every "(use X)" target
              exists, mutual deferrals surfaced for judgement, no bundled-skill
              name collisions, resume derivations contiguous, protocols name
              real agents, no orphaned reference files, SKILL.md size

Run from anywhere:  python3 agent-team-workspace/validate-definitions.py
"""
import os, re, glob, yaml, io, sys, collections
os.chdir(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..'))
HAVE_YAML=True
fails=[]; warns=[]; checks=0
def ck(c,label,detail=""):
    global checks; checks+=1
    if not c: fails.append(f"{label}: {detail}")
def warn(c,label,detail=""):
    if not c: warns.append(f"{label}: {detail}")
def frontmatter(p):
    t=io.open(p,encoding='utf-8').read()
    m=re.match(r'^---\n(.*?)\n---\n',t,re.S)
    return yaml.safe_load(m.group(1)) if m else None
def fm_body(p):
    t=io.open(p,encoding='utf-8').read()
    m=re.match(r'^---\n(.*?)\n---\n',t,re.S)
    return yaml.safe_load(m.group(1)), t[m.end():]

print("=" * 62)
print("1. AGENTS — frontmatter parses, name present")
agents=sorted(glob.glob('.claude/agents/**/*.md', recursive=True))
names=collections.Counter()
for a in agents:
    fm=frontmatter(a)
    ok=ck(fm is not None, "no frontmatter", a)
    if ok and HAVE_YAML:
        n=fm.get('name')
        ck(bool(n), "no name field", a)
        ck(bool(fm.get('description')), "no description", a)
        if n: names[n]+=1
        # file basename should match name (repo convention, not required by loader)
        if n and os.path.basename(a)[:-3]!=n:
            print(f"   note: {a} declares name '{n}' (basename differs)")
print(f"   {len(agents)} agent files parsed")

print("2. AGENTS — no duplicate `name` (loader picks by read order if dupes)")
dupes=[n for n,c in names.items() if c>1]
ck(not dupes, "DUPLICATE agent names", str(dupes))
print(f"   {len(names)} unique names, {len(dupes)} duplicates")

print("3. AGENTS — every hook command resolves and is executable")
hooks=0
for a in agents:
    for cmd in re.findall(r'command:\s*"([^"]+)"', open(a,encoding='utf-8').read()):
        hooks+=1
        p=cmd.replace('${CLAUDE_PROJECT_DIR}','.')
        ck(os.path.isfile(p), "hook target missing", f"{a} -> {cmd}")
        ck(os.access(p, os.X_OK), "hook not executable", p)
print(f"   {hooks} hook command(s) checked")

print("3b. AGENTS — frontmatter grants carry the body rule they require")
for a in agents:
    front, body = fm_body(a)
    if front.get('maxTurns'):
        ck(re.search(r'turn budget|runs low', body, re.I),
           "maxTurns without a cutoff degradation rule (agent-factory Pass A #7)", a)
    if front.get('memory'):
        ck(re.search(r'never hold|process lessons', body, re.I),
           "memory grant without a scope rule (agent-factory Step 2)", a)
        # the rule is stated two ways in this repo: "files outrank memory"
        # and "the files win" — accept either, and any file-first phrasing
        ck(re.search(r'outrank|files win|file[^.]{0,40}wins?\b', body, re.I),
           "memory grant without a files-outrank-memory rule", a)
print(f"   {sum(1 for a in agents if fm_body(a)[0].get('maxTurns'))} capped, "
      f"{sum(1 for a in agents if fm_body(a)[0].get('memory'))} with memory")

print("4. SKILLS — <name>/SKILL.md present, frontmatter parses")
skdirs=sorted(d for d in glob.glob('.claude/skills/*') if os.path.isdir(d))
for d in skdirs:
    s=os.path.join(d,'SKILL.md')
    if ck(os.path.isfile(s), "missing SKILL.md", d):
        fm=frontmatter(s)
        ck(fm is not None, "no frontmatter", s)
        if fm and HAVE_YAML:
            ck(bool(fm.get('description')), "no description", s)
            ck(len(fm.get('description',''))<=1536, "description over 1536 cap", s)
print(f"   {len(skdirs)} skills checked")

print("5. ALL FILES — every referenced repo path resolves")
ROOTS=('agent-team-workspace/','.claude/','golang/','python/','rust/','spaced_repetition/','fixtures/','notes/','system_design/','.github/')
TOK=re.compile(r'[A-Za-z0-9_.<>*-]+(?:/[A-Za-z0-9_.<>*-]+)+')
scan=agents+[os.path.join(d,'SKILL.md') for d in skdirs] \
     +sorted(glob.glob('.claude/skills/*/references/*.md')) \
     +sorted(glob.glob('agent-team-workspace/protocols/*.md')) \
     +sorted(glob.glob('agent-team-workspace/agent-specs/*.md'))+['CLAUDE.md']
missing=set()
for f in scan:
    for line in open(f,encoding='utf-8'):
        for chunk in re.findall(r'`([^`]*)`|\]\(([^)]*)\)', line):
            s=chunk[0] or chunk[1]
            if s.startswith('http'): continue
            for t in TOK.findall(s):
                if t.startswith(ROOTS) and '<' not in t and '*' not in t and not os.path.exists(t):
                    missing.add((f,t))
for f,t in sorted(missing): ck(False, "dangling path", f"{f} -> {t}")
print(f"   {len(scan)} files scanned, {len(missing)} dangling")

skills={os.path.basename(os.path.dirname(p)):p for p in glob.glob('.claude/skills/*/SKILL.md')}
agents={}
for p in glob.glob('.claude/agents/**/*.md',recursive=True):
    agents[fm_body(p)[0]['name']]=p
BUNDLED={'code-review','simplify','security-review','init','run','loop','learn','design','dataviz',
         'docx','pdf','pptx','xlsx','morning','skill-creator','mcp-builder','doc-coauthoring',
         'import-memory','update-config','keybindings-help','claude-api','session-start-hook',
         'fewer-permission-prompts','canvas-design','artifact-design','artifact-diagramming',
         'artifact-capabilities'}
# Only genuine Claude Code built-ins belong here. This repo's own agents are
# resolved from .claude/agents/ above; listing them here too would vouch for an
# agent after it is deleted, which is exactly the dangling route this check exists
# to catch.
BUILTIN_AGENTS={'general-purpose','Explore','Plan','claude','statusline-setup','claude-code-guide'}

print("="*64)
print("6. ROUTER CONTRACT — description carries 'Use when' and 'Not for'")
for n,p in sorted(skills.items()):
    d=fm_body(p)[0].get('description','')
    warn('Use when' in d or 'Use ' in d, "no 'Use when' clause", n)
    warn('Not for' in d, "no 'Not for' boundary", n)
for n,p in sorted(agents.items()):
    d=fm_body(p)[0].get('description','')
    warn('Not for' in d or 'Use ' in d, "thin router description", n)
print(f"   {len(skills)} skills + {len(agents)} agents inspected")

print("7. NEGATIVE SPACE — every 'use X' target exists")
known=set(skills)|set(agents)|BUNDLED|BUILTIN_AGENTS
edges=[]
for n,p in sorted(list(skills.items())+list(agents.items())):
    d=fm_body(p)[0].get('description','')
    for m in re.finditer(r'\(use ([^)]+)\)', d):
        raw=m.group(1).strip()
        # '/' separates alternatives ("golang-pro/rust-pro") but also splits file
        # paths, so drop any path-with-extension before treating '/' as a separator.
        raw=re.sub(r'\S*/\S*\.\w+', ' ', raw)
        # A parenthetical may name several targets: "use pr-loop, golang-pro, or
        # rust-pro". Filtering candidates to ones already known would make this
        # check unable to fail whenever any one target resolves, so decide what
        # is a name by its shape instead: a kebab-case token always is, and a
        # bare single word is only when the fragment is nothing else. Prose like
        # "whose code-bar-raiser owns rounds" is skipped on the leading word.
        for frag in re.split(r',| or |/', raw):
            frag=frag.strip().strip('`.')
            if not frag: continue
            tok=frag.split(' ')[0]
            kebab=re.fullmatch(r'[a-z][a-z0-9]*(-[a-z0-9]+)+', tok)
            if not (kebab or frag==tok): continue
            edges.append((n,tok))
            ck(tok in known, "routes to unknown target", f"{n} -> '{tok}'")
print(f"   {len(edges)} routing edges checked")

print("8. ROUTING LOOPS — A defers to B while B defers to A")
es=set(edges)
def deferrals(n):
    """[(topic, target)] parsed from this description's 'Not for T (use X)' clauses."""
    src=skills.get(n) or agents.get(n)
    if not src: return []
    d=fm_body(src)[0].get('description','')
    out=[]
    for t,x in re.findall(r'Not for ([^(]+)\(use ([^)]+)\)', d):
        for tgt in re.split(r',| or |/', x):
            out.append((set(w for w in t.lower().split() if len(w)>4), tgt.strip().strip('`.').split(' ')[0]))
    return out
seen=set()
for a,b in sorted(es):
    if (b,a) in es and a in known and b in known and (b,a) not in seen:
        seen.add((a,b))
        # a real loop: a defers topic T to b, AND b defers an overlapping topic back to a
        ta=[t for t,x in deferrals(a) if x==b]
        tb=[t for t,x in deferrals(b) if x==a]
        clash=[sorted(x&y) for x in ta for y in tb if x&y]
        # word overlap cannot settle whether two topics are the same topic;
        # surface the pair with its topics and let a reader judge.
        flag = "REVIEW" if clash else "ok    "
        ta_s = "; ".join(" ".join(sorted(t)) for t in ta) or "-"
        tb_s = "; ".join(" ".join(sorted(t)) for t in tb) or "-"
        print(f"   {flag} {a} -> {b} on [{ta_s}] | {b} -> {a} on [{tb_s}]")
print(f"   {len(es)} unique edges")

print("9. NAME COLLISIONS — project skill shadowing a bundled skill")
for n in sorted(skills):
    ck(n not in BUNDLED, "shadows a bundled skill", n)
print(f"   {len(skills)} skill names checked against {len(BUNDLED)} bundled")

print("10. RESUME DERIVATIONS — numbering contiguous, no duplicates")
for p in ['.claude/skills/scoping/SKILL.md',
          'agent-team-workspace/protocols/design-review-loop-agent-team-prompt.md',
          'agent-team-workspace/protocols/pr-loop-agent-team-prompt.md']:
    t=io.open(p,encoding='utf-8').read()
    m=re.search(r'(first matching state|take the first matching state|Read them and take the first matching state)(.*?)\n\n[A-Z#]', t, re.S)
    if not m: warn(False,"no resume derivation found",p); continue
    nums=[int(x) for x in re.findall(r'^\s*(\d+)\.\s', m.group(2), re.M)]
    ck(len(nums)==len(set(nums)), "duplicate resume state numbers", f"{p} {nums}")
    ck(nums==list(range(nums[0], nums[0]+len(nums))), "non-contiguous resume states", f"{p} {nums}")
    print(f"   {os.path.basename(p):<45} states {nums[0]}..{nums[-1]} ({len(nums)})")

print("11. PROTOCOLS — every agent they name exists")
for p in glob.glob('agent-team-workspace/protocols/*.md'):
    t=io.open(p,encoding='utf-8').read()
    for a in set(re.findall(r'`(research-investigator|design-bar-raiser|coding-expert|code-bar-raiser|ai-writing-auditor)`', t)):
        ck(a in agents, "protocol names missing agent", f"{os.path.basename(p)} -> {a}")
print(f"   {len(glob.glob('agent-team-workspace/protocols/*.md'))} protocols checked")

print("12. REFERENCES — no orphaned reference files")
for d in glob.glob('.claude/skills/*/references'):
    sk=os.path.join(os.path.dirname(d),'SKILL.md')
    body=io.open(sk,encoding='utf-8').read()
    for f in glob.glob(os.path.join(d,'*.md')):
        rel='references/'+os.path.basename(f)
        ck(rel in body, "orphaned reference file (never loaded)", f)
print(f"   {len(glob.glob('.claude/skills/*/references/*.md'))} reference files checked")

print("13. SIZE — SKILL.md under the documented 500-line guidance")
for n,p in sorted(skills.items()):
    L=sum(1 for _ in io.open(p,encoding='utf-8'))
    warn(L<=500, "SKILL.md over 500 lines", f"{n} = {L}")
    if L>250: print(f"   note: {n} is {L} lines")


print("="*64)
print(f"{checks-len(fails)}/{checks} hard checks passed")
for f in fails: print("  FAIL", f)
print(f"{len(warns)} advisories")
for w in warns: print("  warn", w)
sys.exit(1 if fails else 0)
