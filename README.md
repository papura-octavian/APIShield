# APIShield

> ML-powered reverse proxy that detects and blocks **OWASP API Top 10** attacks
> in real time. Language-agnostic, container-native, benchmarked against
> ModSecurity + OWASP CRS.

![Status](https://img.shields.io/badge/status-design--phase-orange)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)
![Python](https://img.shields.io/badge/Python-3.12-3776AB?logo=python)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)

---

## What is this

**APIShield** is a distributed system that protects any HTTP REST API by
intercepting traffic through a **Go reverse proxy** and forwarding decisions
to a **Python detection engine** that combines three layers:

1. **Static rules** — fast, deterministic checks for known attack patterns.
2. **ML classifiers** — supervised models trained per attack category
   (Random Forest, XGBoost, char-level CNN).
3. **Anomaly detection** — Isolation Forest / Autoencoder for unknown threats.

It targets 5 categories from the **OWASP API Security Top 10 (2023)**:

| # | Attack | Primary detection |
|---|---|---|
| 1 | **BOLA** — Broken Object Level Authorization | Random Forest on per-user access patterns |
| 2 | **Broken Authentication** | Random Forest on auth-failure patterns |
| 3 | **Unrestricted Resource Consumption** | XGBoost + token bucket |
| 4 | **BFLA** — Broken Function Level Authorization | Random Forest on role/endpoint mapping |
| 5 | **SSRF** — Server-Side Request Forgery | Char-level CNN / n-gram classifier on URLs |

A runtime **honeypot** layer also exposes trap endpoints (`/admin`, `/.env`, `/wp-login.php`)
to flag scanners — but training data comes 100% from public datasets, not honeypot capture.

---

## Why this project

- Demonstrates **applied ML in security** — not a tutorial-grade Kaggle notebook.
- Built with **Go** (proxy core) and **Python** (ML engine) — real distributed system,
  not a single service.
- Includes a **head-to-head benchmark vs. ModSecurity + OWASP CRS** with published results.
- Fully containerized with **docker-compose** — one command to bring up the whole stack.
- **Language-agnostic protection**: the proxy doesn't care if the protected API is
  .NET, Node.js, Python, Java, or anything else over HTTP.

---

## Architecture

```
                   ┌────────────────────────┐
   Client ────────▶│   Go Reverse Proxy     │
                   │   (httputil + gRPC)    │
                   └───────────┬────────────┘
                               │ verdict?
                               ▼
                   ┌────────────────────────┐
                   │  Detection Engine      │
                   │  (Python + FastAPI)    │
                   │                        │
                   │   Rules ──┐            │
                   │           ▼            │
                   │     ML Classifier ──┐  │
                   │                     ▼  │
                   │       Anomaly Detector │
                   └───────────┬────────────┘
                               │
                       ┌───────┴────────┐
                       ▼                ▼
                 ┌──────────┐     ┌──────────┐
                 │PostgreSQL│     │  Redis   │
                 │ (events) │     │ (state)  │
                 └──────────┘     └──────────┘
                       ▲
                       │
                ┌─────────────────┐
                │   Dashboard     │
                │  (Angular +     │
                │   ASP.NET API)  │
                └─────────────────┘
```

**5 containers**, orchestrated with `docker-compose`:
proxy · detection-engine · postgres · redis · dashboard.

---

## Tech stack

| Layer | Tech |
|---|---|
| Proxy core | Go (`net/http`, `httputil.ReverseProxy`) |
| Detection engine | Python 3.12 · FastAPI · scikit-learn · NumPy · pandas |
| Inter-service comms | gRPC (proto3) |
| Storage | PostgreSQL 16 · Redis 7 |
| Dashboard | Angular · ASP.NET Core 10 · Chart.js |
| Demo APIs | ASP.NET Core 10 (MiniShop) · Node.js (MiniBlog, optional) |
| Containerization | Docker · docker-compose |
| CI/CD | GitHub Actions |
| Geolocation | MaxMind GeoLite2 |

---

## Status

> ⚠️ **Currently in design phase.** No implementation code yet.

This repo currently holds:

- `docs/superpowers/specs/api-security-proxy-design_eng.md` — full design
  document in English (architecture, ML algorithms per attack, evaluation
  methodology, ordered execution phases, risks).
- `docs/superpowers/specs/api-security-proxy-design_ro.md` — same document in
  Romanian.

Implementation is split into 6 ordered phases (see the design doc):

1. **Foundation** — repo, docker-compose, demo API, proxy skeleton, storage.
2. **Rules** — rule-based detection for all 5 attack categories.
3. **ML** — data prep, training, anomaly detector, 3-layer pipeline integration.
4. **Dashboard** — Angular + .NET API with live traffic and analytics.
5. **Testing & benchmark** — unit + e2e + load tests, ModSecurity comparison.
6. **Polish** — perf optimization, final graphs, full CI/CD.

---

## Datasets

Training uses **only public datasets** (no scraped or proprietary data):

- [CSE-CIC-IDS2018](https://www.unb.ca/cic/datasets/ids-2018.html) — Canadian Institute for Cybersecurity.
- [HTTP CSIC 2010](https://www.tic.itefi.csic.es/dataset/) — ~25,000 labeled web attacks.
- [OWASP API Security Project](https://owasp.org/www-project-api-security/) — labeled examples.

Synthetic traffic is generated against MiniShop for attack categories not well
covered by public datasets (e.g. BFLA on custom roles).

---

## Performance targets

| Metric | Target |
|---|---|
| Latency overhead per request (P95) | < 10 ms |
| Throughput on a modest VPS | > 1000 RPS |
| F1 score per attack category | > 0.85 |
| Categories with measurable advantage over ModSecurity | ≥ 2 (likely BOLA + BFLA) |

---

## Roadmap

- [x] Project design finalized
- [ ] Phase 1 — Foundation
- [ ] Phase 2 — Rule engine
- [ ] Phase 3 — ML pipeline
- [ ] Phase 4 — Dashboard
- [ ] Phase 5 — Testing & benchmark
- [ ] Phase 6 — Polish & release

---

## License

MIT — see [`LICENSE`](LICENSE).

---

## Author

**Papură Corneliu-Octavian** — Computer Science student, University of Craiova.
Built as a personal portfolio project.
