package skycrypttypes

import "encoding/json"

// JSONUnmarshal decodes the values inside the custom UnmarshalJSON methods in this package. It defaults to
// encoding/json. A program that decodes profiles with a faster encoding/json-compatible decoder, such as
// sonic.ConfigStd, can set it to that decoder's Unmarshal once at startup, before any decoding: otherwise every
// custom field drops back to encoding/json's reflection path, which was most of the cost of decoding a profile.
var JSONUnmarshal = json.Unmarshal
