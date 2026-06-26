# Product Audit Skill

This skill helps Product Owners perform **structured product audits** using a checklist-based framework. 

It guides POs to verify product behaviour, detect inconsistencies, and identify potential risks across product flows.

---

## Preface

The skill is currently defined best for **Back-end Product** with clear confluence documentation. For user-facing product, the development is in place to integrate browser capture for auditing UI/UX.

The skill is a draft for your product audit. Please ***validate*** and ***update*** the file to best-fit your product current behaviour.

---

## Setup

Before using the skill, complete the following:

1. **Run your local AI agent**

Ask the agent to load the skill from the repository:

> Load the Product Audit skill from `/product-audit`

The agent will read the files in this folder and initialize the audit checklist.

2. **Configure Confluence API Token (recommended)**

Set up your Confluence API token so the agent can retrieve relevant documentation during audits.

3. **Update team context**

Update:

```

/references/team-context.md

```

Include team structure, system components, and important internal terminology so the agent can interpret audit checks correctly.

4. **Get your Confluence documents ready for a faster execution**

Make sure that the information including all technical implementation, UI/UX flows are documented on Confluence. 
The skill is equipped to with *Fetch Confluence* script in order to read information from the latest Confluence documentations.

---

## How to Use

1. Run your local AI agent  
2. Load the skill from `/product-audit`  
3. Ask the agent to perform a **product audit for a specific feature or flow**

The agent will guide you using the checklist defined in:

```

/product-audit/SKILL.md

```

---

## Reference

All audit logic and checklist items are defined in:

```

SKILL.md

```

Please read this file for the full audit checklist and failure examples.

