package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

var issueURLPattern = regexp.MustCompile(`/issues/(\d+)$`)
var slugPattern = regexp.MustCompile(`[^a-zA-Z0-9]+`)
var multiDashPattern = regexp.MustCompile(`-{2,}`)

func Main(args []string) int {
	if len(args) < 1 {
		return exitWithError(errors.New("usage: harness-state <paths|profile|adapter|missions|runs|signals|contracts|home|doctor|issue|claim|drive|ship|rc|meetings|teammates|maintenance> ..."))
	}

	var err error
	command := args[0]
	rest := args[1:]

	switch command {
	case "paths":
		err = cmdPaths(rest)
	case "profile":
		err = cmdProfile(rest)
	case "adapter":
		err = cmdAdapter(rest)
	case "missions":
		err = cmdMissions(rest)
	case "runs":
		err = cmdRuns(rest)
	case "signals":
		err = cmdSignals(rest)
	case "contracts":
		err = cmdContracts(rest)
	case "home":
		err = cmdHome(rest)
	case "doctor":
		err = cmdDoctor(rest)
	case "issue":
		err = cmdIssue(rest)
	case "claim":
		err = cmdClaim(rest)
	case "drive":
		err = cmdDrive(rest)
	case "ship":
		err = cmdShip(rest)
	case "rc":
		err = cmdRC(rest)
	case "meetings":
		err = cmdMeetings(rest)
	case "teammates":
		err = cmdTeammates(rest)
	case "maintenance":
		err = cmdMaintenance(rest)
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}

	return exitWithError(err)
}

func exitWithError(err error) int {
	if err == nil {
		return 0
	}
	fmt.Fprintln(os.Stderr, err.Error())
	return 1
}

func printJSON(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func nowTime() time.Time {
	return time.Now().UTC()
}

func parseISO(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			utc := parsed.UTC()
			return &utc, nil
		}
	}
	if strings.HasSuffix(value, "Z") {
		alt := strings.TrimSuffix(value, "Z") + "+00:00"
		for _, layout := range layouts {
			if parsed, err := time.Parse(layout, alt); err == nil {
				utc := parsed.UTC()
				return &utc, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid iso time: %s", value)
}

func repoRoot() string {
	if env := strings.TrimSpace(os.Getenv("HARNESS_REPO_ROOT")); env != "" {
		if abs, err := filepath.Abs(env); err == nil {
			return abs
		}
		return env
	}

	cwd, _ := os.Getwd()
	if out, err := runCommandAt(cwd, []string{"git", "rev-parse", "--show-toplevel"}); err == nil && strings.TrimSpace(out) != "" {
		return strings.TrimSpace(out)
	}

	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		if out, err := runCommandAt(dir, []string{"git", "rev-parse", "--show-toplevel"}); err == nil && strings.TrimSpace(out) != "" {
			return strings.TrimSpace(out)
		}
		return filepath.Dir(dir)
	}

	if cwd != "" {
		return cwd
	}
	return "."
}

func definitionsRoot() string {
	return filepath.Join(repoRoot(), "harness")
}

func runtimeRoot() string {
	return filepath.Join(repoRoot(), ".harness")
}

func stateLockPath() string {
	return filepath.Join(runtimeRoot(), ".lock")
}

func ensureRuntimeDirs() error {
	dirs := []string{
		filepath.Join(runtimeRoot(), "missions"),
		filepath.Join(runtimeRoot(), "runs"),
		filepath.Join(runtimeRoot(), "signals"),
		filepath.Join(runtimeRoot(), "contracts"),
		filepath.Join(runtimeRoot(), "meetings", "active"),
		filepath.Join(runtimeRoot(), "meetings", "archive"),
		filepath.Join(runtimeRoot(), "rc", "active"),
		filepath.Join(runtimeRoot(), "rc", "archive"),
		filepath.Join(runtimeRoot(), "scheduler", "snapshots"),
		filepath.Join(runtimeRoot(), "maintenance"),
		filepath.Join(runtimeRoot(), "cache"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func withHarnessLock(fn func() error) error {
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	lockPath := stateLockPath()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return err
	}
	handle, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer handle.Close()
	if err := syscall.Flock(int(handle.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(handle.Fd()), syscall.LOCK_UN)
	return fn()
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			result[key] = normalizeValue(item)
		}
		return result
	case map[any]any:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			result[fmt.Sprint(key)] = normalizeValue(item)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for index, item := range typed {
			result[index] = normalizeValue(item)
		}
		return result
	default:
		return typed
	}
}

func loadData(path string, defaultValue any) any {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return defaultValue
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return defaultValue
	}
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return defaultValue
	}

	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var jsonValue any
	if err := decoder.Decode(&jsonValue); err == nil {
		return normalizeValue(jsonValue)
	}

	var yamlValue any
	if err := yaml.Unmarshal(raw, &yamlValue); err == nil {
		return normalizeValue(yamlValue)
	}

	return defaultValue
}

func dumpData(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func runCommand(cmd []string) (string, error) {
	return runCommandAt(repoRoot(), cmd)
}

func runCommandAt(cwd string, cmd []string) (string, error) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Dir = cwd
	command.Stderr = io.Discard
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func runCommandWithStderr(cmd []string, cwd string) (string, error) {
	if cwd == "" {
		cwd = repoRoot()
	}
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Dir = cwd
	output, err := command.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		if trimmed != "" {
			return "", fmt.Errorf("%w: %s", err, trimmed)
		}
		return "", err
	}
	return trimmed, nil
}

func slugify(value string) string {
	slug := strings.TrimSpace(strings.ToLower(value))
	slug = slugPattern.ReplaceAllString(slug, "-")
	slug = multiDashPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 48 {
		slug = slug[:48]
		slug = strings.Trim(slug, "-")
	}
	if slug == "" {
		return "mission"
	}
	return slug
}

func stableSlugID(value string) string {
	slug := strings.TrimSpace(strings.ToLower(value))
	slug = slugPattern.ReplaceAllString(slug, "-")
	slug = multiDashPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "item"
	}
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(value))
	sum := fmt.Sprintf("%08x", hasher.Sum32())
	if len(slug) > 40 {
		slug = strings.Trim(slug[:40], "-")
	}
	return slug + "-" + sum
}

func branchFromPattern(pattern, issueNumber, slug string) string {
	result := strings.ReplaceAll(pattern, "{issue}", issueNumber)
	return strings.ReplaceAll(result, "{slug}", slug)
}

func cloneValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			result[key] = cloneValue(item)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for index, item := range typed {
			result[index] = cloneValue(item)
		}
		return result
	default:
		return typed
	}
}

func mergeDicts(base, override any) any {
	baseMap, baseOK := toMap(base)
	overrideMap, overrideOK := toMap(override)
	if baseOK && overrideOK {
		result := make(map[string]any, len(baseMap)+len(overrideMap))
		for key, value := range baseMap {
			result[key] = cloneValue(value)
		}
		for key, value := range overrideMap {
			if existing, ok := result[key]; ok {
				result[key] = mergeDicts(existing, value)
			} else {
				result[key] = cloneValue(value)
			}
		}
		return result
	}
	return cloneValue(override)
}

func parsePipeList(raw string) []string {
	parts := strings.Split(raw, "|")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			items = append(items, value)
		}
	}
	return items
}

func isBlank(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return typed == ""
	case bool:
		return !typed
	case int:
		return typed == 0
	case int64:
		return typed == 0
	case float64:
		return typed == 0
	case json.Number:
		return typed == "0"
	case []any:
		return len(typed) == 0
	case map[string]any:
		return len(typed) == 0
	default:
		return false
	}
}

func fillIfBlank(mapping map[string]any, key string, value any) {
	if isBlank(mapping[key]) {
		mapping[key] = value
	}
}

func getKey(data any, key string, defaultValue any) any {
	current := data
	for _, part := range strings.Split(key, ".") {
		mapping, ok := toMap(current)
		if !ok {
			return defaultValue
		}
		next, ok := mapping[part]
		if !ok {
			return defaultValue
		}
		current = next
	}
	if current == nil {
		return defaultValue
	}
	return current
}

func readStdinJSON() (map[string]any, error) {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(raw)) == "" {
		return map[string]any{}, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	result := map[string]any{}
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return normalizeValue(result).(map[string]any), nil
}

func toMap(value any) (map[string]any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		return typed, true
	default:
		return nil, false
	}
}

func toList(value any) []any {
	switch typed := value.(type) {
	case nil:
		return []any{}
	case []any:
		return typed
	default:
		return []any{}
	}
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case json.Number:
		return typed.String()
	default:
		return fmt.Sprint(typed)
	}
}

func intValue(value any) int {
	switch typed := value.(type) {
	case nil:
		return 0
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		parsed, _ := typed.Int64()
		return int(parsed)
	case string:
		parsed, _ := strconv.Atoi(strings.TrimSpace(typed))
		return parsed
	default:
		return 0
	}
}

func floatValue(value any) float64 {
	switch typed := value.(type) {
	case nil:
		return 0
	case float64:
		return typed
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case json.Number:
		parsed, _ := typed.Float64()
		return parsed
	case string:
		parsed, _ := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return parsed
	default:
		return 0
	}
}

func boolValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		return typed == "true"
	default:
		return false
	}
}

func stringSlice(value any) []string {
	items := toList(value)
	result := make([]string, 0, len(items))
	for _, item := range items {
		text := strings.TrimSpace(stringValue(item))
		if text != "" {
			result = append(result, text)
		}
	}
	return result
}

func profilePaths(name string) []string {
	paths := []string{}
	localOverride := filepath.Join(runtimeRoot(), "profile.local.yml")
	if exists(localOverride) {
		paths = append(paths, localOverride)
	}
	profilesDir := filepath.Join(definitionsRoot(), "profiles")
	if name != "" {
		candidate := filepath.Join(profilesDir, name+".yml")
		if exists(candidate) {
			paths = append(paths, candidate)
		}
		return paths
	}

	entries, _ := filepath.Glob(filepath.Join(profilesDir, "*.yml"))
	sort.Strings(entries)
	if len(entries) == 1 {
		paths = append(paths, entries[0])
	}
	defaultHint := filepath.Join(profilesDir, "default.yml")
	if exists(defaultHint) {
		paths = append(paths, defaultHint)
	}
	return paths
}

func resolveProfile(name string) (map[string]any, string) {
	profilesDir := filepath.Join(definitionsRoot(), "profiles")
	basePath := ""
	if name != "" {
		explicit := filepath.Join(profilesDir, name+".yml")
		if exists(explicit) {
			basePath = explicit
		}
	} else {
		entries, _ := filepath.Glob(filepath.Join(profilesDir, "*.yml"))
		sort.Strings(entries)
		if len(entries) == 1 {
			basePath = entries[0]
		} else {
			defaultHint := filepath.Join(profilesDir, "default.yml")
			if exists(defaultHint) {
				basePath = defaultHint
			}
		}
	}

	base := map[string]any{}
	if basePath != "" {
		if loaded, ok := toMap(loadData(basePath, map[string]any{})); ok {
			base = loaded
		}
	}
	localOverride, _ := toMap(loadData(filepath.Join(runtimeRoot(), "profile.local.yml"), map[string]any{}))
	merged, _ := toMap(mergeDicts(base, localOverride))
	return merged, basePath
}

func adapterManifestPath() string {
	return filepath.Join(definitionsRoot(), "adapter-manifest.yml")
}

func adapterManifest() map[string]any {
	value, _ := toMap(loadData(adapterManifestPath(), map[string]any{}))
	return value
}

func adapterProtocol() map[string]any {
	value, _ := toMap(adapterManifest()["protocol"])
	if value == nil {
		return map[string]any{}
	}
	return value
}

func adapterConfig() map[string]any {
	value, _ := toMap(adapterManifest()["adapter"])
	if value == nil {
		return map[string]any{}
	}
	return value
}

func manifestList(value any) []string {
	return stringSlice(value)
}

func adapterPackageFiles() []string {
	return manifestList(adapterConfig()["package_files"])
}

func adapterSmokeScripts() []string {
	return manifestList(adapterConfig()["smoke_scripts"])
}

func adapterBootstrapConfig() map[string]any {
	value, _ := toMap(adapterConfig()["bootstrap"])
	if value == nil {
		return map[string]any{}
	}
	return value
}

func adapterProfileContract() map[string]any {
	value, _ := toMap(adapterConfig()["profile_contract"])
	if value == nil {
		return map[string]any{}
	}
	return value
}

func repoRelativeString(path string) string {
	rel, err := filepath.Rel(repoRoot(), path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return path
	}
	return rel
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func adapterValidate() map[string]any {
	manifestPath := adapterManifestPath()
	protocol := adapterProtocol()
	adapter := adapterConfig()
	bootstrap := adapterBootstrapConfig()
	profileContract := adapterProfileContract()
	packageFiles := adapterPackageFiles()
	smokeScripts := adapterSmokeScripts()
	profile, basePath := resolveProfile("")
	profileOverridePath := filepath.Join(runtimeRoot(), "profile.local.yml")
	issues := []any{}

	if !exists(manifestPath) {
		issues = append(issues, issueRecord("critical", "adapter-manifest-missing", repoRelativeString(manifestPath), "Adapter manifest is missing.", "Restore harness/adapter-manifest.yml before using bootstrap, sync, or smoke flows."))
	}
	if stringValue(protocol["protocol_id"]) == "" {
		issues = append(issues, issueRecord("warning", "protocol-id-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare protocol_id.", "Declare protocol identity so protocol, adapter, and profile layers stay explicit."))
	}
	if stringValue(protocol["protocol_version"]) == "" {
		issues = append(issues, issueRecord("warning", "protocol-version-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare protocol_version.", "Declare protocol_version to anchor runtime behavior to a concrete protocol snapshot."))
	}
	if stringValue(adapter["adapter_id"]) == "" {
		issues = append(issues, issueRecord("critical", "adapter-id-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare adapter_id.", "Set adapter.adapter_id so sync/bootstrap metadata can identify the package unambiguously."))
	}
	if stringValue(adapter["adapter_version"]) == "" {
		issues = append(issues, issueRecord("critical", "adapter-version-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare adapter_version.", "Set adapter.adapter_version so bootstrap and sync can track adapter drift."))
	}
	if stringValue(adapter["runtime"]) == "" {
		issues = append(issues, issueRecord("critical", "adapter-runtime-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare runtime.", "Set adapter.runtime so the package contract stays explicit about its target runtime."))
	}
	if len(packageFiles) == 0 {
		issues = append(issues, issueRecord("critical", "adapter-package-empty", repoRelativeString(manifestPath), "Adapter manifest does not declare package_files.", "List the reusable adapter files in adapter.package_files so bootstrap and sync can stay manifest-driven."))
	}
	for _, relPath := range packageFiles {
		if !exists(filepath.Join(repoRoot(), relPath)) {
			issues = append(issues, issueRecord("critical", "adapter-package-file-missing", relPath, fmt.Sprintf("Manifest-declared package file is missing: %s", relPath), "Restore the file or remove it from adapter.package_files."))
		}
	}
	if len(smokeScripts) == 0 {
		issues = append(issues, issueRecord("warning", "adapter-smoke-empty", repoRelativeString(manifestPath), "Adapter manifest does not declare smoke scripts.", "Declare smoke entrypoints so the reusable package documents its verification contract."))
	}
	for _, relPath := range smokeScripts {
		normalized := relPath
		if strings.HasPrefix(normalized, "./") {
			normalized = strings.TrimPrefix(normalized, "./")
		}
		if !exists(filepath.Join(repoRoot(), normalized)) {
			issues = append(issues, issueRecord("critical", "adapter-smoke-missing", relPath, fmt.Sprintf("Manifest-declared smoke script is missing: %s", relPath), "Restore the script or update adapter.smoke_scripts."))
		}
	}
	publicContractPath := strings.TrimSpace(stringValue(adapter["public_contract_path"]))
	if publicContractPath == "" {
		issues = append(issues, issueRecord("critical", "adapter-public-contract-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare public_contract_path.", "Set adapter.public_contract_path so the package boundary remains machine-addressable."))
	} else if !exists(filepath.Join(repoRoot(), publicContractPath)) {
		issues = append(issues, issueRecord("critical", "adapter-public-contract-file-missing", publicContractPath, fmt.Sprintf("Public contract file is missing: %s", publicContractPath), "Restore the file or update adapter.public_contract_path."))
	}
	releaseCriteriaPath := strings.TrimSpace(stringValue(adapter["release_criteria_path"]))
	if releaseCriteriaPath == "" {
		issues = append(issues, issueRecord("warning", "adapter-release-criteria-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare release_criteria_path.", "Set adapter.release_criteria_path so the 1.0.0 freeze bar remains machine-addressable."))
	} else if !exists(filepath.Join(repoRoot(), releaseCriteriaPath)) {
		issues = append(issues, issueRecord("warning", "adapter-release-criteria-file-missing", releaseCriteriaPath, fmt.Sprintf("Release criteria file is missing: %s", releaseCriteriaPath), "Restore the file or update adapter.release_criteria_path."))
	}

	profileTemplate := strings.TrimSpace(stringValue(bootstrap["profile_template"]))
	generatedProfilePath := strings.TrimSpace(stringValue(bootstrap["generated_profile_path"]))
	metadataPath := strings.TrimSpace(stringValue(bootstrap["metadata_path"]))
	if profileTemplate == "" {
		issues = append(issues, issueRecord("critical", "bootstrap-profile-template-missing", repoRelativeString(manifestPath), "Adapter bootstrap contract does not declare a profile template path.", "Set adapter.bootstrap.profile_template so bootstrap can generate a repo-local profile deterministically."))
	} else if !exists(filepath.Join(repoRoot(), profileTemplate)) {
		issues = append(issues, issueRecord("critical", "bootstrap-profile-template-file-missing", profileTemplate, fmt.Sprintf("Bootstrap profile template is missing: %s", profileTemplate), "Restore the file or update adapter.bootstrap.profile_template."))
	}
	if generatedProfilePath == "" {
		issues = append(issues, issueRecord("critical", "bootstrap-generated-profile-path-missing", repoRelativeString(manifestPath), "Adapter bootstrap contract does not declare generated_profile_path.", "Set adapter.bootstrap.generated_profile_path so bootstrap knows where to place the generated repo profile."))
	}
	if metadataPath == "" {
		issues = append(issues, issueRecord("warning", "bootstrap-metadata-path-missing", repoRelativeString(manifestPath), "Adapter bootstrap contract does not declare metadata_path.", "Set adapter.bootstrap.metadata_path so sync has a stable remembered-source contract."))
	}

	committedProfilesRoot := strings.TrimSpace(stringValue(profileContract["committed_profiles_root"]))
	runtimeOverridePath := strings.TrimSpace(stringValue(profileContract["runtime_local_override_path"]))
	defaultProfileHint := strings.TrimSpace(stringValue(profileContract["default_profile_hint"]))
	if committedProfilesRoot == "" {
		issues = append(issues, issueRecord("warning", "profile-contract-root-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare committed_profiles_root.", "Declare the committed profile root so the adapter/profile split is machine-readable."))
	} else if !exists(filepath.Join(repoRoot(), committedProfilesRoot)) {
		issues = append(issues, issueRecord("warning", "profile-contract-root-missing-on-disk", committedProfilesRoot, fmt.Sprintf("Committed profiles root is missing: %s", committedProfilesRoot), "Restore the directory or update adapter.profile_contract.committed_profiles_root."))
	}
	if runtimeOverridePath == "" {
		issues = append(issues, issueRecord("warning", "profile-contract-override-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare runtime_local_override_path.", "Declare the runtime-local override path so live state stays separate from committed profile definitions."))
	}
	if defaultProfileHint == "" {
		issues = append(issues, issueRecord("warning", "profile-contract-default-hint-missing", repoRelativeString(manifestPath), "Adapter manifest does not declare default_profile_hint.", "Declare the default profile hint so bootstrapped repos expose a stable committed profile location."))
	}
	if basePath == "" && !exists(profileOverridePath) {
		ref := committedProfilesRoot
		if ref == "" {
			ref = "harness/profiles"
		}
		issues = append(issues, issueRecord("warning", "profile-unresolved", ref, "No committed or runtime-local profile could be resolved.", "Add a committed profile under harness/profiles or a runtime override under .harness/profile.local.yml."))
	}

	sortIssueList(issues)

	return map[string]any{
		"manifest_path":               repoRelativeString(manifestPath),
		"adapter_id":                  stringValue(adapter["adapter_id"]),
		"adapter_version":             stringValue(adapter["adapter_version"]),
		"runtime":                     stringValue(adapter["runtime"]),
		"public_contract_path":        publicContractPath,
		"release_criteria_path":       releaseCriteriaPath,
		"protocol_id":                 stringValue(protocol["protocol_id"]),
		"protocol_version":            stringValue(protocol["protocol_version"]),
		"package_file_count":          len(packageFiles),
		"smoke_script_count":          len(smokeScripts),
		"profile_template":            profileTemplate,
		"generated_profile_path":      generatedProfilePath,
		"metadata_path":               metadataPath,
		"committed_profiles_root":     committedProfilesRoot,
		"runtime_local_override_path": runtimeOverridePath,
		"default_profile_hint":        defaultProfileHint,
		"resolved_profile_path":       repoRelativeString(basePath),
		"resolved_profile_id":         stringValue(getKey(profile, "profile.profile_id", "")),
		"issues":                      issues,
		"valid":                       len(issues) == 0,
	}
}

func sortIssueList(issues []any) {
	order := map[string]int{"critical": 0, "warning": 1, "info": 2}
	sort.Slice(issues, func(i, j int) bool {
		left, _ := toMap(issues[i])
		right, _ := toMap(issues[j])
		leftSeverity := order[stringValue(left["severity"])]
		rightSeverity := order[stringValue(right["severity"])]
		if leftSeverity != rightSeverity {
			return leftSeverity < rightSeverity
		}
		leftKind := stringValue(left["kind"])
		rightKind := stringValue(right["kind"])
		if leftKind != rightKind {
			return leftKind < rightKind
		}
		return stringValue(left["ref"]) < stringValue(right["ref"])
	})
}

func registryPath() string {
	return filepath.Join(runtimeRoot(), "runs", "registry.yml")
}

func missionsRoot() string {
	return filepath.Join(runtimeRoot(), "missions")
}

func missionRegistryPath() string {
	return filepath.Join(missionsRoot(), "registry.yml")
}

func missionPath(missionID string) string {
	return filepath.Join(missionsRoot(), missionID+".yml")
}

func signalLogPath() string {
	return filepath.Join(runtimeRoot(), "signals", "log.yml")
}

func missionTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "mission.yml")
}

func contractTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "contract.yml")
}

func runTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "run.yml")
}

func signalTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "signal.yml")
}

func meetingTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "meeting.yml")
}

func rcTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "rc.yml")
}

func issueTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "issue.md")
}

func prTemplatePath() string {
	return filepath.Join(definitionsRoot(), "templates", "pr.md")
}

func policiesPath() string {
	return filepath.Join(definitionsRoot(), "policies", "defaults.yml")
}

func maintenanceBacklogPath() string {
	return filepath.Join(runtimeRoot(), "maintenance", "backlog.yml")
}

func meetingsActiveRoot() string {
	return filepath.Join(runtimeRoot(), "meetings", "active")
}

func meetingsArchiveRoot() string {
	return filepath.Join(runtimeRoot(), "meetings", "archive")
}

func rcActiveRoot() string {
	return filepath.Join(runtimeRoot(), "rc", "active")
}

func rcArchiveRoot() string {
	return filepath.Join(runtimeRoot(), "rc", "archive")
}

func meetingPath(meetingID string, archived bool) string {
	root := meetingsActiveRoot()
	if archived {
		root = meetingsArchiveRoot()
	}
	return filepath.Join(root, meetingID+".yml")
}

func rcPath(rcID string, archived bool) string {
	root := rcActiveRoot()
	if archived {
		root = rcArchiveRoot()
	}
	return filepath.Join(root, rcID+".yml")
}

func newMeetingID(missionID string) string {
	timestamp := nowTime().UTC()
	return fmt.Sprintf("meeting-%s-%09d-%s", timestamp.Format("20060102150405"), timestamp.Nanosecond(), missionID)
}

func readRegistry() map[string]any {
	value, _ := toMap(loadData(registryPath(), map[string]any{
		"version":    1,
		"updated_at": nowISO(),
		"runs":       []any{},
	}))
	return value
}

func writeRegistry(data map[string]any) error {
	data["updated_at"] = nowISO()
	if err := dumpData(registryPath(), data); err != nil {
		return err
	}
	return syncMissionsFromRegistry(data)
}

func readMissionRegistry() map[string]any {
	value, _ := toMap(loadData(missionRegistryPath(), map[string]any{
		"version":    1,
		"updated_at": nowISO(),
		"missions":   []any{},
	}))
	return value
}

func writeMissionRegistry(data map[string]any) error {
	data["updated_at"] = nowISO()
	return dumpData(missionRegistryPath(), data)
}

func readSignalLog() map[string]any {
	value, _ := toMap(loadData(signalLogPath(), map[string]any{
		"version":    1,
		"updated_at": nowISO(),
		"signals":    []any{},
	}))
	return value
}

func writeSignalLog(data map[string]any) error {
	data["updated_at"] = nowISO()
	return dumpData(signalLogPath(), data)
}

func readMaintenanceBacklog() map[string]any {
	value, _ := toMap(loadData(maintenanceBacklogPath(), map[string]any{
		"version":    1,
		"updated_at": nowISO(),
		"meetings":   []any{},
		"tasks":      []any{},
	}))
	return value
}

func writeMaintenanceBacklog(data map[string]any) error {
	data["updated_at"] = nowISO()
	return dumpData(maintenanceBacklogPath(), data)
}

func readMeeting(meetingID string) map[string]any {
	activePath := meetingPath(meetingID, false)
	if exists(activePath) {
		value, _ := toMap(loadData(activePath, map[string]any{}))
		return value
	}
	archivedPath := meetingPath(meetingID, true)
	if exists(archivedPath) {
		value, _ := toMap(loadData(archivedPath, map[string]any{}))
		return value
	}
	return nil
}

func readMission(missionID string) map[string]any {
	path := missionPath(missionID)
	if !exists(path) {
		return nil
	}
	value, _ := toMap(loadData(path, map[string]any{}))
	mission, _ := toMap(value["mission"])
	if len(mission) == 0 {
		return nil
	}
	return mission
}

func listMissions() []any {
	registry := readMissionRegistry()
	missions := []any{}
	for _, item := range toList(registry["missions"]) {
		mission, _ := toMap(item)
		if len(mission) > 0 {
			missions = append(missions, mission)
		}
	}
	sort.Slice(missions, func(i, j int) bool {
		left, _ := toMap(missions[i])
		right, _ := toMap(missions[j])
		leftActive := activeMission(left)
		rightActive := activeMission(right)
		if leftActive != rightActive {
			return leftActive
		}
		leftUpdated := firstNonBlank(stringValue(left["updated_at"]), stringValue(left["created_at"]))
		rightUpdated := firstNonBlank(stringValue(right["updated_at"]), stringValue(right["created_at"]))
		if leftUpdated != rightUpdated {
			return leftUpdated > rightUpdated
		}
		return stringValue(left["mission_id"]) < stringValue(right["mission_id"])
	})
	return missions
}

func readActiveMeeting(meetingID string) map[string]any {
	activePath := meetingPath(meetingID, false)
	if !exists(activePath) {
		return nil
	}
	value, _ := toMap(loadData(activePath, map[string]any{}))
	return value
}

func readRC(rcID string) map[string]any {
	activePath := rcPath(rcID, false)
	if exists(activePath) {
		value, _ := toMap(loadData(activePath, map[string]any{}))
		return value
	}
	archivedPath := rcPath(rcID, true)
	if exists(archivedPath) {
		value, _ := toMap(loadData(archivedPath, map[string]any{}))
		return value
	}
	return nil
}

func readActiveRC(rcID string) map[string]any {
	activePath := rcPath(rcID, false)
	if !exists(activePath) {
		return nil
	}
	value, _ := toMap(loadData(activePath, map[string]any{}))
	return value
}

func listMeetings(includeArchived bool) []any {
	_ = ensureRuntimeDirs()
	paths, _ := filepath.Glob(filepath.Join(meetingsActiveRoot(), "*.yml"))
	sort.Strings(paths)
	if includeArchived {
		archived, _ := filepath.Glob(filepath.Join(meetingsArchiveRoot(), "*.yml"))
		sort.Strings(archived)
		paths = append(paths, archived...)
	}
	meetings := []any{}
	for _, path := range paths {
		data, _ := toMap(loadData(path, map[string]any{}))
		meeting, _ := toMap(data["meeting"])
		if len(meeting) > 0 {
			meetings = append(meetings, meeting)
		}
	}
	sort.Slice(meetings, func(i, j int) bool {
		left, _ := toMap(meetings[i])
		right, _ := toMap(meetings[j])
		return stringValue(left["created_at"]) > stringValue(right["created_at"])
	})
	return meetings
}

func listRCs(includeArchived bool) []any {
	_ = ensureRuntimeDirs()
	paths, _ := filepath.Glob(filepath.Join(rcActiveRoot(), "*.yml"))
	sort.Strings(paths)
	if includeArchived {
		archived, _ := filepath.Glob(filepath.Join(rcArchiveRoot(), "*.yml"))
		sort.Strings(archived)
		paths = append(paths, archived...)
	}
	rcs := []any{}
	for _, path := range paths {
		data, _ := toMap(loadData(path, map[string]any{}))
		rc, _ := toMap(data["rc"])
		if len(rc) > 0 {
			rcs = append(rcs, rc)
		}
	}
	sort.Slice(rcs, func(i, j int) bool {
		left, _ := toMap(rcs[i])
		right, _ := toMap(rcs[j])
		return stringValue(left["created_at"]) > stringValue(right["created_at"])
	})
	return rcs
}

func activeRC() map[string]any {
	for _, item := range listRCs(false) {
		rc, _ := toMap(item)
		if stringValue(rc["state"]) == "active" {
			return rc
		}
	}
	return nil
}

func latestClosedMeetingForRun(runID string) map[string]any {
	for _, item := range listMeetings(true) {
		meeting, _ := toMap(item)
		if stringValue(meeting["run_id"]) == runID && stringValue(meeting["state"]) == "closed" {
			return meeting
		}
	}
	return nil
}

func governingClosedMeetingForRun(run map[string]any) map[string]any {
	runID := strings.TrimSpace(stringValue(run["run_id"]))
	if adoptedMeetingID := strings.TrimSpace(stringValue(run["adopted_from_meeting_id"])); adoptedMeetingID != "" {
		if data := readMeeting(adoptedMeetingID); data != nil {
			meeting, _ := toMap(data["meeting"])
			if len(meeting) > 0 && stringValue(meeting["run_id"]) == runID && stringValue(meeting["state"]) == "closed" {
				return meeting
			}
		}
	}
	return latestClosedMeetingForRun(runID)
}

func meetingActionCountsAsBacklog(meeting, registry map[string]any) bool {
	if stringValue(meeting["state"]) == "active" {
		return true
	}
	if stringValue(meeting["state"]) != "closed" || !meetingHasOutput(meeting) {
		return true
	}
	runID := strings.TrimSpace(stringValue(meeting["run_id"]))
	if runID == "" {
		return true
	}
	index := findRunIndex(registry, runID)
	if index < 0 {
		return true
	}
	run, _ := toMap(toList(registry["runs"])[index])
	return strings.TrimSpace(stringValue(run["adopted_from_meeting_id"])) != stringValue(meeting["meeting_id"])
}

func openMeetingActions() []any {
	actions := []any{}
	registry := readRegistry()
	for _, item := range listMeetings(true) {
		meeting, _ := toMap(item)
		if !meetingActionCountsAsBacklog(meeting, registry) {
			continue
		}
		for _, actionItem := range toList(meeting["action_items"]) {
			actionMap, isMap := toMap(actionItem)
			if isMap && stringValue(actionMap["state"]) != "" && stringValue(actionMap["state"]) != "open" {
				continue
			}
			summary := textValue(actionItem)
			if summary == "" {
				continue
			}
			entry := map[string]any{
				"meeting_id":    stringValue(meeting["meeting_id"]),
				"meeting_state": stringValue(meeting["state"]),
				"run_id":        stringValue(meeting["run_id"]),
				"mission_id":    stringValue(meeting["mission_id"]),
				"summary":       summary,
				"action_id":     "",
			}
			if isMap {
				entry["action_id"] = stringValue(actionMap["action_id"])
			}
			actions = append(actions, entry)
		}
	}
	return actions
}

func findRunIndex(registry map[string]any, runID string) int {
	for index, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if stringValue(run["run_id"]) == runID {
			return index
		}
	}
	return -1
}

func findSignalIndex(log map[string]any, signalID string) int {
	for index, item := range toList(log["signals"]) {
		signal, _ := toMap(item)
		if stringValue(signal["signal_id"]) == signalID {
			return index
		}
	}
	return -1
}

func appendUniqueText(list []any, value string) []any {
	if strings.TrimSpace(value) == "" {
		return list
	}
	for _, item := range list {
		if stringValue(item) == value {
			return list
		}
	}
	return append(list, value)
}

func dedupeStrings(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	seen := map[string]bool{}
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		result = append(result, trimmed)
	}
	return result
}

func containsString(items []string, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}
	for _, item := range items {
		if strings.TrimSpace(item) == target {
			return true
		}
	}
	return false
}

func profileHighRiskScopeSelector(raw string) string {
	scope := strings.TrimSpace(raw)
	if scope == "" {
		return ""
	}
	scope = strings.TrimPrefix(scope, "./")
	if strings.ContainsAny(scope, "*?[") {
		return "glob:" + strings.TrimSpace(strings.TrimSuffix(scope, "/"))
	}
	if strings.HasSuffix(scope, "/") {
		return "dir:" + strings.TrimSpace(strings.TrimSuffix(scope, "/"))
	}
	if strings.Contains(scope, "/") || strings.Contains(scope, ".") {
		return "file:" + strings.TrimSpace(scope)
	}
	return ""
}

func matchedProfileHighRiskScopes(profile map[string]any, ownership []string, releaseAdjacent bool) []string {
	matches := []string{}
	for _, item := range anyListToStrings(toList(getKey(profile, "profile.high_risk_scopes", []any{}))) {
		scope := strings.TrimSpace(item)
		if scope == "" {
			continue
		}
		lower := strings.ToLower(scope)
		if strings.Contains(lower, "release") && strings.Contains(lower, "rc") {
			if releaseAdjacent {
				matches = append(matches, scope)
			}
			continue
		}
		scopeSelector := profileHighRiskScopeSelector(scope)
		if scopeSelector == "" {
			continue
		}
		for _, claim := range dedupeStrings(ownership) {
			if strings.TrimSpace(claim) == "" {
				continue
			}
			if selectorsConflict(scopeSelector, nil, claim, nil) || selectorContains(scopeSelector, claim) || selectorContains(claim, scopeSelector) {
				matches = append(matches, scope)
				break
			}
		}
	}
	return dedupeStrings(matches)
}

func matchedMeetingTriggerRefs(profile map[string]any, decisionKind string, bundleState string, unresolvedConflictCount int, releaseAdjacent bool, profileHighRiskScopeRefs []string, freshSignalKinds []string) []string {
	matches := []string{}
	for _, item := range anyListToStrings(toList(getKey(profile, "profile.meeting_triggers", []any{}))) {
		trigger := strings.TrimSpace(item)
		if trigger == "" {
			continue
		}
		lower := strings.ToLower(trigger)
		switch {
		case strings.Contains(lower, "conflicting") && strings.Contains(lower, "evidence"):
			if decisionKind == "establish_bundle_governor" || decisionKind == "resolve_parallel_root_runs" || bundleState == "escalation-required" || unresolvedConflictCount > 0 {
				matches = append(matches, trigger)
			}
		case strings.Contains(lower, "blast") || strings.Contains(lower, "architecture"):
			if len(profileHighRiskScopeRefs) > 0 || bundleState == "conflicted" || bundleState == "escalation-required" {
				matches = append(matches, trigger)
			}
		case strings.Contains(lower, "release"):
			if releaseAdjacent || containsString(profileHighRiskScopeRefs, "release and RC flows") || containsString(freshSignalKinds, "interface_change") || containsString(freshSignalKinds, "requires_meeting") || containsString(freshSignalKinds, "review_feedback") {
				matches = append(matches, trigger)
			}
		}
	}
	return dedupeStrings(matches)
}

func containsReleaseMeetingTrigger(triggerRefs []string) bool {
	for _, item := range triggerRefs {
		if strings.Contains(strings.ToLower(strings.TrimSpace(item)), "release") {
			return true
		}
	}
	return false
}

func recommendedDecisionProtocol(decisionKind, decisionSurface string, releaseAdjacent bool, effectiveRiskLevel string, meetingTriggerRefs, freshSignalKinds []string) (string, string) {
	surface := strings.TrimSpace(decisionSurface)
	if surface != "council" && surface != "active-council" {
		return "decision-record", "Decision does not require a council, so the default recommendation is to record the bounded decision and its consequences."
	}
	highRisk := schedulerRiskRank(effectiveRiskLevel) >= schedulerRiskRank("high")
	if decisionKind == "review_gate_decision" ||
		((releaseAdjacent || highRisk || containsReleaseMeetingTrigger(meetingTriggerRefs)) &&
			(decisionKind == "review_follow_up_signals" ||
				decisionKind == "establish_bundle_governor" ||
				containsString(freshSignalKinds, "interface_change") ||
				containsString(freshSignalKinds, "requires_meeting") ||
				containsString(freshSignalKinds, "review_feedback"))) {
		return "incident-command", "Decision is release-adjacent or urgent enough that the council should optimize for bounded command, action ownership, and checkpointed stabilization."
	}
	return "evidence-review", "Decision benefits from bounded evidence review and challenge/synthesis rather than incident command."
}

func upsertMeetingRoleSpec(meeting map[string]any, role, mission string, allowedSources, forbiddenSources []string) {
	role = strings.TrimSpace(role)
	if role == "" {
		return
	}
	roleSpecs := toList(meeting["role_specs"])
	allowedValues := stringListToAny(dedupeStrings(allowedSources))
	forbiddenValues := stringListToAny(dedupeStrings(forbiddenSources))
	updated := false
	for index, item := range roleSpecs {
		spec, ok := toMap(item)
		if !ok || stringValue(spec["role"]) != role {
			continue
		}
		if strings.TrimSpace(mission) != "" {
			spec["mission"] = mission
		}
		if len(allowedSources) > 0 {
			spec["allowed_sources"] = allowedValues
		} else if _, ok := spec["allowed_sources"]; !ok {
			spec["allowed_sources"] = []any{}
		}
		if len(forbiddenSources) > 0 {
			spec["forbidden_sources"] = forbiddenValues
		} else if _, ok := spec["forbidden_sources"]; !ok {
			spec["forbidden_sources"] = []any{}
		}
		if stringValue(spec["created_at"]) == "" {
			spec["created_at"] = nowISO()
		}
		spec["updated_at"] = nowISO()
		roleSpecs[index] = spec
		updated = true
		break
	}
	if !updated {
		roleSpecs = append(roleSpecs, map[string]any{
			"role":              role,
			"mission":           mission,
			"allowed_sources":   allowedValues,
			"forbidden_sources": forbiddenValues,
			"created_at":        nowISO(),
			"updated_at":        nowISO(),
		})
	}
	meeting["role_specs"] = roleSpecs
	meeting["roles"] = appendUniqueText(toList(meeting["roles"]), role)
}

func ownershipClaims(run map[string]any) []string {
	selectors := ownershipSelectorKeys(run)
	if len(selectors) > 0 {
		return selectors
	}
	return dedupeStrings(anyListToStrings(toList(run["ownership_claims"])))
}

func anyListToStrings(items []any) []string {
	result := []string{}
	for _, item := range items {
		value := strings.TrimSpace(stringValue(item))
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func toMapSlice(items []any) []map[string]any {
	result := []map[string]any{}
	for _, item := range items {
		if mapping, ok := toMap(item); ok {
			result = append(result, mapping)
		}
	}
	return result
}

func normalizedStringSet(items []string) []string {
	values := dedupeStrings(items)
	sort.Strings(values)
	return values
}

func normalizeSemanticSelectorValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	parts := strings.Split(trimmed, "/")
	normalized := []string{}
	for _, part := range parts {
		slug := slugify(part)
		if slug != "" {
			normalized = append(normalized, slug)
		}
	}
	if len(normalized) == 0 {
		return slugify(trimmed)
	}
	return strings.Join(normalized, "/")
}

func hierarchicalPrefixMatch(prefix, candidate string) bool {
	if prefix == "" || candidate == "" {
		return false
	}
	return candidate == prefix || strings.HasPrefix(candidate, prefix+"/")
}

func parseOwnershipSelector(raw string) map[string]any {
	value := strings.TrimSpace(raw)
	kind := "opaque"
	selectorValue := value
	if parts := strings.SplitN(value, ":", 2); len(parts) == 2 && strings.TrimSpace(parts[0]) != "" && strings.TrimSpace(parts[1]) != "" {
		kind = strings.TrimSpace(parts[0])
		selectorValue = strings.TrimSpace(parts[1])
	}
	switch kind {
	case "directory", "folder":
		kind = "dir"
	case "path":
		kind = "file"
	case "component":
		kind = "module"
	case "pattern":
		kind = "glob"
	}
	switch kind {
	case "file", "dir", "module", "glob":
		selectorValue = strings.TrimSpace(strings.TrimSuffix(selectorValue, "/"))
		selectorValue = strings.TrimPrefix(selectorValue, "./")
	case "scope", "interface":
		selectorValue = normalizeSemanticSelectorValue(selectorValue)
	}
	family := "opaque"
	match := "exact"
	bounded := true
	switch kind {
	case "mission", "branch":
		family = "anchor"
		bounded = false
	case "scope", "interface":
		family = "semantic"
	case "file":
		family = "path"
	case "dir", "module":
		family = "path"
		match = "prefix"
	case "glob":
		family = "path"
		match = "glob"
	}
	if selectorValue != "" {
		value = kind + ":" + selectorValue
	}
	return map[string]any{
		"selector": value,
		"kind":     kind,
		"value":    selectorValue,
		"family":   family,
		"match":    match,
		"bounded":  bounded,
	}
}

func selectorRecordsFromClaims(claims []string) []any {
	records := []any{}
	for _, claim := range dedupeStrings(claims) {
		records = append(records, parseOwnershipSelector(claim))
	}
	return records
}

func ownershipSelectorKeys(mapping map[string]any) []string {
	selectors := []string{}
	for _, item := range toList(mapping["ownership_selectors"]) {
		record, ok := toMap(item)
		if !ok {
			continue
		}
		selector := strings.TrimSpace(stringValue(record["selector"]))
		if selector == "" {
			kind := strings.TrimSpace(stringValue(record["kind"]))
			value := strings.TrimSpace(stringValue(record["value"]))
			if kind != "" && value != "" {
				selector = kind + ":" + value
			}
		}
		if selector != "" {
			selectors = append(selectors, selector)
		}
	}
	return dedupeStrings(selectors)
}

func ownershipExclusions(mapping map[string]any) []string {
	selectors := []string{}
	for _, item := range toList(mapping["ownership_exclusion_selectors"]) {
		record, ok := toMap(item)
		if !ok {
			continue
		}
		selector := strings.TrimSpace(stringValue(record["selector"]))
		if selector == "" {
			kind := strings.TrimSpace(stringValue(record["kind"]))
			value := strings.TrimSpace(stringValue(record["value"]))
			if kind != "" && value != "" {
				selector = kind + ":" + value
			}
		}
		if selector != "" {
			selectors = append(selectors, selector)
		}
	}
	if len(selectors) > 0 {
		return dedupeStrings(selectors)
	}
	return dedupeStrings(anyListToStrings(toList(mapping["ownership_exclusions"])))
}

func selectorRecordFor(raw string) map[string]any {
	return parseOwnershipSelector(raw)
}

func globStaticPrefix(pattern string) string {
	segments := strings.Split(pattern, "/")
	staticSegments := []string{}
	for _, segment := range segments {
		if strings.ContainsAny(segment, "*?[") {
			break
		}
		if segment != "" {
			staticSegments = append(staticSegments, segment)
		}
	}
	return strings.Join(staticSegments, "/")
}

func selectorContains(containerRaw, memberRaw string) bool {
	container := selectorRecordFor(containerRaw)
	member := selectorRecordFor(memberRaw)
	containerFamily := stringValue(container["family"])
	memberFamily := stringValue(member["family"])
	containerKind := stringValue(container["kind"])
	memberKind := stringValue(member["kind"])
	containerValue := stringValue(container["value"])
	memberValue := stringValue(member["value"])
	if containerFamily == "" || memberFamily == "" || containerFamily != memberFamily {
		return false
	}
	switch containerFamily {
	case "semantic":
		return hierarchicalPrefixMatch(containerValue, memberValue)
	case "path":
		switch containerKind {
		case "file":
			return memberKind == "file" && containerValue == memberValue
		case "dir", "module":
			switch memberKind {
			case "file", "dir", "module":
				return hierarchicalPrefixMatch(containerValue, memberValue)
			case "glob":
				prefix := globStaticPrefix(memberValue)
				return prefix != "" && hierarchicalPrefixMatch(containerValue, prefix)
			}
		}
	}
	return strings.TrimSpace(containerRaw) == strings.TrimSpace(memberRaw)
}

func selectorExcluded(candidate string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if strings.TrimSpace(exclusion) == strings.TrimSpace(candidate) || selectorContains(exclusion, candidate) {
			return true
		}
	}
	return false
}

func selectorsConflict(leftClaim string, leftExclusions []string, rightClaim string, rightExclusions []string) bool {
	if !ownershipSelectorsOverlap(leftClaim, rightClaim) {
		return false
	}
	if selectorExcluded(rightClaim, leftExclusions) {
		return false
	}
	if selectorExcluded(leftClaim, rightExclusions) {
		return false
	}
	return true
}

func moreSpecificSelectorWinner(leftRaw, rightRaw string) (string, bool) {
	left := selectorRecordFor(leftRaw)
	right := selectorRecordFor(rightRaw)
	leftFamily := stringValue(left["family"])
	rightFamily := stringValue(right["family"])
	leftKind := stringValue(left["kind"])
	rightKind := stringValue(right["kind"])
	if leftFamily == "" || rightFamily == "" || leftFamily != rightFamily {
		return "", false
	}
	if strings.TrimSpace(leftRaw) == strings.TrimSpace(rightRaw) {
		return "", false
	}
	if leftFamily == "path" {
		if leftKind == "file" && rightKind != "file" && ownershipSelectorsOverlap(leftRaw, rightRaw) {
			return "left", true
		}
		if rightKind == "file" && leftKind != "file" && ownershipSelectorsOverlap(leftRaw, rightRaw) {
			return "right", true
		}
	}
	leftContainsRight := selectorContains(leftRaw, rightRaw)
	rightContainsLeft := selectorContains(rightRaw, leftRaw)
	switch {
	case leftContainsRight && !rightContainsLeft:
		return "right", true
	case rightContainsLeft && !leftContainsRight:
		return "left", true
	default:
		return "", false
	}
}

func pathOverlap(leftKind, leftValue, rightKind, rightValue string) bool {
	if leftValue == "" || rightValue == "" {
		return false
	}
	sameOrPrefix := func(prefix, candidate string) bool {
		return candidate == prefix || strings.HasPrefix(candidate, prefix+"/")
	}
	globMatchesPath := func(pattern, candidate string) bool {
		matched, err := filepath.Match(pattern, candidate)
		return err == nil && matched
	}
	switch {
	case leftKind == "file" && rightKind == "file":
		return leftValue == rightValue
	case leftKind == "file" && (rightKind == "dir" || rightKind == "module"):
		return sameOrPrefix(rightValue, leftValue)
	case rightKind == "file" && (leftKind == "dir" || leftKind == "module"):
		return sameOrPrefix(leftValue, rightValue)
	case leftKind == "glob" && rightKind == "file":
		return globMatchesPath(leftValue, rightValue)
	case rightKind == "glob" && leftKind == "file":
		return globMatchesPath(rightValue, leftValue)
	case leftKind == "glob" && (rightKind == "dir" || rightKind == "module"):
		prefix := globStaticPrefix(leftValue)
		return prefix == "" || sameOrPrefix(rightValue, prefix) || sameOrPrefix(prefix, rightValue)
	case rightKind == "glob" && (leftKind == "dir" || leftKind == "module"):
		prefix := globStaticPrefix(rightValue)
		return prefix == "" || sameOrPrefix(leftValue, prefix) || sameOrPrefix(prefix, leftValue)
	case leftKind == "glob" && rightKind == "glob":
		leftPrefix := globStaticPrefix(leftValue)
		rightPrefix := globStaticPrefix(rightValue)
		if leftValue == rightValue {
			return true
		}
		if leftPrefix == "" || rightPrefix == "" {
			return false
		}
		return sameOrPrefix(leftPrefix, rightPrefix) || sameOrPrefix(rightPrefix, leftPrefix)
	case (leftKind == "dir" || leftKind == "module") && (rightKind == "dir" || rightKind == "module"):
		return sameOrPrefix(leftValue, rightValue) || sameOrPrefix(rightValue, leftValue)
	default:
		return false
	}
}

func ownershipSelectorsOverlap(leftRaw, rightRaw string) bool {
	left := selectorRecordFor(leftRaw)
	right := selectorRecordFor(rightRaw)
	leftFamily := stringValue(left["family"])
	rightFamily := stringValue(right["family"])
	leftKind := stringValue(left["kind"])
	rightKind := stringValue(right["kind"])
	leftValue := stringValue(left["value"])
	rightValue := stringValue(right["value"])
	if leftFamily == "" || rightFamily == "" {
		return strings.TrimSpace(leftRaw) == strings.TrimSpace(rightRaw)
	}
	if leftFamily != rightFamily {
		return false
	}
	semanticOverlap := func(leftValue, rightValue string) bool {
		if leftValue == "" || rightValue == "" {
			return false
		}
		return hierarchicalPrefixMatch(leftValue, rightValue) || hierarchicalPrefixMatch(rightValue, leftValue)
	}
	switch leftFamily {
	case "anchor", "semantic":
		if leftFamily == "semantic" {
			return semanticOverlap(leftValue, rightValue)
		}
		return leftValue == rightValue
	case "path":
		return pathOverlap(leftKind, leftValue, rightKind, rightValue)
	default:
		return leftKind == rightKind && leftValue == rightValue
	}
}

func boundedOwnershipSelector(raw string) bool {
	return boolValue(parseOwnershipSelector(raw)["bounded"])
}

func setOwnershipSelectors(mapping map[string]any, claims []string) {
	normalized := dedupeStrings(claims)
	mapping["ownership_claims"] = stringListToAny(normalized)
	mapping["ownership_selectors"] = selectorRecordsFromClaims(normalized)
}

func setOwnershipExclusions(mapping map[string]any, claims []string) {
	normalized := dedupeStrings(claims)
	mapping["ownership_exclusions"] = stringListToAny(normalized)
	mapping["ownership_exclusion_selectors"] = selectorRecordsFromClaims(normalized)
}

func ownershipTransferEvents(mapping map[string]any) []any {
	return toList(mapping["ownership_transfers"])
}

func defaultDesiredStateForRun(ownerType string) map[string]any {
	desired := map[string]any{
		"goal_state":          "completed",
		"verification_status": "passed",
		"review_status":       "",
		"checks_status":       "",
		"ship_status":         "",
	}
	if ownerType == "leader" {
		desired["verification_status"] = "passed"
	}
	return desired
}

func normalizeDesiredState(run map[string]any) map[string]any {
	desired, _ := toMap(run["desired_state"])
	if desired == nil {
		desired = map[string]any{}
	}
	defaults := defaultDesiredStateForRun(stringValue(run["owner_type"]))
	for key, value := range defaults {
		if _, ok := desired[key]; !ok {
			desired[key] = value
		}
	}
	run["desired_state"] = desired
	return desired
}

func runDesiredStateSatisfied(run map[string]any) (bool, []string) {
	desired := normalizeDesiredState(run)
	missing := []string{}
	if goal := strings.TrimSpace(stringValue(desired["goal_state"])); goal != "" && stringValue(run["state"]) != goal {
		missing = append(missing, "state:"+goal)
	}
	if expected := strings.TrimSpace(stringValue(desired["verification_status"])); expected != "" && stringValue(run["verification_status"]) != expected {
		missing = append(missing, "verification_status:"+expected)
	}
	if expected := strings.TrimSpace(stringValue(desired["review_status"])); expected != "" && stringValue(run["review_status"]) != expected {
		missing = append(missing, "review_status:"+expected)
	}
	if expected := strings.TrimSpace(stringValue(desired["checks_status"])); expected != "" && stringValue(run["checks_status"]) != expected {
		missing = append(missing, "checks_status:"+expected)
	}
	if expected := strings.TrimSpace(stringValue(desired["ship_status"])); expected != "" && stringValue(run["ship_status"]) != expected {
		missing = append(missing, "ship_status:"+expected)
	}
	return len(missing) == 0, missing
}

func dependencyTargets(run map[string]any) []string {
	return dedupeStrings(anyListToStrings(toList(run["depends_on"])))
}

func unsatisfiedDependencies(run map[string]any, registry map[string]any) []string {
	unsatisfied := []string{}
	for _, dependencyID := range dependencyTargets(run) {
		index := findRunIndex(registry, dependencyID)
		if index < 0 {
			unsatisfied = append(unsatisfied, dependencyID)
			continue
		}
		dependencyRun, _ := toMap(toList(registry["runs"])[index])
		satisfied, _ := runDesiredStateSatisfied(dependencyRun)
		if !satisfied {
			unsatisfied = append(unsatisfied, dependencyID)
		}
	}
	return dedupeStrings(unsatisfied)
}

func overlapPolicy(run map[string]any) string {
	policy := strings.TrimSpace(stringValue(run["overlap_policy"]))
	if policy == "" {
		return "exclusive"
	}
	return policy
}

func overlappingOwnershipRuns(registry map[string]any, candidateClaims []string, candidatePolicy string, ignoreRunIDs []string) []map[string]any {
	if candidatePolicy == "shared" || len(candidateClaims) == 0 {
		return nil
	}
	ignore := map[string]bool{}
	for _, runID := range ignoreRunIDs {
		ignore[runID] = true
	}
	candidateClaims = dedupeStrings(candidateClaims)
	conflicts := []map[string]any{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		runID := stringValue(run["run_id"])
		if ignore[runID] || !activeRun(run) || overlapPolicy(run) == "shared" {
			continue
		}
		shared := []string{}
		runExclusions := ownershipExclusions(run)
		for _, candidate := range candidateClaims {
			for _, claim := range ownershipClaims(run) {
				if selectorsConflict(candidate, nil, claim, runExclusions) {
					shared = append(shared, candidate+"~"+claim)
				}
			}
		}
		if len(shared) == 0 {
			continue
		}
		conflicts = append(conflicts, map[string]any{
			"run_id":           runID,
			"mission_id":       stringValue(run["mission_id"]),
			"owner_id":         stringValue(run["owner_id"]),
			"owner_type":       stringValue(run["owner_type"]),
			"ownership_claims": stringListToAny(dedupeStrings(shared)),
		})
	}
	return conflicts
}

func validAuthorityLevel(level string) bool {
	switch strings.TrimSpace(level) {
	case "L1", "L2", "L3", "L4", "L5":
		return true
	default:
		return false
	}
}

func authorityRank(level string) int {
	switch strings.TrimSpace(level) {
	case "L1":
		return 5
	case "L2":
		return 4
	case "L3":
		return 3
	case "L4":
		return 2
	case "L5":
		return 1
	default:
		return 0
	}
}

func authorityLevelForRank(rank int) string {
	switch rank {
	case 5:
		return "L1"
	case 4:
		return "L2"
	case 3:
		return "L3"
	case 2:
		return "L4"
	case 1:
		return "L5"
	default:
		return ""
	}
}

func isoAfter(left, right string) bool {
	leftTime, err := parseISO(left)
	if err != nil || leftTime == nil {
		return false
	}
	rightTime, err := parseISO(right)
	if err != nil || rightTime == nil {
		return true
	}
	return leftTime.After(*rightTime)
}

func runByID(registry map[string]any, runID string) map[string]any {
	index := findRunIndex(registry, runID)
	if index < 0 {
		return nil
	}
	run, _ := toMap(toList(registry["runs"])[index])
	return run
}

func normalizeSignalKind(kind string) string {
	switch strings.TrimSpace(strings.ToLower(kind)) {
	case "requires_meeting", "meeting_required", "meeting-required":
		return "requires_meeting"
	case "requires_rebase", "requires-rebase":
		return "requires_rebase"
	case "review_feedback", "review-feedback":
		return "review_feedback"
	case "priority_change", "priority-change":
		return "priority_change"
	case "risk_raised", "risk-raised":
		return "risk_raised"
	case "dependency_discovered", "dependency-discovered":
		return "dependency_discovered"
	case "dependency_cleared", "dependency-cleared":
		return "dependency_cleared"
	case "interface_change", "interface-change":
		return "interface_change"
	case "release_ready", "release-ready":
		return "release_ready"
	default:
		return strings.TrimSpace(kind)
	}
}

func signalSupersession(signal map[string]any, registry map[string]any) (bool, map[string]any) {
	runID := stringValue(signal["run_id"])
	if runID == "" {
		return false, nil
	}
	run := runByID(registry, runID)
	if run == nil {
		return false, nil
	}
	emittedAt := stringValue(signal["emitted_at"])
	signalLevel := authorityRank(stringValue(signal["authority_level"]))
	if signalLevel == 0 {
		signalLevel = authorityRank("L4")
	}
	runLevel := authorityRank("L2")
	runUpdated := firstNonBlank(stringValue(run["updated_at"]), stringValue(run["last_reported_at"]), stringValue(run["started_at"]))
	switch normalizeSignalKind(stringValue(signal["kind"])) {
	case "blocked":
		if stringValue(run["state"]) != "blocked" && stringValue(run["state"]) != "waiting" && runLevel >= signalLevel && isoAfter(runUpdated, emittedAt) {
			return true, map[string]any{
				"resolution_kind": "authority_superseded",
				"resolution_note": "Fresher run state supersedes the blocked signal.",
				"superseded_by":   "run:" + runID,
			}
		}
	case "requires_meeting":
		if meeting := latestClosedMeetingForRun(runID); meeting != nil {
			closedAt := firstNonBlank(stringValue(meeting["closed_at"]), stringValue(meeting["updated_at"]))
			if isoAfter(closedAt, emittedAt) {
				return true, map[string]any{
					"resolution_kind": "authority_superseded",
					"resolution_note": "Closed meeting outcome supersedes the prior meeting-request signal.",
					"superseded_by":   "meeting:" + stringValue(meeting["meeting_id"]),
				}
			}
		}
	case "progress":
		if (stringValue(run["state"]) == "completed" || stringValue(run["state"]) == "blocked") && runLevel >= signalLevel && isoAfter(runUpdated, emittedAt) {
			return true, map[string]any{
				"resolution_kind": "authority_superseded",
				"resolution_note": "Fresher run state supersedes the earlier progress signal.",
				"superseded_by":   "run:" + runID,
			}
		}
	}
	return false, nil
}

func activeSignalsForRun(signalLog map[string]any, runID string) []map[string]any {
	signals := []map[string]any{}
	for _, item := range toList(signalLog["signals"]) {
		signal, _ := toMap(item)
		if !activeSignal(signal) || stringValue(signal["run_id"]) != runID {
			continue
		}
		signals = append(signals, signal)
	}
	sort.Slice(signals, func(i, j int) bool {
		return isoAfter(
			firstNonBlank(stringValue(signals[i]["emitted_at"]), stringValue(signals[i]["updated_at"])),
			firstNonBlank(stringValue(signals[j]["emitted_at"]), stringValue(signals[j]["updated_at"])),
		)
	})
	return signals
}

func loadContractForRun(run map[string]any) map[string]any {
	contractRef := strings.TrimSpace(stringValue(run["contract_ref"]))
	if contractRef == "" {
		return nil
	}
	root, _ := toMap(loadData(contractRef, map[string]any{}))
	contract, _ := toMap(root["contract"])
	return contract
}

func syncContractControlPlaneFields(contract, run map[string]any) map[string]any {
	if contract == nil {
		return nil
	}
	contract["branch"] = stringValue(run["branch"])
	contract["worktree_path"] = stringValue(run["worktree_path"])
	contract["desired_state"] = cloneValue(normalizeDesiredState(run))
	setOwnershipSelectors(contract, ownershipClaims(run))
	setOwnershipExclusions(contract, ownershipExclusions(run))
	freshness, _ := toMap(contract["freshness"])
	if freshness == nil {
		freshness = map[string]any{}
	}
	freshness["source_checked_at"] = nowISO()
	freshness["drift_status"] = "current"
	contract["freshness"] = freshness
	return contract
}

func persistContractForRun(run map[string]any) error {
	contractRef := strings.TrimSpace(stringValue(run["contract_ref"]))
	if contractRef == "" {
		return nil
	}
	contract := loadContractForRun(run)
	if contract == nil {
		return nil
	}
	contract = syncContractControlPlaneFields(contract, run)
	return dumpData(contractRef, map[string]any{
		"version":  1,
		"contract": contract,
	})
}

func worktreeAuthoritySnapshot(run map[string]any) map[string]any {
	worktreePath := strings.TrimSpace(stringValue(run["worktree_path"]))
	if worktreePath == "" || !exists(worktreePath) {
		return nil
	}
	snapshot := map[string]any{
		"worktree_path": worktreePath,
	}
	if out, err := runCommandAt(worktreePath, []string{"git", "branch", "--show-current"}); err == nil {
		snapshot["branch"] = strings.TrimSpace(out)
	}
	if out, err := runCommandAt(worktreePath, []string{"git", "status", "--porcelain"}); err == nil {
		filtered := []string{}
		for _, line := range strings.Split(out, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || trimmed == "?? .mission" {
				continue
			}
			filtered = append(filtered, trimmed)
		}
		snapshot["dirty"] = len(filtered) > 0
		if len(filtered) > 0 {
			snapshot["dirty_entries"] = stringListToAny(filtered)
		}
	}
	return snapshot
}

func runAuthorityAudit(run map[string]any, registry map[string]any, signalLog map[string]any) map[string]any {
	runID := stringValue(run["run_id"])
	runUpdatedAt := firstNonBlank(
		stringValue(run["updated_at"]),
		stringValue(run["last_reported_at"]),
		stringValue(run["started_at"]),
		stringValue(run["created_at"]),
	)
	runLevel := "L2"
	runClaims := []any{
		map[string]any{
			"claim_ref":   "execution_state",
			"value":       stringValue(run["state"]),
			"source_ref":  "run:" + runID,
			"updated_at":  runUpdatedAt,
			"directness":  "direct",
			"authority":   runLevel,
			"specificity": "run_state",
		},
	}

	desiredSatisfied, desiredGaps := runDesiredStateSatisfied(run)
	desiredGapRecords := []any{}
	for _, gap := range desiredGaps {
		desiredGapRecords = append(desiredGapRecords, map[string]any{
			"claim_ref":         gap,
			"resolution_kind":   "desired_state_unmet",
			"authority_level":   "L3",
			"source_ref":        "desired_state:" + runID,
			"winning_source":    "run:" + runID,
			"note":              fmt.Sprintf("Current run state does not yet satisfy desired state %s.", gap),
			"desired_state_gap": true,
		})
	}

	adoptionGaps := []any{}
	if meeting := governingClosedMeetingForRun(run); meeting != nil && meetingHasOutput(meeting) {
		outcome := meetingOutcome(meeting)
		if summary := stringValue(outcome["summary"]); summary != "" && strings.TrimSpace(stringValue(run["latest_summary"])) != summary {
			adoptionGaps = append(adoptionGaps, map[string]any{
				"claim_ref":       "latest_summary",
				"resolution_kind": "adoption_pending",
				"authority_level": "L3",
				"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
				"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
				"current_value":   stringValue(run["latest_summary"]),
				"expected_value":  summary,
				"note":            "Closed meeting summary has not been adopted into the run state.",
			})
		}
		if nextAction := meetingAdoptedNextAction(meeting); nextAction != "" && strings.TrimSpace(stringValue(run["next_expected_action"])) != nextAction {
			adoptionGaps = append(adoptionGaps, map[string]any{
				"claim_ref":       "next_expected_action",
				"resolution_kind": "adoption_pending",
				"authority_level": "L3",
				"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
				"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
				"current_value":   stringValue(run["next_expected_action"]),
				"expected_value":  nextAction,
				"note":            "Closed meeting next action or primary adopted follow-up has not been adopted into the run state.",
			})
		}
		if strings.TrimSpace(stringValue(run["adopted_from_meeting_id"])) != stringValue(meeting["meeting_id"]) {
			adoptionGaps = append(adoptionGaps, map[string]any{
				"claim_ref":       "adopted_from_meeting_id",
				"resolution_kind": "adoption_pending",
				"authority_level": "L3",
				"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
				"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
				"current_value":   stringValue(run["adopted_from_meeting_id"]),
				"expected_value":  stringValue(meeting["meeting_id"]),
				"note":            "Closed meeting adoption metadata has not been written back into the run state.",
			})
		}
		orchestration, _ := toMap(meeting["orchestration"])
		outputPacket, _ := toMap(orchestration["output_packet"])
		if len(outputPacket) > 0 {
			expectedKind := stringValue(outputPacket["output_packet_kind"])
			if expectedKind != "" && strings.TrimSpace(stringValue(run["adopted_output_packet_kind"])) != expectedKind {
				adoptionGaps = append(adoptionGaps, map[string]any{
					"claim_ref":       "adopted_output_packet_kind",
					"resolution_kind": "adoption_pending",
					"authority_level": "L3",
					"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
					"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
					"current_value":   stringValue(run["adopted_output_packet_kind"]),
					"expected_value":  expectedKind,
					"note":            "Closed meeting output packet kind has not been adopted into the run state.",
				})
			}
			expectedFollowUpCount := intValue(outputPacket["follow_up_packet_count"])
			if intValue(run["adopted_output_follow_up_packet_count"]) != expectedFollowUpCount {
				adoptionGaps = append(adoptionGaps, map[string]any{
					"claim_ref":       "adopted_output_follow_up_packet_count",
					"resolution_kind": "adoption_pending",
					"authority_level": "L3",
					"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
					"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
					"current_value":   intValue(run["adopted_output_follow_up_packet_count"]),
					"expected_value":  expectedFollowUpCount,
					"note":            "Closed meeting follow-up packet count has not been adopted into the run state.",
				})
			}
			expectedPrimaryFollowUp := meetingPrimaryFollowUpPacket(meeting)
			if expectedPrimaryFollowUp != nil {
				currentPrimaryFollowUp, _ := toMap(run["adopted_primary_follow_up_packet"])
				if stringValue(currentPrimaryFollowUp["follow_up_packet_id"]) != stringValue(expectedPrimaryFollowUp["follow_up_packet_id"]) {
					adoptionGaps = append(adoptionGaps, map[string]any{
						"claim_ref":       "adopted_primary_follow_up_packet",
						"resolution_kind": "adoption_pending",
						"authority_level": "L3",
						"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
						"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
						"current_value":   stringValue(currentPrimaryFollowUp["follow_up_packet_id"]),
						"expected_value":  stringValue(expectedPrimaryFollowUp["follow_up_packet_id"]),
						"note":            "Closed meeting primary follow-up packet has not been adopted into the run state.",
					})
				}
			}
		}
	}

	contractGaps := []any{}
	if contract := loadContractForRun(run); contract != nil {
		contractDesired, _ := toMap(contract["desired_state"])
		runDesired := normalizeDesiredState(run)
		for _, field := range []string{"goal_state", "verification_status", "review_status", "checks_status", "ship_status"} {
			if stringValue(contractDesired[field]) != stringValue(runDesired[field]) {
				contractGaps = append(contractGaps, map[string]any{
					"claim_ref":       "desired_state." + field,
					"resolution_kind": "run_state_overrides_contract",
					"source_ref":      "contract:" + runID + ".desired_state." + field,
					"winning_source":  "run:" + runID,
					"losing_source":   "contract:" + runID + ".desired_state." + field,
					"authority_level": "L2",
					"current_value":   stringValue(contractDesired[field]),
					"expected_value":  stringValue(runDesired[field]),
					"note":            fmt.Sprintf("Contract desired state for %s is stale relative to the run control plane.", field),
				})
			}
		}
		if stringValue(contract["branch"]) != stringValue(run["branch"]) {
			contractGaps = append(contractGaps, map[string]any{
				"claim_ref":       "branch",
				"resolution_kind": "run_state_overrides_contract",
				"source_ref":      "contract:" + runID + ".branch",
				"winning_source":  "run:" + runID,
				"losing_source":   "contract:" + runID + ".branch",
				"authority_level": "L2",
				"current_value":   stringValue(contract["branch"]),
				"expected_value":  stringValue(run["branch"]),
				"note":            "Contract branch does not match the current run branch.",
			})
		}
		if stringValue(contract["worktree_path"]) != stringValue(run["worktree_path"]) {
			contractGaps = append(contractGaps, map[string]any{
				"claim_ref":       "worktree_path",
				"resolution_kind": "run_state_overrides_contract",
				"source_ref":      "contract:" + runID + ".worktree_path",
				"winning_source":  "run:" + runID,
				"losing_source":   "contract:" + runID + ".worktree_path",
				"authority_level": "L2",
				"current_value":   stringValue(contract["worktree_path"]),
				"expected_value":  stringValue(run["worktree_path"]),
				"note":            "Contract worktree path does not match the current run provisioning state.",
			})
		}
		if strings.Join(ownershipClaims(contract), "|") != strings.Join(ownershipClaims(run), "|") {
			contractGaps = append(contractGaps, map[string]any{
				"claim_ref":       "ownership_claims",
				"resolution_kind": "run_state_overrides_contract",
				"source_ref":      "contract:" + runID + ".ownership_claims",
				"winning_source":  "run:" + runID,
				"losing_source":   "contract:" + runID + ".ownership_claims",
				"authority_level": "L2",
				"current_value":   strings.Join(ownershipClaims(contract), "|"),
				"expected_value":  strings.Join(ownershipClaims(run), "|"),
				"note":            "Contract ownership claims no longer match the run ownership boundary.",
			})
		}
	}

	liveStateGaps := []any{}
	if worktreePath := strings.TrimSpace(stringValue(run["worktree_path"])); worktreePath != "" && !exists(worktreePath) {
		liveStateGaps = append(liveStateGaps, map[string]any{
			"claim_ref":       "live_worktree_missing",
			"resolution_kind": "live_state_overrides_run",
			"source_ref":      "run:" + runID,
			"winning_source":  "worktree:" + worktreePath,
			"losing_source":   "run:" + runID + ".worktree_path",
			"authority_level": "L1",
			"current_value":   worktreePath,
			"expected_value":  "present_on_disk",
			"note":            "Run expects a provisioned worktree, but the referenced live worktree path is missing on disk.",
		})
	} else if liveState := worktreeAuthoritySnapshot(run); liveState != nil {
		if liveBranch := strings.TrimSpace(stringValue(liveState["branch"])); liveBranch != "" && liveBranch != stringValue(run["branch"]) {
			liveStateGaps = append(liveStateGaps, map[string]any{
				"claim_ref":       "live_branch",
				"resolution_kind": "live_state_overrides_run",
				"source_ref":      "worktree:" + stringValue(run["worktree_path"]),
				"winning_source":  "worktree:" + stringValue(run["worktree_path"]),
				"losing_source":   "run:" + runID + ".branch",
				"authority_level": "L1",
				"current_value":   stringValue(run["branch"]),
				"expected_value":  liveBranch,
				"note":            "Live worktree branch does not match the run branch record.",
			})
		}
		if boolValue(liveState["dirty"]) {
			desired := normalizeDesiredState(run)
			cleanIntent := stringValue(run["state"]) == "completed" ||
				stringValue(desired["goal_state"]) == "done" ||
				stringValue(desired["goal_state"]) == "completed" ||
				stringValue(run["verification_status"]) == "passed" ||
				stringValue(desired["verification_status"]) == "passed" ||
				stringValue(run["ship_status"]) == "approved" ||
				stringValue(run["ship_status"]) == "merged" ||
				stringValue(desired["ship_status"]) == "approved" ||
				stringValue(desired["ship_status"]) == "merged"
			if cleanIntent {
				liveStateGaps = append(liveStateGaps, map[string]any{
					"claim_ref":       "live_dirty",
					"resolution_kind": "live_state_overrides_run",
					"source_ref":      "worktree:" + stringValue(run["worktree_path"]),
					"winning_source":  "worktree:" + stringValue(run["worktree_path"]),
					"losing_source":   "run:" + runID,
					"authority_level": "L1",
					"current_value":   "dirty",
					"expected_value":  "clean",
					"note":            "Live worktree is dirty even though run/desired state implies the execution should be clean enough for governed verification or handoff.",
				})
			}
		}
	}

	activeSignals := activeSignalsForRun(signalLog, runID)
	signalSources := []any{}
	conflicts := []any{}
	supersededSignals := []any{}

	shipStatus := firstNonBlank(stringValue(run["ship_status"]), "unstarted")
	if shipStatus == "approved" || shipStatus == "merged" {
		if stringValue(run["verification_status"]) != "passed" {
			conflicts = append(conflicts, map[string]any{
				"claim_ref":       "verification_status",
				"source_ref":      "run:" + runID + ".verification_status",
				"winning_source":  "run:" + runID + ".verification_status",
				"losing_source":   "run:" + runID + ".ship_status",
				"authority_level": "L1",
				"resolution_kind": "gate_state_overrides_ship_status",
				"current_value":   shipStatus,
				"expected_value":  "verification_status=passed",
				"note":            "Ship status claims approval/merge readiness, but verification_status is not passed.",
			})
		}
		if stringValue(run["pr_ref"]) != "" && stringValue(run["review_status"]) != "approved" {
			conflicts = append(conflicts, map[string]any{
				"claim_ref":       "review_status",
				"source_ref":      "run:" + runID + ".review_status",
				"winning_source":  "run:" + runID + ".review_status",
				"losing_source":   "run:" + runID + ".ship_status",
				"authority_level": "L1",
				"resolution_kind": "gate_state_overrides_ship_status",
				"current_value":   shipStatus,
				"expected_value":  "review_status=approved",
				"note":            "Ship status claims approval/merge readiness, but review_status is not approved.",
			})
		}
		if stringValue(run["pr_ref"]) != "" && stringValue(run["checks_status"]) != "passed" {
			conflicts = append(conflicts, map[string]any{
				"claim_ref":       "checks_status",
				"source_ref":      "run:" + runID + ".checks_status",
				"winning_source":  "run:" + runID + ".checks_status",
				"losing_source":   "run:" + runID + ".ship_status",
				"authority_level": "L1",
				"resolution_kind": "gate_state_overrides_ship_status",
				"current_value":   shipStatus,
				"expected_value":  "checks_status=passed",
				"note":            "Ship status claims approval/merge readiness, but checks_status is not passed.",
			})
		}
	}
	for _, signal := range activeSignals {
		signalID := stringValue(signal["signal_id"])
		emittedAt := firstNonBlank(stringValue(signal["emitted_at"]), stringValue(signal["updated_at"]))
		signalLevel := firstNonBlank(stringValue(signal["authority_level"]), "L4")
		signalRank := authorityRank(signalLevel)
		signalSources = append(signalSources, map[string]any{
			"source_ref":       "signal:" + signalID,
			"kind":             normalizeSignalKind(stringValue(signal["kind"])),
			"summary":          stringValue(signal["summary"]),
			"authority_level":  signalLevel,
			"authority_rank":   signalRank,
			"updated_at":       emittedAt,
			"directness":       "local_signal",
			"resolution_state": stringValue(signal["state"]),
		})
		if superseded, resolution := signalSupersession(signal, registry); superseded {
			resolution["signal_id"] = signalID
			resolution["signal_kind"] = normalizeSignalKind(stringValue(signal["kind"]))
			resolution["authority_level"] = signalLevel
			resolution["emitted_at"] = emittedAt
			supersededSignals = append(supersededSignals, resolution)
			continue
		}

		conflict := map[string]any(nil)
		switch normalizeSignalKind(stringValue(signal["kind"])) {
		case "blocked":
			if stringValue(run["state"]) != "blocked" && stringValue(run["state"]) != "waiting" {
				if signalRank > authorityRank(runLevel) || (signalRank == authorityRank(runLevel) && isoAfter(emittedAt, runUpdatedAt)) {
					conflict = map[string]any{
						"claim_ref":       "execution_state",
						"signal_id":       signalID,
						"signal_kind":     normalizeSignalKind(stringValue(signal["kind"])),
						"source_ref":      "signal:" + signalID,
						"winning_source":  "signal:" + signalID,
						"losing_source":   "run:" + runID,
						"authority_level": signalLevel,
						"resolution_kind": "run_state_overridden",
						"current_value":   stringValue(run["state"]),
						"expected_value":  "blocked|waiting",
						"note":            "Active blocked signal conflicts with the current run state and should be resolved or adopted.",
					}
				}
			}
		case "requires_meeting":
			hasActiveMeeting := false
			for _, item := range listMeetings(false) {
				meeting, _ := toMap(item)
				if stringValue(meeting["run_id"]) == runID && stringValue(meeting["state"]) == "active" {
					hasActiveMeeting = true
					break
				}
			}
			if !hasActiveMeeting {
				if signalRank > authorityRank(runLevel) || (signalRank == authorityRank(runLevel) && isoAfter(emittedAt, runUpdatedAt)) {
					conflict = map[string]any{
						"claim_ref":       "meeting_required",
						"signal_id":       signalID,
						"signal_kind":     normalizeSignalKind(stringValue(signal["kind"])),
						"source_ref":      "signal:" + signalID,
						"winning_source":  "signal:" + signalID,
						"losing_source":   "run:" + runID,
						"authority_level": signalLevel,
						"resolution_kind": "meeting_request_unadopted",
						"current_value":   stringValue(run["state"]),
						"expected_value":  "waiting_with_active_meeting",
						"note":            "Active requires_meeting signal has not been converted into an active meeting or waiting run state.",
					}
				}
			}
		}
		if conflict != nil {
			conflicts = append(conflicts, conflict)
		}
	}

	return map[string]any{
		"run_id":                            runID,
		"canonical_source":                  "run:" + runID,
		"canonical_authority_level":         runLevel,
		"run_claims":                        runClaims,
		"signal_sources":                    signalSources,
		"artifact_authority_conflict_count": len(contractGaps) + len(liveStateGaps),
		"adoption_gap_count":                len(adoptionGaps),
		"contract_gap_count":                len(contractGaps),
		"higher_authority_conflict_count":   len(conflicts),
		"superseded_signal_count":           len(supersededSignals),
		"desired_state_gap_count":           len(desiredGapRecords),
		"live_state_gap_count":              len(liveStateGaps),
		"desired_state_satisfied":           desiredSatisfied,
		"adoption_gaps":                     adoptionGaps,
		"artifact_authority_conflicts":      append(append([]any{}, contractGaps...), liveStateGaps...),
		"contract_gaps":                     contractGaps,
		"higher_authority_conflicts":        conflicts,
		"superseded_signals":                supersededSignals,
		"desired_state_gaps":                desiredGapRecords,
		"live_state_gaps":                   liveStateGaps,
	}
}

func activeRun(run map[string]any) bool {
	state := stringValue(run["state"])
	return state != "completed" && state != "abandoned"
}

func activeMission(mission map[string]any) bool {
	state := stringValue(mission["state"])
	return state != "done" && state != "canceled"
}

func activeSignal(signal map[string]any) bool {
	return stringValue(signal["state"]) == "active"
}

func activeTeammateRun(run map[string]any) bool {
	return activeRun(run) && stringValue(run["owner_type"]) == "teammate"
}

func childRunsForParent(registry map[string]any, parentRunID string) []map[string]any {
	runs := []map[string]any{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if stringValue(run["parent_run_id"]) == parentRunID && activeTeammateRun(run) {
			runs = append(runs, run)
		}
	}
	return runs
}

func defaultDesiredStateForMission() map[string]any {
	return map[string]any{
		"goal_state":          "done",
		"verification_status": "passed",
		"review_status":       "",
		"checks_status":       "",
		"ship_status":         "",
	}
}

func normalizeMissionDesiredState(mission map[string]any) map[string]any {
	desired, _ := toMap(mission["desired_state"])
	if desired == nil {
		desired = map[string]any{}
	}
	defaults := defaultDesiredStateForMission()
	for key, value := range defaults {
		if _, ok := desired[key]; !ok {
			desired[key] = value
		}
	}
	if stringValue(desired["goal_state"]) == "completed" {
		desired["goal_state"] = "done"
	}
	mission["desired_state"] = desired
	return desired
}

func runSortTimestamp(run map[string]any) string {
	return firstNonBlank(
		stringValue(run["updated_at"]),
		stringValue(run["last_reported_at"]),
		stringValue(run["started_at"]),
		stringValue(run["created_at"]),
	)
}

func missionRuns(registry map[string]any, missionID string) []map[string]any {
	runs := []map[string]any{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if stringValue(run["mission_id"]) == missionID {
			runs = append(runs, run)
		}
	}
	sort.Slice(runs, func(i, j int) bool {
		left := runSortTimestamp(runs[i])
		right := runSortTimestamp(runs[j])
		if left != right {
			return left > right
		}
		return stringValue(runs[i]["run_id"]) > stringValue(runs[j]["run_id"])
	})
	return runs
}

func authoritativeMissionRun(runs []map[string]any) map[string]any {
	for _, run := range runs {
		if stringValue(run["owner_type"]) != "teammate" && activeRun(run) {
			return run
		}
	}
	for _, run := range runs {
		if stringValue(run["owner_type"]) != "teammate" {
			return run
		}
	}
	if len(runs) > 0 {
		return runs[0]
	}
	return nil
}

func missionStateFromRuns(runs []map[string]any) string {
	if len(runs) == 0 {
		return "queued"
	}
	anyBlocked := false
	anyReview := false
	anyWIP := false
	anyClaimed := false
	anyCanceled := false
	allCompleted := true
	for _, run := range runs {
		state := stringValue(run["state"])
		if state != "completed" && state != "abandoned" {
			allCompleted = false
		}
		switch state {
		case "blocked":
			anyBlocked = true
		case "ready", "draft":
			anyClaimed = true
		case "executing", "waiting", "paused":
			anyWIP = true
		case "abandoned":
			anyCanceled = true
		}
		if stringValue(run["pr_ref"]) != "" || stringValue(run["review_status"]) == "pending" || stringValue(run["review_status"]) == "changes_requested" || mergeReady(run) {
			anyReview = true
		}
	}
	switch {
	case allCompleted:
		return "done"
	case anyBlocked:
		return "blocked"
	case anyReview:
		return "review"
	case anyWIP:
		return "wip"
	case anyClaimed:
		return "claimed"
	case anyCanceled:
		return "canceled"
	default:
		return "queued"
	}
}

func findMissionIndex(registry map[string]any, missionID string) int {
	for index, item := range toList(registry["missions"]) {
		mission, ok := toMap(item)
		if ok && stringValue(mission["mission_id"]) == missionID {
			return index
		}
	}
	return -1
}

func uniqueMissionIDsFromRuns(registry map[string]any) []string {
	missionIDs := []string{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		missionID := strings.TrimSpace(stringValue(run["mission_id"]))
		if missionID != "" {
			missionIDs = append(missionIDs, missionID)
		}
	}
	return dedupeStrings(missionIDs)
}

func missionDependencyTargetsFromRuns(runs []map[string]any, registry map[string]any) []string {
	targets := []string{}
	for _, run := range runs {
		runMissionID := strings.TrimSpace(stringValue(run["mission_id"]))
		for _, dependencyID := range dependencyTargets(run) {
			if index := findRunIndex(registry, dependencyID); index >= 0 {
				dependencyRun, _ := toMap(toList(registry["runs"])[index])
				if dependencyMissionID := strings.TrimSpace(stringValue(dependencyRun["mission_id"])); dependencyMissionID != "" {
					if dependencyMissionID == runMissionID {
						continue
					}
					targets = append(targets, dependencyMissionID)
					continue
				}
			}
			targets = append(targets, dependencyID)
		}
	}
	return dedupeStrings(targets)
}

func missionDependenciesUnsatisfied(mission map[string]any, missionRegistry map[string]any) []string {
	unsatisfied := []string{}
	for _, dep := range toList(mission["depends_on"]) {
		dependencyID := strings.TrimSpace(stringValue(dep))
		if dependencyID == "" {
			continue
		}
		index := findMissionIndex(missionRegistry, dependencyID)
		if index < 0 {
			unsatisfied = append(unsatisfied, dependencyID)
			continue
		}
		dependencyMission, _ := toMap(toList(missionRegistry["missions"])[index])
		if stringValue(dependencyMission["state"]) != "done" {
			unsatisfied = append(unsatisfied, dependencyID)
		}
	}
	return dedupeStrings(unsatisfied)
}

func missionDependsOnTarget(missionRegistry map[string]any, missionID, targetID string, visited map[string]bool) bool {
	missionID = strings.TrimSpace(missionID)
	targetID = strings.TrimSpace(targetID)
	if missionID == "" || targetID == "" {
		return false
	}
	if missionID == targetID {
		return true
	}
	if visited[missionID] {
		return false
	}
	visited[missionID] = true
	index := findMissionIndex(missionRegistry, missionID)
	if index < 0 {
		return false
	}
	mission, _ := toMap(toList(missionRegistry["missions"])[index])
	for _, dep := range toList(mission["depends_on"]) {
		dependencyID := strings.TrimSpace(stringValue(dep))
		if dependencyID == "" {
			continue
		}
		if dependencyID == targetID || missionDependsOnTarget(missionRegistry, dependencyID, targetID, visited) {
			return true
		}
	}
	return false
}

func meetingActivityByMission(registry map[string]any) (map[string]int, map[string]int, map[string][]string, map[string][]string) {
	activeMeetingCountByMission := map[string]int{}
	openMeetingActionCountByMission := map[string]int{}
	activeMeetingIDsByDecisionID := map[string][]string{}
	activeMeetingIDsByBundleID := map[string][]string{}
	runMissionByID := map[string]string{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		runID := strings.TrimSpace(stringValue(run["run_id"]))
		missionID := strings.TrimSpace(stringValue(run["mission_id"]))
		if runID != "" && missionID != "" {
			runMissionByID[runID] = missionID
		}
	}
	for _, item := range listMeetings(true) {
		meeting, _ := toMap(item)
		missionID := strings.TrimSpace(stringValue(meeting["mission_id"]))
		if missionID == "" {
			missionID = runMissionByID[strings.TrimSpace(stringValue(meeting["run_id"]))]
		}
		if missionID == "" {
			continue
		}
		if stringValue(meeting["state"]) == "active" {
			activeMeetingCountByMission[missionID]++
			if decisionID := strings.TrimSpace(stringValue(meeting["source_decision_id"])); decisionID != "" {
				activeMeetingIDsByDecisionID[decisionID] = dedupeStrings(append(activeMeetingIDsByDecisionID[decisionID], stringValue(meeting["meeting_id"])))
			}
			sourceScope := firstNonBlank(
				strings.TrimSpace(stringValue(meeting["source_decision_scope"])),
				strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_scope", ""))),
			)
			if sourceScope == "bundle" {
				bundleID := firstNonBlank(
					strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.bundle_id", ""))),
					strings.TrimSpace(stringValue(meeting["bundle_id"])),
				)
				if bundleID != "" {
					activeMeetingIDsByBundleID[bundleID] = dedupeStrings(append(activeMeetingIDsByBundleID[bundleID], stringValue(meeting["meeting_id"])))
				}
			}
		}
		for _, actionItem := range toList(meeting["action_items"]) {
			actionMap, ok := toMap(actionItem)
			if ok && stringValue(actionMap["state"]) != "" && stringValue(actionMap["state"]) != "open" {
				continue
			}
			summary := textValue(actionItem)
			if summary == "" || !meetingActionCountsAsBacklog(meeting, registry) {
				continue
			}
			openMeetingActionCountByMission[missionID]++
		}
	}
	return activeMeetingCountByMission, openMeetingActionCountByMission, activeMeetingIDsByDecisionID, activeMeetingIDsByBundleID
}

func explicitBundleGovernorSignals(signalLog map[string]any) map[string][]map[string]any {
	byBundle := map[string][]map[string]any{}
	for _, item := range toList(signalLog["signals"]) {
		signal, ok := toMap(item)
		if !ok || !activeSignal(signal) {
			continue
		}
		if strings.TrimSpace(stringValue(signal["kind"])) != "bundle_governor_established" {
			continue
		}
		bundleID := strings.TrimSpace(stringValue(signal["bundle_id"]))
		governingMissionID := strings.TrimSpace(stringValue(signal["governing_mission_id"]))
		if bundleID == "" || governingMissionID == "" {
			continue
		}
		byBundle[bundleID] = append(byBundle[bundleID], signal)
	}
	return byBundle
}

func activeBundleCouncilHosts() map[string]map[string][]string {
	byBundle := map[string]map[string][]string{}
	for _, item := range listMeetings(true) {
		meeting, ok := toMap(item)
		if !ok || stringValue(meeting["state"]) != "active" {
			continue
		}
		sourceScope := firstNonBlank(
			strings.TrimSpace(stringValue(meeting["source_decision_scope"])),
			strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_scope", ""))),
		)
		if sourceScope != "bundle" {
			continue
		}
		bundleID := firstNonBlank(
			strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.bundle_id", ""))),
			strings.TrimSpace(stringValue(meeting["bundle_id"])),
		)
		hostMissionID := strings.TrimSpace(stringValue(meeting["mission_id"]))
		meetingID := strings.TrimSpace(stringValue(meeting["meeting_id"]))
		if bundleID == "" || hostMissionID == "" || meetingID == "" {
			continue
		}
		packet := byBundle[bundleID]
		if len(packet) == 0 {
			packet = map[string][]string{
				"meeting_ids":      {},
				"host_mission_ids": {},
			}
		}
		packet["meeting_ids"] = dedupeStrings(append(packet["meeting_ids"], meetingID))
		packet["host_mission_ids"] = dedupeStrings(append(packet["host_mission_ids"], hostMissionID))
		byBundle[bundleID] = packet
	}
	return byBundle
}

func resolveActiveBundleGovernorSignals(signalLog map[string]any, bundleID, resolutionNote string) []string {
	signals := toList(signalLog["signals"])
	resolvedIDs := []string{}
	timestamp := nowISO()
	for index, item := range signals {
		signal, ok := toMap(item)
		if !ok || !activeSignal(signal) {
			continue
		}
		if strings.TrimSpace(stringValue(signal["kind"])) != "bundle_governor_established" || strings.TrimSpace(stringValue(signal["bundle_id"])) != bundleID {
			continue
		}
		signal["state"] = "resolved"
		signal["resolution_note"] = resolutionNote
		signal["resolved_at"] = timestamp
		signal["updated_at"] = timestamp
		signals[index] = signal
		resolvedIDs = append(resolvedIDs, stringValue(signal["signal_id"]))
	}
	signalLog["signals"] = signals
	return dedupeStrings(resolvedIDs)
}

func meetingBundleGovernorAdoption(meeting map[string]any) map[string]any {
	sourceScope := firstNonBlank(
		strings.TrimSpace(stringValue(meeting["source_decision_scope"])),
		strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_scope", ""))),
	)
	sourceKind := firstNonBlank(
		strings.TrimSpace(stringValue(meeting["source_decision_kind"])),
		strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_kind", ""))),
	)
	if sourceScope != "bundle" || sourceKind != "establish_bundle_governor" {
		return nil
	}
	packet, _ := toMap(meeting["source_decision_packet"])
	bundleID := firstNonBlank(
		strings.TrimSpace(stringValue(packet["bundle_id"])),
		strings.TrimSpace(stringValue(meeting["bundle_id"])),
	)
	if bundleID == "" {
		return nil
	}
	candidateMissionIDs := []string{}
	for _, item := range toList(packet["candidate_governor_packets"]) {
		candidate, _ := toMap(item)
		if missionID := strings.TrimSpace(stringValue(candidate["mission_id"])); missionID != "" {
			candidateMissionIDs = append(candidateMissionIDs, missionID)
		}
	}
	candidateMissionIDs = append(candidateMissionIDs, anyListToStrings(toList(packet["bundle_mission_ids"]))...)
	candidateMissionIDs = append(candidateMissionIDs, anyListToStrings(toList(packet["mission_ids"]))...)
	if governingMissionID := strings.TrimSpace(stringValue(packet["bundle_governing_mission_id"])); governingMissionID != "" {
		candidateMissionIDs = append(candidateMissionIDs, governingMissionID)
	}
	candidateMissionIDs = dedupeStrings(candidateMissionIDs)

	var selected map[string]any
	selectedRank := -1
	selectedAt := ""
	for _, item := range toList(meeting["locked_decisions"]) {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		governingMissionID := firstNonBlank(
			strings.TrimSpace(stringValue(decision["governing_mission_id"])),
			strings.TrimSpace(stringValue(decision["mission_id"])),
		)
		if governingMissionID == "" {
			continue
		}
		if len(candidateMissionIDs) > 0 && !containsString(candidateMissionIDs, governingMissionID) {
			continue
		}
		rank := authorityRank(firstNonBlank(stringValue(decision["authority_level"]), "L4"))
		lockedAt := firstNonBlank(stringValue(decision["locked_at"]), stringValue(decision["recorded_at"]))
		if selected == nil || rank > selectedRank || (rank == selectedRank && isoAfter(lockedAt, selectedAt)) {
			selected = decision
			selectedRank = rank
			selectedAt = lockedAt
		}
	}
	if selected == nil {
		governingMissionID := strings.TrimSpace(stringValue(packet["bundle_governing_mission_id"]))
		if governingMissionID == "" {
			return nil
		}
		selected = map[string]any{
			"governing_mission_id": governingMissionID,
			"authority_level":      "L2",
			"summary":              fmt.Sprintf("Bundle governor remains anchored to %s.", governingMissionID),
			"reason":               "Scheduler decision packet already names the bundle governor.",
			"evidence_refs":        []any{},
		}
	}
	governingMissionID := firstNonBlank(
		strings.TrimSpace(stringValue(selected["governing_mission_id"])),
		strings.TrimSpace(stringValue(selected["mission_id"])),
	)
	if governingMissionID == "" {
		return nil
	}
	return map[string]any{
		"adoption_kind":         "bundle_governor_established",
		"bundle_id":             bundleID,
		"governing_mission_id":  governingMissionID,
		"authority_level":       firstNonBlank(stringValue(selected["authority_level"]), "L2"),
		"summary":               firstNonBlank(strings.TrimSpace(stringValue(selected["summary"])), fmt.Sprintf("Bundle governance anchored to %s.", governingMissionID)),
		"reason":                strings.TrimSpace(stringValue(selected["reason"])),
		"evidence_refs":         cloneValue(toList(selected["evidence_refs"])),
		"candidate_mission_ids": stringListToAny(candidateMissionIDs),
	}
}

func adoptClosedMeetingAuthority(meeting map[string]any) (map[string]any, error) {
	adoption := meetingBundleGovernorAdoption(meeting)
	if len(adoption) == 0 {
		return nil, nil
	}
	bundleID := stringValue(adoption["bundle_id"])
	governingMissionID := stringValue(adoption["governing_mission_id"])
	if bundleID == "" || governingMissionID == "" {
		return nil, nil
	}

	log := readSignalLog()
	meetingID := stringValue(meeting["meeting_id"])
	resolutionNote := fmt.Sprintf("Superseded by closed council meeting %s adopting %s as bundle governor for %s.", meetingID, governingMissionID, bundleID)
	resolvedSignalIDs := resolveActiveBundleGovernorSignals(log, bundleID, resolutionNote)

	templateRoot, _ := toMap(loadData(signalTemplatePath(), map[string]any{}))
	signalTemplate, _ := toMap(templateRoot["signal"])
	signalID := fmt.Sprintf("sig-%s-bundle-governor-%s", nowTime().UTC().Format("20060102150405"), slugify(bundleID))
	evidenceRefs := append([]string{meetingPath(meetingID, true)}, anyListToStrings(toList(adoption["evidence_refs"]))...)
	signal, _ := toMap(mergeDicts(signalTemplate, map[string]any{
		"signal_id":            signalID,
		"kind":                 "bundle_governor_established",
		"source_type":          "leader",
		"source_id":            firstNonBlank(stringValue(meeting["leader_id"]), "leader-main"),
		"mission_id":           governingMissionID,
		"applies_to":           append(stringListToAny([]string{bundleID}), stringListToAny(anyListToStrings(toList(getKey(meeting, "source_decision_packet.bundle_mission_ids", []any{}))))...),
		"scope_type":           "mission-graph",
		"title":                "Bundle governor adopted from council",
		"summary":              firstNonBlank(stringValue(adoption["summary"]), fmt.Sprintf("Closed council adopted %s as the bundle governor for %s.", governingMissionID, bundleID)),
		"authority_level":      firstNonBlank(stringValue(adoption["authority_level"]), "L2"),
		"confidence":           "high",
		"bundle_id":            bundleID,
		"bundle_mission_ids":   cloneValue(getKey(meeting, "source_decision_packet.bundle_mission_ids", []any{})),
		"governing_mission_id": governingMissionID,
		"selector_refs":        cloneValue(getKey(meeting, "source_decision_packet.selector_refs", []any{})),
		"resolution_kind":      "explicit_bundle_governor",
		"required_action":      fmt.Sprintf("Keep bundle authority anchored to %s until this adopted governor decision is cleared or superseded.", governingMissionID),
		"evidence_refs":        stringListToAny(dedupeStrings(evidenceRefs)),
		"adopted_from_meeting": meetingID,
		"adoption_reason":      stringValue(adoption["reason"]),
	}))
	fillIfBlank(signal, "emitted_at", nowISO())
	fillIfBlank(signal, "state", "active")
	log["signals"] = append(toList(log["signals"]), signal)
	if err := writeSignalLog(log); err != nil {
		return nil, err
	}

	adoption["signal_id"] = stringValue(signal["signal_id"])
	adoption["resolved_prior_signal_ids"] = stringListToAny(resolvedSignalIDs)
	adoption["signal"] = signal
	return adoption, nil
}

func selectExplicitBundleGovernorSignal(signals []map[string]any, componentSet map[string]bool) (map[string]any, []string) {
	candidates := []map[string]any{}
	for _, signal := range signals {
		governingMissionID := strings.TrimSpace(stringValue(signal["governing_mission_id"]))
		if governingMissionID == "" || !componentSet[governingMissionID] {
			continue
		}
		candidates = append(candidates, signal)
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		left := candidates[i]
		right := candidates[j]
		leftAuthority := authorityRank(firstNonBlank(stringValue(left["authority_level"]), "L4"))
		rightAuthority := authorityRank(firstNonBlank(stringValue(right["authority_level"]), "L4"))
		if leftAuthority != rightAuthority {
			return leftAuthority > rightAuthority
		}
		leftEmitted := firstNonBlank(stringValue(left["emitted_at"]), stringValue(left["updated_at"]))
		rightEmitted := firstNonBlank(stringValue(right["emitted_at"]), stringValue(right["updated_at"]))
		if leftEmitted != rightEmitted {
			return isoAfter(leftEmitted, rightEmitted)
		}
		return stringValue(left["signal_id"]) < stringValue(right["signal_id"])
	})
	selected := candidates[0]
	selectedMissionID := strings.TrimSpace(stringValue(selected["governing_mission_id"]))
	selectedAuthority := authorityRank(firstNonBlank(stringValue(selected["authority_level"]), "L4"))
	selectedEmitted := firstNonBlank(stringValue(selected["emitted_at"]), stringValue(selected["updated_at"]))
	notes := []string{}
	supersededSignals := []string{}
	conflictingSignals := []string{}
	for _, signal := range candidates[1:] {
		governingMissionID := strings.TrimSpace(stringValue(signal["governing_mission_id"]))
		authority := authorityRank(firstNonBlank(stringValue(signal["authority_level"]), "L4"))
		emittedAt := firstNonBlank(stringValue(signal["emitted_at"]), stringValue(signal["updated_at"]))
		if governingMissionID == selectedMissionID {
			supersededSignals = append(supersededSignals, stringValue(signal["signal_id"]))
			continue
		}
		if authority < selectedAuthority || (authority == selectedAuthority && isoAfter(selectedEmitted, emittedAt)) {
			supersededSignals = append(supersededSignals, stringValue(signal["signal_id"]))
			continue
		}
		conflictingSignals = append(conflictingSignals, stringValue(signal["signal_id"]))
	}
	if len(supersededSignals) > 0 {
		notes = append(notes, fmt.Sprintf("Explicit bundle governor decision %s supersedes older or lower-authority governor signals: %s.", stringValue(selected["signal_id"]), strings.Join(dedupeStrings(supersededSignals), ", ")))
	}
	if len(conflictingSignals) > 0 {
		notes = append(notes, fmt.Sprintf("Conflicting explicit bundle governor signals remain active alongside %s: %s.", stringValue(selected["signal_id"]), strings.Join(dedupeStrings(conflictingSignals), ", ")))
	}
	return selected, notes
}

func invalidExplicitBundleGovernorSignals(missionRegistry, registry, signalLog map[string]any) []map[string]any {
	graph, _ := toMap(missionGraphAuthorityAudit(missionRegistry, registry, false))
	activeBundleCouncils := activeBundleCouncilHosts()
	bundlesByID := map[string]map[string]any{}
	for _, item := range toList(graph["bundles"]) {
		bundle, ok := toMap(item)
		if !ok {
			continue
		}
		bundleID := strings.TrimSpace(stringValue(bundle["bundle_id"]))
		if bundleID == "" {
			continue
		}
		bundlesByID[bundleID] = bundle
	}
	invalid := []map[string]any{}
	for _, item := range toList(signalLog["signals"]) {
		signal, ok := toMap(item)
		if !ok || !activeSignal(signal) || strings.TrimSpace(stringValue(signal["kind"])) != "bundle_governor_established" {
			continue
		}
		signalID := strings.TrimSpace(stringValue(signal["signal_id"]))
		bundleID := strings.TrimSpace(stringValue(signal["bundle_id"]))
		governingMissionID := strings.TrimSpace(stringValue(signal["governing_mission_id"]))
		reason := ""
		recommendedAction := ""
		switch {
		case bundleID == "" || governingMissionID == "":
			reason = "Explicit bundle-governor signal is missing bundle or governing-mission metadata."
			recommendedAction = "Clear or repair the invalid explicit governor signal before trusting bundle governance."
		default:
			bundle, exists := bundlesByID[bundleID]
			if !exists {
				reason = fmt.Sprintf("Explicit bundle governor still points at %s, but bundle %s no longer exists in the active mission graph.", governingMissionID, bundleID)
				recommendedAction = "Clear the stale explicit governor signal so the control plane stops carrying orphaned governance."
			} else if !containsString(anyListToStrings(toList(bundle["mission_ids"])), governingMissionID) {
				reason = fmt.Sprintf("Explicit bundle governor still points at %s, but that mission no longer belongs to active bundle %s.", governingMissionID, bundleID)
				recommendedAction = "Clear or replace the explicit governor so bundle governance matches the active mission graph."
			} else if councilPacket := activeBundleCouncils[bundleID]; len(councilPacket) > 0 && !containsString(councilPacket["host_mission_ids"], governingMissionID) {
				reason = fmt.Sprintf(
					"Explicit bundle governor points at %s, but active bundle council(s) %s are currently anchored on %s for %s.",
					governingMissionID,
					strings.Join(councilPacket["meeting_ids"], ", "),
					strings.Join(councilPacket["host_mission_ids"], ", "),
					bundleID,
				)
				recommendedAction = "Clear or replace the explicit governor, or re-materialize the bundle council so governance has one consistent anchor."
			}
		}
		if reason == "" {
			continue
		}
		invalid = append(invalid, map[string]any{
			"signal_id":            signalID,
			"bundle_id":            bundleID,
			"governing_mission_id": governingMissionID,
			"mission_id":           strings.TrimSpace(stringValue(signal["mission_id"])),
			"reason":               reason,
			"recommended_action":   recommendedAction,
		})
	}
	return invalid
}

func syncMissionsFromRegistry(registry map[string]any) error {
	templateRoot, _ := toMap(loadData(missionTemplatePath(), map[string]any{}))
	template, _ := toMap(templateRoot["mission"])
	if template == nil {
		template = map[string]any{}
	}

	missionRegistry := readMissionRegistry()
	missionEntries := []any{}
	currentMissionIDs := map[string]bool{}
	for _, missionID := range uniqueMissionIDsFromRuns(registry) {
		mission, err := materializeMission(template, registry, missionID)
		if err != nil {
			return err
		}
		if mission == nil {
			continue
		}
		currentMissionIDs[missionID] = true
		missionEntries = append(missionEntries, mission)
		if err := dumpData(missionPath(missionID), map[string]any{"version": 1, "mission": mission}); err != nil {
			return err
		}
	}
	if existingPaths, err := filepath.Glob(filepath.Join(missionsRoot(), "*.yml")); err == nil {
		for _, path := range existingPaths {
			if path == missionRegistryPath() {
				continue
			}
			missionID := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			if currentMissionIDs[missionID] {
				continue
			}
			if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	sort.Slice(missionEntries, func(i, j int) bool {
		left, _ := toMap(missionEntries[i])
		right, _ := toMap(missionEntries[j])
		leftUpdated := firstNonBlank(stringValue(left["updated_at"]), stringValue(left["created_at"]))
		rightUpdated := firstNonBlank(stringValue(right["updated_at"]), stringValue(right["created_at"]))
		if leftUpdated != rightUpdated {
			return leftUpdated > rightUpdated
		}
		return stringValue(left["mission_id"]) < stringValue(right["mission_id"])
	})
	missionRegistry["missions"] = missionEntries
	return writeMissionRegistry(missionRegistry)
}

func materializeMission(template map[string]any, registry map[string]any, missionID string) (map[string]any, error) {
	runs := missionRuns(registry, missionID)
	if len(runs) == 0 {
		return nil, nil
	}
	latestRun := runs[0]
	authoritativeRun := authoritativeMissionRun(runs)
	if authoritativeRun == nil {
		authoritativeRun = latestRun
	}
	earliestRun := runs[len(runs)-1]
	activeRunIDs := []string{}
	completedRunIDs := []string{}
	teammateRunIDs := []string{}
	ownership := []string{}
	transferEvents := []any{}
	dependsOn := missionDependencyTargetsFromRuns(runs, registry)
	latestSummary := ""
	nextAction := ""
	branch := ""
	worktreePath := ""
	contractRef := ""
	for _, run := range runs {
		if activeRun(run) {
			activeRunIDs = append(activeRunIDs, stringValue(run["run_id"]))
		} else if stringValue(run["state"]) == "completed" {
			completedRunIDs = append(completedRunIDs, stringValue(run["run_id"]))
		}
		if stringValue(run["owner_type"]) == "teammate" {
			teammateRunIDs = append(teammateRunIDs, stringValue(run["run_id"]))
		}
		ownership = append(ownership, ownershipClaims(run)...)
		for _, event := range ownershipTransferEvents(run) {
			transferEvents = append(transferEvents, cloneValue(event))
		}
		if latestSummary == "" {
			latestSummary = stringValue(run["latest_summary"])
		}
		if nextAction == "" {
			nextAction = stringValue(run["next_expected_action"])
		}
		if branch == "" {
			branch = stringValue(run["branch"])
		}
		if worktreePath == "" {
			worktreePath = stringValue(run["worktree_path"])
		}
		if contractRef == "" {
			contractRef = stringValue(run["contract_ref"])
		}
	}
	controlSummary := firstNonBlank(stringValue(authoritativeRun["latest_summary"]), latestSummary)
	controlNextAction := firstNonBlank(stringValue(authoritativeRun["next_expected_action"]), nextAction)
	controlBranch := firstNonBlank(stringValue(authoritativeRun["branch"]), branch)
	controlWorktreePath := firstNonBlank(stringValue(authoritativeRun["worktree_path"]), worktreePath)
	controlContractRef := firstNonBlank(stringValue(authoritativeRun["contract_ref"]), contractRef)
	contract := map[string]any{}
	if controlContractRef != "" {
		contractRoot, _ := toMap(loadData(controlContractRef, map[string]any{}))
		contract, _ = toMap(contractRoot["contract"])
	}
	mission, _ := toMap(mergeDicts(template, map[string]any{
		"mission_id":            missionID,
		"issue_number":          intValue(firstNonBlank(stringValue(authoritativeRun["issue_number"]), stringValue(latestRun["issue_number"]), stringValue(contract["issue_number"]))),
		"title":                 firstNonBlank(stringValue(contract["title"]), stringValue(authoritativeRun["objective"]), stringValue(latestRun["objective"]), missionID),
		"objective":             firstNonBlank(stringValue(contract["objective"]), stringValue(authoritativeRun["objective"]), stringValue(latestRun["objective"]), missionID),
		"state":                 missionStateFromRuns(runs),
		"priority":              firstNonBlank(stringValue(authoritativeRun["priority"]), stringValue(latestRun["priority"])),
		"risk_level":            firstNonBlank(stringValue(authoritativeRun["risk_level"]), stringValue(latestRun["risk_level"])),
		"branch":                firstNonBlank(stringValue(contract["branch"]), controlBranch),
		"worktree_path":         firstNonBlank(stringValue(contract["worktree_path"]), controlWorktreePath),
		"contract_ref":          controlContractRef,
		"success_criteria":      cloneValue(toList(contract["success_criteria"])),
		"verification_commands": cloneValue(toList(contract["verification_commands"])),
		"context_files":         cloneValue(toList(contract["context_files"])),
		"depends_on":            stringListToAny(dedupeStrings(dependsOn)),
		"ownership_claims":      stringListToAny(dedupeStrings(ownership)),
		"ownership_selectors":   selectorRecordsFromClaims(ownership),
		"ownership_transfers":   transferEvents,
		"active_run_ids":        stringListToAny(dedupeStrings(activeRunIDs)),
		"completed_run_ids":     stringListToAny(dedupeStrings(completedRunIDs)),
		"teammate_run_ids":      stringListToAny(dedupeStrings(teammateRunIDs)),
		"latest_run_id":         stringValue(latestRun["run_id"]),
		"authoritative_run_id":  stringValue(authoritativeRun["run_id"]),
		"latest_summary":        controlSummary,
		"next_expected_action":  controlNextAction,
		"created_at":            firstNonBlank(stringValue(getKey(contract, "freshness.created_at", "")), stringValue(earliestRun["created_at"]), stringValue(earliestRun["started_at"])),
		"updated_at":            firstNonBlank(runSortTimestamp(latestRun), nowISO()),
		"completed_at":          stringValue(latestRun["completed_at"]),
		"run_count":             len(runs),
		"freshness":             cloneValue(getKey(contract, "freshness", map[string]any{})),
	}))
	if desired, ok := toMap(contract["desired_state"]); ok && len(desired) > 0 {
		mission["desired_state"] = cloneValue(desired)
	} else if desired, ok := toMap(latestRun["desired_state"]); ok && len(desired) > 0 {
		mission["desired_state"] = cloneValue(desired)
	}
	normalizeMissionDesiredState(mission)
	if stringValue(mission["state"]) != "done" {
		mission["completed_at"] = ""
	} else if stringValue(mission["completed_at"]) == "" {
		mission["completed_at"] = firstNonBlank(stringValue(latestRun["completed_at"]), nowISO())
	}
	return mission, nil
}

func missionAuthorityAudit(mission map[string]any, registry map[string]any) map[string]any {
	missionID := stringValue(mission["mission_id"])
	runs := missionRuns(registry, missionID)
	if len(runs) == 0 {
		return map[string]any{
			"mission_id":                missionID,
			"canonical_source":          "mission:" + missionID,
			"mission_gap_count":         0,
			"adjudication_count":        0,
			"unresolved_conflict_count": 0,
			"mission_gaps":              []any{},
			"adjudications":             []any{},
			"unresolved_conflicts":      []any{},
			"latest_run_id":             "",
			"authoritative_run_id":      "",
			"active_root_run_count":     0,
			"active_teammate_run_count": 0,
		}
	}

	templateRoot, _ := toMap(loadData(missionTemplatePath(), map[string]any{}))
	template, _ := toMap(templateRoot["mission"])
	canonicalMission, err := materializeMission(template, registry, missionID)
	if err != nil || canonicalMission == nil {
		canonicalMission = map[string]any{}
	}

	latestRun := runs[0]
	authoritativeRun := authoritativeMissionRun(runs)
	if authoritativeRun == nil {
		authoritativeRun = latestRun
	}

	missionGaps := []any{}
	appendStringGap := func(field string, recorded string, expected string, note string) {
		if strings.TrimSpace(recorded) == strings.TrimSpace(expected) {
			return
		}
		missionGaps = append(missionGaps, map[string]any{
			"claim_ref":       field,
			"resolution_kind": "run_aggregate_overrides_mission_record",
			"source_ref":      "mission:" + missionID + "." + field,
			"winning_source":  "mission-canonical:" + missionID,
			"losing_source":   "mission:" + missionID + "." + field,
			"authority_level": "L2",
			"current_value":   recorded,
			"expected_value":  expected,
			"note":            note,
		})
	}
	appendListGap := func(field string, recorded []string, expected []string, note string) {
		left := normalizedStringSet(recorded)
		right := normalizedStringSet(expected)
		if strings.Join(left, "|") == strings.Join(right, "|") {
			return
		}
		missionGaps = append(missionGaps, map[string]any{
			"claim_ref":       field,
			"resolution_kind": "run_aggregate_overrides_mission_record",
			"source_ref":      "mission:" + missionID + "." + field,
			"winning_source":  "mission-canonical:" + missionID,
			"losing_source":   "mission:" + missionID + "." + field,
			"authority_level": "L2",
			"current_value":   strings.Join(left, "|"),
			"expected_value":  strings.Join(right, "|"),
			"note":            note,
		})
	}

	appendStringGap("state", stringValue(mission["state"]), stringValue(canonicalMission["state"]), "Mission state no longer matches the run-derived control plane.")
	appendStringGap("branch", stringValue(mission["branch"]), stringValue(canonicalMission["branch"]), "Mission branch no longer matches the authoritative mission run.")
	appendStringGap("worktree_path", stringValue(mission["worktree_path"]), stringValue(canonicalMission["worktree_path"]), "Mission worktree path no longer matches the authoritative mission run.")
	appendStringGap("latest_run_id", stringValue(mission["latest_run_id"]), stringValue(canonicalMission["latest_run_id"]), "Mission latest_run_id no longer matches the freshest mission run.")
	appendStringGap("authoritative_run_id", stringValue(mission["authoritative_run_id"]), stringValue(canonicalMission["authoritative_run_id"]), "Mission authoritative_run_id no longer matches the control-plane governing run.")
	appendStringGap("latest_summary", stringValue(mission["latest_summary"]), stringValue(canonicalMission["latest_summary"]), "Mission latest_summary no longer matches the authoritative run summary.")
	appendStringGap("next_expected_action", stringValue(mission["next_expected_action"]), stringValue(canonicalMission["next_expected_action"]), "Mission next_expected_action no longer matches the authoritative run.")
	appendListGap("active_run_ids", anyListToStrings(toList(mission["active_run_ids"])), anyListToStrings(toList(canonicalMission["active_run_ids"])), "Mission active_run_ids no longer match the active runs in registry.")
	appendListGap("ownership_claims", ownershipClaims(mission), ownershipClaims(canonicalMission), "Mission ownership claims no longer match the aggregated run ownership boundary.")

	missionDesired := normalizeMissionDesiredState(mission)
	canonicalDesired, _ := toMap(canonicalMission["desired_state"])
	for _, field := range []string{"goal_state", "verification_status", "review_status", "checks_status", "ship_status"} {
		appendStringGap("desired_state."+field, stringValue(missionDesired[field]), stringValue(canonicalDesired[field]), fmt.Sprintf("Mission desired state for %s no longer matches the authoritative mission control plane.", field))
	}

	adjudications := []any{}
	if stringValue(latestRun["run_id"]) != "" && stringValue(authoritativeRun["run_id"]) != "" && stringValue(latestRun["run_id"]) != stringValue(authoritativeRun["run_id"]) {
		resolutionKind := "authoritative_run_controls_mission"
		note := "A non-authoritative sibling run updated more recently, but mission control remains governed by the authoritative run."
		if stringValue(latestRun["owner_type"]) == "teammate" && stringValue(authoritativeRun["owner_type"]) != "teammate" {
			resolutionKind = "leader_run_controls_mission"
			note = "A teammate run updated more recently, but mission control remains governed by the leader-facing run."
		}
		adjudications = append(adjudications, map[string]any{
			"claim_ref":       "mission_control_run",
			"resolution_kind": resolutionKind,
			"source_ref":      "run:" + stringValue(authoritativeRun["run_id"]),
			"winning_source":  "run:" + stringValue(authoritativeRun["run_id"]),
			"losing_source":   "run:" + stringValue(latestRun["run_id"]),
			"authority_level": "L2",
			"current_value":   stringValue(latestRun["run_id"]),
			"expected_value":  stringValue(authoritativeRun["run_id"]),
			"note":            note,
		})
	}

	activeRootRuns := []map[string]any{}
	activeTeammateCount := 0
	for _, run := range runs {
		if activeRun(run) {
			if stringValue(run["owner_type"]) == "teammate" {
				activeTeammateCount++
			} else {
				activeRootRuns = append(activeRootRuns, run)
			}
		}
	}
	unresolvedConflicts := []any{}
	if len(activeRootRuns) > 1 {
		for _, run := range activeRootRuns {
			if stringValue(run["run_id"]) == stringValue(authoritativeRun["run_id"]) {
				continue
			}
			unresolvedConflicts = append(unresolvedConflicts, map[string]any{
				"claim_ref":       "mission_control_run",
				"resolution_kind": "parallel_root_runs_require_resolution",
				"source_ref":      "run:" + stringValue(authoritativeRun["run_id"]),
				"winning_source":  "run:" + stringValue(authoritativeRun["run_id"]),
				"losing_source":   "run:" + stringValue(run["run_id"]),
				"authority_level": "L2",
				"current_value":   stringValue(run["run_id"]),
				"expected_value":  stringValue(authoritativeRun["run_id"]),
				"note":            "Multiple active non-teammate runs exist for the same mission; the leader must close or merge the competing root run.",
			})
		}
	}

	return map[string]any{
		"mission_id":                missionID,
		"canonical_source":          "mission:" + missionID,
		"latest_run_id":             stringValue(latestRun["run_id"]),
		"authoritative_run_id":      stringValue(authoritativeRun["run_id"]),
		"mission_gap_count":         len(missionGaps),
		"adjudication_count":        len(adjudications),
		"unresolved_conflict_count": len(unresolvedConflicts),
		"active_root_run_count":     len(activeRootRuns),
		"active_teammate_run_count": activeTeammateCount,
		"mission_gaps":              missionGaps,
		"adjudications":             adjudications,
		"unresolved_conflicts":      unresolvedConflicts,
	}
}

func missionGraphAuthorityAudit(missionRegistry map[string]any, registry map[string]any, includeInactive bool) map[string]any {
	signalLog := readSignalLog()
	explicitGovernorSignalsByBundle := explicitBundleGovernorSignals(signalLog)
	missions := []map[string]any{}
	missionByID := map[string]map[string]any{}
	for _, item := range toList(missionRegistry["missions"]) {
		mission, _ := toMap(item)
		if len(mission) == 0 {
			continue
		}
		if includeInactive || activeMission(mission) {
			missions = append(missions, mission)
			missionID := strings.TrimSpace(stringValue(mission["mission_id"]))
			if missionID != "" {
				missionByID[missionID] = mission
			}
		}
	}
	sort.Slice(missions, func(i, j int) bool {
		left := stringValue(missions[i]["mission_id"])
		right := stringValue(missions[j]["mission_id"])
		return left < right
	})
	activeMeetingCountByMission, openMeetingActionCountByMission, _, _ := meetingActivityByMission(registry)

	missionAudits := []any{}
	missionIDs := []string{}
	missionIDSet := map[string]bool{}
	missionGapCount := 0
	missionAdjudicationCount := 0
	missionUnresolvedCount := 0
	for _, mission := range missions {
		audit, _ := toMap(missionAuthorityAudit(mission, registry))
		missionAudits = append(missionAudits, audit)
		missionID := stringValue(mission["mission_id"])
		missionIDs = append(missionIDs, missionID)
		missionIDSet[missionID] = true
		missionGapCount += intValue(audit["mission_gap_count"])
		missionAdjudicationCount += intValue(audit["adjudication_count"])
		missionUnresolvedCount += intValue(audit["unresolved_conflict_count"])
	}

	crossMissionConflicts := []any{}
	crossMissionAdjudications := []any{}
	bundleAdjacency := map[string]map[string]bool{}
	bundlePairKinds := map[string][]string{}
	addBundleEdge := func(leftID, rightID, kind string) {
		leftID = strings.TrimSpace(leftID)
		rightID = strings.TrimSpace(rightID)
		kind = strings.TrimSpace(kind)
		if leftID == "" || rightID == "" || leftID == rightID || kind == "" {
			return
		}
		if bundleAdjacency[leftID] == nil {
			bundleAdjacency[leftID] = map[string]bool{}
		}
		if bundleAdjacency[rightID] == nil {
			bundleAdjacency[rightID] = map[string]bool{}
		}
		bundleAdjacency[leftID][rightID] = true
		bundleAdjacency[rightID][leftID] = true
		ordered := dedupeStrings([]string{leftID, rightID})
		sort.Strings(ordered)
		if len(ordered) != 2 {
			return
		}
		pairKey := ordered[0] + "::" + ordered[1]
		bundlePairKinds[pairKey] = dedupeStrings(append(bundlePairKinds[pairKey], kind))
	}
	for index, leftMission := range missions {
		if !activeMission(leftMission) {
			continue
		}
		leftID := stringValue(leftMission["mission_id"])
		for _, rightMission := range missions[index+1:] {
			if !activeMission(rightMission) {
				continue
			}
			rightID := stringValue(rightMission["mission_id"])
			leftDependsOnRight := missionDependsOnTarget(missionRegistry, leftID, rightID, map[string]bool{})
			rightDependsOnLeft := missionDependsOnTarget(missionRegistry, rightID, leftID, map[string]bool{})
			shared := []string{}
			decisiveWinnerID := ""
			decisiveSelectorRefs := []string{}
			unresolvedSelectorRefs := []string{}
			leftExclusions := ownershipExclusions(leftMission)
			rightExclusions := ownershipExclusions(rightMission)
			for _, leftClaim := range ownershipClaims(leftMission) {
				for _, rightClaim := range ownershipClaims(rightMission) {
					if selectorsConflict(leftClaim, leftExclusions, rightClaim, rightExclusions) {
						selectorRef := leftClaim + "~" + rightClaim
						shared = append(shared, selectorRef)
						if winner, ok := moreSpecificSelectorWinner(leftClaim, rightClaim); ok {
							winnerMissionID := leftID
							if winner == "right" {
								winnerMissionID = rightID
							}
							if decisiveWinnerID == "" {
								decisiveWinnerID = winnerMissionID
							}
							if decisiveWinnerID == winnerMissionID {
								decisiveSelectorRefs = append(decisiveSelectorRefs, selectorRef)
							} else {
								unresolvedSelectorRefs = append(unresolvedSelectorRefs, selectorRef)
							}
							continue
						}
						unresolvedSelectorRefs = append(unresolvedSelectorRefs, selectorRef)
					}
				}
			}
			if len(shared) == 0 {
				if leftDependsOnRight || rightDependsOnLeft {
					addBundleEdge(leftID, rightID, "dependency")
				}
				continue
			}
			addBundleEdge(leftID, rightID, "ownership_overlap")
			switch {
			case leftDependsOnRight && !rightDependsOnLeft:
				addBundleEdge(leftID, rightID, "dependency")
				addBundleEdge(leftID, rightID, "dependency_governed_ownership")
				crossMissionAdjudications = append(crossMissionAdjudications, map[string]any{
					"claim_ref":       "mission_ownership_boundary",
					"resolution_kind": "dependency_governed_mission_ownership",
					"source_ref":      "mission:" + rightID,
					"winning_source":  "mission:" + rightID,
					"losing_source":   "mission:" + leftID,
					"authority_level": "L2",
					"mission_ids":     stringListToAny([]string{leftID, rightID}),
					"selector_refs":   stringListToAny(dedupeStrings(shared)),
					"note":            "Dependent mission may overlap the governing mission only as a bounded dependency; the dependency target keeps authority over the shared ownership boundary.",
				})
				continue
			case rightDependsOnLeft && !leftDependsOnRight:
				addBundleEdge(leftID, rightID, "dependency")
				addBundleEdge(leftID, rightID, "dependency_governed_ownership")
				crossMissionAdjudications = append(crossMissionAdjudications, map[string]any{
					"claim_ref":       "mission_ownership_boundary",
					"resolution_kind": "dependency_governed_mission_ownership",
					"source_ref":      "mission:" + leftID,
					"winning_source":  "mission:" + leftID,
					"losing_source":   "mission:" + rightID,
					"authority_level": "L2",
					"mission_ids":     stringListToAny([]string{leftID, rightID}),
					"selector_refs":   stringListToAny(dedupeStrings(shared)),
					"note":            "Dependent mission may overlap the governing mission only as a bounded dependency; the dependency target keeps authority over the shared ownership boundary.",
				})
				continue
			}
			leftCouncilAnchored := activeMeetingCountByMission[leftID] > 0 || openMeetingActionCountByMission[leftID] > 0
			rightCouncilAnchored := activeMeetingCountByMission[rightID] > 0 || openMeetingActionCountByMission[rightID] > 0
			switch {
			case leftCouncilAnchored && !rightCouncilAnchored:
				addBundleEdge(leftID, rightID, "council")
				addBundleEdge(leftID, rightID, "council_governed_ownership")
				crossMissionAdjudications = append(crossMissionAdjudications, map[string]any{
					"claim_ref":       "mission_ownership_boundary",
					"resolution_kind": "council_governed_mission_ownership",
					"source_ref":      "mission:" + leftID,
					"winning_source":  "mission:" + leftID,
					"losing_source":   "mission:" + rightID,
					"authority_level": "L2",
					"mission_ids":     stringListToAny([]string{leftID, rightID}),
					"selector_refs":   stringListToAny(dedupeStrings(shared)),
					"note":            "An active council already governs this overlapping ownership boundary; keep execution anchored to the council-host mission until the meeting outcome is adopted.",
				})
				continue
			case rightCouncilAnchored && !leftCouncilAnchored:
				addBundleEdge(leftID, rightID, "council")
				addBundleEdge(leftID, rightID, "council_governed_ownership")
				crossMissionAdjudications = append(crossMissionAdjudications, map[string]any{
					"claim_ref":       "mission_ownership_boundary",
					"resolution_kind": "council_governed_mission_ownership",
					"source_ref":      "mission:" + rightID,
					"winning_source":  "mission:" + rightID,
					"losing_source":   "mission:" + leftID,
					"authority_level": "L2",
					"mission_ids":     stringListToAny([]string{leftID, rightID}),
					"selector_refs":   stringListToAny(dedupeStrings(shared)),
					"note":            "An active council already governs this overlapping ownership boundary; keep execution anchored to the council-host mission until the meeting outcome is adopted.",
				})
				continue
			}
			if leftDependsOnRight || rightDependsOnLeft {
				addBundleEdge(leftID, rightID, "dependency")
			}
			if len(unresolvedSelectorRefs) == 0 && decisiveWinnerID != "" && len(decisiveSelectorRefs) > 0 {
				losingMissionID := rightID
				if decisiveWinnerID == rightID {
					losingMissionID = leftID
				}
				addBundleEdge(leftID, rightID, "selector_specificity_governed_ownership")
				crossMissionAdjudications = append(crossMissionAdjudications, map[string]any{
					"claim_ref":       "mission_ownership_boundary",
					"resolution_kind": "selector_specificity_governed_mission_ownership",
					"source_ref":      "mission:" + decisiveWinnerID,
					"winning_source":  "mission:" + decisiveWinnerID,
					"losing_source":   "mission:" + losingMissionID,
					"authority_level": "L2",
					"mission_ids":     stringListToAny([]string{leftID, rightID}),
					"selector_refs":   stringListToAny(dedupeStrings(decisiveSelectorRefs)),
					"note":            "The more specific mission ownership selector governs the broader overlapping selector until the boundary is narrowed or explicitly transferred.",
				})
				continue
			}
			addBundleEdge(leftID, rightID, "parallel_conflict")
			crossMissionConflicts = append(crossMissionConflicts, map[string]any{
				"claim_ref":       "mission_ownership_boundary",
				"resolution_kind": "parallel_mission_ownership_conflict",
				"source_ref":      "mission:" + leftID,
				"winning_source":  "",
				"losing_source":   "",
				"authority_level": "L2",
				"mission_ids":     stringListToAny([]string{leftID, rightID}),
				"selector_refs":   stringListToAny(dedupeStrings(shared)),
				"note":            "Active missions have overlapping bounded ownership selectors without a higher-level adjudication.",
			})
		}
	}

	dependencyAdjudications := []any{}
	for _, mission := range missions {
		if !activeMission(mission) {
			continue
		}
		missionID := stringValue(mission["mission_id"])
		for _, dependencyID := range missionDependenciesUnsatisfied(mission, missionRegistry) {
			if missionIDSet[dependencyID] {
				addBundleEdge(missionID, dependencyID, "dependency")
				addBundleEdge(missionID, dependencyID, "dependency_block")
			}
			dependencyAdjudications = append(dependencyAdjudications, map[string]any{
				"claim_ref":       "mission_dependency",
				"resolution_kind": "dependency_mission_blocks_execution",
				"source_ref":      "mission:" + dependencyID,
				"winning_source":  "mission:" + dependencyID,
				"losing_source":   "mission:" + missionID,
				"authority_level": "L2",
				"mission_id":      missionID,
				"dependency_id":   dependencyID,
				"note":            "Dependent mission cannot govern itself as ready-for-unconstrained execution until its dependency mission is done.",
			})
		}
	}

	bundles := []any{}
	visitedBundles := map[string]bool{}
	conflictedBundleCount := 0
	governedBundleCount := 0
	escalatedBundleCount := 0
	for _, missionID := range missionIDs {
		if visitedBundles[missionID] || len(bundleAdjacency[missionID]) == 0 {
			continue
		}
		queue := []string{missionID}
		visitedBundles[missionID] = true
		component := []string{}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			component = append(component, current)
			neighbors := []string{}
			for neighbor := range bundleAdjacency[current] {
				neighbors = append(neighbors, neighbor)
			}
			sort.Strings(neighbors)
			for _, neighbor := range neighbors {
				if visitedBundles[neighbor] {
					continue
				}
				visitedBundles[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
		sort.Strings(component)
		bundleID := "bundle-" + slugify(strings.Join(component, "-"))
		componentSet := map[string]bool{}
		for _, id := range component {
			componentSet[id] = true
		}
		bundleRoots := []string{}
		dependentMissionIDs := []string{}
		for _, id := range component {
			index := findMissionIndex(missionRegistry, id)
			if index < 0 {
				continue
			}
			mission, _ := toMap(toList(missionRegistry["missions"])[index])
			internalDependencyCount := 0
			for _, dep := range toList(mission["depends_on"]) {
				dependencyID := strings.TrimSpace(stringValue(dep))
				if dependencyID != "" && componentSet[dependencyID] {
					internalDependencyCount++
				}
			}
			if internalDependencyCount == 0 {
				bundleRoots = append(bundleRoots, id)
			} else {
				dependentMissionIDs = append(dependentMissionIDs, id)
			}
		}
		sort.Strings(bundleRoots)
		sort.Strings(dependentMissionIDs)
		councilAnchorMissionIDs := []string{}
		for _, id := range component {
			if activeMeetingCountByMission[id] > 0 || openMeetingActionCountByMission[id] > 0 {
				councilAnchorMissionIDs = append(councilAnchorMissionIDs, id)
			}
		}
		sort.Strings(councilAnchorMissionIDs)
		candidateGovernorMissionIDs := dedupeStrings(bundleRoots)
		if len(candidateGovernorMissionIDs) == 0 {
			candidateGovernorMissionIDs = dedupeStrings(component)
		}
		candidateGovernorPackets := []any{}
		for _, candidateMissionID := range candidateGovernorMissionIDs {
			mission := missionByID[candidateMissionID]
			releaseAdjacent := false
			activeRunCount := 0
			for _, run := range missionRuns(registry, candidateMissionID) {
				if activeRun(run) {
					activeRunCount++
				}
				if runReleaseAdjacent(run) {
					releaseAdjacent = true
				}
			}
			candidateGovernorPackets = append(candidateGovernorPackets, map[string]any{
				"mission_id":           candidateMissionID,
				"priority":             stringValue(mission["priority"]),
				"risk_level":           stringValue(mission["risk_level"]),
				"release_adjacent":     releaseAdjacent,
				"active_meeting_count": activeMeetingCountByMission[candidateMissionID],
				"open_meeting_actions": openMeetingActionCountByMission[candidateMissionID],
				"active_run_count":     activeRunCount,
				"is_root_candidate":    containsString(bundleRoots, candidateMissionID),
			})
		}
		sort.Slice(candidateGovernorPackets, func(i, j int) bool {
			left, _ := toMap(candidateGovernorPackets[i])
			right, _ := toMap(candidateGovernorPackets[j])
			leftPriority := schedulerPriorityRank(stringValue(left["priority"]))
			rightPriority := schedulerPriorityRank(stringValue(right["priority"]))
			if leftPriority != rightPriority {
				return leftPriority > rightPriority
			}
			leftRisk := schedulerRiskRank(stringValue(left["risk_level"]))
			rightRisk := schedulerRiskRank(stringValue(right["risk_level"]))
			if leftRisk != rightRisk {
				return leftRisk > rightRisk
			}
			if boolValue(left["release_adjacent"]) != boolValue(right["release_adjacent"]) {
				return boolValue(left["release_adjacent"])
			}
			if boolValue(left["is_root_candidate"]) != boolValue(right["is_root_candidate"]) {
				return boolValue(left["is_root_candidate"])
			}
			return stringValue(left["mission_id"]) < stringValue(right["mission_id"])
		})
		candidateGovernorMissionIDs = []string{}
		for _, item := range candidateGovernorPackets {
			packet, _ := toMap(item)
			candidateGovernorMissionIDs = append(candidateGovernorMissionIDs, stringValue(packet["mission_id"]))
		}
		edgeKinds := []string{}
		edgeCount := 0
		bundleSelectorRefs := []string{}
		resolutionKinds := []string{}
		governingMissionVotes := map[string]int{}
		for pairKey, kinds := range bundlePairKinds {
			pair := strings.Split(pairKey, "::")
			if len(pair) != 2 || !componentSet[pair[0]] || !componentSet[pair[1]] {
				continue
			}
			edgeCount++
			edgeKinds = append(edgeKinds, kinds...)
		}
		conflictCount := 0
		for _, item := range crossMissionConflicts {
			conflict, _ := toMap(item)
			inBundle := true
			for _, id := range anyListToStrings(toList(conflict["mission_ids"])) {
				if !componentSet[id] {
					inBundle = false
					break
				}
			}
			if inBundle {
				conflictCount++
				resolutionKinds = append(resolutionKinds, stringValue(conflict["resolution_kind"]))
				bundleSelectorRefs = append(bundleSelectorRefs, anyListToStrings(toList(conflict["selector_refs"]))...)
			}
		}
		adjudicationCount := 0
		for _, item := range crossMissionAdjudications {
			adjudication, _ := toMap(item)
			inBundle := true
			for _, id := range anyListToStrings(toList(adjudication["mission_ids"])) {
				if !componentSet[id] {
					inBundle = false
					break
				}
			}
			if inBundle {
				adjudicationCount++
				resolutionKinds = append(resolutionKinds, stringValue(adjudication["resolution_kind"]))
				bundleSelectorRefs = append(bundleSelectorRefs, anyListToStrings(toList(adjudication["selector_refs"]))...)
				if winningMissionID := strings.TrimPrefix(stringValue(adjudication["winning_source"]), "mission:"); componentSet[winningMissionID] {
					governingMissionVotes[winningMissionID]++
				}
			}
		}
		for _, item := range dependencyAdjudications {
			adjudication, _ := toMap(item)
			missionID := stringValue(adjudication["mission_id"])
			dependencyID := stringValue(adjudication["dependency_id"])
			if componentSet[missionID] && componentSet[dependencyID] {
				adjudicationCount++
				resolutionKinds = append(resolutionKinds, stringValue(adjudication["resolution_kind"]))
			}
		}
		state := "linked"
		governingMissionID := ""
		governedMissionIDs := []string{}
		losingMissionIDs := []string{}
		authorityBasis := []string{}
		resolutionKind := "linked_bundle"
		requiredAction := "Inspect linked missions before widening parallel execution."
		governanceNotes := []string{}
		explicitGovernorSignalID := ""
		governanceSourceKind := ""
		adoptedFromMeetingID := ""
		if len(governingMissionVotes) == 1 {
			for missionID := range governingMissionVotes {
				governingMissionID = missionID
			}
		}
		switch {
		case conflictCount > 0:
			candidateMissionIDs := candidateGovernorMissionIDs

			selectedWinnerID := ""
			hasBundleAuthoritySignal := false
			explicitGovernorSignal, explicitGovernorNotes := selectExplicitBundleGovernorSignal(explicitGovernorSignalsByBundle[bundleID], componentSet)
			governanceNotes = append(governanceNotes, explicitGovernorNotes...)
			if explicitGovernorSignal != nil {
				selectedWinnerID = strings.TrimSpace(stringValue(explicitGovernorSignal["governing_mission_id"]))
				explicitGovernorSignalID = stringValue(explicitGovernorSignal["signal_id"])
				hasBundleAuthoritySignal = true
				adoptedFromMeetingID = strings.TrimSpace(stringValue(explicitGovernorSignal["adopted_from_meeting"]))
				if adoptedFromMeetingID != "" {
					governanceSourceKind = "council-adopted-governor"
					authorityBasis = append(authorityBasis, "council_adopted_bundle_governor")
				} else {
					governanceSourceKind = "explicit-governor"
					authorityBasis = append(authorityBasis, "explicit_bundle_governor")
				}
			}
			maxVoteCount := 0
			voteLeaders := []string{}
			if selectedWinnerID == "" {
				for _, candidateMissionID := range candidateMissionIDs {
					votes := governingMissionVotes[candidateMissionID]
					if votes > maxVoteCount {
						maxVoteCount = votes
						voteLeaders = []string{candidateMissionID}
					} else if votes > 0 && votes == maxVoteCount {
						voteLeaders = append(voteLeaders, candidateMissionID)
					}
				}
				if maxVoteCount > 0 {
					hasBundleAuthoritySignal = true
				}
				if maxVoteCount > 0 && len(dedupeStrings(voteLeaders)) == 1 {
					selectedWinnerID = voteLeaders[0]
					authorityBasis = append(authorityBasis, "bundle_adjudication_votes")
				}
			}

			if selectedWinnerID == "" && len(candidateMissionIDs) >= 2 {
				type bundleAuthorityCandidate struct {
					missionID          string
					priority           int
					risk               int
					releaseAdjacent    bool
					openMeetingActions int
					activeMeetingCount int
					isRootCandidate    bool
				}
				candidates := []bundleAuthorityCandidate{}
				for _, item := range candidateGovernorPackets {
					packet, _ := toMap(item)
					candidates = append(candidates, bundleAuthorityCandidate{
						missionID:          stringValue(packet["mission_id"]),
						priority:           schedulerPriorityRank(stringValue(packet["priority"])),
						risk:               schedulerRiskRank(stringValue(packet["risk_level"])),
						releaseAdjacent:    boolValue(packet["release_adjacent"]),
						openMeetingActions: intValue(packet["open_meeting_actions"]),
						activeMeetingCount: intValue(packet["active_meeting_count"]),
						isRootCandidate:    boolValue(packet["is_root_candidate"]),
					})
				}
				sort.Slice(candidates, func(i, j int) bool {
					if candidates[i].priority != candidates[j].priority {
						return candidates[i].priority > candidates[j].priority
					}
					if candidates[i].risk != candidates[j].risk {
						return candidates[i].risk > candidates[j].risk
					}
					if candidates[i].releaseAdjacent != candidates[j].releaseAdjacent {
						return candidates[i].releaseAdjacent
					}
					if candidates[i].openMeetingActions != candidates[j].openMeetingActions {
						return candidates[i].openMeetingActions > candidates[j].openMeetingActions
					}
					if candidates[i].activeMeetingCount != candidates[j].activeMeetingCount {
						return candidates[i].activeMeetingCount > candidates[j].activeMeetingCount
					}
					if candidates[i].isRootCandidate != candidates[j].isRootCandidate {
						return candidates[i].isRootCandidate
					}
					return candidates[i].missionID < candidates[j].missionID
				})
				best := candidates[0]
				second := candidates[0]
				if len(candidates) > 1 {
					second = candidates[1]
				}
				hasAuthoritySignal := best.priority > 0 ||
					best.risk > 0 ||
					best.releaseAdjacent ||
					best.openMeetingActions > 0 ||
					best.activeMeetingCount > 0
				hasUniqueTopTuple := len(candidates) == 1 ||
					best.priority != second.priority ||
					best.risk != second.risk ||
					best.releaseAdjacent != second.releaseAdjacent ||
					best.openMeetingActions != second.openMeetingActions ||
					best.activeMeetingCount != second.activeMeetingCount ||
					best.isRootCandidate != second.isRootCandidate
				if hasAuthoritySignal {
					hasBundleAuthoritySignal = true
				}
				if hasAuthoritySignal && hasUniqueTopTuple {
					selectedWinnerID = best.missionID
					if best.priority != second.priority {
						authorityBasis = append(authorityBasis, "priority")
					}
					if best.risk != second.risk {
						authorityBasis = append(authorityBasis, "risk_level")
					}
					if best.releaseAdjacent != second.releaseAdjacent {
						authorityBasis = append(authorityBasis, "release_adjacent")
					}
					if best.openMeetingActions != second.openMeetingActions || best.activeMeetingCount != second.activeMeetingCount {
						authorityBasis = append(authorityBasis, "governance_activity")
					}
					if best.isRootCandidate != second.isRootCandidate {
						authorityBasis = append(authorityBasis, "root_candidate")
					}
					if len(authorityBasis) == 0 {
						authorityBasis = append(authorityBasis, "candidate_governor_tuple")
					}
				} else if hasAuthoritySignal {
					tiedMissionIDs := []string{}
					for _, candidate := range candidates {
						if candidate.priority == best.priority &&
							candidate.risk == best.risk &&
							candidate.releaseAdjacent == best.releaseAdjacent &&
							candidate.openMeetingActions == best.openMeetingActions &&
							candidate.activeMeetingCount == best.activeMeetingCount &&
							candidate.isRootCandidate == best.isRootCandidate {
							tiedMissionIDs = append(tiedMissionIDs, candidate.missionID)
						}
					}
					governanceNotes = append(governanceNotes, fmt.Sprintf("No unique full candidate-governor tuple exists across bundle candidates: %s.", strings.Join(dedupeStrings(tiedMissionIDs), ", ")))
				}
			}

			if selectedWinnerID != "" {
				state = "governed"
				governingMissionID = selectedWinnerID
				governedMissionIDs = append(governedMissionIDs, selectedWinnerID)
				for _, id := range component {
					if id != selectedWinnerID {
						losingMissionIDs = append(losingMissionIDs, id)
					}
				}
				if containsString(authorityBasis, "explicit_bundle_governor") {
					resolutionKind = "explicit_bundle_governor"
					requiredAction = fmt.Sprintf("Keep bundle authority anchored to %s until the explicit bundle-governor decision is cleared or superseded.", selectedWinnerID)
					governanceNotes = append(governanceNotes, fmt.Sprintf("Leader explicitly established %s as the bundle governor via signal %s.", selectedWinnerID, firstNonBlank(explicitGovernorSignalID, "bundle-governor signal")))
				}
				if containsString(authorityBasis, "council_adopted_bundle_governor") {
					resolutionKind = "council_adopted_bundle_governor"
					requiredAction = fmt.Sprintf("Keep bundle authority anchored to %s until a fresher bundle decision supersedes council adoption from %s.", selectedWinnerID, firstNonBlank(adoptedFromMeetingID, "the adopted council"))
					governanceNotes = append(governanceNotes, fmt.Sprintf("Closed bundle council %s adopted %s as the durable governor via signal %s.", firstNonBlank(adoptedFromMeetingID, "the adopted council"), selectedWinnerID, firstNonBlank(explicitGovernorSignalID, "bundle-governor signal")))
				} else {
					if resolutionKind == "linked_bundle" {
						resolutionKind = "bundle_governor_adjudicated"
						requiredAction = fmt.Sprintf("Keep bundle authority anchored to %s and serialize the other overlapping missions behind it unless a council overrides the current bundle governor.", selectedWinnerID)
						governanceNotes = append(governanceNotes, fmt.Sprintf("No explicit governing dependency exists; %s currently governs the bundle by %s.", selectedWinnerID, strings.Join(dedupeStrings(authorityBasis), "/")))
					}
				}
				governedBundleCount++
			} else if hasBundleAuthoritySignal {
				state = "escalation-required"
				if len(councilAnchorMissionIDs) > 1 {
					resolutionKind = "multi_council_bundle_governor_required"
					requiredAction = fmt.Sprintf("Consolidate overlapping active councils on %s or establish one governing mission before parallel execution.", strings.Join(councilAnchorMissionIDs, ", "))
					governanceNotes = append(governanceNotes, fmt.Sprintf("Multiple active council anchors currently govern the same overlap bundle without a single durable governor: %s.", strings.Join(councilAnchorMissionIDs, ", ")))
				} else {
					resolutionKind = "bundle_governor_required"
					requiredAction = "Open a council or add an explicit dependency/ownership transfer to establish a governing mission before parallel execution."
				}
				if len(bundleRoots) > 1 {
					governanceNotes = append(governanceNotes, "Multiple root missions still claim overlapping bundle authority without a decisive governing dependency.")
				}
				escalatedBundleCount++
			} else {
				if len(councilAnchorMissionIDs) > 1 {
					state = "escalation-required"
					resolutionKind = "multi_council_bundle_governor_required"
					requiredAction = fmt.Sprintf("Consolidate overlapping active councils on %s or establish one governing mission before parallel execution.", strings.Join(councilAnchorMissionIDs, ", "))
					governanceNotes = append(governanceNotes, fmt.Sprintf("Multiple active council anchors currently govern the same overlap bundle without a single durable governor: %s.", strings.Join(councilAnchorMissionIDs, ", ")))
					escalatedBundleCount++
				} else {
					state = "conflicted"
					resolutionKind = "bundle_conflict"
					requiredAction = "Narrow selectors, serialize missions, or add a governing dependency/council before parallel execution."
				}
				if len(councilAnchorMissionIDs) == 1 {
					requiredAction = fmt.Sprintf("Use the active council on %s to govern the remaining bundle conflict, or narrow selectors before widening execution.", councilAnchorMissionIDs[0])
					governanceNotes = append(governanceNotes, fmt.Sprintf("Bundle already has one active council anchor on %s, but unresolved overlap remains outside that governed lane.", councilAnchorMissionIDs[0]))
				}
				if len(bundleRoots) > 1 {
					governanceNotes = append(governanceNotes, "Multiple root missions still claim overlapping bundle authority.")
				}
				if state == "conflicted" {
					conflictedBundleCount++
				}
			}
		case adjudicationCount > 0:
			state = "governed"
			switch {
			case governingMissionID != "" && len(councilAnchorMissionIDs) == 1 && governingMissionID == councilAnchorMissionIDs[0]:
				governedMissionIDs = append(governedMissionIDs, governingMissionID)
				for _, id := range component {
					if id != governingMissionID {
						losingMissionIDs = append(losingMissionIDs, id)
					}
				}
				resolutionKind = "council_governed_bundle"
				requiredAction = fmt.Sprintf("Keep bundle authority anchored to %s until its active council conclusions are adopted.", governingMissionID)
			case governingMissionID != "" && len(bundleRoots) == 1 && governingMissionID == bundleRoots[0]:
				governedMissionIDs = append(governedMissionIDs, bundleRoots[0])
				for _, id := range component {
					if id != bundleRoots[0] {
						losingMissionIDs = append(losingMissionIDs, id)
					}
				}
				resolutionKind = "single_root_governed_bundle"
				requiredAction = fmt.Sprintf("Keep bundle authority anchored to %s while dependent missions stay bounded.", governingMissionID)
			case governingMissionID != "":
				governedMissionIDs = append(governedMissionIDs, governingMissionID)
				for _, id := range component {
					if id != governingMissionID {
						losingMissionIDs = append(losingMissionIDs, id)
					}
				}
				resolutionKind = "single_winner_governed_bundle"
				requiredAction = fmt.Sprintf("Keep bundle authority anchored to %s, and narrow or transfer broader overlapping missions before widening execution.", governingMissionID)
			default:
				governedMissionIDs = append(governedMissionIDs, bundleRoots...)
				resolutionKind = "multi_source_governed_bundle"
				requiredAction = "Bundle has adjudications but no single governing root; inspect mission-level authority before widening execution."
			}
			governedBundleCount++
		}
		if len(councilAnchorMissionIDs) > 0 {
			governanceNotes = append(governanceNotes, fmt.Sprintf("Active council anchor missions: %s", strings.Join(councilAnchorMissionIDs, ", ")))
		}
		if len(bundleSelectorRefs) > 0 {
			governanceNotes = append(governanceNotes, fmt.Sprintf("Shared governed selectors: %s", strings.Join(dedupeStrings(bundleSelectorRefs), ", ")))
		}
		bundles = append(bundles, map[string]any{
			"bundle_id":                   bundleID,
			"mission_ids":                 stringListToAny(component),
			"mission_count":               len(component),
			"edge_count":                  edgeCount,
			"edge_kinds":                  stringListToAny(dedupeStrings(edgeKinds)),
			"conflict_count":              conflictCount,
			"adjudication_count":          adjudicationCount,
			"bundle_state":                state,
			"resolution_kind":             resolutionKind,
			"resolution_kinds":            stringListToAny(dedupeStrings(resolutionKinds)),
			"selector_refs":               stringListToAny(dedupeStrings(bundleSelectorRefs)),
			"required_action":             requiredAction,
			"governance_notes":            stringListToAny(governanceNotes),
			"authority_basis":             stringListToAny(dedupeStrings(authorityBasis)),
			"candidate_governor_ids":      stringListToAny(candidateGovernorMissionIDs),
			"candidate_governor_packets":  cloneValue(candidateGovernorPackets),
			"explicit_governor_signal_id": explicitGovernorSignalID,
			"governance_source_kind":      governanceSourceKind,
			"adopted_from_meeting_id":     adoptedFromMeetingID,
			"canonical_mission_id":        component[0],
			"root_mission_ids":            stringListToAny(bundleRoots),
			"dependent_mission_ids":       stringListToAny(dependentMissionIDs),
			"governed_mission_ids":        stringListToAny(dedupeStrings(governedMissionIDs)),
			"losing_mission_ids":          stringListToAny(dedupeStrings(losingMissionIDs)),
			"governing_mission_id":        governingMissionID,
		})
	}

	return map[string]any{
		"mission_count":                     len(missions),
		"mission_ids":                       stringListToAny(missionIDs),
		"mission_gap_count":                 missionGapCount,
		"mission_adjudication_count":        missionAdjudicationCount,
		"mission_unresolved_conflict_count": missionUnresolvedCount,
		"cross_mission_conflict_count":      len(crossMissionConflicts),
		"cross_mission_adjudication_count":  len(crossMissionAdjudications),
		"dependency_adjudication_count":     len(dependencyAdjudications),
		"bundle_count":                      len(bundles),
		"conflicted_bundle_count":           conflictedBundleCount,
		"governed_bundle_count":             governedBundleCount,
		"escalated_bundle_count":            escalatedBundleCount,
		"mission_authority_audits":          missionAudits,
		"bundles":                           bundles,
		"cross_mission_adjudications":       crossMissionAdjudications,
		"cross_mission_conflicts":           crossMissionConflicts,
		"dependency_adjudications":          dependencyAdjudications,
	}
}

func siblingLinked(run map[string]any, targetRunID string) bool {
	for _, item := range append(toList(run["communication_peers"]), toList(run["siblings"])...) {
		if stringValue(item) == targetRunID {
			return true
		}
	}
	return false
}

func sourceTypeForRun(run map[string]any) string {
	ownerType := stringValue(run["owner_type"])
	switch ownerType {
	case "leader", "teammate", "automation", "script":
		return ownerType
	default:
		if ownerType != "" {
			return "teammate"
		}
		return "leader"
	}
}

func allowedRunTransitions() map[string]map[string]bool {
	return map[string]map[string]bool{
		"draft":     {"ready": true, "abandoned": true},
		"ready":     {"executing": true, "waiting": true, "blocked": true, "completed": true, "abandoned": true},
		"executing": {"ready": true, "waiting": true, "blocked": true, "completed": true, "abandoned": true},
		"waiting":   {"ready": true, "executing": true, "blocked": true, "abandoned": true},
		"blocked":   {"ready": true, "waiting": true, "abandoned": true},
		"completed": {},
		"abandoned": {},
	}
}

func validateRunTransition(run map[string]any, newState string) error {
	current := stringValue(run["state"])
	if current == "" {
		current = "draft"
	}
	if current == newState {
		return nil
	}
	allowed := allowedRunTransitions()[current]
	if !allowed[newState] {
		return fmt.Errorf("invalid run state transition: %s -> %s", current, newState)
	}
	return nil
}

func policyValue(key string, defaultValue any) any {
	policies, _ := toMap(loadData(policiesPath(), map[string]any{}))
	return getKey(policies["policies"], key, defaultValue)
}

func meetingMaxDebateRounds() int {
	rounds := intValue(policyValue("meetings.max_debate_rounds", 3))
	if rounds < 2 {
		return 2
	}
	return rounds
}

func recordVerification(run map[string]any, ref string) {
	refs := appendUniqueText(toList(run["verification_refs"]), ref)
	run["verification_refs"] = refs
	run["verification_status"] = "passed"
	run["updated_at"] = nowISO()
}

func recordReview(run map[string]any, state, note string) {
	refs := toList(run["review_refs"])
	if note != "" {
		refs = appendUniqueText(refs, note)
	}
	run["review_refs"] = refs
	run["review_status"] = state
	run["updated_at"] = nowISO()
}

func recordChecks(run map[string]any, state, ref string) {
	refs := toList(run["checks_refs"])
	if ref != "" {
		refs = appendUniqueText(refs, ref)
	}
	run["checks_refs"] = refs
	run["checks_status"] = state
	run["updated_at"] = nowISO()
}

func recordShipRef(run map[string]any, ref string) {
	run["ship_refs"] = appendUniqueText(toList(run["ship_refs"]), ref)
	run["updated_at"] = nowISO()
}

func mergeReady(run map[string]any) bool {
	return stringValue(run["verification_status"]) == "passed" &&
		stringValue(run["review_status"]) == "approved" &&
		stringValue(run["checks_status"]) == "passed" &&
		(strings.EqualFold(stringValue(run["ship_status"]), "in_review") || stringValue(run["ship_status"]) == "approved")
}

func markRunMerged(run map[string]any, summary string) error {
	run["ship_status"] = "merged"
	run["shipped_at"] = nowISO()
	recordShipRef(run, "merged")
	if stringValue(run["state"]) != "completed" {
		if err := validateRunTransition(run, "completed"); err != nil {
			return err
		}
		run["state"] = "completed"
		run["completed_at"] = nowISO()
	}
	run["latest_summary"] = summary
	run["next_expected_action"] = "Consider whether this shipped work should enter an RC."
	run["updated_at"] = nowISO()
	return nil
}

func issueBodyFromTemplate(title, objective string, successCriteria, notes []string) (string, error) {
	templateBytes, err := os.ReadFile(issueTemplatePath())
	if err != nil {
		return "", err
	}
	criteriaBlock := "- TBD"
	if len(successCriteria) > 0 {
		lines := make([]string, 0, len(successCriteria))
		for _, item := range successCriteria {
			lines = append(lines, "- "+item)
		}
		criteriaBlock = strings.Join(lines, "\n")
	}
	notesBlock := "- None"
	if len(notes) > 0 {
		lines := make([]string, 0, len(notes))
		for _, item := range notes {
			lines = append(lines, "- "+item)
		}
		notesBlock = strings.Join(lines, "\n")
	}
	body := string(templateBytes)
	body = strings.ReplaceAll(body, "__TITLE__", title)
	body = strings.ReplaceAll(body, "__OBJECTIVE__", objective)
	body = strings.ReplaceAll(body, "__SUCCESS_CRITERIA__", criteriaBlock)
	body = strings.ReplaceAll(body, "__NOTES__", notesBlock)
	return body, nil
}

func prPayloadForRun(run map[string]any, baseBranch, overrideTitle string) (map[string]any, error) {
	title := strings.TrimSpace(overrideTitle)
	if title == "" {
		title = fmt.Sprintf("[%s] %s", stringValue(run["mission_id"]), strings.TrimSpace(stringValue(run["objective"])))
		title = strings.TrimSpace(title)
	}
	summary := strings.TrimSpace(stringValue(run["latest_summary"]))
	if summary == "" {
		summary = strings.TrimSpace(stringValue(run["objective"]))
	}
	if summary == "" {
		summary = "Mission delivery"
	}
	verificationRefs := stringSlice(run["verification_refs"])
	verificationBlock := "- Not yet recorded"
	if len(verificationRefs) > 0 {
		lines := make([]string, 0, len(verificationRefs))
		for _, item := range verificationRefs {
			lines = append(lines, "- "+item)
		}
		verificationBlock = strings.Join(lines, "\n")
	}
	templateBytes, err := os.ReadFile(prTemplatePath())
	if err != nil {
		return nil, err
	}
	body := string(templateBytes)
	body = strings.ReplaceAll(body, "__SUMMARY__", summary)
	body = strings.ReplaceAll(body, "__VERIFICATION__", verificationBlock)
	body = strings.ReplaceAll(body, "__BRANCH__", stringValue(run["branch"]))
	body = strings.ReplaceAll(body, "__RUN_ID__", stringValue(run["run_id"]))
	body = strings.ReplaceAll(body, "__MISSION_ID__", stringValue(run["mission_id"]))
	return map[string]any{
		"title":       title,
		"body":        body,
		"base_branch": baseBranch,
		"branch":      stringValue(run["branch"]),
	}, nil
}

func openGitHubPR(run map[string]any, baseBranch, title, body string, draft bool, reviewers []string) (map[string]any, error) {
	return openGitHubPRForBranch(stringValue(run["branch"]), baseBranch, title, body, draft, reviewers)
}

func openGitHubPRForBranch(branch, baseBranch, title, body string, draft bool, reviewers []string) (map[string]any, error) {
	cmd := []string{"gh", "pr", "create", "--base", baseBranch, "--head", branch, "--title", title, "--body", body}
	if draft {
		cmd = append(cmd, "--draft")
	}
	for _, reviewer := range reviewers {
		cmd = append(cmd, "--reviewer", reviewer)
	}
	output, err := runCommandWithStderr(cmd, "")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	prURL := ""
	if len(lines) > 0 {
		prURL = strings.TrimSpace(lines[len(lines)-1])
	}
	prRef := ""
	if prURL != "" {
		parts := strings.Split(strings.TrimRight(prURL, "/"), "/")
		prRef = parts[len(parts)-1]
	}
	return map[string]any{"pr_ref": prRef, "pr_url": prURL}, nil
}

func currentBranch() string {
	output, err := runCommand([]string{"git", "branch", "--show-current"})
	if err != nil {
		return ""
	}
	return output
}

func branchExists(branch string) bool {
	_, err := runCommand([]string{"git", "show-ref", "--verify", "refs/heads/" + branch})
	return err == nil
}

func gitStatusPorcelain() string {
	output, err := runCommand([]string{"git", "status", "--porcelain"})
	if err != nil {
		return ""
	}
	return output
}

func gitIsDirty() bool {
	return strings.TrimSpace(gitStatusPorcelain()) != ""
}

func ensureCleanCheckoutContext() error {
	if gitIsDirty() {
		return errors.New("current worktree has uncommitted changes; refusing branch checkout")
	}
	return nil
}

func worktreeRoot(profile map[string]any) string {
	configured := stringValue(getKey(profile, "profile.worktree_root", ""))
	if configured != "" {
		if filepath.IsAbs(configured) {
			return configured
		}
		return filepath.Join(repoRoot(), configured)
	}
	parent := filepath.Dir(repoRoot())
	return filepath.Join(parent, filepath.Base(repoRoot())+"-worktrees")
}

func provisionWorktreePath(profile map[string]any, slug string) string {
	return filepath.Join(worktreeRoot(profile), slug)
}

func gitCheckoutNewBranch(branch, baseBranch string) error {
	if currentBranch() == branch {
		return nil
	}
	if err := ensureCleanCheckoutContext(); err != nil {
		return err
	}
	_, _ = runCommand([]string{"git", "fetch", "origin", baseBranch})
	if branchExists(branch) {
		_, err := runCommand([]string{"git", "checkout", branch})
		return err
	}
	if _, err := runCommand([]string{"git", "checkout", baseBranch}); err != nil {
		return err
	}
	_, _ = runCommand([]string{"git", "pull", "--ff-only", "origin", baseBranch})
	_, err := runCommand([]string{"git", "checkout", "-b", branch})
	return err
}

func gitCheckoutBranch(branch string) error {
	_, err := runCommand([]string{"git", "checkout", branch})
	return err
}

func remoteBranchExists(remote, branch string) bool {
	command := exec.Command("git", "ls-remote", "--exit-code", "--heads", remote, branch)
	command.Dir = repoRoot()
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	return command.Run() == nil
}

func gitPushBranch(remote, branch string) error {
	_, err := runCommandWithStderr([]string{"git", "push", "-u", remote, branch}, "")
	return err
}

func gitDeleteRemoteBranch(remote, branch string) error {
	if branch == "" || !remoteBranchExists(remote, branch) {
		return nil
	}
	_, err := runCommandWithStderr([]string{"git", "push", remote, "--delete", branch}, "")
	return err
}

func gitResolveCommitish(commitish string) (string, error) {
	return runCommand([]string{"git", "rev-parse", "--verify", commitish})
}

func gitCherryPick(commitish string) error {
	_, err := runCommandWithStderr([]string{"git", "cherry-pick", commitish}, "")
	return err
}

func applyCherryPickBranch(branch, baseBranch, resolvedCommit, remote string, applyChanges, pushChanges bool) (map[string]any, error) {
	previousBranch := currentBranch()
	applied := false
	pushed := false
	if applyChanges {
		if err := gitCheckoutNewBranch(branch, baseBranch); err != nil {
			return nil, err
		}
		defer func() {
			if previousBranch != "" && previousBranch != branch && branchExists(previousBranch) {
				_ = gitCheckoutBranch(previousBranch)
			}
		}()
		if err := gitCherryPick(resolvedCommit); err != nil {
			return nil, err
		}
		applied = true
		if pushChanges {
			if err := gitPushBranch(remote, branch); err != nil {
				return nil, err
			}
			pushed = true
		}
	}
	return map[string]any{
		"previous_branch": previousBranch,
		"applied":         applied,
		"pushed":          pushed,
	}, nil
}

func gitCreateWorktree(branch, path, baseBranch string) error {
	if exists(path) {
		return fmt.Errorf("worktree path already exists: %s", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	_, _ = runCommand([]string{"git", "fetch", "origin", baseBranch})
	if branchExists(branch) {
		_, err := runCommand([]string{"git", "worktree", "add", path, branch})
		return err
	}
	if _, err := runCommand([]string{"git", "worktree", "add", "-b", branch, path, "origin/" + baseBranch}); err == nil {
		return nil
	}
	_, err := runCommand([]string{"git", "worktree", "add", "-b", branch, path, baseBranch})
	return err
}

func gitMergeBranchIntoBase(branch, baseBranch, strategy string) error {
	if strategy != "merge" && strategy != "squash" {
		return errors.New("git merge strategy must be one of: merge, squash")
	}
	if err := ensureCleanCheckoutContext(); err != nil {
		return err
	}
	current := currentBranch()
	if _, err := runCommand([]string{"git", "checkout", baseBranch}); err != nil {
		return err
	}
	if strategy == "squash" {
		if _, err := runCommandWithStderr([]string{"git", "merge", "--squash", branch}, ""); err != nil {
			_ = restoreBranch(current, baseBranch)
			return err
		}
		if _, err := runCommandWithStderr([]string{"git", "commit", "-m", fmt.Sprintf("Squash merge %s into %s", branch, baseBranch)}, ""); err != nil {
			_ = restoreBranch(current, baseBranch)
			return err
		}
		return nil
	}
	if _, err := runCommandWithStderr([]string{"git", "merge", "--no-ff", branch, "-m", fmt.Sprintf("Merge %s into %s", branch, baseBranch)}, ""); err != nil {
		_ = restoreBranch(current, baseBranch)
		return err
	}
	return nil
}

func restoreBranch(current, baseBranch string) error {
	if current != "" && current != baseBranch && branchExists(current) {
		return gitCheckoutBranch(current)
	}
	return nil
}

func gitDeleteBranch(branch, fallbackBranch string) error {
	if branch == "" {
		return nil
	}
	current := currentBranch()
	if current == branch {
		if err := ensureCleanCheckoutContext(); err != nil {
			return err
		}
		if _, err := runCommand([]string{"git", "checkout", fallbackBranch}); err != nil {
			return err
		}
	}
	if !branchExists(branch) {
		return nil
	}
	_, err := runCommand([]string{"git", "branch", "-D", branch})
	return err
}

func gitRemoveWorktree(path string) error {
	if !exists(path) {
		return nil
	}
	mainRepo, _ := runCommand([]string{"git", "worktree", "list", "--porcelain"})
	if path != "" && strings.Contains(mainRepo, path) {
		_, err := runCommand([]string{"git", "worktree", "remove", "--force", path})
		return err
	}
	return fmt.Errorf("path exists but is not a registered worktree: %s", path)
}

func ghMergePR(prRef, strategy string, deleteBranch bool) error {
	strategyFlag := map[string]string{"merge": "--merge", "squash": "--squash", "rebase": "--rebase"}[strategy]
	if strategyFlag == "" {
		return errors.New("gh merge strategy must be one of: merge, squash, rebase")
	}
	cmd := []string{"gh", "pr", "merge", prRef, strategyFlag}
	if deleteBranch {
		cmd = append(cmd, "--delete-branch")
	} else {
		cmd = append(cmd, "--delete-branch=false")
	}
	_, err := runCommandWithStderr(cmd, "")
	return err
}

func defaultAuthorityHierarchy() []any {
	return []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"}
}

func normalizeDecisionProtocol(raw string) string {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "", "evidence-review", "evidence_review", "evidence":
		return "evidence-review"
	case "incident-command", "incident_command", "incident":
		return "incident-command"
	case "decision-record", "decision_record", "decision", "adr":
		return "decision-record"
	default:
		return ""
	}
}

func defaultMeetingRoles(protocol string) []any {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return []any{"incident-commander", "operations-lead", "risk-assessor"}
	case "decision-record":
		return []any{"context-owner", "alternative-challenger", "consequence-checker"}
	default:
		return []any{"fact-verifier", "architecture-challenger", "completeness-checker"}
	}
}

func defaultMeetingExpectedOutput(protocol string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return "Produce current situation, chosen immediate action, owner, stop/go constraints, and the next checkpoint."
	case "decision-record":
		return "Produce the decision, alternatives considered, consequences, unresolved questions, and follow-up actions."
	default:
		return "Produce accepted findings, rejected findings, unresolved questions, and action items."
	}
}

func defaultMeetingRoleSpecs(protocol string) []any {
	timestamp := nowISO()
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return []any{
			map[string]any{
				"role":              "incident-commander",
				"mission":           "Keep the decision bounded to the current incident, choose the immediate control action, and prevent unmanaged blast-radius expansion.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Speculation about future design work used to delay immediate stabilizing action"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
			map[string]any{
				"role":              "operations-lead",
				"mission":           "Describe the current operational state, candidate actions, dependencies, and concrete execution constraints for the immediate response window.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Stale run summaries used to override fresher execution evidence"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
			map[string]any{
				"role":              "risk-assessor",
				"mission":           "Surface rollback risk, user impact, and stop/go constraints before the meeting commits to an immediate action.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Unbounded architecture debates that do not change the immediate decision horizon"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
		}
	case "decision-record":
		return []any{
			map[string]any{
				"role":              "context-owner",
				"mission":           "State the current context, constraints, and why the decision must be made now without overstating certainty.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Retrofitted rationale that does not match current evidence"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
			map[string]any{
				"role":              "alternative-challenger",
				"mission":           "Challenge default framing, surface viable alternatives, and make tradeoffs explicit before the decision is recorded.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Purely rhetorical challenge without concrete alternative or consequence"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
			map[string]any{
				"role":              "consequence-checker",
				"mission":           "Verify consequences, rollout implications, and follow-up obligations so the recorded decision is actually actionable.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Already rejected alternatives reused without updated evidence"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
		}
	default:
		return []any{
			map[string]any{
				"role":              "fact-verifier",
				"mission":           "Verify factual claims against the highest-authority live state before the council accepts them.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"L4+: stale notes, legacy artifacts, or summaries used to override fresher evidence"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
			map[string]any{
				"role":              "architecture-challenger",
				"mission":           "Challenge hidden assumptions, blast radius, and design coupling without overriding higher-authority facts.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Lower-authority notes used as truth", "Unbounded speculation without evidence"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
			map[string]any{
				"role":              "completeness-checker",
				"mission":           "Find missing decision inputs, missing rollout details, and missing operator-facing guidance before synthesis.",
				"allowed_sources":   []any{"L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"},
				"forbidden_sources": []any{"Already rejected findings reused without re-verification"},
				"created_at":        timestamp,
				"updated_at":        timestamp,
			},
		}
	}
}

func protocolUsesStructuredDebate(protocol string) bool {
	return normalizeDecisionProtocol(protocol) == "evidence-review"
}

func protocolDirectSynthesisCoverageNote(protocol, claimRef string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return fmt.Sprintf("claim %q stays out of live cross-examination because incident-command prefers bounded stabilization and explicit owner/stop-go output over extended rebuttal rounds.", claimRef)
	case "decision-record":
		return fmt.Sprintf("claim %q stays out of live cross-examination because decision-record prefers bounded capture of tradeoffs and follow-up obligations over extended rebuttal rounds.", claimRef)
	default:
		return fmt.Sprintf("claim %q stays out of live cross-examination because this protocol defers unresolved conflict to synthesis.", claimRef)
	}
}

func protocolBoundedRoundCoverageNote(protocol, claimRef string, roundNumber int) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return fmt.Sprintf("claim %q stops after round %d because incident-command keeps the council bounded and hands unresolved disagreement to synthesis/action ownership instead of further debate rounds.", claimRef, roundNumber)
	case "decision-record":
		return fmt.Sprintf("claim %q stops after round %d because decision-record records the bounded disagreement and consequences instead of escalating into more rebuttal rounds.", claimRef, roundNumber)
	default:
		return fmt.Sprintf("claim %q stops after round %d because the active protocol caps live debate and defers the unresolved lane to synthesis.", claimRef, roundNumber)
	}
}

func protocolDeferredSynthesisRequirement(protocol string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return "No additional live debate is required; synthesize the incident command record with chosen action, owner, and stop/go constraints."
	case "decision-record":
		return "No additional live debate is required; synthesize the bounded decision record with alternatives, consequences, and follow-up actions."
	default:
		return "No additional live debate is required; synthesize the bounded council record."
	}
}

func protocolIndependentReviewRequirement(protocol string, openWaveActionCount int) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		if openWaveActionCount > 0 {
			return fmt.Sprintf("Collect %d bounded incident assessment(s) covering current state, immediate action, owner, and stop/go constraints.", openWaveActionCount)
		}
		return "Collect bounded incident assessments covering current state, immediate action, owner, and stop/go constraints."
	case "decision-record":
		if openWaveActionCount > 0 {
			return fmt.Sprintf("Collect %d bounded decision-record assessment(s) covering context, alternatives, consequences, and follow-up obligations.", openWaveActionCount)
		}
		return "Collect bounded decision-record assessments covering context, alternatives, consequences, and follow-up obligations."
	default:
		if openWaveActionCount > 0 {
			return fmt.Sprintf("Collect %d independent-review position(s) from bounded participants.", openWaveActionCount)
		}
		return "Collect independent-review positions from bounded participants."
	}
}

func protocolOutputRequirement(protocol string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return "Adopt or close the meeting outcome after recording chosen action, owner, stop/go constraints, and the next checkpoint."
	case "decision-record":
		return "Adopt or close the meeting outcome after recording the decision, alternatives, consequences, unresolved questions, and follow-up actions."
	default:
		return "Adopt or close the meeting outcome into run state."
	}
}

func protocolAutoSynthesisReviewRequirement(protocol string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return "Review the auto-synthesized incident command output or orchestrate again to confirm owner, stop/go constraints, and next checkpoint."
	case "decision-record":
		return "Review the auto-synthesized decision record or orchestrate again to confirm alternatives, consequences, and follow-up actions."
	default:
		return "Review auto-synthesized council output or orchestrate again to advance toward output."
	}
}

func protocolDeferredActionSummary(protocol, claimRef string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return fmt.Sprintf("Incident commander judgment and bounded action record required for claim: %s", claimRef)
	case "decision-record":
		return fmt.Sprintf("Decision record must capture unresolved tradeoff and follow-up for claim: %s", claimRef)
	default:
		return fmt.Sprintf("Leader judgment or explicit synthesis note required for claim: %s", claimRef)
	}
}

func protocolOutputPacketKind(protocol string) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return "incident_command_output"
	case "decision-record":
		return "decision_record_output"
	default:
		return "evidence_review_output"
	}
}

func summaryStrings(items []any) []string {
	result := []string{}
	for _, item := range items {
		text := strings.TrimSpace(textValue(item))
		if text != "" {
			result = append(result, text)
		}
	}
	return dedupeStrings(result)
}

func meetingOpenActionSummaries(meeting map[string]any) []string {
	summaries := []string{}
	for _, item := range toList(meeting["action_items"]) {
		action, ok := toMap(item)
		if ok {
			state := strings.TrimSpace(stringValue(action["state"]))
			if state != "" && state != "open" {
				continue
			}
		}
		if text := strings.TrimSpace(textValue(item)); text != "" {
			summaries = append(summaries, text)
		}
	}
	return dedupeStrings(summaries)
}

func meetingPositionSummariesByStance(meeting map[string]any, stance string) []string {
	summaries := []string{}
	for _, item := range toList(meeting["positions"]) {
		position, ok := toMap(item)
		if !ok || strings.TrimSpace(stringValue(position["stance"])) != stance {
			continue
		}
		if text := strings.TrimSpace(stringValue(position["summary"])); text != "" {
			summaries = append(summaries, text)
		}
	}
	return dedupeStrings(summaries)
}

func meetingLockedDecisionSummaries(meeting map[string]any) []string {
	summaries := []string{}
	for _, item := range toList(meeting["locked_decisions"]) {
		if text := strings.TrimSpace(textValue(item)); text != "" {
			summaries = append(summaries, text)
		}
	}
	return dedupeStrings(summaries)
}

func meetingParticipantRoles(meeting map[string]any) []string {
	roles := []string{}
	for _, item := range toList(meeting["participants"]) {
		participant, ok := toMap(item)
		if !ok {
			continue
		}
		role := firstNonBlank(stringValue(participant["role"]), stringValue(participant["owner_id"]), stringValue(participant["run_id"]))
		if strings.TrimSpace(role) != "" {
			roles = append(roles, role)
		}
	}
	return dedupeStrings(roles)
}

func meetingDebateLaneDisposition(route map[string]any) string {
	switch strings.TrimSpace(stringValue(route["resolution_kind"])) {
	case "higher_authority_prevails", "latest_response_prevails":
		switch strings.TrimSpace(stringValue(route["prevailing_stance"])) {
		case "support":
			return "accepted"
		case "challenge":
			return "rejected"
		}
	case "protocol_prefers_direct_synthesis", "protocol_bound_reached", "max_rounds_reached":
		return "unresolved"
	}
	status := strings.TrimSpace(stringValue(route["status"]))
	if status == "" {
		return "unresolved"
	}
	if strings.HasPrefix(status, "awaiting_") {
		return "open"
	}
	return "unresolved"
}

func meetingDebateLaneLiveDebateOpen(route map[string]any) bool {
	status := strings.TrimSpace(stringValue(route["status"]))
	return status == "awaiting_cross_examination" || strings.HasPrefix(status, "awaiting_round_")
}

func meetingDebateLaneSummary(route map[string]any) string {
	claimRef := firstNonBlank(strings.TrimSpace(stringValue(route["claim_ref"])), strings.TrimSpace(stringValue(route["claim_key"])), "unspecified claim")
	resolutionKind := strings.TrimSpace(stringValue(route["resolution_kind"]))
	prevailingRole := firstNonBlank(strings.TrimSpace(stringValue(route["prevailing_role"])), strings.TrimSpace(stringValue(route["prevailing_stance"])), "the prevailing lane")
	overruledRole := firstNonBlank(strings.TrimSpace(stringValue(route["overruled_role"])), strings.TrimSpace(stringValue(route["overruled_stance"])), "the opposing lane")
	switch resolutionKind {
	case "higher_authority_prevails":
		return fmt.Sprintf("claim %q resolved for %s because higher-authority evidence outranked %s without requiring live cross-examination.", claimRef, prevailingRole, overruledRole)
	case "latest_response_prevails":
		return fmt.Sprintf("claim %q resolved for %s after the latest governed response prevailed over %s.", claimRef, prevailingRole, overruledRole)
	case "protocol_prefers_direct_synthesis":
		return fmt.Sprintf("claim %q stays unresolved in live debate because the active protocol defers this lane into bounded synthesis.", claimRef)
	case "protocol_bound_reached":
		return fmt.Sprintf("claim %q stays unresolved after bounded live debate and now requires synthesis to record the remaining tradeoff.", claimRef)
	case "max_rounds_reached":
		return fmt.Sprintf("claim %q reached the live debate round cap and now requires synthesis or leader judgment to close the lane.", claimRef)
	}
	status := strings.TrimSpace(stringValue(route["status"]))
	if status == "awaiting_cross_examination" {
		return fmt.Sprintf("claim %q still awaits governed cross-examination before synthesis can close the lane.", claimRef)
	}
	if strings.HasPrefix(status, "awaiting_round_") {
		return fmt.Sprintf("claim %q still awaits a governed follow-up response before synthesis can close the lane.", claimRef)
	}
	return fmt.Sprintf("claim %q remains an active governed debate lane.", claimRef)
}

func meetingDebateLaneOutput(meeting map[string]any) map[string]any {
	debatePlan, _ := toMap(meetingDebatePlan(meeting))
	findings, _ := toMap(meeting["findings"])
	type claimFinding struct {
		bucket  string
		finding map[string]any
	}
	findingsByClaimKey := map[string]claimFinding{}
	for _, bucket := range []string{"accepted", "rejected", "unresolved"} {
		for _, item := range toList(findings[bucket]) {
			finding, ok := toMap(item)
			if !ok {
				continue
			}
			claimKey := firstNonBlank(
				strings.TrimSpace(stringValue(finding["generated_claim_key"])),
				meetingClaimGroupKey(stringValue(finding["claim_ref"]), stringValue(finding["dimension"]), firstNonBlank(stringValue(finding["summary"]), stringValue(finding["finding_id"]))),
			)
			if claimKey == "" {
				continue
			}
			if _, exists := findingsByClaimKey[claimKey]; !exists {
				findingsByClaimKey[claimKey] = claimFinding{bucket: bucket, finding: finding}
			}
		}
	}
	actionSummaryByClaimKey := map[string]string{}
	for _, item := range toList(meeting["action_items"]) {
		action, ok := toMap(item)
		if !ok {
			continue
		}
		claimKey := strings.TrimSpace(stringValue(action["generated_claim_key"]))
		summary := strings.TrimSpace(stringValue(action["summary"]))
		if claimKey != "" && summary != "" && actionSummaryByClaimKey[claimKey] == "" {
			actionSummaryByClaimKey[claimKey] = summary
		}
	}
	packets := []any{}
	acceptedCount := 0
	rejectedCount := 0
	unresolvedCount := 0
	openCount := 0
	seenClaimKeys := map[string]bool{}
	for _, item := range toList(debatePlan["routes"]) {
		route, ok := toMap(item)
		if !ok {
			continue
		}
		claimKey := firstNonBlank(
			strings.TrimSpace(stringValue(route["claim_key"])),
			meetingClaimGroupKey(stringValue(route["claim_ref"]), stringValue(route["dimension"]), stringValue(route["route_id"])),
		)
		seenClaimKeys[claimKey] = true
		laneStatus := stringValue(route["status"])
		resolutionKind := stringValue(route["resolution_kind"])
		requiredFollowUp := stringValue(route["required_follow_up"])
		summary := meetingDebateLaneSummary(route)
		disposition := meetingDebateLaneDisposition(route)
		if claimFinding, ok := findingsByClaimKey[claimKey]; ok {
			disposition = claimFinding.bucket
			if reason := strings.TrimSpace(stringValue(claimFinding.finding["reason"])); reason != "" {
				summary = reason
			} else if findingSummary := strings.TrimSpace(stringValue(claimFinding.finding["summary"])); findingSummary != "" {
				summary = findingSummary
			}
			if requiredFollowUp == "" {
				requiredFollowUp = actionSummaryByClaimKey[claimKey]
			}
			if resolutionKind == "" {
				switch claimFinding.bucket {
				case "accepted", "rejected":
					resolutionKind = "synthesized_" + claimFinding.bucket + "_finding"
					if laneStatus == "" || laneStatus == "routed" {
						laneStatus = "finding_recorded"
					}
				case "unresolved":
					if routeDebateRound(route) >= intValue(route["max_debate_rounds"]) {
						resolutionKind = "max_rounds_reached"
						if laneStatus == "" || laneStatus == "routed" {
							laneStatus = "max_rounds_reached"
						}
					} else {
						resolutionKind = "synthesized_unresolved_finding"
						if laneStatus == "" || laneStatus == "routed" {
							laneStatus = "finding_recorded"
						}
					}
				}
			}
		}
		liveDebateOpen := laneStatus == "awaiting_cross_examination" || strings.HasPrefix(laneStatus, "awaiting_round_")
		if liveDebateOpen {
			openCount++
		}
		switch disposition {
		case "accepted":
			acceptedCount++
		case "rejected":
			rejectedCount++
		case "open", "unresolved":
			unresolvedCount++
		}
		packets = append(packets, map[string]any{
			"route_id":                       stringValue(route["route_id"]),
			"claim_ref":                      stringValue(route["claim_ref"]),
			"dimension":                      firstNonBlank(stringValue(route["dimension"]), "factual"),
			"lane_status":                    laneStatus,
			"resolution_kind":                resolutionKind,
			"output_disposition":             disposition,
			"live_debate_open":               liveDebateOpen,
			"required_follow_up":             requiredFollowUp,
			"linked_cross_examination_count": intValue(route["linked_cross_examination_count"]),
			"current_round":                  routeDebateRound(route),
			"max_debate_rounds":              intValue(route["max_debate_rounds"]),
			"assigned_role":                  stringValue(route["assigned_role"]),
			"assigned_run_id":                stringValue(route["assigned_run_id"]),
			"support_roles":                  cloneValue(getKey(route, "support_roles", []any{})),
			"challenge_roles":                cloneValue(getKey(route, "challenge_roles", []any{})),
			"prevailing_stance":              stringValue(route["prevailing_stance"]),
			"prevailing_role":                stringValue(route["prevailing_role"]),
			"prevailing_run_id":              stringValue(route["prevailing_run_id"]),
			"overruled_stance":               stringValue(route["overruled_stance"]),
			"overruled_role":                 stringValue(route["overruled_role"]),
			"overruled_run_id":               stringValue(route["overruled_run_id"]),
			"summary":                        summary,
		})
	}
	for claimKey, claimFinding := range findingsByClaimKey {
		if seenClaimKeys[claimKey] {
			continue
		}
		disposition := claimFinding.bucket
		summary := firstNonBlank(
			strings.TrimSpace(stringValue(claimFinding.finding["reason"])),
			strings.TrimSpace(stringValue(claimFinding.finding["summary"])),
			strings.TrimSpace(stringValue(claimFinding.finding["claim_ref"])),
			claimKey,
		)
		switch disposition {
		case "accepted":
			acceptedCount++
		case "rejected":
			rejectedCount++
		case "open", "unresolved":
			unresolvedCount++
		}
		packets = append(packets, map[string]any{
			"route_id":                       "",
			"claim_ref":                      stringValue(claimFinding.finding["claim_ref"]),
			"dimension":                      firstNonBlank(stringValue(claimFinding.finding["dimension"]), "factual"),
			"lane_status":                    "finding_recorded",
			"resolution_kind":                "synthesized_" + disposition + "_finding",
			"output_disposition":             disposition,
			"live_debate_open":               false,
			"required_follow_up":             actionSummaryByClaimKey[claimKey],
			"linked_cross_examination_count": 0,
			"current_round":                  0,
			"max_debate_rounds":              meetingMaxDebateRounds(),
			"assigned_role":                  stringValue(claimFinding.finding["source_role"]),
			"assigned_run_id":                stringValue(claimFinding.finding["source_run_id"]),
			"support_roles":                  []any{},
			"challenge_roles":                []any{},
			"prevailing_stance":              "",
			"prevailing_role":                "",
			"prevailing_run_id":              "",
			"overruled_stance":               "",
			"overruled_role":                 "",
			"overruled_run_id":               "",
			"summary":                        summary,
		})
	}
	return map[string]any{
		"debate_lane_count":            len(packets),
		"accepted_debate_lane_count":   acceptedCount,
		"rejected_debate_lane_count":   rejectedCount,
		"unresolved_debate_lane_count": unresolvedCount,
		"open_debate_lane_count":       openCount,
		"debate_lane_packets":          packets,
	}
}

func meetingOutputFollowUpExecutionKind(protocol string, lane map[string]any, action map[string]any) string {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return "incident_checkpoint"
	case "decision-record":
		return "decision_record_follow_up"
	default:
		if lane != nil {
			disposition := strings.TrimSpace(stringValue(lane["output_disposition"]))
			if disposition == "unresolved" || disposition == "open" {
				return "resolve_governed_lane"
			}
		}
		if action != nil && strings.TrimSpace(stringValue(action["summary"])) != "" {
			return "complete_output_follow_up"
		}
		return "adopt_or_close_output"
	}
}

func meetingOutputFollowUpPackets(meeting map[string]any, protocol string, debateLaneOutput map[string]any) map[string]any {
	laneByClaimKey := map[string]map[string]any{}
	seenLaneKeys := []string{}
	for _, item := range toList(debateLaneOutput["debate_lane_packets"]) {
		lane, ok := toMap(item)
		if !ok {
			continue
		}
		claimRef := strings.TrimSpace(stringValue(lane["claim_ref"]))
		dimension := firstNonBlank(stringValue(lane["dimension"]), "factual")
		if claimRef == "" {
			continue
		}
		claimKey := meetingClaimGroupKey(claimRef, dimension, claimRef)
		laneByClaimKey[claimKey] = lane
		seenLaneKeys = append(seenLaneKeys, claimKey)
	}

	packets := []any{}
	openCount := 0
	coveredClaimKeys := map[string]bool{}

	for _, item := range toList(meeting["action_items"]) {
		action, ok := toMap(item)
		if !ok {
			continue
		}
		state := firstNonBlank(strings.TrimSpace(stringValue(action["state"])), "open")
		claimKey := strings.TrimSpace(stringValue(action["generated_claim_key"]))
		lane, _ := toMap(laneByClaimKey[claimKey])
		summary := firstNonBlank(
			textValue(action),
			stringValue(lane["required_follow_up"]),
			stringValue(lane["summary"]),
		)
		packet := map[string]any{
			"follow_up_packet_id": firstNonBlank(
				strings.TrimSpace(stringValue(action["action_id"])),
				"follow-up-"+stableSlugID(strings.Join([]string{claimKey, summary}, "::")),
			),
			"source_kind":        "action_item",
			"execution_kind":     meetingOutputFollowUpExecutionKind(protocol, lane, action),
			"decision_protocol":  normalizeDecisionProtocol(firstNonBlank(protocol, "evidence-review")),
			"state":              state,
			"summary":            summary,
			"action_id":          strings.TrimSpace(stringValue(action["action_id"])),
			"claim_key":          claimKey,
			"claim_ref":          strings.TrimSpace(stringValue(lane["claim_ref"])),
			"dimension":          firstNonBlank(stringValue(lane["dimension"]), "factual"),
			"lane_status":        strings.TrimSpace(stringValue(lane["lane_status"])),
			"output_disposition": strings.TrimSpace(stringValue(lane["output_disposition"])),
			"required_follow_up": firstNonBlank(strings.TrimSpace(stringValue(lane["required_follow_up"])), summary),
			"assigned_role":      strings.TrimSpace(stringValue(lane["assigned_role"])),
			"assigned_run_id":    strings.TrimSpace(stringValue(lane["assigned_run_id"])),
		}
		if state == "open" {
			openCount++
		}
		if claimKey != "" {
			coveredClaimKeys[claimKey] = true
		}
		packets = append(packets, packet)
	}

	for _, claimKey := range dedupeStrings(seenLaneKeys) {
		if claimKey == "" || coveredClaimKeys[claimKey] {
			continue
		}
		lane, _ := toMap(laneByClaimKey[claimKey])
		disposition := strings.TrimSpace(stringValue(lane["output_disposition"]))
		if disposition != "unresolved" && disposition != "open" {
			continue
		}
		summary := firstNonBlank(
			strings.TrimSpace(stringValue(lane["required_follow_up"])),
			strings.TrimSpace(stringValue(lane["summary"])),
			strings.TrimSpace(stringValue(lane["claim_ref"])),
		)
		packet := map[string]any{
			"follow_up_packet_id": "follow-up-" + stableSlugID(strings.Join([]string{claimKey, summary}, "::")),
			"source_kind":         "debate_lane",
			"execution_kind":      meetingOutputFollowUpExecutionKind(protocol, lane, nil),
			"decision_protocol":   normalizeDecisionProtocol(firstNonBlank(protocol, "evidence-review")),
			"state":               "open",
			"summary":             summary,
			"action_id":           "",
			"claim_key":           claimKey,
			"claim_ref":           strings.TrimSpace(stringValue(lane["claim_ref"])),
			"dimension":           firstNonBlank(stringValue(lane["dimension"]), "factual"),
			"lane_status":         strings.TrimSpace(stringValue(lane["lane_status"])),
			"output_disposition":  disposition,
			"required_follow_up":  firstNonBlank(strings.TrimSpace(stringValue(lane["required_follow_up"])), summary),
			"assigned_role":       strings.TrimSpace(stringValue(lane["assigned_role"])),
			"assigned_run_id":     strings.TrimSpace(stringValue(lane["assigned_run_id"])),
		}
		openCount++
		packets = append(packets, packet)
	}

	return map[string]any{
		"follow_up_packet_count":      len(packets),
		"open_follow_up_packet_count": openCount,
		"follow_up_packets":           packets,
	}
}

func meetingAutoOutputPacket(meeting map[string]any, protocol string) map[string]any {
	findings, _ := toMap(meeting["findings"])
	accepted := summaryStrings(toList(findings["accepted"]))
	rejected := summaryStrings(toList(findings["rejected"]))
	unresolved := summaryStrings(toList(findings["unresolved"]))
	openActions := meetingOpenActionSummaries(meeting)
	supportSummaries := meetingPositionSummariesByStance(meeting, "support")
	challengeSummaries := meetingPositionSummariesByStance(meeting, "challenge")
	lockedDecisions := meetingLockedDecisionSummaries(meeting)
	participantRoles := meetingParticipantRoles(meeting)
	debateLaneOutput, _ := toMap(meetingDebateLaneOutput(meeting))
	followUpPacketOutput, _ := toMap(meetingOutputFollowUpPackets(meeting, protocol, debateLaneOutput))
	summary := firstNonBlank(
		firstNonBlank(accepted...),
		firstNonBlank(supportSummaries...),
		firstNonBlank(lockedDecisions...),
		strings.TrimSpace(stringValue(meeting["question"])),
	)
	nextAction := firstNonBlank(firstNonBlank(openActions...), protocolOutputRequirement(protocol))
	completionGaps := []string{}
	packet := map[string]any{
		"output_packet_kind":           protocolOutputPacketKind(protocol),
		"decision_protocol":            normalizeDecisionProtocol(firstNonBlank(protocol, "evidence-review")),
		"generated_by":                 "orchestrate_auto_synthesis",
		"status":                       "provisional",
		"generated_at":                 nowISO(),
		"summary":                      summary,
		"next_action":                  nextAction,
		"accepted_findings":            stringListToAny(accepted),
		"rejected_findings":            stringListToAny(rejected),
		"unresolved_findings":          stringListToAny(unresolved),
		"follow_up_actions":            stringListToAny(openActions),
		"locked_decisions":             stringListToAny(lockedDecisions),
		"participant_roles":            stringListToAny(participantRoles),
		"debate_lane_count":            intValue(debateLaneOutput["debate_lane_count"]),
		"accepted_debate_lane_count":   intValue(debateLaneOutput["accepted_debate_lane_count"]),
		"rejected_debate_lane_count":   intValue(debateLaneOutput["rejected_debate_lane_count"]),
		"unresolved_debate_lane_count": intValue(debateLaneOutput["unresolved_debate_lane_count"]),
		"open_debate_lane_count":       intValue(debateLaneOutput["open_debate_lane_count"]),
		"debate_lane_packets":          cloneValue(getKey(debateLaneOutput, "debate_lane_packets", []any{})),
		"follow_up_packet_count":       intValue(followUpPacketOutput["follow_up_packet_count"]),
		"open_follow_up_packet_count":  intValue(followUpPacketOutput["open_follow_up_packet_count"]),
		"follow_up_packets":            cloneValue(getKey(followUpPacketOutput, "follow_up_packets", []any{})),
	}
	if intValue(debateLaneOutput["open_debate_lane_count"]) > 0 {
		completionGaps = append(completionGaps, "output packet still has live debate lanes open")
	}
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		chosenAction := firstNonBlank(firstNonBlank(supportSummaries...), firstNonBlank(openActions...), summary)
		owner := firstNonBlank(firstNonBlank(participantRoles...), "leader")
		stopGo := firstNonBlank(firstNonBlank(challengeSummaries...), firstNonBlank(lockedDecisions...), "Record explicit stop/go constraints before adopting the incident outcome.")
		nextCheckpoint := firstNonBlank(firstNonBlank(openActions...), "Close or adopt the incident outcome after recording owner and stop/go constraints.")
		if len(supportSummaries) == 0 {
			completionGaps = append(completionGaps, "incident output lacks a synthesized chosen action from support evidence")
		}
		if len(challengeSummaries) == 0 {
			completionGaps = append(completionGaps, "incident output lacks a synthesized stop/go constraint from challenging evidence")
		}
		packet["summary"] = chosenAction
		packet["chosen_action"] = chosenAction
		packet["owner"] = owner
		packet["stop_go_constraints"] = stopGo
		packet["next_checkpoint"] = nextCheckpoint
		packet["unresolved_questions"] = stringListToAny(unresolved)
	case "decision-record":
		decision := firstNonBlank(firstNonBlank(accepted...), firstNonBlank(supportSummaries...), summary)
		alternatives := dedupeStrings(append(append([]string{}, challengeSummaries...), rejected...))
		consequences := dedupeStrings(append(append([]string{}, openActions...), unresolved...))
		if len(alternatives) == 0 {
			completionGaps = append(completionGaps, "decision record output lacks explicit alternative capture")
		}
		if len(consequences) == 0 {
			completionGaps = append(completionGaps, "decision record output lacks explicit consequence or follow-up capture")
		}
		packet["summary"] = decision
		packet["decision"] = decision
		packet["alternatives"] = stringListToAny(alternatives)
		packet["consequences"] = stringListToAny(consequences)
		packet["unresolved_questions"] = stringListToAny(unresolved)
	default:
		verdict := firstNonBlank(firstNonBlank(accepted...), firstNonBlank(rejected...), summary)
		if len(accepted)+len(rejected)+len(unresolved) == 0 {
			completionGaps = append(completionGaps, "evidence-review output lacks synthesized findings")
		}
		packet["summary"] = verdict
		packet["verdict"] = verdict
	}
	packet["completion_gaps"] = stringListToAny(completionGaps)
	packet["completion_gap_count"] = len(completionGaps)
	return packet
}

func routeDebateRound(route map[string]any) int {
	round := intValue(route["latest_debate_round"])
	if round > 0 {
		return round
	}
	if intValue(route["linked_cross_examination_count"]) > 0 {
		return intValue(route["linked_cross_examination_count"])
	}
	return 1
}

func runReleaseAdjacent(run map[string]any) bool {
	shipStatus := stringValue(run["ship_status"])
	reviewStatus := stringValue(run["review_status"])
	checksStatus := stringValue(run["checks_status"])
	return stringValue(run["pr_ref"]) != "" ||
		stringValue(run["rc_id"]) != "" ||
		(reviewStatus != "" && reviewStatus != "unknown") ||
		(checksStatus != "" && checksStatus != "unknown") ||
		(shipStatus != "" && shipStatus != "unstarted")
}

func inferMeetingDecisionProtocol(run map[string]any, question string, requestedProtocol string) (string, string) {
	if normalized := normalizeDecisionProtocol(requestedProtocol); normalized != "" {
		return normalized, fmt.Sprintf("Explicit protocol override: %s.", normalized)
	}
	lowerQuestion := strings.ToLower(strings.TrimSpace(question))
	releaseAdjacent := runReleaseAdjacent(run)
	priority := strings.ToUpper(firstNonBlank(stringValue(run["priority"]), stringValue(run["mission_priority"])))
	riskLevel := strings.ToLower(firstNonBlank(stringValue(run["risk_level"]), stringValue(run["effective_risk_level"])))
	if (releaseAdjacent || priority == "P0" || riskLevel == "high") &&
		(strings.Contains(lowerQuestion, "incident") ||
			strings.Contains(lowerQuestion, "rollback") ||
			strings.Contains(lowerQuestion, "restore") ||
			strings.Contains(lowerQuestion, "hotfix") ||
			strings.Contains(lowerQuestion, "release") ||
			strings.Contains(lowerQuestion, "production") ||
			strings.Contains(lowerQuestion, "rc") ||
			strings.Contains(lowerQuestion, "promote") ||
			strings.Contains(lowerQuestion, "deploy")) {
		return "incident-command", "Release-adjacent or urgent-risk decision requires incident-style coordination and bounded action ownership."
	}
	if strings.Contains(lowerQuestion, "decision") ||
		strings.Contains(lowerQuestion, "choose") ||
		strings.Contains(lowerQuestion, "tradeoff") ||
		strings.Contains(lowerQuestion, "record") ||
		strings.Contains(lowerQuestion, "architecture") ||
		strings.Contains(lowerQuestion, "approach") {
		return "decision-record", "Question is primarily about recording a bounded design choice and its consequences."
	}
	return "evidence-review", "Default evidence-review protocol fits contested factual or technical judgment."
}

func materializeMeeting(meetingID string, run map[string]any, question, requestedBy, requestedProtocol string) (map[string]any, error) {
	templateRoot, _ := toMap(loadData(meetingTemplatePath(), map[string]any{}))
	template, _ := toMap(templateRoot["meeting"])
	timestamp := nowISO()
	protocol, protocolReason := inferMeetingDecisionProtocol(run, question, requestedProtocol)
	meeting, _ := toMap(mergeDicts(template, map[string]any{
		"meeting_id":               meetingID,
		"mission_id":               stringValue(run["mission_id"]),
		"run_id":                   stringValue(run["run_id"]),
		"decision_protocol":        protocol,
		"decision_protocol_reason": protocolReason,
		"question":                 question,
		"scope":                    firstNonBlank(stringValue(run["objective"]), stringValue(run["mission_id"])),
		"leader_id":                stringValue(run["leader_id"]),
		"requested_by":             requestedBy,
		"state":                    "active",
		"authority_hierarchy":      defaultAuthorityHierarchy(),
		"locked_decisions":         []any{},
		"roles":                    defaultMeetingRoles(protocol),
		"role_specs":               defaultMeetingRoleSpecs(protocol),
		"packet": map[string]any{
			"expected_output": defaultMeetingExpectedOutput(protocol),
			"context_refs":    []any{stringValue(run["contract_ref"])},
			"known_context":   []any{stringValue(run["status_reason"])},
			"run_summary":     stringValue(run["latest_summary"]),
			"status_reason":   stringValue(run["status_reason"]),
		},
		"created_at": timestamp,
		"updated_at": timestamp,
		"closed_at":  "",
	}))
	ensureMeetingGovernanceState(meeting)
	if err := dumpData(meetingPath(meetingID, false), map[string]any{"version": 1, "meeting": meeting}); err != nil {
		return nil, err
	}
	return meeting, nil
}

func schedulerDecisionMeetingHostMissionID(decision map[string]any, requestedMissionID string) (string, error) {
	scope := firstNonBlank(stringValue(decision["decision_scope"]), "mission")
	bundleMissionIDs := dedupeStrings(append(
		anyListToStrings(toList(decision["bundle_mission_ids"])),
		anyListToStrings(toList(decision["mission_ids"]))...,
	))
	if requestedMissionID != "" {
		switch scope {
		case "bundle":
			if !containsString(bundleMissionIDs, requestedMissionID) {
				return "", fmt.Errorf("mission %s does not belong to scheduler decision %s", requestedMissionID, stringValue(decision["decision_id"]))
			}
		default:
			if requestedMissionID != stringValue(decision["mission_id"]) {
				return "", fmt.Errorf("mission %s does not match mission-scoped decision %s", requestedMissionID, stringValue(decision["decision_id"]))
			}
		}
		return requestedMissionID, nil
	}
	if scope != "bundle" {
		if missionID := strings.TrimSpace(stringValue(decision["mission_id"])); missionID != "" {
			return missionID, nil
		}
	}
	if governingMissionID := strings.TrimSpace(stringValue(decision["bundle_governing_mission_id"])); governingMissionID != "" {
		return governingMissionID, nil
	}
	for _, item := range toList(decision["candidate_governor_packets"]) {
		packet, _ := toMap(item)
		if missionID := strings.TrimSpace(stringValue(packet["mission_id"])); missionID != "" {
			return missionID, nil
		}
	}
	if len(bundleMissionIDs) > 0 {
		return bundleMissionIDs[0], nil
	}
	if missionID := strings.TrimSpace(stringValue(decision["mission_id"])); missionID != "" {
		return missionID, nil
	}
	return "", fmt.Errorf("scheduler decision %s does not resolve to a host mission", stringValue(decision["decision_id"]))
}

func schedulerDecisionMeetingHostRun(registry, missionRegistry map[string]any, missionID string) (map[string]any, error) {
	mission := readMission(missionID)
	if mission == nil {
		if index := findMissionIndex(missionRegistry, missionID); index >= 0 {
			mission, _ = toMap(toList(missionRegistry["missions"])[index])
		}
	}
	if mission != nil {
		if authoritativeRunID := strings.TrimSpace(stringValue(mission["authoritative_run_id"])); authoritativeRunID != "" {
			if index := findRunIndex(registry, authoritativeRunID); index >= 0 {
				run, _ := toMap(toList(registry["runs"])[index])
				return run, nil
			}
		}
	}
	runs := missionRuns(registry, missionID)
	run := authoritativeMissionRun(runs)
	if run == nil {
		return nil, fmt.Errorf("mission %s has no authoritative run for decision materialization", missionID)
	}
	return run, nil
}

func schedulerDecisionMeetingQuestion(decision map[string]any) string {
	scope := firstNonBlank(stringValue(decision["decision_scope"]), "mission")
	decisionKind := firstNonBlank(stringValue(decision["decision_kind"]), "decision")
	recommendedAction := strings.TrimSpace(stringValue(decision["recommended_action"]))
	note := strings.TrimSpace(stringValue(decision["note"]))
	switch scope {
	case "bundle":
		bundleID := firstNonBlank(strings.TrimSpace(stringValue(decision["bundle_id"])), "the affected mission bundle")
		if recommendedAction != "" {
			return fmt.Sprintf("Resolve scheduler decision %s for %s. Recommended action: %s", decisionKind, bundleID, recommendedAction)
		}
		if note != "" {
			return fmt.Sprintf("Resolve scheduler decision %s for %s. %s", decisionKind, bundleID, note)
		}
		return fmt.Sprintf("Resolve scheduler decision %s for %s.", decisionKind, bundleID)
	default:
		missionID := firstNonBlank(strings.TrimSpace(stringValue(decision["mission_id"])), "the affected mission")
		if recommendedAction != "" {
			return fmt.Sprintf("Resolve scheduler decision %s for %s. Recommended action: %s", decisionKind, missionID, recommendedAction)
		}
		if note != "" {
			return fmt.Sprintf("Resolve scheduler decision %s for %s. %s", decisionKind, missionID, note)
		}
		return fmt.Sprintf("Resolve scheduler decision %s for %s.", decisionKind, missionID)
	}
}

func attachSchedulerDecisionToMeeting(meeting, decision, snapshot map[string]any, hostMissionID string) map[string]any {
	meeting["source_decision_id"] = stringValue(decision["decision_id"])
	meeting["source_decision_scope"] = firstNonBlank(stringValue(decision["decision_scope"]), "mission")
	meeting["source_decision_kind"] = stringValue(decision["decision_kind"])
	meeting["source_decision_surface"] = stringValue(decision["recommended_decision_surface"])
	meeting["source_decision_snapshot_id"] = stringValue(snapshot["snapshot_id"])
	meeting["source_decision_snapshot_ref"] = firstNonBlank(stringValue(snapshot["latest_snapshot_path"]), schedulerLatestSnapshotPath())
	meeting["source_decision_packet"] = cloneValue(decision)

	packet, _ := toMap(meeting["packet"])
	if packet == nil {
		packet = map[string]any{}
	}
	contextRefs := toList(packet["context_refs"])
	contextRefs = appendUniqueText(contextRefs, firstNonBlank(stringValue(snapshot["latest_snapshot_path"]), schedulerLatestSnapshotPath()))
	if snapshotPath := strings.TrimSpace(stringValue(snapshot["snapshot_path"])); snapshotPath != "" {
		contextRefs = appendUniqueText(contextRefs, snapshotPath)
	}
	if hostMissionID != "" {
		contextRefs = appendUniqueText(contextRefs, missionPath(hostMissionID))
	}
	if scope := strings.TrimSpace(stringValue(decision["decision_scope"])); scope == "bundle" {
		if bundleID := strings.TrimSpace(stringValue(decision["bundle_id"])); bundleID != "" {
			contextRefs = appendUniqueText(contextRefs, bundleID)
		}
		bundleMissionIDs := dedupeStrings(append(
			anyListToStrings(toList(decision["bundle_mission_ids"])),
			anyListToStrings(toList(decision["mission_ids"]))...,
		))
		for _, missionID := range bundleMissionIDs {
			contextRefs = appendUniqueText(contextRefs, missionPath(missionID))
		}
	}
	packet["context_refs"] = contextRefs
	knownContext := toList(packet["known_context"])
	for _, item := range []string{
		strings.TrimSpace(stringValue(decision["note"])),
		strings.TrimSpace(stringValue(decision["recommended_action"])),
		strings.TrimSpace(stringValue(decision["recommended_protocol_reason"])),
	} {
		if item != "" {
			knownContext = appendUniqueText(knownContext, item)
		}
	}
	if scope := strings.TrimSpace(stringValue(decision["decision_scope"])); scope == "bundle" &&
		strings.TrimSpace(stringValue(decision["decision_kind"])) == "establish_bundle_governor" {
		packet["candidate_governor_packets"] = cloneValue(getKey(decision, "candidate_governor_packets", []any{}))
		for _, item := range toList(decision["candidate_governor_packets"]) {
			candidate, ok := toMap(item)
			if !ok {
				continue
			}
			missionID := strings.TrimSpace(stringValue(candidate["mission_id"]))
			if missionID == "" {
				continue
			}
			note := fmt.Sprintf(
				"Candidate governor %s: priority=%s, risk=%s, release_adjacent=%t, open_meeting_actions=%d, active_meetings=%d, root_candidate=%t.",
				missionID,
				firstNonBlank(strings.TrimSpace(stringValue(candidate["priority"])), "unknown"),
				firstNonBlank(strings.TrimSpace(stringValue(candidate["risk_level"])), "unknown"),
				boolValue(candidate["release_adjacent"]),
				intValue(candidate["open_meeting_actions"]),
				intValue(candidate["active_meeting_count"]),
				boolValue(candidate["is_root_candidate"]),
			)
			knownContext = appendUniqueText(knownContext, note)
		}
	}
	packet["known_context"] = knownContext
	meeting["packet"] = packet
	meeting["updated_at"] = nowISO()
	return meeting
}

func seedMeetingParticipantsFromSchedulerDecision(meeting, decision, registry, missionRegistry map[string]any, hostMissionID string) map[string]any {
	protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
	if protocol == "" {
		protocol = "evidence-review"
	}
	roles := anyListToStrings(toList(meeting["roles"]))
	if len(roles) == 0 {
		roles = anyListToStrings(defaultMeetingRoles(protocol))
	}
	if len(roles) == 0 {
		roles = []string{"participant"}
	}

	participants := toList(meeting["participants"])
	existingRunIDs := map[string]bool{}
	for _, item := range participants {
		participant, ok := toMap(item)
		if !ok {
			continue
		}
		runID := strings.TrimSpace(stringValue(participant["run_id"]))
		if runID != "" {
			existingRunIDs[runID] = true
		}
	}

	seedRun := func(run map[string]any, role string) {
		runID := strings.TrimSpace(stringValue(run["run_id"]))
		if runID == "" || existingRunIDs[runID] {
			return
		}
		participants = append(participants, map[string]any{
			"run_id":     runID,
			"owner_id":   stringValue(run["owner_id"]),
			"owner_type": stringValue(run["owner_type"]),
			"role":       firstNonBlank(strings.TrimSpace(role), stringValue(run["owner_id"]), "participant"),
			"joined_at":  nowISO(),
		})
		existingRunIDs[runID] = true
	}

	hostRun, err := schedulerDecisionMeetingHostRun(registry, missionRegistry, hostMissionID)
	if err == nil && hostRun != nil {
		seedRun(hostRun, roles[0])
	}

	if firstNonBlank(stringValue(decision["decision_scope"]), "mission") == "bundle" &&
		stringValue(decision["decision_kind"]) == "establish_bundle_governor" &&
		len(participants) < len(roles) {
		candidateMissionIDs := []string{}
		for _, item := range toList(decision["candidate_governor_packets"]) {
			packet, _ := toMap(item)
			if missionID := strings.TrimSpace(stringValue(packet["mission_id"])); missionID != "" {
				candidateMissionIDs = append(candidateMissionIDs, missionID)
			}
		}
		candidateMissionIDs = append(candidateMissionIDs, anyListToStrings(toList(decision["bundle_mission_ids"]))...)
		if governingMissionID := strings.TrimSpace(stringValue(decision["bundle_governing_mission_id"])); governingMissionID != "" {
			candidateMissionIDs = append(candidateMissionIDs, governingMissionID)
		}
		candidateMissionIDs = dedupeStrings(candidateMissionIDs)
		roleIndex := minInt(len(participants), len(roles)-1)
		for _, missionID := range candidateMissionIDs {
			if len(participants) >= len(roles) {
				break
			}
			run, err := schedulerDecisionMeetingHostRun(registry, missionRegistry, missionID)
			if err != nil || run == nil {
				continue
			}
			seedRun(run, roles[roleIndex])
			roleIndex = minInt(len(participants), len(roles)-1)
		}
	}

	meeting["participants"] = participants
	return meeting
}

func materializeMeetingFromSchedulerDecision(snapshot, decision, registry, missionRegistry map[string]any, requestedMissionID string) (map[string]any, map[string]any, map[string]any, error) {
	hostMissionID, err := schedulerDecisionMeetingHostMissionID(decision, requestedMissionID)
	if err != nil {
		return nil, nil, nil, err
	}
	decisionID := strings.TrimSpace(stringValue(decision["decision_id"]))
	if decisionID != "" {
		for _, item := range listMeetings(true) {
			activeMeeting, ok := toMap(item)
			if !ok || stringValue(activeMeeting["state"]) != "active" {
				continue
			}
			if strings.TrimSpace(stringValue(activeMeeting["source_decision_id"])) != decisionID {
				continue
			}
			if requestedMissionID != "" {
				meetingMissionID := firstNonBlank(
					strings.TrimSpace(stringValue(activeMeeting["mission_id"])),
					strings.TrimSpace(stringValue(getKey(activeMeeting, "source_decision_packet.mission_id", ""))),
				)
				if sourceScope := firstNonBlank(
					strings.TrimSpace(stringValue(activeMeeting["source_decision_scope"])),
					strings.TrimSpace(stringValue(getKey(activeMeeting, "source_decision_packet.decision_scope", ""))),
				); sourceScope == "bundle" {
					bundleMissionIDs := dedupeStrings(append(
						anyListToStrings(toList(getKey(activeMeeting, "source_decision_packet.bundle_mission_ids", []any{}))),
						anyListToStrings(toList(getKey(activeMeeting, "source_decision_packet.mission_ids", []any{})))...,
					))
					if !containsString(bundleMissionIDs, requestedMissionID) {
						continue
					}
				} else if meetingMissionID != requestedMissionID {
					continue
				}
			}
			activeMeeting = attachSchedulerDecisionToMeeting(activeMeeting, decision, snapshot, hostMissionID)
			activeMeeting = seedMeetingParticipantsFromSchedulerDecision(activeMeeting, decision, registry, missionRegistry, hostMissionID)
			if err := writeActiveMeeting(activeMeeting); err != nil {
				return nil, nil, nil, err
			}
			runID := strings.TrimSpace(stringValue(activeMeeting["run_id"]))
			run := map[string]any(nil)
			if runID != "" {
				if index := findRunIndex(registry, runID); index >= 0 {
					run, _ = toMap(toList(registry["runs"])[index])
				}
			}
			entry := map[string]any{
				"meeting_id":                  stringValue(activeMeeting["meeting_id"]),
				"run_id":                      stringValue(activeMeeting["run_id"]),
				"mission_id":                  stringValue(activeMeeting["mission_id"]),
				"question":                    stringValue(activeMeeting["question"]),
				"state":                       stringValue(activeMeeting["state"]),
				"meeting_ref":                 meetingPath(stringValue(activeMeeting["meeting_id"]), false),
				"created_at":                  firstNonBlank(stringValue(activeMeeting["created_at"]), nowISO()),
				"updated_at":                  firstNonBlank(stringValue(activeMeeting["updated_at"]), nowISO()),
				"source_decision_id":          decisionID,
				"source_decision_scope":       firstNonBlank(stringValue(activeMeeting["source_decision_scope"]), stringValue(decision["decision_scope"])),
				"source_decision_kind":        firstNonBlank(stringValue(activeMeeting["source_decision_kind"]), stringValue(decision["decision_kind"])),
				"source_decision_surface":     firstNonBlank(stringValue(activeMeeting["source_decision_surface"]), stringValue(decision["recommended_decision_surface"])),
				"source_decision_snapshot_id": firstNonBlank(stringValue(activeMeeting["source_decision_snapshot_id"]), stringValue(snapshot["snapshot_id"])),
				"reused_existing_meeting":     true,
			}
			return activeMeeting, entry, run, nil
		}
	}
	run, err := schedulerDecisionMeetingHostRun(registry, missionRegistry, hostMissionID)
	if err != nil {
		return nil, nil, nil, err
	}
	question := schedulerDecisionMeetingQuestion(decision)
	meetingID := newMeetingID(stringValue(run["mission_id"]))
	meeting, err := materializeMeeting(meetingID, run, question, "leader-main", stringValue(decision["recommended_decision_protocol"]))
	if err != nil {
		return nil, nil, nil, err
	}
	meeting = attachSchedulerDecisionToMeeting(meeting, decision, snapshot, hostMissionID)
	meeting = seedMeetingParticipantsFromSchedulerDecision(meeting, decision, registry, missionRegistry, hostMissionID)
	if err := writeActiveMeeting(meeting); err != nil {
		return nil, nil, nil, err
	}
	backlog := readMaintenanceBacklog()
	entry := map[string]any{
		"meeting_id":                  meetingID,
		"run_id":                      stringValue(run["run_id"]),
		"mission_id":                  stringValue(run["mission_id"]),
		"question":                    question,
		"state":                       "active",
		"meeting_ref":                 meetingPath(meetingID, false),
		"created_at":                  firstNonBlank(stringValue(meeting["created_at"]), nowISO()),
		"updated_at":                  firstNonBlank(stringValue(meeting["updated_at"]), nowISO()),
		"source_decision_id":          stringValue(decision["decision_id"]),
		"source_decision_scope":       firstNonBlank(stringValue(decision["decision_scope"]), "mission"),
		"source_decision_kind":        stringValue(decision["decision_kind"]),
		"source_decision_surface":     stringValue(decision["recommended_decision_surface"]),
		"source_decision_snapshot_id": stringValue(snapshot["snapshot_id"]),
		"reused_existing_meeting":     false,
	}
	backlog["meetings"] = append(toList(backlog["meetings"]), entry)
	if err := writeMaintenanceBacklog(backlog); err != nil {
		return nil, nil, nil, err
	}
	index := findRunIndex(registry, stringValue(run["run_id"]))
	if index >= 0 {
		runs := toList(registry["runs"])
		run["meeting_ref"] = meetingID
		run["state"] = "waiting"
		run["status_reason"] = "Scheduler decision materialized into meeting: " + stringValue(decision["decision_id"])
		run["next_expected_action"] = "Chair the decision meeting and adopt its conclusions back into scheduler-governed state."
		run["updated_at"] = nowISO()
		runs[index] = run
		registry["runs"] = runs
		if err := writeRegistry(registry); err != nil {
			return nil, nil, nil, err
		}
	}
	return meeting, entry, run, nil
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func updateContractFromRun(run map[string]any) error {
	contractRef := stringValue(run["contract_ref"])
	if contractRef == "" || !exists(contractRef) {
		return nil
	}
	data, _ := toMap(loadData(contractRef, map[string]any{}))
	contract, _ := toMap(data["contract"])
	contract["branch"] = stringValue(run["branch"])
	contract["worktree_path"] = stringValue(run["worktree_path"])
	freshness, _ := toMap(contract["freshness"])
	if freshness == nil {
		freshness = map[string]any{}
	}
	freshness["source_checked_at"] = nowISO()
	contract["freshness"] = freshness
	version := data["version"]
	if version == nil {
		version = 1
	}
	return dumpData(contractRef, map[string]any{"version": version, "contract": contract})
}

func reprovisionRun(run, profile map[string]any, mode string, applyChanges bool) (map[string]any, error) {
	if mode != "none" && mode != "branch" && mode != "worktree" {
		return nil, errors.New("reprovision mode must be one of: none, branch, worktree")
	}
	branch := stringValue(run["branch"])
	if branch == "" {
		return nil, errors.New("run has no branch; cannot reprovision")
	}
	baseBranch := stringValue(getKey(profile, "profile.base_branch", "main"))
	current := currentBranch()
	objective := firstNonBlank(stringValue(run["objective"]), stringValue(run["mission_id"]))
	slug := slugify(objective)
	worktreePath := ""
	if mode == "worktree" {
		worktreePath = provisionWorktreePath(profile, slug)
		run["worktree_path"] = worktreePath
	} else if mode == "none" {
		run["worktree_path"] = ""
	}

	provisioning := map[string]any{
		"mode":            mode,
		"apply_requested": applyChanges,
		"applied":         false,
		"base_branch":     baseBranch,
		"previous_branch": current,
		"target_branch":   branch,
	}
	if worktreePath != "" {
		provisioning["target_worktree_path"] = worktreePath
	}

	if applyChanges && mode != "none" {
		switch mode {
		case "branch":
			if err := gitCheckoutNewBranch(branch, baseBranch); err != nil {
				return nil, err
			}
			provisioning["applied"] = true
			provisioning["current_branch_after_apply"] = currentBranch()
		case "worktree":
			if err := gitCreateWorktree(branch, worktreePath, baseBranch); err != nil {
				return nil, err
			}
			if err := os.WriteFile(filepath.Join(worktreePath, ".mission"), []byte(fmt.Sprintf("%d\n", intValue(run["issue_number"]))), 0o644); err != nil {
				return nil, err
			}
			provisioning["applied"] = true
		}
	}

	if mode == "none" {
		run["mode"] = "solo"
	} else {
		run["mode"] = "managed"
	}
	run["provisioning"] = provisioning
	if err := updateContractFromRun(run); err != nil {
		return nil, err
	}
	return run, nil
}

func updateBacklogMeeting(meetingID string, updates map[string]any) (map[string]any, error) {
	backlog := readMaintenanceBacklog()
	meetings := toList(backlog["meetings"])
	for index, item := range meetings {
		entry, _ := toMap(item)
		if stringValue(entry["meeting_id"]) == meetingID {
			merged, _ := toMap(mergeDicts(entry, updates))
			meetings[index] = merged
			backlog["meetings"] = meetings
			if err := writeMaintenanceBacklog(backlog); err != nil {
				return nil, err
			}
			return merged, nil
		}
	}
	newEntry, _ := toMap(mergeDicts(map[string]any{"meeting_id": meetingID, "created_at": nowISO()}, updates))
	meetings = append(meetings, newEntry)
	backlog["meetings"] = meetings
	if err := writeMaintenanceBacklog(backlog); err != nil {
		return nil, err
	}
	return newEntry, nil
}

func writeActiveMeeting(meeting map[string]any) error {
	ensureMeetingGovernanceState(meeting)
	meeting["updated_at"] = nowISO()
	return dumpData(meetingPath(stringValue(meeting["meeting_id"]), false), map[string]any{"version": 1, "meeting": meeting})
}

func meetingHasOutput(meeting map[string]any) bool {
	findings, _ := toMap(meeting["findings"])
	return len(toList(findings["accepted"])) > 0 ||
		len(toList(findings["rejected"])) > 0 ||
		len(toList(findings["unresolved"])) > 0 ||
		len(toList(meeting["action_items"])) > 0
}

func textValue(item any) string {
	if mapping, ok := toMap(item); ok {
		return strings.TrimSpace(stringValue(mapping["summary"]))
	}
	return strings.TrimSpace(stringValue(item))
}

func meetingOutcome(meeting map[string]any) map[string]any {
	orchestration, _ := toMap(meeting["orchestration"])
	outputPacket, _ := toMap(orchestration["output_packet"])
	findings, _ := toMap(meeting["findings"])
	accepted := []string{}
	for _, item := range toList(findings["accepted"]) {
		text := textValue(item)
		if text != "" {
			accepted = append(accepted, text)
		}
	}
	actions := []string{}
	for _, item := range toList(meeting["action_items"]) {
		if mapping, ok := toMap(item); ok && stringValue(mapping["state"]) != "" && stringValue(mapping["state"]) != "open" {
			continue
		}
		text := textValue(item)
		if text != "" {
			actions = append(actions, text)
		}
	}
	summary := ""
	if len(accepted) > 0 {
		summary = accepted[0]
	} else {
		summary = firstNonBlank(strings.TrimSpace(stringValue(meeting["resolution_note"])), strings.TrimSpace(stringValue(meeting["question"])))
	}
	summary = firstNonBlank(strings.TrimSpace(stringValue(outputPacket["summary"])), summary)
	nextAction := "Resume execution with the meeting outcome applied."
	if len(actions) > 0 {
		nextAction = actions[0]
	}
	nextAction = firstNonBlank(strings.TrimSpace(stringValue(outputPacket["next_action"])), nextAction)
	return map[string]any{
		"summary":     summary,
		"next_action": nextAction,
	}
}

func meetingPrimaryFollowUpPacket(meeting map[string]any) map[string]any {
	orchestration, _ := toMap(meeting["orchestration"])
	outputPacket, _ := toMap(orchestration["output_packet"])
	return primaryFollowUpPacketFromList(toList(outputPacket["follow_up_packets"]))
}

func primaryFollowUpPacketFromList(packets []any) map[string]any {
	for _, item := range packets {
		packet, ok := toMap(item)
		if !ok {
			continue
		}
		state := firstNonBlank(strings.TrimSpace(stringValue(packet["state"])), "open")
		if state == "open" {
			return packet
		}
	}
	for _, item := range packets {
		packet, ok := toMap(item)
		if ok {
			return packet
		}
	}
	return nil
}

func followUpPacketOpenCount(packets []any) int {
	openCount := 0
	for _, item := range packets {
		packet, ok := toMap(item)
		if !ok {
			continue
		}
		state := firstNonBlank(strings.TrimSpace(stringValue(packet["state"])), "open")
		if state == "open" {
			openCount++
		}
	}
	return openCount
}

func meetingCurrentFollowUpPackets(meeting map[string]any) []any {
	protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
	if protocol == "" {
		protocol = "evidence-review"
	}
	debateLaneOutput, _ := toMap(meetingDebateLaneOutput(meeting))
	followUpPacketOutput, _ := toMap(meetingOutputFollowUpPackets(meeting, protocol, debateLaneOutput))
	return toList(followUpPacketOutput["follow_up_packets"])
}

func refreshMeetingOutputPacket(meeting map[string]any) map[string]any {
	orchestration, _ := toMap(meeting["orchestration"])
	if len(orchestration) == 0 {
		return meeting
	}
	outputPacket, _ := toMap(orchestration["output_packet"])
	if len(outputPacket) == 0 {
		return meeting
	}
	protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
	if protocol == "" {
		protocol = "evidence-review"
	}
	orchestration["output_packet"] = meetingAutoOutputPacket(meeting, protocol)
	meeting["orchestration"] = orchestration
	return meeting
}

func syncAdoptedFollowUpLifecycleForMeeting(registry map[string]any, meeting map[string]any) map[string]any {
	meetingID := strings.TrimSpace(stringValue(meeting["meeting_id"]))
	if meetingID == "" {
		return registry
	}
	packets := meetingCurrentFollowUpPackets(meeting)
	packetCount := len(packets)
	openPacketCount := followUpPacketOpenCount(packets)
	currentPrimary := primaryFollowUpPacketFromList(packets)
	currentPrimaryAction := adoptedFollowUpActionText(currentPrimary)
	runs := toList(registry["runs"])
	updated := false
	for index, item := range runs {
		run, _ := toMap(item)
		if strings.TrimSpace(stringValue(run["adopted_from_meeting_id"])) != meetingID {
			continue
		}
		previousPrimary := runPrimaryAdoptedFollowUpPacket(run)
		previousPrimaryAction := adoptedFollowUpActionText(previousPrimary)
		run["adopted_output_follow_up_packet_count"] = packetCount
		run["adopted_output_follow_up_packets"] = cloneValue(packets)
		if outputPacket := runAdoptedOutputPacket(run); len(outputPacket) > 0 {
			outputPacket["follow_up_packet_count"] = packetCount
			outputPacket["open_follow_up_packet_count"] = openPacketCount
			outputPacket["follow_up_packets"] = cloneValue(packets)
			run["adopted_output_packet"] = outputPacket
		}
		if currentPrimary != nil {
			run["adopted_primary_follow_up_packet"] = cloneValue(currentPrimary)
		} else {
			delete(run, "adopted_primary_follow_up_packet")
		}
		switch {
		case currentPrimaryAction != "" && strings.TrimSpace(stringValue(run["next_expected_action"])) == previousPrimaryAction && currentPrimaryAction != previousPrimaryAction:
			run["next_expected_action"] = currentPrimaryAction
		case currentPrimary == nil && previousPrimaryAction != "" && strings.TrimSpace(stringValue(run["next_expected_action"])) == previousPrimaryAction:
			run["next_expected_action"] = "Continue governed execution after the adopted follow-up is complete."
		}
		run["updated_at"] = nowISO()
		runs[index] = run
		updated = true
	}
	if updated {
		registry["runs"] = runs
	}
	return registry
}

func meetingAdoptedNextAction(meeting map[string]any) string {
	outcome := meetingOutcome(meeting)
	if primaryFollowUp := meetingPrimaryFollowUpPacket(meeting); primaryFollowUp != nil {
		if nextAction := adoptedFollowUpActionText(primaryFollowUp); nextAction != "" {
			return nextAction
		}
	}
	return stringValue(outcome["next_action"])
}

func applyMeetingAdoptionToRun(run map[string]any, meeting map[string]any) map[string]any {
	outcome := meetingOutcome(meeting)
	if stringValue(outcome["summary"]) != "" {
		run["latest_summary"] = outcome["summary"]
	}
	if nextAction := meetingAdoptedNextAction(meeting); nextAction != "" {
		run["next_expected_action"] = nextAction
	}
	run["adopted_from_meeting_id"] = stringValue(meeting["meeting_id"])

	orchestration, _ := toMap(meeting["orchestration"])
	outputPacket, _ := toMap(orchestration["output_packet"])
	if len(outputPacket) > 0 {
		run["adopted_output_packet_kind"] = stringValue(outputPacket["output_packet_kind"])
		run["adopted_output_packet"] = cloneValue(outputPacket)
		run["adopted_output_follow_up_packet_count"] = intValue(outputPacket["follow_up_packet_count"])
		run["adopted_output_follow_up_packets"] = cloneValue(getKey(outputPacket, "follow_up_packets", []any{}))
		if primaryFollowUp := meetingPrimaryFollowUpPacket(meeting); primaryFollowUp != nil {
			run["adopted_primary_follow_up_packet"] = cloneValue(primaryFollowUp)
		} else {
			delete(run, "adopted_primary_follow_up_packet")
		}
	} else {
		delete(run, "adopted_output_packet_kind")
		delete(run, "adopted_output_packet")
		delete(run, "adopted_output_follow_up_packet_count")
		delete(run, "adopted_output_follow_up_packets")
		delete(run, "adopted_primary_follow_up_packet")
	}
	return run
}

func runAdoptedOutputPacket(run map[string]any) map[string]any {
	outputPacket, _ := toMap(run["adopted_output_packet"])
	return outputPacket
}

func runAdoptedFollowUpPackets(run map[string]any) []any {
	packets := toList(run["adopted_output_follow_up_packets"])
	if len(packets) > 0 {
		return packets
	}
	outputPacket := runAdoptedOutputPacket(run)
	return toList(outputPacket["follow_up_packets"])
}

func runOpenAdoptedFollowUpPacketCount(run map[string]any) int {
	openCount := 0
	for _, item := range runAdoptedFollowUpPackets(run) {
		packet, ok := toMap(item)
		if !ok {
			continue
		}
		state := firstNonBlank(strings.TrimSpace(stringValue(packet["state"])), "open")
		if state == "open" {
			openCount++
		}
	}
	return openCount
}

func runPrimaryAdoptedFollowUpPacket(run map[string]any) map[string]any {
	packet, _ := toMap(run["adopted_primary_follow_up_packet"])
	if len(packet) > 0 {
		return packet
	}
	for _, item := range runAdoptedFollowUpPackets(run) {
		packet, ok := toMap(item)
		if !ok {
			continue
		}
		state := firstNonBlank(strings.TrimSpace(stringValue(packet["state"])), "open")
		if state == "open" {
			return packet
		}
	}
	for _, item := range runAdoptedFollowUpPackets(run) {
		packet, ok := toMap(item)
		if ok {
			return packet
		}
	}
	return nil
}

func adoptedFollowUpActionText(packet map[string]any) string {
	if len(packet) == 0 {
		return ""
	}
	return firstNonBlank(
		strings.TrimSpace(stringValue(packet["summary"])),
		strings.TrimSpace(stringValue(packet["required_follow_up"])),
	)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ensureMeetingGovernanceState(meeting map[string]any) map[string]any {
	if meeting == nil {
		return nil
	}
	protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
	if protocol == "" {
		protocol = "evidence-review"
	}
	meeting["decision_protocol"] = protocol
	if strings.TrimSpace(stringValue(meeting["decision_protocol_reason"])) == "" {
		meeting["decision_protocol_reason"] = "Legacy or unspecified meetings default to the evidence-review protocol."
	}
	if _, ok := toMap(meeting["orchestration"]); !ok {
		meeting["orchestration"] = map[string]any{}
	}
	packet, _ := toMap(meeting["packet"])
	if packet == nil {
		packet = map[string]any{}
	}
	if strings.TrimSpace(stringValue(packet["expected_output"])) == "" {
		packet["expected_output"] = defaultMeetingExpectedOutput(protocol)
	}
	meeting["packet"] = packet
	if len(toList(meeting["phase_history"])) == 0 {
		meeting["phase_history"] = []any{map[string]any{
			"phase":      firstNonBlank(stringValue(meeting["phase"]), "setup"),
			"entered_at": firstNonBlank(stringValue(meeting["created_at"]), nowISO()),
			"reason":     "meeting_materialized",
		}}
	}
	return meeting
}

func setMeetingPhase(meeting map[string]any, phase, reason string) {
	if meeting == nil {
		return
	}
	ensureMeetingGovernanceState(meeting)
	current := firstNonBlank(stringValue(meeting["phase"]), "setup")
	next := firstNonBlank(strings.TrimSpace(phase), current)
	orchestration, _ := toMap(meeting["orchestration"])
	if orchestration == nil {
		orchestration = map[string]any{}
	}
	timestamp := nowISO()
	if current != next {
		history := append(toList(meeting["phase_history"]), map[string]any{
			"phase":      next,
			"entered_at": timestamp,
			"reason":     firstNonBlank(strings.TrimSpace(reason), "phase_transition"),
		})
		meeting["phase_history"] = history
	}
	meeting["phase"] = next
	orchestration["last_phase_reason"] = firstNonBlank(strings.TrimSpace(reason), stringValue(orchestration["last_phase_reason"]))
	orchestration["last_phase_at"] = timestamp
	meeting["orchestration"] = orchestration
}

func meetingCoverageAudit(meeting map[string]any) map[string]any {
	participants := toList(meeting["participants"])
	positions := toList(meeting["positions"])
	independentByRun := map[string]int{}
	crossByRun := map[string]int{}
	linkedCrossByRun := map[string]int{}
	participantRoles := map[string]string{}
	participantRunIDs := []string{}

	for _, item := range participants {
		participant, ok := toMap(item)
		if !ok {
			continue
		}
		runID := stringValue(participant["run_id"])
		if runID == "" {
			continue
		}
		participantRunIDs = append(participantRunIDs, runID)
		participantRoles[runID] = firstNonBlank(stringValue(participant["role"]), stringValue(participant["owner_id"]), runID)
	}

	for _, item := range positions {
		position, ok := toMap(item)
		if !ok {
			continue
		}
		runID := stringValue(position["run_id"])
		if runID == "" {
			continue
		}
		switch stringValue(position["phase"]) {
		case "independent_review":
			independentByRun[runID]++
		case "cross_examination":
			crossByRun[runID]++
			if strings.TrimSpace(stringValue(position["responds_to_run_id"])) != "" ||
				strings.TrimSpace(stringValue(position["responds_to_position_id"])) != "" ||
				strings.TrimSpace(stringValue(position["responds_to_role"])) != "" {
				linkedCrossByRun[runID]++
			}
		}
	}

	missingIndependent := []string{}
	missingCross := []string{}
	for _, runID := range participantRunIDs {
		role := firstNonBlank(participantRoles[runID], runID)
		if independentByRun[runID] == 0 {
			missingIndependent = append(missingIndependent, role)
		}
		if crossByRun[runID] == 0 {
			missingCross = append(missingCross, role)
		}
	}
	sort.Strings(missingIndependent)
	sort.Strings(missingCross)

	coverageNotes := []string{}
	if len(missingIndependent) > 0 {
		coverageNotes = append(coverageNotes, fmt.Sprintf("participants missing independent review: %s", strings.Join(missingIndependent, ", ")))
	}
	if len(missingCross) > 0 {
		coverageNotes = append(coverageNotes, fmt.Sprintf("participants without cross-examination positions: %s", strings.Join(missingCross, ", ")))
	}

	requiredIndependent := minInt(len(participantRunIDs), 2)
	blockingGaps := []string{}
	if len(participantRunIDs) > 0 && len(independentByRun) < requiredIndependent {
		blockingGaps = append(blockingGaps, fmt.Sprintf("need at least %d distinct independent-review participants before debate can advance", requiredIndependent))
	}
	linkedCrossCount := 0
	for _, count := range linkedCrossByRun {
		linkedCrossCount += count
	}
	if len(crossByRun) > 0 && linkedCrossCount == 0 {
		blockingGaps = append(blockingGaps, "cross-examination positions exist but none explicitly respond to another role or position")
	}

	return map[string]any{
		"participant_count":                len(participantRunIDs),
		"distinct_independent_reviewers":   len(independentByRun),
		"distinct_cross_examiners":         len(crossByRun),
		"linked_cross_examination_count":   linkedCrossCount,
		"participants_missing_independent": stringListToAny(missingIndependent),
		"participants_missing_cross_exam":  stringListToAny(missingCross),
		"coverage_notes":                   stringListToAny(coverageNotes),
		"blocking_gaps":                    stringListToAny(blockingGaps),
	}
}

func meetingWavePlan(meeting map[string]any, coverage map[string]any) map[string]any {
	participants := toList(meeting["participants"])
	positions := toList(meeting["positions"])
	independentByRun := map[string]int{}
	for _, item := range positions {
		position, ok := toMap(item)
		if !ok || stringValue(position["phase"]) != "independent_review" {
			continue
		}
		runID := strings.TrimSpace(stringValue(position["run_id"]))
		if runID != "" {
			independentByRun[runID]++
		}
	}

	waveActions := []any{}
	for _, item := range participants {
		participant, ok := toMap(item)
		if !ok {
			continue
		}
		runID := strings.TrimSpace(stringValue(participant["run_id"]))
		role := firstNonBlank(stringValue(participant["role"]), stringValue(participant["owner_id"]), runID)
		if runID == "" || independentByRun[runID] > 0 {
			continue
		}
		waveActions = append(waveActions, map[string]any{
			"wave_id":         "wave-independent-review-" + slugify(runID),
			"wave_kind":       "independent_review",
			"assigned_role":   role,
			"assigned_run_id": runID,
			"expected_phase":  "independent_review",
			"state":           "open",
			"summary":         fmt.Sprintf("%s still owes an independent-review position before debate can advance cleanly.", firstNonBlank(role, runID)),
		})
	}

	return map[string]any{
		"wave_action_count":      len(waveActions),
		"open_wave_action_count": len(waveActions),
		"wave_actions":           waveActions,
		"coverage_note_count":    len(toList(coverage["coverage_notes"])),
	}
}

func preferredMeetingPosition(positions []map[string]any, stance string) map[string]any {
	var best map[string]any
	bestRank := -1
	bestRecordedAt := ""
	for _, position := range positions {
		if strings.TrimSpace(stringValue(position["stance"])) != stance {
			continue
		}
		rank := authorityRank(firstNonBlank(stringValue(position["authority_level"]), "L4"))
		recordedAt := stringValue(position["recorded_at"])
		if best == nil || rank > bestRank || (rank == bestRank && isoAfter(recordedAt, bestRecordedAt)) {
			best = position
			bestRank = rank
			bestRecordedAt = recordedAt
		}
	}
	return best
}

func meetingClaimGroupKey(claimRef, dimension, fallback string) string {
	claim := firstNonBlank(strings.TrimSpace(claimRef), strings.TrimSpace(fallback))
	if claim == "" {
		claim = "unspecified-claim"
	}
	normalized := slugify(strings.ToLower(claim))
	if normalized == "" {
		normalized = slugify(strings.ToLower(firstNonBlank(fallback, claim)))
	}
	if normalized == "" {
		normalized = "unspecified-claim"
	}
	return firstNonBlank(strings.TrimSpace(dimension), "factual") + "::" + normalized
}

func meetingHasRecordedFindingForClaim(meeting map[string]any, claimRef, dimension string) bool {
	findings, _ := toMap(meeting["findings"])
	if findings == nil {
		return false
	}
	for _, bucket := range []string{"accepted", "rejected", "unresolved"} {
		for _, item := range toList(findings[bucket]) {
			finding, ok := toMap(item)
			if !ok {
				continue
			}
			if strings.TrimSpace(stringValue(finding["claim_ref"])) == strings.TrimSpace(claimRef) &&
				firstNonBlank(stringValue(finding["dimension"]), "factual") == firstNonBlank(strings.TrimSpace(dimension), "factual") {
				return true
			}
		}
	}
	return false
}

func meetingAutoLockBundleGovernorDecision(meeting, registry map[string]any) map[string]any {
	sourceScope := firstNonBlank(
		strings.TrimSpace(stringValue(meeting["source_decision_scope"])),
		strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_scope", ""))),
	)
	sourceKind := firstNonBlank(
		strings.TrimSpace(stringValue(meeting["source_decision_kind"])),
		strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_kind", ""))),
	)
	if sourceScope != "bundle" || sourceKind != "establish_bundle_governor" {
		return nil
	}
	for _, item := range toList(meeting["locked_decisions"]) {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		if strings.TrimSpace(stringValue(decision["governing_mission_id"])) != "" {
			return nil
		}
	}
	packet, _ := toMap(meeting["source_decision_packet"])
	bundleID := firstNonBlank(
		strings.TrimSpace(stringValue(packet["bundle_id"])),
		strings.TrimSpace(stringValue(meeting["bundle_id"])),
	)
	candidateMissionIDs := dedupeStrings(append(
		anyListToStrings(toList(packet["bundle_mission_ids"])),
		anyListToStrings(toList(packet["mission_ids"]))...,
	))
	if governingMissionID := strings.TrimSpace(stringValue(packet["bundle_governing_mission_id"])); governingMissionID != "" {
		candidateMissionIDs = append(candidateMissionIDs, governingMissionID)
	}
	candidateMissionIDs = dedupeStrings(candidateMissionIDs)
	if bundleID == "" || len(candidateMissionIDs) == 0 {
		return nil
	}
	runMissionByID := map[string]string{}
	for _, item := range toList(registry["runs"]) {
		run, ok := toMap(item)
		if !ok {
			continue
		}
		runID := strings.TrimSpace(stringValue(run["run_id"]))
		missionID := strings.TrimSpace(stringValue(run["mission_id"]))
		if runID != "" && missionID != "" {
			runMissionByID[runID] = missionID
		}
	}
	findings, _ := toMap(meeting["findings"])
	type candidate struct {
		missionID      string
		authorityRank  int
		authorityLevel string
		recordedAt     string
		summary        string
		reason         string
		evidenceRefs   []any
	}
	var selected *candidate
	for _, item := range toList(findings["accepted"]) {
		finding, ok := toMap(item)
		if !ok {
			continue
		}
		governingMissionID := strings.TrimSpace(stringValue(finding["governing_mission_id"]))
		if governingMissionID == "" {
			governingMissionID = strings.TrimSpace(runMissionByID[stringValue(finding["source_run_id"])])
		}
		if governingMissionID == "" || !containsString(candidateMissionIDs, governingMissionID) {
			continue
		}
		current := candidate{
			missionID:      governingMissionID,
			authorityRank:  authorityRank(firstNonBlank(stringValue(finding["authority_level"]), "L4")),
			authorityLevel: firstNonBlank(stringValue(finding["authority_level"]), "L4"),
			recordedAt:     firstNonBlank(stringValue(finding["recorded_at"]), nowISO()),
			summary:        firstNonBlank(strings.TrimSpace(stringValue(finding["summary"])), fmt.Sprintf("Bundle governance anchored to %s.", governingMissionID)),
			reason:         strings.TrimSpace(stringValue(finding["reason"])),
			evidenceRefs:   toList(finding["evidence_refs"]),
		}
		if selected == nil || current.authorityRank > selected.authorityRank || (current.authorityRank == selected.authorityRank && isoAfter(current.recordedAt, selected.recordedAt)) {
			selected = &current
		}
	}
	if selected == nil {
		return nil
	}
	return map[string]any{
		"decision_id":          "decision-" + nowTime().UTC().Format("20060102150405"),
		"summary":              firstNonBlank(selected.summary, fmt.Sprintf("Anchor bundle governance to %s.", selected.missionID)),
		"reason":               firstNonBlank(selected.reason, fmt.Sprintf("Auto-locked from accepted council finding adopting %s as the bundle governor.", selected.missionID)),
		"authority_level":      selected.authorityLevel,
		"evidence_refs":        cloneValue(selected.evidenceRefs),
		"locked_at":            nowISO(),
		"governing_mission_id": selected.missionID,
		"bundle_id":            bundleID,
		"generated_by":         "orchestrate_auto_synthesis",
	}
}

func meetingDebatePlan(meeting map[string]any) map[string]any {
	ensureMeetingGovernanceState(meeting)
	protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
	if protocol == "" {
		protocol = "evidence-review"
	}
	orchestration, _ := toMap(meeting["orchestration"])
	previousDebatePlan, _ := toMap(orchestration["debate_plan"])
	previousCurrentWaveBatchID := strings.TrimSpace(stringValue(previousDebatePlan["current_wave_batch_id"]))
	positions := toList(meeting["positions"])
	participants := toList(meeting["participants"])
	maxDebateRounds := meetingMaxDebateRounds()
	claimGroups := map[string][]map[string]any{}
	claimLabels := map[string]string{}
	dimensions := map[string]string{}
	linkedCrossByClaim := map[string]int{}
	participantRunByRole := map[string]string{}

	for _, item := range participants {
		participant, ok := toMap(item)
		if !ok {
			continue
		}
		role := strings.TrimSpace(stringValue(participant["role"]))
		runID := strings.TrimSpace(stringValue(participant["run_id"]))
		if role != "" && runID != "" && participantRunByRole[role] == "" {
			participantRunByRole[role] = runID
		}
	}

	for _, item := range positions {
		position, ok := toMap(item)
		if !ok {
			continue
		}
		claimRef := strings.TrimSpace(stringValue(position["claim_ref"]))
		if claimRef == "" {
			continue
		}
		dimension := firstNonBlank(stringValue(position["dimension"]), "factual")
		key := meetingClaimGroupKey(claimRef, dimension, stringValue(position["position_id"]))
		claimLabels[key] = claimRef
		dimensions[key] = dimension
		claimGroups[key] = append(claimGroups[key], position)
		if stringValue(position["phase"]) == "cross_examination" &&
			(strings.TrimSpace(stringValue(position["responds_to_run_id"])) != "" ||
				strings.TrimSpace(stringValue(position["responds_to_role"])) != "" ||
				strings.TrimSpace(stringValue(position["responds_to_position_id"])) != "") {
			linkedCrossByClaim[key]++
		}
	}

	groupKeys := []string{}
	for key := range claimGroups {
		groupKeys = append(groupKeys, key)
	}
	sort.Strings(groupKeys)

	routes := []any{}
	routeActions := []any{}
	roundActions := []any{}
	blockingGaps := []string{}
	coverageNotes := []string{}
	conflictClaimCount := 0
	resolvedRouteCount := 0
	resolvedWithoutCrossExaminationCount := 0
	authorityResolvedRouteCount := 0
	protocolDeferredRouteCount := 0
	protocolBoundedRouteCount := 0
	plannedRoundCount := 0
	resolvedRoundCount := 0

	for _, key := range groupKeys {
		group := claimGroups[key]
		supportRoles := []string{}
		challengeRoles := []string{}
		supportIDs := []string{}
		challengeIDs := []string{}
		linkedCrossPositions := []map[string]any{}
		for _, position := range group {
			switch strings.TrimSpace(stringValue(position["stance"])) {
			case "support":
				supportRoles = append(supportRoles, firstNonBlank(stringValue(position["role"]), stringValue(position["run_id"])))
				supportIDs = append(supportIDs, stringValue(position["position_id"]))
			case "challenge":
				challengeRoles = append(challengeRoles, firstNonBlank(stringValue(position["role"]), stringValue(position["run_id"])))
				challengeIDs = append(challengeIDs, stringValue(position["position_id"]))
			}
			if stringValue(position["phase"]) == "cross_examination" &&
				(strings.TrimSpace(stringValue(position["responds_to_run_id"])) != "" ||
					strings.TrimSpace(stringValue(position["responds_to_role"])) != "" ||
					strings.TrimSpace(stringValue(position["responds_to_position_id"])) != "") {
				linkedCrossPositions = append(linkedCrossPositions, position)
			}
		}
		supportRoles = dedupeStrings(supportRoles)
		challengeRoles = dedupeStrings(challengeRoles)
		if len(supportRoles) == 0 || len(challengeRoles) == 0 {
			continue
		}
		leadSupport := preferredMeetingPosition(group, "support")
		leadChallenge := preferredMeetingPosition(group, "challenge")
		leadSupportRank := authorityRank(firstNonBlank(stringValue(leadSupport["authority_level"]), "L4"))
		leadChallengeRank := authorityRank(firstNonBlank(stringValue(leadChallenge["authority_level"]), "L4"))
		assignedRole := firstNonBlank(stringValue(leadChallenge["role"]), firstNonBlank(challengeRoles...))
		assignedRunID := firstNonBlank(stringValue(leadChallenge["run_id"]), participantRunByRole[assignedRole])
		targetPositionID := firstNonBlank(stringValue(leadSupport["position_id"]), stringValue(leadChallenge["position_id"]))
		targetRunID := firstNonBlank(stringValue(leadSupport["run_id"]), stringValue(leadChallenge["run_id"]))
		targetRole := firstNonBlank(stringValue(leadSupport["role"]), stringValue(leadChallenge["role"]))
		routeID := firstNonBlank("route-"+slugify(key), fmt.Sprintf("route-%d", conflictClaimCount+1))
		conflictClaimCount++
		hasRecordedFinding := meetingHasRecordedFindingForClaim(meeting, claimLabels[key], dimensions[key])
		route := map[string]any{
			"route_id":                       routeID,
			"claim_key":                      key,
			"claim_ref":                      claimLabels[key],
			"dimension":                      dimensions[key],
			"support_roles":                  stringListToAny(supportRoles),
			"challenge_roles":                stringListToAny(challengeRoles),
			"supporting_position_ids":        stringListToAny(dedupeStrings(supportIDs)),
			"challenging_position_ids":       stringListToAny(dedupeStrings(challengeIDs)),
			"assigned_role":                  assignedRole,
			"assigned_run_id":                assignedRunID,
			"expected_stance":                "challenge",
			"target_position_id":             targetPositionID,
			"target_run_id":                  targetRunID,
			"target_role":                    targetRole,
			"max_debate_rounds":              maxDebateRounds,
			"linked_cross_examination_count": linkedCrossByClaim[key],
			"status":                         "awaiting_cross_examination",
		}
		if linkedCrossByClaim[key] > 0 {
			route["status"] = "routed"
			resolvedRouteCount++
		} else if leadSupportRank != leadChallengeRank {
			prevailingStance := "support"
			prevailingLead := leadSupport
			prevailingRole := stringValue(leadSupport["role"])
			prevailingRunID := stringValue(leadSupport["run_id"])
			prevailingPositionID := stringValue(leadSupport["position_id"])
			overruledStance := "challenge"
			overruledLead := leadChallenge
			if leadChallengeRank > leadSupportRank {
				prevailingStance = "challenge"
				prevailingLead = leadChallenge
				prevailingRole = stringValue(leadChallenge["role"])
				prevailingRunID = stringValue(leadChallenge["run_id"])
				prevailingPositionID = stringValue(leadChallenge["position_id"])
				overruledStance = "support"
				overruledLead = leadSupport
			}
			resolvedRouteCount++
			resolvedWithoutCrossExaminationCount++
			authorityResolvedRouteCount++
			route["status"] = "authority_resolved"
			route["resolution_kind"] = "higher_authority_prevails"
			route["prevailing_stance"] = prevailingStance
			route["prevailing_role"] = prevailingRole
			route["prevailing_run_id"] = prevailingRunID
			route["prevailing_position_id"] = prevailingPositionID
			route["overruled_stance"] = overruledStance
			route["overruled_role"] = stringValue(overruledLead["role"])
			route["overruled_run_id"] = stringValue(overruledLead["run_id"])
			route["overruled_position_id"] = stringValue(overruledLead["position_id"])
			route["required_follow_up"] = "synthesis_can_materialize_verdict"
			route["resolved_without_cross_examination"] = true
			coverageNotes = append(coverageNotes, fmt.Sprintf("claim %q does not need cross-examination because %s authority from %s already outranks %s authority from %s.", claimLabels[key], firstNonBlank(stringValue(prevailingLead["authority_level"]), "L4"), firstNonBlank(prevailingRole, prevailingStance+" lane"), firstNonBlank(stringValue(overruledLead["authority_level"]), "L4"), firstNonBlank(stringValue(overruledLead["role"]), overruledStance+" lane")))
		} else {
			if protocolUsesStructuredDebate(protocol) {
				note := fmt.Sprintf("claim %q lacks linked cross-examination between supporting and challenging roles", claimLabels[key])
				if assignedRole != "" {
					note = fmt.Sprintf("%s; %s should challenge %s on this lane", note, assignedRole, firstNonBlank(targetRole, claimLabels[key]))
				}
				blockingGaps = append(blockingGaps, note)
				routeActions = append(routeActions, map[string]any{
					"route_id":                routeID,
					"claim_ref":               claimLabels[key],
					"dimension":               dimensions[key],
					"assigned_role":           assignedRole,
					"assigned_run_id":         assignedRunID,
					"expected_phase":          "cross_examination",
					"expected_stance":         "challenge",
					"responds_to_position_id": targetPositionID,
					"responds_to_run_id":      targetRunID,
					"responds_to_role":        targetRole,
					"state":                   "open",
					"summary":                 fmt.Sprintf("%s should challenge %s on claim %q", firstNonBlank(assignedRole, "bounded challenger"), firstNonBlank(targetRole, "the supporting lane"), claimLabels[key]),
				})
			} else {
				protocolDeferredRouteCount++
				route["status"] = "protocol_direct_synthesis"
				route["resolution_kind"] = "protocol_prefers_direct_synthesis"
				route["required_follow_up"] = "synthesis_must_record_tradeoff"
				route["decision_protocol"] = protocol
				coverageNotes = append(coverageNotes, protocolDirectSynthesisCoverageNote(protocol, claimLabels[key]))
			}
		}
		if linkedCrossByClaim[key] > 0 && !hasRecordedFinding {
			plannedRoundCount++
			latestLinked := map[string]any(nil)
			latestRecordedAt := ""
			for _, position := range linkedCrossPositions {
				recordedAt := stringValue(position["recorded_at"])
				if latestLinked == nil || isoAfter(recordedAt, latestRecordedAt) {
					latestLinked = position
					latestRecordedAt = recordedAt
				}
			}
			latestStance := strings.TrimSpace(stringValue(latestLinked["stance"]))
			if latestLinked != nil {
				route["latest_stance"] = latestStance
				route["latest_position_id"] = stringValue(latestLinked["position_id"])
				route["latest_run_id"] = stringValue(latestLinked["run_id"])
				route["latest_role"] = stringValue(latestLinked["role"])
			}
			if linkedCrossByClaim[key] >= 2 {
				route["latest_debate_round"] = linkedCrossByClaim[key]
			}
			if linkedCrossByClaim[key] >= maxDebateRounds {
				resolvedRoundCount++
				route["status"] = "max_rounds_reached"
				route["resolution_kind"] = "max_rounds_reached"
				route["required_follow_up"] = "leader_judgment_required"
				coverageNotes = append(coverageNotes, fmt.Sprintf("claim %q reached the debate round cap (%d); hand final adjudication to synthesis.", claimLabels[key], maxDebateRounds))
			} else if latestLinked != nil {
				opposingStance := "support"
				opposingRoles := supportRoles
				opposingLead := leadSupport
				if latestStance == "support" {
					opposingStance = "challenge"
					opposingRoles = challengeRoles
					opposingLead = leadChallenge
				}
				latestAuthorityRank := authorityRank(firstNonBlank(stringValue(latestLinked["authority_level"]), "L4"))
				opposingAuthorityRank := authorityRank(firstNonBlank(stringValue(opposingLead["authority_level"]), "L4"))
				opposingRecordedAt := stringValue(opposingLead["recorded_at"])
				if latestAuthorityRank > opposingAuthorityRank || (linkedCrossByClaim[key] >= 2 && latestAuthorityRank == opposingAuthorityRank && isoAfter(latestRecordedAt, opposingRecordedAt)) {
					resolvedRoundCount++
					route["status"] = "debate_resolved"
					route["resolution_kind"] = "latest_response_prevails"
					route["prevailing_stance"] = latestStance
					route["prevailing_role"] = stringValue(latestLinked["role"])
					route["prevailing_run_id"] = stringValue(latestLinked["run_id"])
					route["prevailing_position_id"] = stringValue(latestLinked["position_id"])
					route["overruled_stance"] = opposingStance
					route["overruled_role"] = stringValue(opposingLead["role"])
					route["overruled_run_id"] = stringValue(opposingLead["run_id"])
					route["overruled_position_id"] = stringValue(opposingLead["position_id"])
					route["required_follow_up"] = "synthesis_can_materialize_verdict"
					coverageNotes = append(coverageNotes, fmt.Sprintf("claim %q debate resolves after round %d because the latest %s response now outranks the opposing lane.", claimLabels[key], linkedCrossByClaim[key], latestStance))
				} else if !protocolUsesStructuredDebate(protocol) {
					protocolBoundedRouteCount++
					route["status"] = "protocol_bound_reached"
					route["resolution_kind"] = "protocol_bound_reached"
					route["required_follow_up"] = "synthesis_must_record_tradeoff"
					route["decision_protocol"] = protocol
					coverageNotes = append(coverageNotes, protocolBoundedRoundCoverageNote(protocol, claimLabels[key], linkedCrossByClaim[key]))
				} else {
					assignedRole := firstNonBlank(stringValue(opposingLead["role"]), firstNonBlank(opposingRoles...))
					assignedRunID := firstNonBlank(stringValue(opposingLead["run_id"]), participantRunByRole[assignedRole])
					nextRoundNumber := linkedCrossByClaim[key] + 1
					roundActions = append(roundActions, map[string]any{
						"round_id":                fmt.Sprintf("%s-round-%d", routeID, nextRoundNumber),
						"route_id":                routeID,
						"claim_ref":               claimLabels[key],
						"dimension":               dimensions[key],
						"assigned_role":           assignedRole,
						"assigned_run_id":         assignedRunID,
						"expected_phase":          "cross_examination",
						"expected_stance":         opposingStance,
						"responds_to_position_id": stringValue(latestLinked["position_id"]),
						"responds_to_run_id":      stringValue(latestLinked["run_id"]),
						"responds_to_role":        stringValue(latestLinked["role"]),
						"round_number":            nextRoundNumber,
						"state":                   "open",
						"summary":                 fmt.Sprintf("%s should answer round %d with a %s response on claim %q before synthesis.", firstNonBlank(assignedRole, "bounded respondent"), nextRoundNumber, opposingStance, claimLabels[key]),
					})
					route["status"] = fmt.Sprintf("awaiting_round_%d", nextRoundNumber)
					route["resolution_kind"] = "needs_follow_up_round"
					coverageNotes = append(coverageNotes, fmt.Sprintf("claim %q needs a round %d %s response before synthesis.", claimLabels[key], nextRoundNumber, opposingStance))
				}
			}
		}
		coverageNotes = append(coverageNotes, fmt.Sprintf("claim %q: support=%s challenge=%s assigned=%s", claimLabels[key], strings.Join(supportRoles, ", "), strings.Join(challengeRoles, ", "), firstNonBlank(assignedRole, "unassigned")))
		routes = append(routes, route)
	}

	type waveBatch struct {
		Key          string
		Phase        string
		RoundNumber  int
		AssignedRole string
		AssignedRun  string
		ActionIDs    []string
		ClaimRefs    []string
	}
	waveBatchByKey := map[string]*waveBatch{}
	appendWaveBatch := func(action map[string]any) {
		assignedRole := strings.TrimSpace(stringValue(action["assigned_role"]))
		assignedRun := strings.TrimSpace(stringValue(action["assigned_run_id"]))
		expectedPhase := firstNonBlank(strings.TrimSpace(stringValue(action["expected_phase"])), "cross_examination")
		roundNumber := intValue(action["round_number"])
		if roundNumber <= 0 {
			roundNumber = 1
		}
		actorKey := firstNonBlank(assignedRun, assignedRole, "unassigned")
		key := fmt.Sprintf("%s::%d::%s", expectedPhase, roundNumber, actorKey)
		batch := waveBatchByKey[key]
		if batch == nil {
			batch = &waveBatch{
				Key:          key,
				Phase:        expectedPhase,
				RoundNumber:  roundNumber,
				AssignedRole: assignedRole,
				AssignedRun:  assignedRun,
			}
			waveBatchByKey[key] = batch
		}
		actionID := firstNonBlank(stringValue(action["round_id"]), stringValue(action["route_id"]))
		if actionID != "" {
			batch.ActionIDs = append(batch.ActionIDs, actionID)
		}
		if claimRef := strings.TrimSpace(stringValue(action["claim_ref"])); claimRef != "" {
			batch.ClaimRefs = append(batch.ClaimRefs, claimRef)
		}
	}
	for _, item := range routeActions {
		action, _ := toMap(item)
		appendWaveBatch(action)
	}
	for _, item := range roundActions {
		action, _ := toMap(item)
		appendWaveBatch(action)
	}
	waveBatchKeys := []string{}
	for key := range waveBatchByKey {
		waveBatchKeys = append(waveBatchKeys, key)
	}
	sort.Strings(waveBatchKeys)
	waveBatchList := []*waveBatch{}
	for _, key := range waveBatchKeys {
		if batch := waveBatchByKey[key]; batch != nil {
			waveBatchList = append(waveBatchList, batch)
		}
	}
	phaseOrder := func(phase string) int {
		switch strings.TrimSpace(phase) {
		case "cross_examination":
			return 0
		case "independent_review":
			return 1
		case "synthesis":
			return 2
		case "output":
			return 3
		default:
			return 9
		}
	}
	sort.SliceStable(waveBatchList, func(i, j int) bool {
		left := waveBatchList[i]
		right := waveBatchList[j]
		if left == nil || right == nil {
			return left != nil
		}
		if phaseOrder(left.Phase) != phaseOrder(right.Phase) {
			return phaseOrder(left.Phase) < phaseOrder(right.Phase)
		}
		if left.RoundNumber != right.RoundNumber {
			return left.RoundNumber < right.RoundNumber
		}
		leftActor := firstNonBlank(left.AssignedRun, left.AssignedRole, "unassigned")
		rightActor := firstNonBlank(right.AssignedRun, right.AssignedRole, "unassigned")
		if leftActor != rightActor {
			return leftActor < rightActor
		}
		leftClaims := strings.Join(dedupeStrings(left.ClaimRefs), "|")
		rightClaims := strings.Join(dedupeStrings(right.ClaimRefs), "|")
		if leftClaims != rightClaims {
			return leftClaims < rightClaims
		}
		return left.Key < right.Key
	})
	waveBatches := []any{}
	maxWaveLaneCount := 0
	multiLaneBatchCount := 0
	openWaveBatchCount := 0
	readyWaveBatchCount := 0
	blockedWaveBatchCount := 0
	resolvedWaveBatchCount := 0
	currentWaveBatchID := ""
	currentWaveBatchLaneCount := 0
	currentWaveBatchSummary := ""
	batchSelectionPolicy := "debate_batch_priority"
	previousCurrentWaveBatchStillOpen := false
	if previousCurrentWaveBatchID != "" {
		for _, batch := range waveBatchList {
			if batch == nil {
				continue
			}
			actionIDs := dedupeStrings(batch.ActionIDs)
			if len(actionIDs) == 0 {
				continue
			}
			batchID := "debate-batch-" + stableSlugID(batch.Key)
			if batchID == previousCurrentWaveBatchID {
				previousCurrentWaveBatchStillOpen = true
				break
			}
		}
	}
	for _, batch := range waveBatchList {
		if batch == nil {
			continue
		}
		laneCount := len(dedupeStrings(batch.ClaimRefs))
		if laneCount > maxWaveLaneCount {
			maxWaveLaneCount = laneCount
		}
		if laneCount > 1 {
			multiLaneBatchCount++
		}
		actionIDs := dedupeStrings(batch.ActionIDs)
		state := "resolved"
		current := false
		batchID := "debate-batch-" + stableSlugID(batch.Key)
		if len(actionIDs) > 0 {
			openWaveBatchCount++
			if previousCurrentWaveBatchStillOpen {
				if batchID == previousCurrentWaveBatchID {
					state = "ready"
					current = true
					readyWaveBatchCount++
					batchSelectionPolicy = "sticky_current_wave_batch"
				} else {
					state = "blocked"
					blockedWaveBatchCount++
				}
			} else if currentWaveBatchID == "" {
				state = "ready"
				current = true
				readyWaveBatchCount++
			} else {
				state = "blocked"
				blockedWaveBatchCount++
			}
		} else {
			resolvedWaveBatchCount++
		}
		summary := fmt.Sprintf("%s should cover %d governed lane(s) in %s round %d.", firstNonBlank(batch.AssignedRole, batch.AssignedRun, "bounded participant"), laneCount, batch.Phase, batch.RoundNumber)
		if current {
			currentWaveBatchID = batchID
			currentWaveBatchLaneCount = laneCount
			currentWaveBatchSummary = summary
		}
		waveBatches = append(waveBatches, map[string]any{
			"batch_id":        batchID,
			"expected_phase":  batch.Phase,
			"round_number":    batch.RoundNumber,
			"assigned_role":   batch.AssignedRole,
			"assigned_run_id": batch.AssignedRun,
			"state":           state,
			"current":         current,
			"action_count":    len(actionIDs),
			"open_action_count": func() int {
				if state == "resolved" {
					return 0
				}
				return len(actionIDs)
			}(),
			"resolved_action_count": func() int {
				if state == "resolved" {
					return len(actionIDs)
				}
				return 0
			}(),
			"action_ids": stringListToAny(actionIDs),
			"claim_refs": stringListToAny(dedupeStrings(batch.ClaimRefs)),
			"lane_count": laneCount,
			"summary":    summary,
		})
	}

	return map[string]any{
		"max_debate_rounds":                        maxDebateRounds,
		"conflict_claim_count":                     conflictClaimCount,
		"planned_route_count":                      len(routes),
		"resolved_route_count":                     resolvedRouteCount,
		"resolved_without_cross_examination_count": resolvedWithoutCrossExaminationCount,
		"authority_resolved_route_count":           authorityResolvedRouteCount,
		"protocol_deferred_route_count":            protocolDeferredRouteCount,
		"protocol_bound_route_count":               protocolBoundedRouteCount,
		"planned_round_count":                      plannedRoundCount,
		"resolved_round_count":                     resolvedRoundCount,
		"unrouted_claim_count":                     len(blockingGaps),
		"open_route_action_count":                  len(routeActions),
		"open_round_action_count":                  len(roundActions),
		"wave_batch_count":                         len(waveBatches),
		"open_wave_batch_count":                    openWaveBatchCount,
		"ready_wave_batch_count":                   readyWaveBatchCount,
		"blocked_wave_batch_count":                 blockedWaveBatchCount,
		"resolved_wave_batch_count":                resolvedWaveBatchCount,
		"batch_selection_policy":                   batchSelectionPolicy,
		"previous_current_wave_batch_id":           previousCurrentWaveBatchID,
		"current_wave_batch_id":                    currentWaveBatchID,
		"current_wave_batch_lane_count":            currentWaveBatchLaneCount,
		"current_wave_batch_summary":               currentWaveBatchSummary,
		"multi_lane_batch_count":                   multiLaneBatchCount,
		"max_wave_lane_count":                      maxWaveLaneCount,
		"routes":                                   routes,
		"route_actions":                            routeActions,
		"round_actions":                            roundActions,
		"wave_batches":                             waveBatches,
		"blocking_gaps":                            stringListToAny(blockingGaps),
		"coverage_notes":                           stringListToAny(coverageNotes),
	}
}

func meetingAutoSynthesize(meeting map[string]any) map[string]any {
	protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
	if protocol == "" {
		protocol = "evidence-review"
	}
	registry := readRegistry()
	debatePlan, _ := toMap(meetingDebatePlan(meeting))
	routeByClaimKey := map[string]map[string]any{}
	for _, item := range toList(debatePlan["routes"]) {
		route, ok := toMap(item)
		if !ok {
			continue
		}
		key := strings.TrimSpace(stringValue(route["claim_key"]))
		if key != "" {
			routeByClaimKey[key] = route
		}
	}
	findings, _ := toMap(meeting["findings"])
	if findings == nil {
		findings = map[string]any{}
	}

	manualClaims := map[string]bool{}
	for _, bucket := range []string{"accepted", "rejected", "unresolved"} {
		cleaned := []any{}
		for _, item := range toList(findings[bucket]) {
			finding, ok := toMap(item)
			if !ok {
				continue
			}
			if stringValue(finding["generated_by"]) == "orchestrate_auto_synthesis" {
				continue
			}
			key := meetingClaimGroupKey(stringValue(finding["claim_ref"]), stringValue(finding["dimension"]), firstNonBlank(stringValue(finding["summary"]), stringValue(finding["finding_id"])))
			manualClaims[key] = true
			cleaned = append(cleaned, finding)
		}
		findings[bucket] = cleaned
	}
	meeting["findings"] = findings

	actions := []any{}
	for _, item := range toList(meeting["action_items"]) {
		action, ok := toMap(item)
		if !ok {
			continue
		}
		if stringValue(action["generated_by"]) == "orchestrate_auto_synthesis" {
			continue
		}
		actions = append(actions, action)
	}
	meeting["action_items"] = actions

	positionGroups := map[string][]map[string]any{}
	claimLabels := map[string]string{}
	for _, item := range toList(meeting["positions"]) {
		position, ok := toMap(item)
		if !ok {
			continue
		}
		claimRef := strings.TrimSpace(stringValue(position["claim_ref"]))
		if claimRef == "" {
			continue
		}
		dimension := firstNonBlank(stringValue(position["dimension"]), "factual")
		key := meetingClaimGroupKey(claimRef, dimension, stringValue(position["position_id"]))
		claimLabels[key] = claimRef
		positionGroups[key] = append(positionGroups[key], position)
	}

	groupKeys := []string{}
	for key := range positionGroups {
		groupKeys = append(groupKeys, key)
	}
	sort.Strings(groupKeys)

	generatedClaims := []string{}
	generatedFindingIDs := []string{}
	generatedActionIDs := []string{}
	generatedSummaries := []string{}
	skippedClaims := []string{}
	claimIndex := 0

	for _, key := range groupKeys {
		if manualClaims[key] {
			skippedClaims = append(skippedClaims, key)
			continue
		}
		positions := positionGroups[key]
		claimRef := claimLabels[key]
		if claimRef == "" {
			claimRef = key
		}
		dimension := "factual"
		if len(positions) > 0 {
			dimension = firstNonBlank(stringValue(positions[0]["dimension"]), "factual")
		}

		supportingIDs := []string{}
		challengingIDs := []string{}
		evidenceRefs := []string{}
		supportRank := 0
		challengeRank := 0
		supportLevel := ""
		challengeLevel := ""
		supportRecordedAt := ""
		challengeRecordedAt := ""
		supportRole := ""
		challengeRole := ""
		supportRunID := ""
		challengeRunID := ""

		for _, position := range positions {
			rank := authorityRank(firstNonBlank(stringValue(position["authority_level"]), "L4"))
			level := authorityLevelForRank(rank)
			evidenceRefs = append(evidenceRefs, anyListToStrings(toList(position["evidence_refs"]))...)
			switch strings.TrimSpace(stringValue(position["stance"])) {
			case "support":
				supportingIDs = append(supportingIDs, stringValue(position["position_id"]))
				if rank > supportRank || (rank == supportRank && isoAfter(stringValue(position["recorded_at"]), firstNonBlank(supportRecordedAt, "1970-01-01T00:00:00Z"))) {
					supportRank = rank
					supportLevel = level
					supportRecordedAt = stringValue(position["recorded_at"])
					supportRole = stringValue(position["role"])
					supportRunID = stringValue(position["run_id"])
				}
			case "challenge":
				challengingIDs = append(challengingIDs, stringValue(position["position_id"]))
				if rank > challengeRank || (rank == challengeRank && isoAfter(stringValue(position["recorded_at"]), firstNonBlank(challengeRecordedAt, "1970-01-01T00:00:00Z"))) {
					challengeRank = rank
					challengeLevel = level
					challengeRecordedAt = stringValue(position["recorded_at"])
					challengeRole = stringValue(position["role"])
					challengeRunID = stringValue(position["run_id"])
				}
			}
		}

		if len(supportingIDs) == 0 && len(challengingIDs) == 0 {
			continue
		}

		claimIndex++
		findingID := fmt.Sprintf("finding-%s-%02d", nowTime().UTC().Format("20060102150405"), claimIndex)
		baseFinding := map[string]any{
			"finding_id":               findingID,
			"summary":                  claimRef,
			"authority_level":          firstNonBlank(supportLevel, challengeLevel, "L4"),
			"evidence_refs":            stringListToAny(dedupeStrings(evidenceRefs)),
			"source_run_id":            "",
			"source_role":              "council-synthesis",
			"claim_ref":                claimRef,
			"dimension":                dimension,
			"supporting_position_ids":  stringListToAny(dedupeStrings(supportingIDs)),
			"challenging_position_ids": stringListToAny(dedupeStrings(challengingIDs)),
			"recorded_at":              nowISO(),
			"generated_by":             "orchestrate_auto_synthesis",
			"generated_claim_key":      key,
		}

		route, _ := toMap(routeByClaimKey[key])
		routeResolution := stringValue(route["resolution_kind"])
		routePrevailingStance := stringValue(route["prevailing_stance"])
		if (routeResolution == "latest_response_prevails" || routeResolution == "higher_authority_prevails") && routePrevailingStance != "" {
			baseFinding["authority_level"] = firstNonBlank(
				map[string]string{
					"support":   supportLevel,
					"challenge": challengeLevel,
				}[routePrevailingStance],
				stringValue(baseFinding["authority_level"]),
				"L4",
			)
			if routePrevailingStance == "support" {
				baseFinding["source_run_id"] = stringValue(route["prevailing_run_id"])
				baseFinding["source_role"] = firstNonBlank(stringValue(route["prevailing_role"]), supportRole, "supporting council participant")
				if routeResolution == "higher_authority_prevails" {
					baseFinding["reason"] = fmt.Sprintf("Auto-synthesized: governed claim lane resolved in favor of support without cross-examination because higher-authority support from %s outranked challenge from %s.", firstNonBlank(stringValue(route["prevailing_role"]), supportRole, "supporting council participant"), firstNonBlank(stringValue(route["overruled_role"]), challengeRole, "challenging council participant"))
				} else {
					baseFinding["reason"] = fmt.Sprintf("Auto-synthesized: governed debate route resolved in favor of support after the latest response from %s prevailed over %s.", firstNonBlank(stringValue(route["prevailing_role"]), supportRole, "supporting council participant"), firstNonBlank(stringValue(route["overruled_role"]), challengeRole, "challenging council participant"))
				}
				findings["accepted"] = append(toList(findings["accepted"]), baseFinding)
			} else {
				baseFinding["source_run_id"] = stringValue(route["prevailing_run_id"])
				baseFinding["source_role"] = firstNonBlank(stringValue(route["prevailing_role"]), challengeRole, "challenging council participant")
				if routeResolution == "higher_authority_prevails" {
					baseFinding["reason"] = fmt.Sprintf("Auto-synthesized: governed claim lane resolved in favor of challenge without cross-examination because higher-authority challenge from %s outranked support from %s.", firstNonBlank(stringValue(route["prevailing_role"]), challengeRole, "challenging council participant"), firstNonBlank(stringValue(route["overruled_role"]), supportRole, "supporting council participant"))
				} else {
					baseFinding["reason"] = fmt.Sprintf("Auto-synthesized: governed debate route resolved in favor of challenge after the latest response from %s prevailed over %s.", firstNonBlank(stringValue(route["prevailing_role"]), challengeRole, "challenging council participant"), firstNonBlank(stringValue(route["overruled_role"]), supportRole, "supporting council participant"))
				}
				findings["rejected"] = append(toList(findings["rejected"]), baseFinding)
			}
			generatedClaims = append(generatedClaims, key)
			generatedFindingIDs = append(generatedFindingIDs, findingID)
			generatedSummaries = append(generatedSummaries, claimRef)
			continue
		}

		switch {
		case routeResolution == "protocol_prefers_direct_synthesis" || routeResolution == "protocol_bound_reached":
			baseFinding["authority_level"] = firstNonBlank(supportLevel, challengeLevel, "L4")
			if routeResolution == "protocol_prefers_direct_synthesis" {
				baseFinding["reason"] = fmt.Sprintf("Auto-synthesized unresolved conflict: %s", protocolDirectSynthesisCoverageNote(protocol, claimRef))
			} else {
				baseFinding["reason"] = protocolBoundedRoundCoverageNote(protocol, claimRef, routeDebateRound(route))
			}
			findings["unresolved"] = append(toList(findings["unresolved"]), baseFinding)
			actionID := fmt.Sprintf("action-%s-%02d", nowTime().UTC().Format("20060102150405"), claimIndex)
			actions = append(actions, map[string]any{
				"action_id":           actionID,
				"summary":             protocolDeferredActionSummary(protocol, claimRef),
				"recorded_at":         nowISO(),
				"state":               "open",
				"generated_by":        "orchestrate_auto_synthesis",
				"generated_claim_key": key,
				"decision_protocol":   protocol,
			})
			generatedActionIDs = append(generatedActionIDs, actionID)
		case len(supportingIDs) > 0 && len(challengingIDs) == 0:
			baseFinding["source_run_id"] = supportRunID
			baseFinding["source_role"] = firstNonBlank(supportRole, "supporting council participant")
			baseFinding["reason"] = fmt.Sprintf("Auto-synthesized from %d supporting council position(s) with no recorded challenge.", len(supportingIDs))
			findings["accepted"] = append(toList(findings["accepted"]), baseFinding)
		case len(challengingIDs) > 0 && len(supportingIDs) == 0:
			baseFinding["source_run_id"] = challengeRunID
			baseFinding["source_role"] = firstNonBlank(challengeRole, "challenging council participant")
			baseFinding["reason"] = fmt.Sprintf("Auto-synthesized from %d challenging council position(s) with no recorded support.", len(challengingIDs))
			findings["rejected"] = append(toList(findings["rejected"]), baseFinding)
		case supportRank > challengeRank:
			baseFinding["authority_level"] = firstNonBlank(supportLevel, stringValue(baseFinding["authority_level"]), "L4")
			baseFinding["source_run_id"] = supportRunID
			baseFinding["source_role"] = firstNonBlank(supportRole, "supporting council participant")
			baseFinding["reason"] = fmt.Sprintf("Auto-synthesized: higher-authority support from %s overrides challenge from %s.", firstNonBlank(supportRole, "supporting council participant"), firstNonBlank(challengeRole, "challenging council participant"))
			findings["accepted"] = append(toList(findings["accepted"]), baseFinding)
		case challengeRank > supportRank:
			baseFinding["authority_level"] = firstNonBlank(challengeLevel, stringValue(baseFinding["authority_level"]), "L4")
			baseFinding["source_run_id"] = challengeRunID
			baseFinding["source_role"] = firstNonBlank(challengeRole, "challenging council participant")
			baseFinding["reason"] = fmt.Sprintf("Auto-synthesized: higher-authority challenge from %s overrides support from %s.", firstNonBlank(challengeRole, "challenging council participant"), firstNonBlank(supportRole, "supporting council participant"))
			findings["rejected"] = append(toList(findings["rejected"]), baseFinding)
		default:
			baseFinding["authority_level"] = firstNonBlank(supportLevel, challengeLevel, "L4")
			baseFinding["reason"] = fmt.Sprintf("Auto-synthesized unresolved conflict: same-authority support and challenge remain on claim %q.", claimRef)
			findings["unresolved"] = append(toList(findings["unresolved"]), baseFinding)
			actionID := fmt.Sprintf("action-%s-%02d", nowTime().UTC().Format("20060102150405"), claimIndex)
			actions = append(actions, map[string]any{
				"action_id":           actionID,
				"summary":             protocolDeferredActionSummary(protocol, claimRef),
				"recorded_at":         nowISO(),
				"state":               "open",
				"generated_by":        "orchestrate_auto_synthesis",
				"generated_claim_key": key,
				"decision_protocol":   protocol,
			})
			generatedActionIDs = append(generatedActionIDs, actionID)
		}

		generatedClaims = append(generatedClaims, key)
		generatedFindingIDs = append(generatedFindingIDs, findingID)
		generatedSummaries = append(generatedSummaries, claimRef)
	}

	meeting["findings"] = findings
	meeting["action_items"] = actions

	filteredLockedDecisions := []any{}
	for _, item := range toList(meeting["locked_decisions"]) {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		if stringValue(decision["generated_by"]) == "orchestrate_auto_synthesis" {
			continue
		}
		filteredLockedDecisions = append(filteredLockedDecisions, decision)
	}
	meeting["locked_decisions"] = filteredLockedDecisions
	autoLockedDecision := meetingAutoLockBundleGovernorDecision(meeting, registry)
	if len(autoLockedDecision) > 0 {
		meeting["locked_decisions"] = append(toList(meeting["locked_decisions"]), autoLockedDecision)
	}
	outputPacket := meetingAutoOutputPacket(meeting, protocol)

	status := "noop"
	if len(generatedFindingIDs) > 0 || len(generatedActionIDs) > 0 || len(outputPacket) > 0 {
		status = "generated"
	}

	return map[string]any{
		"status":                  status,
		"generated_claim_count":   len(generatedClaims),
		"generated_finding_count": len(generatedFindingIDs),
		"generated_action_count":  len(generatedActionIDs),
		"generated_claims":        stringListToAny(generatedClaims),
		"generated_finding_ids":   stringListToAny(generatedFindingIDs),
		"generated_action_ids":    stringListToAny(generatedActionIDs),
		"generated_summaries":     stringListToAny(generatedSummaries),
		"generated_locked_decision_count": func() int {
			if len(autoLockedDecision) == 0 {
				return 0
			}
			return 1
		}(),
		"output_packet":      outputPacket,
		"output_packet_kind": stringValue(outputPacket["output_packet_kind"]),
		"skipped_claims":     stringListToAny(dedupeStrings(skippedClaims)),
		"last_generated_at":  nowISO(),
	}
}

func meetingSynthesisAudit(meeting map[string]any) map[string]any {
	coverage, _ := toMap(meetingCoverageAudit(meeting))
	debatePlan, _ := toMap(meetingDebatePlan(meeting))
	findings, _ := toMap(meeting["findings"])
	accepted := toList(findings["accepted"])
	rejected := toList(findings["rejected"])
	unresolved := toList(findings["unresolved"])
	actions := toList(meeting["action_items"])
	totalFindings := 0
	tracedFindings := 0

	evidenceGaps := []string{}
	for _, bucketItem := range [][]any{accepted, rejected, unresolved} {
		for _, item := range bucketItem {
			finding, ok := toMap(item)
			if !ok {
				continue
			}
			summary := firstNonBlank(stringValue(finding["summary"]), stringValue(finding["finding_id"]))
			if len(toList(finding["evidence_refs"])) == 0 {
				evidenceGaps = append(evidenceGaps, fmt.Sprintf("finding missing evidence refs: %s", summary))
			}
			if strings.TrimSpace(stringValue(finding["authority_level"])) == "" {
				evidenceGaps = append(evidenceGaps, fmt.Sprintf("finding missing authority level: %s", summary))
			}
			if stringValue(finding["finding_id"]) == "" {
				evidenceGaps = append(evidenceGaps, fmt.Sprintf("finding missing finding_id: %s", summary))
			}
			totalFindings++
			if len(toList(finding["supporting_position_ids"])) > 0 || len(toList(finding["challenging_position_ids"])) > 0 {
				tracedFindings++
			} else if len(toList(meeting["positions"])) > 0 {
				evidenceGaps = append(evidenceGaps, fmt.Sprintf("finding missing synthesis trace to meeting positions: %s", summary))
			}
		}
	}
	if len(unresolved) > 0 && len(actions) == 0 {
		evidenceGaps = append(evidenceGaps, "unresolved findings exist without any action item or escalation note")
	}

	blockingGaps := []string{}
	for _, item := range toList(coverage["blocking_gaps"]) {
		text := strings.TrimSpace(stringValue(item))
		if text != "" {
			blockingGaps = append(blockingGaps, text)
		}
	}
	for _, item := range toList(debatePlan["blocking_gaps"]) {
		text := strings.TrimSpace(stringValue(item))
		if text != "" {
			blockingGaps = append(blockingGaps, text)
		}
	}
	coverageNotes := []any{}
	for _, item := range toList(coverage["coverage_notes"]) {
		text := strings.TrimSpace(stringValue(item))
		if text != "" {
			coverageNotes = append(coverageNotes, text)
		}
	}
	for _, item := range toList(debatePlan["coverage_notes"]) {
		text := strings.TrimSpace(stringValue(item))
		if text != "" {
			coverageNotes = append(coverageNotes, text)
		}
	}
	if len(accepted)+len(rejected)+len(unresolved)+len(actions) == 0 {
		blockingGaps = append(blockingGaps, "no synthesized findings or action items recorded yet")
	}
	authorityAudit, _ := toMap(meetingAuthorityAudit(meeting))
	if intValue(authorityAudit["unresolved_conflict_count"]) > 0 {
		blockingGaps = append(blockingGaps, "same-claim findings still have unresolved authority conflicts")
	}

	status := "ready_for_output"
	switch {
	case len(blockingGaps) > 0:
		status = "blocked"
	case len(evidenceGaps) > 0:
		status = "needs_supplement"
	}

	return map[string]any{
		"status": status,
		"tally": map[string]any{
			"accepted_findings":         len(accepted),
			"rejected_findings":         len(rejected),
			"unresolved_findings":       len(unresolved),
			"action_items":              len(actions),
			"traced_findings":           tracedFindings,
			"untraced_findings":         maxInt(totalFindings-tracedFindings, 0),
			"total_positions":           len(toList(meeting["positions"])),
			"phase_history_count":       len(toList(meeting["phase_history"])),
			"participant_count":         coverage["participant_count"],
			"linked_cross_examinations": coverage["linked_cross_examination_count"],
			"planned_routes":            debatePlan["planned_route_count"],
			"resolved_routes":           debatePlan["resolved_route_count"],
			"open_route_actions":        debatePlan["open_route_action_count"],
			"planned_rounds":            debatePlan["planned_round_count"],
			"resolved_rounds":           debatePlan["resolved_round_count"],
			"open_round_actions":        debatePlan["open_round_action_count"],
			"unrouted_claims":           debatePlan["unrouted_claim_count"],
		},
		"blocking_gaps":   stringListToAny(blockingGaps),
		"coverage_notes":  coverageNotes,
		"evidence_gaps":   stringListToAny(evidenceGaps),
		"authority_audit": authorityAudit,
		"generated_at":    nowISO(),
	}
}

func meetingGovernancePolicyGaps(meeting map[string]any) []string {
	gaps := []string{}
	if boolValue(policyValue("meetings.require_authority_hierarchy", true)) && len(toList(meeting["authority_hierarchy"])) == 0 {
		gaps = append(gaps, "missing authority hierarchy")
	}
	if boolValue(policyValue("meetings.require_rejected_findings_log", true)) {
		findings, _ := toMap(meeting["findings"])
		for _, item := range toList(findings["rejected"]) {
			finding, ok := toMap(item)
			if ok && strings.TrimSpace(stringValue(finding["reason"])) == "" {
				gaps = append(gaps, "rejected finding missing reason")
				break
			}
		}
	}
	return dedupeStrings(gaps)
}

func meetingAuthorityAudit(meeting map[string]any) map[string]any {
	findings, _ := toMap(meeting["findings"])
	groups := map[string][]map[string]any{}
	for _, bucket := range []string{"accepted", "rejected", "unresolved"} {
		for _, item := range toList(findings[bucket]) {
			finding, ok := toMap(item)
			if !ok {
				continue
			}
			dimension := firstNonBlank(stringValue(finding["dimension"]), "factual")
			claimRef := firstNonBlank(strings.TrimSpace(stringValue(finding["claim_ref"])), strings.TrimSpace(stringValue(finding["summary"])))
			if claimRef == "" {
				claimRef = stringValue(finding["finding_id"])
			}
			normalizedClaim := slugify(strings.ToLower(claimRef))
			if normalizedClaim == "" {
				normalizedClaim = stringValue(finding["finding_id"])
			}
			groupKey := meetingClaimGroupKey(claimRef, dimension, stringValue(finding["finding_id"]))
			entry := map[string]any{
				"finding_id":      stringValue(finding["finding_id"]),
				"bucket":          bucket,
				"summary":         stringValue(finding["summary"]),
				"claim_ref":       claimRef,
				"dimension":       dimension,
				"authority_level": firstNonBlank(stringValue(finding["authority_level"]), "L4"),
				"authority_rank":  authorityRank(firstNonBlank(stringValue(finding["authority_level"]), "L4")),
				"recorded_at":     stringValue(finding["recorded_at"]),
				"source_role":     stringValue(finding["source_role"]),
				"source_run_id":   stringValue(finding["source_run_id"]),
			}
			groups[groupKey] = append(groups[groupKey], entry)
		}
	}

	conflicts := []any{}
	unresolvedConflictCount := 0
	for _, entries := range groups {
		distinctBuckets := map[string]bool{}
		for _, entry := range entries {
			distinctBuckets[stringValue(entry["bucket"])] = true
		}
		if len(distinctBuckets) <= 1 {
			continue
		}
		sort.Slice(entries, func(i, j int) bool {
			left := intValue(entries[i]["authority_rank"])
			right := intValue(entries[j]["authority_rank"])
			if left != right {
				return left > right
			}
			return isoAfter(stringValue(entries[i]["recorded_at"]), stringValue(entries[j]["recorded_at"]))
		})
		winner := entries[0]
		losers := []any{}
		losingBuckets := map[string]bool{}
		for _, entry := range entries[1:] {
			losers = append(losers, entry)
			losingBuckets[stringValue(entry["bucket"])] = true
		}
		resolutionKind := "authority_wins"
		note := fmt.Sprintf("Higher-authority %s finding from %s overrides lower-authority conflicting findings on the same %s claim.", stringValue(winner["bucket"]), firstNonBlank(stringValue(winner["source_role"]), stringValue(winner["finding_id"])), stringValue(winner["dimension"]))
		if len(entries) > 1 && intValue(entries[1]["authority_rank"]) == intValue(winner["authority_rank"]) && stringValue(entries[1]["bucket"]) != stringValue(winner["bucket"]) {
			resolutionKind = "tie_needs_judgment"
			note = fmt.Sprintf("Same-authority findings disagree on the same %s claim; leader judgment or explicit synthesis note is required.", stringValue(winner["dimension"]))
			unresolvedConflictCount++
		}
		conflicts = append(conflicts, map[string]any{
			"claim_ref":           stringValue(winner["claim_ref"]),
			"dimension":           stringValue(winner["dimension"]),
			"winning_finding":     winner,
			"losing_findings":     losers,
			"losing_bucket_count": len(losingBuckets),
			"resolution_kind":     resolutionKind,
			"note":                note,
		})
	}

	return map[string]any{
		"conflict_count":            len(conflicts),
		"unresolved_conflict_count": unresolvedConflictCount,
		"conflicts":                 conflicts,
	}
}

func emitSignal(kind string, sourceRun map[string]any, title, summary string, appliesTo []string, scopeType, confidence, authorityLevel string, evidenceRefs []string, freshHours float64) (map[string]any, error) {
	if authorityLevel != "" && !validAuthorityLevel(authorityLevel) {
		return nil, fmt.Errorf("invalid authority level: %s", authorityLevel)
	}
	log := readSignalLog()
	templateRoot, _ := toMap(loadData(signalTemplatePath(), map[string]any{}))
	template, _ := toMap(templateRoot["signal"])
	emittedAt := nowISO()
	freshUntil := nowTime().Add(time.Duration(freshHours * float64(time.Hour))).UTC().Format(time.RFC3339Nano)
	merged, _ := toMap(mergeDicts(template, map[string]any{
		"signal_id":       fmt.Sprintf("sig-%s-%s-%s", nowTime().UTC().Format("20060102150405"), slugify(kind), slugify(firstNonBlank(stringValue(sourceRun["owner_id"]), stringValue(sourceRun["run_id"])))),
		"kind":            kind,
		"source_type":     sourceTypeForRun(sourceRun),
		"source_id":       firstNonBlank(stringValue(sourceRun["owner_id"]), stringValue(sourceRun["leader_id"])),
		"run_id":          stringValue(sourceRun["run_id"]),
		"mission_id":      stringValue(sourceRun["mission_id"]),
		"applies_to":      stringListToAny(appliesTo),
		"scope_type":      scopeType,
		"title":           title,
		"summary":         summary,
		"evidence_refs":   stringListToAny(evidenceRefs),
		"authority_level": authorityLevel,
		"confidence":      confidence,
		"emitted_at":      emittedAt,
		"fresh_until":     freshUntil,
		"state":           "active",
	}))
	signals := append(toList(log["signals"]), merged)
	log["signals"] = signals
	if err := writeSignalLog(log); err != nil {
		return nil, err
	}
	return merged, nil
}

func stringListToAny(items []string) []any {
	result := make([]any, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result
}

func materializeRC(rcID, version, branch, baseBranch, productionBranch string, runIDs []string) (map[string]any, error) {
	templateRoot, _ := toMap(loadData(rcTemplatePath(), map[string]any{}))
	template, _ := toMap(templateRoot["rc"])
	timestamp := nowISO()
	rc, _ := toMap(mergeDicts(template, map[string]any{
		"rc_id":               rcID,
		"version":             version,
		"state":               "active",
		"branch":              branch,
		"base_branch":         baseBranch,
		"production_branch":   productionBranch,
		"run_ids":             stringListToAny(runIDs),
		"verification_status": "unknown",
		"verification_refs":   []any{},
		"fix_refs":            []any{},
		"notes":               []any{},
		"created_at":          timestamp,
		"updated_at":          timestamp,
		"promoted_at":         "",
		"closed_at":           "",
	}))
	if err := dumpData(rcPath(rcID, false), map[string]any{"version": 1, "rc": rc}); err != nil {
		return nil, err
	}
	return rc, nil
}

func issueRecord(severity, kind, ref, summary, recommendation string) map[string]any {
	return map[string]any{
		"severity":       severity,
		"kind":           kind,
		"ref":            ref,
		"summary":        summary,
		"recommendation": recommendation,
	}
}

func noteForRC(summary string, extra map[string]any) map[string]any {
	note := map[string]any{"recorded_at": nowISO(), "summary": summary}
	for key, value := range extra {
		if value == nil {
			continue
		}
		if isBlank(value) {
			continue
		}
		note[key] = value
	}
	return note
}

func applyRCRunFollowup(registry map[string]any, runIDs []string, summary, nextAction, rcID string) {
	runs := toList(registry["runs"])
	for _, runID := range runIDs {
		index := findRunIndex(registry, runID)
		if index < 0 {
			continue
		}
		run, _ := toMap(runs[index])
		if rcID != "" {
			run["rc_id"] = rcID
		}
		run["latest_summary"] = summary
		run["next_expected_action"] = nextAction
		run["updated_at"] = nowISO()
	}
	registry["runs"] = runs
}

func ageHours(value string) *float64 {
	parsed, err := parseISO(value)
	if err != nil || parsed == nil {
		return nil
	}
	hours := mathRound((nowTime().Sub(*parsed)).Hours(), 2)
	return &hours
}

func mathRound(value float64, digits int) float64 {
	pow := mathPow10(digits)
	return float64(int(value*pow+0.5)) / pow
}

func mathPow10(power int) float64 {
	value := 1.0
	for i := 0; i < power; i++ {
		value *= 10
	}
	return value
}

func doctorReport() map[string]any {
	_ = ensureRuntimeDirs()
	profile, _ := resolveProfile("")
	adapter := adapterValidate()
	registry := readRegistry()
	missionRegistry := readMissionRegistry()
	signals := readSignalLog()
	meetings := []map[string]any{}
	for _, item := range listMeetings(false) {
		meeting, _ := toMap(item)
		if stringValue(meeting["state"]) == "active" {
			meetings = append(meetings, meeting)
		}
	}
	meetingActions := openMeetingActions()
	activeRCs := []map[string]any{}
	for _, item := range listRCs(false) {
		rc, _ := toMap(item)
		if stringValue(rc["state"]) == "active" {
			activeRCs = append(activeRCs, rc)
		}
	}
	activeRuns := []map[string]any{}
	activeTeammates := []map[string]any{}
	activeMissions := []map[string]any{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if activeRun(run) {
			activeRuns = append(activeRuns, run)
			if activeTeammateRun(run) {
				activeTeammates = append(activeTeammates, run)
			}
		}
	}
	for _, item := range toList(missionRegistry["missions"]) {
		mission, _ := toMap(item)
		if activeMission(mission) {
			activeMissions = append(activeMissions, mission)
		}
	}
	activeSignals := []map[string]any{}
	for _, item := range toList(signals["signals"]) {
		signal, _ := toMap(item)
		if activeSignal(signal) {
			activeSignals = append(activeSignals, signal)
		}
	}
	activeRunIDs := map[string]bool{}
	activeMissionIDs := map[string]bool{}
	for _, run := range activeRuns {
		activeRunIDs[stringValue(run["run_id"])] = true
	}
	for _, mission := range activeMissions {
		activeMissionIDs[stringValue(mission["mission_id"])] = true
	}
	runStaleAfter := floatValue(policyValue("cleanup.stale_run_after_hours", 24))
	signalStaleAfter := floatValue(policyValue("cleanup.stale_signal_after_hours", 12))
	meetingStaleAfter := floatValue(policyValue("cleanup.stale_meeting_after_hours", 6))
	issues := []any{}
	issues = append(issues, toList(adapter["issues"])...)

	for _, run := range activeRuns {
		runID := stringValue(run["run_id"])
		if stringValue(run["created_at"]) == "" {
			issues = append(issues, issueRecord("warning", "run-missing-created-at", runID, "Run is active but missing created_at.", "Backfill timestamps or recreate the run cleanly."))
		}
		reportAge := ageHours(stringValue(run["last_reported_at"]))
		if reportAge == nil {
			issues = append(issues, issueRecord("warning", "run-missing-last-reported-at", runID, "Run is active but missing last_reported_at.", "Have the owner emit a report or repair the run state."))
		} else if *reportAge > runStaleAfter {
			issues = append(issues, issueRecord("warning", "run-stale", runID, fmt.Sprintf("Run has not reported for %.2f hours.", *reportAge), "Review whether the run is still active, blocked, or abandoned."))
		}
		if stringValue(run["state"]) == "blocked" {
			issues = append(issues, issueRecord("warning", "run-blocked", runID, firstNonBlank(stringValue(run["status_reason"]), "Run is blocked."), "Resolve the blocker or convene a meeting."))
		}
		if stringValue(run["pr_ref"]) != "" && stringValue(run["review_status"]) == "changes_requested" {
			issues = append(issues, issueRecord("warning", "pr-changes-requested", runID, "PR has requested changes.", "Address review feedback, push updates, and continue the review loop."))
		}
		if stringValue(run["pr_ref"]) != "" && stringValue(run["checks_status"]) == "failed" {
			issues = append(issues, issueRecord("warning", "pr-checks-failed", runID, "PR checks are failing.", "Fix the failing checks before merging."))
		}
		latestClosedMeeting := governingClosedMeetingForRun(run)
		if latestClosedMeeting != nil && meetingHasOutput(latestClosedMeeting) {
			outcome := meetingOutcome(latestClosedMeeting)
			if stringValue(outcome["summary"]) != "" && strings.TrimSpace(stringValue(run["latest_summary"])) != stringValue(outcome["summary"]) {
				issues = append(issues, issueRecord("warning", "run-missing-meeting-summary", runID, "Run has a closed meeting with accepted output that is not reflected in latest_summary.", "Adopt the meeting outcome into the run state."))
			}
			if expectedNextAction := meetingAdoptedNextAction(latestClosedMeeting); expectedNextAction != "" && strings.TrimSpace(stringValue(run["next_expected_action"])) != expectedNextAction {
				issues = append(issues, issueRecord("warning", "run-missing-meeting-next-action", runID, "Run has a closed meeting with adopted follow-up guidance that is not reflected in next_expected_action.", "Adopt the meeting outcome into the run state."))
			}
		}
		worktreePath := stringValue(run["worktree_path"])
		if worktreePath != "" && !exists(worktreePath) {
			issues = append(issues, issueRecord("warning", "worktree-missing", runID, fmt.Sprintf("Run references missing worktree path: %s", worktreePath), "Re-provision the worktree or update the run record."))
		}
		if stringValue(run["owner_type"]) == "teammate" {
			if stringValue(run["parent_run_id"]) == "" {
				issues = append(issues, issueRecord("warning", "teammate-missing-parent", runID, "Teammate run is missing parent_run_id.", "Repair the delegation link or recreate the teammate assignment."))
			}
			if stringValue(run["assigned_scope"]) == "" {
				issues = append(issues, issueRecord("warning", "teammate-missing-scope", runID, "Teammate run has no assigned scope.", "Record the bounded assignment scope explicitly."))
			}
			if len(toList(run["reporting_expectations"])) == 0 {
				issues = append(issues, issueRecord("warning", "teammate-missing-reporting-contract", runID, "Teammate run has no reporting expectations.", "Record the teammate reporting contract before trusting the assignment."))
			}
		}
		if len(ownershipClaims(run)) == 0 {
			issues = append(issues, issueRecord("warning", "run-missing-ownership-claims", runID, "Run has no explicit ownership claims.", "Record ownership selectors so overlap can be governed explicitly."))
		} else if len(toList(run["ownership_selectors"])) == 0 {
			issues = append(issues, issueRecord("warning", "run-missing-ownership-selectors", runID, "Run has legacy ownership claims but no normalized ownership_selectors.", "Normalize runtime state so selector-aware ownership transfer can work safely."))
		}
		if stringValue(run["mission_id"]) == "" {
			issues = append(issues, issueRecord("critical", "run-missing-mission-id", runID, "Run is missing mission_id.", "Repair the run or recreate it under a mission."))
		} else if !activeMissionIDs[stringValue(run["mission_id"])] && findMissionIndex(missionRegistry, stringValue(run["mission_id"])) < 0 {
			issues = append(issues, issueRecord("critical", "run-mission-missing", runID, "Run references a mission that is missing from mission registry.", "Rebuild mission state from the run registry before continuing."))
		}
		unsatisfiedDeps := unsatisfiedDependencies(run, registry)
		if len(unsatisfiedDeps) > 0 {
			if stringValue(run["state"]) != "waiting" {
				issues = append(issues, issueRecord("warning", "run-dependencies-unsatisfied", runID, fmt.Sprintf("Run depends on unsatisfied runs: %s", strings.Join(unsatisfiedDeps, ", ")), "Move the run to waiting, or satisfy its dependencies first."))
			}
		}
		if stringValue(run["state"]) == "completed" {
			if satisfied, missing := runDesiredStateSatisfied(run); !satisfied {
				issues = append(issues, issueRecord("warning", "run-desired-state-unmet", runID, fmt.Sprintf("Completed run does not satisfy desired state: %s", strings.Join(missing, ", ")), "Either continue the lifecycle or adjust the run's declared desired state."))
			}
		}
		authorityAudit, _ := toMap(runAuthorityAudit(run, registry, signals))
		if intValue(authorityAudit["contract_gap_count"]) > 0 {
			issues = append(issues, issueRecord("warning", "run-contract-drift", runID, "Run control-plane state has drifted from its contract snapshot.", "Run runs authority-audit and refresh the contract so planning truth matches the active run."))
		}
		if intValue(authorityAudit["live_state_gap_count"]) > 0 {
			issues = append(issues, issueRecord("warning", "run-live-state-drift", runID, "Live worktree state disagrees with the run record.", "Reconcile the run record with live git state before trusting the control plane."))
		}
		if intValue(authorityAudit["higher_authority_conflict_count"]) > 0 {
			conflictRefs := []string{}
			for _, item := range toList(authorityAudit["higher_authority_conflicts"]) {
				conflict, ok := toMap(item)
				if !ok {
					continue
				}
				conflictRefs = append(conflictRefs, firstNonBlank(stringValue(conflict["claim_ref"]), stringValue(conflict["signal_id"])))
			}
			issues = append(issues, issueRecord("warning", "run-authority-conflict", runID, fmt.Sprintf("Run has unresolved higher-authority conflicts across active signals: %s", strings.Join(dedupeStrings(conflictRefs), ", ")), "Run runs authority-audit and either adopt or resolve the conflicting higher-authority signal."))
		}
	}

	runIDsByMission := map[string][]string{}
	activeRunIDsByMission := map[string][]string{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		missionID := stringValue(run["mission_id"])
		if missionID == "" {
			continue
		}
		runIDsByMission[missionID] = append(runIDsByMission[missionID], stringValue(run["run_id"]))
		if activeRun(run) {
			activeRunIDsByMission[missionID] = append(activeRunIDsByMission[missionID], stringValue(run["run_id"]))
		}
	}
	for _, mission := range toMapSlice(toList(missionRegistry["missions"])) {
		missionID := stringValue(mission["mission_id"])
		if missionID == "" {
			issues = append(issues, issueRecord("warning", "mission-missing-id", "mission-registry", "Mission entry is missing mission_id.", "Repair or regenerate mission registry."))
			continue
		}
		if len(runIDsByMission[missionID]) == 0 {
			issues = append(issues, issueRecord("warning", "mission-without-runs", missionID, "Mission exists without any corresponding runs.", "Reconcile mission registry or recreate a run for this mission."))
		}
		activeIDs := normalizedStringSet(activeRunIDsByMission[missionID])
		recordedActive := normalizedStringSet(anyListToStrings(toList(mission["active_run_ids"])))
		if strings.Join(activeIDs, "|") != strings.Join(recordedActive, "|") {
			issues = append(issues, issueRecord("warning", "mission-active-run-mismatch", missionID, "Mission active_run_ids do not match the active runs in run registry.", "Regenerate mission state from runs before proceeding."))
		}
		if latestRunID := stringValue(mission["latest_run_id"]); latestRunID == "" {
			issues = append(issues, issueRecord("warning", "mission-missing-latest-run", missionID, "Mission has no latest_run_id.", "Repair mission linkage so leader surfaces can reconcile correctly."))
		} else if findRunIndex(registry, latestRunID) < 0 {
			issues = append(issues, issueRecord("warning", "mission-latest-run-missing", missionID, "Mission latest_run_id does not exist in run registry.", "Regenerate mission state from runs before proceeding."))
		}
		if len(toList(mission["ownership_claims"])) == 0 {
			issues = append(issues, issueRecord("warning", "mission-missing-ownership-claims", missionID, "Mission has no explicit ownership claims.", "Carry ownership selectors up to the mission control plane."))
		} else if len(toList(mission["ownership_selectors"])) == 0 {
			issues = append(issues, issueRecord("warning", "mission-missing-ownership-selectors", missionID, "Mission has legacy ownership claims but no normalized ownership_selectors.", "Regenerate mission state so selector-aware governance can work safely."))
		}
		if activeMission(mission) && len(toList(mission["depends_on"])) > 0 && stringValue(mission["state"]) != "blocked" && stringValue(mission["state"]) != "claimed" && stringValue(mission["state"]) != "wip" && stringValue(mission["state"]) != "review" {
			issues = append(issues, issueRecord("info", "mission-dependency-state-unusual", missionID, "Mission has dependencies but is not in an active dependency-aware state.", "Verify the mission control-plane state is still correct."))
		}
		if stringValue(mission["state"]) == "done" && len(recordedActive) > 0 {
			issues = append(issues, issueRecord("warning", "mission-done-with-active-runs", missionID, "Mission is marked done but still has active runs.", "Close or correct the active runs before trusting the mission state."))
		}
		if desired, ok := toMap(mission["desired_state"]); !ok || stringValue(desired["goal_state"]) == "" {
			issues = append(issues, issueRecord("warning", "mission-missing-desired-state", missionID, "Mission is missing desired_state.goal_state.", "Record mission desired state explicitly for control-plane reconciliation."))
		}
		authorityAudit, _ := toMap(missionAuthorityAudit(mission, registry))
		if intValue(authorityAudit["mission_gap_count"]) > 0 {
			issues = append(issues, issueRecord("warning", "mission-control-plane-drift", missionID, "Mission record has drifted from the authoritative run-derived control plane.", "Run missions authority-audit and regenerate mission state before trusting leader surfaces."))
		}
		if intValue(authorityAudit["unresolved_conflict_count"]) > 0 {
			conflictRefs := []string{}
			for _, item := range toList(authorityAudit["unresolved_conflicts"]) {
				conflict, ok := toMap(item)
				if !ok {
					continue
				}
				conflictRefs = append(conflictRefs, firstNonBlank(stringValue(conflict["claim_ref"]), stringValue(conflict["losing_source"])))
			}
			issues = append(issues, issueRecord("warning", "mission-authority-conflict", missionID, fmt.Sprintf("Mission has unresolved multi-run authority conflicts: %s", strings.Join(dedupeStrings(conflictRefs), ", ")), "Close or merge the competing root runs so the mission control plane has a single governing run."))
		}
	}
	graphAuthorityAudit, _ := toMap(missionGraphAuthorityAudit(missionRegistry, registry, false))
	for _, item := range invalidExplicitBundleGovernorSignals(missionRegistry, registry, signals) {
		invalid, _ := toMap(item)
		issues = append(issues, issueRecord(
			"warning",
			"invalid-explicit-bundle-governor",
			stringValue(invalid["signal_id"]),
			stringValue(invalid["reason"]),
			firstNonBlank(stringValue(invalid["recommended_action"]), "Clear or repair the stale explicit governor signal before trusting bundle governance."),
		))
	}
	if intValue(graphAuthorityAudit["cross_mission_conflict_count"]) > 0 {
		for _, item := range toList(graphAuthorityAudit["cross_mission_conflicts"]) {
			conflict, ok := toMap(item)
			if !ok {
				continue
			}
			issues = append(issues, issueRecord(
				"warning",
				"mission-graph-authority-conflict",
				strings.Join(anyListToStrings(toList(conflict["mission_ids"])), "::"),
				"Active missions overlap across bounded ownership selectors.",
				"Run missions authority-audit --all, then narrow selectors, serialize the work, or transfer ownership boundaries explicitly.",
			))
		}
	}

	exclusiveRuns := []map[string]any{}
	for _, run := range activeRuns {
		if overlapPolicy(run) != "shared" {
			exclusiveRuns = append(exclusiveRuns, run)
		}
	}
	seenOwnershipOverlap := map[string]bool{}
	for index, leftRun := range exclusiveRuns {
		leftID := stringValue(leftRun["run_id"])
		for _, rightRun := range exclusiveRuns[index+1:] {
			rightID := stringValue(rightRun["run_id"])
			shared := []string{}
			leftExclusions := ownershipExclusions(leftRun)
			rightExclusions := ownershipExclusions(rightRun)
			for _, leftClaim := range ownershipClaims(leftRun) {
				for _, rightClaim := range ownershipClaims(rightRun) {
					if selectorsConflict(leftClaim, leftExclusions, rightClaim, rightExclusions) {
						shared = append(shared, leftClaim+"~"+rightClaim)
					}
				}
			}
			if len(shared) == 0 {
				continue
			}
			key := leftID + "::" + rightID + "::" + strings.Join(dedupeStrings(shared), "|")
			if seenOwnershipOverlap[key] {
				continue
			}
			seenOwnershipOverlap[key] = true
			issues = append(issues, issueRecord("critical", "ownership-overlap", leftID, fmt.Sprintf("Active exclusive runs overlap semantically with %s via %s", rightID, strings.Join(dedupeStrings(shared), ", ")), "Narrow ownership selectors, transfer them formally, or explicitly mark overlap as shared."))
		}
	}

	for _, signal := range activeSignals {
		signalID := stringValue(signal["signal_id"])
		runID := stringValue(signal["run_id"])
		if runID != "" && !activeRunIDs[runID] {
			issues = append(issues, issueRecord("critical", "signal-orphaned", signalID, "Active signal references a missing run.", "Run maintenance reconcile before proceeding."))
		}
		if stringValue(signal["emitted_at"]) == "" {
			issues = append(issues, issueRecord("warning", "signal-missing-emitted-at", signalID, "Signal is active but missing emitted_at.", "Repair signal metadata or re-emit the signal."))
		}
		if emittedAge := ageHours(stringValue(signal["emitted_at"])); emittedAge != nil && *emittedAge > signalStaleAfter {
			issues = append(issues, issueRecord("warning", "signal-stale", signalID, fmt.Sprintf("Signal has been active for %.2f hours.", *emittedAge), "Resolve or refresh the signal if it still matters."))
		}
		if freshUntil, err := parseISO(stringValue(signal["fresh_until"])); err == nil && freshUntil != nil && freshUntil.Before(nowTime()) {
			issues = append(issues, issueRecord("warning", "signal-expired", signalID, "Signal freshness window has expired.", "Resolve or refresh the signal."))
		}
		if level := stringValue(signal["authority_level"]); level != "" && !validAuthorityLevel(level) {
			issues = append(issues, issueRecord("warning", "signal-invalid-authority", signalID, "Signal uses an invalid authority level.", "Use one of L1-L5."))
		}
		if superseded, resolution := signalSupersession(signal, registry); superseded {
			issues = append(issues, issueRecord("warning", "signal-authority-superseded", signalID, firstNonBlank(stringValue(resolution["resolution_note"]), "Active signal is superseded by fresher higher-authority state."), "Run maintenance reconcile or signals adjudicate to resolve the superseded signal."))
		}
	}

	for _, meeting := range meetings {
		meetingID := stringValue(meeting["meeting_id"])
		runID := stringValue(meeting["run_id"])
		ensureMeetingGovernanceState(meeting)
		if runID != "" && !activeRunIDs[runID] {
			issues = append(issues, issueRecord("critical", "meeting-orphaned", meetingID, "Active meeting references a missing run.", "Run maintenance reconcile before proceeding."))
		}
		if len(toList(meeting["authority_hierarchy"])) == 0 {
			issues = append(issues, issueRecord("warning", "meeting-missing-authority", meetingID, "Meeting has no authority hierarchy.", "Repair the meeting packet before using its conclusions."))
		}
		if len(toList(meeting["roles"])) == 0 {
			issues = append(issues, issueRecord("warning", "meeting-missing-roles", meetingID, "Meeting has no assigned roles.", "Assign bounded roles before continuing the meeting."))
		}
		if len(toList(meeting["role_specs"])) == 0 {
			issues = append(issues, issueRecord("warning", "meeting-missing-role-specs", meetingID, "Meeting has roles but no role_specs with allowed/forbidden source boundaries.", "Record role specs before trusting council output."))
		}
		if stringValue(meeting["phase"]) == "" {
			issues = append(issues, issueRecord("warning", "meeting-missing-phase", meetingID, "Meeting has no active phase marker.", "Set the meeting phase before continuing the council."))
		}
		if len(toList(meeting["phase_history"])) == 0 {
			issues = append(issues, issueRecord("warning", "meeting-missing-phase-history", meetingID, "Meeting has no phase history trace.", "Re-orchestrate the meeting so the council history is durable."))
		}
		packet, _ := toMap(meeting["packet"])
		if strings.TrimSpace(stringValue(packet["expected_output"])) == "" {
			issues = append(issues, issueRecord("warning", "meeting-missing-expected-output", meetingID, "Meeting packet has no expected_output.", "Define the expected council output before continuing."))
		}
		if len(toList(packet["context_refs"])) == 0 && len(toList(packet["known_context"])) == 0 {
			issues = append(issues, issueRecord("warning", "meeting-missing-context", meetingID, "Meeting packet has neither context_refs nor known_context.", "Record bounded context before continuing the council."))
		}
		if phase := stringValue(meeting["phase"]); phase == "cross_examination" || phase == "synthesis" || phase == "output" {
			if len(toList(meeting["locked_decisions"])) == 0 {
				issues = append(issues, issueRecord("warning", "meeting-missing-locked-decisions", meetingID, "Meeting reached an advanced phase without any locked decisions.", "Record locked decisions or explicitly note that none exist."))
			}
		}
		for _, item := range toList(meeting["role_specs"]) {
			spec, ok := toMap(item)
			if !ok {
				continue
			}
			role := firstNonBlank(stringValue(spec["role"]), meetingID)
			if len(toList(spec["allowed_sources"])) == 0 {
				issues = append(issues, issueRecord("warning", "meeting-role-missing-allowed-sources", role, "Meeting role spec has no allowed_sources.", "Record explicit allowed sources for the role."))
			}
			if len(toList(spec["forbidden_sources"])) == 0 {
				issues = append(issues, issueRecord("warning", "meeting-role-missing-forbidden-sources", role, "Meeting role spec has no forbidden_sources.", "Record explicit forbidden sources for the role."))
			}
		}
		findings, _ := toMap(meeting["findings"])
		for _, item := range toList(findings["rejected"]) {
			finding, ok := toMap(item)
			if ok && strings.TrimSpace(stringValue(finding["reason"])) == "" {
				issues = append(issues, issueRecord("warning", "meeting-rejected-finding-missing-reason", meetingID, "Rejected meeting finding is missing a reason.", "Record why the council rejected the finding."))
			}
			if ok && stringValue(finding["authority_level"]) != "" && !validAuthorityLevel(stringValue(finding["authority_level"])) {
				issues = append(issues, issueRecord("warning", "meeting-finding-invalid-authority", meetingID, "Meeting finding uses an invalid authority level.", "Use one of L1-L5."))
			}
		}
		coverageAudit, _ := toMap(meetingCoverageAudit(meeting))
		debatePlan, _ := toMap(meetingDebatePlan(meeting))
		authorityAudit, _ := toMap(meetingAuthorityAudit(meeting))
		orchestration, _ := toMap(meeting["orchestration"])
		synthesis, _ := toMap(orchestration["synthesis"])
		if phase := stringValue(meeting["phase"]); phase == "cross_examination" || phase == "synthesis" || phase == "output" {
			if len(toList(coverageAudit["blocking_gaps"])) > 0 {
				issues = append(issues, issueRecord("warning", "meeting-debate-blocked", meetingID, fmt.Sprintf("Meeting still has debate blocking gaps: %s", strings.Join(anyListToStrings(toList(coverageAudit["blocking_gaps"])), "; ")), "Fill the debate gaps or lower the meeting phase."))
			}
			if wavePlan, ok := toMap(orchestration["wave_plan"]); ok && intValue(wavePlan["open_wave_action_count"]) > 0 {
				openWaves := []string{}
				for _, item := range toList(wavePlan["wave_actions"]) {
					action, ok := toMap(item)
					if !ok || stringValue(action["state"]) != "" && stringValue(action["state"]) != "open" {
						continue
					}
					openWaves = append(openWaves, firstNonBlank(stringValue(action["summary"]), stringValue(action["assigned_role"])))
				}
				if len(openWaves) > 0 {
					issues = append(issues, issueRecord("warning", "meeting-wave-action-open", meetingID, fmt.Sprintf("Meeting still has open debate-wave obligations: %s", strings.Join(openWaves, "; ")), "Record the missing independent-review contribution before trusting the council progression."))
				}
			}
			if len(toList(debatePlan["blocking_gaps"])) > 0 {
				issues = append(issues, issueRecord("warning", "meeting-debate-unrouted", meetingID, fmt.Sprintf("Meeting still has unrouted conflicting claims: %s", strings.Join(anyListToStrings(toList(debatePlan["blocking_gaps"])), "; ")), "Add linked cross-examination for each conflicting claim before trusting synthesis."))
			}
			if intValue(debatePlan["open_route_action_count"]) > 0 {
				openRoutes := []string{}
				for _, item := range toList(debatePlan["route_actions"]) {
					action, ok := toMap(item)
					if !ok || stringValue(action["state"]) != "" && stringValue(action["state"]) != "open" {
						continue
					}
					openRoutes = append(openRoutes, firstNonBlank(stringValue(action["summary"]), stringValue(action["claim_ref"])))
				}
				if len(openRoutes) > 0 {
					issues = append(issues, issueRecord("warning", "meeting-route-action-open", meetingID, fmt.Sprintf("Meeting still has open debate route obligations: %s", strings.Join(openRoutes, "; ")), "Record the expected linked cross-examination or explicitly narrow the conflicting claim lane."))
				}
			}
			if intValue(debatePlan["open_round_action_count"]) > 0 {
				openRounds := []string{}
				for _, item := range toList(debatePlan["round_actions"]) {
					action, ok := toMap(item)
					if !ok || stringValue(action["state"]) != "open" {
						continue
					}
					openRounds = append(openRounds, stringValue(action["summary"]))
				}
				if len(openRounds) > 0 {
					issues = append(issues, issueRecord("warning", "meeting-round-action-open", meetingID, fmt.Sprintf("Meeting still has open follow-up debate obligations: %s", strings.Join(openRounds, "; ")), "Record the next-round response before trusting synthesis on that claim lane."))
				}
			}
			openBatchCount := intValue(debatePlan["open_wave_batch_count"])
			readyBatchCount := intValue(debatePlan["ready_wave_batch_count"])
			currentBatchID := strings.TrimSpace(stringValue(debatePlan["current_wave_batch_id"]))
			if (intValue(debatePlan["open_route_action_count"]) > 0 || intValue(debatePlan["open_round_action_count"]) > 0) && openBatchCount == 0 {
				issues = append(issues, issueRecord("warning", "meeting-batch-plan-missing", meetingID, "Meeting has open debate obligations but no active wave batch plan.", "Re-run meetings orchestrate and repair debate batching before trusting council control flow."))
			}
			if openBatchCount > 0 && currentBatchID == "" {
				issues = append(issues, issueRecord("warning", "meeting-current-batch-missing", meetingID, "Meeting has open debate batches but no deterministic current batch.", "Repair debate batch selection so one active batch is chosen before execution continues."))
			}
			if openBatchCount > 0 && readyBatchCount != 1 {
				issues = append(issues, issueRecord("warning", "meeting-batch-selection-ambiguous", meetingID, fmt.Sprintf("Meeting has %d open debate batch(es) but %d ready batch(es).", openBatchCount, readyBatchCount), "Ensure debate batching exposes exactly one current ready batch and leaves the rest blocked."))
			}
		}
		if len(toList(meeting["positions"])) > 0 && len(toList(coverageAudit["blocking_gaps"])) > 0 && intValue(coverageAudit["linked_cross_examination_count"]) == 0 && stringValue(meeting["phase"]) != "independent_review" {
			issues = append(issues, issueRecord("warning", "meeting-cross-exam-unlinked", meetingID, "Cross-examination positions do not explicitly respond to another role or position.", "Record linked cross-examination positions with responds-to metadata."))
		}
		if phase := stringValue(meeting["phase"]); phase == "synthesis" || phase == "output" {
			if synthesis == nil || strings.TrimSpace(stringValue(synthesis["status"])) == "" {
				issues = append(issues, issueRecord("warning", "meeting-synthesis-missing-audit", meetingID, "Meeting reached synthesis/output without a synthesis audit trace.", "Run meetings orchestrate so synthesis state is materialized."))
			} else {
				if len(toList(synthesis["blocking_gaps"])) > 0 {
					issues = append(issues, issueRecord("warning", "meeting-synthesis-blocked", meetingID, fmt.Sprintf("Synthesis still has blocking gaps: %s", strings.Join(anyListToStrings(toList(synthesis["blocking_gaps"])), "; ")), "Resolve synthesis blockers before treating the meeting as settled."))
				}
				if stringValue(meeting["phase"]) == "output" && stringValue(synthesis["status"]) != "ready_for_output" {
					issues = append(issues, issueRecord("warning", "meeting-output-with-incomplete-synthesis", meetingID, "Meeting is in output phase but synthesis is not ready_for_output.", "Supplement evidence or re-open synthesis before final adoption."))
				}
			}
			if intValue(debatePlan["open_wave_batch_count"]) > 0 {
				issues = append(issues, issueRecord("warning", "meeting-output-with-open-batch", meetingID, "Meeting reached synthesis/output while an active debate batch still remains open.", "Finish the current debate batch before treating the council as ready for adoption."))
			}
		}
		if len(toList(findings["unresolved"])) > 0 && len(toList(meeting["action_items"])) == 0 {
			issues = append(issues, issueRecord("warning", "meeting-unresolved-without-action", meetingID, "Meeting has unresolved findings but no action items.", "Record a follow-up action or explicit escalation for unresolved questions."))
		}
		if intValue(authorityAudit["unresolved_conflict_count"]) > 0 {
			issues = append(issues, issueRecord("warning", "meeting-unresolved-authority-conflict", meetingID, "Same-claim findings still disagree at the same authority level.", "Resolve the conflict in synthesis or record an explicit leader judgment."))
		}
		if meetingAge := ageHours(stringValue(meeting["created_at"])); meetingAge != nil {
			if *meetingAge > 0.25 && len(toList(meeting["participants"])) == 0 {
				issues = append(issues, issueRecord("warning", "meeting-missing-participants", meetingID, "Meeting is active but has no participant roster yet.", "Add bounded participants before continuing independent review."))
			}
			if *meetingAge > 0.5 && len(toList(meeting["participants"])) > 0 && len(toList(meeting["positions"])) == 0 {
				issues = append(issues, issueRecord("warning", "meeting-missing-positions", meetingID, "Meeting has participants but no recorded positions yet.", "Record independent review or cross-examination positions before synthesis."))
			}
			if *meetingAge > meetingStaleAfter {
				issues = append(issues, issueRecord("warning", "meeting-stale", meetingID, fmt.Sprintf("Meeting has been open for %.2f hours.", *meetingAge), "Close the meeting or refresh its packet and participants."))
			}
			if *meetingAge > 0.5 && !meetingHasOutput(meeting) {
				issues = append(issues, issueRecord("warning", "meeting-no-output", meetingID, "Meeting is active but has not recorded findings or action items yet.", "Record council findings or close the meeting if it is no longer needed."))
			}
		}
	}

	if len(activeRCs) > 1 {
		issues = append(issues, issueRecord("critical", "rc-multiple-active", "rc-control-plane", "More than one active RC exists.", "Close or promote surplus RCs before proceeding."))
	}

	for _, rc := range activeRCs {
		rcID := stringValue(rc["rc_id"])
		if len(toList(rc["run_ids"])) == 0 {
			issues = append(issues, issueRecord("warning", "rc-missing-runs", rcID, "Active RC has no associated shipped runs.", "Attach shipped runs or close the RC."))
		}
		if stringValue(rc["verification_status"]) != "passed" {
			issues = append(issues, issueRecord("warning", "rc-unverified", rcID, "Active RC has not been marked verified.", "Record RC verification before promotion."))
		}
	}

	sortIssueList(issues)
	critical := countSeverity(issues, "critical")
	warning := countSeverity(issues, "warning")
	info := countSeverity(issues, "info")
	nextAction := "Runtime state looks healthy."
	if critical > 0 {
		nextAction = "Repair critical runtime drift before new execution."
	} else if warning > 0 {
		nextAction = "Triage warnings, then continue execution."
	}

	return map[string]any{
		"profile_id":            stringValue(getKey(profile, "profile.profile_id", "")),
		"adapter":               adapter,
		"issue_count":           len(issues),
		"critical_count":        critical,
		"warning_count":         warning,
		"info_count":            info,
		"active_mission_count":  len(activeMissions),
		"active_teammate_count": len(activeTeammates),
		"issues":                issues,
		"next_action":           nextAction,
		"open_meeting_actions":  len(meetingActions),
	}
}

func countSeverity(issues []any, severity string) int {
	count := 0
	for _, item := range issues {
		mapping, _ := toMap(item)
		if stringValue(mapping["severity"]) == severity {
			count++
		}
	}
	return count
}

func cmdPaths(args []string) error {
	if len(args) < 1 {
		return errors.New("paths requires a subcommand")
	}
	switch args[0] {
	case "repo-root":
		fmt.Println(repoRoot())
	case "definitions-root":
		fmt.Println(definitionsRoot())
	case "runtime-root":
		fmt.Println(runtimeRoot())
	case "ensure-runtime-dirs":
		if err := ensureRuntimeDirs(); err != nil {
			return err
		}
		fmt.Println(runtimeRoot())
	case "profile-path":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		_, basePath := resolveProfile(name)
		fmt.Println(basePath)
	default:
		return fmt.Errorf("unknown paths subcommand: %s", args[0])
	}
	return nil
}

func cmdProfile(args []string) error {
	if len(args) < 1 {
		return errors.New("profile requires a subcommand")
	}
	name := ""
	if args[0] == "path" && len(args) > 1 {
		name = args[1]
	}
	profile, basePath := resolveProfile(name)
	switch args[0] {
	case "show":
		return printJSON(profile)
	case "path":
		fmt.Println(basePath)
	case "id":
		fmt.Println(stringValue(getKey(profile, "profile.profile_id", "")))
	case "get":
		if len(args) < 2 {
			return errors.New("profile get requires KEY")
		}
		key := args[1]
		defaultValue := ""
		if len(args) > 2 {
			defaultValue = args[2]
		}
		if !strings.HasPrefix(key, "profile.") {
			key = "profile." + key
		}
		value := getKey(profile, key, defaultValue)
		switch value.(type) {
		case map[string]any, []any:
			return printJSON(value)
		default:
			fmt.Println(stringValue(value))
		}
	default:
		return fmt.Errorf("unknown profile subcommand: %s", args[0])
	}
	return nil
}

func cmdAdapter(args []string) error {
	if len(args) < 1 {
		return errors.New("adapter requires a subcommand")
	}
	switch args[0] {
	case "show":
		return printJSON(adapterManifest())
	case "version":
		fmt.Println(stringValue(adapterConfig()["adapter_version"]))
	case "id":
		fmt.Println(stringValue(adapterConfig()["adapter_id"]))
	case "package-files":
		for _, relPath := range adapterPackageFiles() {
			fmt.Println(relPath)
		}
	case "smoke-scripts":
		for _, relPath := range adapterSmokeScripts() {
			fmt.Println(relPath)
		}
	case "profile-template":
		fmt.Println(stringValue(adapterBootstrapConfig()["profile_template"]))
	case "generated-profile-path":
		fmt.Println(stringValue(adapterBootstrapConfig()["generated_profile_path"]))
	case "metadata-path":
		fmt.Println(stringValue(adapterBootstrapConfig()["metadata_path"]))
	case "public-contract-path":
		fmt.Println(stringValue(adapterConfig()["public_contract_path"]))
	case "release-criteria-path":
		fmt.Println(stringValue(adapterConfig()["release_criteria_path"]))
	case "validate":
		return printJSON(adapterValidate())
	default:
		return fmt.Errorf("unknown adapter subcommand: %s", args[0])
	}
	return nil
}

func cmdMissions(args []string) error {
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("missions requires a subcommand")
	}
	registry := readMissionRegistry()
	switch args[0] {
	case "path":
		fmt.Println(missionRegistryPath())
	case "list":
		includeAll := false
		for _, arg := range args[1:] {
			if arg == "--all" {
				includeAll = true
			}
		}
		missions := []any{}
		for _, item := range listMissions() {
			mission, _ := toMap(item)
			if includeAll || activeMission(mission) {
				missions = append(missions, mission)
			}
		}
		return printJSON(missions)
	case "get":
		if len(args) < 2 {
			return errors.New("missions get requires MISSION_ID")
		}
		missionID := strings.TrimSpace(args[1])
		if missionID == "" {
			return errors.New("missions get requires MISSION_ID")
		}
		if mission := readMission(missionID); mission != nil {
			return printJSON(mission)
		}
		index := findMissionIndex(registry, missionID)
		if index < 0 {
			return fmt.Errorf("mission not found: %s", missionID)
		}
		mission, _ := toMap(toList(registry["missions"])[index])
		return printJSON(mission)
	case "graph":
		graph := []any{}
		for _, item := range listMissions() {
			mission, _ := toMap(item)
			graph = append(graph, map[string]any{
				"mission_id":           stringValue(mission["mission_id"]),
				"state":                stringValue(mission["state"]),
				"depends_on":           cloneValue(toList(mission["depends_on"])),
				"ownership_claims":     cloneValue(toList(mission["ownership_claims"])),
				"active_run_ids":       cloneValue(toList(mission["active_run_ids"])),
				"latest_run_id":        stringValue(mission["latest_run_id"]),
				"authoritative_run_id": stringValue(mission["authoritative_run_id"]),
				"next_expected_action": stringValue(mission["next_expected_action"]),
			})
		}
		return printJSON(map[string]any{
			"mission_count": len(graph),
			"missions":      graph,
		})
	case "schedule":
		includeAll := false
		for _, arg := range args[1:] {
			if arg == "--all" {
				includeAll = true
			}
		}
		return withHarnessLock(func() error {
			profile, _ := resolveProfile("")
			snapshot := schedulerSummary(profile, readMissionRegistry(), readRegistry(), readSignalLog(), includeAll)
			written, err := writeSchedulerSnapshot(snapshot)
			if err != nil {
				return err
			}
			return printJSON(written)
		})
	case "authority-audit":
		if len(args) >= 2 && args[1] == "--all" {
			return printJSON(missionGraphAuthorityAudit(registry, readRegistry(), false))
		}
		if len(args) < 2 {
			return errors.New("missions authority-audit requires MISSION_ID or --all")
		}
		missionID := strings.TrimSpace(args[1])
		if missionID == "" {
			return errors.New("missions authority-audit requires MISSION_ID or --all")
		}
		mission := readMission(missionID)
		if mission == nil {
			index := findMissionIndex(registry, missionID)
			if index < 0 {
				return fmt.Errorf("mission not found: %s", missionID)
			}
			mission, _ = toMap(toList(registry["missions"])[index])
		}
		return printJSON(missionAuthorityAudit(mission, readRegistry()))
	case "establish-governor":
		if len(args) < 3 {
			return errors.New("missions establish-governor requires BUNDLE_ID GOVERNING_MISSION_ID [NOTE]")
		}
		bundleID := strings.TrimSpace(args[1])
		governingMissionID := strings.TrimSpace(args[2])
		note := strings.TrimSpace(strings.Join(args[3:], " "))
		if bundleID == "" || governingMissionID == "" {
			return errors.New("missions establish-governor requires BUNDLE_ID GOVERNING_MISSION_ID [NOTE]")
		}
		return withHarnessLock(func() error {
			missionRegistry := readMissionRegistry()
			runtimeRegistry := readRegistry()
			graph, _ := toMap(missionGraphAuthorityAudit(missionRegistry, runtimeRegistry, false))
			var bundle map[string]any
			for _, item := range toList(graph["bundles"]) {
				candidate, _ := toMap(item)
				if stringValue(candidate["bundle_id"]) == bundleID {
					bundle = candidate
					break
				}
			}
			if len(bundle) == 0 {
				return fmt.Errorf("bundle not found: %s", bundleID)
			}
			if !containsString(anyListToStrings(toList(bundle["mission_ids"])), governingMissionID) {
				return fmt.Errorf("mission %s does not belong to bundle %s", governingMissionID, bundleID)
			}
			log := readSignalLog()
			signalTemplateRoot, _ := toMap(loadData(signalTemplatePath(), map[string]any{}))
			signalTemplate, _ := toMap(signalTemplateRoot["signal"])
			signalID := fmt.Sprintf("sig-%s-bundle-governor-%s", nowTime().UTC().Format("20060102150405"), slugify(bundleID))
			appliesTo := append([]any{}, toList(bundle["mission_ids"])...)
			appliesTo = append(appliesTo, bundleID)
			signal, _ := toMap(mergeDicts(signalTemplate, map[string]any{
				"signal_id":            signalID,
				"kind":                 "bundle_governor_established",
				"source_type":          "leader",
				"source_id":            "leader-main",
				"mission_id":           governingMissionID,
				"applies_to":           appliesTo,
				"scope_type":           "mission-graph",
				"title":                "Bundle governor established",
				"summary":              firstNonBlank(note, fmt.Sprintf("Leader established %s as the explicit governor for %s.", governingMissionID, bundleID)),
				"authority_level":      "L2",
				"confidence":           "high",
				"bundle_id":            bundleID,
				"bundle_mission_ids":   cloneValue(toList(bundle["mission_ids"])),
				"governing_mission_id": governingMissionID,
				"selector_refs":        cloneValue(toList(bundle["selector_refs"])),
				"resolution_kind":      "explicit_bundle_governor",
				"required_action":      fmt.Sprintf("Keep bundle authority anchored to %s until the explicit governor decision is cleared or superseded.", governingMissionID),
			}))
			fillIfBlank(signal, "emitted_at", nowISO())
			fillIfBlank(signal, "state", "active")
			log["signals"] = append(toList(log["signals"]), signal)
			if err := writeSignalLog(log); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"signal": signal,
				"bundle": bundle,
			})
		})
	case "clear-governor":
		if len(args) < 2 {
			return errors.New("missions clear-governor requires BUNDLE_ID [NOTE]")
		}
		bundleID := strings.TrimSpace(args[1])
		note := strings.TrimSpace(strings.Join(args[2:], " "))
		if bundleID == "" {
			return errors.New("missions clear-governor requires BUNDLE_ID [NOTE]")
		}
		return withHarnessLock(func() error {
			log := readSignalLog()
			signals := toList(log["signals"])
			resolvedIDs := []string{}
			for index, item := range signals {
				signal, _ := toMap(item)
				if !activeSignal(signal) || stringValue(signal["kind"]) != "bundle_governor_established" || stringValue(signal["bundle_id"]) != bundleID {
					continue
				}
				signal["state"] = "resolved"
				signal["resolution_note"] = firstNonBlank(note, fmt.Sprintf("Leader cleared explicit governor for %s.", bundleID))
				signal["updated_at"] = nowISO()
				signals[index] = signal
				resolvedIDs = append(resolvedIDs, stringValue(signal["signal_id"]))
			}
			if len(resolvedIDs) == 0 {
				return fmt.Errorf("no active explicit bundle governor found for %s", bundleID)
			}
			log["signals"] = signals
			if err := writeSignalLog(log); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"bundle_id":           bundleID,
				"resolved_signal_ids": stringListToAny(dedupeStrings(resolvedIDs)),
				"resolution_note":     firstNonBlank(note, fmt.Sprintf("Leader cleared explicit governor for %s.", bundleID)),
			})
		})
	case "show-template":
		return printJSON(loadData(missionTemplatePath(), map[string]any{}))
	default:
		return fmt.Errorf("unknown missions subcommand: %s", args[0])
	}
	return nil
}

func cmdRuns(args []string) error {
	if len(args) < 1 {
		return errors.New("runs requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	registry := readRegistry()
	switch subcommand {
	case "init":
		return withHarnessLock(func() error {
			if err := writeRegistry(registry); err != nil {
				return err
			}
			fmt.Println(registryPath())
			return nil
		})
	case "list":
		return printJSON(registry)
	case "active":
		missionID := ""
		if len(args) > 1 {
			missionID = args[1]
		}
		runs := []any{}
		for _, item := range toList(registry["runs"]) {
			run, _ := toMap(item)
			if activeRun(run) && (missionID == "" || stringValue(run["mission_id"]) == missionID) {
				runs = append(runs, run)
			}
		}
		return printJSON(runs)
	case "get":
		if len(args) < 2 {
			return errors.New("runs get requires RUN_ID")
		}
		index := findRunIndex(registry, args[1])
		if index < 0 {
			return fmt.Errorf("run not found: %s", args[1])
		}
		return printJSON(toList(registry["runs"])[index])
	case "authority-audit":
		if len(args) < 2 {
			return errors.New("runs authority-audit requires RUN_ID")
		}
		index := findRunIndex(registry, args[1])
		if index < 0 {
			return fmt.Errorf("run not found: %s", args[1])
		}
		run, _ := toMap(toList(registry["runs"])[index])
		return printJSON(runAuthorityAudit(run, registry, readSignalLog()))
	case "add":
		payload, err := readStdinJSON()
		if err != nil {
			return err
		}
		return withHarnessLock(func() error {
			registry = readRegistry()
			templateRoot, _ := toMap(loadData(runTemplatePath(), map[string]any{}))
			template, _ := toMap(templateRoot["run"])
			merged, _ := toMap(mergeDicts(template, payload))
			if stringValue(merged["run_id"]) == "" {
				return errors.New("runs add requires run_id")
			}
			if findRunIndex(registry, stringValue(merged["run_id"])) >= 0 {
				return fmt.Errorf("run already exists: %s", stringValue(merged["run_id"]))
			}
			timestamp := nowISO()
			fillIfBlank(merged, "created_at", timestamp)
			if len(ownershipClaims(merged)) > 0 {
				setOwnershipSelectors(merged, ownershipClaims(merged))
			}
			state := stringValue(merged["state"])
			if state == "ready" || state == "executing" || state == "paused" || state == "waiting" || state == "blocked" {
				fillIfBlank(merged, "started_at", timestamp)
				fillIfBlank(merged, "last_reported_at", timestamp)
			}
			merged["updated_at"] = timestamp
			registry["runs"] = append(toList(registry["runs"]), merged)
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(merged)
		})
	case "update":
		if len(args) < 2 {
			return errors.New("runs update requires RUN_ID")
		}
		payload, err := readStdinJSON()
		if err != nil {
			return err
		}
		return withHarnessLock(func() error {
			registry = readRegistry()
			index := findRunIndex(registry, args[1])
			if index < 0 {
				return fmt.Errorf("run not found: %s", args[1])
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			merged, _ := toMap(mergeDicts(run, payload))
			if len(ownershipClaims(merged)) > 0 {
				setOwnershipSelectors(merged, ownershipClaims(merged))
			}
			merged["updated_at"] = nowISO()
			if err := persistContractForRun(merged); err != nil {
				return err
			}
			runs[index] = merged
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(merged)
		})
	case "set-state":
		if len(args) < 3 {
			return errors.New("runs set-state requires RUN_ID STATE [REASON]")
		}
		return withHarnessLock(func() error {
			registry = readRegistry()
			index := findRunIndex(registry, args[1])
			if index < 0 {
				return fmt.Errorf("run not found: %s", args[1])
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, args[2]); err != nil {
				return err
			}
			run["state"] = args[2]
			if len(args) > 3 {
				run["status_reason"] = args[3]
			}
			run["updated_at"] = nowISO()
			if args[2] == "completed" || args[2] == "abandoned" {
				run["completed_at"] = nowISO()
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "set-desired":
		if len(args) < 2 {
			return errors.New("runs set-desired requires RUN_ID [--goal=STATE] [--verification=STATE] [--review=STATE] [--checks=STATE] [--ship=STATE]")
		}
		runID := args[1]
		goalState := ""
		verificationStatus := ""
		reviewStatus := ""
		checksStatus := ""
		shipStatus := ""
		dependencyPolicy := ""
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--goal="):
				goalState = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--goal=")), goalState)
			case strings.HasPrefix(arg, "--verification="):
				verificationStatus = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--verification=")), verificationStatus)
			case strings.HasPrefix(arg, "--review="):
				reviewStatus = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--review=")), reviewStatus)
			case strings.HasPrefix(arg, "--checks="):
				checksStatus = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--checks=")), checksStatus)
			case strings.HasPrefix(arg, "--ship="):
				shipStatus = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--ship=")), shipStatus)
			case strings.HasPrefix(arg, "--dependency-policy="):
				dependencyPolicy = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--dependency-policy=")), dependencyPolicy)
			}
		}
		return withHarnessLock(func() error {
			registry = readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			desired := normalizeDesiredState(run)
			if goalState != "" {
				desired["goal_state"] = goalState
			}
			if verificationStatus != "" {
				desired["verification_status"] = verificationStatus
			}
			if reviewStatus != "" {
				desired["review_status"] = reviewStatus
			}
			if checksStatus != "" {
				desired["checks_status"] = checksStatus
			}
			if shipStatus != "" {
				desired["ship_status"] = shipStatus
			}
			run["desired_state"] = desired
			if dependencyPolicy != "" {
				run["dependency_policy"] = dependencyPolicy
			}
			run["updated_at"] = nowISO()
			if err := persistContractForRun(run); err != nil {
				return err
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "transfer-ownership":
		if len(args) < 3 {
			return errors.New("runs transfer-ownership requires SOURCE_RUN_ID TARGET_RUN_ID --owns=SELECTOR [--reason=TEXT] [--mode=move|copy|split]")
		}
		sourceRunID := args[1]
		targetRunID := args[2]
		selectors := []string{}
		reason := "Leader reassigned ownership selector."
		mode := "move"
		for _, arg := range args[3:] {
			switch {
			case strings.HasPrefix(arg, "--owns="):
				selectors = append(selectors, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--owns=")))...)
			case strings.HasPrefix(arg, "--reason="):
				reason = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--reason=")), reason)
			case strings.HasPrefix(arg, "--mode="):
				mode = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--mode=")), mode)
			}
		}
		if mode != "move" && mode != "copy" && mode != "split" {
			return errors.New("runs transfer-ownership --mode must be one of: move, copy, split")
		}
		if len(selectors) == 0 {
			return errors.New("runs transfer-ownership requires at least one --owns selector")
		}
		selectors = dedupeStrings(selectors)
		for _, selector := range selectors {
			if !boundedOwnershipSelector(selector) {
				return fmt.Errorf("runs transfer-ownership only supports bounded selectors, not ownership anchors: %s", selector)
			}
		}
		return withHarnessLock(func() error {
			registry = readRegistry()
			sourceIndex := findRunIndex(registry, sourceRunID)
			targetIndex := findRunIndex(registry, targetRunID)
			if sourceIndex < 0 || targetIndex < 0 {
				return errors.New("runs transfer-ownership requires two existing runs")
			}
			runs := toList(registry["runs"])
			sourceRun, _ := toMap(runs[sourceIndex])
			targetRun, _ := toMap(runs[targetIndex])
			if !activeRun(sourceRun) || !activeRun(targetRun) {
				return errors.New("runs transfer-ownership requires both source and target runs to be active")
			}
			if stringValue(sourceRun["mission_id"]) == "" || stringValue(sourceRun["mission_id"]) != stringValue(targetRun["mission_id"]) {
				return errors.New("runs transfer-ownership currently requires source and target runs to belong to the same mission")
			}
			if mode == "copy" {
				if overlapPolicy(sourceRun) != "shared" || overlapPolicy(targetRun) != "shared" {
					return errors.New("runs transfer-ownership --mode=copy requires both source and target runs to allow shared overlap")
				}
			}
			sourceSelectors := ownershipClaims(sourceRun)
			sourceExclusions := ownershipExclusions(sourceRun)
			targetSelectors := ownershipClaims(targetRun)
			sourceSet := map[string]bool{}
			for _, selector := range sourceSelectors {
				sourceSet[selector] = true
			}
			missing := []string{}
			splitParents := map[string][]string{}
			for _, selector := range selectors {
				if mode == "split" {
					if sourceSet[selector] {
						missing = append(missing, selector+" (exact selector already owned; use --mode=move for exact handoff)")
						continue
					}
					for _, sourceSelector := range sourceSelectors {
						if strings.TrimSpace(sourceSelector) != strings.TrimSpace(selector) && selectorContains(sourceSelector, selector) {
							splitParents[selector] = append(splitParents[selector], sourceSelector)
						}
					}
					if len(splitParents[selector]) == 0 {
						missing = append(missing, selector+" (split requires a broader source selector)")
					}
					continue
				}
				if !sourceSet[selector] {
					missing = append(missing, selector)
				}
			}
			if len(missing) > 0 {
				return fmt.Errorf("source run does not own requested selectors: %s", strings.Join(missing, ", "))
			}
			remainingSource := []string{}
			transferSet := map[string]bool{}
			for _, selector := range selectors {
				transferSet[selector] = true
			}
			for _, selector := range sourceSelectors {
				if mode == "copy" || mode == "split" || !transferSet[selector] {
					remainingSource = append(remainingSource, selector)
				}
			}
			nextSourceExclusions := dedupeStrings(sourceExclusions)
			if mode == "split" {
				nextSourceExclusions = dedupeStrings(append(nextSourceExclusions, selectors...))
			}
			targetSelectors = dedupeStrings(append(targetSelectors, selectors...))
			setOwnershipSelectors(sourceRun, remainingSource)
			setOwnershipExclusions(sourceRun, nextSourceExclusions)
			setOwnershipSelectors(targetRun, targetSelectors)
			setOwnershipExclusions(targetRun, ownershipExclusions(targetRun))
			if err := persistContractForRun(sourceRun); err != nil {
				return err
			}
			if err := persistContractForRun(targetRun); err != nil {
				return err
			}
			event := map[string]any{
				"transfer_id":     "transfer-" + nowTime().UTC().Format("20060102150405"),
				"selector_refs":   stringListToAny(selectors),
				"from_run_id":     sourceRunID,
				"to_run_id":       targetRunID,
				"reason":          reason,
				"mode":            mode,
				"authority_level": "L2",
				"transferred_at":  nowISO(),
			}
			sourceRun["ownership_transfers"] = append(ownershipTransferEvents(sourceRun), cloneValue(event))
			targetRun["ownership_transfers"] = append(ownershipTransferEvents(targetRun), cloneValue(event))
			if mode == "copy" {
				sourceRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership selectors shared with another bounded run.")
				sourceRun["next_expected_action"] = "Continue within the shared bounded ownership scope without drifting across other claims."
				targetRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership selectors shared in.")
				targetRun["next_expected_action"] = "Proceed inside the explicitly shared bounded ownership scope."
			} else if mode == "split" {
				sourceRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership subtree delegated while the broader boundary stays owned with exclusions.")
				sourceRun["next_expected_action"] = "Continue within the remaining parent boundary without reclaiming the delegated child subtree."
				targetRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership subtree transferred in under a split parent boundary.")
				targetRun["next_expected_action"] = "Proceed inside the delegated child boundary without widening back into the parent scope."
			} else {
				sourceRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership selectors transferred away.")
				sourceRun["next_expected_action"] = "Continue within the remaining bounded ownership scope."
				targetRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership selectors transferred in.")
				targetRun["next_expected_action"] = "Proceed with the newly transferred bounded ownership scope."
			}
			sourceRun["updated_at"] = nowISO()
			targetRun["updated_at"] = nowISO()
			runs[sourceIndex] = sourceRun
			runs[targetIndex] = targetRun
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"transfer":   event,
				"source_run": sourceRun,
				"target_run": targetRun,
			})
		})
	default:
		return fmt.Errorf("unknown runs subcommand: %s", subcommand)
	}
}

func cmdSignals(args []string) error {
	if len(args) < 1 {
		return errors.New("signals requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	log := readSignalLog()
	switch subcommand {
	case "init":
		return withHarnessLock(func() error {
			if err := writeSignalLog(log); err != nil {
				return err
			}
			fmt.Println(signalLogPath())
			return nil
		})
	case "list":
		return printJSON(log)
	case "active":
		missionID := ""
		if len(args) > 1 {
			missionID = args[1]
		}
		signals := []any{}
		for _, item := range toList(log["signals"]) {
			signal, _ := toMap(item)
			if activeSignal(signal) && (missionID == "" || stringValue(signal["mission_id"]) == missionID) {
				signals = append(signals, signal)
			}
		}
		return printJSON(signals)
	case "emit":
		payload, err := readStdinJSON()
		if err != nil {
			return err
		}
		return withHarnessLock(func() error {
			log = readSignalLog()
			templateRoot, _ := toMap(loadData(signalTemplatePath(), map[string]any{}))
			template, _ := toMap(templateRoot["signal"])
			merged, _ := toMap(mergeDicts(template, payload))
			if stringValue(merged["signal_id"]) == "" {
				return errors.New("signals emit requires signal_id")
			}
			if findSignalIndex(log, stringValue(merged["signal_id"])) >= 0 {
				return fmt.Errorf("signal already exists: %s", stringValue(merged["signal_id"]))
			}
			merged["kind"] = normalizeSignalKind(stringValue(merged["kind"]))
			fillIfBlank(merged, "emitted_at", nowISO())
			fillIfBlank(merged, "state", "active")
			log["signals"] = append(toList(log["signals"]), merged)
			if err := writeSignalLog(log); err != nil {
				return err
			}
			return printJSON(merged)
		})
	case "set-state":
		if len(args) < 3 {
			return errors.New("signals set-state requires SIGNAL_ID STATE [NOTE]")
		}
		return withHarnessLock(func() error {
			log = readSignalLog()
			index := findSignalIndex(log, args[1])
			if index < 0 {
				return fmt.Errorf("signal not found: %s", args[1])
			}
			signals := toList(log["signals"])
			signal, _ := toMap(signals[index])
			signal["state"] = args[2]
			if len(args) > 3 {
				signal["resolution_note"] = args[3]
			}
			signals[index] = signal
			log["signals"] = signals
			if err := writeSignalLog(log); err != nil {
				return err
			}
			return printJSON(signal)
		})
	case "adjudicate":
		targetAll := false
		signalID := ""
		for _, arg := range args[1:] {
			if arg == "--all" {
				targetAll = true
			} else if !strings.HasPrefix(arg, "--") && signalID == "" {
				signalID = arg
			}
		}
		if !targetAll && signalID == "" {
			return errors.New("signals adjudicate requires SIGNAL_ID or --all")
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			log = readSignalLog()
			signals := toList(log["signals"])
			resolved := []any{}
			for index, item := range signals {
				signal, _ := toMap(item)
				if stringValue(signal["state"]) != "active" {
					continue
				}
				if !targetAll && stringValue(signal["signal_id"]) != signalID {
					continue
				}
				superseded, resolution := signalSupersession(signal, registry)
				if !superseded {
					continue
				}
				signal["state"] = "resolved"
				signal["resolved_at"] = nowISO()
				for key, value := range resolution {
					signal[key] = value
				}
				signals[index] = signal
				resolved = append(resolved, signal)
			}
			log["signals"] = signals
			if len(resolved) > 0 {
				if err := writeSignalLog(log); err != nil {
					return err
				}
			}
			return printJSON(map[string]any{
				"resolved_count": len(resolved),
				"signals":        resolved,
			})
		})
	default:
		return fmt.Errorf("unknown signals subcommand: %s", subcommand)
	}
}

func cmdContracts(args []string) error {
	if len(args) < 1 {
		return errors.New("contracts requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	switch args[0] {
	case "materialize":
		if len(args) < 2 {
			return errors.New("contracts materialize requires OUTPUT_PATH")
		}
		outputPath := args[1]
		if !filepath.IsAbs(outputPath) {
			outputPath = filepath.Join(repoRoot(), outputPath)
		}
		payload, err := readStdinJSON()
		if err != nil {
			return err
		}
		templateRoot, _ := toMap(loadData(contractTemplatePath(), map[string]any{}))
		template, _ := toMap(templateRoot["contract"])
		merged, _ := toMap(mergeDicts(template, payload))
		freshness, _ := toMap(merged["freshness"])
		if freshness == nil {
			freshness = map[string]any{}
		}
		if stringValue(freshness["created_at"]) == "" {
			freshness["created_at"] = nowISO()
		}
		merged["freshness"] = freshness
		if err := dumpData(outputPath, map[string]any{"version": 1, "contract": merged}); err != nil {
			return err
		}
		fmt.Println(outputPath)
	case "show-template":
		return printJSON(loadData(contractTemplatePath(), map[string]any{}))
	default:
		return fmt.Errorf("unknown contracts subcommand: %s", args[0])
	}
	return nil
}

func schedulerSnapshotsRoot() string {
	return filepath.Join(runtimeRoot(), "scheduler", "snapshots")
}

func schedulerLatestSnapshotPath() string {
	return filepath.Join(schedulerSnapshotsRoot(), "latest.json")
}

func schedulerSnapshotPath(snapshotID string) string {
	return filepath.Join(schedulerSnapshotsRoot(), snapshotID+".json")
}

func latestSchedulerSnapshot() map[string]any {
	payload, _ := toMap(loadData(schedulerLatestSnapshotPath(), map[string]any{}))
	snapshot, _ := toMap(payload["snapshot"])
	return snapshot
}

func schedulerDecisionByID(snapshot map[string]any, decisionID string) map[string]any {
	for _, item := range toList(snapshot["required_human_decisions"]) {
		decision, _ := toMap(item)
		if stringValue(decision["decision_id"]) == decisionID {
			return decision
		}
	}
	for _, item := range listMeetings(true) {
		meeting, _ := toMap(item)
		if stringValue(meeting["state"]) != "active" {
			continue
		}
		if stringValue(meeting["source_decision_id"]) != decisionID {
			continue
		}
		decision, _ := toMap(cloneValue(meeting["source_decision_packet"]))
		if len(decision) == 0 {
			decision = map[string]any{}
		}
		decision["decision_id"] = firstNonBlank(stringValue(decision["decision_id"]), decisionID)
		decision = projectActiveMeetingProgressIntoDecision(decision, meeting)
		return decision
	}
	return nil
}

func schedulerPriorityRank(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "p0", "critical", "urgent", "highest":
		return 5
	case "p1", "high":
		return 4
	case "p2", "medium", "normal":
		return 3
	case "p3", "low":
		return 2
	case "p4", "deferred":
		return 1
	default:
		return 0
	}
}

func schedulerRiskRank(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "critical":
		return 5
	case "high":
		return 4
	case "medium":
		return 3
	case "low":
		return 2
	default:
		return 0
	}
}

func schedulerMaintenanceActions(missionRegistry, registry, signalLog map[string]any) ([]any, map[string]int) {
	actions := []any{}
	tally := map[string]int{
		"run":     0,
		"signal":  0,
		"meeting": 0,
	}
	runStaleAfter := floatValue(policyValue("cleanup.stale_run_after_hours", 24))
	signalStaleAfter := floatValue(policyValue("cleanup.stale_signal_after_hours", 12))
	meetingStaleAfter := floatValue(policyValue("cleanup.stale_meeting_after_hours", 6))

	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if !activeRun(run) {
			continue
		}
		runID := stringValue(run["run_id"])
		missionID := stringValue(run["mission_id"])
		reportAge := ageHours(stringValue(run["last_reported_at"]))
		switch {
		case reportAge == nil:
			actions = append(actions, map[string]any{
				"target_kind":        "run",
				"target_id":          runID,
				"mission_id":         missionID,
				"action_kind":        "refresh_stale_run",
				"reason":             "Active run is missing last_reported_at, so scheduler cannot trust its freshness.",
				"recommended_action": "Have the owner emit a fresh report or repair the run metadata before widening execution.",
			})
			tally["run"]++
		case *reportAge > runStaleAfter:
			actions = append(actions, map[string]any{
				"target_kind":        "run",
				"target_id":          runID,
				"mission_id":         missionID,
				"action_kind":        "refresh_stale_run",
				"reason":             fmt.Sprintf("Active run has not reported for %.2f hours.", *reportAge),
				"recommended_action": "Refresh the run report or retire the stale run before widening execution.",
			})
			tally["run"]++
		}
	}

	for _, item := range toList(signalLog["signals"]) {
		signal, _ := toMap(item)
		if !activeSignal(signal) {
			continue
		}
		signalID := stringValue(signal["signal_id"])
		missionID := stringValue(signal["mission_id"])
		if emittedAge := ageHours(stringValue(signal["emitted_at"])); emittedAge != nil && *emittedAge > signalStaleAfter {
			actions = append(actions, map[string]any{
				"target_kind":        "signal",
				"target_id":          signalID,
				"mission_id":         missionID,
				"action_kind":        "resolve_stale_signal",
				"reason":             fmt.Sprintf("Active signal has been open for %.2f hours.", *emittedAge),
				"recommended_action": "Resolve or refresh the stale signal before letting it influence future scheduling.",
			})
			tally["signal"]++
			continue
		}
		if freshUntil, err := parseISO(stringValue(signal["fresh_until"])); err == nil && freshUntil != nil && freshUntil.Before(nowTime()) {
			actions = append(actions, map[string]any{
				"target_kind":        "signal",
				"target_id":          signalID,
				"mission_id":         missionID,
				"action_kind":        "resolve_stale_signal",
				"reason":             "Active signal has expired freshness and should not keep steering the scheduler silently.",
				"recommended_action": "Resolve the expired signal or re-emit it with fresh evidence before widening execution.",
			})
			tally["signal"]++
		}
	}
	for _, issue := range invalidExplicitBundleGovernorSignals(missionRegistry, registry, signalLog) {
		invalid, _ := toMap(issue)
		actions = append(actions, map[string]any{
			"target_kind":          "signal",
			"target_id":            stringValue(invalid["signal_id"]),
			"mission_id":           stringValue(invalid["mission_id"]),
			"bundle_id":            stringValue(invalid["bundle_id"]),
			"governing_mission_id": stringValue(invalid["governing_mission_id"]),
			"action_kind":          "clear_invalid_bundle_governor",
			"reason":               stringValue(invalid["reason"]),
			"recommended_action":   stringValue(invalid["recommended_action"]),
		})
		tally["signal"]++
	}

	for _, item := range listMeetings(true) {
		meeting, _ := toMap(item)
		if stringValue(meeting["state"]) != "active" {
			continue
		}
		meetingID := stringValue(meeting["meeting_id"])
		missionID := stringValue(meeting["mission_id"])
		if meetingAge := ageHours(stringValue(meeting["updated_at"])); meetingAge != nil && *meetingAge > meetingStaleAfter {
			actions = append(actions, map[string]any{
				"target_kind":        "meeting",
				"target_id":          meetingID,
				"mission_id":         missionID,
				"action_kind":        "refresh_stale_meeting",
				"reason":             fmt.Sprintf("Active meeting has not been refreshed for %.2f hours.", *meetingAge),
				"recommended_action": "Refresh or close the stale meeting before trusting it as current governance state.",
			})
			tally["meeting"]++
		}
	}

	return actions, tally
}

func aggregateSchedulerHumanDecisions(decisions []any) []any {
	bundleOrder := []string{}
	bundlePackets := map[string]map[string]any{}
	aggregated := []any{}
	for _, item := range decisions {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		bundleID := strings.TrimSpace(stringValue(decision["bundle_id"]))
		decisionKind := strings.TrimSpace(stringValue(decision["decision_kind"]))
		decisionSurface := strings.TrimSpace(stringValue(decision["recommended_decision_surface"]))
		if bundleID == "" || decisionKind == "" || decisionKind != "establish_bundle_governor" {
			if stringValue(decision["decision_scope"]) == "" {
				decision["decision_scope"] = "mission"
			}
			decision["decision_id"] = schedulerDecisionPacketID(decision)
			aggregated = append(aggregated, decision)
			continue
		}
		packetKey := bundleID + "::" + decisionKind + "::" + decisionSurface
		packet, exists := bundlePackets[packetKey]
		if !exists {
			packet = map[string]any{
				"decision_scope":                "bundle",
				"bundle_id":                     bundleID,
				"decision_kind":                 decisionKind,
				"recommended_decision_surface":  decisionSurface,
				"recommended_decision_protocol": stringValue(decision["recommended_decision_protocol"]),
				"recommended_protocol_reason":   stringValue(decision["recommended_protocol_reason"]),
				"recommended_action":            stringValue(decision["recommended_action"]),
				"bundle_required_action":        stringValue(decision["bundle_required_action"]),
				"bundle_state":                  stringValue(decision["bundle_state"]),
				"bundle_resolution_kind":        stringValue(decision["bundle_resolution_kind"]),
				"bundle_governing_mission_id":   stringValue(decision["bundle_governing_mission_id"]),
				"candidate_governor_ids":        cloneValue(getKey(decision, "candidate_governor_ids", []any{})),
				"candidate_governor_packets":    cloneValue(getKey(decision, "candidate_governor_packets", []any{})),
				"meeting_trigger_refs":          []any{},
				"mission_ids":                   []any{},
				"mission_count":                 0,
				"priority":                      "",
				"effective_risk_level":          "",
				"release_adjacent":              false,
				"updated_at":                    "",
				"note":                          "",
				"notes":                         []any{},
			}
			bundlePackets[packetKey] = packet
			bundleOrder = append(bundleOrder, packetKey)
		}
		packet["meeting_trigger_refs"] = stringListToAny(dedupeStrings(append(
			anyListToStrings(toList(packet["meeting_trigger_refs"])),
			anyListToStrings(toList(decision["meeting_trigger_refs"]))...,
		)))
		packet["mission_ids"] = stringListToAny(dedupeStrings(append(
			append(
				anyListToStrings(toList(packet["mission_ids"])),
				stringValue(decision["mission_id"]),
			),
			anyListToStrings(toList(decision["bundle_mission_ids"]))...,
		)))
		notes := dedupeStrings(append(
			anyListToStrings(toList(packet["notes"])),
			strings.TrimSpace(stringValue(decision["note"])),
		))
		packet["notes"] = stringListToAny(notes)
		packet["mission_count"] = len(anyListToStrings(toList(packet["mission_ids"])))
		if schedulerPriorityRank(stringValue(decision["priority"])) > schedulerPriorityRank(stringValue(packet["priority"])) {
			packet["priority"] = stringValue(decision["priority"])
		}
		if schedulerRiskRank(stringValue(decision["effective_risk_level"])) > schedulerRiskRank(stringValue(packet["effective_risk_level"])) {
			packet["effective_risk_level"] = stringValue(decision["effective_risk_level"])
		}
		if boolValue(decision["release_adjacent"]) {
			packet["release_adjacent"] = true
		}
		if protocolRank := decisionProtocolRank(stringValue(decision["recommended_decision_protocol"])); protocolRank > decisionProtocolRank(stringValue(packet["recommended_decision_protocol"])) {
			packet["recommended_decision_protocol"] = stringValue(decision["recommended_decision_protocol"])
			packet["recommended_protocol_reason"] = stringValue(decision["recommended_protocol_reason"])
		}
		if stringValue(decision["updated_at"]) > stringValue(packet["updated_at"]) {
			packet["updated_at"] = stringValue(decision["updated_at"])
		}
		packet["recommended_action"] = firstNonBlank(
			stringValue(packet["bundle_required_action"]),
			stringValue(packet["recommended_action"]),
			stringValue(decision["recommended_action"]),
		)
		packet["note"] = fmt.Sprintf(
			"Bundle %s requires one governing-mission decision across %d overlapping mission(s): %s",
			bundleID,
			packet["mission_count"],
			strings.Join(notes, " "),
		)
		packet["decision_id"] = schedulerDecisionPacketID(packet)
	}
	for _, packetKey := range bundleOrder {
		packet := bundlePackets[packetKey]
		packet["decision_id"] = schedulerDecisionPacketID(packet)
		aggregated = append(aggregated, packet)
	}
	return aggregated
}

func recoverActiveBundleDecisionPackets(decisions []any) []any {
	seenDecisionIDs := map[string]bool{}
	for _, item := range decisions {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		if decisionID := strings.TrimSpace(stringValue(decision["decision_id"])); decisionID != "" {
			seenDecisionIDs[decisionID] = true
		}
	}

	activeByDecisionID := map[string]map[string]any{}
	for _, item := range listMeetings(true) {
		meeting, ok := toMap(item)
		if !ok || stringValue(meeting["state"]) != "active" {
			continue
		}
		decisionID := strings.TrimSpace(stringValue(meeting["source_decision_id"]))
		if decisionID == "" || seenDecisionIDs[decisionID] {
			continue
		}
		sourceScope := firstNonBlank(
			strings.TrimSpace(stringValue(meeting["source_decision_scope"])),
			strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.decision_scope", ""))),
		)
		if sourceScope != "bundle" {
			continue
		}
		packet, _ := toMap(cloneValue(meeting["source_decision_packet"]))
		if len(packet) == 0 {
			continue
		}
		activeEntry, exists := activeByDecisionID[decisionID]
		if !exists {
			packet["decision_id"] = firstNonBlank(stringValue(packet["decision_id"]), decisionID)
			packet["recommended_decision_surface"] = "active-council"
			packet["active_meeting_ids"] = []any{}
			packet["active_meeting_count"] = 0
			activeEntry = packet
			activeByDecisionID[decisionID] = activeEntry
		}
		activeEntry = projectActiveMeetingProgressIntoDecision(activeEntry, meeting)
		activeEntry["updated_at"] = firstNonBlank(
			stringValue(activeEntry["updated_at"]),
			stringValue(meeting["updated_at"]),
			stringValue(meeting["created_at"]),
		)
	}

	for _, packet := range activeByDecisionID {
		decisions = append(decisions, packet)
	}
	return decisions
}

func decisionSurfaceRank(surface string) int {
	switch strings.TrimSpace(surface) {
	case "active-council":
		return 3
	case "council":
		return 2
	case "leader-decision":
		return 1
	default:
		return 0
	}
}

func decisionProtocolRank(protocol string) int {
	switch normalizeDecisionProtocol(protocol) {
	case "incident-command":
		return 3
	case "evidence-review":
		return 2
	case "decision-record":
		return 1
	default:
		return 0
	}
}

func decisionScopeRank(scope string) int {
	switch strings.TrimSpace(scope) {
	case "bundle":
		return 2
	case "mission":
		return 1
	default:
		return 0
	}
}

func meetingPhaseProgressRank(phase string) int {
	switch strings.TrimSpace(phase) {
	case "output":
		return 5
	case "synthesis":
		return 4
	case "cross_examination":
		return 3
	case "independent_review":
		return 2
	case "setup":
		return 1
	default:
		return 0
	}
}

func schedulerDecisionProgressRank(decision map[string]any) int {
	if decision == nil {
		return 0
	}
	if strings.TrimSpace(stringValue(decision["recommended_decision_surface"])) != "active-council" {
		return 0
	}
	return meetingPhaseProgressRank(stringValue(decision["active_decision_phase"]))
}

func schedulerDecisionKey(decision map[string]any) string {
	scope := firstNonBlank(stringValue(decision["decision_scope"]), "mission")
	if scope == "bundle" {
		if bundleID := strings.TrimSpace(stringValue(decision["bundle_id"])); bundleID != "" {
			return "bundle:" + bundleID
		}
	}
	if missionID := strings.TrimSpace(stringValue(decision["mission_id"])); missionID != "" {
		return "mission:" + missionID
	}
	return ""
}

func meetingCurrentWaveBatch(wavePlan, debatePlan map[string]any) map[string]any {
	if debatePlan != nil {
		currentBatchID := strings.TrimSpace(stringValue(debatePlan["current_wave_batch_id"]))
		if currentBatchID != "" {
			for _, item := range toList(debatePlan["wave_batches"]) {
				batch, ok := toMap(item)
				if !ok {
					continue
				}
				if strings.TrimSpace(stringValue(batch["batch_id"])) == currentBatchID {
					return batch
				}
			}
		}
	}
	if wavePlan == nil || intValue(wavePlan["open_wave_action_count"]) == 0 {
		return nil
	}
	waveActions := toList(wavePlan["wave_actions"])
	summary := "Complete the current independent-review wave obligations before debate can advance."
	if len(waveActions) == 1 {
		summary = firstNonBlank(strings.TrimSpace(stringValue(getKey(waveActions[0], "summary", ""))), summary)
	}
	assignedRole := ""
	assignedRunID := ""
	if len(waveActions) == 1 {
		assignedRole = stringValue(getKey(waveActions[0], "assigned_role", ""))
		assignedRunID = stringValue(getKey(waveActions[0], "assigned_run_id", ""))
	}
	return map[string]any{
		"batch_id":        "wave-batch-independent-review",
		"expected_phase":  "independent_review",
		"round_number":    0,
		"assigned_role":   assignedRole,
		"assigned_run_id": assignedRunID,
		"action_count":    intValue(wavePlan["open_wave_action_count"]),
		"claim_refs":      []any{},
		"summary":         summary,
	}
}

func meetingCurrentWaveActions(wavePlan, debatePlan map[string]any) []any {
	currentBatch := meetingCurrentWaveBatch(wavePlan, debatePlan)
	if currentBatch == nil {
		return []any{}
	}
	currentBatchID := strings.TrimSpace(stringValue(currentBatch["batch_id"]))
	if currentBatchID == "wave-batch-independent-review" {
		actions := []any{}
		for _, item := range toList(wavePlan["wave_actions"]) {
			action, ok := toMap(item)
			if !ok {
				continue
			}
			state := strings.TrimSpace(stringValue(action["state"]))
			if state != "" && state != "open" {
				continue
			}
			actions = append(actions, cloneValue(action))
		}
		return actions
	}
	actionByID := map[string]map[string]any{}
	for _, bucket := range [][]any{toList(debatePlan["route_actions"]), toList(debatePlan["round_actions"])} {
		for _, item := range bucket {
			action, ok := toMap(item)
			if !ok {
				continue
			}
			actionID := firstNonBlank(strings.TrimSpace(stringValue(action["round_id"])), strings.TrimSpace(stringValue(action["route_id"])))
			if actionID == "" {
				continue
			}
			actionByID[actionID] = action
		}
	}
	actions := []any{}
	for _, item := range toList(currentBatch["action_ids"]) {
		actionID := strings.TrimSpace(stringValue(item))
		if actionID == "" {
			continue
		}
		if action, ok := actionByID[actionID]; ok {
			actions = append(actions, cloneValue(action))
		}
	}
	return actions
}

func meetingBlockedWaveBatchSummaries(debatePlan map[string]any) []any {
	if debatePlan == nil {
		return []any{}
	}
	blocked := []any{}
	for _, item := range toList(debatePlan["wave_batches"]) {
		batch, ok := toMap(item)
		if !ok || strings.TrimSpace(stringValue(batch["state"])) != "blocked" {
			continue
		}
		blocked = append(blocked, map[string]any{
			"batch_id":        stringValue(batch["batch_id"]),
			"expected_phase":  stringValue(batch["expected_phase"]),
			"round_number":    intValue(batch["round_number"]),
			"assigned_role":   stringValue(batch["assigned_role"]),
			"assigned_run_id": stringValue(batch["assigned_run_id"]),
			"action_count":    intValue(batch["action_count"]),
			"lane_count":      intValue(batch["lane_count"]),
			"claim_refs":      cloneValue(getKey(batch, "claim_refs", []any{})),
			"summary":         stringValue(batch["summary"]),
		})
	}
	return blocked
}

func projectActiveMeetingProgressIntoDecision(decision, meeting map[string]any) map[string]any {
	if decision == nil || meeting == nil {
		return decision
	}
	meetingID := strings.TrimSpace(stringValue(meeting["meeting_id"]))
	activeMeetingIDs := dedupeStrings(append(
		anyListToStrings(toList(decision["active_meeting_ids"])),
		meetingID,
	))
	decision["active_meeting_ids"] = stringListToAny(activeMeetingIDs)
	decision["active_meeting_count"] = len(activeMeetingIDs)
	decision["recommended_decision_surface"] = "active-council"
	decision["active_decision_phase"] = firstNonBlank(strings.TrimSpace(stringValue(meeting["phase"])), "setup")
	decision["active_decision_progress_rank"] = schedulerDecisionProgressRank(decision)
	decision["active_decision_protocol"] = normalizeDecisionProtocol(firstNonBlank(stringValue(meeting["decision_protocol"]), "evidence-review"))
	orchestration, _ := toMap(meeting["orchestration"])
	if orchestration != nil {
		if nextRequirement := strings.TrimSpace(stringValue(orchestration["next_requirement"])); nextRequirement != "" {
			decision["active_decision_next_requirement"] = nextRequirement
			decision["recommended_action"] = nextRequirement
		}
		if wavePlan, _ := toMap(orchestration["wave_plan"]); wavePlan != nil {
			decision["active_decision_open_wave_action_count"] = intValue(wavePlan["open_wave_action_count"])
		}
		if debatePlan, _ := toMap(orchestration["debate_plan"]); debatePlan != nil {
			decision["active_decision_open_route_action_count"] = intValue(debatePlan["open_route_action_count"])
			decision["active_decision_open_round_action_count"] = intValue(debatePlan["open_round_action_count"])
			wavePlan, _ := toMap(orchestration["wave_plan"])
			if currentBatch := meetingCurrentWaveBatch(wavePlan, debatePlan); currentBatch != nil {
				decision["active_decision_current_wave_batch_id"] = stringValue(currentBatch["batch_id"])
				decision["active_decision_current_wave_batch_phase"] = stringValue(currentBatch["expected_phase"])
				decision["active_decision_current_wave_batch_round_number"] = intValue(currentBatch["round_number"])
				decision["active_decision_current_wave_batch_assigned_role"] = stringValue(currentBatch["assigned_role"])
				decision["active_decision_current_wave_batch_assigned_run_id"] = stringValue(currentBatch["assigned_run_id"])
				decision["active_decision_current_wave_batch_action_count"] = intValue(currentBatch["action_count"])
				decision["active_decision_current_wave_batch_claim_refs"] = cloneValue(getKey(currentBatch, "claim_refs", []any{}))
				if summary := strings.TrimSpace(stringValue(currentBatch["summary"])); summary != "" {
					decision["active_decision_current_wave_batch_summary"] = summary
				}
			}
			decision["active_decision_current_wave_actions"] = meetingCurrentWaveActions(wavePlan, debatePlan)
			decision["active_decision_blocked_wave_batches"] = meetingBlockedWaveBatchSummaries(debatePlan)
			if summary := strings.TrimSpace(stringValue(debatePlan["current_wave_batch_summary"])); summary != "" {
				decision["active_decision_current_wave_batch_summary"] = summary
			}
		}
		if gaps := toList(orchestration["last_gap_report"]); len(gaps) > 0 {
			decision["active_decision_gap_count"] = len(gaps)
		}
		if outputPacket, _ := toMap(orchestration["output_packet"]); outputPacket != nil {
			decision["active_decision_output_packet_kind"] = stringValue(outputPacket["output_packet_kind"])
			decision["active_decision_output_packet"] = cloneValue(outputPacket)
			decision["active_decision_output_debate_lane_count"] = intValue(outputPacket["debate_lane_count"])
			decision["active_decision_output_open_debate_lane_count"] = intValue(outputPacket["open_debate_lane_count"])
			decision["active_decision_output_unresolved_debate_lane_count"] = intValue(outputPacket["unresolved_debate_lane_count"])
			decision["active_decision_output_debate_lane_packets"] = cloneValue(getKey(outputPacket, "debate_lane_packets", []any{}))
			decision["active_decision_output_follow_up_packet_count"] = intValue(outputPacket["follow_up_packet_count"])
			decision["active_decision_output_open_follow_up_packet_count"] = intValue(outputPacket["open_follow_up_packet_count"])
			decision["active_decision_output_follow_up_packets"] = cloneValue(getKey(outputPacket, "follow_up_packets", []any{}))
		}
	}
	meetingUpdatedAt := firstNonBlank(strings.TrimSpace(stringValue(meeting["updated_at"])), strings.TrimSpace(stringValue(meeting["created_at"])))
	if meetingUpdatedAt != "" {
		decision["active_decision_last_updated_at"] = meetingUpdatedAt
		if meetingAge := ageHours(meetingUpdatedAt); meetingAge != nil {
			decision["active_decision_age_hours"] = fmt.Sprintf("%.2f", *meetingAge)
			meetingStaleAfter := floatValue(policyValue("cleanup.stale_meeting_after_hours", 6))
			if *meetingAge > meetingStaleAfter {
				decision["active_decision_stale"] = true
				decision["recommended_action"] = fmt.Sprintf(
					"Refresh or close the stale active council(s) for this decision before trusting them as current governance state: %s.",
					strings.Join(activeMeetingIDs, ", "),
				)
			}
		}
	}
	decision["recommended_action"] = firstNonBlank(
		stringValue(decision["recommended_action"]),
		fmt.Sprintf("Review the active council(s) already materialized for this decision: %s.", strings.Join(activeMeetingIDs, ", ")),
	)
	return decision
}

func schedulerDecisionPacketID(decision map[string]any) string {
	parts := []string{
		firstNonBlank(stringValue(decision["decision_scope"]), "mission"),
		firstNonBlank(stringValue(decision["decision_kind"]), "decision"),
	}
	if bundleID := strings.TrimSpace(stringValue(decision["bundle_id"])); bundleID != "" {
		parts = append(parts, "bundle:"+bundleID)
	} else if missionID := strings.TrimSpace(stringValue(decision["mission_id"])); missionID != "" {
		parts = append(parts, "mission:"+missionID)
	}
	return "decision-" + stableSlugID(strings.Join(parts, "::"))
}

func promoteMaterializedDecisionPackets(decisions []any, activeMeetingIDsByDecisionID map[string][]string) []any {
	promoted := append([]any{}, decisions...)
	for index, item := range promoted {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		decisionID := stringValue(decision["decision_id"])
		activeMeetingIDs := dedupeStrings(activeMeetingIDsByDecisionID[decisionID])
		if len(activeMeetingIDs) == 0 {
			decision["active_meeting_ids"] = []any{}
			decision["active_meeting_count"] = 0
			promoted[index] = decision
			continue
		}
		for _, item := range listMeetings(true) {
			meeting, _ := toMap(item)
			if !containsString(activeMeetingIDs, stringValue(meeting["meeting_id"])) {
				continue
			}
			decision = projectActiveMeetingProgressIntoDecision(decision, meeting)
			break
		}
		promoted[index] = decision
	}
	return promoted
}

func schedulerDecisionClass(decision map[string]any) string {
	return fmt.Sprintf(
		"%d|%d|%d|%d|%t",
		decisionSurfaceRank(stringValue(decision["recommended_decision_surface"])),
		decisionScopeRank(stringValue(decision["decision_scope"])),
		schedulerPriorityRank(stringValue(decision["priority"])),
		schedulerRiskRank(stringValue(decision["effective_risk_level"])),
		boolValue(decision["release_adjacent"]),
	)
}

func schedulerFollowUpSignalUrgency(kind string) int {
	switch normalizeSignalKind(kind) {
	case "risk_raised":
		return 4
	case "review_feedback":
		return 3
	case "priority_change":
		return 2
	case "requires_meeting", "interface_change":
		return 1
	default:
		return 0
	}
}

func schedulerFollowUpSignalUrgencyRank(kinds []string) int {
	best := 0
	for _, kind := range dedupeStrings(kinds) {
		if urgency := schedulerFollowUpSignalUrgency(kind); urgency > best {
			best = urgency
		}
	}
	return best
}

func schedulerPrimaryFollowUpSignalKind(kinds []string) string {
	bestKind := ""
	bestUrgency := 0
	for _, kind := range dedupeStrings(kinds) {
		if urgency := schedulerFollowUpSignalUrgency(kind); urgency > bestUrgency {
			bestKind = normalizeSignalKind(kind)
			bestUrgency = urgency
		}
	}
	return bestKind
}

func schedulerMaintenanceKey(action map[string]any) string {
	if bundleID := strings.TrimSpace(stringValue(action["bundle_id"])); bundleID != "" {
		return "bundle:" + bundleID
	}
	targetKind := firstNonBlank(strings.TrimSpace(stringValue(action["target_kind"])), "unknown")
	targetID := firstNonBlank(strings.TrimSpace(stringValue(action["target_id"])), strings.TrimSpace(stringValue(action["mission_id"])))
	if targetID == "" {
		return ""
	}
	return targetKind + ":" + targetID
}

func schedulerMaintenancePriority(action map[string]any) int {
	switch strings.TrimSpace(stringValue(action["action_kind"])) {
	case "clear_invalid_bundle_governor":
		return 4
	case "refresh_stale_meeting":
		return 3
	case "resolve_stale_signal":
		return 2
	case "refresh_stale_run":
		return 1
	default:
		return 0
	}
}

func schedulerMaintenanceClass(action map[string]any) string {
	return fmt.Sprintf(
		"%d|%t|%s",
		schedulerMaintenancePriority(action),
		strings.TrimSpace(stringValue(action["bundle_id"])) != "",
		firstNonBlank(strings.TrimSpace(stringValue(action["target_kind"])), "unknown"),
	)
}

func sortSchedulerMaintenanceActions(actions []any, previousPrimaryMaintenanceKey string) []any {
	sorted := append([]any{}, actions...)
	sort.Slice(sorted, func(i, j int) bool {
		left, _ := toMap(sorted[i])
		right, _ := toMap(sorted[j])
		leftPriority := schedulerMaintenancePriority(left)
		rightPriority := schedulerMaintenancePriority(right)
		if leftPriority != rightPriority {
			return leftPriority > rightPriority
		}
		leftBundle := strings.TrimSpace(stringValue(left["bundle_id"])) != ""
		rightBundle := strings.TrimSpace(stringValue(right["bundle_id"])) != ""
		if leftBundle != rightBundle {
			return leftBundle
		}
		if previousPrimaryMaintenanceKey != "" {
			leftKey := schedulerMaintenanceKey(left)
			rightKey := schedulerMaintenanceKey(right)
			if leftKey == previousPrimaryMaintenanceKey && rightKey != previousPrimaryMaintenanceKey {
				return true
			}
			if rightKey == previousPrimaryMaintenanceKey && leftKey != previousPrimaryMaintenanceKey {
				return false
			}
		}
		leftMission := strings.TrimSpace(stringValue(left["mission_id"]))
		rightMission := strings.TrimSpace(stringValue(right["mission_id"]))
		if leftMission != rightMission {
			return leftMission < rightMission
		}
		leftTargetKind := firstNonBlank(strings.TrimSpace(stringValue(left["target_kind"])), "unknown")
		rightTargetKind := firstNonBlank(strings.TrimSpace(stringValue(right["target_kind"])), "unknown")
		if leftTargetKind != rightTargetKind {
			return leftTargetKind < rightTargetKind
		}
		leftTargetID := firstNonBlank(strings.TrimSpace(stringValue(left["target_id"])), strings.TrimSpace(stringValue(left["bundle_id"])))
		rightTargetID := firstNonBlank(strings.TrimSpace(stringValue(right["target_id"])), strings.TrimSpace(stringValue(right["bundle_id"])))
		return leftTargetID < rightTargetID
	})
	return sorted
}

func schedulerMaintenanceIsCritical(action map[string]any) bool {
	switch strings.TrimSpace(stringValue(action["action_kind"])) {
	case "clear_invalid_bundle_governor", "refresh_stale_meeting":
		return true
	default:
		return false
	}
}

func schedulerPrimaryControlPlanePacket(requiredHumanDecisions []any, primaryRequiredHumanDecision map[string]any, maintenanceActions []any, primaryMaintenanceAction map[string]any) (string, map[string]any, string, int) {
	if len(maintenanceActions) > 0 && schedulerMaintenanceIsCritical(primaryMaintenanceAction) {
		policy := "critical_maintenance_only"
		if len(requiredHumanDecisions) > 0 {
			policy = "critical_maintenance_over_decision"
		}
		packet, _ := toMap(cloneValue(primaryMaintenanceAction))
		return "maintenance", packet, policy, maxInt(len(maintenanceActions)-1, 0) + len(requiredHumanDecisions)
	}
	if len(requiredHumanDecisions) > 0 {
		policy := "human_decision_only"
		if len(maintenanceActions) > 0 {
			policy = "human_decision_over_noncritical_maintenance"
		}
		packet, _ := toMap(cloneValue(primaryRequiredHumanDecision))
		return "human-decision", packet, policy, maxInt(len(requiredHumanDecisions)-1, 0) + len(maintenanceActions)
	}
	if len(maintenanceActions) > 0 {
		packet, _ := toMap(cloneValue(primaryMaintenanceAction))
		return "maintenance", packet, "maintenance_only", maxInt(len(maintenanceActions)-1, 0)
	}
	return "", map[string]any{}, "none", 0
}

func sortSchedulerHumanDecisions(decisions []any, previousPrimaryDecisionKey string) []any {
	sorted := append([]any{}, decisions...)
	sort.Slice(sorted, func(i, j int) bool {
		left, _ := toMap(sorted[i])
		right, _ := toMap(sorted[j])
		if boolValue(left["active_decision_stale"]) != boolValue(right["active_decision_stale"]) {
			return boolValue(left["active_decision_stale"])
		}
		leftSurface := decisionSurfaceRank(stringValue(left["recommended_decision_surface"]))
		rightSurface := decisionSurfaceRank(stringValue(right["recommended_decision_surface"]))
		if leftSurface != rightSurface {
			return leftSurface > rightSurface
		}
		leftProgress := schedulerDecisionProgressRank(left)
		rightProgress := schedulerDecisionProgressRank(right)
		if leftProgress != rightProgress {
			return leftProgress > rightProgress
		}
		leftScope := decisionScopeRank(stringValue(left["decision_scope"]))
		rightScope := decisionScopeRank(stringValue(right["decision_scope"]))
		if leftScope != rightScope {
			return leftScope > rightScope
		}
		leftFollowUpUrgency := intValue(left["follow_up_signal_urgency_rank"])
		rightFollowUpUrgency := intValue(right["follow_up_signal_urgency_rank"])
		if leftFollowUpUrgency != rightFollowUpUrgency {
			return leftFollowUpUrgency > rightFollowUpUrgency
		}
		leftPriority := schedulerPriorityRank(stringValue(left["priority"]))
		rightPriority := schedulerPriorityRank(stringValue(right["priority"]))
		if leftPriority != rightPriority {
			return leftPriority > rightPriority
		}
		leftRisk := schedulerRiskRank(stringValue(left["effective_risk_level"]))
		rightRisk := schedulerRiskRank(stringValue(right["effective_risk_level"]))
		if leftRisk != rightRisk {
			return leftRisk > rightRisk
		}
		if boolValue(left["release_adjacent"]) != boolValue(right["release_adjacent"]) {
			return boolValue(left["release_adjacent"])
		}
		if previousPrimaryDecisionKey != "" {
			leftKey := schedulerDecisionKey(left)
			rightKey := schedulerDecisionKey(right)
			if leftKey == previousPrimaryDecisionKey && rightKey != previousPrimaryDecisionKey {
				return true
			}
			if rightKey == previousPrimaryDecisionKey && leftKey != previousPrimaryDecisionKey {
				return false
			}
		}
		leftUpdated := stringValue(left["updated_at"])
		rightUpdated := stringValue(right["updated_at"])
		if leftUpdated != rightUpdated {
			return leftUpdated > rightUpdated
		}
		leftKey := firstNonBlank(stringValue(left["bundle_id"]), stringValue(left["mission_id"]))
		rightKey := firstNonBlank(stringValue(right["bundle_id"]), stringValue(right["mission_id"]))
		return leftKey < rightKey
	})
	return sorted
}

func schedulerDecisionQueueMaps(decisions []any) (map[string]int, map[string]bool) {
	queuePositionByKey := map[string]int{}
	primaryByKey := map[string]bool{}
	for index, item := range decisions {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		scope := firstNonBlank(stringValue(decision["decision_scope"]), "mission")
		switch scope {
		case "bundle":
			bundleID := strings.TrimSpace(stringValue(decision["bundle_id"]))
			if bundleID == "" {
				continue
			}
			queuePositionByKey["bundle:"+bundleID] = index + 1
			if index == 0 {
				primaryByKey["bundle:"+bundleID] = true
			}
		default:
			missionID := strings.TrimSpace(stringValue(decision["mission_id"]))
			if missionID == "" {
				continue
			}
			queuePositionByKey["mission:"+missionID] = index + 1
			if index == 0 {
				primaryByKey["mission:"+missionID] = true
			}
		}
	}
	return queuePositionByKey, primaryByKey
}

func schedulerBatchClass(entry map[string]any) string {
	return fmt.Sprintf(
		"%d|%d|%t",
		schedulerPriorityRank(stringValue(entry["priority"])),
		schedulerRiskRank(firstNonBlank(stringValue(entry["effective_risk_level"]), stringValue(entry["risk_level"]))),
		boolValue(entry["release_adjacent"]),
	)
}

func stickyParallelBatchSelection(runnableMissionIDs []string, entryByMission map[string]map[string]any, previousSelectedMissionIDs []string, maxParallelRuns int) ([]string, bool) {
	if len(runnableMissionIDs) < 2 || len(previousSelectedMissionIDs) < 2 || maxParallelRuns < 2 {
		return nil, false
	}
	topClass := schedulerBatchClass(entryByMission[runnableMissionIDs[0]])
	selected := []string{}
	for _, missionID := range previousSelectedMissionIDs {
		if !containsString(runnableMissionIDs, missionID) {
			return nil, false
		}
		if schedulerBatchClass(entryByMission[missionID]) != topClass {
			return nil, false
		}
		if len(selected) < maxParallelRuns {
			selected = append(selected, missionID)
		}
	}
	if len(selected) < 2 {
		return nil, false
	}
	for _, missionID := range runnableMissionIDs {
		if len(selected) >= maxParallelRuns {
			break
		}
		if schedulerBatchClass(entryByMission[missionID]) != topClass {
			break
		}
		if !containsString(selected, missionID) {
			selected = append(selected, missionID)
		}
	}
	if len(selected) < 2 {
		return nil, false
	}
	return selected, true
}

func schedulerSummary(profile, missionRegistry, registry, signalLog map[string]any, includeInactive bool) map[string]any {
	now := nowTime().UTC().Format("20060102150405")
	generatedAt := nowISO()
	graphAuthority, _ := toMap(missionGraphAuthorityAudit(missionRegistry, registry, includeInactive))
	previousSchedulerSnapshot := latestSchedulerSnapshot()
	previousForegroundMissionID := ""
	previousSelectedMissionIDs := dedupeStrings(anyListToStrings(toList(previousSchedulerSnapshot["selected_mission_ids"])))
	if len(previousSelectedMissionIDs) > 0 {
		previousForegroundMissionID = previousSelectedMissionIDs[0]
	}
	previousPrimaryDecision, _ := toMap(previousSchedulerSnapshot["primary_required_human_decision"])
	previousPrimaryDecisionKey := schedulerDecisionKey(previousPrimaryDecision)
	previousPrimaryMaintenance, _ := toMap(previousSchedulerSnapshot["primary_maintenance_action"])
	previousPrimaryMaintenanceKey := schedulerMaintenanceKey(previousPrimaryMaintenance)
	maxParallelRuns := intValue(getKey(profile, "profile.max_parallel_runs", 0))
	if maxParallelRuns <= 0 {
		maxParallelRuns = 1
	}
	requireExplicitParallelJustification := boolValue(policyValue("execution.require_explicit_parallel_justification", true))
	requirePrimaryHumanDecisionBeforeNewForeground := boolValue(policyValue("execution.require_primary_human_decision_before_new_foreground", true))
	requirePrimaryMaintenanceBeforeNewForeground := boolValue(policyValue("execution.require_primary_maintenance_before_new_foreground", true))

	bundleInfoByMission := map[string]map[string]any{}
	conflictMissionSet := map[string]bool{}
	escalatedBundleMissionSet := map[string]bool{}
	governedLoserMissionSet := map[string]bool{}
	ambiguousGovernedMissionSet := map[string]bool{}
	for _, item := range toList(graphAuthority["bundles"]) {
		bundle, _ := toMap(item)
		bundleID := strings.TrimSpace(stringValue(bundle["bundle_id"]))
		bundleState := strings.TrimSpace(stringValue(bundle["bundle_state"]))
		resolutionKind := strings.TrimSpace(stringValue(bundle["resolution_kind"]))
		governingMissionID := strings.TrimSpace(stringValue(bundle["governing_mission_id"]))
		requiredAction := strings.TrimSpace(stringValue(bundle["required_action"]))
		missionIDs := dedupeStrings(anyListToStrings(toList(bundle["mission_ids"])))
		for _, missionID := range missionIDs {
			role := "member"
			switch {
			case bundleState == "governed" && governingMissionID != "" && missionID == governingMissionID:
				role = "governor"
			case bundleState == "governed" && governingMissionID != "" && missionID != governingMissionID:
				role = "governed"
			case bundleState == "governed":
				role = "governed-ambiguous"
			case bundleState == "escalation-required":
				role = "escalation-required"
			case bundleState == "conflicted":
				role = "conflicted"
			}
			bundleInfoByMission[missionID] = map[string]any{
				"bundle_id":                   bundleID,
				"bundle_state":                bundleState,
				"bundle_role":                 role,
				"resolution_kind":             resolutionKind,
				"governing_mission_id":        governingMissionID,
				"governance_source_kind":      stringValue(bundle["governance_source_kind"]),
				"adopted_from_meeting_id":     stringValue(bundle["adopted_from_meeting_id"]),
				"required_action":             requiredAction,
				"bundle_mission_ids":          stringListToAny(missionIDs),
				"candidate_governor_ids":      cloneValue(getKey(bundle, "candidate_governor_ids", []any{})),
				"candidate_governor_packets":  cloneValue(getKey(bundle, "candidate_governor_packets", []any{})),
				"explicit_governor_signal_id": stringValue(bundle["explicit_governor_signal_id"]),
				"governed_mission_ids":        cloneValue(bundle["governed_mission_ids"]),
				"losing_mission_ids":          cloneValue(bundle["losing_mission_ids"]),
				"authority_basis":             cloneValue(bundle["authority_basis"]),
				"selector_refs":               cloneValue(bundle["selector_refs"]),
				"governance_notes":            cloneValue(bundle["governance_notes"]),
			}
			switch bundleState {
			case "conflicted":
				conflictMissionSet[missionID] = true
			case "escalation-required":
				escalatedBundleMissionSet[missionID] = true
			case "governed":
				if governingMissionID == "" {
					ambiguousGovernedMissionSet[missionID] = true
				} else if missionID != governingMissionID {
					governedLoserMissionSet[missionID] = true
				}
			}
		}
	}

	runMissionByID := map[string]string{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		runID := strings.TrimSpace(stringValue(run["run_id"]))
		missionID := strings.TrimSpace(stringValue(run["mission_id"]))
		if runID != "" && missionID != "" {
			runMissionByID[runID] = missionID
		}
	}

	maintenanceActions, maintenanceTally := schedulerMaintenanceActions(missionRegistry, registry, signalLog)
	maintenanceActions = sortSchedulerMaintenanceActions(maintenanceActions, previousPrimaryMaintenanceKey)

	activeMeetingCountByMission, openMeetingActionCountByMission, activeMeetingIDsByDecisionID, activeMeetingIDsByBundleID := meetingActivityByMission(registry)

	for _, item := range toList(graphAuthority["cross_mission_conflicts"]) {
		conflict, _ := toMap(item)
		for _, missionID := range anyListToStrings(toList(conflict["mission_ids"])) {
			if strings.TrimSpace(missionID) != "" {
				conflictMissionSet[missionID] = true
			}
		}
	}

	followUpSignalsByMission := map[string][]any{}
	staleFollowUpSignalsByMission := map[string][]any{}
	freshFollowUpSignalKindsByMission := map[string][]string{}
	for _, item := range toList(signalLog["signals"]) {
		signal, _ := toMap(item)
		if !activeSignal(signal) {
			continue
		}
		kind := normalizeSignalKind(stringValue(signal["kind"]))
		if kind != "requires_meeting" &&
			kind != "interface_change" &&
			kind != "dependency_discovered" &&
			kind != "blocked" &&
			kind != "review_feedback" &&
			kind != "priority_change" &&
			kind != "risk_raised" &&
			kind != "requires_rebase" {
			continue
		}
		missionIDs := []string{}
		if missionID := strings.TrimSpace(stringValue(signal["mission_id"])); missionID != "" {
			missionIDs = append(missionIDs, missionID)
		}
		for _, ref := range anyListToStrings(toList(signal["applies_to"])) {
			ref = strings.TrimSpace(ref)
			if strings.HasPrefix(ref, "issue-") {
				missionIDs = append(missionIDs, ref)
			}
		}
		entry := map[string]any{
			"signal_id":       stringValue(signal["signal_id"]),
			"kind":            kind,
			"mission_id":      firstNonBlank(strings.TrimSpace(stringValue(signal["mission_id"])), ""),
			"run_id":          stringValue(signal["run_id"]),
			"authority_level": stringValue(signal["authority_level"]),
			"summary":         stringValue(signal["summary"]),
			"fresh_until":     stringValue(signal["fresh_until"]),
			"resolution_note": stringValue(signal["resolution_note"]),
			"source_type":     stringValue(signal["source_type"]),
			"source_id":       stringValue(signal["source_id"]),
		}
		isFresh := true
		if freshUntil, err := parseISO(stringValue(signal["fresh_until"])); err == nil && freshUntil != nil && freshUntil.Before(nowTime()) {
			isFresh = false
		}
		for _, missionID := range dedupeStrings(missionIDs) {
			if missionID == "" {
				continue
			}
			if isFresh {
				followUpSignalsByMission[missionID] = append(followUpSignalsByMission[missionID], entry)
				freshFollowUpSignalKindsByMission[missionID] = append(freshFollowUpSignalKindsByMission[missionID], kind)
			} else {
				staleFollowUpSignalsByMission[missionID] = append(staleFollowUpSignalsByMission[missionID], entry)
			}
		}
	}

	entryByMission := map[string]map[string]any{}
	missionOrder := []string{}
	rescheduleActions := []any{}
	rawRequiredHumanDecisions := []any{}

	for _, item := range toList(missionRegistry["missions"]) {
		mission, _ := toMap(item)
		missionID := strings.TrimSpace(stringValue(mission["mission_id"]))
		if missionID == "" {
			continue
		}
		if !includeInactive && !activeMission(mission) {
			continue
		}
		runs := missionRuns(registry, missionID)
		authoritativeRun := authoritativeMissionRun(runs)
		activeRunIDs := []string{}
		activeTeammateCount := 0
		releaseAdjacent := false
		reviewBlocked := false
		checksFailed := false
		mergeReadyPresent := false
		for _, run := range runs {
			if activeRun(run) {
				activeRunIDs = append(activeRunIDs, stringValue(run["run_id"]))
				if stringValue(run["owner_type"]) == "teammate" {
					activeTeammateCount++
				}
			}
			shipStatus := stringValue(run["ship_status"])
			reviewStatus := stringValue(run["review_status"])
			checksStatus := stringValue(run["checks_status"])
			if stringValue(run["pr_ref"]) != "" ||
				stringValue(run["rc_id"]) != "" ||
				(reviewStatus != "" && reviewStatus != "unknown") ||
				(checksStatus != "" && checksStatus != "unknown") ||
				(shipStatus != "" && shipStatus != "unstarted") {
				releaseAdjacent = true
			}
			if reviewStatus == "pending" || reviewStatus == "changes_requested" {
				reviewBlocked = true
			}
			if checksStatus == "failed" {
				checksFailed = true
			}
			if mergeReady(run) {
				mergeReadyPresent = true
			}
		}

		classification := "runnable"
		reasons := []string{}
		recommendedAction := firstNonBlank(stringValue(mission["next_expected_action"]), "Review mission state and continue governed execution.")
		decisionKind := ""
		unsatisfiedMissionDeps := missionDependenciesUnsatisfied(mission, missionRegistry)
		unsatisfiedRunDeps := []string{}
		authoritativeRunID := ""
		authoritativeRunState := ""
		authoritativeRunMode := ""
		authoritativeRunSurface := ""
		adoptedOutputPacketKind := ""
		adoptedOutputFollowUpPacketCount := 0
		adoptedOutputOpenFollowUpPacketCount := 0
		adoptedPrimaryFollowUpPacket := map[string]any{}
		adoptedPrimaryFollowUpAction := ""
		updatedAt := firstNonBlank(stringValue(mission["updated_at"]), stringValue(mission["created_at"]))
		bundleInfo, _ := toMap(bundleInfoByMission[missionID])
		bundleID := stringValue(bundleInfo["bundle_id"])
		bundleState := stringValue(bundleInfo["bundle_state"])
		bundleRole := stringValue(bundleInfo["bundle_role"])
		bundleResolutionKind := stringValue(bundleInfo["resolution_kind"])
		bundleGovernorID := stringValue(bundleInfo["governing_mission_id"])
		bundleGovernanceSourceKind := stringValue(bundleInfo["governance_source_kind"])
		bundleAdoptedFromMeetingID := stringValue(bundleInfo["adopted_from_meeting_id"])
		bundleRequiredAction := stringValue(bundleInfo["required_action"])
		bundleCandidateGovernorIDs := anyListToStrings(toList(bundleInfo["candidate_governor_ids"]))
		bundleCandidateGovernorPackets := cloneValue(getKey(bundleInfo, "candidate_governor_packets", []any{}))
		if authoritativeRun != nil {
			authoritativeRunID = stringValue(authoritativeRun["run_id"])
			authoritativeRunState = stringValue(authoritativeRun["state"])
			authoritativeRunMode = stringValue(authoritativeRun["mode"])
			authoritativeRunSurface = stringValue(authoritativeRun["execution_surface"])
			adoptedOutputPacketKind = strings.TrimSpace(stringValue(authoritativeRun["adopted_output_packet_kind"]))
			adoptedOutputFollowUpPacketCount = len(runAdoptedFollowUpPackets(authoritativeRun))
			adoptedOutputOpenFollowUpPacketCount = runOpenAdoptedFollowUpPacketCount(authoritativeRun)
			if packet := runPrimaryAdoptedFollowUpPacket(authoritativeRun); packet != nil {
				adoptedPrimaryFollowUpPacket, _ = toMap(cloneValue(packet))
				adoptedPrimaryFollowUpAction = adoptedFollowUpActionText(packet)
			}
			updatedAt = firstNonBlank(runSortTimestamp(authoritativeRun), updatedAt)
			unsatisfiedRunDeps = unsatisfiedDependencies(authoritativeRun, registry)
		}
		parallelJustification := firstNonBlank(stringValue(mission["parallel_justification"]), stringValue(getKey(authoritativeRun, "parallel_justification", "")))
		parallelJustified := strings.TrimSpace(parallelJustification) != ""
		ownershipRefs := ownershipClaims(mission)
		if len(ownershipRefs) == 0 && authoritativeRun != nil {
			ownershipRefs = ownershipClaims(authoritativeRun)
		}
		profileHighRiskScopeRefs := matchedProfileHighRiskScopes(profile, ownershipRefs, releaseAdjacent)
		missionAudit, _ := toMap(missionAuthorityAudit(mission, registry))
		effectiveRiskLevel := stringValue(mission["risk_level"])
		if authoritativeRun != nil && schedulerRiskRank(effectiveRiskLevel) == 0 {
			effectiveRiskLevel = firstNonBlank(effectiveRiskLevel, stringValue(authoritativeRun["risk_level"]))
		}
		freshFollowUpSignalKinds := dedupeStrings(freshFollowUpSignalKindsByMission[missionID])
		followUpSignalUrgencyRank := schedulerFollowUpSignalUrgencyRank(freshFollowUpSignalKinds)
		primaryFollowUpSignalKind := schedulerPrimaryFollowUpSignalKind(freshFollowUpSignalKinds)
		if containsString(freshFollowUpSignalKinds, "risk_raised") && schedulerRiskRank(effectiveRiskLevel) < schedulerRiskRank("high") {
			effectiveRiskLevel = "high"
		}
		if len(profileHighRiskScopeRefs) > 0 && schedulerRiskRank(effectiveRiskLevel) < schedulerRiskRank("high") {
			effectiveRiskLevel = "high"
		}
		meetingTriggerRefs := []string{}
		recommendedDecisionSurface := ""
		recommendedDecisionProtocolValue := ""
		recommendedProtocolReason := ""

		if !activeMission(mission) {
			classification = "completed"
			reasons = append(reasons, "Mission is no longer active in the control plane.")
			recommendedAction = "No active scheduling action required."
		} else {
			switch {
			case len(unsatisfiedMissionDeps) > 0 || len(unsatisfiedRunDeps) > 0:
				classification = "blocked-by-dependency"
				reasons = append(reasons, fmt.Sprintf("Mission dependencies are not yet satisfied: %s", strings.Join(dedupeStrings(append(unsatisfiedMissionDeps, unsatisfiedRunDeps...)), ", ")))
				recommendedAction = "Wait for dependency missions/runs to satisfy their desired state before continuing."
				if authoritativeRunID != "" && authoritativeRunState != "waiting" {
					rescheduleActions = append(rescheduleActions, map[string]any{
						"mission_id":      missionID,
						"run_id":          authoritativeRunID,
						"action_kind":     "move_to_waiting",
						"target_state":    "waiting",
						"current_state":   authoritativeRunState,
						"reason":          "Run is still foregrounded even though scheduler sees unsatisfied dependencies.",
						"dependency_refs": stringListToAny(dedupeStrings(append(unsatisfiedMissionDeps, unsatisfiedRunDeps...))),
					})
				}
			case ambiguousGovernedMissionSet[missionID]:
				classification = "needs-human-decision"
				reasons = append(reasons, fmt.Sprintf("Mission is in governed bundle %s, but the control plane does not yet have a single governing mission.", bundleID))
				recommendedAction = firstNonBlank(bundleRequiredAction, "Inspect bundle authority before widening execution.")
				decisionKind = "inspect_bundle_authority"
			case escalatedBundleMissionSet[missionID]:
				classification = "needs-human-decision"
				reasons = append(reasons, fmt.Sprintf("Mission is in escalation-required bundle %s.", bundleID))
				recommendedAction = firstNonBlank(bundleRequiredAction, "Open a council or add an explicit dependency/ownership transfer before widening execution.")
				decisionKind = "establish_bundle_governor"
			case governedLoserMissionSet[missionID]:
				classification = "deferred"
				reasons = append(reasons, fmt.Sprintf("Mission is currently governed by %s inside bundle %s.", firstNonBlank(bundleGovernorID, "another mission"), bundleID))
				recommendedAction = firstNonBlank(bundleRequiredAction, fmt.Sprintf("Defer this mission behind bundle governor %s until bundle authority changes.", firstNonBlank(bundleGovernorID, "the governing mission")))
				if authoritativeRunID != "" && (authoritativeRunState == "ready" || authoritativeRunState == "executing") {
					rescheduleActions = append(rescheduleActions, map[string]any{
						"mission_id":    missionID,
						"run_id":        authoritativeRunID,
						"action_kind":   "defer_until_bundle_governor_clears",
						"target_state":  "waiting",
						"current_state": authoritativeRunState,
						"reason":        fmt.Sprintf("Scheduler is deferring this mission behind bundle governor %s.", firstNonBlank(bundleGovernorID, "the governing mission")),
					})
				}
			case conflictMissionSet[missionID]:
				classification = "blocked-by-overlap"
				reasons = append(reasons, "Mission is in a live cross-mission ownership conflict bundle.")
				recommendedAction = "Narrow selectors, serialize missions, or add a governing dependency/council before execution resumes."
				decisionKind = "resolve_overlap_conflict"
			case containsString(freshFollowUpSignalKinds, "blocked"):
				classification = "blocked-by-policy"
				reasons = append(reasons, "Mission has a fresh blocked signal that has not yet been adopted into runtime state.")
				recommendedAction = "Adopt or resolve the fresh blocker signal before continuing."
				if authoritativeRunID != "" && (authoritativeRunState == "ready" || authoritativeRunState == "executing") {
					rescheduleActions = append(rescheduleActions, map[string]any{
						"mission_id":    missionID,
						"run_id":        authoritativeRunID,
						"action_kind":   "pause_for_fresh_blocker_signal",
						"target_state":  "waiting",
						"current_state": authoritativeRunState,
						"reason":        "Scheduler detected a fresh blocker signal that has not yet been adopted into run state.",
					})
				}
			case containsString(freshFollowUpSignalKinds, "requires_rebase"):
				classification = "blocked-by-policy"
				reasons = append(reasons, "Mission has a fresh requires-rebase signal that should be handled before execution continues.")
				recommendedAction = "Rebase or refresh the governed run state before continuing this mission."
				if authoritativeRunID != "" && (authoritativeRunState == "ready" || authoritativeRunState == "executing") {
					rescheduleActions = append(rescheduleActions, map[string]any{
						"mission_id":    missionID,
						"run_id":        authoritativeRunID,
						"action_kind":   "pause_for_fresh_rebase_signal",
						"target_state":  "waiting",
						"current_state": authoritativeRunState,
						"reason":        "Scheduler detected a fresh requires-rebase signal that has not yet been adopted into runtime state.",
					})
				}
			case containsString(freshFollowUpSignalKinds, "dependency_discovered"):
				classification = "blocked-by-dependency"
				reasons = append(reasons, "Mission has a fresh dependency-discovered signal that has not yet been converted into dependency state.")
				recommendedAction = "Adopt the fresh dependency signal before continuing this mission."
				if authoritativeRunID != "" && (authoritativeRunState == "ready" || authoritativeRunState == "executing") {
					rescheduleActions = append(rescheduleActions, map[string]any{
						"mission_id":    missionID,
						"run_id":        authoritativeRunID,
						"action_kind":   "pause_for_fresh_dependency_signal",
						"target_state":  "waiting",
						"current_state": authoritativeRunState,
						"reason":        "Scheduler detected a fresh dependency-discovered signal that has not yet been adopted into mission state.",
					})
				}
			case containsString(freshFollowUpSignalKinds, "requires_meeting") ||
				containsString(freshFollowUpSignalKinds, "interface_change") ||
				containsString(freshFollowUpSignalKinds, "review_feedback") ||
				containsString(freshFollowUpSignalKinds, "priority_change") ||
				containsString(freshFollowUpSignalKinds, "risk_raised"):
				classification = "needs-human-decision"
				reasons = append(reasons, "Mission has fresh follow-up control-plane signals that should be reviewed before widening execution.")
				switch primaryFollowUpSignalKind {
				case "risk_raised":
					recommendedAction = "Review the fresh risk-raised signal and either tighten policy or convene a governed decision before widening execution."
				case "review_feedback":
					recommendedAction = "Review the fresh review feedback and adopt it into governed review state before widening execution."
				case "priority_change":
					recommendedAction = "Review the fresh priority-change signal and adopt it into mission state before widening execution."
				default:
					recommendedAction = "Review the fresh coordination signals and either adopt them into state or materialize the required meeting."
				}
				decisionKind = "review_follow_up_signals"
				if authoritativeRunID != "" && (authoritativeRunState == "ready" || authoritativeRunState == "executing") {
					rescheduleActions = append(rescheduleActions, map[string]any{
						"mission_id":    missionID,
						"run_id":        authoritativeRunID,
						"action_kind":   "pause_for_fresh_follow_up_signal",
						"target_state":  "waiting",
						"current_state": authoritativeRunState,
						"reason":        "Scheduler detected fresh coordination signals that have not yet been adopted into runtime state.",
					})
				}
			case intValue(missionAudit["unresolved_conflict_count"]) > 0:
				classification = "needs-human-decision"
				reasons = append(reasons, "Mission has unresolved competing root-run authority.")
				recommendedAction = "Close, merge, or explicitly govern the competing root run before continuing."
				decisionKind = "resolve_parallel_root_runs"
			case openMeetingActionCountByMission[missionID] > 0:
				classification = "needs-human-decision"
				reasons = append(reasons, fmt.Sprintf("Mission has %d open meeting action(s).", openMeetingActionCountByMission[missionID]))
				recommendedAction = "Review open meeting actions and translate them into accepted execution state."
				decisionKind = "adopt_meeting_actions"
			case bundleID != "" && len(activeMeetingIDsByBundleID[bundleID]) > 0:
				classification = "needs-human-decision"
				reasons = append(reasons, fmt.Sprintf("Bundle %s already has active council coverage: %s.", bundleID, strings.Join(dedupeStrings(activeMeetingIDsByBundleID[bundleID]), ", ")))
				recommendedAction = fmt.Sprintf("Review the active bundle council(s) for %s before widening execution: %s.", bundleID, strings.Join(dedupeStrings(activeMeetingIDsByBundleID[bundleID]), ", "))
			case activeMeetingCountByMission[missionID] > 0:
				classification = "needs-human-decision"
				reasons = append(reasons, fmt.Sprintf("Mission has %d active meeting(s).", activeMeetingCountByMission[missionID]))
				recommendedAction = "Review the active council before widening execution."
				decisionKind = "review_active_council"
			case stringValue(mission["state"]) == "review" || reviewBlocked || checksFailed || mergeReadyPresent:
				classification = "needs-human-decision"
				reasons = append(reasons, "Mission is in review/checks/merge-gate handling.")
				recommendedAction = "Resolve review and checks or merge the ready PR before continuing other work."
				decisionKind = "review_gate_decision"
			case stringValue(mission["state"]) == "blocked" || authoritativeRunState == "blocked":
				classification = "blocked-by-policy"
				reasons = append(reasons, firstNonBlank(stringValue(mission["status_reason"]), stringValue(getKey(authoritativeRun, "status_reason", "")), "Mission is blocked by current policy or runtime state."))
				recommendedAction = "Clear the blocking policy, signal, or state before continuing."
			case authoritativeRunState == "waiting" || authoritativeRunState == "paused":
				classification = "deferred"
				reasons = append(reasons, fmt.Sprintf("Authoritative run is currently %s.", authoritativeRunState))
				recommendedAction = firstNonBlank(stringValue(getKey(authoritativeRun, "next_expected_action", "")), recommendedAction)
			default:
				classification = "runnable"
				reasons = append(reasons, "Mission is clear to run under current dependency, overlap, and authority policy.")
				if bundleState == "governed" && bundleGovernorID == missionID {
					reasons = append(reasons, fmt.Sprintf("Mission currently governs bundle %s under %s.", bundleID, bundleResolutionKind))
					recommendedAction = firstNonBlank(bundleRequiredAction, fmt.Sprintf("Keep bundle authority anchored to %s while governed missions stay bounded.", missionID))
				}
			}
		}
		if classification == "runnable" && adoptedOutputOpenFollowUpPacketCount > 0 {
			reasons = append(reasons, fmt.Sprintf("Mission carries %d adopted council follow-up packet(s) in run state.", adoptedOutputOpenFollowUpPacketCount))
			recommendedAction = firstNonBlank(
				adoptedPrimaryFollowUpAction,
				recommendedAction,
			)
		}
		if decisionKind != "" {
			meetingTriggerRefs = matchedMeetingTriggerRefs(profile, decisionKind, bundleState, intValue(missionAudit["unresolved_conflict_count"]), releaseAdjacent, profileHighRiskScopeRefs, freshFollowUpSignalKinds)
			switch {
			case activeMeetingCountByMission[missionID] > 0:
				recommendedDecisionSurface = "active-council"
			case len(meetingTriggerRefs) > 0:
				recommendedDecisionSurface = "council"
			default:
				recommendedDecisionSurface = "leader-decision"
			}
			recommendedDecisionProtocolValue, recommendedProtocolReason = recommendedDecisionProtocol(
				decisionKind,
				recommendedDecisionSurface,
				releaseAdjacent,
				effectiveRiskLevel,
				meetingTriggerRefs,
				freshFollowUpSignalKinds,
			)
		}
		if len(profileHighRiskScopeRefs) > 0 {
			reasons = append(reasons, fmt.Sprintf("Mission touches profile high-risk scopes: %s.", strings.Join(profileHighRiskScopeRefs, ", ")))
		}

		if decisionKind != "" {
			rawRequiredHumanDecisions = append(rawRequiredHumanDecisions, map[string]any{
				"decision_scope":                 "mission",
				"mission_id":                     missionID,
				"bundle_id":                      bundleID,
				"bundle_state":                   bundleState,
				"bundle_resolution_kind":         bundleResolutionKind,
				"bundle_governing_mission_id":    bundleGovernorID,
				"bundle_governance_source_kind":  bundleGovernanceSourceKind,
				"bundle_adopted_from_meeting_id": bundleAdoptedFromMeetingID,
				"bundle_required_action":         bundleRequiredAction,
				"bundle_mission_ids":             cloneValue(getKey(bundleInfo, "bundle_mission_ids", []any{})),
				"candidate_governor_ids":         stringListToAny(bundleCandidateGovernorIDs),
				"candidate_governor_packets":     cloneValue(bundleCandidateGovernorPackets),
				"decision_kind":                  decisionKind,
				"recommended_action":             recommendedAction,
				"recommended_decision_surface":   recommendedDecisionSurface,
				"recommended_decision_protocol":  recommendedDecisionProtocolValue,
				"recommended_protocol_reason":    recommendedProtocolReason,
				"meeting_trigger_refs":           stringListToAny(dedupeStrings(meetingTriggerRefs)),
				"fresh_follow_up_signal_kinds":   stringListToAny(freshFollowUpSignalKinds),
				"fresh_follow_up_signal_count":   len(followUpSignalsByMission[missionID]),
				"follow_up_signal_urgency_rank":  followUpSignalUrgencyRank,
				"note":                           firstNonBlank(strings.Join(dedupeStrings(reasons), " "), recommendedAction),
				"priority":                       stringValue(mission["priority"]),
				"effective_risk_level":           effectiveRiskLevel,
				"release_adjacent":               releaseAdjacent,
				"updated_at":                     updatedAt,
			})
		}

		entry := map[string]any{
			"mission_id":                                 missionID,
			"mission_state":                              stringValue(mission["state"]),
			"classification":                             classification,
			"priority":                                   stringValue(mission["priority"]),
			"risk_level":                                 stringValue(mission["risk_level"]),
			"parallel_justified":                         parallelJustified,
			"parallel_justification":                     parallelJustification,
			"effective_risk_level":                       effectiveRiskLevel,
			"profile_high_risk":                          len(profileHighRiskScopeRefs) > 0,
			"profile_high_risk_scope_refs":               stringListToAny(profileHighRiskScopeRefs),
			"recommended_decision_surface":               recommendedDecisionSurface,
			"recommended_decision_protocol":              recommendedDecisionProtocolValue,
			"recommended_protocol_reason":                recommendedProtocolReason,
			"meeting_trigger_refs":                       stringListToAny(meetingTriggerRefs),
			"release_adjacent":                           releaseAdjacent,
			"authoritative_run_id":                       authoritativeRunID,
			"authoritative_run_state":                    authoritativeRunState,
			"authoritative_run_mode":                     authoritativeRunMode,
			"authoritative_execution_surface":            authoritativeRunSurface,
			"adopted_output_packet_kind":                 adoptedOutputPacketKind,
			"adopted_output_follow_up_packet_count":      adoptedOutputFollowUpPacketCount,
			"adopted_output_open_follow_up_packet_count": adoptedOutputOpenFollowUpPacketCount,
			"adopted_primary_follow_up_packet":           cloneValue(adoptedPrimaryFollowUpPacket),
			"adopted_primary_follow_up_action":           adoptedPrimaryFollowUpAction,
			"active_run_ids":                             stringListToAny(dedupeStrings(activeRunIDs)),
			"active_run_count":                           len(dedupeStrings(activeRunIDs)),
			"teammate_run_count":                         activeTeammateCount,
			"unsatisfied_mission_dependencies":           stringListToAny(unsatisfiedMissionDeps),
			"unsatisfied_run_dependencies":               stringListToAny(unsatisfiedRunDeps),
			"meeting_count":                              activeMeetingCountByMission[missionID],
			"open_meeting_action_count":                  openMeetingActionCountByMission[missionID],
			"follow_up_signal_count":                     len(followUpSignalsByMission[missionID]),
			"follow_up_signals":                          cloneValue(followUpSignalsByMission[missionID]),
			"fresh_follow_up_signal_kinds":               stringListToAny(freshFollowUpSignalKinds),
			"follow_up_signal_urgency_rank":              followUpSignalUrgencyRank,
			"stale_follow_up_signal_count":               len(staleFollowUpSignalsByMission[missionID]),
			"stale_follow_up_signals":                    cloneValue(staleFollowUpSignalsByMission[missionID]),
			"bundle_id":                                  bundleID,
			"bundle_state":                               bundleState,
			"bundle_role":                                bundleRole,
			"bundle_resolution_kind":                     bundleResolutionKind,
			"bundle_governing_mission_id":                bundleGovernorID,
			"bundle_governance_source_kind":              bundleGovernanceSourceKind,
			"bundle_adopted_from_meeting_id":             bundleAdoptedFromMeetingID,
			"bundle_required_action":                     bundleRequiredAction,
			"candidate_governor_ids":                     stringListToAny(bundleCandidateGovernorIDs),
			"candidate_governor_packets":                 cloneValue(bundleCandidateGovernorPackets),
			"latest_summary":                             stringValue(mission["latest_summary"]),
			"recommended_action":                         recommendedAction,
			"reasons":                                    stringListToAny(dedupeStrings(reasons)),
			"updated_at":                                 updatedAt,
		}
		entryByMission[missionID] = entry
		missionOrder = append(missionOrder, missionID)
	}
	requiredHumanDecisions := sortSchedulerHumanDecisions(
		promoteMaterializedDecisionPackets(
			recoverActiveBundleDecisionPackets(aggregateSchedulerHumanDecisions(rawRequiredHumanDecisions)),
			activeMeetingIDsByDecisionID,
		),
		previousPrimaryDecisionKey,
	)
	decisionQueuePositionByKey, primaryDecisionByKey := schedulerDecisionQueueMaps(requiredHumanDecisions)
	decisionIDByKey := map[string]string{}
	decisionPacketByKey := map[string]map[string]any{}
	for _, item := range requiredHumanDecisions {
		decision, ok := toMap(item)
		if !ok {
			continue
		}
		if key := schedulerDecisionKey(decision); key != "" {
			decisionIDByKey[key] = stringValue(decision["decision_id"])
			decisionPacketByKey[key] = decision
		}
	}

	for _, missionID := range missionOrder {
		entry := entryByMission[missionID]
		queueKey := "mission:" + missionID
		if bundleID := stringValue(entry["bundle_id"]); bundleID != "" && stringValue(entry["bundle_state"]) == "escalation-required" {
			queueKey = "bundle:" + bundleID
		}
		if queuePosition, ok := decisionQueuePositionByKey[queueKey]; ok {
			entry["human_decision_queue_position"] = queuePosition
			entry["primary_human_decision"] = primaryDecisionByKey[queueKey]
			entry["human_decision_id"] = decisionIDByKey[queueKey]
			if packet, ok := decisionPacketByKey[queueKey]; ok {
				entry["active_decision_meeting_ids"] = cloneValue(getKey(packet, "active_meeting_ids", []any{}))
				entry["active_decision_meeting_count"] = intValue(packet["active_meeting_count"])
				if surface := strings.TrimSpace(stringValue(packet["recommended_decision_surface"])); surface != "" {
					entry["recommended_decision_surface"] = surface
				}
			}
			if bundleID := strings.TrimSpace(stringValue(entry["bundle_id"])); bundleID != "" {
				activeBundleMeetingIDs := dedupeStrings(activeMeetingIDsByBundleID[bundleID])
				entry["bundle_active_decision_meeting_ids"] = stringListToAny(activeBundleMeetingIDs)
				entry["bundle_active_decision_meeting_count"] = len(activeBundleMeetingIDs)
				if len(activeBundleMeetingIDs) > 0 {
					reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Bundle %s already has active council coverage: %s.", bundleID, strings.Join(activeBundleMeetingIDs, ", ")))
					entry["reasons"] = stringListToAny(dedupeStrings(reasons))
					entry["recommended_action"] = firstNonBlank(
						stringValue(getKey(decisionPacketByKey[queueKey], "recommended_action", "")),
						fmt.Sprintf("Review the active bundle council(s) for %s before widening execution: %s.", bundleID, strings.Join(activeBundleMeetingIDs, ", ")),
					)
					entry["recommended_decision_surface"] = "active-council"
					entry["bundle_control_state"] = "active-council"
					if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
						currentState := stringValue(entry["authoritative_run_state"])
						if currentState == "ready" || currentState == "executing" {
							holdReason := firstNonBlank(
								stringValue(getKey(decisionPacketByKey[queueKey], "recommended_action", "")),
								stringValue(getKey(decisionPacketByKey[queueKey], "active_decision_next_requirement", "")),
								fmt.Sprintf("Scheduler is holding this mission behind the active bundle council for %s.", bundleID),
							)
							rescheduleActions = append(rescheduleActions, map[string]any{
								"mission_id":                  missionID,
								"run_id":                      runID,
								"action_kind":                 "hold_for_active_bundle_council",
								"target_state":                "waiting",
								"current_state":               currentState,
								"bundle_id":                   bundleID,
								"active_decision_meeting_ids": stringListToAny(activeBundleMeetingIDs),
								"reason":                      holdReason,
							})
						}
					}
				}
			}
			if queuePosition > 1 && stringValue(entry["classification"]) == "needs-human-decision" {
				reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Another governed human decision currently outranks this mission in the scheduler portfolio (queue position %d).", queuePosition))
				entry["reasons"] = stringListToAny(dedupeStrings(reasons))
				entry["recommended_action"] = firstNonBlank(
					stringValue(entry["recommended_action"]),
					fmt.Sprintf("Hold this mission behind the current primary human decision until queue position %d is reached.", queuePosition),
				)
				if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
					currentState := stringValue(entry["authoritative_run_state"])
					if currentState == "ready" || currentState == "executing" {
						rescheduleActions = append(rescheduleActions, map[string]any{
							"mission_id":              missionID,
							"run_id":                  runID,
							"action_kind":             "defer_pending_primary_human_decision",
							"target_state":            "waiting",
							"current_state":           currentState,
							"human_decision_position": queuePosition,
							"decision_id":             decisionIDByKey[queueKey],
							"reason":                  "Scheduler is deferring this mission behind a higher-priority human decision in the portfolio queue.",
						})
					}
				}
			}
		} else {
			entry["human_decision_queue_position"] = 0
			entry["primary_human_decision"] = false
			entry["human_decision_id"] = ""
			entry["active_decision_meeting_ids"] = []any{}
			entry["active_decision_meeting_count"] = 0
			if bundleID := strings.TrimSpace(stringValue(entry["bundle_id"])); bundleID != "" {
				activeBundleMeetingIDs := dedupeStrings(activeMeetingIDsByBundleID[bundleID])
				entry["bundle_active_decision_meeting_ids"] = stringListToAny(activeBundleMeetingIDs)
				entry["bundle_active_decision_meeting_count"] = len(activeBundleMeetingIDs)
				if len(activeBundleMeetingIDs) > 0 {
					entry["bundle_control_state"] = "active-council"
					reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Bundle %s already has active council coverage: %s.", bundleID, strings.Join(activeBundleMeetingIDs, ", ")))
					entry["reasons"] = stringListToAny(dedupeStrings(reasons))
					entry["recommended_action"] = firstNonBlank(
						stringValue(getKey(decisionPacketByKey["bundle:"+bundleID], "recommended_action", "")),
						fmt.Sprintf("Review the active bundle council(s) for %s before widening execution: %s.", bundleID, strings.Join(activeBundleMeetingIDs, ", ")),
					)
					if surface := strings.TrimSpace(stringValue(entry["recommended_decision_surface"])); surface == "council" || surface == "leader-decision" || surface == "" {
						entry["recommended_decision_surface"] = "active-council"
					}
					if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
						currentState := stringValue(entry["authoritative_run_state"])
						if currentState == "ready" || currentState == "executing" {
							holdReason := firstNonBlank(
								stringValue(getKey(decisionPacketByKey["bundle:"+bundleID], "recommended_action", "")),
								stringValue(getKey(decisionPacketByKey["bundle:"+bundleID], "active_decision_next_requirement", "")),
								fmt.Sprintf("Scheduler is holding this mission behind the active bundle council for %s.", bundleID),
							)
							rescheduleActions = append(rescheduleActions, map[string]any{
								"mission_id":                  missionID,
								"run_id":                      runID,
								"action_kind":                 "hold_for_active_bundle_council",
								"target_state":                "waiting",
								"current_state":               currentState,
								"bundle_id":                   bundleID,
								"active_decision_meeting_ids": stringListToAny(activeBundleMeetingIDs),
								"reason":                      holdReason,
							})
						}
					}
				}
			}
		}
		if strings.TrimSpace(stringValue(entry["bundle_control_state"])) == "" {
			entry["bundle_control_state"] = firstNonBlank(stringValue(entry["bundle_state"]), "none")
		}
		if _, ok := entry["bundle_active_decision_meeting_ids"]; !ok {
			entry["bundle_active_decision_meeting_ids"] = []any{}
			entry["bundle_active_decision_meeting_count"] = 0
		}
		entryByMission[missionID] = entry
	}

	sort.Strings(missionOrder)
	runnableMissionIDs := []string{}
	for _, missionID := range missionOrder {
		entry := entryByMission[missionID]
		if stringValue(entry["classification"]) == "runnable" {
			runnableMissionIDs = append(runnableMissionIDs, missionID)
		}
	}
	sort.Slice(runnableMissionIDs, func(i, j int) bool {
		left := entryByMission[runnableMissionIDs[i]]
		right := entryByMission[runnableMissionIDs[j]]
		leftPriority := schedulerPriorityRank(stringValue(left["priority"]))
		rightPriority := schedulerPriorityRank(stringValue(right["priority"]))
		if leftPriority != rightPriority {
			return leftPriority > rightPriority
		}
		leftRisk := schedulerRiskRank(firstNonBlank(stringValue(left["effective_risk_level"]), stringValue(left["risk_level"])))
		rightRisk := schedulerRiskRank(firstNonBlank(stringValue(right["effective_risk_level"]), stringValue(right["risk_level"])))
		if leftRisk != rightRisk {
			return leftRisk > rightRisk
		}
		if boolValue(left["release_adjacent"]) != boolValue(right["release_adjacent"]) {
			return boolValue(left["release_adjacent"])
		}
		leftID := runnableMissionIDs[i]
		rightID := runnableMissionIDs[j]
		if previousForegroundMissionID != "" {
			if leftID == previousForegroundMissionID && rightID != previousForegroundMissionID {
				return true
			}
			if rightID == previousForegroundMissionID && leftID != previousForegroundMissionID {
				return false
			}
		}
		leftUpdated := stringValue(left["updated_at"])
		rightUpdated := stringValue(right["updated_at"])
		if leftUpdated != rightUpdated {
			return leftUpdated > rightUpdated
		}
		return runnableMissionIDs[i] < runnableMissionIDs[j]
	})

	selectedMissionIDs := []string{}
	serializedBehindMissionID := ""
	batchAdmissionPolicy := "single"
	selectionPolicy := "freshest_priority_risk"
	batchCapacityDeferredMissionIDs := []string{}
	if len(runnableMissionIDs) == 1 {
		selectedMissionIDs = append(selectedMissionIDs, runnableMissionIDs[0])
		if runnableMissionIDs[0] == previousForegroundMissionID {
			selectionPolicy = "sticky_foreground"
		}
	} else if len(runnableMissionIDs) > 1 {
		serialize := false
		serializeReason := ""
		for _, missionID := range runnableMissionIDs {
			entry := entryByMission[missionID]
			if schedulerPriorityRank(stringValue(entry["priority"])) >= schedulerPriorityRank("p0") {
				serialize = true
				serializeReason = "p0_priority"
				break
			}
			if boolValue(entry["release_adjacent"]) || schedulerRiskRank(firstNonBlank(stringValue(entry["effective_risk_level"]), stringValue(entry["risk_level"]))) >= 4 {
				serialize = true
				serializeReason = "high_risk"
				break
			}
		}
		if serialize {
			if serializeReason == "p0_priority" {
				batchAdmissionPolicy = "serialize_p0_priority"
			} else {
				batchAdmissionPolicy = "serialize_high_risk"
			}
			selectedMissionIDs = append(selectedMissionIDs, runnableMissionIDs[0])
			if selectedMissionIDs[0] == previousForegroundMissionID {
				selectionPolicy = "sticky_foreground"
			} else if serializeReason == "p0_priority" {
				selectionPolicy = "priority_forced"
			} else {
				selectionPolicy = "risk_forced"
			}
			serializedBehindMissionID = runnableMissionIDs[0]
			for _, missionID := range runnableMissionIDs[1:] {
				entry := entryByMission[missionID]
				entry["classification"] = "deferred"
				reasonText := fmt.Sprintf("Serialized behind higher-risk or release-adjacent mission %s.", serializedBehindMissionID)
				actionKind := "defer_until_foreground_risk_clears"
				actionReason := fmt.Sprintf("Scheduler is serializing this mission behind %s due to higher-risk foreground work.", serializedBehindMissionID)
				recommendedAction := fmt.Sprintf("Defer this mission until %s clears its foreground risk.", serializedBehindMissionID)
				if serializeReason == "p0_priority" {
					reasonText = fmt.Sprintf("Serialized behind P0 mission %s.", serializedBehindMissionID)
					actionKind = "defer_until_p0_clears"
					actionReason = fmt.Sprintf("Scheduler is serializing this mission behind P0 mission %s.", serializedBehindMissionID)
					recommendedAction = fmt.Sprintf("Defer this mission until P0 mission %s clears.", serializedBehindMissionID)
				}
				reasons := append(anyListToStrings(toList(entry["reasons"])), reasonText)
				entry["reasons"] = stringListToAny(dedupeStrings(reasons))
				entry["recommended_action"] = recommendedAction
				if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
					currentState := stringValue(entry["authoritative_run_state"])
					if currentState == "ready" || currentState == "executing" {
						rescheduleActions = append(rescheduleActions, map[string]any{
							"mission_id":    missionID,
							"run_id":        runID,
							"action_kind":   actionKind,
							"target_state":  "waiting",
							"current_state": currentState,
							"reason":        actionReason,
						})
					}
				}
				entryByMission[missionID] = entry
			}
		} else {
			candidateMissionIDs := append([]string{}, runnableMissionIDs...)
			if len(candidateMissionIDs) > maxParallelRuns {
				candidateMissionIDs = candidateMissionIDs[:maxParallelRuns]
			}
			if requireExplicitParallelJustification {
				allCandidatesJustified := len(candidateMissionIDs) > 0
				for _, missionID := range candidateMissionIDs {
					if !boolValue(entryByMission[missionID]["parallel_justified"]) {
						allCandidatesJustified = false
						break
					}
				}
				if !allCandidatesJustified {
					batchAdmissionPolicy = "serialize_missing_parallel_justification"
					selectedMissionIDs = append(selectedMissionIDs, runnableMissionIDs[0])
					if selectedMissionIDs[0] == previousForegroundMissionID {
						selectionPolicy = "sticky_foreground"
					}
					serializedBehindMissionID = runnableMissionIDs[0]
					for _, missionID := range runnableMissionIDs[1:] {
						entry := entryByMission[missionID]
						entry["classification"] = "deferred"
						reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Deferred because the current scheduler batch lacks explicit parallel justification for all candidate missions; foreground remains on %s.", serializedBehindMissionID))
						entry["reasons"] = stringListToAny(dedupeStrings(reasons))
						entry["recommended_action"] = fmt.Sprintf("Record explicit parallel justification before widening this mission alongside %s.", serializedBehindMissionID)
						if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
							currentState := stringValue(entry["authoritative_run_state"])
							if currentState == "ready" || currentState == "executing" {
								rescheduleActions = append(rescheduleActions, map[string]any{
									"mission_id":    missionID,
									"run_id":        runID,
									"action_kind":   "defer_due_missing_parallel_justification",
									"target_state":  "waiting",
									"current_state": currentState,
									"reason":        fmt.Sprintf("Scheduler kept the foreground batch on %s because explicit parallel justification is missing for at least one candidate mission.", serializedBehindMissionID),
								})
							}
						}
						entryByMission[missionID] = entry
					}
				}
			}
			if len(selectedMissionIDs) == 0 {
				batchAdmissionPolicy = "bounded_parallel_capacity"
				if stickySelection, ok := stickyParallelBatchSelection(runnableMissionIDs, entryByMission, previousSelectedMissionIDs, maxParallelRuns); ok {
					selectedMissionIDs = append(selectedMissionIDs, stickySelection...)
					selectionPolicy = "sticky_parallel_batch"
					for _, missionID := range runnableMissionIDs {
						if !containsString(selectedMissionIDs, missionID) {
							batchCapacityDeferredMissionIDs = append(batchCapacityDeferredMissionIDs, missionID)
						}
					}
				} else {
					selectedMissionIDs = append(selectedMissionIDs, runnableMissionIDs...)
					if len(selectedMissionIDs) > 0 && selectedMissionIDs[0] == previousForegroundMissionID {
						selectionPolicy = "sticky_foreground"
					} else {
						selectionPolicy = "parallel_priority_risk"
					}
				}
			}
			if len(selectedMissionIDs) > maxParallelRuns {
				batchCapacityDeferredMissionIDs = append(batchCapacityDeferredMissionIDs, selectedMissionIDs[maxParallelRuns:]...)
				selectedMissionIDs = selectedMissionIDs[:maxParallelRuns]
			}
			if len(batchCapacityDeferredMissionIDs) > 0 {
				batchCapacityDeferredMissionIDs = dedupeStrings(batchCapacityDeferredMissionIDs)
				for _, missionID := range batchCapacityDeferredMissionIDs {
					entry := entryByMission[missionID]
					entry["classification"] = "deferred"
					reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Deferred by scheduler capacity because the current bounded parallel limit is %d.", maxParallelRuns))
					entry["reasons"] = stringListToAny(dedupeStrings(reasons))
					entry["recommended_action"] = fmt.Sprintf("Defer this mission until one of the %d admitted parallel slots clears.", maxParallelRuns)
					if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
						currentState := stringValue(entry["authoritative_run_state"])
						if currentState == "ready" || currentState == "executing" {
							rescheduleActions = append(rescheduleActions, map[string]any{
								"mission_id":    missionID,
								"run_id":        runID,
								"action_kind":   "defer_due_parallel_capacity",
								"target_state":  "waiting",
								"current_state": currentState,
								"reason":        fmt.Sprintf("Scheduler bounded the parallel batch to %d admitted mission(s).", maxParallelRuns),
							})
						}
					}
					entryByMission[missionID] = entry
				}
			}
		}
	}

	selectedRunIDs := []string{}
	if requirePrimaryHumanDecisionBeforeNewForeground && len(requiredHumanDecisions) > 0 && len(selectedMissionIDs) > 0 {
		currentPrimaryDecision, _ := toMap(requiredHumanDecisions[0])
		keepExistingForeground := previousForegroundMissionID != "" && containsString(selectedMissionIDs, previousForegroundMissionID)
		primarySelectedEntry := entryByMission[selectedMissionIDs[0]]
		primarySelectedUrgent := schedulerPriorityRank(stringValue(primarySelectedEntry["priority"])) >= schedulerPriorityRank("p0") ||
			boolValue(primarySelectedEntry["release_adjacent"]) ||
			schedulerRiskRank(firstNonBlank(stringValue(primarySelectedEntry["effective_risk_level"]), stringValue(primarySelectedEntry["risk_level"]))) >= schedulerRiskRank("high")
		if !keepExistingForeground && !primarySelectedUrgent {
			heldMissionIDs := append([]string{}, selectedMissionIDs...)
			selectedMissionIDs = []string{}
			batchAdmissionPolicy = "hold_for_primary_human_decision"
			selectionPolicy = "decision_first_hold"
			primaryDecisionNote := firstNonBlank(
				stringValue(currentPrimaryDecision["decision_kind"]),
				firstNonBlank(stringValue(currentPrimaryDecision["bundle_id"]), stringValue(currentPrimaryDecision["mission_id"]), "the current primary human decision"),
			)
			for _, missionID := range heldMissionIDs {
				entry := entryByMission[missionID]
				if stringValue(entry["classification"]) == "runnable" {
					entry["classification"] = "deferred"
				}
				reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Scheduler is holding new foreground admission until the current primary human decision (%s) is resolved or explicitly deferred.", primaryDecisionNote))
				entry["reasons"] = stringListToAny(dedupeStrings(reasons))
				entry["recommended_action"] = firstNonBlank(
					stringValue(currentPrimaryDecision["recommended_action"]),
					"Resolve the current primary human decision before starting new foreground execution.",
				)
				if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
					currentState := stringValue(entry["authoritative_run_state"])
					if currentState == "ready" || currentState == "executing" {
						rescheduleActions = append(rescheduleActions, map[string]any{
							"mission_id":    missionID,
							"run_id":        runID,
							"action_kind":   "hold_for_primary_human_decision",
							"target_state":  "waiting",
							"current_state": currentState,
							"reason":        "Scheduler is not admitting new foreground work while a higher-priority human decision remains unresolved.",
						})
					}
				}
				entryByMission[missionID] = entry
			}
		}
	}
	primaryMaintenanceAction := map[string]any{}
	maintenanceSelectionPolicy := "maintenance_priority"
	if len(maintenanceActions) > 0 {
		primaryMaintenanceAction, _ = toMap(maintenanceActions[0])
		primaryKey := schedulerMaintenanceKey(primaryMaintenanceAction)
		if previousPrimaryMaintenanceKey != "" && primaryKey == previousPrimaryMaintenanceKey && len(maintenanceActions) > 1 {
			primaryClass := schedulerMaintenanceClass(primaryMaintenanceAction)
			for _, item := range maintenanceActions[1:] {
				peer, _ := toMap(item)
				if schedulerMaintenanceClass(peer) == primaryClass {
					maintenanceSelectionPolicy = "sticky_primary_maintenance"
					break
				}
			}
		}
	}
	if requirePrimaryMaintenanceBeforeNewForeground && len(maintenanceActions) > 0 && len(selectedMissionIDs) > 0 {
		primarySelectedEntry := entryByMission[selectedMissionIDs[0]]
		primarySelectedUrgent := schedulerPriorityRank(stringValue(primarySelectedEntry["priority"])) >= schedulerPriorityRank("p0") ||
			boolValue(primarySelectedEntry["release_adjacent"]) ||
			schedulerRiskRank(firstNonBlank(stringValue(primarySelectedEntry["effective_risk_level"]), stringValue(primarySelectedEntry["risk_level"]))) >= schedulerRiskRank("high")
		if schedulerMaintenanceIsCritical(primaryMaintenanceAction) && !primarySelectedUrgent {
			heldMissionIDs := append([]string{}, selectedMissionIDs...)
			selectedMissionIDs = []string{}
			batchAdmissionPolicy = "hold_for_primary_maintenance"
			selectionPolicy = "maintenance_first_hold"
			primaryMaintenanceNote := firstNonBlank(
				stringValue(primaryMaintenanceAction["action_kind"]),
				firstNonBlank(stringValue(primaryMaintenanceAction["bundle_id"]), stringValue(primaryMaintenanceAction["target_id"]), "the current primary maintenance action"),
			)
			for _, missionID := range heldMissionIDs {
				entry := entryByMission[missionID]
				if stringValue(entry["classification"]) == "runnable" {
					entry["classification"] = "deferred"
				}
				reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Scheduler is holding new foreground admission until the current primary maintenance action (%s) is repaired or explicitly deferred.", primaryMaintenanceNote))
				entry["reasons"] = stringListToAny(dedupeStrings(reasons))
				entry["recommended_action"] = firstNonBlank(
					stringValue(primaryMaintenanceAction["recommended_action"]),
					"Repair the current primary maintenance action before starting new foreground execution.",
				)
				if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
					currentState := stringValue(entry["authoritative_run_state"])
					if currentState == "ready" || currentState == "executing" {
						rescheduleActions = append(rescheduleActions, map[string]any{
							"mission_id":    missionID,
							"run_id":        runID,
							"action_kind":   "hold_for_primary_maintenance",
							"target_state":  "waiting",
							"current_state": currentState,
							"reason":        "Scheduler is not admitting new foreground work while a higher-priority maintenance action remains unresolved.",
						})
					}
				}
				entryByMission[missionID] = entry
			}
		}
	}
	for _, missionID := range selectedMissionIDs {
		if runID := stringValue(entryByMission[missionID]["authoritative_run_id"]); runID != "" {
			selectedRunIDs = append(selectedRunIDs, runID)
		}
	}

	executionMode := "solo"
	if len(selectedMissionIDs) > 1 {
		executionMode = "scheduled_parallel_missions"
	} else if len(selectedMissionIDs) == 1 {
		selected := entryByMission[selectedMissionIDs[0]]
		if intValue(selected["teammate_run_count"]) > 0 || intValue(selected["active_run_count"]) > 1 {
			executionMode = "managed_delegation"
		}
	}

	entries := []any{}
	classificationTally := map[string]int{
		"runnable":              0,
		"deferred":              0,
		"blocked-by-dependency": 0,
		"blocked-by-policy":     0,
		"blocked-by-overlap":    0,
		"completed":             0,
		"needs-human-decision":  0,
	}
	bundleRoleTally := map[string]int{
		"none":                0,
		"governor":            0,
		"governed":            0,
		"governed-ambiguous":  0,
		"escalation-required": 0,
		"conflicted":          0,
		"member":              0,
	}
	activeCouncilBundleIDs := []string{}
	for _, missionID := range missionOrder {
		entry := entryByMission[missionID]
		classification := stringValue(entry["classification"])
		if _, ok := classificationTally[classification]; ok {
			classificationTally[classification]++
		}
		bundleRole := firstNonBlank(stringValue(entry["bundle_role"]), "none")
		if _, ok := bundleRoleTally[bundleRole]; ok {
			bundleRoleTally[bundleRole]++
		}
		if intValue(entry["bundle_active_decision_meeting_count"]) > 0 {
			if bundleID := strings.TrimSpace(stringValue(entry["bundle_id"])); bundleID != "" {
				activeCouncilBundleIDs = append(activeCouncilBundleIDs, bundleID)
			}
		}
		entry["selected"] = false
		for _, selectedMissionID := range selectedMissionIDs {
			if missionID == selectedMissionID {
				entry["selected"] = true
				break
			}
		}
		entries = append(entries, entry)
		if classification == "runnable" && stringValue(entry["authoritative_run_state"]) == "waiting" {
			rescheduleActions = append(rescheduleActions, map[string]any{
				"mission_id":    missionID,
				"run_id":        stringValue(entry["authoritative_run_id"]),
				"action_kind":   "resume_candidate",
				"target_state":  "ready",
				"current_state": stringValue(entry["authoritative_run_state"]),
				"reason":        "Scheduler sees the mission as runnable even though the authoritative run is still waiting.",
			})
		}
	}

	nextAction := "No runnable mission remains; review blockers, meetings, or dependencies."
	councilDecisionCount := 0
	decisionScopeTally := map[string]int{}
	decisionSurfaceTally := map[string]int{}
	for _, item := range requiredHumanDecisions {
		decision, _ := toMap(item)
		if stringValue(decision["recommended_decision_surface"]) == "council" || stringValue(decision["recommended_decision_surface"]) == "active-council" {
			councilDecisionCount++
		}
		scope := firstNonBlank(stringValue(decision["decision_scope"]), "mission")
		decisionScopeTally[scope]++
		surface := firstNonBlank(stringValue(decision["recommended_decision_surface"]), "leader-decision")
		decisionSurfaceTally[surface]++
	}
	primaryRequiredHumanDecision := map[string]any{}
	decisionSelectionPolicy := "portfolio_priority_risk"
	if len(requiredHumanDecisions) > 0 {
		primaryRequiredHumanDecision, _ = toMap(requiredHumanDecisions[0])
		primaryKey := schedulerDecisionKey(primaryRequiredHumanDecision)
		if previousPrimaryDecisionKey != "" && primaryKey == previousPrimaryDecisionKey && len(requiredHumanDecisions) > 1 {
			primaryClass := schedulerDecisionClass(primaryRequiredHumanDecision)
			for _, item := range requiredHumanDecisions[1:] {
				peer, _ := toMap(item)
				if schedulerDecisionClass(peer) == primaryClass {
					decisionSelectionPolicy = "sticky_primary_human_decision"
					break
				}
			}
		}
	}
	primaryControlPlaneKind, primaryControlPlanePacket, primaryControlPlaneSelectionPolicy, deferredControlPlanePacketCount := schedulerPrimaryControlPlanePacket(requiredHumanDecisions, primaryRequiredHumanDecision, maintenanceActions, primaryMaintenanceAction)
	if primaryControlPlaneKind == "maintenance" && schedulerMaintenanceIsCritical(primaryMaintenanceAction) && len(requiredHumanDecisions) > 0 {
		primaryMaintenanceNote := firstNonBlank(
			stringValue(primaryMaintenanceAction["action_kind"]),
			firstNonBlank(stringValue(primaryMaintenanceAction["bundle_id"]), stringValue(primaryMaintenanceAction["target_id"]), "the current primary maintenance action"),
		)
		for _, missionID := range missionOrder {
			entry := entryByMission[missionID]
			if intValue(entry["human_decision_queue_position"]) <= 0 {
				entryByMission[missionID] = entry
				continue
			}
			entry["deferred_by_primary_control_plane_kind"] = "maintenance"
			entry["deferred_by_primary_control_plane_action"] = firstNonBlank(stringValue(primaryMaintenanceAction["action_kind"]), "maintenance")
			reasons := append(anyListToStrings(toList(entry["reasons"])), fmt.Sprintf("Scheduler is deferring this governed human-decision work behind the current primary maintenance action (%s).", primaryMaintenanceNote))
			entry["reasons"] = stringListToAny(dedupeStrings(reasons))
			if runID := stringValue(entry["authoritative_run_id"]); runID != "" {
				currentState := stringValue(entry["authoritative_run_state"])
				if currentState == "ready" || currentState == "executing" {
					rescheduleAction := map[string]any{
						"mission_id":    missionID,
						"run_id":        runID,
						"action_kind":   "defer_pending_primary_maintenance",
						"target_state":  "waiting",
						"current_state": currentState,
						"reason":        "Scheduler is deferring this governed human-decision work behind a higher-priority maintenance action.",
					}
					if decisionID := strings.TrimSpace(stringValue(entry["human_decision_id"])); decisionID != "" {
						rescheduleAction["decision_id"] = decisionID
					}
					rescheduleActions = append(rescheduleActions, rescheduleAction)
				}
			}
			entryByMission[missionID] = entry
		}
	}
	primaryDecisionAction := strings.TrimSpace(stringValue(primaryRequiredHumanDecision["recommended_action"]))
	primaryDecisionSurface := strings.TrimSpace(stringValue(primaryRequiredHumanDecision["recommended_decision_surface"]))
	primaryMaintenanceActionText := strings.TrimSpace(stringValue(primaryMaintenanceAction["recommended_action"]))
	switch {
	case primaryControlPlaneKind == "maintenance" && len(selectedMissionIDs) > 0:
		nextAction = firstNonBlank(primaryMaintenanceActionText, "Keep the current foreground batch bounded and repair the primary governance maintenance action before widening execution.")
	case primaryControlPlaneKind == "maintenance":
		nextAction = firstNonBlank(primaryMaintenanceActionText, "Repair the primary governance maintenance action before widening execution.")
	case primaryControlPlaneKind == "human-decision" && councilDecisionCount > 0 && len(selectedMissionIDs) > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Open or review councils for the required human decisions before widening execution; keep the current foreground batch bounded in the meantime.")
		if nextAction == primaryDecisionAction && len(selectedMissionIDs) > 0 && primaryDecisionSurface == "active-council" {
			nextAction = nextAction + " Keep the current foreground batch bounded in the meantime."
		}
	case primaryControlPlaneKind == "human-decision" && councilDecisionCount > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Open or review councils for the required human decisions before widening execution.")
	case primaryControlPlaneKind == "human-decision" && len(selectedMissionIDs) > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Resolve required human decisions before widening execution; keep the current foreground batch bounded in the meantime.")
		if nextAction == primaryDecisionAction && len(selectedMissionIDs) > 0 && primaryDecisionSurface == "active-council" {
			nextAction = nextAction + " Keep the current foreground batch bounded in the meantime."
		}
	case primaryControlPlaneKind == "human-decision":
		nextAction = firstNonBlank(primaryDecisionAction, "Resolve required human decisions before widening execution.")
	case councilDecisionCount > 0 && len(selectedMissionIDs) > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Open or review councils for the required human decisions before widening execution; keep the current foreground batch bounded in the meantime.")
		if nextAction == primaryDecisionAction && len(selectedMissionIDs) > 0 && primaryDecisionSurface == "active-council" {
			nextAction = nextAction + " Keep the current foreground batch bounded in the meantime."
		}
	case councilDecisionCount > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Open or review councils for the required human decisions before widening execution.")
	case len(requiredHumanDecisions) > 0 && len(selectedMissionIDs) > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Resolve required human decisions before widening execution; keep the current foreground batch bounded in the meantime.")
		if nextAction == primaryDecisionAction && len(selectedMissionIDs) > 0 && primaryDecisionSurface == "active-council" {
			nextAction = nextAction + " Keep the current foreground batch bounded in the meantime."
		}
	case len(requiredHumanDecisions) > 0:
		nextAction = firstNonBlank(primaryDecisionAction, "Resolve required human decisions before widening execution.")
	case len(maintenanceActions) > 0 && schedulerMaintenanceIsCritical(primaryMaintenanceAction) && len(selectedMissionIDs) > 0:
		nextAction = "Keep the current foreground batch bounded and repair the primary governance maintenance action before widening execution."
	case len(maintenanceActions) > 0 && schedulerMaintenanceIsCritical(primaryMaintenanceAction):
		nextAction = "Repair the primary governance maintenance action before widening execution."
	case len(maintenanceActions) > 0 && len(selectedMissionIDs) > 0:
		nextAction = "Keep the current foreground batch bounded and run maintenance on stale control-plane state before widening execution."
	case len(maintenanceActions) > 0:
		nextAction = "Run maintenance on stale control-plane state before widening execution."
	case len(selectedMissionIDs) == 1 && strings.TrimSpace(stringValue(entryByMission[selectedMissionIDs[0]]["adopted_primary_follow_up_action"])) != "":
		nextAction = stringValue(entryByMission[selectedMissionIDs[0]]["adopted_primary_follow_up_action"])
	case len(selectedMissionIDs) > 1:
		if len(batchCapacityDeferredMissionIDs) > 0 {
			nextAction = fmt.Sprintf("Scheduler admitted %d missions into the bounded parallel batch and deferred %d behind the current parallel capacity.", len(selectedMissionIDs), len(batchCapacityDeferredMissionIDs))
		} else {
			nextAction = fmt.Sprintf("Scheduler selected %d missions for a governed parallel batch.", len(selectedMissionIDs))
		}
	case len(selectedMissionIDs) == 1 && executionMode == "managed_delegation":
		nextAction = fmt.Sprintf("Scheduler recommends managed delegation on %s.", selectedMissionIDs[0])
	case len(selectedMissionIDs) == 1:
		nextAction = fmt.Sprintf("Scheduler selected %s as the foreground mission.", selectedMissionIDs[0])
	case classificationTally["blocked-by-dependency"] > 0:
		nextAction = "Satisfy dependencies before resuming blocked missions."
	}

	snapshot := map[string]any{
		"generated_at":                           generatedAt,
		"profile_id":                             stringValue(getKey(profile, "profile.profile_id", "")),
		"execution_mode":                         executionMode,
		"selection_policy":                       selectionPolicy,
		"previous_foreground_mission_id":         previousForegroundMissionID,
		"previous_selected_mission_ids":          stringListToAny(previousSelectedMissionIDs),
		"mission_count":                          len(entries),
		"selected_mission_count":                 len(selectedMissionIDs),
		"selected_run_count":                     len(selectedRunIDs),
		"selected_mission_ids":                   stringListToAny(selectedMissionIDs),
		"selected_run_ids":                       stringListToAny(selectedRunIDs),
		"max_parallel_runs":                      maxParallelRuns,
		"batch_admission_policy":                 batchAdmissionPolicy,
		"capacity_deferred_mission_ids":          stringListToAny(batchCapacityDeferredMissionIDs),
		"capacity_deferred_count":                len(batchCapacityDeferredMissionIDs),
		"classification_tally":                   classificationTally,
		"bundle_role_tally":                      bundleRoleTally,
		"active_council_bundle_count":            len(dedupeStrings(activeCouncilBundleIDs)),
		"active_council_bundle_ids":              stringListToAny(dedupeStrings(activeCouncilBundleIDs)),
		"mission_classifications":                entries,
		"required_human_decision_count":          len(requiredHumanDecisions),
		"deferred_human_decision_count":          maxInt(len(requiredHumanDecisions)-1, 0),
		"council_decision_count":                 councilDecisionCount,
		"required_human_decision_scope_tally":    decisionScopeTally,
		"required_human_decision_surface_tally":  decisionSurfaceTally,
		"decision_selection_policy":              decisionSelectionPolicy,
		"previous_primary_human_decision_key":    previousPrimaryDecisionKey,
		"primary_required_human_decision":        primaryRequiredHumanDecision,
		"required_human_decisions":               requiredHumanDecisions,
		"maintenance_action_count":               len(maintenanceActions),
		"deferred_maintenance_action_count":      maxInt(len(maintenanceActions)-1, 0),
		"maintenance_action_tally":               maintenanceTally,
		"maintenance_selection_policy":           maintenanceSelectionPolicy,
		"previous_primary_maintenance_key":       previousPrimaryMaintenanceKey,
		"primary_maintenance_action":             primaryMaintenanceAction,
		"maintenance_actions":                    maintenanceActions,
		"primary_control_plane_kind":             primaryControlPlaneKind,
		"primary_control_plane_selection_policy": primaryControlPlaneSelectionPolicy,
		"primary_control_plane_packet":           primaryControlPlanePacket,
		"deferred_control_plane_packet_count":    deferredControlPlanePacketCount,
		"reschedule_action_count":                len(rescheduleActions),
		"reschedule_actions":                     rescheduleActions,
		"follow_up_signal_count": len(dedupeStrings(func() []string {
			ids := []string{}
			for _, entry := range entries {
				missionEntry, _ := toMap(entry)
				for _, signal := range toList(missionEntry["follow_up_signals"]) {
					signalMap, _ := toMap(signal)
					if id := stringValue(signalMap["signal_id"]); id != "" {
						ids = append(ids, id)
					}
				}
			}
			return ids
		}())),
		"next_action":          nextAction,
		"snapshot_id":          "schedule-" + now,
		"snapshot_path":        "",
		"latest_snapshot_path": "",
	}
	return snapshot
}

func writeSchedulerSnapshot(snapshot map[string]any) (map[string]any, error) {
	snapshotID := firstNonBlank(stringValue(snapshot["snapshot_id"]), "schedule-"+nowTime().UTC().Format("20060102150405"))
	snapshot["snapshot_id"] = snapshotID
	snapshotPath := schedulerSnapshotPath(snapshotID)
	latestPath := schedulerLatestSnapshotPath()
	payload := map[string]any{
		"version":  1,
		"snapshot": snapshot,
	}
	if err := dumpData(snapshotPath, payload); err != nil {
		return nil, err
	}
	if err := dumpData(latestPath, payload); err != nil {
		return nil, err
	}
	snapshot["snapshot_path"] = snapshotPath
	snapshot["latest_snapshot_path"] = latestPath
	return snapshot, nil
}

func cmdHome(args []string) error {
	_ = args
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	profile, _ := resolveProfile("")
	registry := readRegistry()
	missionRegistry := readMissionRegistry()
	signals := readSignalLog()
	meetings := []any{}
	activeMissions := []any{}
	blockedMissions := []any{}
	reviewMissions := []any{}
	dependencyBlockedMissions := []any{}
	activeRuns := []any{}
	activeTeammates := []any{}
	blockedRuns := []any{}
	waitingRuns := []any{}
	dependencyBlockedRuns := []any{}
	reviewBlockedRuns := []any{}
	checksFailedRuns := []any{}
	mergeReadyRuns := []any{}
	ownershipOverlapClaims := map[string][]string{}
	for _, item := range listMeetings(false) {
		meeting, _ := toMap(item)
		if stringValue(meeting["state"]) == "active" {
			meetings = append(meetings, meeting)
		}
	}
	for _, item := range toList(missionRegistry["missions"]) {
		mission, _ := toMap(item)
		if !activeMission(mission) {
			continue
		}
		activeMissions = append(activeMissions, mission)
		switch stringValue(mission["state"]) {
		case "blocked":
			blockedMissions = append(blockedMissions, mission)
		case "review":
			reviewMissions = append(reviewMissions, mission)
		}
		if len(toList(mission["depends_on"])) > 0 && len(toList(mission["active_run_ids"])) > 0 {
			for _, dep := range toList(mission["depends_on"]) {
				dependencyID := stringValue(dep)
				index := findMissionIndex(missionRegistry, dependencyID)
				if index < 0 {
					dependencyBlockedMissions = append(dependencyBlockedMissions, mission)
					break
				}
				dependencyMission, _ := toMap(toList(missionRegistry["missions"])[index])
				if stringValue(dependencyMission["state"]) != "done" {
					dependencyBlockedMissions = append(dependencyBlockedMissions, mission)
					break
				}
			}
		}
	}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if !activeRun(run) {
			continue
		}
		activeRuns = append(activeRuns, run)
		if activeTeammateRun(run) {
			activeTeammates = append(activeTeammates, run)
		}
		if stringValue(run["state"]) == "blocked" {
			blockedRuns = append(blockedRuns, run)
		}
		if stringValue(run["state"]) == "waiting" {
			waitingRuns = append(waitingRuns, run)
		}
		if len(unsatisfiedDependencies(run, registry)) > 0 {
			dependencyBlockedRuns = append(dependencyBlockedRuns, run)
		}
		if stringValue(run["review_status"]) == "changes_requested" {
			reviewBlockedRuns = append(reviewBlockedRuns, run)
		}
		if stringValue(run["checks_status"]) == "failed" {
			checksFailedRuns = append(checksFailedRuns, run)
		}
		if mergeReady(run) {
			mergeReadyRuns = append(mergeReadyRuns, run)
		}
		if overlapPolicy(run) != "shared" {
			for _, claim := range ownershipClaims(run) {
				ownershipOverlapClaims[claim] = append(ownershipOverlapClaims[claim], stringValue(run["run_id"]))
			}
		}
	}
	overlappingClaims := []any{}
	seenOverlapPairs := map[string]bool{}
	activeExclusiveRuns := []map[string]any{}
	for _, item := range activeRuns {
		run, _ := toMap(item)
		if overlapPolicy(run) != "shared" {
			activeExclusiveRuns = append(activeExclusiveRuns, run)
		}
	}
	for claim, runIDs := range ownershipOverlapClaims {
		if len(runIDs) > 1 {
			overlappingClaims = append(overlappingClaims, map[string]any{
				"claim":   claim,
				"run_ids": stringListToAny(dedupeStrings(runIDs)),
			})
		}
	}
	for index, leftRun := range activeExclusiveRuns {
		leftID := stringValue(leftRun["run_id"])
		for _, rightRun := range activeExclusiveRuns[index+1:] {
			rightID := stringValue(rightRun["run_id"])
			shared := []string{}
			leftExclusions := ownershipExclusions(leftRun)
			rightExclusions := ownershipExclusions(rightRun)
			for _, leftClaim := range ownershipClaims(leftRun) {
				for _, rightClaim := range ownershipClaims(rightRun) {
					if selectorsConflict(leftClaim, leftExclusions, rightClaim, rightExclusions) {
						shared = append(shared, leftClaim+"~"+rightClaim)
					}
				}
			}
			if len(shared) == 0 {
				continue
			}
			pairKey := leftID + "::" + rightID + "::" + strings.Join(dedupeStrings(shared), "|")
			if seenOverlapPairs[pairKey] {
				continue
			}
			seenOverlapPairs[pairKey] = true
			overlappingClaims = append(overlappingClaims, map[string]any{
				"claim":         "semantic-overlap",
				"run_ids":       stringListToAny([]string{leftID, rightID}),
				"selector_refs": stringListToAny(dedupeStrings(shared)),
			})
		}
	}
	activeSignals := []any{}
	coordinationSignals := []any{}
	for _, item := range toList(signals["signals"]) {
		signal, _ := toMap(item)
		if !activeSignal(signal) {
			continue
		}
		activeSignals = append(activeSignals, signal)
		if stringValue(signal["scope_type"]) == "sibling" {
			coordinationSignals = append(coordinationSignals, signal)
			continue
		}
		switch normalizeSignalKind(stringValue(signal["kind"])) {
		case "requires_meeting", "interface_change", "dependency_discovered", "review_feedback", "priority_change", "risk_raised", "requires_rebase":
			coordinationSignals = append(coordinationSignals, signal)
		}
	}
	activeRunIDs := map[string]bool{}
	for _, item := range activeRuns {
		run, _ := toMap(item)
		activeRunIDs[stringValue(run["run_id"])] = true
	}
	orphanSignals := []any{}
	for _, item := range activeSignals {
		signal, _ := toMap(item)
		if runID := stringValue(signal["run_id"]); runID != "" && !activeRunIDs[runID] {
			orphanSignals = append(orphanSignals, signal)
		}
	}
	orphanMeetings := []any{}
	for _, item := range meetings {
		meeting, _ := toMap(item)
		if runID := stringValue(meeting["run_id"]); runID != "" && !activeRunIDs[runID] {
			orphanMeetings = append(orphanMeetings, meeting)
		}
	}
	meetingActions := openMeetingActions()
	activeRCRecord := activeRC()
	missionGraphEdges := 0
	for _, item := range activeMissions {
		mission, _ := toMap(item)
		missionGraphEdges += len(toList(mission["depends_on"]))
	}
	graphAuthorityAudit, _ := toMap(missionGraphAuthorityAudit(missionRegistry, registry, false))
	scheduler := schedulerSummary(profile, missionRegistry, registry, signals, false)
	selectedRunIDs := map[string]bool{}
	for _, item := range toList(scheduler["selected_run_ids"]) {
		if runID := strings.TrimSpace(stringValue(item)); runID != "" {
			selectedRunIDs[runID] = true
		}
	}
	adoptedOutputRunCount := 0
	adoptedOutputFollowUpPacketCount := 0
	adoptedOutputOpenFollowUpPacketCount := 0
	primaryAdoptedOutputRunID := ""
	primaryAdoptedOutputPacketKind := ""
	primaryAdoptedFollowUpPacket := map[string]any{}
	primaryAdoptedRunSelected := false
	primaryAdoptedUpdatedAt := ""
	for _, item := range activeRuns {
		run, _ := toMap(item)
		outputPacket := runAdoptedOutputPacket(run)
		if len(outputPacket) == 0 && strings.TrimSpace(stringValue(run["adopted_output_packet_kind"])) == "" {
			continue
		}
		adoptedOutputRunCount++
		adoptedOutputFollowUpPacketCount += len(runAdoptedFollowUpPackets(run))
		adoptedOutputOpenFollowUpPacketCount += runOpenAdoptedFollowUpPacketCount(run)
		runID := stringValue(run["run_id"])
		runSelected := selectedRunIDs[runID]
		runUpdatedAt := firstNonBlank(runSortTimestamp(run), stringValue(run["updated_at"]))
		shouldPromote := false
		switch {
		case primaryAdoptedOutputRunID == "":
			shouldPromote = true
		case runSelected && !primaryAdoptedRunSelected:
			shouldPromote = true
		case runSelected == primaryAdoptedRunSelected && runUpdatedAt > primaryAdoptedUpdatedAt:
			shouldPromote = true
		}
		if shouldPromote {
			primaryAdoptedOutputRunID = runID
			primaryAdoptedOutputPacketKind = firstNonBlank(strings.TrimSpace(stringValue(run["adopted_output_packet_kind"])), strings.TrimSpace(stringValue(outputPacket["output_packet_kind"])))
			if packet := runPrimaryAdoptedFollowUpPacket(run); packet != nil {
				primaryAdoptedFollowUpPacket, _ = toMap(cloneValue(packet))
			} else {
				primaryAdoptedFollowUpPacket = map[string]any{}
			}
			primaryAdoptedRunSelected = runSelected
			primaryAdoptedUpdatedAt = runUpdatedAt
		}
	}
	nextAction := "Continue active execution and monitor signals."
	switch {
	case len(orphanSignals) > 0 || len(orphanMeetings) > 0:
		nextAction = "Reconcile orphaned runtime objects before proceeding."
	case intValue(scheduler["required_human_decision_count"]) > 0:
		nextAction = stringValue(scheduler["next_action"])
	case intValue(graphAuthorityAudit["cross_mission_conflict_count"]) > 0:
		nextAction = "Resolve cross-mission authority conflicts before widening execution."
	case len(overlappingClaims) > 0:
		nextAction = "Resolve ownership overlap before delegating or continuing parallel execution."
	case len(blockedMissions) > 0:
		nextAction = "Mission graph contains blocked work; resolve blockers before expanding execution."
	case len(activeRuns) == 0:
		nextAction = "No active runs. Claim a mission."
	case len(dependencyBlockedRuns) > 0:
		nextAction = "Some runs are ahead of their dependency graph; satisfy or serialize those dependencies."
	case len(dependencyBlockedMissions) > 0:
		nextAction = "Mission dependencies are not yet satisfied; avoid widening execution until the graph clears."
	case len(reviewMissions) > 0:
		nextAction = "At least one mission is in review; close review loops before opening more execution."
	case len(reviewBlockedRuns) > 0 || len(checksFailedRuns) > 0:
		nextAction = "Triage PR feedback and failing checks before continuing other work."
	case len(mergeReadyRuns) > 0:
		nextAction = "One or more PRs are merge-ready."
	case len(meetingActions) > 0:
		nextAction = "Review open meeting action items and translate them into active execution."
	case len(blockedRuns) > 0:
		nextAction = "Review blockers or convene a meeting."
	case len(meetings) > 0:
		nextAction = "Review active meetings and convert conclusions into run updates."
	case intValue(scheduler["required_human_decision_count"]) == 0 &&
		intValue(scheduler["maintenance_action_count"]) == 0 &&
		intValue(scheduler["selected_mission_count"]) == 1 &&
		adoptedFollowUpActionText(primaryAdoptedFollowUpPacket) != "":
		nextAction = adoptedFollowUpActionText(primaryAdoptedFollowUpPacket)
	case activeRCRecord != nil:
		nextAction = "Review the active RC and decide whether to promote, fix, or close it."
	case intValue(scheduler["selected_mission_count"]) > 1:
		nextAction = "Scheduler found a governed parallel batch; review the selected missions and widen execution deliberately."
	case stringValue(scheduler["execution_mode"]) == "managed_delegation" && intValue(scheduler["selected_mission_count"]) == 1:
		nextAction = "Scheduler recommends managed delegation on the selected foreground mission."
	case len(waitingRuns) > 0:
		nextAction = "Review waiting runs and dependency state."
	}
	return printJSON(map[string]any{
		"profile_id":                                       stringValue(getKey(profile, "profile.profile_id", "")),
		"active_mission_count":                             len(activeMissions),
		"blocked_mission_count":                            len(blockedMissions),
		"review_mission_count":                             len(reviewMissions),
		"dependency_mission_count":                         len(dependencyBlockedMissions),
		"mission_graph_edge_count":                         missionGraphEdges,
		"mission_bundle_count":                             intValue(graphAuthorityAudit["bundle_count"]),
		"mission_conflicted_bundle_count":                  intValue(graphAuthorityAudit["conflicted_bundle_count"]),
		"mission_governed_bundle_count":                    intValue(graphAuthorityAudit["governed_bundle_count"]),
		"mission_escalated_bundle_count":                   intValue(graphAuthorityAudit["escalated_bundle_count"]),
		"mission_graph_conflict_count":                     intValue(graphAuthorityAudit["cross_mission_conflict_count"]),
		"mission_dependency_adjudication_count":            intValue(graphAuthorityAudit["dependency_adjudication_count"]),
		"active_run_count":                                 len(activeRuns),
		"active_teammate_count":                            len(activeTeammates),
		"blocked_run_count":                                len(blockedRuns),
		"dependency_blocked_count":                         len(dependencyBlockedRuns),
		"review_blocked_count":                             len(reviewBlockedRuns),
		"checks_failed_count":                              len(checksFailedRuns),
		"merge_ready_count":                                len(mergeReadyRuns),
		"ownership_overlap_count":                          len(overlappingClaims),
		"active_signal_count":                              len(activeSignals),
		"coordination_signal_count":                        len(coordinationSignals),
		"active_meeting_count":                             len(meetings),
		"active_rc_count":                                  boolToInt(activeRCRecord != nil),
		"open_meeting_action_count":                        len(meetingActions),
		"orphan_signal_count":                              len(orphanSignals),
		"orphan_meeting_count":                             len(orphanMeetings),
		"missions":                                         activeMissions,
		"dependency_blocked_missions":                      dependencyBlockedMissions,
		"mission_authority_audit":                          graphAuthorityAudit,
		"runs":                                             activeRuns,
		"dependency_blocked_runs":                          dependencyBlockedRuns,
		"signals":                                          activeSignals,
		"meetings":                                         meetings,
		"overlapping_claims":                               overlappingClaims,
		"active_rc":                                        mapOrEmpty(activeRCRecord),
		"meeting_actions":                                  meetingActions,
		"orphan_signals":                                   orphanSignals,
		"orphan_meetings":                                  orphanMeetings,
		"scheduler_execution_mode":                         stringValue(scheduler["execution_mode"]),
		"scheduler_selected_mission_count":                 intValue(scheduler["selected_mission_count"]),
		"scheduler_required_human_decision_count":          intValue(scheduler["required_human_decision_count"]),
		"scheduler_deferred_human_decision_count":          intValue(scheduler["deferred_human_decision_count"]),
		"scheduler_council_decision_count":                 intValue(scheduler["council_decision_count"]),
		"scheduler_maintenance_action_count":               intValue(scheduler["maintenance_action_count"]),
		"scheduler_deferred_maintenance_action_count":      intValue(scheduler["deferred_maintenance_action_count"]),
		"scheduler_reschedule_action_count":                intValue(scheduler["reschedule_action_count"]),
		"scheduler_tally":                                  cloneValue(getKey(scheduler, "classification_tally", map[string]any{})),
		"scheduler_bundle_role_tally":                      cloneValue(getKey(scheduler, "bundle_role_tally", map[string]any{})),
		"scheduler_active_council_bundle_count":            intValue(scheduler["active_council_bundle_count"]),
		"scheduler_active_council_bundle_ids":              cloneValue(getKey(scheduler, "active_council_bundle_ids", []any{})),
		"scheduler_required_human_decision_scope_tally":    cloneValue(getKey(scheduler, "required_human_decision_scope_tally", map[string]any{})),
		"scheduler_required_human_decision_surface_tally":  cloneValue(getKey(scheduler, "required_human_decision_surface_tally", map[string]any{})),
		"scheduler_decision_selection_policy":              stringValue(scheduler["decision_selection_policy"]),
		"scheduler_previous_primary_human_decision_key":    stringValue(scheduler["previous_primary_human_decision_key"]),
		"scheduler_primary_required_human_decision":        cloneValue(getKey(scheduler, "primary_required_human_decision", map[string]any{})),
		"scheduler_maintenance_selection_policy":           stringValue(scheduler["maintenance_selection_policy"]),
		"scheduler_previous_primary_maintenance_key":       stringValue(scheduler["previous_primary_maintenance_key"]),
		"scheduler_primary_maintenance_action":             cloneValue(getKey(scheduler, "primary_maintenance_action", map[string]any{})),
		"scheduler_primary_control_plane_kind":             stringValue(scheduler["primary_control_plane_kind"]),
		"scheduler_primary_control_plane_selection_policy": stringValue(scheduler["primary_control_plane_selection_policy"]),
		"scheduler_primary_control_plane_packet":           cloneValue(getKey(scheduler, "primary_control_plane_packet", map[string]any{})),
		"scheduler_deferred_control_plane_packet_count":    intValue(scheduler["deferred_control_plane_packet_count"]),
		"scheduler_maintenance_action_tally":               cloneValue(getKey(scheduler, "maintenance_action_tally", map[string]any{})),
		"scheduler_previous_selected_mission_ids":          cloneValue(toList(scheduler["previous_selected_mission_ids"])),
		"scheduler_selected_mission_ids":                   cloneValue(toList(scheduler["selected_mission_ids"])),
		"scheduler_next_action":                            stringValue(scheduler["next_action"]),
		"adopted_output_run_count":                         adoptedOutputRunCount,
		"adopted_output_follow_up_packet_count":            adoptedOutputFollowUpPacketCount,
		"adopted_output_open_follow_up_packet_count":       adoptedOutputOpenFollowUpPacketCount,
		"primary_adopted_output_run_id":                    primaryAdoptedOutputRunID,
		"primary_adopted_output_packet_kind":               primaryAdoptedOutputPacketKind,
		"primary_adopted_follow_up_packet":                 cloneValue(primaryAdoptedFollowUpPacket),
		"next_action":                                      nextAction,
	})
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func mapOrEmpty(value map[string]any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	return value
}

func cmdDoctor(args []string) error {
	_ = args
	return printJSON(doctorReport())
}

func cmdIssue(args []string) error {
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.New("issue requires a subcommand")
	}
	subcommand := args[0]
	if subcommand != "body" && subcommand != "draft" && subcommand != "create" {
		return fmt.Errorf("unknown issue subcommand: %s", subcommand)
	}
	titleParts := []string{}
	objective := ""
	criteria := []string{}
	notes := []string{}
	labels := []string{}
	assignees := []string{}
	draftOnly := subcommand == "body" || subcommand == "draft"
	for _, arg := range args[1:] {
		switch {
		case strings.HasPrefix(arg, "--objective="):
			objective = strings.TrimSpace(strings.TrimPrefix(arg, "--objective="))
		case strings.HasPrefix(arg, "--criteria="):
			criteria = append(criteria, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--criteria=")))...)
		case strings.HasPrefix(arg, "--note="):
			if note := strings.TrimSpace(strings.TrimPrefix(arg, "--note=")); note != "" {
				notes = append(notes, note)
			}
		case strings.HasPrefix(arg, "--label="):
			if label := strings.TrimSpace(strings.TrimPrefix(arg, "--label=")); label != "" {
				labels = append(labels, label)
			}
		case strings.HasPrefix(arg, "--assignee="):
			if assignee := strings.TrimSpace(strings.TrimPrefix(arg, "--assignee=")); assignee != "" {
				assignees = append(assignees, assignee)
			}
		case arg == "--draft":
			draftOnly = true
		default:
			titleParts = append(titleParts, arg)
		}
	}
	title := strings.TrimSpace(strings.Join(titleParts, " "))
	if title == "" {
		return fmt.Errorf("issue %s requires TITLE", subcommand)
	}
	if objective == "" {
		objective = title
	}
	body, err := issueBodyFromTemplate(title, objective, criteria, notes)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"title":            title,
		"objective":        objective,
		"success_criteria": stringListToAny(criteria),
		"notes":            stringListToAny(notes),
		"labels":           stringListToAny(labels),
		"assignees":        stringListToAny(assignees),
		"body":             body,
		"mode":             ternaryString(draftOnly, "draft", "create"),
	}
	if subcommand == "body" {
		fmt.Println(body)
		return nil
	}
	if draftOnly {
		return printJSON(payload)
	}
	cmd := []string{"gh", "issue", "create", "--title", title, "--body", body}
	for _, label := range labels {
		cmd = append(cmd, "--label", label)
	}
	for _, assignee := range assignees {
		cmd = append(cmd, "--assignee", assignee)
	}
	output, err := runCommandWithStderr(cmd, "")
	if err != nil {
		return fmt.Errorf("issue create failed: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	issueURL := ""
	if len(lines) > 0 {
		issueURL = strings.TrimSpace(lines[len(lines)-1])
	}
	issueNumber := 0
	if matches := issueURLPattern.FindStringSubmatch(issueURL); len(matches) == 2 {
		issueNumber, _ = strconv.Atoi(matches[1])
	}
	payload["issue_url"] = issueURL
	payload["issue_number"] = issueNumber
	return printJSON(payload)
}

func ternaryString(condition bool, whenTrue, whenFalse string) string {
	if condition {
		return whenTrue
	}
	return whenFalse
}

func cmdClaim(args []string) error {
	if len(args) < 1 {
		return errors.New("claim requires ISSUE_NUMBER [TITLE_OVERRIDE]")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	provisionMode := "none"
	applyChanges := false
	dependsOn := []string{}
	ownershipClaimsInput := []string{}
	overlapPolicyValue := "exclusive"
	positionals := []string{}
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--provision="):
			provisionMode = strings.TrimPrefix(arg, "--provision=")
		case arg == "--apply":
			applyChanges = true
		case strings.HasPrefix(arg, "--depends-on="):
			dependsOn = append(dependsOn, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--depends-on=")))...)
		case strings.HasPrefix(arg, "--owns="):
			ownershipClaimsInput = append(ownershipClaimsInput, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--owns=")))...)
		case arg == "--allow-overlap":
			overlapPolicyValue = "shared"
		default:
			positionals = append(positionals, arg)
		}
	}
	if len(positionals) == 0 {
		return errors.New("claim requires ISSUE_NUMBER [TITLE_OVERRIDE]")
	}
	issueNumber := positionals[0]
	titleOverride := ""
	if len(positionals) > 1 {
		titleOverride = strings.TrimSpace(strings.Join(positionals[1:], " "))
	}
	profile, _ := resolveProfile("")
	issueTitle := titleOverride
	issueURL := ""
	issueBody := ""
	if issueTitle == "" {
		output, err := runCommand([]string{"gh", "issue", "view", issueNumber, "--json", "title,body,url"})
		if err != nil {
			return errors.New("claim could not fetch issue data; provide a title override")
		}
		decoder := json.NewDecoder(strings.NewReader(output))
		decoder.UseNumber()
		payload := map[string]any{}
		if err := decoder.Decode(&payload); err != nil {
			return errors.New("claim could not fetch issue data; provide a title override")
		}
		issueTitle = strings.TrimSpace(stringValue(payload["title"]))
		issueBody = strings.TrimSpace(stringValue(payload["body"]))
		issueURL = strings.TrimSpace(stringValue(payload["url"]))
	}

	var summary map[string]any
	err := withHarnessLock(func() error {
		registry := readRegistry()
		missionID := "issue-" + issueNumber
		for _, item := range toList(registry["runs"]) {
			run, _ := toMap(item)
			if activeRun(run) && stringValue(run["mission_id"]) == missionID {
				return fmt.Errorf("active run already exists for %s", missionID)
			}
		}
		slug := slugify(issueTitle)
		branchPattern := stringValue(getKey(profile, "profile.branch_pattern", "mission/{issue}-{slug}"))
		branch := branchFromPattern(branchPattern, issueNumber, slug)
		ownershipClaims := dedupeStrings(append([]string{"mission:" + missionID, "branch:" + branch}, ownershipClaimsInput...))
		contractPath := filepath.Join(runtimeRoot(), "contracts", "issue-"+issueNumber+".yml")
		runID := fmt.Sprintf("run-%s-claim-%s", nowTime().UTC().Format("20060102150405"), issueNumber)
		baseBranch := stringValue(getKey(profile, "profile.base_branch", "main"))
		plannedWorktreePath := ""
		current := currentBranch()
		if provisionMode == "worktree" {
			plannedWorktreePath = provisionWorktreePath(profile, slug)
		} else if provisionMode != "none" && provisionMode != "branch" {
			return errors.New("claim --provision must be one of: none, branch, worktree")
		}

		templateRoot, _ := toMap(loadData(contractTemplatePath(), map[string]any{}))
		template, _ := toMap(templateRoot["contract"])
		discoveries := []string{}
		if issueURL != "" {
			discoveries = append(discoveries, issueURL)
		}
		if issueBody != "" {
			bodySnippet := issueBody
			if len(bodySnippet) > 240 {
				bodySnippet = bodySnippet[:240]
			}
			discoveries = append(discoveries, bodySnippet)
		}
		contractPayload, _ := toMap(mergeDicts(template, map[string]any{
			"mission_id":            missionID,
			"issue_number":          intValue(issueNumber),
			"title":                 issueTitle,
			"objective":             issueTitle,
			"success_criteria":      []any{},
			"sub_tasks":             []any{},
			"depends_on":            stringListToAny(dedupeStrings(dependsOn)),
			"desired_state":         defaultDesiredStateForRun("leader"),
			"branch":                branch,
			"worktree_path":         plannedWorktreePath,
			"context_files":         []any{},
			"verification_commands": getKey(profile, "profile.default_verify_commands", []any{}),
			"notes": map[string]any{
				"assumptions_invalidated": []any{},
				"manual_followups":        []any{},
				"discoveries":             stringListToAny(discoveries),
			},
		}))
		setOwnershipSelectors(contractPayload, ownershipClaims)
		freshness, _ := toMap(contractPayload["freshness"])
		if freshness == nil {
			freshness = map[string]any{}
		}
		if stringValue(freshness["created_at"]) == "" {
			freshness["created_at"] = nowISO()
		}
		contractPayload["freshness"] = freshness
		if err := dumpData(contractPath, map[string]any{"version": 1, "contract": contractPayload}); err != nil {
			return err
		}

		runTemplateRoot, _ := toMap(loadData(runTemplatePath(), map[string]any{}))
		runTemplate, _ := toMap(runTemplateRoot["run"])
		runPayload, _ := toMap(mergeDicts(runTemplate, map[string]any{
			"run_id":               runID,
			"mission_id":           missionID,
			"issue_number":         intValue(issueNumber),
			"leader_id":            "leader-main",
			"owner_type":           "leader",
			"owner_id":             "leader-main",
			"state":                "ready",
			"mode":                 ternaryString(provisionMode == "none", "solo", "managed"),
			"execution_surface":    "foreground-skill",
			"overlap_policy":       overlapPolicyValue,
			"branch":               branch,
			"worktree_path":        plannedWorktreePath,
			"contract_ref":         contractPath,
			"depends_on":           stringListToAny(dedupeStrings(dependsOn)),
			"dependency_policy":    "all",
			"desired_state":        defaultDesiredStateForRun("leader"),
			"objective":            issueTitle,
			"latest_summary":       "Claimed and ready for execution",
			"next_expected_action": "Start drive flow and run the narrowest useful verification.",
		}))
		setOwnershipSelectors(runPayload, ownershipClaims)
		if conflicts := overlappingOwnershipRuns(registry, ownershipClaims, overlapPolicyValue, nil); len(conflicts) > 0 {
			return fmt.Errorf("ownership overlap detected; retry with --allow-overlap or change --owns selectors: %v", conflicts)
		}
		timestamp := nowISO()
		fillIfBlank(runPayload, "created_at", timestamp)
		fillIfBlank(runPayload, "started_at", timestamp)
		fillIfBlank(runPayload, "last_reported_at", timestamp)
		runPayload["updated_at"] = timestamp

		provisioning := map[string]any{
			"mode":            provisionMode,
			"apply_requested": applyChanges,
			"applied":         false,
			"base_branch":     baseBranch,
			"previous_branch": current,
		}
		if provisionMode == "branch" || provisionMode == "worktree" {
			provisioning["target_branch"] = branch
		}
		if provisionMode == "worktree" {
			provisioning["target_worktree_path"] = plannedWorktreePath
		}
		if applyChanges && provisionMode != "none" {
			if err := func() error {
				if provisionMode == "branch" {
					return gitCheckoutNewBranch(branch, baseBranch)
				}
				if err := gitCreateWorktree(branch, plannedWorktreePath, baseBranch); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(plannedWorktreePath, ".mission"), []byte(issueNumber+"\n"), 0o644)
			}(); err != nil {
				provisioning["error"] = err.Error()
				runPayload["state"] = "blocked"
				runPayload["status_reason"] = "Provisioning failed: " + err.Error()
				runPayload["next_expected_action"] = "Clean the current worktree or switch to worktree provisioning."
				log := readSignalLog()
				signalTemplateRoot, _ := toMap(loadData(signalTemplatePath(), map[string]any{}))
				signalTemplate, _ := toMap(signalTemplateRoot["signal"])
				signal, _ := toMap(mergeDicts(signalTemplate, map[string]any{
					"signal_id":       fmt.Sprintf("sig-%s-claim-failed-%s", nowTime().UTC().Format("20060102150405"), issueNumber),
					"kind":            "blocked",
					"source_type":     "leader",
					"source_id":       "leader-main",
					"run_id":          runID,
					"mission_id":      missionID,
					"applies_to":      []any{missionID},
					"scope_type":      "local",
					"title":           "Claim provisioning failed",
					"summary":         err.Error(),
					"authority_level": "L1",
					"confidence":      "high",
					"fresh_until":     nowISO(),
				}))
				fillIfBlank(signal, "emitted_at", nowISO())
				fillIfBlank(signal, "state", "active")
				log["signals"] = append(toList(log["signals"]), signal)
				if err := writeSignalLog(log); err != nil {
					return err
				}
			} else {
				provisioning["applied"] = true
				if provisionMode == "branch" {
					provisioning["current_branch_after_apply"] = currentBranch()
				}
			}
		}

		runPayload["provisioning"] = provisioning
		registry["runs"] = append(toList(registry["runs"]), runPayload)
		if err := writeRegistry(registry); err != nil {
			return err
		}
		summary = map[string]any{
			"mission_id":    missionID,
			"issue_number":  intValue(issueNumber),
			"title":         issueTitle,
			"branch":        branch,
			"mission_path":  missionPath(missionID),
			"contract_path": contractPath,
			"run_id":        runID,
			"issue_url":     issueURL,
			"provisioning":  provisioning,
			"next_action":   runPayload["next_expected_action"],
		}
		return nil
	})
	if err != nil {
		return err
	}
	return printJSON(summary)
}

func cmdDrive(args []string) error {
	if len(args) < 1 {
		return errors.New("drive requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	switch subcommand {
	case "report":
		if len(args) < 3 {
			return errors.New("drive report requires RUN_ID SUMMARY")
		}
		runID := args[1]
		summary := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			run["latest_summary"] = summary
			run["last_reported_at"] = nowISO()
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "start":
		if len(args) < 2 {
			return errors.New("drive start requires RUN_ID [SUMMARY]")
		}
		runID := args[1]
		summary := "Execution started"
		if len(args) > 2 {
			summary = strings.TrimSpace(strings.Join(args[2:], " "))
		}
		return updateDriveRunState(runID, "executing", summary, "")
	case "resume":
		if len(args) < 2 {
			return errors.New("drive resume requires RUN_ID [SUMMARY]")
		}
		runID := args[1]
		summary := "Execution resumed"
		if len(args) > 2 {
			summary = strings.TrimSpace(strings.Join(args[2:], " "))
		}
		return updateDriveRunState(runID, "ready", summary, "")
	case "wait":
		if len(args) < 2 {
			return errors.New("drive wait requires RUN_ID [REASON] [--on=RUN_ID|RUN_ID]")
		}
		runID := args[1]
		reasonTokens := []string{}
		dependsOn := []string{}
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--on="):
				dependsOn = append(dependsOn, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--on=")))...)
			default:
				reasonTokens = append(reasonTokens, arg)
			}
		}
		reason := strings.TrimSpace(strings.Join(reasonTokens, " "))
		if reason == "" {
			reason = "Waiting on dependency resolution."
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, "waiting"); err != nil {
				return err
			}
			run["state"] = "waiting"
			run["status_reason"] = reason
			if len(dependsOn) > 0 {
				run["depends_on"] = stringListToAny(dedupeStrings(append(dependencyTargets(run), dependsOn...)))
			}
			run["last_reported_at"] = nowISO()
			run["updated_at"] = nowISO()
			run["next_expected_action"] = "Monitor dependencies and resume when the required runs satisfy their desired state."
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "complete":
		if len(args) < 2 {
			return errors.New("drive complete requires RUN_ID [SUMMARY]")
		}
		runID := args[1]
		summary := "Execution completed"
		if len(args) > 2 {
			summary = strings.TrimSpace(strings.Join(args[2:], " "))
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, "completed"); err != nil {
				return err
			}
			run["state"] = "completed"
			run["latest_summary"] = summary
			run["status_reason"] = ""
			run["completed_at"] = nowISO()
			run["last_reported_at"] = nowISO()
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "adopt-meeting":
		if len(args) < 3 {
			return errors.New("drive adopt-meeting requires RUN_ID MEETING_ID")
		}
		runID := args[1]
		meetingID := args[2]
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			data := readMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			if stringValue(meeting["state"]) != "closed" {
				return errors.New("drive adopt-meeting requires a closed meeting")
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			run = applyMeetingAdoptionToRun(run, meeting)
			run["status_reason"] = firstNonBlank(stringValue(meeting["resolution_note"]), stringValue(run["status_reason"]))
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"run": run,
				"meeting": map[string]any{
					"meeting_id":      stringValue(meeting["meeting_id"]),
					"resolution_note": stringValue(meeting["resolution_note"]),
				},
			})
		})
	case "reprovision":
		if len(args) < 2 {
			return errors.New("drive reprovision requires RUN_ID [--mode=none|branch|worktree] [--apply]")
		}
		runID := ""
		mode := ""
		applyChanges := false
		for _, arg := range args[1:] {
			switch {
			case strings.HasPrefix(arg, "--mode="):
				mode = strings.TrimPrefix(arg, "--mode=")
			case arg == "--apply":
				applyChanges = true
			case runID == "":
				runID = arg
			}
		}
		if runID == "" {
			return errors.New("drive reprovision requires RUN_ID")
		}
		return withHarnessLock(func() error {
			profile, _ := resolveProfile("")
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			currentMode := stringValue(getKey(run, "provisioning.mode", ""))
			if currentMode == "" {
				if stringValue(run["worktree_path"]) != "" {
					currentMode = "worktree"
				} else {
					currentMode = "none"
				}
			}
			targetMode := mode
			if targetMode == "" {
				targetMode = currentMode
			}
			updatedRun, err := reprovisionRun(run, profile, targetMode, applyChanges)
			if err != nil {
				run["status_reason"] = "Reprovisioning failed: " + err.Error()
				run["next_expected_action"] = "Review reprovision failure and choose a safer provisioning mode."
				run["updated_at"] = nowISO()
				runs[index] = run
				registry["runs"] = runs
				_ = writeRegistry(registry)
				return err
			}
			updatedRun["status_reason"] = ""
			updatedRun["next_expected_action"] = "Start drive flow and run the narrowest useful verification."
			updatedRun["updated_at"] = nowISO()
			if err := persistContractForRun(updatedRun); err != nil {
				return err
			}
			runs[index] = updatedRun
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(updatedRun)
		})
	case "set-next":
		if len(args) < 3 {
			return errors.New("drive set-next requires RUN_ID NEXT_ACTION")
		}
		runID := args[1]
		nextAction := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			run["next_expected_action"] = nextAction
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "block":
		if len(args) < 3 {
			return errors.New("drive block requires RUN_ID REASON")
		}
		runID := args[1]
		reason := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, "blocked"); err != nil {
				return err
			}
			run["state"] = "blocked"
			run["status_reason"] = reason
			run["last_reported_at"] = nowISO()
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			log := readSignalLog()
			signalTemplateRoot, _ := toMap(loadData(signalTemplatePath(), map[string]any{}))
			signalTemplate, _ := toMap(signalTemplateRoot["signal"])
			signal, _ := toMap(mergeDicts(signalTemplate, map[string]any{
				"signal_id":       fmt.Sprintf("sig-%s-blocked-%s", nowTime().UTC().Format("20060102150405"), runID),
				"kind":            "blocked",
				"source_type":     "teammate",
				"source_id":       stringValue(run["owner_id"]),
				"run_id":          runID,
				"mission_id":      stringValue(run["mission_id"]),
				"applies_to":      []any{stringValue(run["mission_id"])},
				"scope_type":      "local",
				"title":           "Run blocked",
				"summary":         reason,
				"authority_level": "L2",
				"confidence":      "high",
				"fresh_until":     nowISO(),
			}))
			fillIfBlank(signal, "emitted_at", nowISO())
			fillIfBlank(signal, "state", "active")
			log["signals"] = append(toList(log["signals"]), signal)
			if err := writeSignalLog(log); err != nil {
				return err
			}
			return printJSON(map[string]any{"run": run, "signal": signal})
		})
	case "meeting":
		if len(args) < 3 {
			return errors.New("drive meeting requires RUN_ID QUESTION [--protocol=evidence-review|incident-command|decision-record]")
		}
		runID := args[1]
		protocol := ""
		questionTokens := []string{}
		for _, arg := range args[2:] {
			if strings.HasPrefix(arg, "--protocol=") {
				protocol = strings.TrimSpace(strings.TrimPrefix(arg, "--protocol="))
			} else {
				questionTokens = append(questionTokens, arg)
			}
		}
		question := strings.TrimSpace(strings.Join(questionTokens, " "))
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, "waiting"); err != nil {
				return err
			}
			backlog := readMaintenanceBacklog()
			meetingID := newMeetingID(stringValue(run["mission_id"]))
			meeting, err := materializeMeeting(meetingID, run, question, stringValue(run["owner_id"]), protocol)
			if err != nil {
				return err
			}
			entry := map[string]any{
				"meeting_id":  meetingID,
				"run_id":      runID,
				"mission_id":  stringValue(run["mission_id"]),
				"question":    question,
				"state":       "active",
				"meeting_ref": meetingPath(meetingID, false),
				"created_at":  firstNonBlank(stringValue(meeting["created_at"]), nowISO()),
				"updated_at":  firstNonBlank(stringValue(meeting["updated_at"]), nowISO()),
			}
			backlog["meetings"] = append(toList(backlog["meetings"]), entry)
			if err := writeMaintenanceBacklog(backlog); err != nil {
				return err
			}
			run["state"] = "waiting"
			run["status_reason"] = "Meeting active: " + question
			run["next_expected_action"] = "Chair the meeting and convert conclusions into run updates."
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{"meeting": meeting, "backlog_entry": entry, "run": run})
		})
	case "meeting-close":
		if len(args) < 3 {
			return errors.New("drive meeting-close requires RUN_ID MEETING_ID [NOTE]")
		}
		runID := args[1]
		meetingID := args[2]
		note := "Meeting closed"
		if len(args) > 3 {
			note = strings.TrimSpace(strings.Join(args[3:], " "))
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			data := readMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			meeting["state"] = "closed"
			meeting["updated_at"] = nowISO()
			meeting["closed_at"] = nowISO()
			meeting["resolution_note"] = note
			if err := dumpData(meetingPath(meetingID, true), map[string]any{"version": 1, "meeting": meeting}); err != nil {
				return err
			}
			_ = os.Remove(meetingPath(meetingID, false))
			if _, err := updateBacklogMeeting(meetingID, map[string]any{
				"state":           "closed",
				"closed_at":       meeting["closed_at"],
				"resolution_note": note,
				"meeting_ref":     meetingPath(meetingID, true),
				"updated_at":      nowISO(),
			}); err != nil {
				return err
			}
			authorityAdoption, err := adoptClosedMeetingAuthority(meeting)
			if err != nil {
				return err
			}
			if len(authorityAdoption) > 0 {
				meeting["authority_adoption"] = authorityAdoption
				if err := dumpData(meetingPath(meetingID, true), map[string]any{"version": 1, "meeting": meeting}); err != nil {
					return err
				}
				if _, err := updateBacklogMeeting(meetingID, map[string]any{
					"authority_adoption": authorityAdoption,
					"updated_at":         nowISO(),
				}); err != nil {
					return err
				}
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, "ready"); err != nil {
				return err
			}
			run["state"] = "ready"
			run["status_reason"] = note
			run = applyMeetingAdoptionToRun(run, meeting)
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"meeting":            meeting,
				"run":                run,
				"authority_adoption": authorityAdoption,
			})
		})
	default:
		return fmt.Errorf("unknown drive subcommand: %s", subcommand)
	}
}

func updateDriveRunState(runID, newState, summary, reason string) error {
	return withHarnessLock(func() error {
		registry := readRegistry()
		index := findRunIndex(registry, runID)
		if index < 0 {
			return fmt.Errorf("run not found: %s", runID)
		}
		runs := toList(registry["runs"])
		run, _ := toMap(runs[index])
		if err := validateRunTransition(run, newState); err != nil {
			return err
		}
		run["state"] = newState
		run["latest_summary"] = summary
		run["status_reason"] = reason
		run["last_reported_at"] = nowISO()
		run["updated_at"] = nowISO()
		runs[index] = run
		registry["runs"] = runs
		if err := writeRegistry(registry); err != nil {
			return err
		}
		return printJSON(run)
	})
}

func cmdShip(args []string) error {
	if len(args) < 1 {
		return errors.New("ship requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	switch subcommand {
	case "status":
		if len(args) < 2 {
			return errors.New("ship status requires RUN_ID")
		}
		registry := readRegistry()
		index := findRunIndex(registry, args[1])
		if index < 0 {
			return fmt.Errorf("run not found: %s", args[1])
		}
		run, _ := toMap(toList(registry["runs"])[index])
		return printJSON(map[string]any{
			"run_id":               stringValue(run["run_id"]),
			"mission_id":           stringValue(run["mission_id"]),
			"verification_status":  firstNonBlank(stringValue(run["verification_status"]), "unknown"),
			"verification_refs":    toList(run["verification_refs"]),
			"review_status":        firstNonBlank(stringValue(run["review_status"]), "unknown"),
			"review_refs":          toList(run["review_refs"]),
			"checks_status":        firstNonBlank(stringValue(run["checks_status"]), "unknown"),
			"checks_refs":          toList(run["checks_refs"]),
			"ship_status":          firstNonBlank(stringValue(run["ship_status"]), "unstarted"),
			"ship_refs":            toList(run["ship_refs"]),
			"pr_ref":               stringValue(run["pr_ref"]),
			"pr_url":               stringValue(run["pr_url"]),
			"rc_id":                stringValue(run["rc_id"]),
			"merge_ready":          mergeReady(run),
			"next_expected_action": stringValue(run["next_expected_action"]),
		})
	case "verify":
		if len(args) < 3 {
			return errors.New("ship verify requires RUN_ID REF")
		}
		runID := args[1]
		ref := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			recordVerification(run, ref)
			run["latest_summary"] = "Verification recorded for ship readiness"
			run["next_expected_action"] = "Mark the run ship-ready or keep iterating if more work remains."
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "ready":
		if len(args) < 2 {
			return errors.New("ship ready requires RUN_ID [SUMMARY]")
		}
		runID := args[1]
		summary := "Run is ready to ship"
		if len(args) > 2 {
			summary = strings.TrimSpace(strings.Join(args[2:], " "))
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if boolValue(policyValue("admissions.require_pre_ship_validation", true)) && stringValue(run["verification_status"]) != "passed" {
				return errors.New("ship ready requires verification_status=passed")
			}
			run["ship_status"] = "ready"
			run["latest_summary"] = summary
			run["next_expected_action"] = "Push the branch and create or update the PR."
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "push":
		if len(args) < 2 {
			return errors.New("ship push requires RUN_ID [REMOTE]")
		}
		runID := args[1]
		remote := "origin"
		if len(args) > 2 {
			remote = args[2]
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			branch := stringValue(run["branch"])
			if branch == "" {
				return errors.New("ship push requires a run branch")
			}
			if _, err := runCommandWithStderr([]string{"git", "push", "-u", remote, branch}, ""); err != nil {
				return fmt.Errorf("ship push failed: %w", err)
			}
			run["ship_status"] = "pushed"
			recordShipRef(run, "push:"+remote+"/"+branch)
			run["latest_summary"] = fmt.Sprintf("Branch pushed to %s: %s", remote, branch)
			run["next_expected_action"] = "Create or update the PR."
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "pr-draft":
		if len(args) < 2 {
			return errors.New("ship pr-draft requires RUN_ID [--title=...] [--base=...]")
		}
		runID := args[1]
		overrideTitle := ""
		overrideBase := ""
		for _, arg := range args[2:] {
			if strings.HasPrefix(arg, "--title=") {
				overrideTitle = strings.TrimSpace(strings.TrimPrefix(arg, "--title="))
			} else if strings.HasPrefix(arg, "--base=") {
				overrideBase = strings.TrimSpace(strings.TrimPrefix(arg, "--base="))
			}
		}
		registry := readRegistry()
		index := findRunIndex(registry, runID)
		if index < 0 {
			return fmt.Errorf("run not found: %s", runID)
		}
		run, _ := toMap(toList(registry["runs"])[index])
		baseBranch := overrideBase
		if baseBranch == "" {
			baseBranch = stringValue(getKey(resolveProfileOrEmpty(), "profile.base_branch", "main"))
		}
		payload, err := prPayloadForRun(run, baseBranch, overrideTitle)
		if err != nil {
			return err
		}
		payload["run_id"] = runID
		payload["mode"] = "draft"
		return printJSON(payload)
	default:
		// remaining ship subcommands are appended below
	}
	return cmdShipRemainder(args)
}

func resolveProfileOrEmpty() map[string]any {
	profile, _ := resolveProfile("")
	return profile
}

func cmdShipRemainder(args []string) error {
	subcommand := args[0]
	switch subcommand {
	case "pr-open":
		if len(args) < 2 {
			return errors.New("ship pr-open requires RUN_ID [--title=...] [--base=...] [--draft] [--reviewer=...]")
		}
		runID := args[1]
		overrideTitle := ""
		overrideBase := ""
		reviewers := []string{}
		draft := false
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--title="):
				overrideTitle = strings.TrimSpace(strings.TrimPrefix(arg, "--title="))
			case strings.HasPrefix(arg, "--base="):
				overrideBase = strings.TrimSpace(strings.TrimPrefix(arg, "--base="))
			case strings.HasPrefix(arg, "--reviewer="):
				if reviewer := strings.TrimSpace(strings.TrimPrefix(arg, "--reviewer=")); reviewer != "" {
					reviewers = append(reviewers, reviewer)
				}
			case arg == "--draft":
				draft = true
			}
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			shipStatus := firstNonBlank(stringValue(run["ship_status"]), "unstarted")
			if shipStatus != "ready" && shipStatus != "pushed" && shipStatus != "in_review" {
				return errors.New("ship pr-open requires ship ready or pushed state")
			}
			baseBranch := overrideBase
			if baseBranch == "" {
				baseBranch = stringValue(getKey(resolveProfileOrEmpty(), "profile.base_branch", "main"))
			}
			payload, err := prPayloadForRun(run, baseBranch, overrideTitle)
			if err != nil {
				return err
			}
			prResult, err := openGitHubPR(run, baseBranch, stringValue(payload["title"]), stringValue(payload["body"]), draft, reviewers)
			if err != nil {
				return fmt.Errorf("ship pr-open failed: %w", err)
			}
			run["ship_status"] = "in_review"
			run["pr_ref"] = stringValue(prResult["pr_ref"])
			run["pr_url"] = stringValue(prResult["pr_url"])
			run["review_status"] = "requested"
			run["checks_status"] = "pending"
			recordShipRef(run, "pr:"+stringValue(prResult["pr_ref"]))
			run["latest_summary"] = "PR opened via GitHub: " + stringValue(prResult["pr_ref"])
			run["next_expected_action"] = "Track review and checks until the PR is merge-ready."
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"run": run,
				"pr": map[string]any{
					"title":       payload["title"],
					"base_branch": payload["base_branch"],
					"pr_ref":      prResult["pr_ref"],
					"pr_url":      prResult["pr_url"],
					"draft":       draft,
					"reviewers":   stringListToAny(reviewers),
				},
			})
		})
	case "pr":
		if len(args) < 3 {
			return errors.New("ship pr requires RUN_ID PR_REF [PR_URL]")
		}
		runID := args[1]
		prRef := args[2]
		prURL := ""
		if len(args) > 3 {
			prURL = args[3]
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if firstNonBlank(stringValue(run["ship_status"]), "unstarted") == "unstarted" {
				return errors.New("ship pr requires ship ready state first")
			}
			run["ship_status"] = "in_review"
			run["pr_ref"] = prRef
			if prURL != "" {
				run["pr_url"] = prURL
			}
			run["review_status"] = "requested"
			run["checks_status"] = "pending"
			recordShipRef(run, "pr:"+prRef)
			run["latest_summary"] = "PR opened or updated: " + prRef
			run["next_expected_action"] = "Track review and checks until the PR is merge-ready."
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "review":
		if len(args) < 3 {
			return errors.New("ship review requires RUN_ID STATE [NOTE]")
		}
		runID := args[1]
		state := args[2]
		note := strings.TrimSpace(strings.Join(args[3:], " "))
		allowed := map[string]bool{"requested": true, "changes_requested": true, "approved": true, "commented": true}
		if !allowed[state] {
			return errors.New("ship review state must be one of: approved, changes_requested, commented, requested")
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			shipStatus := firstNonBlank(stringValue(run["ship_status"]), "unstarted")
			if shipStatus != "in_review" && shipStatus != "approved" {
				return errors.New("ship review requires a PR in review state first")
			}
			recordReview(run, state, firstNonBlank(note, state))
			switch state {
			case "approved":
				run["ship_status"] = "approved"
				run["latest_summary"] = firstNonBlank(note, "PR approved")
				if stringValue(run["checks_status"]) == "passed" {
					run["next_expected_action"] = "Merge the PR when checks are passing."
				} else {
					run["next_expected_action"] = "Wait for checks to pass before merging."
				}
			case "changes_requested":
				run["ship_status"] = "in_review"
				run["latest_summary"] = firstNonBlank(note, "Review requested changes")
				run["next_expected_action"] = "Address requested changes, push updates, and continue review."
			default:
				run["ship_status"] = "in_review"
				run["latest_summary"] = firstNonBlank(note, "Review state updated: "+state)
				run["next_expected_action"] = "Continue review and monitor checks."
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "checks":
		if len(args) < 3 {
			return errors.New("ship checks requires RUN_ID STATE [REF]")
		}
		runID := args[1]
		state := args[2]
		ref := strings.TrimSpace(strings.Join(args[3:], " "))
		allowed := map[string]bool{"pending": true, "failed": true, "passed": true}
		if !allowed[state] {
			return errors.New("ship checks state must be one of: failed, passed, pending")
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			shipStatus := firstNonBlank(stringValue(run["ship_status"]), "unstarted")
			if shipStatus != "in_review" && shipStatus != "approved" {
				return errors.New("ship checks requires a PR in review state first")
			}
			recordChecks(run, state, firstNonBlank(ref, state))
			switch {
			case state == "failed":
				run["ship_status"] = "in_review"
				run["latest_summary"] = firstNonBlank(ref, "Checks failed")
				run["next_expected_action"] = "Fix the failing checks and push an update."
			case state == "passed" && stringValue(run["review_status"]) == "approved":
				run["ship_status"] = "approved"
				run["latest_summary"] = firstNonBlank(ref, "Checks passed and PR is merge-ready")
				run["next_expected_action"] = "Merge the PR."
			default:
				run["ship_status"] = "in_review"
				run["latest_summary"] = firstNonBlank(ref, "Checks state updated: "+state)
				run["next_expected_action"] = "Continue review and checks until merge-ready."
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "merge":
		if len(args) < 2 {
			return errors.New("ship merge requires RUN_ID [--via=git|gh] [--strategy=merge|squash|rebase] [--delete-branch] [SUMMARY]")
		}
		runID := args[1]
		mergeVia := "git"
		strategy := "merge"
		deleteBranch := false
		summaryTokens := []string{}
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--via="):
				mergeVia = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--via=")), "git")
			case strings.HasPrefix(arg, "--strategy="):
				strategy = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--strategy=")), "merge")
			case arg == "--delete-branch":
				deleteBranch = true
			default:
				summaryTokens = append(summaryTokens, arg)
			}
		}
		summary := strings.TrimSpace(strings.Join(summaryTokens, " "))
		if summary == "" {
			summary = "Merged cleanly"
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if boolValue(policyValue("admissions.require_pre_ship_validation", true)) && stringValue(run["verification_status"]) != "passed" {
				return errors.New("ship merge requires verification_status=passed")
			}
			if stringValue(run["pr_ref"]) != "" {
				if stringValue(run["review_status"]) != "approved" {
					return errors.New("ship merge requires review_status=approved when a PR is tracked")
				}
				if stringValue(run["checks_status"]) != "passed" {
					return errors.New("ship merge requires checks_status=passed when a PR is tracked")
				}
			}
			profile, _ := resolveProfile("")
			baseBranch := stringValue(getKey(profile, "profile.base_branch", "main"))
			var err error
			if mergeVia == "gh" {
				if stringValue(run["pr_ref"]) == "" {
					return errors.New("ship merge via gh requires a tracked PR")
				}
				err = ghMergePR(stringValue(run["pr_ref"]), strategy, deleteBranch)
			} else {
				err = gitMergeBranchIntoBase(stringValue(run["branch"]), baseBranch, strategy)
			}
			if err != nil {
				return fmt.Errorf("ship merge failed: %w", err)
			}
			if err := markRunMerged(run, summary); err != nil {
				return err
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"run":                     run,
				"merged_via":              mergeVia,
				"strategy":                strategy,
				"delete_branch_requested": deleteBranch,
			})
		})
	case "merged":
		if len(args) < 2 {
			return errors.New("ship merged requires RUN_ID [SUMMARY]")
		}
		runID := args[1]
		summary := "PR merged"
		if len(args) > 2 {
			summary = strings.TrimSpace(strings.Join(args[2:], " "))
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if boolValue(policyValue("admissions.require_pre_ship_validation", true)) && stringValue(run["verification_status"]) != "passed" {
				return errors.New("ship merged requires verification_status=passed")
			}
			if stringValue(run["pr_ref"]) != "" {
				if stringValue(run["review_status"]) != "approved" {
					return errors.New("ship merged requires review_status=approved when a PR is tracked")
				}
				if stringValue(run["checks_status"]) != "passed" {
					return errors.New("ship merged requires checks_status=passed when a PR is tracked")
				}
			}
			if err := markRunMerged(run, summary); err != nil {
				return err
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(run)
		})
	case "cleanup":
		if len(args) < 2 {
			return errors.New("ship cleanup requires RUN_ID [--delete-branch] [--delete-worktree]")
		}
		runID := args[1]
		deleteBranch := false
		deleteWorktree := false
		for _, arg := range args[2:] {
			if arg == "--delete-branch" {
				deleteBranch = true
			}
			if arg == "--delete-worktree" {
				deleteWorktree = true
			}
		}
		if !deleteBranch && !deleteWorktree {
			deleteWorktree = true
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			profile, _ := resolveProfile("")
			cleaned := []string{}
			if deleteWorktree && stringValue(run["worktree_path"]) != "" {
				if err := gitRemoveWorktree(stringValue(run["worktree_path"])); err != nil {
					return fmt.Errorf("ship cleanup failed removing worktree: %w", err)
				}
				cleaned = append(cleaned, "worktree:"+stringValue(run["worktree_path"]))
				run["worktree_path"] = ""
			}
			if deleteBranch && stringValue(run["branch"]) != "" {
				if err := gitDeleteBranch(stringValue(run["branch"]), stringValue(getKey(profile, "profile.base_branch", "main"))); err != nil {
					return fmt.Errorf("ship cleanup failed deleting branch: %w", err)
				}
				cleaned = append(cleaned, "branch:"+stringValue(run["branch"]))
			}
			for _, ref := range cleaned {
				recordShipRef(run, "cleanup:"+ref)
			}
			if len(cleaned) > 0 {
				run["latest_summary"] = "Post-merge cleanup completed"
			} else {
				run["latest_summary"] = "No cleanup action was needed"
			}
			if stringValue(run["ship_status"]) == "merged" {
				run["next_expected_action"] = "Run is cleanly delivered."
			}
			run["updated_at"] = nowISO()
			if err := updateContractFromRun(run); err != nil {
				return err
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{"run": run, "cleaned": stringListToAny(cleaned)})
		})
	default:
		return fmt.Errorf("unknown ship subcommand: %s", subcommand)
	}
}

func cmdRC(args []string) error {
	if len(args) < 1 {
		return errors.New("rc requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	switch subcommand {
	case "list":
		includeArchived := false
		for _, arg := range args[1:] {
			if arg == "--all" {
				includeArchived = true
			}
		}
		return printJSON(listRCs(includeArchived))
	case "get":
		if len(args) < 2 {
			return errors.New("rc get requires RC_ID")
		}
		data := readRC(args[1])
		if data == nil {
			return fmt.Errorf("rc not found: %s", args[1])
		}
		return printJSON(getKey(data, "rc", map[string]any{}))
	default:
		return cmdRCRemainder(args)
	}
}

func cmdRCRemainder(args []string) error {
	subcommand := args[0]
	switch subcommand {
	case "create":
		if len(args) < 3 {
			return errors.New("rc create requires VERSION RUN_ID...")
		}
		version := args[1]
		runIDs := []string{}
		applyChanges := false
		for _, arg := range args[2:] {
			if arg == "--apply" {
				applyChanges = true
			} else if !strings.HasPrefix(arg, "--") {
				runIDs = append(runIDs, arg)
			}
		}
		if len(runIDs) == 0 {
			return errors.New("rc create requires at least one RUN_ID")
		}
		var rc map[string]any
		return withHarnessLock(func() error {
			if activeRC() != nil {
				return errors.New("an active RC already exists")
			}
			profile, _ := resolveProfile("")
			registry := readRegistry()
			for _, runID := range runIDs {
				index := findRunIndex(registry, runID)
				if index < 0 {
					return fmt.Errorf("run not found: %s", runID)
				}
				run, _ := toMap(toList(registry["runs"])[index])
				if stringValue(run["ship_status"]) != "merged" {
					return fmt.Errorf("run is not merged and cannot enter RC: %s", runID)
				}
			}
			branchPattern := stringValue(getKey(profile, "profile.rc_branch_pattern", "rc/{version}"))
			branch := strings.ReplaceAll(branchPattern, "{version}", version)
			baseBranch := stringValue(getKey(profile, "profile.base_branch", "main"))
			productionBranch := stringValue(getKey(profile, "profile.production_branch", baseBranch))
			rcID := fmt.Sprintf("rc-%s-%s", nowTime().UTC().Format("20060102150405"), slugify(version))
			if applyChanges {
				if err := gitCheckoutNewBranch(branch, baseBranch); err != nil {
					return err
				}
			}
			created, err := materializeRC(rcID, version, branch, baseBranch, productionBranch, runIDs)
			if err != nil {
				return err
			}
			rc = created
			notes := toList(rc["notes"])
			notes = append(notes, noteForRC("RC created for version "+version, map[string]any{"apply_changes": applyChanges}))
			rc["notes"] = notes
			rc["updated_at"] = nowISO()
			if err := dumpData(rcPath(rcID, false), map[string]any{"version": 1, "rc": rc}); err != nil {
				return err
			}
			applyRCRunFollowup(registry, runIDs, "RC created for "+version, "Verify active RC "+rcID+" and decide on promote, fix, or close.", rcID)
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(rc)
		})
	case "verify":
		if len(args) < 3 {
			return errors.New("rc verify requires RC_ID REF")
		}
		rcID := args[1]
		ref := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			data := readActiveRC(rcID)
			if data == nil {
				return fmt.Errorf("active rc not found: %s", rcID)
			}
			rc, _ := toMap(data["rc"])
			refs := appendUniqueText(toList(rc["verification_refs"]), ref)
			rc["verification_refs"] = refs
			rc["verification_status"] = "passed"
			rc["updated_at"] = nowISO()
			rc["notes"] = append(toList(rc["notes"]), noteForRC("Verification recorded: "+ref, nil))
			if err := dumpData(rcPath(rcID, false), map[string]any{"version": 1, "rc": rc}); err != nil {
				return err
			}
			registry := readRegistry()
			applyRCRunFollowup(registry, stringSlice(rc["run_ids"]), "RC "+rcID+" verified", "Promote the RC or continue staged validation if more evidence is needed.", rcID)
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(rc)
		})
	case "cancel":
		if len(args) < 2 {
			return errors.New("rc cancel requires RC_ID [NOTE] [--apply] [--remote=origin]")
		}
		rcID := args[1]
		noteTokens := []string{}
		applyChanges := false
		remote := "origin"
		for _, arg := range args[2:] {
			switch {
			case arg == "--apply":
				applyChanges = true
			case strings.HasPrefix(arg, "--remote="):
				remote = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--remote=")), "origin")
			default:
				noteTokens = append(noteTokens, arg)
			}
		}
		note := strings.TrimSpace(strings.Join(noteTokens, " "))
		if note == "" {
			note = "RC canceled"
		}
		return withHarnessLock(func() error {
			data := readActiveRC(rcID)
			if data == nil {
				return fmt.Errorf("active rc not found: %s", rcID)
			}
			rc, _ := toMap(data["rc"])
			branch := stringValue(rc["branch"])
			localDeleted := false
			remoteDeleted := false
			if applyChanges && branch != "" {
				if err := gitDeleteBranch(branch, stringValue(getKey(rc, "base_branch", "main"))); err != nil {
					return fmt.Errorf("rc cancel failed deleting local branch: %w", err)
				}
				localDeleted = true
				if err := gitDeleteRemoteBranch(remote, branch); err != nil {
					return fmt.Errorf("rc cancel failed deleting remote branch: %w", err)
				}
				remoteDeleted = true
			}
			rc["state"] = "closed"
			rc["updated_at"] = nowISO()
			rc["closed_at"] = nowISO()
			rc["notes"] = append(toList(rc["notes"]), noteForRC(note, map[string]any{
				"apply_changes":         applyChanges,
				"local_branch_deleted":  localDeleted,
				"remote_branch_deleted": remoteDeleted,
				"remote":                remote,
			}))
			if err := dumpData(rcPath(rcID, true), map[string]any{"version": 1, "rc": rc}); err != nil {
				return err
			}
			_ = os.Remove(rcPath(rcID, false))
			registry := readRegistry()
			applyRCRunFollowup(registry, stringSlice(rc["run_ids"]), "RC "+rcID+" canceled", "Fix on "+firstNonBlank(stringValue(rc["base_branch"]), "main")+" and cut a fresh RC if release still matters.", "")
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"rc":                    rc,
				"apply_changes":         applyChanges,
				"local_branch_deleted":  localDeleted,
				"remote_branch_deleted": remoteDeleted,
				"remote":                remote,
			})
		})
	case "fix":
		if len(args) < 3 {
			return errors.New("rc fix requires RC_ID COMMITISH [--apply] [--push] [--remote=origin]")
		}
		rcID := args[1]
		commitish := args[2]
		applyChanges := false
		pushChanges := false
		remote := "origin"
		for _, arg := range args[3:] {
			switch {
			case arg == "--apply":
				applyChanges = true
			case arg == "--push":
				pushChanges = true
			case strings.HasPrefix(arg, "--remote="):
				remote = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--remote=")), "origin")
			}
		}
		return withHarnessLock(func() error {
			data := readActiveRC(rcID)
			if data == nil {
				return fmt.Errorf("active rc not found: %s", rcID)
			}
			rc, _ := toMap(data["rc"])
			resolvedCommit, err := gitResolveCommitish(commitish)
			if err != nil {
				return fmt.Errorf("rc fix could not resolve commitish: %w", err)
			}
			previousBranch := currentBranch()
			applied := false
			pushed := false
			if applyChanges {
				if err := gitCheckoutNewBranch(stringValue(rc["branch"]), stringValue(getKey(rc, "base_branch", "main"))); err != nil {
					return fmt.Errorf("rc fix failed: %w", err)
				}
				defer func() {
					if previousBranch != "" && previousBranch != stringValue(rc["branch"]) && branchExists(previousBranch) {
						_ = gitCheckoutBranch(previousBranch)
					}
				}()
				if err := gitCherryPick(resolvedCommit); err != nil {
					return fmt.Errorf("rc fix failed: %w", err)
				}
				applied = true
				if pushChanges {
					if err := gitPushBranch(remote, stringValue(rc["branch"])); err != nil {
						return fmt.Errorf("rc fix failed: %w", err)
					}
					pushed = true
				}
			}
			rc["fix_refs"] = appendUniqueText(toList(rc["fix_refs"]), resolvedCommit)
			rc["verification_status"] = "unknown"
			rc["updated_at"] = nowISO()
			rc["notes"] = append(toList(rc["notes"]), noteForRC("RC fix cherry-picked from "+commitish, map[string]any{
				"commitish":       commitish,
				"resolved_commit": resolvedCommit,
				"apply_changes":   applyChanges,
				"applied":         applied,
				"pushed":          pushed,
				"remote":          remote,
			}))
			if err := dumpData(rcPath(rcID, false), map[string]any{"version": 1, "rc": rc}); err != nil {
				return err
			}
			registry := readRegistry()
			applyRCRunFollowup(registry, stringSlice(rc["run_ids"]), "RC "+rcID+" received a surgical fix", "Re-verify the RC after the cherry-picked fix before promotion.", rcID)
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"rc":              rc,
				"commitish":       commitish,
				"resolved_commit": resolvedCommit,
				"apply_changes":   applyChanges,
				"applied":         applied,
				"pushed":          pushed,
				"remote":          remote,
			})
		})
	case "hotfix":
		return rcHotfixOrBackport(args, true)
	case "backport":
		return rcHotfixOrBackport(args, false)
	case "promote", "close":
		if len(args) < 2 {
			return fmt.Errorf("rc %s requires RC_ID [NOTE]", subcommand)
		}
		rcID := args[1]
		note := "RC closed"
		if subcommand == "promote" {
			note = "RC promoted"
		}
		if len(args) > 2 {
			note = strings.TrimSpace(strings.Join(args[2:], " "))
		}
		return withHarnessLock(func() error {
			data := readActiveRC(rcID)
			if data == nil {
				return fmt.Errorf("active rc not found: %s", rcID)
			}
			rc, _ := toMap(data["rc"])
			if subcommand == "promote" && boolValue(policyValue("admissions.require_rc_validation", true)) && stringValue(rc["verification_status"]) != "passed" {
				return errors.New("rc promote requires verification_status=passed")
			}
			if subcommand == "promote" {
				rc["state"] = "promoted"
				rc["promoted_at"] = nowISO()
			} else {
				rc["state"] = "closed"
			}
			rc["updated_at"] = nowISO()
			rc["closed_at"] = nowISO()
			rc["notes"] = append(toList(rc["notes"]), noteForRC(note, nil))
			if err := dumpData(rcPath(rcID, true), map[string]any{"version": 1, "rc": rc}); err != nil {
				return err
			}
			_ = os.Remove(rcPath(rcID, false))
			registry := readRegistry()
			nextAction := "RC closed without promotion. Decide whether to re-cut or stop the release."
			if subcommand == "promote" {
				nextAction = "Release candidate promoted. Track production health and clean up release branches."
			}
			applyRCRunFollowup(registry, stringSlice(rc["run_ids"]), "RC "+rcID+" "+map[string]string{"promote": "promoted", "close": "closed"}[subcommand], nextAction, "")
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(rc)
		})
	default:
		return fmt.Errorf("unknown rc subcommand: %s", subcommand)
	}
}

func rcHotfixOrBackport(args []string, hotfix bool) error {
	if len(args) < 3 {
		if hotfix {
			return errors.New("rc hotfix requires VERSION COMMITISH [--apply] [--push] [--pr-open] [--draft] [--remote=origin] [--reviewer=...]")
		}
		return errors.New("rc backport requires VERSION COMMITISH [--apply] [--push] [--pr-open] [--draft] [--remote=origin] [--reviewer=...] [--target=branch]")
	}
	version := args[1]
	commitish := args[2]
	applyChanges := false
	pushChanges := false
	prOpen := false
	draft := false
	remote := "origin"
	reviewers := []string{}
	overrideTarget := ""
	for _, arg := range args[3:] {
		switch {
		case arg == "--apply":
			applyChanges = true
		case arg == "--push":
			pushChanges = true
		case arg == "--pr-open":
			prOpen = true
		case arg == "--draft":
			draft = true
		case strings.HasPrefix(arg, "--remote="):
			remote = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--remote=")), "origin")
		case strings.HasPrefix(arg, "--reviewer="):
			if reviewer := strings.TrimSpace(strings.TrimPrefix(arg, "--reviewer=")); reviewer != "" {
				reviewers = append(reviewers, reviewer)
			}
		case !hotfix && strings.HasPrefix(arg, "--target="):
			overrideTarget = strings.TrimSpace(strings.TrimPrefix(arg, "--target="))
		}
	}
	profile, _ := resolveProfile("")
	productionBranch := stringValue(getKey(profile, "profile.production_branch", "main"))
	baseBranch := stringValue(getKey(profile, "profile.base_branch", productionBranch))
	targetBranch := baseBranch
	branchPattern := stringValue(getKey(profile, "profile.backport_branch_pattern", "backport/{version}"))
	idPrefix := "backport"
	if hotfix {
		branchPattern = stringValue(getKey(profile, "profile.hotfix_branch_pattern", "hotfix/{version}"))
		targetBranch = productionBranch
		idPrefix = "hotfix"
	} else if overrideTarget != "" {
		targetBranch = overrideTarget
	}
	branch := strings.ReplaceAll(branchPattern, "{version}", version)
	recordID := fmt.Sprintf("%s-%s-%s", idPrefix, nowTime().UTC().Format("20060102150405"), slugify(version))
	resolvedCommit, err := gitResolveCommitish(commitish)
	if err != nil {
		if hotfix {
			return fmt.Errorf("rc hotfix could not resolve commitish: %w", err)
		}
		return fmt.Errorf("rc backport could not resolve commitish: %w", err)
	}
	applied := false
	pushed := false
	prResult := map[string]any{}
	if applyChanges {
		branchResult, err := applyCherryPickBranch(branch, targetBranch, resolvedCommit, remote, applyChanges, pushChanges)
		if err != nil {
			if hotfix {
				return fmt.Errorf("rc hotfix failed: %w", err)
			}
			return fmt.Errorf("rc backport failed: %w", err)
		}
		applied = boolValue(branchResult["applied"])
		pushed = boolValue(branchResult["pushed"])
		if prOpen {
			body := ""
			title := ""
			if hotfix {
				body = fmt.Sprintf("## Summary\n\n- Hotfix version: `%s`\n- Source ref: `%s`\n- Resolved commit: `%s`\n\n## Backport\n\n- Required: %s\n- Target base branch: `%s`\n", version, commitish, resolvedCommit, ternaryString(baseBranch != productionBranch, "yes", "no"), baseBranch)
				title = fmt.Sprintf("[hotfix %s] Cherry-pick %s", version, commitish)
			} else {
				body = fmt.Sprintf("## Summary\n\n- Backport version: `%s`\n- Source ref: `%s`\n- Resolved commit: `%s`\n- Production branch: `%s`\n- Target branch: `%s`\n", version, commitish, resolvedCommit, productionBranch, targetBranch)
				title = fmt.Sprintf("[backport %s] Cherry-pick %s", version, commitish)
			}
			result, err := openGitHubPRForBranch(branch, targetBranch, title, body, draft, reviewers)
			if err != nil {
				if hotfix {
					return fmt.Errorf("rc hotfix failed: %w", err)
				}
				return fmt.Errorf("rc backport failed: %w", err)
			}
			prResult = result
		}
	}
	payload := map[string]any{
		"version":         version,
		"branch":          branch,
		"commitish":       commitish,
		"resolved_commit": resolvedCommit,
		"apply_changes":   applyChanges,
		"applied":         applied,
		"pushed":          pushed,
		"remote":          remote,
		"pr_opened":       len(prResult) > 0,
		"pr_ref":          stringValue(prResult["pr_ref"]),
		"pr_url":          stringValue(prResult["pr_url"]),
		"reviewers":       stringListToAny(reviewers),
	}
	if hotfix {
		payload["hotfix_id"] = recordID
		payload["production_branch"] = productionBranch
		payload["base_branch"] = baseBranch
		payload["backport_required"] = baseBranch != productionBranch
		if baseBranch != productionBranch {
			payload["next_expected_action"] = fmt.Sprintf("After production merge, backport the hotfix into %s.", baseBranch)
		} else {
			payload["next_expected_action"] = "No backport required in single-branch mode."
		}
	} else {
		payload["backport_id"] = recordID
		payload["target_branch"] = targetBranch
		payload["production_branch"] = productionBranch
		payload["next_expected_action"] = fmt.Sprintf("Merge the backport into %s after validating it there.", targetBranch)
	}
	return printJSON(payload)
}

func cmdMeetings(args []string) error {
	if len(args) < 1 {
		return errors.New("meetings requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	switch subcommand {
	case "list":
		includeArchived := false
		for _, arg := range args[1:] {
			if arg == "--all" {
				includeArchived = true
			}
		}
		return printJSON(listMeetings(includeArchived))
	case "active":
		meetings := []any{}
		missionID := ""
		if len(args) > 1 && !strings.HasPrefix(args[1], "--") {
			missionID = args[1]
		}
		for _, item := range listMeetings(false) {
			meeting, _ := toMap(item)
			if stringValue(meeting["state"]) == "active" && (missionID == "" || stringValue(meeting["mission_id"]) == missionID) {
				meetings = append(meetings, meeting)
			}
		}
		return printJSON(meetings)
	case "get":
		if len(args) < 2 {
			return errors.New("meetings get requires MEETING_ID")
		}
		data := readMeeting(args[1])
		if data == nil {
			return fmt.Errorf("meeting not found: %s", args[1])
		}
		return printJSON(getKey(data, "meeting", map[string]any{}))
	case "summary":
		if len(args) < 2 {
			return errors.New("meetings summary requires MEETING_ID")
		}
		data := readMeeting(args[1])
		if data == nil {
			return fmt.Errorf("meeting not found: %s", args[1])
		}
		meeting, _ := toMap(data["meeting"])
		ensureMeetingGovernanceState(meeting)
		outcome := meetingOutcome(meeting)
		openActionCount := 0
		for _, item := range toList(meeting["action_items"]) {
			mapping, ok := toMap(item)
			if !ok || stringValue(mapping["state"]) == "" || stringValue(mapping["state"]) == "open" {
				openActionCount++
			}
		}
		findings, _ := toMap(meeting["findings"])
		orchestration, _ := toMap(meeting["orchestration"])
		synthesis, _ := toMap(orchestration["synthesis"])
		autoSynthesis, _ := toMap(orchestration["auto_synthesis"])
		debatePlan, _ := toMap(orchestration["debate_plan"])
		wavePlan, _ := toMap(orchestration["wave_plan"])
		outputPacket, _ := toMap(orchestration["output_packet"])
		currentWaveBatch := meetingCurrentWaveBatch(wavePlan, debatePlan)
		currentWaveActions := meetingCurrentWaveActions(wavePlan, debatePlan)
		blockedWaveBatches := meetingBlockedWaveBatchSummaries(debatePlan)
		authorityAudit, _ := toMap(meetingAuthorityAudit(meeting))
		governanceGaps := meetingGovernancePolicyGaps(meeting)
		return printJSON(map[string]any{
			"meeting_id":                               stringValue(meeting["meeting_id"]),
			"state":                                    stringValue(meeting["state"]),
			"decision_protocol":                        stringValue(meeting["decision_protocol"]),
			"decision_protocol_reason":                 stringValue(meeting["decision_protocol_reason"]),
			"question":                                 stringValue(meeting["question"]),
			"outcome_summary":                          outcome["summary"],
			"next_action":                              outcome["next_action"],
			"phase_history_count":                      len(toList(meeting["phase_history"])),
			"governance_gap_count":                     len(governanceGaps),
			"governance_gaps":                          stringListToAny(governanceGaps),
			"debate_round":                             intValue(orchestration["debate_round"]),
			"synthesis_status":                         stringValue(synthesis["status"]),
			"blocking_gap_count":                       len(toList(synthesis["blocking_gaps"])),
			"coverage_note_count":                      len(toList(synthesis["coverage_notes"])),
			"authority_conflict_count":                 intValue(authorityAudit["conflict_count"]),
			"unresolved_authority_conflict_count":      intValue(authorityAudit["unresolved_conflict_count"]),
			"role_spec_count":                          len(toList(meeting["role_specs"])),
			"locked_decisions":                         len(toList(meeting["locked_decisions"])),
			"accepted_count":                           len(toList(findings["accepted"])),
			"rejected_count":                           len(toList(findings["rejected"])),
			"unresolved_count":                         len(toList(findings["unresolved"])),
			"traced_finding_count":                     intValue(getKey(synthesis, "tally.traced_findings", 0)),
			"untraced_finding_count":                   intValue(getKey(synthesis, "tally.untraced_findings", 0)),
			"auto_generated_finding_count":             intValue(autoSynthesis["generated_finding_count"]),
			"auto_generated_action_count":              intValue(autoSynthesis["generated_action_count"]),
			"auto_generated_claim_count":               intValue(autoSynthesis["generated_claim_count"]),
			"planned_route_count":                      intValue(debatePlan["planned_route_count"]),
			"resolved_route_count":                     intValue(debatePlan["resolved_route_count"]),
			"open_route_action_count":                  intValue(debatePlan["open_route_action_count"]),
			"resolved_without_cross_examination_count": intValue(debatePlan["resolved_without_cross_examination_count"]),
			"authority_resolved_route_count":           intValue(debatePlan["authority_resolved_route_count"]),
			"protocol_deferred_route_count":            intValue(debatePlan["protocol_deferred_route_count"]),
			"protocol_bound_route_count":               intValue(debatePlan["protocol_bound_route_count"]),
			"planned_round_count":                      intValue(debatePlan["planned_round_count"]),
			"resolved_round_count":                     intValue(debatePlan["resolved_round_count"]),
			"open_round_action_count":                  intValue(debatePlan["open_round_action_count"]),
			"wave_batch_count":                         intValue(debatePlan["wave_batch_count"]),
			"open_wave_batch_count":                    intValue(debatePlan["open_wave_batch_count"]),
			"ready_wave_batch_count":                   intValue(debatePlan["ready_wave_batch_count"]),
			"blocked_wave_batch_count":                 intValue(debatePlan["blocked_wave_batch_count"]),
			"resolved_wave_batch_count":                intValue(debatePlan["resolved_wave_batch_count"]),
			"batch_selection_policy":                   stringValue(debatePlan["batch_selection_policy"]),
			"previous_current_wave_batch_id":           stringValue(debatePlan["previous_current_wave_batch_id"]),
			"current_wave_batch_id":                    stringValue(debatePlan["current_wave_batch_id"]),
			"current_wave_batch_summary":               stringValue(debatePlan["current_wave_batch_summary"]),
			"current_wave_batch_lane_count":            intValue(debatePlan["current_wave_batch_lane_count"]),
			"current_wave_batch_assigned_role":         stringValue(getKey(currentWaveBatch, "assigned_role", "")),
			"current_wave_batch_assigned_run_id":       stringValue(getKey(currentWaveBatch, "assigned_run_id", "")),
			"current_wave_batch_phase":                 stringValue(getKey(currentWaveBatch, "expected_phase", "")),
			"current_wave_batch_round_number":          intValue(getKey(currentWaveBatch, "round_number", 0)),
			"current_wave_batch_action_count":          intValue(getKey(currentWaveBatch, "action_count", 0)),
			"current_wave_batch_claim_refs":            cloneValue(getKey(currentWaveBatch, "claim_refs", []any{})),
			"current_wave_actions":                     currentWaveActions,
			"blocked_wave_batches":                     blockedWaveBatches,
			"output_packet_kind":                       stringValue(getKey(outputPacket, "output_packet_kind", "")),
			"output_packet_completion_gap_count":       intValue(getKey(outputPacket, "completion_gap_count", 0)),
			"output_packet":                            cloneValue(outputPacket),
			"output_debate_lane_count":                 intValue(getKey(outputPacket, "debate_lane_count", 0)),
			"output_accepted_debate_lane_count":        intValue(getKey(outputPacket, "accepted_debate_lane_count", 0)),
			"output_rejected_debate_lane_count":        intValue(getKey(outputPacket, "rejected_debate_lane_count", 0)),
			"output_unresolved_debate_lane_count":      intValue(getKey(outputPacket, "unresolved_debate_lane_count", 0)),
			"output_open_debate_lane_count":            intValue(getKey(outputPacket, "open_debate_lane_count", 0)),
			"output_debate_lane_packets":               cloneValue(getKey(outputPacket, "debate_lane_packets", []any{})),
			"output_follow_up_packet_count":            intValue(getKey(outputPacket, "follow_up_packet_count", 0)),
			"output_open_follow_up_packet_count":       intValue(getKey(outputPacket, "open_follow_up_packet_count", 0)),
			"output_follow_up_packets":                 cloneValue(getKey(outputPacket, "follow_up_packets", []any{})),
			"multi_lane_batch_count":                   intValue(debatePlan["multi_lane_batch_count"]),
			"max_wave_lane_count":                      intValue(debatePlan["max_wave_lane_count"]),
			"max_debate_rounds":                        intValue(debatePlan["max_debate_rounds"]),
			"open_wave_action_count":                   intValue(wavePlan["open_wave_action_count"]),
			"unrouted_claim_count":                     intValue(debatePlan["unrouted_claim_count"]),
			"open_action_count":                        openActionCount,
		})
	case "authority-audit":
		if len(args) < 2 {
			return errors.New("meetings authority-audit requires MEETING_ID")
		}
		data := readMeeting(args[1])
		if data == nil {
			return fmt.Errorf("meeting not found: %s", args[1])
		}
		meeting, _ := toMap(data["meeting"])
		return printJSON(meetingAuthorityAudit(meeting))
	case "materialize-primary-decision":
		requestedMissionID := ""
		for _, arg := range args[1:] {
			if strings.HasPrefix(arg, "--mission=") {
				requestedMissionID = strings.TrimSpace(strings.TrimPrefix(arg, "--mission="))
			}
		}
		return withHarnessLock(func() error {
			snapshot := latestSchedulerSnapshot()
			if len(snapshot) == 0 {
				return errors.New("no latest scheduler snapshot found; run missions schedule first")
			}
			decision, _ := toMap(snapshot["primary_required_human_decision"])
			if len(decision) == 0 || stringValue(decision["decision_id"]) == "" {
				return errors.New("latest scheduler snapshot has no primary human-decision packet to materialize")
			}
			registry := readRegistry()
			missionRegistry := readMissionRegistry()
			meeting, entry, run, err := materializeMeetingFromSchedulerDecision(snapshot, decision, registry, missionRegistry, requestedMissionID)
			if err != nil {
				return err
			}
			return printJSON(map[string]any{
				"decision":      decision,
				"meeting":       meeting,
				"backlog_entry": entry,
				"run":           run,
			})
		})
	default:
		return cmdMeetingsRemainder(args)
	}
}

func cmdMeetingsRemainder(args []string) error {
	subcommand := args[0]
	switch subcommand {
	case "materialize-decision":
		if len(args) < 2 {
			return errors.New("meetings materialize-decision requires DECISION_ID [--mission=MISSION_ID]")
		}
		decisionID := strings.TrimSpace(args[1])
		requestedMissionID := ""
		for _, arg := range args[2:] {
			if strings.HasPrefix(arg, "--mission=") {
				requestedMissionID = strings.TrimSpace(strings.TrimPrefix(arg, "--mission="))
			}
		}
		return withHarnessLock(func() error {
			snapshot := latestSchedulerSnapshot()
			if len(snapshot) == 0 {
				return errors.New("no latest scheduler snapshot found; run missions schedule first")
			}
			decision := schedulerDecisionByID(snapshot, decisionID)
			if len(decision) == 0 {
				return fmt.Errorf("scheduler decision not found: %s", decisionID)
			}
			registry := readRegistry()
			missionRegistry := readMissionRegistry()
			meeting, entry, run, err := materializeMeetingFromSchedulerDecision(snapshot, decision, registry, missionRegistry, requestedMissionID)
			if err != nil {
				return err
			}
			return printJSON(map[string]any{
				"decision":      decision,
				"meeting":       meeting,
				"backlog_entry": entry,
				"run":           run,
			})
		})
	case "add-role":
		if len(args) < 3 {
			return errors.New("meetings add-role requires MEETING_ID ROLE")
		}
		meetingID := args[1]
		roleTokens := []string{}
		mission := ""
		allowedSources := []string{}
		forbiddenSources := []string{}
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--mission="):
				mission = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--mission=")), mission)
			case strings.HasPrefix(arg, "--allow="):
				allowedSources = append(allowedSources, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--allow=")))...)
			case strings.HasPrefix(arg, "--forbid="):
				forbiddenSources = append(forbiddenSources, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--forbid=")))...)
			default:
				roleTokens = append(roleTokens, arg)
			}
		}
		role := strings.TrimSpace(strings.Join(roleTokens, " "))
		if role == "" {
			return errors.New("meetings add-role requires ROLE")
		}
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			upsertMeetingRoleSpec(meeting, role, mission, allowedSources, forbiddenSources)
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "set-phase":
		if len(args) < 3 {
			return errors.New("meetings set-phase requires MEETING_ID PHASE")
		}
		allowed := map[string]bool{"setup": true, "independent_review": true, "cross_examination": true, "synthesis": true, "output": true, "closed": true}
		if !allowed[args[2]] {
			return errors.New("meetings set-phase phase must be one of: closed, cross_examination, independent_review, output, setup, synthesis")
		}
		return withHarnessLock(func() error {
			data := readActiveMeeting(args[1])
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", args[1])
			}
			meeting, _ := toMap(data["meeting"])
			setMeetingPhase(meeting, args[2], "manual phase override")
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "set-protocol":
		if len(args) < 3 {
			return errors.New("meetings set-protocol requires MEETING_ID PROTOCOL")
		}
		protocol := normalizeDecisionProtocol(args[2])
		if protocol == "" {
			return errors.New("meetings set-protocol protocol must be one of: evidence-review, incident-command, decision-record")
		}
		return withHarnessLock(func() error {
			data := readActiveMeeting(args[1])
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", args[1])
			}
			meeting, _ := toMap(data["meeting"])
			oldProtocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
			packet, _ := toMap(meeting["packet"])
			if packet == nil {
				packet = map[string]any{}
			}
			if firstNonBlank(oldProtocol, "evidence-review") != protocol {
				if stringValue(packet["expected_output"]) == "" || stringValue(packet["expected_output"]) == defaultMeetingExpectedOutput(oldProtocol) {
					packet["expected_output"] = defaultMeetingExpectedOutput(protocol)
				}
				meeting["roles"] = defaultMeetingRoles(protocol)
				meeting["role_specs"] = defaultMeetingRoleSpecs(protocol)
			}
			meeting["decision_protocol"] = protocol
			meeting["decision_protocol_reason"] = fmt.Sprintf("Manual protocol override: %s.", protocol)
			meeting["packet"] = packet
			meeting["updated_at"] = nowISO()
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "orchestrate":
		if len(args) < 2 {
			return errors.New("meetings orchestrate requires MEETING_ID")
		}
		meetingID := args[1]
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			ensureMeetingGovernanceState(meeting)
			protocol := normalizeDecisionProtocol(stringValue(meeting["decision_protocol"]))
			if protocol == "" {
				protocol = "evidence-review"
			}
			oldPhase := firstNonBlank(stringValue(meeting["phase"]), "setup")
			packet, _ := toMap(meeting["packet"])
			roleSpecs := toList(meeting["role_specs"])
			participants := toList(meeting["participants"])
			coverageAudit, _ := toMap(meetingCoverageAudit(meeting))
			wavePlan, _ := toMap(meetingWavePlan(meeting, coverageAudit))
			debatePlan, _ := toMap(meetingDebatePlan(meeting))
			governanceGaps := meetingGovernancePolicyGaps(meeting)
			openDebateBatchCount := intValue(debatePlan["open_wave_batch_count"])
			blockedDebateBatchCount := intValue(debatePlan["blocked_wave_batch_count"])
			currentDebateBatchSummary := strings.TrimSpace(stringValue(debatePlan["current_wave_batch_summary"]))
			positions := toList(meeting["positions"])
			gaps := []string{}
			if len(roleSpecs) == 0 {
				gaps = append(gaps, "missing role specs")
			}
			if len(participants) == 0 {
				gaps = append(gaps, "missing participants")
			}
			if strings.TrimSpace(stringValue(packet["expected_output"])) == "" {
				gaps = append(gaps, "missing expected output")
			}
			independentReviewRuns := map[string]bool{}
			crossExamRuns := map[string]bool{}
			for _, item := range positions {
				position, ok := toMap(item)
				if !ok {
					continue
				}
				switch stringValue(position["phase"]) {
				case "independent_review":
					independentReviewRuns[stringValue(position["run_id"])] = true
				case "cross_examination":
					crossExamRuns[stringValue(position["run_id"])] = true
				}
			}
			newPhase := oldPhase
			nextRequirement := ""
			currentDebateBatchRequirement := func(prefix, fallback string) string {
				if currentDebateBatchSummary == "" {
					if blockedDebateBatchCount > 0 {
						return fmt.Sprintf("%s %d additional debate batch(es) remain blocked.", fallback, blockedDebateBatchCount)
					}
					return fallback
				}
				message := fmt.Sprintf("%s %s", prefix, currentDebateBatchSummary)
				if blockedDebateBatchCount > 0 {
					message = fmt.Sprintf("%s %d additional debate batch(es) remain blocked behind it.", message, blockedDebateBatchCount)
				}
				return message
			}
			protocolPrefersDirectSynthesis := func() bool {
				return !protocolUsesStructuredDebate(protocol) &&
					(intValue(debatePlan["protocol_deferred_route_count"]) > 0 || intValue(debatePlan["protocol_bound_route_count"]) > 0)
			}
			switch oldPhase {
			case "", "setup":
				gaps = append(gaps, governanceGaps...)
				if len(gaps) == 0 {
					newPhase = "independent_review"
					nextRequirement = protocolIndependentReviewRequirement(protocol, intValue(wavePlan["open_wave_action_count"]))
				} else {
					nextRequirement = "Add packet/participant prerequisites before independent review."
				}
			case "independent_review":
				gaps = append(gaps, governanceGaps...)
				for _, item := range toList(coverageAudit["blocking_gaps"]) {
					text := strings.TrimSpace(stringValue(item))
					if text != "" {
						gaps = append(gaps, text)
					}
				}
				if len(independentReviewRuns) >= 2 {
					if intValue(debatePlan["open_route_action_count"]) == 0 && intValue(debatePlan["open_round_action_count"]) == 0 {
						newPhase = "synthesis"
						if intValue(debatePlan["planned_route_count"]) == 0 {
							nextRequirement = "No conflicting claim lanes remain; synthesize the independent-review record."
						} else if protocolPrefersDirectSynthesis() {
							nextRequirement = protocolDeferredSynthesisRequirement(protocol)
						} else if intValue(debatePlan["resolved_without_cross_examination_count"]) > 0 {
							nextRequirement = "No unresolved debate lanes remain; synthesize the authority-resolved council record."
						} else {
							nextRequirement = "Synthesize the current council record."
						}
					} else {
						newPhase = "cross_examination"
						if intValue(debatePlan["open_route_action_count"]) > 0 {
							nextRequirement = fmt.Sprintf("Execute %d planned cross-examination route(s) before synthesis.", intValue(debatePlan["open_route_action_count"]))
						} else {
							nextRequirement = "Collect cross-examination positions that challenge or refine earlier claims."
						}
					}
				} else {
					if intValue(wavePlan["open_wave_action_count"]) > 0 {
						nextRequirement = fmt.Sprintf("Still waiting on %d independent-review wave obligation(s) before cross-examination.", intValue(wavePlan["open_wave_action_count"]))
					} else {
						nextRequirement = "Need at least two independent-review positions before cross-examination."
					}
				}
			case "cross_examination":
				gaps = append(gaps, governanceGaps...)
				for _, item := range toList(coverageAudit["blocking_gaps"]) {
					text := strings.TrimSpace(stringValue(item))
					if text != "" {
						gaps = append(gaps, text)
					}
				}
				for _, item := range toList(debatePlan["blocking_gaps"]) {
					text := strings.TrimSpace(stringValue(item))
					if text != "" {
						gaps = append(gaps, text)
					}
				}
				if intValue(debatePlan["open_route_action_count"]) == 0 &&
					intValue(debatePlan["open_round_action_count"]) == 0 &&
					(intValue(debatePlan["planned_route_count"]) == 0 || intValue(debatePlan["resolved_route_count"]) >= intValue(debatePlan["planned_route_count"])) {
					newPhase = "synthesis"
					if intValue(debatePlan["resolved_without_cross_examination_count"]) > 0 {
						nextRequirement = "No unresolved debate lanes remain; synthesize the authority-resolved council record."
					} else if protocolPrefersDirectSynthesis() {
						nextRequirement = protocolDeferredSynthesisRequirement(protocol)
					} else {
						nextRequirement = "Convert debated positions into accepted, rejected, unresolved findings, and action items."
					}
				} else if len(crossExamRuns) >= 1 && intValue(coverageAudit["linked_cross_examination_count"]) > 0 {
					if intValue(debatePlan["open_round_action_count"]) > 0 {
						nextRequirement = fmt.Sprintf("Execute %d follow-up debate round(s) before synthesis.", intValue(debatePlan["open_round_action_count"]))
					} else {
						newPhase = "synthesis"
						if protocolPrefersDirectSynthesis() {
							nextRequirement = protocolDeferredSynthesisRequirement(protocol)
						} else {
							nextRequirement = "Convert debated positions into accepted, rejected, unresolved findings, and action items."
						}
					}
				} else {
					if intValue(debatePlan["open_route_action_count"]) > 0 {
						nextRequirement = currentDebateBatchRequirement("Complete current debate batch before synthesis.", fmt.Sprintf("Still waiting on %d planned route obligation(s) before synthesis.", intValue(debatePlan["open_route_action_count"])))
					} else if intValue(debatePlan["open_round_action_count"]) > 0 {
						nextRequirement = currentDebateBatchRequirement("Complete current debate batch before synthesis.", fmt.Sprintf("Still waiting on %d follow-up debate round(s) before synthesis.", intValue(debatePlan["open_round_action_count"])))
					} else {
						nextRequirement = "Need at least one linked cross-examination position before synthesis."
					}
				}
			case "synthesis":
				gaps = append(gaps, governanceGaps...)
				synthesis := meetingSynthesisAudit(meeting)
				if len(governanceGaps) > 0 {
					nextRequirement = fmt.Sprintf("Repair council governance gaps before output: %s.", strings.Join(governanceGaps, "; "))
				} else if openDebateBatchCount > 0 {
					nextRequirement = currentDebateBatchRequirement("Complete current debate batch before output.", "Complete outstanding debate batches before output.")
				} else if stringValue(synthesis["status"]) == "ready_for_output" {
					newPhase = "output"
					nextRequirement = protocolOutputRequirement(protocol)
				} else {
					nextRequirement = "Supplement synthesis gaps before output."
				}
			case "output":
				nextRequirement = protocolOutputRequirement(protocol)
			default:
				nextRequirement = "Unknown phase; repair the meeting phase explicitly."
			}
			orchestration, _ := toMap(meeting["orchestration"])
			if orchestration == nil {
				orchestration = map[string]any{}
			}
			if newPhase != oldPhase {
				setMeetingPhase(meeting, newPhase, fmt.Sprintf("orchestrate advanced from %s", oldPhase))
			}
			currentPhase := firstNonBlank(stringValue(meeting["phase"]), newPhase, oldPhase)
			autoSynthesis := map[string]any{
				"status": "idle",
			}
			if currentPhase == "synthesis" || currentPhase == "output" {
				autoSynthesis = meetingAutoSynthesize(meeting)
			}
			synthesis := meetingSynthesisAudit(meeting)
			if currentPhase == "synthesis" && len(governanceGaps) > 0 {
				nextRequirement = fmt.Sprintf("Repair council governance gaps before output: %s.", strings.Join(governanceGaps, "; "))
			} else if currentPhase == "synthesis" && openDebateBatchCount > 0 {
				nextRequirement = currentDebateBatchRequirement("Complete current debate batch before output.", "Complete outstanding debate batches before output.")
			} else if currentPhase == "synthesis" && stringValue(synthesis["status"]) == "ready_for_output" {
				setMeetingPhase(meeting, "output", "orchestrate advanced from synthesis")
				currentPhase = "output"
				nextRequirement = protocolOutputRequirement(protocol)
			} else if currentPhase == "synthesis" && stringValue(autoSynthesis["status"]) == "generated" {
				nextRequirement = protocolAutoSynthesisReviewRequirement(protocol)
			}
			orchestration["last_gap_report"] = stringListToAny(gaps)
			orchestration["next_requirement"] = nextRequirement
			orchestration["last_orchestrated_at"] = nowISO()
			orchestration["coverage"] = coverageAudit
			orchestration["wave_plan"] = wavePlan
			orchestration["debate_plan"] = debatePlan
			orchestration["auto_synthesis"] = autoSynthesis
			orchestration["synthesis"] = synthesis
			orchestration["output_packet"] = cloneValue(getKey(autoSynthesis, "output_packet", map[string]any{}))
			orchestration["debate_round"] = maxInt(intValue(orchestration["debate_round"]), intValue(coverageAudit["linked_cross_examination_count"]))
			meeting["orchestration"] = orchestration
			if len(governanceGaps) > 0 {
				meeting["governance_gaps"] = stringListToAny(governanceGaps)
			} else {
				meeting["governance_gaps"] = []any{}
			}
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(map[string]any{
				"meeting":          meeting,
				"old_phase":        oldPhase,
				"new_phase":        firstNonBlank(stringValue(meeting["phase"]), newPhase),
				"gaps":             stringListToAny(gaps),
				"next_requirement": nextRequirement,
				"auto_synthesis":   autoSynthesis,
				"orchestration":    meeting["orchestration"],
			})
		})
	case "add-participant":
		if len(args) < 4 {
			return errors.New("meetings add-participant requires MEETING_ID RUN_ID ROLE")
		}
		meetingID := args[1]
		runID := args[2]
		role := strings.TrimSpace(strings.Join(args[3:], " "))
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			participants := toList(meeting["participants"])
			existsParticipant := false
			for _, item := range participants {
				participant, ok := toMap(item)
				if ok && stringValue(participant["run_id"]) == runID {
					existsParticipant = true
					break
				}
			}
			if !existsParticipant {
				participants = append(participants, map[string]any{
					"run_id":     runID,
					"owner_id":   stringValue(run["owner_id"]),
					"owner_type": stringValue(run["owner_type"]),
					"role":       role,
					"joined_at":  nowISO(),
				})
			}
			meeting["participants"] = participants
			meeting["roles"] = appendUniqueText(toList(meeting["roles"]), role)
			if stringValue(meeting["phase"]) == "" {
				meeting["phase"] = "setup"
			}
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			run["meeting_ref"] = meetingID
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "add-context":
		if len(args) < 3 {
			return errors.New("meetings add-context requires MEETING_ID REF")
		}
		meetingID := args[1]
		ref := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			packet, _ := toMap(meeting["packet"])
			if packet == nil {
				packet = map[string]any{}
			}
			packet["context_refs"] = appendUniqueText(toList(packet["context_refs"]), ref)
			meeting["packet"] = packet
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "add-known-context":
		if len(args) < 3 {
			return errors.New("meetings add-known-context requires MEETING_ID TEXT")
		}
		meetingID := args[1]
		note := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			packet, _ := toMap(meeting["packet"])
			if packet == nil {
				packet = map[string]any{}
			}
			packet["known_context"] = appendUniqueText(toList(packet["known_context"]), note)
			meeting["packet"] = packet
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "lock-decision":
		if len(args) < 3 {
			return errors.New("meetings lock-decision requires MEETING_ID TEXT")
		}
		meetingID := args[1]
		textTokens := []string{}
		reason := ""
		authorityLevel := "L2"
		evidenceRefs := []string{}
		governingMissionID := ""
		bundleID := ""
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--reason="):
				reason = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--reason=")), reason)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), authorityLevel)
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
			case strings.HasPrefix(arg, "--governing-mission="):
				governingMissionID = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--governing-mission=")), governingMissionID)
			case strings.HasPrefix(arg, "--bundle="):
				bundleID = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--bundle=")), bundleID)
			default:
				textTokens = append(textTokens, arg)
			}
		}
		summary := strings.TrimSpace(strings.Join(textTokens, " "))
		if summary == "" {
			return errors.New("meetings lock-decision requires TEXT")
		}
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			lockedDecision := map[string]any{
				"decision_id":     "decision-" + nowTime().UTC().Format("20060102150405"),
				"summary":         summary,
				"reason":          reason,
				"authority_level": authorityLevel,
				"evidence_refs":   stringListToAny(evidenceRefs),
				"locked_at":       nowISO(),
			}
			if governingMissionID != "" {
				lockedDecision["governing_mission_id"] = governingMissionID
			}
			if firstNonBlank(bundleID, strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.bundle_id", "")))) != "" {
				lockedDecision["bundle_id"] = firstNonBlank(bundleID, strings.TrimSpace(stringValue(getKey(meeting, "source_decision_packet.bundle_id", ""))))
			}
			locked := append(toList(meeting["locked_decisions"]), lockedDecision)
			meeting["locked_decisions"] = locked
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "add-position":
		if len(args) < 5 {
			return errors.New("meetings add-position requires MEETING_ID RUN_ID PHASE STANCE TEXT")
		}
		meetingID := args[1]
		runID := args[2]
		phase := args[3]
		stance := args[4]
		positionTokens := []string{}
		evidenceRefs := []string{}
		authorityLevel := "L2"
		claimRef := ""
		dimension := "factual"
		respondsToRunID := ""
		respondsToRole := ""
		respondsToPositionID := ""
		for _, arg := range args[5:] {
			switch {
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), "L2")
			case strings.HasPrefix(arg, "--claim="):
				claimRef = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--claim=")), claimRef)
			case strings.HasPrefix(arg, "--dimension="):
				dimension = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--dimension=")), dimension)
			case strings.HasPrefix(arg, "--responds-to-run="):
				respondsToRunID = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--responds-to-run=")), respondsToRunID)
			case strings.HasPrefix(arg, "--responds-to-role="):
				respondsToRole = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--responds-to-role=")), respondsToRole)
			case strings.HasPrefix(arg, "--responds-to-position="):
				respondsToPositionID = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--responds-to-position=")), respondsToPositionID)
			default:
				positionTokens = append(positionTokens, arg)
			}
		}
		summary := strings.TrimSpace(strings.Join(positionTokens, " "))
		if summary == "" {
			return errors.New("meetings add-position requires TEXT")
		}
		allowedPhases := map[string]bool{"independent_review": true, "cross_examination": true, "synthesis": true}
		if !allowedPhases[phase] {
			return errors.New("meetings add-position phase must be one of: cross_examination, independent_review, synthesis")
		}
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			run, _ := toMap(toList(registry["runs"])[index])
			participants := toList(meeting["participants"])
			role := firstNonBlank(stringValue(run["role"]), stringValue(run["owner_id"]))
			existing := map[string]any(nil)
			for _, item := range participants {
				participant, ok := toMap(item)
				if ok && stringValue(participant["run_id"]) == runID {
					existing = participant
					break
				}
			}
			if existing != nil {
				role = firstNonBlank(stringValue(existing["role"]), role)
			} else {
				participants = append(participants, map[string]any{
					"run_id":     runID,
					"owner_id":   stringValue(run["owner_id"]),
					"owner_type": stringValue(run["owner_type"]),
					"role":       role,
					"joined_at":  nowISO(),
				})
			}
			meeting["participants"] = participants
			positions := toList(meeting["positions"])
			positions = append(positions, map[string]any{
				"position_id":             "position-" + nowTime().UTC().Format("20060102150405"),
				"run_id":                  runID,
				"owner_id":                stringValue(run["owner_id"]),
				"role":                    role,
				"phase":                   phase,
				"stance":                  stance,
				"summary":                 summary,
				"claim_ref":               strings.TrimSpace(claimRef),
				"dimension":               firstNonBlank(dimension, "factual"),
				"evidence_refs":           stringListToAny(evidenceRefs),
				"authority_level":         authorityLevel,
				"responds_to_run_id":      respondsToRunID,
				"responds_to_role":        respondsToRole,
				"responds_to_position_id": respondsToPositionID,
				"recorded_at":             nowISO(),
			})
			meeting["positions"] = positions
			setMeetingPhase(meeting, phase, fmt.Sprintf("position recorded by %s", runID))
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "set-expected-output":
		if len(args) < 3 {
			return errors.New("meetings set-expected-output requires MEETING_ID TEXT")
		}
		meetingID := args[1]
		expectedOutput := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			packet, _ := toMap(meeting["packet"])
			if packet == nil {
				packet = map[string]any{}
			}
			packet["expected_output"] = expectedOutput
			meeting["packet"] = packet
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "add-finding":
		if len(args) < 4 {
			return errors.New("meetings add-finding requires MEETING_ID accepted|rejected|unresolved TEXT")
		}
		meetingID := args[1]
		findingKind := args[2]
		if findingKind != "accepted" && findingKind != "rejected" && findingKind != "unresolved" {
			return errors.New("meetings add-finding kind must be accepted, rejected, or unresolved")
		}
		summaryTokens := []string{}
		reason := ""
		authorityLevel := "L2"
		evidenceRefs := []string{}
		sourceRunID := ""
		sourceRole := ""
		claimRef := ""
		dimension := "factual"
		supportingPositionIDs := []string{}
		challengingPositionIDs := []string{}
		for _, arg := range args[3:] {
			switch {
			case strings.HasPrefix(arg, "--reason="):
				reason = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--reason=")), reason)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), authorityLevel)
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
			case strings.HasPrefix(arg, "--run="):
				sourceRunID = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--run=")), sourceRunID)
			case strings.HasPrefix(arg, "--role="):
				sourceRole = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--role=")), sourceRole)
			case strings.HasPrefix(arg, "--claim="):
				claimRef = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--claim=")), claimRef)
			case strings.HasPrefix(arg, "--dimension="):
				dimension = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--dimension=")), dimension)
			case strings.HasPrefix(arg, "--supports="):
				supportingPositionIDs = append(supportingPositionIDs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--supports=")))...)
			case strings.HasPrefix(arg, "--challenges="):
				challengingPositionIDs = append(challengingPositionIDs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--challenges=")))...)
			default:
				summaryTokens = append(summaryTokens, arg)
			}
		}
		summary := strings.TrimSpace(strings.Join(summaryTokens, " "))
		if summary == "" {
			return errors.New("meetings add-finding requires TEXT")
		}
		if authorityLevel != "" && !validAuthorityLevel(authorityLevel) {
			return fmt.Errorf("invalid authority level: %s", authorityLevel)
		}
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			findings, _ := toMap(meeting["findings"])
			if findings == nil {
				findings = map[string]any{}
			}
			bucket := append(toList(findings[findingKind]), map[string]any{
				"finding_id":               "finding-" + nowTime().UTC().Format("20060102150405"),
				"summary":                  summary,
				"reason":                   reason,
				"authority_level":          authorityLevel,
				"evidence_refs":            stringListToAny(evidenceRefs),
				"source_run_id":            sourceRunID,
				"source_role":              sourceRole,
				"claim_ref":                claimRef,
				"dimension":                dimension,
				"supporting_position_ids":  stringListToAny(dedupeStrings(supportingPositionIDs)),
				"challenging_position_ids": stringListToAny(dedupeStrings(challengingPositionIDs)),
				"recorded_at":              nowISO(),
			})
			findings[findingKind] = bucket
			meeting["findings"] = findings
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "add-action":
		if len(args) < 3 {
			return errors.New("meetings add-action requires MEETING_ID TEXT")
		}
		meetingID := args[1]
		actionText := strings.TrimSpace(strings.Join(args[2:], " "))
		return withHarnessLock(func() error {
			data := readActiveMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("active meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			actions := append(toList(meeting["action_items"]), map[string]any{
				"action_id":   "action-" + nowTime().UTC().Format("20060102150405"),
				"summary":     actionText,
				"recorded_at": nowISO(),
				"state":       "open",
			})
			meeting["action_items"] = actions
			if err := writeActiveMeeting(meeting); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	case "resolve-action":
		if len(args) < 3 {
			return errors.New("meetings resolve-action requires MEETING_ID ACTION_ID [STATE]")
		}
		meetingID := args[1]
		actionID := args[2]
		actionState := "done"
		if len(args) > 3 {
			actionState = args[3]
		}
		if actionState != "done" && actionState != "cancelled" && actionState != "open" {
			return errors.New("meetings resolve-action state must be one of: open, done, cancelled")
		}
		return withHarnessLock(func() error {
			data := readMeeting(meetingID)
			if data == nil {
				return fmt.Errorf("meeting not found: %s", meetingID)
			}
			meeting, _ := toMap(data["meeting"])
			actions := toList(meeting["action_items"])
			found := false
			for index, item := range actions {
				action, ok := toMap(item)
				if ok && stringValue(action["action_id"]) == actionID {
					action["state"] = actionState
					action["updated_at"] = nowISO()
					if actionState == "done" || actionState == "cancelled" {
						action["closed_at"] = nowISO()
					}
					actions[index] = action
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("action not found: %s", actionID)
			}
			meeting["action_items"] = actions
			meeting = refreshMeetingOutputPacket(meeting)
			if exists(meetingPath(meetingID, false)) {
				if err := writeActiveMeeting(meeting); err != nil {
					return err
				}
			} else {
				meeting["updated_at"] = nowISO()
				if err := dumpData(meetingPath(meetingID, true), map[string]any{"version": 1, "meeting": meeting}); err != nil {
					return err
				}
			}
			registry := readRegistry()
			registry = syncAdoptedFollowUpLifecycleForMeeting(registry, meeting)
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(meeting)
		})
	default:
		return fmt.Errorf("unknown meetings subcommand: %s", subcommand)
	}
}

func cmdTeammates(args []string) error {
	if len(args) < 1 {
		return errors.New("teammates requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	subcommand := args[0]
	registry := readRegistry()
	switch subcommand {
	case "list":
		includeAll := false
		scope := ""
		for _, arg := range args[1:] {
			if arg == "--all" {
				includeAll = true
			} else if !strings.HasPrefix(arg, "--") && scope == "" {
				scope = arg
			}
		}
		runs := []any{}
		for _, item := range toList(registry["runs"]) {
			run, _ := toMap(item)
			if stringValue(run["owner_type"]) != "teammate" {
				continue
			}
			if !includeAll && !activeRun(run) {
				continue
			}
			if scope != "" && stringValue(run["parent_run_id"]) != scope && stringValue(run["run_id"]) != scope && stringValue(run["mission_id"]) != scope {
				continue
			}
			runs = append(runs, run)
		}
		return printJSON(runs)
	case "get":
		if len(args) < 2 {
			return errors.New("teammates get requires RUN_ID")
		}
		index := findRunIndex(registry, args[1])
		if index < 0 {
			return fmt.Errorf("teammate run not found: %s", args[1])
		}
		run, _ := toMap(toList(registry["runs"])[index])
		if stringValue(run["owner_type"]) != "teammate" {
			return fmt.Errorf("teammate run not found: %s", args[1])
		}
		return printJSON(run)
	default:
		return cmdTeammatesRemainder(args)
	}
}

func cmdTeammatesRemainder(args []string) error {
	subcommand := args[0]
	switch subcommand {
	case "assign":
		if len(args) < 5 {
			return errors.New("teammates assign requires PARENT_RUN_ID TEAMMATE_ID ROLE OBJECTIVE")
		}
		parentRunID := args[1]
		teammateID := args[2]
		role := args[3]
		positionals := []string{}
		assignedScope := ""
		executionSurface := "worktree-thread"
		mode := "managed"
		dependsOn := []string{}
		ownershipClaimsInput := []string{}
		overlapPolicyValue := "exclusive"
		siblingIDs := []string{}
		authorityRefs := []string{}
		forbiddenSources := []string{}
		reportingExpectations := []string{}
		for _, arg := range args[4:] {
			switch {
			case strings.HasPrefix(arg, "--scope="):
				assignedScope = strings.TrimSpace(strings.TrimPrefix(arg, "--scope="))
			case strings.HasPrefix(arg, "--surface="):
				executionSurface = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--surface=")), executionSurface)
			case strings.HasPrefix(arg, "--mode="):
				mode = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--mode=")), mode)
			case strings.HasPrefix(arg, "--depends-on="):
				dependsOn = append(dependsOn, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--depends-on=")))...)
			case strings.HasPrefix(arg, "--owns="):
				ownershipClaimsInput = append(ownershipClaimsInput, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--owns=")))...)
			case arg == "--allow-overlap":
				overlapPolicyValue = "shared"
			case strings.HasPrefix(arg, "--sibling="):
				siblingIDs = append(siblingIDs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--sibling=")))...)
			case strings.HasPrefix(arg, "--allow="):
				authorityRefs = append(authorityRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--allow=")))...)
			case strings.HasPrefix(arg, "--forbid="):
				forbiddenSources = append(forbiddenSources, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--forbid=")))...)
			case strings.HasPrefix(arg, "--report="):
				reportingExpectations = append(reportingExpectations, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--report=")))...)
			default:
				positionals = append(positionals, arg)
			}
		}
		objective := strings.TrimSpace(strings.Join(positionals, " "))
		if objective == "" {
			return errors.New("teammates assign requires OBJECTIVE")
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			parentIndex := findRunIndex(registry, parentRunID)
			if parentIndex < 0 {
				return fmt.Errorf("parent run not found: %s", parentRunID)
			}
			runs := toList(registry["runs"])
			parentRun, _ := toMap(runs[parentIndex])
			if !activeRun(parentRun) {
				return errors.New("parent run is not active")
			}
			for _, item := range runs {
				run, _ := toMap(item)
				if activeTeammateRun(run) && stringValue(run["parent_run_id"]) == parentRunID && stringValue(run["owner_id"]) == teammateID {
					return fmt.Errorf("active teammate assignment already exists for %s", teammateID)
				}
			}
			teammateRunID := fmt.Sprintf("run-%s-teammate-%s", nowTime().UTC().Format("20060102150405"), slugify(teammateID))
			contractPath := filepath.Join(runtimeRoot(), "contracts", teammateRunID+".yml")
			contractTemplateRoot, _ := toMap(loadData(contractTemplatePath(), map[string]any{}))
			contractTemplate, _ := toMap(contractTemplateRoot["contract"])
			successCriteria := reportingExpectations
			if len(successCriteria) == 0 {
				successCriteria = []string{"Report progress, blockers, and verification outcomes."}
			}
			contractPayload, _ := toMap(mergeDicts(contractTemplate, map[string]any{
				"mission_id":            stringValue(parentRun["mission_id"]),
				"issue_number":          intValue(parentRun["issue_number"]),
				"title":                 objective,
				"objective":             objective,
				"success_criteria":      stringListToAny(successCriteria),
				"sub_tasks":             []any{},
				"depends_on":            stringListToAny(dedupeStrings(dependsOn)),
				"desired_state":         defaultDesiredStateForRun("teammate"),
				"branch":                stringValue(parentRun["branch"]),
				"worktree_path":         stringValue(parentRun["worktree_path"]),
				"context_files":         stringListToAny(authorityRefs),
				"verification_commands": []any{},
				"notes": map[string]any{
					"delegated_from":    parentRunID,
					"teammate_id":       teammateID,
					"role":              role,
					"assigned_scope":    firstNonBlank(assignedScope, objective),
					"forbidden_sources": stringListToAny(forbiddenSources),
				},
			}))
			setOwnershipSelectors(contractPayload, dedupeStrings(append([]string{"scope:" + slugify(firstNonBlank(assignedScope, objective))}, ownershipClaimsInput...)))
			freshness, _ := toMap(contractPayload["freshness"])
			if freshness == nil {
				freshness = map[string]any{}
			}
			fillIfBlank(freshness, "created_at", nowISO())
			contractPayload["freshness"] = freshness
			if err := dumpData(contractPath, map[string]any{"version": 1, "contract": contractPayload}); err != nil {
				return err
			}
			runTemplateRoot, _ := toMap(loadData(runTemplatePath(), map[string]any{}))
			runTemplate, _ := toMap(runTemplateRoot["run"])
			timestamp := nowISO()
			siblingRuns := childRunsForParent(registry, parentRunID)
			siblingRunIDs := append([]string{}, siblingIDs...)
			communicationPeers := []string{}
			for _, run := range siblingRuns {
				runID := stringValue(run["run_id"])
				siblingRunIDs = append(siblingRunIDs, runID)
				communicationPeers = append(communicationPeers, runID)
			}
			authorityRefList := authorityRefs
			if len(authorityRefList) == 0 {
				authorityRefList = []string{stringValue(parentRun["contract_ref"])}
			}
			ownershipClaims := dedupeStrings(append([]string{"scope:" + slugify(firstNonBlank(assignedScope, objective))}, ownershipClaimsInput...))
			reportingList := reportingExpectations
			if len(reportingList) == 0 {
				reportingList = []string{"Report progress, blockers, and verification outcomes."}
			}
			teammateRun, _ := toMap(mergeDicts(runTemplate, map[string]any{
				"run_id":                 teammateRunID,
				"mission_id":             stringValue(parentRun["mission_id"]),
				"issue_number":           intValue(parentRun["issue_number"]),
				"parent_run_id":          parentRunID,
				"leader_id":              stringValue(parentRun["leader_id"]),
				"owner_type":             "teammate",
				"owner_id":               teammateID,
				"delegated_by":           firstNonBlank(stringValue(parentRun["owner_id"]), stringValue(parentRun["leader_id"])),
				"role":                   role,
				"state":                  "ready",
				"mode":                   mode,
				"execution_surface":      executionSurface,
				"overlap_policy":         overlapPolicyValue,
				"branch":                 stringValue(parentRun["branch"]),
				"worktree_path":          stringValue(parentRun["worktree_path"]),
				"contract_ref":           contractPath,
				"objective":              objective,
				"assigned_scope":         firstNonBlank(assignedScope, objective),
				"authority_refs":         stringListToAny(authorityRefList),
				"forbidden_sources":      stringListToAny(forbiddenSources),
				"reporting_expectations": stringListToAny(reportingList),
				"depends_on":             stringListToAny(dependsOn),
				"dependency_policy":      "all",
				"desired_state":          defaultDesiredStateForRun("teammate"),
				"siblings":               stringListToAny(siblingRunIDs),
				"communication_peers":    stringListToAny(communicationPeers),
				"latest_summary":         "Delegated and ready for bounded execution",
				"next_expected_action":   "Acknowledge the assignment, gather local evidence, and report back.",
				"created_at":             timestamp,
				"started_at":             timestamp,
				"last_reported_at":       timestamp,
				"updated_at":             timestamp,
			}))
			setOwnershipSelectors(teammateRun, ownershipClaims)
			if conflicts := overlappingOwnershipRuns(registry, ownershipClaims, overlapPolicyValue, []string{parentRunID}); len(conflicts) > 0 {
				return fmt.Errorf("ownership overlap detected for teammate assignment; retry with --allow-overlap or change --owns selectors: %v", conflicts)
			}
			runs = append(runs, teammateRun)
			for _, sibling := range siblingRuns {
				sibling["siblings"] = appendUniqueText(toList(sibling["siblings"]), teammateRunID)
				sibling["communication_peers"] = appendUniqueText(toList(sibling["communication_peers"]), teammateRunID)
				sibling["updated_at"] = nowISO()
			}
			parentRun["mode"] = "parallel"
			parentRun["latest_summary"] = "Delegated bounded task to " + teammateID
			parentRun["next_expected_action"] = "Monitor teammate reports, coordination signals, and meeting triggers."
			parentRun["updated_at"] = nowISO()
			runs[parentIndex] = parentRun
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			assignmentSignal, err := emitSignal("progress", teammateRun, "Delegated teammate assignment: "+teammateID, fmt.Sprintf("%s assigned as %s: %s", teammateID, role, objective), []string{parentRunID, teammateRunID, stringValue(parentRun["mission_id"])}, "sibling", "high", "L2", []string{contractPath}, 12.0)
			if err != nil {
				return err
			}
			return printJSON(map[string]any{"run": teammateRun, "signal": assignmentSignal})
		})
	case "link":
		if len(args) < 3 {
			return errors.New("teammates link requires RUN_A RUN_B [REASON]")
		}
		runAID := args[1]
		runBID := args[2]
		reason := "Leader approved teammate coordination."
		if len(args) > 3 {
			reason = strings.TrimSpace(strings.Join(args[3:], " "))
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			indexA := findRunIndex(registry, runAID)
			indexB := findRunIndex(registry, runBID)
			if indexA < 0 || indexB < 0 {
				return errors.New("teammates link requires two existing runs")
			}
			runs := toList(registry["runs"])
			runA, _ := toMap(runs[indexA])
			runB, _ := toMap(runs[indexB])
			for _, key := range []string{"siblings", "communication_peers"} {
				runA[key] = appendUniqueText(toList(runA[key]), runBID)
				runB[key] = appendUniqueText(toList(runB[key]), runAID)
			}
			runA["updated_at"] = nowISO()
			runB["updated_at"] = nowISO()
			runs[indexA] = runA
			runs[indexB] = runB
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{"linked": []any{runAID, runBID}, "reason": reason})
		})
	case "report":
		if len(args) < 3 {
			return errors.New("teammates report requires RUN_ID SUMMARY")
		}
		runID := args[1]
		kind := "progress"
		title := ""
		confidence := "medium"
		scopeType := "local"
		authorityLevel := "L2"
		freshHours := 12.0
		evidenceRefs := []string{}
		appliesTo := []string{}
		nextAction := ""
		summaryTokens := []string{}
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--kind="):
				kind = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--kind=")), kind)
			case strings.HasPrefix(arg, "--title="):
				title = strings.TrimSpace(strings.TrimPrefix(arg, "--title="))
			case strings.HasPrefix(arg, "--confidence="):
				confidence = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--confidence=")), confidence)
			case strings.HasPrefix(arg, "--scope="):
				scopeType = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--scope=")), scopeType)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), authorityLevel)
			case strings.HasPrefix(arg, "--fresh-hours="):
				freshHours, _ = strconv.ParseFloat(firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--fresh-hours=")), "12"), 64)
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
			case strings.HasPrefix(arg, "--applies="):
				appliesTo = append(appliesTo, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--applies=")))...)
			case strings.HasPrefix(arg, "--next="):
				nextAction = strings.TrimSpace(strings.TrimPrefix(arg, "--next="))
			default:
				summaryTokens = append(summaryTokens, arg)
			}
		}
		summary := strings.TrimSpace(strings.Join(summaryTokens, " "))
		if summary == "" {
			return errors.New("teammates report requires SUMMARY")
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if stringValue(run["owner_type"]) != "teammate" {
				return errors.New("teammates report requires a teammate-owned run")
			}
			run["latest_summary"] = summary
			run["last_reported_at"] = nowISO()
			if nextAction != "" {
				run["next_expected_action"] = nextAction
			}
			run["updated_at"] = nowISO()
			targets := appliesTo
			if len(targets) == 0 {
				targets = []string{stringValue(run["parent_run_id"]), stringValue(run["mission_id"])}
			}
			signal, err := emitSignal(kind, run, firstNonBlank(title, "Teammate report: "+stringValue(run["owner_id"])), summary, targets, scopeType, confidence, authorityLevel, evidenceRefs, freshHours)
			if err != nil {
				return err
			}
			run["active_signals"] = appendUniqueText(toList(run["active_signals"]), stringValue(signal["signal_id"]))
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{"run": run, "signal": signal})
		})
	case "message":
		if len(args) < 5 {
			return errors.New("teammates message requires SOURCE_RUN_ID TARGET_RUN_ID KIND SUMMARY")
		}
		sourceRunID := args[1]
		targetRunID := args[2]
		kind := args[3]
		title := ""
		confidence := "medium"
		authorityLevel := "L2"
		freshHours := 6.0
		evidenceRefs := []string{}
		summaryTokens := []string{}
		for _, arg := range args[4:] {
			switch {
			case strings.HasPrefix(arg, "--title="):
				title = strings.TrimSpace(strings.TrimPrefix(arg, "--title="))
			case strings.HasPrefix(arg, "--confidence="):
				confidence = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--confidence=")), confidence)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), authorityLevel)
			case strings.HasPrefix(arg, "--fresh-hours="):
				freshHours, _ = strconv.ParseFloat(firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--fresh-hours=")), "6"), 64)
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
			default:
				summaryTokens = append(summaryTokens, arg)
			}
		}
		summary := strings.TrimSpace(strings.Join(summaryTokens, " "))
		if summary == "" {
			return errors.New("teammates message requires SUMMARY")
		}
		return withHarnessLock(func() error {
			registry := readRegistry()
			sourceIndex := findRunIndex(registry, sourceRunID)
			targetIndex := findRunIndex(registry, targetRunID)
			if sourceIndex < 0 || targetIndex < 0 {
				return errors.New("teammates message requires valid source and target runs")
			}
			runs := toList(registry["runs"])
			sourceRun, _ := toMap(runs[sourceIndex])
			targetRun, _ := toMap(runs[targetIndex])
			if stringValue(sourceRun["owner_type"]) != "teammate" {
				return errors.New("teammates message source must be a teammate run")
			}
			if stringValue(targetRun["owner_type"]) != "teammate" {
				return errors.New("teammates message target must be a teammate run")
			}
			if !siblingLinked(sourceRun, targetRunID) {
				return errors.New("teammates message requires leader-approved linkage via teammates link or shared delegation")
			}
			signal, err := emitSignal(kind, sourceRun, firstNonBlank(title, "Teammate message to "+stringValue(targetRun["owner_id"])), summary, []string{targetRunID, stringValue(targetRun["mission_id"])}, "sibling", confidence, authorityLevel, evidenceRefs, freshHours)
			if err != nil {
				return err
			}
			sourceRun["active_signals"] = appendUniqueText(toList(sourceRun["active_signals"]), stringValue(signal["signal_id"]))
			targetRun["active_signals"] = appendUniqueText(toList(targetRun["active_signals"]), stringValue(signal["signal_id"]))
			sourceRun["last_reported_at"] = nowISO()
			sourceRun["updated_at"] = nowISO()
			targetRun["updated_at"] = nowISO()
			runs[sourceIndex] = sourceRun
			runs[targetIndex] = targetRun
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{"signal": signal, "source_run": sourceRun, "target_run": targetRun})
		})
	case "escalate":
		if len(args) < 3 {
			return errors.New("teammates escalate requires RUN_ID QUESTION [--materialize] [--protocol=evidence-review|incident-command|decision-record]")
		}
		runID := args[1]
		materialize := false
		protocol := ""
		questionTokens := []string{}
		for _, arg := range args[2:] {
			if arg == "--materialize" {
				materialize = true
			} else if strings.HasPrefix(arg, "--protocol=") {
				protocol = strings.TrimSpace(strings.TrimPrefix(arg, "--protocol="))
			} else {
				questionTokens = append(questionTokens, arg)
			}
		}
		question := strings.TrimSpace(strings.Join(questionTokens, " "))
		return withHarnessLock(func() error {
			registry := readRegistry()
			index := findRunIndex(registry, runID)
			if index < 0 {
				return fmt.Errorf("run not found: %s", runID)
			}
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if stringValue(run["owner_type"]) != "teammate" {
				return errors.New("teammates escalate requires a teammate-owned run")
			}
			signal, err := emitSignal("requires_meeting", run, "Meeting requested by "+stringValue(run["owner_id"]), question, []string{stringValue(run["parent_run_id"]), stringValue(run["mission_id"])}, "sibling", "high", "L2", []string{stringValue(run["contract_ref"])}, 12.0)
			if err != nil {
				return err
			}
			run["status_reason"] = question
			run["last_reported_at"] = nowISO()
			run["updated_at"] = nowISO()
			run["active_signals"] = appendUniqueText(toList(run["active_signals"]), stringValue(signal["signal_id"]))
			payload := map[string]any{"signal": signal, "run": run}
			if materialize {
				meetingID := newMeetingID(stringValue(run["mission_id"]))
				meeting, err := materializeMeeting(meetingID, run, question, stringValue(run["owner_id"]), protocol)
				if err != nil {
					return err
				}
				meeting["phase"] = "setup"
				if err := writeActiveMeeting(meeting); err != nil {
					return err
				}
				run["meeting_ref"] = meetingID
				run["state"] = "waiting"
				run["next_expected_action"] = "Participate in the meeting with bounded evidence and role discipline."
				payload["meeting"] = meeting
			}
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(payload)
		})
	default:
		return fmt.Errorf("unknown teammates subcommand: %s", subcommand)
	}
}

func cmdMaintenance(args []string) error {
	if len(args) < 1 {
		return errors.New("maintenance requires a subcommand")
	}
	if err := ensureRuntimeDirs(); err != nil {
		return err
	}
	switch args[0] {
	case "normalize":
		return withHarnessLock(func() error {
			registry := readRegistry()
			missionRegistry := readMissionRegistry()
			log := readSignalLog()
			normalizedRuns := []string{}
			normalizedSignals := []string{}
			normalizedContracts := []string{}
			emptySignalLog := map[string]any{"signals": []any{}}
			runs := toList(registry["runs"])
			for index, item := range runs {
				run, _ := toMap(item)
				changed := false
				timestamp := firstNonBlank(stringValue(run["updated_at"]), nowISO())
				if stringValue(run["created_at"]) == "" {
					run["created_at"] = timestamp
					changed = true
				}
				state := stringValue(run["state"])
				if (state == "ready" || state == "executing" || state == "waiting" || state == "blocked") && stringValue(run["started_at"]) == "" {
					run["started_at"] = firstNonBlank(stringValue(run["created_at"]), timestamp)
					changed = true
				}
				if (state == "ready" || state == "executing" || state == "waiting" || state == "blocked") && stringValue(run["last_reported_at"]) == "" {
					run["last_reported_at"] = firstNonBlank(stringValue(run["updated_at"]), stringValue(run["started_at"]), nowISO())
					changed = true
				}
				if len(ownershipClaims(run)) == 0 {
					claims := []string{}
					if missionID := stringValue(run["mission_id"]); missionID != "" {
						claims = append(claims, "mission:"+missionID)
					}
					if branch := stringValue(run["branch"]); branch != "" {
						claims = append(claims, "branch:"+branch)
					}
					if assignedScope := strings.TrimSpace(stringValue(run["assigned_scope"])); assignedScope != "" {
						claims = append(claims, "scope:"+slugify(assignedScope))
					}
					if len(claims) > 0 {
						setOwnershipSelectors(run, claims)
						changed = true
					}
				} else if len(toList(run["ownership_selectors"])) == 0 {
					setOwnershipSelectors(run, ownershipClaims(run))
					changed = true
				}
				if _, ok := toMap(run["desired_state"]); !ok {
					run["desired_state"] = defaultDesiredStateForRun(stringValue(run["owner_type"]))
					changed = true
				}
				if strings.TrimSpace(stringValue(run["dependency_policy"])) == "" {
					run["dependency_policy"] = "all"
					changed = true
				}
				if changed {
					run["updated_at"] = nowISO()
					normalizedRuns = append(normalizedRuns, stringValue(run["run_id"]))
					runs[index] = run
				}
				authorityAudit, _ := toMap(runAuthorityAudit(run, registry, emptySignalLog))
				if intValue(authorityAudit["contract_gap_count"]) > 0 {
					if err := persistContractForRun(run); err != nil {
						return err
					}
					normalizedContracts = append(normalizedContracts, stringValue(run["run_id"]))
				}
			}
			registry["runs"] = runs
			signals := toList(log["signals"])
			for index, item := range signals {
				signal, _ := toMap(item)
				normalizedKind := normalizeSignalKind(stringValue(signal["kind"]))
				if normalizedKind != stringValue(signal["kind"]) {
					signal["kind"] = normalizedKind
					normalizedSignals = append(normalizedSignals, stringValue(signal["signal_id"]))
				}
				if stringValue(signal["emitted_at"]) == "" {
					signal["emitted_at"] = firstNonBlank(stringValue(signal["fresh_until"]), nowISO())
					normalizedSignals = append(normalizedSignals, stringValue(signal["signal_id"]))
				}
				signals[index] = signal
			}
			log["signals"] = signals
			needsMissionSync := len(toList(missionRegistry["missions"])) == 0 && len(runs) > 0
			if !needsMissionSync {
				for _, item := range toList(missionRegistry["missions"]) {
					mission, _ := toMap(item)
					if stringValue(mission["authoritative_run_id"]) == "" {
						needsMissionSync = true
						break
					}
				}
			}
			if len(normalizedRuns) > 0 || needsMissionSync {
				if err := writeRegistry(registry); err != nil {
					return err
				}
			}
			if len(normalizedSignals) > 0 {
				if err := writeSignalLog(log); err != nil {
					return err
				}
			}
			return printJSON(map[string]any{
				"normalized_run_ids":    stringListToAny(normalizedRuns),
				"normalized_signal_ids": stringListToAny(normalizedSignals),
				"normalized_contracts":  stringListToAny(dedupeStrings(normalizedContracts)),
				"run_count":             len(normalizedRuns),
				"signal_count":          len(normalizedSignals),
				"contract_count":        len(dedupeStrings(normalizedContracts)),
			})
		})
	case "reconcile":
		return withHarnessLock(func() error {
			registry := readRegistry()
			log := readSignalLog()
			activeRunIDs := map[string]bool{}
			for _, item := range toList(registry["runs"]) {
				run, _ := toMap(item)
				if activeRun(run) {
					activeRunIDs[stringValue(run["run_id"])] = true
				}
			}
			reconciledSignals := []string{}
			adjudicatedSignals := []string{}
			signals := toList(log["signals"])
			for index, item := range signals {
				signal, _ := toMap(item)
				if stringValue(signal["state"]) != "active" {
					continue
				}
				runID := stringValue(signal["run_id"])
				if runID != "" && !activeRunIDs[runID] {
					signal["state"] = "resolved"
					signal["resolution_note"] = "Reconciled orphaned signal: referenced run is missing."
					signal["resolved_at"] = nowISO()
					signals[index] = signal
					reconciledSignals = append(reconciledSignals, stringValue(signal["signal_id"]))
					continue
				}
				if superseded, resolution := signalSupersession(signal, registry); superseded {
					signal["state"] = "resolved"
					signal["resolved_at"] = nowISO()
					for key, value := range resolution {
						signal[key] = value
					}
					signals[index] = signal
					adjudicatedSignals = append(adjudicatedSignals, stringValue(signal["signal_id"]))
				}
			}
			log["signals"] = signals
			reconciledMeetings := []string{}
			for _, item := range listMeetings(false) {
				meeting, _ := toMap(item)
				if stringValue(meeting["state"]) != "active" {
					continue
				}
				runID := stringValue(meeting["run_id"])
				if runID != "" && !activeRunIDs[runID] {
					meeting["state"] = "closed"
					meeting["updated_at"] = nowISO()
					meeting["closed_at"] = nowISO()
					meeting["resolution_note"] = "Reconciled orphaned meeting: referenced run is missing."
					if err := dumpData(meetingPath(stringValue(meeting["meeting_id"]), true), map[string]any{"version": 1, "meeting": meeting}); err != nil {
						return err
					}
					_ = os.Remove(meetingPath(stringValue(meeting["meeting_id"]), false))
					if _, err := updateBacklogMeeting(stringValue(meeting["meeting_id"]), map[string]any{
						"state":           "closed",
						"closed_at":       meeting["closed_at"],
						"resolution_note": meeting["resolution_note"],
						"meeting_ref":     meetingPath(stringValue(meeting["meeting_id"]), true),
						"updated_at":      nowISO(),
					}); err != nil {
						return err
					}
					reconciledMeetings = append(reconciledMeetings, stringValue(meeting["meeting_id"]))
				}
			}
			if len(reconciledSignals) > 0 || len(adjudicatedSignals) > 0 {
				if err := writeSignalLog(log); err != nil {
					return err
				}
			}
			return printJSON(map[string]any{
				"reconciled_signal_ids":  stringListToAny(reconciledSignals),
				"adjudicated_signal_ids": stringListToAny(adjudicatedSignals),
				"reconciled_meeting_ids": stringListToAny(reconciledMeetings),
				"signal_count":           len(reconciledSignals) + len(adjudicatedSignals),
				"meeting_count":          len(reconciledMeetings),
			})
		})
	default:
		return fmt.Errorf("unknown maintenance subcommand: %s", args[0])
	}
}
