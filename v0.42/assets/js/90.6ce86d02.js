(window.webpackJsonp=window.webpackJsonp||[]).push([[90],{698:function(e,t,o){"use strict";o.r(t);var n=o(1),a=Object(n.a)({},(function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[o("h1",{attrs:{id:"adr-027-deterministic-protobuf-serialization"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#adr-027-deterministic-protobuf-serialization"}},[e._v("#")]),e._v(" ADR 027: Deterministic Protobuf Serialization")]),e._v(" "),o("h2",{attrs:{id:"changelog"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#changelog"}},[e._v("#")]),e._v(" Changelog")]),e._v(" "),o("ul",[o("li",[e._v("2020-08-07: Initial Draft")]),e._v(" "),o("li",[e._v("2020-09-01: Further clarify rules")])]),e._v(" "),o("h2",{attrs:{id:"status"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#status"}},[e._v("#")]),e._v(" Status")]),e._v(" "),o("p",[e._v("Proposed")]),e._v(" "),o("h2",{attrs:{id:"abstract"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#abstract"}},[e._v("#")]),e._v(" Abstract")]),e._v(" "),o("p",[e._v("Fully deterministic structure serialization, which works across many languages and clients,\nis needed when signing messages. We need to be sure that whenever we serialize\na data structure, no matter in which supported language, the raw bytes\nwill stay the same.\n"),o("a",{attrs:{href:"https://developers.google.com/protocol-buffers/docs/proto3",target:"_blank",rel:"noopener noreferrer"}},[e._v("Protobuf"),o("OutboundLink")],1),e._v("\nserialization is not bijective (i.e. there exist a practically unlimited number of\nvalid binary representations for a given protobuf document)"),o("sup",[e._v("1")]),e._v(".")]),e._v(" "),o("p",[e._v("This document describes a deterministic serialization scheme for\na subset of protobuf documents, that covers this use case but can be reused in\nother cases as well.")]),e._v(" "),o("h3",{attrs:{id:"context"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#context"}},[e._v("#")]),e._v(" Context")]),e._v(" "),o("p",[e._v("For signature verification in Cosmos SDK, the signer and verifier need to agree on\nthe same serialization of a "),o("code",[e._v("SignDoc")]),e._v(" as defined in\n"),o("RouterLink",{attrs:{to:"/architecture/adr-020-protobuf-transaction-encoding.html"}},[e._v("ADR-020")]),e._v(" without transmitting the\nserialization.")],1),e._v(" "),o("p",[e._v("Currently, for block signatures we are using a workaround: we create a new "),o("a",{attrs:{href:"https://github.com/cosmos/cosmos-sdk/blob/9e85e81e0e8140067dd893421290c191529c148c/proto/cosmos/tx/v1beta1/tx.proto#L30",target:"_blank",rel:"noopener noreferrer"}},[e._v("TxRaw"),o("OutboundLink")],1),e._v("\ninstance (as defined in "),o("a",{attrs:{href:"https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-020-protobuf-transaction-encoding.md#transactions",target:"_blank",rel:"noopener noreferrer"}},[e._v("adr-020-protobuf-transaction-encoding"),o("OutboundLink")],1),e._v(")\nby converting all "),o("a",{attrs:{href:"https://github.com/cosmos/cosmos-sdk/blob/9e85e81e0e8140067dd893421290c191529c148c/proto/cosmos/tx/v1beta1/tx.proto#L13",target:"_blank",rel:"noopener noreferrer"}},[e._v("Tx"),o("OutboundLink")],1),e._v("\nfields to bytes on the client side. This adds an additional manual\nstep when sending and signing transactions.")]),e._v(" "),o("h3",{attrs:{id:"decision"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#decision"}},[e._v("#")]),e._v(" Decision")]),e._v(" "),o("p",[e._v("The following encoding scheme is to be used by other ADRs,\nand in particular for "),o("code",[e._v("SignDoc")]),e._v(" serialization.")]),e._v(" "),o("h2",{attrs:{id:"specification"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#specification"}},[e._v("#")]),e._v(" Specification")]),e._v(" "),o("h3",{attrs:{id:"scope"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#scope"}},[e._v("#")]),e._v(" Scope")]),e._v(" "),o("p",[e._v("This ADR defines a protobuf3 serializer. The output is a valid protobuf\nserialization, such that every protobuf parser can parse it.")]),e._v(" "),o("p",[e._v("No maps are supported in version 1 due to the complexity of defining a\ndeterministic serialization. This might change in future. Implementations must\nreject documents containing maps as invalid input.")]),e._v(" "),o("h3",{attrs:{id:"background-protobuf3-encoding"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#background-protobuf3-encoding"}},[e._v("#")]),e._v(" Background - Protobuf3 Encoding")]),e._v(" "),o("p",[e._v("Most numeric types in protobuf3 are encoded as\n"),o("a",{attrs:{href:"https://developers.google.com/protocol-buffers/docs/encoding#varints",target:"_blank",rel:"noopener noreferrer"}},[e._v("varints"),o("OutboundLink")],1),e._v(".\nVarints are at most 10 bytes, and since each varint byte has 7 bits of data,\nvarints are a representation of "),o("code",[e._v("uint70")]),e._v(" (70-bit unsigned integer). When\nencoding, numeric values are casted from their base type to "),o("code",[e._v("uint70")]),e._v(", and when\ndecoding, the parsed "),o("code",[e._v("uint70")]),e._v(" is casted to the appropriate numeric type.")]),e._v(" "),o("p",[e._v("The maximum valid value for a varint that complies with protobuf3 is\n"),o("code",[e._v("FF FF FF FF FF FF FF FF FF 7F")]),e._v(" (i.e. "),o("code",[e._v("2**70 -1")]),e._v("). If the field type is\n"),o("code",[e._v("{,u,s}int64")]),e._v(", the highest 6 bits of the 70 are dropped during decoding,\nintroducing 6 bits of malleability. If the field type is "),o("code",[e._v("{,u,s}int32")]),e._v(", the\nhighest 38 bits of the 70 are dropped during decoding, introducing 38 bits of\nmalleability.")]),e._v(" "),o("p",[e._v("Among other sources of non-determinism, this ADR eliminates the possibility of\nencoding malleability.")]),e._v(" "),o("h3",{attrs:{id:"serialization-rules"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#serialization-rules"}},[e._v("#")]),e._v(" Serialization rules")]),e._v(" "),o("p",[e._v("The serialization is based on the\n"),o("a",{attrs:{href:"https://developers.google.com/protocol-buffers/docs/encoding",target:"_blank",rel:"noopener noreferrer"}},[e._v("protobuf3 encoding"),o("OutboundLink")],1),e._v("\nwith the following additions:")]),e._v(" "),o("ol",[o("li",[e._v("Fields must be serialized only once in ascending order")]),e._v(" "),o("li",[e._v("Extra fields or any extra data must not be added")]),e._v(" "),o("li",[o("a",{attrs:{href:"https://developers.google.com/protocol-buffers/docs/proto3#default",target:"_blank",rel:"noopener noreferrer"}},[e._v("Default values"),o("OutboundLink")],1),e._v("\nmust be omitted")]),e._v(" "),o("li",[o("code",[e._v("repeated")]),e._v(" fields of scalar numeric types must use\n"),o("a",{attrs:{href:"https://developers.google.com/protocol-buffers/docs/encoding#packed",target:"_blank",rel:"noopener noreferrer"}},[e._v("packed encoding"),o("OutboundLink")],1)]),e._v(" "),o("li",[e._v("Varint encoding must not be longer than needed:\n"),o("ul",[o("li",[e._v("No trailing zero bytes (in little endian, i.e. no leading zeroes in big\nendian). Per rule 3 above, the default value of "),o("code",[e._v("0")]),e._v(" must be omitted, so\nthis rule does not apply in such cases.")]),e._v(" "),o("li",[e._v("The maximum value for a varint must be "),o("code",[e._v("FF FF FF FF FF FF FF FF FF 01")]),e._v(".\nIn other words, when decoded, the highest 6 bits of the 70-bit unsigned\ninteger must be "),o("code",[e._v("0")]),e._v(". (10-byte varints are 10 groups of 7 bits, i.e.\n70 bits, of which only the lowest 70-6=64 are useful.)")]),e._v(" "),o("li",[e._v("The maximum value for 32-bit values in varint encoding must be "),o("code",[e._v("FF FF FF FF 0F")]),e._v("\nwith one exception (below). In other words, when decoded, the highest 38\nbits of the 70-bit unsigned integer must be "),o("code",[e._v("0")]),e._v(".\n"),o("ul",[o("li",[e._v("The one exception to the above is "),o("em",[e._v("negative")]),e._v(" "),o("code",[e._v("int32")]),e._v(", which must be\nencoded using the full 10 bytes for sign extension"),o("sup",[e._v("2")]),e._v(".")])])]),e._v(" "),o("li",[e._v("The maximum value for Boolean values in varint encoding must be "),o("code",[e._v("01")]),e._v(" (i.e.\nit must be "),o("code",[e._v("0")]),e._v(" or "),o("code",[e._v("1")]),e._v("). Per rule 3 above, the default value of "),o("code",[e._v("0")]),e._v(" must\nbe omitted, so if a Boolean is included it must have a value of "),o("code",[e._v("1")]),e._v(".")])])])]),e._v(" "),o("p",[e._v("While rule number 1. and 2. should be pretty straight forward and describe the\ndefault behavior of all protobuf encoders the author is aware of, the 3rd rule\nis more interesting. After a protobuf3 deserialization you cannot differentiate\nbetween unset fields and fields set to the default value"),o("sup",[e._v("3")]),e._v(". At\nserialization level however, it is possible to set the fields with an empty\nvalue or omitting them entirely. This is a significant difference to e.g. JSON\nwhere a property can be empty ("),o("code",[e._v('""')]),e._v(", "),o("code",[e._v("0")]),e._v("), "),o("code",[e._v("null")]),e._v(" or undefined, leading to 3\ndifferent documents.")]),e._v(" "),o("p",[e._v("Omitting fields set to default values is valid because the parser must assign\nthe default value to fields missing in the serialization"),o("sup",[e._v("4")]),e._v(". For scalar\ntypes, omitting defaults is required by the spec"),o("sup",[e._v("5")]),e._v(". For "),o("code",[e._v("repeated")]),e._v("\nfields, not serializing them is the only way to express empty lists. Enums must\nhave a first element of numeric value 0, which is the default"),o("sup",[e._v("6")]),e._v(". And\nmessage fields default to unset"),o("sup",[e._v("7")]),e._v(".")]),e._v(" "),o("p",[e._v("Omitting defaults allows for some amount of forward compatibility: users of\nnewer versions of a protobuf schema produce the same serialization as users of\nolder versions as long as newly added fields are not used (i.e. set to their\ndefault value).")]),e._v(" "),o("h3",{attrs:{id:"implementation"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#implementation"}},[e._v("#")]),e._v(" Implementation")]),e._v(" "),o("p",[e._v("There are three main implementation strategies, ordered from the least to the\nmost custom development:")]),e._v(" "),o("ul",[o("li",[o("p",[o("strong",[e._v("Use a protobuf serializer that follows the above rules by default.")]),e._v(" E.g.\n"),o("a",{attrs:{href:"https://pkg.go.dev/github.com/gogo/protobuf/gogoproto",target:"_blank",rel:"noopener noreferrer"}},[e._v("gogoproto"),o("OutboundLink")],1),e._v(" is known to\nbe compliant by in most cases, but not when certain annotations such as\n"),o("code",[e._v("nullable = false")]),e._v(" are used. It might also be an option to configure an\nexisting serializer accordingly.")])]),e._v(" "),o("li",[o("p",[o("strong",[e._v("Normalize default values before encoding them.")]),e._v(" If your serializer follows\nrule 1. and 2. and allows you to explicitly unset fields for serialization,\nyou can normalize default values to unset. This can be done when working with\n"),o("a",{attrs:{href:"https://www.npmjs.com/package/protobufjs",target:"_blank",rel:"noopener noreferrer"}},[e._v("protobuf.js"),o("OutboundLink")],1),e._v(":")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"js",base64:"Y29uc3QgYnl0ZXMgPSBTaWduRG9jLmVuY29kZSh7CiAgYm9keUJ5dGVzOiBib2R5Lmxlbmd0aCAmZ3Q7IDAgPyBib2R5IDogbnVsbCwgLy8gbm9ybWFsaXplIGVtcHR5IGJ5dGVzIHRvIHVuc2V0CiAgYXV0aEluZm9CeXRlczogYXV0aEluZm8ubGVuZ3RoICZndDsgMCA/IGF1dGhJbmZvIDogbnVsbCwgLy8gbm9ybWFsaXplIGVtcHR5IGJ5dGVzIHRvIHVuc2V0CiAgY2hhaW5JZDogY2hhaW5JZCB8fCBudWxsLCAvLyBub3JtYWxpemUgJnF1b3Q7JnF1b3Q7IHRvIHVuc2V0CiAgYWNjb3VudE51bWJlcjogYWNjb3VudE51bWJlciB8fCBudWxsLCAvLyBub3JtYWxpemUgMCB0byB1bnNldAogIGFjY291bnRTZXF1ZW5jZTogYWNjb3VudFNlcXVlbmNlIHx8IG51bGwsIC8vIG5vcm1hbGl6ZSAwIHRvIHVuc2V0Cn0pLmZpbmlzaCgpOwo="}})],1),e._v(" "),o("li",[o("p",[o("strong",[e._v("Use a hand-written serializer for the types you need.")]),e._v(" If none of the above\nways works for you, you can write a serializer yourself. For SignDoc this\nwould look something like this in Go, building on existing protobuf utilities:")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"aWYgIXNpZ25Eb2MuYm9keV9ieXRlcy5lbXB0eSgpIHsKICAgIGJ1Zi5Xcml0ZVVWYXJJbnQ2NCgweEEpIC8vIHdpcmUgdHlwZSBhbmQgZmllbGQgbnVtYmVyIGZvciBib2R5X2J5dGVzCiAgICBidWYuV3JpdGVVVmFySW50NjQoc2lnbkRvYy5ib2R5X2J5dGVzLmxlbmd0aCgpKQogICAgYnVmLldyaXRlQnl0ZXMoc2lnbkRvYy5ib2R5X2J5dGVzKQp9CgppZiAhc2lnbkRvYy5hdXRoX2luZm8uZW1wdHkoKSB7CiAgICBidWYuV3JpdGVVVmFySW50NjQoMHgxMikgLy8gd2lyZSB0eXBlIGFuZCBmaWVsZCBudW1iZXIgZm9yIGF1dGhfaW5mbwogICAgYnVmLldyaXRlVVZhckludDY0KHNpZ25Eb2MuYXV0aF9pbmZvLmxlbmd0aCgpKQogICAgYnVmLldyaXRlQnl0ZXMoc2lnbkRvYy5hdXRoX2luZm8pCn0KCmlmICFzaWduRG9jLmNoYWluX2lkLmVtcHR5KCkgewogICAgYnVmLldyaXRlVVZhckludDY0KDB4MWEpIC8vIHdpcmUgdHlwZSBhbmQgZmllbGQgbnVtYmVyIGZvciBjaGFpbl9pZAogICAgYnVmLldyaXRlVVZhckludDY0KHNpZ25Eb2MuY2hhaW5faWQubGVuZ3RoKCkpCiAgICBidWYuV3JpdGVCeXRlcyhzaWduRG9jLmNoYWluX2lkKQp9CgppZiBzaWduRG9jLmFjY291bnRfbnVtYmVyICE9IDAgewogICAgYnVmLldyaXRlVVZhckludDY0KDB4MjApIC8vIHdpcmUgdHlwZSBhbmQgZmllbGQgbnVtYmVyIGZvciBhY2NvdW50X251bWJlcgogICAgYnVmLldyaXRlVVZhckludChzaWduRG9jLmFjY291bnRfbnVtYmVyKQp9CgppZiBzaWduRG9jLmFjY291bnRfc2VxdWVuY2UgIT0gMCB7CiAgICBidWYuV3JpdGVVVmFySW50NjQoMHgyOCkgLy8gd2lyZSB0eXBlIGFuZCBmaWVsZCBudW1iZXIgZm9yIGFjY291bnRfc2VxdWVuY2UKICAgIGJ1Zi5Xcml0ZVVWYXJJbnQoc2lnbkRvYy5hY2NvdW50X3NlcXVlbmNlKQp9Cg=="}})],1)]),e._v(" "),o("h3",{attrs:{id:"test-vectors"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#test-vectors"}},[e._v("#")]),e._v(" Test vectors")]),e._v(" "),o("p",[e._v("Given the protobuf definition "),o("code",[e._v("Article.proto")])]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"protobuf",base64:"cGFja2FnZSBibG9nOwpzeW50YXggPSAmcXVvdDtwcm90bzMmcXVvdDs7CgplbnVtIFR5cGUgewogIFVOU1BFQ0lGSUVEID0gMDsKICBJTUFHRVMgPSAxOwogIE5FV1MgPSAyOwp9OwoKZW51bSBSZXZpZXcgewogIFVOU1BFQ0lGSUVEID0gMDsKICBBQ0NFUFRFRCA9IDE7CiAgUkVKRUNURUQgPSAyOwp9OwoKbWVzc2FnZSBBcnRpY2xlIHsKICBzdHJpbmcgdGl0bGUgPSAxOwogIHN0cmluZyBkZXNjcmlwdGlvbiA9IDI7CiAgdWludDY0IGNyZWF0ZWQgPSAzOwogIHVpbnQ2NCB1cGRhdGVkID0gNDsKICBib29sIHB1YmxpYyA9IDU7CiAgYm9vbCBwcm9tb3RlZCA9IDY7CiAgVHlwZSB0eXBlID0gNzsKICBSZXZpZXcgcmV2aWV3ID0gODsKICByZXBlYXRlZCBzdHJpbmcgY29tbWVudHMgPSA5OwogIHJlcGVhdGVkIHN0cmluZyBiYWNrbGlua3MgPSAxMDsKfTsK"}}),e._v(" "),o("p",[e._v("serializing the values")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"yaml",base64:"dGl0bGU6ICZxdW90O1RoZSB3b3JsZCBuZWVkcyBjaGFuZ2Ug8J+MsyZxdW90OwpkZXNjcmlwdGlvbjogJnF1b3Q7JnF1b3Q7CmNyZWF0ZWQ6IDE1OTY4MDYxMTEwODAKdXBkYXRlZDogMApwdWJsaWM6IHRydWUKcHJvbW90ZWQ6IGZhbHNlCnR5cGU6IFR5cGUuTkVXUwpyZXZpZXc6IFJldmlldy5VTlNQRUNJRklFRApjb21tZW50czogWyZxdW90O05pY2Ugb25lJnF1b3Q7LCAmcXVvdDtUaGFuayB5b3UmcXVvdDtdCmJhY2tsaW5rczogW10K"}}),e._v(" "),o("p",[e._v("must result in the serialization")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"MGExYjU0Njg2NTIwNzc2ZjcyNmM2NDIwNmU2NTY1NjQ3MzIwNjM2ODYxNmU2NzY1MjBmMDlmOGNiMzE4ZThiZWJlYzhiYzJlMjgwMTM4MDI0YTA4NGU2OTYzNjUyMDZmNmU2NTRhMDk1NDY4NjE2ZTZiMjA3OTZmNzUK"}}),e._v(" "),o("p",[e._v("When inspecting the serialized document, you see that every second field is\nomitted:")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"JCBlY2hvIDBhMWI1NDY4NjUyMDc3NmY3MjZjNjQyMDZlNjU2NTY0NzMyMDYzNjg2MTZlNjc2NTIwZjA5ZjhjYjMxOGU4YmViZWM4YmMyZTI4MDEzODAyNGEwODRlNjk2MzY1MjA2ZjZlNjU0YTA5NTQ2ODYxNmU2YjIwNzk2Zjc1IHwgeHhkIC1yIC1wIHwgcHJvdG9jIC0tZGVjb2RlX3JhdwoxOiAmcXVvdDtUaGUgd29ybGQgbmVlZHMgY2hhbmdlIFwzNjBcMjM3XDIxNFwyNjMmcXVvdDsKMzogMTU5NjgwNjExMTA4MAo1OiAxCjc6IDIKOTogJnF1b3Q7TmljZSBvbmUmcXVvdDsKOTogJnF1b3Q7VGhhbmsgeW91JnF1b3Q7Cg=="}}),e._v(" "),o("h2",{attrs:{id:"consequences"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#consequences"}},[e._v("#")]),e._v(" Consequences")]),e._v(" "),o("p",[e._v("Having such an encoding available allows us to get deterministic serialization\nfor all protobuf documents we need in the context of Cosmos SDK signing.")]),e._v(" "),o("h3",{attrs:{id:"positive"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#positive"}},[e._v("#")]),e._v(" Positive")]),e._v(" "),o("ul",[o("li",[e._v("Well defined rules that can be verified independent of a reference\nimplementation")]),e._v(" "),o("li",[e._v("Simple enough to keep the barrier to implement transaction signing low")]),e._v(" "),o("li",[e._v("It allows us to continue to use 0 and other empty values in SignDoc, avoiding\nthe need to work around 0 sequences. This does not imply the change from\nhttps://github.com/cosmos/cosmos-sdk/pull/6949 should not be merged, but not\ntoo important anymore.")])]),e._v(" "),o("h3",{attrs:{id:"negative"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#negative"}},[e._v("#")]),e._v(" Negative")]),e._v(" "),o("ul",[o("li",[e._v("When implementing transaction signing, the encoding rules above must be\nunderstood and implemented.")]),e._v(" "),o("li",[e._v("The need for rule number 3. adds some complexity to implementations.")]),e._v(" "),o("li",[e._v("Some data structures may require custom code for serialization. Thus\nthe code is not very portable - it will require additional work for each\nclient implementing serialization to properly handle custom data structures.")])]),e._v(" "),o("h3",{attrs:{id:"neutral"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#neutral"}},[e._v("#")]),e._v(" Neutral")]),e._v(" "),o("h3",{attrs:{id:"usage-in-sdk"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#usage-in-sdk"}},[e._v("#")]),e._v(" Usage in SDK")]),e._v(" "),o("p",[e._v('For the reasons mentioned above ("Negative" section) we prefer to keep workarounds\nfor shared data structure. Example: the aforementioned '),o("code",[e._v("TxRaw")]),e._v(" is using raw bytes\nas a workaround. This allows them to use any valid Protobuf library without\nthe need of implementing a custom serializer that adheres to this standard (and related risks of bugs).")]),e._v(" "),o("h2",{attrs:{id:"references"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#references"}},[e._v("#")]),e._v(" References")]),e._v(" "),o("ul",[o("li",[o("sup",[e._v("1")]),e._v(" "),o("em",[e._v("When a message is serialized, there is no guaranteed order for\nhow its known or unknown fields should be written. Serialization order is an\nimplementation detail and the details of any particular implementation may\nchange in the future. Therefore, protocol buffer parsers must be able to parse\nfields in any order.")]),e._v(" from\nhttps://developers.google.com/protocol-buffers/docs/encoding#order")]),e._v(" "),o("li",[o("sup",[e._v("2")]),e._v(" https://developers.google.com/protocol-buffers/docs/encoding#signed_integers")]),e._v(" "),o("li",[o("sup",[e._v("3")]),e._v(" "),o("em",[e._v("Note that for scalar message fields, once a message is parsed\nthere's no way of telling whether a field was explicitly set to the default\nvalue (for example whether a boolean was set to false) or just not set at all:\nyou should bear this in mind when defining your message types. For example,\ndon't have a boolean that switches on some behavior when set to false if you\ndon't want that behavior to also happen by default.")]),e._v(" from\nhttps://developers.google.com/protocol-buffers/docs/proto3#default")]),e._v(" "),o("li",[o("sup",[e._v("4")]),e._v(" "),o("em",[e._v("When a message is parsed, if the encoded message does not\ncontain a particular singular element, the corresponding field in the parsed\nobject is set to the default value for that field.")]),e._v(" from\nhttps://developers.google.com/protocol-buffers/docs/proto3#default")]),e._v(" "),o("li",[o("sup",[e._v("5")]),e._v(" "),o("em",[e._v("Also note that if a scalar message field is set to its default,\nthe value will not be serialized on the wire.")]),e._v(" from\nhttps://developers.google.com/protocol-buffers/docs/proto3#default")]),e._v(" "),o("li",[o("sup",[e._v("6")]),e._v(" "),o("em",[e._v("For enums, the default value is the first defined enum value,\nwhich must be 0.")]),e._v(" from\nhttps://developers.google.com/protocol-buffers/docs/proto3#default")]),e._v(" "),o("li",[o("sup",[e._v("7")]),e._v(" "),o("em",[e._v("For message fields, the field is not set. Its exact value is\nlanguage-dependent.")]),e._v(" from\nhttps://developers.google.com/protocol-buffers/docs/proto3#default")]),e._v(" "),o("li",[e._v("Encoding rules and parts of the reasoning taken from\n"),o("a",{attrs:{href:"https://github.com/regen-network/canonical-proto3",target:"_blank",rel:"noopener noreferrer"}},[e._v("canonical-proto3 Aaron Craelius"),o("OutboundLink")],1)])])],1)}),[],!1,null,null,null);t.default=a.exports}}]);