# Core Principles

7 prinsip inti NAEOS yang memandu semua keputusan teknis dan arsitektural.

## 1. Specification is the Single Source of Truth

Spesifikasi mendefinisikan apa yang ada, apa yang valid, dan apa yang dibutuhkan. Kode, dokumentasi, dan keputusan engineering harus tersinkronisasi dengannya.

**Implikasi:**
- Semua perubahan dimulai dari spesifikasi
- Kode adalah hasil dari spesifikasi, bukan sebaliknya
- Dokumentasi harus konsisten dengan spesifikasi

## 2. Architecture Precedes Implementation

Keputusan arsitektur harus dibuat dan didokumentasikan sebelum kode ditulis.

**Implikasi:**
- Buat ADR (Architecture Decision Record) untuk setiap keputusan
- Dokumentasikan trade-offs
- Review arsitektur sebelum implementasi

## 3. Everything is Declarative

Semua konfigurasi dan perilaku didefinisikan secara deklaratif.

**Implikasi:**
- Spesifikasi YAML/JSON adalah satu-satunya cara mendefinisikan sistem
- Tidak ada imperative configuration
- Semua state bisa direkonstruksi dari spesifikasi

## 4. Everything is an Artifact

Semua output pipeline adalah artifacts yang bisa di-version, review, dan audit.

**Implikasi:**
- Kode, config, docs adalah artifacts
- Semua artifacts didokumentasikan
- Artifacts bisa di-cache dan di-share

## 5. Everything is Versioned

Semua komponen menggunakan semantic versioning.

**Implikasi:**
- Spesifikasi punya versi
- Artifacts punya versi
- Dependencies di-lock

## 6. Everything is Extensible

Sistem dirancang untuk diperluas tanpa modifikasi inti.

**Implikasi:**
- Plugin system untuk extension
- Custom adapters untuk output baru
- Custom profiles untuk industri baru

## 7. Small Kernel, Powerful Extensions

Kernel minimalis, semua fitur ada di extensions.

**Implikasi:**
- Kernel hanya handle core services
- Semua fitur ada di packages terpisah
- Extensions bisa di-install/uninstall

## 8. Human Review Before Automation

Otomatisasi harus didahului oleh review manusia.

**Implikasi:**
- Review artifacts sebelum deploy
- Policy evaluation untuk governance
- Audit trail untuk semua perubahan
