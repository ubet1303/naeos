export interface FeatureData {
  icon: string;
  title: string;
  description: string;
  details: string[];
  code?: string;
  codeLang?: string;
}

export const featuresEN: FeatureData[] = [
  {
    icon: "⚡",
    title: "Core Pipeline",
    description: "An end-to-end pipeline that transforms YAML/JSON specifications into production-ready artifacts through 9 stages.",
    details: [
      "Parser — YAML/JSON parsing with variable interpolation, cross-references, and file includes",
      "Normalizer — data normalization and type coercion",
      "Resolver — cross-reference and import resolution with depth-limiting",
      "NEIR Builder — unified project model (NAEOS Engineering Intermediate Representation)",
      "Validator — comprehensive validation: circular deps, port conflicts, module boundaries",
      "Scheduler — DAG-based task scheduling with parallel execution",
      "Generator — multi-language code generation (Go, TypeScript, Python, Java, Rust)",
      "Renderer — template rendering with 10+ built-in templates",
      "Reviewer — policy evaluation and artifact review",
    ],
    code: `# Pipeline configuration
pipeline:
  name: my-project
  mode: development
  verbose: true
  output_dir: ./out
  cache:
    enabled: true
    max_age: 5m`,
    codeLang: "yaml",
  },
  {
    icon: "📋",
    title: "Spec Language v2",
    description: "A powerful specification language designed for describing software systems of any scale.",
    details: [
      '${var} — variable interpolation with nested access (e.g., ${modules.auth.name})',
      '$env{VAR} — environment variable resolution',
      '$ref{path} — cross-reference resolution across the spec tree',
      '$include{file} — multi-file spec composition',
      '$import{path} — modular spec fragments with caching',
      '$fn{name(args)} — custom functions: upper, lower, slug, default, len, coalesce',
      '$if{condition} / $endif — conditional sections',
      'Schema versioning with auto-check (minimum v0.1.0)',
    ],
    code: `# Spec Language v2 features
project: my-app
modules:
  - name: auth
    path: ./internal/auth
    description: Authentication module
    dependencies: [core]
services:
  - name: gateway
    kind: http
    port: 8080
    endpoints:
      - method: POST
        path: /auth/login
        action: login
architecture:
  pattern: hexagonal
generation:
  languages: [go, typescript]`,
    codeLang: "yaml",
  },
  {
    icon: "🌐",
    title: "Multi-Language Code Generation",
    description: "Generate idiomatic code in 5 languages with framework-specific adapters for each ecosystem.",
    details: [
      "Go — standard net/http handlers with clean architecture structure",
      "TypeScript — Express.js with type-safe routes and middleware",
      "Python — FastAPI with async endpoints and Pydantic schemas",
      "Java — Spring Boot with REST controllers and dependency injection",
      "Rust — Actix-Web or Axum 0.7 with type-safe handlers",
      "Framework-specific Dockerfile, docker-compose, and CI generation per language",
      "Consistent project structure (README, Makefile, .gitignore) across all languages",
    ],
    code: `# Generate for multiple targets
naeos run --config config.yaml --input spec.yaml \\
  --language go --language typescript \\
  --language python --language java \\
  --language rust`,
    codeLang: "bash",
  },
  {
    icon: "🤖",
    title: "AI Integration",
    description: "Compile NEIR models into AI instruction sets for 6 major coding assistants.",
    details: [
      "GitHub Copilot — .github/copilot-instructions.md with project context",
      "Claude Code — CLAUDE.md with architecture and coding standards",
      "Cursor — .cursorrules with framework-specific rules",
      "Gemini CLI — .gemini/CONFIG.md with project guidelines",
      "Codex — AGENTS.md with system-level instructions",
      "OpenCode — AGENTS.md with task-oriented directives",
      "MCP Server — Model Context Protocol for AI agent tool integration",
      "Context Bundles — LLM-optimized project summaries with dependency graphs",
    ],
    code: `# Compile to all AI adapters
naeos compile --all --input-file spec.yaml

# Or target specific adapters
naeos compile --copilot --claude --cursor \\
  --gemini --codex --opencode`,
    codeLang: "bash",
  },
  {
    icon: "🛒",
    title: "Marketplace",
    description: "An ecosystem for sharing and discovering industry profiles, plugins, and templates.",
    details: [
      "Profile Marketplace — publish, search, and download industry profiles",
      "Plugin Marketplace — install, uninstall, and search WASM-based plugins",
      "5 Built-in Profiles: SaaS, AI Agent, FinTech, Healthcare, Government",
      "+ 5 Additional: EdTech Platform, Ecommerce Engine, IoT Backend, Media Streaming, Blockchain Node",
      "Remote registry subscription for community profiles",
      "Plugin SDK with WASM runtime via wazero (Go WebAssembly interpreter)",
    ],
    code: `# Browse marketplace
naeos marketplace list
naeos marketplace search --query "ai-agent"
naeos marketplace install naeos-ai-agent --profile

# Manage plugins
naeos plugin install naeos-plugin-graphql
naeos plugin list`,
    codeLang: "bash",
  },
  {
    icon: "🛡️",
    title: "Governance & Security",
    description: "Built-in policy enforcement, artifact review, and security analysis for every pipeline run.",
    details: [
      "Policy Evaluator — 7 operators, 5 default rules for automated governance",
      "Artifact Review — quality gates for generated code before output",
      "Audit Trail — traceability with memory and file auditors (dual persistence)",
      "Security Analysis — SAST-style scanning with SARIF output",
      "RBAC — admin, developer, and viewer roles with full route permission mapping",
      "Encryption at Rest — AES-256-GCM for sensitive data storage",
      "Rate Limiting — configurable with X-RateLimit-* response headers",
    ],
    code: `# Run security audit
naeos security audit --config config.yaml

# Compliance export
naeos compliance export --format json

# Set secrets
naeos security set-secret --key DB_PASSWORD --value secret123`,
    codeLang: "bash",
  },
];

export const featuresID: FeatureData[] = [
  {
    icon: "⚡",
    title: "Pipeline Inti",
    description: "Pipeline end-to-end yang mengubah spesifikasi YAML/JSON menjadi artefak siap-produksi melalui 9 tahap.",
    details: [
      "Parser — parsing YAML/JSON dengan interpolasi variabel, referensi silang, dan include file",
      "Normalizer — normalisasi data dan koersi tipe",
      "Resolver — resolusi referensi silang dan import dengan pembatasan kedalaman",
      "NEIR Builder — model proyek terpadu (NAEOS Engineering Intermediate Representation)",
      "Validator — validasi komprehensif: dependensi sirkuler, konflik port, batasan modul",
      "Scheduler — penjadwalan tugas berbasis DAG dengan eksekusi paralel",
      "Generator — generasi kode multi-bahasa (Go, TypeScript, Python, Java, Rust)",
      "Renderer — rendering template dengan 10+ template bawaan",
      "Reviewer — evaluasi kebijakan dan review artefak",
    ],
    code: `# Konfigurasi pipeline
pipeline:
  name: proyek-saya
  mode: development
  verbose: true
  output_dir: ./out
  cache:
    enabled: true
    max_age: 5m`,
    codeLang: "yaml",
  },
  {
    icon: "📋",
    title: "Spec Language v2",
    description: "Bahasa spesifikasi yang dirancang untuk mendeskripsikan sistem perangkat lunak dalam berbagai skala.",
    details: [
      '${var} — interpolasi variabel dengan akses bersarang (contoh: ${modules.auth.name})',
      '$env{VAR} — resolusi variabel lingkungan',
      '$ref{path} — resolusi referensi silang di seluruh pohon spesifikasi',
      '$include{file} — komposisi spesifikasi multi-file',
      '$import{path} — fragmen spesifik modular dengan caching',
      '$fn{name(args)} — fungsi kustom: upper, lower, slug, default, len, coalesce',
      '$if{condition} / $endif — bagian kondisional',
      'Versioning skema dengan pengecekan otomatis (minimum v0.1.0)',
    ],
    code: `# Fitur Spec Language v2
project: aplikasi-saya
modules:
  - name: auth
    path: ./internal/auth
    description: Modul autentikasi
    dependencies: [core]
services:
  - name: gateway
    kind: http
    port: 8080
    endpoints:
      - method: POST
        path: /auth/login
        action: login
architecture:
  pattern: hexagonal
generation:
  languages: [go, typescript]`,
    codeLang: "yaml",
  },
  {
    icon: "🌐",
    title: "Generasi Kode Multi-Bahasa",
    description: "Hasilkan kode idiomatis dalam 5 bahasa dengan adapter khusus untuk setiap ekosistem framework.",
    details: [
      "Go — handler net/http standar dengan struktur arsitektur bersih",
      "TypeScript — Express.js dengan route type-safe dan middleware",
      "Python — FastAPI dengan endpoint async dan skema Pydantic",
      "Java — Spring Boot dengan REST controller dan dependency injection",
      "Rust — Actix-Web atau Axum 0.7 dengan handler type-safe",
      "Dockerfile, docker-compose, dan CI spesifik-framework per bahasa",
      "Struktur proyek konsisten (README, Makefile, .gitignore) di semua bahasa",
    ],
    code: `# Generate untuk banyak target
naeos run --config config.yaml --input spec.yaml \\
  --language go --language typescript \\
  --language python --language java \\
  --language rust`,
    codeLang: "bash",
  },
  {
    icon: "🤖",
    title: "Integrasi AI",
    description: "Kompilasi model NEIR menjadi set instruksi AI untuk 6 asisten coding utama.",
    details: [
      "GitHub Copilot — .github/copilot-instructions.md dengan konteks proyek",
      "Claude Code — CLAUDE.md dengan arsitektur dan standar coding",
      "Cursor — .cursorrules dengan aturan spesifik-framework",
      "Gemini CLI — .gemini/CONFIG.md dengan panduan proyek",
      "Codex — AGENTS.md dengan instruksi tingkat sistem",
      "OpenCode — AGENTS.md dengan direktif berbasis tugas",
      "MCP Server — Model Context Protocol untuk integrasi alat AI agent",
      "Context Bundles — ringkasan proyek yang dioptimalkan untuk LLM",
    ],
    code: `# Kompilasi ke semua adapter AI
naeos compile --all --input-file spec.yaml

# Atau target adapter spesifik
naeos compile --copilot --claude --cursor \\
  --gemini --codex --opencode`,
    codeLang: "bash",
  },
  {
    icon: "🛒",
    title: "Marketplace",
    description: "Ekosistem untuk berbagi dan menemukan profil industri, plugin, dan template.",
    details: [
      "Profile Marketplace — publikasi, pencarian, dan unduh profil industri",
      "Plugin Marketplace — instal, hapus, dan cari plugin berbasis WASM",
      "5 Profil Bawaan: SaaS, AI Agent, FinTech, Healthcare, Government",
      "+ 5 Tambahan: EdTech, Ecommerce Engine, IoT Backend, Media Streaming, Blockchain Node",
      "Langganan registry jarak jauh untuk profil komunitas",
      "Plugin SDK dengan runtime WASM via wazero (interpreter WebAssembly Go)",
    ],
    code: `# Jelajahi marketplace
naeos marketplace list
naeos marketplace search --query "ai-agent"
naeos marketplace install naeos-ai-agent --profile

# Kelola plugin
naeos plugin install naeos-plugin-graphql
naeos plugin list`,
    codeLang: "bash",
  },
  {
    icon: "🛡️",
    title: "Tata Kelola & Keamanan",
    description: "Penegakan kebijakan, review artefak, dan analisis keamanan terintegrasi untuk setiap eksekusi pipeline.",
    details: [
      "Policy Evaluator — 7 operator, 5 aturan default untuk tata kelola otomatis",
      "Artifact Review — gerbang kualitas untuk kode yang dihasilkan sebelum output",
      "Audit Trail — ketertelusuran dengan auditor memori dan file (persistensi ganda)",
      "Analisis Keamanan — pemindaian gaya SAST dengan output SARIF",
      "RBAC — peran admin, developer, dan viewer dengan pemetaan izin rute lengkap",
      "Enkripsi Saat Istirahat — AES-256-GCM untuk penyimpanan data sensitif",
      "Rate Limiting — dapat dikonfigurasi dengan header respons X-RateLimit-*",
    ],
    code: `# Jalankan audit keamanan
naeos security audit --config config.yaml

# Ekspor kepatuhan
naeos compliance export --format json

# Atur secrets
naeos security set-secret --key DB_PASSWORD --value secret123`,
    codeLang: "bash",
  },
];
