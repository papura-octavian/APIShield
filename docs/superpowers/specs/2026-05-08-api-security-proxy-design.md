# API Security Proxy — Design (proiect personal)

**Autor:** Papură Corneliu-Octavian
**Tip proiect:** personal / portfolio (GitHub, CV)
**Data:** 2026-05-08

---

## 1. Ce este

Sistem distribuit care protejează un API REST arbitrar prin interceptarea traficului cu un reverse proxy în Go și rutarea deciziei către un motor de detecție Python care combină trei straturi: reguli statice, clasificatori ML, detector de anomalii.

Acoperă 5 categorii din OWASP API Security Top 10 (2023):
- BOLA (Broken Object Level Authorization)
- Broken Authentication
- Unrestricted Resource Consumption
- BFLA (Broken Function Level Authorization)
- SSRF (Server-Side Request Forgery)

## 2. De ce (pentru CV)

- **Go** — limbaj nou pe CV (autorul știe deja .NET solid).
- **ML aplicat în security**, nu doar tutoriale Kaggle.
- **Sistem distribuit real**, multi-container, cu observabilitate.
- **Comparație cu ModSecurity + OWASP CRS** — diferențiator vizibil în README.
- **Agnostic de limbaj**: proxy-ul protejează orice API HTTP (demo pe .NET, optional Node).

## 3. Arhitectură

5 containere Docker, orchestrate cu `docker-compose`:

1. **Proxy core** (Go) — interceptează HTTP, serializează cererea, cere verdict, forward la backend dacă OK.
2. **Detection engine** (Python + FastAPI + scikit-learn) — reguli → clasificator ML → anomaly detector.
3. **Honeypot endpoints** — URL-uri capcană (`/admin`, `/.env`, `/wp-login.php`) integrate în proxy (runtime feature, NU sursă de training).
4. **Storage** — PostgreSQL (events) + Redis (rate-limit state, cache).
5. **Dashboard** — Angular + ASP.NET Core API: live traffic, hartă atacatori, request inspector, performance.

### Flux decizional

```
Cerere HTTP ─▶ Proxy ─▶ Engine: [Rules] ─HIT─▶ verdict
                                   │ miss
                                   ▼
                            [ML Classifier] ─HIGH─▶ verdict
                                   │ low
                                   ▼
                            [Anomaly Det.] ─ANOM─▶ verdict
                                   │ ok
                                   ▼
                              CLEAN ─▶ forward la API
```

### Comunicare

- Proxy ↔ Engine: gRPC (proto3) pentru latență scăzută.
- Engine ↔ Storage: async (channels Go / asyncio Python).
- Dashboard ↔ Storage: REST API .NET.

## 4. Stack

| Componentă | Tehnologie |
|---|---|
| Proxy core | Go (`net/http`, `httputil.ReverseProxy`) |
| Detection engine | Python 3.12, FastAPI, scikit-learn, NumPy, pandas |
| Comunicare | gRPC (proto3) |
| Storage events | PostgreSQL 16 |
| Cache / rate-limit | Redis 7 |
| Dashboard backend | ASP.NET Core 10 |
| Dashboard frontend | Angular + Chart.js / ApexCharts |
| Demo API principal | ASP.NET Core 10 (extindere TaskManagerAPI → MiniShop) |
| Demo API secundar (opțional) | Node.js + Express (MiniBlog) |
| Containerizare | Docker + docker-compose |
| CI/CD | GitHub Actions |
| Geolocation IP | MaxMind GeoLite2 |

## 5. Date pentru ML

100% dataset-uri publice:
- **CSE-CIC-IDS2018** — Canadian Institute for Cybersecurity.
- **HTTP CSIC 2010** — ~25.000 atacuri etichetate (SQLi, XSS, CRLF etc.).
- **OWASP API Security Project** — exemple etichetate.

Trafic sintetic suplimentar pe MiniShop pentru atacuri specifice neacoperite (ex: BFLA pe roluri custom).

Honeypot-ul **nu** alimentează antrenarea — doar capturează atacuri în runtime pentru dashboard.

## 6. Algoritmi ML

| Atac | Algoritm | Features principale |
|---|---|---|
| BOLA | Random Forest | ID-uri accesate per user, raport propriu/altul, viteza incrementare ID |
| Broken Auth | Random Forest | Tentative eșuate / fereastră timp, distribuție user-agent, dispersie geografică |
| Resource Consumption | XGBoost + token bucket | Rate per endpoint per user, burst patterns |
| BFLA | Random Forest | Mapping role/endpoint, scanare patterns |
| SSRF | Char-level CNN sau n-gram classifier pe URL | URL embeddings, schema, IP target |
| Anomalii (toate) | Isolation Forest sau Autoencoder | Toate features-urile combinate |

## 7. Demo API — MiniShop (.NET)

Extensie a `TaskManagerAPI` existent:

- `/api/auth/login`, `/api/auth/register` — Broken Auth
- `/api/users/{id}/profile`, `/api/users/{id}/orders` — BOLA
- `/api/admin/users`, `/api/admin/products` — BFLA
- `/api/products/search?query=...` — Resource Consumption
- `/api/integrations/preview?url=...` — SSRF

Opțional: MiniBlog (Node.js) pentru a demonstra agnosticismul de limbaj.

## 8. Testing și benchmark

- **Unit tests** — xUnit (.NET), pytest (Python), `go test` (Go).
- **ML evaluation** — k-fold cross-validation (k=5); precision, recall, F1, ROC-AUC, PR-AUC.
- **End-to-end attack sims** — OWASP ZAP automate + scripturi proprii + Locust/k6.
- **Comparație vs ModSecurity + OWASP CRS** — aceleași teste pe ambele sisteme; rezultate publicate în README cu grafice.

### Ținte

- Latency overhead: < 10ms (P95).
- Throughput: > 1000 RPS pe VPS modest.
- F1 per categorie: > 0.85.
- Avantaj demonstrabil față de ModSecurity pe ≥ 2 categorii (probabil BOLA + BFLA).

## 9. Structura repo

```
api-security-proxy/
├── proxy/                    # Go + httputil.ReverseProxy
├── detection-engine/         # Python + FastAPI + scikit-learn
├── dashboard/                # Angular + .NET API
├── demo-apis/
│   ├── minishop-net/
│   └── miniblog-node/        # opțional
├── attack-scripts/           # Python + OWASP ZAP automation
├── data/{raw,processed}/
├── benchmarks/modsecurity-config/
├── docs/
├── docker-compose.yml
└── .github/workflows/
```

## 10. Etape de execuție (ordonate, fără timing)

### Etapa 0 — Pregătire
1. Citire scurtă: 3-5 articole pe API security + ML (suficient pentru orientare, fără bibliografie academică).
2. Descărcare și verificare dataset-uri publice.
3. Învățare Go (Tour of Go + Go by Example).
4. Docker basics + docker-compose.
5. Refresh scikit-learn (clasificatori + cross-validation).
6. Prototip minimal: proxy Go ↔ endpoint Python (verificare flow gRPC).

### Etapa 1 — Fundație
1. Setup repo + structură.
2. `docker-compose.yml` minimal funcțional.
3. MiniShop API (extindere TaskManagerAPI).
4. Proxy core: interceptare, forward, structured logging.
5. Storage layer (Postgres + Redis).
6. Detection engine schelet (FastAPI + integrare cu proxy).
7. Logging end-to-end funcțional.

### Etapa 2 — Reguli (fără ML)
Rule engine pentru BOLA, Broken Auth, Resource Consumption, BFLA, SSRF + honeypot endpoints + unit tests.

### Etapa 3 — ML
1. Pregătire date (curățare, etichetare, feature engineering).
2. EDA în Jupyter.
3. Antrenare classifier per atac.
4. Evaluare cu k-fold.
5. Anomaly detector.
6. Save/load modele în engine.
7. Integrare pipeline 3 straturi.

### Etapa 4 — Dashboard
Backend API + Angular cu: live traffic, attacks map, attack types breakdown, top attackers + manual block, request inspector + FP marking, performance metrics.

### Etapa 5 — Testing și benchmark
Unit tests complete, OWASP ZAP automate, scripturi proprii, load tests, e2e în CI, setup ModSecurity + CRS, comparație și colectare date.

### Etapa 6 — Polish
1. Optimizări (pprof, py-spy).
2. Grafice finale pentru README.
3. README cu screenshots + diagrame + tabele de rezultate.
4. CI/CD complet.
5. Demo live (`docker-compose up` → totul funcțional).

## 11. Riscuri

1. **Învățare Go întârzie startul** → tutorial intensiv ÎNAINTE de Etapa 1.
2. **Dataset-uri publice nu acoperă tot** → completare cu trafic sintetic.
3. **Multe false positives** → calibrare praguri + secțiune onestă în README.
4. **Comparație nefavorabilă cu ModSecurity per total** → focus pe categoriile cu avantaj clar (BOLA, BFLA).

## 12. Tehnologii noi învățate (puncte CV)

- Go (DevOps, SRE, Cloud).
- gRPC (microservicii moderne).
- Docker + docker-compose.
- Redis.
- GitHub Actions (CI/CD).
- MLOps light.
- MaxMind GeoLite2.
