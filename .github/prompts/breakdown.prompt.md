---
name: breakdown
description: Break down increment implementation into actionable tasks
argument-hint: increment name
---

# Rule: Generating a Task List

This document represents the **HOW - Detailed** - breaking down implementation into actionable, step-by-step tasks that can be followed by developers or AI agents.

## Goal

To guide an AI assistant in creating a **lean, incremental task list** in Markdown format that breaks down a small feature into actionable implementation steps, just as a developer would plan their work. Tasks should be minimal, testable, and deliver value incrementally.

## Prerequisites

Before generating tasks, verify that all prerequisite files exist:

**Required:**
1. **`CONSTITUTION.md`** - Defines technical decisions, frameworks, and development principles
2. **`[increment-name]/increment.md`** - Defines WHAT to build and acceptance criteria
3. **`[increment-name]/design.md`** - Suggests initial approach (technical starting point)

**If any file does not exist:**
- Inform the user which prerequisite is missing
- Do not proceed until all three files exist

**If all files exist:**
- Read `CONSTITUTION.md` to understand the broad technical foundation, testing philosophy, and deployment constraints
- Read `[increment-name]/increment.md` to understand what to build (acceptance criteria)
- Read `[increment-name]/design.md` to understand initial approach (starting point that complies with constitution)
- Generate tasks that implement the feature following the design's initial approach, staying flexible to adjust as needed

## Task Structure Guidelines

**Keep tasks minimal and incremental:**
- Each task should be completable in **15-30 minutes**
- Each task should deliver **testable progress**
- Tasks should map to **acceptance criteria** from the feature
- Tasks should follow **initial approach** from the design (but stay flexible)
- Focus on **behavior** (what works), not files (what changes)
- Prioritize "make it work" over "make it perfect"

**Task organization:**
- Start with setup (branch, dependencies)
- Build incrementally (component → integration → verification)
- Each step should be verifiable
- End with deployment readiness

**Definition of Done for each task:**
- Code written
- Locally verified (manual or automated test)
- Meets acceptance criteria from feature

## Process

1.  **Verify Prerequisites:** Check for `CONSTITUTION.md`, `[increment-name]/increment.md`, and `[increment-name]/design.md`
2.  **Read Context:**
    - **Read `CONSTITUTION.md`** to understand testing philosophy, deployment approach, and technical constraints
    - **Read `[increment-name]/increment.md`** to understand acceptance criteria (Given/When/Then scenarios)
    - **Read `[increment-name]/design.md`** to understand initial approach, component boundaries, and data flow
3.  **Generate Tasks:** Create minimal, incremental tasks that follow the design's initial approach to satisfy the increment's acceptance criteria ("make it work" first)
4.  **Save Task List:** Save as `[increment-name]/breakdown.md`

## Before Generating Tasks - Self Check

Ask yourself:
- [ ] Did I verify all three prerequisite files exist?
- [ ] Did I read CONSTITUTION.md for principles and constraints?
- [ ] Did I read increment.md for acceptance criteria?
- [ ] Did I read design.md for technical approach?

If any checkbox is unchecked, STOP and complete that step first.

Note: This prompt typically doesn't require clarifying questions since all context should be in the three documents. Generate tasks directly after verification.

**Note:** Tasks are generated in one pass—no two-phase approval needed. Keep them minimal and focused.

## Output Format

The generated task list should follow this structure:

```markdown
# Tasks: [Feature Name]

## Relevant Files

List files that will be created or modified for this feature. This scopes the work and prevents unintended changes.

- `path/to/component.tsx` - Brief description (e.g., Upload form component)
- `path/to/component.test.tsx` - Unit tests for component
- `api/profile/picture.ts` - Brief description (e.g., Upload endpoint)
- `api/profile/picture.test.ts` - Unit tests for endpoint
- `lib/utils/validation.ts` - Brief description (e.g., File validation helpers)

**Note:** Only list files directly related to this increment. Don't include framework files or unrelated code.

## Progress Tracking

**As you complete each task, check it off by changing `- [ ]` to `- [x]`.**

Update after completing each task to track progress.

## Implementation Tasks

- [ ] **Setup**
  - [ ] Create feature branch: `git checkout -b feature/[increment-name]`
  - [ ] [Any dependencies or configuration needed]

- [ ] **[Behavior/Capability 1]** (maps to acceptance criteria)
  - [ ] [Specific implementation step]
  - [ ] [Verification step]

- [ ] **[Behavior/Capability 2]** (maps to acceptance criteria)
  - [ ] [Specific implementation step]
  - [ ] [Verification step]

- [ ] **Integration & Verification**
  - [ ] [Integration step if needed]
  - [ ] Verify all acceptance criteria pass
  - [ ] [Manual testing step if per constitution]

- [ ] **Quick Review**
  - [ ] Verify all acceptance criteria pass
  - [ ] Check alignment with design approach
  - [ ] Remove any dead code or debug statements

- [ ] **Deploy**
  - [ ] Commit changes with clear message
  - [ ] [Push/merge per constitution deployment strategy]
```

## Example Output

```markdown
# Tasks: Upload Profile Picture

## Relevant Files

- `components/profile/UploadPhoto.tsx` - Upload form component with file input and preview
- `components/profile/UploadPhoto.test.tsx` - Unit tests for upload component
- `pages/api/profile/picture.ts` - POST endpoint for image upload
- `pages/api/profile/picture.test.ts` - Unit tests for upload endpoint
- `lib/storage/s3.ts` - S3 upload helper functions
- `components/profile/ProfileImage.tsx` - Component to display profile image

## Progress Tracking

**As you complete each task, check it off by changing `- [ ]` to `- [x]`.**

🛑 **STOP - BEFORE GENERATING TASK LIST**

Self-check before proceeding:
- [ ] I will generate ALL tasks as UNCHECKED `- [ ]` 
- [ ] I will NOT check tasks off until I actually complete them
- [ ] I understand this file tracks future work, not completed work
- [ ] Tasks follow the structure: Setup → Implementation → Integration → Review → Deploy

If any checkbox is unchecked, DO NOT PROCEED. These tasks represent work to be done, not work already completed.

## Implementation Tasks

- [ ] **Setup**
  - [ ] Create feature branch: `git checkout -b feature/upload-profile-picture`
  - [ ] Install S3 client library if needed (check ADR)

- [ ] **File upload component**
  - [ ] Create upload form with file input (JPG/PNG only)
  - [ ] Add client-side validation: file type and size < 5MB
  - [ ] Display selected file preview
  - [ ] Verify: Can select file, see preview, validation errors show

- [ ] **Upload API endpoint**
  - [ ] Create POST /api/profile/picture endpoint
  - [ ] Validate file type and size server-side
  - [ ] Upload to S3 with user ID as key
  - [ ] Return public image URL
  - [ ] Verify: curl upload works, returns URL, rejects invalid files

- [ ] **Display uploaded image**
  - [ ] Update profile UI to show image from URL
  - [ ] Handle missing image state (placeholder)
  - [ ] Verify: Image displays after upload, persists on refresh

- [ ] **Error handling**
  - [ ] Show error message on upload failure
  - [ ] Allow retry on error
  - [ ] Verify: Error scenarios display correctly

- [ ] **Integration & Verification**
  - [ ] Test complete flow: select → upload → display
  - [ ] Verify acceptance criteria: under 5MB works, over 5MB fails, image persists
  - [ ] Manual test: Upload actual image via UI

- [ ] **Quick Review**
  - [ ] Verify all acceptance criteria pass
  - [ ] Check alignment with design approach
  - [ ] Remove any dead code or debug statements

- [ ] **Deploy**
  - [ ] Commit: "Add profile picture upload feature"
  - [ ] Push to branch, create PR (per constitution CI/CD)
```

## Output

- **Format:** Markdown (`.md`)
- **Location:** `[increment-name]/` directory
- **Filename:** `breakdown.md`

## Final Instructions

1. **Verify all prerequisites exist** (`CONSTITUTION.md`, `increment.md`, `design.md`)
2. **Read all three documents** to understand principles, requirements, and initial approach
3. **Generate behavior-focused tasks** that map to acceptance criteria
4. **Keep tasks small** (15-30 min each, testable progress)
5. **Include verification steps** for each capability
6. **Follow constitutional principles** for testing and deployment
7. **Prioritize "make it work"** - get acceptance criteria passing first
8. **No explicit references** to constitution/feature/design in the generated tasks—just implement them
9. **Focus on what works**, not what files change
