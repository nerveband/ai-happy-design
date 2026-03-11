# Shiptypes.com -- Baseline & Principles

Source: [shiptypes.com](https://shiptypes.com) by Boris Tane (Cloudflare)

> "Ship types, not docs."
> Types are the UX of your codebase.

## Core Principles

### 1. Schema-First Architecture
Define contracts once in a schema language. Everything else generated from that single schema.
**AHD Status**: PARTIAL -- 23/130 commands have schemas. Target: 100%.

### 2. Types as Executable Documentation
A type definition is an executable contract that cannot drift. Documentation is a lossy copy.
**AHD Status**: PARTIAL -- Go schemas are rich but describe.go is manual. Target: Generate describe.go from schemas.

### 3. The Drift Problem
Any system where two things must stay in sync manually WILL fall out of sync. Eliminate duplicates.
**AHD Status**: FAIL -- 6 independent artifact groups. Target: Schema as single source, everything derived.

### 4. AI-Native Codebase
Agent with types gets it right first call. Agent with only docs takes 3-4 attempts.
**AHD Status**: PARTIAL -- Good tooling but incomplete schema coverage. Target: Every command typed.

### 5. Every Surface Area Typed
All boundaries should be typed.
**AHD Status**: FAIL -- Plugin handlers use `params: any`. Target: Generated TypeScript interfaces.

### 6. Compiler as Reviewer
Can't ship code that violates the contract.
**AHD Status**: FAIL -- No build check that schemas match handlers. Target: Coverage test.

### 7. Breaking Changes Hard
Found at build time, not production.
**AHD Status**: FAIL -- Param rename silently breaks. Target: Coverage test catches.

## Action Items (Priority Order)
1. Extend schemas to all ~130 commands
2. Add schema<>handler coverage test
3. Generate describe.go from schemas
4. Unify schema packages (schema + commonschema)
5. Type plugin handler params (future)
6. Type response envelope (future)
