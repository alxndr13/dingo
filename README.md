# Dingo

An opinionated templating CLI.

## Background

At my last employer, I was introduced to [Kapitan](kapitan.dev), which became the foundation for many projects by enabling multiple environments regardless of technology. We mainly used it to generate Terraform HCL (I’m not a fan of modules) and bash scripts when we didn’t use custom Go programs.

What I disliked about Kapitan was its Python ecosystem and the hassle of installing it on Docker images or new machines. You always had to ensure the exact Python version, or something would break.

When tasked with building a new foundation for a client, my colleagues and I chose to build on a custom Go library inspired by Kapitan. Since it was just a library, we had to add many features and wrap it in a CLI tool to generate:

-  Terraform code
-  HAProxy configurations
-  PostgreSQL scripts
-  Kubernetes manifests, and more

We even presented our "Infrastructure as Data" approach [at a local German conference](https://www.continuouslifecycle.de/veranstaltung-21248-0-wie-wir-mit-infrastructure-as-data-eine-plattform-gebaut-haben.html).

As I moved between projects and employers, I repeatedly needed multiple environments that were easy to create and maintain by modifying data—like structures like updating versions, adding users, or adjusting permissions. Being lazy, I wanted to delegate these tasks with safeguards in place.

That’s how Dingo was born.

## My idea and implementation of a useful general-purpose templating CLI

I wanted a general-purpose templating CLI like Kapitan but written in Go, so I could distribute binaries easily and integrate tightly with the cloud-native ecosystem.

For my MVP, I focused on:

-  Familiar data management using a `base` and multiple `overlays`, terms borrowed from Kustomize.
-  Templates written in Go Templates with extras like the [`sprig` library](https://masterminds.github.io/sprig/).
-  Secrets stored in GCP Secrets Manager and resolved on the fly.
-  Data validation using [cue](https://cuelang.org/) to separate validation from CLI code, avoiding recompilation for schema updates.

## Terms of Dingo

Dingo has four main concepts:

-  **Data (Base and Overlays):** Users insert data here; overlays override base data.
-  **Templates:** Go Templates enriched with data.
-  **Decryptor (optional):** Replaces values wrapped in `$$` with secrets when a decryptor is selected.
-  **Schema:** Defines allowed data fields, their structure, and valid values.

---

to be continued.. i need to go to bed 😴
