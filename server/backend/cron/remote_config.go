package cron

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Raw JSON objects preserve unrelated parameters, groups, conditions and metadata.
type configObject map[string]json.RawMessage

func setStatusParameters(template configObject, state tumonlineStatus) (bool, error) {
	parameters := configObject{}
	if raw, ok := template["parameters"]; ok {
		if err := json.Unmarshal(raw, &parameters); err != nil || parameters == nil {
			return false, fmt.Errorf("invalid Remote Config parameters")
		}
	}
	groups := map[string]configObject{}
	if raw, ok := template["parameterGroups"]; ok {
		if err := json.Unmarshal(raw, &groups); err != nil {
			return false, err
		}
	}
	changed := false
	for key, value := range map[string]struct{ value, valueType string }{
		"isTUMOnlineMaintenanceMode": {strconv.FormatBool(state.Maintenance), "BOOLEAN"},
		"isTUMOnlineDegraded":        {strconv.FormatBool(state.Degraded), "BOOLEAN"},
		"tumOnlineStatusMessage_de":  {state.MessageDE, "STRING"},
		"tumOnlineStatusMessage_en":  {state.MessageEN, "STRING"},
	} {
		target := parameters
		groupName := ""
		for name, group := range groups {
			var groupParameters configObject
			if err := json.Unmarshal(group["parameters"], &groupParameters); err != nil {
				return false, err
			}
			if _, ok := groupParameters[key]; ok {
				if _, duplicate := target[key]; duplicate {
					return false, fmt.Errorf("duplicate Remote Config key %s", key)
				}
				target, groupName = groupParameters, name
			}
		}
		parameter := configObject{}
		if raw, ok := target[key]; ok {
			if err := json.Unmarshal(raw, &parameter); err != nil || parameter == nil {
				return false, fmt.Errorf("invalid Remote Config parameter %s", key)
			}
		}
		// Conditional values would mask the synchronized defaults for some users.
		var conditional configObject
		if raw, ok := parameter["conditionalValues"]; ok {
			if err := json.Unmarshal(raw, &conditional); err != nil {
				return false, err
			}
			if len(conditional) != 0 {
				return false, fmt.Errorf("remove conditional values from managed Remote Config key %s before enabling sync", key)
			}
		}
		var previous map[string]any
		if raw, ok := parameter["defaultValue"]; ok {
			if err := json.Unmarshal(raw, &previous); err != nil {
				return false, err
			}
		}
		var previousType string
		if raw, ok := parameter["valueType"]; ok {
			if err := json.Unmarshal(raw, &previousType); err != nil {
				return false, err
			}
		}
		if previous["value"] == value.value && previous["useInAppDefault"] != true && previousType == value.valueType {
			continue
		}
		changed = true
		parameter["defaultValue"], _ = json.Marshal(map[string]string{"value": value.value})
		parameter["valueType"], _ = json.Marshal(value.valueType)
		target[key], _ = json.Marshal(parameter)
		if groupName != "" {
			groups[groupName]["parameters"], _ = json.Marshal(target)
		}
	}
	if changed {
		template["parameters"], _ = json.Marshal(parameters)
		if len(groups) > 0 {
			template["parameterGroups"], _ = json.Marshal(groups)
		}
	}
	return changed, nil
}

func publishTUMOnlineStatus(ctx context.Context, client *http.Client, endpoint string, state tumonlineStatus) error {
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return err
		}
		// Firebase requires compression to reliably return the ETag.
		req.Header.Set("Accept-Encoding", "gzip")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		template, err := readConfigTemplate(resp)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		changed, err := setStatusParameters(template, state)
		if err != nil || !changed {
			return err
		}
		etag := resp.Header.Get("ETag")
		if etag == "" || etag == "*" {
			return fmt.Errorf("Firebase returned no usable ETag")
		}
		delete(template, "version") // Read-only server metadata.
		body, err := json.Marshal(template)
		if err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag)
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusPreconditionFailed {
			continue // Re-read and merge against the latest template.
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Firebase publish returned HTTP %d", resp.StatusCode)
		}
		return nil
	}
	return fmt.Errorf("Firebase template changed during all three publish attempts")
}

func readConfigTemplate(resp *http.Response) (configObject, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Firebase fetch returned HTTP %d", resp.StatusCode)
	}
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer func() { _ = gz.Close() }()
		reader = gz
	}
	body, err := io.ReadAll(io.LimitReader(reader, (8<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(body) > 8<<20 {
		return nil, fmt.Errorf("Firebase template exceeds size limit")
	}
	var template configObject
	if err := json.Unmarshal(body, &template); err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("invalid Firebase template")
	}
	return template, nil
}
