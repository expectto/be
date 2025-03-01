package be_testified

// We need to de-gomegafy our matchers. Here's the list of matchers that depends on gomega:

/*
•	be_math Package (matchers_math.go):
	• GreaterThan, GreaterThanEqual, LessThan, LessThanEqual, Approx:
	• InRange:
•	be_time Package (matchers_time.go):
	• LaterThan, LaterThanEqual, EarlierThan, EarlierThanEqual, Eq, Approx:
•	be_string Package (matchers_string.go):
	• EmptyString, NonEmptyString, Float, Titled, LowerCaseOnly, UpperCaseOnly, ContainingSubstring, ContainingOnlyCharacters, ContainingCharacters, MatchWildcard, ValidEmail:
•	be_url Package (matchers_url.go):
	• URL, HavingHost, HavingHostname, HavingScheme, HavingPort, HavingPath, HavingRawQuery, HavingSearchParam, HavingMultipleSearchParam, HavingUsername, HavingUserinfo, HavingPassword, WithHttps, WithHttp, NotHavingScheme, NotHavingPort:
•	be_http Package (matchers_http.go):
	• Request, HavingMethod, HavingURL, etc.:
•	Internal PSI Matchers (under internal/psi_matchers/):
	• all_matcher.go (AllMatcher), any_matcher.go (AnyMatcher):
	• eq_matcher.go (EqMatcher):
	• never_matcher.go (NeverMatcher):
	• not_matcher.go (NotMatcher):
	• have_length_matcher.go (HaveLengthMatcher):
	• assignable_to_matcher.go (AssignableToMatcher):
	• implements_matcher.go (ImplementsMatcher):
	• jwt_token_matcher.go (JwtTokenMatcher):
	• kind_matcher.go (KindMatcher):
	• string_template_matcher.go (StringTemplateMatcher):
	• url_field_matcher.go (UrlFieldMatcher):
*/
