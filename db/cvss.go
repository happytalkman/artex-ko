package db

import (
	"fmt"
	"math"
	"strings"
)

// CalcCVSS31Score parses a CVSS v3.1 base vector string in the form
// AV:_/AC:_/PR:_/UI:_/S:_/C:_/I:_/A:_ and returns the numeric base score
// (0.0–10.0) according to the CVSS v3.1 specification.
// Returns an error when the vector is malformed or contains unknown metric values.
func CalcCVSS31Score(vector string) (float64, error) {
	// Strip optional "CVSS:3.1/" prefix
	vector = strings.TrimPrefix(vector, "CVSS:3.1/")
	vector = strings.TrimPrefix(vector, "CVSS:3.0/")

	parts := strings.Split(vector, "/")
	metrics := make(map[string]string, len(parts))
	for _, p := range parts {
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			return 0, fmt.Errorf("cvss: malformed metric %q", p)
		}
		metrics[strings.ToUpper(kv[0])] = strings.ToUpper(kv[1])
	}

	get := func(key string) (string, error) {
		v, ok := metrics[key]
		if !ok {
			return "", fmt.Errorf("cvss: missing metric %s", key)
		}
		return v, nil
	}

	avS, err := get("AV")
	if err != nil {
		return 0, err
	}
	acS, err := get("AC")
	if err != nil {
		return 0, err
	}
	prS, err := get("PR")
	if err != nil {
		return 0, err
	}
	uiS, err := get("UI")
	if err != nil {
		return 0, err
	}
	sS, err := get("S")
	if err != nil {
		return 0, err
	}
	cS, err := get("C")
	if err != nil {
		return 0, err
	}
	iS, err := get("I")
	if err != nil {
		return 0, err
	}
	aS, err := get("A")
	if err != nil {
		return 0, err
	}

	// Attack Vector
	avMap := map[string]float64{"N": 0.85, "A": 0.62, "L": 0.55, "P": 0.2}
	av, ok := avMap[avS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown AV value %q", avS)
	}

	// Attack Complexity
	acMap := map[string]float64{"L": 0.77, "H": 0.44}
	ac, ok := acMap[acS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown AC value %q", acS)
	}

	// Scope (needed for PR lookup)
	scopeChanged := sS == "C"
	if sS != "U" && sS != "C" {
		return 0, fmt.Errorf("cvss: unknown S value %q", sS)
	}

	// Privileges Required (depends on Scope)
	var prMap map[string]float64
	if scopeChanged {
		prMap = map[string]float64{"N": 0.85, "L": 0.68, "H": 0.50}
	} else {
		prMap = map[string]float64{"N": 0.85, "L": 0.62, "H": 0.27}
	}
	pr, ok := prMap[prS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown PR value %q", prS)
	}

	// User Interaction
	uiMap := map[string]float64{"N": 0.85, "R": 0.62}
	ui, ok := uiMap[uiS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown UI value %q", uiS)
	}

	// Impact metrics
	impactMap := map[string]float64{"N": 0.0, "L": 0.22, "H": 0.56}
	cImp, ok := impactMap[cS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown C value %q", cS)
	}
	iImp, ok := impactMap[iS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown I value %q", iS)
	}
	aImp, ok := impactMap[aS]
	if !ok {
		return 0, fmt.Errorf("cvss: unknown A value %q", aS)
	}

	// ISCBase = 1 - (1-C)*(1-I)*(1-A)
	iscBase := 1.0 - (1.0-cImp)*(1.0-iImp)*(1.0-aImp)

	// ISS (Impact Sub-Score)
	var iss float64
	if scopeChanged {
		iss = 7.52*(iscBase-0.029) - 3.25*math.Pow(iscBase-0.02, 15)
	} else {
		iss = 6.42 * iscBase
	}

	if iss <= 0 {
		return 0.0, nil
	}

	// Exploitability Sub-Score
	ess := 8.22 * av * ac * pr * ui

	// Base Score
	var base float64
	if scopeChanged {
		base = math.Min(1.08*(iss+ess), 10.0)
	} else {
		base = math.Min(iss+ess, 10.0)
	}

	// Round up to one decimal place
	return roundUp1(base), nil
}

// roundUp1 rounds a float64 up to one decimal place (CVSS 3.1 rounding rule).
func roundUp1(x float64) float64 {
	// Multiply by 10, ceil, divide by 10.
	rounded := math.Ceil(x*10) / 10
	// Guard against floating-point noise: snap values like 9.80000000001 to 9.8
	return math.Round(rounded*10) / 10
}
