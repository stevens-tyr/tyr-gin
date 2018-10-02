package tyrgin

import "sync"

// Aggregate execute all statusEndpoint StatusCheck() functions asynchronously and return the
// overall status by returning the highest severity item in the following order:
// CRIT, WARN, OK
func Aggregate(statusEndpoints []StatusEndpoint, typeFilter string) string {

	// Type filter needs to be of a certain format.
	if len(typeFilter) > 0 {
		if typeFilter != "internal" && typeFilter != "external" {
			sl := StatusList{
				StatusList: []Status{
					{
						Description: "Invalid type",
						Result:      CRITICAL,
						Details:     "Unknown check type given for aggregate check",
					},
				},
			}

			return SerializeStatusList(sl)
		}
	}

	// If typeFilter is none loop through and check if typeFilters are
	// internal or external and if so add it to list of StatusEndpoints.
	s := statusEndpoints
	if typeFilter != "" {

		s = []StatusEndpoint{}
		for _, statusEndpoint := range statusEndpoints {

			if typeFilter == "internal" {
				if statusEndpoint.Type == "internal" {
					s = append(s, statusEndpoint)
				}
			} else if typeFilter == "external" {
				if statusEndpoint.Type != "internal" {
					s = append(s, statusEndpoint)
				}
			}

		}

	}

	responses := make(chan StatusList)

	// Concurrent problems make Wait Group the size of the slice.
	var wg sync.WaitGroup
	wg.Add(len(s))

	// Check the status of each.
	for _, statusEndpoint := range s {
		go func(statusEndpoint StatusEndpoint) {
			responses <- statusEndpoint.StatusCheck.CheckStatus(statusEndpoint.Name)
		}(statusEndpoint)
	}

	// Make a list for each type.
	var crits []StatusList
	var warns []StatusList
	var oks []StatusList

	// Collect each of the status checks and store them into appropiate list.
	go func() {
		for r := range responses {
			switch r.StatusList[0].Result {
			case CRITICAL:
				crits = append(crits, r)
			case WARNING:
				warns = append(warns, r)
			case OK:
				oks = append(oks, r)
			default:
				panic("Invalid AlertLevel")
			}
			wg.Done()
		}
	}()

	wg.Wait()
	close(responses)

	// Default list if no critical or warnings.
	sl := StatusList{
		StatusList: []Status{
			{
				Description: "Aggregate Check",
				Result:      OK,
				Details:     "",
			},
		},
	}

	// Critical is higher precedence than warning.
	if len(crits) > 0 {
		sl = crits[0]
	} else if len(warns) > 0 {
		sl = warns[0]
	}

	return SerializeStatusList(sl)
}
