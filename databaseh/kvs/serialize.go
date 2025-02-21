// Package kvs implements a key/value store. These functions provide for
// object serialization/deserialization to the kvs.
package kvs

import (
	"encoding/json"

	"github.com/paulfdunn/go-helper/osh/v2/runtimeh"
)

// Deserialize deserializes an object from the KVS.
// Note that the caller needs to call with a pointer to the object to deserialize.
// If fields in the persisted object are non-nil, they will overwrite fields in the
// provided object, otherwise values in the provide object will be in the returned object.
// If the key is not in the KVS, obj is unchanged and there is no error.
func (kvs KVS) Deserialize(key string, obj interface{}) error {
	var b []byte
	var err error
	if b, err = kvs.Get(key); err != nil {
		return runtimeh.SourceInfoError("", err)
	}

	if b == nil {
		// No data and no error means there was no object in the KVS; obj is unchanged
		return nil
	}

	// Merge persisted data into obj.
	if err := json.Unmarshal(b, &obj); err != nil {
		return runtimeh.SourceInfoError("", err)
	}

	return nil
}

// Serialize serializes an object into the KVS.
func (kvs KVS) Serialize(key string, obj interface{}) error {
	var b []byte
	var err error
	if b, err = json.Marshal(obj); err != nil {
		return runtimeh.SourceInfoError("", err)
	}

	return runtimeh.SourceInfoError("", kvs.Set(key, b))
}
