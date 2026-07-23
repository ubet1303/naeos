package compliance

import (
	"fmt"
	"sort"
	"time"
)

type Framework string

const (
	FrameworkSOC2     Framework = "SOC2"
	FrameworkHIPAA    Framework = "HIPAA"
	FrameworkISO27001 Framework = "ISO27001"
)

type ControlStatus string

const (
	StatusPass    ControlStatus = "pass"
	StatusFail    ControlStatus = "fail"
	StatusPartial ControlStatus = "partial"
	StatusNA      ControlStatus = "not_applicable"
)

type Control struct {
	ID          string        `json:"id"`
	Framework   Framework     `json:"framework"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      ControlStatus `json:"status"`
	Evidence    string        `json:"evidence,omitempty"`
	Notes       string        `json:"notes,omitempty"`
}

type Report struct {
	Framework   Framework `json:"framework"`
	GeneratedAt time.Time `json:"generated_at"`
	Controls    []Control `json:"controls"`
	Summary     Summary   `json:"summary"`
}

type Summary struct {
	Total         int     `json:"total"`
	Pass          int     `json:"pass"`
	Fail          int     `json:"fail"`
	Partial       int     `json:"partial"`
	NotApplicable int     `json:"not_applicable"`
	Score         float64 `json:"score"`
}

func summarize(controls []Control) Summary {
	s := Summary{Total: len(controls)}
	for _, c := range controls {
		switch c.Status {
		case StatusPass:
			s.Pass++
		case StatusFail:
			s.Fail++
		case StatusPartial:
			s.Partial++
		case StatusNA:
			s.NotApplicable++
		}
	}
	applicable := s.Total - s.NotApplicable
	if applicable > 0 {
		s.Score = float64(s.Pass) / float64(applicable) * 100
	}
	return s
}

// SOC2 Controls

func GenerateSOC2Report(checks map[string]ControlStatus) Report {
	controls := []Control{
		{ID: "CC6.1", Framework: FrameworkSOC2, Title: "Logical Access Controls", Description: "The entity implements logical access security measures"},
		{ID: "CC6.2", Framework: FrameworkSOC2, Title: "Authentication", Description: "Prior to issuing system credentials, the entity registers and authorizes new internal users"},
		{ID: "CC6.3", Framework: FrameworkSOC2, Title: "Access Removal", Description: "The entity authorizes, modifies, or removes access to data, software, functions, and other protected information assets based on roles"},
		{ID: "CC7.1", Framework: FrameworkSOC2, Title: "Vulnerability Management", Description: "To meet its objectives, the entity uses detection and monitoring procedures"},
		{ID: "CC7.2", Framework: FrameworkSOC2, Title: "Incident Response", Description: "The entity monitors system components and the operation of those components for anomalies"},
		{ID: "CC8.1", Framework: FrameworkSOC2, Title: "Change Management", Description: "The entity authorizes, designs, develops or acquires, configures, documents, tests, approves, and implements changes"},
		{ID: "A1.2", Framework: FrameworkSOC2, Title: "Environmental Protections", Description: "The entity authorizes, designs, develops or acquires, configures, documents, tests, approves, and implements environmental protections"},
		{ID: "C1.1", Framework: FrameworkSOC2, Title: "Data Classification", Description: "The entity classifies information and the types of data processed"},
	}

	return buildReport(FrameworkSOC2, controls, checks)
}

// HIPAA Controls

func GenerateHIPAAReport(checks map[string]ControlStatus) Report {
	controls := []Control{
		{ID: "§164.308(a)(1)", Framework: FrameworkHIPAA, Title: "Security Management Process", Description: "Implement policies and procedures to prevent, detect, contain, and correct security violations"},
		{ID: "§164.308(a)(3)", Framework: FrameworkHIPAA, Title: "Workforce Security", Description: "Implement policies and procedures to ensure that workforce members have appropriate access"},
		{ID: "§164.308(a)(4)", Framework: FrameworkHIPAA, Title: "Information Access Management", Description: "Implement policies and procedures for authorizing access to ePHI"},
		{ID: "§164.308(a)(5)", Framework: FrameworkHIPAA, Title: "Security Awareness Training", Description: "Implement security awareness and training program for all workforce members"},
		{ID: "§164.310(a)(1)", Framework: FrameworkHIPAA, Title: "Facility Access Controls", Description: "Implement policies and procedures to limit physical access to electronic information systems"},
		{ID: "§164.312(a)(1)", Framework: FrameworkHIPAA, Title: "Access Control", Description: "Implement technical policies and procedures for electronic information systems that maintain ePHI"},
		{ID: "§164.312(b)", Framework: FrameworkHIPAA, Title: "Audit Controls", Description: "Implement hardware, software, and/or procedural mechanisms that record and examine activity in systems"},
		{ID: "§164.312(c)(1)", Framework: FrameworkHIPAA, Title: "Integrity", Description: "Implement policies and procedures to protect ePHI from improper alteration or destruction"},
		{ID: "§164.312(d)", Framework: FrameworkHIPAA, Title: "Person or Entity Authentication", Description: "Implement procedures to verify that a person or entity seeking access to ePHI is the one claimed"},
		{ID: "§164.312(e)(1)", Framework: FrameworkHIPAA, Title: "Transmission Security", Description: "Implement technical security measures to guard against unauthorized access during electronic transmission"},
	}

	return buildReport(FrameworkHIPAA, controls, checks)
}

func buildReport(fw Framework, controls []Control, checks map[string]ControlStatus) Report {
	for i := range controls {
		if status, ok := checks[controls[i].ID]; ok {
			controls[i].Status = status
		} else {
			controls[i].Status = StatusPartial
			controls[i].Notes = "Not assessed"
		}
	}

	sort.Slice(controls, func(i, j int) bool {
		return controls[i].ID < controls[j].ID
	})

	return Report{
		Framework:   fw,
		GeneratedAt: time.Now().UTC(),
		Controls:    controls,
		Summary:     summarize(controls),
	}
}

// FormatReport renders a human-readable report.
func FormatReport(r Report) string {
	var out string
	out += fmt.Sprintf("=== %s Compliance Report ===\n", r.Framework)
	out += fmt.Sprintf("Generated: %s\n\n", r.GeneratedAt.Format(time.RFC3339))

	for _, c := range r.Controls {
		status := string(c.Status)
		out += fmt.Sprintf("[%s] %s — %s\n", status, c.ID, c.Title)
		if c.Notes != "" {
			out += fmt.Sprintf("  Note: %s\n", c.Notes)
		}
	}

	out += "\n--- Summary ---\n"
	out += fmt.Sprintf("Total: %d | Pass: %d | Fail: %d | Partial: %d | N/A: %d\n",
		r.Summary.Total, r.Summary.Pass, r.Summary.Fail, r.Summary.Partial, r.Summary.NotApplicable)
	out += fmt.Sprintf("Score: %.1f%%\n", r.Summary.Score)
	return out
}
