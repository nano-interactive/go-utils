package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/nano-interactive/go-utils"
)

const jsonNull = "null"

var jsonNullBytes []byte // points to jsonNull

func init() {
	jsonNullBytes = utils.UnsafeBytes(jsonNull)
	bson.DefaultRegistry.RegisterTypeMapEntry(bson.TypeObjectID, reflect.TypeOf(ObjectID{}))
}

type (
	NullByte    sql.NullByte
	NullBool    sql.NullBool
	NullFloat32 struct {
		Float32 float32
		Valid   bool
	}
	NullFloat64 sql.NullFloat64
	NullTime    sql.NullTime
	NullString  sql.NullString
	NullInt16   sql.NullInt16
	NullInt32   sql.NullInt32
	NullInt64   sql.NullInt64
	NullUint16  struct {
		Uint16 uint16
		Valid  bool
	}
	NullUint32 struct {
		Uint32 uint32
		Valid  bool
	}
	NullUint64 struct {
		Uint64 uint64
		Valid  bool
	}
)

func (o *NullByte) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b byte
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Byte = b
	}

	return nil
}

func (o NullByte) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Byte)
}

func (o *NullBool) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
		o.Bool = false
	default:
		var b bool
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Bool = b
	}

	return nil
}
func (o *NullBool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Compare(jsonNullBytes, data) == 0 {
		o.Valid = false
		o.Bool = false
		return nil
	}

	switch utils.UnsafeString(data) {
	case "true":
		o.Valid = true
		o.Bool = true
	case "false":
		o.Valid = true
		o.Bool = false
	default:
		o.Valid = false
		o.Bool = false
	}

	return nil
}
func (o NullBool) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Bool)
}

func (o *NullFloat32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Float32)
}

func (o *NullFloat32) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b float32
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Float32 = b
	}

	return nil
}

func (o NullFloat64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Float64)
}

func (o *NullFloat64) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b float64
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Float64 = b
	}

	return nil
}

func (o NullString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.String)
}

func (o *NullString) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b string
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.String = b
	}

	return nil
}
func (o NullInt16) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Int16)
}
func (o *NullInt16) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b int16
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Int16 = b
	}

	return nil
}
func (o NullInt32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Int32)
}
func (o *NullInt32) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b int32
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Int32 = b
	}

	return nil
}
func (o NullInt64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Int64)
}
func (o *NullInt64) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b int64
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Int64 = b
	}

	return nil
}
func (o NullUint16) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Uint16)
}
func (o *NullUint16) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b uint16
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Uint16 = b
	}

	return nil
}
func (o NullUint32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Uint32)
}
func (o *NullUint32) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b uint32
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Uint32 = b
	}

	return nil
}
func (o NullUint64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !o.Valid {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(o.Uint64)
}
func (o *NullUint64) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	switch t {
	case bson.TypeNull:
		o.Valid = false
	default:
		var b uint64
		if err := bson.UnmarshalValue(t, bytes, &b); err != nil {
			return err
		}
		o.Valid = true
		o.Uint64 = b
	}

	return nil
}

func (o NullByte) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Byte)
}

func (o NullFloat32) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Float32)
}

func (o NullFloat64) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Float64)
}

func (o NullTime) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Time)
}

func (o NullString) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.String)
}

func (o NullInt16) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Int16)
}

func (o NullInt32) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Int32)
}

func (o NullInt64) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Int64)
}

func (o NullUint16) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Uint16)
}

func (o NullUint32) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Uint32)
}

func (o NullUint64) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return jsonNullBytes, nil
	}

	return json.Marshal(o.Uint64)
}

func (o *NullByte) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b byte
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Byte = b
	return nil
}

func (o *NullFloat32) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b float32
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Float32 = b
	return nil
}

func (o *NullFloat64) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b float64
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Float64 = b
	return nil
}

func (o *NullTime) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b time.Time
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Time = b
	return nil
}

func (o *NullString) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b string
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.String = b
	return nil
}

func (o *NullInt16) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b int16
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Int16 = b
	return nil
}

func (o *NullInt32) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b int32
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Int32 = b
	return nil
}

func (o *NullInt64) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b int64
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Int64 = b
	return nil
}

func (o *NullUint16) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b uint16
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Uint16 = b
	return nil
}

func (o *NullUint32) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b uint32
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Uint32 = b
	return nil
}

func (o *NullUint64) UnmarshalJSON(value []byte) error {
	if bytes.Compare(jsonNullBytes, value) == 0 {
		o.Valid = false
		return nil
	}

	var b uint64
	if err := json.Unmarshal(value, &b); err != nil {
		return err
	}

	o.Valid = true
	o.Uint64 = b
	return nil
}

func (n NullFloat32) Value() (driver.Value, error) {
	t := sql.NullFloat64{
		Float64: float64(n.Float32),
		Valid:   n.Valid,
	}

	return t.Value()
}

func (n *NullFloat32) Scan(value any) error {
	var t sql.NullFloat64
	if err := t.Scan(value); err != nil {
		return err
	}

	n.Valid = t.Valid
	n.Float32 = float32(t.Float64)
	return nil
}

func (n NullUint16) Value() (driver.Value, error) {
	t := sql.NullInt32{
		Int32: int32(n.Uint16),
		Valid: n.Valid,
	}

	return t.Value()
}

func (n *NullUint16) Scan(value any) error {
	var t sql.NullInt32
	if err := t.Scan(value); err != nil {
		return err
	}

	n.Valid = t.Valid
	n.Uint16 = uint16(t.Int32)
	return nil
}

func (n NullUint32) Value() (driver.Value, error) {
	t := sql.NullInt64{
		Int64: int64(n.Uint32),
		Valid: n.Valid,
	}

	return t.Value()
}

func (n *NullUint32) Scan(value any) error {
	var t sql.NullInt64
	if err := t.Scan(value); err != nil {
		return err
	}

	n.Valid = t.Valid
	n.Uint32 = uint32(t.Int64)
	return nil
}
