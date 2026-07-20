export interface CLIGroup {
  group: string;
  description: string;
  commands: CLICommand[];
}

export interface CLICommand {
  name: string;
  description: string;
  usage: string;
  flags?: { flag: string; short?: string; type: string; default: string; desc: string }[];
}

export const cliEN: CLIGroup[] = [
  {
    group: "Core",
    description: "Essential pipeline commands for daily use",
    commands: [
      {
        name: "run",
        description: "Execute the full NAEOS pipeline — parse, validate, generate, and render",
        usage: "naeos run --config config.yaml --input-file spec.yaml",
        flags: [
          { flag: "--config", type: "string", default: '""', desc: "Path to JSON or YAML config file" },
          { flag: "--input", type: "string", default: '""', desc: "Inline specification text" },
          { flag: "--input-file", type: "string", default: '""', desc: "Path to specification file" },
          { flag: "--output", type: "string", default: "text", desc: "Output format: text, json, yaml" },
          { flag: "--language", type: "string", default: '""', desc: "Target language (repeatable)" },
        ],
      },
      {
        name: "validate",
        description: "Validate a specification file for correctness",
        usage: "naeos validate --config config.yaml --input-file spec.yaml",
        flags: [
          { flag: "--config", type: "string", default: '""', desc: "Path to config file" },
          { flag: "--input", type: "string", default: '""', desc: "Inline specification text" },
          { flag: "--input-file", type: "string", default: '""', desc: "Path to spec file" },
          { flag: "--output", type: "string", default: "text", desc: "Output format: text, json, yaml" },
        ],
      },
      {
        name: "compile",
        description: "Compile NEIR into AI instruction sets for coding assistants",
        usage: "naeos compile --all --input-file spec.yaml",
        flags: [
          { flag: "--input-file", type: "string", default: '""', desc: "Path to spec file" },
          { flag: "--all", type: "bool", default: "false", desc: "Compile for all AI adapters" },
          { flag: "--copilot", type: "bool", default: "false", desc: "Compile for GitHub Copilot" },
          { flag: "--claude", type: "bool", default: "false", desc: "Compile for Claude Code" },
          { flag: "--cursor", type: "bool", default: "false", desc: "Compile for Cursor" },
          { flag: "--gemini", type: "bool", default: "false", desc: "Compile for Gemini CLI" },
          { flag: "--codex", type: "bool", default: "false", desc: "Compile for Codex" },
          { flag: "--opencode", type: "bool", default: "false", desc: "Compile for OpenCode" },
        ],
      },
      {
        name: "context",
        description: "Generate AI context bundles — LLM-optimized project summaries",
        usage: "naeos context --input-file spec.yaml",
      },
    ],
  },
  {
    group: "Development",
    description: "Development workflow commands",
    commands: [
      {
        name: "init",
        description: "Generate a default NAEOS config file",
        usage: "naeos init -o config.yaml",
        flags: [
          { flag: "--output", short: "-o", type: "string", default: "config.example.yaml", desc: "Output path" },
        ],
      },
      {
        name: "scaffold",
        description: "Generate a starter project scaffold with spec, Makefile, and .gitignore",
        usage: "naeos scaffold --name my-project --language go",
        flags: [
          { flag: "--name", type: "string", default: '""', desc: "Project name" },
          { flag: "--output", type: "string", default: '""', desc: "Output directory" },
          { flag: "--language", type: "string", default: '""', desc: "Target language (repeatable)" },
        ],
      },
      {
        name: "test",
        description: "Run tests for generated code across all target languages",
        usage: "naeos test --config config.yaml",
      },
      {
        name: "docgen",
        description: "Auto-generate API and module documentation from specs",
        usage: "naeos docgen --input-file spec.yaml",
      },
      {
        name: "diff",
        description: "Compare two specifications with colorized output",
        usage: "naeos diff --old spec-v1.yaml --new spec-v2.yaml",
      },
      {
        name: "watch",
        description: "Hot-reload pipeline on specification file changes",
        usage: "naeos watch --config config.yaml --input-file spec.yaml",
      },
    ],
  },
  {
    group: "Management",
    description: "Project and artifact management",
    commands: [
      {
        name: "marketplace",
        description: "Browse, search, and install profiles and plugins",
        usage: "naeos marketplace list",
      },
      {
        name: "profile",
        description: "Manage industry profiles (list, install, create)",
        usage: "naeos profile list",
      },
      {
        name: "plugin",
        description: "Manage WASM-based plugins (list, install, uninstall)",
        usage: "naeos plugin list",
      },
      {
        name: "artifacts",
        description: "Manage the artifact store (list, export, import)",
        usage: "naeos artifacts list",
      },
      {
        name: "template",
        description: "List and inspect prompt templates",
        usage: "naeos template list",
      },
      {
        name: "migrate",
        description: "Run schema migrations for specification versions",
        usage: "naeos migrate --from v0.1 --to v0.2",
      },
      {
        name: "rollback",
        description: "Rollback project changes to a previous state",
        usage: "naeos rollback --revision <id>",
      },
    ],
  },
  {
    group: "System",
    description: "System diagnostics and utilities",
    commands: [
      {
        name: "doctor",
        description: "Run system diagnostics and health checks",
        usage: "naeos doctor --config config.yaml",
      },
      {
        name: "mcp",
        description: "Start the Model Context Protocol server for AI agent integration",
        usage: "naeos mcp",
      },
      {
        name: "version",
        description: "Display NAEOS version information",
        usage: "naeos version",
      },
      {
        name: "completion",
        description: "Generate shell completion scripts (bash, zsh, fish)",
        usage: "naeos completion bash",
      },
      {
        name: "security",
        description: "Run security audits and manage secrets",
        usage: "naeos security audit --config config.yaml",
      },
      {
        name: "compliance",
        description: "Export compliance reports in JSON or CSV format",
        usage: "naeos compliance export --format json",
      },
      {
        name: "status",
        description: "Show pipeline and system status",
        usage: "naeos status",
      },
    ],
  },
];

export const cliID: CLIGroup[] = [
  {
    group: "Inti",
    description: "Perintah pipeline esensial untuk penggunaan sehari-hari",
    commands: [
      {
        name: "run",
        description: "Jalankan pipeline NAEOS lengkap — parse, validasi, generate, dan render",
        usage: "naeos run --config config.yaml --input-file spec.yaml",
        flags: [
          { flag: "--config", type: "string", default: '""', desc: "Path ke file konfigurasi JSON/YAML" },
          { flag: "--input", type: "string", default: '""', desc: "Teks spesifikasi inline" },
          { flag: "--input-file", type: "string", default: '""', desc: "Path ke file spesifikasi" },
          { flag: "--output", type: "string", default: "text", desc: "Format output: text, json, yaml" },
          { flag: "--language", type: "string", default: '""', desc: "Bahasa target (dapat diulang)" },
        ],
      },
      {
        name: "validate",
        description: "Validasi file spesifikasi untuk kebenaran",
        usage: "naeos validate --config config.yaml --input-file spec.yaml",
        flags: [
          { flag: "--config", type: "string", default: '""', desc: "Path ke file konfigurasi" },
          { flag: "--input", type: "string", default: '""', desc: "Teks spesifikasi inline" },
          { flag: "--input-file", type: "string", default: '""', desc: "Path ke file spesifikasi" },
          { flag: "--output", type: "string", default: "text", desc: "Format output: text, json, yaml" },
        ],
      },
      {
        name: "compile",
        description: "Kompilasi NEIR menjadi set instruksi AI untuk asisten coding",
        usage: "naeos compile --all --input-file spec.yaml",
        flags: [
          { flag: "--input-file", type: "string", default: '""', desc: "Path ke file spesifikasi" },
          { flag: "--all", type: "bool", default: "false", desc: "Kompilasi untuk semua adapter AI" },
        ],
      },
      {
        name: "context",
        description: "Hasilkan bundle konteks AI — ringkasan proyek yang dioptimalkan untuk LLM",
        usage: "naeos context --input-file spec.yaml",
      },
    ],
  },
  {
    group: "Pengembangan",
    description: "Perintah alur kerja pengembangan",
    commands: [
      {
        name: "init",
        description: "Hasilkan file konfigurasi NAEOS default",
        usage: "naeos init -o config.yaml",
      },
      {
        name: "scaffold",
        description: "Hasilkan scaffold proyek awal dengan spesifikasi, Makefile, dan .gitignore",
        usage: "naeos scaffold --name my-project --language go",
      },
      {
        name: "test",
        description: "Jalankan tes untuk kode yang dihasilkan di semua bahasa target",
        usage: "naeos test --config config.yaml",
      },
      {
        name: "docgen",
        description: "Hasilkan dokumentasi API dan modul secara otomatis dari spesifikasi",
        usage: "naeos docgen --input-file spec.yaml",
      },
      {
        name: "diff",
        description: "Bandingkan dua spesifikasi dengan output berwarna",
        usage: "naeos diff --old spec-v1.yaml --new spec-v2.yaml",
      },
      {
        name: "watch",
        description: "Pipeline hot-reload saat file spesifikasi berubah",
        usage: "naeos watch --config config.yaml --input-file spec.yaml",
      },
    ],
  },
  {
    group: "Manajemen",
    description: "Manajemen proyek dan artefak",
    commands: [
      {
        name: "marketplace",
        description: "Jelajahi, cari, dan instal profil serta plugin",
        usage: "naeos marketplace list",
      },
      {
        name: "profile",
        description: "Kelola profil industri (daftar, instal, buat)",
        usage: "naeos profile list",
      },
      {
        name: "plugin",
        description: "Kelola plugin berbasis WASM (daftar, instal, hapus)",
        usage: "naeos plugin list",
      },
      {
        name: "artifacts",
        description: "Kelola penyimpanan artefak (daftar, ekspor, impor)",
        usage: "naeos artifacts list",
      },
      {
        name: "template",
        description: "Lihat dan inspeksi template prompt",
        usage: "naeos template list",
      },
      {
        name: "migrate",
        description: "Jalankan migrasi skema untuk versi spesifikasi",
        usage: "naeos migrate --from v0.1 --to v0.2",
      },
      {
        name: "rollback",
        description: "Kembalikan perubahan proyek ke keadaan sebelumnya",
        usage: "naeos rollback --revision <id>",
      },
    ],
  },
  {
    group: "Sistem",
    description: "Diagnostik sistem dan utilitas",
    commands: [
      {
        name: "doctor",
        description: "Jalankan diagnostik sistem dan pemeriksaan kesehatan",
        usage: "naeos doctor --config config.yaml",
      },
      {
        name: "mcp",
        description: "Mulai server Model Context Protocol untuk integrasi agen AI",
        usage: "naeos mcp",
      },
      {
        name: "version",
        description: "Tampilkan informasi versi NAEOS",
        usage: "naeos version",
      },
      {
        name: "completion",
        description: "Hasilkan skrip completion shell (bash, zsh, fish)",
        usage: "naeos completion bash",
      },
      {
        name: "security",
        description: "Jalankan audit keamanan dan kelola secrets",
        usage: "naeos security audit --config config.yaml",
      },
      {
        name: "compliance",
        description: "Ekspor laporan kepatuhan dalam format JSON atau CSV",
        usage: "naeos compliance export --format json",
      },
      {
        name: "status",
        description: "Tampilkan status pipeline dan sistem",
        usage: "naeos status",
      },
    ],
  },
];
