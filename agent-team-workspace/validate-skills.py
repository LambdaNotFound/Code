import os, re, io, glob, yaml, subprocess, sys
os.chdir('/home/user/Code')
OLD2NEW = {"scoping":"scope-problem","pr-review":"review-pr","rust-expert":"write-rust",
           "agent-factory":"build-agent","design-loop":"run-design-loop","pr-loop":"run-pr-loop",
           "review-code":"review-code","summarize-changes":"summarize-changes"}
skills = sorted(os.path.basename(os.path.dirname(p)) for p in glob.glob('.claude/skills/*/SKILL.md'))
known = set(skills) | {os.path.basename(a)[:-3] for a in glob.glob('.claude/agents/**/*.md', recursive=True)} | {
    'golang-pro','rust-pro','doc-coauthoring','architect-reviewer','code-reviewer','leetcode-reviewer'}
fails=[]
def ck(c, sk, what):
    if not c: fails.append(f"{sk}: {what}")
    return c

print(f"{'SKILL':<19}{'fm':<4}{'desc':<6}{'refs':<6}{'routes':<8}{'paths':<7}{'stale':<7}{'lines':<7}{'git'}")
print("-"*72)
for sk in skills:
    p=f'.claude/skills/{sk}/SKILL.md'
    txt=io.open(p,encoding='utf-8').read()
    m=re.match(r'^---\n(.*?)\n---\n', txt, re.S)
    fm = yaml.safe_load(m.group(1)) if m else None
    body = txt[m.end():] if m else txt
    ck(fm is not None, sk, "frontmatter does not parse")
    desc = (fm or {}).get('description','')
    ck(bool(desc) and len(desc)<=1536, sk, "description missing or over cap")

    refs = re.findall(r'\]\((references/[^)]+)\)', txt)
    refs_ok = all(os.path.isfile(f'.claude/skills/{sk}/{r}') for r in refs)
    ck(refs_ok, sk, "a references/ link does not resolve")
    nref = len(glob.glob(f'.claude/skills/{sk}/references/*.md'))
    ck(nref == len(set(refs)), sk, f"reference count mismatch: {nref} files vs {len(set(refs))} links")

    routes=[]
    for mm in re.finditer(r'\(use ([^)]+)\)', desc):
        raw=re.sub(r'\S*/\S*\.\w+',' ',mm.group(1))
        for frag in re.split(r',| or |/', raw):
            frag=frag.strip().strip('`.')
            if not frag: continue
            tok=frag.split(' ')[0]
            if re.fullmatch(r'[a-z][a-z0-9]*(-[a-z0-9]+)+', tok) or frag==tok:
                routes.append(tok)
    bad=[r for r in routes if r not in known]
    ck(not bad, sk, f"routes to unknown: {bad}")

    paths=set()
    for chunk in re.findall(r'`([^`]*)`|\]\(([^)]*)\)', txt):
        s2=chunk[0] or chunk[1]
        if s2.startswith('http'): continue
        for t in re.findall(r'[A-Za-z0-9_.<>*-]+(?:/[A-Za-z0-9_.<>*-]+)+', s2):
            if t.startswith(('agent-team-workspace/','.claude/')) and '<' not in t and '*' not in t:
                paths.add(t)
    badp=[t for t in paths if not os.path.exists(t)]
    ck(not badp, sk, f"dangling path: {badp}")

    # only a skill-reference form is stale; "scoping" is also an ordinary gerund,
    # so a bare-word match flags legitimate prose
    probe = txt.replace('pr-loop-agent-team-prompt','').replace('run-pr-loop','').replace('run-design-loop','')
    stale=[o for o in OLD2NEW if o!=OLD2NEW[o]
           and re.search(rf'(?<!run-)(?<!/)(?<![a-z-])(/{o}\b|`{o}`|skills/{o}\b|\(use {o}\b)', probe)]
    ck(not stale, sk, f"stale old skill name: {stale}")

    # content integrity: body must not have shrunk vs the pre-rename version in git
    old = [k for k,v in OLD2NEW.items() if v==sk][0]
    prev = subprocess.run(['git','show',f'HEAD:.claude/skills/{old}/SKILL.md'],
                          capture_output=True, text=True)
    gitmark = "?"
    if prev.returncode==0:
        d = len(txt.splitlines()) - len(prev.stdout.splitlines())
        gitmark = "same" if d==0 else f"{d:+d}"
        ck(abs(d)<=3, sk, f"line count moved {d} vs pre-rename — content may have been lost")
    print(f"{sk:<19}{'ok':<4}{len(desc):<6}{len(set(refs)):<6}{len(routes):<8}{len(paths):<7}"
          f"{'none':<7}{len(txt.splitlines()):<7}{gitmark}")

print("-"*72)
print(f"{len(skills)} skills checked, {len(fails)} failures")
for f in fails: print("  FAIL", f)
sys.exit(1 if fails else 0)
