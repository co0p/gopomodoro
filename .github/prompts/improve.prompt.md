---
name: improve
description: Create a refactor plan to improve code after increment works
argument-hint: increment name
---

# Rule: Generating a Refactor Plan

This document represents **Make It Good** - improving code quality, extracting patterns, and cleaning up after getting the feature working.

## Goal

To guide an AI assistant in creating a focused refactor plan that improves code quality after a feature is working, just as a craftsperson would refine their work. This is the "make it good" phase of Kent Beck's "make it work, make it good, make it fast" cycle.

## Prerequisites

Before generating a refactor plan, verify that the feature is working:

**Required:**
1. **`CONSTITUTION.md`** - Defines technical decisions, frameworks, and development principles
2. **`[increment-name]/increment.md`** - Defines acceptance criteria
3. **`[increment-name]/breakdown.md`** - Implementation tasks (should be completed)
4. **Working implementation** - All acceptance criteria passing

**If prerequisites don't exist or feature isn't working:**
- Inform the user to complete implementation first
- Refactoring comes AFTER getting it working
- Do not proceed until acceptance criteria pass

**If feature is working:**
- Review the **entire codebase** - not just files modified in this feature
- Scan for patterns, duplication, and inconsistencies across all project files
- Identify improvement opportunities both in new code and existing code
- Look for opportunities to extract shared utilities or consolidate patterns
## Process

1.  **Verify Prerequisites:** Check that feature works and all files exist
2.  **Scan Entire Codebase:** Review all project files, not just feature files
3.  **Identify Patterns:** Look for duplication, inconsistencies, and opportunities across the whole project
4.  **Review Implementation:** Examine actual code, identify improvements
5.  **Check Constitution:** Ensure refactorings align with project principles
6.  **Generate Refactor Plan:** Create focused improvement tasks
## Before Generating Refactor Plan - Self Check

Ask yourself:
- [ ] Did I verify the feature is working (acceptance criteria pass)?
- [ ] Did I scan the **entire codebase** for patterns and duplication?
- [ ] Did I review all project files, not just the feature implementation?
- [ ] Did I look for opportunities to extract shared utilities?
- [ ] Did I check CONSTITUTION.md for alignment?
- [ ] Did I identify specific code smells or improvement opportunities?

If any checkbox is unchecked, STOP and complete that step first.

Note: This prompt typically doesn't require user questions since you're analyzing existing code. Generate the refactor plan directly after comprehensive code review.

## What to Look For

When reviewing the **entire codebase**, identify:

### Code Smells (Whole Project Scope)
- **Duplication:** Repeated logic across **any files** in the project
- **Inconsistent patterns:** Same problem solved differently in different places
- **Long functions:** Functions doing too much
- **Unclear names:** Variables/functions that don't reveal intent
- **Magic values:** Hard-coded numbers or strings scattered across files
- **Deep nesting:** Too many levels of if/for/while

### Design Improvements (Cross-File)
- **Missing abstractions:** Patterns that appear in multiple files but aren't extracted
- **Shared utilities opportunity:** Common logic that could be extracted to shared module
- **Tight coupling:** Components that know too much about each other
- **Violated principles:** Code that contradicts constitutional guidelines
- **Unclear boundaries:** Responsibilities mixed between components
- **Inconsistent conventions:** Different naming or structure patterns across files
- **Undocumented patterns:** Design decisions made during implementation that should be codified for future features

### Test Quality
- **Missing edge cases:** Scenarios not covered
- **Brittle tests:** Tests that break on minor changes
- **Duplicate test setup:** Repeated arrange code
- **Unclear assertions:** Tests that don't clearly show intent

### Documentation Gaps
- **Missing comments:** Complex logic without explanation
- **Outdated docs:** README or design docs that don't match reality
### Documentation Gaps
- **Missing comments:** Complex logic without explanation
## Refactor Plan Structure

Keep refactor plans **small and focused**. Each refactor should:
- Be completable in **15-30 minutes**
- Maintain all passing acceptance criteria
- Make one improvement at a time
- Be independently committable

```markdown
# Refactor: [Feature Name]

## Current State
[Brief description of what works and what could be better]

## Refactoring Goals
- [What quality improvements we're targeting]
- [What patterns or principles we're applying]

## Project-Level Design Decisions (If Any)
[Only include this section if you discovered patterns/decisions that should guide future features]

### Suggested ADR: [Decision Name]
**Pattern Discovered:** [What pattern or decision emerged during implementation]  
**Project Impact:** [Why this should hold true for the entire codebase]  
**Proposed Location:** `design/[number]-[decision-name].md`  
**To Be Referenced:** Future features should check this alongside CONSTITUTION.md

[Repeat for each project-level decision discovered]

## Refactor Tasksssing acceptance criteria
- Make one improvement at a time
- Be independently committable

```markdown
## Success Criteria
- [ ] All acceptance criteria still pass
- [ ] Code follows constitutional principles
- [ ] [Specific quality improvement achieved]
- [ ] No new functionality added
- [ ] Project-level patterns documented (if discovered)

## Completion
- [ ] Commit changes: "Refactor: [brief description of improvements]"
- [ ] All refactorings applied and tested
```Refactoring Goals
- [What quality improvements we're targeting]
- [What patterns or principles we're applying]

## Refactor Tasks

### 1. [Improvement Name]
**Why:** [What's wrong/unclear in current code]  
**What:** [Specific refactoring to apply]  
## Refactoring Goals
- Extract validation into shared module (DRY)
- Separate file handling from business logic (Single Responsibility)
- Centralize error messages (maintainability)
- Improve test coverage for edge cases

## Project-Level Design Decisions
During implementation, we established a pattern that should guide future file upload features.

### Suggested ADR: Client-Server Validation Parity
**Pattern Discovered:** We implemented identical validation on both client (immediate feedback) and server (security), using shared validation functions  
**Project Impact:** All future upload features should follow this pattern - validation logic must be shared, not duplicated  
**Proposed Location:** `design/001-shared-validation-pattern.md`  
**To Be Referenced:** When implementing any future upload or form validation features

## Refactor Tasksh]  
**Verify:** [Confirmation step]

[3-7 refactorings max]

## Success Criteria
- [ ] All acceptance criteria still pass
- [ ] Code follows constitutional principles
- [ ] [Specific quality improvement achieved]
- [ ] No new functionality added

## Completion
- [ ] Commit changes: "Refactor: [brief description of improvements]"
- [ ] All refactorings applied and tested
```

## Template Example

```markdown
# Refactor: Upload Profile Picture

## Current State
Feature works - users can upload profile pictures successfully. However, validation logic is duplicated between client and server, file handling code is mixed with business logic, and error messages are hard-coded strings.

## Refactoring Goals
- Extract validation into shared module (DRY)
- Separate file handling from business logic (Single Responsibility)
- Centralize error messages (maintainability)
- Improve test coverage for edge cases

## Refactor Tasks

### 1. Extract File Validation
**Why:** Validation logic duplicated in `UploadPhoto.tsx` and `picture.ts` (10 lines each)  
**What:** Extract to `lib/validation/fileValidation.ts`  
**How:** Create `validateImageFile(file)` function, import in both places  
**Verify:** Run existing tests, both client and server validation still work

### 2. Separate Storage Concerns
**Why:** API handler mixes S3 upload with request handling logic (35-line function)  
**What:** Extract S3 operations to `lib/storage/profileImages.ts`  
## Success Criteria
- [ ] All acceptance criteria still pass
- [ ] No validation logic duplicated
- [ ] File handling separated from business logic
- [ ] Error messages centralized and consistent
- [ ] Edge cases covered by tests
- [ ] No new features added (strictly improvements)
- [ ] Shared validation pattern documented in ADR for future reference

## Completion
- [ ] Commit changes: "Refactor: [brief description of improvements]"
- [ ] All refactorings applied and tested
```erify:** Error scenarios display same messages as before

### 4. Add Edge Case Tests
**Why:** Missing tests for corrupted files, network failures  
**What:** Add test cases for error scenarios  
**How:** Use existing test framework, mock failure conditions  
**Verify:** New tests pass, coverage report shows improvement

### 5. Improve Function Names
**Why:** `handleUpload` and `processFile` don't clearly indicate what they do  
**What:** Rename to `uploadProfileImageToS3` and `validateAndPrepareImageFile`  
**How:** IDE rename refactor  
**Verify:** All references updated, tests pass

## Success Criteria
- [ ] All acceptance criteria still pass
- [ ] No validation logic duplicated
- [ ] File handling separated from business logic
- [ ] Error messages centralized and consistent
- [ ] Edge cases covered by tests
- [ ] No new features added (strictly improvements)
```

## Refactoring Principles
## Final Instructions

1. **Verify feature is working** - all acceptance criteria pass
2. **Scan entire codebase** - review ALL project files, not just feature files
3. **Look for patterns** - identify duplication and inconsistencies across the whole project
4. **Review actual implementation code** - don't assume, inspect everything
5. **Check constitutional alignment** - ensure improvements follow principles
6. **Identify project-level decisions** - patterns/conventions that should guide future features
7. **Suggest ADRs when needed** - if decisions have project-wide impact, propose documenting them
8. **Focus on clarity and quality** - make code easier to understand
9. **Keep refactorings small** - 15-30 min each, independently committable
10. **Maintain functionality** - no new features, only improvements
11. **Stay pragmatic** - perfect is the enemy of good
12. **Document why** - explain the problem each refactor solves
13. **Consider extraction** - if patterns appear across files, suggest shared utilities

**Remember:** The CONSTITUTION defines WHY and foundational tech choices. ADRs in `design/` capture HOW patterns that emerged during implementation and should be followed consistently.
❌ **Do everything at once** - incremental is safer

## Target Audience

Refactor plans are for **developers improving working code** who need:
- Clear direction on what to improve
- Understanding of why improvements matter
- Confidence that refactoring is safe
- Incremental steps that maintain functionality
## Final Instructions

1. **Verify feature is working** - all acceptance criteria pass
2. **Scan entire codebase** - review ALL project files, not just feature files
3. **Look for patterns** - identify duplication and inconsistencies across the whole project
4. **Review actual implementation code** - don't assume, inspect everything
5. **Check constitutional alignment** - ensure improvements follow principles
6. **Focus on clarity and quality** - make code easier to understand
7. **Keep refactorings small** - 15-30 min each, independently committable
8. **Maintain functionality** - no new features, only improvements
9. **Stay pragmatic** - perfect is the enemy of good
10. **Document why** - explain the problem each refactor solves
11. **Consider extraction** - if patterns appear across files, suggest shared utilitiess
2. **Review actual implementation code** - don't assume, inspect
3. **Check constitutional alignment** - ensure improvements follow principles
4. **Focus on clarity and quality** - make code easier to understand
5. **Keep refactorings small** - 15-30 min each, independently committable
6. **Maintain functionality** - no new features, only improvements
7. **Stay pragmatic** - perfect is the enemy of good
8. **Document why** - explain the problem each refactor solves
