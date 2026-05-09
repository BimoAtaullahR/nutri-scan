# Backend API

The Backend API owns NutriScan's user-facing product workflows, persistence, scan orchestration, nudge decisions, and trend APIs.

## Language

**Scan**:
A user-initiated attempt to analyze a food image before eating.
_Avoid_: detection request, image check

**Scan Lifecycle**:
The backend-owned progression of a **Scan** through creation, processing, completion, or failure.
_Avoid_: AI status, upload status

**Nudge Decision**:
The backend-owned recommendation shown to the user after a **Scan**, based on inference results and product rules.
_Avoid_: AI recommendation, model advice

**Anonymous User**:
A user represented by a generated identity before NutriScan introduces full account authentication.
_Avoid_: guest session, temporary account

**User Profile**:
Personal attributes attached to a user for personalization, such as height and weight, without implying a registered account.
_Avoid_: account, medical record

**BMI Category**:
A coarse personalization category derived from height and weight for preventive feedback.
_Avoid_: diagnosis, medical classification

## Relationships

- A **Scan** has exactly one **Scan Lifecycle**
- A **Scan** may produce one **Nudge Decision**
- A **Nudge Decision** is based on inference results from AI/ML Inference
- A **Weekly Energy Trend** is calculated from completed **Scans**
- An **Anonymous User** can own many **Scans**
- An **Anonymous User** may have one **User Profile**
- A **Nudge Decision** may consider a **User Profile**
- A **User Profile** may produce one **BMI Category**

## Backend Ownership

- **Scan** owns scan orchestration and persistence.
- **Nudge** owns product rules for preventive nudges.
- **Trend** owns weekly energy trend reporting.
- **Inference Client** adapts Backend API workflows to AI/ML Inference contracts.

## Example Dialogue

> **Dev:** "If AI inference times out, is the **Scan** deleted?"
> **Domain expert:** "No — the **Scan Lifecycle** moves to failed so the app can show a retry path and the backend can inspect failures."

## Flagged Ambiguities

- "AI recommendation" was used to describe the user-facing nudge — resolved: this is a backend-owned **Nudge Decision**, not direct model advice.
- "guest" was used for unauthenticated usage — resolved: MVP uses **Anonymous User** as a durable generated identity, not a throwaway session.
- "account" was used for personalization data — resolved: MVP uses **User Profile** for personal attributes, separate from authentication.
- "BMI" was used as if clinical advice — resolved: MVP may use **BMI Category** only for coarse preventive personalization, not diagnosis.
