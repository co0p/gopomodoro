---
name: increment
description: Create an increment specification with Gherkin acceptance criteria
argument-hint: increment name or user capability
---

# Rule: Generating an Increment Specification

This document represents the **WHAT** - describing individual increments from the user's perspective without technical implementation details.

## Goal

To guide an AI assistant in creating a lightweight, focused increment specification in Markdown format, just as a product owner would define user value. Each increment is a **small, testable addition** (subfunction) that validates a specific assumption through Gherkin-style acceptance criteria.

## Prerequisites

Before generating a feature, verify that a `CONSTITUTION.md` file exists in the project root. The constitution defines the project's core principles, technical decisions, and development philosophy.

**If `CONSTITUTION.md` does not exist:**
- Inform the user that a constitution is required first
- Suggest using `/create-constitution` to generate one
- Do not proceed with feature generation until the constitution exists

**If `CONSTITUTION.md` exists:**
- Read and understand the project's core principles and technical foundation
- Ensure the feature respects constitutional constraints (tech stack, principles, limitations)
- Use the constitution to guide what's feasible and align with project values

## Process

1.  **Verify Prerequisites:** Check for `CONSTITUTION.md` and read it if present.
2.  **Receive Initial Prompt:** The user provides a brief description of a feature or capability.
3.  **STOP - Ask Clarifying Questions:**
    
    **DO NOT GENERATE THE FEATURE YET.**
    
    Ask only 2-3 essential questions. Provide options in letter/number lists for easy selection. Wait for user answers.
    
4.  **Generate Feature:** (Only after receiving user's answers) Create a small, testable increment with Gherkin-style acceptance criteria.
5.  **Save Feature:** Save as `[increment-name]/increment.md`.

## Before Generating Feature - Self Check

Ask yourself:
- [ ] Did I verify CONSTITUTION.md exists and read it?
- [ ] Did I receive the user's feature description?
- [ ] Did I ask 2-3 clarifying questions?
- [ ] Did I receive user's answers (e.g., "1A, 2C, 3X")?

If any checkbox is unchecked, STOP and complete that step first.

## Clarifying Questions (Guidelines)

Ask only 2-3 critical questions to write a clear, testable feature.

### Focus Areas:
- What specific capability/action?
- What assumption are we testing?
- What does success look like?

### Formatting Requirements

- **Number all questions** (1, 2, 3, etc.)
- **List options as A, B, C, D, etc.**
- **Always include option X: "Skip this question / I don't know yet"**
- Keep it simple for responses like "1A, 2C, 3X"

### Example Questions

```
1. What specific capability are we building?
   A. Upload/select something
   B. Display/show something
   C. Process/transform something
   D. Save/persist something
   X. Skip this question / I don't know yet

2. What assumption are we testing?
   A. Users want this capability
   B. This approach is technically feasible
   C. This will improve a specific metric
   X. Skip this question / I don't know yet

3. What does success look like?
   A. User completes the action successfully
   B. Specific metric improves
   C. Feature works within time/size constraints
   X. Skip this question / I don't know yet
```

## Increment Structure

Use this structure for all increments—small, testable additions with Gherkin-style acceptance criteria.

```markdown
# [Increment Title]

## Job Story
**When** [situation]  
**I want to** [action]  
**So I can** [outcome]

**Assumption Being Tested:** [Specific hypothesis for this increment]

## Acceptance Criteria

- **Given** [precondition]  
  **When** [action]  
  **Then** [observable outcome]

- **Given** [precondition]  
  **When** [action]  
  **Then** [observable outcome]

- **Given** [error condition]  
  **When** [action]  
  **Then** [error handling outcome]

## Success Signal
[How we know this increment works - metric or observation]

## Out of Scope
- [What this increment does NOT include]
```

---

## Template Example

```markdown
# Upload Profile Picture

## Job Story
**When** I'm setting up my profile for the first time  
**I want to** upload a recognizable photo  
**So I can** build trust with other users before my first consultation

**Assumption Being Tested:** Users will complete profile setup if photo upload takes less than 30 seconds.

## Acceptance Criteria

- **Given** I am on my profile page  
  **When** I click "Upload Photo" and select a valid JPG/PNG file under 5MB  
  **Then** I see my new profile picture displayed immediately

- **Given** I am on my profile page  
  **When** I select a file over 5MB  
  **Then** I see an error message "File too large. Maximum 5MB." and can try again

- **Given** I have uploaded a photo  
  **When** I refresh the page  
  **Then** my profile picture persists (saved to my account)

## Success Signal
- 90% of users who start upload complete it successfully
- Average upload time < 15 seconds

## Out of Scope
- Photo editing/cropping (users edit before upload)
- Multiple photo uploads (one profile picture only)
- Camera integration (file selection only)
```

## Important: Separation of Concerns

**Features capture the WHAT, not the HOW:**
- **Features** define small, testable increments with Gherkin-style acceptance criteria
- **Architecture Decision Records (ADRs)** capture technical implementation decisions
- **Tasks** break down implementation into actionable steps

Keep features focused on user behavior and observable outcomes. Technical implementation details belong in ADRs and tasks.

**Example structure:**
```
/upload-profile-picture/
  increment.md           # WHAT: User capability (Gherkin)
  design.md            # HOW (Initial): Starting approach
  breakdown.md         # HOW (Detailed): Implementation steps
  improve.md           # Make it good: Clean up after it works
```

## Target Audience

Features are read by developers and testers who need clear, testable acceptance criteria for incremental work. They should be understandable without technical jargon, focusing on observable user behavior and outcomes.

## Output

*   **Format:** Markdown (`.md`)
*   **Location:** `/tasks/[increment-name]/` or appropriate project directory
*   **Filename:** `increment.md`

## Final Instructions

1. **Verify `CONSTITUTION.md` exists** before proceeding
2. **Read the constitution** to understand project values and constraints
3. **Ask 2-3 clarifying questions** only when essential
4. **Use Gherkin-style acceptance criteria** (Given/When/Then)
5. **Keep it minimal** - one screen max, small testable increment
6. **Focus on observable behavior** - no technical implementation details
7. **Make it testable** - clear success criteria that can be verified
8. **One increment, one assumption** - tight focus on single capability
