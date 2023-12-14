package types

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/tinylib/msgp/msgp"
)

var timeBase = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()

type CQLUUID gocql.UUID

func (u CQLUUID) MarshalJSON() ([]byte, error) {
	return utils.UnsafeBytes(u.toString(true)), nil
}

func (u *CQLUUID) UnmarshalJSON(data []byte) error {
	str := utils.UnsafeString(data)[1 : len(utils.UnsafeString(data))-1]
	if len(str) > 36 {
		return fmt.Errorf("invalid JSON CQLUUID %s", str)
	}

	parsed, err := gocql.ParseUUID(str)
	if err == nil {
		copy(u[:], parsed[:])
	}

	return err
}

func (u CQLUUID) MarshalText() ([]byte, error) {
	return utils.UnsafeBytes(u.String()), nil
}

func (u *CQLUUID) UnmarshalText(text []byte) error {
	uuid, err := gocql.ParseUUID(utils.UnsafeString(text))
	if err != nil {
		return err
	}

	*u = CQLUUID(uuid)

	return nil
}

func (u *CQLUUID) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	t := info.Type()
	if t == gocql.TypeUUID || t == gocql.TypeTimeUUID {
		uuid, err := gocql.UUIDFromBytes(data)
		if err != nil {
			return err
		}
		*u = CQLUUID(uuid)
		return nil
	}

	var uuid gocql.UUID

	err := gocql.Unmarshal(info, data, &uuid)
	if err != nil {
		return err
	}

	*u = CQLUUID(uuid)
	return nil
}

func (u CQLUUID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, gocql.UUID(u))
}

// EncodeMsg implements msgp.Encodable
func (u *CQLUUID) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.Append(0x81, 0xa4, 0x55, 0x55, 0x49, 0x44)
	if err != nil {
		return
	}
	err = en.WriteBytes(u[:])
	if err != nil {
		err = msgp.WrapError(err, "CQLUUID")
		return
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (u *CQLUUID) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "CQLUUID":
			var uuid [16]byte
			_, err = dc.ReadBytes(uuid[:])

			if err != nil {
				err = msgp.WrapError(err, "CQLUUID")
				return
			}
			*u = uuid
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (u *CQLUUID) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, u.Msgsize())
	o = append(o, 0x81, 0xa4, 0x55, 0x55, 0x49, 0x44)
	return msgp.AppendBytes(o, u[:]), nil
}

// UnmarshalMsg implements msgp.Unmarshaler
func (u *CQLUUID) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "CQLUUID":
			var uuidBytes [16]byte
			var value []byte
			value, bts, err = msgp.ReadBytesBytes(bts, uuidBytes[:])
			if err != nil {
				return nil, err
			}
			uuid, err := gocql.UUIDFromBytes(value)
			if err != nil {
				err = msgp.WrapError(err, "CQLUUID")
				return nil, err
			}
			*u = CQLUUID(uuid)
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts

	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (u *CQLUUID) Msgsize() int {
	return 1 + 5 + 16 // nolint:gomnd
}

var offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34} // nolint:all

const hexString = "0123456789abcdef"

func (u CQLUUID) toString(quote bool) string {
	var uuid [36]byte

	for i, b := range u {
		uuid[offsets[i]] = hexString[b>>4]
		uuid[offsets[i]+1] = hexString[b&0xF]
	}

	uuid[8] = '-'
	uuid[13] = '-'
	uuid[18] = '-'
	uuid[23] = '-'

	if !quote {
		return utils.UnsafeString(uuid[:])
	}

	var data [38]byte
	data[0] = '"'
	copy(data[1:], uuid[0:])
	data[37] = '"'
	return utils.UnsafeString(data[:])
}

// String returns the CQLUUID in it's canonical form, a 32 digit hexadecimal
// number in the form of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u CQLUUID) UnsafeString() string {
	return u.toString(false)
}

func (u CQLUUID) String() string {
	return string(u.toString(false))
}

// Bytes returns the raw byte slice for this UUID. A CQLUUID is always 128 bits
// (16 bytes) long.
func (u CQLUUID) Bytes() []byte {
	return u[:]
}

// Variant returns the variant of this CQLUUID. This package will only generate
// UUIDs in the IETF variant.
func (u CQLUUID) Variant() int {
	x := u[8]
	if x&0x80 == 0 {
		return gocql.VariantNCSCompat
	}
	if x&0x40 == 0 {
		return gocql.VariantIETF
	}
	if x&0x20 == 0 {
		return gocql.VariantMicrosoft
	}
	return gocql.VariantFuture
}

// Version extracts the version of this CQLUUID variant. The RFC 4122 describes
// five kinds of UUIDs.
func (u CQLUUID) Version() int {
	return int(u[6] & 0xF0 >> 4)
}

// Node extracts the MAC address of the node who generated this CQLUUID. It will
// return nil if the CQLUUID is not a time based UUID (version 1).
func (u CQLUUID) Node() []byte {
	if u.Version() != 1 {
		return nil
	}
	return u[10:]
}

// Clock extracts the clock sequence of this CQLUUID. It will return zero if the
// CQLUUID is not a time based UUID (version 1).
func (u CQLUUID) Clock() uint32 {
	if u.Version() != 1 {
		return 0
	}

	// Clock sequence is the lower 14bits of u[8:10]
	return uint32(u[8]&0x3F)<<8 | uint32(u[9])
}

// Timestamp extracts the timestamp information from a time based CQLUUID
// (version 1).
func (u CQLUUID) Timestamp() int64 {
	if u.Version() != 1 {
		return 0
	}
	return int64(uint64(u[0])<<24|uint64(u[1])<<16|
		uint64(u[2])<<8|uint64(u[3])) +
		int64(uint64(u[4])<<40|uint64(u[5])<<32) +
		int64(uint64(u[6]&0x0F)<<56|uint64(u[7])<<48)
}

// Time is like Timestamp, except that it returns a time.Time.
func (u CQLUUID) Time() time.Time {
	if u.Version() != 1 {
		return time.Time{}
	}
	t := u.Timestamp()
	sec := t / 1e7
	nsec := (t % 1e7) * 100
	return time.Unix(sec+timeBase, nsec).UTC()
}
