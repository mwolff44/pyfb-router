package models

// SetCallerID sets the callerId and corresponding headers regarding the GW rules
func SetCallerID(callerID string, pai bool, ppi bool, pid bool, addPlusInCaller bool, socket string) (string, string) {
	var extraHeaders string
	// Set callerID regarding GW rules

	// Add + in front of CallerID for outbound call
	if addPlusInCaller == true {
		callerID = "+" + callerID
	}
	// If anonymous
	// If PAI, set it
	// RFC 3325
	if pai == true {
		extraHeaders = extraHeaders + "P-Asserted-Identity: <sip:" + callerID + "@" + socket + ">\r\n"
	}

	// If PPI, set it
	// RFC 3325
	if ppi == true {
		extraHeaders = extraHeaders + "P-Preferred-Identity: <sip:" + callerID + "@" + socket + ">\r\n"
	}

	// If RPID, set it
	// https://tools.ietf.org/html/draft-ietf-sip-privacy-00
	if pid == true {
		extraHeaders = extraHeaders + "Remote-Party-ID: <sip:" + callerID + "@" + socket + ">\r\n"
	}

	return callerID, extraHeaders
}

// SetCallee sets the CalleeId regarding the GW rules - it concatenates gateway prefix with
// ratecard prefix with calleeID and at the end the suffix
func SetCallee(calleeID string, gwprefix string, suffix string, ratecardprefix string) string {
	// Set Callee regarding GW rules

	calleeID = gwprefix + ratecardprefix + calleeID + suffix

	return calleeID
}

// ConstructFromHeader constructs the SIP FROM header
func ConstructFromHeader(calleridInFrom bool, newCallerID string, socket string, username string) HeaderDetail {
	var fromHeader HeaderDetail
	if calleridInFrom == true {
		fromHeader = HeaderDetail{Display: newCallerID, URI: "sip:" + newCallerID + "@" + socket}
	} else {
		fromHeader = HeaderDetail{Display: "", URI: "sip:" + username + "@" + socket}
	}

	return fromHeader
}
