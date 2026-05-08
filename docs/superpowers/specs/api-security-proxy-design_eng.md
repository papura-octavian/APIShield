# API Security Proxy — Design (personal project)

**Author:** Papură Corneliu-Octavian
**Project type:** personal / portfolio (GitHub, CV)

---

## 1. What it is

A distributed system that protects an arbitrary REST API by intercepting traffic with a Go reverse proxy and routing the decision to a Python detection engine that combines three layers: static rules, ML classifiers, and an anomaly detector.

Covers 5 categories from the OWASP API Security Top 10 (2023):
- BOLA (Broken Object Level Authorization)
- Broken Authentication
- Unrestricted Resource Consumption
- BFLA (Broken Function Level Authorization)
- SSRF (Server-Side Request Forgery)

## 2. Why (for the CV)

- **Go** — a new language on the CV (the author already has solid .NET experience).
- **ML applied to security**, not yet another Kaggle tutorial.
- **A real distributed system**, multi-container, with observability.
- **Comparison against ModSecurity + OWASP CRS** — a visible differentiator in the README.
- **Language-agnostic**: the proxy protects any HTTP API (demo on .NET, optionally Node).

## 3. Architecture

5 Docker containers, orchestrated with `docker-compose`:

1. **Proxy core** (Go) — intercepts HTTP, serializes the request, asks for a verdict, forwards to the backend if OK.
2. **Detection engine** (Python + FastAPI + scikit-learn) — rules → ML classifier → anomaly detector.
3. **Honeypot endpoints** — trap URLs (`/admin`, `/.env`, `/wp-login.php`) integrated into the proxy (runtime feature, NOT a training data source).
4. **Storage** — PostgreSQL (events) + Redis (rate-limit state, cache).
5. **Dashboard** — Angular + ASP.NET Core API: live traffic, attacker map, request inspector, performance metrics.

### Decision flow

```
HTTP request ─▶ Proxy ─▶ Engine: [Rules] ─HIT─▶ verdict
                                    │ miss
                                    ▼
                             [ML Classifier] ─HIGH─▶ verdict
                                    │ low
                                    ▼
                             [Anomaly Det.] ─ANOM─▶ verdict
                                    │ ok
                                    ▼
                               CLEAN ─▶ forward to API
```

### Communication

- Proxy ↔ Engine: gRPC (proto3) for low latency.
- Engine ↔ Storage: async (Go channels / Python asyncio).
- Dashboard ↔ Storage: .NET REST API.

## 4. Stack

| Component | Technology |
|---|---|
| Proxy core | Go (`net/http`, `httputil.ReverseProxy`) |
| Detection engine | Python 3.12, FastAPI, scikit-learn, NumPy, pandas |
| Inter-service comms | gRPC (proto3) |
| Event storage | PostgreSQL 16 |
| Cache / rate-limit | Redis 7 |
| Dashboard backend | ASP.NET Core 10 |
| Dashboard frontend | Angular + Chart.js / ApexCharts |
| Primary demo API | ASP.NET Core 10 (extension of TaskManagerAPI → MiniShop) |
| Secondary demo API (optional) | Node.js + Express (MiniBlog) |
| Containerization | Docker + docker-compose |
| CI/CD | GitHub Actions |
| IP geolocation | MaxMind GeoLite2 |

## 5. ML data

100% public datasets:
- **CSE-CIC-IDS2018** — Canadian Institute for Cybersecurity.
- **HTTP CSIC 2010** — ~25,000 labeled attacks (SQLi, XSS, CRLF, etc.).
- **OWASP API Security Project** — labeled examples.

Additional synthetic traffic generated against MiniShop for attacks not well covered by public datasets (e.g. BFLA on custom roles).

The honeypot **does not** feed training — it only captures attacks at runtime for the dashboard.

## 6. ML algorithms

| Attack | Algorithm | Main features |
|---|---|---|
| BOLA | Random Forest | IDs accessed per user, own/other ratio, ID enumeration speed |
| Broken Auth | Random Forest | Failed attempts / time window, user-agent distribution, geographic dispersion |
| Resource Consumption | XGBoost + token bucket | Rate per endpoint per user, burst patterns |
| BFLA | Random Forest | Role/endpoint mapping, scanning patterns |
| SSRF | Char-level CNN or n-gram URL classifier | URL embeddings, schema, target IP |
| Anomalies (all) | Isolation Forest or Autoencoder | All combined features |

## 7. Demo API — MiniShop (.NET)

Extension of the existing `TaskManagerAPI`:

- `/api/auth/login`, `/api/auth/register` — Broken Auth
- `/api/users/{id}/profile`, `/api/users/{id}/orders` — BOLA
- `/api/admin/users`, `/api/admin/products` — BFLA
- `/api/products/search?query=...` — Resource Consumption
- `/api/integrations/preview?url=...` — SSRF

Optional: MiniBlog (Node.js) to demonstrate language agnosticism.

## 8. Testing and benchmarking

- **Unit tests** — xUnit (.NET), pytest (Python), `go test` (Go).
- **ML evaluation** — k-fold cross-validation (k=5); precision, recall, F1, ROC-AUC, PR-AUC.
- **End-to-end attack sims** — automated OWASP ZAP + custom scripts + Locust/k6.
- **vs. ModSecurity + OWASP CRS** — same tests run on both systems; results published in the README with charts.

### Targets

- Per-request latency overhead: < 10 ms (P95).
- Throughput: > 1000 RPS on a modest VPS.
- F1 per attack category: > 0.85.
- Demonstrable advantage over ModSecurity in ≥ 2 categories (likely BOLA + BFLA).

## 9. Repository layout

```
api-security-proxy/
├── proxy/                    # Go + httputil.ReverseProxy
├── detection-engine/         # Python + FastAPI + scikit-learn
├── dashboard/                # Angular + .NET API
├── demo-apis/
│   ├── minishop-net/
│   └── miniblog-node/        # optional
├── attack-scripts/           # Python + OWASP ZAP automation
├── data/{raw,processed}/
├── benchmarks/modsecurity-config/
├── docs/
├── docker-compose.yml
└── .github/workflows/
```

## 10. Execution phases (ordered, no timing)

### Phase 0 — Preparation
1. Light reading: 3-5 articles on API security + ML (just enough for orientation, no academic bibliography).
2. Download and verify public datasets.
3. Learn Go (Tour of Go + Go by Example).
4. Docker basics + docker-compose.
5. scikit-learn refresh (classifiers + cross-validation).
6. Minimal prototype: Go proxy ↔ Python endpoint (gRPC flow check).

### Phase 1 — Foundation
1. Repo setup + structure.
2. Minimal working `docker-compose.yml`.
3. MiniShop API (extension of TaskManagerAPI).
4. Proxy core: interception, forwarding, structured logging.
5. Storage layer (Postgres + Redis).
6. Detection engine skeleton (FastAPI + integration with proxy).
7. End-to-end logging working.

### Phase 2 — Rules (no ML)
Rule engine for BOLA, Broken Auth, Resource Consumption, BFLA, SSRF + honeypot endpoints + unit tests.

### Phase 3 — ML
1. Data prep (cleaning, labeling, feature engineering).
2. EDA in Jupyter.
3. Train classifier per attack.
4. k-fold evaluation.
5. Anomaly detector.
6. Save/load models in the engine.
7. Integrate the 3-layer pipeline.

### Phase 4 — Dashboard
Backend API + Angular with: live traffic, attacks map, attack types breakdown, top attackers + manual block, request inspector + FP marking, performance metrics.

### Phase 5 — Testing and benchmarking
Full unit tests, automated OWASP ZAP, custom scripts, load tests, e2e in CI, ModSecurity + CRS setup, comparison and data collection.

### Phase 6 — Polish
1. Optimization (pprof, py-spy).
2. Final charts for the README.
3. README with screenshots + diagrams + result tables.
4. Full CI/CD.
5. Live demo (`docker-compose up` → everything functional).

## 11. Risks

1. **Learning Go delays the start** → intensive tutorial BEFORE Phase 1.
2. **Public datasets don't cover everything** → fill gaps with synthetic traffic.
3. **Many false positives** → threshold calibration + an honest section in the README.
4. **Unfavorable comparison vs. ModSecurity overall** → focus on categories with a clear advantage (BOLA, BFLA).

## 12. New technologies learned (CV value)

- Go (DevOps, SRE, Cloud).
- gRPC (modern microservices).
- Docker + docker-compose.
- Redis.
- GitHub Actions (CI/CD).
- Light MLOps.
- MaxMind GeoLite2.
