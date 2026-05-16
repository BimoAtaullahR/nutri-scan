# Backend API

The Backend API owns NutriScan's user-facing product workflows, persistence, scan orchestration, nudge decisions, and trend APIs.

## Language

**Scan**:
A user-initiated attempt to analyze a food image before eating.
_Avoid_: detection request, image check

**Scan Lifecycle**:
The backend-owned progression of a **Scan** through creation, processing, completion, or failure.
_Avoid_: AI status, upload status

**Sync-First Scan**:
A scan workflow that tries to return completed feedback within the request window, while preserving scan lifecycle polling when inference is slow or fails.
_Avoid_: blocking-only scan, background-only scan

**Nudge Decision**:
The backend-owned recommendation shown to the user after a **Scan**, based on inference results and product rules.
_Avoid_: AI recommendation, model advice

**Generic Nudge Decision**:
A nudge decision produced without a user profile, using only scan and inference information.
_Avoid_: incomplete recommendation, unpersonalized AI advice

**Personalized Nudge Decision**:
A nudge decision that also considers a user profile for coarse personalization.
_Avoid_: medical recommendation, diagnosis-based advice

**Nudge Action**:
The user-facing action proposed by a nudge decision, such as eating as planned, setting aside a portion, or reviewing a low-confidence scan.
_Avoid_: nutrient prescription, medical instruction

**Estimated Prevented Energy**:
Approximate energy the user may avoid if they follow a nudge action.
_Avoid_: exact calories saved, clinical outcome

**Review Food Nudge**:
A nudge action used when inference confidence is too low and the user should review or retry the food result.
_Avoid_: failed scan, model error

**Scan Image**:
The short-lived image bytes submitted for a scan and forwarded to AI/ML Inference without default long-term storage.
_Avoid_: stored food photo, image archive

**Anonymous User**:
A user represented by a generated identity before NutriScan introduces full account authentication.
_Avoid_: guest session, temporary account

**User Profile**:
Personal attributes attached to a user for personalization, such as height and weight, without implying a registered account.
_Avoid_: account, medical record

**BMI Category**:
A coarse personalization category derived from height and weight for preventive feedback.
_Avoid_: diagnosis, medical classification

**Core Scan Loop**:
The product loop where a user scans food, receives estimated energy feedback and a nudge, records a response, and later sees the result reflected in trends.
_Avoid_: full app completion, all-screen backend coverage

**Daily Energy Summary**:
A day-level snapshot of estimated energy eaten, remaining energy, and MVP placeholder burned energy for the current user.
_Avoid_: nutrition dashboard, macro report

**Meal Energy Summary**:
A per-meal grouping of completed scan estimated energy for breakfast, lunch, dinner, and snack.
_Avoid_: meal plan, exact macro breakdown

**Meal Type**:
The user's intended meal bucket for a scan: breakfast, lunch, dinner, or snack.
_Avoid_: food category, eating schedule

## Relationships

- A **Scan** has exactly one **Scan Lifecycle**
- A **Sync-First Scan** starts as processing and may complete inside the create request
- A **Scan** may produce one **Nudge Decision**
- A **Nudge Decision** is based on inference results from AI/ML Inference
- A **Nudge Decision** may be a **Generic Nudge Decision** when no **User Profile** exists
- A **Nudge Decision** may be a **Personalized Nudge Decision** when a **User Profile** exists
- A **Nudge Decision** has one **Nudge Action**
- A **Nudge Decision** may include **Estimated Prevented Energy**
- A **Review Food Nudge** is a valid **Nudge Decision** outcome for low-confidence inference
- A **Scan** starts from one **Scan Image**
- A **Weekly Energy Trend** is calculated from completed **Scans**
- An **Anonymous User** can own many **Scans**
- An **Anonymous User** may have one **User Profile**
- A **Nudge Decision** may consider a **User Profile**
- A **User Profile** may produce one **BMI Category**
- The **Core Scan Loop** starts with a **Scan** and can produce a **Nudge Decision**, a nudge response, and **Weekly Energy Trend** data
- A **Daily Energy Summary** is calculated from completed **Scans** for one day
- A **Daily Energy Summary** includes a **Meal Energy Summary**
- A **Scan** has one **Meal Type**, either provided by the Mobile App or assigned by Backend API fallback rules

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
- "backend 50% progress" was used as if it could mean all endpoints or all screens — resolved: the 50% backend target means a usable **Core Scan Loop**, not complete backend coverage for every designed screen.
- "Profile required before scan" was considered for personalization — resolved: scanning is allowed before a **User Profile** exists, and the backend returns a **Generic Nudge Decision** until personalization data is available.
- "Homepage backend scope" was broader in the Figma design than the 50% backend target — resolved: the target includes **Daily Energy Summary** and **Meal Energy Summary**, while water tracking, achievement journaling, and accurate macro tracking remain out of scope.
- "Meal bucket" was used for homepage grouping — resolved: the canonical term is **Meal Type**, preferably sent by the Mobile App and only assigned by fallback time rules when absent.
- "Aura plate nutrient breakdown" was implied by Figma labels — resolved: the 50% backend target returns **Nudge Decisions** based on estimated energy, not accurate protein, carbohydrate, vitamin, or mineral reports.
- "Scan response mode" was ambiguous — resolved: the 50% backend target uses **Sync-First Scan** so mobile can get immediate feedback when inference is fast and poll when it is not.
- "Low-confidence scan" was considered as a failure — resolved: low confidence completes the **Scan Lifecycle** with a **Review Food Nudge**; failed scans are reserved for technical failures.
- "Uploaded food image" was considered as persisted product data — resolved: the canonical term is **Scan Image**, and backend forwards it to AI/ML Inference without default long-term storage.
