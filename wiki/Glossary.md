# Glossary

| Istilah | Definisi |
|---------|----------|
| **NAEOS** | Nusantara Engineering & Architecture Operating System — platform engineering system open-source |
| **NEIR** | Nusantara Engineering Intermediate Representation — model intermediasi yang merepresentasikan proyek secara unified |
| **Spec** | Spesifikasi — dokumen YAML/JSON yang mendefinisikan proyek |
| **Pipeline** | Rantai pemrosesan: parse → normalize → resolve → build → validate → generate |
| **Kernel** | Komponen inti sistem yang mengelola service registry, event bus, dan telemetry |
| **Artifact** | Output yang dihasilkan oleh pipeline (kode, config, docs) |
| **Profile** | Profil industri yang menyediakan template siap pakai (SaaS, FinTech, Healthcare, dll) |
| **Adapter** | Penghasil output untuk target tertentu (Copilot, Claude, Cursor, Gemini, Codex, OpenCode) |
| **Compiler** | Komponen yang mentransformasikan NEIR ke instruction sets untuk AI tools |
| **Context Bundle** | Ringkasan proyek dalam format markdown/plain text untuk konsumsi LLM |
| **Governance** | Sistem tata kelola yang mengatur kebijakan, validasi, dan review |
| **Policy** | Aturan yang dievaluasi selama pipeline berjalan (exists, not_empty, contains, gt, lt, in) |
| **Schema Version** | Versi SemVer spesifikasi (minimum: 0.1.0) |
| **Module** | Unit kode dalam proyek dengan name, path, dependencies |
| **Service** | Komponen runtime (http, grpc, worker, cli, job) |
| **Endpoint** | Titik masuk API (method + path + action) |
| **DAG** | Directed Acyclic Graph — struktur data untuk dependency resolution |
| **Artifact Store** | Penyimpanan artifact dengan content-hash dedup |
| **Migration** | Proses upgrade schema versi spesifikasi |
