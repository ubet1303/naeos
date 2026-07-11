# Profiles

NAEOS menyediakan profil industri yang sudah dikonfigurasi sebelumnya untuk mempercepat setup proyek.

## Profil Tersedia

| Profil | Deskripsi |
|--------|-----------|
| **SaaS** | Aplikasi Software-as-a-Service dengan multi-tenant, billing, analytics |
| **AI Agent** | Sistem AI agent dengan LLM integration, tool calling, memory |
| **FinTech** | Sistem keuangan dengan compliance, audit trail, enkripsi |
| **Healthcare** | Sistem kesehatan dengan HIPAA compliance, audit logging |
| **Government** | Sistem pemerintahan dengan security clearance, audit |

## Penggunaan CLI

```bash
# Lihat semua profil
naeos profile list

# Lihat detail profil
naeos profile show saas

# Cari profil
naeos profile search "multi-tenant"

# Apply profil ke spesifikasi
naeos profile apply saas
```

## Penggunaan Go

```go
import "github.com/NAEOS-foundation/naeos/internal/profiles"

registry := profiles.NewRegistry()

// List semua profil
allProfiles := registry.List()
for _, p := range allProfiles {
    fmt.Printf("%s: %s\n", p.Name, p.Description)
}

// Cari profil
results := registry.Search("fintech")

// Apply profil ke spesifikasi
spec := registry.Apply("saas")
```

## Struktur Profil

```go
type Profile struct {
    Name        string
    Description string
    Industry    string
    Modules     []Module
    Services    []Service
    Architecture *Architecture
    Security    *Security
    Deployment  *Deployment
    Testing     *Testing
}
```

## Konversi ke Spec YAML

```go
profile := registry.Get("saas")
specYAML := profile.ToSpecYAML()
// Returns YAML string yang bisa langsung digunakan
```

## Custom Profil

Buat profil kustom dengan mendaftarkan ke registry:

```go
registry := profiles.NewRegistry()

customProfile := &profiles.Profile{
    Name:        "my-industry",
    Description: "Custom profile for my industry",
    Industry:    "manufacturing",
    Modules: []profiles.Module{
        {Name: "inventory", Path: "./inventory"},
        {Name: "supply-chain", Path: "./supply-chain"},
    },
}

registry.Register(customProfile)
```
