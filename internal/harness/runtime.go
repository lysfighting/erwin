package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func openMeetingActions() []any {
	actions := []any{}
	for _, item := range listMeetings(true) {
		meeting, _ := toMap(item)
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
	}
	switch kind {
	case "file", "dir", "module":
		selectorValue = strings.TrimSpace(strings.TrimSuffix(selectorValue, "/"))
		selectorValue = strings.TrimPrefix(selectorValue, "./")
	case "scope", "interface":
		selectorValue = slugify(selectorValue)
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

func selectorRecordFor(raw string) map[string]any {
	return parseOwnershipSelector(raw)
}

func pathOverlap(leftKind, leftValue, rightKind, rightValue string) bool {
	if leftValue == "" || rightValue == "" {
		return false
	}
	sameOrPrefix := func(prefix, candidate string) bool {
		return candidate == prefix || strings.HasPrefix(candidate, prefix+"/")
	}
	switch {
	case leftKind == "file" && rightKind == "file":
		return leftValue == rightValue
	case leftKind == "file" && (rightKind == "dir" || rightKind == "module"):
		return sameOrPrefix(rightValue, leftValue)
	case rightKind == "file" && (leftKind == "dir" || leftKind == "module"):
		return sameOrPrefix(leftValue, rightValue)
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
	switch leftFamily {
	case "anchor", "semantic":
		return leftKind == rightKind && leftValue == rightValue
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
		for _, candidate := range candidateClaims {
			for _, claim := range ownershipClaims(run) {
				if ownershipSelectorsOverlap(candidate, claim) {
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
	switch stringValue(signal["kind"]) {
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
		snapshot["dirty"] = strings.TrimSpace(out) != ""
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
			"claim_ref":        gap,
			"resolution_kind":  "desired_state_unmet",
			"authority_level":  "L3",
			"source_ref":       "desired_state:" + runID,
			"winning_source":   "run:" + runID,
			"note":             fmt.Sprintf("Current run state does not yet satisfy desired state %s.", gap),
			"desired_state_gap": true,
		})
	}

	adoptionGaps := []any{}
	if meeting := latestClosedMeetingForRun(runID); meeting != nil && meetingHasOutput(meeting) {
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
		if nextAction := stringValue(outcome["next_action"]); nextAction != "" && strings.TrimSpace(stringValue(run["next_expected_action"])) != nextAction {
			adoptionGaps = append(adoptionGaps, map[string]any{
				"claim_ref":       "next_expected_action",
				"resolution_kind": "adoption_pending",
				"authority_level": "L3",
				"source_ref":      "meeting:" + stringValue(meeting["meeting_id"]),
				"winning_source":  "meeting:" + stringValue(meeting["meeting_id"]),
				"current_value":   stringValue(run["next_expected_action"]),
				"expected_value":  nextAction,
				"note":            "Closed meeting next action has not been adopted into the run state.",
			})
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
					"resolution_kind": "contract_stale",
					"source_ref":      "contract:" + runID,
					"winning_source":  "run:" + runID,
					"current_value":   stringValue(contractDesired[field]),
					"expected_value":  stringValue(runDesired[field]),
					"note":            fmt.Sprintf("Contract desired state for %s is stale relative to the run control plane.", field),
				})
			}
		}
		if stringValue(contract["branch"]) != stringValue(run["branch"]) {
			contractGaps = append(contractGaps, map[string]any{
				"claim_ref":       "branch",
				"resolution_kind": "contract_stale",
				"source_ref":      "contract:" + runID,
				"winning_source":  "run:" + runID,
				"current_value":   stringValue(contract["branch"]),
				"expected_value":  stringValue(run["branch"]),
				"note":            "Contract branch does not match the current run branch.",
			})
		}
		if stringValue(contract["worktree_path"]) != stringValue(run["worktree_path"]) {
			contractGaps = append(contractGaps, map[string]any{
				"claim_ref":       "worktree_path",
				"resolution_kind": "contract_stale",
				"source_ref":      "contract:" + runID,
				"winning_source":  "run:" + runID,
				"current_value":   stringValue(contract["worktree_path"]),
				"expected_value":  stringValue(run["worktree_path"]),
				"note":            "Contract worktree path does not match the current run provisioning state.",
			})
		}
		if strings.Join(ownershipClaims(contract), "|") != strings.Join(ownershipClaims(run), "|") {
			contractGaps = append(contractGaps, map[string]any{
				"claim_ref":       "ownership_claims",
				"resolution_kind": "contract_stale",
				"source_ref":      "contract:" + runID,
				"winning_source":  "run:" + runID,
				"current_value":   strings.Join(ownershipClaims(contract), "|"),
				"expected_value":  strings.Join(ownershipClaims(run), "|"),
				"note":            "Contract ownership claims no longer match the run ownership boundary.",
			})
		}
	}

	liveStateGaps := []any{}
	if liveState := worktreeAuthoritySnapshot(run); liveState != nil {
		if liveBranch := strings.TrimSpace(stringValue(liveState["branch"])); liveBranch != "" && liveBranch != stringValue(run["branch"]) {
			liveStateGaps = append(liveStateGaps, map[string]any{
				"claim_ref":       "live_branch",
				"resolution_kind": "live_state_disagrees",
				"source_ref":      "worktree:" + stringValue(run["worktree_path"]),
				"winning_source":  "worktree:" + stringValue(run["worktree_path"]),
				"current_value":   stringValue(run["branch"]),
				"expected_value":  liveBranch,
				"note":            "Live worktree branch does not match the run branch record.",
			})
		}
	}

	activeSignals := activeSignalsForRun(signalLog, runID)
	signalSources := []any{}
	conflicts := []any{}
	supersededSignals := []any{}
	for _, signal := range activeSignals {
		signalID := stringValue(signal["signal_id"])
		emittedAt := firstNonBlank(stringValue(signal["emitted_at"]), stringValue(signal["updated_at"]))
		signalLevel := firstNonBlank(stringValue(signal["authority_level"]), "L4")
		signalRank := authorityRank(signalLevel)
		signalSources = append(signalSources, map[string]any{
			"source_ref":       "signal:" + signalID,
			"kind":             stringValue(signal["kind"]),
			"summary":          stringValue(signal["summary"]),
			"authority_level":  signalLevel,
			"authority_rank":   signalRank,
			"updated_at":       emittedAt,
			"directness":       "local_signal",
			"resolution_state": stringValue(signal["state"]),
		})
		if superseded, resolution := signalSupersession(signal, registry); superseded {
			resolution["signal_id"] = signalID
			resolution["signal_kind"] = stringValue(signal["kind"])
			resolution["authority_level"] = signalLevel
			resolution["emitted_at"] = emittedAt
			supersededSignals = append(supersededSignals, resolution)
			continue
		}

		conflict := map[string]any(nil)
		switch stringValue(signal["kind"]) {
		case "blocked":
			if stringValue(run["state"]) != "blocked" && stringValue(run["state"]) != "waiting" {
				if signalRank > authorityRank(runLevel) || (signalRank == authorityRank(runLevel) && isoAfter(emittedAt, runUpdatedAt)) {
					conflict = map[string]any{
						"claim_ref":        "execution_state",
						"signal_id":        signalID,
						"signal_kind":      stringValue(signal["kind"]),
						"source_ref":       "signal:" + signalID,
						"winning_source":   "signal:" + signalID,
						"losing_source":    "run:" + runID,
						"authority_level":  signalLevel,
						"resolution_kind":  "run_state_overridden",
						"current_value":    stringValue(run["state"]),
						"expected_value":   "blocked|waiting",
						"note":             "Active blocked signal conflicts with the current run state and should be resolved or adopted.",
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
						"signal_kind":     stringValue(signal["kind"]),
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
		"run_id":                        runID,
		"canonical_source":              "run:" + runID,
		"canonical_authority_level":     runLevel,
		"run_claims":                    runClaims,
		"signal_sources":                signalSources,
		"adoption_gap_count":            len(adoptionGaps),
		"contract_gap_count":            len(contractGaps),
		"higher_authority_conflict_count": len(conflicts),
		"superseded_signal_count":       len(supersededSignals),
		"desired_state_gap_count":       len(desiredGapRecords),
		"live_state_gap_count":          len(liveStateGaps),
		"desired_state_satisfied":       desiredSatisfied,
		"adoption_gaps":                 adoptionGaps,
		"contract_gaps":                 contractGaps,
		"higher_authority_conflicts":    conflicts,
		"superseded_signals":            supersededSignals,
		"desired_state_gaps":            desiredGapRecords,
		"live_state_gaps":               liveStateGaps,
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
		for _, dependencyID := range dependencyTargets(run) {
			if index := findRunIndex(registry, dependencyID); index >= 0 {
				dependencyRun, _ := toMap(toList(registry["runs"])[index])
				if dependencyMissionID := strings.TrimSpace(stringValue(dependencyRun["mission_id"])); dependencyMissionID != "" {
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
	runs := []map[string]any{}
	for _, item := range toList(registry["runs"]) {
		run, _ := toMap(item)
		if stringValue(run["mission_id"]) == missionID {
			runs = append(runs, run)
		}
	}
	if len(runs) == 0 {
		return nil, nil
	}
	sort.Slice(runs, func(i, j int) bool {
		left := runSortTimestamp(runs[i])
		right := runSortTimestamp(runs[j])
		if left != right {
			return left > right
		}
		return stringValue(runs[i]["run_id"]) > stringValue(runs[j]["run_id"])
	})
	latestRun := runs[0]
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
	contract := map[string]any{}
	if contractRef != "" {
		contractRoot, _ := toMap(loadData(contractRef, map[string]any{}))
		contract, _ = toMap(contractRoot["contract"])
	}
	mission, _ := toMap(mergeDicts(template, map[string]any{
		"mission_id":            missionID,
		"issue_number":          intValue(firstNonBlank(stringValue(latestRun["issue_number"]), stringValue(contract["issue_number"]))),
		"title":                 firstNonBlank(stringValue(contract["title"]), stringValue(latestRun["objective"]), missionID),
		"objective":             firstNonBlank(stringValue(contract["objective"]), stringValue(latestRun["objective"]), missionID),
		"state":                 missionStateFromRuns(runs),
		"priority":              stringValue(latestRun["priority"]),
		"risk_level":            stringValue(latestRun["risk_level"]),
		"branch":                firstNonBlank(stringValue(contract["branch"]), branch),
		"worktree_path":         firstNonBlank(stringValue(contract["worktree_path"]), worktreePath),
		"contract_ref":          contractRef,
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
		"latest_summary":        latestSummary,
		"next_expected_action":  nextAction,
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

func defaultMeetingRoles() []any {
	return []any{"fact-verifier", "architecture-challenger", "completeness-checker"}
}

func defaultMeetingRoleSpecs() []any {
	timestamp := nowISO()
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

func materializeMeeting(meetingID string, run map[string]any, question, requestedBy string) (map[string]any, error) {
	templateRoot, _ := toMap(loadData(meetingTemplatePath(), map[string]any{}))
	template, _ := toMap(templateRoot["meeting"])
	timestamp := nowISO()
	meeting, _ := toMap(mergeDicts(template, map[string]any{
		"meeting_id":          meetingID,
		"mission_id":          stringValue(run["mission_id"]),
		"run_id":              stringValue(run["run_id"]),
		"question":            question,
		"scope":               firstNonBlank(stringValue(run["objective"]), stringValue(run["mission_id"])),
		"leader_id":           stringValue(run["leader_id"]),
		"requested_by":        requestedBy,
		"state":               "active",
		"authority_hierarchy": defaultAuthorityHierarchy(),
		"locked_decisions":    []any{},
		"roles":               defaultMeetingRoles(),
		"role_specs":          defaultMeetingRoleSpecs(),
		"packet": map[string]any{
			"expected_output": "Produce accepted findings, rejected findings, unresolved questions, and action items.",
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
	nextAction := "Resume execution with the meeting outcome applied."
	if len(actions) > 0 {
		nextAction = actions[0]
	}
	return map[string]any{
		"summary":     summary,
		"next_action": nextAction,
	}
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
	if _, ok := toMap(meeting["orchestration"]); !ok {
		meeting["orchestration"] = map[string]any{}
	}
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

func meetingSynthesisAudit(meeting map[string]any) map[string]any {
	coverage, _ := toMap(meetingCoverageAudit(meeting))
	findings, _ := toMap(meeting["findings"])
	accepted := toList(findings["accepted"])
	rejected := toList(findings["rejected"])
	unresolved := toList(findings["unresolved"])
	actions := toList(meeting["action_items"])

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
	coverageNotes := []any{}
	for _, item := range toList(coverage["coverage_notes"]) {
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
			"total_positions":           len(toList(meeting["positions"])),
			"phase_history_count":       len(toList(meeting["phase_history"])),
			"participant_count":         coverage["participant_count"],
			"linked_cross_examinations": coverage["linked_cross_examination_count"],
		},
		"blocking_gaps":   stringListToAny(blockingGaps),
		"coverage_notes":  coverageNotes,
		"evidence_gaps":   stringListToAny(evidenceGaps),
		"authority_audit": authorityAudit,
		"generated_at":    nowISO(),
	}
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
			groupKey := dimension + "::" + normalizedClaim
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
		latestClosedMeeting := latestClosedMeetingForRun(runID)
		if latestClosedMeeting != nil && meetingHasOutput(latestClosedMeeting) {
			outcome := meetingOutcome(latestClosedMeeting)
			if stringValue(outcome["summary"]) != "" && strings.TrimSpace(stringValue(run["latest_summary"])) != stringValue(outcome["summary"]) {
				issues = append(issues, issueRecord("warning", "run-missing-meeting-summary", runID, "Run has a closed meeting with accepted output that is not reflected in latest_summary.", "Adopt the meeting outcome into the run state."))
			}
			if stringValue(outcome["next_action"]) != "" && strings.TrimSpace(stringValue(run["next_expected_action"])) != stringValue(outcome["next_action"]) {
				issues = append(issues, issueRecord("warning", "run-missing-meeting-next-action", runID, "Run has a closed meeting with action items that are not reflected in next_expected_action.", "Adopt the meeting outcome into the run state."))
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
		if len(toList(mission["depends_on"])) > 0 && stringValue(mission["state"]) != "blocked" && stringValue(mission["state"]) != "claimed" && stringValue(mission["state"]) != "wip" && stringValue(mission["state"]) != "review" {
			issues = append(issues, issueRecord("info", "mission-dependency-state-unusual", missionID, "Mission has dependencies but is not in an active dependency-aware state.", "Verify the mission control-plane state is still correct."))
		}
		if stringValue(mission["state"]) == "done" && len(recordedActive) > 0 {
			issues = append(issues, issueRecord("warning", "mission-done-with-active-runs", missionID, "Mission is marked done but still has active runs.", "Close or correct the active runs before trusting the mission state."))
		}
		if desired, ok := toMap(mission["desired_state"]); !ok || stringValue(desired["goal_state"]) == "" {
			issues = append(issues, issueRecord("warning", "mission-missing-desired-state", missionID, "Mission is missing desired_state.goal_state.", "Record mission desired state explicitly for control-plane reconciliation."))
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
			for _, leftClaim := range ownershipClaims(leftRun) {
				for _, rightClaim := range ownershipClaims(rightRun) {
					if ownershipSelectorsOverlap(leftClaim, rightClaim) {
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
		authorityAudit, _ := toMap(meetingAuthorityAudit(meeting))
		orchestration, _ := toMap(meeting["orchestration"])
		synthesis, _ := toMap(orchestration["synthesis"])
		if phase := stringValue(meeting["phase"]); phase == "cross_examination" || phase == "synthesis" || phase == "output" {
			if len(toList(coverageAudit["blocking_gaps"])) > 0 {
				issues = append(issues, issueRecord("warning", "meeting-debate-blocked", meetingID, fmt.Sprintf("Meeting still has debate blocking gaps: %s", strings.Join(anyListToStrings(toList(coverageAudit["blocking_gaps"])), "; ")), "Fill the debate gaps or lower the meeting phase."))
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
				"next_expected_action": stringValue(mission["next_expected_action"]),
			})
		}
		return printJSON(map[string]any{
			"mission_count": len(graph),
			"missions":      graph,
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
			return errors.New("runs transfer-ownership requires SOURCE_RUN_ID TARGET_RUN_ID --owns=SELECTOR [--reason=TEXT]")
		}
		sourceRunID := args[1]
		targetRunID := args[2]
		selectors := []string{}
		reason := "Leader reassigned ownership selector."
		for _, arg := range args[3:] {
			switch {
			case strings.HasPrefix(arg, "--owns="):
				selectors = append(selectors, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--owns=")))...)
			case strings.HasPrefix(arg, "--reason="):
				reason = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--reason=")), reason)
			}
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
			sourceSelectors := ownershipClaims(sourceRun)
			targetSelectors := ownershipClaims(targetRun)
			sourceSet := map[string]bool{}
			for _, selector := range sourceSelectors {
				sourceSet[selector] = true
			}
			missing := []string{}
			for _, selector := range selectors {
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
				if !transferSet[selector] {
					remainingSource = append(remainingSource, selector)
				}
			}
			targetSelectors = dedupeStrings(append(targetSelectors, selectors...))
			setOwnershipSelectors(sourceRun, remainingSource)
			setOwnershipSelectors(targetRun, targetSelectors)
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
				"mode":            "move",
				"authority_level": "L2",
				"transferred_at":  nowISO(),
			}
			sourceRun["ownership_transfers"] = append(ownershipTransferEvents(sourceRun), cloneValue(event))
			targetRun["ownership_transfers"] = append(ownershipTransferEvents(targetRun), cloneValue(event))
			sourceRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership selectors transferred away.")
			sourceRun["next_expected_action"] = "Continue within the remaining bounded ownership scope."
			targetRun["latest_summary"] = firstNonBlank(strings.TrimSpace(reason), "Ownership selectors transferred in.")
			targetRun["next_expected_action"] = "Proceed with the newly transferred bounded ownership scope."
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
			for _, leftClaim := range ownershipClaims(leftRun) {
				for _, rightClaim := range ownershipClaims(rightRun) {
					if ownershipSelectorsOverlap(leftClaim, rightClaim) {
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
		switch stringValue(signal["kind"]) {
		case "requires_meeting", "interface_change", "dependency_discovered":
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
	nextAction := "Continue active execution and monitor signals."
	switch {
	case len(orphanSignals) > 0 || len(orphanMeetings) > 0:
		nextAction = "Reconcile orphaned runtime objects before proceeding."
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
	case activeRCRecord != nil:
		nextAction = "Review the active RC and decide whether to promote, fix, or close it."
	case len(waitingRuns) > 0:
		nextAction = "Review waiting runs and dependency state."
	}
	return printJSON(map[string]any{
		"profile_id":                  stringValue(getKey(profile, "profile.profile_id", "")),
		"active_mission_count":        len(activeMissions),
		"blocked_mission_count":       len(blockedMissions),
		"review_mission_count":        len(reviewMissions),
		"dependency_mission_count":    len(dependencyBlockedMissions),
		"mission_graph_edge_count":    missionGraphEdges,
		"active_run_count":            len(activeRuns),
		"active_teammate_count":       len(activeTeammates),
		"blocked_run_count":           len(blockedRuns),
		"dependency_blocked_count":    len(dependencyBlockedRuns),
		"review_blocked_count":        len(reviewBlockedRuns),
		"checks_failed_count":         len(checksFailedRuns),
		"merge_ready_count":           len(mergeReadyRuns),
		"ownership_overlap_count":     len(overlappingClaims),
		"active_signal_count":         len(activeSignals),
		"coordination_signal_count":   len(coordinationSignals),
		"active_meeting_count":        len(meetings),
		"active_rc_count":             boolToInt(activeRCRecord != nil),
		"open_meeting_action_count":   len(meetingActions),
		"orphan_signal_count":         len(orphanSignals),
		"orphan_meeting_count":        len(orphanMeetings),
		"missions":                    activeMissions,
		"dependency_blocked_missions": dependencyBlockedMissions,
		"runs":                        activeRuns,
		"dependency_blocked_runs":     dependencyBlockedRuns,
		"signals":                     activeSignals,
		"meetings":                    meetings,
		"overlapping_claims":          overlappingClaims,
		"active_rc":                   mapOrEmpty(activeRCRecord),
		"meeting_actions":             meetingActions,
		"orphan_signals":              orphanSignals,
		"orphan_meetings":             orphanMeetings,
		"next_action":                 nextAction,
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
			outcome := meetingOutcome(meeting)
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if stringValue(outcome["summary"]) != "" {
				run["latest_summary"] = outcome["summary"]
			}
			if stringValue(outcome["next_action"]) != "" {
				run["next_expected_action"] = outcome["next_action"]
			}
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
			return errors.New("drive meeting requires RUN_ID QUESTION")
		}
		runID := args[1]
		question := strings.TrimSpace(strings.Join(args[2:], " "))
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
			meetingID := fmt.Sprintf("meeting-%s-%s", nowTime().UTC().Format("20060102150405"), stringValue(run["mission_id"]))
			meeting, err := materializeMeeting(meetingID, run, question, stringValue(run["owner_id"]))
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
			outcome := meetingOutcome(meeting)
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
			runs := toList(registry["runs"])
			run, _ := toMap(runs[index])
			if err := validateRunTransition(run, "ready"); err != nil {
				return err
			}
			run["state"] = "ready"
			run["status_reason"] = note
			if stringValue(outcome["summary"]) != "" {
				run["latest_summary"] = outcome["summary"]
			}
			run["next_expected_action"] = outcome["next_action"]
			run["updated_at"] = nowISO()
			runs[index] = run
			registry["runs"] = runs
			if err := writeRegistry(registry); err != nil {
				return err
			}
			return printJSON(map[string]any{"meeting": meeting, "run": run})
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
		authorityAudit, _ := toMap(meetingAuthorityAudit(meeting))
		return printJSON(map[string]any{
			"meeting_id":                          stringValue(meeting["meeting_id"]),
			"state":                               stringValue(meeting["state"]),
			"question":                            stringValue(meeting["question"]),
			"outcome_summary":                     outcome["summary"],
			"next_action":                         outcome["next_action"],
			"phase_history_count":                 len(toList(meeting["phase_history"])),
			"debate_round":                        intValue(orchestration["debate_round"]),
			"synthesis_status":                    stringValue(synthesis["status"]),
			"blocking_gap_count":                  len(toList(synthesis["blocking_gaps"])),
			"coverage_note_count":                 len(toList(synthesis["coverage_notes"])),
			"authority_conflict_count":            intValue(authorityAudit["conflict_count"]),
			"unresolved_authority_conflict_count": intValue(authorityAudit["unresolved_conflict_count"]),
			"role_spec_count":                     len(toList(meeting["role_specs"])),
			"locked_decisions":                    len(toList(meeting["locked_decisions"])),
			"accepted_count":                      len(toList(findings["accepted"])),
			"rejected_count":                      len(toList(findings["rejected"])),
			"unresolved_count":                    len(toList(findings["unresolved"])),
			"open_action_count":                   openActionCount,
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
	default:
		return cmdMeetingsRemainder(args)
	}
}

func cmdMeetingsRemainder(args []string) error {
	subcommand := args[0]
	switch subcommand {
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
			oldPhase := firstNonBlank(stringValue(meeting["phase"]), "setup")
			packet, _ := toMap(meeting["packet"])
			roleSpecs := toList(meeting["role_specs"])
			participants := toList(meeting["participants"])
			coverageAudit, _ := toMap(meetingCoverageAudit(meeting))
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
			switch oldPhase {
			case "", "setup":
				if len(gaps) == 0 {
					newPhase = "independent_review"
					nextRequirement = "Collect independent-review positions from bounded participants."
				} else {
					nextRequirement = "Add packet/participant prerequisites before independent review."
				}
			case "independent_review":
				for _, item := range toList(coverageAudit["blocking_gaps"]) {
					text := strings.TrimSpace(stringValue(item))
					if text != "" {
						gaps = append(gaps, text)
					}
				}
				if len(independentReviewRuns) >= 2 {
					newPhase = "cross_examination"
					nextRequirement = "Collect cross-examination positions that challenge or refine earlier claims."
				} else {
					nextRequirement = "Need at least two independent-review positions before cross-examination."
				}
			case "cross_examination":
				for _, item := range toList(coverageAudit["blocking_gaps"]) {
					text := strings.TrimSpace(stringValue(item))
					if text != "" {
						gaps = append(gaps, text)
					}
				}
				if len(crossExamRuns) >= 1 && intValue(coverageAudit["linked_cross_examination_count"]) > 0 {
					newPhase = "synthesis"
					nextRequirement = "Convert debated positions into accepted, rejected, unresolved findings, and action items."
				} else {
					nextRequirement = "Need at least one linked cross-examination position before synthesis."
				}
			case "synthesis":
				synthesis := meetingSynthesisAudit(meeting)
				if stringValue(synthesis["status"]) == "ready_for_output" {
					newPhase = "output"
					nextRequirement = "Adopt or close the meeting outcome into run state."
				} else {
					nextRequirement = "Supplement synthesis gaps before output."
				}
			case "output":
				nextRequirement = "Close the meeting or adopt its outcome."
			default:
				nextRequirement = "Unknown phase; repair the meeting phase explicitly."
			}
			orchestration, _ := toMap(meeting["orchestration"])
			if orchestration == nil {
				orchestration = map[string]any{}
			}
			orchestration["last_gap_report"] = stringListToAny(gaps)
			orchestration["next_requirement"] = nextRequirement
			orchestration["last_orchestrated_at"] = nowISO()
			orchestration["coverage"] = coverageAudit
			orchestration["synthesis"] = meetingSynthesisAudit(meeting)
			orchestration["debate_round"] = maxInt(intValue(orchestration["debate_round"]), intValue(coverageAudit["linked_cross_examination_count"]))
			meeting["orchestration"] = orchestration
			if newPhase != oldPhase {
				setMeetingPhase(meeting, newPhase, fmt.Sprintf("orchestrate advanced from %s", oldPhase))
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
		for _, arg := range args[2:] {
			switch {
			case strings.HasPrefix(arg, "--reason="):
				reason = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--reason=")), reason)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), authorityLevel)
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
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
			locked := append(toList(meeting["locked_decisions"]), map[string]any{
				"decision_id":     "decision-" + nowTime().UTC().Format("20060102150405"),
				"summary":         summary,
				"reason":          reason,
				"authority_level": authorityLevel,
				"evidence_refs":   stringListToAny(evidenceRefs),
				"locked_at":       nowISO(),
			})
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
		respondsToRunID := ""
		respondsToRole := ""
		respondsToPositionID := ""
		for _, arg := range args[5:] {
			switch {
			case strings.HasPrefix(arg, "--evidence="):
				evidenceRefs = append(evidenceRefs, parsePipeList(strings.TrimSpace(strings.TrimPrefix(arg, "--evidence=")))...)
			case strings.HasPrefix(arg, "--authority="):
				authorityLevel = firstNonBlank(strings.TrimSpace(strings.TrimPrefix(arg, "--authority=")), "L2")
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
				"finding_id":      "finding-" + nowTime().UTC().Format("20060102150405"),
				"summary":         summary,
				"reason":          reason,
				"authority_level": authorityLevel,
				"evidence_refs":   stringListToAny(evidenceRefs),
				"source_run_id":   sourceRunID,
				"source_role":     sourceRole,
				"claim_ref":       claimRef,
				"dimension":       dimension,
				"recorded_at":     nowISO(),
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
			return errors.New("teammates escalate requires RUN_ID QUESTION [--materialize]")
		}
		runID := args[1]
		materialize := false
		questionTokens := []string{}
		for _, arg := range args[2:] {
			if arg == "--materialize" {
				materialize = true
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
				meetingID := fmt.Sprintf("meeting-%s-%s", nowTime().UTC().Format("20060102150405"), stringValue(run["mission_id"]))
				meeting, err := materializeMeeting(meetingID, run, question, stringValue(run["owner_id"]))
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
				if stringValue(signal["emitted_at"]) == "" {
					signal["emitted_at"] = firstNonBlank(stringValue(signal["fresh_until"]), nowISO())
					normalizedSignals = append(normalizedSignals, stringValue(signal["signal_id"]))
					signals[index] = signal
				}
			}
			log["signals"] = signals
			needsMissionSync := len(toList(missionRegistry["missions"])) == 0 && len(runs) > 0
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
