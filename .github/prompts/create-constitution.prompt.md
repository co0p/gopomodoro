---
name: create-constitution
description: Generate a project constitution defining core development principles
argument-hint: optional project type or tech stack
---

# Rule: Generating a Project Constitution

This document represents the **WHY** - the strategic foundation that defines core development principles and non-negotiable standards for the project.

## Goal

To guide an AI assistant in creating a lightweight project constitution that defines the core development principles and non-negotiable standards for the project, just as a CTO or architect would establish technical direction. The constitution serves as the foundational reference for how the team builds software and makes technical decisions.

## Process

1.  **Receive Initial Prompt:** The user requests to create or update the project constitution.
2.  **Analyze Project Context:** Review existing project files (especially `README.md` and existing `CONSTITUTION.md`) to understand the technical landscape and infer appropriate principles.
3.  **Suggest Principles:** Based on the project context, propose 3-5 core principles with rationale.
4.  **STOP - Ask Clarifying Questions:**
    
    **DO NOT GENERATE THE CONSTITUTION YET.**
    
    Ask clarifying questions only if critical information is missing or if the suggested principles need refinement. Wait for user answers.
    
5.  **Generate Constitution:** (Only after receiving answers or user confirmation to proceed)
6.  **Save Constitution:** Save the document as `CONSTITUTION.md` in the project root.

## Before Generating Constitution - Self Check

Ask yourself:
- [ ] Did I analyze project context (README.md, existing files)?
- [ ] Did I suggest principles with rationale?
- [ ] Did I ask clarifying questions if needed?
- [ ] Did I receive user's answers or confirmation to proceed?
- [ ] Will the constitution cover at least 3 different pillars?
- [ ] Does each principle map to a specific pillar?
- [ ] Are principles declarative, testable, and specific?

If any checkbox is unchecked, STOP and complete that step first.

## Context Analysis

Before generating the constitution, the AI should analyze:

*   **Technical Foundation:** If `README.md` exists, extract key values and standards if mentioned
*   **Project Type:** Web app, library, CLI tool, etc. (inferred from package.json, README, or file structure)
*   **Tech Stack:** Identify specific technologies in use:
    *   **Languages:** JavaScript, TypeScript, Python, Java, Go, Rust, etc.
    *   **Frontend Framework:** React, Vue, Angular, Svelte, Next.js, etc.
    *   **Backend Framework:** Express, NestJS, FastAPI, Django, Spring Boot, etc.
    *   **Database:** PostgreSQL, MySQL, MongoDB, Redis, SQLite, etc.
    *   **Build Tools:** Vite, webpack, esbuild, Turbopack, etc.
    *   **Testing Framework:** Jest, Vitest, Pytest, JUnit, Playwright, Cypress, etc.
    *   **Infrastructure:** Docker, Kubernetes, AWS, Vercel, Railway, etc.
*   **Existing Patterns:** Coding conventions, testing approaches, documentation practices
*   **Team Size/Maturity:** Solo project vs. team project (inferred from contributors, docs complexity)

## Suggesting Principles - The 6 Pillars

A strong constitution covers the **6 core pillars** of modern software engineering. These guide AI decision-making across architecture, implementation, and trade-offs.

**Before suggesting principles, ask the user:**

```
How comprehensive should your constitution be?
A. Cover all 6 pillars (comprehensive guidance for all aspects)
B. Focus on 3-5 key principles (lighter, focused on project priorities)
X. You decide based on project context
```

- **Option A (All 6 Pillars):** Include one principle from each pillar for complete coverage
- **Option B (3-5 Principles):** Suggest the most relevant principles based on project context, covering at least 3 different pillars
- **Option X (AI Decision):** Analyze project context and choose the appropriate level of detail

### The 6 Pillars

**1. Delivery Velocity**
How fast to ship vs. how polished? Iteration philosophy, MVP definition, acceptable quality thresholds.
- Examples: "Ship fast and iterate" vs. "Plan thoroughly before building"
- Guides: Feature scope, when to refactor, release cadence

**2. Test Strategy**
What to test, when to test, how much coverage is enough?
- Examples: "Critical paths only" vs. "Comprehensive coverage" vs. "TDD always"
- Guides: Test writing, refactoring confidence, deployment decisions

**3. Design Integrity**
How to structure code? Dependency rules, SOLID principles, architectural boundaries.
- Examples: "Dependencies point inward" vs. "Pragmatic coupling OK" vs. "Strict layering"
- Guides: Where to put logic, when to create abstractions, module boundaries

**4. Simplicity First**
When to add abstraction? YAGNI application, refactoring triggers, complexity tolerance.
- Examples: "Build only what's needed" vs. "Design for future" vs. "Three strikes then refactor"
- Guides: Premature optimization, abstraction timing, code evolution

**5. Technical Debt Boundaries**
When are shortcuts acceptable? How to track and pay down debt?
- Examples: "Quick hacks in prototypes OK" vs. "Never compromise quality" vs. "Debt must be labeled"
- Guides: Shortcut decisions, refactoring priority, quality bar

**6. Dependency Discipline**
When to add libraries? How to isolate third-party code? Framework philosophy.
- Examples: "Minimize dependencies ruthlessly" vs. "Use best tools" vs. "Wrap all external APIs"
- Guides: Library selection, vendor coupling, upgrade strategy

### Principle Selection Guidelines

**Good principles are:**
*   **Declarative:** State what IS, not what we hope to be
*   **Testable:** Can verify if we're following it
*   **Specific:** Concrete enough to guide decisions
*   **Limiting:** Rules out certain approaches or patterns
*   **Justifiable:** Has a clear "why"

**Examples of strong principles:**
*   ✅ "We write tests before fixing bugs" (specific, testable, **Test Strategy**)
*   ✅ "Dependencies point inward—business logic never imports from frameworks" (limiting, **Design Integrity**)
*   ✅ "No abstraction until third occurrence" (guides timing, **Simplicity First**)

**Examples of weak principles:**
*   ❌ "We value quality" (too vague, no pillar)
*   ❌ "We try to write good code" (not limiting)
*   ❌ "We might add tests later" (not declarative)


### Technology Questions

When asking about technology, be specific about the stack. **Always include option X for each question to allow skipping.**

```
1. What is your primary programming language?
   A. TypeScript/JavaScript
   B. Python
   C. Java/Kotlin
   D. Go
   E. Rust
   F. Other: ___________
   X. Skip this question / I don't know yet

2. What frontend framework are you using?
   A. React
   B. Vue
   C. Angular
   D. Svelte
   E. Next.js
   F. None (vanilla JS/HTML)
   G. Other: ___________
   X. Skip this question / I don't know yet

3. What backend framework are you using?
   A. Express
   B. NestJS
   C. FastAPI
   D. Django
   E. Spring Boot
   F. None (serverless/static)
   G. Other: ___________
   X. Skip this question / I don't know yet

4. What database(s) are you using?
   A. PostgreSQL
   B. MySQL
   C. MongoDB
   D. SQLite
   E. Redis
   F. in-memory
   G. None yet
   X. Skip this question / I don't know yet

5. What testing framework are you using?
   A. Jest
   B. Vitest
   C. Pytest
   D. JUnit
   E. Playwright/Cypress
   F. None yet
   G. Other: ___________
   X. Skip this question / I don't know yet
```

### Principle & Philosophy Questions

### Formatting Requirements

- **Number all questions** (1, 2, 3, etc.)
- **List options for each question as A, B, C, D, etc.** for easy reference
- **Always include option X: "Skip this question / I don't know yet"** for questions the user cannot answer
- Keep questions focused on missing critical context only
- Allow the user to select multiple answers, AB for A and B for example

### Example Questions

Ask questions that map to the **6 pillars**, focusing on what's not obvious from the codebase:

```
1. What is your delivery velocity philosophy? (Pillar: Delivery Velocity)
   A. Ship fast and iterate (MVP → feedback → improve)
   B. Plan thoroughly, then build (get it right first time)
   C. Balanced (fast for experiments, careful for core)
   X. Skip this question / I don't know yet

2. What's your testing philosophy? (Pillar: Test Strategy)
   A. TDD always (write tests first)
   B. Test critical paths only (pragmatic coverage)
   C. Comprehensive coverage (high confidence)
   D. Minimal testing (move fast, fix when broken)
   X. Skip this question / I don't know yet

3. What are your dependency rules? (Pillar: Design Integrity)
   A. Strict layering (dependencies point inward)
   B. Pragmatic coupling (optimize for simplicity)
   C. Not defined yet
   X. Skip this question / I don't know yet

4. When do you add abstraction? (Pillar: Simplicity First)
   A. Three strikes rule (third occurrence)
   B. Plan ahead (design for future needs)
   C. Never abstract (duplication is fine)
   D. When it hurts (refactor reactively)
   X. Skip this question / I don't know yet

5. How do you handle technical debt? (Pillar: Technical Debt Boundaries)
   A. Quick hacks allowed, label and schedule cleanup
   B. Never take shortcuts (quality always)
   C. Fix when it blocks progress
   X. Skip this question / I don't know yet

6. What's your dependency philosophy? (Pillar: Dependency Discipline)
   A. Minimize ruthlessly (prefer standard library)
   B. Use best tools available (pragmatic)
   C. Wrap third-party code (isolation)
   X. Skip this question / I don't know yet
```

## Constitution Structure

The generated constitution should include the following sections:

### 1. Project Name & Purpose
Brief statement of what the project is and who it serves.

### 2. Core Principles (3-5 principles)

**Cover at least 3 of the 6 pillars** based on what matters most for this project:
- **Delivery Velocity** - Speed vs. quality trade-offs
- **Test Strategy** - What and how to test
- **Design Integrity** - SOLID, dependency rules, boundaries
- **Simplicity First** - YAGNI, abstraction timing, refactoring
- **Technical Debt Boundaries** - Acceptable shortcuts, paydown strategy
- **Dependency Discipline** - Library philosophy, isolation

Each principle should include:
*   **Principle Name:** Short, memorable title
*   **Statement:** Clear declaration of the rule/standard
*   **Rationale:** Why this principle exists (1-2 sentences)
*   **In Practice:** Concrete example or implication

**Template for each principle:**
```markdown
#### [Principle Number]. [Principle Name]

**Statement:** [Clear, declarative statement of the principle]

**Rationale:** [Why this matters - the business/technical reason]

**In Practice:**
- [Concrete example or implication 1]
- [Concrete example or implication 2]
- [Concrete example or implication 3]
```


## Output Format

```markdown
# Project Constitution

## About This Project

[1-2 sentences: what the project does and who it serves]

---

## Core Principles

### 1. [Principle Name] _(Pillar: [Pillar Name])_

**Statement:** [Declarative statement]

**Rationale:** [Why this exists]

**In Practice:**
- [Implication 1]
- [Implication 2]

### 2. [Principle Name] _(Pillar: [Pillar Name])_

[Same structure...]

### 3. [Principle Name] _(Pillar: [Pillar Name])_

[Same structure...]

[4-5 if needed...]

### Pillar Coverage

_This constitution addresses the following pillars:_
- ✓ [Pillar Name] (Principle #)
- ✓ [Pillar Name] (Principle #)
- ✓ [Pillar Name] (Principle #)
- [Additional pillars if covered]

## Technical Decicions

### Languages
- [Declarative statement]: [Why we chose this]
- [Declarative statement]: [Why we chose this]


### Frameworks
- [Declarative statement]: [Why we chose this]
- [Declarative statement]: [Why we chose this]

### Deployment
- [Declarative statement]: [Why we chose this]
- [Declarative statement]: [Why we chose this]


---


**Last Updated:** [Current Date in YYYY-MM-DD format]
```

## Example Output

See `/examples/constitution-example.md` for a complete sample constitution.

## Target Audience

The constitution is read by:
*   **AI Agents** making technical decisions during increment implementation
*   **Developers** making daily coding decisions
*   **Code reviewers** evaluating if PRs align with project values
*   **New contributors** understanding project philosophy
*   **Project leads** making architectural decisions

Keep language clear and avoid jargon. Principles should be memorable and easy to reference.

## Final Instructions

1. **Analyze first, ask second:** Review project context before asking questions
2. **Suggest, don't assume:** Propose principles based on analysis, but allow user to refine
3. **Keep it lightweight:** 3-5 principles maximum, each principle fits on screen
4. **Make it actionable:** Every principle should guide real decisions
5. **Validate pillar coverage:** Ensure at least 3 different pillars are addressed
6. **Label each principle:** Explicitly indicate which pillar each principle addresses
7. **Include pillar coverage summary:** Add checklist showing which pillars are covered
8. **Save as `CONSTITUTION.md`** in project root
9. **Do NOT implement changes** based on the constitution - only create the document

🛑 **STOP - BEFORE SAVING CONSTITUTION**

Final validation:
- [ ] Constitution includes 3-5 principles
- [ ] At least 3 different pillars are covered
- [ ] Each principle clearly labeled with its pillar
- [ ] Pillar coverage summary included
- [ ] All principles are declarative, testable, and specific
- [ ] Technical decisions section includes tech stack choices

If any checkbox is unchecked, DO NOT SAVE. Fix the issues first.
