package schema

// Leaf ...
//
// NB: To build or rebuild, use https://jsonschema.net/
var Leaf string = `{
  "$schema": "http://json-schema.org/draft-07/schema",
  "$id": "http://example.com/example.json",
  "type": "object",
  "title": "The root schema",
  "description": "The root schema comprises the entire JSON document.",
  "default": {},
  "examples": [
    {
      "id": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
      "position": 0,
      "size": 100,
      "origin": "",
      "previous": "",
      "next": ""
    }
  ],
  "required": [
    "id",
    "position",
    "size"
  ],
  "additionalProperties": true,
  "properties": {
    "id": {
      "$id": "#/properties/id",
      "type": "string",
      "title": "The id schema",
      "description": "An explanation about the purpose of this instance.",
      "default": "",
      "examples": [
        "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
      ]
    },
    "position": {
      "$id": "#/properties/position",
      "type": "integer",
      "title": "The position schema",
      "description": "An explanation about the purpose of this instance.",
      "default": 0,
      "examples": [
        0
      ]
    },
    "size": {
      "$id": "#/properties/size",
      "type": "integer",
      "title": "The size schema",
      "description": "An explanation about the purpose of this instance.",
      "default": 0,
      "examples": [
        100
      ]
    },
    "origin": {
      "$id": "#/properties/origin",
      "type": "string",
      "title": "The origin schema",
      "description": "An explanation about the purpose of this instance.",
      "default": "",
      "examples": [
        ""
      ]
    },
    "previous": {
      "$id": "#/properties/previous",
      "type": "string",
      "title": "The previous schema",
      "description": "An explanation about the purpose of this instance.",
      "default": "",
      "examples": [
        ""
      ]
    },
    "next": {
      "$id": "#/properties/next",
      "type": "string",
      "title": "The next schema",
      "description": "An explanation about the purpose of this instance.",
      "default": "",
      "examples": [
        ""
      ]
    }
  }
}`
