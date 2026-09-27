package knowledge

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Section headings of an opportunity record, in the order they are expected.
const (
	sectionSummary       = "Summary"
	sectionHypothesis    = "Hypothesis"
	sectionEvidence      = "Evidence"
	sectionUncertainties = "Uncertainties"
	sectionDecision      = "Decision"
)

var sections = []string{sectionSummary, sectionHypothesis, sectionEvidence, sectionUncertainties, sectionDecision}

var frontMatterKeys = []string{"id", "title", "state"}

// ParseOpportunity parses an opportunity record written as Markdown.
//
// The record starts with a front matter block holding id, title and state,
// followed by the level-two sections Summary, Hypothesis, Evidence,
// Uncertainties and Decision. Text sections are joined into one line.
// Evidence and Uncertainties are "- " bullet lists. Unknown keys, unknown
// sections and missing sections are errors. source names the record in errors.
func ParseOpportunity(source string, data []byte) (Opportunity, error) {
	o, err := parseOpportunity(source, data)
	if err != nil {
		return Opportunity{}, fmt.Errorf("%s: %w", source, err)
	}
	return o, nil
}

func parseOpportunity(source string, data []byte) (Opportunity, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	meta, body, err := splitFrontMatter(lines)
	if err != nil {
		return Opportunity{}, err
	}
	content, err := splitSections(body)
	if err != nil {
		return Opportunity{}, err
	}

	evidence, err := parseList(sectionEvidence, content[sectionEvidence])
	if err != nil {
		return Opportunity{}, err
	}
	uncertainties, err := parseList(sectionUncertainties, content[sectionUncertainties])
	if err != nil {
		return Opportunity{}, err
	}

	return NewOpportunity(Opportunity{
		ID:            ID(meta["id"]),
		Title:         meta["title"],
		Summary:       joinText(content[sectionSummary]),
		State:         State(meta["state"]),
		Hypothesis:    joinText(content[sectionHypothesis]),
		Evidence:      evidence,
		Uncertainties: uncertainties,
		Decision:      joinText(content[sectionDecision]),
		Source:        source,
	})
}

// splitFrontMatter reads the "---" delimited key: value block at the top.
func splitFrontMatter(lines []string) (map[string]string, []string, error) {
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil, nil, errors.New("front matter must start on the first line with ---")
	}
	meta := map[string]string{}
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "---" {
			for _, key := range frontMatterKeys {
				if _, ok := meta[key]; !ok {
					return nil, nil, fmt.Errorf("front matter: missing key %q", key)
				}
			}
			return meta, lines[i+1:], nil
		}
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, nil, fmt.Errorf("front matter line %d: expected key: value", i+1)
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !slices.Contains(frontMatterKeys, key) {
			return nil, nil, fmt.Errorf("front matter line %d: unknown key %q", i+1, key)
		}
		if _, dup := meta[key]; dup {
			return nil, nil, fmt.Errorf("front matter line %d: duplicate key %q", i+1, key)
		}
		meta[key] = value
	}
	return nil, nil, errors.New("front matter: missing closing ---")
}

// splitSections groups body lines under their "## " heading.
func splitSections(body []string) (map[string][]string, error) {
	content := map[string][]string{}
	current := ""
	for _, line := range body {
		if heading, ok := strings.CutPrefix(line, "## "); ok {
			heading = strings.TrimSpace(heading)
			if !slices.Contains(sections, heading) {
				return nil, fmt.Errorf("unknown section %q: expected one of %v", heading, sections)
			}
			if _, dup := content[heading]; dup {
				return nil, fmt.Errorf("duplicate section %q", heading)
			}
			current = heading
			content[current] = []string{}
			continue
		}
		if current == "" {
			if strings.TrimSpace(line) != "" {
				return nil, fmt.Errorf("text %q before first section", strings.TrimSpace(line))
			}
			continue
		}
		content[current] = append(content[current], line)
	}
	for _, s := range sections {
		if _, ok := content[s]; !ok {
			return nil, fmt.Errorf("missing section %q", s)
		}
	}
	return content, nil
}

func parseList(section string, lines []string) ([]string, error) {
	items := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		item, ok := strings.CutPrefix(line, "- ")
		if !ok {
			return nil, fmt.Errorf("section %q: expected \"- \" list item, got %q", section, line)
		}
		items = append(items, strings.TrimSpace(item))
	}
	return items, nil
}

func joinText(lines []string) string {
	return strings.Join(strings.Fields(strings.Join(lines, " ")), " ")
}
