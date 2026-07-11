# CLI Reference

Referensi lengkap untuk semua perintah NAEOS CLI.

## Global Flags

| Flag | Deskripsi | Default |
|------|-----------|---------|
| `--verbose` | Aktifkan logging verbose | `false` |
| `--dry-run` | Preview tanpa tulis ke disk | `false` |

## Perintah

### Pipeline

| Perintah | Deskripsi |
|----------|-----------|
| `naeos run` | Jalankan pipeline lengkap |
| `naeos validate` | Validasi spesifikasi |
| `naeos lint` | Lint spesifikasi |

### Code Generation

| Perintah | Deskripsi |
|----------|-----------|
| `naeos create` | Buat project baru dari spesifikasi |
| `naeos scaffold` | Generate scaffold kode |
| `naeos export` | Export artifacts |

### AI Integration

| Perintah | Deskripsi |
|----------|-----------|
| `naeos compile` | Compile ke AI instruction sets |
| `naeos context` | Generate AI context bundle |
| `naeos ai suggest` | Dapatkan saran AI |
| `naeos ai explain` | Jelaskan konsep NAEOS |

### Project Management

| Perintah | Deskripsi |
|----------|-----------|
| `naeos init` | Inisialisasi config file |
| `naeos status` | Lihat status pipeline |
| `naeos diff` | Bandingkan spesifikasi |
| `naeos watch` | Watch perubahan file |
| `naeos workspace` | Kelola workspace |

### Profiles & Artifacts

| Perintah | Deskripsi |
|----------|-----------|
| `naeos profile` | Kelola profil industri |
| `naeos artifacts` | Kelola artifact store |
| `naeos marketplace` | Marketplace profil |

### Governance

| Perintah | Deskripsi |
|----------|-----------|
| `naeos audit` | Audit spesifikasi |
| `naeos review` | Review artifacts |
| `naeos migrate` | Schema migration |

### System

| Perintah | Deskripsi |
|----------|-----------|
| `naeos doctor` | Cek kesehatan sistem |
| `naeos repair` | Repair spesifikasi |
| `naeos version` | Tampilkan versi |
| `naeos completion` | Generate shell completion |

### Advanced

| Perintah | Deskripsi |
|----------|-----------|
| `naeos kernel` | Inspeksi kernel |
| `naeos plugin` | Kelola plugins |
| `naeos template` | Kelola templates |
| `naeos lock` | Lock dependencies |
| `naeos rollback` | Rollback perubahan |

## Detail Perintah

### naeos run

```bash
naeos run [flags]

Flags:
  --config string       Path ke config file
  --input string        Spesifikasi inline
  --input-file string   Path ke file spesifikasi
  --output string       Format output (text, json, yaml)
  --output-file string  Path ke output file
  --language string     Target language (go, typescript, python, java, rust)
  --dry-run             Preview tanpa tulis ke disk
```

### naeos compile

```bash
naeos compile [flags]

Flags:
  --target string       Target tool (copilot, claude, cursor, gemini, codex, opencode)
  --all                 Compile ke semua target
  --input string        Spesifikasi inline
  --input-file string   Path ke file spesifikasi
  --output string       Directory untuk output files
```

### naeos context

```bash
naeos context [flags]

Flags:
  --input string        Spesifikasi inline
  --input-file string   Path ke file spesifikasi
  --output string       Format output (markdown, plain, json, yaml)
  --output-file string  Path ke output file
```

### naeos profile

```bash
naeos profile [command]

Commands:
  list                  Lihat semua profil
  show [name]           Lihat detail profil
  search [query]        Cari profil
  apply [name]          Apply profil ke spesifikasi
```

### naeos artifacts

```bash
naeos artifacts [command]

Commands:
  list                  Lihat semua artifacts
  info [path]           Lihat detail artifact
  dedup                 Deduplicate artifacts
  summary               Ringkasan artifacts
```

### naeos migrate

```bash
naeos migrate [command]

Commands:
  run                   Jalankan migration
  plan                  Rencana migration
  versions              Lihat versi tersedia
```

### naeos doctor

```bash
naeos doctor [flags]

Flags:
  --quick               Quick check saja
  --spec string         Path ke spesifikasi untuk divalidasi
```

### naeos diff

```bash
naeos diff [flags]

Flags:
  --input string        Spesifikasi baseline
  --input-file string   Path ke baseline file
  --target string       Spesifikasi target
  --target-file string  Path ke target file
```
