---
name: design
description: Sketch an initial design approach for a feature increment
argument-hint: increment name or brief description
---

# Rule: Generating an Initial Design Document

This document represents the **HOW (Initial)** - sketching initial design ideas, potential approaches, and technical considerations before implementation.

## Goal

To guide an AI assistant in creating a focused Architecture Decision Record that documents the key technical decisions for implementing a **small, testable increment**, just as a senior developer would sketch technical approaches. The ADR explains WHICH technical approaches were chosen and WHY, focusing on component boundaries, data flow, and trade-offs specific to this feature.

## Prerequisites

Before generating an ADR, verify that both prerequisite files exist:

**Required:**
1. **`CONSTITUTION.md`** - Defines technical decisions, frameworks, and development principles
2. **`[increment-name]/increment.md`** - Defines WHAT the user wants to accomplish

**If either file does not exist:**
- Inform the user which prerequisite is missing
- Suggest creating the constitution first (using `/create-constitution`)
- Then creating the feature (using `/create-feature`)
- Do not proceed until both files exist

**If both files exist:**
- Read the constitution to understand the broad technical foundation and constraints
- Read the feature to understand what needs to be built (user goals and acceptance criteria)
- Sketch initial design ideas that comply with the constitution while implementing the feature

## Process

1.  **Verify Prerequisites:** Check for `CONSTITUTION.md` and the relevant `increment.md`.
2.  **Receive Initial Prompt:** The user requests a design for a specific increment.
3.  **Analyze Context:** Review constitution and increment to understand constraints and requirements.
4.  **STOP - Ask Clarifying Questions:**
    
    **DO NOT GENERATE THE DESIGN YET.**
    
    Ask only 2-3 essential technical questions not answered by constitution or feature. Wait for user answers.
    
5.  **Generate Design:** (Only after receiving answers) Create a lightweight design sketch for this small increment.
6.  **Save Design:** Save as `[increment-name]/design.md`.

## Before Generating Design - Self Check

Ask yourself:
- [ ] Did I verify CONSTITUTION.md and increment.md exist?
- [ ] Did I read both documents to understand constraints and requirements?
- [ ] Did I ask 2-3 technical clarifying questions?
- [ ] Did I receive user's answers?

If any checkbox is unchecked, STOP and complete that step first.

## Clarifying Questions (Guidelines)

Ask only 2-3 technical questions not already answered by the constitution or use case. Focus on **this specific increment's** architecture.

### Focus Areas:
- **Component boundaries:** How to split responsibilities for this feature
- **Data flow:** How information moves through the system
- **Integration patterns:** How to connect to external services
- **Storage decisions:** What persists vs. what's ephemeral

**Important:** Don't ask about tech stack, frameworks, or principles—those are in the constitution. Focus on this feature's specific technical decisions.

### Formatting Requirements

- **Number all questions** (1, 2, 3)
- **List options as A, B, C, D, etc.**
- **Always include option X: "Skip this question / I don't know yet"**
- Keep it simple for responses like "1A, 2C"

### Example Questions

```
1. How should data flow for this feature?
   A. Client → API → Database → Client
   B. Client → API (async job) → Client polls
   C. Client → Third-party API directly
   X. Skip this question / I don't know yet

2. Where should this feature's state live?
   A. Client-side only (component state)
   B. Backend session/cache
   C. Database (persistent)
   D. Hybrid (client + backend)
   X. Skip this question / I don't know yet

3. How should this integrate with external services?
   A. Direct API calls
   B. Queue-based async processing
   C. Webhook callbacks
   X. Skip this question / I don't know yet
```

## Design Document Structure

Keep designs focused on **initial approach for small increments**. Aim for one screen or less. This is a sketch, not final decisions.

```markdown
# Design: [Feature Name]

**Date:** [YYYY-MM-DD]  
**Status:** Initial Sketch

## Design Summary
[2-3 sentences: What we're building and the initial technical approach to try]

## Initial Approach

### [Design Consideration]
**Approach:** [What we'll try first]  
**Rationale:** [Why this seems reasonable - technical reasoning]  
**Trade-offs:** [What we're accepting vs. what we're gaining]  
**Alternatives to Consider:** [Other approaches we might pivot to]

### [Design Consideration]
**Approach:** [What we'll try]  
**Rationale:** [Why]  
**Trade-offs:** [Gains vs. costs]  
**Alternatives to Consider:** [Other options]

[Repeat for 2-5 key considerations only]

## Architecture Overview

**Components:**
- [Component name]: [Responsibility for this feature]
- [Component name]: [Responsibility]

**Data Flow:**
[Brief description or simple diagram of how data moves]

**Integration Points:**
- [External service/API]: [How we integrate]

**State Management:**
[Where state lives and why for this feature]

## Implementation Constraints
- [Technical constraint or limitation from these decisions]
- [Another constraint developers need to know]

## Open Questions
- [Technical unknown to resolve during implementation]
- [Deferred decision with reason]
```

---

## Template Example

```markdown
# Design: Upload Profile Picture

**Date:** 2025-11-24  
**Status:** Initial Sketch

## Design Summary
We'll try client-side image upload with backend validation and storage in cloud object storage. Prioritizes simple implementation over advanced features like cropping or compression.

## Initial Approach

### Image Upload Pattern
**Approach:** Client uploads directly to backend API, backend stores in S3-compatible storage  
**Rationale:** Simple request/response pattern fits constitutional "speed over complexity". No presigned URLs or direct-to-S3 needed for MVP.  
**Trade-offs:** Backend handles file upload (bandwidth cost) but simpler to implement and secure. Accept slower uploads for cleaner architecture.  
**Alternatives to Consider:** Direct client→S3 upload (presigned URLs) - might revisit if upload speed becomes an issue.

### Image Validation
**Approach:** Backend validates file type and size before storage  
**Rationale:** Server-side validation is trustworthy; client-side can be bypassed. Prevents malicious uploads.  
**Trade-offs:** Upload completes then fails validation (wasted bandwidth) but more secure. Accept this for MVP.  
**Alternatives to Consider:** Client-only validation - skip due to security; dual validation - unnecessary complexity for now.

### Storage Location
**Approach:** Cloud object storage (S3/compatible), not database  
**Rationale:** Images are large binary data; database storage would bloat backup size and slow queries. Object storage is designed for this.  
**Trade-offs:** External dependency on S3 but appropriate tool for the job. Cost is negligible at current scale.  
**Alternatives to Consider:** Database blob storage - avoid due to performance; filesystem storage - avoid due to scaling/backup complexity.

## Architecture Overview

**Components:**
- ProfileUploadForm: Handles file selection and upload UI
- ProfileAPI: Validates and stores image, returns URL
- ObjectStorage: Persists image files

**Data Flow:**
User selects file → Client validates size/type → POST /api/profile/picture → Backend validates → Store in S3 → Return public URL → Client displays

**Integration Points:**
- AWS S3 (or compatible): Image storage with public read access

**State Management:**
- Profile picture URL stored in user profile (database)
- Actual image in S3 with predictable key: `profile-pictures/{userId}.jpg`

## Implementation Constraints
- Images must be JPG/PNG only (validation at API boundary)
- Max file size: 5MB (enforced backend, suggested frontend)
- No image transformation (cropping/resizing) - users prepare images before upload
- Public read access required (no authentication on image URLs)

## Open Questions
- CDN needed? Defer until traffic shows need
- Image compression? Defer to post-MVP if users upload large files
- Multiple image formats (WebP)? Defer until browser support is issue
```

## Target Audience

Design documents are read by **developers implementing the feature** who need:
- An initial technical direction to start with
- Understanding of trade-offs and alternatives
- Constraints to respect during implementation

Keep designs **lightweight and flexible**—focused on getting started, not final decisions. Expect to learn and adjust during implementation.

## No Code in Design Documents

**Important:** Design documents sketch approaches, not implementations.
- Actual code belongs in the project repository
- Design documents suggest starting points, not prescribe final solutions
- Focus on initial architectural ideas, component boundaries, and data flow
- Expect to evolve these ideas during implementation

## Output

*   **Format:** Markdown (`.md`)
*   **Location:** `[increment-name]/` directory
*   **Filename:** `design.md`

## Final Instructions

1. **Verify `CONSTITUTION.md` and `increment.md` exist** before proceeding
2. **Read both documents** to understand principles and requirements
3. **Focus on initial approach for this increment** - not the entire system
4. **Ask 2-3 clarifying questions** only about architectural uncertainties
5. **Sketch 2-5 key considerations** - component boundaries, data flow, integration, storage
6. **Make trade-offs explicit** - what we're accepting vs. what we're gaining
7. **Keep it concise** - aim for one screen, lightweight starting point
8. **No implementation code** - only initial architectural ideas and rationale
9. **Stay flexible** - this is a starting point, not final decisions
